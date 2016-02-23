package main

type taxonomy struct {
	Terms []term `xml:"term"`
}

type term struct {
	CanonicalName string   `xml:"canonicalName"`
	ID            string   `xml:"id,attr"`
	Children      children `xml:"children"`
}

type children struct {
	Terms []term `xml:"term"`
}
