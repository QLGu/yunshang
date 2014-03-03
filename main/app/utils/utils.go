package utils

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"time"

	"github.com/itang/gotang"
	"github.com/nu7hatch/gouuid"
	"github.com/robfig/revel"
)

func Uuid() string {
	u4, err := uuid.NewV4()
	gotang.AssertNoError(err)
	return u4.String()
}

func Sha1(content string) string {
	h := sha1.New()
	io.WriteString(h, content)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func RandomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return min + r.Intn(max-min)
}

func RenderTemplate(templatePath string, data interface{}) string {
	template, err := revel.MainTemplateLoader.Template(templatePath)
	gotang.AssertNoError(err)

	var b bytes.Buffer
	err = template.Render(&b, data)
	gotang.AssertNoError(err)

	return b.String()
}

func DoIOWithTimeout(f func() error, t time.Duration) error {
	timeout := time.NewTicker(t)
	defer timeout.Stop()
	done := make(chan error)
	go func() {
		done <- f()
	}()
	select {
	case <-timeout.C:
		return errors.New("timeout")
	case err := <-done:
		return err
	}
}
