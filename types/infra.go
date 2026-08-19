package types

// Infrastructure: credentials, applications, user-defined links,
// SNMP configuration, and Grafana endpoints.

// Credential is the Secure Credentials Vault payload for
// Client.CreateCredential and Client.UpdateCredential.
type Credential struct {
	// Alias is the unique credential alias. Required for create.
	Alias string `json:"alias,omitempty"`
	// Username is the credential username.
	Username string `json:"username,omitempty"`
	// Password is the credential password.
	Password string `json:"password,omitempty"`
	// Attributes holds additional key/value attributes.
	Attributes map[string]string `json:"attributes,omitempty"`
}

// Application is the application payload for
// Client.CreateApplication.
type Application struct {
	// Name is the unique application name. Required.
	Name string `json:"name,omitempty"`
	// MonitoredServices lists monitored service ID references,
	// e.g. [{"id": 201}].
	MonitoredServices []map[string]any `json:"monitoredServices,omitempty"`
}

// UserDefinedLink is the user-defined link payload for
// Client.CreateUserDefinedLink.
type UserDefinedLink struct {
	// NodeIDA is the database ID of the A-side node. Required.
	NodeIDA int `json:"nodeIdA,omitempty"`
	// NodeIDZ is the database ID of the Z-side node. Required.
	NodeIDZ int `json:"nodeIdZ,omitempty"`
	// ComponentLabelA is the label for the A-side component
	// (e.g. interface name).
	ComponentLabelA string `json:"componentLabelA,omitempty"`
	// ComponentLabelZ is the label for the Z-side component.
	ComponentLabelZ string `json:"componentLabelZ,omitempty"`
	// LinkID is the unique link identifier string.
	LinkID string `json:"linkId,omitempty"`
	// LinkLabel is the human-readable link label.
	LinkLabel string `json:"linkLabel,omitempty"`
	// Owner is the username of the link owner.
	Owner string `json:"owner,omitempty"`
}

// SnmpConfig is the SNMP configuration payload for
// Client.SetSnmpConfig.
//
// Fields correspond to <definition/> attributes in snmp-config.xsd.
// Unset fields inherit from the OpenNMS default SNMP configuration.
type SnmpConfig struct {
	// Version is the SNMP version: "v1", "v2c", or "v3".
	Version string `json:"version,omitempty"`
	// Community is the community string (v1/v2c only).
	Community string `json:"community,omitempty"`
	// Port is the UDP port. Default 161.
	Port int `json:"port,omitempty"`
	// Timeout is the timeout in milliseconds.
	Timeout int `json:"timeout,omitempty"`
	// Retries is the retry count before giving up.
	Retries int `json:"retries,omitempty"`
	// MaxVarsPerPDU is the maximum variables per GETBULK PDU.
	MaxVarsPerPDU int `json:"maxVarsPerPdu,omitempty"`
	// MaxRepetitions is the maximum repetitions for GETBULK.
	MaxRepetitions int `json:"maxRepetitions,omitempty"`
	// MaxRequestSize is the maximum PDU size in bytes.
	MaxRequestSize int `json:"maxRequestSize,omitempty"`
	// ProxyHost is the proxy host IP for this definition.
	ProxyHost string `json:"proxyHost,omitempty"`
	// SecurityName is the SNMPv3 security name (username).
	SecurityName string `json:"securityName,omitempty"`
	// SecurityLevel is the SNMPv3 security level —
	// 1 noAuthNoPriv, 2 authNoPriv, 3 authPriv.
	SecurityLevel int `json:"securityLevel,omitempty"`
	// AuthPassphrase is the SNMPv3 authentication passphrase.
	AuthPassphrase string `json:"authPassphrase,omitempty"`
	// AuthProtocol is the SNMPv3 auth protocol: "MD5" or "SHA".
	AuthProtocol string `json:"authProtocol,omitempty"`
	// PrivPassphrase is the SNMPv3 privacy passphrase.
	PrivPassphrase string `json:"privPassphrase,omitempty"`
	// PrivProtocol is the SNMPv3 privacy protocol: "DES",
	// "AES128", "AES192", or "AES256".
	PrivProtocol string `json:"privProtocol,omitempty"`
	// ContextName is the SNMPv3 context name.
	ContextName string `json:"contextName,omitempty"`
}

// GrafanaEndpoint is the Grafana endpoint payload for
// Client.CreateGrafanaEndpoint, Client.UpdateGrafanaEndpoint, and
// Client.VerifyGrafanaEndpoint.
type GrafanaEndpoint struct {
	// ID is the numeric identifier (required when updating).
	ID int `json:"id,omitempty"`
	// UID is the unique identifier of the endpoint. Required.
	UID string `json:"uid,omitempty"`
	// URL is the base URL of the Grafana instance. Required.
	URL string `json:"url,omitempty"`
	// APIKey is the Grafana API key. Required.
	APIKey string `json:"apiKey,omitempty"`
	// Description is an optional description.
	Description string `json:"description,omitempty"`
	// ConnectTimeout is the connect timeout in milliseconds.
	ConnectTimeout int `json:"connectTimeout,omitempty"`
	// ReadTimeout is the read timeout in milliseconds.
	ReadTimeout int `json:"readTimeout,omitempty"`
}
