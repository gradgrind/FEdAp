#include "baseconstraint.h"
#include "structures.h"

//BaseConstraint::BaseConstraint(){}

//ResourcePlacement::ResourcePlacement(){}

bool ResourcePlacement::hardConstraint(
    PlacementData* pdata, TimeSlot t)
{
    //TODO: should this constraint be inline?
    // Check the slot in the week's vector for this resource
    //return wkdata.at(resource_index).at(t) == 0;
}

void LessonGroup::find_possible_slots(PlacementData* pdata)
{
    //TODO
    auto slots_per_week = pdata->slots_per_week;
    int i = 0;
    while (i < slots_per_week) {
        for (const auto rsrc : resources) {
            auto rweek = pdata->resource_timetables[rsrc];
            if (rweek[i] != 0)
                goto next_slot;
        }
        // no collisions => possible slot
        possible_slots.emplace_back(i);
    next_slot:
        ++i;
    }
}

void LessonGroup::find_possible_placements()
{
    //TODO
}
