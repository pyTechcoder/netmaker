package ncutils

import (
	"net"

	"github.com/gravitl/netmaker/models"
)

func IfaceDelta(currentNode *models.Node, newNode *models.Node) bool {
	// single comparison statements
	if newNode.Endpoint != currentNode.Endpoint ||
		newNode.LocalAddress != currentNode.LocalAddress ||
		newNode.PublicKey != currentNode.PublicKey ||
		newNode.Address != currentNode.Address ||
		newNode.IsEgressGateway != currentNode.IsEgressGateway ||
		newNode.IsIngressGateway != currentNode.IsIngressGateway ||
		newNode.IsRelay != currentNode.IsRelay ||
		newNode.UDPHolePunch != currentNode.UDPHolePunch ||
		newNode.IsPending != currentNode.IsPending ||
		newNode.PersistentKeepalive != currentNode.PersistentKeepalive ||
		newNode.DNSOn != currentNode.DNSOn ||
		len(newNode.ExcludedAddrs) != len(currentNode.ExcludedAddrs) ||
		len(newNode.AllowedIPs) != len(currentNode.AllowedIPs) {
		return true
	}

	// multi-comparison statements
	if newNode.IsDualStack == "yes" {
		if newNode.Address6 != currentNode.Address6 {
			return true
		}
	}

	if newNode.IsEgressGateway == "yes" {
		if len(currentNode.EgressGatewayRanges) != len(newNode.EgressGatewayRanges) {
			return true
		}
		for _, address := range newNode.EgressGatewayRanges {
			if !StringSliceContains(currentNode.EgressGatewayRanges, address) {
				return true
			}
		}
	}

	if newNode.IsRelay == "yes" {
		if len(currentNode.RelayAddrs) != len(newNode.RelayAddrs) {
			return true
		}
		for _, address := range newNode.RelayAddrs {
			if !StringSliceContains(currentNode.RelayAddrs, address) {
				return true
			}
		}
	}

	for _, address := range newNode.AllowedIPs {
		if !StringSliceContains(currentNode.AllowedIPs, address) {
			return true
		}
	}
	return false
}

// StringSliceContains - sees if a string slice contains a string element
func StringSliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IPNetSliceContains - sees if a string slice contains a string element
func IPNetSliceContains(slice []net.IPNet, item net.IPNet) bool {
	for _, s := range slice {
		if s.String() == item.String() {
			return true
		}
	}
	return false
}

// IfaceExists - return true if you can find the iface
func IfaceExists(ifacename string) bool {
	localnets, err := net.Interfaces()
	if err != nil {
		return false
	}
	for _, localnet := range localnets {
		if ifacename == localnet.Name {
			return true
		}
	}
	return false
}
