package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

const testUUID = "bba39990-c78d-3629-ae83-808c333c6dbc"
const getSectionsResponse = `[{"apiUrl":"http://localhost:8080/transformers/sections/bba39990-c78d-3629-ae83-808c333c6dbc"}]`
const getSectionByUUIDResponse = `{"uuid":"bba39990-c78d-3629-ae83-808c333c6dbc","alternativeIdentifiers":{"TME":["MTE3-U3ViamVjdHM="],"uuids":["bba39990-c78d-3629-ae83-808c333c6dbc"]},"prefLabel":"Global Sections","type":"Section"}`
const getSectionsCountResponse = `1`
const getSectionsIdsResponse = `{"id":"bba39990-c78d-3629-ae83-808c333c6dbc"}`

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
		{"Success - get section by uuid", newRequest("GET", fmt.Sprintf("/transformers/sections/%s", testUUID)), &dummyService{found: true, sections: []section{getDummySection(testUUID, "Global Sections", "MTE3-U3ViamVjdHM=")}}, http.StatusOK, "application/json", getSectionByUUIDResponse},
		{"Not found - get section by uuid", newRequest("GET", fmt.Sprintf("/transformers/sections/%s", testUUID)), &dummyService{found: false, sections: []section{section{}}}, http.StatusNotFound, "application/json", ""},
		{"Success - get sections", newRequest("GET", "/transformers/sections"), &dummyService{found: true, sections: []section{section{UUID: testUUID}}}, http.StatusOK, "application/json", getSectionsResponse},
		{"Not found - get sections", newRequest("GET", "/transformers/sections"), &dummyService{found: false, sections: []section{}}, http.StatusNotFound, "application/json", ""},
		{"Test Section Count", newRequest("GET", "/transformers/sections/__count"), &dummyService{found: true, sections: []section{section{UUID: testUUID}}}, http.StatusOK, "text/plain", getSectionsCountResponse},
		{"Test Section Ids", newRequest("GET", "/transformers/sections/__ids"), &dummyService{found: true, sections: []section{section{UUID: testUUID}}}, http.StatusOK, "text/plain", getSectionsIdsResponse},
	}

	for _, test := range tests {
		rec := httptest.NewRecorder()
		router(test.dummyService).ServeHTTP(rec, test.req)
		assert.True(test.statusCode == rec.Code, fmt.Sprintf("%s: Wrong response code, was %d, should be %d", test.name, rec.Code, test.statusCode))
		assert.Equal(strings.TrimSpace(test.body), strings.TrimSpace(rec.Body.String()), fmt.Sprintf("%s: Wrong body", test.name))
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
	m.HandleFunc("/transformers/sections/__count", h.getCount).Methods("GET")
	m.HandleFunc("/transformers/sections/__ids", h.getIds).Methods("GET")
	m.HandleFunc("/transformers/sections/__reload", h.reload).Methods("POST")
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

func (s *dummyService) checkConnectivity() error {
	return nil
}

func (s *dummyService) getSectionCount() int {
	return len(s.sections)
}

func (s *dummyService) getSectionIds() []string {
	i := 0
	keys := make([]string, len(s.sections))

	for _, t := range s.sections {
		keys[i] = t.UUID
		i++
	}
	return keys
}

func (s *dummyService) reload() error {
	return nil
}
