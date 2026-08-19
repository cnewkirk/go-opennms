package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	opennms "github.com/cnewkirk/go-opennms"
)

// writeTag returns a per-suite unique tag matching the Python
// smoke test's f"smoke-{int(time.time())}".
func writeTag() string {
	return fmt.Sprintf("smoke-%d", time.Now().Unix())
}

// writeCleanup runs a best-effort cleanup; errors are ignored and it
// never affects pass/fail counts.
func writeCleanup(fn func() error) {
	_ = fn()
}

// writeIDString renders a decoded JSON id (string or number) as a
// string, or "" when absent.
func writeIDString(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case float64:
		return strconv.Itoa(int(x))
	case int:
		return strconv.Itoa(x)
	}
	return ""
}

// writeIDInt converts a decoded JSON id (string or number) to int,
// or 0 when absent.
func writeIDInt(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	case string:
		n, _ := strconv.Atoi(x)
		return n
	}
	return 0
}

// writeBoolPtr returns a pointer to v when it is a JSON bool, else
// nil (mirroring Python's "omit when absent" kwargs).
func writeBoolPtr(v any) *bool {
	if b, ok := v.(bool); ok {
		return &b
	}
	return nil
}

// runWriteChecks exercises write endpoints (create/update/delete),
// cleaning up everything it creates.
func runWriteChecks(ctx context.Context, c *opennms.Client, r *runner) {
	writeOps(ctx, c, r)
	writeNodeLifecycle(ctx, c, r)
	writeIdentity(ctx, c, r)
	writeSchedOutageAssoc(ctx, c, r)
	writeProvisioning(ctx, c, r)
	writeMonitoringLocations(ctx, c, r)
	writeServiceEntities(ctx, c, r)
	writeClassificationScv(ctx, c, r)
	writeConfigsNbi(ctx, c, r)
	writeReportsGraphmlMisc(ctx, c, r)
}

func writeOps(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section(fmt.Sprintf("write operations  [tag: %s]", tag))

	// Events ── fire-and-forget; no cleanup needed
	r.run("create_event (internal test UEI)", nil, func() (any, error) {
		return nil, c.CreateEvent(ctx, map[string]any{
			"uei":      "uei.opennms.org/internal/test",
			"source":   "smoke_test.py",
			"severity": "Normal",
			"parms": []map[string]any{
				{"parmName": "smoke-tag", "value": tag},
			},
		})
	})

	// Alarms ── ack then immediately unack; only if an unacked alarm exists
	_, aid := first(func(limit int) (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "alarm", "")
	if alarmID := writeIDInt(aid); alarmID != 0 {
		alreadyAcked := true // play it safe
		if alarm, err := c.GetAlarm(ctx, alarmID); err == nil {
			alreadyAcked = alarm != nil && alarm["ackUser"] != nil
		}
		if !alreadyAcked {
			r.run(fmt.Sprintf("ack_alarm    id=%d", alarmID), nil, func() (any, error) {
				return nil, c.AckAlarm(ctx, alarmID, "")
			})
			r.run(fmt.Sprintf("unack_alarm  id=%d", alarmID), nil, func() (any, error) {
				return nil, c.UnackAlarm(ctx, alarmID)
			})
		} else {
			r.skip(fmt.Sprintf("ack_alarm / unack_alarm  id=%d", alarmID),
				"already acknowledged – skipping to avoid side-effects")
		}
	} else {
		r.skip("ack_alarm / unack_alarm", "no alarms")
	}

	// Categories ── create / get / delete
	catName := fmt.Sprintf("Smoke-Test-%s", tag)
	if _, ok := r.run(fmt.Sprintf("create_category  (%s)", catName), nil, func() (any, error) {
		return c.CreateCategory(ctx, map[string]any{"name": catName})
	}); ok {
		r.run(fmt.Sprintf("get_category     (%s)", catName), nil, func() (any, error) {
			return c.GetCategory(ctx, catName)
		})
		r.run(fmt.Sprintf("delete_category  (%s)", catName), nil, func() (any, error) {
			return nil, c.DeleteCategory(ctx, catName)
		})
	} else {
		r.skip(fmt.Sprintf("get_category / delete_category  (%s)", catName),
			"create failed")
	}

	// Groups ── create / get / delete
	grpName := fmt.Sprintf("smoke-test-%s", tag)
	if _, ok := r.run(fmt.Sprintf("create_group  (%s)", grpName), nil, func() (any, error) {
		return c.CreateGroup(ctx, map[string]any{
			"name": grpName, "comments": "smoke test – safe to delete",
		})
	}); ok {
		r.run(fmt.Sprintf("get_group     (%s)", grpName), nil, func() (any, error) {
			return c.GetGroup(ctx, grpName)
		})
		r.run(fmt.Sprintf("delete_group  (%s)", grpName), nil, func() (any, error) {
			return nil, c.DeleteGroup(ctx, grpName)
		})
	} else {
		r.skip(fmt.Sprintf("get_group / delete_group  (%s)", grpName),
			"create failed")
	}

	// Scheduled outages ── create / get / delete
	soName := fmt.Sprintf("smoke-test-%s", tag)
	if _, ok := r.run(fmt.Sprintf("create_sched_outage  (%s)", soName), nil, func() (any, error) {
		return c.CreateSchedOutage(ctx, map[string]any{
			"name": soName,
			"type": "specific",
			"time": []map[string]any{
				{"begins": "01-Jan-2000 00:00:00", "ends": "01-Jan-2000 00:00:01"},
			},
		})
	}); ok {
		r.run(fmt.Sprintf("get_sched_outage     (%s)", soName), nil, func() (any, error) {
			return c.GetSchedOutage(ctx, soName)
		})
		r.run(fmt.Sprintf("delete_sched_outage  (%s)", soName), nil, func() (any, error) {
			return nil, c.DeleteSchedOutage(ctx, soName)
		})
	} else {
		r.skip(fmt.Sprintf("get_sched_outage / delete_sched_outage  (%s)", soName),
			"create failed")
	}

	// Requisitions ── create / get / delete (no import, so no real nodes created)
	reqName := fmt.Sprintf("smoke-test-%s", tag)
	if _, ok := r.run(fmt.Sprintf("create_requisition  (%s)", reqName), nil, func() (any, error) {
		return nil, c.CreateRequisition(ctx, map[string]any{
			"foreign-source": reqName, "node": []any{},
		})
	}); ok {
		r.run(fmt.Sprintf("get_requisition     (%s)", reqName), nil, func() (any, error) {
			return c.GetRequisition(ctx, reqName)
		})
		r.run(fmt.Sprintf("delete_requisition  (%s)", reqName), nil, func() (any, error) {
			return nil, c.DeleteRequisition(ctx, reqName)
		})
	} else {
		r.skip(fmt.Sprintf("get_requisition / delete_requisition  (%s)", reqName),
			"create failed")
	}

	// Maps ── create / update / delete (removed in OpenNMS Horizon 16;
	// only works on pre-16 servers)
	res, _ := r.warn("create_map", "maps REST API removed in OpenNMS Horizon 16",
		nil, func() (any, error) {
			return c.CreateMap(ctx, map[string]any{
				"name":     fmt.Sprintf("Smoke Test %s", tag),
				"mapWidth": 1920, "mapHeight": 1080,
			})
		})
	mid := 0
	if m, ok := res.(map[string]any); ok && m != nil {
		mid = writeIDInt(m["id"])
	}
	if mid != 0 {
		r.run(fmt.Sprintf("update_map  id=%d", mid), nil, func() (any, error) {
			return nil, c.UpdateMap(ctx, mid, map[string]any{
				"name":     fmt.Sprintf("Smoke Test %s (updated)", tag),
				"mapWidth": 1920, "mapHeight": 1080,
			})
		})
		r.run(fmt.Sprintf("delete_map  id=%d", mid), nil, func() (any, error) {
			return nil, c.DeleteMap(ctx, mid)
		})
	} else {
		r.skip("update_map / delete_map", "maps API unavailable")
	}
}

func writeNodeLifecycle(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: node lifecycle")

	// Rescan the stable self-monitor node; scanning (or emptying) a
	// throwaway REST-created node triggers async delete propagation.
	r.run("rescan_node  id=1", nil, func() (any, error) {
		return nil, c.RescanNode(ctx, "1")
	})

	r.run("create_node", nil, func() (any, error) {
		return c.CreateNode(ctx, map[string]any{
			"label": fmt.Sprintf("%s-node", tag), "type": "A", "location": "Default",
		})
	})
	nid := ""
	if nodes, err := c.GetNodes(ctx, &opennms.ListOptions{}, nil); err == nil && nodes != nil {
		items, _ := nodes["node"].([]any)
		for _, item := range items {
			node, _ := item.(map[string]any)
			if node != nil && node["label"] == fmt.Sprintf("%s-node", tag) {
				nid = writeIDString(node["id"])
			}
		}
	}
	if nid == "" {
		r.skip("node lifecycle suite", "create_node produced no node")
		return
	}

	r.run(fmt.Sprintf("update_node  id=%s", nid), nil, func() (any, error) {
		return nil, c.UpdateNode(ctx, nid, map[string]string{
			"label": fmt.Sprintf("%s-node", tag),
		})
	})

	ip, ip2 := "10.254.0.1", "10.254.0.2"
	r.run("create_node_ip_interface", nil, func() (any, error) {
		return c.CreateNodeIpInterface(ctx, nid, map[string]any{
			"ipAddress": ip, "isManaged": "M", "snmpPrimary": "N",
		})
	})
	r.run("update_node_ip_interface", nil, func() (any, error) {
		return nil, c.UpdateNodeIpInterface(ctx, nid, ip,
			map[string]string{"isManaged": "M"})
	})
	r.run("create_node_ip_interface (second)", nil, func() (any, error) {
		return c.CreateNodeIpInterface(ctx, nid, map[string]any{
			"ipAddress": ip2, "isManaged": "M", "snmpPrimary": "N",
		})
	})

	r.run("create_node_snmp_interface", nil, func() (any, error) {
		return c.CreateNodeSnmpInterface(ctx, nid, map[string]any{
			"ifIndex": 991, "ifName": fmt.Sprintf("%s0", tag), "ifType": 6,
		})
	})
	r.run("update_node_snmp_interface", nil, func() (any, error) {
		return nil, c.UpdateNodeSnmpInterface(ctx, nid, 991,
			map[string]string{"ifAlias": tag})
	})
	r.run("delete_node_snmp_interface", nil, func() (any, error) {
		return nil, c.DeleteNodeSnmpInterface(ctx, nid, 991)
	})

	cat := fmt.Sprintf("%s-cat", tag)
	r.run("create_category", nil, func() (any, error) {
		return c.CreateCategory(ctx, map[string]any{"name": cat})
	})
	r.run("add_node_category", nil, func() (any, error) {
		return c.AddNodeCategory(ctx, nid, map[string]any{"name": cat})
	})
	r.run("update_node_category", nil, func() (any, error) {
		return nil, c.UpdateNodeCategory(ctx, nid, cat,
			map[string]string{"name": cat})
	})
	r.run("delete_node_category", nil, func() (any, error) {
		return nil, c.DeleteNodeCategory(ctx, nid, cat)
	})
	r.run("associate_category_with_node", nil, func() (any, error) {
		return nil, c.AssociateCategoryWithNode(ctx, cat, writeIDInt(nid))
	})
	r.run("dissociate_category_from_node", nil, func() (any, error) {
		return nil, c.DissociateCategoryFromNode(ctx, cat, writeIDInt(nid))
	})
	r.run("update_category", nil, func() (any, error) {
		return nil, c.UpdateCategory(ctx, cat,
			map[string]string{"description": "smoke"})
	})

	r.run("update_node_asset_record", nil, func() (any, error) {
		return nil, c.UpdateNodeAssetRecord(ctx, nid,
			map[string]string{"building": tag})
	})

	r.warn("add_node_hardware_inventory",
		"hardware inventory root entity may need SNMP data",
		nil, func() (any, error) {
			return c.AddNodeHardwareInventory(ctx, nid, map[string]any{
				"entPhysicalIndex": 1, "entPhysicalName": tag,
				"entPhysicalClass": 3,
			})
		})
	r.warn("update_node_hardware_entity", "requires the entity created above",
		nil, func() (any, error) {
			return nil, c.UpdateNodeHardwareEntity(ctx, nid, 1,
				map[string]string{"entPhysicalAlias": tag})
		})
	r.warn("delete_node_hardware_entity", "requires the entity created above",
		nil, func() (any, error) {
			return nil, c.DeleteNodeHardwareEntity(ctx, nid, 1)
		})

	meta := []map[string]any{{"context": "X-smoke", "key": "k1", "value": "v1"}}
	r.run("set_node_metadata", nil, func() (any, error) {
		return nil, c.SetNodeMetadata(ctx, nid, meta)
	})
	r.run("set_node_metadata_value", nil, func() (any, error) {
		return nil, c.SetNodeMetadataValue(ctx, nid, "X-smoke", "k2", "v2")
	})
	r.run("delete_node_metadata_key", nil, func() (any, error) {
		return nil, c.DeleteNodeMetadataKey(ctx, nid, "X-smoke", "k2")
	})
	r.run("delete_node_metadata_context", nil, func() (any, error) {
		return nil, c.DeleteNodeMetadataContext(ctx, nid, "X-smoke")
	})
	r.run("set_interface_metadata", nil, func() (any, error) {
		return nil, c.SetInterfaceMetadata(ctx, nid, ip, meta)
	})
	r.run("set_interface_metadata_value", nil, func() (any, error) {
		return nil, c.SetInterfaceMetadataValue(ctx, nid, ip, "X-smoke", "k2", "v2")
	})
	r.run("delete_interface_metadata_key", nil, func() (any, error) {
		return nil, c.DeleteInterfaceMetadataKey(ctx, nid, ip, "X-smoke", "k2")
	})
	r.run("delete_interface_metadata_context", nil, func() (any, error) {
		return nil, c.DeleteInterfaceMetadataContext(ctx, nid, ip, "X-smoke")
	})

	// Deleting the last service (or interface) of a node starts
	// async delete propagation, so keep a second service and a second
	// interface alive until the node itself is deleted.
	r.run("create_node_ip_service", nil, func() (any, error) {
		return c.CreateNodeIpService(ctx, nid, ip, map[string]any{
			"serviceType": map[string]any{"name": "ICMP"}, "status": "A",
		})
	})
	r.run("create_node_ip_service (second)", nil, func() (any, error) {
		return c.CreateNodeIpService(ctx, nid, ip, map[string]any{
			"serviceType": map[string]any{"name": "SNMP"}, "status": "A",
		})
	})
	r.run("set_service_metadata", nil, func() (any, error) {
		return nil, c.SetServiceMetadata(ctx, nid, ip, "ICMP", meta)
	})
	r.run("set_service_metadata_value", nil, func() (any, error) {
		return nil, c.SetServiceMetadataValue(ctx, nid, ip, "ICMP", "X-smoke", "k2", "v2")
	})
	r.run("delete_service_metadata_key", nil, func() (any, error) {
		return nil, c.DeleteServiceMetadataKey(ctx, nid, ip, "ICMP", "X-smoke", "k2")
	})
	r.run("delete_service_metadata_context", nil, func() (any, error) {
		return nil, c.DeleteServiceMetadataContext(ctx, nid, ip, "ICMP", "X-smoke")
	})
	r.run("delete_node_ip_service", nil, func() (any, error) {
		return nil, c.DeleteNodeIpService(ctx, nid, ip, "ICMP")
	})
	r.run("delete_node_ip_interface (second)", nil, func() (any, error) {
		return nil, c.DeleteNodeIpInterface(ctx, nid, ip2)
	})
	r.run(fmt.Sprintf("delete_node  id=%s", nid), nil, func() (any, error) {
		return nil, c.DeleteNode(ctx, nid)
	})
	r.run("delete_category", nil, func() (any, error) {
		return nil, c.DeleteCategory(ctx, cat)
	})
}

func writeIdentity(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: users, groups, roles")

	user := fmt.Sprintf("%s-user", tag)
	grp := fmt.Sprintf("%s-grp", tag)
	cat := fmt.Sprintf("%s-gcat", tag)
	r.run("create_user", nil, func() (any, error) {
		return c.CreateUser(ctx, map[string]any{
			"user-id": user, "password": "smoke-pw",
		}, true)
	})
	r.run("update_user", nil, func() (any, error) {
		return nil, c.UpdateUser(ctx, user,
			map[string]string{"fullName": "Smoke User"})
	})
	r.run("assign_role_to_user", nil, func() (any, error) {
		return nil, c.AssignRoleToUser(ctx, user, "ROLE_READONLY")
	})
	r.run("revoke_role_from_user", nil, func() (any, error) {
		return nil, c.RevokeRoleFromUser(ctx, user, "ROLE_READONLY")
	})
	r.run("create_group", nil, func() (any, error) {
		return c.CreateGroup(ctx, map[string]any{"name": grp})
	})
	r.run("update_group", nil, func() (any, error) {
		return nil, c.UpdateGroup(ctx, grp, map[string]string{"comments": "smoke"})
	})
	r.run("add_user_to_group", nil, func() (any, error) {
		return nil, c.AddUserToGroup(ctx, grp, user)
	})
	r.run("remove_user_from_group", nil, func() (any, error) {
		return nil, c.RemoveUserFromGroup(ctx, grp, user)
	})
	r.run("create_category (group assoc)", nil, func() (any, error) {
		return c.CreateCategory(ctx, map[string]any{"name": cat})
	})
	r.run("add_category_to_group", nil, func() (any, error) {
		return nil, c.AddCategoryToGroup(ctx, grp, cat)
	})
	r.run("remove_category_from_group", nil, func() (any, error) {
		return nil, c.RemoveCategoryFromGroup(ctx, grp, cat)
	})
	r.run("associate_category_with_group", nil, func() (any, error) {
		return nil, c.AssociateCategoryWithGroup(ctx, cat, grp)
	})
	r.run("dissociate_category_from_group", nil, func() (any, error) {
		return nil, c.DissociateCategoryFromGroup(ctx, cat, grp)
	})
	r.run("delete_category (group assoc)", nil, func() (any, error) {
		return nil, c.DeleteCategory(ctx, cat)
	})
	r.run("delete_group", nil, func() (any, error) {
		return nil, c.DeleteGroup(ctx, grp)
	})
	r.run("delete_user", nil, func() (any, error) {
		return nil, c.DeleteUser(ctx, user)
	})
}

func writeSchedOutageAssoc(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: scheduled outage daemon associations")

	name := fmt.Sprintf("%s-outage", tag)
	r.run("create_sched_outage", nil, func() (any, error) {
		return c.CreateSchedOutage(ctx, map[string]any{
			"name": name, "type": "specific",
			"time": []map[string]any{
				{"begins": "01-Jan-2030 00:00:00", "ends": "01-Jan-2030 01:00:00"},
			},
			"interface": []map[string]any{{"address": "10.254.0.99"}},
		})
	})
	r.run("associate_sched_outage_notifd", nil, func() (any, error) {
		return nil, c.AssociateSchedOutageNotifd(ctx, name)
	})
	r.run("dissociate_sched_outage_notifd", nil, func() (any, error) {
		return nil, c.DissociateSchedOutageNotifd(ctx, name)
	})
	r.warn("associate_sched_outage_collectd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.AssociateSchedOutageCollectd(ctx, name, "example1")
		})
	r.warn("dissociate_sched_outage_collectd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.DissociateSchedOutageCollectd(ctx, name, "example1")
		})
	r.warn("associate_sched_outage_pollerd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.AssociateSchedOutagePollerd(ctx, name, "example1")
		})
	r.warn("dissociate_sched_outage_pollerd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.DissociateSchedOutagePollerd(ctx, name, "example1")
		})
	r.warn("associate_sched_outage_threshd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.AssociateSchedOutageThreshd(ctx, name, "example1")
		})
	r.warn("dissociate_sched_outage_threshd", "package name is config-defined",
		nil, func() (any, error) {
			return nil, c.DissociateSchedOutageThreshd(ctx, name, "example1")
		})
	r.run("delete_sched_outage", nil, func() (any, error) {
		return nil, c.DeleteSchedOutage(ctx, name)
	})
}

func writeProvisioning(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: requisitions and foreign sources")

	req := fmt.Sprintf("%s-req", tag)
	fid := fmt.Sprintf("%s-n1", tag)
	ip := "10.254.1.1"
	r.run("create_requisition", nil, func() (any, error) {
		return nil, c.CreateRequisition(ctx, map[string]any{"foreign-source": req})
	})
	r.run("create_requisition_node", nil, func() (any, error) {
		return nil, c.CreateRequisitionNode(ctx, req, map[string]any{
			"foreign-id": fid, "node-label": fmt.Sprintf("%s-rnode", tag),
		})
	})
	r.run("create_requisition_node_interface", nil, func() (any, error) {
		return nil, c.CreateRequisitionNodeInterface(ctx, req, fid, map[string]any{
			"ip-addr": ip, "status": 1, "snmp-primary": "N",
		})
	})
	r.run("create_requisition_node_service", nil, func() (any, error) {
		return nil, c.CreateRequisitionNodeService(ctx, req, fid, ip,
			map[string]any{"service-name": "ICMP"})
	})
	r.run("add_requisition_node_category", nil, func() (any, error) {
		return nil, c.AddRequisitionNodeCategory(ctx, req, fid,
			map[string]any{"name": "Production"})
	})
	r.run("set_requisition_node_asset", nil, func() (any, error) {
		return nil, c.SetRequisitionNodeAsset(ctx, req, fid,
			map[string]any{"name": "building", "value": tag})
	})
	r.run("update_requisition", nil, func() (any, error) {
		return nil, c.UpdateRequisition(ctx, req,
			map[string]string{"foreign-source": req})
	})
	r.run("update_requisition_node", nil, func() (any, error) {
		return nil, c.UpdateRequisitionNode(ctx, req, fid,
			map[string]string{"node-label": fmt.Sprintf("%s-rnode2", tag)})
	})
	r.run("update_requisition_node_interface", nil, func() (any, error) {
		return nil, c.UpdateRequisitionNodeInterface(ctx, req, fid, ip,
			map[string]string{"descr": "smoke"})
	})
	r.run("delete_requisition_node_service", nil, func() (any, error) {
		return nil, c.DeleteRequisitionNodeService(ctx, req, fid, ip, "ICMP")
	})
	r.run("delete_requisition_node_category", nil, func() (any, error) {
		return nil, c.DeleteRequisitionNodeCategory(ctx, req, fid, "Production")
	})
	r.run("delete_requisition_node_asset", nil, func() (any, error) {
		return nil, c.DeleteRequisitionNodeAsset(ctx, req, fid, "building")
	})
	r.run("delete_requisition_node_interface", nil, func() (any, error) {
		return nil, c.DeleteRequisitionNodeInterface(ctx, req, fid, ip)
	})
	r.run("delete_requisition_node", nil, func() (any, error) {
		return nil, c.DeleteRequisitionNode(ctx, req, fid)
	})
	r.run("import_requisition", nil, func() (any, error) {
		return nil, c.ImportRequisition(ctx, req, true)
	})

	fs := fmt.Sprintf("%s-req", tag)
	r.run("create_foreign_source", nil, func() (any, error) {
		return c.CreateForeignSource(ctx, map[string]any{
			"name": fs, "scan-interval": "12w",
		})
	})
	r.run("update_foreign_source", nil, func() (any, error) {
		return nil, c.UpdateForeignSource(ctx, fs,
			map[string]any{"scan-interval": "6w"})
	})
	r.run("add_foreign_source_detector", nil, func() (any, error) {
		return c.AddForeignSourceDetector(ctx, fs, map[string]any{
			"name":  "ICMP",
			"class": "org.opennms.netmgt.provision.detector.icmp.IcmpDetector",
		})
	})
	r.run("delete_foreign_source_detector", nil, func() (any, error) {
		return nil, c.DeleteForeignSourceDetector(ctx, fs, "ICMP")
	})
	r.run("add_foreign_source_policy", nil, func() (any, error) {
		return c.AddForeignSourcePolicy(ctx, fs, map[string]any{
			"name": "no-discovered-ips",
			"class": "org.opennms.netmgt.provision.persist.policies" +
				".MatchingIpInterfacePolicy",
			"parameter": []map[string]any{
				{"key": "action", "value": "DO_NOT_PERSIST"},
				{"key": "matchBehavior", "value": "NO_PARAMETERS"},
			},
		})
	})
	r.run("delete_foreign_source_policy", nil, func() (any, error) {
		return nil, c.DeleteForeignSourcePolicy(ctx, fs, "no-discovered-ips")
	})
	r.run("delete_foreign_source", nil, func() (any, error) {
		return nil, c.DeleteForeignSource(ctx, fs)
	})
	writeCleanup(func() error {
		return c.DeleteForeignSource(ctx, "deployed/"+fs)
	})
	r.run("delete_requisition", nil, func() (any, error) {
		return nil, c.DeleteRequisition(ctx, req)
	})
	r.run("delete_deployed_requisition", nil, func() (any, error) {
		return nil, c.DeleteDeployedRequisition(ctx, req)
	})
}

func writeMonitoringLocations(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: monitoring locations")

	loc := fmt.Sprintf("%s-loc", tag)
	r.run("create_monitoring_location", nil, func() (any, error) {
		return nil, c.CreateMonitoringLocation(ctx, map[string]any{
			"location-name": loc, "monitoring-area": "smoke",
		})
	})
	r.run("update_monitoring_location", nil, func() (any, error) {
		return nil, c.UpdateMonitoringLocation(ctx, loc,
			map[string]string{"monitoring-area": "smoke2"})
	})
	r.run("delete_monitoring_location", nil, func() (any, error) {
		return nil, c.DeleteMonitoringLocation(ctx, loc)
	})
}

// writeFindBusinessService looks up a business service id by name,
// resolving list entries that are location URLs (mirrors the Python
// best-effort lookup; any error aborts the search).
func writeFindBusinessService(ctx context.Context, c *opennms.Client, name string) int {
	id := 0
	result, err := c.GetBusinessServices(ctx)
	if err != nil || result == nil {
		return 0
	}
	items, _ := result["business-services"].([]any)
	for _, b := range items {
		var detail map[string]any
		if s, isStr := b.(string); isStr {
			n, convErr := strconv.Atoi(s[strings.LastIndex(s, "/")+1:])
			if convErr != nil {
				return id
			}
			d, getErr := c.GetBusinessService(ctx, n)
			if getErr != nil {
				return id
			}
			detail = d
		} else {
			detail, _ = b.(map[string]any)
		}
		if detail != nil && detail["name"] == name {
			id = writeIDInt(detail["id"])
		}
	}
	return id
}

func writeServiceEntities(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: BSM, situations, applications, links")

	bs := fmt.Sprintf("%s-bs", tag)
	r.run("create_business_service", nil, func() (any, error) {
		return c.CreateBusinessService(ctx, map[string]any{
			"name": bs, "attributes": map[string]any{"attribute": []any{}},
			"reduce-function": map[string]any{"type": "HighestSeverity"},
		})
	})
	_, _ = c.GetBusinessServices(ctx)
	bid := writeFindBusinessService(ctx, c, bs)
	if bid != 0 {
		r.run("update_business_service", nil, func() (any, error) {
			return nil, c.UpdateBusinessService(ctx, bid, map[string]any{
				"name":            bs,
				"reduce-function": map[string]any{"type": "HighestSeverity"},
			})
		})
		r.run("add_reduction_key_edge", nil, func() (any, error) {
			return c.AddReductionKeyEdge(ctx, bid, map[string]any{
				"reduction-key": fmt.Sprintf("%s-rk", tag),
				"map-function":  map[string]any{"type": "Identity"}, "weight": 1,
			})
		})
		r.warn("add_ip_service_edge", "requires monitored service with id=1",
			nil, func() (any, error) {
				return c.AddIpServiceEdge(ctx, bid, map[string]any{
					"ip-service-id": 1,
					"map-function":  map[string]any{"type": "Identity"}, "weight": 1,
				})
			})
		bs2 := fmt.Sprintf("%s-bs2", tag)
		r.run("create_business_service (child)", nil, func() (any, error) {
			return c.CreateBusinessService(ctx, map[string]any{
				"name":            bs2,
				"reduce-function": map[string]any{"type": "HighestSeverity"},
			})
		})
		cid := writeFindBusinessService(ctx, c, bs2)
		if cid != 0 {
			r.run("add_child_edge", nil, func() (any, error) {
				return c.AddChildEdge(ctx, bid, map[string]any{
					"child-id":     cid,
					"map-function": map[string]any{"type": "Identity"}, "weight": 1,
				})
			})
		}
		edgeRemoved := false
		if detail, err := c.GetBusinessService(ctx, bid); err == nil && detail != nil {
			edges, _ := detail["reduction-key-edges"].([]any)
			for _, edge := range edges {
				var eid any
				switch e := edge.(type) {
				case float64:
					eid = e
				case map[string]any:
					eid = e["id"]
				}
				if eid != nil {
					edgeID := writeIDInt(eid)
					r.run("remove_business_service_edge", nil, func() (any, error) {
						return nil, c.RemoveBusinessServiceEdge(ctx, bid, edgeID)
					})
					edgeRemoved = true
					break
				}
			}
		}
		if !edgeRemoved {
			r.skip("remove_business_service_edge", "no edge id found")
		}
		r.run("reload_business_service_daemon", nil, func() (any, error) {
			return nil, c.ReloadBusinessServiceDaemon(ctx)
		})
		if cid != 0 {
			r.run("delete_business_service (child)", nil, func() (any, error) {
				return nil, c.DeleteBusinessService(ctx, cid)
			})
		}
		r.run("delete_business_service", nil, func() (any, error) {
			return nil, c.DeleteBusinessService(ctx, bid)
		})
	} else {
		r.skip("business service mutations", "create returned no id")
	}

	// Situations need existing alarms to group
	_, aid := first(func(limit int) (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "alarm", "")
	if alarmID := writeIDInt(aid); alarmID != 0 {
		r.warn("create_situation", "requires situation support for the alarm set",
			nil, func() (any, error) {
				return c.CreateSituation(ctx, []int{alarmID},
					fmt.Sprintf("%s-situation", tag), "")
			})
		_, sid := first(func(limit int) (any, error) {
			return c.GetSituations(ctx, &opennms.ListOptions{Limit: limit})
		}, "alarm", "")
		if sitID := writeIDInt(sid); sitID != 0 {
			r.warn(fmt.Sprintf("accept_situation  id=%d", sitID),
				"situation lifecycle", nil, func() (any, error) {
					return nil, c.AcceptSituation(ctx, sitID)
				})
			r.warn(fmt.Sprintf("clear_situation  id=%d", sitID),
				"situation lifecycle", nil, func() (any, error) {
					return nil, c.ClearSituation(ctx, sitID)
				})
		}
	} else {
		r.skip("create_situation", "no alarms present")
	}
	_, aid2 := first(func(limit int) (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "alarm", "")
	_, sid2 := first(func(limit int) (any, error) {
		return c.GetSituations(ctx, &opennms.ListOptions{Limit: limit})
	}, "alarm", "")
	alarmID2, sitID2 := writeIDInt(aid2), writeIDInt(sid2)
	if sitID2 != 0 && alarmID2 != 0 {
		r.warn("add_alarms_to_situation", "situation lifecycle",
			nil, func() (any, error) {
				return c.AddAlarmsToSituation(ctx, sitID2, []int{alarmID2}, "")
			})
		r.warn("remove_alarms_from_situation", "situation lifecycle",
			nil, func() (any, error) {
				return nil, c.RemoveAlarmsFromSituation(ctx, sitID2, []int{alarmID2})
			})
		r.warn("clear_situation_alarms", "situation lifecycle",
			nil, func() (any, error) {
				return nil, c.ClearSituationAlarms(ctx, sitID2, nil)
			})
	} else {
		r.skip("situation alarm mutations", "no situation present")
	}
	if alarmID2 != 0 {
		r.run("create_ack (alarm)", nil, func() (any, error) {
			return c.CreateAck(ctx, "ack", alarmID2, 0)
		})
		r.run("create_ack (unack)", nil, func() (any, error) {
			return c.CreateAck(ctx, "unack", alarmID2, 0)
		})
	} else {
		r.skip("create_ack", "no alarms")
	}
	// The Python original probes with a reduction-key string where the
	// Go client takes the situation's alarm id; 0 preserves the
	// always-warn probe without a live situation.
	r.warn("submit_situation_feedback", "requires situation-feedback feature",
		nil, func() (any, error) {
			return nil, c.SubmitSituationFeedback(ctx, 0, []any{})
		})

	app := fmt.Sprintf("%s-app", tag)
	r.run("create_application", nil, func() (any, error) {
		return c.CreateApplication(ctx, map[string]any{"name": app})
	})
	appID := 0
	if apps, err := c.GetApplications(ctx, &opennms.ListOptions{}); err == nil && apps != nil {
		items, _ := apps["application"].([]any)
		for _, item := range items {
			a, _ := item.(map[string]any)
			if a != nil && a["name"] == app {
				appID = writeIDInt(a["id"])
			}
		}
	}
	if appID != 0 {
		r.run("delete_application", nil, func() (any, error) {
			return nil, c.DeleteApplication(ctx, appID)
		})
	} else {
		r.skip("delete_application", "application id not found")
	}

	r.warn("create_user_defined_link", "requires two existing nodes",
		nil, func() (any, error) {
			return c.CreateUserDefinedLink(ctx, map[string]any{
				"node-id-a": 1, "node-id-z": 1, "link-id": tag,
				"owner": "smoke",
			})
		})
	if links, err := c.GetUserDefinedLinks(ctx); err == nil && links != nil {
		items, _ := links["user-defined-link"].([]any)
		for _, item := range items {
			link, _ := item.(map[string]any)
			if link != nil && link["link-id"] == tag {
				linkID := writeIDInt(link["db-id"])
				r.run("delete_user_defined_link", nil, func() (any, error) {
					return nil, c.DeleteUserDefinedLink(ctx, linkID)
				})
			}
		}
	}
}

func writeClassificationScv(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: flow classification and credentials vault")

	grpID, ruleID := 0, 0
	r.run("create_classification_group", nil, func() (any, error) {
		return nil, c.CreateClassificationGroup(ctx, map[string]any{
			"name": fmt.Sprintf("%s-cgrp", tag), "enabled": true,
			"description": "smoke",
		})
	})
	if groups, err := c.GetClassificationGroups(ctx, nil); err == nil && groups != nil {
		items, _ := groups["classificationGroup"].([]any)
		for _, item := range items {
			g, _ := item.(map[string]any)
			if g != nil && g["name"] == fmt.Sprintf("%s-cgrp", tag) {
				grpID = writeIDInt(g["id"])
			}
		}
	}
	if grpID != 0 {
		r.run("update_classification_group", nil, func() (any, error) {
			return nil, c.UpdateClassificationGroup(ctx, grpID, map[string]any{
				"name": fmt.Sprintf("%s-cgrp", tag), "enabled": false,
				"description": "smoke2",
			})
		})
		r.run("create_classification_rule", nil, func() (any, error) {
			return nil, c.CreateClassificationRule(ctx, map[string]any{
				"name": fmt.Sprintf("%s-rule", tag), "dstPort": "9999",
				"protocol": "tcp", "omnidirectional": false,
				"group": map[string]any{"id": grpID},
			})
		})
		if rules, err := c.GetClassificationRules(ctx, nil); err == nil && rules != nil {
			items, _ := rules["classificationRule"].([]any)
			for _, item := range items {
				rule, _ := item.(map[string]any)
				if rule != nil && rule["name"] == fmt.Sprintf("%s-rule", tag) {
					ruleID = writeIDInt(rule["id"])
				}
			}
		}
		if ruleID != 0 {
			r.run("update_classification_rule", nil, func() (any, error) {
				return nil, c.UpdateClassificationRule(ctx, ruleID, map[string]any{
					"name": fmt.Sprintf("%s-rule", tag), "dstPort": "9998",
					"protocol": "tcp", "omnidirectional": false,
					"group": map[string]any{"id": grpID},
				})
			})
			r.run("delete_classification_rule", nil, func() (any, error) {
				return nil, c.DeleteClassificationRule(ctx, ruleID)
			})
		}
		r.warn("import_classification_rules", "CSV import format varies by version",
			nil, func() (any, error) {
				return nil, c.ImportClassificationRules(ctx, grpID,
					"name;protocol;srcAddress;srcPort;dstAddress;dstPort;"+
						"exporterFilter;omnidirectional\n"+
						fmt.Sprintf("%s-csv;tcp;;;;9997;;false\n", tag))
			})
		r.warn("delete_classification_rules (group)", "bulk delete of the group rules",
			nil, func() (any, error) {
				return nil, c.DeleteClassificationRules(ctx, grpID)
			})
		r.run("delete_classification_group", nil, func() (any, error) {
			return nil, c.DeleteClassificationGroup(ctx, grpID)
		})
	} else {
		r.skip("classification mutations", "group id not found")
	}
	r.run("classify", nil, func() (any, error) {
		return c.Classify(ctx, map[string]any{
			"protocol": "tcp", "dstPort": "443", "srcAddress": "10.0.0.1",
			"srcPort": "55555", "dstAddress": "10.0.0.2",
			"exporterAddress": "10.0.0.3",
		})
	})

	alias := fmt.Sprintf("%s-cred", tag)
	r.run("create_credential", nil, func() (any, error) {
		return nil, c.CreateCredential(ctx, map[string]any{
			"alias": alias, "username": "smoke", "password": "pw",
		})
	})
	r.run("update_credential", nil, func() (any, error) {
		return nil, c.UpdateCredential(ctx, alias, map[string]any{
			"alias": alias, "username": "smoke2", "password": "pw2",
		})
	})
	r.run("delete_credential", nil, func() (any, error) {
		return nil, c.DeleteCredential(ctx, alias)
	})
}

func writeConfigsNbi(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: config management and northbounders")

	cfg, _ := r.warn("get_config (provisiond)", "cm API payload for round-trip update",
		nil, func() (any, error) {
			return c.GetConfig(ctx, "provisiond", "default")
		})
	if m, ok := cfg.(map[string]any); ok && m != nil {
		r.run("update_config (no-op round-trip)", nil, func() (any, error) {
			return nil, c.UpdateConfig(ctx, "provisiond", "default", m)
		})
	} else {
		r.skip("update_config", "no provisiond config payload")
	}
	r.warn("create_config", "most cm schemas are single-instance",
		nil, func() (any, error) {
			return nil, c.CreateConfig(ctx, "provisiond", tag,
				map[string]any{"importThreads": 8})
		})
	r.warn("delete_config", "cleanup of the instance above",
		nil, func() (any, error) {
			return nil, c.DeleteConfig(ctx, "provisiond", tag)
		})
	r.warn("delete_config_part", "requires a part path in the schema",
		nil, func() (any, error) {
			return nil, c.DeleteConfigPart(ctx, "provisiond", "default",
				"nonexistent-part")
		})

	dest := fmt.Sprintf("%s-dest", tag)
	r.run("create_email_nbi_destination", nil, func() (any, error) {
		return nil, c.CreateEmailNbiDestination(ctx, map[string]any{"name": dest})
	})
	r.run("update_email_nbi_destination", nil, func() (any, error) {
		return nil, c.UpdateEmailNbiDestination(ctx, dest,
			map[string]string{"firstOccurrenceOnly": "true"})
	})
	r.run("delete_email_nbi_destination", nil, func() (any, error) {
		return nil, c.DeleteEmailNbiDestination(ctx, dest)
	})
	status, _ := r.warn("get_email_nbi_status", "read for status round-trip",
		nil, func() (any, error) {
			return c.GetEmailNbiStatus(ctx)
		})
	enabled := false
	if m, ok := status.(map[string]any); ok && m != nil {
		enabled, _ = m["enabled"].(bool)
	}
	r.run("set_email_nbi_status (restore)", nil, func() (any, error) {
		return nil, c.SetEmailNbiStatus(ctx, enabled)
	})
	ecfg, _ := r.warn("get_email_nbi_config", "read for config round-trip",
		nil, func() (any, error) {
			return c.GetEmailNbiConfig(ctx)
		})
	if m, ok := ecfg.(map[string]any); ok && m != nil {
		r.run("update_email_nbi_config (no-op)", nil, func() (any, error) {
			return nil, c.UpdateEmailNbiConfig(ctx, m)
		})
	}

	sink := fmt.Sprintf("%s-sink", tag)
	r.run("create_snmptrap_nbi_trapsink", nil, func() (any, error) {
		return nil, c.CreateSnmptrapNbiTrapsink(ctx, map[string]any{
			"name": sink, "ip-address": "127.0.0.1", "port": 1162,
		})
	})
	r.run("update_snmptrap_nbi_trapsink", nil, func() (any, error) {
		return nil, c.UpdateSnmptrapNbiTrapsink(ctx, sink,
			map[string]string{"port": "1163"})
	})
	r.run("delete_snmptrap_nbi_trapsink", nil, func() (any, error) {
		return nil, c.DeleteSnmptrapNbiTrapsink(ctx, sink)
	})
	tstat, _ := r.warn("get_snmptrap_nbi_status", "read for status round-trip",
		nil, func() (any, error) {
			return c.GetSnmptrapNbiStatus(ctx)
		})
	trapEnabled := false
	if m, ok := tstat.(map[string]any); ok && m != nil {
		trapEnabled, _ = m["enabled"].(bool)
	}
	r.run("set_snmptrap_nbi_status (restore)", nil, func() (any, error) {
		return nil, c.SetSnmptrapNbiStatus(ctx, trapEnabled)
	})
	tcfg, _ := r.warn("get_snmptrap_nbi_config", "read for config round-trip",
		nil, func() (any, error) {
			return c.GetSnmptrapNbiConfig(ctx)
		})
	if m, ok := tcfg.(map[string]any); ok && m != nil {
		r.run("update_snmptrap_nbi_config (no-op)", nil, func() (any, error) {
			return nil, c.UpdateSnmptrapNbiConfig(ctx, m)
		})
	}

	sdest := fmt.Sprintf("%s-sdest", tag)
	r.run("create_syslog_nbi_destination", nil, func() (any, error) {
		return nil, c.CreateSyslogNbiDestination(ctx, map[string]any{
			"destination-name": sdest, "host": "127.0.0.1", "port": 1514,
		})
	})
	r.run("update_syslog_nbi_destination", nil, func() (any, error) {
		return nil, c.UpdateSyslogNbiDestination(ctx, sdest,
			map[string]string{"port": "1515"})
	})
	r.run("delete_syslog_nbi_destination", nil, func() (any, error) {
		return nil, c.DeleteSyslogNbiDestination(ctx, sdest)
	})
	sstat, _ := r.warn("get_syslog_nbi_status", "read for status round-trip",
		nil, func() (any, error) {
			return c.GetSyslogNbiStatus(ctx)
		})
	syslogEnabled := false
	if m, ok := sstat.(map[string]any); ok && m != nil {
		syslogEnabled, _ = m["enabled"].(bool)
	}
	r.run("set_syslog_nbi_status (restore)", nil, func() (any, error) {
		return nil, c.SetSyslogNbiStatus(ctx, syslogEnabled)
	})
	scfg, _ := r.warn("get_syslog_nbi_config", "read for config round-trip",
		nil, func() (any, error) {
			return c.GetSyslogNbiConfig(ctx)
		})
	if m, ok := scfg.(map[string]any); ok && m != nil {
		r.run("update_syslog_nbi_config (no-op)", nil, func() (any, error) {
			return nil, c.UpdateSyslogNbiConfig(ctx, m)
		})
	}

	r.warn("create_javamail_readmail", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.CreateJavamailReadmail(ctx,
				map[string]any{"name": fmt.Sprintf("%s-rm", tag)})
		})
	r.warn("update_javamail_readmail", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.UpdateJavamailReadmail(ctx, fmt.Sprintf("%s-rm", tag),
				map[string]string{"host": "127.0.0.1"})
		})
	r.warn("delete_javamail_readmail", "cleanup",
		nil, func() (any, error) {
			return nil, c.DeleteJavamailReadmail(ctx, fmt.Sprintf("%s-rm", tag))
		})
	r.warn("create_javamail_sendmail", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.CreateJavamailSendmail(ctx,
				map[string]any{"name": fmt.Sprintf("%s-sm", tag)})
		})
	r.warn("update_javamail_sendmail", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.UpdateJavamailSendmail(ctx, fmt.Sprintf("%s-sm", tag),
				map[string]string{"host": "127.0.0.1"})
		})
	r.warn("delete_javamail_sendmail", "cleanup",
		nil, func() (any, error) {
			return nil, c.DeleteJavamailSendmail(ctx, fmt.Sprintf("%s-sm", tag))
		})
	r.warn("create_javamail_end2end", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.CreateJavamailEnd2end(ctx,
				map[string]any{"name": fmt.Sprintf("%s-e2e", tag)})
		})
	r.warn("update_javamail_end2end", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.UpdateJavamailEnd2end(ctx, fmt.Sprintf("%s-e2e", tag),
				map[string]string{"readMailConfigName": fmt.Sprintf("%s-rm", tag)})
		})
	r.warn("delete_javamail_end2end", "cleanup",
		nil, func() (any, error) {
			return nil, c.DeleteJavamailEnd2end(ctx, fmt.Sprintf("%s-e2e", tag))
		})
	r.warn("set_javamail_default_config", "javamail config API absent on some versions",
		nil, func() (any, error) {
			return nil, c.SetJavamailDefaultConfig(ctx, map[string]any{
				"defaultReadConfigName": "default",
				"defaultSendConfigName": "default",
			})
		})
}

func writeReportsGraphmlMisc(ctx context.Context, c *opennms.Client, r *runner) {
	tag := writeTag()
	r.section("write: reports, graphml, grafana, settings")

	template := ""
	if templates, err := c.GetReportTemplates(ctx); err == nil && len(templates) > 0 {
		if t, ok := templates[0].(map[string]any); ok && t != nil {
			template = writeIDString(t["id"])
		}
	}
	if template != "" {
		trigger := fmt.Sprintf("%s-trigger", tag)
		r.warn("schedule_report", "requires schedulable template parameters",
			nil, func() (any, error) {
				return nil, c.ScheduleReport(ctx, template, "PDF",
					"0 0 6 1 1 ? 2099", []any{},
					map[string]any{"instanceId": trigger, "persist": true})
			})
		r.warn("update_scheduled_report", "requires the trigger created above",
			nil, func() (any, error) {
				return nil, c.UpdateScheduledReport(ctx, trigger,
					map[string]any{"cronExpression": "0 0 7 1 1 ? 2099"})
			})
		r.warn("delete_scheduled_report", "cleanup",
			nil, func() (any, error) {
				return nil, c.DeleteScheduledReport(ctx, trigger)
			})
		r.run("delete_scheduled_reports", nil, func() (any, error) {
			return nil, c.DeleteScheduledReports(ctx)
		})
		r.warn("run_report", "server-side render can be heavy or need params",
			nil, func() (any, error) {
				return c.RunReport(ctx, template, "PDF", []any{})
			})
		r.warn("deliver_report", "delivery needs renderable template",
			nil, func() (any, error) {
				return nil, c.DeliverReport(ctx, template, "PDF", []any{},
					map[string]any{
						"instanceId": fmt.Sprintf("%s-deliver", tag),
						"persist":    true,
					})
			})
		r.run("delete_persisted_reports", nil, func() (any, error) {
			return nil, c.DeletePersistedReports(ctx)
		})
	} else {
		r.skip("report write suite", "no report templates")
	}

	gname := fmt.Sprintf("%s-graph", tag)
	r.run("create_graphml", nil, func() (any, error) {
		return nil, c.CreateGraphml(ctx, gname,
			`<?xml version="1.0" encoding="UTF-8"?>`+
				`<graphml xmlns="http://graphml.graphdrawing.org/xmlns">`+
				`<key id="label" for="all" attr.name="label"`+
				` attr.type="string"/>`+
				fmt.Sprintf(`<graph id="%s"><data key="label">smoke</data>`, gname)+
				`<node id="n0"><data key="label">n0</data></node>`+
				`</graph></graphml>`)
	})
	r.run("delete_graphml", nil, func() (any, error) {
		return nil, c.DeleteGraphml(ctx, gname)
	})

	guid := fmt.Sprintf("%s-grafana", tag)
	r.run("create_grafana_endpoint", nil, func() (any, error) {
		return nil, c.CreateGrafanaEndpoint(ctx, map[string]any{
			"uid": guid, "url": "http://127.0.0.1:3000",
			"apiKey": "smoke",
		})
	})
	gid := 0
	if endpoints, err := c.GetGrafanaEndpoints(ctx); err == nil {
		for _, item := range endpoints {
			e, _ := item.(map[string]any)
			if e != nil && e["uid"] == guid {
				gid = writeIDInt(e["id"])
			}
		}
	}
	if gid != 0 {
		r.run("update_grafana_endpoint", nil, func() (any, error) {
			return nil, c.UpdateGrafanaEndpoint(ctx, gid, map[string]any{
				"id": gid, "uid": guid, "url": "http://127.0.0.1:3000",
				"apiKey": "smoke2",
			})
		})
		r.warn("verify_grafana_endpoint", "verification calls the Grafana URL",
			nil, func() (any, error) {
				return nil, c.VerifyGrafanaEndpoint(ctx, map[string]any{
					"url": "http://127.0.0.1:3000", "apiKey": "smoke",
				})
			})
		r.run("delete_grafana_endpoint", nil, func() (any, error) {
			return nil, c.DeleteGrafanaEndpoint(ctx, gid)
		})
	} else {
		r.skip("grafana endpoint mutations", "endpoint id not found")
	}
	r.run("delete_grafana_endpoints", nil, func() (any, error) {
		return nil, c.DeleteGrafanaEndpoints(ctx)
	})

	geo, _ := r.warn("get_geocoding_config", "read for round-trip restore",
		nil, func() (any, error) {
			return c.GetGeocodingConfig(ctx)
		})
	r.run("reset_geocoding_config", nil, func() (any, error) {
		return nil, c.ResetGeocodingConfig(ctx)
	})
	if m, ok := geo.(map[string]any); ok && m != nil {
		if geocoderID, _ := m["activeGeocoderId"].(string); geocoderID != "" {
			r.run("set_active_geocoder (restore)", nil, func() (any, error) {
				return nil, c.SetActiveGeocoder(ctx, geocoderID)
			})
		}
	}
	r.warn("configure_geocoder", "geocoder config keys vary by provider",
		nil, func() (any, error) {
			return nil, c.ConfigureGeocoder(ctx, "nominatim",
				map[string]string{"userAgent": "python-opennms-smoke"})
		})

	stats, _ := r.warn("get_usage_statistics_status", "read for round-trip restore",
		nil, func() (any, error) {
			return c.GetUsageStatisticsStatus(ctx)
		})
	if m, ok := stats.(map[string]any); ok && m != nil {
		r.run("set_usage_statistics_status (restore)", nil, func() (any, error) {
			return nil, c.SetUsageStatisticsStatus(ctx,
				writeBoolPtr(m["enabled"]),
				writeBoolPtr(m["initialNoticeAcknowledged"]))
		})
	}
	pstat, _ := r.warn("get_product_update_status", "read for round-trip restore",
		nil, func() (any, error) {
			return c.GetProductUpdateStatus(ctx)
		})
	if m, ok := pstat.(map[string]any); ok && m != nil {
		r.run("set_product_update_status (restore)", nil, func() (any, error) {
			return nil, c.SetProductUpdateStatus(ctx,
				writeBoolPtr(m["optedIn"]),
				writeBoolPtr(m["noticeAcknowledged"]))
		})
	}
	r.warn("submit_product_update_enrollment",
		"returns 500 when enrollment is disabled (documented)",
		nil, func() (any, error) {
			return nil, c.SubmitProductUpdateEnrollment(ctx, map[string]any{
				"consent": false, "email": "smoke@example.invalid",
			})
		})

	r.run("set_snmp_config", nil, func() (any, error) {
		return nil, c.SetSnmpConfig(ctx, "10.254.99.1", map[string]any{
			"readCommunity": "smoke", "version": "v2c",
		})
	})
	r.warn("delete_persisted_report", "requires an existing persisted report id",
		nil, func() (any, error) {
			return nil, c.DeletePersistedReport(ctx, 999999)
		})
	r.warn("update_eventconf_event", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.UpdateEventconfEvent(ctx, "smoke", "1", map[string]any{
				"uei": "uei.opennms.org/smoke", "event-label": "smoke",
				"descr": "smoke", "logmsg": map[string]any{"content": "smoke"},
				"severity": "Normal",
			})
		})
	r.run("query_geolocations", nil, func() (any, error) {
		return c.QueryGeolocations(ctx, "", "", nil)
	})
	r.warn("get_measurements_multi", "requires collected time-series data",
		nil, func() (any, error) {
			return c.GetMeasurementsMulti(ctx, map[string]any{
				"start": 0, "end": 1,
				"source": []map[string]any{{
					"aggregation": "AVERAGE", "attribute": "loadavg1",
					"label": "l", "resourceId": "node[1].nodeSnmp[]",
				}},
			})
		})
	r.warn("get_graph_view", "requires the topology container",
		nil, func() (any, error) {
			return c.GetGraphView(ctx, "bsm", "bsm", 1, nil)
		})

	_, eid := first(func(limit int) (any, error) {
		return c.GetEvents(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "event", "")
	if eventID := writeIDInt(eid); eventID != 0 {
		r.run(fmt.Sprintf("ack_event  id=%d", eventID), nil, func() (any, error) {
			return nil, c.AckEvent(ctx, eventID)
		})
		r.run(fmt.Sprintf("unack_event  id=%d", eventID), nil, func() (any, error) {
			return nil, c.UnackEvent(ctx, eventID)
		})
	} else {
		r.skip("ack_event / unack_event", "no events")
	}
	r.run("bulk_ack_events", nil, func() (any, error) {
		return nil, c.BulkAckEvents(ctx, map[string]string{"limit": "1"})
	})
	r.run("bulk_unack_events", nil, func() (any, error) {
		return nil, c.BulkUnackEvents(ctx, map[string]string{"limit": "1"})
	})
	r.run("bulk_ack_alarms", nil, func() (any, error) {
		return nil, c.BulkAckAlarms(ctx, map[string]string{"limit": "1"})
	})
	r.run("bulk_unack_alarms", nil, func() (any, error) {
		return nil, c.BulkUnackAlarms(ctx, map[string]string{"limit": "1"})
	})
	r.warn("bulk_clear_alarms", "clears matching alarms permanently",
		nil, func() (any, error) {
			return nil, c.BulkClearAlarms(ctx, map[string]string{"limit": "1"})
		})
	r.warn("bulk_escalate_alarms", "escalates matching alarms",
		nil, func() (any, error) {
			return nil, c.BulkEscalateAlarms(ctx, map[string]string{"limit": "1"})
		})
	_, aid := first(func(limit int) (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "alarm", "")
	if alarmID := writeIDInt(aid); alarmID != 0 {
		r.warn(fmt.Sprintf("escalate_alarm  id=%d", alarmID),
			"raises severity permanently", nil, func() (any, error) {
				return nil, c.EscalateAlarm(ctx, alarmID)
			})
		r.warn(fmt.Sprintf("clear_alarm  id=%d", alarmID),
			"clears the alarm permanently", nil, func() (any, error) {
				return nil, c.ClearAlarm(ctx, alarmID)
			})
	} else {
		r.skip("clear_alarm / escalate_alarm", "no alarms")
	}

	r.warn("discover (127.0.0.1 one-shot)", "submits a discovery scan job",
		nil, func() (any, error) {
			return nil, c.Discover(ctx, map[string]any{
				"specifics": []map[string]any{{
					"ip": "127.0.0.1", "location": "Default",
					"retries": 1, "timeout": 2000,
				}},
			})
		})
	r.warn("update_ifservices", "parameters depend on existing services",
		nil, func() (any, error) {
			return nil, c.UpdateIfservices(ctx,
				map[string]string{"services": "ICMP", "status": "R"})
		})
	r.warn("backup_device_config", "requires DeviceConfig-enabled service",
		nil, func() (any, error) {
			return c.BackupDeviceConfig(ctx, []map[string]any{
				{"ipAddress": "1", "location": ""},
			})
		})
	r.warn("delete_resource", "requires the resource to exist",
		nil, func() (any, error) {
			return nil, c.DeleteResource(ctx, "node[999999].nodeSnmp[]")
		})
	r.warn("trigger_destination_path", "requires a configured destination path",
		nil, func() (any, error) {
			return nil, c.TriggerDestinationPath(ctx, "smoke-nonexistent-path")
		})
	r.warn("upload_filesystem_contents", "requires FILESYSTEM EDITOR role",
		nil, func() (any, error) {
			return nil, c.UploadFilesystemContents(ctx, "smoke-test.xml",
				[]byte("<x/>"))
		})
	r.warn("delete_filesystem_file", "requires FILESYSTEM EDITOR role",
		nil, func() (any, error) {
			return nil, c.DeleteFilesystemFile(ctx, "smoke-test.xml")
		})
	r.warn("upload_eventconf", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.UploadEventconf(ctx,
				[]byte("<events xmlns='http://xmlns.opennms.org/xsd/eventconf'/>"))
		})
	r.warn("create_eventconf_event", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return c.CreateEventconfEvent(ctx, "smoke", map[string]any{
				"uei": fmt.Sprintf("uei.opennms.org/%s", tag), "event-label": tag,
				"descr": "smoke", "logmsg": map[string]any{"content": "smoke"},
				"severity": "Normal",
			})
		})
	r.warn("set_eventconf_sources_status", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.SetEventconfSourcesStatus(ctx,
				map[string]any{"sources": []any{}})
		})
	r.warn("set_eventconf_events_status", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.SetEventconfEventsStatus(ctx, "smoke",
				map[string]any{"events": []any{}})
		})
	r.warn("delete_eventconf_events", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.DeleteEventconfEvents(ctx, "smoke", nil)
		})
	r.warn("delete_eventconf_sources", "eventconf v2 API requires Horizon 35+",
		nil, func() (any, error) {
			return nil, c.DeleteEventconfSources(ctx,
				map[string]any{"sources": []any{"smoke"}})
		})

	// KSC reports have no DELETE endpoint; use a timestamp-unique ID
	// and accept the leftover on throwaway instances only.
	kscID := int(time.Now().Unix())%100000 + 10000
	r.warn(fmt.Sprintf("create_ksc_report  id=%d", kscID),
		"KSC API has no DELETE; leaves a report behind",
		nil, func() (any, error) {
			return nil, c.CreateKscReport(ctx, map[string]any{
				"id": kscID, "label": fmt.Sprintf("%s-ksc", tag),
			})
		})
	r.warn(fmt.Sprintf("add_graph_to_ksc_report  id=%d", kscID),
		"requires the report above and a valid resource",
		nil, func() (any, error) {
			return nil, c.AddGraphToKscReport(ctx, kscID, "mib2.bits",
				"node[1].nodeSnmp[]", tag, "")
		})
}
