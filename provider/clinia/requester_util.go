package clinia

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"

	"github.com/clinia/models-client-go/cliniamodel/common"
	"github.com/clinia/models-client-go/cliniamodel/requestergrpc"
)

func makeRequester(ctx context.Context, baseURL *string) (common.Requester, error) {
	if baseURL == nil || strings.TrimSpace(*baseURL) == "" {
		return nil, fmt.Errorf("clinia: BaseURL is required")
	}
	raw := strings.TrimSpace(*baseURL)
	if !strings.Contains(raw, "://") {
		raw = "http://" + raw
	}
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("clinia: invalid BaseURL: %w", err)
	}
	hostPort := u.Host
	if hostPort == "" {
		hostPort = strings.TrimPrefix(u.Path, "/")
	}
	host, portStr, err := net.SplitHostPort(hostPort)
	if err != nil {
		return nil, fmt.Errorf("clinia: BaseURL must include host:port (got %q)", hostPort)
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("clinia: invalid port in BaseURL: %w", err)
	}

	scheme := common.HTTP
	cfg := common.RequesterConfig{Host: common.Host{Url: host, Port: port, Scheme: scheme}}
	return requestergrpc.NewRequester(ctx, cfg)
}
