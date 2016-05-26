package percentage

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"github.com/dustin/go-humanize"
)

type Op string

const (
	HUNDRED = 100.0
	MAX     = -92233720368547760
	EOF     = rune(0)
)

type Expr struct {
	left  float64
	Op    Op
	right float64
}

func NewExpr(exp string) (*Expr, error) {
	e, err := parse(exp)
	if err != nil {
		return e, err
	}
	return e, nil
}

func (e *Expr) eval() (float64, error) {
	switch e.Op {
	case "+":
		evaled := toFixed(e.left + percentCalc(e.left, e.right))
		if evaled == MAX {
			return -1, errors.New("too big!")
		}
		return evaled, nil
	case "-":
		evaled := toFixed(e.left - percentCalc(e.left, e.right))
		if evaled == MAX {
			return -1, errors.New("too big!")
		}
		return evaled, nil
	case "*", "x", "X":
		evaled := toFixed(e.left * percentCalc(e.left, e.right))
		if evaled == MAX {
			return -1, errors.New("too big!")
		}
		return evaled, nil
	case "/":
		evaled := toFixed(e.left / percentCalc(e.left, e.right))
		if evaled == MAX {
			return -1, errors.New("too big!")
		}
		return evaled, nil
	case "of":
		evaled := toFixed(percentCalc(e.right, e.left))
		if evaled == MAX {
			return -1, errors.New("too big!")
		}
		return evaled, nil
	case "in":
		evaled, err := inCalc(e.left, e.right)
		if err != nil {
			return -1, err
		}
		return toFixed(evaled), nil
	}
	return -1.0, errors.New("No operator match") // dead
}

func (e *Expr) PrintExpr() string {
	switch e.Op {
	case "+":
		return fmt.Sprintf("%s + %s%%", humanize.Commaf(e.left), humanize.Commaf(e.right))
	case "-":
		return fmt.Sprintf("%s - %s%%", humanize.Commaf(e.left), humanize.Commaf(e.right))
	case "*", "x", "X":
		return fmt.Sprintf("%s × %s%%", humanize.Commaf(e.left), humanize.Commaf(e.right))
	case "/":
		return fmt.Sprintf("%s ÷ %s%%", humanize.Commaf(e.left), humanize.Commaf(e.right))
	case "of":
		return fmt.Sprintf("%s%% of %s", humanize.Commaf(e.left), humanize.Commaf(e.right))
	case "in":
		return fmt.Sprintf("%s in %s", humanize.Commaf(e.left), humanize.Commaf(e.right))
	}
	return "-1.0" // dead
}

func (e *Expr) PrintValue() string {
	evaled, err := e.eval()
	if err != nil {
		return err.Error()
	}
	if e.Op == "in" {
		return fmt.Sprintf("%s%%", humanize.Commaf(evaled))
	}
	return fmt.Sprintf("%s", humanize.Commaf(evaled))
}

type Scanner struct {
	r *bufio.Reader
}

func NewScanner(r io.Reader) *Scanner {
	return &Scanner{bufio.NewReader(r)}
}

func (s *Scanner) read() rune {
	ch, _, err := s.r.ReadRune()
	if err != nil {
		return EOF
	}
	return ch
}

func (s *Scanner) unread() { s.r.UnreadRune() }

func isWhitespace(ch rune) bool { return ch == ' ' || ch == '\t' || ch == '\n' || ch == '%' } // % is a whitespace to us

func isDigit(ch rune) bool { return (ch >= '0' && ch <= '9') || ch == '.' || ch == '-' }

func parse(exp string) (*Expr, error) {
	scanner := NewScanner(strings.NewReader(exp))

	parsedExpr := &Expr{}

	left := new(bytes.Buffer)
	right := new(bytes.Buffer)

	lastWasWhitespace := false // this mess is caused to distinguish between negative numbers and the minus operator
	for {
		ch := scanner.read()
		if isWhitespace(ch) {
			lastWasWhitespace = true
			continue
		}

		if isDigit(ch) {
			if ch == '-' && lastWasWhitespace {
				scanner.unread()
				break
			}

			left.WriteRune(ch)
			lastWasWhitespace = false
			continue
		}

		scanner.unread()
		break
	}

OPLOOP:
	for {
		ch := scanner.read()
		if isWhitespace(ch) {
			continue
		}

		switch ch {
		case '+', '-', '*', 'X', 'x', '/':
			parsedExpr.Op = Op(ch)
			break OPLOOP
		case 'o', 'O':
			ch = scanner.read()
			if ch == 'f' || ch == 'F' {
				parsedExpr.Op = Op("of")
				break OPLOOP
			}
		case 'i', 'I':
			ch := scanner.read()
			if ch == 'n' || ch == 'N' {
				parsedExpr.Op = Op("in")
				break OPLOOP
			}
		}
		scanner.unread()
		break
	}

	for {
		ch := scanner.read()
		if isWhitespace(ch) {
			continue
		}

		if isDigit(ch) {
			right.WriteRune(ch)
			continue
		}

		scanner.unread()
		break
	}

	if parsedExpr.Op == "" {
		return parsedExpr, errors.New("Error: No operator")
	}
	leftFloat, err := strconv.ParseFloat(left.String(), 64)
	if err != nil {
		return parsedExpr, errors.New("Error: Check your left operand")
	}
	rightFloat, err := strconv.ParseFloat(right.String(), 64)
	if err != nil {
		return parsedExpr, errors.New("Error: Check your right operand")
	}

	parsedExpr.left = leftFloat
	parsedExpr.right = rightFloat
	return parsedExpr, nil
}

func percentCalc(x, p float64) float64 {
	return (x / HUNDRED) * p
}

func inCalc(x, y float64) (float64, error) {
	if y == 0 {
		return -1, errors.New("NaN")
	}
	return (x / y) * HUNDRED, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// rounds to a fixed percision of 2
func toFixed(num float64) float64 {
	output := math.Pow(10, 2.0)
	return float64(round(num*output)) / output
}
