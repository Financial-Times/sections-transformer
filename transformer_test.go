package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTransform(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		name    string
		term    term
		section section
	}{
		{"Trasform term to section", term{CanonicalName: "Blogs", ID: "MTIx-U2VjdGlvbnM="}, section{UUID: "d22ff01b-f7b0-3c84-8543-c92e346b4585", CanonicalName: "Blogs", TmeIdentifier: "MTIx-U2VjdGlvbnM=", Type: "Section"}},
	}

	for _, test := range tests {
		expectedSection := transformSection(test.term)
		assert.Equal(test.section, expectedSection, fmt.Sprintf("%s: Expected section incorrect", test.name))
	}

}
