package ratelimit

import (
	"reflect"
	"testing"

	"github.com/mholt/caddy"
	"github.com/mholt/caddy/caddyhttp/httpserver"
)

func TestSetup(t *testing.T) {

	c := caddy.NewTestController("http", `ratelimit / 2 2`)
	err := setup(c)
	if err != nil {
		t.Errorf("Expected no errors, got: %v", err)
	}
	mids := httpserver.GetConfig(c).Middleware()
	if len(mids) == 0 {
		t.Fatal("Expected middleware, got 0 instead")
	}
}

func TestRateLimitParse(t *testing.T) {

	tests := []struct {
		input     string
		shouldErr bool
		expected  []Rule
	}{
		{
			`ratelimit / 2 0`, false, []Rule{
				{2, 0, []string{"/"}},
			},
		},
		{
			`ratelimit / 2 1`, false, []Rule{
				{2, 1, []string{"/"}},
			},
		},
		{
			`ratelimit / notFloat64 0`, true, []Rule{},
		},
		{
			`ratelimit / 2 0.1`, true, []Rule{},
		},
		{
			`ratelimit 2 2 {
                /resource0
                /resource1
            }`, false, []Rule{
				{2, 2, []string{"/resource0", "/resource1"}},
			},
		},
	}

	for i, test := range tests {
		actual, err := rateLimitParse(caddy.NewTestController("http", test.input))

		if err == nil && test.shouldErr {
			t.Errorf("Test %d didn't error, but it should have", i)
		} else if err != nil && !test.shouldErr {
			t.Errorf("Test %d errored, but it shouldn't have; got '%v'", i, err)
		}

		if len(actual) != len(test.expected) {
			t.Fatalf("Test %d expected %d rules, but got %d", i, len(test.expected), len(actual))
		}

		for j, expectedRule := range test.expected {
			actualRule := actual[j]

			if actualRule.Rate != expectedRule.Rate {
				t.Errorf("Test %d, rule %d: Expected rate '%d', got '%d'", expectedRule.Rate, actualRule.Rate)
			}
			if actualRule.Burst != expectedRule.Burst {
				t.Errorf("Test %d, rule %d: Expected burst '%d', got '%d'", expectedRule.Burst, actualRule.Burst)
			}
			if !reflect.DeepEqual(actualRule.Resources, expectedRule.Resources) {
				t.Errorf("Test %d, rule %d: Expected resource '%v', got '%v'", expectedRule.Resources, actualRule.Resources)
			}
		}
	}
}
