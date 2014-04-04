package utils

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"math/rand"
	"time"
	"reflect"

	"github.com/disintegration/imaging"
	"github.com/itang/gotang"
	"github.com/nu7hatch/gouuid"
	"github.com/revel/revel"
	"github.com/revel/revel/cache"
)

// UUID
func Uuid() string {
	u4, err := uuid.NewV4()
	gotang.AssertNoError(err, "")

	return u4.String()
}

// SHA1
func Sha1(content string) string {
	h := sha1.New()
	io.WriteString(h, content)

	return fmt.Sprintf("%x", h.Sum(nil))
}

// 随机字符串
func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randInt(65, 90))
	}

	return string(bytes)
}

func randInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max - min)
}

// 显示模板
func RenderTemplate(templatePath string, data interface{}) string {
	template, err := revel.MainTemplateLoader.Template(templatePath)
	gotang.AssertNoError(err, "")

	var b bytes.Buffer
	err = template.Render(&b, data)
	gotang.AssertNoError(err, "")

	return b.String()
}

// 生成并保存缩略图
func MakeAndSaveThumbnail(fromFile string, toFile string, w, h int) error {
	tnImage, err := MakeThumbnail(fromFile, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

// 生成并保存缩略图
func MakeAndSaveThumbnailFromReader(reader io.Reader, toFile string, w, h int) error {
	tnImage, err := MakeThumbnailFromReader(reader, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

func MakeAndSaveFromReader(reader io.Reader, toFile string, t string, w, h int) error {
	tnImage, err := MakeFromReader(reader, t, w, h)
	if err != nil {
		return err
	}
	return imaging.Save(tnImage, toFile)
}

// 生成缩略图
func MakeThumbnail(fromFile string, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := imaging.Open(fromFile)
	if err != nil {
		return nil, err
	}

	image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	return
}

func MakeThumbnailFromReader(reader io.Reader, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	return
}

func MakeFromReader(reader io.Reader, t string, w, h int) (image *image.NRGBA, err error) {
	srcImage, err := Open(reader)
	if err != nil {
		return nil, err
	}

	switch t {
	case "thumbnail":
		image = imaging.Thumbnail(srcImage, w, h, imaging.Lanczos)
	case "resize":
		image = imaging.Resize(srcImage, w, h, imaging.Lanczos)
	case "fit":
		image = imaging.Fit(srcImage, w, h, imaging.Lanczos)
	default:
		panic("只支持thumbnail, resize, fit")
	}

	return
}

// Open loads an image from file
func Open(reader io.Reader) (img image.Image, err error) {
	img, _, err = image.Decode(reader)
	if err != nil {
		return
	}

	img = toNRGBA(img)
	return
}

// This function used internally to convert any image type to NRGBA if needed.
func toNRGBA(img image.Image) *image.NRGBA {
	srcBounds := img.Bounds()
	if srcBounds.Min.X == 0 && srcBounds.Min.Y == 0 {
		if src0, ok := img.(*image.NRGBA); ok {
			return src0
		}
	}
	return imaging.Clone(img)
}

func ToJSON(o interface{}) string {
	b, err := json.Marshal(o)
	gotang.AssertNoError(err, "ToJSON")

	return string(b)
}

func FromJSON(s string, o interface{}) {
	err := json.Unmarshal([]byte(s), o)
	gotang.AssertNoError(err, "FromJSON)")
}

//panicable
type CacheDataLoader func(string) interface{}

func Cache(key string, target interface{}, loader CacheDataLoader) {
	CacheWithExpires(key, target, loader, cache.FOREVER)
}

var cacheKeys []string

func GetCacheKeys() []string {
	return cacheKeys
}

func CacheWithExpires(key string, target interface{}, loader CacheDataLoader, expires time.Duration) {
	if err := cache.Get(key, target); err != nil {
		values := loader(key)
		setValueToAddress(target, values)
		cacheKeys = append(cacheKeys, key)
		go cache.Set(key, values, expires)
	}
}

func ClearCache(key string) {
	go cache.Delete(key)
}

func setValueToAddress(target interface{}, value interface{}) {
	p := reflect.ValueOf(target)
	gotang.Assert(p.Type().Kind() == reflect.Ptr, "target should be Pointer")

	v := p.Elem()

	gotang.Assert(v.CanSet(), "target should be CanSet")
	v.Set(reflect.ValueOf(value))
}
