package ratelimit

import (
	"fmt"
	"strconv"
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func TestAllowNAndRetryAfter(t *testing.T) {

	tests := []struct {
		keys             []string
		rule             Rule
		qps              int
		shouldRetryAfter time.Duration
		shouldErr        bool
		expected         bool
	}{
		{
			[]string{"127.0.0.1", "get", "", "/"}, Rule{Methods: "", Status: "", Rate: 2, Burst: 2, Unit: "second"}, 2, 0, false, true,
		},
		{
			[]string{"127.0.0.1", "get,post", "*", "/"}, Rule{Methods: "get,post", Status: "*", Rate: 1, Burst: 2, Unit: "minute"}, 2, 30 * time.Second, false, true,
		},
		{
			[]string{"127.0.0.1", "*", "", "/"}, Rule{Methods: "*", Status: "", Rate: 1, Burst: 0, Unit: "hour"}, 1, rate.InfDuration, false, false,
		},
		{
			[]string{"127.0.0.1", "", "*", "/"}, Rule{Methods: "", Status: "*", Rate: 0, Burst: 0}, 2, 0, false, true,
		},
		{
			[]string{"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6ImNhZGR5LXJhdGUtbGltaXQiLCJpYXQiOjE1MTYyMzkwMjJ9.-i8APTDIVPfX4XjLQIEJH-GugwBoDizdUJjDgnOJpCI", "", "*", "/"}, Rule{Methods: "", Status: "*", Rate: 0, Burst: 0}, 2, 0, false, true,
		},
	}

	for i, test := range tests {
		test.keys = append(test.keys, strconv.Itoa(i))
		actual := cl.AllowN(test.keys, test.rule, test.qps)
		retryAfter := cl.RetryAfter(test.keys)
		if retryAfter < test.shouldRetryAfter {
			t.Errorf("E! test %d: shouldRetryAfter %d, got %d", i, test.shouldRetryAfter, retryAfter)
		}
		if actual != test.expected {
			t.Errorf("E! test %d: expected %t, got %t", i, test.expected, actual)
		}
	}

	// spawn multiple goroutines to test concurrent read/write in map
	num := make([]int, 1000)
	for range num {
		go func() {
			for {
				cl.AllowN(tests[0].keys, tests[0].rule, tests[0].qps)
			}
		}()
		go func() {
			for {
				cl.AllowN(tests[0].keys, tests[0].rule, tests[0].qps)
			}
		}()
	}
}

func BenchmarkSingleKey(b *testing.B) {

	keys := []string{"127.0.0.1", "get", "", "/"}
	for i := 1; i <= 8; i *= 2 {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				benchmarkAllowNAndRetryAfter(keys)
			}
		})
	}
}

func BenchmarkRandomKey(b *testing.B) {

	for i := 1; i <= 8; i *= 2 {
		b.Run(fmt.Sprintf("%d", i), func(b *testing.B) {
			for n := 0; n < b.N; n++ {
				keys := []string{"127.0.0.1", "get", "", "/" + strconv.Itoa(i) + "-" + strconv.Itoa(n)}
				benchmarkAllowNAndRetryAfter(keys)
			}
		})
	}
}

func benchmarkAllowNAndRetryAfter(keys []string) {

	cl.AllowN(keys, Rule{Methods: "get", Status: "", Rate: 2, Burst: 2, Unit: "second"}, 1)
	cl.RetryAfter(keys)
}
