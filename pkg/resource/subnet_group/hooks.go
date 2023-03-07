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

package subnet_group

import (
	"context"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

// getTags gets tags from given SubnetGroup.
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

// updateTags updates tags of given SubnetGroup to desired tags.
func (rm *resourceManager) updateTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()

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
		rlog.Debug("removing tags from parameter group", "tags", toDelete)
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
		rlog.Debug("adding tags to parameter group", "tags", toAdd)
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
