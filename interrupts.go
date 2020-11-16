package main

func (cpu *gameboyCPU) ISR() {
	//Interrupt service routine

	interruptEnable := cpu.gb.mmu.readbyte(0xFFFF)
	interruptFlags := cpu.gb.mmu.readbyte(0xFF0F)

	//Handle an interupt if interrupt exists
	if interruptEnable != 0 && interruptFlags != 0 {
		cpu.SP -= 2
		cpu.gb.mmu.writeword(cpu.SP, cpu.PC)

		if bitSet(interruptEnable, 0) && bitSet(interruptFlags, 0) {
			cpu.VBlank()
		} else if bitSet(interruptEnable, 1) && bitSet(interruptFlags, 1) {
			cpu.LCDSTAT()
		} else if bitSet(interruptEnable, 2) && bitSet(interruptFlags, 2) {
			cpu.TIMER()
		} else if bitSet(interruptEnable, 3) && bitSet(interruptFlags, 3) {
			cpu.SERIAL()
		} else if bitSet(interruptEnable, 4) && bitSet(interruptFlags, 4) {
			cpu.JOYPAD()
		}
	}

	cpu.cycles += 20 //Take a total of 5 machine cycles
}

/*
The general gist of handling these 5 interrupts are as follows.
1) Disable interrupt request held in the IF register
2) Jump to a INT VEC
3) Disable IME to prevent any more interrupts from being serviced
*/
func (cpu *gameboyCPU) VBlank() {
	//V-blank
	cpu.gb.mmu.writebyte(0xFF0F, 0xFF-0x01)
	cpu.PC = 0x0040
	cpu.IME = false
	println("VBLANK")
}

func (cpu *gameboyCPU) LCDSTAT() {
	//LCD STAT
	cpu.gb.mmu.writebyte(0xFF0F, 0xFF-0x02)
	cpu.PC = 0x0048
	cpu.IME = false
	println("LCD STAT")
}

func (cpu *gameboyCPU) TIMER() {
	//Timer
	cpu.gb.mmu.writebyte(0xFF0F, 0xFF-0x04)
	cpu.PC = 0x0050
	cpu.IME = false
	println("TIMER")
}

func (cpu *gameboyCPU) SERIAL() {
	//Serial
	cpu.gb.mmu.writebyte(0xFF0F, 0xFF-0x08)
	cpu.PC = 0x0058
	cpu.IME = false
	println("SERIAL")
}

func (cpu *gameboyCPU) JOYPAD() {
	//Joypad
	cpu.gb.mmu.writebyte(0xFF0F, 0xFF-0x10)
	cpu.PC = 0x0060
	cpu.IME = false
	println("JOYPAD")
}
