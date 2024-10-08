package tinyweb

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/GuoChengH/tinyweb/render"
)

const ANY = "ANY"

type HandleFunc func(ctx *Context)

type MiddlewareFunc func(handleFunc HandleFunc) HandleFunc

type routerGroup struct {
	name               string
	handleFuncMap      map[string]map[string]HandleFunc
	middlewaresFuncMap map[string]map[string][]MiddlewareFunc
	handlerMethodMap   map[string][]string
	treeNode           *treeNode
	middlewares        []MiddlewareFunc
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(name string, method string, h HandleFunc, ctx *Context) {
	// 组中间件
	if r.middlewares != nil {
		for _, middleware := range r.middlewares {
			h = middleware(h)
		}
	}
	// 组路由级别
	middlewareFuncs := r.middlewaresFuncMap[name][method]
	if middlewareFuncs != nil {
		for _, middlewareFunc := range middlewareFuncs {
			h = middlewareFunc(h)
		}
	}
	h(ctx)

}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:               name,
		middlewaresFuncMap: make(map[string]map[string][]MiddlewareFunc),
		handleFuncMap:      make(map[string]map[string]HandleFunc),
		handlerMethodMap:   make(map[string][]string),
		treeNode:           &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}
func (r *routerGroup) handle(name string, method string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	_, ok := r.handleFuncMap[name]
	if !ok {
		r.handleFuncMap[name] = make(map[string]HandleFunc)
		r.middlewaresFuncMap[name] = make(map[string][]MiddlewareFunc)
	}
	_, ok = r.handleFuncMap[name][method]
	if ok {
		panic("method already exists")
	}
	r.handleFuncMap[name][method] = handleFunc
	r.middlewaresFuncMap[name][method] = append(r.middlewaresFuncMap[name][method], middlewares...)
	// 这里注册路由name为
	r.treeNode.Put(name)
}
func (r *routerGroup) Any(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, ANY, handleFunc, middlewares...)
}
func (r *routerGroup) Get(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodGet, handleFunc, middlewares...)
}
func (r *routerGroup) Post(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodPost, handleFunc, middlewares...)
}
func (r *routerGroup) Put(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodPut, handleFunc, middlewares...)
}
func (r *routerGroup) Delete(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodDelete, handleFunc, middlewares...)
}
func (r *routerGroup) Patch(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodPatch, handleFunc, middlewares...)
}
func (r *routerGroup) Options(name string, handleFunc HandleFunc, middlewares ...MiddlewareFunc) {
	r.handle(name, http.MethodOptions, handleFunc, middlewares...)
}

type Engine struct {
	router
	funcMap    template.FuncMap
	HTMLRender render.HTMLRender
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
}
func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}
func (e *Engine) LoadTemplate(pattern string) {
	t := template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
	e.HTMLRender = render.HTMLRender{Template: t}
}
func (e *Engine) SetHTMLTemplate(t *template.Template) {
	e.HTMLRender = render.HTMLRender{Template: t}
}
func (e *Engine) Run() {
	http.Handle("/", e)
	err := http.ListenAndServe(":8111", nil)
	if err != nil {
		panic(err)
	}
}
func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e.httpRequestHandle(w, r)
}
func (e *Engine) httpRequestHandle(w http.ResponseWriter, r *http.Request) {
	method := r.Method
	for _, group := range e.routerGroups {
		routerName := SubStringLast(r.RequestURI, group.name)
		node := group.treeNode.Get(routerName)
		if node != nil && node.isEnd {
			// 路由匹配
			ctx := &Context{
				W:      w,
				R:      r,
				engine: e,
			}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.methodHandle(node.routerName, ANY, handle, ctx)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(node.routerName, method, handle, ctx)
				return
			}
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			return
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s %s not found \n", r.RequestURI, method)
}
