package controller

import (
	"fmt"
	"io"
	"log"

	"github.com/batt/battcaverna-ha/devices"
)

type deviceEntry struct {
	dev  devices.Device
	tags []string
}

type Controller struct {
	devices []*deviceEntry
	read    io.Reader
	write   io.Writer
}

func NewController(r io.Reader, w io.Writer) *Controller {
	return &Controller{read: r, write: w}
}

func (c *Controller) RegisterDevice(dev devices.Device, tags []string) {
	newdev := &deviceEntry{dev: dev, tags: tags}

	err := dev.Setup()
	if err != nil {
		log.Println("Error in device setup", err)
	}
	c.devices = append(c.devices, newdev)

}

type deviceMsg struct {
	msgType string `json:"type"`
}

func (c *Controller) Run() {

	for _, d := range c.devices {
		go func() {
			for {
				err := d.dev.Loop(c.write)
				if err != nil {
					log.Println("Error in executing dev", d.dev, err)
				}
			}
		}()
	}

	buf := make([]byte, 2048)
	for {
		n, err := c.read.Read(buf)
		if err != nil {
			log.Println("Error reading", err)
		}

		fmt.Println(n, "bytes received", string(buf[:n]))

		/*
			var msg deviceMsg
			err = json.Unmarshal(buf, &msg)
			fmt.println(msg)
		*/
	}
}
