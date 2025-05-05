package utils

import (
	"io"
	"net/http"
	"strings"
)

func HttpRequestString(url string, others ...string) (string, error) {
	client := http.DefaultClient
	var req *http.Request
	method := "GET"
	var body io.Reader
	if len(others) > 0 && others[0] != "" {
		method = others[0]
	}
	if len(others) > 1 && others[1] != "" {
		body = strings.NewReader(others[1])
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return "", err
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	resultByte, err := io.ReadAll(resp.Body)
	return string(resultByte), err
}
