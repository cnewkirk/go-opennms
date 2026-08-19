package types

// Scheduled outages, KSC reports, maps, and monitoring locations.

// SchedOutageTime is a time window entry within a SchedOutage.
type SchedOutageTime struct {
	// Day is the day specifier. For "weekly" outages: lowercase
	// weekday name (e.g. "saturday"). For "monthly": integer
	// day-of-month.
	Day string `json:"day,omitempty"`
	// Begins is the start time "HH:MM:SS" (daily/weekly/monthly)
	// or full datetime "DD-Mon-YYYY HH:MM:SS" (specific).
	Begins string `json:"begins,omitempty"`
	// Ends is the end time in the same format as Begins.
	Ends string `json:"ends,omitempty"`
}

// SchedOutage is the scheduled outage payload for
// Client.CreateSchedOutage.
type SchedOutage struct {
	// Name is the unique scheduled outage name. Required.
	Name string `json:"name,omitempty"`
	// Type is the outage type: "weekly", "monthly", "specific",
	// or "daily".
	Type string `json:"type,omitempty"`
	// Time lists the time windows during which the outage applies.
	Time []SchedOutageTime `json:"time,omitempty"`
	// Node lists the nodes to suppress during the outage,
	// e.g. [{"id": 1}, {"id": 2}].
	Node []map[string]any `json:"node,omitempty"`
	// Interface lists the IP interfaces to suppress,
	// e.g. [{"address": "192.168.0.1"}].
	Interface []map[string]any `json:"interface,omitempty"`
}

// KscGraph is a graph entry within a KscReport.
type KscGraph struct {
	// Title is the graph title.
	Title string `json:"title,omitempty"`
	// ResourceID is the OpenNMS resource ID,
	// e.g. "node[1].interfaceSnmp[eth0-04013f75f101]".
	ResourceID string `json:"resourceId,omitempty"`
	// Timespan is the timespan identifier, e.g. "7_day",
	// "1_month".
	Timespan string `json:"timespan,omitempty"`
	// Graphtype is the graph type (RRD graph definition name),
	// e.g. "mib2.bits".
	Graphtype string `json:"graphtype,omitempty"`
	// GraphIndex is the position of this graph within the report.
	// Pointer so index 0 survives omitempty.
	GraphIndex *int `json:"graphIndex,omitempty"`
}

// KscReport is the KSC report payload for Client.CreateKscReport.
type KscReport struct {
	// ID is the report ID (set to 0 for new reports). Pointer so
	// the documented 0 survives omitempty.
	ID *int `json:"id,omitempty"`
	// Label is the report display name. Required.
	Label string `json:"label,omitempty"`
	// ShowTimespanButton shows the timespan selection button.
	ShowTimespanButton *bool `json:"show_timespan_button,omitempty"`
	// ShowGraphtypeButton shows the graph type selection button.
	ShowGraphtypeButton *bool `json:"show_graphtype_button,omitempty"`
	// GraphsPerLine is the number of graphs per row.
	GraphsPerLine int `json:"graphs_per_line,omitempty"`
	// Graphs lists the graph definitions.
	Graphs []KscGraph `json:"graphs,omitempty"`
}

// MapPayload is the map payload for Client.CreateMap and
// Client.UpdateMap.
//
// It corresponds to the Map TypedDict in python-opennms's types.py;
// the struct is named MapPayload here because the Map identifier is
// taken by this package's Map conversion function.
type MapPayload struct {
	// Name is the unique map name. Required for create.
	Name string `json:"name,omitempty"`
	// MapWidth is the map canvas width in pixels.
	MapWidth int `json:"mapWidth,omitempty"`
	// MapHeight is the map canvas height in pixels.
	MapHeight int `json:"mapHeight,omitempty"`
	// AccessMode is the access control mode: "RW" (read-write) or
	// "RO" (read-only).
	AccessMode string `json:"accessMode,omitempty"`
	// Owner is the username of the map owner.
	Owner string `json:"owner,omitempty"`
}

// MonitoringLocation is the monitoring location payload for
// Client.CreateMonitoringLocation.
type MonitoringLocation struct {
	// LocationName is the unique location identifier. Required.
	LocationName string `json:"location-name,omitempty"`
	// MonitoringArea is a geographic or logical area label.
	MonitoringArea string `json:"monitoring-area,omitempty"`
	// Priority is the display sort order. Pointer so priority 0
	// survives omitempty.
	Priority *int `json:"priority,omitempty"`
	// Tags holds optional location tags.
	Tags []string `json:"tags,omitempty"`
}
