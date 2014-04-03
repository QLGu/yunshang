package controllers

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/itang/gotang"
	gio "github.com/itang/gotang/io"
	"github.com/itang/yunshang/main/app/models/entity"
	"github.com/itang/yunshang/main/app/routes"
	"github.com/itang/yunshang/main/app/utils"
	"github.com/lunny/xorm"
	"github.com/revel/revel"
)

// 相关Actions
type News struct {
	AppController
}

// 产品主页
func (c News) Index() revel.Result {
	gs_tws := c.newsApi().FindTWNews("1", 4)
	gs_no_tws := c.newsApi().FindNoTWNews("1", 18)

	hys := c.newsApi().FindNews("2", 22)
	jss := c.newsApi().FindNews("3", 22)

	c.setChannel("news/index")

	return c.Render(gs_tws, gs_no_tws, hys, jss)
}

func (c News) View(id int64) revel.Result {
	if id == 0 {
		return c.NotFound("此文章不存在！")
	}

	news, exists := c.newsApi().GetNewsById(id)
	if !exists {
		return c.NotFound("此文章不存在！")
	}
	if !news.Enabled && !isAdmin(c.Session) {
		return c.NotFound("此文章不存在！")
	}

	newsDetail, err := c.newsApi().GetNewsDetail(id)
	if err != nil {
		newsDetail = ""
	}
	files := c.newsApi().FindNewsMaterial(id)

	prevs := c.newsApi().GetPrevDisplayNews(id)
	nexts := c.newsApi().GetNextDisplayNews(id)

	if news.IsServiceArticle() {
		c.setChannel("services/view")
	} else if news.IsAboutArticle() {
		c.setChannel("about/view")
	} else {
		c.setChannel("news/view")
	}
	return c.Render(news, newsDetail, files, prevs, nexts)
}

func (c News) DetailSummary(id int64) revel.Result {
	newsDetail, err := c.newsApi().GetNewsDetail(id)
	if err != nil {
		newsDetail = ""
	}
	return c.RenderHtml(truncStr(newsDetail, 100, "..."))
}

func (c News) List(code string) revel.Result {
	news := c.newsApi().FindNews(code, 22)

	ct, _ := c.newsApi().GetCategoryByCodeString(code)

	c.setChannel("news/list")
	return c.Render(news, ct)
}

func (c News) SdImages(id int64) revel.Result {
	images := c.newsApi().FindNewsImages(id, entity.NTScheDiag)
	return c.RenderJson(Success("", images))
}

func (c News) SdImage(file string) revel.Result {
	targetFile, err := c.newsApi().GetNewsImageFile(file, entity.NTScheDiag)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

func (c News) ImagePics(id int64) revel.Result {
	images := c.newsApi().FindNewsImages(id, entity.NTPics)
	return c.RenderJson(Success("", images))
}

func (c News) ImagePicsList(id int64) revel.Result {
	var images []entity.NewsParam
	c.db.Where("type=? and news_id=?", entity.NTPics, id).Find(&images)
	var ret = ""
	for _, v := range images {
		ret += fmt.Sprintf("?file=%sue_separate_ue", v.Value)
	}
	return c.RenderText(ret)
}

func (c News) ImagePic(file string) revel.Result {
	targetFile, err := c.newsApi().GetNewsImageFile(file, entity.NTPics)
	if err != nil {
		return c.NotFound("No found file " + file)
	}
	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

func (c News) MFiles(id int64) revel.Result {
	files := c.newsApi().FindNewsMaterial(id)
	return c.RenderJson(Success("", files))
}

// 材料
func (c News) MFile(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.news.m", "data/news/m")
	f := filepath.Join(dir, filepath.Base(file))
	targetFile, err := os.Open(f)
	if err != nil {
		return c.NotFound("No File Found！")
	}

	return c.RenderFile(targetFile, "")
}

// param file： 图片标识： {{id}}.jpg
func (c News) Image(file string) revel.Result {
	dir := revel.Config.StringDefault("dir.data.news.logo", "data/news/logo")

	imageFile := filepath.Join(dir, filepath.Base(file))
	if !(gio.Exists(imageFile) && gio.IsFile(imageFile)) {
		imageFile = filepath.Join("public/img", "default.png")
	}

	targetFile, err := os.Open(imageFile)
	if err != nil {
		return c.NotFound("No Image Found！")
	}

	c.Response.ContentType = "image/jpg"
	return c.RenderFile(targetFile, "")
}

///////////////////////

func (c Admin) deleteNewsParam(id int64) revel.Result {
	if ret := c.checkErrorAsJsonResult(c.newsApi().DeleteNewsParam(id)); ret != nil {
		return ret
	}

	return c.RenderJson(Success("删除成功！", ""))
}

func (c Admin) DeleteNewsSdImage(id int64) revel.Result {
	return c.deleteNewsParam(id)
}

func (c Admin) DeleteNewsImagePic(id int64) revel.Result {
	return c.deleteNewsParam(id)
}

func (c Admin) DeleteNewsMFile(id int64) revel.Result {
	return c.deleteNewsParam(id)
	//TODO delete file?
}

func (c Admin) NewNews(id int64) revel.Result {
	var (
		p      entity.News
		detail = ""
	)

	if id == 0 { // new

	} else { //edit
		p, _ = c.newsApi().GetNewsById(id)
		detail, _ = c.newsApi().GetNewsDetail(p.Id)
	}
	return c.Render(p, detail)
}

func (c Admin) DoNewNews(p entity.News) revel.Result {
	c.Validation.Required(p.Title).Message("请填写名称").Key("name")

	if ret := c.doValidate(routes.Admin.NewNews(p.Id)); ret != nil {
		return ret
	}

	id, err := c.newsApi().SaveNews(p)
	if err != nil {
		c.Flash.Error("保存新闻失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存新闻成功！")
	}

	return c.Redirect(routes.Admin.NewNews(id))
}

func (c Admin) ToggleNewsEnabled(id int64) revel.Result {
	api := c.newsApi()
	p, ok := api.GetNewsById(id)
	if !ok {
		return c.RenderJson(Error("新闻不存在", nil))
	}

	err := api.ToggleNewsEnabled(&p)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(Success("发布成功！", nil))
		}
		return c.RenderJson(Success("取消发布成功！", nil))
	}
}

///////////////////////////

func (c Admin) News() revel.Result {
	c.setChannel("news/index")
	return c.Render()
}

func (c Admin) NewsData(filter_status string, filter_tag string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
		if len(filter_tag) > 0 {
			session.And("tags like ?", "%"+filter_tag+"%")
		}
	})
	page := c.newsApi().FindAllNewsForPage(ps)
	return c.renderDTJson(page)
}

func (c Admin) NewsCategories() revel.Result {
	c.setChannel("news/category")
	return c.Render()
}

func (c Admin) NewsCategoriesData() revel.Result {
	page := c.newsApi().FindAllCategoriesForPage(c.pageSearcher())
	return c.renderDTJson(page)
}

func (c Admin) NewNewsCategory(id int64) revel.Result {
	var p entity.NewsCategory
	if id == 0 { // new
		//p = entity.Provider{}
	} else { //edit
		p, _ = c.newsApi().GetCategoryById(id)
	}
	return c.Render(p)
}

func (c Admin) DoNewNewsCategory(p entity.NewsCategory) revel.Result {
	c.Validation.Required(p.Name).Message("请填写名称")

	if ret := c.doValidate(routes.Admin.NewNewsCategory(p.Id)); ret != nil {
		return ret
	}

	id, err := c.newsApi().SaveCategory(p)
	if err != nil {
		c.Flash.Error("保存分类失败，请重试！" + err.Error())
	} else {
		c.Flash.Success("保存分类成功！")
	}

	return c.Redirect(routes.Admin.NewNewsCategory(id))
}

func (c Admin) ToggleNewsCategoryEnabled(id int64) revel.Result {
	api := c.newsApi()
	p, ok := api.GetCategoryById(id)
	if !ok {
		return c.RenderJson(Error("分类不存在", nil))
	}

	err := api.ToggleCategoryEnabled(&p)
	if err != nil {
		return c.RenderJson(Error(err.Error(), nil))
	} else {
		if p.Enabled {
			return c.RenderJson(Success("激活成功！", nil))
		}
		return c.RenderJson(Success("禁用成功！", nil))
	}
}

func (c Admin) NewsComments() revel.Result {
	c.setChannel("news/comments")
	return c.Render()
}

func (c Admin) NewsCommentsData(filter_status string) revel.Result {
	ps := c.pageSearcherWithCalls(func(session *xorm.Session) {
		switch filter_status {
		case "true":
			session.And("enabled=?", true)
		case "false":
			session.And("enabled=?", false)
		}
		session.And("target_type=?", entity.CT_ARTICLE)
	})
	page := c.newsApi().NewsCommentsForPage(ps)
	return c.renderDTJson(page)
}

// 上传Logo
func (c Admin) UploadNewsLogo(id int64, image *os.File) revel.Result {
	c.Validation.Required(image != nil)
	if c.Validation.HasErrors() {
		return c.RenderJson(Error("请选择图片", nil))
	}
	p, exists := c.productApi().GetProductById(id)
	if !exists {
		return c.RenderJson(Error("操作失败，新闻不存在", nil))
	}
	to := filepath.Join(revel.Config.StringDefault("dir.data.news。logo", "data/news/logo"), fmt.Sprintf("%d.jpg", p.Id))

	err := utils.MakeAndSaveFromReader(image, to, "thumbnail", 200, 200)
	if ret := c.checkUploadError(err, "保存上传图片报错！"); ret != nil {
		return ret
	}

	return c.RenderJson(Success("上传成功", nil))
}

func (c Admin) UploadNewsImage(id int64, t int) revel.Result {
	var (
		dir, ct string
		count   int
	)

	if t == entity.NTScheDiag {
		dir = "data/news/sd/"
		ct = "thumbnail"
	} else if t == entity.NTPics {
		dir = "data/news/pics/"
		ct = "fit"
	} else {
		return c.RenderJson(Error("上传失败！ 类型不对", nil))
	}

	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			p := entity.NewsParam{Type: t, Name: fileHeader.Filename, NewsId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to := utils.Uuid() + ".jpg"
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, ct, 200, 200)
			gotang.AssertNoError(err, "生成图片出错！")

			count += 1
		}
	}

	if count == 0 {
		return c.RenderJson(Error("请选择要上传的图片", nil))
	}

	return c.RenderJson(Success("上传成功！", nil))
}

func (c Admin) UploadNewsImageForUEditor(id int64) revel.Result {
	dir := "data/news/pics/"
	ct := "fit"
	t := entity.NTPics

	var Original = ""
	var Url = ""
	var Title = ""
	var State = ""
	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			p := entity.NewsParam{Type: t, Name: fileHeader.Filename, NewsId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to := utils.Uuid() + ".jpg"
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			from, _ := fileHeader.Open()
			err = utils.MakeAndSaveFromReader(from, dir+to, ct, 200, 200)
			gotang.AssertNoError(err, "生成图片出错！")

			Original = fileHeader.Filename
			Title = Original
			State = "SUCCESS"
			Url = "?file=" + to
		}
	}

	ret := struct {
		Original string `json:"original"`
		Url      string `json:"url"`
		Title    string `json:"title"`
		State    string `json:"state"`
	}{Original, Url, Title, State}
	return c.RenderJson(ret)
}

func (c Admin) UploadNewsMaterial(id int64) revel.Result {
	count := 0
	for _, fileHeaders := range c.Params.Files {
		for _, fileHeader := range fileHeaders {
			to := ""
			p := entity.NewsParam{Type: entity.NTMaterial, Name: fileHeader.Filename, Value: to, NewsId: id}
			e, err := c.db.Insert(&p)
			gotang.Assert(e == 1, "New")
			gotang.AssertNoError(err, `Insert`)

			to = fmt.Sprintf("%d-%s", p.Id, fileHeader.Filename)
			p.Value = to
			c.db.Id(p.Id).Cols("value").Update(&p)

			out, err := os.Create("data/news/m/" + to)
			gotang.AssertNoError(err, `os.Create`)

			in, err := fileHeader.Open()
			gotang.AssertNoError(err, `fileHeader.Open()`)

			io.Copy(out, in)

			out.Close()
			in.Close()
			count += 1
		}
	}
	if count == 0 {
		return c.RenderJson(Error("请选择要上传的文件", nil))
	}

	return c.RenderJson(Success("上传成功！", nil))
}

func (c Admin) SaveNewsDetail(id int64, content string) revel.Result {
	err := c.newsApi().SaveNewsDetail(id, content)
	if err != nil {
		return c.RenderJson(Error("保存信息出错,"+err.Error(), nil))
	}

	return c.RenderJson(Success("保存信息成功！", nil))
}
