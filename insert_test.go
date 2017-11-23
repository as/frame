package frame

import (
	"testing"
)

func TestInsertOneChar(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(string("1"), 0)
	h.Untick()
	//etch.WriteFile(t, `testdata/TestInsertOneChar.expected.png`, have)
	check(t, have, "TestInsertOneChar", testMode)
}

func TestInsert10Chars(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for i := 0; i < 10; i++ {
		h.Insert(string("1"), 0)
	}
	check(t, have, "TestInsert10Chars", testMode)
}

func TestInsert22Chars2Lines(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 2; j++ {
		for i := 0; i < 10; i++ {
			h.Insert(string("1"), h.Len())
		}
		h.Insert(string("\n"), h.Len())
	}
	check(t, have, "TestInsert22Chars2Lines", testMode)
}

func TestInsert1000(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 1000; j++ {
		h.Insert(string{byte(j)}, int64(j))
	}
	check(t, have, "TestInsert1000", testMode)
}

func TestInsertTabSpaceNewline(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 5; j++ {
		h.Insert(string("abc\t \n\n\t $\n"), int64(j))
	}
	check(t, have, "TestInsertTabSpaceNewline", testMode)
}
