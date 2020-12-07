package main

var (
	timerIncrementSpeed = [4]int{1024,16,64,256}
)

type timers struct {
	timerClock int //My rudimentary way of dealing with timers hehe

	DIV uint8  //Incremented at a flat rate of 16384 hz
	TIMA uint8 //Incremented at rate specified by Timer control
	TMA uint8  //When TIMA overflows, this data is loaded
	TAC uint8  //Timer control
}

func (timers *timers) handleTimers(){
	//Called at rate of 16384hz to properly handle timers
	//Not the prettiest solution but hey it works
	

	// timers.timerClock++
	// timers.timerClock %= 16384 
}