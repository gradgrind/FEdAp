# Lesson Placement

The _Lessons_ arise in connection with the _Courses_. A _Course_ can have a set of _Lessons_ associated with it, which all share the _Course's_ resources, etc. The _Lessons_ can have different durations.

The specified durations of the _Lessons_ are treated as hard constraints, i.e. they must be respected. Constraints can be added to limit whether and how multiple _Lessons_ of a course on one day are permissible. There may well be one or more "days-between" constraints associated with the _Course_.

A _TtUnit_ is a set of _Lessons_ which are connected by constraints limiting how they can be placed relative to each other. The idea is to construct a list of possible placements for the whole set (i.e. for the _TtUnit_), considered as a unit. The danger is that some of these lists could get rather large, but this will need to be tested.

## Fixed placements

First the blocked time slots for student-groups, teachers, and perhaps rooms, must be entered into their week-tables (as -1). Then each activity must be placed, first checking that the placement is possible as regards its resources. This can already call for flexibility in regard to the room choices.

## Preparations for placement of unfixed activities

Then, every remaining course is assessed for each slot, ignoring the lesson duration. This provides a list of possible slots, perhaps arranged in day lists, allowing further filtering to allow for longer durations and (hard) "days-between" constraints.

After this, any other hard constraints affecting relative placement of _Lessons__ should be processed to produce a list of possible placements for each set of linked _Lessons_ (i.e. for each _TtUnit_.

Actually, by giving weights to the lesson-set placements also soft constraints could be included in this processing. Should only unconstrained placements be used until these have all been somehow rejected or should also the weighted ones be included from the start (with correspondingly lower probability of being chosen first)?

## Building a _TtUnit_

This might work best by performing actual placements, then backtracking when necessary.

 - Start with the _Lessons_ of a single course, placing one after the other in all possible combinations. Placements can be rejected if hard "days-between" constraints are broken. If soft constraints are broken, their penalties can be accumulated and the combination's total can be associated with it. Repetitions can be avoided by starting with the longest _Lessons_ and for each duration only checking each slot once. This results in a list of time-slot tuples for the course. After each cycle the placements must be undone.

 - `TODO`: It may be possible to extend this procedure with further constraints on relationships between lessons or courses, if the interactions don't get too complicated ... This needs investigating.
