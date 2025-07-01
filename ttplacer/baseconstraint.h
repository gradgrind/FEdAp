#ifndef BASECONSTRAINT_H
#define BASECONSTRAINT_H

#include <vector>
typedef int TimeSlot;

struct PlacementData
{
    std::vector<int> ResourceTimetable;
};

class BaseConstraint
{
public:
    //BaseConstraint();

    virtual bool hardConstraint(PlacementData* pdata, TimeSlot t) = 0;
    virtual int softConstraint(PlacementData* pdata, TimeSlot t) = 0;
};

class ResourcePlacement : BaseConstraint
{
    int resourceIndex; //TODO: actually, there could be a pointer to
    // the actual vector for the week rather than this indirect approach ...
    bool hardConstraint(PlacementData* pdata, TimeSlot t) override;
};

#endif // BASECONSTRAINT_H
