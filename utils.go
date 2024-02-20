package main

import "strings"

func getStatType(buf []byte) rsyslogType {
	line := string(buf)
	if strings.Contains(line, "processed") {
		return rsyslogAction
	} else if strings.Contains(line, "\"name\": \"omkafka\"") {
		// Not checking for just omkafka here as multiple actions may/will contain that word.
		// omkafka lines have a submitted field, so they need to be filtered before rsyslogInput
		return rsyslogOmkafka
	} else if strings.Contains(line, "submitted") {
		return rsyslogInput
	} else if strings.Contains(line, "called.recvmmsg") {
		return rsyslogInputIMDUP
	} else if strings.Contains(line, "enqueued") {
		return rsyslogQueue
	} else if strings.Contains(line, "utime") {
		return rsyslogResource
	} else if strings.Contains(line, "dynstats") {
		return rsyslogDynStat
	} else if strings.Contains(line, "dynafile cache") {
		return rsyslogDynafileCache
	} else if strings.Contains(line, "omfwd") {
		return rsyslogForward
	} else if strings.Contains(line, "mmkubernetes") {
		return rsyslogKubernetes
	}
	return rsyslogUnknown
}
