package main

import (
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type sectionService interface {
	getSections() ([]sectionLink, bool)
	getSectionByUUID(uuid string) (section, bool)
}

type sectionServiceImpl struct {
	repository repository
	baseURL    string
	sectionsMap  map[string]section
	sectionLinks []sectionLink
}

func newSectionService(repo repository, baseURL string) (sectionService, error) {

	s := &sectionServiceImpl{repository: repo, baseURL: baseURL}
	err := s.init()
	if err != nil {
		return &sectionServiceImpl{}, err
	}
	return s, nil
}

func (s *sectionServiceImpl) init() error {
	s.sectionsMap = make(map[string]section)
	tax, err := s.repository.getSectionsTaxonomy()
	if err != nil {
		return err
	}
	s.initSectionsMap(tax.Terms)
	return nil
}

func (s *sectionServiceImpl) getSections() ([]sectionLink, bool) {
	if len(s.sectionLinks) > 0 {
		return s.sectionLinks, true
	}
	return s.sectionLinks, false
}

func (s *sectionServiceImpl) getSectionByUUID(uuid string) (section, bool) {
	section, found := s.sectionsMap[uuid]
	return section, found
}

func (s *sectionServiceImpl) initSectionsMap(terms []term) {
	for _, t := range terms {
		top := transformSection(t)
		s.sectionsMap[top.UUID] = top
		s.sectionLinks = append(s.sectionLinks, sectionLink{APIURL: s.baseURL + top.UUID})
		s.initSectionsMap(t.Children.Terms)
	}
}
