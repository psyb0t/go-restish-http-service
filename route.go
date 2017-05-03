package httpservice

import (
	"encoding/json"
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

func (r *Route) SuccessResponse() {
	r.writer.WriteHeader(200)
}

func (r *Route) SuccessObjectResponse(object interface{}) error {
	r.writer.WriteHeader(200)

	json_object, err := json.Marshal(object)
	if err != nil {
		return err
	}

	r.writer.Write(json_object)

	return nil
}

func (r *Route) ErrorResponse(reason string) {
	r.writer.WriteHeader(500)
	r.writer.Write([]byte(reason))
}
