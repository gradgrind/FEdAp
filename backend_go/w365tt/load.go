package w365tt

import (
	"gradgrind/backend/base"
	"gradgrind/backend/ttbase"
	"path/filepath"
	"strings"
)

func LoadFile(path string) {
	abspath, err := filepath.Abs(path)
	if err != nil {
		base.Report(`<Error>Couldn't resolve file path: %s>`, path)
		return
	}

	stempath := strings.TrimSuffix(abspath, filepath.Ext(abspath))
	stempath = strings.TrimSuffix(stempath, "_w365")

	db := base.NewDb()
	LoadJSON(db, abspath)
	db.PrepareDb()
	ttinfo := ttbase.MakeTtInfo(db)
	ttinfo.PrepareCoreData()
}
