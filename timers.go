package main

var (
	
	frequency = [4]int{1024,16,64,256} //In T-cycles
)

type timers struct {
	cpu *gameboyCPU

	clock      int
	divClock   int
	DIV        uint8 //Incremented at a flat rate of 16384 hz
	TIMA       uint8 //Incremented at rate specified by Timer control
	TMA        uint8 //When TIMA overflows, this data is loaded
	TAC        uint8 //Timer control
}

func initTimers(cpu *gameboyCPU) *timers {
	timers := new(timers)
	timers.cpu = cpu
	
	return timers
}

func (timers *timers) tick() {
	if bitSet(timers.TAC,2){
		if timers.clock == frequency[timers.TAC & 0x3] {
			if timers.TIMA == 0xFF {
				timers.TIMA = timers.TMA
				timers.cpu.requestTimer()
			} else {
				timers.TIMA++
			}
			timers.clock = 0
		} 
		timers.clock++
	}

	if timers.divClock == 256 {
		timers.DIV++
		timers.divClock = 0
	}
	timers.divClock++
}
