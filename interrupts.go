package main


var (
	interrupts = [5]string{"VBLANK","LCDSTAT","TIMER","SERIAL","JOYPAD"}
)
func (cpu *gameboyCPU) handleInterrupts() {
	//Interrupt service routine
	if cpu.IE != 0 && cpu.IF != 0 {
		cpu.HALT = false
	}

	if !cpu.IME {
		return
	}

	//Handle an interupt if interrupt exists
	if cpu.IE != 0 && cpu.IF != 0 {
		for bit := 0; bit < 5; bit++ {
			if bitSet(cpu.IE,uint8(bit)) && bitSet(cpu.IF,uint8(bit)){
				cpu.PUSH(cpu.PC)
				cpu.IF &^= (1 << bit)
				cpu.PC = uint16(0x0040 + (0x8 * bit))
				cpu.IME = false
				cpu.cycles += 20
				
				if isDebugging {
					cpu.gb.debug.printConsole(interrupts[bit] + "\n","cyan")
				}
			}
		}
	}
}
/*
	The general gist of handling these 5 interrupts are as follows.
	1) Disable interrupt request held in the IF register
	2) Jump to a INT VEC
	3) Disable IME to prevent any more interrupts from being serviced
*/