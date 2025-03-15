# Courses Editor

The courses editor shows a table of courses filtered by class or teacher. It may be useful to offer filtering on lesson subject, too.

## Courses, Supercourses and Subcourses

Three types of course are supported. A standard course has a subject, one or more teachers, one or more student groups, a room specification and a number of lesson durations (the durations of the individual lessons in a school week). Apart from the subject these are all optional, allowing some flexibility in regard to what a "course" actually is – it could also cover teachers' meetings, lunch-breaks, etc. There will also be other optional parameters to cover constraints.

Courses can also be grouped together, for example when there are particular time-slots in the week in which the courses are taught in a concentrated way for a number of weeks, followed by some weeks of another course, and so on throughout the year. For the timetable, these "supercourses" need their own names, in the same way that standard courses have a subject name. The "subcourses" which are the individual teaching blocks of the supercourses, are much like standard courses, except that they don't have lesson durations, these being supplied by the supercourse to which they belong. Also timetable constraints for these courses will be a function of the supercourse, though the subcourses might well have configuration parameters for uses outside of the timetable module.

The filtered course tables contain standard courses and subcourses. The tables are read-only. To allow editing of the entries there is also a panel, rather like an HTML form. In this panel there are also controls to show and edit the supercourse of a subcourse. Certain properties of a subcourse (especially the lesson durations) are actually properties of the supercourse, so where these are shown as belonging to the subcourse, the values will be shared by all subcourses of this supercourse. Thus, when a subcourse is edited it may not be enough to redisplay just the subcourse itself, it might be necessary to redisplay all the other subcourses of the supercourse as well. The extra cost of redisplaying the whole table might not be so great, so that it might be sensible to redisplay the whole table when a change is made.

A complication is that a subcourse can "belong" to more than one supercourse.
It is understandable in that the subgroup basically – as far as
the timetabling is concerned – functions as a set of "resources"
(teachers, groups, etc.). The supercourses provide the lessons (and a
name/tag). Thus the panel specifying the supercourses must allow for more than one entry. Selection of one of these would then show the lessons (and perhaps constraints) associated with this supercourse. In the main course table, the lessons field would show the lessons from all supercourses associated with a subcourse.

## Table fields

The table has the following fields (columns):
    subject, groups, teachers, lesson lengths, properties(?), rooms

## Editor form fields
The side panel ("form") allows editing of the selected course and has, in addition to the table's editable fields and the supercourse information, the following fields:
    block-name, block-subject, constraints

## Summary information

There might also be summary information about the number of lessons for
the teacher or groups set in the filter determining the courses shown.

## Editing operations

Apart from editing the various fields of a course, it should be possible to delete courses and add new ones. Further, a course may be added to a supercourse (possibly defining a new supercourse). If the supercourse is new, it has no lessons and those of the subcourse (if any) could be transferred to the supercourse. Otherwise the course's lessons should be discarded, perhaps with a warning. Also removing a course from a supercourse should be possible. The lessons of that supercourse would no longer be relevant for the subcourse, however, if it has no other supercourse, the lessons could be copied (or transferred if the change results in the supercourse being discarded completely).
