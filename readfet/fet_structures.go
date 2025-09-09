// This doesn't cover all possible nodes in the FET file, especially the
// constraints, of which there are rather a lot, are not complete.
package readfet

import (
	"encoding/xml"
)

type Fet struct {
	Version          string `xml:"version,attr"`
	Mode             string
	Institution_Name string
	Comments         string // this can be a source reference
	Days_List        fetDaysList
	Hours_List       fetHoursList
	Teachers_List    fetTeachersList
	Subjects_List    fetSubjectsList
	Rooms_List       fetRoomsList
	Students_List    fetStudentsList
	//Buildings_List
	Activity_Tags_List     fetActivityTags
	Activities_List        fetActivitiesList
	Time_Constraints_List  timeConstraints
	Space_Constraints_List spaceConstraints
}

type fetDay struct {
	XMLName   xml.Name `xml:"Day"`
	Name      string
	Long_Name string
}

type fetDaysList struct {
	XMLName        xml.Name `xml:"Days_List"`
	Number_of_Days int
	Day            []fetDay
}

type fetHour struct {
	XMLName   xml.Name `xml:"Hour"`
	Name      string
	Long_Name string
}

type fetHoursList struct {
	XMLName         xml.Name `xml:"Hours_List"`
	Number_of_Hours int
	Hour            []fetHour
}

type fetTeacher struct {
	XMLName   xml.Name `xml:"Teacher"`
	Name      string
	Long_Name string
	Comments  string
}

type fetTeachersList struct {
	XMLName xml.Name `xml:"Teachers_List"`
	Teacher []fetTeacher
}

type fetSubject struct {
	XMLName   xml.Name `xml:"Subject"`
	Name      string
	Long_Name string
	Comments  string
}

type fetSubjectsList struct {
	XMLName xml.Name `xml:"Subjects_List"`
	Subject []fetSubject
}

type fetRoomsList struct {
	XMLName xml.Name `xml:"Rooms_List"`
	Room    []fetRoom
}

type fetRoom struct {
	XMLName                      xml.Name `xml:"Room"`
	Name                         string   // e.g. k3 ...
	Long_Name                    string
	Capacity                     int           // 30000
	Virtual                      bool          // false
	Number_of_Sets_of_Real_Rooms int           `xml:",omitempty"`
	Set_of_Real_Rooms            []realRoomSet `xml:",omitempty"`
	Comments                     string
}

type realRoomSet struct {
	Number_of_Real_Rooms int // normally 1, I suppose
	Real_Room            []string
}

type fetCategory struct {
	//XMLName             xml.Name `xml:"Category"`
	Number_of_Divisions int
	Division            []string
}

type fetSubgroup struct {
	Name string // 13.m.MaE
	//Number_of_Students int // 0
	//Comments string // ""
}

type fetGroup struct {
	Name string // 13.K
	//Number_of_Students int // 0
	//Comments string // ""
	Subgroup []fetSubgroup
}

type fetClass struct {
	//XMLName  xml.Name `xml:"Year"`
	Name      string
	Long_Name string
	Comments  string
	//Number_of_Students int (=0)
	// The information regarding categories, divisions of each category,
	// and separator is only used in the dialog to divide the year
	// automatically by categories.
	Number_of_Categories int
	Separator            string // CLASS_GROUP_SEP
	Category             []fetCategory
	Group                []fetGroup
}

type fetStudentsList struct {
	XMLName xml.Name `xml:"Students_List"`
	Year    []fetClass
}

type fetActivity struct {
	XMLName           xml.Name `xml:"Activity"`
	Id                int
	Teacher           []string `xml:",omitempty"`
	Subject           string
	Activity_Tag      string   `xml:",omitempty"`
	Students          []string `xml:",omitempty"`
	Active            bool
	Total_Duration    int
	Duration          int
	Activity_Group_Id int
	Comments          string
}

type fetActivitiesList struct {
	XMLName  xml.Name `xml:"Activities_List"`
	Activity []fetActivity
}

type fetActivityTag struct {
	XMLName   xml.Name `xml:"Activity_Tag"`
	Name      string
	Printable bool
}

type fetActivityTags struct {
	XMLName      xml.Name `xml:"Activity_Tags_List"`
	Activity_Tag []fetActivityTag
}

type timeConstraints struct {
	XMLName xml.Name `xml:"Time_Constraints_List"`
	//
	ConstraintBasicCompulsoryTime          basicTimeConstraint
	ConstraintStudentsSetNotAvailableTimes []studentsNotAvailable
	ConstraintTeacherNotAvailableTimes     []teacherNotAvailable

	ConstraintActivityPreferredStartingTime    []startingTime
	ConstraintActivityPreferredTimeSlots       []activityPreferredTimes
	ConstraintActivitiesPreferredTimeSlots     []preferredSlots
	ConstraintActivitiesPreferredStartingTimes []preferredStarts
	ConstraintMinDaysBetweenActivities         []minDaysBetweenActivities
	ConstraintActivityEndsStudentsDay          []lessonEndsDay
	ConstraintActivitiesSameStartingTime       []sameStartingTime

	ConstraintStudentsSetMaxGapsPerDay                  []maxGapsPerDay
	ConstraintStudentsSetMaxGapsPerWeek                 []maxGapsPerWeek
	ConstraintStudentsSetMinHoursDaily                  []minLessonsPerDay
	ConstraintStudentsSetMaxHoursDaily                  []maxLessonsPerDay
	ConstraintStudentsSetIntervalMaxDaysPerWeek         []maxDaysinIntervalPerWeek
	ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour []maxLateStarts
	ConstraintStudentsSetMaxHoursDailyInInterval        []lunchBreak

	ConstraintTeacherMaxDaysPerWeek          []maxDaysT
	ConstraintTeacherMaxGapsPerDay           []maxGapsPerDayT
	ConstraintTeacherMaxGapsPerWeek          []maxGapsPerWeekT
	ConstraintTeacherMaxHoursDailyInInterval []lunchBreakT
	ConstraintTeacherMinHoursDaily           []minLessonsPerDayT
	ConstraintTeacherMaxHoursDaily           []maxLessonsPerDayT
	ConstraintTeacherIntervalMaxDaysPerWeek  []maxDaysinIntervalPerWeekT
}

type basicTimeConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsoryTime"`
	Weight_Percentage int
	Active            bool
}

type spaceConstraints struct {
	XMLName                          xml.Name `xml:"Space_Constraints_List"`
	ConstraintBasicCompulsorySpace   basicSpaceConstraint
	ConstraintActivityPreferredRooms []roomChoice
	ConstraintActivityPreferredRoom  []placedRoom
	ConstraintRoomNotAvailableTimes  []roomNotAvailable
}

type basicSpaceConstraint struct {
	XMLName           xml.Name `xml:"ConstraintBasicCompulsorySpace"`
	Weight_Percentage int
	Active            bool
}

type preferredSlots struct {
	XMLName                        xml.Name `xml:"ConstraintActivitiesPreferredTimeSlots"`
	Weight_Percentage              string
	Teacher                        string
	Students                       string
	Subject                        string
	Activity_Tag                   string
	Duration                       string
	Number_of_Preferred_Time_Slots int
	Preferred_Time_Slot            []preferredTime
	Active                         bool
}

type preferredTime struct {
	//XMLName                       xml.Name `xml:"Preferred_Time_Slot"`
	Preferred_Day  string
	Preferred_Hour string
}

type preferredStarts struct {
	XMLName                            xml.Name `xml:"ConstraintActivitiesPreferredStartingTimes"`
	Weight_Percentage                  string
	Teacher                            string
	Students                           string
	Subject                            string
	Activity_Tag                       string
	Duration                           string
	Number_of_Preferred_Starting_Times int
	Preferred_Starting_Time            []preferredStart
	Active                             bool
}

type preferredStart struct {
	//XMLName                       xml.Name `xml:"Preferred_Starting_Time"`
	Preferred_Starting_Day  string
	Preferred_Starting_Hour string
}

type lessonEndsDay struct {
	XMLName           xml.Name `xml:"ConstraintActivityEndsStudentsDay"`
	Weight_Percentage string
	Activity_Id       int
	Active            bool
}

type activityPreferredTimes struct {
	XMLName                        xml.Name `xml:"ConstraintActivityPreferredTimeSlots"`
	Weight_Percentage              string
	Activity_Id                    int
	Number_of_Preferred_Time_Slots int
	Preferred_Time_Slot            []preferredTime
	Active                         bool
}

type sameStartingTime struct {
	XMLName              xml.Name `xml:"ConstraintActivitiesSameStartingTime"`
	Weight_Percentage    string
	Number_of_Activities int
	Activity_Id          []int
	Active               bool
}

type startingTime struct {
	XMLName            xml.Name `xml:"ConstraintActivityPreferredStartingTime"`
	Weight_Percentage  string
	Activity_Id        int
	Day                string
	Hour               string
	Permanently_Locked bool
	Active             bool
}

type minDaysBetweenActivities struct {
	XMLName                 xml.Name `xml:"ConstraintMinDaysBetweenActivities"`
	Weight_Percentage       string
	Consecutive_If_Same_Day bool
	Number_of_Activities    int
	Activity_Id             []int
	MinDays                 int
	Active                  bool
}

type studentsNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintStudentsSetNotAvailableTimes"`
	Weight_Percentage             int
	Students                      string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}

type teacherNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintTeacherNotAvailableTimes"`
	Weight_Percentage             int
	Teacher                       string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}

type notAvailableTime struct {
	XMLName xml.Name `xml:"Not_Available_Time"`
	Day     string
	Hour    string
}

type lunchBreakT struct {
	XMLName             xml.Name `xml:"ConstraintTeacherMaxHoursDailyInInterval"`
	Weight_Percentage   string
	Teacher             string
	Interval_Start_Hour string
	Interval_End_Hour   string
	Maximum_Hours_Daily int
	Active              bool
}

type maxGapsPerDayT struct {
	XMLName           xml.Name `xml:"ConstraintTeacherMaxGapsPerDay"`
	Weight_Percentage string
	Teacher           string
	Max_Gaps          int
	Active            bool
}

type maxGapsPerWeekT struct {
	XMLName           xml.Name `xml:"ConstraintTeacherMaxGapsPerWeek"`
	Weight_Percentage string
	Teacher           string
	Max_Gaps          int
	Active            bool
}

type minLessonsPerDayT struct {
	XMLName             xml.Name `xml:"ConstraintTeacherMinHoursDaily"`
	Weight_Percentage   string
	Teacher             string
	Minimum_Hours_Daily int
	Allow_Empty_Days    bool
	Active              bool
}

type maxLessonsPerDayT struct {
	XMLName             xml.Name `xml:"ConstraintTeacherMaxHoursDaily"`
	Weight_Percentage   string
	Teacher             string
	Maximum_Hours_Daily int
	Active              bool
}

type maxDaysT struct {
	XMLName           xml.Name `xml:"ConstraintTeacherMaxDaysPerWeek"`
	Weight_Percentage string
	Teacher           string
	Max_Days_Per_Week int
	Active            bool
}

type maxDaysinIntervalPerWeekT struct {
	XMLName             xml.Name `xml:"ConstraintTeacherIntervalMaxDaysPerWeek"`
	Weight_Percentage   string
	Teacher             string
	Interval_Start_Hour string
	Interval_End_Hour   string
	// Interval_End_Hour void ("") means the end of the day (which has no name)
	Max_Days_Per_Week int
	Active            bool
}

type lunchBreak struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMaxHoursDailyInInterval"`
	Weight_Percentage   string
	Students            string
	Interval_Start_Hour string
	Interval_End_Hour   string
	Maximum_Hours_Daily int
	Active              bool
}

type maxGapsPerDay struct {
	XMLName           xml.Name `xml:"ConstraintStudentsSetMaxGapsPerDay"`
	Weight_Percentage string
	Max_Gaps          int
	Students          string
	Active            bool
}

type maxGapsPerWeek struct {
	XMLName           xml.Name `xml:"ConstraintStudentsSetMaxGapsPerWeek"`
	Weight_Percentage string
	Max_Gaps          int
	Students          string
	Active            bool
}

type minLessonsPerDay struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMinHoursDaily"`
	Weight_Percentage   string
	Minimum_Hours_Daily int
	Students            string
	Allow_Empty_Days    bool
	Active              bool
}

type maxLessonsPerDay struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetMaxHoursDaily"`
	Weight_Percentage   string
	Maximum_Hours_Daily int
	Students            string
	Active              bool
}

type maxDaysinIntervalPerWeek struct {
	XMLName             xml.Name `xml:"ConstraintStudentsSetIntervalMaxDaysPerWeek"`
	Weight_Percentage   string
	Students            string
	Interval_Start_Hour string
	Interval_End_Hour   string
	// Interval_End_Hour void ("") means the end of the day (which has no name)
	Max_Days_Per_Week int
	Active            bool
}

type maxLateStarts struct {
	XMLName                       xml.Name `xml:"ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour"`
	Weight_Percentage             string
	Max_Beginnings_At_Second_Hour int
	Students                      string
	Active                        bool
}

type placedRoom struct {
	XMLName              xml.Name `xml:"ConstraintActivityPreferredRoom"`
	Weight_Percentage    int
	Activity_Id          int
	Room                 string
	Number_of_Real_Rooms int      `xml:",omitempty"`
	Real_Room            []string `xml:",omitempty"`
	Permanently_Locked   bool     // false
	Active               bool     // true
}

type roomChoice struct {
	XMLName                   xml.Name `xml:"ConstraintActivityPreferredRooms"`
	Weight_Percentage         int
	Activity_Id               int
	Number_of_Preferred_Rooms int
	Preferred_Room            []string
	Active                    bool // true
}

type roomNotAvailable struct {
	XMLName                       xml.Name `xml:"ConstraintRoomNotAvailableTimes"`
	Weight_Percentage             int
	Room                          string
	Number_of_Not_Available_Times int
	Not_Available_Time            []notAvailableTime
	Active                        bool
}
