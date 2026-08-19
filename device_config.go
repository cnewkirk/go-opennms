package opennms

// Device Configuration REST API – /rest/device-config.

import (
	"context"
	"fmt"
	"strconv"
	"strings"
)

// GetDeviceConfigs lists all device configurations (sorted by
// lastUpdated by default).
//
// deviceName filters by device hostname, ipAddress by IP address, and
// configType by config type string; "" skips each filter.
// createdAfter/createdBefore restrict to configs created after/before
// the given ms epoch; 0 skips each bound.
func (c *Client) GetDeviceConfigs(ctx context.Context, opts *ListOptions, deviceName, ipAddress, configType string, createdAfter, createdBefore int) (map[string]any, error) {
	params := listParams(opts)
	if deviceName != "" {
		params.Set("deviceName", deviceName)
	}
	if ipAddress != "" {
		params.Set("ipAddress", ipAddress)
	}
	if configType != "" {
		params.Set("configType", configType)
	}
	if createdAfter != 0 {
		params.Set("createdAfter", strconv.Itoa(createdAfter))
	}
	if createdBefore != 0 {
		params.Set("createdBefore", strconv.Itoa(createdBefore))
	}
	return c.getObject(ctx, "device-config", params, false)
}

// GetDeviceConfig returns the device configuration with the given
// database config ID.
func (c *Client) GetDeviceConfig(ctx context.Context, configID int) (map[string]any, error) {
	return c.getObject(ctx, fmt.Sprintf("device-config/%d", configID), nil, false)
}

// GetDeviceConfigByInterface returns all configs for a specific
// interface ID.
func (c *Client) GetDeviceConfigByInterface(ctx context.Context, interfaceID int) (map[string]any, error) {
	return c.getObject(ctx,
		fmt.Sprintf("device-config/interface/%d", interfaceID), nil, false)
}

// GetLatestDeviceConfigs returns the latest config for all devices.
//
// search is a search term for device name / IP; status filters by
// backup status. "" skips each filter.
func (c *Client) GetLatestDeviceConfigs(ctx context.Context, opts *ListOptions, search, status string) (map[string]any, error) {
	params := listParams(opts)
	if search != "" {
		params.Set("search", search)
	}
	if status != "" {
		params.Set("status", status)
	}
	return c.getObject(ctx, "device-config/latest", params, false)
}

// DownloadDeviceConfigs downloads configs for one or more config IDs
// and returns the raw config text.
func (c *Client) DownloadDeviceConfigs(ctx context.Context, configIDs []int) (string, error) {
	ids := make([]string, len(configIDs))
	for i, id := range configIDs {
		ids[i] = strconv.Itoa(id)
	}
	params := mergeFilters(nil, map[string]string{"id": strings.Join(ids, ",")})
	return asString(c.get(ctx, "device-config/download", params, false))
}

// BackupDeviceConfig triggers a backup retrieval for one or more
// interfaces.
//
// Each backup request map has keys "ipAddress" (interface IP address),
// "location" (monitoring location name, e.g. "Default"), "serviceName"
// (e.g. "DeviceConfig-default"), and "blocking" (whether to wait for
// the backup to complete).
func (c *Client) BackupDeviceConfig(ctx context.Context, backups []map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "device-config/backup", backups, nil, false))
}
