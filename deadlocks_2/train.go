package main

import (
	"sort"
	"sync"
	"time"
)

var (
	controller = sync.Mutex{}
	cond       = sync.NewCond(&controller)
)

func MoveTrain1(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1

		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersections2(train.Id, crossing.Position, crossing.Position+trainLength, crossings)
			}

			back := train.Front - train.TrainLength

			if back == crossing.Position {
				crossing.Intersection.LockedBy = -1
				crossing.Intersection.Mutex.Unlock()
			}
		}

		time.Sleep(30 * time.Millisecond)
	}
}

func MoveTrain2(train *Train, distance int, crossings []*Crossing) {
	for train.Front < distance {
		train.Front += 1

		for _, crossing := range crossings {
			if train.Front == crossing.Position {
				lockIntersections2(train.Id, crossing.Position, crossing.Position+trainLength, crossings)
			}

			back := train.Front - train.TrainLength

			if back == crossing.Position {
				controller.Lock()
				crossing.Intersection.LockedBy = -1
				cond.Broadcast()
				controller.Unlock()
			}
		}

		time.Sleep(30 * time.Millisecond)
	}
}

func allFree(intersectionsToLock []*Intersection) bool {
	for _, intersection := range intersectionsToLock {
		if intersection.LockedBy > 0 {
			return false
		}
	}

	return true
}

func lockIntersections2(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	var intersectionsToLock []*Intersection

	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}

	controller.Lock()

	for !allFree(intersectionsToLock) {
		cond.Wait()
	}

	for _, intersection := range intersectionsToLock {
		intersection.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}

	controller.Unlock()
}

func lockIntersections1(id, reserveStart, reserveEnd int, crossings []*Crossing) {
	var intersectionsToLock []*Intersection

	for _, crossing := range crossings {
		if reserveEnd >= crossing.Position && reserveStart <= crossing.Position && crossing.Intersection.LockedBy != id {
			intersectionsToLock = append(intersectionsToLock, crossing.Intersection)
		}
	}

	sort.Slice(intersectionsToLock, func(i, j int) bool {
		return intersectionsToLock[i].Id < intersectionsToLock[j].Id
	})

	for _, intersection := range intersectionsToLock {
		intersection.Mutex.Lock()
		intersection.LockedBy = id
		time.Sleep(10 * time.Millisecond)
	}
}
