package power_hpc

import (
	"encoding/binary"
	"fmt"
	. "github.com/influxdata/influxdb-comparisons/bulk_data_gen/common"
	"log"
	"math/rand"
	"os"
	"time"
)

type PowerHpcSimulator struct {
	madePoints int64
	madeValues int64
	maxPoints  int64

	hosts     []Host
	hostIndex int
}

func (g *PowerHpcSimulator) SeenPoints() int64 {
	return g.madePoints
}

func (g *PowerHpcSimulator) SeenValues() int64 {
	return g.madeValues
}

func (g *PowerHpcSimulator) Total() int64 {
	return g.maxPoints
}

func (g *PowerHpcSimulator) Finished() bool {
	return g.madePoints >= g.maxPoints
}

type PowerHpcSimulatorConfig struct {
	Start time.Time
	End   time.Time

	SamplingRate int64
	ChannelCount int64
	HostCount    int64
}

func (d *PowerHpcSimulatorConfig) ToSimulator() *PowerHpcSimulator {
	channelNames := make([][]byte, d.ChannelCount)
	for i := 0; i < len(channelNames); i++ {
		channelNames[i] = []byte(fmt.Sprintf("socket_%d", i))
	}

	interval := time.Duration(int64(time.Second.Nanoseconds() / d.SamplingRate))
	duration := d.End.Sub(d.Start)
	maxPoints := (duration.Nanoseconds() / interval.Nanoseconds()) * d.ChannelCount * d.HostCount
	log.Printf("max points: %d, channels: %d, hosts: %d", maxPoints, d.ChannelCount, d.HostCount)

	file, err := os.Open("/home/tilsche/metricq/hta_lmg_4k.bin")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fi, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	fakeValues := make([]float64, fi.Size()/8)
	if err := binary.Read(file, binary.LittleEndian, &fakeValues); err != nil {
		log.Fatal(err)
	}

	hosts := make([]Host, d.HostCount)
	for i := 0; i < len(hosts); i++ {
		offset := time.Duration(rand.Intn(int(interval.Nanoseconds())))
		timestampStart := d.Start.Add(offset)
		timestampEnd := d.End.Add(offset)
		name := []byte(fmt.Sprintf("h%d", i))

		fakeValueIndices := make([]int, d.ChannelCount)
		for i := 0; i < int(d.ChannelCount); i++ {
			fakeValueIndices[i] = rand.Intn(len(fakeValues))
		}
		hosts[i] = Host{
			channels:         channelNames,
			interval:         interval,
			timestampStart:   timestampStart,
			timestampNow:     timestampStart,
			timestampEnd:     timestampEnd,
			hostname:         name,
			fakeValues:       fakeValues,
			fakeValueIndices: fakeValueIndices,
		}
	}

	sim := &PowerHpcSimulator{
		madePoints: 0,
		madeValues: 0,
		maxPoints:  maxPoints,
		hosts:      hosts,
		hostIndex:  0,
	}
	log.Printf("Actual power sampling interval %v\n", interval)
	return sim
}

// Next advances a Point to the next state in the generator.
func (d *PowerHpcSimulator) Next(p *Point) {
	host := &d.hosts[d.hostIndex]
	d.hostIndex++
	if d.hostIndex == len(d.hosts) {
		d.hostIndex = 0
	}

	host.Next(p)

	d.madePoints++
	d.madeValues += int64(len(host.channels))
	return
}
