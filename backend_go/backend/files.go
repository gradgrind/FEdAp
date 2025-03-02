package backend

import (
	"gradgrind/backend/base"
	"gradgrind/backend/ttbase"
	"os"
	"path/filepath"
)

func LoadFile(
	cmd map[string]any,
	loader func(*base.DbTopLevel, string),
) bool {
	path0, ok := cmd["FILEPATH"]
	var path string
	if ok {
		path, ok = path0.(string)
	}
	if !ok || path == "" {
		base.Report(`<Error>Opening data file: no file path>`)
		return false
	}
	// Actually, the supplied path should already be absolute ...
	abspath, err := filepath.Abs(path)
	if err != nil {
		base.Report(`<Error>Couldn't resolve file path: %s\n  +++ %v>`,
			path, err)
		return false
	}

	// Make a new database structure
	db := base.NewDb()
	// Load the data using ths supplied function
	loader(db, abspath)
	db.PrepareDb()
	if db.Invalid {
		return false
	}
	DB = db

	// By default the timetable data is loaded into TtData and further
	// checked by testing fixed allocations, etc. These steps can, however,
	// be skipped.
	tt, ok := cmd["TIMETABLE"]
	if ok && tt.(string) == "NO" {
		TtData = nil
	} else {
		TtData = ttbase.MakeTtInfo(DB)
		if !ok || tt.(string) != "LOAD_ONLY" {
			TtData.PrepareCoreData()
		}
	}
	return true
}

func SaveFile(filePath string, data []byte) bool {
	f, err := os.Create(filePath)
	if err == nil {
		defer f.Close()
		if _, err = f.Write(data); err == nil {
			return true
		}
	}
	base.Report(`<Error>Couldn't write output to: %s\n  +++ %v>`,
		filePath, err)
	return false
}
