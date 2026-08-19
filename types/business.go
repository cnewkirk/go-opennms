package types

// Business Services (v2) and flow classifications.

// BsFunction is a map or reduce function reference within a
// Business Service definition.
type BsFunction struct {
	// Type is the function name, e.g. "Identity", "Increase",
	// "HighestSeverity", "Threshold".
	Type string `json:"type,omitempty"`
	// Properties holds function-specific properties,
	// e.g. {"threshold": "0.5"} for "Threshold".
	Properties map[string]string `json:"properties,omitempty"`
}

// BusinessService is the business service payload for
// Client.CreateBusinessService and Client.UpdateBusinessService.
type BusinessService struct {
	// Name is the unique business service name. Required.
	Name string `json:"name,omitempty"`
	// Attributes holds arbitrary key/value metadata.
	Attributes map[string]string `json:"attributes,omitempty"`
	// ReduceFunction is the reduce function configuration.
	ReduceFunction *BsFunction `json:"reduceFunction,omitempty"`
}

// BsIpServiceEdge is the IP-service edge payload for
// Client.AddIpServiceEdge.
type BsIpServiceEdge struct {
	// IPServiceID is the database ID of the monitored IP service.
	// Required.
	IPServiceID int `json:"ipServiceId,omitempty"`
	// MapFunction is the map function configuration.
	MapFunction *BsFunction `json:"mapFunction,omitempty"`
	// Weight is the edge weight (higher weight = greater
	// influence).
	Weight int `json:"weight,omitempty"`
}

// BsReductionKeyEdge is the reduction-key edge payload for
// Client.AddReductionKeyEdge.
type BsReductionKeyEdge struct {
	// ReductionKey is the alarm reduction key to monitor.
	// Required.
	ReductionKey string `json:"reductionKey,omitempty"`
	// MapFunction is the map function configuration.
	MapFunction *BsFunction `json:"mapFunction,omitempty"`
	// Weight is the edge weight.
	Weight int `json:"weight,omitempty"`
}

// BsChildEdge is the child-service edge payload for
// Client.AddChildEdge.
type BsChildEdge struct {
	// ChildID is the database ID of the child business service.
	// Required.
	ChildID int `json:"childId,omitempty"`
	// MapFunction is the map function configuration.
	MapFunction *BsFunction `json:"mapFunction,omitempty"`
	// Weight is the edge weight.
	Weight int `json:"weight,omitempty"`
}

// ClassificationRule is the classification rule payload for
// Client.CreateClassificationRule and
// Client.UpdateClassificationRule.
type ClassificationRule struct {
	// Name is the application name this rule identifies. Required.
	Name string `json:"name,omitempty"`
	// DstPort is the destination port or range, e.g. "443" or
	// "8080-8090".
	DstPort string `json:"dstPort,omitempty"`
	// SrcPort is the source port or range.
	SrcPort string `json:"srcPort,omitempty"`
	// DstAddress is the destination IP address or CIDR,
	// e.g. "10.0.0.0/8".
	DstAddress string `json:"dstAddress,omitempty"`
	// SrcAddress is the source IP address or CIDR.
	SrcAddress string `json:"srcAddress,omitempty"`
	// Protocol is the protocol: "tcp", "udp", "icmp".
	Protocol string `json:"protocol,omitempty"`
	// ExporterFilter is a FIQL filter selecting which exporters
	// this rule applies to.
	ExporterFilter string `json:"exporterFilter,omitempty"`
	// GroupID is the ID of the classification group this rule
	// belongs to.
	GroupID int `json:"groupId,omitempty"`
	// Position is the sort position within the group. Pointer so
	// position 0 survives omitempty.
	Position *int `json:"position,omitempty"`
}

// ClassificationGroup is the classification group payload for
// Client.CreateClassificationGroup and
// Client.UpdateClassificationGroup.
type ClassificationGroup struct {
	// Name is the unique group name. Required.
	Name string `json:"name,omitempty"`
	// Enabled reports whether this group's rules are active.
	Enabled *bool `json:"enabled,omitempty"`
}

// ClassifyRequest is the flow classification request for
// Client.Classify.
type ClassifyRequest struct {
	// SrcAddress is the source IP address.
	SrcAddress string `json:"srcAddress,omitempty"`
	// SrcPort is the source port number.
	SrcPort int `json:"srcPort,omitempty"`
	// DstAddress is the destination IP address.
	DstAddress string `json:"dstAddress,omitempty"`
	// DstPort is the destination port number.
	DstPort int `json:"dstPort,omitempty"`
	// Protocol is the protocol integer or name.
	Protocol string `json:"protocol,omitempty"`
	// ExporterAddress is the IP address of the flow exporter.
	ExporterAddress string `json:"exporterAddress,omitempty"`
}
