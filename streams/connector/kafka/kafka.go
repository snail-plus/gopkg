// Copyright 2024 eve.  All rights reserved.

package kafka

import (
	"context"
	"gitee.com/eve_3/gopkg/log"
	"os"
	"os/signal"
	"syscall"

	"github.com/segmentio/kafka-go"

	"gitee.com/eve_3/gopkg/streams"
	"gitee.com/eve_3/gopkg/streams/flow"
)

// KafkaSource represents an Apache Kafka source connector.
type KafkaSource struct {
	r         *kafka.Reader
	out       chan any
	ctx       context.Context
	cancelCtx context.CancelFunc
}

// NewKafkaSource returns a new KafkaSource instance.
func NewKafkaSource(ctx context.Context, config kafka.ReaderConfig) (*KafkaSource, error) {
	out := make(chan any)
	cctx, cancel := context.WithCancel(ctx)

	sink := &KafkaSource{
		r:         kafka.NewReader(config),
		out:       out,
		ctx:       cctx,
		cancelCtx: cancel,
	}

	go sink.init()
	return sink, nil
}

// init starts the main loop.
func (ks *KafkaSource) init() {
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
	go ks.consume()

	select {
	case <-sigchan:
		ks.cancelCtx()
	case <-ks.ctx.Done():
	}

	close(ks.out)
	ks.r.Close()
}

func (ks *KafkaSource) consume() {
	for {
		// the `ReadMessage` method blocks until we receive the next event
		msg, err := ks.r.ReadMessage(ks.ctx)
		if err != nil {
			log.Errorf("Failed to read message, %v", err)
		}
		ks.out <- msg
	}
}

// Via streams data through the given flow.
func (ks *KafkaSource) Via(_flow streams.Flow) streams.Flow {
	flow.DoStream(ks, _flow)
	return _flow
}

// Out returns an output channel for sending data.
func (ks *KafkaSource) Out() <-chan any {
	return ks.out
}

// KafkaSink represents an Apache Kafka sink connector.
type KafkaSink struct {
	ctx context.Context
	w   *kafka.Writer
	in  chan any
}

// NewKafkaSink returns a new KafkaSink instance.
func NewKafkaSink(ctx context.Context, config kafka.WriterConfig) (*KafkaSink, error) {
	sink := &KafkaSink{
		ctx: ctx,
		w:   kafka.NewWriter(config),
		in:  make(chan any),
	}

	go sink.init()
	return sink, nil
}

// init starts the main loop.
func (ks *KafkaSink) init() {
	for msg := range ks.in {
		var km kafka.Message
		switch m := msg.(type) {
		case []byte:
			km.Value = m
		case string:
			km.Value = []byte(m)
		case *kafka.Message:
			km = *m
		default:
			log.Infof("Unsupported message type, message: %v", m)
			continue
		}
		if err := ks.w.WriteMessages(ks.ctx, km); err != nil {
			log.Errorf("Failed to write message， %v", err)
		}
	}

	ks.w.Close()
}

// In returns an input channel for receiving data.
func (ks *KafkaSink) In() chan<- any {
	return ks.in
}
