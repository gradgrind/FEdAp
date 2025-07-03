# Sketch of my latest placement attempts

## autoplace.go

Collect – `CollectCourseLessons(ttinfo)` in "support.go" – the lessons that need placing and sort them according to a measure of their "placeability", based on the number of slots into which they could be placed `-> alist`.

The idea is that those lessons with fewer available slots should probably be placed first. These should end up at the beginning of the list.

Now, it turns out that I wanted to organize the list as a stack so that deplaced activities could be added as the next items to be processed. So the activities list actually has the supposedly most difficult activities at the end ... (i.e. the list is reversed).

Build a map associating each non-placed activity with its position in the initial list `-> posmap`. TODO: what is this used for?

### func (pmon *placementMonitor) basicLoop()

After some set-up actions, `(pmon *placementMonitor) basicLoop()` in "noncolliding.go" is called:
Try to place the topmost unplaced activity (repeatedly).
Try all possible placements until one is found that doesn't require the removal of another activity. Start searching at a random slot, only testing those in the activity's "PossibleSlots" list.
Repeat until no more activities can be placed. `pmon.bestState` is updated if – and only if – there is an improvement. That means that each position on the journey must be evaluated.
Return true if pmon.bestState has been updated.
