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
	NDays        int
	NHours       int
	HoursPerWeek int
	// `ActivitySlots` is an array of activities + 1 entries, each entry being
	// the time slot in which the corresponding activity has been placed, or
	// -1 if unplaced. There is no activity with index 0.
	ActivitySlots []TimeSlot
	// `ResourceWeeks` contains the allocations of the "resources" (atomic
	// groups, teachers, rooms) to activities (indexes). This is organized
	// as an array of "week-chunks" (`HoursPerWeek` entries), one for each
	// resource in the `Resources array`.
	ResourceWeeks []ActivityIndex

	// `Resources` is an array mapping resource indexes to their corresponding
	// atomic group, teacher or room nodes (it contains pointers).
	Resources    []any
	RoomIndex    map[NodeRef]ResourceIndex
	TeacherIndex map[NodeRef]ResourceIndex

	// `AtomicGroups` maps a class or group NodeRef to its list of atomic
	// group indexes.
	AtomicGroups map[NodeRef][]ResourceIndex
	// `ClassDivisions` maps a class NodeRef to its list of divisions, each of
	// these being a list of group NodeRefs.
	ClassDivisions map[NodeRef][][]NodeRef

	// Set up by `GatherCourseInfo`
	Activities []*TtActivity
	//?? ActivityCourses []*TtCourseInfo
	//?? CourseInfo      map[NodeRef]*TtCourseInfo // key is Course or SuperCourse

	//???
	DayIndex     map[string]int
	HourIndex    map[string]int
	GroupIndexes map[string][]ResourceIndex
	VirtualRooms map[string]TtVirtualRoom

	//?? basic_activity_groups map[int]*BasicActivityGroup
	//?? CollectedBags         map[*BasicActivityGroup]*BagCollection
}

// BasicSetup performs the initialization of a TtData structure, collecting
// "resources" (atomic student groups, teachers and rooms) and "activities".
func BasicSetup(db *base.DbTopLevel) *TtData {
	days := len(db.Days)
	hours := len(db.Hours)
	tt_data := &TtData{
		NDays:        days,
		NHours:       hours,
		HoursPerWeek: days * hours,
		//?? ActivitySlots: slices.Repeat([]TimeSlot{-1}, activities+1),
	}

	course_info := CollectCourses(db)
	class_divisions := FilterDivisions(db, course_info)

	// Atomic groups: an atomic group is a "resource", it is an ordered list
	// of single groups, one from each division.
	// The atomic groups take the lowest resource indexes (starting at 0).
	// `AtomicGroups` maps the classes and groups to a list of their resource
	// indexes.
	tt_data.MakeAtomicGroups(db, class_divisions)

	// Add teachers and rooms to resource array
	tt_data.TeacherResources(db)
	tt_data.RoomResources(db)
	tt_data.ResourceWeeks = make([]ActivityIndex,
		(len(tt_data.Resources))*days*hours)

	// Get the activities for the timetable
	tt_data.MakeActivities(db, course_info)
	// ... initially all unplaced
	tt_data.ActivitySlots = slices.Repeat(
		[]TimeSlot{-1},
		len(tt_data.Activities))

	// Add the pseudo activities due to the NotAvailable lists of classes,
	// teachers and rooms.
	tt_data.BlockResources(db)

	/* TODO
	// Get preliminary constraint info – needed for the call to addActivity
	ttinfo.processConstraints()

	// Add the remaining Activity information
	ttinfo.addActivityInfo(t2tt, r2tt, g2ags)
	*/

	return tt_data
}

func (tt_data *TtData) BlockResource(resource ResourceIndex, slot TimeSlot) {
	tt_data.ResourceWeeks[int(resource)*tt_data.HoursPerWeek+int(slot)] = -1
}

func (tt_data *TtData) TeacherResources(db *base.DbTopLevel) {
	tt_data.TeacherIndex = map[NodeRef]ResourceIndex{}
	for _, t := range db.Teachers {
		i := len(tt_data.Resources)
		tt_data.TeacherIndex[t.Id] = i
		tt_data.Resources = append(tt_data.Resources, t)
	}
}

func (tt_data *TtData) RoomResources(db *base.DbTopLevel) {
	tt_data.RoomIndex = map[NodeRef]ResourceIndex{}
	for _, r := range db.Rooms {
		i := len(tt_data.Resources)
		tt_data.RoomIndex[r.Id] = i
		tt_data.Resources = append(tt_data.Resources, r)
	}
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
