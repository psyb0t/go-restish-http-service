package restishhttpservice

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type Route struct {
	Params  httprouter.Params
	Service *HttpService
	writer  http.ResponseWriter
	request *http.Request
}

func (r *Route) SetHeader(name, value string) {
	r.writer.Header().Set(name, value)
}

func (r *Route) GetHeader(name string) string {
	return r.request.Header.Get(name)
}

func (r *Route) GetBody() ([]byte, error) {
	body_data, err := ioutil.ReadAll(r.request.Body)
	if err != nil {
		return nil, err
	}

	return body_data, nil
}

func (r *Route) IsAuthorized() bool {
	auth_key := r.GetHeader("X-Auth-Key")
	if auth_key == "" {
		return false
	}

	if !r.Service.Auth.IsValidKey(auth_key) {
		return false
	}

	return true
}

func (r *Route) SuccessResponse() {
	r.SetHeader("Content-Type", "text/plain")
	r.writer.WriteHeader(200)
}

func (r *Route) SuccessObjectResponse(object interface{}) error {
	r.writer.WriteHeader(200)

	json_object, err := json.Marshal(object)
	if err != nil {
		return err
	}

	r.SetHeader("Content-Type", "application/json")

	r.writer.Write(json_object)

	return nil
}

func (r *Route) SuccessStringResponse(text string) {
	r.writer.WriteHeader(200)
	r.SetHeader("Content-Type", "text/plain")
	r.writer.Write([]byte(text))
}

func (r *Route) ForbiddenResponse() {
	r.writer.WriteHeader(403)
	r.SetHeader("Content-Type", "text/plain")
	r.writer.Write([]byte("Permission denied!"))
}

func (r *Route) ErrorResponse(reason string) {
	r.writer.WriteHeader(500)
	r.SetHeader("Content-Type", "text/plain")
	r.writer.Write([]byte(reason))
}

func (r *Route) RedirectResponse(location string) {
	r.SetHeader("Location", location)
	r.writer.WriteHeader(301)

}
