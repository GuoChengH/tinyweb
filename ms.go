package tinyweb

import (
	"fmt"
	"net/http"
)

const ANY = "ANY"

type HandleFunc func(ctx *Context)

type MiddlewareFunc func(handleFunc HandleFunc) HandleFunc

type routerGroup struct {
	name             string
	handleFuncMap    map[string]map[string]HandleFunc
	handlerMethodMap map[string][]string
	treeNode         *treeNode
	middlewares      []MiddlewareFunc
}

func (r *routerGroup) Use(middlewareFunc ...MiddlewareFunc) {
	r.middlewares = append(r.middlewares, middlewareFunc...)
}

func (r *routerGroup) methodHandle(h HandleFunc, ctx *Context) {
	// 中间件

	if r.middlewares != nil {
		for _, middleware := range r.middlewares {
			h = middleware(h)
		}
	}
	h(ctx)

}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:             name,
		handleFuncMap:    make(map[string]map[string]HandleFunc),
		handlerMethodMap: make(map[string][]string),
		treeNode:         &treeNode{name: "/", children: make([]*treeNode, 0)},
	}
	r.routerGroups = append(r.routerGroups, group)
	return group
}
func (r *routerGroup) handle(name string, method string, handleFunc HandleFunc) {
	_, ok := r.handleFuncMap[name]
	if !ok {
		r.handleFuncMap[name] = make(map[string]HandleFunc)
	}
	_, ok = r.handleFuncMap[name][method]
	if ok {
		panic("method already exists")
	}
	r.handleFuncMap[name][method] = handleFunc
	// 这里注册路由name为
	r.treeNode.Put(name)
}
func (r *routerGroup) Any(name string, handleFunc HandleFunc) {
	r.handle(name, ANY, handleFunc)
}
func (r *routerGroup) Get(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodGet, handleFunc)
}
func (r *routerGroup) Post(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPost, handleFunc)
}
func (r *routerGroup) Put(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPut, handleFunc)
}
func (r *routerGroup) Delete(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodDelete, handleFunc)
}
func (r *routerGroup) Patch(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodPatch, handleFunc)
}
func (r *routerGroup) Options(name string, handleFunc HandleFunc) {
	r.handle(name, http.MethodOptions, handleFunc)
}

type Engine struct {
	router
}

func New() *Engine {
	return &Engine{
		router: router{},
	}
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
			ctx := &Context{w, r}
			handle, ok := group.handleFuncMap[node.routerName][ANY]
			if ok {
				group.methodHandle(handle, ctx)
				return
			}
			handle, ok = group.handleFuncMap[node.routerName][method]
			if ok {
				group.methodHandle(handle, ctx)
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
