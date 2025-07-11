package ttengine

import (
	"fedap/ttbase"
	"math/rand/v2"
	"slices"
)

func (pmon *placementMonitor) placeConditional() bool {
	// Force a placement of the next activity if one of the possibilities
	// leads – after a call to "placeNonColliding" – to an improved score.
	// Depending an a probability function a worsened state might be accepted.
	// On failure (non-acceptance), restore entry state and return false.
	// Note that pmon.bestState is not necessarily changed, even if true
	// is returned.
	ttinfo := pmon.ttinfo
	aix := pmon.unplaced[len(pmon.unplaced)-1]
	a := ttinfo.Activities[aix]
	nposs := len(a.PossibleSlots)
	var dpen Penalty
	var dpenx Penalty
	//var dpen1 Penalty
	//i0 := rand.IntN(nposs) // In this version the random starting point
	//i := i0				 // doesn't help.
	i0 := 0
	i := 0

	// Save entry state.
	state0 := pmon.saveState()
	var state1 *ttState
	var ttot Penalty = 0
	var tlist []Penalty
	var alist []*ttState

	//TODO: Initial threshold = ?
	var threshold Penalty = 5

	var clashes []ttbase.ActivityIndex
	for {
		slot := a.PossibleSlots[i]
		clashes = ttinfo.FindClashes(aix, slot)
		if len(clashes) == 0 {
			// Only accept slots where a replacement is necessary.
			goto nextslot
		}
		for _, aixx := range clashes {
			if pmon.doNotRemove(aixx) {
				goto nextslot
			}
		}
		for _, aixx := range clashes {
			ttinfo.UnplaceActivity(aixx)
		}
		dpen = pmon.place(aix, slot)
		for _, aixx := range clashes {
			dpen += pmon.evaluate1(aixx)
		}
		// Update penalty info
		for r, p := range pmon.pendingPenalties {
			pmon.resourcePenalties[r] = p
		}
		pmon.score += dpen
		// Remove from "unplaced" list
		pmon.unplaced = pmon.unplaced[:len(pmon.unplaced)-1]
		// ... and add removed activities
		pmon.unplaced = append(pmon.unplaced, clashes...)

		if pmon.placeNonColliding(-1) {
			// It's possible that a return here can increase the tendency to
			// get stuck with unplaced activities.
			// Perhaps only return when there are none?
			//return true
			//
			if len(pmon.unplaced) == 0 {
				//return true
			}
		}
		// Allow more flexible acceptance.
		dpenx = dpen + PENALTY_UNPLACED_ACTIVITY*Penalty(
			len(pmon.unplaced)-len(state0.unplaced))
		// Decide whether to accept.
		if dpenx < 0 { // <= or < ?
			return true // It seems a bit better with this return, at least
			// in the phase of reducing unplaced activities.
			//dpenx = 0
		}
		{
			dfac := dpenx / threshold
			// The traditional exponential function seems no better,
			// this function may be a little faster?
			t := N1 / (dfac*dfac + N2)
			//t := Penalty(math.Exp(float64(-dfac)) * float64(N0))
			//if t != 0 && Penalty(rand.IntN(N0)) < t {
			//TODO: A different probability of acceptance (< t*N, with N
			// <1 or >1) may be helpful under some circumstances.
			// Perhaps it should be variable?

			//++return true
			//}

			//
			ttot += t
			tlist = append(tlist, ttot)
			alist = append(alist, pmon.saveState())
			//
		}

		// Restore state.
		pmon.restoreState(state0)
	nextslot:
		i += 1
		if i == nposs {
			i = 0
		}
		if i == i0 {
			// All slots have been tested.

			//if state1 == nil {
			if ttot == 0 {
				break
			}
			j, _ := slices.BinarySearch(tlist, Penalty(rand.IntN(int(ttot))))
			state1 = alist[int(j)]

			pmon.restoreState(state1)
			return true
		}
	}
	return false
}
