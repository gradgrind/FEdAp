package base

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
)

const (
	Error   = `<>Error>`
	Warning = `<>Warning>`
	Notice  = `<>Notice>`
	InfoTag = `<>Info>`
	Bug     = `<>Bug>`
)

var (
	tregexp  *regexp.Regexp
	logbase  *LogBase
	nlregexp *regexp.Regexp
)

func init() {
	tregexp = regexp.MustCompile("(?s)^<([a-zA-Z]*)>(.+)>$")
	nlregexp = regexp.MustCompile(`[\t \f\r]*[\n][\t\n \f\r]*>?`)
}

type I18nMessage struct {
	tag  string
	text string
}

type LogBase struct {
	Logger  *log.Logger
	LangMap map[string]I18nMessage
	Channel chan map[string]any
}

func OpenLog(ochan chan map[string]any, logpath string) {
	var file *os.File
	if logpath == "" {
		file = os.Stderr
	} else {
		os.Remove(logpath)
		var err error
		file, err = os.OpenFile(logpath, os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}
	logbase = &LogBase{
		Logger:  log.New(file, "++", 0),
		LangMap: map[string]I18nMessage{},
		Channel: ochan,
	}
}

// TODO: Read message file
// An attempt could be made to read the system language (e.g. using
// https://github.com/Xuanwo/go-locale), or else the language could be set
// explicitly (perhaps using a persistent settings file).
// Language files are basically just key/value look-up tables, a key-line
// (the source text) being followed by its value line (the translation).
// Key line:
//
//	<<preamble>Message
//
// Value line:
//
//	>>Translated Message
//
// The value lines can be omitted if no translation is available.
// Lines may not be split, even if they are very long, but newlines can be
// specified by "\n" (also tab with "\t"). The %-escapes (of the go "fmt"
// package) should be preserved, because there will be corresponding
// arguments.
// The preamble consists of a line number and the message type, separated by
// a ":". The look-up key should be constructed to be the same as in the
// source, i.e. thus:
//
//	<message type>Message>
//
// There are also package/file lines:
//
//	::package-name@file-name
//
// These serve only as documentation, so that the source lines where the
// messages originate can be found easily.
func ReadMessages(path string) bool {
	readFile, err := os.Open(path)
	if err != nil {
		Report(`<Error>Reading translations: %v>`, err)
		return false
	}
	defer readFile.Close()
	logbase.LangMap = map[string]I18nMessage{}

	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)
	var sourceLine *I18nMessage
	l := 0
	for fileScanner.Scan() {
		l++
		line := fileScanner.Text()
		if strings.HasPrefix(line, "<<") {
			sl := strings.SplitN(line[2:], ">", 2)
			if len(sl) == 2 {
				pl := strings.Split(sl[0], ":")
				if len(pl) == 2 {
					sourceLine = &I18nMessage{pl[1], sl[1]}
					continue
				}
			}
			Report(
				`<Error>Invalid line (%d) in messages file, %s:\n  -- %s>`,
				l, path, line)
			continue
		}
		if strings.HasPrefix(line, ">>") {
			if sourceLine != nil {
				key := "<" + sourceLine.tag + ">" + sourceLine.text + ">"
				logbase.LangMap[key] = I18nMessage{sourceLine.tag, line[2:]}
				sourceLine = nil
			}
		}
	}
	return true
}

// I18N looks up a message in the message catalogue, performing value
// substitutions. Note that for this to work in conjunction with the "gettr"
// utility the message string must be a single back-tick enclosed string
// where the only permitted "escape" characters are "\n" and "\t". Control
// characters will be removed, newlines including all surrounding whitespace.
// If a new line starts with ">" this will be stripped (this is to allow
// indentation to be used, preserving the whitespace only after the ">").
// The message needs to have the structure defined by the regular expression
// [tregexp]. The main part and the tag prefix are returned separately.
// The main part can have formatting escapes, whose matching parameters need
// to be provided as arguments to this function.
func I18N(msg string, args ...any) (string, string) {
	// Preprocess message string (merge lines, remove control characters)
	msg = nlregexp.ReplaceAllString(msg, "")
	msg = strings.Map(func(r rune) rune {
		if unicode.IsPrint(r) {
			return r
		}
		return -1
	}, msg)
	// Look up message
	msgt, ok := logbase.LangMap[msg]
	if !ok {
		// Add the untranslated string to the message map
		rm := tregexp.FindStringSubmatch(msg)
		if rm == nil {
			Report(`<Bug>Invalid message string: %#v>`, msg)
			panic("Bug")
		}
		msgt = I18nMessage{rm[1], rm[2]}
		logbase.LangMap[msg] = msgt
	}
	return fmt.Sprintf(msgt.text, args...), msgt.tag
}

// Tr looks up the given string in the message map, performing argument
// substitution, but ignoring the message's type tag.
// See function [I18N] for further details.
func Tr(text string, args ...any) string {
	tr, _ := I18N(text, args...)
	return tr
}

// Report logs a message. This uses [I18N()] to support translations by
// looking up the messages in [logbase.LangMap]. The messages must have a
// prefix enclosed in angle brackets to indicate the type of the message.
// They must also be terminated by a right angle bracket.
// The currently supported prefixes are:
//
//	"<Error>", "<Warning>", "<Notice>: These will force the process
//			pop-up to open.
//	"<Info>": This is like "Notice", but the process window will not open
//	      if the operation completes quickly enough.
//	"<Bug>": This will cause a special message pop-up to be shown.
func Report(msg string, args ...any) {
	// Look up message
	msgt, tag := I18N(msg, args...)
	logbase.Logger.Println(tag + ">" + msgt)

	// Send to back-end interface
	if logbase.Channel != nil {
		logbase.Channel <- map[string]any{
			"REPORT": tag,
			"TEXT":   msgt,
			"TR":     Tr("<>" + tag + ">"),
		}
	}
}
