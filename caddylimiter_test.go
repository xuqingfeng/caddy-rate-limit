package ratelimit

import (
	"strconv"
	"testing"
    "time"

    "golang.org/x/time/rate"
)

func TestAllowAndRetryAfter(t *testing.T) {

	tests := []struct {
		keys      []string
		rule      Rule
		qps       int
        shouldRetryAfter time.Duration
		shouldErr bool
		expected  bool
	}{
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 2, Burst: 2, Unit: "second"}, 2, 0, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 2, Unit: "minute"}, 2, 1 * time.Minute, false, true,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 1, Burst: 0, Unit: "hour"}, 1, rate.InfDuration, false, false,
		},
		{
			[]string{"127.0.0.1", "/"}, Rule{Rate: 0, Burst: 0}, 2, rate.InfDuration, false, true,
		},
	}

	for i, test := range tests {
		test.keys = append(test.keys, strconv.Itoa(i))
		t.Logf("keys: %v", test.keys)
		actual := cl.AllowN(test.keys, test.rule, test.qps)
        retryAfter := cl.RetryAfter(test.keys)
        if retryAfter < test.shouldRetryAfter {
            t.Errorf("Test %d: shouldeRetryAfter %d, got %d", i, test.shouldRetryAfter, retryAfter)
        }
		if actual != test.expected {
			t.Errorf("Test %d: expected %t, got %t", i, test.expected, actual)
		}
	}
}