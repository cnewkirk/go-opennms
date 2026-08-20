package opennms

// Foreign Sources REST API – /rest/foreignSources.

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// foreignSourcesEscape escapes each /-separated segment of a foreign
// source name, preserving slashes so the "deployed/<name>" namespace
// stays addressable (as it is in python-opennms).
func foreignSourcesEscape(name string) string {
	segments := strings.Split(name, "/")
	for i, s := range segments {
		segments[i] = url.PathEscape(s)
	}
	return strings.Join(segments, "/")
}

// Foreign sources

// GetForeignSources lists all active (pending + deployed) foreign
// sources.
func (c *Client) GetForeignSources(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources", nil, false)
}

// GetForeignSource returns a specific foreign source by name.
func (c *Client) GetForeignSource(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources/"+foreignSourcesEscape(name), nil, false)
}

// GetDefaultForeignSource returns the default foreign source
// definition.
func (c *Client) GetDefaultForeignSource(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources/default", nil, false)
}

// GetDeployedForeignSources lists all deployed foreign sources.
func (c *Client) GetDeployedForeignSources(ctx context.Context) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources/deployed", nil, false)
}

// GetDeployedForeignSourceCount returns the count of deployed foreign
// sources.
func (c *Client) GetDeployedForeignSourceCount(ctx context.Context) (int, error) {
	return c.getCount(ctx, "foreignSources/deployed/count", false)
}

// CreateForeignSource creates a new foreign source.
//
// foreignSource is the definition map, e.g.:
//
//	map[string]any{
//	    "name":          "Servers",
//	    "scan-interval": "1d",
//	    "detectors": []any{
//	        map[string]any{
//	            "name":      "ICMP",
//	            "class":     "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
//	            "parameter": []any{},
//	        },
//	    },
//	    "policies": []any{},
//	}
func (c *Client) CreateForeignSource(ctx context.Context, foreignSource map[string]any) (map[string]any, error) {
	return asObject(c.post(ctx, "foreignSources", foreignSource, nil, false))
}

// UpdateForeignSource updates an existing foreign source. The server
// accepts this update form-encoded only; the foreignSource values are
// stringified into form fields (e.g. "scan-interval": "12h").
func (c *Client) UpdateForeignSource(ctx context.Context, name string, foreignSource map[string]any) error {
	form := url.Values{}
	for k, v := range foreignSource {
		form.Set(k, fmt.Sprint(v))
	}
	_, err := c.putForm(ctx, "foreignSources/"+foreignSourcesEscape(name), form, nil, false)
	return err
}

// DeleteForeignSource deletes a foreign source.
func (c *Client) DeleteForeignSource(ctx context.Context, name string) error {
	_, err := c.del(ctx, "foreignSources/"+foreignSourcesEscape(name), nil, nil, false, "")
	return err
}

// Detectors

// GetForeignSourceDetectors lists detectors for the foreign source.
func (c *Client) GetForeignSourceDetectors(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources/"+foreignSourcesEscape(name)+"/detectors", nil, false)
}

// GetForeignSourceDetector returns a specific detector from the
// foreign source.
func (c *Client) GetForeignSourceDetector(ctx context.Context, name, detector string) (map[string]any, error) {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/detectors/" + url.PathEscape(detector)
	return c.getObject(ctx, path, nil, false)
}

// AddForeignSourceDetector adds a detector to the foreign source.
//
// detector is the definition map, e.g.:
//
//	map[string]any{
//	    "name":      "HTTP",
//	    "class":     "org.opennms.netmgt.provision.detector.web.HttpDetector",
//	    "parameter": []any{},
//	}
func (c *Client) AddForeignSourceDetector(ctx context.Context, name string, detector map[string]any) (map[string]any, error) {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/detectors"
	return asObject(c.post(ctx, path, detector, nil, false))
}

// DeleteForeignSourceDetector removes a detector from the foreign
// source.
func (c *Client) DeleteForeignSourceDetector(ctx context.Context, name, detector string) error {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/detectors/" + url.PathEscape(detector)
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}

// Policies

// GetForeignSourcePolicies lists policies for the foreign source.
func (c *Client) GetForeignSourcePolicies(ctx context.Context, name string) (map[string]any, error) {
	return c.getObject(ctx, "foreignSources/"+foreignSourcesEscape(name)+"/policies", nil, false)
}

// GetForeignSourcePolicy returns a specific policy from the foreign
// source.
func (c *Client) GetForeignSourcePolicy(ctx context.Context, name, policy string) (map[string]any, error) {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/policies/" + url.PathEscape(policy)
	return c.getObject(ctx, path, nil, false)
}

// AddForeignSourcePolicy adds a policy to the foreign source.
//
// policy is the definition map, e.g.:
//
//	map[string]any{
//	    "name":  "Do Not Persist Discovered IPs",
//	    "class": "org.opennms.netmgt.provision.persist.policies.MatchingIpInterfacePolicy",
//	    "parameter": []any{
//	        map[string]any{"key": "action", "value": "DO_NOT_PERSIST"},
//	    },
//	}
func (c *Client) AddForeignSourcePolicy(ctx context.Context, name string, policy map[string]any) (map[string]any, error) {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/policies"
	return asObject(c.post(ctx, path, policy, nil, false))
}

// DeleteForeignSourcePolicy removes a policy from the foreign source.
func (c *Client) DeleteForeignSourcePolicy(ctx context.Context, name, policy string) error {
	path := "foreignSources/" + foreignSourcesEscape(name) + "/policies/" + url.PathEscape(policy)
	_, err := c.del(ctx, path, nil, nil, false, "")
	return err
}
