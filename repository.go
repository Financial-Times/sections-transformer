package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

type repository interface {
	getSectionsTaxonomy() (taxonomy, error)
}

type tmeRepository struct {
	httpClient              httpClient
	principalHeader         string
	structureServiceBaseURL string
}

func newTmeRepository(client httpClient, structureServiceBaseURL string, principalHeader string) repository {
	return &tmeRepository{httpClient: client, principalHeader: principalHeader, structureServiceBaseURL: structureServiceBaseURL}
}

func (t *tmeRepository) getSectionsTaxonomy() (taxonomy, error) {
	req, err := http.NewRequest("GET", t.structureServiceBaseURL+"/metadata-services/structure/1.0/taxonomies/sections/terms?includeDisabledTerms=false", nil)
	if err != nil {
		return taxonomy{}, err
	}
	req.Header.Set("ClientUserPrincipal", t.principalHeader)
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return taxonomy{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return taxonomy{}, fmt.Errorf("Structure service returned %d", resp.StatusCode)
	}

	contents, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return taxonomy{}, err
	}

	tax := taxonomy{}
	err = xml.Unmarshal(contents, &tax)
	if err != nil {
		return taxonomy{}, err
	}
	return tax, nil
}
