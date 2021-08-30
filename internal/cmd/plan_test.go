package cmd

import (
	"log"
	"os"
	"testing"
)

// Example_printResults runs printResults() with default values
func Example_printResults_default() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	plan := newPlan(NewConfig())
	plan.printResults()

	// Output:
	// Printing results:
	// 1:65535 policy DROP
}

// Example_printResults runs printResults() with only dropped packets
func Example_printResults_drop() {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	config := NewConfig()
	config.PortRange = "1024:1032"
	plan := newPlan(config)
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032 policy DROP
}

// Example_printResults runs printResults() with dropped packets and
// some passing packets
func Example_printResults_dropPass() {
	// init
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	config := NewConfig()
	config.PortRange = "1024:1032"
	plan := newPlan(config)

	// create result messages
	r := &MessageResult{
		Result: ResultPass,
	}
	results := []*MessageResult{r}

	// set results for some items
	for i := uint32(3); i < 6; i++ {
		plan.items[i].receiverResults = results
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032 policy DROP
	// 1027	pass
	// 1028	pass
	// 1029	pass
}

// Example_printResults runs printResults() with dropped packets and
// some tcp resetted (rejected) packets
func Example_printResults_dropTCPReset() {
	// init
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	config := NewConfig()
	config.PortRange = "1024:1032"
	plan := newPlan(config)

	// create result messages
	r := &MessageResult{
		Result: ResultTCPReset,
	}
	results := []*MessageResult{r}

	// set results for some items
	for i := uint32(3); i < 6; i++ {
		plan.items[i].senderResults = results
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032 policy DROP
	// 1027	reject
	// 1028	reject
	// 1029	reject
}

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
