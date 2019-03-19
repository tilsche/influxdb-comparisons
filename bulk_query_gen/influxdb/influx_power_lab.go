package influxdb

import (
	"fmt"
	bulkQuerygen "github.com/influxdata/influxdb-comparisons/bulk_query_gen"
	"time"
)

type InfluxPowerLab struct {
	InfluxCommon
	queryInterval time.Duration
}

func NewInfluxPowerLab(dbConfig bulkQuerygen.DatabaseConfig, interval bulkQuerygen.TimeInterval, duration time.Duration, scaleVar int) bulkQuerygen.QueryGenerator {

	if _, ok := dbConfig[bulkQuerygen.DatabaseName]; !ok {
		panic("need influx database name")
	}

	return &InfluxPowerLab{
		InfluxCommon:  *newInfluxCommon(InfluxQL, dbConfig[bulkQuerygen.DatabaseName], interval, scaleVar),
		queryInterval: duration,
	}
}

func (d *InfluxPowerLab) Dispatch(i int) bulkQuerygen.Query {
	q := bulkQuerygen.NewHTTPQuery() // from pool
	bulkQuerygen.PowerLabDispatchAll(d, i, q, d.ScaleVar)
	return q
}

func (d *InfluxPowerLab) MinMaxMeanPower(qi bulkQuerygen.Query) {
	interval := d.AllInterval.RandWindow(d.queryInterval)
	groupInterval := d.queryInterval / 1000

	if d.language != InfluxQL {
		panic("Wrong language")
	}

	query := fmt.Sprintf(
		`SELECT max(*), min(*), mean(*)
FROM power
WHERE (Category = 'test' and System = 'bench' and Host = 'c0')
and time >= '%s' and time < '%s'
GROUP by time(%s)`, interval.StartString(), interval.EndString(), groupInterval)
	humanLabel := fmt.Sprintf("InfluxDB min/max/mean, rand %s by %s", d, d.queryInterval, groupInterval)

	q := qi.(*bulkQuerygen.HTTPQuery)
	d.getHttpQuery(humanLabel, interval.StartString(), query, q)
}
