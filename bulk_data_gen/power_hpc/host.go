package power_hpc

import (
	. "github.com/influxdata/influxdb-comparisons/bulk_data_gen/common"
	"time"
)

type Host struct {
	channels [][]byte

	interval       time.Duration
	timestampNow   time.Time
	timestampStart time.Time
	timestampEnd   time.Time

	hostname         []byte
	fakeValues       []float64
	fakeValueIndices []int
}

func (h *Host) Next(p *Point) {
	p.AppendTag([]byte("Category"), []byte("test"))
	p.AppendTag([]byte("System"), []byte("bench"))
	p.AppendTag([]byte("Host"), h.hostname)

	p.SetMeasurementName([]byte("power"))
	p.SetTimestamp(&h.timestampNow)

	for i, channel := range h.channels {
		p.AppendField(channel, h.fakeValues[h.fakeValueIndices[i]])
		h.fakeValueIndices[i]++
		h.fakeValueIndices[i] %= len(h.fakeValues)
	}

	h.timestampNow = h.timestampNow.Add(h.interval)
	return
}
