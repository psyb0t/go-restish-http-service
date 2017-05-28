package restishhttpservice

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type RouteMethod func(*Route)

type HttpService struct {
	ListenAddress string
	DB            *DB
	Auth          *Auth
	serverName    string
	headers       map[string]string
	router        *httprouter.Router
}

func New(address, db_path string) *HttpService {
	service := &HttpService{
		ListenAddress: address,
		router:        httprouter.New(),
		DB:            NewDB(db_path),
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

func (s *HttpService) AddRoute(method, path string, fn RouteMethod) error {
	handler := func(w http.ResponseWriter, r *http.Request,
		p httprouter.Params) {

		w.Header().Set("Server", s.serverName)

		for key, val := range s.headers {
			w.Header().Set(key, val)
		}

		fn(&Route{
			Params:  p,
			Service: s,
			writer:  w,
			request: r,
		})
	}

	switch method {
	case "GET":
		s.router.GET(path, handler)
		return nil
	case "POST":
		s.router.POST(path, handler)
		return nil
	case "PUT":
		s.router.PUT(path, handler)
		return nil
	case "DELETE":
		s.router.DELETE(path, handler)
		return nil
	}

	return errors.New("Unsupported route method")
}

func (s *HttpService) Start() {
	panic(http.ListenAndServe(s.ListenAddress, s.router))
}
