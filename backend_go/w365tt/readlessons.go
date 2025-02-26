package w365tt

import (
	"gradgrind/backend/base"
)

func (db *DbTopLevel) readLessons(newdb *base.DbTopLevel) {
	for _, e := range db.Lessons {
		// The course must be Course or Supercourse.
		_, ok := db.CourseMap[e.Course]
		if !ok {
			base.Report(
				`<Error>Lesson %s:\n  Invalid course: %s>`,
				e.Id, e.Course)
			continue
		}
		// Check the Rooms.
		reflist := []Ref{}
		for _, rref := range e.Rooms {
			_, ok := db.RealRooms[rref]
			if ok {
				reflist = append(reflist, rref)
			} else {
				base.Report(
					`<Error>Invalid Room in Lesson %s:\n  -- %s>`,
					e.Id, rref)
			}
		}
		// Check for flag "SubstitutionService"
		flags := []string{}
		for _, f := range e.Flags {
			if f == "SubstitutionService" {
				c, ok := newdb.Elements[e.Course].(*base.Course)
				if ok {
					if len(c.Groups) == 1 {
						cgref := c.Groups[0]
						cg := newdb.Elements[cgref].(*base.Group)
						if cg.Tag != "" {
							base.Report(
								`<Error>SubstitutionService lesson (%s) not
								> specified for full stand-in class>`,
								e.Id)
							continue
						}
						cl := newdb.Elements[cg.Class].(*base.Class)
						// Mark class as stand-ins object
						cl.Tag = ""
						cl.Letter = ""
						cl.Year = -1
					}
				} else {
					base.Report(
						`<Error>SubstitutionService lesson (%s) not in
						> standard course: %s>`,
						e.Id, e.Course)
				}
				continue
			}
			// Other flags are just copied over
			flags = append(flags, f)
		}
		n := newdb.NewLesson(e.Id)
		n.Course = e.Course
		n.Duration = e.Duration
		n.Day = e.Day
		n.Hour = e.Hour
		n.Fixed = e.Fixed
		n.Rooms = reflist
		n.Flags = flags
		n.Background = e.Background
		n.Footnote = e.Footnote
	}
}
