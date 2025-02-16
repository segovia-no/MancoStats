package pubgDAL

import (
	"encoding/json"
	"net/http"
)

type PUBGRequester struct {
	reqHeaders map[string]string
}

func NewPUBGRequester(apiToken string) *PUBGRequester {
	return &PUBGRequester{
		reqHeaders: genRequestHeaders(apiToken),
	}
}

func genRequestHeaders(apiToken string) map[string]string {
	reqHeaders := make(map[string]string)
	reqHeaders["Authorization"] = "Bearer " + apiToken
	reqHeaders["accept"] = "application/vnd.api+json"

	return reqHeaders
}

func (r *PUBGRequester) Get(url string, targetStruct any) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	for k, v := range r.reqHeaders {
		req.Header.Add(k, v)
	}

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
