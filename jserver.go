package jserver

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/koron/go-dproxy"
)

// JSONRouter is router for JSON files.
type JSONRouter struct {
	Router *mux.Router
}

// Add path that returns data from JSON file. Retruns an error.
func (r *JSONRouter) Add(path string, file string) error {
	var data interface{}
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	json.Unmarshal(raw, &data)
	dataMap := map[string]interface{}{}
	keys := make(map[string]struct{})
	for _, v := range data.([]interface{}) {
		proxy := dproxy.New(v)

		id, _ := proxy.M("id").String()
		dataMap[id] = v

		data, _ := proxy.Map()
		for k := range data {
			keys[k] = struct{}{}
		}
	}

	// /path
	r.Router.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		parameters := r.URL.Query()
		filteredData := filter(data.([]interface{}), func(record interface{}) bool {
			for pkey, pValue := range parameters {
				_, ok := keys[pkey]
				if !ok {
					continue
				}
				dValue, _ := dproxy.New(record).M(pkey).Value()
				switch value := dValue.(type) {
				case float64:
					return strconv.Itoa(int(value)) == pValue[0]
				case string:
					return value == pValue[0]
				}
			}
			return true
		})

		buff, _ := json.Marshal(filteredData)
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

func (r *JSONRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.Router.ServeHTTP(w, req)
}

func filter(list []interface{}, f func(interface{}) bool) []interface{} {
	result := make([]interface{}, 0)
	for _, v := range list {
		if f(v) {
			result = append(result, v)
		}
	}
	return result
}

// NewJSONRouter returns a new router instance.
func NewJSONRouter() *JSONRouter {
	router := &JSONRouter{}
	router.Router = mux.NewRouter()
	return router
}
