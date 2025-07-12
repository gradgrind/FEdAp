package fet

/*

PRE
ConstraintStudentsSetNotAvailableTimes
ConstraintTeacherNotAvailableTimes
ConstraintActivityPreferredStartingTime
ConstraintMinDaysBetweenActivities

++ same starting time
++ preferred starting times / slots (i.e. where there is a choice)

DURING?
ConstraintStudentsSetIntervalMaxDaysPerWeek
ConstraintStudentsSetMaxHoursDailyInInterval
ConstraintTeacherMaxHoursDaily
ConstraintTeacherIntervalMaxDaysPerWeek

++ (partly) activity begins/ends day
++ gap between activities within one day
++

POST
ConstraintStudentsSetMaxGapsPerWeek
ConstraintStudentsSetMinHoursDaily
ConstraintStudentsSetEarlyMaxBeginningsAtSecondHour
ConstraintTeacherMaxGapsPerDay
ConstraintTeacherMaxHoursDaily

++ (rest) activity begins/ends day

*/
