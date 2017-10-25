package jserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/koron/go-dproxy"
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
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(buff))
	})
	r.Router.HandleFunc(path+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		targetID := vars["id"]
		w.WriteHeader(http.StatusOK)
		for _, v := range x.([]interface{}) {
			id, _ := dproxy.New(v).M("id").String()
			if targetID == id {
				buff, _ := json.Marshal(v)
				fmt.Fprintf(w, string(buff))
				return
			}
		}
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "{\"error\":\"Not Found.\"}")
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
