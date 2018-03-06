package box

import (
	"testing"

	. "github.com/as/font"
)

func runwith(s string) *Run {
	r := NewRun(5, 5000, NewGoMono(fsize))
	r.Boxscan([]byte(s), 1024)
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

/*
func TestStartCell(t *testing.T) {
	r := runwith("\n\n\n\n\n\n\n\n\n\n10\t12\t14\t15\t\nabcdefg\n\nzzzzzzzzzzzzzzzzz")
	for i := 0; i < 10; i++ {
		checkbox(t, "TestStartCell", r.StartCell(i), i)
	}
	checkbox(t, "TestStartCell", r.StartCell(10), 10)
	checkbox(t, "TestStartCell", r.StartCell(11), 10)

	checkbox(t, "TestStartCell", r.StartCell(12), 10)
	checkbox(t, "TestStartCell", r.StartCell(13), 10)

	checkbox(t, "TestStartCell", r.StartCell(16), 10)
	checkbox(t, "TestStartCell", r.StartCell(23), 23)
	checkbox(t, "TestStartCell", r.StartCell(27), 27)
}
*/
func TestStartCell2(t *testing.T) {
	r := runwith("\nAAA\tBBB\tCCC")
	//r.DumpBoxes()
	checkbox(t, "10", r.EndCell(3), 6)
	checkbox(t, "20", r.EndLine(3), 6)
	checkbox(t, `\nAAA\tBBB\tCCC`, r.StartCell(3), 1)
	r = runwith("AAA\tBBB\tCCC")
	//r.DumpBoxes()
	checkbox(t, "10", r.EndCell(3), 5)
	checkbox(t, "20", r.EndLine(3), 5)
	checkbox(t, `AAA\tBBB\tCCC`, r.StartCell(2), 0)
}

func TestEndCell(t *testing.T) {
	r := runwith("\n\n\n\n\n\n\n\n\n\n10\t12\t14\t15\t\nabcdefg\n\nzzzzzzzzzzzzzzzzz")

	for i := 0; i < 10; i++ {
		checkbox(t, "TestEndCell", r.EndCell(i), i)
	}
	checkbox(t, "TestEndCell", r.EndCell(10), 18)
	checkbox(t, "TestEndCell", r.EndCell(11), 18)

	checkbox(t, "TestEndCell", r.EndCell(12), 18)
	checkbox(t, "TestEndCell", r.EndCell(13), 18)

	checkbox(t, "TestEndCell", r.EndCell(16), 18)
	checkbox(t, "TestEndCell", r.EndCell(23), 23)
	checkbox(t, "TestEndCell", r.EndCell(27), 27)
}

func TestNextCell(t *testing.T) {
	r := runwith("\n\n\n\n\n\n\n\n\n\n10\t12\t14\t15\t\nabcdefg\n\nzzzzzzzzzzzzzzzzz")

	for i := 0; i < 10; i++ {
		checkbox(t, "TestNextCell", r.NextCell(i), 10)
	}
	checkbox(t, "TestNextCell", r.NextCell(10), 23)
	checkbox(t, "TestNextCell", r.NextCell(11), 23)

	checkbox(t, "TestNextCell", r.NextCell(12), 23)
	checkbox(t, "TestNextCell", r.NextCell(13), 23)

	checkbox(t, "TestNextCell", r.NextCell(16), 23)
	checkbox(t, "TestNextCell", r.NextCell(23), 23)
}

func TestNextCell2(t *testing.T) {
	r := runwith("\tfmt.Println(\"hello world\")\n}\none\ttwo\three")
	checkbox(t, "10", r.EndCell(1), 1)
	checkbox(t, "20", r.StartCell(6), 5)
	checkbox(t, "30", r.EndLine(6), 10)
	checkbox(t, "40", r.EndCell(6), 10)
	checkbox(t, "50", r.StartCell(6), 5)
	checkbox(t, "60", r.NextCell(1), 5)
}

/*
func TestStretch1(t *testing.T) {
	r := runwith("AAA\tBBB\tCCC")
	r.Stretch(4)
	t.Fail()
}

func TestStretch2(t *testing.T) {
	r := runwith("AAA\tBBB\tCCC\n")
	r.Stretch(4)
	t.Fail()
}

func TestStretch3(t *testing.T) {
	r := runwith("AAA\tBBB\tCCC\nDDD\tEEE\tFFF")
	r.Stretch(4)
	t.Fail()
}
*/
