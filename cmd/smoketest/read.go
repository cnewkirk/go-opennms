package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"

	opennms "github.com/cnewkirk/go-opennms"
)

// readStr renders any value the way Python's str() renders it for
// the detail column.
func readStr(v any) string {
	return fmt.Sprintf("%v", v)
}

// readID renders a decoded JSON id (float64) without float
// formatting artifacts; other types fall back to %v.
func readID(v any) string {
	if f, ok := v.(float64); ok && f == math.Trunc(f) {
		return strconv.FormatInt(int64(f), 10)
	}
	return fmt.Sprintf("%v", v)
}

// readTruthy mirrors Python truthiness for values pulled out of
// decoded JSON.
func readTruthy(v any) bool {
	switch x := v.(type) {
	case nil:
		return false
	case float64:
		return x != 0
	case string:
		return x != ""
	case bool:
		return x
	}
	return true
}

// readItems extracts a list of items from a decoded response: a
// wrapped map (trying keys in order) or a bare list.
func readItems(result any, keys ...string) []any {
	switch v := result.(type) {
	case map[string]any:
		for _, k := range keys {
			if items, ok := v[k].([]any); ok {
				return items
			}
		}
	case []any:
		return v
	}
	return nil
}

// readField returns item[key] when item is an object, else nil.
func readField(item any, key string) any {
	if m, ok := item.(map[string]any); ok {
		return m[key]
	}
	return nil
}

// readStrField returns item[key] as a string, or "" when absent or
// not a string.
func readStrField(item any, key string) string {
	if m, ok := item.(map[string]any); ok {
		if s, ok := m[key].(string); ok {
			return s
		}
	}
	return ""
}

// readFieldStr renders the first present key of an object result as
// a string; "" for non-objects or when no key is present.
func readFieldStr(r any, keys ...string) string {
	m, ok := r.(map[string]any)
	if !ok {
		return ""
	}
	for _, k := range keys {
		if v, found := m[k]; found {
			return fmt.Sprintf("%v", v)
		}
	}
	return ""
}

// readFirstNode fetches the first node and its id.
func readFirstNode(ctx context.Context, c *opennms.Client) (map[string]any, any) {
	return first(func(limit int) (any, error) {
		return c.GetNodes(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "node", "")
}

// runReadChecks exercises every read-only endpoint group, in the
// same order as smoke_test.py's main().
func runReadChecks(ctx context.Context, c *opennms.Client, r *runner) {
	readInfo(ctx, c, r)
	readAlarms(ctx, c, r)
	readEvents(ctx, c, r)
	readAcks(ctx, c, r)
	readNotifications(ctx, c, r)
	readNodes(ctx, c, r)
	readOutages(ctx, c, r)
	readRequisitions(ctx, c, r)
	readForeignSources(ctx, c, r)
	readSnmpConfig(ctx, c, r)
	readGroups(ctx, c, r)
	readUsers(ctx, c, r)
	readCategories(ctx, c, r)
	readSchedOutages(ctx, c, r)
	readKscReports(ctx, c, r)
	readResources(ctx, c, r)
	readMeasurements(ctx, c, r)
	readHeatmap(ctx, c, r)
	readMaps(ctx, c, r)
	readGraphs(ctx, c, r)
	readFlows(ctx, c, r)
	readDeviceConfig(ctx, c, r)
	readSituations(ctx, c, r)
	readBusinessServices(ctx, c, r)
	readEnlinkd(ctx, c, r)
	readV2Interfaces(ctx, c, r)
	readMonitoringLocations(ctx, c, r)
	readMinions(ctx, c, r)
	readIfservices(ctx, c, r)
	readAvailability(ctx, c, r)
	readHealth(ctx, c, r)
	readWhoami(ctx, c, r)
	readMonitoringSystems(ctx, c, r)
	readPrefabGraphs(ctx, c, r)
	readFlowDscp(ctx, c, r)
	readBusinessServiceFunctions(ctx, c, r)
	readClassifications(ctx, c, r)
	readSituationFeedback(ctx, c, r)
	readUserDefinedLinks(ctx, c, r)
	readApplications(ctx, c, r)
	readPerspectivePoller(ctx, c, r)
	readForeignSourcesConfig(ctx, c, r)
	readRequisitionNames(ctx, c, r)
	readSnmpMetadata(ctx, c, r)
	readProvisiond(ctx, c, r)
	readEventconf(ctx, c, r)
	readAssetSuggestions(ctx, c, r)
	readScv(ctx, c, r)
	readConfigMgmt(ctx, c, r)
	readSnmptrapNbi(ctx, c, r)
	readEmailNbi(ctx, c, r)
	readSyslogNbi(ctx, c, r)
	readJavamailConfig(ctx, c, r)
	readPagination(ctx, c, r)
	readExceptions(ctx, c, r)
}

func readInfo(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("info")
	r.run("get_info", func(res any) string {
		return readFieldStr(res, "displayVersion", "version")
	}, func() (any, error) {
		return c.GetInfo(ctx)
	})
}

func readAlarms(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("alarms")
	r.run("get_alarms", n("alarm"), func() (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: 5}, nil)
	})
	r.warn("get_alarm_count", "", readStr, func() (any, error) {
		return c.GetAlarmCount(ctx)
	})
	r.run("get_alarm_stats", nil, func() (any, error) {
		return c.GetAlarmStats(ctx, nil)
	})
	r.run("get_alarm_stats_by_severity", nil, func() (any, error) {
		return c.GetAlarmStatsBySeverity(ctx, nil)
	})
	r.warn("get_alarm_history",
		"requires opennms-alarm-history-elastic Karaf feature",
		nil, func() (any, error) {
			return c.GetAlarmHistory(ctx, 0)
		})
	r.run("get_alarms_v2", n("alarm"), func() (any, error) {
		return c.GetAlarmsV2(ctx, "", &opennms.ListOptions{Limit: 5})
	})

	_, aid := first(func(limit int) (any, error) {
		return c.GetAlarms(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "alarm", "")
	if readTruthy(aid) {
		r.run(fmt.Sprintf("get_alarm                id=%s", readID(aid)),
			func(res any) string {
				return readFieldStr(res, "severity")
			}, func() (any, error) {
				return c.GetAlarm(ctx, toInt(aid))
			})
		r.run(fmt.Sprintf("get_alarm_v2             id=%s", readID(aid)),
			nil, func() (any, error) {
				return c.GetAlarmV2(ctx, toInt(aid))
			})
		r.warn(fmt.Sprintf("get_alarm_history_at     id=%s", readID(aid)),
			"requires opennms-alarm-history-elastic Karaf feature",
			nil, func() (any, error) {
				return c.GetAlarmHistoryAt(ctx, toInt(aid), 0)
			})
		r.warn(fmt.Sprintf("get_alarm_history_states id=%s", readID(aid)),
			"requires opennms-alarm-history-elastic Karaf feature",
			nil, func() (any, error) {
				return c.GetAlarmHistoryStates(ctx, toInt(aid))
			})
	} else {
		for _, lbl := range []string{"get_alarm", "get_alarm_v2",
			"get_alarm_history_at", "get_alarm_history_states"} {
			r.skip(lbl, "no alarms")
		}
	}
}

func readEvents(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("events")
	// Unfiltered event queries can time out on large systems (33M+
	// rows). Filter by the lowest-ID node (typically the
	// self-monitor) to hit an indexed column and avoid full table
	// scans.
	_, nid := first(func(limit int) (any, error) {
		return c.GetNodes(ctx, &opennms.ListOptions{
			Limit: limit, OrderBy: "id", Order: "asc",
		}, nil)
	}, "node", "")
	var nodeFilter map[string]string
	if readTruthy(nid) {
		nodeFilter = map[string]string{"node.id": readID(nid)}
	}
	result, _ := r.run("get_events", n("event"), func() (any, error) {
		return c.GetEvents(ctx, &opennms.ListOptions{Limit: 5}, nodeFilter)
	})
	r.warn("get_event_count", "", readStr, func() (any, error) {
		return c.GetEventCount(ctx)
	})

	var eid any
	if events := readItems(result, "event"); len(events) > 0 {
		eid = readField(events[0], "id")
	}
	if readTruthy(eid) {
		r.run(fmt.Sprintf("get_event  id=%s", readID(eid)), nil,
			func() (any, error) {
				return c.GetEvent(ctx, toInt(eid))
			})
	} else {
		r.skip("get_event", "no events")
	}
}

func readAcks(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("acknowledgements")
	r.run("get_acks", n("ack"), func() (any, error) {
		return c.GetAcks(ctx, &opennms.ListOptions{Limit: 5}, nil)
	})
	r.warn("get_ack_count", "", readStr, func() (any, error) {
		return c.GetAckCount(ctx)
	})

	_, ackID := first(func(limit int) (any, error) {
		return c.GetAcks(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "ack", "")
	if readTruthy(ackID) {
		r.run(fmt.Sprintf("get_ack  id=%s", readID(ackID)), nil,
			func() (any, error) {
				return c.GetAck(ctx, toInt(ackID))
			})
	} else {
		r.skip("get_ack", "no acks")
	}
}

func readNotifications(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("notifications")
	r.run("get_notifications", n("notification"), func() (any, error) {
		return c.GetNotifications(ctx, &opennms.ListOptions{Limit: 5}, nil)
	})
	r.warn("get_notification_count", "", readStr, func() (any, error) {
		return c.GetNotificationCount(ctx)
	})

	_, nid := first(func(limit int) (any, error) {
		return c.GetNotifications(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "notification", "")
	if readTruthy(nid) {
		r.run(fmt.Sprintf("get_notification  id=%s", readID(nid)), nil,
			func() (any, error) {
				return c.GetNotification(ctx, toInt(nid))
			})
	} else {
		r.skip("get_notification", "no notifications")
	}
}

func readNodes(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("nodes")
	r.run("get_nodes", n("node"), func() (any, error) {
		return c.GetNodes(ctx, &opennms.ListOptions{Limit: 5}, nil)
	})
	r.run("get_node_count", readStr, func() (any, error) {
		return c.GetNodeCount(ctx)
	})

	_, nid := readFirstNode(ctx, c)
	if readTruthy(nid) {
		nidStr := readID(nid)
		r.run(fmt.Sprintf("get_node                    id=%s", nidStr),
			func(res any) string {
				return readFieldStr(res, "label")
			}, func() (any, error) {
				return c.GetNode(ctx, nidStr)
			})
		ifacesR, _ := r.run(
			fmt.Sprintf("get_node_ip_interfaces      id=%s", nidStr),
			n("ipInterface"), func() (any, error) {
				return c.GetNodeIpInterfaces(ctx, nidStr,
					&opennms.ListOptions{Limit: 3}, nil)
			})
		var ip, svc string
		if ifaces := readItems(ifacesR, "ipInterface"); len(ifaces) > 0 {
			ip = readStrField(ifaces[0], "ipAddress")
			if ip != "" {
				r.run(fmt.Sprintf("get_node_ip_interface       ip=%s", ip),
					nil, func() (any, error) {
						return c.GetNodeIpInterface(ctx, nidStr, ip)
					})
				svcR, _ := r.run(
					fmt.Sprintf("get_node_ip_services        ip=%s", ip),
					nil, func() (any, error) {
						return c.GetNodeIpServices(ctx, nidStr, ip)
					})
				if svcs := readItems(svcR, "service"); len(svcs) > 0 {
					svc = readStrField(
						readField(svcs[0], "serviceType"), "name")
					if svc != "" {
						r.run(fmt.Sprintf("get_node_ip_service     svc=%s",
							svc), nil, func() (any, error) {
							return c.GetNodeIpService(ctx, nidStr, ip, svc)
						})
					}
				}
			}
		}
		snmpR, _ := r.run(
			fmt.Sprintf("get_node_snmp_interfaces    id=%s", nidStr),
			n("snmpInterface"), func() (any, error) {
				return c.GetNodeSnmpInterfaces(ctx, nidStr,
					&opennms.ListOptions{Limit: 3}, nil)
			})
		if snmps := readItems(snmpR, "snmpInterface"); len(snmps) > 0 {
			ifIndex := readField(snmps[0], "ifIndex")
			if ifIndex != nil {
				r.run(fmt.Sprintf("get_node_snmp_interface  ifIndex=%s",
					readID(ifIndex)), nil, func() (any, error) {
					return c.GetNodeSnmpInterface(ctx, nidStr,
						toInt(ifIndex))
				})
			}
		}
		catsR, _ := r.run(
			fmt.Sprintf("get_node_categories         id=%s", nidStr),
			nil, func() (any, error) {
				return c.GetNodeCategories(ctx, nidStr)
			})
		if cats := readItems(catsR, "category"); len(cats) > 0 {
			cname := readStrField(cats[0], "name")
			if cname != "" {
				r.run(fmt.Sprintf("get_node_category       (%s)", cname),
					nil, func() (any, error) {
						return c.GetNodeCategory(ctx, nidStr, cname)
					})
			}
		}
		r.run(fmt.Sprintf("get_node_asset_record       id=%s", nidStr),
			nil, func() (any, error) {
				return c.GetNodeAssetRecord(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_hardware_inventory id=%s", nidStr),
			"requires opennms-plugin-provisioning-snmp-hardware-inventory",
			nil, func() (any, error) {
				return c.GetNodeHardwareInventory(ctx, nidStr)
			})
		r.run(fmt.Sprintf("get_node_metadata           id=%s", nidStr),
			nil, func() (any, error) {
				return c.GetNodeMetadata(ctx, nidStr)
			})
		if ip != "" {
			r.warn(fmt.Sprintf("get_interface_metadata    ip=%s", ip),
				"may 404 if no metadata on interface",
				nil, func() (any, error) {
					return c.GetInterfaceMetadata(ctx, nidStr, ip)
				})
			if svc != "" {
				r.warn(fmt.Sprintf("get_service_metadata     svc=%s", svc),
					"may 404 if no metadata on service",
					nil, func() (any, error) {
						return c.GetServiceMetadata(ctx, nidStr, ip, svc)
					})
			}
		}
		r.run(fmt.Sprintf("get_node_outages            id=%s", nidStr),
			nil, func() (any, error) {
				return c.GetNodeOutages(ctx, toInt(nid))
			})
		r.warn(fmt.Sprintf("get_resources_for_node      id=%s", nidStr),
			"", nil, func() (any, error) {
				return c.GetResourcesForNode(ctx, nidStr)
			})
	} else {
		for _, lbl := range []string{"get_node", "get_node_ip_interfaces",
			"get_node_ip_interface", "get_node_ip_services",
			"get_node_snmp_interfaces", "get_node_snmp_interface",
			"get_node_categories", "get_node_category",
			"get_node_asset_record", "get_node_hardware_inventory",
			"get_node_metadata", "get_interface_metadata",
			"get_node_outages", "get_resources_for_node"} {
			r.skip(lbl, "no nodes")
		}
	}
}

func readOutages(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("outages")
	r.run("get_outages", n("outage"), func() (any, error) {
		return c.GetOutages(ctx, &opennms.ListOptions{Limit: 5}, nil)
	})
	r.warn("get_outage_count", "", readStr, func() (any, error) {
		return c.GetOutageCount(ctx)
	})

	_, oid := first(func(limit int) (any, error) {
		return c.GetOutages(ctx, &opennms.ListOptions{Limit: limit}, nil)
	}, "outage", "")
	if readTruthy(oid) {
		r.run(fmt.Sprintf("get_outage  id=%s", readID(oid)), nil,
			func() (any, error) {
				return c.GetOutage(ctx, toInt(oid))
			})
	} else {
		r.skip("get_outage", "no outages")
	}
}

func readRequisitions(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("requisitions")
	result, _ := r.run("get_requisitions", n("model-import"),
		func() (any, error) {
			return c.GetRequisitions(ctx)
		})
	r.warn("get_requisition_count", "", readStr, func() (any, error) {
		return c.GetRequisitionCount(ctx)
	})
	r.run("get_deployed_requisitions", n("model-import"),
		func() (any, error) {
			return c.GetDeployedRequisitions(ctx)
		})
	r.warn("get_deployed_requisition_count", "", readStr,
		func() (any, error) {
			return c.GetDeployedRequisitionCount(ctx)
		})
	reqs := readItems(result, "model-import", "requisition")
	if len(reqs) > 0 {
		rname := readStrField(reqs[0], "foreign-source")
		if rname == "" {
			rname = readStrField(reqs[0], "name")
		}
		if rname != "" {
			r.run(fmt.Sprintf("get_requisition       (%s)", rname), nil,
				func() (any, error) {
					return c.GetRequisition(ctx, rname)
				})
			nodesR, _ := r.run(
				fmt.Sprintf("get_requisition_nodes (%s)", rname), nil,
				func() (any, error) {
					return c.GetRequisitionNodes(ctx, rname)
				})
			if rnodes := readItems(nodesR, "node"); len(rnodes) > 0 {
				fid := readStrField(rnodes[0], "foreign-id")
				if fid != "" {
					r.run(fmt.Sprintf(
						"get_requisition_node            (%s)", fid),
						nil, func() (any, error) {
							return c.GetRequisitionNode(ctx, rname, fid)
						})
					r.run(fmt.Sprintf(
						"get_requisition_node_interfaces (%s)", fid),
						nil, func() (any, error) {
							return c.GetRequisitionNodeInterfaces(
								ctx, rname, fid)
						})
					r.run(fmt.Sprintf(
						"get_requisition_node_categories (%s)", fid),
						nil, func() (any, error) {
							return c.GetRequisitionNodeCategories(
								ctx, rname, fid)
						})
					r.run(fmt.Sprintf(
						"get_requisition_node_assets     (%s)", fid),
						nil, func() (any, error) {
							return c.GetRequisitionNodeAssets(
								ctx, rname, fid)
						})
				}
			}
		}
	} else {
		r.skip("get_requisition / get_requisition_nodes",
			"no requisitions")
	}
}

func readForeignSources(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("foreign sources")
	result, _ := r.run("get_foreign_sources", nil, func() (any, error) {
		return c.GetForeignSources(ctx)
	})
	r.run("get_deployed_foreign_sources", nil, func() (any, error) {
		return c.GetDeployedForeignSources(ctx)
	})
	r.warn("get_deployed_foreign_source_count", "", readStr,
		func() (any, error) {
			return c.GetDeployedForeignSourceCount(ctx)
		})
	r.run("get_default_foreign_source", nil, func() (any, error) {
		return c.GetDefaultForeignSource(ctx)
	})
	sources := readItems(result, "foreignSource")
	if len(sources) > 0 {
		sname := readStrField(sources[0], "name")
		if sname != "" {
			r.run(fmt.Sprintf("get_foreign_source           (%s)", sname),
				nil, func() (any, error) {
					return c.GetForeignSource(ctx, sname)
				})
			detR, _ := r.run(
				fmt.Sprintf("get_foreign_source_detectors (%s)", sname),
				nil, func() (any, error) {
					return c.GetForeignSourceDetectors(ctx, sname)
				})
			if dets := readItems(detR, "detector"); len(dets) > 0 {
				dname := readStrField(dets[0], "name")
				if dname != "" {
					r.run(fmt.Sprintf(
						"get_foreign_source_detector  (%s)", dname),
						nil, func() (any, error) {
							return c.GetForeignSourceDetector(
								ctx, sname, dname)
						})
				}
			}
			polR, _ := r.run(
				fmt.Sprintf("get_foreign_source_policies  (%s)", sname),
				nil, func() (any, error) {
					return c.GetForeignSourcePolicies(ctx, sname)
				})
			if pols := readItems(polR, "policy"); len(pols) > 0 {
				pname := readStrField(pols[0], "name")
				if pname != "" {
					r.run(fmt.Sprintf(
						"get_foreign_source_policy    (%s)", pname),
						nil, func() (any, error) {
							return c.GetForeignSourcePolicy(
								ctx, sname, pname)
						})
				}
			}
		}
	} else {
		r.skip("get_foreign_source / detectors / policies",
			"no foreign sources")
	}
}

func readSnmpConfig(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("snmp config")
	// 127.0.0.1 may not be explicitly configured; a 404 is expected
	// on clean servers.
	r.run("get_snmp_config (127.0.0.1)", nil, func() (any, error) {
		return c.GetSnmpConfig(ctx, "127.0.0.1", "")
	})
}

func readGroups(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("groups")
	result, _ := r.run("get_groups", n("group"), func() (any, error) {
		return c.GetGroups(ctx)
	})
	groups := readItems(result, "group")
	if len(groups) > 0 {
		name := readStrField(groups[0], "name")
		r.run(fmt.Sprintf("get_group            (%s)", name), nil,
			func() (any, error) {
				return c.GetGroup(ctx, name)
			})
		r.run(fmt.Sprintf("get_group_users      (%s)", name), nil,
			func() (any, error) {
				return c.GetGroupUsers(ctx, name)
			})
		r.run(fmt.Sprintf("get_group_categories (%s)", name), nil,
			func() (any, error) {
				return c.GetGroupCategories(ctx, name)
			})
	} else {
		r.skip("get_group / get_group_users / get_group_categories",
			"no groups")
	}
}

func readUsers(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("users")
	result, _ := r.run("get_users", n("user"), func() (any, error) {
		return c.GetUsers(ctx)
	})
	users := readItems(result, "user")
	if len(users) > 0 {
		uname := readStrField(users[0], "user-id")
		r.run(fmt.Sprintf("get_user  (%s)", uname), nil,
			func() (any, error) {
				return c.GetUser(ctx, uname)
			})
	} else {
		r.skip("get_user", "no users")
	}
}

func readCategories(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("categories")
	result, _ := r.run("get_categories", n("category"),
		func() (any, error) {
			return c.GetCategories(ctx)
		})
	cats := readItems(result, "category")
	if len(cats) > 0 {
		name := readStrField(cats[0], "name")
		r.run(fmt.Sprintf("get_category  (%s)", name), nil,
			func() (any, error) {
				return c.GetCategory(ctx, name)
			})
	} else {
		r.skip("get_category", "no categories")
	}
}

func readSchedOutages(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("scheduled outages")
	result, _ := r.run("get_sched_outages", n("schedules"),
		func() (any, error) {
			return c.GetSchedOutages(ctx)
		})
	schedules := readItems(result, "schedules")
	if len(schedules) > 0 {
		name := readStrField(schedules[0], "name")
		r.run(fmt.Sprintf("get_sched_outage  (%s)", name), nil,
			func() (any, error) {
				return c.GetSchedOutage(ctx, name)
			})
	} else {
		r.skip("get_sched_outage", "no scheduled outages")
	}
}

func readKscReports(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("ksc reports")
	result, _ := r.run("get_ksc_reports", n("kscReport"),
		func() (any, error) {
			return c.GetKscReports(ctx)
		})
	r.warn("get_ksc_report_count", "", readStr, func() (any, error) {
		return c.GetKscReportCount(ctx)
	})
	var reports []any
	if m, ok := result.(map[string]any); ok {
		reports, _ = m["kscReport"].([]any)
	}
	if len(reports) > 0 {
		rid := readField(reports[0], "id")
		r.run(fmt.Sprintf("get_ksc_report  id=%s", readID(rid)), nil,
			func() (any, error) {
				return c.GetKscReport(ctx, toInt(rid))
			})
	} else {
		r.skip("get_ksc_report", "no KSC reports")
	}
}

func readResources(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("resources")
	r.warn("get_resources", "", nil, func() (any, error) {
		return c.GetResources(ctx, 1)
	})
}

func readMeasurements(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("measurements")
	// Measurements require a node with collected SNMP performance
	// data. Skip gracefully if no suitable resources can be found.
	var nodes []any
	if nr, err := c.GetNodes(ctx, &opennms.ListOptions{Limit: 1},
		nil); err == nil {
		nodes = readItems(nr, "node")
	}
	if len(nodes) == 0 {
		r.skip("get_measurements", "no nodes")
		return
	}
	nid := readField(nodes[0], "id")
	var ifaceRes map[string]any
	if rr, err := c.GetResourcesForNode(ctx, readID(nid)); err == nil {
		var children []any
		if nodeRes, ok := rr["resource"].(map[string]any); ok {
			if ch, ok := nodeRes["children"].(map[string]any); ok {
				children, _ = ch["resource"].([]any)
			}
		}
		for _, child := range children {
			m, ok := child.(map[string]any)
			if !ok {
				continue
			}
			if id, _ := m["id"].(string); strings.Contains(
				id, "interfaceSnmp") {
				ifaceRes = m
				break
			}
		}
	}
	if ifaceRes == nil {
		r.skip("get_measurements",
			"no SNMP interface resources on first node")
		return
	}
	rid, _ := ifaceRes["id"].(string)
	r.run("get_measurements  ifInOctets", nil, func() (any, error) {
		return c.GetMeasurements(ctx, rid, "ifInOctets", nil)
	})
}

func readHeatmap(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("heatmap")
	r.run("get_heatmap_outages_categories", nil, func() (any, error) {
		return c.GetHeatmapOutagesCategories(ctx)
	})
	r.run("get_heatmap_outages_foreign_sources", nil, func() (any, error) {
		return c.GetHeatmapOutagesForeignSources(ctx)
	})
	r.run("get_heatmap_outages_monitored_services", nil,
		func() (any, error) {
			return c.GetHeatmapOutagesMonitoredServices(ctx)
		})
	r.run("get_heatmap_alarms_categories", nil, func() (any, error) {
		return c.GetHeatmapAlarmsCategories(ctx)
	})
	r.run("get_heatmap_alarms_foreign_sources", nil, func() (any, error) {
		return c.GetHeatmapAlarmsForeignSources(ctx)
	})
	r.run("get_heatmap_alarms_monitored_services", nil,
		func() (any, error) {
			return c.GetHeatmapAlarmsMonitoredServices(ctx)
		})
	// nodes_by methods require a grouping key; discover one from
	// categories
	var catName string
	if cats, err := c.GetCategories(ctx); err == nil {
		if items := readItems(cats, "category"); len(items) > 0 {
			catName = readStrField(items[0], "name")
		}
	}
	if catName != "" {
		r.run(fmt.Sprintf("get_heatmap_outages_nodes_by_category (%s)",
			catName), nil, func() (any, error) {
			return c.GetHeatmapOutagesNodesByCategory(ctx, catName)
		})
		r.run(fmt.Sprintf("get_heatmap_alarms_nodes_by_category  (%s)",
			catName), nil, func() (any, error) {
			return c.GetHeatmapAlarmsNodesByCategory(ctx, catName)
		})
	} else {
		r.skip("heatmap nodes_by_category", "no categories")
	}
	var fsName string
	if fss, err := c.GetForeignSources(ctx); err == nil {
		if items := readItems(fss, "foreignSource"); len(items) > 0 {
			fsName = readStrField(items[0], "name")
		}
	}
	if fsName != "" {
		r.run(fmt.Sprintf(
			"get_heatmap_outages_nodes_by_foreign_source (%s)", fsName),
			nil, func() (any, error) {
				return c.GetHeatmapOutagesNodesByForeignSource(ctx, fsName)
			})
		r.run(fmt.Sprintf(
			"get_heatmap_alarms_nodes_by_foreign_source  (%s)", fsName),
			nil, func() (any, error) {
				return c.GetHeatmapAlarmsNodesByForeignSource(ctx, fsName)
			})
	} else {
		r.skip("heatmap nodes_by_foreign_source", "no foreign sources")
	}
}

func readMaps(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("maps")
	result, _ := r.warn("get_maps",
		"SVG maps may not be available in all versions",
		n("map"), func() (any, error) {
			return c.GetMaps(ctx)
		})
	maps := readItems(result, "map")
	if len(maps) > 0 {
		mid := readField(maps[0], "id")
		r.run(fmt.Sprintf("get_map          id=%s", readID(mid)), nil,
			func() (any, error) {
				return c.GetMap(ctx, toInt(mid))
			})
		r.run(fmt.Sprintf("get_map_elements id=%s", readID(mid)), nil,
			func() (any, error) {
				return c.GetMapElements(ctx, toInt(mid))
			})
	} else {
		r.skip("get_map / get_map_elements", "no maps")
	}
}

func readGraphs(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("topology graphs")
	result, _ := r.run("get_graph_containers", nil, func() (any, error) {
		return c.GetGraphContainers(ctx)
	})
	// Response key varies by OpenNMS version
	containers := readItems(result, "graphContainers", "container")
	if len(containers) > 0 {
		cid := readID(readField(containers[0], "id"))
		r.run(fmt.Sprintf("get_graph_container  (%s)", cid), nil,
			func() (any, error) {
				return c.GetGraphContainer(ctx, cid)
			})
		graphs, _ := readField(containers[0], "graphs").([]any)
		if len(graphs) > 0 {
			ns := readStrField(graphs[0], "namespace")
			r.run(fmt.Sprintf("get_graph  (%s/%s)", cid, ns), nil,
				func() (any, error) {
					return c.GetGraph(ctx, cid, ns)
				})
		} else {
			r.skip("get_graph", "container has no graphs")
		}
	} else {
		r.skip("get_graph_container / get_graph", "no graph containers")
	}
}

func readFlows(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("flows")
	flowNote := "requires flow persistence (Elasticsearch/OpenSearch)"
	r.warn("get_flow_count", flowNote, readStr, func() (any, error) {
		return c.GetFlowCount(ctx)
	})
	r.run("get_flow_exporters", n("exporters"), func() (any, error) {
		return c.GetFlowExporters(ctx)
	})
	r.warn("get_flow_applications", flowNote, nil, func() (any, error) {
		return c.GetFlowApplications(ctx, &opennms.FlowOptions{TopN: 5})
	})
	r.warn("get_flow_applications_enumerate", flowNote, nil,
		func() (any, error) {
			return c.GetFlowApplicationsEnumerate(ctx,
				&opennms.FlowOptions{Limit: 5})
		})
	r.warn("get_flow_conversations", flowNote, nil, func() (any, error) {
		return c.GetFlowConversations(ctx, &opennms.FlowOptions{TopN: 5})
	})
	r.warn("get_flow_conversations_enumerate", flowNote, nil,
		func() (any, error) {
			return c.GetFlowConversationsEnumerate(ctx,
				&opennms.FlowOptions{Limit: 5})
		})
	r.warn("get_flow_hosts", flowNote, nil, func() (any, error) {
		return c.GetFlowHosts(ctx, &opennms.FlowOptions{TopN: 5})
	})
	r.warn("get_flow_hosts_enumerate", flowNote, nil, func() (any, error) {
		return c.GetFlowHostsEnumerate(ctx, &opennms.FlowOptions{Limit: 5})
	})
}

func readDeviceConfig(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("device config")
	r.run("get_device_configs", n("deviceConfig"), func() (any, error) {
		return c.GetDeviceConfigs(ctx, &opennms.ListOptions{Limit: 5},
			"", "", "", 0, 0)
	})
	r.run("get_latest_device_configs", n("deviceConfig"),
		func() (any, error) {
			return c.GetLatestDeviceConfigs(ctx,
				&opennms.ListOptions{Limit: 5}, "", "")
		})
}

func readSituations(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("situations (v2)")
	r.warn("get_situations", "requires Alarmd situation correlation",
		n("alarm"), func() (any, error) {
			return c.GetSituations(ctx, &opennms.ListOptions{Limit: 5})
		})
}

func readBusinessServices(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("business services (v2)")
	result, _ := r.run("get_business_services", n("business-service"),
		func() (any, error) {
			return c.GetBusinessServices(ctx)
		})
	bsvcs := readItems(result, "business-service", "business-services")
	if len(bsvcs) > 0 {
		bid := readField(bsvcs[0], "id")
		if readTruthy(bid) {
			r.run(fmt.Sprintf("get_business_service  id=%s", readID(bid)),
				nil, func() (any, error) {
					return c.GetBusinessService(ctx, toInt(bid))
				})
		}
	} else {
		r.skip("get_business_service", "no business services")
	}
}

func readEnlinkd(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("enlinkd (v2)")
	_, nid := readFirstNode(ctx, c)
	enlinkdNote := "requires enlinkd daemon"
	if readTruthy(nid) {
		nidStr := readID(nid)
		r.warn(fmt.Sprintf("get_node_enlinkd         id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeEnlinkd(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_lldp_links      id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeLldpLinks(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_cdp_links       id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeCdpLinks(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_ospf_links      id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeOspfLinks(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_isis_links      id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeIsisLinks(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_bridge_links    id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeBridgeLinks(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_lldp_element    id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeLldpElement(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_cdp_element     id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeCdpElement(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_ospf_element    id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeOspfElement(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_isis_element    id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeIsisElement(ctx, nidStr)
			})
		r.warn(fmt.Sprintf("get_node_bridge_elements id=%s", nidStr),
			enlinkdNote, nil, func() (any, error) {
				return c.GetNodeBridgeElements(ctx, nidStr)
			})
	} else {
		r.skip("get_node_enlinkd", "no nodes")
	}
}

func readV2Interfaces(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("ip/snmp interfaces (v2)")
	r.run("get_ip_interfaces", n("ipInterface"), func() (any, error) {
		return c.GetIpInterfaces(ctx, "", &opennms.ListOptions{Limit: 5})
	})
	r.run("get_snmp_interfaces", n("snmpInterface"), func() (any, error) {
		return c.GetSnmpInterfaces(ctx, "", &opennms.ListOptions{Limit: 5})
	})
}

func readMonitoringLocations(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("monitoring locations")
	r.run("get_monitoring_locations", n("location"), func() (any, error) {
		return c.GetMonitoringLocations(ctx, nil)
	})
	r.run("get_monitoring_location_count", readStr, func() (any, error) {
		return c.GetMonitoringLocationCount(ctx)
	})
	r.run("get_default_monitoring_location", nil, func() (any, error) {
		return c.GetDefaultMonitoringLocation(ctx)
	})
}

func readMinions(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("minions")
	r.run("get_minions", n("minion"), func() (any, error) {
		return c.GetMinions(ctx, nil)
	})
	r.run("get_minion_count", readStr, func() (any, error) {
		return c.GetMinionCount(ctx)
	})
}

func readIfservices(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("ifservices")
	r.run("get_ifservices", n("service"), func() (any, error) {
		return c.GetIfservices(ctx, nil)
	})
	r.run("get_ifservices_v2", n("service"), func() (any, error) {
		return c.GetIfservicesV2(ctx, "", &opennms.ListOptions{Limit: 5})
	})
}

func readAvailability(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("availability")
	r.run("get_availability", nil, func() (any, error) {
		return c.GetAvailability(ctx)
	})
	_, nid := readFirstNode(ctx, c)
	if readTruthy(nid) {
		r.run(fmt.Sprintf("get_availability_node  id=%s", readID(nid)),
			nil, func() (any, error) {
				return c.GetAvailabilityNode(ctx, toInt(nid))
			})
	} else {
		r.skip("get_availability_node", "no nodes")
	}
}

func readHealth(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("health")
	r.run("get_health", nil, func() (any, error) {
		return c.GetHealth(ctx, "")
	})
	r.warn("get_health_probe",
		"returns 599 if any health check is unhealthy",
		nil, func() (any, error) {
			return c.GetHealthProbe(ctx)
		})
}

func readWhoami(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("whoami")
	r.run("get_whoami", func(res any) string {
		return readFieldStr(res, "id")
	}, func() (any, error) {
		return c.GetWhoami(ctx)
	})
}

func readMonitoringSystems(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("monitoring systems")
	r.warn("get_monitoring_system",
		"may not be available on all versions",
		nil, func() (any, error) {
			return c.GetMonitoringSystem(ctx)
		})
}

func readPrefabGraphs(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("prefab graphs")
	r.run("get_prefab_graph_names", nil, func() (any, error) {
		return c.GetPrefabGraphNames(ctx)
	})
	_, nid := readFirstNode(ctx, c)
	if readTruthy(nid) {
		r.warn(fmt.Sprintf("get_prefab_graphs_for_node  id=%s",
			readID(nid)),
			"may time out on large node inventories",
			nil, func() (any, error) {
				return c.GetPrefabGraphsForNode(ctx, readID(nid))
			})
	} else {
		r.skip("get_prefab_graphs_for_node", "no nodes")
	}
}

func readFlowDscp(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("flow DSCP")
	note := "requires flow persistence (Elasticsearch/OpenSearch)"
	r.warn("get_flow_dscp", note, nil, func() (any, error) {
		return c.GetFlowDscp(ctx, &opennms.FlowOptions{TopN: 5})
	})
	r.warn("get_flow_dscp_enumerate", note, nil, func() (any, error) {
		return c.GetFlowDscpEnumerate(ctx, &opennms.FlowOptions{Limit: 5})
	})
	r.warn("get_flow_graph_url", note, nil, func() (any, error) {
		return c.GetFlowGraphUrl(ctx)
	})
}

func readBusinessServiceFunctions(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("business service functions (v2)")
	mapR, _ := r.run("get_map_functions", nil, func() (any, error) {
		return c.GetMapFunctions(ctx)
	})
	reduceR, _ := r.run("get_reduce_functions", nil, func() (any, error) {
		return c.GetReduceFunctions(ctx)
	})
	// Extract a function name to test individual lookup
	if mfuncs := readItems(mapR, "functions"); len(mfuncs) > 0 {
		mname, ok := mfuncs[0].(string)
		if !ok {
			mname = readStrField(mfuncs[0], "name")
		}
		if mname != "" {
			r.run(fmt.Sprintf("get_map_function  (%s)", mname), nil,
				func() (any, error) {
					return c.GetMapFunction(ctx, mname)
				})
		}
	}
	if rfuncs := readItems(reduceR, "functions"); len(rfuncs) > 0 {
		rname, ok := rfuncs[0].(string)
		if !ok {
			rname = readStrField(rfuncs[0], "name")
		}
		if rname != "" {
			r.run(fmt.Sprintf("get_reduce_function  (%s)", rname), nil,
				func() (any, error) {
					return c.GetReduceFunction(ctx, rname)
				})
		}
	}
}

func readClassifications(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("classifications")
	clsNote := "requires flow classification plugin"
	rulesR, _ := r.warn("get_classification_rules", clsNote, nil,
		func() (any, error) {
			return c.GetClassificationRules(ctx, nil)
		})
	groupsR, _ := r.warn("get_classification_groups", clsNote, nil,
		func() (any, error) {
			return c.GetClassificationGroups(ctx, nil)
		})
	r.warn("get_classification_protocols", clsNote, nil,
		func() (any, error) {
			return c.GetClassificationProtocols(ctx)
		})
	// Drill into first rule if available
	if rules := readItems(rulesR, "rule", "classification"); len(rules) > 0 {
		rid := readField(rules[0], "id")
		if readTruthy(rid) {
			r.warn(fmt.Sprintf("get_classification_rule   id=%s",
				readID(rid)), clsNote, nil, func() (any, error) {
				return c.GetClassificationRule(ctx, toInt(rid))
			})
		}
	}
	// Drill into first group if available
	if groups := readItems(groupsR, "group"); len(groups) > 0 {
		gid := readField(groups[0], "id")
		if readTruthy(gid) {
			r.warn(fmt.Sprintf("get_classification_group  id=%s",
				readID(gid)), clsNote, nil, func() (any, error) {
				return c.GetClassificationGroup(ctx, toInt(gid))
			})
		}
	}
}

func readSituationFeedback(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("situation feedback")
	r.warn("get_situation_feedback_tags",
		"requires situation correlation + feedback feature",
		nil, func() (any, error) {
			return c.GetSituationFeedbackTags(ctx, "")
		})
}

func readUserDefinedLinks(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("user-defined links (v2)")
	result, _ := r.run("get_user_defined_links", nil, func() (any, error) {
		return c.GetUserDefinedLinks(ctx)
	})
	links := readItems(result, "user-defined-link")
	if len(links) > 0 {
		lid := readField(links[0], "id")
		if readTruthy(lid) {
			r.run(fmt.Sprintf("get_user_defined_link  id=%s", readID(lid)),
				nil, func() (any, error) {
					return c.GetUserDefinedLink(ctx, toInt(lid))
				})
		}
	} else {
		r.skip("get_user_defined_link", "no user-defined links")
	}
}

func readApplications(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("applications (v2)")
	result, _ := r.run("get_applications", n("application"),
		func() (any, error) {
			return c.GetApplications(ctx, nil)
		})
	apps := readItems(result, "application")
	if len(apps) > 0 {
		aid := readField(apps[0], "id")
		if readTruthy(aid) {
			r.run(fmt.Sprintf("get_application  id=%s", readID(aid)), nil,
				func() (any, error) {
					return c.GetApplication(ctx, toInt(aid))
				})
		}
	} else {
		r.skip("get_application", "no applications")
	}
}

func readPerspectivePoller(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("perspective poller (v2)")
	r.warn("get_perspective_poller_status (app 1)",
		"requires perspective poller + application with id=1",
		nil, func() (any, error) {
			return c.GetPerspectivePollerStatus(ctx, 1, 0, 0)
		})
}

func readForeignSourcesConfig(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("foreign sources config")
	r.run("get_foreign_source_config_policies", nil, func() (any, error) {
		return c.GetForeignSourceConfigPolicies(ctx)
	})
	r.run("get_foreign_source_config_detectors", nil, func() (any, error) {
		return c.GetForeignSourceConfigDetectors(ctx)
	})
	r.run("get_foreign_source_config_assets", nil, func() (any, error) {
		return c.GetForeignSourceConfigAssets(ctx)
	})
	r.run("get_foreign_source_config_categories", nil,
		func() (any, error) {
			return c.GetForeignSourceConfigCategories(ctx)
		})
}

func readRequisitionNames(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("requisition names")
	r.run("get_requisition_names", nil, func() (any, error) {
		return c.GetRequisitionNames(ctx)
	})
}

func readSnmpMetadata(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("snmp metadata (v2)")
	_, nid := readFirstNode(ctx, c)
	if readTruthy(nid) {
		r.run(fmt.Sprintf("get_snmp_metadata  id=%s", readID(nid)), nil,
			func() (any, error) {
				return c.GetSnmpMetadata(ctx, toInt(nid))
			})
	} else {
		r.skip("get_snmp_metadata", "no nodes")
	}
}

func readProvisiond(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("provisiond (v2)")
	r.warn("get_provisiond_status", "requires provisiond v2 API support",
		nil, func() (any, error) {
			return c.GetProvisiondStatus(ctx)
		})
}

func readEventconf(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("eventconf (v2)")
	ecNote := "requires eventconf v2 API support"
	namesR, _ := r.warn("get_eventconf_source_names", ecNote, nil,
		func() (any, error) {
			return c.GetEventconfSourceNames(ctx)
		})
	r.warn("get_eventconf_filter", ecNote, nil, func() (any, error) {
		return c.GetEventconfFilter(ctx, "", "", nil)
	})
	r.warn("get_eventconf_filter_sources", ecNote, nil,
		func() (any, error) {
			return c.GetEventconfFilterSources(ctx, "", nil)
		})
	// Drill into first source if source names returned
	var srcName string
	if names, ok := namesR.([]any); ok && len(names) > 0 {
		if readTruthy(names[0]) {
			srcName = fmt.Sprintf("%v", names[0])
		}
	}
	if srcName != "" {
		r.warn(fmt.Sprintf("get_eventconf_source          (%s)", srcName),
			ecNote, nil, func() (any, error) {
				return c.GetEventconfSource(ctx, srcName)
			})
		r.warn(fmt.Sprintf("get_eventconf_filter_events   (%s)", srcName),
			ecNote, nil, func() (any, error) {
				return c.GetEventconfFilterEvents(ctx, srcName, nil)
			})
		r.warn(fmt.Sprintf("download_eventconf_events     (%s)", srcName),
			ecNote, nil, func() (any, error) {
				return c.DownloadEventconfEvents(ctx, srcName)
			})
	}
}

func readAssetSuggestions(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("asset suggestions")
	r.run("get_asset_suggestions", nil, func() (any, error) {
		return c.GetAssetSuggestions(ctx)
	})
}

func readScv(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("secure credentials vault")
	result, _ := r.warn("get_credentials",
		"requires SCV REST API support", nil, func() (any, error) {
			return c.GetCredentials(ctx)
		})
	creds := readItems(result, "credential")
	if len(creds) > 0 {
		alias, ok := creds[0].(string)
		if !ok {
			alias = readStrField(creds[0], "alias")
		}
		if alias != "" {
			r.warn(fmt.Sprintf("get_credential  (%s)", alias),
				"requires SCV REST API support", nil,
				func() (any, error) {
					return c.GetCredential(ctx, alias)
				})
		}
	}
}

func readConfigMgmt(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("configuration management")
	namesR, _ := r.run("get_config_names", nil, func() (any, error) {
		return c.GetConfigNames(ctx)
	})
	r.warn("get_config_schemas",
		"may return empty body on some versions",
		nil, func() (any, error) {
			return c.GetConfigSchemas(ctx)
		})
	cfgNames, _ := namesR.([]any)
	if len(cfgNames) > 0 {
		cname := fmt.Sprintf("%v", cfgNames[0])
		r.warn(fmt.Sprintf("get_config_schema  (%s)", cname),
			"may return empty body on some versions", nil,
			func() (any, error) {
				return c.GetConfigSchema(ctx, cname)
			})
		r.warn(fmt.Sprintf("get_config_ids     (%s)", cname),
			"may return empty body on some versions", nil,
			func() (any, error) {
				return c.GetConfigIds(ctx, cname)
			})
	}
}

func readSnmptrapNbi(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("SNMP trap NBI config")
	r.warn("get_snmptrap_nbi_config", "requires SNMP trap NBI plugin",
		nil, func() (any, error) {
			return c.GetSnmptrapNbiConfig(ctx)
		})
	r.warn("get_snmptrap_nbi_status", "requires SNMP trap NBI plugin",
		nil, func() (any, error) {
			return c.GetSnmptrapNbiStatus(ctx)
		})
}

func readEmailNbi(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("email NBI config")
	r.warn("get_email_nbi_config", "requires email NBI plugin", nil,
		func() (any, error) {
			return c.GetEmailNbiConfig(ctx)
		})
	r.warn("get_email_nbi_status", "requires email NBI plugin", nil,
		func() (any, error) {
			return c.GetEmailNbiStatus(ctx)
		})
}

func readSyslogNbi(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("syslog NBI config")
	r.warn("get_syslog_nbi_config", "requires syslog NBI plugin", nil,
		func() (any, error) {
			return c.GetSyslogNbiConfig(ctx)
		})
	r.warn("get_syslog_nbi_status", "requires syslog NBI plugin", nil,
		func() (any, error) {
			return c.GetSyslogNbiStatus(ctx)
		})
}

func readJavamailConfig(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("javamail config")
	r.warn("get_javamail_default_config",
		"may not be available on all versions",
		nil, func() (any, error) {
			return c.GetJavamailDefaultConfig(ctx)
		})
	r.run("get_javamail_readmails", nil, func() (any, error) {
		return c.GetJavamailReadmails(ctx)
	})
	r.run("get_javamail_sendmails", nil, func() (any, error) {
		return c.GetJavamailSendmails(ctx)
	})
	r.run("get_javamail_end2ends", nil, func() (any, error) {
		return c.GetJavamailEnd2ends(ctx)
	})
}

func readPagination(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("pagination")
	// Paginate nodes with a tiny page size and verify we get the
	// full count.
	totalR, _ := r.run("get_node_count", nil, func() (any, error) {
		return c.GetNodeCount(ctx)
	})
	total, isInt := totalR.(int)
	if isInt && total > 0 {
		capN := min(total, 15) // don't fetch thousands
		count := 0
		var pageErr error
		for _, err := range opennms.Paginate(
			func(limit, offset int) (map[string]any, error) {
				return c.GetNodes(ctx, &opennms.ListOptions{
					Limit: limit, Offset: offset,
				}, nil)
			}, "node", 5) {
			if err != nil {
				pageErr = err
				break
			}
			count++
			if count >= capN {
				break
			}
		}
		if pageErr != nil {
			r.fail("paginate(get_nodes, page_size=5)", pageErr)
		} else {
			r.ok("paginate(get_nodes, page_size=5)",
				fmt.Sprintf("%d items (cap %d)", count, capN))
		}
	} else {
		r.skip("paginate(get_nodes)", "no nodes or count unavailable")
	}
}

func readExceptions(ctx context.Context, c *opennms.Client, r *runner) {
	r.section("exception hierarchy")
	// Verify that a 404 comes back as ErrNotFound, not a bare error.
	_, err := c.GetNode(ctx, "999999999")
	switch {
	case err == nil:
		r.fail("NotFoundError on missing node",
			errors.New("expected NotFoundError but call succeeded"))
	case errors.Is(err, opennms.ErrNotFound):
		r.ok("NotFoundError on missing node",
			"get_node(999999999) raised NotFoundError")
	default:
		r.fail("NotFoundError on missing node",
			fmt.Errorf("expected NotFoundError, got %T: %v", err, err))
	}
}
