package bulk_query_gen

import (
	"math"
	"time"
)

type PowerLab interface {
	MinMaxMeanPower(Query, time.Duration)

	Dispatch(int) Query
}

func PowerLabDispatchAll(d PowerLab, iteration int, q Query, scaleVar int) {
	if scaleVar < 0 {
		panic("logic error: bad scalevar")
	}
	d.MinMaxMeanPower(q, time.Millisecond*time.Duration(math.Pow10(scaleVar)))
}
