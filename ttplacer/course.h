#ifndef COURSE_H
#define COURSE_H

#include "structures.h"
#include <vector>

struct Course
{
    Course();
};

struct SubCourse : BasicStructure
{
    std::vector<Ref> SuperCourses;
    // Maybe rather: std::vector<SuperCourse*> superCourses; ...
    Ref Subject;
    std::vector<Ref> Groups;
    std::vector<Ref> Teachers;
    Ref Room; // Room, RoomGroup or RoomChoiceGroup Element
    SubCourse();
};

struct SuperCourse
{
    SuperCourse();
};

#endif // COURSE_H
