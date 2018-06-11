package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {

	fmt.Println("Let's go for a walk!")
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go start(wg)

	wg.Wait()
	fmt.Println("Alarm is armed.")

}

func start(wg *sync.WaitGroup) {
	
	wgin := &sync.WaitGroup{}
	wgin.Add(2)
	go countTime("Bob", [2]int{60, 90},"getting ready" , wgin)
	go countTime("Alice", [2]int{60, 90},"getting ready" , wgin)
	wgin.Wait()

	wgin.Add(3)
	fmt.Println("Arming alarm")


	go func() {
		fmt.Println("Alarm is counting down")
		wgin.Done()
	}()
	go countTime("Bob", [2]int{35, 45},"putting on shoes" , wgin)
	go countTime("Alice", [2]int{35, 45},"putting on shoes" , wgin)


	wgin.Wait()
	fmt.Println("Exiting and locking the door")

	wg.Done()

}

func countTime(name string, randRange [2]int, stage string, wg *sync.WaitGroup) {
	fmt.Println(strings.Join([]string{name, "started", stage}, " "))
	t := random(randRange[0], randRange[1])
	stringList := []string{name, "spent", strconv.Itoa(t), "seconds",  stage}
	fmt.Println(strings.Join(stringList, " "))
	wg.Done()

}

func random(int1, int2 int) int {
	randDiff := rand.Intn(int2-int1)
	addSleep()
	return randDiff + int1
}

func addSleep(){
	rand.Seed(time.Now().Unix())
	time.Sleep(time.Second * time.Duration(rand.Intn(5)))
}