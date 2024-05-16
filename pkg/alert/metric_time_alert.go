package alert

import (
	"time"

	"github.com/julnicolas/httpmon/pkg/metrics"
)

type MetricsTimeAlert struct {
	BaseAlert
}

func NewMetricsTimeAlert(period time.Duration, rule AlertRule) *MetricsTimeAlert {
	o := &MetricsTimeAlert{
		BaseAlert: *NewBaseAlert(period, rule),
	}
	o.timer = NewMetricsTimer(period)

	return o
}

func (o *MetricsTimeAlert) Eval(m metrics.Metric) State {
	o.timer.(*MetricsTimer).Metric(m)
	return o.BaseAlert.Eval(m)
}

func (o *MetricsTimeAlert) Description() string {
	return "MetricsTimeAlert"
}

func (o *MetricsTimeAlert) Name() NameT {
	return "MetricsTimeAlert"
}

func (o *MetricsTimeAlert) DeepCopy() Alert {
	base := o.BaseAlert.DeepCopy()
	n := new(MetricsTimeAlert)
	n.BaseAlert = *base.(*BaseAlert)

	return n
}
