package helper

import (
	"fmt"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/ipam"
	"github.com/metal-stack/metal-api/cmd/metal-api/internal/metal"
)

func GetNetworkUsage(nw *metal.Network, ipamer ipam.IPAMer) *metal.NetworkUsage {
	usage := &metal.NetworkUsage{}
	if nw == nil {
		return usage
	}
	for _, prefix := range nw.Prefixes {
		u, err := ipamer.PrefixUsage(prefix.String())
		if err != nil {
			continue
		}
		usage.AvailableIPs = usage.AvailableIPs + u.AvailableIPs
		usage.UsedIPs = usage.UsedIPs + u.UsedIPs
		usage.AvailablePrefixes = usage.AvailablePrefixes + u.AvailablePrefixes
		usage.UsedPrefixes = usage.UsedPrefixes + u.UsedPrefixes
	}
	return usage
}

func CheckAnyIPOfPrefixesInUse(ips []metal.IP, prefixes metal.Prefixes) error {
	for _, ip := range ips {
		for _, p := range prefixes {
			if ip.ParentPrefixCidr == p.String() {
				return fmt.Errorf("prefix %s has ip %s in use", p.String(), ip.IPAddress)
			}
		}
	}
	return nil
}
