package cmd

import (
	"log"
	"os"
	"testing"
)

// getExamplePrintResultsPlan is a init helper for the printResult examples
func getExamplePrintResultsPlan(pr string) *plan {
	log.SetFlags(0)
	log.SetOutput(os.Stdout)
	config := NewConfig()
	if pr != "" {
		config.PortRange = pr
	}
	return newPlan(config)
}

// Example_printResults_default runs printResults() with default values
func Example_printResults_default() {
	plan := getExamplePrintResultsPlan("")
	plan.printResults()

	// Output:
	// Printing results:
	// 1:65535	drop
}

// Example_printResults_drop runs printResults() with only dropped packets
func Example_printResults_drop() {
	plan := getExamplePrintResultsPlan("1024:1032")
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032	drop
}

// Example_printResults_dropPass runs printResults() with dropped packets and
// some passing packets
func Example_printResults_dropPass() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

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
	// 1024:1026	drop
	// 1027:1029	pass
	// 1030:1032	drop
}

// Example_printResults_dropTCPReset runs printResults() with dropped packets
// and some tcp resetted (rejected) packets
func Example_printResults_dropTCPReset() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

	// create result messages
	r := &MessageResult{
		Result: ResultTCPReset,
	}
	results := []*MessageResult{r}

	// set results for some items
	for i := uint32(3); i < 6; i++ {
		plan.items[i].SenderResults = results
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1026	drop
	// 1027:1029	reject
	// 1030:1032	drop
}

// Example_printResults_rejectTCPReset runs printResults() with
// tcp resetted (rejected) packets
func Example_printResults_rejectTCPReset() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

	// create result messages
	r := &MessageResult{
		Result: ResultTCPReset,
	}
	results := []*MessageResult{r}

	// set results for all items
	for _, i := range plan.items {
		i.SenderResults = results
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032	reject
}

// Example_printResults_pass runs printResults() with passing packets
func Example_printResults_pass() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

	// create result messages
	r := &MessageResult{
		Result: ResultPass,
	}
	results := []*MessageResult{r}

	// set results for all items
	for _, i := range plan.items {
		i.receiverResults = results
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1032	pass
}

// Example_printResults_even runs printResults() with the same amount of
// packets in each category
func Example_printResults_even() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

	// create reset result messages
	rr := &MessageResult{
		Result: ResultTCPReset,
	}
	rresults := []*MessageResult{rr}

	// create pass result messages
	pr := &MessageResult{
		Result: ResultPass,
	}
	presults := []*MessageResult{pr}

	// set reject results for some items
	for i := uint32(3); i < 6; i++ {
		plan.items[i].SenderResults = rresults
	}

	// set pass results for some items
	for i := uint32(6); i < 9; i++ {
		plan.items[i].receiverResults = presults
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024:1026	drop
	// 1027:1029	reject
	// 1030:1032	pass
}

// Example_printResults_interleaved runs printResults() with each packet in a
// different category
func Example_printResults_interleaved() {
	// init
	plan := getExamplePrintResultsPlan("1024:1032")

	// create reset result messages
	rr := &MessageResult{
		Result: ResultTCPReset,
	}
	rresults := []*MessageResult{rr}

	// create pass result messages
	pr := &MessageResult{
		Result: ResultPass,
	}
	presults := []*MessageResult{pr}

	for i, item := range plan.items {
		switch i % 3 {
		case 1:
			// set reject results
			item.SenderResults = rresults
		case 2:
			// set pass results
			item.receiverResults = presults
		}
	}

	// check output
	plan.printResults()

	// Output:
	// Printing results:
	// 1024	drop
	// 1025	reject
	// 1026	pass
	// 1027	drop
	// 1028	reject
	// 1029	pass
	// 1030	drop
	// 1031	reject
	// 1032	pass
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
