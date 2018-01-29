package frame

import (
	"fmt"
	"image"
	"os"
	"testing"

	"github.com/as/etch"
)

// If adding new graphical tests, change to modeSaveResult
const testMode = modeSaveResult

const (
	modeSaveResult = iota
	modeCheckResult
)

func TestMain(m *testing.M) {
	v := m.Run()
	if testMode == modeCheckResult {
		v = 1
		fmt.Println("*** DANGER ***")
		fmt.Println("*** testMode == modeSaveResult ")
		fmt.Println("*** change to testMode = modeCheckResult in etch_test.go")
		fmt.Println()
	}
	os.Exit(v)
}

func check(t *testing.T, have image.Image, name string, mode int) {
	wantfile := fmt.Sprintf("testdata/%s.expected.png", name)
	if mode == modeSaveResult {
		etch.WriteFile(t, wantfile, have)
	}
	etch.AssertFile(t, have, wantfile, fmt.Sprintf("%s.png", name))
}
