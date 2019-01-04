package main

import (
	"errors"
	"fmt"
	"log"
	"net"

	"github.com/coreos/go-iptables/iptables"
)

// OpenAccess opens ipv4 access for remoteAddr.
// will return error if anything goes wrong.
func OpenAccess(remoteAddr string) error {
	ipv4, err := parseIpv4(remoteAddr)
	if err != nil {
		return err
	}

	if err := validate(ipv4); err != nil {
		return err
	}

	return openFirewallFor(ipv4)
}

func parseIpv4(remoteAddr string) (*string, error) {
	s, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		log.Printf("split address from port failed for %v", remoteAddr)
		return nil, err
	}
	ipObj := net.ParseIP(s)
	if ipObj == nil {
		log.Printf("Failed to create ip from %v", s)
		return nil, errors.New("Failed to parse ip address")
	}
	ip4 := ipObj.To4()
	if ip4 == nil {
		log.Printf("the ip is not ipv4, can only handle ipv4!")
		return nil, errors.New(fmt.Sprintf("ip %v is not ipv4, can't handle it!", ipObj))
	}
	ip := ip4.String()
	return &ip, nil
}

// @TODO implement
func validate(ipaddr *string) error {
	log.Printf("TBD validate ipaddress is ok etc etc")
	return nil
}

func openFirewallFor(ipaddr *string) error {

	vpnNet := "192.168.123.0/24"
	client := fmt.Sprintf("%v/32", *ipaddr)

	var err error
	ipt, err := iptables.New()
	if err != nil {
		log.Printf("Failed to create iptables wrapper %v", err)
		return err
	}
	err = ipt.AppendUnique("nat", "POSTROUTING", "-s", client, "-d", vpnNet, "-j", "MASQUERADE")
	if err != nil {
		log.Printf("failed to add to nat %v", err)
		return err
	}
	err = ipt.AppendUnique("filter", "FORWARD", "-s", client, "-d", vpnNet, "-j", "ACCEPT")
	if err != nil {
		log.Printf("failed to add to nat %v", err)
		return err
	}
	return nil
}
