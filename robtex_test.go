package robtex

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, pattern, filename string) *Client {
	t.Helper()

	mux := http.NewServeMux()

	server := httptest.NewServer(mux)
	t.Cleanup(server.Close)

	client := New("")
	client.httpClient = server.Client()
	client.baseURL, _ = url.Parse(server.URL)

	mux.HandleFunc(pattern, func(rw http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodGet {
			http.Error(rw, "unsupported method", http.StatusMethodNotAllowed)
			return
		}

		file, err := os.Open(filepath.FromSlash(path.Join("./fixtures", filename)))
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
		defer func() { _ = file.Close() }()

		_, err = io.Copy(rw, file)
		if err != nil {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	return client
}

func TestClient_IPQuery(t *testing.T) {
	client := setup(t, "/ipquery/199.19.54.1", "ipquery.json")

	response, err := client.IPQuery(context.Background(), "199.19.54.1")
	require.NoError(t, err)

	expected := &IPQueryResponse{
		Status:   "ok",
		City:     "Toronto",
		Country:  "Canada",
		AS:       12041,
		ASName:   "AFILIAS-NST Afilias Limited",
		Whois:    "Afilias Canada, Corp. (ACC-308)",
		Route:    "REACH (Customer Route)",
		BGPRoute: "199.19.54.0/24",
		ActiveForwardDNS: []Item{
			{O: "b0.org.afilias-nst.org", Timestamp: 1498717654},
		},
		ActiveDNSHistory: []Item{},
		PassiveReverseDNS: []Item{
			{O: "b0.org.afilias-nst.org", Timestamp: 1501358830},
			{O: "b0.nic.ngo", Timestamp: 1501358850},
			{O: "b0.org.afilias-nst.info", Timestamp: 1493870770},
			{O: "b0.nic.ong", Timestamp: 1500758606},
			{O: "b0.org.afilias-nst.org", Timestamp: 1452591322},
		},
		PassiveDNSHistory: []Item{},
	}

	assert.Equal(t, expected, response)
}

func TestClient_ASQuery(t *testing.T) {
	client := setup(t, "/asquery/1234", "asquery.json")

	response, err := client.ASQuery(context.Background(), "1234")
	require.NoError(t, err)

	expected := &ASQueryResponse{
		Status: "ok",
		Nets: []Prefix{
			{N: "132.171.0.0/16", InBGP: 1},
			{N: "132.171.0.0/17", InBGP: 1},
			{N: "132.171.128.0/17", InBGP: 1},
			{N: "193.110.32.0/21", InBGP: 1},
			{N: "193.110.32.0/22", InBGP: 1},
			{N: "193.110.36.0/22", InBGP: 1},
		},
	}

	assert.Equal(t, expected, response)
}

func TestClient_PassiveDNSForward(t *testing.T) {
	client := setup(t, "/pdns/forward/a.iana-servers.net", "ppdns-forward.ldjson")

	response, err := client.PassiveDNSForward(context.Background(), "a.iana-servers.net")
	require.NoError(t, err)

	expected := []PassiveDNS{
		{RRName: "a.iana-servers.net", RRData: "2001:500:8c::53", RRType: "AAAA", TimeFirst: 1441242410, TimeLast: 1460542918, Count: 18},
		{RRName: "a.iana-servers.net", RRData: "2001:500:8f::53", RRType: "AAAA", TimeFirst: 1460751956, TimeLast: 1501399246, Count: 18},
		{RRName: "a.iana-servers.net", RRData: "199.43.132.53", RRType: "A", TimeFirst: 1441242410, TimeLast: 1460542918, Count: 18},
		{RRName: "a.iana-servers.net", RRData: "199.43.135.53", RRType: "A", TimeFirst: 1460751956, TimeLast: 1501399246, Count: 18},
	}

	assert.Equal(t, expected, response)
}

func TestClient_PassiveDNSReverse(t *testing.T) {
	client := setup(t, "/pdns/reverse/199.43.132.53", "ppdns-reverse.ldjson")

	response, err := client.PassiveDNSReverse(context.Background(), "199.43.132.53")
	require.NoError(t, err)

	expected := []PassiveDNS{
		{RRName: "a.iana-servers.org", RRData: "199.43.132.53", RRType: "A", TimeFirst: 1439620242, TimeLast: 1460165924, Count: 18},
		{RRName: "a.icann-servers.net", RRData: "199.43.132.53", RRType: "A", TimeFirst: 1448005462, TimeLast: 1456189254, Count: 2},
		{RRName: "a.iana-servers.org", RRData: "199.43.132.53", RRType: "A", TimeFirst: 1439620242, TimeLast: 1460165924, Count: 18},
		{RRName: "a.iana-servers.net", RRData: "199.43.132.53", RRType: "A", TimeFirst: 1441242410, TimeLast: 1460542918, Count: 18},
	}

	assert.Equal(t, expected, response)
}
