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

package acl

import (
	"context"
	"github.com/pkg/errors"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	svcutil "github.com/aws-controllers-k8s/memorydb-controller/pkg/util"
)

var (
	resourceStatusActive string = "active"
)

// validateACLNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable
func (rm *resourceManager) validateACLNeedsUpdate(
	latest *resource,
) error {

	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != resourceStatusActive {
		return requeue.NeededAfter(
			errors.New("ACL cannot be modified as its status is not 'active'."),
			requeue.DefaultRequeueAfterDuration)
	}

	return nil
}

// getEvents gets events from ACL in service.
func (rm *resourceManager) getEvents(
	ctx context.Context,
	r *resource,
) ([]*svcapitypes.Event, error) {
	input := svcutil.NewDescribeEventsInput(*r.ko.Spec.Name, svcsdk.SourceTypeAcl, svcutil.MaxEvents)
	resp, err := rm.sdkapi.DescribeEventsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeEvents", err)
	if err != nil {
		rm.log.V(1).Info("Error during DescribeEvents-ACL", "error", err)
		return nil, err
	}
	return svcutil.EventsFromDescribe(resp), nil
}

// isACLActive returns true when the status of the given ACL is set to `active`
func (rm *resourceManager) isACLActive(
	latest *resource,
) bool {
	latestStatus := latest.ko.Status.Status
	return latestStatus != nil && *latestStatus == resourceStatusActive
}

// getTags gets tags from given ACL.
func (rm *resourceManager) getTags(
	ctx context.Context,
	resourceARN string,
) ([]*svcapitypes.Tag, error) {
	resp, err := rm.sdkapi.ListTagsWithContext(
		ctx,
		&svcsdk.ListTagsInput{
			ResourceArn: &resourceARN,
		},
	)
	rm.metrics.RecordAPICall("GET", "ListTags", err)
	if err != nil {
		return nil, err
	}
	tags := resourceTagsFromSDKTags(resp.TagList)
	return tags, nil
}

// updateTags updates tags of given ACL to desired tags.
func (rm *resourceManager) updateTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.updateTags")
	defer func(err error) { exit(err) }(err)

	arn := (*string)(latest.ko.Status.ACKResourceMetadata.ARN)

	desiredTags := ToACKTags(desired.ko.Spec.Tags)
	latestTags := ToACKTags(latest.ko.Spec.Tags)

	added, _, removed := ackcompare.GetTagsDifference(latestTags, desiredTags)

	toAdd := FromACKTags(added)
	toRemove := FromACKTags(removed)

	var toDelete []*string
	for _, removedElement := range toRemove {
		toDelete = append(toDelete, removedElement.Key)
	}

	if len(toDelete) > 0 {
		rlog.Debug("removing tags from ACL", "tags", toDelete)
		_, err = rm.sdkapi.UntagResourceWithContext(
			ctx,
			&svcsdk.UntagResourceInput{
				ResourceArn: arn,
				TagKeys:     toDelete,
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "UntagResource", err)
		if err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		rlog.Debug("adding tags to ACL", "tags", toAdd)
		_, err = rm.sdkapi.TagResourceWithContext(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: arn,
				Tags:        sdkTagsFromResourceTags(toAdd),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "TagResource", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func sdkTagsFromResourceTags(
	rTags []*svcapitypes.Tag,
) []*svcsdk.Tag {
	tags := make([]*svcsdk.Tag, len(rTags))
	for i := range rTags {
		tags[i] = &svcsdk.Tag{
			Key:   rTags[i].Key,
			Value: rTags[i].Value,
		}
	}
	return tags
}

func resourceTagsFromSDKTags(
	sdkTags []*svcsdk.Tag,
) []*svcapitypes.Tag {
	tags := make([]*svcapitypes.Tag, len(sdkTags))
	for i := range sdkTags {
		tags[i] = &svcapitypes.Tag{
			Key:   sdkTags[i].Key,
			Value: sdkTags[i].Value,
		}
	}
	return tags
}
