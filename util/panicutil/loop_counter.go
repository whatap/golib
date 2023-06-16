package panicutil

const (
	MAX_COUNTERS                            = 8192
	C_AGENT_BOOT_BootWindowsService         = 10
	C_AGENT_BOOT_sendAgentInfoAsync         = 11
	C_AGENT_BOOT_applyTags                  = 12
	C_AGENT_BOOT_runStackDumpPeriodically   = 13
	C_AGENT_BOOT_runStackDumpPeriodically_2 = 14
	C_AGENT_BOOT_runStackDumpPeriodically_3 = 15
	C_Config_run                            = 16
	C_Config_SetValues                      = 17
	C_Config_SetValues_2                    = 18
	C_Config_SetValues_3                    = 19
	C_Config_SetValues_4                    = 20
	C_Config_SetValues_5                    = 21
	C_Config_SearchKey                      = 22
	C_Config_apply_agentless                = 23
	C_Config_UpdateScripts                  = 24
	C_CONTROLHANDLER_runControl             = 25
	C_CONTROLHANDLER_process                = 26
	C_CONTROLHANDLER_process_2              = 27
	C_CONTROLHANDLER_process_3              = 28
	C_NETSTAT                               = 29
	C_NETSTAT_2                             = 30
	C_NETSTAT3                              = 31
	C_UPDATEPEERLIST                        = 32
	C_UPDATEPEERLIST_2                      = 33
	C_UPDATEPEERLIST_3                      = 34
	C_GETCOMMANDLIST                        = 35
	C_GETCOMMANDLIST_2                      = 36
	C_EXECUTE                               = 37
	C_EXECUTE_2                             = 38
	C_COUNTERMANAGER_StartCounterManager    = 39
	C_COUNTERMANAGER_StartCounterManager_2  = 40
	C_COUNTERMANAGER_StartCounterManager_3  = 41
	C_COUNTERMANAGER_POLL                   = 42
	C_COUNTERMANAGER_POLL_2                 = 43
	C_LOGACTIONMANAGER_HandleLogEvent       = 44
	C_TASKDISK_process                      = 45
	C_TASKLOGEVENT_process                  = 46
	C_TASKNETSTAT_Init                      = 47
	C_TASKPROC_process                      = 48
	C_TASKPROC_process_2                    = 49
	C_TASKPROC_process_3                    = 50
	C_DATATEXT_initial                      = 51
	C_DATATEXT_process                      = 52
	C_SECURITYMASTER_run                    = 53
	C_AGENTLESSCHECKMAIN_process            = 54
	C_AGENTLESSCHECKMAIN_process_2          = 55
	C_AGENTLESSCHECKMAIN_process_3          = 56
	C_AGENTLESSCHECKMAIN_process_4          = 57
	C_AGENTLESSCHECKMAIN_process_5          = 58
	C_AGENTLESSCHECKMAIN_TestRun            = 59
	C_EXECUTE_check                         = 60
	C_ATTRMAIN_process                      = 61
	C_MEMCACHED_GETPERF                     = 62
	C_MEMCACHED_GETPERF_2                   = 63
	C_MEMCACHED_GETPERF_3                   = 64
	C_MEMCACHED_GETPERF_4                   = 65
	C_MEMCACHED_GETPERF_5                   = 66
	C_MEMCACHED_GETPERF_6                   = 67
	C_MEMCACHED_contains                    = 68
	C_REDIS_GETPERF                         = 69
	C_REDIS_GETPERF_2                       = 70
	C_REDIS_GETPERF_3                       = 71
	C_REDIS_containes                       = 72
	C_WEBCHECK_LISTURLS                     = 73
	C_WEBCHECK_send                         = 74
	C_WEBCHECK_check                        = 75
	C_DOWNCHECKMAIN_STARTHELLO_MAIN         = 76
	C_DOWNCHECKMAIN_STARTHELLO              = 77
	C_DOWNCHECKMAIN_STARTHELLO_2            = 78
	C_DOWNCHECKMAIN_PROCESS                 = 79
	C_DOWNCHECKMAIN_PROCESS2                = 80
	C_DOWNCHECKMAIN_CHECK                   = 81
)

var (
	CountIds = []int{C_AGENT_BOOT_BootWindowsService,
		C_AGENT_BOOT_sendAgentInfoAsync,
		C_AGENT_BOOT_applyTags,
		C_AGENT_BOOT_runStackDumpPeriodically,
		C_AGENT_BOOT_runStackDumpPeriodically_2,
		C_AGENT_BOOT_runStackDumpPeriodically_3,
		C_Config_run,
		C_Config_SetValues,
		C_Config_SetValues_2,
	}
	cyclecounts = make([]int64, MAX_COUNTERS)
)

func Cycle(loopid int) {
	cyclecounts[loopid] += 1
}
