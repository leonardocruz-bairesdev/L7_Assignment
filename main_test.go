package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/qa/get/ID1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGet2(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/qa/get/ID5", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetAll(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/qa/getAll", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	//fmt.println(w.Body.String())
	//assert.Equal(t, '[{"ID":"ID1","Question":"This is my updated question No. 1","Answer":"","Answered":false,"User":"User1"},{"ID":"ID2","Question":"This is my question No. 2 with answer removed","Answer":"","Answered":false,"User":"User1"}]', w.Body.String())
}

func TestGetAllByUser(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/qa/getAllByUser/User1", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

}

func TestGetAllByUser2(t *testing.T) {
	router := NewRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/qa/getAllByUser/User4", nil)
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

}
