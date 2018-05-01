package drivers

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

type pinMock struct {
	name string
	val  bool
	r    io.Reader
	w    io.Writer
}

func NewPinMock(name string, r io.Reader, w io.Writer) *pinMock {
	p := &pinMock{name, false, r, w}
	return p
}

var stream []byte

func (p *pinMock) Value() bool {
	if p.r != nil {
		var b [1]byte
		p.r.Read(b[:])
		p.val = true
		if b[0] == '0' {
			p.val = false
		}
	}
	return p.val
}

func (p *pinMock) SetValue(v bool) {
	p.val = v
	val := "0"
	if v {
		val = "1"
	}
	p.w.Write([]byte(fmt.Sprint(p.name, val)))
}

func TestSipo(t *testing.T) {
	var wbuf []byte
	buffer := bytes.NewBuffer(wbuf)
	var rbuf []byte
	rbuffer := bytes.NewBuffer(rbuf)

	clk := NewPinMock("CK", nil, buffer)
	miso := NewPinMock("MI", rbuffer, nil)
	mosi := NewPinMock("MO", nil, buffer)
	load := NewPinMock("LD", nil, buffer)
	s := NewSipo(clk, miso, mosi, load)
	check := "CK0MO0LD0"
	if buffer.String() != check {
		t.Fatalf("want %v, got %v\n", check, buffer)
	}
	L := "MO0CK1CK0"
	H := "MO1CK1CK0"
	LOAD := "LD1LD0"
	//send 0x00
	check += LOAD
	check += L + L + L + L + L + L + L + L
	check += LOAD

	// read 0x00
	rbuffer.WriteString("00000000")
	in := s.TransferByte(0x00)
	if buffer.String() != check {
		t.Fatalf("want %v, got %v\n", check, buffer)
	}

	if in != 0x00 {
		t.Fatalf("want %v, got %v\n", 0, in)
	}

	//send 0xff
	check += LOAD
	check += H + H + H + H + H + H + H + H
	check += LOAD

	//read 0xff
	rbuffer.WriteString("11111111")
	in = s.TransferByte(0xFF)
	if buffer.String() != check {
		t.Fatalf("want %v, got %v\n", check, buffer)
	}

	if in != 0xff {
		t.Fatalf("want %v, got %v\n", 0xff, in)
	}

	//send 0x55
	check += LOAD
	check += H + L + H + L + H + L + H + L
	check += LOAD

	//read 0x55
	rbuffer.WriteString("10101010")
	in = s.TransferByte(0x55)
	if buffer.String() != check {
		t.Fatalf("want %v, got %v\n", check, buffer)
	}

	if in != 0x55 {
		t.Fatalf("want %v, got %v\n", 0x55, in)
	}

	//send 0x23 0x34
	check3 := []byte{0x23, 0x34}
	check += LOAD
	check += H + H + L + L + L + H + L + L
	check += L + L + H + L + H + H + L + L
	check += LOAD

	// read LSB first 0x23 0x34 -> 00100011(0x23) 00110100(0x34)
	// NOTE: LSB first!!
	rbuffer.WriteString("1100010000101100")

	in3 := s.Transfer(check3)
	if buffer.String() != check {
		t.Fatalf("want %v, got %v\n", check, buffer)
	}

	if !bytes.Equal(in3, check3) {
		t.Fatalf("want %v, got %v\n", check3, in3)
	}
}
