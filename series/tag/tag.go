// Licensed to LinDB under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. LinDB licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package tag

import (
	"sort"
	"strings"

	protoMetricsV1 "github.com/lindb/lindb/proto/gen/v1/metrics"
)

type KeyValues []*protoMetricsV1.KeyValue

func (kvs KeyValues) Len() int           { return len(kvs) }
func (kvs KeyValues) Less(i, j int) bool { return kvs[i].Key < kvs[j].Key }
func (kvs KeyValues) Swap(i, j int)      { kvs[i], kvs[j] = kvs[j], kvs[i] }

// DeDup sorts keyvalues and removes the duplicates
func (kvs KeyValues) DeDup() KeyValues {
	if len(kvs) < 2 {
		return kvs
	}
	sort.Sort(kvs)
	var (
		fast = 1
		slow = 0
	)
	for fast < kvs.Len() {
		// move to next
		if kvs[fast].Key != kvs[slow].Key {
			slow++
		}
		kvs[slow] = kvs[fast]
		fast++
	}
	return kvs[:slow+1]
}

// Map transforms the KeyValues into map
func (kvs KeyValues) Map() map[string]string {
	var m = make(map[string]string)
	for idx := range kvs {
		m[kvs[idx].Key] = kvs[idx].Value
	}
	return m
}

// Clone returns a copy of keyValues
func (kvs KeyValues) Clone() KeyValues {
	var dst = make([]*protoMetricsV1.KeyValue, len(kvs))
	for i := range kvs {
		dst[i] = &protoMetricsV1.KeyValue{
			Key:   kvs[i].Key,
			Value: kvs[i].Value,
		}
	}
	return dst
}

// Merge merges another keyvalue list into a new one
func (kvs KeyValues) Merge(other KeyValues) KeyValues {
	if len(other) == 0 {
		return kvs.Clone()
	}
	m := kvs.Map()
	for _, item := range other {
		m[item.Key] = item.Value
	}
	merged := make(KeyValues, len(m))
	idx := 0
	for key, value := range m {
		merged[idx] = &protoMetricsV1.KeyValue{
			Key:   key,
			Value: value,
		}
		idx++
	}
	sort.Sort(merged)
	return merged
}

func KeyValuesFromMap(tags map[string]string) KeyValues {
	if tags == nil {
		return nil
	}
	var kvs KeyValues
	for k, v := range tags {
		kvs = append(kvs, &protoMetricsV1.KeyValue{Key: k, Value: v})
	}
	return kvs
}

func ConcatKeyValues(kvs KeyValues) string {
	if len(kvs) == 0 {
		return ""
	}
	sort.Sort(kvs)
	tagKeysLen := len(kvs)
	var b strings.Builder
	b.Grow(128)
	for idx := range kvs {
		b.WriteString(kvs[idx].Key)
		b.WriteString("=")
		b.WriteString(kvs[idx].Value)
		if idx != tagKeysLen-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

// Concat concats map-tags to string
func Concat(tags map[string]string) string {
	if tags == nil {
		return ""
	}
	tagKeys := make([]string, 0, len(tags))
	var b strings.Builder
	b.Grow(128)
	for key := range tags {
		tagKeys = append(tagKeys, key)
	}
	sort.Strings(tagKeys)
	tagKeysLen := len(tagKeys)
	for idx, tagKey := range tagKeys {
		b.WriteString(tagKey)
		b.WriteString("=")
		b.WriteString(tags[tagKey])
		if idx != tagKeysLen-1 {
			b.WriteString(",")
		}
	}
	return b.String()
}

// ConcatTagValues cancats the tag values to string
func ConcatTagValues(tagValues []string) string {
	if len(tagValues) == 0 {
		return ""
	}
	return strings.Join(tagValues, ",")
}

// SplitTagValues splits the string of tag values to array
func SplitTagValues(tags string) []string {
	if tags == "" {
		return []string{}
	}
	return strings.Split(tags, ",")
}
