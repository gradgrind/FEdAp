package ttprint

import (
	"gradgrind/backend/base"
	"gradgrind/backend/ttbase"
)

func getRooms(
	ttinfo *ttbase.TtInfo,
	pagemap map[Ref][]xPage,
) []ttPage {
	data := getRoomData(ttinfo)
	pages := []ttPage{}
	for _, e := range ttinfo.Db.Rooms {
		tiles, ok := data[e.Id]
		if !ok {
			continue
		}
		page := ttPage{
			"Short":      e.Tag,
			"Activities": tiles,
		}
		page.extendPage(pagemap[e.Id])
		pages = append(pages, page)
	}
	return pages
}

func getOneRoom(
	ttinfo *ttbase.TtInfo,
	pagemap map[Ref][]xPage,
	e *base.Room,
) []ttPage {
	data := getRoomData(ttinfo)
	tiles, ok := data[e.Id]
	if !ok {
		tiles = []Tile{} // Avoid none in JSON if table empty
	}
	page := ttPage{
		"Short":      e.Tag,
		"Activities": tiles,
	}
	page.extendPage(pagemap[e.Id])
	return []ttPage{page}
}

func getRoomData(ttinfo *ttbase.TtInfo) map[Ref][]Tile {
	db := ttinfo.Db
	// Generate the tiles.
	roomTiles := map[Ref][]Tile{}
	type rdata struct { // for SuperCourses
		groups   map[Ref]bool
		teachers map[Ref]bool
	}

	for cref, cinfo := range ttinfo.CourseInfo {
		subject := ttinfo.Ref2Tag[cinfo.Subject]
		// For SuperCourses gather the resources from the relevant SubCourses.
		sc, ok := db.Elements[cref].(*base.SuperCourse)
		if ok {
			rmap := map[Ref]rdata{}
			for _, subref := range sc.SubCourses {
				sub := db.Elements[subref].(*base.SubCourse)
				if sub.Room != "" {
					rlist := []Ref{}
					r0 := db.Elements[sub.Room]
					r, ok := r0.(*base.Room)
					if ok {
						rlist = append(rlist, r.Id)
					} else {
						r, ok := r0.(*base.RoomGroup)
						if ok {
							rlist = append(rlist, r.Rooms...)
						} else {
							r, ok := r0.(*base.RoomChoiceGroup)
							if ok {
								rlist = append(rlist, r.Rooms...)
							} else {
								base.Report(
									`<Bug>Invalid room in course %s:\n  %+v>`,
									ttinfo.View(cinfo), r0)
								panic("Bug")
							}
						}
					}
					for _, rref := range rlist {
						rdata1, ok := rmap[rref]
						if !ok {
							rdata1 = rdata{
								map[Ref]bool{},
								map[Ref]bool{},
							}
						}
						for _, tref1 := range sub.Teachers {
							rdata1.teachers[tref1] = true
						}
						for _, gref := range sub.Groups {
							rdata1.groups[gref] = true
						}
						rmap[rref] = rdata1
					}
				}
			}

			// The actual rooms are associated with the lessons
			for _, l := range cinfo.Lessons {
				if l.Day < 0 {
					continue
				}
				for _, rref := range l.Rooms {
					tlist := []Ref{}
					glist := []Ref{}
					rdata1, ok := rmap[rref]
					if ok {
						for t := range rdata1.teachers {
							tlist = append(tlist, t)
						}
						for g := range rdata1.groups {
							glist = append(glist, g)
						}
					}
					gstrings := splitGroups(ttinfo.SortList(glist))
					tstrings := ttinfo.SortList(tlist)
					tile := Tile{
						Day:      l.Day,
						Hour:     l.Hour,
						Duration: l.Duration,
						//Fraction: 1,
						//Offset: 0,
						//Total:    1,
						Subject:    subject,
						Groups:     gstrings,
						Teachers:   tstrings,
						Background: l.Background,
						Footnote:   l.Footnote,
					}
					roomTiles[rref] = append(roomTiles[rref], tile)
				}
			}
		} else {
			// A normal Course
			glist := []Ref{}
			glist = append(glist, cinfo.Groups...)
			gstrings := splitGroups(ttinfo.SortList(glist))
			tlist := []Ref{}
			tlist = append(tlist, cinfo.Teachers...)
			tstrings := ttinfo.SortList(tlist)

			// The rooms are associated with the lessons
			for _, l := range cinfo.Lessons {
				if l.Day < 0 {
					continue
				}
				for _, rref := range l.Rooms {
					tile := Tile{
						Day:      l.Day,
						Hour:     l.Hour,
						Duration: l.Duration,
						//Fraction: 1,
						//Offset:   0,
						//Total:    1,
						Subject:    subject,
						Groups:     gstrings,
						Teachers:   tstrings,
						Background: l.Background,
						Footnote:   l.Footnote,
					}
					roomTiles[rref] = append(roomTiles[rref], tile)
				}
			}
		}
	}
	return roomTiles
}
