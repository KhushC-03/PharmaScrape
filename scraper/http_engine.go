package scraper

import (
	"net/http"
	"net/http/cookiejar"
)

func CreateSession() (*http.Client, error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
}
