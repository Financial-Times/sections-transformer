package main

type section struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	PrimaryType            string                 `json:"type"`
	TypeHierarchy          []string               `json:"types"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	Uuids []string `json:"uuids,omitempty"`
}

type sectionLink struct {
	APIURL string `json:"apiUrl"`
}

var primaryType = "Section"
var sectionTypes = []string{"Thing", "Concept", "Classification", "Section"}
