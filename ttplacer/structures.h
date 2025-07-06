#ifndef STRUCTURES_H
#define STRUCTURES_H

#include <string>
#include <vector>
using Ref = std::string;
using TimeSlot = int;

struct PlacementData
{
    std::vector<std::vector<int>> resource_timetables;
};

struct BasicStructure
{                     //TODO: Should probably be abstract
    const Ref Id;     // UUID, unique identifier
    std::string Tag;  // Short name or descriptor, unique if not empty
    std::string Name; // Long name or descriptor, not necessarily unique
};

struct Subject : BasicStructure
{};

struct GeneralRoom : BasicStructure
{
    virtual bool isReal() = 0;
};

struct Room : GeneralRoom
{
    bool isReal() override { return true; }
    std::vector<TimeSlot> NotAvailable;
};

//TODO: Can a RoomGroup contain a RoomChoiceGroup, and vice versa?
struct RoomGroup : GeneralRoom
{
    bool isReal() override { return false; }
    std::vector<Room *> Rooms;
};

struct RoomChoiceGroup : GeneralRoom
{
    bool isReal() override { return false; }
    std::vector<Room *> Rooms;
};

struct different_days
{
    int ndays;
    int weight; // 0 – 100
};

struct LessonGroup : BasicStructure
{
    std::vector<int> durations; // sorted, longest first
    std::vector<int> resources;
    std::vector<different_days> spread;

    std::vector<TimeSlot> possible_slots;
    std::vector<std::vector<TimeSlot>> posssible_placements;

    void find_possible_slots(PlacementData *pdata);
    void find_possible_placements();
};

#endif // STRUCTURES_H
