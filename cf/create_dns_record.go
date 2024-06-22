package cf

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func (cf *CF) CreateDNSRecord(realIP string) error {
	zoneID, err := cf.getZoneID()
	if err != nil {
		return fmt.Errorf("createDNSRecord: %s", err)
	}

	if cf.createDNSRecordReq == nil {
		cf.createDNSRecordReq, err = http.NewRequest("POST", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records", nil)
		if err != nil {
			return fmt.Errorf("createDNSRecord: %s", err)
		}
		cf.createDNSRecordReq.Header.Set("Authorization", "Bearer "+cf.apiKey)
		cf.createDNSRecordReq.Header.Set("Content-Type", "application/json")
	}

	body := struct {
		Content string `json:"content"`
		Name    string `json:"name"`
		Type    string `json:"type"`
		Comment string `json:"comment"`
	}{
		Content: realIP,
		Name:    cf.ddnsDomain,
		Type:    "A",
		Comment: "DDNS record",
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("createDNSRecord: %s", err)
	}

	cf.createDNSRecordReq.Body = io.NopCloser(bytes.NewReader(jsonBody))

	resp, err := http.DefaultClient.Do(cf.createDNSRecordReq)
	if err != nil {
		return fmt.Errorf("createDNSRecord: %s", err)
	}
	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("createDNSRecord: %s", err)
	}
	respBodyStruct := struct {
		Success bool     `json:"success"`
		Errors  []string `json:"errors"`
	}{}
	err = json.Unmarshal(respBytes, &respBodyStruct)
	if err != nil {
		return fmt.Errorf("createDNSRecord: %s", err)
	}
	if !respBodyStruct.Success {
		return fmt.Errorf("createDNSRecord: %s", strings.Join(respBodyStruct.Errors, ", "))
	}

	return nil
}
