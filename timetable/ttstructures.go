// Package timetable deals with timetables. It uses its own data structures,
// which are designed around the need to handle constraints and avoid
// "collisions" (e.g. a teacher in two classes at the same time).
package timetable

import (
	"fedap/base"
	"slices"
)

type NodeRef = base.Ref // node reference (UUID)

type ActivityIndex int16
type ResourceIndex = int
type TimeSlot int16

type TtActivity struct {
	Id       ActivityIndex
	Duration int16
	Fixed    bool
	//BasicActivityGroup *BasicActivityGroup
	Resources   []ResourceIndex
	RoomChoices [][]ResourceIndex
	//TODO: Constraints to be applied when placing manually?
	//Constraints []TtConstraint
}

type TtRoom struct {
	Id       NodeRef
	Tag      string
	Resource ResourceIndex
}

type TtVirtualRoom struct {
	RoomIndexes       []ResourceIndex
	RoomChoiceIndexes [][]ResourceIndex
}

type Placement struct {
	Activity ActivityIndex
	Slot     TimeSlot
}

type TimetableUnit struct {
	Activities []ActivityIndex
	Placements [][]TimeSlot
	//TODO: Constraints to be applied after a TimetableUnit has been placed?
	//Constraints []TtConstraint
	Next int
}

// A TtData is the top-level structure for the timetable data.
type TtData struct {
	NDays         int
	NHours        int
	HoursPerWeek  int
	ActivitySlots []TimeSlot
	ResourceWeeks []ActivityIndex

	AtomicGroups   map[NodeRef][]*AtomicGroup
	ClassDivisions map[NodeRef][][]NodeRef

	// Set up by `GatherCourseInfo`
	Activities      []*TtActivity
	ActivityCourses []*TtCourseInfo
	CourseInfo      map[NodeRef]*TtCourseInfo // key is Course or SuperCourse

	//???
	Resources    []any
	DayIndex     map[string]int
	HourIndex    map[string]int
	TeacherIndex map[string]ResourceIndex
	GroupIndexes map[string][]ResourceIndex
	RoomIndex    map[string]ResourceIndex
	VirtualRooms map[string]TtVirtualRoom

	basic_activity_groups map[int]*BasicActivityGroup
	CollectedBags         map[*BasicActivityGroup]*BagCollection
}

// BasicSetup performs the initialization of a TtData structure, collecting
// "resources" (teachers, student groups and rooms) and "activities".
func BasicSetup(db *base.DbTopLevel) *TtData {
	days := len(db.Days)
	hours := len(db.Hours)

	course_info := CollectCourses(db)
	class_divisions := FilterDivisions(db, course_info)

	// Atomic groups: an atomic group is an ordered list of single groups,
	// one from each division.
	atomicGroups, atomicGroupIndex := MakeAtomicGroups(db, class_divisions)

	//TODO: What does this do? Could it be combined with MakeAtomicGroups?
	ags := []*AtomicGroup{}
	g2ags := map[NodeRef][]ResourceIndex{}
	for _, cl := range db.Classes {
		for _, ag := range atomicGroups[cl.ClassGroup] {
			ags = append(ags, ag)
			// Add to the Group -> index list map
			g2ags[cl.ClassGroup] = append(g2ags[cl.ClassGroup], ag.Index)
			for _, gref := range ag.Groups {
				g2ags[gref] = append(g2ags[gref], ag.Index)
			}
		}
	}
	// Sort the AtomicGroups
	slices.SortFunc(ags, func(a, b *AtomicGroup) int {
		if a.Index < b.Index {
			return -1
		}
		return 1
	})
	lg := len(ags)
	lr := len(db.Rooms)
	lt := len(db.Teachers)

	ttdata := &TtData{
		NDays:         days,
		NHours:        hours,
		HoursPerWeek:  days * hours,
		ActivitySlots: slices.Repeat([]TimeSlot{-1}, activities+1),
		ResourceWeeks: make([]ActivityIndex, (lg+lr+lt)*days*hours),
	}
	return ttdata
}

func (ttdata *TtData) BlockResource(resource ResourceIndex, slot TimeSlot) {
	ttdata.ResourceWeeks[int(resource)*ttdata.HoursPerWeek+int(slot)] = -1
}

func (ttdata *TtData) PrepareRooms(db *base.DbTopLevel) {
	lr := len(db.Rooms)
	r2tt := map[NodeRef]ResourceIndex{}
	for _, r := range db.Rooms {
		r2tt[r.Id] = i
		ttdata.Resources[i] = r
		i++
	}

}

func xxx() {
	// Allocate a vector for pointers to all Resources: teachers, (atomic)
	// student groups and (real) rooms.
	// Allocate a vector for pointers to all Activities, keeping the first
	// entry free (0 should be an invalid ActivityIndex).
	// Allocate a vector for a week of time slots for each Resource. Each
	// cell represents a timetable slot for a single resource. If it is
	// occupied – by an ActivityIndex – that indicates which Activity is
	// using the Resource at this time. A value of -1 indicates that the
	// time slot is blocked for this Resource.

	lt := len(db.Teachers)
	lr := len(db.Rooms)

	ags := []*AtomicGroup{}
	g2ags := map[Ref][]ResourceIndex{}
	for _, cl := range db.Classes {
		for _, ag := range ttinfo.AtomicGroups[cl.ClassGroup] {
			ags = append(ags, ag)
			// Add to the Group -> index list map
			g2ags[cl.ClassGroup] = append(g2ags[cl.ClassGroup], ag.Index)
			for _, gref := range ag.Groups {
				g2ags[gref] = append(g2ags[gref], ag.Index)
			}
		}
	}

	// Sort the AtomicGroups
	slices.SortFunc(ags, func(a, b *AtomicGroup) int {
		if a.Index < b.Index {
			return -1
		}
		return 1
	})

	lg := len(ags)

	// If using a single vector for all slots:
	ttinfo.Resources = make([]any, lt+lr+lg)
	ttinfo.TtSlots = make([]ActivityIndex, (lt+lr+lg)*ttinfo.SlotsPerWeek)

	// The slice cells are initialized to 0 or nil, according to slice type.

	// Copy the AtomicGroups to the beginning of the Resources slice.
	i := 0
	for _, ag := range ags {
		if ag.Index != i {
			base.Bug.Fatalf("Atomic group index != resource index:\n"+
				"  -- %d: %+v\n", i, ag)
		}
		ttinfo.Resources[i] = ag
		//fmt.Printf(" :: %+v\n", ag)
		i++
	}
	ttinfo.NAtomicGroups = i

	t2tt := map[Ref]ResourceIndex{}
	for _, t := range db.Teachers {
		t2tt[t.Id] = i
		ttinfo.Resources[i] = t
		i++
	}
	r2tt := map[Ref]ResourceIndex{}
	for _, r := range db.Rooms {
		r2tt[r.Id] = i
		ttinfo.Resources[i] = r
		i++
	}

	// Add the pseudo activities due to the NotAvailable lists
	ttinfo.addBlockers(t2tt, r2tt)

	// Get preliminary constraint info – needed for the call to addActivity
	ttinfo.processConstraints()

	// Add the remaining Activity information
	ttinfo.addActivityInfo(t2tt, r2tt, g2ags)

}

//

/* TODO: Split into components?
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
*/
