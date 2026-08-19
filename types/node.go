package types

// Node API (keys are camelCase).

// Node is the node payload for Client.CreateNode and
// Client.UpdateNode.
type Node struct {
	// Label is the human-readable node display name.
	// Required for create.
	Label string `json:"label,omitempty"`
	// Type is the node type. "A" (active) is the standard value.
	Type string `json:"type,omitempty"`
	// ForeignSource is the requisition / foreign source name.
	ForeignSource string `json:"foreignSource,omitempty"`
	// ForeignID is the unique ID within the foreign source.
	ForeignID string `json:"foreignId,omitempty"`
	// Location is the monitoring location name.
	// Defaults to "Default".
	Location string `json:"location,omitempty"`
	// SysName is the SNMP sysName.
	SysName string `json:"sysName,omitempty"`
	// SysDescription is the SNMP sysDescr.
	SysDescription string `json:"sysDescription,omitempty"`
	// SysContact is the SNMP sysContact.
	SysContact string `json:"sysContact,omitempty"`
	// SysLocation is the SNMP sysLocation.
	SysLocation string `json:"sysLocation,omitempty"`
}

// NodeIpInterface is the IP interface payload for
// Client.CreateNodeIpInterface and Client.UpdateNodeIpInterface.
type NodeIpInterface struct {
	// IPAddress is the IP address of the interface.
	// Required for create.
	IPAddress string `json:"ipAddress,omitempty"`
	// IsManaged is the management status. "M" managed, "U"
	// unmanaged, "D" deleted. Defaults to "M".
	IsManaged string `json:"isManaged,omitempty"`
	// SnmpPrimary is the SNMP primary flag. "P" primary, "S"
	// secondary, "N" not eligible. Defaults to "N".
	SnmpPrimary string `json:"snmpPrimary,omitempty"`
	// HostName is the reverse-DNS hostname for this interface.
	HostName string `json:"hostName,omitempty"`
}

// NodeSnmpInterface is the SNMP interface payload for
// Client.CreateNodeSnmpInterface and Client.UpdateNodeSnmpInterface.
type NodeSnmpInterface struct {
	// IfIndex is the SNMP ifIndex. Required for create.
	IfIndex int `json:"ifIndex,omitempty"`
	// IfName is the interface name, e.g. "GigabitEthernet0/0".
	IfName string `json:"ifName,omitempty"`
	// IfDescr is the interface description from SNMP ifDescr.
	IfDescr string `json:"ifDescr,omitempty"`
	// IfAlias is the interface alias from SNMP ifAlias.
	IfAlias string `json:"ifAlias,omitempty"`
	// IfType is the SNMP ifType integer (6 = ethernetCsmacd).
	IfType int `json:"ifType,omitempty"`
	// Collect is the collection flag. "C" collect, "N" do not
	// collect.
	Collect string `json:"collect,omitempty"`
	// Poll is the poll flag. "P" poll, "N" do not poll.
	Poll string `json:"poll,omitempty"`
}

// NodeAssetRecord is the asset record payload for
// Client.UpdateNodeAssetRecord.
//
// All fields are optional. Pass only the fields you want to change.
type NodeAssetRecord struct {
	// Category is the asset category, e.g. "Routers".
	Category string `json:"category,omitempty"`
	// Manufacturer is the hardware manufacturer.
	Manufacturer string `json:"manufacturer,omitempty"`
	// Vendor is the vendor name.
	Vendor string `json:"vendor,omitempty"`
	// ModelNumber is the model number.
	ModelNumber string `json:"modelNumber,omitempty"`
	// SerialNumber is the serial number.
	SerialNumber string `json:"serialNumber,omitempty"`
	// Description is a free-text description.
	Description string `json:"description,omitempty"`
	// OperatingSystem is the operating system name and version.
	OperatingSystem string `json:"operatingSystem,omitempty"`
	// Rack is the rack identifier.
	Rack string `json:"rack,omitempty"`
	// Building is the building identifier.
	Building string `json:"building,omitempty"`
	// Floor is the floor identifier.
	Floor string `json:"floor,omitempty"`
	// Room is the room identifier.
	Room string `json:"room,omitempty"`
	// Country is the country code.
	Country string `json:"country,omitempty"`
}

// HardwareEntity is the hardware inventory entity payload for
// Client.AddNodeHardwareInventory and Client.UpdateNodeHardwareEntity.
type HardwareEntity struct {
	// EntityPhysicalIndex is the ENTITY-MIB entPhysicalIndex.
	// Required for create.
	EntityPhysicalIndex int `json:"entityPhysicalIndex,omitempty"`
	// EntPhysicalDescr is the physical description.
	EntPhysicalDescr string `json:"entPhysicalDescr,omitempty"`
	// EntPhysicalClass is the physical class integer
	// (e.g. 3 = chassis).
	EntPhysicalClass int `json:"entPhysicalClass,omitempty"`
	// EntPhysicalName is the physical name.
	EntPhysicalName string `json:"entPhysicalName,omitempty"`
	// EntPhysicalSerialNum is the serial number.
	EntPhysicalSerialNum string `json:"entPhysicalSerialNum,omitempty"`
	// EntPhysicalMfgName is the manufacturer name.
	EntPhysicalMfgName string `json:"entPhysicalMfgName,omitempty"`
	// EntPhysicalModelName is the model name.
	EntPhysicalModelName string `json:"entPhysicalModelName,omitempty"`
	// EntPhysicalIsFRU reports whether this entity is a
	// field-replaceable unit.
	EntPhysicalIsFRU *bool `json:"entPhysicalIsFRU,omitempty"`
	// Children holds child hardware entities.
	Children []map[string]any `json:"children,omitempty"`
}
