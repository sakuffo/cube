// Definition of what a scheduler does
// Determine a set of candidate workers on which a task could run
// Score the candidate workers from best to worst
// Pick the worker with the best score

package scheduler

type Scheduler interface {
	SelectCandidateNodes()
	Score()
	Pick()
}