package tinyweb

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	"encoding/json"
	"encoding/xml"

	"github.com/GuoChengH/tinyweb/render"
)

type Context struct {
	W      http.ResponseWriter
	R      *http.Request
	engine *Engine
}

func (c *Context) HTML(status int, html string) error {
	// 状态是200 OK 默认不设置的话，如果调用了writerHeader 默认返回也是200 OK
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	c.W.WriteHeader(status)
	_, err := c.W.Write([]byte(html))
	return err
}

func (c *Context) HTMLTemplate(name string, data any, filenames ...string) error {
	// 状态是200 OK 默认不设置的话，如果调用了writerHeader 默认返回也是200 OK
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t, err := t.ParseFiles(filenames...)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) HTMLTemplateGlob(name string, data any, pattern string) error {
	// 状态是200 OK 默认不设置的话，如果调用了writerHeader 默认返回也是200 OK
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	t := template.New(name)
	t, err := t.ParseGlob(pattern)
	if err != nil {
		return err
	}
	err = t.Execute(c.W, data)
	return err
}

func (c *Context) Template(name string, data any) error {
	// 状态是200 OK 默认不设置的话，如果调用了writerHeader 默认返回也是200 OK
	c.W.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := c.engine.HTMLRender.Template.ExecuteTemplate(c.W, name, data)
	return err
}

func (c *Context) JSON(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/json; charset=utf-8")
	c.W.WriteHeader(status)
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = c.W.Write(jsonData)
	return err
}
func (c *Context) XML(status int, data any) error {
	c.W.Header().Set("Content-Type", "application/xml; charset=utf-8")
	// c.W.WriteHeader(status)
	// xmlData, err := xml.Marshal(data)
	// if err != nil {
	// 	return err
	// }
	// _, err = c.W.Write(xmlData)
	err := xml.NewEncoder(c.W).Encode(data)
	return err
}

func (c *Context) File(filename string) {
	http.ServeFile(c.W, c.R, filename)
}

// filepath是相对文件系统的路径
func (c *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		c.W.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		c.W.Header().Set("Content-Disposition", `attachment;filename*=utf-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(c.W, c.R, filepath)
}

// fs 一般是 http.Dir() 可以再相对目录下查找文件
func (c *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		c.R.URL.Path = old
	}(c.R.URL.Path)

	c.R.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(c.W, c.R)
}

func (c *Context) Redirect(status int, url string) {
	if (status < http.StatusMultipleChoices || status > http.StatusPermanentRedirect) && status != http.StatusCreated {
		panic(fmt.Sprintf("invalid redirect status code: %d", status))
	}
	http.Redirect(c.W, c.R, url, status)
}

func (c *Context) String(status int, format string, values ...any) error {
	err := c.Render(c.W, &render.String{Format: format, Data: values})
	c.W.WriteHeader(status)
	return err
}

func (c *Context) Render(w http.ResponseWriter, r render.Render) error {
	return r.Render(w)
}
