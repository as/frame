package main

import (
	"strings"
	"fmt"
	"strconv"
)

type Cmd struct{
	kind string
	data interface{}
}

type Addr struct{
	Q0, Q1 int64
}
func (a Addr) String() string{
	return fmt.Sprintf("[%d:%d]", a.Q0, a.Q1)
}

func addrparse(s string) (addr Addr){
	const Max = int64(^0)
	s = strings.TrimSpace(s)
	if s == "," {
		return Addr{0, Max}
	}
	q0, ok := number(s)
	if ok{
		return Addr{q0, q0}
	}
	i := strings.IndexAny(s, ",")
	if i < 0{
		return Addr{}
	}
	q0, ok = number(s[:i])
	if !ok{
		return Addr{}
	}
	s = s[i:]
	if len(s) == 0{
		return Addr{q0, Max}
	}
	q1, ok := number(s[1:])
	if !ok{
		return Addr{}
	}
	return Addr{q0, q1}
}
func number(c string) (int64, bool){
	i, err := strconv.Atoi(strings.TrimSpace(c))
	if err != nil{
		fmt.Println(err)
		return 0, false
	}
	return int64(i), true
}
func zcmdparse(s string) (cmd Cmd){
	s = s+" "
	i := strings.IndexAny(s, " \t")
	if i < 0{
		return 
	}
	c := string(s[:i])
	switch {
	case c == "a":    // append
	case c == "i":    // insert
	case c == "d":	// delete
	case c == "x":	// eXtract
	case c == "y":	// extract complement
	case c == "p":	// print
	case c[0] == '#': // char offset
	case c[0] == '/': // regexp
	case c[0] == '?': // regexp (back)
	default:
		cmd.kind = "s"
		cmd.data = addrparse(s)
		return
	}
	data := string(s[i:])
	if len(data)==0{
		return Cmd{c, ""}
	}
	return Cmd{c, data[1:]}
}
