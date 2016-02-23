package main

type section struct {
	UUID          string `json:"uuid"`
	CanonicalName string `json:"canonicalName"`
	TmeIdentifier string `json:"tmeIdentifier"`
	Type          string `json:"type"`
}

type sectionLink struct {
	APIURL string `json:"apiUrl"`
}
