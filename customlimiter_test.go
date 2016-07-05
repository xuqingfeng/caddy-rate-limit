package ratelimit

import (
	"strconv"
	"testing"
)

func TestAllow(t *testing.T) {

	tests := []struct {
		keys      []string
		rule      Rule
		shouldErr bool
		expected  bool
	}{
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 2, Burst: 2}, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 2}, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 0}, false, false,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 0, Burst: 0}, false, false,
		},
	}

	for i, test := range tests {
		test.keys = append(test.keys, strconv.Itoa(i))
		t.Logf("keys: %v", test.keys)
		cl.Allow(test.keys, test.rule)
		// second time
        actual := cl.Allow(test.keys, test.rule)
		if actual != test.expected {
			t.Errorf("Test %d: Expected %t, got %t", i, test.expected, actual)
		}
	}
}
