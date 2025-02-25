package base

import (
	"fmt"
	"log"
	"os"
)

// TODO: Replace these loggers by [Report], using logbase
var (
	//	Message *log.Logger
	//	Warning *log.Logger
	//	Error   *log.Logger
	//	Bug     *log.Logger

	logbase *LogBase
)

type LogBase struct {
	Logger   *log.Logger
	LangMap  map[string]string
	Fallback map[string]string
}

func OpenLog(logpath string) {
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

	logbase.Logger = log.New(file, "++", log.Lshortfile)

	//Message = log.New(file, "*INFO* ", log.Lshortfile)
	//Warning = log.New(file, "*WARNING* ", log.Lshortfile)
	//Error = log.New(file, "*ERROR* ", log.Lshortfile)
	//Bug = log.New(file, "*BUG* ", log.Lshortfile)
}

// I18N looks up a message in the message catalogue, performing value
// substitutions.
func I18N(msg string, args ...any) string {
	// Look up message
	msgt, ok := logbase.LangMap[msg]
	if !ok {
		msgt, ok = logbase.Fallback[msg]
		if !ok {
			Report("<Bug>Unknown message: %[1]s ::: %+[2]v>", msg, args)
			panic("Bug")
		}
	}
	return fmt.Sprintf(msgt, args...)
}

// TODO
// Report logs a message. The keys should have a prefix to indicate the
// type of the error and also the messages themselves should have an
// appropriate prefix:
//
//	"ERROR_"   -> "[Error] ..."
//	"INFO_"    -> "[Info] ..."
//	"WARNING_" -> "[Warning] ..."
//	"BUG_"     -> "[Bug] ..."
//
// The message prefixes may be translated.
func Report(msg string, args ...any) {
	// Look up message
	msgt := I18N(msg, args...)
	logbase.Logger.Println(msg + "#" + msgt)
}

// Tr adds message strings to the Fallback map of logbase, initializing
// logbase if necessary.
// It can be called from init functions.
func Tr(trmap map[string]string) {
	lg := logbase
	if lg == nil {
		lg = &LogBase{
			//Logger: log.New(file, "++", log.Lshortfile),
			//TODO: load the data from somewhere ...
			LangMap:  map[string]string{},
			Fallback: map[string]string{},
		}
		logbase = lg
	}
	for k, v := range trmap {
		if _, nok := lg.Fallback[k]; nok {
			Report("<Bug>Message defined twice: %s>", k)
			panic("Bug")
		}
		lg.Fallback[k] = v
	}
}
