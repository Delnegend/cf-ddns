package cf

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var NO_DNS_RECORD_FOUND = "NO_DNS_RECORD_FOUND"

type DNSRecord struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

// Get the CF Record ID and value for the DNS record for DDNS domain
func (cf *CF) GetDDNSRecordInfo() (DNSRecord, error) {
	zoneID, err := cf.getZoneID()
	if err != nil {
		return DNSRecord{}, fmt.Errorf("getRecordID: %s", err)
	}

	if cf.getDNSRecordReq == nil {
		cf.getDNSRecordReq, err = http.NewRequest("GET", "https://api.cloudflare.com/client/v4/zones/"+zoneID+"/dns_records", nil)
		if err != nil {
			return DNSRecord{}, fmt.Errorf("getRecordID: %s", err)
		}
		cf.getDNSRecordReq.Header.Set("Authorization", "Bearer "+cf.apiKey)
	}

	resp, err := http.DefaultClient.Do(cf.getDNSRecordReq)
	if err != nil {
		return DNSRecord{}, fmt.Errorf("getRecordID: %s", err)
	}
	defer resp.Body.Close()

	bodyStr, err := io.ReadAll(resp.Body)
	if err != nil {
		return DNSRecord{}, fmt.Errorf("getRecordID: %s", err)
	}

	var body struct {
		Result []DNSRecord `json:"result"`
	}
	err = json.Unmarshal(bodyStr, &body)
	if err != nil {
		return DNSRecord{}, fmt.Errorf("getRecordID: %s", err)
	}

	for _, dnsRecord := range body.Result {
		if dnsRecord.Name == cf.ddnsDomain {
			return dnsRecord, nil
		}
	}

	return DNSRecord{Name: NO_DNS_RECORD_FOUND}, fmt.Errorf("getRecordID: no record found for %s", cf.ddnsDomain)
}
