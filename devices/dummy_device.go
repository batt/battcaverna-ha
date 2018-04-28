package devices

import (
	"io"
	"log"
	"time"
)

type DummyDevice struct {
}

func (d *DummyDevice) Setup() error {
	log.Println("DummyDevice Setup")
	return nil
}

func (d *DummyDevice) Loop(w io.Writer) error {
	log.Println("DummyDevice Loop")
	time.Sleep(1 * time.Second)
	w.Write([]byte("DummyDevice Loop\n"))

	return nil
}

func (d *DummyDevice) SetState(buf []byte) error {
	log.Println("DummyDevice SetState", string(buf))
	return nil
}

func (d *DummyDevice) TearDown() error {
	log.Println("DummyDevice TearDown")
	return nil
}
