package contentLoader

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

type VideoMeta struct {
	Id          string
	Title       string
	Type        string
	DownloadUrl string
}

type VideoMetaExtractor interface {
	ExtractVideoMeta(url string) (*VideoMeta, error)
	IsServingUrl(url string) bool
}

type ContentLoader struct {
	extractors []VideoMetaExtractor
}

// register sources here and remove factory
func New(extractors ...VideoMetaExtractor) *ContentLoader {
	return &ContentLoader{
		extractors: extractors,
	}
}

// AK TODO make it injectable, requires refactoring of sources
func (d *ContentLoader) DownloadContent(dUrl string) ([]byte, error) {
	log.Info("Start downloading ", time.Now())
	defer log.Info("Finish downloading ", time.Now())
	// AK TODO temp solutions
	//tr := &http.Transport{
	//TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // <--- Problem
	//}

	// AK TODO not sure if the client is reausable
	client := &http.Client{
		Timeout: 50 * time.Second,
	}

	resp, err := client.Get(dUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("the request has failde with statuscode %d. Data: %s", resp.StatusCode, body)
	}

	return body, nil
}

func (l *ContentLoader) ExtractVideoMeta(url string) (*VideoMeta, error) {
	extractor, err := l.GetExtractor(url)

	if err != nil {
		return nil, err
	}

	return extractor.ExtractVideoMeta(url)
}

func (l *ContentLoader) CanExtractVideoMeta(url string) bool {
	for _, extractor := range l.extractors {
		if extractor.IsServingUrl(url) {
			return true
		}
	}
	return false
}

func (l *ContentLoader) GetExtractor(url string) (VideoMetaExtractor, error) {
	for _, extractor := range l.extractors {
		if extractor.IsServingUrl(url) {
			return extractor, nil
		}
	}
	return nil, fmt.Errorf("No extractor found for url %s. Video can't be converted. Please register an extractor.", url)
}
