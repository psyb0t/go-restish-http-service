package simplehttp

import (
	"errors"
	"fmt"
	"net/http"
	"path"

	"github.com/julienschmidt/httprouter"
)

type RouteMethod func(*Route)
type RouteMiddleware func(*Route)

type HttpService struct {
	ListenAddress string
	DB            *DB
	Auth          *Auth
	serverName    string
	headers       map[string]string
	router        *httprouter.Router
}

type routeGroupItem struct {
	Path       string
	HttpMethod string
	Method     RouteMethod
	Middleware []RouteMiddleware
}

func NewRouteGroupItem(path string, http_method string,
	method RouteMethod, middleware ...RouteMiddleware) *routeGroupItem {

	return &routeGroupItem{
		Path:       path,
		HttpMethod: http_method,
		Method:     method,
		Middleware: middleware,
	}
}

func New(address, db_path string) *HttpService {
	service := &HttpService{
		ListenAddress: address,
		router:        httprouter.New(),
	}

	if db_path != "" {
		service.DB = NewDB(db_path)
	}

	service.Auth = NewAuth(service)

	service.headers = make(map[string]string)

	return service
}

func (s *HttpService) SetServerName(name string) {
	s.serverName = name
}

func (s *HttpService) AddHeader(name, value string) {
	s.headers[name] = value
}

func (s *HttpService) AddRouteGroup(prefix string, routes ...*routeGroupItem) {
	for _, route := range routes {
		s.AddRoute(route.HttpMethod,
			path.Join(prefix, route.Path), route.Method, route.Middleware...)
	}
}

func (s *HttpService) AddRoute(http_method, url_path string,
	fn RouteMethod, middleware ...RouteMiddleware) error {

	handler := func(w http.ResponseWriter, r *http.Request,
		p httprouter.Params) {

		w.Header().Set("Server", s.serverName)

		for key, val := range s.headers {
			w.Header().Set(key, val)
		}

		route, err := NewRoute(p, s, w, r)
		if err != nil {
			route.ErrorResponse(err.Error())
			return
		}

		for _, m := range middleware {
			if m != nil {
				m(route)
			}
		}

		fn(route)
	}

	switch http_method {
	case "GET":
		s.router.GET(url_path, handler)
		return nil
	case "POST":
		s.router.POST(url_path, handler)
		return nil
	case "PUT":
		s.router.PUT(url_path, handler)
		return nil
	case "DELETE":
		s.router.DELETE(url_path, handler)
		return nil
	}

	return errors.New(fmt.Sprintf(
		"Unsupported %s http method for route %s", http_method, url_path))
}

func (s *HttpService) AddStaticRoute(url_path, fs_path string) {
	s.router.ServeFiles(
		path.Join(url_path, "*filepath"), http.Dir(fs_path))
}

func (s *HttpService) Start() {
	panic(http.ListenAndServe(s.ListenAddress, s.router))
}
