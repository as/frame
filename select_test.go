package frame

import (
	"testing"
)

func TestSelect0to0(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 0)
	check(t, have, "TestSelect0to0", testMode)
}

func TestSelect0to1(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 1)
	check(t, have, "TestSelect0to1", testMode)
}
func TestSelectLine(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 12)
	check(t, have, "TestSelectLine", testMode)
}

func TestSelectLinePlus(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 13)
	check(t, have, "TestSelectLinePlus", testMode)
}

func TestSelectLinePlus1(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 13+1)
	check(t, have, "TestSelectLinePlus1", testMode)
}

func TestSelectAll(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, 9999)
	check(t, have, "TestSelectAll", testMode)
}

func TestSelectAllSub1(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, h.Len())
	p0, p1 := h.Dot()
	p1--
	h.Select(p0, p1)
	check(t, have, "TestSelectAllSub1", testMode)
}

func TestSelectAllSubAll(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(0, h.Len())
	h.Select(0, 0)
	check(t, have, "TestSelectAllSubAll", testMode)
}

func TestMidToEnd(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(h.Len()/2, h.Len())
	check(t, have, "TestMidToEnd", testMode)
}
func TestMidToEndThenStartToMid(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(h.Len()/2, h.Len())
	h.Select(0, h.Len()/2)
	check(t, have, "TestMidToEndThenStartToMid", testMode)
}

func TestSelectTabSpaceNewline(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 5; j++ {
		h.Insert([]byte("abc\t \n\n\t $\n"), int64(j))
	}
	h.Select(h.Len()/2, h.Len()-5)
	check(t, have, "TestSelectTabSpaceNewline", testMode)
}
func TestSelectTabSpaceNewlineSub1(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	for j := 0; j < 5; j++ {
		h.Insert([]byte("abc\t \n\n\t $\n"), int64(j))
	}
	h.Select(h.Len()/2, h.Len()-5-1)
	check(t, have, "TestSelectTabSpaceNewlineSub1", testMode)
}
func TestSelectEndLineAndDec(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	h.Select(167+9, 168+9)
	check(t, have, "TestSelectEndLineAndDec", testMode)
}

var testSelectData = []byte(`Hello world.
Your editor doesn't always know best.
	Your empty file directory has been deleted.
func main(){
	for i := 0; i < 100; i++{
		// comment
	}
}
$ Editor (vi or emacs)?
Usenet is like letters to the editor, only without an editor.  - Larry Wall
Type C-h for help; C-x u to undo changes.  ('C-' means use CTRL key.) GNU Emacs comes with ABSOLUTELY NO WARRANTY; type C-h C-w for full details.You may give out copies of Emacs; type C-h C-c to see the conditions.Type C-h t for a tutorial on using Emacs.





`)

// Broken tests that need work

// TODO(as): regenerate without trailing broken tickmark and rerun this test
func TestSelectNone(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	check(t, have, "TestSelectNone", testMode)
}

// TODO(as): regenerate without trailing broken tickmark and rerun this test
func TestSelectNoneUntick(t *testing.T) {
	h, _, have, _ := abtestbg(R)
	h.Insert(testSelectData, 0)
	check(t, have, "TestSelectNoneUntick", testMode)
}
