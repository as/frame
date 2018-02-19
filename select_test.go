package frame

import (
	"github.com/as/etch"
	"testing"
)

func TestSelectFlow(t *testing.T) {
	h, w, have, want := abtestbg(R)
	x := []byte("The quick brown fox jumped over the lazy dog")
	lx := int64(len(x))
	h.Insert(x, 0)
	w.Insert(x, 0)
	for i := int64(0); i < lx; i++ {
		h.Select(0+i, lx)
		h.Select(0+i, 0+i)
		h.Insert([]byte("@"), i*2)
		h.Delete(i*2, i*2+1)
	}
	h.Select(h.Len(), h.Len())
	etch.Assert(t, have, want, "TestSelectFlow.png")
}

func TestTypingFlow(t *testing.T) {
	h, _, _, _ := abtestbg(R)
	for _, c := range []byte("abcde") {
		p0, _ := h.Dot()
		h.Insert([]byte{c}, p0)
	}
	have, _ := h.Dot()
	want := int64(5)
	if have != want {
		t.Logf("typing dot: have %d, want %d\n", have, want)
		t.Fail()
	}
}

func testReg(t *testing.T, name, what string, where int) {
	t.Helper()
	h, _, have, _ := abtestbg(R)
	h.Insert([]byte("abcde"), 0)
	h.Select(1, 4)
	ckDot(t, h, 1, 4)
	h.Insert([]byte(what), int64(where))
	p0, p1 := int64(1), int64(4)
	if where <= 1 {
		p0++
		p1++
	} else if where < 4 {
		p1++
	}
	ckDot(t, h, p0, p1)
	h.Select(p0, p1)
	check(t, have, "TestInsertRegion"+name, testMode)
}

func TestRegion0(t *testing.T) { testReg(t, "0", "0", 0) }
func TestRegion1(t *testing.T) { testReg(t, "1", "1", 1) }
func TestRegion2(t *testing.T) { testReg(t, "2", "2", 2) }
func TestRegion3(t *testing.T) { testReg(t, "3", "3", 3) }
func TestRegion4(t *testing.T) { testReg(t, "4", "4", 4) }
func TestRegion5(t *testing.T) { testReg(t, "5", "5", 5) }

func ckDot(t *testing.T, f *Frame, p0, p1 int64) {
	t.Helper()
	q0, q1 := f.Dot()
	if q0 != int64(p0) || q1 != int64(p1) {
		t.Logf("bad selection: have [%d:%d), want [%d:%d)\n", q0, q1, p0, p1)
		t.Fail()
	}
}

func TestInsertExtendsSelection(t *testing.T) {
	h, _, _, _ := abtestbg(R)
	h.Insert([]byte("abcde"), 0)
	h.Select(1, 4)
	h.Insert([]byte("x"), 2)
	ckDot(t, h, 1, 5)
}
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
