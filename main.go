package main

import (
	"fmt"
	"os"
	"sync"
	"time"
)

type fork struct {
	forkMu sync.Mutex
}

type philo struct {
	id        int64
	isDead    bool
	leftFork  fork
	rightFork fork
}

type simulation struct {
	philos        []philo
	forks         []fork
	numPhilos     int64
	deathTime     time.Duration
	eatTime       time.Duration
	sleepTime     time.Duration
	numEatTimes   int64
	simuStartTime time.Time
	simuEndTime   time.Time
}

func newSimulation(numPhilos int64, deathTime, eatTime, sleepTime time.Duration, numEatTimes int64) *simulation {
	simu := simulation{
		philos:      make([]philo, numPhilos),
		forks:       make([]fork, numPhilos),
		numPhilos:   numPhilos,
		deathTime:   deathTime,
		eatTime:     eatTime,
		sleepTime:   sleepTime,
		numEatTimes: numEatTimes,
	}
	return &simu
}

func main() {
	if len(os.Args) < 5 || len(os.Args) > 6 {
		fmt.Println("usage: ./goPhilosophers <numPhilos> <deathTime> <eatTime> <sleepTime> [numEatTimes]")
		os.Exit(1)
	}
}
