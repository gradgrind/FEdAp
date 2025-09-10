package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func TeachersActivities(
	fet *readfet.Fet,
	activities []*ActivityData,
	stemfile string,
) {
	nhours := len(fet.Hours_List.Hour)
	f := excelize.NewFile()
	overview_headers(fet, f, ALL_TEACHERS)
	// Map teacher to index
	tmap := map[string]int{}
	// Start rows of the detail tables
	row0_students := 1
	row0_subjects := row0_students + nhours + 1 + PERSONAL_TABLES_GAP
	row0_rooms := row0_subjects + nhours + 1 + PERSONAL_TABLES_GAP
	for i, t := range fet.Teachers_List.Teacher {
		n := t.Name
		tmap[n] = i
		// Teacher row header in ALL_TEACHERS sheet
		cr, err := excelize.CoordinatesToCellName(1, i+3)
		if err != nil {
			panic(err)
		}
		f.SetCellStr(ALL_TEACHERS, cr, n)
		// Add personal sheet for teacher
		f.NewSheet(n)
		// Add headers for student group table
		week_headers(fet, f, n, row0_students)
		week_headers(fet, f, n, row0_subjects)
		week_headers(fet, f, n, row0_rooms)
	}
	// Get the data from the activities
	for _, adata := range activities {
		sbj := adata.Subject
		slist := strings.Join(adata.Students, ",")
		rlist := strings.Join(adata.Rooms, ",")
		for _, t := range adata.Teachers {
			tix, ok := tmap[t]
			if !ok {
				panic("Unknown teacher: " + t)
			}
			if adata.Time.Day < 0 {
				continue
			}
			l := adata.Duration

			// Coordinates in ALL_TEACHERS sheet
			row := tix + 3
			col := adata.Time.Day*nhours + adata.Time.Hour + 2

			r1 := row0_students + adata.Time.Hour + 1
			r2 := row0_subjects + adata.Time.Hour + 1
			r3 := row0_rooms + adata.Time.Hour + 1

			for { // for each hour in duration
				// ALL_TEACHERS sheet
				cr, err := excelize.CoordinatesToCellName(col, row)
				if err != nil {
					panic(fmt.Sprintf("Invalid time: %d.%d", adata.Time.Day, adata.Time.Hour))
				}
				f.SetCellStr(ALL_TEACHERS, cr, slist)
				// Individual teacher's sheet
				//  - students
				cr1, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r1)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(t, cr1, slist)
				//  - subjects
				cr2, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r2)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(t, cr2, sbj)
				//  - rooms
				cr3, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r3)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(t, cr3, rlist)
				l--
				if l <= 0 {
					break
				}
				col++
				r1++
				r2++
				r3++
			}
		}
	}
	if err := f.SaveAs(stemfile + "_teachers.xlsx"); err != nil {
		panic(err)
	}
}
