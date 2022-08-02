// Copyright  The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package tqltraces

import (
	"encoding/hex"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/collector/pdata/pcommon"
	"go.opentelemetry.io/collector/pdata/ptrace"

	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/tql"
	"github.com/open-telemetry/opentelemetry-collector-contrib/pkg/telemetryquerylanguage/tql/tqltest"
)

var (
	traceID  = [16]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	traceID2 = [16]byte{16, 15, 14, 13, 12, 11, 10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
	spanID   = [8]byte{1, 2, 3, 4, 5, 6, 7, 8}
	spanID2  = [8]byte{8, 7, 6, 5, 4, 3, 2, 1}
)

func Test_newPathGetSetter(t *testing.T) {
	refSpan, _, _ := createTelemetry()

	newAttrs := pcommon.NewMap()
	newAttrs.UpsertString("hello", "world")

	newEvents := ptrace.NewSpanEventSlice()
	newEvents.AppendEmpty().SetName("new event")

	newLinks := ptrace.NewSpanLinkSlice()
	newLinks.AppendEmpty().SetSpanID(pcommon.NewSpanID(spanID2))

	newStatus := ptrace.NewSpanStatus()
	newStatus.SetMessage("new status")

	newArrStr := pcommon.NewValueSlice()
	newArrStr.SliceVal().AppendEmpty().SetStringVal("new")

	newArrBool := pcommon.NewValueSlice()
	newArrBool.SliceVal().AppendEmpty().SetBoolVal(false)

	newArrInt := pcommon.NewValueSlice()
	newArrInt.SliceVal().AppendEmpty().SetIntVal(20)

	newArrFloat := pcommon.NewValueSlice()
	newArrFloat.SliceVal().AppendEmpty().SetDoubleVal(2.0)

	newArrBytes := pcommon.NewValueSlice()
	newArrBytes.SliceVal().AppendEmpty().SetBytesVal(pcommon.NewImmutableByteSlice([]byte{9, 6, 4}))

	tests := []struct {
		name     string
		path     []tql.Field
		orig     interface{}
		newVal   interface{}
		modified func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource)
	}{
		{
			name: "trace_id",
			path: []tql.Field{
				{
					Name: "trace_id",
				},
			},
			orig:   pcommon.NewTraceID(traceID),
			newVal: pcommon.NewTraceID(traceID2),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetTraceID(pcommon.NewTraceID(traceID2))
			},
		},
		{
			name: "span_id",
			path: []tql.Field{
				{
					Name: "span_id",
				},
			},
			orig:   pcommon.NewSpanID(spanID),
			newVal: pcommon.NewSpanID(spanID2),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetSpanID(pcommon.NewSpanID(spanID2))
			},
		},
		{
			name: "trace_id string",
			path: []tql.Field{
				{
					Name: "trace_id",
				},
				{
					Name: "string",
				},
			},
			orig:   hex.EncodeToString(traceID[:]),
			newVal: hex.EncodeToString(traceID2[:]),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetTraceID(pcommon.NewTraceID(traceID2))
			},
		},
		{
			name: "span_id string",
			path: []tql.Field{
				{
					Name: "span_id",
				},
				{
					Name: "string",
				},
			},
			orig:   hex.EncodeToString(spanID[:]),
			newVal: hex.EncodeToString(spanID2[:]),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetSpanID(pcommon.NewSpanID(spanID2))
			},
		},
		{
			name: "trace_state",
			path: []tql.Field{
				{
					Name: "trace_state",
				},
			},
			orig:   "key1=val1,key2=val2",
			newVal: "key=newVal",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetTraceState("key=newVal")
			},
		},
		{
			name: "trace_state key",
			path: []tql.Field{
				{
					Name:   "trace_state",
					MapKey: tqltest.Strp("key1"),
				},
			},
			orig:   "val1",
			newVal: "newVal",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetTraceState("key1=newVal,key2=val2")
			},
		},
		{
			name: "parent_span_id",
			path: []tql.Field{
				{
					Name: "parent_span_id",
				},
			},
			orig:   pcommon.NewSpanID(spanID2),
			newVal: pcommon.NewSpanID(spanID),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetParentSpanID(pcommon.NewSpanID(spanID))
			},
		},
		{
			name: "name",
			path: []tql.Field{
				{
					Name: "name",
				},
			},
			orig:   "bear",
			newVal: "cat",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetName("cat")
			},
		},
		{
			name: "kind",
			path: []tql.Field{
				{
					Name: "kind",
				},
			},
			orig:   int64(2),
			newVal: int64(3),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetKind(ptrace.SpanKindClient)
			},
		},
		{
			name: "start_time_unix_nano",
			path: []tql.Field{
				{
					Name: "start_time_unix_nano",
				},
			},
			orig:   int64(100_000_000),
			newVal: int64(200_000_000),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetStartTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(200)))
			},
		},
		{
			name: "end_time_unix_nano",
			path: []tql.Field{
				{
					Name: "end_time_unix_nano",
				},
			},
			orig:   int64(500_000_000),
			newVal: int64(200_000_000),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetEndTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(200)))
			},
		},
		{
			name: "attributes",
			path: []tql.Field{
				{
					Name: "attributes",
				},
			},
			orig:   refSpan.Attributes(),
			newVal: newAttrs,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Clear()
				newAttrs.CopyTo(span.Attributes())
			},
		},
		{
			name: "attributes string",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("str"),
				},
			},
			orig:   "val",
			newVal: "newVal",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().UpsertString("str", "newVal")
			},
		},
		{
			name: "attributes bool",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("bool"),
				},
			},
			orig:   true,
			newVal: false,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().UpsertBool("bool", false)
			},
		},
		{
			name: "attributes int",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("int"),
				},
			},
			orig:   int64(10),
			newVal: int64(20),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().UpsertInt("int", 20)
			},
		},
		{
			name: "attributes float",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("double"),
				},
			},
			orig:   float64(1.2),
			newVal: float64(2.4),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().UpsertDouble("double", 2.4)
			},
		},
		{
			name: "attributes bytes",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("bytes"),
				},
			},
			orig:   []byte{1, 3, 2},
			newVal: []byte{2, 3, 4},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().UpsertBytes("bytes", pcommon.NewImmutableByteSlice([]byte{2, 3, 4}))
			},
		},
		{
			name: "attributes array string",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_str"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_str")
				return val.SliceVal()
			}(),
			newVal: []string{"new"},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Upsert("arr_str", newArrStr)
			},
		},
		{
			name: "attributes array bool",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_bool"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_bool")
				return val.SliceVal()
			}(),
			newVal: []bool{false},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Upsert("arr_bool", newArrBool)
			},
		},
		{
			name: "attributes array int",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_int"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_int")
				return val.SliceVal()
			}(),
			newVal: []int64{20},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Upsert("arr_int", newArrInt)
			},
		},
		{
			name: "attributes array float",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_float"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_float")
				return val.SliceVal()
			}(),
			newVal: []float64{2.0},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Upsert("arr_float", newArrFloat)
			},
		},
		{
			name: "attributes array bytes",
			path: []tql.Field{
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_bytes"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_bytes")
				return val.SliceVal()
			}(),
			newVal: [][]byte{{9, 6, 4}},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Attributes().Upsert("arr_bytes", newArrBytes)
			},
		},
		{
			name: "dropped_attributes_count",
			path: []tql.Field{
				{
					Name: "dropped_attributes_count",
				},
			},
			orig:   int64(10),
			newVal: int64(20),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetDroppedAttributesCount(20)
			},
		},
		{
			name: "events",
			path: []tql.Field{
				{
					Name: "events",
				},
			},
			orig:   refSpan.Events(),
			newVal: newEvents,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Events().RemoveIf(func(_ ptrace.SpanEvent) bool {
					return true
				})
				newEvents.CopyTo(span.Events())
			},
		},
		{
			name: "dropped_events_count",
			path: []tql.Field{
				{
					Name: "dropped_events_count",
				},
			},
			orig:   int64(20),
			newVal: int64(30),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetDroppedEventsCount(30)
			},
		},
		{
			name: "links",
			path: []tql.Field{
				{
					Name: "links",
				},
			},
			orig:   refSpan.Links(),
			newVal: newLinks,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Links().RemoveIf(func(_ ptrace.SpanLink) bool {
					return true
				})
				newLinks.CopyTo(span.Links())
			},
		},
		{
			name: "dropped_links_count",
			path: []tql.Field{
				{
					Name: "dropped_links_count",
				},
			},
			orig:   int64(30),
			newVal: int64(40),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.SetDroppedLinksCount(40)
			},
		},
		{
			name: "status",
			path: []tql.Field{
				{
					Name: "status",
				},
			},
			orig:   refSpan.Status(),
			newVal: newStatus,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				newStatus.CopyTo(span.Status())
			},
		},
		{
			name: "status code",
			path: []tql.Field{
				{
					Name: "status",
				},
				{
					Name: "code",
				},
			},
			orig:   int64(ptrace.StatusCodeOk),
			newVal: int64(ptrace.StatusCodeError),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Status().SetCode(ptrace.StatusCodeError)
			},
		},
		{
			name: "status message",
			path: []tql.Field{
				{
					Name: "status",
				},
				{
					Name: "message",
				},
			},
			orig:   "good span",
			newVal: "bad span",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				span.Status().SetMessage("bad span")
			},
		},
		{
			name: "resource attributes",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name: "attributes",
				},
			},
			orig:   refSpan.Attributes(),
			newVal: newAttrs,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Clear()
				newAttrs.CopyTo(resource.Attributes())
			},
		},
		{
			name: "resource attributes string",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("str"),
				},
			},
			orig:   "val",
			newVal: "newVal",
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().UpsertString("str", "newVal")
			},
		},
		{
			name: "resource attributes bool",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("bool"),
				},
			},
			orig:   true,
			newVal: false,
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().UpsertBool("bool", false)
			},
		},
		{
			name: "resource attributes int",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("int"),
				},
			},
			orig:   int64(10),
			newVal: int64(20),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().UpsertInt("int", 20)
			},
		},
		{
			name: "resource attributes float",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("double"),
				},
			},
			orig:   float64(1.2),
			newVal: float64(2.4),
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().UpsertDouble("double", 2.4)
			},
		},
		{
			name: "resource attributes bytes",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("bytes"),
				},
			},
			orig:   []byte{1, 3, 2},
			newVal: []byte{2, 3, 4},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().UpsertBytes("bytes", pcommon.NewImmutableByteSlice([]byte{2, 3, 4}))
			},
		},
		{
			name: "resource attributes array string",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_str"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_str")
				return val.SliceVal()
			}(),
			newVal: []string{"new"},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Upsert("arr_str", newArrStr)
			},
		},
		{
			name: "resource attributes array bool",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_bool"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_bool")
				return val.SliceVal()
			}(),
			newVal: []bool{false},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Upsert("arr_bool", newArrBool)
			},
		},
		{
			name: "resource attributes array int",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_int"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_int")
				return val.SliceVal()
			}(),
			newVal: []int64{20},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Upsert("arr_int", newArrInt)
			},
		},
		{
			name: "resource attributes array float",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_float"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_float")
				return val.SliceVal()
			}(),
			newVal: []float64{2.0},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Upsert("arr_float", newArrFloat)
			},
		},
		{
			name: "resource attributes array bytes",
			path: []tql.Field{
				{
					Name: "resource",
				},
				{
					Name:   "attributes",
					MapKey: tqltest.Strp("arr_bytes"),
				},
			},
			orig: func() pcommon.Slice {
				val, _ := refSpan.Attributes().Get("arr_bytes")
				return val.SliceVal()
			}(),
			newVal: [][]byte{{9, 6, 4}},
			modified: func(span ptrace.Span, il pcommon.InstrumentationScope, resource pcommon.Resource) {
				resource.Attributes().Upsert("arr_bytes", newArrBytes)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			accessor, err := newPathGetSetter(tt.path)
			assert.NoError(t, err)

			span, il, resource := createTelemetry()

			got := accessor.Get(SpanTransformContext{
				Span:                 span,
				InstrumentationScope: il,
				Resource:             resource,
			})
			assert.Equal(t, tt.orig, got)

			accessor.Set(SpanTransformContext{
				Span:                 span,
				InstrumentationScope: il,
				Resource:             resource,
			}, tt.newVal)

			exSpan, exIl, exRes := createTelemetry()
			tt.modified(exSpan, exIl, exRes)

			assert.Equal(t, exSpan, span)
			assert.Equal(t, exIl, il)
			assert.Equal(t, exRes, resource)
		})
	}
}

func createTelemetry() (ptrace.Span, pcommon.InstrumentationScope, pcommon.Resource) {
	span := ptrace.NewSpan()
	span.SetTraceID(pcommon.NewTraceID(traceID))
	span.SetSpanID(pcommon.NewSpanID(spanID))
	span.SetTraceState("key1=val1,key2=val2")
	span.SetParentSpanID(pcommon.NewSpanID(spanID2))
	span.SetName("bear")
	span.SetKind(ptrace.SpanKindServer)
	span.SetStartTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(100)))
	span.SetEndTimestamp(pcommon.NewTimestampFromTime(time.UnixMilli(500)))
	span.Attributes().UpsertString("str", "val")
	span.Attributes().UpsertBool("bool", true)
	span.Attributes().UpsertInt("int", 10)
	span.Attributes().UpsertDouble("double", 1.2)
	span.Attributes().UpsertBytes("bytes", pcommon.NewImmutableByteSlice([]byte{1, 3, 2}))

	arrStr := pcommon.NewValueSlice()
	arrStr.SliceVal().AppendEmpty().SetStringVal("one")
	arrStr.SliceVal().AppendEmpty().SetStringVal("two")
	span.Attributes().Upsert("arr_str", arrStr)

	arrBool := pcommon.NewValueSlice()
	arrBool.SliceVal().AppendEmpty().SetBoolVal(true)
	arrBool.SliceVal().AppendEmpty().SetBoolVal(false)
	span.Attributes().Upsert("arr_bool", arrBool)

	arrInt := pcommon.NewValueSlice()
	arrInt.SliceVal().AppendEmpty().SetIntVal(2)
	arrInt.SliceVal().AppendEmpty().SetIntVal(3)
	span.Attributes().Upsert("arr_int", arrInt)

	arrFloat := pcommon.NewValueSlice()
	arrFloat.SliceVal().AppendEmpty().SetDoubleVal(1.0)
	arrFloat.SliceVal().AppendEmpty().SetDoubleVal(2.0)
	span.Attributes().Upsert("arr_float", arrFloat)

	arrBytes := pcommon.NewValueSlice()
	arrBytes.SliceVal().AppendEmpty().SetBytesVal(pcommon.NewImmutableByteSlice([]byte{1, 2, 3}))
	arrBytes.SliceVal().AppendEmpty().SetBytesVal(pcommon.NewImmutableByteSlice([]byte{2, 3, 4}))
	span.Attributes().Upsert("arr_bytes", arrBytes)

	span.SetDroppedAttributesCount(10)

	span.Events().AppendEmpty().SetName("event")
	span.SetDroppedEventsCount(20)

	span.Links().AppendEmpty().SetTraceID(pcommon.NewTraceID(traceID))
	span.SetDroppedLinksCount(30)

	span.Status().SetCode(ptrace.StatusCodeOk)
	span.Status().SetMessage("good span")

	il := pcommon.NewInstrumentationScope()
	il.SetName("library")
	il.SetVersion("version")

	resource := pcommon.NewResource()
	span.Attributes().CopyTo(resource.Attributes())

	return span, il, resource
}

func Test_ParseEnum(t *testing.T) {
	tests := []struct {
		name string
		want tql.Enum
	}{
		{
			name: "SPAN_KIND_UNSPECIFIED",
			want: tql.Enum(ptrace.SpanKindUnspecified),
		},
		{
			name: "SPAN_KIND_INTERNAL",
			want: tql.Enum(ptrace.SpanKindInternal),
		},
		{
			name: "SPAN_KIND_SERVER",
			want: tql.Enum(ptrace.SpanKindServer),
		},
		{
			name: "SPAN_KIND_CLIENT",
			want: tql.Enum(ptrace.SpanKindClient),
		},
		{
			name: "SPAN_KIND_PRODUCER",
			want: tql.Enum(ptrace.SpanKindProducer),
		},
		{
			name: "SPAN_KIND_CONSUMER",
			want: tql.Enum(ptrace.SpanKindConsumer),
		},
		{
			name: "STATUS_CODE_UNSET",
			want: tql.Enum(ptrace.StatusCodeUnset),
		},
		{
			name: "STATUS_CODE_OK",
			want: tql.Enum(ptrace.StatusCodeOk),
		},
		{
			name: "STATUS_CODE_ERROR",
			want: tql.Enum(ptrace.StatusCodeError),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseEnum((*tql.EnumSymbol)(tqltest.Strp(tt.name)))
			assert.NoError(t, err)
			assert.Equal(t, *actual, tt.want)
		})
	}
}

func Test_ParseEnum_False(t *testing.T) {
	tests := []struct {
		name       string
		enumSymbol *tql.EnumSymbol
	}{
		{
			name:       "unknown enum symbol",
			enumSymbol: (*tql.EnumSymbol)(tqltest.Strp("not an enum")),
		},
		{
			name:       "nil enum symbol",
			enumSymbol: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, err := ParseEnum(tt.enumSymbol)
			assert.Error(t, err)
			assert.Nil(t, actual)
		})
	}
}
