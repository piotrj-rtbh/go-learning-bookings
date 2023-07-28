package main

import "testing"

// running tests: go to proper folder: cd ./cmd/web
// and run: go test
// coverage: go test -cover
// coverage in browser (shows which parts were covered and not covered by the tests):
// 		go test -coverprofile=coverage.out && go tool cover -html=coverage.out

func TestRun(t *testing.T) {
	err := run()
	if err != nil {
		t.Error("failed run()")
	}
}
