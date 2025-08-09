package timetable

import (
	"fedap/base"
	"fmt"
	"strings"
)

const ATOMIC_GROUP_SEP1 = "#"
const ATOMIC_GROUP_SEP2 = "~"

type ClassDivision struct {
	Class     *base.Class
	Divisions [][]NodeRef
}

// Prepare filtered versions of the class Divisions containing only
// those Divisions which have Groups used in Lessons.
func FilterDivisions(
	db *base.DbTopLevel,
	course_info []*CourseInfo,
) []ClassDivision {
	// Collect groups used in activities, using CourseInfo structures.
	usedgroups := map[NodeRef]bool{}
	for _, cinfo := range course_info {
		for _, g := range cinfo.Groups {
			usedgroups[g] = true
		}
	}
	// Filter the class divisions, discarding the division names.
	cdivs := []ClassDivision{}
	for _, c := range db.Classes {
		divs := [][]NodeRef{}
		for _, div := range c.Divisions {
			for _, gref := range div.Groups {
				if usedgroups[gref] {
					divs = append(divs, div.Groups)
					break
				}
			}
		}
		cdivs = append(cdivs, ClassDivision{c, divs})
	}
	return cdivs
}

type AtomicGroup struct {
	Index  ResourceIndex
	Class  NodeRef
	Groups []NodeRef
	Tag    string // A constructed tag to represent the atomic group
}

func MakeAtomicGroups(
	db *base.DbTopLevel,
	class_divisions []ClassDivision,
) (map[NodeRef][]*AtomicGroup, int) {
	// An atomic group is an ordered list of single groups, one from each
	// division.
	atomicGroups := map[NodeRef][]*AtomicGroup{}
	atomicGroupIndex := 0

	// Go through the classes inspecting their Divisions.
	// Build a list-basis for the atomic groups based on the Cartesian product.
	for _, cdivs := range class_divisions {
		cl := cdivs.Class
		if len(cdivs.Divisions) == 0 {
			// Make an atomic group for the class
			cag := &AtomicGroup{
				Index: atomicGroupIndex,
				Class: cl.Id,
				Tag:   cl.Tag + ATOMIC_GROUP_SEP1,
			}
			atomicGroupIndex++
			atomicGroups[cl.ClassGroup] = []*AtomicGroup{cag}
			continue
		}

		// The atomic groups will be built as a list of lists of Refs.
		agrefs := [][]NodeRef{{}}
		for _, dglist := range cdivs.Divisions {
			// Add another division – increases underlying list lengths.
			agrefsx := [][]NodeRef{}
			for _, ag := range agrefs {
				// Extend each of the old list items by appending each
				// group of the new division in turn – multiplies the
				// total number of atomic groups.
				for _, g := range dglist {
					gx := make([]NodeRef, len(ag)+1)
					copy(gx, append(ag, g))
					agrefsx = append(agrefsx, gx)
				}
			}
			agrefs = agrefsx
		}
		//fmt.Printf("  §§§ Divisions in %s: %+v\n", cl.Tag, divs)
		//fmt.Printf("     --> %+v\n", agrefs)

		// Make AtomicGroups
		aglist := []*AtomicGroup{}
		for _, ag := range agrefs {
			glist := []string{}
			for _, gref := range ag {
				gtag := db.Elements[gref].(*base.Group).Tag
				glist = append(glist, gtag)
			}
			ago := &AtomicGroup{
				Index:  atomicGroupIndex,
				Class:  cl.Id,
				Groups: ag,
				Tag: cl.Tag + ATOMIC_GROUP_SEP1 +
					strings.Join(glist, ATOMIC_GROUP_SEP2),
			}
			atomicGroupIndex++
			aglist = append(aglist, ago)
		}

		// Map the individual groups to their atomic groups.
		g2ags := map[NodeRef][]*AtomicGroup{}
		count := 1
		divIndex := len(cdivs.Divisions)
		for divIndex > 0 {
			divIndex--
			divGroups := cdivs.Divisions[divIndex]
			agi := 0 // ag index
			for agi < len(aglist) {
				for _, g := range divGroups {
					for j := 0; j < count; j++ {
						g2ags[g] = append(g2ags[g], aglist[agi])
						agi++
					}
				}
			}
			count *= len(divGroups)
		}

		atomicGroups[cl.ClassGroup] = aglist
		for g, agl := range g2ags {
			atomicGroups[g] = agl
		}
	}
	return atomicGroups, atomicGroupIndex
}

// For testing
func (ttinfo *TtInfo) PrintAtomicGroups() {
	for _, cl := range ttinfo.Db.Classes {
		agls := []string{}
		for _, ag := range ttinfo.AtomicGroups[cl.ClassGroup] {
			agls = append(agls, ag.Tag)
		}
		fmt.Printf("  ++ %s: %+v\n", ttinfo.Ref2Tag[cl.ClassGroup], agls)
		for _, div := range ttinfo.ClassDivisions[cl.Id] {
			for _, g := range div {
				agls := []string{}
				for _, ag := range ttinfo.AtomicGroups[g] {
					agls = append(agls, ag.Tag)
				}
				fmt.Printf("    -- %s: %+v\n", ttinfo.Ref2Tag[g], agls)
			}
		}
	}
}
