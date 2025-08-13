package timetable

import (
	"fedap/base"
	"fmt"
	"slices"
	"strings"
)

// A CourseInfo is an intermediate representation of a course (Course or
// SuperCourse) for the timetable.
type CourseInfo struct {
	Id         NodeRef
	Subject    NodeRef
	Groups     []NodeRef
	Teachers   []NodeRef
	Rooms      []NodeRef // entries can be Room, RoomGroup or RoomChoiceGroup
	Activities []NodeRef
}

// Make a shortish string view of a CourseInfo – can be useful in tests
func View(db *base.DbTopLevel, cinfo *CourseInfo) string {
	tlist := []string{}
	for _, t := range cinfo.Teachers {
		tlist = append(tlist, db.Ref2Tag(t))
	}
	glist := []string{}
	for _, g := range cinfo.Groups {
		glist = append(glist, db.Ref2Tag(g))
	}
	return fmt.Sprintf("<Course %s/%s:%s>",
		strings.Join(glist, ","),
		strings.Join(tlist, ","),
		db.Ref2Tag(cinfo.Subject),
	)
}

// Collect courses (Course and SuperCourse) and their activities.
// Build a list of CourseInfo structures.
func CollectCourses(db *base.DbTopLevel) []*CourseInfo {
	cinfo_list := []*CourseInfo{}

	// Gather the SuperCourses.
	for _, spc := range db.SuperCourses {
		cref := spc.Id
		groups := []NodeRef{}
		teachers := []NodeRef{}
		rooms := []NodeRef{}
		for _, sbcref := range spc.SubCourses {
			sbc := db.Elements[sbcref].(*base.SubCourse)
			// Add groups
			if len(sbc.Groups) != 0 {
				groups = append(groups, sbc.Groups...)
			}
			// Add teachers
			if len(sbc.Teachers) != 0 {
				teachers = append(teachers, sbc.Teachers...)
			}
			// Add rooms
			if sbc.Room != "" {
				rooms = append(rooms, sbc.Room)
			}
		}
		// Eliminate duplicate resources by sorting and then compacting
		slices.Sort(groups)
		slices.Sort(teachers)
		slices.Sort(rooms)
		cinfo_list = append(cinfo_list, &CourseInfo{
			Id:         cref,
			Subject:    spc.Subject,
			Groups:     slices.Compact(groups),
			Teachers:   slices.Compact(teachers),
			Rooms:      slices.Compact(rooms),
			Activities: spc.Lessons,
		})
	}

	// Gather the plain Courses.
	for _, c := range db.Courses {
		cref := c.Id
		rooms := []NodeRef{}
		if c.Room != "" {
			rooms = append(rooms, c.Room)
		}
		cinfo_list = append(cinfo_list, &CourseInfo{
			Id:         cref,
			Subject:    c.Subject,
			Groups:     c.Groups,
			Teachers:   c.Teachers,
			Rooms:      rooms,
			Activities: c.Lessons,
		})
	}

	return cinfo_list
}

// `MakeActivities` creates the `TtActivity` structures ...
func (tt_data *TtData) MakeActivities(db *base.DbTopLevel, cinfo_list []*CourseInfo) {
	tt_data.Activities = []*TtActivity{{}} // first entry is empty
	for _, cinfo := range cinfo_list {
		// Get resource indexes
		resources := []ResourceIndex{}
		for _, r := range cinfo.Groups {
			agilist, ok := tt_data.AtomicGroups[r]
			if !ok {
				base.Bug.Fatalf("Unknown group: %s", r)
			}
			resources = append(resources, agilist...)
		}
		for _, r := range cinfo.Teachers {
			agi, ok := tt_data.TeacherIndex[r]
			if !ok {
				base.Bug.Fatalf("Unknown teacher: %s", r)
			}
			resources = append(resources, agi)
		}
		roomchoices := [][]ResourceIndex{}
		for _, r := range cinfo.Rooms {
			agi, ok := tt_data.RoomIndex[r]
			if ok {
				resources = append(resources, agi)
				continue
			}
			// Not a Room – it can be a RoomGroup or RoomChoiceGroup
			elem, ok := db.Elements[r]
			if !ok {
				base.Bug.Fatalf("Unknown room: %s", r)
			}
			rg, ok := elem.(*base.RoomGroup)
			if ok {
				for _, rr := range rg.Rooms {
					agi, ok = tt_data.RoomIndex[rr]
					if !ok {
						base.Bug.Fatalf(
							"Unknown room in RoomGroup %s: %s", r, rr)
					}
					resources = append(resources, agi)
				}
				continue
			}
			rcg, ok := elem.(*base.RoomChoiceGroup)
			if ok {
				rooms := []ResourceIndex{}
				for _, rr := range rcg.Rooms {
					agi, ok = tt_data.RoomIndex[rr]
					if !ok {
						base.Bug.Fatalf(
							"Unknown room in RoomChoiceGroup %s: %s", r, rr)
					}
					rooms = append(rooms, agi)
				}
				roomchoices = append(roomchoices, rooms)
				continue
			}
			base.Bug.Fatalf("Expecting room element, found: %s", r)
		}
		// Check for duplicate resources
		//TODO: Is this a bug, or can it happen "legally"?
		// If the latter, duplicates would need to be removed.
		rset := map[ResourceIndex]bool{}
		for _, ri := range resources {
			if rset[ri] {
				base.Bug.Fatalf("Duplicate resource: %+v\n in course: %+v",
					tt_data.Resources[ri], cinfo)
			}
		}

		//TODO: What about collecting the BAGs (distinct durations) here?

		// Build a TtActivity for each lesson
		for _, lref := range cinfo.Activities {
			l := db.Elements[lref].(*base.Lesson)
			if slices.Contains(l.Flags, "SubstitutionService") {
				cinfo.Groups = nil
			}
			//p := -1
			//if l.Day >= 0 {
			//	p = l.Day*tt_data.NHours + l.Hour
			//}
			ttl := &TtActivity{
				Id: ActivityIndex(len(tt_data.Activities)),
				//Placement:  p,
				Duration:    int16(l.Duration),
				Fixed:       l.Fixed,
				Resources:   resources,
				RoomChoices: roomchoices,
				//Lesson:     l,
				//CourseInfo: cinfo,
			}
			tt_data.Activities = append(tt_data.Activities, ttl)
		}
	}
}

/* TODO??
type VirtualRoom struct {
	Rooms       []NodeRef   // only ("real") Rooms
	RoomChoices [][]NodeRef // list of ("real") Room lists
}
*/

/* GatherCourseInfo collects the course information which is relevant for
// the timetable, building a TtCourseInfo structure for each course with
// activities.
func (ttdata *TtData) GatherCourseInfo(db *base.DbTopLevel) {
	// Gather the Groups, Teachers and "rooms" for the Courses and
	// SuperCourses with lessons/activities (only).
	// Gather the Activities for these Courses and SuperCourses.
	// Also, the SuperCourses (with lessons/activities) get a list of their
	// SubCourses.

	//ttdata.CourseInfo = map[NodeRef]*TtCourseInfo{}
	courseInfo := map[NodeRef]*CourseInfo{}
	ttdata.Activities = make([]*TtActivity, 1) // 1-based indexing, 0 is invalid

	// Collect courses with lessons/activities.
	roomData := ttdata.collectCourses(db)

	// Prepare the internal room structure, filtering the room lists of
	// the SuperCourses.
	for cref, crlist := range roomData {
		// Join all Rooms and the Rooms from RoomGroups into a "compulsory"
		// list. Then go through the RoomChoiceGroups. If one contains a
		// compulsory room, ignore the choice.
		// The result is a list of Rooms and a list of Room-choice-lists,
		// which can be converted into a tt virtual room.
		rooms := []NodeRef{}
		roomChoices := [][]NodeRef{}
		for _, rref := range crlist {
			rx := db.Elements[rref]
			_, ok := rx.(*base.Room)
			if ok {
				rooms = append(rooms, rref)
			} else {
				rg, ok := rx.(*base.RoomGroup)
				if ok {
					rooms = append(rooms, rg.Rooms...)
				} else {
					rc, ok := rx.(*base.RoomChoiceGroup)
					if !ok {
						base.Bug.Fatalf(
							"Invalid room in course %s:\n  %s\n",
							cref, rref)
					}
					roomChoices = append(roomChoices, rc.Rooms)
				}
			}
		}
		// Remove duplicates in Room list.
		slices.Sort(rooms)
		rooms = slices.Compact(rooms)
		// Filter choice lists.
		roomChoices = slices.DeleteFunc(roomChoices, func(rcl []NodeRef) bool {
			for _, rc := range rcl {
				if slices.Contains(rooms, rc) {
					return true
				}
			}
			return false
		})

		// Add virtual room to CourseInfo item
		cinfo := ttdata.CourseInfo[cref]
		cinfo.Room = TtVirtualRoom{
			Rooms:       rooms,
			RoomChoices: roomChoices,
		}

		// Check the allocated rooms at Lesson.Rooms
		ttinfo.checkAllocatedRooms(cinfo)
	}
}
*/

/*
func (ttinfo *TtInfo) checkAllocatedRooms(cinfo *CourseInfo) {
	// Check the room allocations for the lessons of the given course.
	// Report just one error per course.
	for _, lix := range cinfo.Lessons {
		l := ttinfo.Activities[lix].Lesson
		lrooms := l.Rooms
		// If no rooms are allocated, don't regard this as invalid.
		if len(lrooms) == 0 {
			return
		}
		// First check number of rooms
		vr := cinfo.Room
		if len(vr.Rooms)+len(vr.RoomChoices) != len(lrooms) {
			rlist := []string{}
			for _, rref := range lrooms {
				rlist = append(rlist, ttinfo.Ref2Tag[rref])
			}
			base.Warning.Printf("Lesson in Course %s has wrong number"+
				" of rooms allocated:\n  -- %+v (expected %d)\n",
				cinfo.Id, rlist, len(vr.Rooms)+len(vr.RoomChoices))
			return
		}
		// Check validity of "compulsory" rooms
		lrmap := map[Ref]bool{}
		for _, rref := range lrooms {
			lrmap[rref] = true
		}
		for _, rref := range vr.Rooms {
			if lrmap[rref] {
				delete(lrmap, rref)
			} else {
				base.Warning.Printf("Lesson in Course %s needs room %s\n",
					cinfo.Id, ttinfo.Ref2Tag[rref])
				return
			}
		}
		// Check validity of "chosen" rooms
		cmap := make([]Ref, 0, len(vr.RoomChoices))
		var fx func(i int) bool
		fx = func(i int) bool {
			if i == len(vr.RoomChoices) {
				return true
			}
			for _, rref := range vr.RoomChoices[i] {
				if lrmap[rref] && !slices.Contains(cmap, rref) {
					cmap = append(cmap, rref)
					if fx(i + 1) {
						return true
					}
					cmap = cmap[:i]
				}
			}
			return false
		}
		if !fx(0) {
			rlist := []string{}
			for rref := range lrmap {
				rlist = append(rlist, ttinfo.Ref2Tag[rref])
			}
			base.Warning.Printf("Lesson in Course %s has invalid"+
				" room-choice allocations: %+v\n",
				cinfo.Id, rlist)
			return
		}
	}
}
*/
