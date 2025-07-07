#include "activity.h"

Activity::Activity() {}

/* Test placing activity a in time-slot t
 * 
 * wkdata = pdata.wkdatavecs
 * for (const auto & rsrc : a.resources) {
 *     if (wkdata[rsrc][t] != 0) return false;
 * }
 * 
 * for (const auto & cnstrnt : a.hard_constraints) {
 *     if (!cnstrnt.hardConstraint(pdata, t)) return false;
 * }
 * 
 */

bool Activity::testPlacement(
    PlacementData* pdata, TimeSlot t)
{
    // Check simple necessary "resources"
    auto weekdata = pdata->resource_timetables;
    for (const auto& rsrc : resources) {
        if (weekdata[rsrc][t] != 0)
            return false;
    }
    // Check other hard constraints
    for (const auto& cnstrnt : hard_constraints) {
        if (!cnstrnt->hardConstraint(pdata, t))
            return false;
    }
}

bool placeActivityGroup(
    PlacementData* pdata, TimeSlot t)
{}
