package drivers

import (
	"fmt"
	"log"
	"os"
)

type GPIOPort string

const (
	PortA GPIOPort = "A"
	PortB GPIOPort = "B"
	PortC GPIOPort = "C"
	PortD GPIOPort = "D"
	PortE GPIOPort = "E"
)

var portNumber = map[GPIOPort]int{
	PortA: 0,
	PortB: 32,
	PortC: 64,
	PortD: 96,
	PortE: 128,
}

type Direction string

const (
	In  Direction = "in"
	Out Direction = "out"
)

const gpioPath string = "/sys/class/gpio/"

type Pin struct {
	port GPIOPort
	pin  int
	fp   *os.File
}

func NewPin(port GPIOPort, pin int, d Direction) *Pin {
	p := &Pin{port, pin, nil}
	pindir := gpioPath + "pio" + string(p.port) + fmt.Sprint(p.pin)

	if _, err := os.Stat(pindir); os.IsNotExist(err) {
		/* If pin not already exported */
		f, err := os.OpenFile(gpioPath+"export", os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		// export pin
		_, err = f.WriteString(fmt.Sprint(portNumber[port] + pin))
		if err != nil {
			log.Fatal(err)
		}
	}

	fp, err := os.OpenFile(pindir+"/value", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	p.fp = fp
	p.SetDirection(d)

	return p
}

func (p *Pin) SetDirection(d Direction) {
	f, err := os.OpenFile(gpioPath+"pio"+string(p.port)+fmt.Sprint(p.pin)+"/direction", os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(string(d))
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Pin) SetValue(v bool) {
	val := "0"
	if v {
		val = "1"
	}

	p.fp.Seek(0, os.SEEK_SET)
	_, err := p.fp.WriteString(val)
	if err != nil {
		log.Fatal(err)
	}
}

func (p *Pin) Value() bool {

	p.fp.Seek(0, os.SEEK_SET)

	var buf [1]byte
	_, err := p.fp.Read(buf[:])
	if err != nil {
		log.Fatal(err)
	}

	if buf[0] == '0' {
		return false
	}
	return true
}

func (p *Pin) Close() {
	p.fp.Close()

	f, err := os.OpenFile(gpioPath+"unexport", os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprint(portNumber[p.port] + p.pin))
	if err != nil {
		log.Fatal(err)
	}
}
