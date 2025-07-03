# Constraints

The constraints determine which activities can be placed in which time slot. The most obvious and straightforward are those connected with physical quantities (which can be regarded, in a certain sense, as "resources"), such as student-groups, teachers and rooms. Any of these can only be associated with one of the activities placed in a given time slot. Other constraints specify restrictions on the placements which come from policy or other considerations, for example, maximum lessons per day for a teacher or student group, lessons in one subject being on different days, lunch breaks or other "gaps" (slots with no activity).

## Hard and soft constraints

Some constraints are necessary, in the sense that not fulfilling them would be physically impossible (e.g., a teacher can only be in one lesson at a time) or that school policy requires their fulfilment absolutely (e.g. a certain student group may only have lessons in the morning). These are called hard constraints. Any acceptable timetable must fulfil these constraints.

There are also so-called soft constraints, whose fulfilment is desirable, but not absolutely necessary. These constraints can be "weighted", which allows the specification of the level of importance for each of the constraints.

## Handling constraints

There are various ways of handling the constraints during the generation of a timetable. The approach adopted here is to attempt to first generate a complete timetable taking only the hard constraints into consideration. The result will then be judged on the basis of the soft constraints. Strategies for producing results with better fulfilment of the soft constraints form a more-or-less separate part of the process. Variations may be possible which take some of the soft constraints into account while dealing with the hard constraints.
