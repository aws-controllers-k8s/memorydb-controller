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

package user

import (
	"context"
	"github.com/pkg/errors"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	memorydbutil "github.com/aws-controllers-k8s/memorydb-controller/pkg/util"
)

var (
	resourceStatusActive string = "active"
)

// validateUserNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable and
// to validate if resource update is required.
func (rm *resourceManager) validateUserNeedsUpdate(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {

	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != resourceStatusActive {
		return nil, requeue.NeededAfter(
			errors.New("User cannot be modified as its status is not 'active'."),
			requeue.DefaultRequeueAfterDuration)
	}

	// AccessString might be transformed by the server to a different value when this happens
	// delta would be generated when old access string is supplied. We would need to re-patch the
	// Spec so it is updated with the transformed value.
	annotations := desired.ko.ObjectMeta.GetAnnotations()
	if val, ok := annotations[AnnotationLastRequestedAccessString]; ok && desired.ko.Spec.AccessString != nil {
		if val == *desired.ko.Spec.AccessString && delta.DifferentAt("Spec.AccessString") && len(delta.Differences) == 1 {
			return latest, nil
		}
	}

	return nil, nil
}

// userActive returns true when the status of the given User is set to `active`
func (rm *resourceManager) userActive(
	latest *resource,
) bool {
	latestStatus := latest.ko.Status.Status
	return latestStatus != nil && *latestStatus == resourceStatusActive
}

// getTags gets tags from given ParameterGroup.
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
	tags := memorydbutil.ResourceTagsFromSDKTags(resp.TagList)
	return tags, nil
}

// updateTags updates tags of given ParameterGroup to desired tags.
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
		rlog.Debug("removing tags from user", "tags", toDelete)
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
		rlog.Debug("adding tags to user", "tags", toAdd)
		_, err = rm.sdkapi.TagResourceWithContext(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: arn,
				Tags:        memorydbutil.SDKTagsFromResourceTags(toAdd),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "TagResource", err)
		if err != nil {
			return err
		}
	}

	return nil
}
