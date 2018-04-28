package devices

import "io"

type Device interface {
	Setup() error
	Loop(io.Writer) error
	SetState([]byte) error
	TearDown() error
}
