package types

// Requisition / provisioning API (most keys use hyphens).

// ForeignSourceDetector is a provisioning detector within a
// ForeignSource. Also accepted by Client.AddForeignSourceDetector.
type ForeignSourceDetector struct {
	// Name is the detector display name, e.g. "ICMP". Required.
	Name string `json:"name,omitempty"`
	// Class is the fully-qualified Java detector class name, e.g.
	// "org.opennms.netmgt.provision.detector.icmp.IcmpDetector".
	// Required.
	Class string `json:"class,omitempty"`
	// Parameter holds detector parameters as
	// [{"key": ..., "value": ...}].
	Parameter []map[string]any `json:"parameter,omitempty"`
}

// ForeignSourcePolicy is a provisioning policy within a
// ForeignSource. Also accepted by Client.AddForeignSourcePolicy.
type ForeignSourcePolicy struct {
	// Name is the policy display name. Required.
	Name string `json:"name,omitempty"`
	// Class is the fully-qualified Java policy class name, e.g.
	// "org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy".
	// Required.
	Class string `json:"class,omitempty"`
	// Parameter holds policy parameters as
	// [{"key": ..., "value": ...}].
	Parameter []map[string]any `json:"parameter,omitempty"`
}

// ForeignSource is the foreign source definition for
// Client.CreateForeignSource and Client.UpdateForeignSource.
type ForeignSource struct {
	// Name is the foreign source name. Required.
	Name string `json:"name,omitempty"`
	// ScanInterval is the rescan interval in OpenNMS duration
	// format, e.g. "1d" (daily), "1w" (weekly), "-1" (no rescan).
	ScanInterval string `json:"scan-interval,omitempty"`
	// Detectors holds the provisioning detectors.
	Detectors []ForeignSourceDetector `json:"detectors,omitempty"`
	// Policies holds the provisioning policies.
	Policies []ForeignSourcePolicy `json:"policies,omitempty"`
}

// RequisitionService is a monitored-service entry within a
// RequisitionInterface. Also accepted by
// Client.CreateRequisitionNodeService.
type RequisitionService struct {
	// ServiceName is the service name matching a provisioning
	// detector, e.g. "ICMP", "SNMP", "HTTP".
	ServiceName string `json:"service-name,omitempty"`
}

// RequisitionInterface is an IP interface entry within a
// RequisitionNode. Also accepted by
// Client.CreateRequisitionNodeInterface.
type RequisitionInterface struct {
	// IPAddr is the IP address of the interface. Required.
	IPAddr string `json:"ip-addr,omitempty"`
	// SnmpPrimary is the SNMP primary flag — "P" (primary), "S"
	// (secondary), "N" (not eligible). Default "N".
	SnmpPrimary string `json:"snmp-primary,omitempty"`
	// Status is the management status — 1 managed, 3 unmanaged.
	Status int `json:"status,omitempty"`
	// MonitoredService lists the services to monitor.
	MonitoredService []RequisitionService `json:"monitored-service,omitempty"`
}

// RequisitionAsset is an asset record entry within a
// RequisitionNode. Also accepted by Client.SetRequisitionNodeAsset.
type RequisitionAsset struct {
	// Name is the asset field name, e.g. "manufacturer",
	// "serialNumber", "description".
	Name string `json:"name,omitempty"`
	// Value is the asset field value.
	Value string `json:"value,omitempty"`
}

// RequisitionNode is the node payload for
// Client.CreateRequisitionNode.
type RequisitionNode struct {
	// ForeignID is the unique node identifier within the
	// requisition. Must be stable across imports. Required.
	ForeignID string `json:"foreign-id,omitempty"`
	// NodeLabel is the human-readable node display name. Required.
	NodeLabel string `json:"node-label,omitempty"`
	// Location is the monitoring location name.
	// Defaults to "Default".
	Location string `json:"location,omitempty"`
	// Interface holds the IP interfaces to provision.
	Interface []RequisitionInterface `json:"interface,omitempty"`
	// Category holds surveillance categories,
	// e.g. [{"name": "Production"}].
	Category []map[string]any `json:"category,omitempty"`
	// Asset holds asset record fields.
	Asset []RequisitionAsset `json:"asset,omitempty"`
	// MetaData holds metadata entries with "context", "key", and
	// "value" keys.
	MetaData []map[string]any `json:"meta-data,omitempty"`
}

// Requisition is the requisition payload for
// Client.CreateRequisition.
type Requisition struct {
	// ForeignSource is the requisition (foreign source) name.
	// Required.
	ForeignSource string `json:"foreign-source,omitempty"`
	// Node lists the nodes to include in the requisition.
	Node []RequisitionNode `json:"node,omitempty"`
}
