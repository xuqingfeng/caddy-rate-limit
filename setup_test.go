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
			`ratelimit / 2 0 second`, false, []Rule{
				{2, 0, []string{"/"}, "second"},
			},
		},
		{
			`ratelimit / 2 1 badUnit`, false, []Rule{
				{2, 1, []string{"/"}, "badUnit"},
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
                /resource0
                /resource1
            }`, false, []Rule{
				{2, 2, []string{"/resource0", "/resource1"}, "second"},
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
				t.Errorf("Test %d, rule %d: Expected rate '%f', got '%f'", i, j, expectedRule.Rate, actualRule.Rate)
			}
			if actualRule.Burst != expectedRule.Burst {
				t.Errorf("Test %d, rule %d: Expected burst '%d', got '%d'", i, j, expectedRule.Burst, actualRule.Burst)
			}
			if actualRule.Unit != expectedRule.Unit {
				t.Errorf("Test %d, rule %d: Expected unit '%s', got '%s'", i, j, expectedRule.Unit, actualRule.Unit)
			}
			if !reflect.DeepEqual(actualRule.Resources, expectedRule.Resources) {
				t.Errorf("Test %d, rule %d: Expected resource '%v', got '%v'", i, j, expectedRule.Resources, actualRule.Resources)
			}
		}
	}
}
