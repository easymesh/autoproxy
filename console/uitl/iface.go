package util

import (
	"fmt"
	"io/ioutil"
	"net"
	"strings"
)

type InterfaceInfo struct {
	Ifname  string
	IP      net.IP
	IPNet   net.IPNet
	Mac     string
	Flag    net.Flags
}

func getIPv4(addrs []net.Addr) (net.IP, net.IPNet) {
	for _, v := range addrs {
		ip, ipnet, err := net.ParseCIDR(v.String())
		if err != nil {
			continue
		}
		if ip.IsLoopback() {
			continue
		}
		if ipnet.IP.To4() == nil {
			continue
		}
		return ip, *ipnet
	}
	return nil, net.IPNet{}
}

func InterfaceGetByIP(ip string) (*InterfaceInfo, error) {
	ifaces, err := InterfaceGet()
	if err != nil {
		return nil, err
	}
	for _, v := range ifaces {
		if v.IP.String() == ip {
			return &v, nil
		}
	}
	return nil, fmt.Errorf("interface %s not exist", ip)
}

func InterfaceGet() ([]InterfaceInfo, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	var output []InterfaceInfo

	for _, i := range ifaces {
		if strings.Compare(i.Name, "lo") == 0 {
			continue
		}
		byName, err := net.InterfaceByName(i.Name)
		if err != nil {
			continue
		}
		temp := InterfaceInfo{
			Ifname:  i.Name,
			Mac:     i.HardwareAddr.String(),
			Flag:    i.Flags,
		}
		addresses, err := byName.Addrs()
		if len(addresses) > 0 {
			temp.IP, temp.IPNet = getIPv4(addresses)
		}
		output = append(output, temp)
	}
	return output, nil
}

func HostnameGet() string {
	hostname, err := ioutil.ReadFile("/etc/hostname")
	if err != nil {
		return "unkown"
	}
	return strings.ReplaceAll(string(hostname),"\n","")
}