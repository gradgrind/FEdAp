# Terminology

I try to use a consistent naming scheme for the various elements used in this application, especially in regard to the construction of timetables. The terms I have chosen may not fit every usage of the application ideally, but I decided to have a fixed naming scheme in order to keep the description of the program as straightforward as possible – to avoid making descriptions of potentially complicated algorithms even more complicated. The interpretation of these terms can then be adapted to the concrete situation in which the application is used.

Most elements have a "name" and a "tag", the latter being a short form which is useful in documents where space is at a premium, such as a timetable. The tag is unique within a particular element type, serving as an identifier for the element. The name need not necessarily be unique. For example, a student might have a unique identification number – this would be the tag. It does occur occasionally that two students have the same name, but their tags will be different. If also the teachers have identification numbers, it may be that a teacher has the same number as a student, but this is permitted because the elements are of different types.

This application is aimed primarily at schools, so, on the whole, corresponding terms will be used. Even the educational arena has different terms for the various elements, depending, for example on country or fashion. Is a group of students/pupils to be called a year(-group), grade, class, or whatever? Some of these terms can easily become ambiguous, so I will try to define their meanings in the context of this application.

## Students: Year-groups, classes and other groupings ...

The "students" are the (mostly young) people who are subjected to the endeavours of the "teaching establishment". They are generally divided into age groups which share "activities".

### Classes

The "class" is the primary student grouping. Its tag comprises a number (the year/grade) followed by short text component. Both components are optional, but at least one must be present. Typical classes might be something like "3A" or "11X". The term "year" is reserved for the differentiation of school years, where it will indeed generally be called "school-year".

### Groups and Divisions

A class may be further divided into smaller teaching groups. I call the way in which a class is divided into groups a "division". The groups within a division share no students, so they can be taught in parallel. A simple division of a class would be into the two groups "A" and "B". The class may be divided in more than one way to satisfy the requirements of different teaching areas. For example, there could also be the division into groups "I", "II" and "III". Divisions have no tag, but can have a descriptive name.

Groups have tags which combine the class with a short group-tag, which is unique within the class, such as "2.A" or "12.S". `TODO`: The exact format of the class and group tags should probably not be too rigid. The use of '.' as a separator should be a configuration option, maybe even the whole way in which the tags are constructed from the components, so that the program can adapt to the existing naming conventions of an institution.

It may also be convenient to define additional groups which are combinations of groups *within a division*. For example, consider a division `{A, B, C}`. In some activities one might have the groups `A` and `B` together, so that a group `X`(= `A + B`) would be useful.

### Atomic-Groups (internal)

When placing student groups in a particular time slot it is necessary to check that no students have more than one activity at a time. One way of doing this is by means of special subgroups representing the smallest combinations of students which can be placed in a time slot independently of each other. These are here called "atomic groups". In an extreme case these could correspond to individual pupils, but a more normal example would be the example above, where a class has two divisions, `{A, B}` and `{I, II, III}`. The atomic groups arising from this situation would be `{A-I, A-II, A-III, B-I, B-II, B-III}`, i.e. the Cartesian product of the divisions.

## Teachers

The "teachers" are the people who carry out the educational endeavours. They are associated for this purpose with groups of students in teaching units called "activities". A person can not be both a teacher and a student in this model – this is a consequence of the model structure and has no educational or political foundations. Like students, they can only be in one place at a time.

## Rooms

The "rooms" are the spaces (presumably physical) available for the execution of the activities – where one or more teachers spend time with one or more groups of students. There may, however, be special cases of "activities" with no teacher ("free study"?) or no students (such as a conference), which still need a room.

To cover more needs, a couple of special room types are supported:

 - *RoomGroup*: This specifies a list of _Rooms_, all of which are necessary.

 - *RoomChoiceGroup*: This specifies a list of _Rooms_ one of which must be available for the _Activity_.

 Each _Course_ or _SubCourse_ (see below) can specify a single _Room_, a _RoomGroup_ or a _RoomChoiceGroup_ (or no room).

## Resources

The three elements described above (classes/groups, teachers and rooms) are treated as "resources" for the purpose of constructing a timetable. An _Activity_ brings a set of these elements (zero or more of each type) together in one or more consecutive time-slots within the (weekly?) program known as the timetable.

## Activities

An "activity" is like a lesson (teaching unit), occupying a number of resources for a certain time. It need not always directly involve teaching, however. This term would here also cover a lunch-break or an administrative meeting. The length of an _Activity_ is here termed "duration". It is basically the same as an "Activity" in `FET`.

## Time-Slots

For the purposes of the timetable, the "week" (which could have pretty well any number of days) is divided into "time-slots". These time-slots determine the times at which an activity can take place. They will generally be of equal lengths, to avoid unfairness and other problems with activity placement, but variation is, in principal, possible and there can be breaks between activities. To avoid placing activities of longer duration (needing multiple time-slots) over longer breaks, it may be necessary to add corresponding constraints (e.g. permitted starting times). A time slot will usually be presented as a combination of day and slot-within-the-day, such as "Mo:2" (second slot on Monday), though its internal representation will probably be rather different (for reasons of efficiency).

## Courses, SuperCourses, SubCourses

_Courses_ couple associated _Resources_ (student-groups, teachers, rooms) with a set of _Activities_, for which suitable time-slots must be found.

Three types of "course" are supported. The "normal" type is just called a _Course_, where the same activities are to take place at the same time every week. Some courses, however, are taught in blocks of, say, a few weeks at the same time every week. Other course blocks occupy the same weekly time-slots, but at other times in the year. This eventuality is covered by what are here called _SuperCourses_ and _SubCourses_. A _SuperCourse_ occupies regular time-slots, just like a normal _Course_, but has no directly associated _Resources_. Associated with the _SuperCourse_ are a number of _SubCourses_, which specify the _Resources_ for each of the real courses which are to take place in the time slots. To enable more elaborate course time arrangements, a _SubCourse_ can be associated with more than one _SuperCourse_.

Each _Course_ and _SuperCourse_ has a list of durations (the unit being a time-slot), determining the _Activities_ belonging to this course. This corresponds roughly to an "Activity_Group" in `FET`.

### When might a _SubCourse_ be associated with more than one _SuperCourse_?

A number of time-slots within the week may be set aside for teaching blocks in which multiple classes are involved. A reason might be that the teachers for the courses which are to take place in these blocks are active in more than one class. It might not be typical, but it is possible that the exact times are slightly different in the various classes. Consider a case in which the normal distribution of times is the first two activities of each day (every day). In one class, however, the time must be different on one day. To cover this case, three _SuperCourses_ could be used, one for the days where all classes have the same times, the other two for the two times on the special day. All _SubCourses_ would be associated with the first of these _SuperCourses_ and also with one of the other two, depending on which class they are aimed at.
