package appstate

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

func (g *AppState) GetRealIP() (string, error) {
	if g.request == nil {
		slog.Debug("CurrentIPGetter.Request: initializing new request")
		g.request, _ = http.NewRequest("GET", "https://cloudflare.com/cdn-cgi/trace", nil)
		g.request.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/126.0.0.0 Safari/537.36")
	} else {
		slog.Debug("CurrentIPGetter.Request: reusing existing request")
	}

	resp, err := http.DefaultClient.Do(g.request)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(body), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "ip=") {
			return strings.TrimPrefix(line, "ip="), nil
		}
	}

	return "", fmt.Errorf("no ip found")
}
