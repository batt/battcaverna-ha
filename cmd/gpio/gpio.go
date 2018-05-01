package main

import (
	"fmt"
	"time"

	"github.com/batt/battcaverna-ha/drivers"
)

func main() {
	extSir := drivers.NewPin(drivers.PortE, 7, drivers.Out)
	defer extSir.Close()

	mainsOn := drivers.NewPin(drivers.PortE, 9, drivers.In)
	defer mainsOn.Close()

	for i := 0; i < 20; i++ {
		extSir.SetValue(true)
		fmt.Println(mainsOn.Value(), extSir.Value())
		time.Sleep(1 * time.Second)
		extSir.SetValue(false)
		fmt.Println(mainsOn.Value(), extSir.Value())
		time.Sleep(1 * time.Second)
	}

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
