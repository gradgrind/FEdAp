package fet

import (
	"fedap/base"
	"fmt"
)

/* It looks like Activity_Group_Ids won't play any significant role here.
 * They could be used to signal shared teachers and groups, but not rooms.
 * That means that the placement tests will need to be done in full for each
 * activity, which will lose some efficiency. On the other hand, activities
 * with duration > 1 can be handled fairly transparently.
 */

type TimeSlot int

type TtActivity struct {
	Id          int
	Placement   TimeSlot
	Fixed       bool
	Resources   []int
	RoomChoices [][]int
}

func (tt_data *TtData) SetupActivities(fetdata *fet) {
	tt_data.MapActivity = map[int]int{}
	for _, a := range fetdata.Activities_List.Activity {
		if !a.Active {
			continue
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

		tt_data.MapActivity[a.Id] = len(tt_data.Activities)
		tt_data.Activities = append(tt_data.Activities, activity)
	}

	//TODO: Need (at least) fixed rooms

	for _, rc := range fetdata.Space_Constraints_List.
		ConstraintActivityPreferredRooms {

		if !rc.Active {
			continue
		}
		if rc.Weight_Percentage == 100 && rc.Number_of_Preferred_Rooms == 1 {
			aix, ok := tt_data.MapActivity[rc.Activity_Id]
			if !ok {
				base.Error.Fatalf("Unknown activity id: %d", rc.Activity_Id)
			}
			a := tt_data.Activities[aix]
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
			aix, ok := tt_data.MapActivity[rc.Activity_Id]
			if !ok {
				base.Error.Fatalf("Unknown activity id: %d", rc.Activity_Id)
			}
			a := tt_data.Activities[aix]
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

//TODO: Blocked time-slots, different-days (and other constraints?)
