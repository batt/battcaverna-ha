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
}
