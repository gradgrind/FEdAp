# Automatic timetabling – an idea

The task is to place a set of "activities" in a timetable such that there are no "resource collisions", all "hard" constraints are satisfied and the penalties arising from unsatisfied "soft" constraints are minimized.

The "resources" are teachers, rooms and groups of students (which are divided into "classes" and their subdivisions, "groups"). A resource can be associated with at most one activity at any given time.

The activities arise in connection with the courses (see _Courses_ and _SuperCourses_). A course can have a set of activities associated with it, which all share the course's resources, etc. The activities can have different durations.

## Modelling the timetable

The timetable is based on a grid of slots (days * slots-per-day). There is such a grid for each resource. The entries are the indexes of the activities to which the resource is assigned at the respective times. Activity indexes start at 1, leaving index 0 free to be used for "no activity". There is also a special activity index (-1) indicating that the resource is blocked (not available) at this time.

Each activity is to be assigned to one of the time slots (or consecutive slots if the activity has a duration greater than 1). To record the assignments there is an array in which each activity has an entry, its slot index. Slot indexes start at 0 (indicating the first time slot on the first day). Currently unplaced activities would have a special slot index (-1).

## Activity placement

It probably simplifies the process considerably if room resources where there is a choice are not handled until a complete timetable is available.

Some use can be made of the relationships between activities to reduce the number of possible combinations, but the blocked slots – for resources – and the fixed activities should be dealt with first.

### Blocked slots

Each resource has a (possibly empty) list of blocked time slots. Before any activities are processed, these blocks should be entered into the resource-allocation grids.

### Fixed placements

Activities with fixed times should then be placed in the timetable (activity slot and resource allocation). For these activities, resource conflicts with other fixed activities and blocked resources are not permitted, but other constraints are not taken into consideration.

### Preparations for placement of unfixed activities

The activities within a course share a common set of resources. Thus, for a basic check on where they can be placed, only one activity must be tested (for each duration). So these activities are first grouped according to duration – I call these `basic activity groups` – and a basic placement test (just for conflicts with blocked times and fixed activities) is carried out. Note that fixed activities are not included. Then the possible combinations are determined. The number of combinations is smaller than the number of permutations (if there is more than one member), so this can already lead to a significant reduction in the number of possible placements to be tested.

The next stage is to build `timetable units` by taking certain constraints into consideration – these are constraints which limit the relative placement of activities. An example would be a constraint specifying that certain activities may not be on the same day. In the simplest case this would connect activities which are already in the same `basic activity groups`, but if the activities are from different `basic activity groups`, these would be tied together in a larger `timetable unit`.

This can get quite complicated, so it seems sensible to split this stage into two passes. To begin with, all the relevant constraints are read and used to bind `basic activity groups` into `timetable units`. The idea is to construct a list of possible placements for the whole set considered as a unit.

After this, all possible activity-slot combinations of each `timetable unit` are tested by performing the placements and subsequently – where there are no collisions – applying the constraints. This will lead to a penalty score, which is stored along with the combination as one possibility for the `timetable unit`, or to the rejection of the combination. As this process involves real placements, the placement state must be saved before each combination is tested and restored afterwards.

It is hoped that this stage can lead to a further significant reduction in the number of combinations to be tested by the main placement algorithm when this sort of constraint is present. There might however, be a problem with combinatorial explosion if there are enough of this type of constraint to cause large `timetable units` to be built. This would need to be investigated.

The result of all this is a collection of `timetable units`, each of which has a list of "placements", each consisting of a list of pairs, (activity index, time slot), and a penalty (>=0). The binding constraints have already been subsumed into these "placements" and are no longer needed in this context. However, if manual placement of individual activities is to be supported, that might conflict with the data structures described here (see below, "Placement of individual activities").

### Further constraints

There may well be constraints that are best applied after a `timetable unit` has been placed in the timetable. There may be a way of attaching these to the `timetable units` in a way that can be processed efficiently. Alternatively, it may make sense to trigger these in connection with placements to particular resource slots.

Then there are constraints, which probably only make sense after all activities have been placed.

## Placement of individual activities

It could be a bit difficult to make the manual placement of individual activities compatible with the use of the `timetable units` for automatic placement. It is probably desirable to allow manual placement, so the best compromise might be to treat it as a separate stage, _after_ automatic placement. It would be based purely on the activities and independent of the `timetable units`, so the original constraints would need to be retained, and attached to the activities where relevant. Where they have already been subsumed into the `timetable units` they would be ignored during automatic processing.

## Placement of `timetable units`

A `timetable unit` has a list of potentially possible placements for its whole set of activities. **TODO?** These could be tested sequentially, "ticking off" the ones that have already been dealt with – if that is compatible with the placement algorithm. This "ticking off" might be handled simply by keeping an index of the "next possibility". It may well be sensible to sort the possibilities, presumably at random, weighting according to the penalty, so that each run will use a different sequence.

The basic placement algorithm would perhaps simply try to place each of the `timetable units`, one after the other, backtracking if at some stage no possibility is found. Various tweaks may be necessary, e.g. bigger jumps when the process seems to be a bit stuck, but I haven't got that far yet. It might be helpful to sort the `timetable units` before starting, so that those with fewer possibilities are handled first.
