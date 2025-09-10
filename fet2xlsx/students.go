package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

func StudentsActivities(
	fet *readfet.Fet,
	activities []*ActivityData,
	stemfile string,
) {
	nhours := len(fet.Hours_List.Hour)
	f := excelize.NewFile()
	overview_headers(fet, f, ALL_STUDENTS)
	// Map student group to index
	gmap := map[string]int{}
	// Start rows of the detail tables
	row0_subjects := 1
	row0_teachers := row0_subjects + nhours + 1 + PERSONAL_TABLES_GAP
	row0_rooms := row0_teachers + nhours + 1 + PERSONAL_TABLES_GAP
	// Show student "Years" and "Groups"
	glist := []string{}
	i := 0
	for _, y := range fet.Students_List.Year {
		glist = append(glist, y.Name)
		for _, g := range y.Group {
			glist = append(glist, g.Name)
		}
	}
	for _, n := range glist {
		gmap[n] = i
		// Group row header in ALL_STUDENTS sheet
		cr, err := excelize.CoordinatesToCellName(1, i+3)
		if err != nil {
			panic(err)
		}
		i++
		f.SetCellStr(ALL_STUDENTS, cr, n)
		// Add personal sheet for group
		f.NewSheet(n)
		// Add headers for student group table
		week_headers(fet, f, n, row0_subjects)
		week_headers(fet, f, n, row0_teachers)
		week_headers(fet, f, n, row0_rooms)
	}
	// Get the data from the activities
	for _, adata := range activities {
		sbj := adata.Subject
		tlist := strings.Join(adata.Teachers, ",")
		rlist := strings.Join(adata.Rooms, ",")
		for _, g := range adata.Students {
			gix, ok := gmap[g]
			if !ok {
				panic("Unknown student group: " + g)
			}
			if adata.Time.Day < 0 {
				continue
			}
			l := adata.Duration

			// Coordinates in ALL_STUDENTS sheet
			row := gix + 3
			col := adata.Time.Day*nhours + adata.Time.Hour + 2

			r1 := row0_subjects + adata.Time.Hour + 1
			r2 := row0_teachers + adata.Time.Hour + 1
			r3 := row0_rooms + adata.Time.Hour + 1

			for { // for each hour in duration
				// ALL_STUDENTS sheet
				cr, err := excelize.CoordinatesToCellName(col, row)
				if err != nil {
					panic(fmt.Sprintf("Invalid time: %d.%d", adata.Time.Day, adata.Time.Hour))
				}
				f.SetCellStr(ALL_STUDENTS, cr, sbj)
				// Individual group's sheet
				//  - subjects
				cr2, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r1)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(g, cr2, sbj)
				//  - teachers
				cr1, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r2)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(g, cr1, tlist)
				//  - rooms
				cr3, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r3)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(g, cr3, rlist)
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
	if err := f.SaveAs(stemfile + "_students.xlsx"); err != nil {
		panic(err)
	}
}
