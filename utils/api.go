package utils

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"

	"futa.express.api.accountant/utils/logs"
)

func GetWithoutJWT(url string, token string) (result string, status int) {
	var jsonStr = []byte("")

	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode
	}
	return string(body), resp.StatusCode
}

func Get(url string, token string) (result string, status int) {
	var jsonStr = []byte("")

	req, _ := http.NewRequest("GET", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode
	}
	return string(body), resp.StatusCode
}

// Post PostApi
func Post(url string, jsonData string, token string) (result string, status int) {
	logs.Infof("Calling URL %v with POST\n", url)
	var jsonStr = []byte(jsonData)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}

	resp, err := client.Do(req)

	if err != nil {
		logs.Errorf(fmt.Sprintf("Call api %v error %v", url, err))
		return "", http.StatusInternalServerError
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", resp.StatusCode
	}
	return string(body), resp.StatusCode
}
