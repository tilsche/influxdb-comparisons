package power_lab

import (
	"encoding/binary"
	"fmt"
	. "github.com/influxdata/influxdb-comparisons/bulk_data_gen/common"
	"log"
	"math/rand"
	"os"
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

	fakeValues       []float64
	fakeValueIndices []int
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
	log.Printf("interval: %s, duration: %s, max points: %d, channes: %d", interval, duration, maxPoints, d.ChannelCount)

	file, err := os.Open("/home/tilsche/metricq/hta_lmg_4k.bin")
	if err != nil {
		log.Fatal(err)
	}
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fakeValues := make([]float64, fi.Size()/8)
	if err := binary.Read(file, binary.LittleEndian, &fakeValues); err != nil {
		log.Fatal(err)
	}
	fakeValueIndices := make([]int, d.ChannelCount)
	for i := 0; i < int(d.ChannelCount); i++ {
		fakeValueIndices[i] = rand.Intn(len(fakeValues))
	}

	sim := &PowerLabSimulator{
		madePoints: 0,
		madeValues: 0,
		maxPoints:  maxPoints,

		simulatedMeasurementIndex: 0,

		channels:         channelNames,
		interval:         interval,
		timestampNow:     d.Start,
		timestampStart:   d.Start,
		timestampEnd:     d.End,
		fakeValues:       fakeValues,
		fakeValueIndices: fakeValueIndices,
	}
	log.Printf("Actual power sampling interval %v\n", interval)
	return sim
}

// Next advances a Point to the next state in the generator.
func (d *PowerLabSimulator) Next(p *Point) {
	p.AppendTag([]byte("System"), []byte("ariel"))

	p.SetMeasurementName([]byte("power"))
	p.SetTimestamp(&d.timestampNow)

	for i, channel := range d.channels {
		p.AppendField(channel, d.fakeValues[d.fakeValueIndices[i]])
		d.fakeValueIndices[i]++
		d.fakeValueIndices[i] %= len(d.fakeValues)
	}

	d.madePoints++
	d.madeValues += int64(len(d.channels))
	d.timestampNow = d.timestampNow.Add(d.interval)
	return
}
