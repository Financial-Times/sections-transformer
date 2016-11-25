package main

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetSections(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseURL  string
		terms    []term
		sections []sectionLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/sections/",
			[]term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "mNGQ2MWQ0NDMtMDc5Mi00NWExLTlkMGQtNWZhZjk0NGExOWU2-Z2VucVz"}},
			[]sectionLink{sectionLink{APIURL: "localhost:8080/transformers/sections/20ddda23-a1bb-3530-88aa-60232583895a"},
				sectionLink{APIURL: "localhost:8080/transformers/sections/cfd7a2d5-bc8f-3585-b98a-db69f7b8cfea"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/sections/", []term{}, []sectionLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSectionService(&repo, test.baseURL, "Sections", 10000)
		expectedSections, found := service.getSections()
		assert.Equal(test.sections, expectedSections, fmt.Sprintf("%s: Expected sections link incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

func TestGetSectionByUuid(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		terms   []term
		uuid    string
		section section
		found   bool
		err     error
	}{
		{"Success", []term{term{CanonicalName: "Z_Archive", RawID: "b8337559-ac08-3404-9025-bad51ebe2fc7"}, term{CanonicalName: "Feature", RawID: "TkdRMk1XUTBORE10TURjNU1pMDBOV0V4TFRsa01HUXROV1poWmprME5HRXhPV1UyLVoyVnVjbVZ6-U2VjdGlvbnM=]"}},
			"20ddda23-a1bb-3530-88aa-60232583895a", getDummySection("20ddda23-a1bb-3530-88aa-60232583895a", "Z_Archive", "YjgzMzc1NTktYWMwOC0zNDA0LTkwMjUtYmFkNTFlYmUyZmM3-U2VjdGlvbnM="), true, nil},
		{"Not found", []term{term{CanonicalName: "Z_Archive", RawID: "845dc7d7-ae89-4fed-a819-9edcbb3fe507"}, term{CanonicalName: "Feature", RawID: "NGQ2MWdefsdfsfcmVz"}},
			"some uuid", section{}, false, nil},
		{"Error on init", []term{}, "some uuid", section{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{terms: test.terms, err: test.err}
		service, err := newSectionService(&repo, "", "Sections", 10000)
		expectedSection, found := service.getSectionByUUID(test.uuid)
		assert.Equal(test.section, expectedSection, fmt.Sprintf("%s: Expected section incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	terms []term
	err   error
}

func (d *dummyRepo) GetTmeTermsFromIndex(startRecord int) ([]interface{}, error) {
	if startRecord > 0 {
		return nil, d.err
	}
	var interfaces = make([]interface{}, len(d.terms))
	for i, data := range d.terms {
		interfaces[i] = data
	}
	return interfaces, d.err
}
func (d *dummyRepo) GetTmeTermById(uuid string) (interface{}, error) {
	return d.terms[0], d.err
}

func getDummySection(uuid string, prefLabel string, tmeID string) section {
	return section{
		UUID:      uuid,
		PrefLabel: prefLabel,
		Type:      "Section",
		AlternativeIdentifiers: alternativeIdentifiers{TME: []string{tmeID}, Uuids: []string{uuid}}}
}
