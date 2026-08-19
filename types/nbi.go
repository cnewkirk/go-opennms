package types

// Northbound interface (NBI) and Javamail configurations.

// SnmpTrapNbiTrapSink is the SNMP trap sink payload for
// Client.CreateSnmptrapNbiTrapsink and
// Client.UpdateSnmptrapNbiTrapsink.
type SnmpTrapNbiTrapSink struct {
	// Name is the unique trap sink name. Required.
	Name string `json:"name,omitempty"`
	// IPAddress is the destination IP address. Required.
	IPAddress string `json:"ipAddress,omitempty"`
	// Port is the destination UDP port. Default 162.
	Port int `json:"port,omitempty"`
	// Community is the SNMP community string. Default "public".
	Community string `json:"community,omitempty"`
}

// SnmpTrapNbiConfig is the SNMP trap NBI configuration payload for
// Client.UpdateSnmptrapNbiConfig.
type SnmpTrapNbiConfig struct {
	// Enabled reports whether SNMP trap forwarding is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Trapsinks lists the trap sink definitions.
	Trapsinks []SnmpTrapNbiTrapSink `json:"trapsinks,omitempty"`
}

// EmailNbiDestination is the email NBI destination payload for
// Client.CreateEmailNbiDestination and
// Client.UpdateEmailNbiDestination.
type EmailNbiDestination struct {
	// Name is the unique destination name. Required.
	Name string `json:"name,omitempty"`
	// FirstOccurrenceOnly, when true, sends only on first
	// occurrence of an alarm (not on re-occurrences).
	FirstOccurrenceOnly *bool `json:"firstOccurrenceOnly,omitempty"`
	// Filters lists alarm filter expressions.
	Filters []string `json:"filters,omitempty"`
}

// EmailNbiConfig is the email NBI configuration payload for
// Client.UpdateEmailNbiConfig.
type EmailNbiConfig struct {
	// Enabled reports whether email alarm forwarding is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Destinations lists the destination definitions.
	Destinations []EmailNbiDestination `json:"destinations,omitempty"`
}

// SyslogNbiDestination is the syslog NBI destination payload for
// Client.CreateSyslogNbiDestination and
// Client.UpdateSyslogNbiDestination.
type SyslogNbiDestination struct {
	// Name is the unique destination name. Required.
	Name string `json:"name,omitempty"`
	// Host is the destination hostname or IP address. Required.
	Host string `json:"host,omitempty"`
	// Port is the destination UDP/TCP port. Default 514.
	Port int `json:"port,omitempty"`
	// FirstOccurrenceOnly, when true, forwards only on first
	// occurrence.
	FirstOccurrenceOnly *bool `json:"firstOccurrenceOnly,omitempty"`
	// Filters lists alarm filter expressions.
	Filters []string `json:"filters,omitempty"`
}

// SyslogNbiConfig is the syslog NBI configuration payload for
// Client.UpdateSyslogNbiConfig.
type SyslogNbiConfig struct {
	// Enabled reports whether syslog alarm forwarding is enabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Destinations lists the destination definitions.
	Destinations []SyslogNbiDestination `json:"destinations,omitempty"`
}

// JavamailDefaultConfig is the default Javamail configuration
// payload for Client.SetJavamailDefaultConfig.
type JavamailDefaultConfig struct {
	// DefaultReadConfigName is the name of the default read-mail
	// configuration.
	DefaultReadConfigName string `json:"defaultReadConfigName,omitempty"`
	// DefaultSendConfigName is the name of the default send-mail
	// configuration.
	DefaultSendConfigName string `json:"defaultSendConfigName,omitempty"`
	// DefaultEnd2endConfigName is the name of the default
	// end-to-end test configuration.
	DefaultEnd2endConfigName string `json:"defaultEnd2endConfigName,omitempty"`
}

// JavamailReadmail is the read-mail configuration payload for
// Client.CreateJavamailReadmail and Client.UpdateJavamailReadmail.
type JavamailReadmail struct {
	// Name is the unique configuration name. Required.
	Name string `json:"name,omitempty"`
	// Host is the mail server hostname.
	Host string `json:"host,omitempty"`
	// Port is the mail server port (e.g. 993 for IMAPS).
	Port int `json:"port,omitempty"`
	// Protocol is the mail protocol: "imap", "imaps", "pop3",
	// "pop3s".
	Protocol string `json:"protocol,omitempty"`
}

// JavamailSendmail is the send-mail configuration payload for
// Client.CreateJavamailSendmail and Client.UpdateJavamailSendmail.
type JavamailSendmail struct {
	// Name is the unique configuration name. Required.
	Name string `json:"name,omitempty"`
	// Host is the SMTP server hostname.
	Host string `json:"host,omitempty"`
	// Port is the SMTP server port (e.g. 25, 465, 587).
	Port int `json:"port,omitempty"`
	// Protocol is the SMTP protocol: "smtp", "smtps".
	Protocol string `json:"protocol,omitempty"`
}

// JavamailEnd2End is the end-to-end mail test configuration payload
// for Client.CreateJavamailEnd2end and Client.UpdateJavamailEnd2end.
type JavamailEnd2End struct {
	// Name is the unique configuration name. Required.
	Name string `json:"name,omitempty"`
	// ReadMailConfigName is the read-mail configuration to use.
	ReadMailConfigName string `json:"readMailConfigName,omitempty"`
	// SendMailConfigName is the send-mail configuration to use.
	SendMailConfigName string `json:"sendMailConfigName,omitempty"`
}
