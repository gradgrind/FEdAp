package backend

import (
	"fmt"
	"gradgrind/backend/base"
	"slices"
	"strconv"
)

func InitView() {
	courseState = &CoursesState{}
	vals := []string{`<>Class>`, `<>Teacher>`, `<>Subject>`}
	res := make([]string, len(vals))
	for i, s := range vals {
		res[i], _ = base.I18N(s)
	}
	Gui("COMBOBOX_SET_ITEMS", "table_type", res)
}

var (
	courseState *CoursesState
)

func init() {
	commandMap["GET_COURSES"] = getCourses
}

// TODO: This is probably not what is needed. Rather I will want the data
// for the courses tables, or rather a single course table?
// I could pass a class, teacher (or subject?) and get the data to be
// displayed?

/* For a class, teacher (or subject?) a list of courses is needed.
The table has the following fields (columns):
    subject, groups, teachers, lesson lengths, properties(?), rooms
The side panel ("form") allows editing of the selected course and has
the additional fields:
    block-name, block-subject, constraints
This information should be fetchable separately, but if the back-end has
a concept of "selected", this will need to be fetched with the list data
for the selected course.
There might also be summary information about the number of lessons for
the teacher or groups.
*/

/* A subcourse can "belong" to more than one supercourse, which complicates
things a bit. It is understandable in that the subcourse basically – as far as
the timetabling is concerned – just functions as a set of "resources"
(teachers, groups, etc.). The supercourses provide only the lessons (and a
name/tag).
Editing possibilities:
    - add course to block (there is a problem with existing lessons in the
    course, if the block is new – has no lessons – they could be transferred,
    otherwise need confirmation that they should be discarded)
    - remove course from block (if block now has no subcourses it can be
    discarded?), option to keep or discard fully freed block?
*/

// TODO?
func ToInt(val any) int {
	i, ok := val.(int)
	if !ok {
		f, ok := val.(float64)
		if !ok {
			panic(fmt.Sprintf("Not an integer: %#v", val))
		}
		i = int(f)
		if float64(i) != f {
			panic(fmt.Sprintf("Invalid (integer) value: %g", f))
		}
	}
	return i
}

func getCourses(cmd map[string]any, outmap map[string]any) bool {
	//TODO: Need the type of table and the unit for which the courses are
	// to be collected.
	tt := ToInt(cmd["TableType"])
	if tt != courseState.tableType {
		// Change the table type
		courseState.tableType = tt
		coursesInit()
	}

	outmap["Courses"] = DB.Courses
	outmap["SuperCourses"] = DB.SuperCourses
	outmap["SubCourses"] = DB.SubCourses
	return true
}

/* The handler for the courses gui ...

1) Full initialization of the view.
 - decide on classes, teachers or subjects
 - set the selection list
 - collect the data for the initial selection

*/

type CoursesState struct {
	// Maintains the state of the current courses view.
	// Note that this is dependent on the current database and needs
	// to be renewed when a new data set is loaded.
	tableType int // index to the table types list, see [CourseViewTypes()]
	tableShow int // index to the items list
	items     []Ref
}

func coursesInit() {
	var clist [][]string
	items := []Ref{}
	switch courseState.tableType {
	case 0: // show a class's courses
		// Set the classes
		clist = make([][]string, len(DB.Classes))
		for i, c := range DB.Classes {
			items = append(items, c.Id)
			clist[i] = []string{c.Tag, c.Name}
		}

	case 1: // show a teacher's courses
		// Set the teachers
		clist = make([][]string, len(DB.Teachers))
		for i, c := range DB.Teachers {
			items = append(items, c.Id)
			clist[i] = []string{c.Tag, c.Name}
		}

	case 2: // show the courses in a particular subject
		// Set the subjects
		clist = make([][]string, len(DB.Subjects))
		for i, c := range DB.Subjects {
			items = append(items, c.Id)
			clist[i] = []string{c.Tag, c.Name}
		}

	default:
		base.Report(`<Bug>coursesInit: courseState.tableType = %d>`,
			courseState.tableType)
		return
	}
	courseState.items = items
	Gui("COMBOBOX_SET_ITEMS_X", "table_show", clist)

	// Then I would need the courses involving the currently selected
	// class (etc.) to set up the table rows. The printing module may help
	// with some ideas here ...
	// The rows could be sent as a list of string-lists.

	// Finally, one of the courses (by default the first?) would be
	// selected (leading to its display in the editor panel).
}

/*
type CourseDisplayData struct {
	Subject string
	SubjectName string
	Teachers []string
	TeacherNames []string
	Groups []string
	Rooms []string
	//RoomNames []string
	LessonLengths []int
	...
}
*/

func getCourseShowData(cp base.CourseInterface) map[string]any {
	courseDisplayData := map[string]any{}
	// Subject
	sbj := getElementNames(cp.GetSubject())
	courseDisplayData["Subject"] = sbj[0]
	courseDisplayData["SubjectName"] = sbj[1]
	// Teachers
	tlist := []string{}
	tlistnames := []string{}
	for _, tref := range cp.GetTeachers() {
		tp := getElementNames(tref)
		tlist = append(tlist, tp[0])
		tlistnames = append(tlistnames, tp[1])
	}
	courseDisplayData["Teachers"] = tlist
	courseDisplayData["TeacherNames"] = tlistnames
	// Groups
	glist := []string{}
	for _, gref := range cp.GetGroups() {
		glist = append(glist, getElementNames(gref)[0])
	}

	//TODO: This might be better outside of this function, so that it
	// can accumulate lesson/block data
	if sub, ok := cp.(*base.SubCourse); ok {
		ulist := []string{}
		for _, i := range sub.Units {
			ulist = append(ulist, strconv.Itoa(i))
		}
		courseDisplayData["Units"] = ulist
	} else if c, ok := cp.(*base.Course); ok {
		ulist := []string{}
		for _, lref := range c.Lessons {
			//l := DB.Elements[lref].(*base.Lesson)
			l := elcast[*base.Lesson](lref)
			ulist = append(ulist, strconv.Itoa(l.Duration))
		}
		courseDisplayData["Units"] = ulist
	} else {
		//TODO
		panic("Unexpected element")
	}

	// Lesson lengths
	// This is more difficult? A normal course has a list of lesson lengths,
	// but what should be shown for a subcourse?
	// What about just calling this lengths, which could mean either lessons
	// or units/weeks? The block-nature would be shown either in this cell
	// (e.g. with a special symbol) or in the Properties column.
	courseDisplayData["Units"] = []string{}

	// Properties
	// What should be shown here?

	courseDisplayData["Properties"] = []string{}
	// Room – can be real room, choice or group
	room := cp.GetRoom()
	courseDisplayData["Room"] = room[0]
	courseDisplayData["RoomName"] = room[1]
	return courseDisplayData

	//TODO: Maybe all field values should be strings rather than lists?

	//TODO: Gather totals for info line
}

// TODO: Is this useful?
func elcast[T any](ref Ref) T {
	ep, ok := DB.Elements[ref]
	if !ok {
		//TODO
		panic("Not an element ref: " + ref)
	}
	e, ok := ep.(T)
	if !ok {
		//TODO
		panic(fmt.Sprintf("Element not of type %T, ref: %s", e, ref))
	}
	return e
}

func coursesForTeacher(tindex int) {
	courses := []map[string]any{}
	//TODO: get tref from tindex
	tref := courseState.items[courseState.tableShow]
	for _, cp := range DB.Courses {
		if slices.Contains(cp.Teachers, tref) {
			// Include this course
			courses = append(courses, getCourseShowData(cp))
		}
	}

	outmap["Courses"] = DB.Courses
	outmap["SuperCourses"] = DB.SuperCourses
	outmap["SubCourses"] = DB.SubCourses

}
