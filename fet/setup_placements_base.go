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
	DayIndex      map[string]int
	HourIndex     map[string]int
	TeacherIndex  map[string]int
	GroupIndexes  map[string][]int
	RoomIndex     map[string]int
	VirtualRooms  map[string]TtVirtualRoom
	Activities    []TtActivity
	MapActivity   map[int]int // translate activity id, FET -> internal
}

type TtTeacher struct {
	Id            Ref
	Tag           string
	ResourceIndex int
}

type TtRoom struct {
	Id            Ref
	Tag           string
	ResourceIndex int
}

type TtVirtualRoom struct {
	RoomIndexes       []int
	RoomChoiceIndexes [][]int
}

type TtAtomicGroup struct {
	//Id            Ref
	Tag           string
	ResourceIndex int
}

func (tt_data *TtData) TimeSlotIndex(day string, hour string) int {
	d := tt_data.DayIndex[day]
	h := tt_data.HourIndex[hour]
	return d*tt_data.NHours + h
}

func (tt_data *TtData) PrepareResources(fetdata *fet) {
	tt_data.NDays = fetdata.Days_List.Number_of_Days
	tt_data.DayIndex = map[string]int{}
	for i, d := range fetdata.Days_List.Day {
		tt_data.DayIndex[d.Name] = i
	}
	tt_data.NHours = fetdata.Hours_List.Number_of_Hours
	tt_data.HourIndex = map[string]int{}
	for i, h := range fetdata.Hours_List.Hour {
		tt_data.HourIndex[h.Name] = i
	}
	slots_per_week := tt_data.NDays * tt_data.NHours
	fmt.Printf("n days = %d, n hours = %d\n", tt_data.NDays, tt_data.NHours)

	tt_data.TeacherIndex = map[string]int{}
	tt_data.GroupIndexes = map[string][]int{}
	tt_data.RoomIndex = map[string]int{}
	tt_data.VirtualRooms = map[string]TtVirtualRoom{}

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

	// Teachers and groups are referenced directly in the Activity nodes.
	// Rooms are not, they are attached to the activities by constraints. The
	// constraints specify a list of rooms from which one is to be chosen.
	// To complicate matters, there are – in addition to the "real" rooms –
	// "virtual" rooms. If the choice list contains just one room, it may be real
	// or virtual, otherwise only real rooms are acceptable.
	// A virtual room contains a list of needed real rooms and a list of room
	// choices, whose components must also be real rooms.
	// FET may well support more complex arrangements which are currently not
	// supported by this scheme.

	for _, r := range fetdata.Rooms_List.Room {
		//TODO: The handling of virtual rooms may need delaying until all the
		// real rooms have been read.
		if r.Virtual {
			// This is basically a list of room choices, where each "choice"
			// may actually be a single room, i.e. no choice!
			item_p := TtVirtualRoom{}
			for _, rset := range r.Set_of_Real_Rooms {
				if rset.Number_of_Real_Rooms == 1 {
					r := rset.Real_Room[0]
					tix, ok := tt_data.RoomIndex[r]
					if ok {
						item_p.RoomIndexes = append(item_p.RoomIndexes, tix)
					} else {
						base.Error.Fatalf("Unknown (real) room: %s", r)
					}
				} else {
					rooms := []int{}
					for _, r := range rset.Real_Room {
						tix, ok := tt_data.RoomIndex[r]
						if ok {
							rooms = append(rooms, tix)
						} else {
							base.Error.Fatalf("Unknown (real) room: %s", r)
						}
					}
					item_p.RoomChoiceIndexes = append(item_p.RoomChoiceIndexes,
						rooms)
				}
			}
			tt_data.VirtualRooms[r.Name] = item_p
			fmt.Printf("ROOM: %s – %v\n", r.Name, item_p)

		} else {
			tid := base.NewId()
			tix := len(tt_data.Resources)
			item_p := &TtRoom{tid, r.Name, tix}
			tt_data.Resources = append(tt_data.Resources, item_p)
			tt_data.RoomIndex[r.Name] = tix
			tt_data.ResourceWeeks = append(tt_data.ResourceWeeks,
				make([]int, slots_per_week))
			fmt.Printf("ROOM: %s – %d\n", r.Name, tix)
		}

	}
}
