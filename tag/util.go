package tag

import (
	"bytes"
	"fmt"
	"image"
	"io"
	"io/ioutil"
	"log"
	"os"
	"text/tabwriter"
	"sort"
	"strings"
	window "github.com/as/ms/win"

	"github.com/as/clip"
	"github.com/as/cursor"
)

func (t *Tag) readfile(s string) (p []byte) {
	var err error
	if isdir(s) {
		fi, err := ioutil.ReadDir(s)
		if err != nil {
			log.Println(err)
			return nil
		}
		sort.SliceStable(fi, func(i, j int) bool{
			if fi[i].IsDir() && !fi[j].IsDir(){
				return true
			}
			ni,nj := fi[i].Name(), fi[j].Name()
			return strings.Compare(ni, nj) < 0
		})
		x := t.Font.MeasureByte('e')
		n := t.Frame.Bounds().Dx()/x
		m := 0
		b := new(bytes.Buffer)
		w := tabwriter.NewWriter(b, 0, 0, 3, ' ', 0)
		for _, v := range fi {
			nm := v.Name()
			if v.IsDir(){
				nm += string(os.PathSeparator)
			}
			entry := fmt.Sprintf("%s\t", nm)
			dm := m+len(entry)
			if dm > n{
				fmt.Fprintf(w, "\n")
				m=-(dm+3)
			}
			fmt.Fprintf(w, entry)
			m+=dm+3
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
