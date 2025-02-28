package backend

import (
	"gradgrind/backend/base"
	"gradgrind/backend/ttbase"
	"os"
	"path/filepath"
)

func LoadFile(path string, loader func(*base.DbTopLevel, string)) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		base.Report(`<Error>Couldn't resolve file path: %s\n  +++ %v>`,
			path, err)
		return
	}

	//stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	//stempath = strings.TrimSuffix(stempath, "_w365")

	db := base.NewDb()
	loader(db, abspath)
	//LoadJSON(db, abspath)
	db.PrepareDb()
	DB = db

	ttinfo := ttbase.MakeTtInfo(db)

	//TODO: This was optional for inputting to print. That may be difficult
	// to handle outside of that dedicated app. How much of it is essential?
	ttinfo.PrepareCoreData()

	TtData = ttinfo
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
