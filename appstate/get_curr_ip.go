package appstate

import (
	"fmt"
	"net"
)

func (as *AppState) GetCurrentIP() (string, error) {
	domain := as.GetDDNSDomain()

	addrs, err := net.LookupHost(domain)
	if err != nil {
		return "", fmt.Errorf("GetCurrentIP: failed to lookup host %s: %s", domain, err)
	}
	if len(addrs) == 0 {
		return "", fmt.Errorf("GetCurrentIP: no addresses found for %s", domain)
	}
	return addrs[0], nil
}
