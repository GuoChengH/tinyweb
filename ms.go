package tinyweb

import (
	"fmt"
	"net/http"
)

const ANY = "ANY"

type HandleFunc func(ctx *Context)

type routerGroup struct {
	name             string
	handleFuncMap    map[string]map[string]HandleFunc
	handlerMethodMap map[string][]string
}

type router struct {
	routerGroups []*routerGroup
}

func (r *router) Group(name string) *routerGroup {
	group := &routerGroup{
		name:             name,
		handleFuncMap:    make(map[string]map[string]HandleFunc),
		handlerMethodMap: make(map[string][]string),
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
	method := r.Method
	for _, group := range e.routerGroups {
		for name, methodHandle := range group.handleFuncMap {
			url := "/" + group.name + "/" + name
			if r.RequestURI == url {
				ctx := &Context{w, r}
				handle, ok := methodHandle[ANY]
				if ok {
					handle(ctx)
					return
				}
				handle, ok = methodHandle[method]
				if ok {
					handle(ctx)
					return
				}
				w.WriteHeader(http.StatusMethodNotAllowed)
				fmt.Fprintf(w, "%s %s not allowed \n", r.RequestURI, method)
			}
		}
	}
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%s %s not found \n", r.RequestURI, method)
}
