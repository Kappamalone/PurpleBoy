package main

var (
	//timerIncrementSpeed = [4]int{1024,16,64,256}
	//These are the amount of timer clock cycles it takes to increment
	//TIMA, depending on the rate specified by TIMA.
	relativeTimerIncreases = [4]int{16,1024,256,64}
)

type timers struct {
	cpu *gameboyCPU

	timerClock int //My rudimentary way of dealing with timers hehe
	DIV uint8  //Incremented at a flat rate of 16384 hz
	TIMA uint8 //Incremented at rate specified by Timer control
	TMA uint8  //When TIMA overflows, this data is loaded
	TAC uint8  //Timer control
}

func initTimers(cpu *gameboyCPU) *timers{
	timers := new(timers)
	timers.cpu = cpu
	
	return timers
}

func (timers *timers) handleTimers(){
	//Called at rate of 273 times per frame (273 * 60 = 16380hz) to properly handle timers
	//Not the prettiest solution but hey it works
	if bitSet(timers.TAC,2) {
		if (timers.timerClock % relativeTimerIncreases[timers.TAC & 0x3]) == 0 {
			if (timers.TIMA + 1) == 0 { //Check if overflow
				timers.TIMA = timers.TMA
				timers.cpu.IF |= 0x4
				timers.cpu.gb.debug.printConsole("Requesting timer interrupt\n","red")
			} else {
				timers.TIMA++
			}
		}
	}
	timers.DIV++
	timers.timerClock++
	timers.timerClock %= 1024
}