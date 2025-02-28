package base

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"unicode"
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
func readMessages(path string) {

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
			Report(`<Bug>Invalid message string: %#v>`)
			panic("Bug")
		}
		msgt = I18nMessage{rm[1], rm[2]}
		logbase.LangMap[msg] = msgt
	}
	return fmt.Sprintf(msgt.text, args...), msgt.tag
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
			"DONE":   "",
			"REPORT": tag,
			"TEXT":   msgt,
		}
	}
}
