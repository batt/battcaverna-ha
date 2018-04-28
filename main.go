package main

import (
	"log"
	"net"

	"github.com/batt/battcaverna-ha/devices"

	"github.com/batt/battcaverna-ha/controller"
)

const (
	maxDatagramSize = 8192
	address         = "224.0.0.5:1701"
)

func main() {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenMulticastUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadBuffer(maxDatagramSize)

	c := controller.NewController(conn)
	d := devices.DummyDevice{}
	c.RegisterDevice(&d, []string{"dummy"})
	c.Run()
}
