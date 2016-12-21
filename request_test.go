package ratelimit

import (
	"io/ioutil"
	"net/http"
	"testing"
    "bytes"
)

func TestIsLocalIpAddress(t *testing.T) {

	tests := []struct {
		input    string
		expected bool
	}{
		{
			"127.0.0.10",
			true,
		},
		{
			"10.1.2.3",
			true,
		},
		{
			"172.16.0.10",
			true,
		},
		{
			"192.168.100.10",
			true,
		},
		{
			"100.100.100.100",
			false,
		},
		{
			"::1",
			true,
		},
		{
			"fdff:ffff:ffff:ffff:ffff:ffff:ffff:ffff",
			true,
		},
		{
			"fe80::",
			false,
		},
	}

	for i, test := range tests {
		if ret := IsLocalIpAddress(test.input, localIpNets); ret != test.expected {
			t.Errorf("Test %d Expected %t, get %t", i, test.expected, ret)
		}
	}
}

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
		// bytes.NewBuffer https://stackoverflow.com/questions/37314715/reading-http-response-body-stream
        buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
        _, err := buf.ReadFrom(resp.Body)
        if err != nil {
            t.Errorf("Test %d errored: '%v'", i, err)
        }
        ip := buf.Bytes()
		if string(ip) != test.expected {
			t.Errorf("Test %d Expected %s, get %s", i, test.expected, ip)
		}
		resp.Body.Close()
	}
}
