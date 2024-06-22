package cf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (cf *CF) UpdateDNSRecord(recordID string, realIP string) error {
	zoneID, err := cf.getZoneID()
	if err != nil {
		return fmt.Errorf("updateDNSRecord: %s", err)
	}

	if cf.updateDNSRecordReq == nil {
		cf.updateDNSRecordReq, err = http.NewRequest("PATCH", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records/"+recordID, nil)
		if err != nil {
			cf.updateDNSRecordReq.Body = nil
			return fmt.Errorf("updateDNSRecord: %s", err)
		}
		cf.updateDNSRecordReq.Header.Set("Authorization", "Bearer "+cf.apiKey)
		cf.updateDNSRecordReq.Header.Set("Content-Type", "application/json")
	}

	body := struct {
		Content string `json:"content"`
		Name    string `json:"name"`
		Proxied bool   `json:"proxied"`
		Type    string `json:"type"`
		Comment string `json:"comment"`
	}{
		Content: realIP,
		Name:    cf.ddnsDomain,
		Proxied: cf.proxy,
		Type:    "A",
		Comment: "DDNS record",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: %s", err)
	}

	cf.updateDNSRecordReq.Body = io.NopCloser(strings.NewReader(string(jsonBody)))

	resp, err := http.DefaultClient.Do(cf.updateDNSRecordReq)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: %s", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: %s", err)
	}
	respBodyStruct := struct {
		Success bool     `json:"success"`
		Errors  []string `json:"errors"`
	}{}
	err = json.Unmarshal(respBytes, &respBodyStruct)
	if err != nil {
		return fmt.Errorf("updateDNSRecord: %s", err)
	}
	if !respBodyStruct.Success {
		return fmt.Errorf("updateDNSRecord: %s", strings.Join(respBodyStruct.Errors, ", "))
	}

	return nil
}
