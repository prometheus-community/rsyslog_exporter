// Copyright 2024 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
