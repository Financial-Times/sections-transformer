package main

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetSections(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name     string
		baseURL  string
		tax      taxonomy
		sections []sectionLink
		found    bool
		err      error
	}{
		{"Success", "localhost:8080/transformers/sections/",
			taxonomy{Terms: []term{term{CanonicalName: "Comment", ID: "MTE2-U2VjdGlvbnM=", Children: children{[]term{term{CanonicalName: "Blogs", ID: "MTIx-U2VjdGlvbnM="}}}}}},
			[]sectionLink{sectionLink{APIURL: "localhost:8080/transformers/sections/38dbd827-fedc-3ebe-919f-e64cf55ea959"},
				sectionLink{APIURL: "localhost:8080/transformers/sections/d22ff01b-f7b0-3c84-8543-c92e346b4585"}}, true, nil},
		{"Error on init", "localhost:8080/transformers/sections/", taxonomy{}, []sectionLink(nil), false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newSectionService(&repo, test.baseURL)
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
		tax     taxonomy
		uuid    string
		section section
		found   bool
		err     error
	}{
		{"Success", taxonomy{Terms: []term{term{CanonicalName: "Comment", ID: "MTE2-U2VjdGlvbnM=", Children: children{[]term{term{CanonicalName: "Blogs", ID: "MTIx-U2VjdGlvbnM="}}}}}},
			"38dbd827-fedc-3ebe-919f-e64cf55ea959", section{UUID: "38dbd827-fedc-3ebe-919f-e64cf55ea959", CanonicalName: "Comment", TmeIdentifier: "MTE2-U2VjdGlvbnM=", Type: "Section"}, true, nil},
		{"Not found", taxonomy{Terms: []term{term{CanonicalName: "Comment", ID: "MTE2-U2VjdGlvbnM=", Children: children{[]term{term{CanonicalName: "Blogs", ID: "MTIx-U2VjdGlvbnM="}}}}}},
			"some uuid", section{}, false, nil},
		{"Error on init", taxonomy{}, "some uuid", section{}, false, errors.New("Error getting taxonomy")},
	}

	for _, test := range tests {
		repo := dummyRepo{tax: test.tax, err: test.err}
		service, err := newSectionService(&repo, "")
		expectedSection, found := service.getSectionByUUID(test.uuid)
		assert.Equal(test.section, expectedSection, fmt.Sprintf("%s: Expected section incorrect", test.name))
		assert.Equal(test.found, found)
		assert.Equal(test.err, err)
	}
}

type dummyRepo struct {
	tax taxonomy
	err error
}

func (d *dummyRepo) getSectionsTaxonomy() (taxonomy, error) {
	return d.tax, d.err
}
