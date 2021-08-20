package cmd

import "testing"

// TestNewPlan tests creating a plan
func TestNewPlan(t *testing.T) {
	test := func(pr string, want int) {
		config := NewConfig()
		config.PortRange = pr
		p := newPlan(config)
		got := len(p.items)
		if got != want {
			t.Errorf("got %d, expected %d\n", got, want)
		}
	}

	// test creating plan with single ports
	test("1", 1)
	test("1024", 1)
	test("65535", 1)

	// test creating plan with port ranges
	test("1:1024", 1024)
	test("1024:32000", 30977)
	test("32000:65535", 33536)

	// test creating plan with maximum port range
	test("1:65535", 65535)

	// test creating plan with invalid port ranges
	test("0", 0)
	test("0:0", 0)
	test("1024:3", 0)
	test("65536", 0)
	test("65555", 0)
	test("100000", 0)
	test("65534:65555", 0)
}
