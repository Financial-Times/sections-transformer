package main

type section struct {
	UUID                   string                 `json:"uuid"`
	AlternativeIdentifiers alternativeIdentifiers `json:"alternativeIdentifiers,omitempty"`
	PrefLabel              string                 `json:"prefLabel"`
	Type                   string                 `json:"type"`
}

type alternativeIdentifiers struct {
	TME   []string `json:"TME,omitempty"`
	Uuids []string `json:"uuids,omitempty"`
}

type sectionLink struct {
	APIURL string `json:"apiUrl"`
}
