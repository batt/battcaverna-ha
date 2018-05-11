package drivers

import (
	"time"
)

type PinMover interface {
	SetValue(bool)
	Value() bool
}

type Sipo struct {
	clk  PinMover
	miso PinMover
	mosi PinMover
	load PinMover
}

const (
	halfPeriod  = 250 * time.Microsecond
	clkDefault  = false
	mosiDefault = false
	loadDefault = true
)

func NewSipo(clk, miso, mosi, load PinMover) *Sipo {
	s := &Sipo{clk, miso, mosi, load}
	s.clk.SetValue(clkDefault)
	s.mosi.SetValue(mosiDefault)
	s.load.SetValue(loadDefault)

	return s
}

func (s *Sipo) pulse(p PinMover) {
	curr := p.Value()
	p.SetValue(!curr)
	time.Sleep(halfPeriod)
	p.SetValue(curr)
	time.Sleep(halfPeriod)
}

func (s *Sipo) clkPulse() {
	s.pulse(s.clk)
}

func (s *Sipo) loadPulse() {
	s.pulse(s.load)
}

func (s *Sipo) transferByte(b byte) byte {
	in := byte(0)
	for i := uint(0); i < 8; i++ {
		s.mosi.SetValue(b&0x80 != 0)
		b <<= 1
		in <<= 1
		if s.miso.Value() {
			in |= 1
		}
		s.clkPulse()
	}

	return in
}

func (s *Sipo) TransferByte(b byte) byte {
	s.loadPulse()
	in := s.transferByte(b)
	s.loadPulse()
	return in
}

func (s *Sipo) Transfer(buf []byte) []byte {
	var in []byte
	s.loadPulse()
	for _, b := range buf {
		c := s.transferByte(b)
		in = append(in, c)
	}
	s.loadPulse()
	return in
}
