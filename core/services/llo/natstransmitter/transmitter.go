package natstransmitter

import "github.com/grafana/dskit/services"

type Transmitter interface {
	llotypes.Transmitter
	services.Service
}

var _ Transmitter = &transmitter{}

type transmitter struct{}

func NewTransmitter() Transmitter {
	return &transmitter{}
}

func (t *transmitter) Transmit() {
	// TODO
}
