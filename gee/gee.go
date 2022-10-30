package gee

import (
	"fmt"
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	//key:请求  value:请求的处理方法
	router map[string]HandlerFunc
}

//实现http.handler中ServeHTTP方法,从而可以将Engine作为http.ListenAndServe的参数启动服务
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	//查找映射表，如果找的到便执行自定义的handler方法，找不到便抛404
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND: %s\n", req.URL)
	}
}

func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	engine.router[key] = handler
}

//用户调用GET、POST方法时，会将url和handler注册到映射表当中。
// GET  请求
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

// POST  请求
func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// 调用ListenAndServe(addr string, handler Handler) 来服务启动
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

//创建并返回一个Engine
func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc, 0),
	}
}
