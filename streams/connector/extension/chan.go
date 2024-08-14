// Copyright 2024 eve.  All rights reserved.

package extension

import (
	"gitlab.evebatterycloud.com/infra/gopkg/streams"
	"gitlab.evebatterycloud.com/infra/gopkg/streams/flow"
)

// ChanSource represents an inbound connector that streams items from a channel.
type ChanSource struct {
	in chan any
}

// NewChanSource returns a new ChanSource instance.
func NewChanSource(in chan any) *ChanSource {
	return &ChanSource{in}
}

// Via streams data through the given flow.
func (cs *ChanSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(cs, _flow)
	return _flow
}

// Out returns an output channel for sending data.
func (cs *ChanSource) Out() <-chan any {
	return cs.in
}

// ChanSink represents an outbound connector that streams items to a channel.
type ChanSink struct {
	Out chan any
}

// NewChanSink returns a new ChanSink instance.
func NewChanSink(out chan any) *ChanSink {
	return &ChanSink{out}
}

// In returns an input channel for receiving data.
func (ch *ChanSink) In() chan<- any {
	return ch.Out
}
