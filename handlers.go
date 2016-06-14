package main

import (
	"encoding/json"
	"fmt"
	"github.com/Financial-Times/go-fthealth/v1a"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"net/http"
)

type sectionsHandler struct {
	service sectionService
}

// HealthCheck does something
func (h *sectionsHandler) HealthCheck() v1a.Check {
	return v1a.Check{
		BusinessImpact:   "Unable to respond to request for the section data from TME",
		Name:             "Check connectivity to TME",
		PanicGuide:       "https://sites.google.com/a/ft.com/ft-technology-service-transition/home/run-book-library/sections-transfomer",
		Severity:         1,
		TechnicalSummary: "Cannot connect to TME to be able to supply sections",
		Checker:          h.checker,
	}
}

// Checker does more stuff
func (h *sectionsHandler) checker() (string, error) {
	err := h.service.checkConnectivity()
	if err == nil {
		return "Connectivity to TME is ok", err
	}
	return "Error connecting to TME", err
}

func newSectionsHandler(service sectionService) sectionsHandler {
	return sectionsHandler{service: service}
}

func (h *sectionsHandler) getSections(writer http.ResponseWriter, req *http.Request) {
	obj, found := h.service.getSections()
	writeJSONResponse(obj, found, writer)
}

func (h *sectionsHandler) getSectionByUUID(writer http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	uuid := vars["uuid"]

	obj, found := h.service.getSectionByUUID(uuid)
	writeJSONResponse(obj, found, writer)
}

//GoodToGo returns a 503 if the healthcheck fails - suitable for use from varnish to check availability of a node
func (h *sectionsHandler) GoodToGo(writer http.ResponseWriter, req *http.Request) {
	if _, err := h.checker(); err != nil {
		writer.WriteHeader(http.StatusServiceUnavailable)
	}
}

func writeJSONResponse(obj interface{}, found bool, writer http.ResponseWriter) {
	writer.Header().Add("Content-Type", "application/json")

	if !found {
		writer.WriteHeader(http.StatusNotFound)
		return
	}

	enc := json.NewEncoder(writer)
	if err := enc.Encode(obj); err != nil {
		log.Errorf("Error on json encoding=%v\n", err)
		writeJSONError(writer, err.Error(), http.StatusInternalServerError)
		return
	}
}

func writeJSONError(w http.ResponseWriter, errorMsg string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprintln(w, fmt.Sprintf("{\"message\": \"%s\"}", errorMsg))
}
