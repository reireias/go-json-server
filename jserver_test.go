package jserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPath(t *testing.T) {
	type Sample struct {
		ID   string `json:"id"`
		Age  int    `json:"age"`
		Name string `json:"name"`
	}

	var s []Sample
	jsonRouter := NewJsonRouter()
	jsonRouter.Add("/test", "./test.json", s)

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `[{"age":14,"id":"1","name":"kaban"},{"age":15,"id":"2","name":"serval"}]` {
		t.Error("response body invalid")
		return
	}
}

func TestID(t *testing.T) {
	type Sample struct {
		ID   string `json:"id"`
		Age  int    `json:"age"`
		Name string `json:"name"`
	}

	var s []Sample
	jsonRouter := NewJsonRouter()
	jsonRouter.Add("/test", "./test.json", s)

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test/1")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `{"age":14,"id":"1","name":"kaban"}` {
		t.Error("response body invalid")
		return
	}
}
