package jserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
)

type JsonRouter struct {
	Router *mux.Router
}

func (r *JsonRouter) Add(path string, file string, x interface{}) error {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	json.Unmarshal(raw, &x)
	r.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		buff, _ := json.Marshal(x)
		fmt.Fprintf(w, string(buff))
	})
	return nil
}

func (r *JsonRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func NewJsonRouter() *JsonRouter {
	router := &JsonRouter{}
	router.Router = mux.NewRouter()
	return router
}
