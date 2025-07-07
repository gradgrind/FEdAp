#ifndef BASECONSTRAINT_H
#define BASECONSTRAINT_H

#include "structures.h"
#include <vector>

class BaseConstraint
{
public:
    //BaseConstraint();

    virtual bool hardConstraint(PlacementData* pdata, TimeSlot t) = 0;
    virtual int softConstraint(PlacementData* pdata, TimeSlot t) = 0;
};

class DifferentDays : BaseConstraint
{
    std::vector<int> activities;

public:
    bool hardConstraint(PlacementData* pdata, TimeSlot t) override;
};

//TODO--?
class ResourcePlacement : BaseConstraint
{
    int resourceIndex; //TODO: actually, there could be a pointer to
    // the actual vector for the week rather than this indirect approach ...
    // ... or not if I want one object for multiple runs ...???
    bool hardConstraint(PlacementData* pdata, TimeSlot t) override;
};

#endif // BASECONSTRAINT_H
