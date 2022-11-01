package gee

import (
	"log"
	"net/http"
	"strings"
)

type HandlerFunc func(c *Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 中间件
	parent      *RouterGroup  // support nesting
	engine      *Engine       // all groups share a Engine instance
}

type Engine struct {
	//key:请求  value:请求的处理方法
	//router map[string]HandlerFunc
	router *router
	groups []*RouterGroup
	*RouterGroup
}

//实现http.handler中ServeHTTP方法,从而可以将Engine作为http.ListenAndServe的参数启动服务
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	/*
		key := req.Method + "-" + req.URL.Path
		//查找映射表，如果找的到便执行自定义的handler方法，找不到便抛404
		if handler, ok := engine.router[key]; ok {
			handler(w, req)
		} else {
			fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
		}
	*/
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	//key := method + "-" + pattern
	//engine.router[key] = handler

	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

//用户调用GET、POST方法时，会将url和handler注册到映射表当中。
// GET 请求
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// POST 请求
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

// 调用ListenAndServe(addr string, handler Handler) 来服务启动
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//创建并返回一个Engine
func New() *Engine {
	//return &Engine{
	//	router: newRouter(),
	//}
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}
