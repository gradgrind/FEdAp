package base

// The structures used for the "database"
//TODO: Currently dealing only with the elements needed for the timetable

type Ref string // Element Id

type BasicElement struct {
	Id Ref
}

type TaggedElement struct {
	BasicElement
	Tag string
}

type TimeSlot struct {
	Day  int
	Hour int
}

type Division struct {
	Name   string
	Groups []Ref
}

type Info struct {
	Institution        string
	FirstAfternoonHour int
	MiddayBreak        []int
	Reference          string
}

type Day struct {
	TaggedElement
	Name string
}

type Hour struct {
	TaggedElement
	Name  string
	Start string
	End   string
}

type Teacher struct {
	TaggedElement
	Name             string
	Firstname        string
	NotAvailable     []TimeSlot
	MinLessonsPerDay int // default = -1
	MaxLessonsPerDay int // default = -1
	MaxDays          int // default = -1
	MaxGapsPerDay    int // default = -1
	MaxGapsPerWeek   int // default = -1
	MaxAfternoons    int // default = -1
	LunchBreak       bool
}

type Subject struct {
	TaggedElement
	Name string
}

type Room struct {
	TaggedElement
	Name         string
	NotAvailable []TimeSlot
}

func (r *Room) IsReal() bool {
	return true
}

type RoomGroup struct {
	TaggedElement
	Name  string
	Rooms []Ref
}

func (r *RoomGroup) IsReal() bool {
	return false
}

type RoomChoiceGroup struct {
	TaggedElement
	Name  string
	Rooms []Ref
}

func (r *RoomChoiceGroup) IsReal() bool {
	return false
}

type Class struct {
	TaggedElement
	Name             string
	Year             int
	Letter           string
	NotAvailable     []TimeSlot
	Divisions        []Division
	MinLessonsPerDay int // default = -1
	MaxLessonsPerDay int // default = -1
	MaxGapsPerDay    int // default = -1
	MaxGapsPerWeek   int // default = -1
	MaxAfternoons    int // default = -1
	LunchBreak       bool
	ForceFirstHour   bool
	ClassGroup       Ref
}

type Group struct {
	TaggedElement
	// These fields do not belong in the JSON object:
	Class Ref `json:"-"`
}

// TODO: Lessons should be ordered, highest duration first
type Course struct {
	BasicElement
	Subject  Ref
	Groups   []Ref
	Teachers []Ref
	Room     Ref // Room, RoomGroup or RoomChoiceGroup Element
	// These fields do not belong in the JSON object:
	Lessons []Ref `json:"-"`
}

func (c *Course) IsSuperCourse() bool {
	return false
}

func (c *Course) AddLesson(lref Ref) {
	c.Lessons = append(c.Lessons, lref)
}

type SuperCourse struct {
	BasicElement
	Subject Ref
	// These fields do not belong in the JSON object:
	SubCourses []Ref `json:"-"`
	Lessons    []Ref `json:"-"`
}

func (c *SuperCourse) IsSuperCourse() bool {
	return true
}

func (c *SuperCourse) AddLesson(lref Ref) {
	c.Lessons = append(c.Lessons, lref)
}

type SubCourse struct {
	BasicElement
	SuperCourses []Ref
	Subject      Ref
	Groups       []Ref
	Teachers     []Ref
	Room         Ref // Room, RoomGroup or RoomChoiceGroup Element
}

type GeneralRoom interface {
	IsReal() bool
}

type Lesson struct {
	BasicElement
	Course     Ref // Course or SuperCourse Elements
	Duration   int
	Day        int
	Hour       int
	Fixed      bool
	Rooms      []Ref
	Flags      []string `json:",omitempty"`
	Background string
	Footnote   string
}

type LessonCourse interface {
	IsSuperCourse() bool
	AddLesson(Ref)
}

type Constraint interface {
	CType() string
}

type PrintTable struct {
	Type          string
	TypstTemplate string
	TypstJson     string
	Pdf           string
	Typst         map[string]any
	Pages         []map[string]any
}

type DbTopLevel struct {
	Info             Info
	PrintTables      []*PrintTable
	Days             []*Day
	Hours            []*Hour
	Teachers         []*Teacher
	Subjects         []*Subject
	Rooms            []*Room
	RoomGroups       []*RoomGroup       `json:",omitempty"`
	RoomChoiceGroups []*RoomChoiceGroup `json:",omitempty"`
	Classes          []*Class
	Groups           []*Group       `json:",omitempty"`
	Courses          []*Course      `json:",omitempty"`
	SuperCourses     []*SuperCourse `json:",omitempty"`
	SubCourses       []*SubCourse   `json:",omitempty"`
	Lessons          []*Lesson      `json:",omitempty"`
	Constraints      []Constraint   `json:",omitempty"`

	// These fields do not belong in the JSON object:
	Elements map[Ref]any `json:"-"`
}

func (db *DbTopLevel) Ref2Tag(ref Ref) string {
	t, ok := db.Elements[ref].(TaggedElement)
	if !ok {
		Bug.Fatalf("No Ref2Tag for %s\n", ref)
	}
	return t.Tag
}
