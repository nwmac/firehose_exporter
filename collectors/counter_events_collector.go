package collectors

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/log"

	"github.com/bosh-prometheus/firehose_exporter/metrics"
	"github.com/bosh-prometheus/firehose_exporter/utils"
)

type CounterEventsCollector struct {
	namespace                  string
	environment                string
	metricsStore               *metrics.Store
	counterEventsCollectorDesc *prometheus.Desc
}

func NewCounterEventsCollector(
	namespace string,
	environment string,
	metricsStore *metrics.Store,
) *CounterEventsCollector {
	counterEventsCollectorDesc := prometheus.NewDesc(
		prometheus.BuildFQName(namespace, counter_events_subsystem, "collector"),
		"Cloud Foundry Firehose counter metrics collector.",
		nil,
		prometheus.Labels{"environment": environment},
	)

	return &CounterEventsCollector{
		namespace:                  namespace,
		environment:                environment,
		metricsStore:               metricsStore,
		counterEventsCollectorDesc: counterEventsCollectorDesc,
	}
}

func (c CounterEventsCollector) Collect(ch chan<- prometheus.Metric) {
	for _, counterEvent := range c.metricsStore.GetCounterEvents() {
		metricName := utils.NormalizeName(counterEvent.Origin) + "_" + utils.NormalizeName(counterEvent.Name) + "_total"
		tcm, err := prometheus.NewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(c.namespace, counter_events_subsystem, metricName),
				fmt.Sprintf("Cloud Foundry Firehose '%s' total counter event from '%s'.", utils.NormalizeNameDesc(counterEvent.Name), utils.NormalizeOriginDesc(counterEvent.Origin)),
				[]string{"origin", "bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_ip"},
				prometheus.Labels{"environment": c.environment},
			),
			prometheus.CounterValue,
			float64(counterEvent.Total),
			counterEvent.Origin,
			counterEvent.Deployment,
			counterEvent.Job,
			counterEvent.Index,
			counterEvent.IP,
		)
		if err != nil {
			log.Errorf("Counter Event `%s` from `%s` discarded: %s", counterEvent.Name, counterEvent.Origin, err)
			continue
		}
		ch <- tcm

		metricName = utils.NormalizeName(counterEvent.Origin) + "_" + utils.NormalizeName(counterEvent.Name) + "_delta"
		dcm, err := prometheus.NewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(c.namespace, counter_events_subsystem, metricName),
				fmt.Sprintf("Cloud Foundry Firehose '%s' delta counter event from '%s'.", utils.NormalizeNameDesc(counterEvent.Name), utils.NormalizeOriginDesc(counterEvent.Origin)),
				[]string{"origin", "bosh_deployment", "bosh_job_name", "bosh_job_id", "bosh_job_ip"},
				prometheus.Labels{"environment": c.environment},
			),
			prometheus.GaugeValue,
			float64(counterEvent.Delta),
			counterEvent.Origin,
			counterEvent.Deployment,
			counterEvent.Job,
			counterEvent.Index,
			counterEvent.IP,
		)
		if err != nil {
			log.Errorf("Counter Event `%s` from `%s` discarded: %s", counterEvent.Name, counterEvent.Origin, err)
			continue
		}
		ch <- dcm
	}
}

func (c CounterEventsCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.counterEventsCollectorDesc
}
