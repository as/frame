package box

import (
	"github.com/as/frame/font"
	"testing"
)

func runwith(s string) *Run {
	r := NewRun(5, 5000, font.NewBasic(fsize))
	r.Bxscan([]byte(s), 1024)
	//	r.DumpBoxes()
	return &r
}
func checkbox(t *testing.T, testname string, havebx, wantbx int) {
	if wantbx != havebx {
		t.Logf("%s: have %d, want = %d\n", testname, havebx, wantbx)
		t.Fail()
	}
}

func TestStartLine1(t *testing.T) {
	r := runwith("the quick brown fox jumps over the lazy dog")
	for i := 0; i < r.Nbox; i++ {
		checkbox(t, "TestStartLine1", r.StartLine(i), 0)
	}
}
func TestStartLine2(t *testing.T) {
	r := runwith("the quick brown fox\njumps over the lazy dog")
	checkbox(t, "TestStartLine2", r.StartLine(0), 0)
	checkbox(t, "TestStartLine2", r.StartLine(1), 0)
	checkbox(t, "TestStartLine2", r.StartLine(2), 2)
}
func TestStartLine2TrailingNL(t *testing.T) {
	r := runwith("the quick brown fox\njumps over the lazy dog\n")
	checkbox(t, "TestStartLine2TrailingNL", r.StartLine(2), 2)
	checkbox(t, "TestStartLine2TrailingNL", r.StartLine(3), 2)
}

func TestStartLineCol(t *testing.T) {
	r := runwith("box0\tbox2\tboxx4\tboxxx6\nbx8\tbxA\tthe\tlazy\tdog\n")
	for i := 0; i <= 7; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 0)
	}
	for i := 8; i < 16; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 8)
	}
	checkbox(t, "TestStartLineCol", r.StartLine(18), 18)
}

func TestStartMultiLineCol(t *testing.T) {
	r := runwith("box0\tbox2\tboxx4\tboxxx6\nbx8\tbxA\tthe\tlazy\tdog\nbox18\tbox20\tboxx22\tboxxx24\nbx26\tb28\tthe\tlazy\tdog\n")
	for i := 0; i <= 7; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 0)
	}
	for i := 8; i < 16; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 8)
	}
	checkbox(t, "TestStartLineCol", r.StartLine(18), 18)
	for i := 18 + 0; i <= 18+7; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 18+0)
	}
	for i := 18 + 8; i < 18+16; i++ {
		checkbox(t, "TestStartLineCol", r.StartLine(i), 18+8)
	}
}

func TestFindCell(t *testing.T) {
	r := runwith("\n\n\n\n\n\n\n\n\n\n10\t12\t14\t15\t\nabcdefg\n\nzzzzzzzzzzzzzzzzz")
	for i := 0; i < 10; i++ {
		checkbox(t, "TestFindCell", r.FindCell(i), i)
	}
	checkbox(t, "TestFindCell", r.FindCell(10), 10)
	checkbox(t, "TestFindCell", r.FindCell(11), 10)

	checkbox(t, "TestFindCell", r.FindCell(12), 10)
	checkbox(t, "TestFindCell", r.FindCell(13), 10)

	checkbox(t, "TestFindCell", r.FindCell(16), 10)
	checkbox(t, "TestFindCell", r.FindCell(23), 23)
	checkbox(t, "TestFindCell", r.FindCell(27), 27)
}
