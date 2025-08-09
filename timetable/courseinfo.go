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
	Rooms      []NodeRef
	Activities []NodeRef
}

/* TODO??
type VirtualRoom struct {
	Rooms       []NodeRef   // only ("real") Rooms
	RoomChoices [][]NodeRef // list of ("real") Room lists
}
*/

// GatherCourseInfo collects the course information which is relevant for
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

/* TODO: make TtActivity structures
for i, cinfo := range cinfo_list {
	// Add lessons to CourseInfo
	llist := clessons[i]
	for _, lref := range llist {
		l := db.Elements[lref].(*base.Lesson)
		if slices.Contains(l.Flags, "SubstitutionService") {
			cinfo.Groups = nil
		}
		// Index of new Activity:
		ttlix := len(ttdata.Activities)
		p := -1
		if l.Day >= 0 {
			p = l.Day*ttdata.NHours + l.Hour
		}
		ttl := &TtActivity{
			Index:      ttlix,
			Placement:  p,
			Duration:   l.Duration,
			Fixed:      l.Fixed,
			Lesson:     l,
			CourseInfo: cinfo,
		}
		ttdata.Activities = append(ttdata.Activities, ttl)
		cinfo.Lessons = append(cinfo.Lessons, ttlix)
	}
}
*/

// Make a shortish string view of a CourseInfo – can be useful in tests
func (ttinfo *TtInfo) View(cinfo *CourseInfo) string {
	tlist := []string{}
	for _, t := range cinfo.Teachers {
		tlist = append(tlist, ttinfo.Ref2Tag[t])
	}
	glist := []string{}
	for _, g := range cinfo.Groups {
		gx, ok := ttinfo.Ref2Tag[g]
		if !ok {
			base.Bug.Fatalf("No Ref2Tag for %s\n", g)
		}
		glist = append(glist, gx)
	}

	return fmt.Sprintf("<Course %s/%s:%s>",
		strings.Join(glist, ","),
		strings.Join(tlist, ","),
		ttinfo.Ref2Tag[cinfo.Subject],
	)
}
