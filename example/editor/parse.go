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
)
type parser struct{
	in chan item
	out chan func()
	last, tok item
	err error
}

type Address interface{
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
	p.Next()
	op, a1 := parseOp(p)
	if op == '\x00' {
		return a0
	}
	return &Compound{a0: a0, a1: a1, op: op}
}

func parseOp(p *parser) (op byte, a Address){
	if p.tok.kind != kindOp{
		return
	}
	v := p.tok.value
	if v == ""{
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

func parseSimpleAddr(p *parser) (a Address) {
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
		return Regexp{re, back, rel}
	case kindLineOffset, kindByteOffset:
		i := p.mustatoi(v)
		if rel < 0{
			i = -i
		}
		if k.kind == kindLineOffset{
			return Line{i, rel}
		}
		return Byte{i, rel}
	case kindDot:
		return Dot{}
	}
	p.Next()
	p.err = fmt.Errorf("bad address")
	return
}


type File interface{
	Insert(p []byte, at int64) (wrote int64)
	Delete(q0, q1 int64)
	SetSelect(q0, q1 int64)
	Dot() (q0, q1 int64)
}

type command struct{
	fn func(File)
	s string
}

func parseCmd(p *parser) (c command) {
	v := p.tok.value
	c.s = v
	switch  v {
	case "a", "i":
		tok := p.Next()
		if tok.kind != kindArg{
			p.fatal(fmt.Errorf("not arg"))
		}
		argv := p.tok.value
		c.fn = func(f File){
			q0, q1 := f.Dot()
			if v == "i"{
				f.Insert([]byte(argv), q0)
			} else {
				f.Insert([]byte(argv), q1)
			}
		}
		return
	case "c":
		println("c")
	case "d":
	case "e":
	case "k":
	case "s":
	case "x":
		println("x")
	case "y":
	}
	return
}

func (p *parser) Next() *item {
	p.last = p.tok
	p.tok = <- p.in
	return &p.tok
}

func parse(i chan item) (*parser)  {
	p := &parser{
		in: i,
	}
	go p.run()
	return p
}

func (p *parser) run() {
	for {
		tok := p.Next()
		if tok.kind == kindEof || p.err != nil{
			break
		}
		addr0 := parseAddr(p)
		p.Next()
		cmd := parseCmd(p)
		fmt.Printf("(%s) %#v and cmd is %#v\n", tok, addr0, cmd)
//		fn(p)()
//		p.Push(fn(p))
	}
}
func cmdparse(s string) (cmd Cmd){
	_, itemc := lex("cmd", s)
	parse(itemc)
	return zcmdparse(s)
}