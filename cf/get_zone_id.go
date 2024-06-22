package cf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (cf *CF) getZoneID() (string, error) {
	if cf.zoneID != "" {
		return cf.zoneID, nil
	}

	var err error
	if cf.getZoneIDReq == nil {
		cf.getZoneIDReq, err = http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones", nil)
		if err != nil {
			return "", fmt.Errorf("getZoneID: %s", err)
		}
		cf.getZoneIDReq.Header.Set("Authorization", "Bearer "+cf.apiKey)
	}

	cf.getZoneIDReq.Header.Set("Authorization", "Bearer "+cf.apiKey)
	resp, err := http.DefaultClient.Do(cf.getZoneIDReq)
	if err != nil {
		return "", fmt.Errorf("getZoneID: %s", err)
	}
	defer resp.Body.Close()

	bodyStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("getZoneID: %s", err)
	}

	var body struct {
		Result []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"result"`
	}
	err = json.Unmarshal(bodyStr, &body)
	if err != nil {
		return "", fmt.Errorf("getZoneID: %s", err)
	}

	for _, zone := range body.Result {
		if zone.Name == cf.domain {
			cf.zoneID = zone.ID
			return zone.ID, nil
		}
	}

	return "", fmt.Errorf("getZoneID: no zone found for %s", cf.domain)
}
