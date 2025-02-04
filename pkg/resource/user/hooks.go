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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/memorydb"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/memorydb/types"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	svcutil "github.com/aws-controllers-k8s/memorydb-controller/pkg/util"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
)

var (
	resourceStatusActive string = "active"
)

// validateUserNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable and
// to validate if resource update is required.
func (rm *resourceManager) validateUserNeedsUpdate(
	ctx context.Context,
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

// getEvents gets events from User in service.
func (rm *resourceManager) getEvents(
	ctx context.Context,
	r *resource,
) ([]*svcapitypes.Event, error) {
	maxResults := int32(svcutil.MaxEvents)
	input := &svcsdk.DescribeEventsInput{
		SourceName: r.ko.Spec.Name,
		SourceType: svcsdktypes.SourceTypeUser,
		MaxResults: &maxResults,
	}
	resp, err := rm.sdkapi.DescribeEvents(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeEvents", err)
	if err != nil {
		rm.log.V(1).Info("Error during DescribeEvents-User", "error", err)
		return nil, err
	}
	events := make([]*svcapitypes.Event, len(resp.Events))
	for i, event := range resp.Events {
		events[i] = &svcapitypes.Event{
			Date:       &metav1.Time{Time: *event.Date},
			Message:    event.Message,
			SourceName: event.SourceName,
			SourceType: (*string)(&event.SourceType),
		}
	}
	return events, nil
}

// isUserActive returns true when the status of the given User is set to `active`
func (rm *resourceManager) isUserActive(
	latest *resource,
) bool {
	latestStatus := latest.ko.Status.Status
	return latestStatus != nil && *latestStatus == resourceStatusActive
}

// getTags gets tags from given User.
func (rm *resourceManager) getTags(
	ctx context.Context,
	resourceARN string,
) ([]*svcapitypes.Tag, error) {
	resp, err := rm.sdkapi.ListTags(
		ctx,
		&svcsdk.ListTagsInput{
			ResourceArn: &resourceARN,
		},
	)
	rm.metrics.RecordAPICall("GET", "ListTags", err)
	if err != nil {
		return nil, err
	}
	tags := make([]*svcapitypes.Tag, 0, len(resp.TagList))
	for _, tag := range resp.TagList {
		tags = append(tags, &svcapitypes.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
	return tags, nil
}

// updateTags updates tags of given User to desired tags.
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

	var toDelete []string
	for _, removedElement := range toRemove {
		toDelete = append(toDelete, *removedElement.Key)
	}

	if len(toDelete) > 0 {
		rlog.Debug("removing tags from user", "tags", toDelete)
		_, err = rm.sdkapi.UntagResource(
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
		_, err = rm.sdkapi.TagResource(
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
) []svcsdktypes.Tag {
	tags := make([]svcsdktypes.Tag, len(rTags))
	for i := range rTags {
		tags[i] = svcsdktypes.Tag{
			Key:   rTags[i].Key,
			Value: rTags[i].Value,
		}
	}
	return tags
}

func resourceTagsFromSDKTags(
	sdkTags []svcsdktypes.Tag,
) []svcapitypes.Tag {
	tags := make([]svcapitypes.Tag, len(sdkTags))
	for i := range sdkTags {
		tags[i] = svcapitypes.Tag{
			Key:   sdkTags[i].Key,
			Value: sdkTags[i].Value,
		}
	}
	return tags
}

func (rm *resourceManager) resolveUserPasswords(
	ctx context.Context,
	secretRefs []*ackv1alpha1.SecretKeyReference,
) ([]*string, error) {
	// If no password references
	if len(secretRefs) == 0 {
		return nil, nil
	}

	// This will hold the actual plaintext passwords we read from each Secret
	var out []*string

	for _, ref := range secretRefs {
		// 1) Use the ACK runtime method to fetch the plaintext Secret value
		//    from the Kubernetes Secret indicated by the SecretKeyReference.
		secretValue, err := rm.rr.SecretValueFromReference(ctx, ref)
		if err != nil {
			return nil, err
		}

		// 2) Append that plaintext password to our out slice.
		//    We have to pass a *string to MemoryDB
		out = append(out, &secretValue)
	}

	return out, nil
}
