package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

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

	var out []byte
	for _, d := range os.Args[1:] {
		data, err := strconv.ParseInt(d, 16, 9)
		if err != nil {
			log.Fatalln(err)
		}
		out = append(out, byte(data))
	}

	in := sipo.Transfer(out)
	fmt.Println(in)
}
