package common

import (
	"errors"
	"regexp"
)

type NewReleaseEvent struct {
	Event       string `json:"event"`
	BuildId     int    `json:"buildId"`
	Branch      string `json:"branch"`
	DownloadUrl string `json:"downloadUrl"`
	ProductName string `json:"productName"`
	Timestamp   string `json:"timestamp"`
	BuildSHA    string `json:"buildSHA"`
}

func (e *NewReleaseEvent) Validate() error {
	urlValidator := regexp.MustCompile(`(?:http(s)?://)?[\w.-]+(?:\.[\w.-]+)+[\w\-_~:/?#[\]@!$&'()*+,;=.]+$`)

	if e.ProductName == "" {
		return errors.New("productName must be specified")
	}
	if e.DownloadUrl == "" {
		return errors.New("downloadUrl must be specified")
	}
	if !urlValidator.MatchString(e.DownloadUrl) {
		return errors.New("downloadUrl does not look like a valid URL")
	}
	if e.Branch == "" {
		return errors.New("branch must be specified")
	}
	return nil
}

type SearchRequest struct {
	Branch           string `json:"branch"`
	ProductName      string `json:"productName"`
	AlwaysShowMaster bool   `json:"alwaysShowMaster"`
}
