package open

import (
	"testing"
	// "github.com/stretchr/testify/assert"
	// "github.com/whatap/golib/net/oneway"
)

func Test_OpenMx(t *testing.T) {

}

func Test_OpenMxHelp(t *testing.T) {
	// oPack := NewOpenMxHelpPack()

	// // # HELP cpu The number of processors available to the Java virtual machine
	// // # TYPE cpu gauge
	// o := NewOpenMxHelp()
	// o.metric = "kong_memory_lua_shared_dict_total_bytes"
	// // o.Put()
}

var sampleJson = `
{
                            "name": "kong_memory_lua_shared_dict_total_bytes",
                            "description": "Total capacity in bytes of a shared_dict",
                            "gauge": {
                                "dataPoints": [
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_cluster_events"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_core_db_cache"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 134217728
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_core_db_cache_miss"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 12582912
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_counters"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_db_cache"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 134217728
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_db_cache_miss"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 12582912
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_healthchecks"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_keyring"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_locks"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 8388608
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_profiling_state"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1572864
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_rate_limiting_counters"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 12582912
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_reports_consumers"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 10485760
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_reports_routes"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_reports_services"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_reports_workspaces"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_secrets"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_vaults_hcv"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_vitals"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_vitals_counters"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 52428800
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "kong_vitals_lists"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 1048576
                                    },
                                    {
                                        "attributes": [
                                            {
                                                "key": "kong_subsystem",
                                                "value": {
                                                    "stringValue": "http"
                                                }
                                            },
                                            {
                                                "key": "node_id",
                                                "value": {
                                                    "stringValue": "424e7127-44ab-44a7-a7c8-236de68e3620"
                                                }
                                            },
                                            {
                                                "key": "shared_dict",
                                                "value": {
                                                    "stringValue": "prometheus_metrics"
                                                }
                                            }
                                        ],
                                        "timeUnixNano": "1733981893796000000",
                                        "asDouble": 5242880
                                    }
                                ]
                            }
                        }
`
