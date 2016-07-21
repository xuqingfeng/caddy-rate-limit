package ratelimit

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetRemoteIP(t *testing.T) {

	tests := []struct {
		input    map[string]string
		expected string
	}{
		{
			make(map[string]string),
			"127.0.0.1",
		},
	}

	defaultClient := &http.Client{}
	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	for i, test := range tests {
		req, err = http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatalf("Test %d errored: '%v'", i, err)
		}
		for k, v := range test.input {
			req.Header.Add(k, v)
		}
		resp, err = defaultClient.Do(req)
		if err != nil {
			t.Errorf("Test %d errored: '%v'", i, err)
		}
		// FIXME: 16/7/21 bytes.NewBuffer
		ip, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Errorf("Test %d errored: '%v'", i, err)
		}
		if string(ip) != test.expected {
			t.Errorf("Test %d Expected %s, get %s", i, test.expected, ip)
		}
		resp.Body.Close()
	}
}
