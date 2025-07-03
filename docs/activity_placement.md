# Activity Placement

## Where do the activities come from?

The activities arise in connection with the _Courses_. A _Course_ can have a set of _Activities_ associated with it, which all share the _Course's_ resources, etc. There may well be one or more "days-between" constraints associated with the _Course_.

## Fixed placements

First the blocked time slots for student-groups, teachers, and perhaps rooms, must be entered into their week-tables (as -1). Then each activity must be placed, first checking that the placement is possible as regards its resources. This can already call for replacement of room choices.

## Preparations for placement of unfixed activities

Then, every remaining activity is assessed for each slot, ignoring the activity duration. This provides a list of possible slots, perhaps arranged in day lists, allowing further filtering to allow for longer durations and (hard) "days-between" constraints.

After this, any other hard constraints affecting relative placement of activities should be processed to produce lists of possible placements for sets of activities which are thus linked.

Actually, by giving weights to the activity-set placements also soft constraints could be included in this processing. Should only unconstrained placements be used until these have all been somehow rejected or should also the weighted ones be included from the start (with correspondingly lower probability of being chosen first)?
