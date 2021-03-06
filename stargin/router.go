package stargin

import (
	"fmt"
	"net/http"
	"strings"
)

//路由注册解析

type Router struct {
	roots    map[string]*node
	handlers map[string]Handlerfunc
}

func newRouter() *Router {
	return &Router{
		roots:    make(map[string]*node),
		handlers: make(map[string]Handlerfunc),
	}
}

func (r *Router) handle(c *Context) {
	if n, params := r.getRoute(c.Method, c.Path); n != nil {
		c.Params = params
		key := c.Method + ">>" + n.pattern
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.Writer.WriteHeader(http.StatusNotFound)
		c.handlers = append(c.handlers, func(c *Context) {
			fmt.Fprintf(c.Writer, "404 NOT FOUND:%s\n", c.Request.URL)
		})
	}
	c.Next()
}

func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' { //只有一个 * 才被允许
				break
			}
		}
	}
	return parts
}

func (r *Router) addRoute(method string, pattern string, handler Handlerfunc) {
	parts := parsePattern(pattern)
	key := method + ">>" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}

	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *Router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	if root, ok := r.roots[method]; ok {
		if n := root.search(searchParts, 0); n != nil {
			parts := parsePattern(n.pattern)
			for index, part := range parts {
				if part[0] == ':' {
					params[part[1:]] = searchParts[index]
				}
				if part[0] == '*' && len(part) > 1 {
					params[part[1:]] = strings.Join(searchParts[index:], "/")
					break
				}
			}
			return n, params
		}
	}
	return nil, nil
}

func (r *Router) getRoutes(method string) []*node {
	if root, ok := r.roots[method]; ok {
		nodes := make([]*node, 0)
		root.travel(&nodes)
	}
	return nil
}
