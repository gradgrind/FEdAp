package ttengine

import (
	"cmp"
	"fedap/ttbase"
	"fmt"
	"slices"
	"time"
)

// The approach used here might be described as "Escalating Radicality". It
// is based on an algorithm something like Simulated Annealing, but it seems
// that the cooling process may not be very helpful. When no further
// improvements can be made, more radical steps are taken, allowing
// increasing worsening of the penalty for a limited period.

// The initial threshold seems to have little effect on the result. Even 0
// seems to produce only a minor deterioration?
const THRESHOLD0 = 1000
const N0 = 1000
const NSTEPS = 1000

// -----------------
const N1 = NSTEPS * NSTEPS
const N2 = N1 / N0

const PENALTY_UNPLACED_ACTIVITY Penalty = 1000

func PlaceLessons(
	ttinfo *ttbase.TtInfo,
	//alist []ttbase.ActivityIndex,
) bool {
	alist := CollectCourseLessons(ttinfo)

	// Seems to improve speed considerably, especially with complex data:
	slices.Reverse(alist)

	// Build a map associating each non-fixed activity with its position
	// in the initial list.
	posmap := make([]int, len(ttinfo.Activities))
	for i, aix := range alist {
		posmap[aix] = i
	}

	var pmon *placementMonitor
	{
		var delta int64 = 7 // This might be a reasonable value?
		pmon = &placementMonitor{
			count:                 delta,
			delta:                 delta,
			added:                 make([]int64, len(ttinfo.Activities)),
			ttinfo:                ttinfo,
			activityPlacementList: posmap,
			unplaced:              alist,
			resourcePenalties:     make([]Penalty, len(ttinfo.Resources)),
			score:                 0,
			pendingPenalties:      map[ttbase.ResourceIndex]Penalty{},
		}
	}
	pmon.initConstraintData()

	// Calculate initial stage 1 penalties
	for r := 0; r < len(ttinfo.Resources); r++ {
		p := pmon.resourcePenalty1(r)
		pmon.resourcePenalties[r] = p
		pmon.score += p
		//fmt.Printf("$ PENALTY %d: %d\n", r, p)
	}

	// Add penalty for unplaced lessons
	fmt.Printf("$ PENALTY %d: %d\n", len(alist),
		pmon.score+PENALTY_UNPLACED_ACTIVITY*Penalty(len(alist)))

	//TODO--
	state0 := pmon.saveState()
	NR := 1
	tsum := 0.0
	var tlist []float64
	for i := NR; i != 0; i-- {
		start := time.Now()

		pmon.basicLoop()

		// calculate the exe time
		elapsed := time.Since(start)
		fmt.Printf("\n#### ELAPSED: %s\n\n", elapsed)
		telapsed := elapsed.Seconds()
		tsum += telapsed
		tlist = append(tlist, telapsed)

		if i != 1 {
			pmon.restoreState(state0)
		}
	}
	tmean := tsum / float64(NR)
	slices.Sort(tlist)
	NR2 := NR / 2
	tmedian := tlist[NR2]
	if NR%2 == 0 {
		tmedian = (tmedian + tlist[NR2+1]) / 2
	}
	fmt.Printf("#+++ MEAN: %.2f s, MEDIAN: %.2f s.\n", tmean, tmedian)
	return false
	//--

	pmon.basicLoop()
	return false
}

func (pmon *placementMonitor) basicLoop() {

	//TODO: This might need to be placed before the call to "basicLoop":
	pmon.initBreakoutData()
	pmon.bestState = pmon.saveState()
	pmon.placeNonColliding(-1)

	//

	var blockslot ttbase.SlotIndex
	var bestscore Penalty
	var bestunplaced int
	var end0 Penalty = 0

	//--
	unplaced := 100000
	radicalcount := 0
	pccount := 0
	var score Penalty = 0
	//maxradicalcount := 0
	//--

	for {
	evaluate:
		//pmon.fullIntegrityCheck()
		//TODO: exit criteria

		if bestscore != pmon.bestState.score ||
			bestunplaced != len(pmon.bestState.unplaced) {
			bestscore = pmon.bestState.score
			bestunplaced = len(pmon.bestState.unplaced)

			//TODO--
			fmt.Printf("NEW SCORE: %d : %d\n", bestunplaced, bestscore)
		}

		if bestunplaced == 0 {
			//return // to exit when all activities have been placed
			if end0 == 0 {
				end0 = bestscore / 4
			}
			if bestscore <= end0 {
				return
			}
		}

		blockslot = -1
		//for len(pmon.unplaced) == 0 { // seems to get stuck ...
		if len(pmon.unplaced) == 0 {
			blockslot = pmon.removeRandomActivity()
			if pmon.placeNonColliding(blockslot) {
				// score improved
				//pmon.printScore("evaluate")
				goto evaluate
			}

			//TODO: Alternative to for loop? But what about best updating?
			if len(pmon.unplaced) == 0 {
				pmon.removeRandomActivity()
			}
		}
		// Get a bit more radical – allow activities to be replaced.
		// Perform just one forced placement, followed by placeNonColliding.
		// An increased penalty may be accepted, depending on a probability
		// function.

		//pmon.printScore("placeConditional")

		if pmon.placeConditional() {

			//--
			if len(pmon.bestState.unplaced) == 0 {
				if score == pmon.bestState.score {
					pccount++
					if pccount%10 == 0 {
						if pccount%1000 == 0 {
							//fmt.Printf(" +++++++++ pccount: %d\n", pccount)
						}
						//TODO: The following is just a bodge, but it seems
						// to improve things a bit ... at least for a bit ...

						// Perhaps I should retain the initial order and use
						// this somehow?

						pmon.restoreState(pmon.bestState)

						//TODO: This is very "random"!
						// Would reorganizing the freed activities according to the number-of-slots
						// criterion help? Why doesn't the "simple" version get past placeConditional()?
						nc := pccount % 10
						for c := 0; c < nc; c++ {
							pmon.removeRandomActivity()
						}
						if nc > 1 {
							//slices.Reverse(pmon.unplaced)
							slices.SortFunc(pmon.unplaced,
								func(a, b ttbase.ActivityIndex) int {
									return cmp.Compare(
										pmon.activityPlacementList[a],
										pmon.activityPlacementList[b])
								})
						}
					}
				} else {
					score = pmon.bestState.score
					pccount = 0
				}
			}

			//--

			continue
		}

		//TODO: This doesn't seem to be a good place? 1) rarely "successful",
		// 2) seems ineffective
		// Test whether the best score has been beaten.
		lcur := len(pmon.unplaced)
		lbest := len(pmon.bestState.unplaced)
		if lcur < lbest || (lcur == lbest && pmon.score < pmon.bestState.score) {
			pmon.bestState = pmon.saveState()
			pmon.printScore("Better")
			//--panic("TODO--")
		}

		//--
		if len(pmon.unplaced) < unplaced {
			//fmt.Printf("*** UNPLACED: %d N: %d\n", unplaced, radicalcount)
			unplaced = len(pmon.unplaced)
			radicalcount = 0
		}
		radicalcount++
		if radicalcount%10 == 0 {
			fmt.Printf("********************** RADICAL: %d N: %d\n",
				unplaced, radicalcount)
		}

		/*
			if radicalcount%20 == 0 {
				pmon.removeRandomActivity()
				if radicalcount%100 == 0 {
					pmon.removeRandomActivity()
					pmon.removeRandomActivity()
					pmon.removeRandomActivity()
					if radicalcount%1000 == 0 {
						fmt.Printf("********************** RADICAL: %d N: %d\n",
							unplaced, radicalcount)
						if radicalcount > maxradicalcount {
							maxradicalcount = radicalcount
							fmt.Println("     +++ NEW MAXIMUM +++")
							fmt.Printf(" -- SCORE: %d : %d\n",
								len(pmon.bestState.unplaced), pmon.bestState.score)
							fmt.Println("**********************")
						}
						pmon.restoreState(pmon.bestState)
						pmon.removeRandomActivity()
						pmon.removeRandomActivity()
						pmon.removeRandomActivity()
						pmon.removeRandomActivity()
						pmon.removeRandomActivity()
						pmon.removeRandomActivity()
						goto evaluate
					}
				}
			}
			//--

		*/

		// Mechanism to escape to other solution areas:
		// Accept a worsening step and follow its progress a while to see
		// if a better solution area can be found.

		// If no improvement has been made, go to the next level, if there is
		// one. If not, choose the next possibility from this level. Once
		// these have all failed to bring an improvement, end this level.

		if !pmon.radicalStep() {
			//TODO: Revert to bestState?
			pmon.restoreState(pmon.bestState)

			pmon.removeRandomActivity() // doesn't help? Does on data x01!
			// Indeed, there it is important.
		}
	}
}

func (pmon *placementMonitor) printScore(msg string) {
	var p Penalty = 0
	for r := 0; r < len(pmon.ttinfo.Resources); r++ {
		p += pmon.resourcePenalty1(r)
	}
	fmt.Printf("§ Score: %s %d\n", msg,
		pmon.score+Penalty(len(pmon.unplaced))*PENALTY_UNPLACED_ACTIVITY)
	if p != pmon.score {
		fmt.Printf("§ ... error: %d != %d\n", p, pmon.score)
		panic("!!!")
	}
}

func (pmon *placementMonitor) printStateScore(msg string, state *ttState) {
	fmt.Printf("§ StateScore: %s %d\n", msg,
		state.score+Penalty(len(state.unplaced))*PENALTY_UNPLACED_ACTIVITY)
}
