package main

import (
	"fmt"
	"strconv"
	"strings"
)

type Cmd struct {
	kind string
	data interface{}
}

type Addr struct {
	Q0, Q1 int64
}

func (a Addr) String() string {
	return fmt.Sprintf("[%d:%d]", a.Q0, a.Q1)
}

func addrparse(s string) (addr Addr) {
	const Max = int64(^0)
	s = strings.TrimSpace(s)
	if s == "," {
		return Addr{0, Max}
	}
	q0, ok := number(s)
	if ok {
		return Addr{q0, q0}
	}
	i := strings.IndexAny(s, ",")
	if i < 0 {
		return Addr{}
	}
	q0, ok = number(s[:i])
	if !ok {
		return Addr{}
	}
	s = s[i:]
	if len(s) == 0 {
		return Addr{q0, Max}
	}
	q1, ok := number(s[1:])
	if !ok {
		return Addr{}
	}
	return Addr{q0, q1}
}
func number(c string) (int64, bool) {
	i, err := strconv.Atoi(strings.TrimSpace(c))
	if err != nil {
		fmt.Println(err)
		return 0, false
	}
	return int64(i), true
}
