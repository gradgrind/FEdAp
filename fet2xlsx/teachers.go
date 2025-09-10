package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

var (
	Teachers     []string       // years and their groups
	TeacherIndex map[string]int // map group to index
)

func GetTeachers(fet *readfet.Fet) {
	TeacherIndex = map[string]int{}
	for i, t := range fet.Teachers_List.Teacher {
		Teachers = append(Teachers, t.Name)
		TeacherIndex[t.Name] = i
	}
}

func TeachersActivities(
	fet *readfet.Fet,
	activities []*ActivityData,
	stemfile string,
) (string, error) {
	nhours := len(fet.Hours_List.Hour)
	f := excelize.NewFile()
	if err := overview_headers(fet, f, ALL_TEACHERS); err != nil {
		return "", err
	}
	// Start rows of the detail tables
	row0_students := 1
	row0_subjects := row0_students + nhours + 1 + PERSONAL_TABLES_GAP
	row0_rooms := row0_subjects + nhours + 1 + PERSONAL_TABLES_GAP
	for i, n := range Teachers {
		// Teacher row header in ALL_TEACHERS sheet
		cr, err := excelize.CoordinatesToCellName(1, i+3)
		if err != nil {
			return "", err
		}
		f.SetCellStr(ALL_TEACHERS, cr, n)
		// Add personal sheet for teacher
		f.NewSheet(n)
		// Add headers for student group table
		if err := week_headers(fet, f, n, row0_students); err != nil {
			return "", err
		}
		if err := week_headers(fet, f, n, row0_subjects); err != nil {
			return "", err
		}
		if err := week_headers(fet, f, n, row0_rooms); err != nil {
			return "", err
		}
	}
	// Get the data from the activities
	for _, adata := range activities {
		sbj := adata.Subject
		slist := strings.Join(adata.Students, ",")
		rlist := strings.Join(adata.Rooms, ",")
		for _, t := range adata.Teachers {
			tix, ok := TeacherIndex[t]
			if !ok {
				return "", fmt.Errorf("unknown teacher: %s", t)
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
					return "", fmt.Errorf(
						"invalid time: %d.%d", adata.Time.Day, adata.Time.Hour)
				}
				f.SetCellStr(ALL_TEACHERS, cr, slist)
				// Individual teacher's sheet
				//  - students
				cr1, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r1)
				if err != nil {
					return "", err
				}
				f.SetCellStr(t, cr1, slist)
				//  - subjects
				cr2, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r2)
				if err != nil {
					return "", err
				}
				f.SetCellStr(t, cr2, sbj)
				//  - rooms
				cr3, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r3)
				if err != nil {
					return "", err
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
	opath := stemfile + "_teachers.xlsx"
	if err := f.SaveAs(opath); err != nil {
		return "", err
	}
	return opath, nil
}
