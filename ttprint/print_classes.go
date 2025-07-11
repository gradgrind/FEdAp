package ttprint

import (
	"fedap/base"
	"fedap/ttbase"
	"slices"
)

func getClasses(
	ttinfo *ttbase.TtInfo,
	pagemap map[base.Ref][]xPage,
) []ttPage {
	data := getClassData(ttinfo)
	pages := []ttPage{}
	for _, e := range ttinfo.Db.Classes {
		if e.Tag == "" {
			continue
		}
		tiles, ok := data[e.Id]
		if !ok {
			continue
		}
		page := ttPage{
			"Name":       e.Name,
			"Short":      e.Tag,
			"Activities": tiles,
		}
		page.extendPage(pagemap[e.Id])
		pages = append(pages, page)
	}
	return pages
}

func getOneClass(
	ttinfo *ttbase.TtInfo,
	pagemap map[base.Ref][]xPage,
	e *base.Class,
) []ttPage {
	data := getClassData(ttinfo)
	tiles, ok := data[e.Id]
	if !ok {
		tiles = []Tile{} // Avoid none in JSON if table empty
	}
	page := ttPage{
		"Name":       e.Name,
		"Short":      e.Tag,
		"Activities": tiles,
	}
	page.extendPage(pagemap[e.Id])
	return []ttPage{page}
}

func getClassData(ttinfo *ttbase.TtInfo) map[base.Ref][]Tile {
	db := ttinfo.Db
	ref2id := ttinfo.Ref2Tag
	type cdata struct { // for SuperCourses
		groups   map[base.Ref]bool
		rooms    map[base.Ref]bool
		teachers map[base.Ref]bool
	}
	// Generate the tiles.
	classTiles := map[base.Ref][]Tile{}
	for cref, cinfo := range ttinfo.CourseInfo {
		subject := ref2id[cinfo.Subject]
		// For SuperCourses gather the resources from the relevant SubCourses.
		sc, ok := db.Elements[cref].(*base.SuperCourse)
		if ok {
			cmap := map[base.Ref]cdata{}
			for _, subref := range sc.SubCourses {
				sub := db.Elements[subref].(*base.SubCourse)
				for _, gref := range sub.Groups {
					g := db.Elements[gref].(*base.Group)
					cref := g.Class
					//c := db.Elements[cref].(*base.Class)

					cdata1, ok := cmap[cref]
					if !ok {
						cdata1 = cdata{
							map[base.Ref]bool{},
							map[base.Ref]bool{},
							map[base.Ref]bool{},
						}
					}

					for _, gref1 := range sub.Groups {
						cdata1.groups[gref1] = true
					}
					for _, tref := range sub.Teachers {
						cdata1.teachers[tref] = true
					}
					if sub.Room != "" {
						cdata1.rooms[sub.Room] = true
					}
					cmap[cref] = cdata1
				}
			}
			// The rooms are associated with the lessons
			for _, lix := range cinfo.Lessons {
				l := ttinfo.Activities[lix].Lesson
				if l.Day < 0 {
					continue
				}
				rooms := l.Rooms
				for cref, cdata1 := range cmap {
					tlist := []base.Ref{}
					for t := range cdata1.teachers {
						tlist = append(tlist, t)
					}
					glist := []base.Ref{}
					for g := range cdata1.groups {
						glist = append(glist, g)
					}
					// Choose the rooms in "rooms" which could be relevant for
					// the list of (general) rooms in cdata1.rooms.
					rlist := []base.Ref{}
					for rref := range cdata1.rooms {
						rx := db.Elements[rref]
						_, ok := rx.(*base.Room)
						if ok {
							if slices.Contains(rooms, rref) {
								rlist = append(rlist, rref)
							}
							continue
						}
						rg, ok := rx.(*base.RoomGroup)
						if ok {
							for _, rr := range rg.Rooms {
								if slices.Contains(rooms, rr) {
									rlist = append(rlist, rr)
								}
							}
							continue
						}
						rc, ok := rx.(*base.RoomChoiceGroup)
						if ok {
							for _, rr := range rc.Rooms {
								if slices.Contains(rooms, rr) {
									rlist = append(rlist, rr)
								}
							}
							continue
						}
						base.Bug.Fatalf("Not a room: %s\n", rref)
					}

					// The groups need special handling, to determine tile
					// fractions (with the groups from the current class)
					chips := ttinfo.SortClassGroups(cref, glist)
					tstrings := ttinfo.SortList(tlist)
					rstrings := ttinfo.SortList(rlist)

					for _, chip := range chips {
						gstrings := append(chip.Groups, chip.ExtraGroups...)
						tile := Tile{
							Day:        l.Day,
							Hour:       l.Hour,
							Duration:   l.Duration,
							Fraction:   chip.Fraction,
							Offset:     chip.Offset,
							Total:      chip.Total,
							Subject:    subject,
							Groups:     gstrings,
							Teachers:   tstrings,
							Rooms:      rstrings,
							Background: l.Background,
							Footnote:   l.Footnote,
						}
						classTiles[cref] = append(classTiles[cref], tile)
					}
				}
			}
		} else {
			// A normal Course
			tlist := []base.Ref{}
			tlist = append(tlist, cinfo.Teachers...)
			tstrings := ttinfo.SortList(tlist)

			// The rooms are associated with the lessons
			for _, lix := range cinfo.Lessons {
				l := ttinfo.Activities[lix].Lesson
				if l.Day < 0 {
					continue
				}
				rlist := []base.Ref{}
				rlist = append(rlist, l.Rooms...)
				rstrings := ttinfo.SortList(rlist)

				// The groups need special handling, to determine tile
				// fractions (with the groups from the current class)
				for _, gref := range cinfo.Groups {
					g := db.Elements[gref].(*base.Group)
					cref := g.Class
					chips := ttinfo.SortClassGroups(cref, cinfo.Groups)

					for _, chip := range chips {
						gstrings := append(chip.Groups, chip.ExtraGroups...)
						tile := Tile{
							Day:        l.Day,
							Hour:       l.Hour,
							Duration:   l.Duration,
							Fraction:   chip.Fraction,
							Offset:     chip.Offset,
							Total:      chip.Total,
							Subject:    subject,
							Groups:     gstrings,
							Teachers:   tstrings,
							Rooms:      rstrings,
							Background: l.Background,
							Footnote:   l.Footnote,
						}
						classTiles[cref] = append(classTiles[cref], tile)
					}
				}
			}
		}
	}
	return classTiles
}
