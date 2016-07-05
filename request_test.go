package ratelimit

import (
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetRemoteIP(t *testing.T) {

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Error(err)
	}
	ip, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		t.Error(err)
	}
	t.Log(string(ip))
	if string(ip) != "127.0.0.1" {
		t.Errorf("Expected 127.0.0.1, get %s", ip)
	}
}
