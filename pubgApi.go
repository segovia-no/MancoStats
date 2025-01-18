package main

import (
	"encoding/json"
	"net/http"
)

var reqHeaders map[string]string

func setupReqHeaders() {
	reqHeaders = make(map[string]string)
	reqHeaders["Authorization"] = "Bearer " + PubgApiToken
	reqHeaders["accept"] = "application/vnd.api+json"
}

func addPubgHeaders(req *http.Request) {
	if reqHeaders == nil {
		setupReqHeaders()
	}
	for k, v := range reqHeaders {
		req.Header.Add(k, v)
	}
}

func PubgApiGET(url string, targetStruct any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	addPubgHeaders(req)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	err = json.NewDecoder(resp.Body).Decode(&targetStruct)
	if err != nil {
		return err
	}
	resp.Body.Close()

	return nil
}
