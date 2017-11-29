package jserver

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListPath(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

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

func TestIDPath(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

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

func TestIDPathNotFound(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test/9999")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `{"error":"Not Found."}` {
		t.Error("response body invalid")
		return
	}
}

func TestFilterString(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test?name=kaban")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `[{"age":14,"id":"1","name":"kaban"}]` {
		t.Error("response body invalid")
		return
	}
}

func TestFilterStringNoRecord(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test?name=invalid")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `[]` {
		t.Error("response body invalid")
		return
	}
}

func TestFilterNumber(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test?age=14")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `[{"age":14,"id":"1","name":"kaban"}]` {
		t.Error("response body invalid")
		return
	}
}

func TestFilterNumberInvalid(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test?age=invalid")
	if err != nil {
		t.Error("unexpected")
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error("unexpected")
		return
	}

	if string(body) != `[]` {
		t.Error("response body invalid")
		return
	}
}

func TestFilterNotExistsKey(t *testing.T) {
	jsonRouter := NewJSONRouter()
	jsonRouter.Add("/test", "./test.json")

	ts := httptest.NewServer(jsonRouter)
	defer ts.Close()

	res, err := http.Get(ts.URL + "/test?invalid=14")
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
