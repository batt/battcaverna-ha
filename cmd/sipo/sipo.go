package main

import (
	"github.com/batt/battcaverna-ha/drivers"
)

func main() {
	clk := drivers.NewPin(drivers.PortB, 20, drivers.Out)
	defer clk.Close()
	miso := drivers.NewPin(drivers.PortB, 18, drivers.In)
	defer miso.Close()
	mosi := drivers.NewPin(drivers.PortB, 19, drivers.Out)
	defer mosi.Close()
	load := drivers.NewPin(drivers.PortB, 23, drivers.Out)
	defer load.Close()
	sipo := drivers.NewSipo(clk, miso, mosi, load)

	for i := byte(0); i < 255; i++ {
		sipo.TransferByte(i)
	}
}
