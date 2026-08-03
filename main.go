package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"
)

type fork struct {
	forkMu sync.Mutex
}

type philo struct {
	id            int64
	leftFork      *fork
	rightFork     *fork
	numTimesEated int64
	lastTimeEat   time.Time
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
	stopOnce      sync.Once
}

func newSimulation() *simulation {
	simu := simulation{}
	return &simu
}

func philoLife(simu *simulation, currPhilo *philo, stop chan bool) {
	if currPhilo.id%2 != 0 {
		time.Sleep(100 * time.Millisecond)
	}

	currPhilo.lastTimeEat = time.Now()
	for {
		select {
		case <-stop:
			return
		default:
		}

		timeSinceStart := time.Since(simu.simuStartTime).Milliseconds()
		if time.Since(currPhilo.lastTimeEat) > simu.deathTime {
			fmt.Printf("%d [philo %d] died\n", timeSinceStart, currPhilo.id)
			simu.stopOnce.Do(func() {
				close(stop)
			})
			break
		}

		currPhilo.leftFork.forkMu.Lock()
		fmt.Printf("%d [philo %d] has taken a fork\n", timeSinceStart, currPhilo.id)
		currPhilo.rightFork.forkMu.Lock()
		fmt.Printf("%d [philo %d] has taken a fork\n", timeSinceStart, currPhilo.id)
		fmt.Printf("%d [philo %d] eat\n", timeSinceStart, currPhilo.id)
		select {
		case <-stop:
			currPhilo.leftFork.forkMu.Unlock()
			currPhilo.rightFork.forkMu.Unlock()
			return
		case <-time.After(simu.eatTime):
		}
		currPhilo.leftFork.forkMu.Unlock()
		currPhilo.rightFork.forkMu.Unlock()
		currPhilo.lastTimeEat = time.Now()
		currPhilo.numTimesEated++

		if currPhilo.numTimesEated == simu.numEatTimes {
			break
		}

		fmt.Printf("%d [philo %d] think\n", timeSinceStart, currPhilo.id)

		fmt.Printf("%d [philo %d] sleep\n", timeSinceStart, currPhilo.id)
		select {
		case <-stop:
			return
		case <-time.After(simu.sleepTime):
		}
	}
}

func main() {
	if len(os.Args) < 5 || len(os.Args) > 6 {
		fmt.Println("usage: ./goPhilosophers <numPhilos> <deathTime> <eatTime> <sleepTime> [numEatTimes]")
		os.Exit(1)
	}

	simu := newSimulation()
	var err error

	simu.numPhilos, err = strconv.ParseInt(os.Args[1], 10, 64)
	if err != nil {
		fmt.Println("error:numPhilos:", err)
		os.Exit(1)
	}

	deathMs, err := strconv.ParseInt(os.Args[2], 10, 64)
	if err != nil {
		fmt.Println("error:deathTime:", err)
		os.Exit(1)
	}
	simu.deathTime = time.Duration(deathMs) * time.Millisecond

	eatMs, err := strconv.ParseInt(os.Args[3], 10, 64)
	if err != nil {
		fmt.Println("error:eatTime:", err)
		os.Exit(1)
	}
	simu.eatTime = time.Duration(eatMs) * time.Millisecond

	sleepMs, err := strconv.ParseInt(os.Args[4], 10, 64)
	if err != nil {
		fmt.Println("error:sleepTime:", err)
		os.Exit(1)
	}
	simu.sleepTime = time.Duration(sleepMs) * time.Millisecond

	if len(os.Args) == 6 {
		simu.numEatTimes, err = strconv.ParseInt(os.Args[5], 10, 64)
		if err != nil {
			fmt.Println("error:numEatTimes:", err)
			os.Exit(1)
		}
	}

	simu.philos = make([]philo, simu.numPhilos)
	simu.forks = make([]fork, simu.numPhilos)

	var wg sync.WaitGroup

	simu.simuStartTime = time.Now()

	stop := make(chan bool)
	for i := range simu.philos {
		currPhilo := &simu.philos[i]

		currPhilo.id = int64(i)
		if currPhilo.id%2 == 0 {
			currPhilo.leftFork = &simu.forks[currPhilo.id]
			currPhilo.rightFork = &simu.forks[(currPhilo.id+1)%simu.numPhilos]
		} else {
			currPhilo.rightFork = &simu.forks[currPhilo.id]
			currPhilo.leftFork = &simu.forks[(currPhilo.id+1)%simu.numPhilos]
		}

		wg.Go(func() {
			philoLife(simu, currPhilo, stop)
		})
	}
	wg.Wait()
}
