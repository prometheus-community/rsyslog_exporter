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

import (
	"encoding/json"
	"fmt"
)

type omkafka struct {
	Name                     string `json:"name"`
	Origin                   string `json:"origin"`
	Submitted                int64  `json:"submitted"`
	MaxOutQSize              int64  `json:"maxoutqsize"`
	Failures                 int64  `json:"failures"`
	TopicDynacacheSkipped    int64  `json:"topicdynacache.skipped"`
	TopicDynacacheMiss       int64  `json:"topicdynacache.miss"`
	TopicDynacacheEvicted    int64  `json:"topicdynacache.evicted"`
	Acked                    int64  `json:"acked"`
	FailuresMsgTooLarge      int64  `json:"failures_msg_too_large"`
	FailuresUnknownTopic     int64  `json:"failures_unknown_topic"`
	FailuresQueueFull        int64  `json:"failures_queue_full"`
	FailuresUnknownPartition int64  `json:"failures_unknown_partition"`
	FailuresOther            int64  `json:"failures_other"`
	ErrorsTimedOut           int64  `json:"errors_timed_out"`
	ErrorsTransport          int64  `json:"errors_transport"`
	ErrorsBrokerDown         int64  `json:"errors_broker_down"`
	ErrorsAuth               int64  `json:"errors_auth"`
	ErrorsSSL                int64  `json:"errors_ssl"`
	ErrorsOther              int64  `json:"errors_other"`
	RttAvgUsec               int64  `json:"rtt_avg_usec"`
	ThrottleAvgMsec          int64  `json:"throttle_avg_msec"`
	IntLatencyAvgUsec        int64  `json:"int_latency_avg_usec"`
}

const (
	messagesDescription       = "number of messages: submitted: messages submitted to omkafka for processing (with both acknowledged deliveries to broker as well as failed or re-submitted from omkafka to librdkafka); failures: messages that librdkafka failed to deliver (broken down into various types in omkafka_failures); acked: messages that were acknowledged by kafka broker. Note that kafka broker provides two levels of delivery acknowledgements depending on topicConfParam: default (acks=1) implies delivery to the leader only while acks=-1 implies delivery to leader as well as replication to all brokers"
	topicDynaCacheDescription = "skipped: dynamic topic cache lookups that find an existing topic and skip creating a new one; miss: dynamic topic cache lookups that fail to find an existing topic and end up creating new one; evicted: dynamic topic cache entry evictions"
	failuresDescription       = "msg_too_large: failed to deliver to the broker because broker considers message to be too large. Note that omkafka may still resubmit to librdkafka depending on resubmitOnFailure option; unknown_topic: failed to deliver to the broker because broker does not recognize the topic; queue_full: dropped by librdkafka when its queue becomes full. Note that default size of librdkafka queue is 100,000 messages; unknown_partition: failed to deliver because broker does not recognize a partition; other: all of the rest of the failures that do not fall in any of the other failure categories"
	errorsDescription         = "timed_out: messages that librdkafka could not deliver within timeout. These errors will cause action to be suspended but messages can be retried depending on retry options; transport: messages that librdkafka could not deliver due to transport errors. These messages can be retried depending on retry options; broker_down: messages that librdkafka could not deliver because it thinks that broker is not accessible. These messages can be retried depending on options; auth: messages that librdkafka could not deliver due to authentication errors. These messages can be retried depending on the options; ssl: messages that librdkafka could not deliver due to ssl errors. These messages can be retried depending on the options; other: rest of librdkafka errors"
)

func newOmkafkaFromJSON(b []byte) (*omkafka, error) {
	var pstat omkafka
	err := json.Unmarshal(b, &pstat)
	if err != nil {
		return nil, fmt.Errorf("failed to decode omkafka stat `%v`: %v", string(b), err)
	}
	return &pstat, nil
}

func (o *omkafka) toPoints() []*point {
	points := make([]*point, 22)

	// A input_submitted metric was always created for omkafka
	// as statType filter matched "submitted". Ensure we still
	// emit that metric for backwards compatibility.
	points[0] = &point{
		Name:        "input_submitted",
		Type:        counter,
		Value:       o.Submitted,
		Description: "messages submitted",
		LabelName:   "input",
		LabelValue:  o.Name,
	}
	points[1] = &point{
		Name:        "omkafka_messages",
		Type:        counter,
		Value:       o.Submitted,
		Description: messagesDescription,
		LabelName:   "type",
		LabelValue:  "submitted",
	}
	points[2] = &point{
		Name:        "omkafka_maxoutqsize",
		Type:        counter,
		Value:       o.MaxOutQSize,
		Description: "high water mark of output queue size",
	}

	points[3] = &point{
		Name:        "omkafka_messages",
		Type:        counter,
		Value:       o.Failures,
		Description: messagesDescription,
		LabelName:   "type",
		LabelValue:  "failures",
	}

	points[4] = &point{
		Name:        "omkafka_topicdynacache",
		Type:        counter,
		Value:       o.TopicDynacacheSkipped,
		Description: topicDynaCacheDescription,
		LabelName:   "type",
		LabelValue:  "skipped",
	}

	points[5] = &point{
		Name:        "omkafka_topicdynacache",
		Type:        counter,
		Value:       o.TopicDynacacheMiss,
		Description: topicDynaCacheDescription,
		LabelName:   "type",
		LabelValue:  "miss",
	}

	points[6] = &point{
		Name:        "omkafka_topicdynacache",
		Type:        counter,
		Value:       o.TopicDynacacheEvicted,
		Description: topicDynaCacheDescription,
		LabelName:   "type",
		LabelValue:  "evicted",
	}

	points[7] = &point{
		Name:        "omkafka_messages",
		Type:        counter,
		Value:       o.Acked,
		Description: messagesDescription,
		LabelName:   "type",
		LabelValue:  "acked",
	}

	points[8] = &point{
		Name:        "omkafka_failures",
		Type:        counter,
		Value:       o.FailuresMsgTooLarge,
		Description: failuresDescription,
		LabelName:   "type",
		LabelValue:  "msg_too_large",
	}

	points[9] = &point{
		Name:        "omkafka_failures",
		Type:        counter,
		Value:       o.FailuresUnknownTopic,
		Description: failuresDescription,
		LabelName:   "type",
		LabelValue:  "unknown_topic",
	}

	points[10] = &point{
		Name:        "omkafka_failures",
		Type:        counter,
		Value:       o.FailuresQueueFull,
		Description: failuresDescription,
		LabelName:   "type",
		LabelValue:  "queue_full",
	}

	points[11] = &point{
		Name:        "omkafka_failures",
		Type:        counter,
		Value:       o.FailuresUnknownPartition,
		Description: failuresDescription,
		LabelName:   "type",
		LabelValue:  "unknown_partition",
	}

	points[12] = &point{
		Name:        "omkafka_failures",
		Type:        counter,
		Value:       o.FailuresOther,
		Description: failuresDescription,
		LabelName:   "type",
		LabelValue:  "other",
	}

	points[13] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsTimedOut,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "timed_out",
	}

	points[14] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsTransport,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "transport",
	}

	points[15] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsBrokerDown,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "broker_down",
	}

	points[16] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsAuth,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "auth",
	}

	points[17] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsSSL,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "ssl",
	}

	points[18] = &point{
		Name:        "omkafka_errors",
		Type:        counter,
		Value:       o.ErrorsOther,
		Description: errorsDescription,
		LabelName:   "type",
		LabelValue:  "other",
	}

	points[19] = &point{
		Name:        "omkafka_rtt_avg_usec_acg",
		Type:        gauge,
		Value:       o.RttAvgUsec,
		Description: "broker round trip time in microseconds averaged over all brokers. It is based on the statistics callback window specified through statistics.interval.ms parameter to librdkafka. Average exclude brokers with less than 100 microseconds rtt",
	}

	points[20] = &point{
		Name:        "omkafka_throttle_avg_msec_avg",
		Type:        gauge,
		Value:       o.ThrottleAvgMsec,
		Description: "broker throttling time in milliseconds averaged over all brokers. This is also a part of window statistics delivered by librdkakfka. Average excludes brokers with zero throttling time",
	}

	points[21] = &point{
		Name:        "omkafka_int_latency_avg_usec_avg",
		Type:        gauge,
		Value:       o.IntLatencyAvgUsec,
		Description: "internal librdkafka producer queue latency in microseconds averaged other all brokers. This is also part of window statistics and average excludes brokers with zero internal latency",
	}

	return points
}
