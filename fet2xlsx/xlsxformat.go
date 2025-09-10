package fet2xlsx

import (
	"fedap/readfet"

	"github.com/xuri/excelize/v2"
)

const (
	ALL_TEACHERS        string = "All teachers"
	ALL_STUDENTS        string = "All students"
	PERSONAL_TABLES_GAP int    = 2
)

func overview_headers(fet *readfet.Fet, f *excelize.File, sheet string) error {
	err := f.SetSheetName("Sheet1", sheet)
	if err != nil {
		return err
	}
	// Days and hours headers, use the "Name" field
	days := fet.Days_List.Day
	hours := fet.Hours_List.Hour
	col := 2
	for _, d := range days {
		cr, err := excelize.CoordinatesToCellName(col, 1)
		if err != nil {
			return err
		}
		f.SetCellStr(sheet, cr, d.Name)
		for _, h := range hours {
			cr, err := excelize.CoordinatesToCellName(col, 2)
			if err != nil {
				return err
			}
			f.SetCellStr(sheet, cr, h.Name)
			col++
		}
	}
	return nil
}

func week_headers(
	fet *readfet.Fet,
	f *excelize.File,
	sheet string,
	r0 int,
) error {
	days := fet.Days_List.Day
	hours := fet.Hours_List.Hour
	cr, err := excelize.CoordinatesToCellName(1, r0)
	if err != nil {
		return err
	}
	f.SetCellStr(sheet, cr, "hour\\day")
	for i, d := range days {
		cr, err := excelize.CoordinatesToCellName(i+2, r0)
		if err != nil {
			return err
		}
		f.SetCellStr(sheet, cr, d.Name)
	}
	for _, h := range hours {
		r0++
		cr, err := excelize.CoordinatesToCellName(1, r0)
		if err != nil {
			return err
		}
		f.SetCellStr(sheet, cr, h.Name)
	}
	return nil
}
