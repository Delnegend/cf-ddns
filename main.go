package main

import (
	"cf-ddns/appstate"
	"cf-ddns/cf"
	"log/slog"
	"os"
	"time"

	"github.com/lmittmann/tint"
)

func main() {
	asInstance, err := appstate.NewAppState()
	if err != nil {
		slog.Error("NewAppState", "err", err)
		return
	}

	slog.SetDefault(slog.New(
		tint.NewHandler(os.Stderr, &tint.Options{
			Level:      asInstance.GetLogLevel(),
			TimeFormat: time.RFC1123Z,
		}),
	))

	slog.Info("AppState initialized",
		"domain", asInstance.GetDomain(),
		"ddnsDomain", asInstance.GetDDNSDomain(),
		"proxy", asInstance.GetProxy(),
		"getCurrIPMethod", asInstance.GetCurrIPMethod(),
		"sleepInterval", asInstance.GetSleepInterval())

	cfInstance, err := cf.NewCF(asInstance)
	if err != nil {
		slog.Error("NewCF", "err", err)
		return
	}

	isFirstRun := true
	for {
		switch isFirstRun {
		case true:
			isFirstRun = false
		case false:
			time.Sleep(asInstance.GetSleepInterval())
		}

		realIP, err := asInstance.GetRealIP()
		if err != nil {
			slog.Warn("can't get current IP", "err", err)
		}

		if asInstance.GetCurrIPMethod() == appstate.GET_CURR_IP_METHOD_NSLOOKUP {
			currIP, err := asInstance.GetCurrentIP()
			if err != nil {
				slog.Error("can't get real IP", "err", err)
				continue
			}
			if realIP == currIP {
				slog.Info("no change in IP", "realIP", realIP, "currIP", currIP)
				continue
			}
			slog.Info("IP changes", "realIP", realIP, "currIP", currIP)
		}

		dnsRecord, err := cfInstance.GetDDNSRecordInfo()
		if (err != nil) && (dnsRecord.Name != cf.NO_DNS_RECORD_FOUND) {
			slog.Error("can't get DNS record info", "err", err)
			continue
		}

		if asInstance.GetCurrIPMethod() == appstate.GET_CURR_IP_METHOD_CF {
			if dnsRecord.Content == realIP {
				slog.Info("no change in IP", "realIP", realIP, "currIP", dnsRecord.Content)
				continue
			}
			slog.Info("IP changes", "realIP", realIP, "currIP", dnsRecord.Content)
		}

		if dnsRecord.Name == cf.NO_DNS_RECORD_FOUND {
			slog.Info("DNS record not found, creating")
			if err := cfInstance.CreateDNSRecord(realIP); err != nil {
				slog.Error("can't create DNS record", "err", err)
				continue
			}
			slog.Info("DNS record created")
			continue
		}

		slog.Info("Found existing DNS record", "dnsRecordID", dnsRecord.ID, "dnsRecordName", dnsRecord.Name)
		if err := cfInstance.UpdateDNSRecord(dnsRecord.ID, realIP); err != nil {
			slog.Error("can't update DNS record", "err", err)
			continue
		}
		slog.Info("DNS record updated")
	}
}
