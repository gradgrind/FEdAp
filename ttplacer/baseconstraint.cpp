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

void LessonGroup::find_possible_slots(
    PlacementData* pdata)
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
    const int slots_per_day = pdata->slots_per_day;
    const auto tsn = possible_slots.size();
    for (size_t i = 0; i < tsn; ++i) {
        int dn = -1; // current duration (to avoid permutations)
        auto ts0 = possible_slots[i];
        auto day0 = ts0 / slots_per_day;
        auto j = i;
        auto ts = ts0;
        std::vector<TimeSlot> tsvec;
        for (const auto d0 : durations) {
            for (int d = 0; d < d0; ++d) {
                tsvec.emplace_back(ts);

                if (++j == tsn)
                    goto fail1;
                ts = possible_slots[j];
                if (ts / slots_per_day != day0)
                    goto fail0;
            }
        }
    fail0: // goto next first slot
    }
fail1:
}
