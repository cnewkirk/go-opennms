package types

// Events, event configuration, and measurements.

// Event is the event payload for Client.CreateEvent.
type Event struct {
	// UEI is the event UEI,
	// e.g. "uei.opennms.org/internal/test". Required.
	UEI string `json:"uei,omitempty"`
	// Source is the event source identifier (name of the
	// generating script or application).
	Source string `json:"source,omitempty"`
	// Severity is the event severity: "Indeterminate", "Cleared",
	// "Normal", "Warning", "Minor", "Major", "Critical".
	Severity string `json:"severity,omitempty"`
	// NodeID is the database ID of the associated node.
	NodeID int `json:"nodeId,omitempty"`
	// Interface is the IP address of the associated interface.
	Interface string `json:"interface,omitempty"`
	// Service is the service name of the associated service.
	Service string `json:"service,omitempty"`
	// IfIndex is the SNMP ifIndex of the associated interface.
	IfIndex int `json:"ifIndex,omitempty"`
	// Description is an HTML-formatted description.
	Description string `json:"description,omitempty"`
	// LogMsg is the short log message text.
	LogMsg string `json:"logMsg,omitempty"`
	// OperInstruct holds operator instructions.
	OperInstruct string `json:"operInstruct,omitempty"`
	// Parms holds event parameters as a list of
	// {"parmName": ..., "value": ...} entries. The value may also
	// be a map: {"value": ..., "type": "string",
	// "encoding": "text"}.
	Parms []map[string]any `json:"parms,omitempty"`
}

// EventConfLogmsg is a log message entry within an EventConfEvent.
type EventConfLogmsg struct {
	// Content is the log message text.
	Content string `json:"content,omitempty"`
	// Dest is the destination: "logndisplay", "logonly",
	// "displayonly", "suppress", "donotpersist".
	Dest string `json:"dest,omitempty"`
}

// EventConfEvent is the event definition payload for
// Client.CreateEventconfEvent and Client.UpdateEventconfEvent.
type EventConfEvent struct {
	// UEI is the event UEI. Required.
	UEI string `json:"uei,omitempty"`
	// Label is the human-readable label.
	Label string `json:"label,omitempty"`
	// Descr is an HTML-formatted description.
	Descr string `json:"descr,omitempty"`
	// Logmsg is the log message configuration.
	Logmsg *EventConfLogmsg `json:"logmsg,omitempty"`
	// Severity is the event severity: "Indeterminate", "Cleared",
	// "Normal", "Warning", "Minor", "Major", "Critical".
	Severity string `json:"severity,omitempty"`
}

// MeasurementSource is a data source entry within a
// MeasurementsQuery.
type MeasurementSource struct {
	// ResourceID is the OpenNMS resource ID,
	// e.g. "node[1].interfaceSnmp[eth0-04013f75f101]". Required.
	ResourceID string `json:"resourceId,omitempty"`
	// Attribute is the RRD attribute name, e.g. "ifInOctets".
	// Required.
	Attribute string `json:"attribute,omitempty"`
	// Label is the column label used in expressions and output.
	// Required.
	Label string `json:"label,omitempty"`
	// Aggregation is the consolidation function: "AVERAGE"
	// (default), "MIN", "MAX", "LAST".
	Aggregation string `json:"aggregation,omitempty"`
	// Transient, when true, excludes this source from the
	// response (used only for intermediate expression values).
	Transient *bool `json:"transient,omitempty"`
}

// MeasurementExpression is a JEXL expression entry within a
// MeasurementsQuery.
type MeasurementExpression struct {
	// Label is the output column label. Required.
	Label string `json:"label,omitempty"`
	// Value is the JEXL expression referencing source labels,
	// e.g. "ifInOctets * 8". Required.
	Value string `json:"value,omitempty"`
	// Transient, when true, excludes this expression from the
	// response output.
	Transient *bool `json:"transient,omitempty"`
}

// MeasurementsQuery is the query payload for
// Client.GetMeasurementsMulti.
type MeasurementsQuery struct {
	// Start is the query start time as a Unix millisecond epoch
	// timestamp.
	Start int `json:"start,omitempty"`
	// End is the query end time as a Unix millisecond epoch
	// timestamp.
	End int `json:"end,omitempty"`
	// Step is the desired step size in milliseconds.
	Step int `json:"step,omitempty"`
	// Maxrows is the maximum number of time-step rows to return.
	Maxrows int `json:"maxrows,omitempty"`
	// Source lists the data sources to fetch.
	Source []MeasurementSource `json:"source,omitempty"`
	// Expression lists the JEXL expressions computed from source
	// labels.
	Expression []MeasurementExpression `json:"expression,omitempty"`
}
