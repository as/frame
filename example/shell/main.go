package main

import (
	//	"github.com/as/clip"
	"golang.org/x/exp/shiny/driver"
	"golang.org/x/exp/shiny/screen"
	"golang.org/x/mobile/event/key"
	"golang.org/x/mobile/event/lifecycle"
	"golang.org/x/mobile/event/mouse"
	"golang.org/x/mobile/event/paint"
	"golang.org/x/image/font"
	"golang.org/x/mobile/event/size"
	"image"
	"fmt"
	"os/exec"
	"strings"

	"github.com/as/frame"
)


var winSize = image.Pt(1920, 1080)

type Fd struct {
	inc  chan []byte
	outc chan []byte
	errc chan []byte
}

func Run(fd *Fd, s string, args ...string) {
	cmd := exec.Command(s, args...)
	in, _ := cmd.StdinPipe()
	out, _ := cmd.StdoutPipe()
	e, _ := cmd.StderrPipe()
	killc := make(chan bool)
	go func() {
		for {
			select{
			case <- killc:
				return
			case v := <- fd.inc:
				_, err := in.Write(v)
				if err != nil {
					return
				}
			}
		}
	}()
	go func() {
		p := make([]byte, 8192)
		for {
			select {
			case <-killc:
				return
			default:
				n, err := out.Read(p)
				if n > 0 {
					fd.outc <- append([]byte{}, p[:n]...)
				}
				if err != nil {
					return
				}
			}
		}
	}()
	go func() {
		p := make([]byte, 8192)
		for {
			select {
			case <-killc:
				return
			default:
				n, err := e.Read(p)
				if n > 0 {
					fd.errc <- append([]byte{}, p[:n]...)
				}
				if err != nil {
					return
				}
			}
		}
	}()
	
	err := cmd.Start()
	if err != nil {
		fd.errc <- []byte(fmt.Sprintf("%s\n", err))
	    close(killc)
		return
	}
	cmd.Wait()
	close(killc)
}

func main() {
	fmt.Print()
	selecting := false
	focused := false

	
	driver.Main(func(src screen.Screen) {
		win, _ := src.NewWindow(&screen.NewWindowOptions{winSize.X, winSize.Y})
		tx, _ := src.NewTexture(winSize)
		buf, _ := src.NewBuffer(winSize)
		fr := frame.New(buf.RGBA(), image.Pt(25, 25), nil)
		measure := func(p []byte) int{
			x := font.MeasureBytes(fr.Font, p)
			return int(x >> 6)
		}
	boxes := frame.NewBoxes(measure)
	for i := 0; i < 5; i++{
		test := make([]byte, 20)
		boxes.WriteAt([]byte(fmt.Sprintf("%d", i)), int64(i))
		boxes.ReadAt(test, int64(i))
		fmt.Printf("%q\n",test)
	}
	for i := 0; i < 5; i++{
		test := make([]byte, 20)
		boxes.ReadAt(test, int64(i))
		fmt.Printf("%q\n",test)
	}
			tick := &frame.Tick{
			Fr: fr, 
			Select: frame.Select{
				Img: image.NewRGBA(fr.Bounds()),
			},
		}
		fr.Tick = tick

		con := make(chan []byte)
		conin := make(chan []byte)

		go func(){
			for p := range con{
				tick.Write(p)
				win.Send(paint.Event{})			
			}
		}()
		
		go func() {
			for {
				select {
				case p := <-conin:
					s := strings.Fields(string(p))
					fmt.Println()
					fmt.Println("conin:", s)
					if len(s) == 1 && s[0] == "clear" || s[0] == "cls" {
						tick.P0 = 0
						tick.P1 = len(tick.Fr.Bytes())
						if tick.P1 >= 0{
							tick.Delete()
						}
					} else if len(s) > 0 {
						var args []string
						if len(s) > 1 {
							args = s[1:]
						}
						tick.P0 = tick.P1
						con <- []byte{'\n'}
						Run(&Fd{inc: conin, outc: con, errc: con}, s[0], args...)
					}
					con <- []byte{';'}
				}
			}
		}()
		apos := image.ZP
		for {
			switch e := win.NextEvent().(type) {
			case key.Event:
				if e.Direction != key.DirPress && e.Direction != key.DirNone {
					break
				}
				switch e.Code {
				case key.CodeRightArrow:
					if e.Modifiers != key.ModShift {
						tick.P0++
					}
					tick.P1++
				case key.CodeLeftArrow:
					if e.Modifiers != key.ModShift {
						tick.P0--
					}
					tick.P1--
				case key.CodeDeleteBackspace:
					tick.Delete()
				case key.CodeReturnEnter:
					if tick.P1 == tick.P0 {
						for tick.P0 >= 0 {
							c := tick.Fr.Bytes()[tick.P0]
							if c == ';' {
								tick.P0++
								break
							}
							if c == '\n' {
								break
							}
							tick.P0--
						}
					}
					conin <- []byte(tick.String())
				default:
					if e.Rune != -1 {
						tick.WriteRune(e.Rune)
					}
				}
				win.Send(paint.Event{})
			case mouse.Event:
				apos = image.Pt(int(e.X), int(e.Y))
					if selecting{
						tick.P1 = fr.IndexOf(apos)
						fr.Dirty = true
					}
				if e.Button == mouse.ButtonLeft {
					if e.Direction == mouse.DirPress {
						tick.Close()
						tick.P1 = fr.IndexOf(apos)
						tick.SelectAt(tick.P1)
						tick.P0 = tick.P1
						fr.Dirty = true
						selecting = true
					}
					if e.Direction == mouse.DirRelease{
						selecting = false
					}
				}
				win.Send(paint.Event{})
			case size.Event, paint.Event:
				if fr.Dirty {
					fr.Redraw(selecting, apos)
					tx.Upload(image.ZP, buf, buf.Bounds())
					win.Copy(buf.Bounds().Min, tx, tx.Bounds(), screen.Over, nil)
				}
				if !focused{
					win.Copy(buf.Bounds().Min, tx, tx.Bounds(), screen.Over, nil)
				}
				win.Publish()
			case lifecycle.Event:
				if e.To == lifecycle.StageDead {
					return
				}
				
				// NT doesn't repaint the window if another window covers it
				if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOff{
					focused = false
				} else if e.Crosses(lifecycle.StageFocused) == lifecycle.CrossOn{
					focused = true
				}
			}
		}
	})
}
