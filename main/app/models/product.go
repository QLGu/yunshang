package models

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	. "github.com/ahmetalpbalkan/go-linq"
	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
	"strconv"
)

func NewProductService(db *xorm.Session) *ProductService {
	return &ProductService{db}
}

/////////////////////////////////////////////////////////////////////////////

type ProductService struct {
	db *xorm.Session
}

func (self ProductService) Total() int64 {
	total, _ := self.db.Count(&entity.Product{})
	return total
}

func (self ProductService) FindAllProductsForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Product{})
	gotang.AssertNoError(err, "")

	var products []entity.Product

	err1 := ps.BuildQuerySession().Find(&products)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, products, ps)
}

func (self ProductService) FindAllAvailableProductsForPage(ps *PageSearcher) (page *PageData) {
	ps.FilterCall = func(db *xorm.Session) {
		db.And("enabled=?", true)
	}

	ps.SearchKeyCall = func(db *xorm.Session) {
		db.And("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Product{})
	gotang.AssertNoError(err, "")

	var products []entity.Product
	err1 := ps.BuildQuerySession().Find(&products)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, products, ps)
}

func (self ProductService) FindAllAvailableProducts() (ps []entity.Product) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
	return
}

func (self ProductService) FindProductsByIds(ids []int64) (ps []entity.Product) {
	_ = self.db.Where("id in (" + strings.Join(asStrSliceFromInt64(ids), ",") + ")").Find(&ps)
	return
}

func (self ProductService) FindAllAvailableProductsByCtCode(ctcode string) (ps []entity.Product) {
	s := self.db.Where("enabled=?", true)
	if len(ctcode) > 0 {
		s.And("category_code like ?", ctcode+"%")
	}
	_ = s.Find(&ps)
	return
}

func (self ProductService) FindAvailableCategoryChainByCode(ctcode string) (ps []entity.ProductCategory) {
	if len(ctcode) == 0 {
		return
	}
	s := self.db.Where("enabled=?", true)

	codes := strings.Split(ctcode, "-")
	println("LENNN", len(codes))
	filterCodes := make([]string, len(codes))
	for i := 0; i < len(codes); i++ {
		e := i + 1
		filterCodes[i] = "'" + strings.Join(codes[0:e], "-") + "'"
	}

	s.And(fmt.Sprintf("code in (%s)", strings.Join(filterCodes, ","))).Asc("id")
	err := s.Find(&ps)
	gotang.AssertNoError(err, "FindAvailableCategoryChainByCode")

	return
}

func (self ProductService) GetProductById(id int64) (p entity.Product, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) GetProductPrefPrice(productId int64, quantity int) float64 {
	prices := self.FindProductPrices(productId)
	e, f, err := From(prices).FirstBy(func(t T) (bool, error) {
		p := t.(entity.ProductPrices)
		if p.EndQuantity != 0 {
			if quantity >= p.StartQuantity && quantity <= p.EndQuantity {
				return true, nil
			}
		} else {
			if quantity >= p.StartQuantity {
				return true, nil
			}
		}
		return false, nil
	})
	gotang.AssertNoError(err, "GetProductPrefPrice")

	if !f {
		return self.GetProductUnitPrice(productId)
	}
	return e.(entity.ProductPrices).Price
}

func (self ProductService) GetCategoryCode(id int64) string {
	c, _ := self.GetCategoryById(id)
	return c.Code
}

func (self ProductService) SaveProduct(p entity.Product) (id int64, err error) {
	if p.Id == 0 { //insert
		p.CategoryCode = self.GetCategoryCode(p.CategoryId)
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}

		p.Code = p.Id + entity.ProductStartDisplayCode //编码
		_, err = self.db.Id(p.Id).Cols("code").Update(&p)
		id = p.Id

		FireEvent(EventObject{Name: EStockLog, SourceId: p.Id, Message: fmt.Sprintf("创建产品入库：%d(%s)", p.StockNumber, p.UnitName)})
		return
	} else { // update
		currDa, ok := self.GetProductById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			if p.CategoryId != currDa.CategoryId || len(currDa.CategoryCode) == 0 { //检测类型是否修改
				p.CategoryCode = self.GetCategoryCode(p.CategoryId) //更新分类编码， 冗余设计
			}

			_, err = self.db.Id(p.Id).Update(&p)
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此产品不存在")
		}
	}
}

func (self ProductService) AddProductStock(id int64, stock int, message string) (newstock int, err error) {
	p, ok := self.GetProductById(id)
	if !ok {
		err = errors.New("产品不存在!")
		return
	}
	p.StockNumber += stock
	_, err = self.db.Id(p.Id).Cols("stock_number").Update(&p)

	FireEvent(EventObject{Name: EStockLog, SourceId: p.Id, Message: fmt.Sprintf("入库：%d(%s), 当前库存%d(%s), %s", stock, p.UnitName, p.StockNumber, p.UnitName, message)})

	newstock = p.StockNumber
	return
}

func (self ProductService) ToggleProductEnabled(p *entity.Product) error {
	p.Enabled = !p.Enabled
	if p.Enabled {
		p.EnabledAt = time.Now()
		_, err := self.db.Id(p.Id).Cols("enabled", "enabled_at").Update(p)
		return err
	}
	p.UnEnabledAt = time.Now()
	_, err := self.db.Id(p.Id).Cols("enabled", "un_enabled_at").Update(p)
	return err
}

func (self ProductService) FindProductPrices(id int64) (prices []entity.ProductPrices) {
	_ = self.db.Where("product_id=?", id).Asc("start_quantity").Find(&prices)
	return
}

func (self ProductService) HasProductSetPrice(id int64) bool {
	ps := self.FindProductPrices(id)
	if len(ps) == 0 {
		return false
	}
	for _, v := range ps {
		if v.Price == 0 {
			return false
		}
	}
	return true
}

func (self ProductService) GetProductUnitPrice(id int64) float64 {
	prices := self.FindProductPrices(id)
	if len(prices) == 0 {
		return 0
	}

	min, err := From(prices).Select(func(t T) (T, error) { return t.(entity.ProductPrices).StartQuantity, nil }).MinInt()
	gotang.AssertNoError(err, "")
	mf := func(t T) (bool, error) { return t.(entity.ProductPrices).StartQuantity == min, nil }
	p, found, err := From(prices).FirstBy(mf)
	gotang.Assert(found, "")
	gotang.AssertNoError(err, "")

	return p.(entity.ProductPrices).Price
}

func (self ProductService) GetProductPricesSplits(productId int64) string {
	return strings.Join(asStrSlice(self.GetProductPricesSplitInts(productId)), " ")
}

func (self ProductService) GetProductPricesSplitInts(productId int64) []int {
	prices := self.FindProductPrices(productId)
	if len(prices) == 0 {
		return nil
	}
	ss, err := From(prices).Select(func(t T) (T, error) { return t.(entity.ProductPrices).StartQuantity, nil }).OrderInts().Results()
	gotang.AssertNoError(err, "")

	return asIntSlice(ss)
}

func (self ProductService) SplitProductPrices(productId int64, start_quantitys string) error {
	var starts []int
	if len(start_quantitys) == 0 {
		return fmt.Errorf("请输入数字")
	}
	startArr := strings.Fields(start_quantitys)
	starts = make([]int, len(startArr))
	for index, v := range startArr {
		start, err := strconv.Atoi(v)
		if err != nil {
			return fmt.Errorf("请输入数字, %v", v)
		}
		starts[index] = start
	}

	_starts, err := From(starts).OrderInts().Results()
	gotang.AssertNoError(err, "")
	starts = asIntSlice(_starts)
	if starts[0] != 1 {
		starts = append([]int{1}, starts...)
	}

	//检查变化
	if eqSlice(starts, self.GetProductPricesSplitInts(productId)) {
		return fmt.Errorf("输入跟已有定价条目无变化！")
	} else {
		//清除已有的
		prices := self.FindProductPrices(productId)
		for _, price := range prices {
			self.db.Delete(price)
		}
	}

	for i := 0; i < len(starts); i++ {
		start_quantity := starts[i]
		var end_quantity int
		if i+1 == len(starts) {
			end_quantity = 0
		} else {
			end_quantity = starts[i+1] - 1
		}

		var p = entity.ProductPrices{
			ProductId:     productId,
			Name:          fmt.Sprintf("%v", start_quantity),
			StartQuantity: start_quantity,
			EndQuantity:   end_quantity,
		}

		_, err := self.db.Insert(&p)
		if err != nil {
			return err
		}
	}
	return nil
}

func (self ProductService) UpdateProductPrice(id int64) (err error) {
	p, exists := self.GetProductById(id)
	gotang.Assert(exists, "")

	p.Price = self.GetProductUnitPrice(id)
	_, err = self.db.Id(p.Id).Cols("price").Update(&p)
	return err
}

func (self ProductService) FindAllProvidersForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.Provider{})
	gotang.AssertNoError(err, "")

	var providers []entity.Provider
	err1 := ps.BuildQuerySession().Find(&providers)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, providers, ps)
}

// 推荐的品牌
func (self ProductService) RecommendProviders() (ps []entity.Provider) {
	_ = self.db.Where("enabled=? and tags=?", true, "推荐").Desc("id").Find(&ps)
	return
}

func (self ProductService) FindAllAvailableProviders() (ps []entity.Provider) {
	_ = self.db.Where("enabled=?", true).Find(&ps)
	return
}

func (self ProductService) GetProviderById(id int64) (p entity.Provider, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) GetProviderByProductId(id int64) (p entity.Provider, ok bool) {
	product, _ := self.GetProductById(id)
	ok, _ = self.db.Where("id=?", product.ProviderId).Get(&p)
	return
}

func (self ProductService) SaveProvider(p entity.Provider) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}
		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetProviderById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.db.Id(p.Id).Update(&p)

			if len(p.Tags) == 0 { //标签更新特殊处理
				currDa, ok = self.GetProviderById(p.Id)
				currDa.Tags = ""
				self.db.Id(currDa.Id).Cols("tags").Update(&currDa)
			}
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此制造商不存在")
		}
	}
}

func (self ProductService) ToggleProviderEnabled(p *entity.Provider) error {
	p.Enabled = !p.Enabled
	_, err := self.db.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self ProductService) DeleteProvider(id int64) error {
	p, _ := self.GetProviderById(id)
	_, err := self.db.Delete(&p)
	return err
}

const detailFilePattern = "data/products/detail/%d.html"

func (self ProductService) SaveProductDetail(id int64, content string) (err error) {
	p, ok := self.GetProductById(id)
	if !ok {
		return errors.New("产品不存在！")
	}

	to := fmt.Sprintf(detailFilePattern, p.Id)
	_, err = os.Create(to)
	if err != nil {
		return
	}

	err = ioutil.WriteFile(to, []byte(content), 0644)
	if err != nil {
		return
	}
	return
}

func (self ProductService) GetProductDetail(id int64) (detail string, err error) {
	from := fmt.Sprintf(detailFilePattern, id)
	f, err := os.Open(from)
	if err != nil {
		return
	}

	r, err := ioutil.ReadAll(f)
	if err != nil {
		return
	}
	detail = string(r)
	return
}

func (self ProductService) GetProductImageFile(file string, t int) (targetFile *os.File, err error) {
	var dir string
	switch t {
	case entity.PTScheDiag:
		dir = revel.Config.StringDefault("dir.data.product.sd", "data/products/sd")
	case entity.PTPics:
		dir = revel.Config.StringDefault("dir.data.product.pics", "data/products/pics")
	default:
		return targetFile, errors.New("不支持类型")
	}
	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err = os.Open(imageFile)
	return
}

func (self ProductService) DeleteProductParams(id int64) (err error) {
	var p entity.ProductParams
	_, err = self.db.Where("id=?", id).Get(&p)
	if err != nil {
		return
	}
	_, err = self.db.Delete(&p)
	if err != nil {
		return
	}

	return nil
}

func (self ProductService) FindProductImages(id int64, t int) (ps []entity.ProductParams) {
	_ = self.db.Where("type=? and product_id=?", t, id).Find(&ps)
	return
}

func (self ProductService) FindPrefProducts(limit int) (ps []entity.Product) {
	session := self.availableQuery().And("tags like ?", "%#最新优惠%")
	if limit > 0 {
		session.Limit(limit)
	}
	_ = session.Find(&ps)
	return
}

func (self ProductService) FindSpecialOfferProducts(limit int) (ps []entity.Product) {
	session := self.availableQuery().And("tags like ?", "%#特价%")
	if limit > 0 {
		session.Limit(limit)
	}
	_ = session.Find(&ps)
	return
}

func (self ProductService) FindHotProducts() (ps []entity.Product) {
	_ = self.availableQuery().And("tags like ?", "%#热门产品%").Find(&ps)
	return
}

////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////
// Categories

func (self ProductService) FindAllCategoriesForPage(ps *PageSearcher) (page *PageData) {
	ps.SearchKeyCall = func(db *xorm.Session) {
		db.Where("name like ?", "%"+ps.Search+"%")
	}

	total, err := ps.BuildCountSession().Count(&entity.ProductCategory{})
	gotang.AssertNoError(err, "")

	var categories []entity.ProductCategory
	err1 := ps.BuildQuerySession().Find(&categories)
	gotang.AssertNoError(err1, "")

	return NewPageData(total, categories, ps)
}

func (self ProductService) availableQuery() *xorm.Session {
	return self.db.Where("enabled=?", true)
}

func (self ProductService) FindAllAvailableCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().Find(&ps)
	return
}

func (self ProductService) FindAvailableTopCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", 0).Find(&ps)
	return
}

func (self ProductService) FindAvailableLeafCategories() (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id !=?", 0). /*.Limit(10)*/ Find(&ps)
	return
}

func (self ProductService) FindAllAvailableCategoriesByParentId(id int64) (ps []entity.ProductCategory) {
	_ = self.availableQuery().And("parent_id=?", id).Find(&ps)
	return
}

func (self ProductService) GetCategoryById(id int64) (p entity.ProductCategory, ok bool) {
	ok, _ = self.db.Where("id=?", id).Get(&p)
	return
}

func (self ProductService) SaveCategory(p entity.ProductCategory) (id int64, err error) {
	if p.Id == 0 { //insert
		_, err = self.db.Insert(&p)
		if err != nil {
			return
		}
		_ = self.updateCategoryCode(p)

		id = p.Id
		return
	} else { // update
		currDa, ok := self.GetCategoryById(p.Id)
		if ok {
			p.DataVersion = currDa.DataVersion
			_, err = self.db.Id(p.Id).Update(&p)
			if p.ParentId != currDa.ParentId {
				_ = self.updateCategoryCode(p)
			}
			return p.Id, err
		} else {
			return 0, fmt.Errorf("此分类不存在")
		}
	}
}

func (self ProductService) updateCategoryCode(p entity.ProductCategory) (err error) {
	p, _ = self.GetCategoryById(p.Id) //Hacked
	if p.ParentId != 0 {
		parent, _ := self.GetCategoryById(p.ParentId)
		p.Code = fmt.Sprintf("%v-%v", parent.Code, p.Id) //编码
	} else {
		p.Code = fmt.Sprintf("%v", p.Id) //编码
	}
	_, err = self.db.Id(p.Id).Cols("code").Update(&p)
	return
}

func (self ProductService) ToggleCategoryEnabled(p *entity.ProductCategory) error {
	p.Enabled = !p.Enabled
	_, err := self.db.Id(p.Id).Cols("enabled").Update(p)
	return err
}

func (self ProductService) FindAllProductStockLogs(id int64) (ps []entity.ProductStockLog) {
	_ = self.db.Where("product_id=?", id).Desc("id").Find(&ps)
	return
}

func (self ProductService) SubProductStockNumbersByOrder(order entity.Order) (err error) {
	ps := NewOrderService(self.db).GetOrderItemsByAdmin(order.Id)

	for _, p := range ps {
		_, err = self.AddProductStock(p.ProductId, -p.Quantity, fmt.Sprintf("订单%d 支付成功，减库存！", order.Code))
		if err != nil {
			return
		}
	}
	return
}

func (self ProductService) FindLatestProducts(limit int) (ps []entity.Product) {
	_ = self.availableQuery().Limit(limit).Desc("id").Find(&ps)
	return
}
