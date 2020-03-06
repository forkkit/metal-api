package network

import (
	"fmt"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/service/helper"
	"github.com/metal-stack/metal-lib/zapup"
	"github.com/prometheus/client_golang/prometheus"
	"go.uber.org/zap"
	"strconv"
	"strings"
)

var (
	usedIpsDesc = prometheus.NewDesc(
		"metal_network_ip_used",
		"The total number of used IPs of the network",
		[]string{"networkId", "prefixes", "destPrefixes", "partitionId", "projectId", "parentNetworkID", "vrf", "isPrivateSuper", "useNat", "isUnderlay"}, nil,
	)
	availableIpsDesc = prometheus.NewDesc(
		"metal_network_ip_available",
		"The total number of available IPs of the network",
		[]string{"networkId", "prefixes", "destPrefixes", "partitionId", "projectId", "parentNetworkID", "vrf", "isPrivateSuper", "useNat", "isUnderlay"}, nil,
	)
	usedPrefixesDesc = prometheus.NewDesc(
		"metal_network_prefix_used",
		"The total number of used prefixes of the network",
		[]string{"networkId", "prefixes", "destPrefixes", "partitionId", "projectId", "parentNetworkID", "vrf", "isPrivateSuper", "useNat", "isUnderlay"}, nil,
	)
	availablePrefixesDesc = prometheus.NewDesc(
		"metal_network_prefix_available",
		"The total number of available prefixes of the network",
		[]string{"networkId", "prefixes", "destPrefixes", "partitionId", "projectId", "parentNetworkID", "vrf", "isPrivateSuper", "useNat", "isUnderlay"}, nil,
	)
)

// networkUsageCollector implements the prometheus collector interface.
type networkUsageCollector struct {
	r *networkResource
}

func (nuc networkUsageCollector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(nuc, ch)
}

func (nuc networkUsageCollector) Collect(ch chan<- prometheus.Metric) {
	// FIXME bad workaround to be able to run make spec
	if nuc.r == nil || nuc.r.ds == nil {
		return
	}
	nws, err := nuc.r.ds.ListNetworks()
	if err != nil {
		zapup.MustRootLogger().Error("Failed to get network usage", zap.Error(err))
		return
	}

	for i := range nws {
		usage := helper.GetNetworkUsage(&nws[i], nuc.r.ipamer)

		privateSuper := fmt.Sprintf("%t", nws[i].PrivateSuper)
		nat := fmt.Sprintf("%t", nws[i].Nat)
		underlay := fmt.Sprintf("%t", nws[i].Underlay)
		prefixes := strings.Join(nws[i].Prefixes.String(), ",")
		destPrefixes := strings.Join(nws[i].DestinationPrefixes.String(), ",")
		vrf := strconv.FormatUint(uint64(nws[i].Vrf), 3)

		metric, err := prometheus.NewConstMetric(
			usedIpsDesc,
			prometheus.CounterValue,
			float64(usage.UsedIPs),
			nws[i].ID,
			prefixes,
			destPrefixes,
			nws[i].PartitionID,
			nws[i].ProjectID,
			nws[i].ParentNetworkID,
			vrf,
			privateSuper,
			nat,
			underlay,
		)
		if err != nil {
			zapup.MustRootLogger().Error("Failed create metric for UsedIPs", zap.Error(err))
			return
		}
		ch <- metric

		metric, err = prometheus.NewConstMetric(
			availableIpsDesc,
			prometheus.CounterValue,
			float64(usage.AvailableIPs),
			nws[i].ID,
			prefixes,
			destPrefixes,
			nws[i].PartitionID,
			nws[i].ProjectID,
			nws[i].ParentNetworkID,
			vrf,
			privateSuper,
			nat,
			underlay,
		)
		if err != nil {
			zapup.MustRootLogger().Error("Failed create metric for AvailableIPs", zap.Error(err))
			return
		}
		ch <- metric
		metric, err = prometheus.NewConstMetric(
			usedPrefixesDesc,
			prometheus.CounterValue,
			float64(usage.UsedPrefixes),
			nws[i].ID,
			prefixes,
			destPrefixes,
			nws[i].PartitionID,
			nws[i].ProjectID,
			nws[i].ParentNetworkID,
			vrf,
			privateSuper,
			nat,
			underlay,
		)
		if err != nil {
			zapup.MustRootLogger().Error("Failed create metric for UsedPrefixes", zap.Error(err))
			return
		}
		ch <- metric
		metric, err = prometheus.NewConstMetric(
			availablePrefixesDesc,
			prometheus.CounterValue,
			float64(usage.AvailablePrefixes),
			nws[i].ID,
			prefixes,
			destPrefixes,
			nws[i].PartitionID,
			nws[i].ProjectID,
			nws[i].ParentNetworkID,
			vrf,
			privateSuper,
			nat,
			underlay,
		)
		if err != nil {
			zapup.MustRootLogger().Error("Failed create metric for AvailablePrefixes", zap.Error(err))
			return
		}
		ch <- metric
	}
}
