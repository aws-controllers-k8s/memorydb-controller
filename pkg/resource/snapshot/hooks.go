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

package snapshot

import (
	"context"
	"errors"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

var (
	resourceStatusAvailable string = "available"
)

func (rm *resourceManager) customDescribeSnapshotSetOutput(
	resp *svcsdk.DescribeSnapshotsOutput,
	ko *svcapitypes.Snapshot,
) (*svcapitypes.Snapshot, error) {
	if len(resp.Snapshots) == 0 {
		return ko, nil
	}
	elem := resp.Snapshots[0]
	rm.customSetOutput(elem, ko)
	return ko, nil
}

func (rm *resourceManager) customCreateSnapshotSetOutput(
	resp *svcsdk.CreateSnapshotOutput,
	ko *svcapitypes.Snapshot,
) (*svcapitypes.Snapshot, error) {
	rm.customSetOutput(resp.Snapshot, ko)
	return ko, nil
}

func (rm *resourceManager) customCopySnapshotSetOutput(
	resp *svcsdk.CopySnapshotOutput,
	ko *svcapitypes.Snapshot,
) *svcapitypes.Snapshot {
	rm.customSetOutput(resp.Snapshot, ko)
	return ko
}

func (rm *resourceManager) customSetOutput(
	respSnapshot *svcsdk.Snapshot,
	ko *svcapitypes.Snapshot,
) {
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
	snapshotStatus := respSnapshot.Status
	syncConditionStatus := corev1.ConditionUnknown
	if snapshotStatus != nil {
		if *snapshotStatus == "available" ||
			*snapshotStatus == "failed" {
			syncConditionStatus = corev1.ConditionTrue
		} else {
			// resource in "creating", "restoring","exporting"
			syncConditionStatus = corev1.ConditionFalse
		}
	}
	var resourceSyncedCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			resourceSyncedCondition = condition
			break
		}
	}
	if resourceSyncedCondition == nil {
		resourceSyncedCondition = &ackv1alpha1.Condition{
			Type:   ackv1alpha1.ConditionTypeResourceSynced,
			Status: syncConditionStatus,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	} else {
		resourceSyncedCondition.Status = syncConditionStatus
	}
}

func (rm *resourceManager) customTryCopySnapshot(
	ctx context.Context,
	r *resource,
) (*resource, error) {
	if r.ko.Spec.SourceSnapshotName == nil {
		return nil, nil
	}
	if r.ko.Spec.ClusterName != nil {
		return nil, ackerr.NewTerminalError(errors.New("Cannot specify ClusterName when SourceSnapshotName is specified"))
	}

	input, err := rm.newCopySnapshotPayload(r)
	if err != nil {
		return nil, err
	}

	resp, respErr := rm.sdkapi.CopySnapshot(input)

	rm.metrics.RecordAPICall("CREATE", "CopySnapshot", respErr)
	if respErr != nil {
		return nil, respErr
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Snapshot.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Snapshot.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}

	if resp.Snapshot.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.Snapshot.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}

	if resp.Snapshot.ClusterConfiguration != nil {
		f1 := &svcapitypes.ClusterConfiguration{}
		if resp.Snapshot.ClusterConfiguration.Description != nil {
			f1.Description = resp.Snapshot.ClusterConfiguration.Description
		}
		if resp.Snapshot.ClusterConfiguration.EngineVersion != nil {
			f1.EngineVersion = resp.Snapshot.ClusterConfiguration.EngineVersion
		}
		if resp.Snapshot.ClusterConfiguration.MaintenanceWindow != nil {
			f1.MaintenanceWindow = resp.Snapshot.ClusterConfiguration.MaintenanceWindow
		}
		if resp.Snapshot.ClusterConfiguration.Name != nil {
			f1.Name = resp.Snapshot.ClusterConfiguration.Name
		}
		if resp.Snapshot.ClusterConfiguration.NodeType != nil {
			f1.NodeType = resp.Snapshot.ClusterConfiguration.NodeType
		}
		if resp.Snapshot.ClusterConfiguration.NumShards != nil {
			f1.NumShards = resp.Snapshot.ClusterConfiguration.NumShards
		}
		if resp.Snapshot.ClusterConfiguration.ParameterGroupName != nil {
			f1.ParameterGroupName = resp.Snapshot.ClusterConfiguration.ParameterGroupName
		}
		if resp.Snapshot.ClusterConfiguration.Port != nil {
			f1.Port = resp.Snapshot.ClusterConfiguration.Port
		}
		if resp.Snapshot.ClusterConfiguration.Shards != nil {
			f1f8 := []*svcapitypes.ShardDetail{}
			for _, f1f8iter := range resp.Snapshot.ClusterConfiguration.Shards {
				f1f8elem := &svcapitypes.ShardDetail{}
				if f1f8iter.Configuration != nil {
					f1f8elemf0 := &svcapitypes.ShardConfiguration{}
					if f1f8iter.Configuration.ReplicaCount != nil {
						f1f8elemf0.ReplicaCount = f1f8iter.Configuration.ReplicaCount
					}
					if f1f8iter.Configuration.Slots != nil {
						f1f8elemf0.Slots = f1f8iter.Configuration.Slots
					}
					f1f8elem.Configuration = f1f8elemf0
				}
				if f1f8iter.Name != nil {
					f1f8elem.Name = f1f8iter.Name
				}
				if f1f8iter.Size != nil {
					f1f8elem.Size = f1f8iter.Size
				}
				if f1f8iter.SnapshotCreationTime != nil {
					f1f8elem.SnapshotCreationTime = &metav1.Time{*f1f8iter.SnapshotCreationTime}
				}
				f1f8 = append(f1f8, f1f8elem)
			}
			f1.Shards = f1f8
		}
		if resp.Snapshot.ClusterConfiguration.SnapshotRetentionLimit != nil {
			f1.SnapshotRetentionLimit = resp.Snapshot.ClusterConfiguration.SnapshotRetentionLimit
		}
		if resp.Snapshot.ClusterConfiguration.SnapshotWindow != nil {
			f1.SnapshotWindow = resp.Snapshot.ClusterConfiguration.SnapshotWindow
		}
		if resp.Snapshot.ClusterConfiguration.SubnetGroupName != nil {
			f1.SubnetGroupName = resp.Snapshot.ClusterConfiguration.SubnetGroupName
		}
		if resp.Snapshot.ClusterConfiguration.TopicArn != nil {
			f1.TopicARN = resp.Snapshot.ClusterConfiguration.TopicArn
		}
		if resp.Snapshot.ClusterConfiguration.VpcId != nil {
			f1.VPCID = resp.Snapshot.ClusterConfiguration.VpcId
		}
		ko.Status.ClusterConfiguration = f1
	} else {
		ko.Status.ClusterConfiguration = nil
	}
	if resp.Snapshot.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.Snapshot.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.Snapshot.Source != nil {
		ko.Status.Source = resp.Snapshot.Source
	} else {
		ko.Status.Source = nil
	}
	if resp.Snapshot.Status != nil {
		ko.Status.Status = resp.Snapshot.Status
	} else {
		ko.Status.Status = nil
	}

	rm.setStatusDefaults(ko)
	// custom set output from response
	rm.customCopySnapshotSetOutput(resp, ko)
	return &resource{ko}, nil
}

// newCopySnapshotPayload returns an SDK-specific struct for the HTTP request
// payload of the CopySnapshot API call
func (rm *resourceManager) newCopySnapshotPayload(
	r *resource,
) (*svcsdk.CopySnapshotInput, error) {
	res := &svcsdk.CopySnapshotInput{}

	if r.ko.Spec.SourceSnapshotName != nil {
		res.SetSourceSnapshotName(*r.ko.Spec.SourceSnapshotName)
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.SetKmsKeyId(*r.ko.Spec.KMSKeyID)
	}

	if r.ko.Spec.Name != nil {
		res.SetTargetSnapshotName(*r.ko.Spec.Name)
	}

	return res, nil
}

// isSnapshotAvailable returns true when the status of the given Snapshot is set to `available`
func (rm *resourceManager) isSnapshotAvailable(
	latest *resource,
) bool {
	latestStatus := latest.ko.Status.Status
	return latestStatus != nil && *latestStatus == resourceStatusAvailable
}

// getTags gets tags from given Snapshot.
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

func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func(err error) { exit(err) }(err)
	if delta.DifferentAt("Spec.Tags") {
		if err = rm.updateTags(ctx, desired, latest); err != nil {
			return nil, err
		}
	}
	return desired, nil
}

// updateTags updates tags of given Snapshot to desired tags.
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
		rlog.Debug("removing tags from snapshot", "tags", toDelete)
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
		rlog.Debug("adding tags to snapshot", "tags", toAdd)
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
