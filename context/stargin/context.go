package stargin

import (
	json2 "encoding/json"
	"fmt"
	"net/http"
)

type M map[string]interface{} //别名

type Context struct {
	Writer http.ResponseWriter
	Request *http.Request
	//response info
	StatusCode int
	//request info
	Path   string
	Method string
}

func newContext(w http.ResponseWriter,r *http.Request) *Context {
	return &Context{
		Writer: w,
		Request: r,
		Path: r.URL.Path,
		Method: r.Method,
	}
}

func (c *Context)PostForm(key string) string  {
	return c.Request.FormValue(key)
}

func (c *Context)Query(key string)string {
	return c.Request.URL.Query().Get(key)
}

func (c *Context)SetStatus(code int) {
	c.StatusCode =code
	c.Writer.WriteHeader(code)
}

func (c *Context)SetHeader(key string,value string){
	c.Writer.Header().Set(key,value)
}

//快速构造String/Data/JSON/HTML响应的方法
func (c *Context)String(code int,key string,value ...interface{})  {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write([]byte(fmt.Sprintf(key , value...)))
}

func (c *Context)Json(code int,json interface{})  {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json2.NewEncoder(c.Writer)
	if err:= encoder.Encode(json);err!=nil{
		panic(err)
	}
}

func (c *Context)Html(code int,html string)  {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.Writer.Write([]byte(html))
}

func (c *Context)Data(code int ,data []byte)  {
	c.SetStatus(code)
	c.Writer.Write(data)
}