package config

import "time"

type Alert struct {
	// Period is the alert loop evaluation period
	// -> alerts are evaluated after every Period time
	Period            time.Duration
	RequestsPerSecond RequestsPerSecond
}

func (o Alert) Default() Alert {
	return Alert{
		Period:            time.Second,
		RequestsPerSecond: RequestsPerSecond{}.Default(),
	}
}

// RequestPerSecond is struct to configure the RequestsPerSecond alert
type RequestsPerSecond struct {
	// Period of time over which the alert rule must be true so that the alert
	// becomes active
	Period time.Duration
	// Threshold of requests per second, if greater the alert is on
	Threshold float64
}

func (o RequestsPerSecond) Default() RequestsPerSecond {
	return RequestsPerSecond{
		// The alert is pending for 1 minute
		// then becomes active (the alert is logged at this
		// time on the UI) for 1 minute.
		// Then it goes back at the inactive state. This state
		// being logged on the UI.
		//
		// In conclusion, the system has detected the alert for 2 minutes
		// In realtime pending alerts are displayed on the left pane of the
		// alert dashboard.
		Period:    time.Minute,
		Threshold: 10,
	}
}
