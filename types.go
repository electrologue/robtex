package robtex

// IPQueryResponse ipquery response type.
type IPQueryResponse struct {
	Status            string `json:"status"`
	City              string `json:"city"`
	Country           string `json:"country"`
	AS                int    `json:"as"`
	ASName            string `json:"asname"`
	Whois             string `json:"whoisdesc"`
	Route             string `json:"routedesc"`
	BGPRoute          string `json:"bgproute"`
	ActiveForwardDNS  []Item `json:"act"`
	ActiveDNSHistory  []Item `json:"acth"`
	PassiveReverseDNS []Item `json:"pas"`
	PassiveDNSHistory []Item `json:"pash"`
}

// Item IP item.
type Item struct {
	O         string `json:"o"`
	Timestamp int    `json:"t"`
}

// ASQueryResponse asquery response type.
type ASQueryResponse struct {
	Status string   `json:"status"`
	Nets   []Prefix `json:"nets"`
}

// Prefix AS prefix.
type Prefix struct {
	N     string `json:"n"`
	InBGP int    `json:"inbgp"`
}

// PassiveDNS pdns response type.
type PassiveDNS struct {
	RRName    string `json:"rrname"`
	RRData    string `json:"rrdata"`
	RRType    string `json:"rrtype"`
	TimeFirst int    `json:"time_first"`
	TimeLast  int    `json:"time_last"`
	Count     int    `json:"count"`
}
