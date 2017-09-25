package frame

import (
	"testing"
)

func TestInsertOneChar(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert([]byte("1"), 0)
	h.Untick()
	//etch.WriteFile(t, `testdata/TestInsertOneChar.expected.png`, have)
	check(t, have, "TestInsertOneChar", modeCheckResult)
}

func TestInsert10Chars(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for i := 0; i < 10; i++ {
		h.Insert([]byte("1"), 0)
	}
	check(t, have, "TestInsert10Chars", modeCheckResult)
}

func TestInsert22Chars2Lines(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 2; j++ {
		for i := 0; i < 10; i++ {
			h.Insert([]byte("1"), h.Len())
		}
		h.Insert([]byte("\n"), h.Len())
	}
	check(t, have, "TestInsert22Chars2Lines", modeCheckResult)
}

func TestInsert1000(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 1000; j++ {
		h.Insert([]byte{byte(j)}, int64(j))
	}
	check(t, have, "TestInsert1000", modeCheckResult)
}

func TestInsertTabSpaceNewline(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 5; j++ {
		h.Insert([]byte("abc\t \n\n\t $\n"), int64(j))
	}
	check(t, have, "TestInsertTabSpaceNewline", modeSaveResult)
}
