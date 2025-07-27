package fet

//TODO: The stuff in files setup_... is just temporary, experimental.
// If it turns out to be useful, it should be reworked to use the base DB structures
// as source rather than FET files. It might be possible to read FET files into the
// the DB structures, to some extent?

import (
	"fedap/base"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type ActivityIndex int16
type ResourceIndex = int
type TimeSlot int16

type TtRoom struct {
	Id       Ref
	Tag      string
	Resource ResourceIndex
}

type TtVirtualRoom struct {
	RoomIndexes       []ResourceIndex
	RoomChoiceIndexes [][]ResourceIndex
}

type TtData struct {
	NDays        int
	NHours       int
	HoursPerWeek int
	Resources    []any
	DayIndex     map[string]int
	HourIndex    map[string]int
	TeacherIndex map[string]ResourceIndex
	GroupIndexes map[string][]ResourceIndex
	RoomIndex    map[string]ResourceIndex
	VirtualRooms map[string]TtVirtualRoom
	Activities   []TtActivity

	ActivitySlots []TimeSlot
	ResourceWeeks []ActivityIndex

	basic_activity_groups map[int]*BasicActivityGroup
	CollectedBags         map[*BasicActivityGroup]*BagCollection
}

// TODO--
func (tt_data *TtData) PrintBags() {
	for i, bag := range tt_data.basic_activity_groups {
		fmt.Printf("BAG %d:\n", i)
		fmt.Printf(" -- %v\n", bag)
	}
}

/* TODO: A "placement state" is essentially defined by the placements of
 * all the activities. As such a simple vector of these placements would
 * suffice. However, also the blocked time-slots need to be considered.
 * Further, there may well be room-choices. These could be stored in the
 * activity nodes themselves.
 * An alternative would be to use the resource-weeks vector, which should
 * contain all the placement information, but is likely to be much larger.
 * To minimize the save/restore times it could be sensible to prefer simple
 * copying over regeneration, so it might be worth saving both vectors and,
 * if possible, not saving changeable stuff in the activities.
 */

// Collect basic-possible-slots for the BasicActivityGroups.
func (tt_data *TtData) BasicBagSlots() {
	slots_per_week := TimeSlot(tt_data.NDays * tt_data.NHours)
	tt_data.HoursPerWeek = int(slots_per_week)
	//bagmap := map[*BasicActivityGroup]*[]*BasicActivityGroup{}
	for _, bag := range tt_data.basic_activity_groups {
		// Find all slots which are in principle possible for the activities
		// in this BAG.

		// This is after the placement of the fixed activities, so some BAGs may
		// not have any activities which still need placing.
		// However, they might still be needed for relative constraints!
		aixs := []int{}
		for _, aix := range bag.Activities {
			a := tt_data.Activities[aix]
			if !a.Fixed {
				aixs = append(aixs, aix)
			}
		}

		if len(aixs) != 0 {
			aix := bag.Activities[0]
			slots := []TimeSlot{}
			for t := range slots_per_week {
				if tt_data.TestPlaceBasic(aix, t) {
					slots = append(slots, TimeSlot(t))
				}
			}
			bag.BasicSlots = SlotCombinations(slots, len(aixs))
		}
	}

	//TODO: What to do with bagmap?
	/* TODO: call baglist_placements() instead of doing the following
	for _, baglist := range bagmap {
		placements := [][]TimeSlot{}
		stem := [][]TimeSlot{}
		for _, bag := range *baglist {
			for _, slots := range bag.BasicSlots {
				state := tt_data.SaveState()
				for i, aix := range bag.Activities {
					slot := slots[i]
					// Test and place the activity in the slot, also checking
					// days-between, etc. and accumulating penalties
					if tt_data.TestPlaceBasic(aix, slot) {
						tt_data.PlaceBasic(aix, slot)
					} else {
						// This combination failed
						goto next
					}
				}

			next:
				tt_data.RestoreStateMove(state)
			}
		}

		//TODO: Assign placements to baglist
	} */

	/*
		done := map[*BasicActivityGroup]bool{}

		sortFunc := func(a, b *BasicActivityGroup) int {
			return cmp.Compare(a.Activities[0], b.Activities[0])
		}
		for _, bag := range slices.SortedFunc(maps.Keys(bagmap), sortFunc) {
			fmt.Printf("\n$ BasicSlots: %+v\n", bag.BasicSlots)

			if done[bag] {
				continue
			}
			baglist := bagmap[bag]
			fmt.Printf("\n+++ %+v: %d\n  --", bag.Activities, len(*baglist))
			for _, bag1 := range *baglist {
				fmt.Printf(" %+v", bag1.Activities)
				done[bag1] = true
			}
		}
	*/
}

func (tt_data *TtData) baglist_placements(
	stem []TimeSlot,
	baglist *[]*BasicActivityGroup,
	bagli int,
) [][]TimeSlot {
	if bagli == len(*baglist) {
		//TODO
	}

	placements := [][]TimeSlot{}

	bag := (*baglist)[bagli]

	for _, slots := range bag.BasicSlots {
		state := tt_data.SaveState()
		for i, aix := range bag.Activities {
			slot := slots[i]
			// Test and place the activity in the slot
			if tt_data.TestPlaceBasic(aix, slot) {
				tt_data.PlaceBasic(aix, slot)
			} else {
				// This combination failed
				goto next
			}
		}
		// All placements OK, goto next bag
		placements = append(placements, tt_data.baglist_placements(
			append(append([]TimeSlot{}, stem...), slots...),
			baglist,
			bagli+1,
		)...)

	next:
		tt_data.RestoreStateMove(state)
	}

	return placements

	//TODO: also check days-between, etc. and accumulate penalties
}

type TtTeacher struct {
	Id       Ref
	Tag      string
	Resource ResourceIndex
}

type TtAtomicGroup struct {
	//Id            Ref
	Tag      string
	Resource ResourceIndex
}

func (tt_data *TtData) TimeSlotIndex(day string, hour string) int {
	d := tt_data.DayIndex[day]
	h := tt_data.HourIndex[hour]
	return d*tt_data.NHours + h
}

func PrepareResources(fetdata *fet) *TtData {
	tt_data := &TtData{
		NDays:                 fetdata.Days_List.Number_of_Days,
		NHours:                fetdata.Hours_List.Number_of_Hours,
		DayIndex:              map[string]int{},
		HourIndex:             map[string]int{},
		TeacherIndex:          map[string]ResourceIndex{},
		GroupIndexes:          map[string][]ResourceIndex{},
		RoomIndex:             map[string]ResourceIndex{},
		VirtualRooms:          map[string]TtVirtualRoom{},
		basic_activity_groups: map[int]*BasicActivityGroup{},
		CollectedBags:         map[*BasicActivityGroup]*BagCollection{},
	}

	for i, d := range fetdata.Days_List.Day {
		tt_data.DayIndex[d.Name] = i
	}
	for i, h := range fetdata.Hours_List.Hour {
		tt_data.HourIndex[h.Name] = i
	}
	slots_per_week := tt_data.NDays * tt_data.NHours
	tt_data.HoursPerWeek = slots_per_week
	fmt.Printf("n days = %d, n hours = %d\n", tt_data.NDays, tt_data.NHours)

	for _, t := range fetdata.Teachers_List.Teacher {
		tid := base.NewId()
		tix := len(tt_data.Resources)
		item_p := &TtTeacher{tid, t.Name, tix}
		tt_data.Resources = append(tt_data.Resources, item_p)
		tt_data.TeacherIndex[t.Name] = tix
		tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
			make([]ActivityIndex, slots_per_week)...)
	}
	fmt.Printf("n teachers: %d\n", len(tt_data.TeacherIndex))

	for _, c := range fetdata.Students_List.Year {
		sep := c.Separator
		agsep := "#"
		if sep == "#" {
			agsep = "$"
		}
		//-- fmt.Printf("§ %d - %#v\n", len(c.Group), c)

		// If there are no groups, the class is the atomic-group (with
		// a slightly changed name)
		// If a group has subgroups, these are the atomic-groups.
		// Otherwise the group will be the subgroup (with
		// a slightly changed name)

		if len(c.Group) == 0 {
			// Make new atomic-group
			atag := c.Name + agsep
			tix := len(tt_data.Resources)
			item_p := TtAtomicGroup{atag, tix}
			tt_data.GroupIndexes[c.Name] = []int{tix}
			tt_data.Resources = append(tt_data.Resources, item_p)
			tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
				make([]ActivityIndex, slots_per_week)...)
			fmt.Printf("§ %s - %d\n", c.Name, tix)

		} else {
			agmap := map[string]int{}

			for _, g := range c.Group {
				if len(g.Subgroup) == 0 {
					// Make new atomic-group
					atag := strings.ReplaceAll(g.Name, sep, agsep)
					tix := len(tt_data.Resources)
					item_p := TtAtomicGroup{atag, tix}
					tt_data.GroupIndexes[g.Name] = []int{tix}
					tt_data.Resources = append(tt_data.Resources, item_p)
					tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
						make([]ActivityIndex, slots_per_week)...)
					fmt.Printf("§ %s - %d\n", g.Name, tix)

				} else {
					// Use the existing subgroups as atomic-groups
					aglist := []int{}
					for _, ag := range g.Subgroup {
						tix, ok := agmap[ag.Name]
						if !ok {
							tix = len(tt_data.Resources)
							agmap[ag.Name] = tix
							item_p := TtAtomicGroup{ag.Name, tix}
							tt_data.Resources = append(tt_data.Resources, item_p)
							tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
								make([]ActivityIndex, slots_per_week)...)
						}
						aglist = append(aglist, tix)
					}
					tt_data.GroupIndexes[g.Name] = aglist
					fmt.Printf("§ %s - %v\n", g.Name, aglist)

				}
			}
			// Now get atomic-groups for whole class
			aglist := slices.Sorted(maps.Values(agmap))
			tt_data.GroupIndexes[c.Name] = aglist
			fmt.Printf("§§ %s - %v\n", c.Name, aglist)
		}
	}

	// Teachers and groups are referenced directly in the FET Activity nodes.
	// Rooms are not, they are attached to the activities by constraints. The
	// constraints specify a list of rooms from which one is to be chosen.
	// To complicate matters, there are – in addition to the "real" rooms –
	// "virtual" rooms. If the choice list contains just one room, it may be real
	// or virtual, otherwise only real rooms are acceptable.
	// A virtual room contains a list of needed real rooms and a list of room
	// choices, whose components must also be real rooms.
	// FET may well support more complex arrangements which are currently not
	// supported by this scheme.

	for _, r := range fetdata.Rooms_List.Room {
		//TODO: The handling of virtual rooms may need delaying until all the
		// real rooms have been read.
		if r.Virtual {
			// This is basically a list of room choices, where each "choice"
			// may actually be a single room, i.e. no choice!
			item_p := TtVirtualRoom{}
			for _, rset := range r.Set_of_Real_Rooms {
				if rset.Number_of_Real_Rooms == 1 {
					r := rset.Real_Room[0]
					tix, ok := tt_data.RoomIndex[r]
					if ok {
						item_p.RoomIndexes = append(item_p.RoomIndexes, tix)
					} else {
						base.Error.Fatalf("Unknown (real) room: %s", r)
					}
				} else {
					rooms := []int{}
					for _, r := range rset.Real_Room {
						tix, ok := tt_data.RoomIndex[r]
						if ok {
							rooms = append(rooms, tix)
						} else {
							base.Error.Fatalf("Unknown (real) room: %s", r)
						}
					}
					item_p.RoomChoiceIndexes = append(item_p.RoomChoiceIndexes,
						rooms)
				}
			}
			tt_data.VirtualRooms[r.Name] = item_p
			fmt.Printf("ROOM: %s – %v\n", r.Name, item_p)

		} else {
			tid := base.NewId()
			tix := len(tt_data.Resources)
			item_p := &TtRoom{tid, r.Name, tix}
			tt_data.Resources = append(tt_data.Resources, item_p)
			tt_data.RoomIndex[r.Name] = tix
			tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
				make([]ActivityIndex, slots_per_week)...)
			fmt.Printf("ROOM: %s – %d\n", r.Name, tix)
		}

	}
	return tt_data
}

func SlotCombinations(iterable []TimeSlot, r int) (rt [][]TimeSlot) {
	pool := iterable
	n := len(pool)

	if r > n {
		return
	}

	indices := make([]int, r)
	for i := range indices {
		indices[i] = i
	}

	result := make([]TimeSlot, r)
	for i, el := range indices {
		result[i] = pool[el]
	}
	s2 := make([]TimeSlot, r)
	copy(s2, result)
	rt = append(rt, s2)

	for {
		i := r - 1
		for ; i >= 0 && indices[i] == i+n-r; i-- {
		}

		if i < 0 {
			return
		}

		indices[i] += 1
		for j := i + 1; j < r; j += 1 {
			indices[j] = indices[j-1] + 1
		}

		for ; i < len(indices); i += 1 {
			result[i] = pool[indices[i]]
		}
		s2 = make([]TimeSlot, r)
		copy(s2, result)
		rt = append(rt, s2)
	}
}
