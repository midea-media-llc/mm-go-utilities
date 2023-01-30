package utils

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"futa.express.api.accountant/utils/logs"
)

// MakeGetRequest sends http GET request
func MakeGetRequest(url url.URL, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, url.String(), nil)

	if err != nil {
		logs.Errorf("Cannot create GET request to external resource %s", url.String())

		return nil, err
	}

	// set header
	for key, val := range headers {
		req.Header.Set(key, val)
	}

	// set query param
	q := req.URL.Query()
	for key, val := range queryParams {
		q.Add(key, val)
	}
	req.URL.RawQuery = q.Encode()

	logs.Debugf("Request to external resource GET %s", url.String())

	resp, err := client.Do(req)

	if err != nil {
		logs.Errorf("Error when make http GET request", err)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logs.Errorf("Cannot parse response body", err)

		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logs.Errorf("Request returns non OK status: %d\n%s\n", resp.StatusCode, string(body))
		return nil, errors.New("request returns non OK status")
	}

	return body, nil
}

// MakePostRequest sends http POST request
func MakePostRequest(url url.URL, payload string, headers map[string]string, queryParams map[string]string) ([]byte, error) {
	client := &http.Client{}
	reqBody := bytes.NewBufferString(payload)
	req, err := http.NewRequest(http.MethodPost, url.String(), reqBody)

	if err != nil {
		logs.Errorf("Cannot create POST request to external resource %s \n", url.String())
	}

	// set header
	if headers != nil {
		for key, val := range headers {
			req.Header.Set(key, val)
		}
		req.Header.Set("Content-Type", "application/json")
	}

	// set query param
	if queryParams != nil {
		q := req.URL.Query()
		for key, val := range queryParams {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}

	logs.Infof("Request to external resource POST %s \n", url.String())

	resp, err := client.Do(req)

	if err != nil {
		logs.Errorf("Error when make http POST request", err)

		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		logs.Errorf("Cannot parse response body", err)

		return nil, err
	}

	if resp.StatusCode >= http.StatusBadRequest {
		logs.Errorf("Request returns non OK status: %d\n%s\n", resp.StatusCode, string(body))

		return nil, errors.New("request returns non OK status")
	}

	return body, nil
}
