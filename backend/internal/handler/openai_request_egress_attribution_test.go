package handler

import (
	"testing"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func TestOpenAIRequestEgressFieldsClassifiesAccountRoutes(t *testing.T) {
	proxyID := int64(15)
	proxy := &service.Proxy{ID: proxyID, Name: "CodeAgent Proxy SG"}

	tests := []struct {
		name      string
		account   *service.Account
		wantClass string
		wantID    *int64
		wantName  string
	}{
		{name: "nil", account: nil, wantClass: "unknown", wantName: "unknown"},
		{name: "canonical direct", account: &service.Account{Name: "provider [direct]"}, wantClass: "direct", wantName: "direct/no_proxy"},
		{name: "generic proxy", account: &service.Account{Name: "provider [proxy]", ProxyID: &proxyID, Proxy: proxy}, wantClass: "proxy", wantID: &proxyID, wantName: "CodeAgent Proxy SG"},
		{name: "residential proxy", account: &service.Account{Name: "provider [residential]", ProxyID: &proxyID, Proxy: proxy}, wantClass: "residential", wantID: &proxyID, wantName: "CodeAgent Proxy SG"},
		{name: "legacy direct", account: &service.Account{Name: "provider"}, wantClass: "legacy_direct", wantName: "direct/no_proxy"},
		{name: "legacy proxy", account: &service.Account{Name: "provider", ProxyID: &proxyID, Proxy: proxy}, wantClass: "legacy_proxy", wantID: &proxyID, wantName: "CodeAgent Proxy SG"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := openAIRequestEgressAttributionForTest(tt.account)
			if got.Class != tt.wantClass || got.ProxyName != tt.wantName {
				t.Fatalf("attribution = %+v, want class=%s name=%s", got, tt.wantClass, tt.wantName)
			}
			if tt.wantID == nil {
				if got.ProxyID != nil {
					t.Fatalf("proxy id = %d, want nil", *got.ProxyID)
				}
				return
			}
			if got.ProxyID == nil || *got.ProxyID != *tt.wantID {
				t.Fatalf("proxy id = %v, want %d", got.ProxyID, *tt.wantID)
			}
		})
	}
}

func openAIRequestEgressAttributionForTest(account *service.Account) openAIRequestEgressAttribution {
	fields := openAIRequestEgressFields(account)
	got := openAIRequestEgressAttribution{ProxyName: "unknown", Class: "unknown"}
	for _, field := range fields {
		switch field.Key {
		case "egress_class":
			got.Class = field.String
		case "proxy_name":
			got.ProxyName = field.String
		case "proxy_id":
			proxyID := field.Integer
			got.ProxyID = &proxyID
		}
	}
	return got
}
