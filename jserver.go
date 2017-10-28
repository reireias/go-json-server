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

func (r *JsonRouter) Add(path string, file string) error {
	var data interface{}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	json.Unmarshal(raw, &data)
	dataMap := map[string]interface{}{}
	for _, v := range data.([]interface{}) {
		id, _ := dproxy.New(v).M("id").String()
		dataMap[id] = v
	}

	// /path
	r.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		buff, _ := json.Marshal(data)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, string(buff))
	})

	// /path/{id}
	r.Router.HandleFunc(path+"/{id}", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		targetID := vars["id"]
		w.WriteHeader(http.StatusOK)
		targetData, ok := dataMap[targetID]
		if ok {
			buff, _ := json.Marshal(targetData)
			fmt.Fprintf(w, string(buff))
		} else {
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprint(w, "{\"error\":\"Not Found.\"}")
		}
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
