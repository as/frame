package tag

import (
	"bytes"
	"fmt"
	window "github.com/as/ms/win"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
	"text/tabwriter"

	"github.com/as/clip"
	"github.com/as/cursor"
)

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func (t *Tag) readfile(s string) (p []byte) {
	var err error
	if isdir(s) {
		fi, err := ioutil.ReadDir(s)
		if err != nil {
			log.Println(err)
			return nil
		}
		sort.SliceStable(fi, func(i, j int) bool {
			if fi[i].IsDir() && !fi[j].IsDir() {
				return true
			}
			ni, nj := fi[i].Name(), fi[j].Name()
			return strings.Compare(ni, nj) < 0
		})
		dx := t.Font.MeasureByte('e')
		x := 0
		b := new(bytes.Buffer)
		w := tabwriter.NewWriter(b, 0, 8, 3, ' ', 0)
		maxx := t.Frame.Bounds().Dx()
		for _, v := range fi {
			nm := v.Name()
			if v.IsDir() {
				nm += string(os.PathSeparator)
			}
			word := fmt.Sprintf("\t%s", nm)
			wordlen := len(word) - 1
			wordpix := wordlen * dx
			advance := max(wordpix, 8*x)
			if x+advance > maxx {
				fmt.Fprintf(w, "\t\n")
				x = -advance
			}
			fmt.Fprintf(w, word)
			x += advance
		}
		w.Flush()
		return b.Bytes()
	}
	p, err = ioutil.ReadFile(s)
	if err != nil {
		log.Println(err)
	}
	return p
}
func writefile(s string, p []byte) {
	fd, err := os.Create(s)
	if err != nil {
		log.Println(err)
	}
	n, err := io.Copy(fd, bytes.NewReader(p))
	if err != nil {
		log.Println(err)
	}
	println("wrote", n, "bytes")
}

func init() {
	var err error
	Clip, err = clip.New()
	if err != nil {
		panic(err)
	}
}
func moveMouse(pt image.Point) {
	cursor.MoveTo(window.ClientAbs().Min.Add(pt))
}
