package fet

import "fedap/base"

func (tt_data *TtData) ResourceBlocking(fetdata *fet) {
	for _, rc := range fetdata.Time_Constraints_List.
		ConstraintTeacherNotAvailableTimes {

		tt := rc.Teacher
		t, ok := tt_data.TeacherIndex[tt]
		if !ok {
			base.Error.Fatalf("Unknown teacher: %s", tt)
		}
		week := tt_data.ResourceWeeks[t]
		for _, dh := range rc.Not_Available_Time {
			week[tt_data.TimeSlotIndex(dh.Day, dh.Hour)] = -1
		}

	}

	for _, rc := range fetdata.Time_Constraints_List.
		ConstraintStudentsSetNotAvailableTimes {

		gg := rc.Students
		gx, ok := tt_data.GroupIndexes[gg]
		if !ok {
			base.Error.Fatalf("Unknown student group: %s", gg)
		}
		for _, dh := range rc.Not_Available_Time {
			slot := tt_data.TimeSlotIndex(dh.Day, dh.Hour)
			for _, g := range gx {
				week := tt_data.ResourceWeeks[g]
				week[slot] = -1
			}
		}
	}

	//TODO? room blocking not currently supported
}
