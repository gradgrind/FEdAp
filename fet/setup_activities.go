package fet

import (
	"fedap/base"
	"fmt"
	"strconv"
)

/* It looks like Activity_Group_Ids won't play any significant role here.
 * They could be used to signal shared teachers and groups, but not rooms.
 * That means that the placement tests will need to be done in full for each
 * activity, which will lose some efficiency. On the other hand, activities
 * with duration > 1 can be handled fairly transparently.
 */

type TimeSlot int

type PreferredTime struct {
	Slot    TimeSlot
	Penalty int
}

type TtActivity struct {
	Id                     int
	Placement              TimeSlot //TODO?
	Fixed                  bool     //TODO?
	Resources              []int
	RoomChoices            [][]int
	DaysBetweenConstraints []ConstraintDaysBetween
	PreferredTimes         []PreferredTime
}

func (a *TtActivity) setTime(t int) {
	a.Placement = TimeSlot(t)
}

// Although FET activities are numbered contiguously from 1, a distinct indexing
// is used here, which means that a translation, MapActivity, is needed. This
// allows skipping of non-active activities, etc.
// TODO: This maybe a bit short-sighted ... wouldn't it be better to use the FET
// numbering, even with gaps?
func (tt_data *TtData) SetupActivities(fetdata *fet) {
	// Find highest activity index
	aix := 0
	for _, a := range fetdata.Activities_List.Activity {
		if a.Id > aix {
			aix = a.Id
		}
	}

	// Set up initial activity groups (from FET activity groups)
	activity_groups := map[int][]int{}

	tt_data.Activities = make([]TtActivity, aix+1)
	for _, a := range fetdata.Activities_List.Activity {
		if !a.Active {
			continue
		}

		if a.Activity_Group_Id == 0 {
			activity_groups[a.Id] = []int{a.Id}
		} else {
			activity_groups[a.Activity_Group_Id] = append(
				activity_groups[a.Activity_Group_Id], a.Id)
		}

		activity := TtActivity{Id: a.Id}
		for _, t := range a.Teacher {
			tix, ok := tt_data.TeacherIndex[t]
			if !ok {
				base.Error.Fatalf("Unknown teacher: %s", t)
			}
			activity.Resources = append(activity.Resources, tix)
		}
		for _, g := range a.Students {
			tixs, ok := tt_data.GroupIndexes[g]
			if !ok {
				base.Error.Fatalf("Unknown group: %s", g)
			}
			activity.Resources = append(activity.Resources, tixs...)
		}

		tt_data.Activities[a.Id] = activity
	}

	tt_data.fet_activity_groups = activity_groups

	// Need (at least) fixed rooms

	for _, rc := range fetdata.Space_Constraints_List.
		ConstraintActivityPreferredRooms {

		if !rc.Active {
			continue
		}
		if rc.Weight_Percentage == 100 && rc.Number_of_Preferred_Rooms == 1 {
			a := tt_data.Activities[rc.Activity_Id]
			rp := rc.Preferred_Room[0]
			r, ok := tt_data.RoomIndex[rp]
			if ok {
				a.Resources = append(a.Resources, r)
				fmt.Printf("Activity fixed room, %d: %v\n", rc.Activity_Id, r)
			} else {
				vr, ok := tt_data.VirtualRooms[rp]
				if !ok {
					base.Error.Fatalf("Unknown room): %s", rp)
				}
				a.Resources = append(a.Resources, vr.RoomIndexes...)
				fmt.Printf("Activity fixed room, %d: %v\n", rc.Activity_Id, vr.RoomIndexes)

				//TODO? RoomChoiceIndexes?
			}
		} // else: TODO? This is a choice list or soft constraint

	}
	for _, rc := range fetdata.Space_Constraints_List.
		ConstraintActivityPreferredRoom {

		if !rc.Active {
			continue
		}
		if rc.Weight_Percentage == 100 {
			a := tt_data.Activities[rc.Activity_Id]
			if len(rc.Real_Room) == 0 {
				// A single real room
				r, ok := tt_data.RoomIndex[rc.Room]
				if !ok {
					base.Error.Fatalf("Unknown (real) room): %s", rc.Room)
				}
				a.Resources = append(a.Resources, r)
				fmt.Printf("Activity fixed room, %d: %v\n", rc.Activity_Id, r)
			} else {
				// Possibly more than one real room
				for _, rr := range rc.Real_Room {
					r, ok := tt_data.RoomIndex[rr]
					if !ok {
						base.Error.Fatalf("Unknown (real) room): %s", rr)
					}
					a.Resources = append(a.Resources, r)
					fmt.Printf("Activity fixed room, %d: %v\n", rc.Activity_Id, r)
				}
			}
		} // else: TODO? This is a soft constraint
	}
}

func (tt_data *TtData) SetupFixedTimes(fetdata *fet) {
	for _, rc := range fetdata.Time_Constraints_List.
		ConstraintActivityPreferredStartingTime {
		if !rc.Active {
			continue
		}
		t := TimeSlot(tt_data.TimeSlotIndex(rc.Preferred_Day, rc.Preferred_Hour))
		a := &tt_data.Activities[rc.Activity_Id]
		wp, _ := strconv.ParseFloat(rc.Weight_Percentage, 64)
		w := 100.0 - wp
		if w >= 0.1 { // everything under 0.1% counts as "hard"
			penalty := int(100.0 / w)
			a.PreferredTimes = append(a.PreferredTimes,
				PreferredTime{t, penalty})
		} else {
			a.Placement = t
			//TODO?
			a.Fixed = true
		}
	}
}

type ConstraintDaysBetween struct {
	Activity           int
	Penalty            int
	MinDays            int
	SameDayConsecutive bool // only relevant if penalty > 0:
	// if true, it is a "hard" constraint
}

func (tt_data *TtData) SetupDaysBetween(fetdata *fet) {
	for _, rc := range fetdata.Time_Constraints_List.
		ConstraintMinDaysBetweenActivities {
		if !rc.Active {
			continue
		}
		penalty := -1
		wp, _ := strconv.ParseFloat(rc.Weight_Percentage, 64)
		w := 100.0 - wp
		if w >= 0.1 { // everything under 0.1% counts as "hard"
			penalty = int(100.0 / w)
		}
		mindays := rc.MinDays
		for _, aix0 := range rc.Activity_Id {
			a := &tt_data.Activities[aix0]
			for _, aix := range rc.Activity_Id {
				if aix != aix0 {
					a.DaysBetweenConstraints = append(a.DaysBetweenConstraints,
						ConstraintDaysBetween{
							aix, penalty, mindays, rc.Consecutive_If_Same_Day})
				}
			}
		}
	}
}
