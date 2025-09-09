package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
	"strings"

	"github.com/xuri/excelize/v2"
)

const (
	ALL_TEACHERS string = "All teachers"
)

func TeachersActivities(
	fet *readfet.Fet,
	activities []*ActivityData,
	stemfile string,
) {
	f := excelize.NewFile()
	err := f.SetSheetName("Sheet1", ALL_TEACHERS)
	if err != nil {
		panic(err)
	}
	col := 2
	for _, d := range fet.Days_List.Day {
		cr, err := excelize.CoordinatesToCellName(col, 1)
		if err != nil {
			panic(err)
		}
		f.SetCellStr(ALL_TEACHERS, cr, d.Name)
		for _, h := range fet.Hours_List.Hour {
			cr, err := excelize.CoordinatesToCellName(col, 2)
			if err != nil {
				panic(err)
			}
			f.SetCellStr(ALL_TEACHERS, cr, h.Name)
			col++
		}
	}
	// Map teacher to index
	tmap := map[string]int{}
	for i, t := range fet.Teachers_List.Teacher {
		n := t.Name
		tmap[n] = i
		cr, err := excelize.CoordinatesToCellName(1, i+3)
		if err != nil {
			panic(err)
		}
		f.SetCellStr(ALL_TEACHERS, cr, n)
		f.NewSheet(n)
		cr, err = excelize.CoordinatesToCellName(1, 1)
		if err != nil {
			panic(err)
		}
		f.SetCellStr(n, cr, "hour\\day")
	}
	for _, adata := range activities {
		//sbj := adata.Subject
		slist := strings.Join(adata.Students, ",")
		for _, t := range adata.Teachers {
			tix, ok := tmap[t]
			if !ok {
				panic("Unknown teacher: " + t)
			}
			if adata.Time.Day < 0 {
				continue
			}
			l := adata.Duration
			row := tix + 3
			col := adata.Time.Day*len(fet.Hours_List.Hour) + adata.Time.Hour + 2

			r1 := adata.Time.Hour + 2
			{
				for i, h := range fet.Hours_List.Hour {
					cr, err := excelize.CoordinatesToCellName(1, i+2)
					if err != nil {
						panic(err)
					}
					f.SetCellStr(t, cr, h.Name)
				}
				for i, d := range fet.Days_List.Day {
					cr, err := excelize.CoordinatesToCellName(i+2, 1)
					if err != nil {
						panic(err)
					}
					f.SetCellStr(t, cr, d.Name)
				}
			}

			for {
				cr, err := excelize.CoordinatesToCellName(col, row)
				if err != nil {
					panic(fmt.Sprintf("Invalid time: %d.%d", adata.Time.Day, adata.Time.Hour))
				}
				f.SetCellStr(ALL_TEACHERS, cr, slist)
				cr1, err := excelize.CoordinatesToCellName(adata.Time.Day+2, r1)
				if err != nil {
					panic(err)
				}
				f.SetCellStr(t, cr1, slist)
				l--
				if l <= 0 {
					break
				}
				col++
				r1++
			}
		}
	}
	if err := f.SaveAs(stemfile + "_teachers.xlsx"); err != nil {
		panic(err)
	}
}
