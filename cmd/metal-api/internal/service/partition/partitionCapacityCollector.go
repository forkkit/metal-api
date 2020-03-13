package partition

import (
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
)

var (
	capacityTotalDesc = prometheus.NewDesc(
		"metal_partition_capacity_total",
		"The total capacity of machines in the partition",
		[]string{"partition", "size"}, nil,
	)
	capacityFreeDesc = prometheus.NewDesc(
		"metal_partition_capacity_free",
		"The capacity of free machines in the partition",
		[]string{"partition", "size"}, nil,
	)
	capacityAllocatedDesc = prometheus.NewDesc(
		"metal_partition_capacity_allocated",
		"The capacity of allocated machines in the partition",
		[]string{"partition", "size"}, nil,
	)
	capacityFaultyDesc = prometheus.NewDesc(
		"metal_partition_capacity_faulty",
		"The capacity of faulty machines in the partition",
		[]string{"partition", "size"}, nil,
	)
)

// partitionCapacityCollector implements the Collector interface.
type partitionCapacityCollector struct {
	r *partitionResource
}

func (pcc partitionCapacityCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(pcc, ch)
}

func (pcc partitionCapacityCollector) Collect(ch chan<- prometheus.Metric) {
	pcs, err := pcc.r.calcPartitionCapacities()
	if err != nil {
		zapup.MustRootLogger().Error("Failed to get partition capacity", zap.Error(err))
		return
	}

	for _, pc := range pcs {
		for _, sc := range pc.ServerCapacities {
			metric, err := prometheus.NewConstMetric(
				capacityTotalDesc,
				prometheus.CounterValue,
				float64(sc.Total),
				pc.Meta.Id,
				sc.Size,
			)
			if err != nil {
				zapup.MustRootLogger().Error("Failed to create metric for totalCapacity", zap.Error(err))
				return
			}
			ch <- metric

			metric, err = prometheus.NewConstMetric(
				capacityFreeDesc,
				prometheus.CounterValue,
				float64(sc.Free),
				pc.Meta.Id,
				sc.Size,
			)
			if err != nil {
				zapup.MustRootLogger().Error("Failed to create metric for freeCapacity", zap.Error(err))
				return
			}
			ch <- metric
			metric, err = prometheus.NewConstMetric(
				capacityAllocatedDesc,
				prometheus.CounterValue,
				float64(sc.Allocated),
				pc.Meta.Id,
				sc.Size,
			)
			if err != nil {
				zapup.MustRootLogger().Error("Failed to create metric for allocatedCapacity", zap.Error(err))
				return
			}
			ch <- metric
			metric, err = prometheus.NewConstMetric(
				capacityFaultyDesc,
				prometheus.CounterValue,
				float64(sc.Faulty),
				pc.Meta.Id,
				sc.Size,
			)
			if err != nil {
				zapup.MustRootLogger().Error("Failed to create metric for faultyCapacity", zap.Error(err))
				return
			}
			ch <- metric
		}
	}
}
