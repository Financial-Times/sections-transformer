package main

import (
	"github.com/Financial-Times/tme-reader/tmereader"
	log "github.com/Sirupsen/logrus"
	"net/http"
)

type httpClient interface {
	Do(req *http.Request) (resp *http.Response, err error)
}

type sectionService interface {
	getSections() ([]sectionLink, bool)
	getSectionByUUID(uuid string) (section, bool)
	checkConnectivity() error
	getSectionCount() int
	getSectionIds() []string
	reload() error
}

type sectionServiceImpl struct {
	repository    tmereader.Repository
	baseURL       string
	sectionsMap   map[string]section
	sectionLinks  []sectionLink
	taxonomyName  string
	maxTmeRecords int
}

func newSectionService(repo tmereader.Repository, baseURL string, taxonomyName string, maxTmeRecords int) (sectionService, error) {
	s := &sectionServiceImpl{repository: repo, baseURL: baseURL, taxonomyName: taxonomyName, maxTmeRecords: maxTmeRecords}
	err := s.init()
	if err != nil {
		return &sectionServiceImpl{}, err
	}
	return s, nil
}

func (s *sectionServiceImpl) init() error {
	s.sectionsMap = make(map[string]section)
	responseCount := 0
	log.Printf("Fetching sections from TME\n")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Printf("Finished fetching sections from TME\n")
			break
		}
		s.initSectionsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d section links\n", len(s.sectionLinks))

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

func (s *sectionServiceImpl) checkConnectivity() error {
	// TODO: Can we just hit an endpoint to check if TME is available? Or do we need to make sure we get genre taxonmies back? Maybe a healthcheck or gtg endpoint?
	// TODO: Can we use a count from our responses while actually in use to trigger a healthcheck?
	//	_, err := s.repository.GetTmeTermsFromIndex(1)
	//	if err != nil {
	//		return err
	//	}
	return nil
}

func (s *sectionServiceImpl) initSectionsMap(terms []interface{}) {
	for _, iTerm := range terms {
		t := iTerm.(term)
		top := transformSection(t, s.taxonomyName)
		s.sectionsMap[top.UUID] = top
		s.sectionLinks = append(s.sectionLinks, sectionLink{APIURL: s.baseURL + top.UUID})
	}
}

func (s *sectionServiceImpl) getSectionCount() int {
	return len(s.sectionLinks)
}

func (s *sectionServiceImpl) getSectionIds() []string {
	i := 0
	keys := make([]string, len(s.sectionsMap))

	for k := range s.sectionsMap {
		keys[i] = k
		i++
	}
	return keys
}

func (s *sectionServiceImpl) reload() error {
	s.sectionsMap = make(map[string]section)
	responseCount := 0
	log.Println("Fetching sections from TME")
	for {
		terms, err := s.repository.GetTmeTermsFromIndex(responseCount)
		if err != nil {
			return err
		}

		if len(terms) < 1 {
			log.Println("Finished fetching topics from TME")
			break
		}
		s.initSectionsMap(terms)
		responseCount += s.maxTmeRecords
	}
	log.Printf("Added %d section links\n", len(s.sectionLinks))

	return nil
}
