package fet

//TODO: This stuff is just temporary ... experimenting ...

import (
	"fedap/base"
	"fmt"
	"maps"
	"slices"
	"strings"
)

type TtData struct {
	NDays         int
	NHours        int
	Resources     []any
	ResourceWeeks [][]int
	TeacherIndex  map[string]int
	GroupIndexes  map[string][]int
}

type TtTeacher struct {
	Id            Ref
	Tag           string
	ResourceIndex int
}

type TtClass struct {
	Id  Ref
	Tag string
	//ResourceIndex int
}

type TtAtomicGroup struct {
	//Id            Ref
	Tag           string
	ResourceIndex int
}

func prepare_placements(fetdata *fet) {
	tt_data := TtData{}
	tt_data.NDays = fetdata.Days_List.Number_of_Days
	tt_data.NHours = fetdata.Hours_List.Number_of_Hours
	slots_per_week := tt_data.NDays * tt_data.NHours
	fmt.Printf("n days = %d, n hours = %d\n", tt_data.NDays, tt_data.NHours)

	tt_data.TeacherIndex = map[string]int{}
	tt_data.GroupIndexes = map[string][]int{}

	for _, t := range fetdata.Teachers_List.Teacher {
		tid := base.NewId()
		tix := len(tt_data.Resources)
		item_p := &TtTeacher{tid, t.Name, tix}
		tt_data.Resources = append(tt_data.Resources, item_p)
		tt_data.TeacherIndex[t.Name] = tix
		tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
			make([]int, slots_per_week))
	}
	fmt.Printf("n teachers: %d\n", len(tt_data.TeacherIndex))

	for _, c := range fetdata.Students_List.Year {
		sep := c.Separator
		agsep := "#"
		if sep == "#" {
			agsep = "$"
		}
		//-- fmt.Printf("§ %d - %#v\n", len(c.Group), c)

		// If there are no groups, the class is the atomic-group (with
		// a slightly changed name)
		// If a group has subgroups, these are the atomic-groups.
		// Otherwise the group will be the subgroup (with
		// a slightly changed name)

		if len(c.Group) == 0 {
			// Make new atomic-group
			atag := c.Name + agsep
			tix := len(tt_data.Resources)
			item_p := TtAtomicGroup{atag, tix}
			tt_data.GroupIndexes[c.Name] = []int{tix}
			tt_data.Resources = append(tt_data.Resources, item_p)
			tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
				make([]int, slots_per_week))
			fmt.Printf("§ %s - %d\n", c.Name, tix)

		} else {
			agmap := map[string]int{}

			for _, g := range c.Group {
				if len(g.Subgroup) == 0 {
					// Make new atomic-group
					atag := strings.ReplaceAll(g.Name, sep, agsep)
					tix := len(tt_data.Resources)
					item_p := TtAtomicGroup{atag, tix}
					tt_data.GroupIndexes[g.Name] = []int{tix}
					tt_data.Resources = append(tt_data.Resources, item_p)
					tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
						make([]int, slots_per_week))
					fmt.Printf("§ %s - %d\n", g.Name, tix)

				} else {
					// Use the existing subgroups as atomic-groups
					aglist := []int{}
					for _, ag := range g.Subgroup {
						tix, ok := agmap[ag.Name]
						if !ok {
							tix = len(tt_data.Resources)
							agmap[ag.Name] = tix
							item_p := TtAtomicGroup{ag.Name, tix}
							tt_data.Resources = append(tt_data.Resources, item_p)
							tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
								make([]int, slots_per_week))
						}
						aglist = append(aglist, tix)
					}
					tt_data.GroupIndexes[g.Name] = aglist
					fmt.Printf("§ %s - %v\n", g.Name, aglist)

				}
			}
			// Now get atomic-groups for whole class
			aglist := slices.Sorted(maps.Values(agmap))
			tt_data.GroupIndexes[c.Name] = aglist
			fmt.Printf("§§ %s - %v\n", c.Name, aglist)
		}
	}

	//TODO: Rooms can be real or virtual. The rooms are attached to the activities
	// indirectly, via constraints. These can specify a list of rooms from which
	// one is to be chosen. If there is more than one, they must (here) be all
	// real. A single room may be real or virtual.
	// All rooms used in a virtual room should be real (here). FET may well support
	// more complex arrangements, but that may be too much here.

	//TODO: Blocked time-slots, different-days (and other constraints?)
}
