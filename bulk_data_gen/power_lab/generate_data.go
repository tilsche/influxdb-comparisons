package power_lab

import (
	"fmt"
	. "github.com/influxdata/influxdb-comparisons/bulk_data_gen/common"
	"log"
	"math/rand"
	"time"
)

type PowerLabSimulator struct {
	madePoints int64
	madeValues int64
	maxPoints  int64

	simulatedMeasurementIndex int

	channels [][]byte

	interval       time.Duration
	timestampNow   time.Time
	timestampStart time.Time
	timestampEnd   time.Time
}

func (g *PowerLabSimulator) SeenPoints() int64 {
	return g.madePoints
}

func (g *PowerLabSimulator) SeenValues() int64 {
	return g.madeValues
}

func (g *PowerLabSimulator) Total() int64 {
	return g.maxPoints
}

func (g *PowerLabSimulator) Finished() bool {
	return g.madePoints >= g.maxPoints
}

type PowerLabSimulatorConfig struct {
	Start time.Time
	End   time.Time

	SamplingRate int64
	ChannelCount int64
}

func (d *PowerLabSimulatorConfig) ToSimulator() *PowerLabSimulator {
	channelNames := make([][]byte, d.ChannelCount)
	for i := 0; i < len(channelNames); i++ {
		channelNames[i] = []byte(fmt.Sprintf("socket_%d", i))
	}

	interval := time.Duration(int64(time.Second.Nanoseconds() / d.SamplingRate))
	duration := d.End.Sub(d.Start)
	maxPoints := duration.Nanoseconds() / interval.Nanoseconds()
	log.Printf("max points: %d, channes: %d", maxPoints, d.ChannelCount)

	sim := &PowerLabSimulator{
		madePoints: 0,
		madeValues: 0,
		maxPoints:  maxPoints,

		simulatedMeasurementIndex: 0,

		channels:       channelNames,
		interval:       interval,
		timestampNow:   d.Start,
		timestampStart: d.Start,
		timestampEnd:   d.End,
	}
	log.Printf("Actual power sampling interval %v\n", interval)
	return sim
}

// Next advances a Point to the next state in the generator.
func (d *PowerLabSimulator) Next(p *Point) {
	p.AppendTag([]byte("System"), []byte("ariel"))

	p.SetMeasurementName([]byte("power"))
	p.SetTimestamp(&d.timestampNow)

	for _, channel := range d.channels {
		p.AppendField(channel, rand.NormFloat64()*20+80)
	}

	d.madePoints++
	d.madeValues += int64(len(d.channels))
	d.timestampNow = d.timestampNow.Add(d.interval)
	return
}
