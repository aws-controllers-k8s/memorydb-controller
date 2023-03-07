// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License"). You may
// not use this file except in compliance with the License. A copy of the
// License is located at
//
//     http://aws.amazon.com/apache2.0/
//
// or in the "license" file accompanying this file. This file is distributed
// on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either
// express or implied. See the License for the specific language governing
// permissions and limitations under the License.

package util

import (
	"time"

	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/samber/lo"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

var (
	// eventsDuration represents the maximum time range worth of
	//events to retrieve.
	eventsDuration = time.Duration(14*24) * time.Hour
	// MaxEvents represents the maximum events can to retrieve.
	MaxEvents = int64(20)
)

func NewDescribeEventsInput(
	sourceName string,
	sourceType string,
	maxResults int64,
) *svcsdk.DescribeEventsInput {
	input := &svcsdk.DescribeEventsInput{}
	input.SetSourceType(sourceType)
	input.SetSourceName(sourceName)
	input.SetMaxResults(maxResults)
	input.SetDuration(int64(eventsDuration.Minutes()))
	return input
}

// EventsFromDescribe returns a slice of Event structs given the
// Output response shape from a DescribeEventsWithContext call
func EventsFromDescribe(
	resp *svcsdk.DescribeEventsOutput,
) []*svcapitypes.Event {
	events := lo.Map(resp.Events, func(respEvent *svcsdk.Event, index int) *svcapitypes.Event {
		event := &svcapitypes.Event{
			Message: respEvent.Message,
		}
		if respEvent.Date != nil {
			eventDate := metav1.NewTime(*respEvent.Date)
			event.Date = &eventDate
		}
		return event
	})

	return events
}
