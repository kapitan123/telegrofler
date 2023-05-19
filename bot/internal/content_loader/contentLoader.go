package contentLoader

import (
	"fmt"
	"io"
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
	client     *http.Client
}

func New(extractors ...VideoMetaExtractor) *ContentLoader {
	client := &http.Client{
		Timeout: 50 * time.Second,
	}

	return &ContentLoader{
		extractors: extractors,
		client:     client,
	}
}

// AK TODO make it injectable, requires refactoring of sources
// AK TODO should have context for termination
func (d *ContentLoader) DownloadContent(dUrl string, w io.Writer) error {
	log.Info("Start downloading ", time.Now())
	defer log.Info("Finish downloading ", time.Now())

	resp, err := d.client.Get(dUrl)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("the request has failde with statuscode %d", resp.StatusCode)
	}

	_, err = io.Copy(w, resp.Body)

	if err != nil {
		return err
	}

	return nil
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
	return nil, fmt.Errorf("no extractor found for url %s. Video can't be converted. Please register an extractor", url)
}