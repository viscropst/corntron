package internal

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cheggaaa/pb/v3"
)

func HttpRequest(url string, others ...string) (io.ReadCloser, int64, error) {
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
		return nil, 0, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		resultByte, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, 0, errors.New(string(resultByte))
	}
	return resp.Body, resp.ContentLength, nil
}

func HttpRequestBytes(url string, others ...string) ([]byte, error) {
	resp, _, err := HttpRequest(url, others...)
	defer CloseFileAndFinishBar(resp, nil)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(resp)
}

func HttpRequestString(url string, others ...string) (string, error) {
	resultByte, err := HttpRequestBytes(url, others...)
	if err != nil {
		return "", err
	}
	return string(resultByte), err
}

func HttpRequestFile(url, filename string, others ...string) error {
	resp, len, err := HttpRequest(url, others...)
	if err != nil {
		return err
	}
	defer CloseFileAndFinishBar(resp, nil)
	bar := pb.Default.Start64(len)
	return IOToFile(resp, filename, bar)
}
