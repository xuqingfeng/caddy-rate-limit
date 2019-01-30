package ratelimit

import (
	"bytes"
	"net"
	"net/http"
	"testing"
)

func TestIsWhitelistIPAddress(t *testing.T) {

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
	testWhitelist := []string{"127.0.0.0/8", "10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16", "::1/128", "fc00::/7"}
	var testWhitelistIPNets []*net.IPNet
	for _, s := range testWhitelist {
		_, ipNet, err := net.ParseCIDR(s)
		if err == nil {
			testWhitelistIPNets = append(testWhitelistIPNets, ipNet)
		}
	}

	for i, test := range tests {
		if ret := IsWhitelistIPAddress(test.input, testWhitelistIPNets); ret != test.expected {
			t.Errorf("E! test %d expected %t, got %t", i, test.expected, ret)
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

	var (
		req  *http.Request
		resp *http.Response
		err  error
	)
	for i, test := range tests {
		req, err = http.NewRequest("GET", server.URL, nil)
		if err != nil {
			t.Fatalf("F! test %d error: %v", i, err)
		}
		for k, v := range test.input {
			req.Header.Add(k, v)
		}
		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			t.Errorf("E! test %d error: %v", i, err)
		}
		// bytes.NewBuffer https://stackoverflow.com/questions/37314715/reading-http-response-body-stream
		buf := bytes.NewBuffer(make([]byte, 0, resp.ContentLength))
		_, err := buf.ReadFrom(resp.Body)
		if err != nil {
			t.Errorf("E! test %d error: %v", i, err)
		}
		ip := buf.Bytes()
		if string(ip) != test.expected {
			t.Errorf("E! test %d expected %s, got %s", i, test.expected, ip)
		}
		resp.Body.Close()
	}
}

func TestMatchMethod(t *testing.T) {

	tests := []struct {
		input       string
		inputMethod string
		expected    bool
	}{
		{"*", "GET", true},
		{
			"get", "GET", true,
		},
		{
			"get,post", "POST", true,
		},
		{
			"put,patch", "DELETE", false,
		},
	}

	var (
		err error
	)

	for i, test := range tests {
		if err != nil {
			t.Fatalf("F! test %d error: %v", i, err)
		}
		if ret := MatchMethod(test.input, test.inputMethod); ret != test.expected {
			t.Errorf("E! test %d expected %t, got %t", i, test.expected, ret)
		}
	}
}

func TestMatchStatus(t *testing.T) {

	tests := []struct {
		input       string
		inputStatus string
		expected    bool
	}{
		{"", "200", false},
		{"*", "200", false},
		{"200,404", "404", true},
	}

	var (
		err error
	)

	for i, test := range tests {
		if err != nil {
			t.Fatalf("F! test %d error: %v", i, err)
		}

		if ret := MatchStatus(test.input, test.inputStatus); ret != test.expected {
			t.Errorf("E! test %d expected %t, got %t", i, test.expected, ret)
		}
	}
}
