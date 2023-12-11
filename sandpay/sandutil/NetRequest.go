package sandutil

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func Do(api, jsonData string) ([]byte, error) {

	client := http.DefaultClient
	req, err := http.NewRequest("POST", api, strings.NewReader(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dataByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}

func DoForm(api string, postData url.Values) ([]byte, error) {

	client := http.DefaultClient
	req, err := http.NewRequest("POST", api, strings.NewReader(postData.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	dataByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return dataByte, nil
}
