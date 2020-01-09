package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	gr "github.com/awesome-fc/golang-runtime"
)

// Makes a simple GET request against the serving endpoint to check the serving service is up and running.
// TODO: implement your own testing logic for predictions/inferences
func runTestCases(ctx *gr.FCContext, evt map[string]string) ([]byte, error) {
	servingEndpoint := evt["servingEndpoint"]

	httpClient := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     50 * time.Second,
		},
	}

	resp, err := httpClient.Get(servingEndpoint)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 404 || resp.StatusCode == 200 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		bodyStr := strings.Replace(string(body), "\n", "", -1)
		return []byte(fmt.Sprintf(`{"httpStatus": %d, "servingStatus": "succeeded", "body": "%s"}`, resp.StatusCode, bodyStr)), nil
	}
	return []byte(fmt.Sprintf(`{"httpStatus": %d, "servingStatus": "succeeded"}`, resp.StatusCode)), nil
}