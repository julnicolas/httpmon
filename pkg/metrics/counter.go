package metrics

// Counter is a common metrics type to count discretely occuring events
// If more elaborated Labels are needed, feel free to implement advanced
// counters in your packages or source files.
type Counter struct {
	time   int64
	name   string
	total  float64
	labels map[string]float64
}

func (o Counter) ScrapeTime() int64 {
	return o.time
}

func (o Counter) Name() string {
	return o.name
}

func (o Counter) Total() float64 {
	return o.total
}

func (o Counter) Labels() interface{} {
	return o.labels
}

func (o Counter) TypedLabels() map[string]float64 {
	return o.labels
}
