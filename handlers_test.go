package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSectionsResponse = "[{\"apiUrl\":\"http://localhost:8080/transformers/sections/bba39990-c78d-3629-ae83-808c333c6dbc\"}]\n"
const getSectionByUUIDResponse = "{\"uuid\":\"bba39990-c78d-3629-ae83-808c333c6dbc\",\"canonicalName\":\"Metals Markets\",\"tmeIdentifier\":\"MTE3-U3ViamVjdHM=\",\"type\":\"Section\"}\n"

func TestHandlers(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name         string
		req          *http.Request
		dummyService sectionService
		statusCode   int
		contentType  string // Contents of the Content-Type header
		body         string
	}{
		{"Success - get section by uuid", newRequest("GET", fmt.Sprintf("/transformers/sections/%s", testUUID)), &dummyService{found: true, sections: []section{section{UUID: testUUID, CanonicalName: "Metals Markets", TmeIdentifier: "MTE3-U3ViamVjdHM=", Type: "Section"}}}, http.StatusOK, "application/json", getSectionByUUIDResponse},
		{"Not found - get section by uuid", newRequest("GET", fmt.Sprintf("/transformers/sections/%s", testUUID)), &dummyService{found: false, sections: []section{section{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get sections", newRequest("GET", "/transformers/sections"), &dummyService{found: true, sections: []section{section{UUID: testUUID}}}, http.StatusOK, "application/json", getSectionsResponse},
		{"Not found - get sections", newRequest("GET", "/transformers/sections"), &dummyService{found: false, sections: []section{}}, http.StatusNotFound, "application/json", ""},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(test.body, rec.Body.String(), fmt.Sprintf("%s: Wrong body", test.name))
	}
}

func newRequest(method, url string) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic(err)
	}
	return req
}

func router(s sectionService) *mux.Router {
	m := mux.NewRouter()
	h := newSectionsHandler(s)
	m.HandleFunc("/transformers/sections", h.getSections).Methods("GET")
	m.HandleFunc("/transformers/sections/{uuid}", h.getSectionByUUID).Methods("GET")
	return m
}

type dummyService struct {
	found    bool
	sections []section
}

func (s *dummyService) getSections() ([]sectionLink, bool) {
	var sectionLinks []sectionLink
	for _, sub := range s.sections {
		sectionLinks = append(sectionLinks, sectionLink{APIURL: "http://localhost:8080/transformers/sections/" + sub.UUID})
	}
	return sectionLinks, s.found
}

func (s *dummyService) getSectionByUUID(uuid string) (section, bool) {
	return s.sections[0], s.found
}
