package fet

import (
	"encoding/xml"
)

type startingTime struct {
	XMLName            xml.Name `xml:"ConstraintActivityPreferredStartingTime"`
	Weight_Percentage  string
	Activity_Id        int
	Preferred_Day      string
	Preferred_Hour     string
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

/* from FET FAQ

Q: Help on ConstraintMinDaysBetweenActivities.

A: It refers to a set of activities and involves a constant, N. For every pair of activities in the set, it does not allow the distance (in days) between them to be less than N. If you specify N=1, then this constraint means that no two activities can be scheduled in the same day. N=2 means that each two activities must be separated by at least one day

Example: 3 activities and N=2. Then, one can place them on Monday, Wednesday and Friday (5 days week).

Example2: 2 activities, N=3. Then, one can place them on Monday and Thursday, on Monday and Friday, then on Tuesday and Friday (5 days week).

The weight is recommended to be between 95.0%-100.0%. The best might be 99.75% or a value a little under 100%, because FET can detect impossible constraints this way and avoid them. The weight is subjective.

You can specify consecutive if same day. Please be careful, even if constraint min days between activities has 0% weight, if you select this consecutive if same day, this consecutive will be forced. You will not be able to find a timetable with the two activities in the same day, separated by break, not available or other activities, even if the constraint has weight 0%, if you select consecutive if same day.

Currently FET can put at most 2 activities in the same day if 'consecutive if same day' is true. FET cannot put 3 or more activities in the same day if 'consecutive if same day' is true.

Important: please do not input unnecessary duplicates. If you input for instance 2 constraints:

1. Activities 1 and 2, min days 1, consecutive if same day=true, weight=95%
2. Activities 1 and 2, min days 1, consecutive if same day=false, weight=95%
(these are different constraints),
then the outcome of these 2 constraints will be a constraint:

Activities 1 and 2, min days 1, consecutive if same day=true, weight=100%-5%*5%=99.75%, very high. This is because of FET algorithm.

You may however add 2 constraints for the same activities if you want 95% with min 2 days and 100% with min 1 day. These are not duplicates.

You might get an impossible timetable with duplicates, so beware. If you need to balance 3 activities in a 5 days week, you can add, in the new version 5.5.8 and higher, directly from the add activity dialog, 2 constraints. You just have to input min days 2, and FET will ask if you want to add a second constraint with min days 1. This way, you can ensure that the activities are balanced better (at least one day apart, usually 2 days apart)
*/

// *** Teacher constraints
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

// for MaxAfternoons
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

// *** Class constraints

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

// for MaxAfternoons
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

// for ForceFirstHour
type maxLateStarts struct {
	XMLName                       xml.Name `xml:"ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour"`
	Weight_Percentage             string
	Max_Beginnings_At_Second_Hour int
	Students                      string
	Active                        bool
}
