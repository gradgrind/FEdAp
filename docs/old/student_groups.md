# Student Groups and their Naming Scheme

For the purpose of the timetable, students are divided into groups which share lessons. This could be done in a variety of ways, but to keep things simple here a particular structure and nomenclature is adopted.

## Student-Groups

The students are divided into administrative and teaching teaching groups which I will call Student-Groups. Each student belongs to one and only one of these. These will generally correspond to the groups that sit in the lessons. These groups will often be year groups (students of a particular age), though this is not necessary. This could be conveyed in the group's name by beginning with the year, e.g. `11A` and `11B`.

## Divisions and Division-Groups

The Student-Groups can also be subdivided into smaller subgroups. These divisions can be performed in multiple ways, allowing different divisions for particular groups of subjects.

Thus, Student-Group `11A` may be divided into smaller groups for language lessons. These subgroups also have a name and are separated from the Student-Group name by a `.`, e.g. `11A.X` and `11A.Y`. Each student in `11A` is in either subgroup `X` or subgroup `Y`. The combination of the subgroups `X` and `Y` is called a Division, the subgroups are call Division-Groups. A Division may be given a descriptive name-tag, but this serves no internal purpose.

This places the limitation on Student-Group and Division-Group names that they cannot contain a `.`. It would be possible to make the separator configurable, but this might make the documentation and examples more confusing, so it was decided to use a fixed separator. For cases where this is in conflict with the desired naming scheme (say the school uses 11.A and 11.B for the year groups), the interface to the timetabling software's internal data could perform appropriate transformations.

## Atomic-Groups

When placing students groups in a particular time slot it is necessary to check that no students have more than one activity at a time. One way of doing this is to produce special subgroups representing the smallest combinations of students which can be placed in a time slot independently of each other. These are here called Atomic-Groups. In an extreme case these could correspond to individual pupils, but a more normal example would be, say, when a Student-Group has two Divisions, one with groups `X` and `Y`, another with groups `P` and `Q`. The Atomic-Groups for this Student-Group would then be `{X-P, X-Q, Y-P, Y-Q}`.
