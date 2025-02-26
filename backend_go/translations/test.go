package main

import (
	"fmt"
	"strconv"
	"strings"
	"regexp"
)

func main() {
	nlregexp := regexp.MustCompile(`[\t \f\r]*[\n][\t\n \f\r]*>?`)
	m0 := `Invalid "message" string:
	
		>\nx\t %#v`
	m0s := nlregexp.ReplaceAllString(m0, "")
	
	m1 := strings.ReplaceAll(m0s, "\\n", "\n")
	m2 := strings.ReplaceAll(m1, "\\t", "\t")
	fmt.Println(m2)
	//fmt.Println(strconv.Unquote("`" + m0 + "`"))
	
	rx, err := regexp.Compile(`[\t \f\r]*[\n][\t \f\r]*>?`)
	if err != nil { panic(err) }

	s0 := `First line\tx
	 > Second line	y
Third line\nFourth line\t	z`
	s0 = rx.ReplaceAllString(s0, "\\n")
	fmt.Println(s0)
	
	sq := "\"" + s0 + "\""
	fmt.Println(sq)
	
	//sq = strings.ReplaceAll(sq, "\n", "\\n")
	s, err := strconv.Unquote(sq)
	fmt.Printf("%q, %v\n", s, err)
	
	fmt.Println(s)
	
	fmt.Println()
	
	s1 := `"Invalid \"message\" string:\n\t %#v"`
	fmt.Println(s1)
	
	fmt.Println(strconv.Quote("Invalid \"message\" string:\n\t %#v"))
	fmt.Println(strconv.Quote(`Invalid \"message\" string:\n\t %#v`))
	fmt.Println(strconv.Quote(`Invalid "message" string:\n\t %#v`))
	s, err = strconv.Unquote(s1)
	fmt.Printf("%q, %v\n", s, err)
    fmt.Println(s)
}
