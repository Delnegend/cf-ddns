package cf

import (
	"cf-ddns/appstate"
	"fmt"
	"net/http"
)

type CF struct {
	domain     string
	ddnsDomain string
	apiKey     string
	zoneID     string
	proxy      bool

	getZoneIDReq       *http.Request
	getDNSRecordReq    *http.Request
	createDNSRecordReq *http.Request
	updateDNSRecordReq *http.Request
}

func NewCF(appState *appstate.AppState) (*CF, error) {
	ddnsDomain := appState.GetDDNSDomain()

	cfInstance := CF{
		domain:     appState.GetDomain(),
		ddnsDomain: ddnsDomain,
		apiKey:     appState.GetApiKey(),
		proxy:      appState.GetProxy(),
	}
	if _, err := cfInstance.getZoneID(); err != nil {
		return nil, fmt.Errorf("NewCF: %s", err)
	}
	return &cfInstance, nil
}
