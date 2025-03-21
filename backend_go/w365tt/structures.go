package w365tt

import (
	"encoding/json"
	"gradgrind/backend/base"
	"gradgrind/backend/ttprint"
)

// The structures used for the "database", adapted to read from W365

type Ref = base.Ref // Element reference

type Info struct {
	Institution        string
	FirstAfternoonHour int
	MiddayBreak        []int
	Reference          string `json:"Scenario"`
}

type Day struct {
	Id   Ref
	Type string
	Name string
	Tag  string `json:"Shortcut"`
}

type Hour struct {
	Id    Ref
	Type  string
	Name  string
	Tag   string `json:"Shortcut"`
	Start string
	End   string
}

type TimeSlot struct {
	Day  int
	Hour int
}

type Teacher struct {
	Id               Ref
	Type             string
	Name             string
	Tag              string `json:"Shortcut"`
	Firstname        string
	NotAvailable     []TimeSlot `json:"Absences"`
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxDays          int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
}

func (t *Teacher) UnmarshalJSON(data []byte) error {
	// Customize defaults for Teacher
	t.MinLessonsPerDay = -1
	t.MaxLessonsPerDay = -1
	t.MaxDays = -1
	t.MaxGapsPerDay = -1
	t.MaxGapsPerWeek = -1
	t.MaxAfternoons = -1

	type tempT Teacher
	return json.Unmarshal(data, (*tempT)(t))
}

type Subject struct {
	Id   Ref
	Type string
	Name string
	Tag  string `json:"Shortcut"`
}

type Room struct {
	Id           Ref
	Type         string
	Name         string
	Tag          string     `json:"Shortcut"`
	NotAvailable []TimeSlot `json:"Absences"`
}

type RoomGroup struct {
	Id    Ref
	Type  string
	Name  string
	Tag   string `json:"Shortcut"`
	Rooms []Ref
}

type Class struct {
	Id               Ref
	Type             string
	Name             string
	Tag              string `json:"Shortcut"`
	Year             int    `json:"Level"`
	Letter           string
	NotAvailable     []TimeSlot `json:"Absences"`
	Divisions        []Division
	MinLessonsPerDay int
	MaxLessonsPerDay int
	MaxGapsPerDay    int
	MaxGapsPerWeek   int
	MaxAfternoons    int
	LunchBreak       bool
	ForceFirstHour   bool
}

func (t *Class) UnmarshalJSON(data []byte) error {
	// Customize defaults for Teacher
	t.MinLessonsPerDay = -1
	t.MaxLessonsPerDay = -1
	t.MaxGapsPerDay = -1
	t.MaxGapsPerWeek = -1
	t.MaxAfternoons = -1

	type tempT Class
	return json.Unmarshal(data, (*tempT)(t))
}

type Group struct {
	Id   Ref
	Type string
	Tag  string `json:"Shortcut"`
}

type Division struct {
	Id     Ref
	Type   string
	Name   string
	Groups []Ref
}

type Course struct {
	Id             Ref
	Type           string
	Subjects       []Ref
	Groups         []Ref
	Teachers       []Ref
	PreferredRooms []Ref
}

type SuperCourse struct {
	Id         Ref
	Type       string
	EpochPlan  Ref
	SubCourses []SubCourse
}

type SubCourse struct {
	Id             Ref
	Type           string
	Subjects       []Ref
	Groups         []Ref
	Teachers       []Ref
	PreferredRooms []Ref
}

type Lesson struct {
	Id         Ref
	Type       string
	Course     Ref // Course or SuperCourse Elements
	Duration   int
	Day        int
	Hour       int
	Fixed      bool
	Rooms      []Ref `json:"LocalRooms"` // only Room Elements
	Flags      []string
	Background string
	Footnote   string
}

type EpochPlan struct {
	Id   Ref
	Type string
	Tag  string `json:"Shortcut"`
	Name string
}

type DbTopLevel struct {
	Info         Info `json:"W365TT"`
	PrintTables  []*ttprint.PrintTable
	FetData      map[string]string
	Days         []*Day
	Hours        []*Hour
	Teachers     []*Teacher
	Subjects     []*Subject
	Rooms        []*Room
	RoomGroups   []*RoomGroup
	Classes      []*Class
	Groups       []*Group
	Courses      []*Course
	SuperCourses []*SuperCourse
	Lessons      []*Lesson
	EpochPlans   []*EpochPlan
	Constraints  []map[string]any

	// These fields do not belong in the JSON object:

	// RealRooms is used to check the validity of a (real) room reference
	// and to access the room's tag/shortcut
	RealRooms map[Ref]string `json:"-"`
	// RoomGroupMap is used to check the validity of a room-group reference
	RoomGroupMap map[Ref]bool `json:"-"`
	// SubjectMap is used to check the validity of a subject reference
	// and to access the subject's tag/shortcut
	SubjectMap map[Ref]string `json:"-"`
	// GroupRefMap is used to check the validity of a group reference and to
	// replace class references in courses by class-group references
	GroupRefMap map[Ref]Ref `json:"-"`
	// TeacherMap is used to check the validity of a teacher reference
	TeacherMap map[Ref]bool `json:"-"`
	CourseMap  map[Ref]bool `json:"-"`
	// SubjectTags maps a subject's tag/shortcut to its reference/Id
	SubjectTags map[string]Ref `json:"-"`
	// RoomTags maps a subject's tag/shortcut to its reference/Id
	RoomTags map[string]Ref `json:"-"`
	// RoomChoiceNames allows reuse of [RoomChoice] elements, it maps the
	// name, constructed from the component tags, to the reference/Id
	RoomChoiceNames map[string]Ref `json:"-"`
}

// Block all afternoons if nAfternnons == 0.
func (dbp *DbTopLevel) handleZeroAfternoons(
	notAvailable []TimeSlot,
	nAfternoons int,
) []base.TimeSlot {
	// Make a bool array and fill this in two passes, then remake list
	namap := make([][]bool, len(dbp.Days))
	nhours := len(dbp.Hours)
	// In the first pass, conditionally block afternoons
	for i := range namap {
		namap[i] = make([]bool, nhours)
		if nAfternoons == 0 {
			for h := dbp.Info.FirstAfternoonHour; h < nhours; h++ {
				namap[i][h] = true
			}
		}
	}
	// In the second pass, include existing blocked hours.
	for _, ts := range notAvailable {
		if ts.Hour < len(dbp.Hours) {
			// Exclude invalid hours
			namap[ts.Day][ts.Hour] = true
		}
		//TODO: else an error message?
	}
	// Build a new base.TimeSlot list
	na := []base.TimeSlot{}
	for d, naday := range namap {
		for h, nahour := range naday {
			if nahour {
				na = append(na, base.TimeSlot{Day: d, Hour: h})
			}
		}
	}
	return na
}
