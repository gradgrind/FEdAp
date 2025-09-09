package fet2xlsx

import (
	"fedap/readfet"
	"fmt"
	"maps"
	"slices"
)

type ActivityData struct {
	Id       int
	Subject  string
	Time     TimeSlot
	Duration int
	Teachers []string
	Students []string
	Rooms    []string
}

type TimeSlot struct {
	Day  int
	Hour int
}

func GetActivityData(fet *readfet.Fet) {
	// Map student "groups" to their "year"
	g2y := map[string]string{}
	for _, s := range fet.Students_List.Year {
		y := s.Name
		for _, g := range s.Group {
			g2y[g.Name] = y
		}
	}
	// Get the rooms
	rmap := map[int][]string{}
	for _, rc := range fet.Space_Constraints_List.
		ConstraintActivityPreferredRoom {

		if len(rc.Real_Room) != 0 {
			rmap[rc.Activity_Id] = rc.Real_Room
		} else {
			rmap[rc.Activity_Id] = []string{rc.Room}
		}
	}
	atimes := ActivityTimes(fet)
	for _, act := range fet.Activities_List.Activity {

		t, ok := atimes[act.Id]
		if !ok {
			t = TimeSlot{-1, -1}
		}
		// Filter the students, not including a group if the
		// class/year is included:
		// 1) Separate the years and the groups
		sset := map[string]struct{}{}
		groups := map[string]string{}
		for _, s := range act.Students {
			y, ok := g2y[s]
			if ok {
				groups[s] = y
			} else {
				sset[s] = struct{}{}
			}
		}
		// 2) Add the groups whose years are not included
		for g, y := range groups {
			_, ok := sset[y]
			if !ok {
				sset[g] = struct{}{}
			}
		}

		adata := &ActivityData{
			Id:       act.Id,
			Subject:  act.Subject,
			Time:     t,
			Duration: act.Duration,
			Teachers: act.Teacher,
			Students: slices.Sorted(maps.Keys(sset)),
		}

		// The rooms
		r, ok := rmap[act.Id]
		if ok {
			adata.Rooms = r
		}

		fmt.Printf(" -- %v\n", adata)
	}
}

func ActivityTimes(fet *readfet.Fet) map[int]TimeSlot {
	dmap := map[string]int{}
	for i, d := range fet.Days_List.Day {
		dmap[d.Name] = i
	}
	hmap := map[string]int{}
	for i, h := range fet.Hours_List.Hour {
		hmap[h.Name] = i
	}
	atimes := map[int]TimeSlot{}

	for _, atime := range fet.Time_Constraints_List.
		ConstraintActivityPreferredStartingTime {

		atimes[atime.Activity_Id] = TimeSlot{
			dmap[atime.Day], hmap[atime.Hour],
		}
	}
	return atimes
}
