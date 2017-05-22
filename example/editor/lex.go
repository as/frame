package main

import(
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type item struct{
	kind Kind
	value string
}
func (i item) String() string {
	return fmt.Sprintf("%s %s", i.kind, i.value)
}

const MaxBytes = 1<<63-1

func max() string{
	return fmt.Sprintf("%v", MaxBytes)
}

type Kind int
const(
	kindOp    Kind = iota
	kindString
	kindSlash
	kindQuest
	kindRel
	kindComma
	kindDot
	kindEof
	kindColon
	kindSemi
	kindHash
	kindErr
	kindRegexp
	kindRegexpBack
	kindByteOffset
	kindLineOffset
	kindCmd
	kindArg
)
const(
	eof = '\x00'
	slash = '/'
	quest = '?'
	comma = ','
	plus = "+"
	dot = '.'
	colon = ':'
	semi = ';'
	hash = '#'
)

type statefn func(*lexer) statefn

type lexer struct {
	name string
	input	string
	start int
	pos int
	width int
	items chan item
	lastop item
	first bool
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name: name,
		input: input,
		items: make(chan item),
		lastop: item{kindOp, "+"},
		first: true,
	}
	go l.run()	// run state machine
	return l, l.items
}

func (l *lexer) run() {
	for state := lexAny; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func (l *lexer) acceptUntil(delim string){
	lim := 8192
	i := 0
	for strings.IndexRune(delim, l.next()) < 0 {
		i++
		if i > lim{
			l.errorf("missing terminating char %q: %q\n", delim, l)
			l.ignore()
			l.emit(kindEof)
			return
		}
	}
	l.backup()	
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) emit(t Kind) {
	l.items <- item{t, l.String()}
	l.start = l.pos
}

func (l *lexer) inject(it item){
	l.items <- it
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) String() string{
	return string(l.input[l.start:l.pos])
}

func space(r rune) bool {
	return unicode.IsSpace(r)
}

func ignoreSpaces(l *lexer) {
	if l.accept(" 	") {
		l.acceptRun(" 	")
		l.ignore()
	}
}

const(
	Ralpha = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	Rdigit = "0123456789"
	Rop = "+-;,"
	Rmod = "#/?"
)

func lexAny(l *lexer) statefn{
	if l.accept(Rdigit+Rop+Rmod){
		l.backup()
		return lexAddr
	}
	l.emit(kindDot)
	return lexCmd
}


func lexAddr(l *lexer) statefn{
	ignoreSpaces(l)
	switch l.peek(){
	case eof:
		l.next()
		l.emit(kindEof)
		return nil
	case ',', ';':
		// LHS is empty so use #0
		if l.first{
			l.inject(item{kindByteOffset, "0"})
			l.first = false
		}
		return lexOp
	case '+', '-':
		if l.first{
			l.first = false
		}
		l.accept("+-")
		l.emit(kindRel)
		return lexAddr
	case slash, quest:
		return lexRegexp
	case dot:
		l.accept(".")
		l.emit(kindDot)
		return lexOp
	case hash:
		l.accept("#")
		l.ignore()
		ignoreSpaces(l)
		if !l.accept(Rdigit){
			return l.errorf("non-numeric offset")
		}
		l.acceptRun(Rdigit)
		l.emit(kindByteOffset)
		return lexOp
	default:
		if l.accept(Rdigit){
			l.acceptRun(Rdigit)
			l.emit(kindLineOffset)
			return lexOp
		}		
	}
	return lexCmd
}

func lexCmd(l *lexer) statefn{
	ignoreSpaces(l)
	if l.peek() == eof{
		l.emit(kindEof)
		return nil
	}
	if !l.accept(Ralpha){
		return l.errorf("bad command")
	}
	l.emit(kindCmd)
	switch l.peek(){
	case eof:
		l.emit(kindEof)
		return nil
	default:
		return lexArg
	}
}

func lexOp(l *lexer) statefn{
	ignoreSpaces(l)
	if l.peek() == eof{
		l.emit(kindEof)
		return nil
	}
	op := ""
	if l.accept(Rop) {
		op = l.String()
		l.emit(kindOp)
	} 
	ignoreSpaces(l)
	if l.accept(Rdigit+Rmod) {
		if op == ""{
			l.inject(l.lastop)
		}
		l.backup()
		l.lastop = item{kindOp, op}
	}
	if op != "" && (l.accept(Ralpha) || l.peek() == eof) {
		l.inject(item{kindByteOffset, max()})
		l.backup()
	}
	return lexAddr
}

func lexArg(l *lexer) statefn{
	r := string(l.next())
	l.ignore()
	l.acceptUntil(r)
	l.emit(kindArg)
	if !l.accept(string(r)){
		return l.errorf("bad delimiter")
	}
	l.ignore()
	return lexCmd
}

func lexRegexp(l *lexer) statefn{
	r := l.next()
	if r != '?' && r != '/'{
		return l.errorf("bad regexp delimiter: %q", l)
	}
	l.ignore()
	l.acceptUntil(string(rune(r)))
	if r == '?'{
		l.emit(kindRegexpBack)
	} else {
		l.emit(kindRegexp)
	}
	if !l.accept(string(rune(r))){
		return l.errorf("bad regexp terminator: %q", l)
	}
	l.ignore()
	return lexOp
}

func (l *lexer) errorf(format string, args ...interface{}) statefn {
	l.items <- item {
		kindErr,
		fmt.Sprintf(format, args...),
	}
	return nil
}

/*

func AlphaNumeric(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

func lexText(l *lexer) statefn {
	panic("lexText")
	ignoreSpaces(l)
	l.acceptRun(Ralpha + Rdigit + Rop + Rmod)
	if l.pos == l.start && l.next() == eof {
		l.emit(itemEOF)
		return nil
	}
	ignoreSpaces(l)
	l.emit(kindText)

	fmt.Printf("%#v", l.input[l.pos:])
	switch r := l.peek(); {
	case r == '=':
		return lexEquals
	case r == '$':
		return lexEnv
	case r == eof:
		if l.pos == l.start {
			println("itemEOF")
			l.emit(itemEOF)
			return nil
		}
	}
	return nil
}

func lexIdentifier(l *lexer) statefn {
	panic("lexIdentifier")
	l.acceptRun("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	if AlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad identifier syntax: %q", l)
	}
	l.emit(kindText)
	return lexInsideAction
}

func lexInsideAction(l *lexer) statefn {
	panic("lexInsideAction")
	// Either num, string, or id
	for {
		if strings.HasPrefix(l.input[l.pos:], rightMeta) {
			return lexRightMeta
		}
		switch r := l.next(); {
		case r == eof || r == '\n':
			return l.errorf("unclosed action")
		case unicode.IsSpace(r):
			l.ignore()
		case r == '|':
			l.emit(kindPipe)
		case r == '+' || r == '-' || '0' <= r && r <= '9':
			l.backup()
			return lexNumber
		case AlphaNumeric(r):
			l.backup()
			return lexIdentifier
		}
	}
}

func lexLeftMeta(l *lexer) statefn {
	l.pos += len(leftMeta)
	l.emit(kindLeftMeta)
	return lexInsideAction
}

func lexRightMeta(l *lexer) statefn {
	l.pos += len(rightMeta)
	l.emit(kindRightMeta)
	return lexText
}

func lexNumber(l *lexer) statefn {
	l.accept("+-")
	digits := "0123456789"
	if l.accept("0") && l.accept("xX") {
		digits += "abcdefABCDEF"
	}
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if l.accept("eE") {
		l.accept("+-")
		l.acceptRun("0123456789")
	}
	// imaginary
	if AlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number syntax: %q", l)
	}
	l.emit(kindNumber)
	return lexInsideAction
}
*/