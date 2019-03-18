package bulk_query_gen

type PowerLab interface {
	MinMaxMeanPower(Query)

	Dispatch(int) Query
}

func PowerLabDispatchAll(d PowerLab, iteration int, q Query, scaleVar int) {
	if scaleVar < 0 {
		panic("logic error: bad scalevar")
	}
	d.MinMaxMeanPower(q)
}
