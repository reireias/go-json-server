package jserver

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

type Friends struct {
	ID     string  `json:"id"`
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float64 `json:"height"`
}

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

	err = checkResponseSize(body, 2)
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 0, "1", 14, "kaban")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 1, "2", 15, "serval")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
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

	err = checkSingleResponse(body, "1", 14, "kaban")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
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
		t.Error("response body invalid: " + string(body))
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

	err = checkResponseSize(body, 1)
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 0, "1", 14, "kaban")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
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
		t.Error("response body invalid: " + string(body))
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

	err = checkResponseSize(body, 1)
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 0, "1", 14, "kaban")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
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
		t.Error("response body invalid: " + string(body))
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

	err = checkResponseSize(body, 2)
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 0, "1", 14, "kaban")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
	err = checkResponse(body, 1, "2", 15, "serval")
	if err != nil {
		t.Error(err.Error() + "\n" + string(body))
		return
	}
}

func checkResponseSize(body []byte, expectSize int) error {
	data, err := unmarshalFriendsJSON(body)
	if err != nil {
		return err
	}
	if len(data) != expectSize {
		return errors.New("response body size invalid")
	}
	return nil
}

func checkResponse(body []byte, index int, id string, age int, name string) error {
	data, err := unmarshalFriendsJSON(body)
	if err != nil {
		return err
	}
	if data[index].ID == id && data[index].Age == age && data[index].Name == name {
		return nil
	}
	return errors.New("response body invalid")
}

func checkSingleResponse(body []byte, id string, age int, name string) error {
	data := Friends{}
	err := json.Unmarshal(body, &data)
	if err != nil {
		return err
	}
	if data.ID == id && data.Age == age && data.Name == name {
		return nil
	}
	return errors.New("response body invalid")
}

func unmarshalFriendsJSON(body []byte) ([]Friends, error) {
	data := []Friends{}
	err := json.Unmarshal(body, &data)
	return data, err
}
