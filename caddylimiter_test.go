package ratelimit

import (
	"strconv"
	"testing"
)

func TestAllow(t *testing.T) {

	tests := []struct {
		keys      []string
		rule      Rule
		qps       int
		shouldErr bool
		expected  bool
	}{
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 2, Burst: 2, Unit: "second"}, 2, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 2, Unit: "minute"}, 2, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 0, Unit: "hour"}, 1, false, false,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 0, Burst: 0}, 2, false, false,
		},
	}

    // TODO: 16/11/25 sleep
	for i, test := range tests {
		test.keys = append(test.keys, strconv.Itoa(i))
		t.Logf("keys: %v", test.keys)
		actual := cl.AllowN(test.keys, test.rule, test.qps)
		if actual != test.expected {
			t.Errorf("Test %d: Expected %t, got %t", i, test.expected, actual)
		}
	}
}
