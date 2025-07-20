package fet

import "fedap/base"

func (tt_data *TtData) ResourceBlocking(fetdata *fet) {
	slots_per_week := tt_data.HoursPerWeek
	for _, rc := range fetdata.Time_Constraints_List.
		ConstraintTeacherNotAvailableTimes {

		tt := rc.Teacher
		t, ok := tt_data.TeacherIndex[tt]
		if !ok {
			base.Error.Fatalf("Unknown teacher: %s", tt)
		}
		for _, dh := range rc.Not_Available_Time {
			tt_data.ResourceWeeks[t*slots_per_week+tt_data.TimeSlotIndex(dh.Day, dh.Hour)] = -1
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
				tt_data.ResourceWeeks[g*slots_per_week+slot] = -1
			}
		}
	}

	//TODO? room blocking not currently supported
}
