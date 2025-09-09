package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
)

func TeachersActivities(fet *readfet.Fet) {
	for _, act := range fet.Activities_List.Activity {
		fmt.Printf(" -- %03d: %v\n", act.Id, act.Teacher)
	}
}
