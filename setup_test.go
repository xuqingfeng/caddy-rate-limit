package ratelimit

import (
	"reflect"
	"testing"

	"github.com/caddyserver/caddy"
	"github.com/caddyserver/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {

	c := caddy.NewTestController("http", `ratelimit get / 2 2 second`)
	err := setup(c)
	if err != nil {
		t.Errorf("E! expected no errors, got: %v", err)
	}
	mids := httpserver.GetConfig(c).Middleware()
	if len(mids) == 0 {
		t.Fatal("E! expected middleware, got 0 instead")
	}
}

func TestRateLimitParse(t *testing.T) {

	tests := []struct {
		input     string
		shouldErr bool
		expected  []Rule
	}{
		{
			`ratelimit get / 2 0 second`, false, []Rule{
				{"get", 2, 0, "second", []string{}, "", "*", []string{"/"}},
			},
		},
		{
			`ratelimit post / 2 1 badUnit`, false, []Rule{
				{"post", 2, 1, "badUnit", []string{}, "", "200", []string{"/"}},
			},
		},
		{
			`ratelimit * / notFloat64 0 second`, true, []Rule{},
		},
		{
			`ratelimit * / 2 0.1 second`, true, []Rule{},
		},
		{
			`ratelimit badMethods / 2 1 second`, false, []Rule{
				{"badMethods", 2, 1, "second", []string{}, "", "403", []string{"/"}},
			},
		},
		{
			`ratelimit put,patch 2 2 second {
							whitelist 127.0.0.1/32
							status 403,405
                            /resource0
                            /resource1
                        }`, false, []Rule{
				{"put,patch", 2, 2, "second", []string{"127.0.0.1/32"}, "", "403,405", []string{"/resource0", "/resource1"}},
			},
		},
		{
			`ratelimit 2 3 minute {
					whitelist asdf
					status xyz
					/resource0
					/resource1
				}`, true, []Rule{},
		},
		{
			`ratelimit 2 3 minute {
					whitelist asdf
					limit_by_header xxx
					status xyz
					/resource0
					/resource1
				}`, true, []Rule{},
		},
	}

	for i, test := range tests {
		actual, err := rateLimitParse(caddy.NewTestController("http", test.input))

		if err == nil && test.shouldErr {
			t.Errorf("E! test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("E! test %d error, but it shouldn't have; got %v", i, err)
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("E! test %d expected %d rules, but got %d", i, len(test.expected), len(actual))
		}

		for j, expectedRule := range test.expected {
			actualRule := actual[j]

			if actualRule.Rate != expectedRule.Rate {
				t.Errorf("E! test %d, rule %d: expected rate %d, got %d", i, j, expectedRule.Rate, actualRule.Rate)
			}
			if actualRule.Burst != expectedRule.Burst {
				t.Errorf("E! test %d, rule %d: expected burst %d, got %d", i, j, expectedRule.Burst, actualRule.Burst)
			}
			if actualRule.Unit != expectedRule.Unit {
				t.Errorf("E! test %d, rule %d: expected unit %s, got %s", i, j, expectedRule.Unit, actualRule.Unit)
			}
			if !reflect.DeepEqual(actualRule.Resources, expectedRule.Resources) {
				t.Errorf("E! test %d, rule %d: expected resource %v, got %v", i, j, expectedRule.Resources, actualRule.Resources)
			}
		}
	}
}
