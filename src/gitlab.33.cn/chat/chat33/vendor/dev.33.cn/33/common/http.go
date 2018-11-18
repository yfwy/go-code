package common

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) ([]byte, error) {
	resp, err := http.Get(url)

	if nil != resp {
		defer resp.Body.Close()
	}

	if nil != err {
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func HTTPPostForm(reqUrl string, headers map[string]string, payload io.Reader) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}

	headers["Content-Type"] = "application/x-www-form-urlencoded"
	return HTTPRequest("POST", reqUrl, headers, payload)
}

// HTTPPostJSON with json []byte
func HTTPPostJSON(reqUrl string, headers map[string]string, payload io.Reader) ([]byte, error) {
	if headers == nil {
		headers = make(map[string]string)
	}
	headers["Content-Type"] = "text/json"
	return HTTPRequest("POST", reqUrl, headers, payload)
}

func HTTPRequest(method string, reqUrl string, headers map[string]string, payload io.Reader) ([]byte, error) {
	req, err := http.NewRequest(method, reqUrl, payload)
	if err != nil {
		return nil, errors.New("make http request error: " + err.Error())
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Add(k, v)
		}
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}

	if err != nil {
		return nil, errors.New("do http request error: " + err.Error())
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("ready http response error: " + err.Error())
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode:%d, Desc:%s", resp.StatusCode, string(body)))
	}

	return body, nil
}
