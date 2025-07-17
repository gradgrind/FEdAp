package fet

import (
	"fedap/base"
	"fmt"
	"slices"
	"strconv"
)

/* It looks like Activity_Group_Ids won't play any significant role here.
 * They could be used to signal shared teachers and groups, but not rooms.
 * That means that the placement tests will need to be done in full for each
 * activity, which will lose some efficiency. On the other hand, activities
 * with duration > 1 can be handled fairly transparently.
 */

func (tt_data *TtData) Day(t TimeSlot) int {
	return int(t) / tt_data.NHours
}

func (tt_data *TtData) Hour(t TimeSlot) int {
	return int(t) % tt_data.NHours
}

//type PreferredTime struct {
//	Slot    TimeSlot
//	Penalty int
//}

type TtActivity struct {
	Id       int
	Duration int
	//Placement TimeSlot
	//PlacementPenalty int
	BasicActivityGroup     *BasicActivityGroup
	Resources              []int
	RoomChoices            [][]int
	DaysBetweenConstraints []ConstraintDaysBetween
	//PreferredTimes         []PreferredTime
	Fixed bool
}

//func (a *TtActivity) setTime(t int) {
//	a.Placement = TimeSlot(t)
//}

type BasicActivityGroup struct {
	Activities           []int
	BasicSlots           []TimeSlot
	LinkedActivityGroups []*BasicActivityGroup
	Duration             int
	Processed            bool
}

// Although FET activities are numbered contiguously, there may be, effectively,
// gaps caused by inactive activities. To keep things fairly transparent, however,
// the FET Id numbers are retained, padding the activities slice with empty
// entries.
func (tt_data *TtData) SetupActivities(fetdata *fet) {
	// Find highest activity index
	aix := 0
	for _, a := range fetdata.Activities_List.Activity {
		if a.Id > aix {
			aix = a.Id
		}
	}

	//TODO: THe BAGs should be set up AFTER the fixed activities, as the fixed
	// activities should not be in BAGs!

	// Set up initial activity groups (from FET activity groups), these are the
	// BasicActivityGroups. They group activities with the same resources and
	// durations. This can possibly be used to reduce the number of placement
	// options. They are used later when building groups of activities to be
	// placed simultaneously.

	activity_groups := map[int][]*BasicActivityGroup{} // key is activity group id
	// This maps the FET activity group to the basic local activity group
	activity_bags := map[int]*BasicActivityGroup{} // key is activity id

	tt_data.Activities = make([]TtActivity, aix+1)
	for _, a := range fetdata.Activities_List.Activity {
		if !a.Active {
			continue
		}

		if a.Activity_Group_Id == 0 {
			activity_bags[a.Id] = &BasicActivityGroup{
				Activities: []int{a.Id},
				//BasicSlots
				//LinkedActivityGroups
				Duration:  a.Duration,
				Processed: false,
			}
			continue
		}
		var bag *BasicActivityGroup
		if a.Activity_Group_Id == a.Id {
			bag = &BasicActivityGroup{
				Activities: []int{a.Id},
				//BasicSlots
				//LinkedActivityGroups
				Duration:  a.Duration,
				Processed: false,
			}
			activity_groups[a.Id] = []*BasicActivityGroup{bag}
		} else {
			// Split FET activity group according to duration
			i := 0
			for _, bag = range activity_groups[a.Activity_Group_Id] {
				if bag.Duration == a.Duration {
					// Add to existing entry
					bag.Activities = append(bag.Activities, a.Id)
					goto done1
				}
				if bag.Duration < a.Duration {
					// New entry before current one (longer duration first)
					break
				}
				i++
			}
			// Insert or append new entry
			bag = &BasicActivityGroup{
				Activities: []int{a.Id},
				//BasicSlots
				//LinkedActivityGroups
				Duration:  a.Duration,
				Processed: false,
			}
			activity_groups[a.Activity_Group_Id] = slices.Insert(
				activity_groups[a.Activity_Group_Id], i, bag)
		done1:
		}
		activity_bags[a.Id] = bag

		activity := TtActivity{Id: a.Id, Duration: a.Duration}
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

	tt_data.basic_activity_groups = activity_bags

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
		// NOTE that the FET constraint allows weights of less than 100%, but here
		// a hard constraint is assumed in all cases.
		t := TimeSlot(tt_data.TimeSlotIndex(rc.Preferred_Day, rc.Preferred_Hour))
		if !tt_data.TestPlaceBasic(rc.Activity_Id, t) {
			base.Error.Fatalf(
				"Fixed activity %d cannot be placed in time-slot %d\n",
				rc.Activity_Id, t)
		}
		tt_data.PlaceBasic(rc.Activity_Id, t)
		a := &tt_data.Activities[rc.Activity_Id]
		a.Fixed = true
	}
}

// This asssumes that the target has been tested, so that the placement is valid.
func (tt_data *TtData) PlaceBasic(activity int, t TimeSlot) {
	a := tt_data.Activities[activity]
	rmap := tt_data.ResourceWeeks
	for range a.Duration {
		for _, r := range a.Resources {
			rmap[r][t] = ActivityIndex(activity)
		}
		t++
	}
}

func (tt_data *TtData) TestPlaceBasic(activity int, t TimeSlot) bool {
	a := tt_data.Activities[activity]

	//TODO: Note that this part of the test won't be necessary for pre-checked
	// activities.
	if a.Duration > 1 {
		// Check that it fits in a day
		if a.Duration+tt_data.Hour(t) > tt_data.NHours {
			return false
		}
	}

	// Check that all the slots are free for all the resources
	rmap := tt_data.ResourceWeeks
	for range a.Duration {
		for _, r := range a.Resources {
			if rmap[r][t] != 0 {
				return false
			}
		}
		t++
	}
	return true
}

type ConstraintDaysBetween struct {
	Activity           ActivityIndex
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
							ActivityIndex(aix),
							penalty, mindays,
							rc.Consecutive_If_Same_Day})
				}
			}
		}
	}
}
