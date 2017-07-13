package ratelimit

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {

	c := caddy.NewTestController("http", `ratelimit / 2 2 second`)
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
			`ratelimit / 2 0 second`, false, []Rule{
				{2, 0, []string{}, []string{"/"}, "second"},
			},
		},
		{
			`ratelimit / 2 1 badUnit`, false, []Rule{
				{2, 1, []string{}, []string{"/"}, "badUnit"},
			},
		},
		{
			`ratelimit / notFloat64 0 second`, true, []Rule{},
		},
		{
			`ratelimit / 2 0.1 second`, true, []Rule{},
		},
		{
			`ratelimit 2 2 second {
							whitelist 127.0.0.1/32
                            /resource0
                            /resource1
                        }`, false, []Rule{
				{2, 2, []string{"127.0.0.1/32"}, []string{"/resource0", "/resource1"}, "second"},
			},
		},
		{
			`ratelimit 2 3 minute {
					whitelist asdf
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
			t.Errorf("E! test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("E! test %d expected %d rules, but got %d", i, len(test.expected), len(actual))
		}

		for j, expectedRule := range test.expected {
			actualRule := actual[j]

			if actualRule.Rate != expectedRule.Rate {
				t.Errorf("E! test %d, rule %d: expected rate '%d', got '%d'", i, j, expectedRule.Rate, actualRule.Rate)
			}
			if actualRule.Burst != expectedRule.Burst {
				t.Errorf("E! test %d, rule %d: expected burst '%d', got '%d'", i, j, expectedRule.Burst, actualRule.Burst)
			}
			if actualRule.Unit != expectedRule.Unit {
				t.Errorf("E! test %d, rule %d: expected unit '%s', got '%s'", i, j, expectedRule.Unit, actualRule.Unit)
			}
			if !reflect.DeepEqual(actualRule.Resources, expectedRule.Resources) {
				t.Errorf("E! test %d, rule %d: expected resource '%v', got '%v'", i, j, expectedRule.Resources, actualRule.Resources)
			}
		}
	}
}
