#ifndef ACTIVITY_H
#define ACTIVITY_H

#include "baseconstraint.h"
#include <vector>

struct Activity
{
    std::vector<int> activity_indexes;
    std::vector<int> resources;
    std::vector<std::vector<int>> room_choices;
    //TODO: room_choices perhaps as a hard_constraints entry?
    std::vector<BaseConstraint*> hard_constraints;

    Activity();

    bool testPlacement(PlacementData* pdata, TimeSlot t);
};

struct ActivityGroup
{
    std::vector<Activity> activities;

    /* Find possible slots.
     * 
     * Are all the activities identical? What about small differences,
     * perhaps different durations? ... I suspect it may be clearer to
     * limit the scope to really identical activities.
     * 
     * Test placement of an activity, taking into account any constraints
     * on their relative placement, especially days between them.
     * 
     * This might be "built in" to the Activity handling (see
     * `activity_indexes`), in that placing a non-fixed member would imply
     * also placing the others.
     * 
     * Fixed activities could be at the end of the list, so that only
     * those before them need to be placed algorithmically.
     *  
     */
};

#endif // ACTIVITY_H
