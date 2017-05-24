package main
// todo: 
//	1) lexer should fill addresses in the following way
//		,	rhs: 0, lhs: max
//        +	rhs: \n, lhs: \n (search backwards)
//		-	rhs:	
//
//

import(
	"regexp"
	"strconv"
	"fmt"
	"strings"
	"bytes"
)
type parser struct{
	in chan item
	out chan func()
	last, tok item
	err error
	stop chan error
	cmd []*command
	addr Address
}

type Address interface{
	Set(f File)
}
type Regexp struct{
	re *regexp.Regexp
	back bool
	rel int
}
type Byte struct{
	Q int64
	rel int
}
type Line struct{
	Q int64
	rel int
}
type Dot struct{
}
type Compound struct{
	a0, a1 Address
	op byte
}

func (p *parser) mustatoi(s string) int64{
	i, err := strconv.Atoi(s)
	if err != nil{
		p.fatal(err)
	}
	return int64(i)
}
func (p *parser) fatal(err error) {
	panic(err)
}

func parseAddr(p *parser) (a Address){
	a0 := parseSimpleAddr(p)
	fmt.Printf("parseAddr:1 %s\n", p.tok)
	p.Next()
	fmt.Printf("parseAddr:2 %s\n", p.tok)
	op, a1 := parseOp(p)
	if op == '\x00' {
		return a0
	}
	p.Next()
	return &Compound{a0: a0, a1: a1, op: op}
}

func parseOp(p *parser) (op byte, a Address){
	fmt.Printf("parseOp:1 %s\n", p.tok)
	if p.tok.kind != kindOp{
		return
	}
	v := p.tok.value
	if v == ""{
		panic("no value"+v)
		return
	}
	if strings.IndexAny(v, "+-;,") == -1 {
		p.fatal(fmt.Errorf("bad op: %q", v))
	}
	p.Next()
	return v[0], parseSimpleAddr(p)
}

func tryRelative(p *parser) int {
	v := p.tok.value
	k := p.tok
	if k.kind == kindRel{
		defer p.Next()
		if v == "+"{
			return 1
		}
		return -1
	}
	return 0
}
// Put
func parseSimpleAddr(p *parser) (a Address) {
	fmt.Printf("parseSimpleAddr:1 %s\n", p.tok)
	back := false
	rel := tryRelative(p)
	v := p.tok.value
	k := p.tok
	fmt.Printf("%s\n", k)
	switch k.kind{
	case kindRegexpBack:
		back = true
		fallthrough
	case kindRegexp:
		re, err := regexp.Compile(v)
		if err != nil{
			p.fatal(err)
			return
		}
		return &Regexp{re, back, rel}
	case kindLineOffset, kindByteOffset:
		i := p.mustatoi(v)
		if rel < 0{
			i = -i
		}
		if k.kind == kindLineOffset{
			return &Line{i, rel}
		}
		return &Byte{i, rel}
	case kindDot:
		return &Dot{}
	}
	p.err = fmt.Errorf("bad address")
	return
}


type File interface{
	Insert(p []byte, at int64) (wrote int64)
	Delete(q0, q1 int64)
	SetSelect(q0, q1 int64)
	Dot() (q0, q1 int64)
	Bytes() []byte
}

type command struct{
	fn func(File)
	s string
	args string
	next *command
}

func (c *command) Next() func(File) {
	return c.next.fn
}
func (c *command) n() *command{
	return c.next
}

func parseArg(p *parser) (arg string){
	fmt.Printf("parseArg: %s\n", p.tok.value)
	p.Next()
	fmt.Printf("parseArg: %s\n", p.tok.value)
	if p.tok.kind != kindArg{
		p.fatal(fmt.Errorf("not arg"))
	}
	return p.tok.value	
}

// Put
func parseCmd(p *parser) (c *command) {
	c = new(command)
	v := p.tok.value
	fmt.Printf("parseCmd: %s\n", v)
	c.s = v
	switch v {
	case "a", "i":
		argv := parseArg(p)
		c.args = argv
		c.fn = func(f File){
			q0, q1 := f.Dot()
			b := []byte(argv)
			if v == "i"{
				fmt.Printf("trace: f.Insert(%q, %d)\n",b,q0)
				f.Insert(b, q0)
			} else {
				fmt.Printf("trace: f.Insert(%q, %d)\n",b,q1)
				f.Insert(b, q1)
			}
		}
		return
	case "c":
		argv := parseArg(p)
		c.args = argv
		c.fn = func(f File){
			q0, q1 := f.Dot()
			f.Delete(q0, q1)
			f.Insert([]byte(argv), q0)
		}
		return
	case "d":
		c.fn = func(f File){
			q0, q1 := f.Dot()
			f.Delete(q0, q1)
		}
		return
	case "e":
	case "g":
	case "k":
	case "s":
	case "v":
	case "x":
		argv := parseArg(p)
		c.args = argv
		re, err := regexp.Compile(argv)
		if err != nil{
			p.fatal(err)
			return
		}
		c.fn = func(f File){
			q0, q1 := f.Dot()
			x0, x1 := q0, q0
			for {
				if x1 > q1{
					break
				}
				buf := bytes.NewReader(f.Bytes()[x1:q1])
				loc := re.FindReaderIndex(buf)
				if loc == nil{
					println("not found")
					break
				}
				x0, x1 = int64(loc[0])+x1, int64(loc[1])+x1
				f.SetSelect(x0,x1)
				a := len(f.Bytes())
				if nextfn := c.Next(); nextfn != nil{
					nextfn(f)
				}
				d0, d1 := f.Dot()
				b := len(f.Bytes()) - a
				x1 += int64(b) + (d1-d0)
				q1 += int64(b)
				x0 = x1
			}
		}
		return
	case "y":
	}
	return nil
}

func (p *parser) Next() *item {
	p.last = p.tok
	p.tok = <- p.in
	return &p.tok
}

func parse(i chan item) (*parser)  {
	p := &parser{
		in: i,
		stop: make(chan error),
	}
	go p.run()
	return p
}

func (p *parser) run() {
	tok := p.Next()
	if tok.kind == kindEof || p.err != nil{
		if tok.kind == kindEof{
			p.fatal(fmt.Errorf("run: unexpected eof"))
			return
		}
		p.fatal(fmt.Errorf("run: %s", p.err))
		return
	}
	p.addr = parseAddr(p)
	for {
		c := parseCmd(p)
		if c == nil{
			break
		}
		p.cmd = append(p.cmd, c)
		fmt.Printf("(%s) %#v and cmd is %#v\n", tok, p.addr, c)
		p.Next()
	}
	p.stop <- p.err
	close(p.stop)
}

func (c *Compound) Set(f File){
	fmt.Printf("compound %#v\n", c)
	c.a0.Set(f)
	q0, _ := f.Dot()
	c.a1.Set(f)
	_, r1 := f.Dot()
	f.SetSelect(q0, r1)
}

func (b *Byte) Set(f File){
	q0, q1 := f.Dot()
	q := b.Q
	if b.rel == -1{
		f.SetSelect(q+q0, q+q0)
	} else if b.rel == 1 {
		f.SetSelect(q+q1, q+q1)
	} else {
		f.SetSelect(q, q)
	}
}
func (r *Regexp) Set(f File){
	_, q1 := f.Dot()
	org := q1
	buf := bytes.NewReader(f.Bytes()[q1:])
	loc := r.re.FindReaderIndex(buf)				
	if loc == nil{
		return
	}
	r0, r1 := int64(loc[0])+org, int64(loc[1])+org
	if r.rel == 1{
		r0=r1
	}
	f.SetSelect(r0,r1)
}
func (r *Line) Set(f File){

}

func (Dot) Set(f File){
}

func compileAddr(a Address) func(f File){
	return a.Set
}

func compile(p *parser) (cmd *command){
	for i := range p.cmd{
		if i+1 == len(p.cmd){
			break
		}
		p.cmd[i].next = p.cmd[i+1]
	}
	fn := func(f File){
		addr := compileAddr(p.addr)
		if addr != nil{
			addr(f)
		}
		if p.cmd != nil && p.cmd[0] != nil && p.cmd[0].fn != nil{
			p.cmd[0].fn(f)
		}
	}
	return &command{fn: fn}
}

func cmdparse(s string) (cmd *command){
	_, itemc := lex("cmd", s)
	p := parse(itemc)
	err := <- p.stop
	if err != nil{
		panic(err)
	}
	return compile(p)
	//return zcmdparse(s)
}