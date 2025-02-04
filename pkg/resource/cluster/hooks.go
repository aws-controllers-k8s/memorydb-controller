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

package cluster

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/memorydb"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/memorydb/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	svcutil "github.com/aws-controllers-k8s/memorydb-controller/pkg/util"
)

var (
	condMsgCurrentlyDeleting            = "cluster is currently being deleted"
	condMsgNoDeleteWhileUpdating        = "cluster is currently being updated. cannot delete"
	resourceStatusAvailable      string = "available"
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("delete is in progress"),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileUpdating = ackrequeue.NeededAfter(
		errors.New("update is in progress"),
		ackrequeue.DefaultRequeueAfterDuration*2,
	)
)

const (
	StatusAvailable    = "available"
	StatusDeleting     = "deleting"
	StatusUpdating     = "updating"
	StatusCreateFailed = "create-failed"
)

// isDeleting returns true if supplied cluster resource state is 'deleting'
func isDeleting(r *resource) bool {
	if r == nil || r.ko.Status.Status == nil {
		return false
	}
	status := *r.ko.Status.Status
	return status == StatusDeleting
}

// isUpdating returns true if supplied cluster resource state is 'modifying'
func isUpdating(r *resource) bool {
	if r == nil || r.ko.Status.Status == nil {
		return false
	}
	status := *r.ko.Status.Status
	return status == StatusUpdating
}

func (rm *resourceManager) setShardDetails(
	ctx context.Context,
	r *resource,
	ko *svcapitypes.Cluster,
) (*svcapitypes.Cluster, error) {

	resp, err := rm.sdkFind(ctx, r)

	if err != nil {
		return nil, err
	}

	ko.Status = resp.ko.Status
	ko.Spec.NumReplicasPerShard = resp.ko.Spec.NumReplicasPerShard
	ko.Spec.NumShards = resp.ko.Spec.NumShards

	return ko, nil
}

func (rm *resourceManager) setAllowedNodeTypeUpdates(
	ctx context.Context,
	ko *svcapitypes.Cluster,
) error {
	if *ko.Status.Status != StatusAvailable {
		return nil
	}

	response, respErr := rm.sdkapi.ListAllowedNodeTypeUpdates(ctx, &svcsdk.ListAllowedNodeTypeUpdatesInput{
		ClusterName: ko.Spec.Name,
	})
	rm.metrics.RecordAPICall("GET", "ListAllowedNodeTypeUpdates", respErr)
	if respErr == nil {
		ko.Status.AllowedScaleDownNodeTypes = aws.StringSlice(response.ScaleDownNodeTypes)
		ko.Status.AllowedScaleUpNodeTypes = aws.StringSlice(response.ScaleUpNodeTypes)
	}

	return respErr
}

// validateClusterNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable and
// to validate if resource update is required.
func (rm *resourceManager) validateClusterNeedsUpdate(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != StatusAvailable {
		return nil, ackrequeue.NeededAfter(
			fmt.Errorf("cluster cannot be updated as its status is not '%s'", StatusAvailable),
			ackrequeue.DefaultRequeueAfterDuration)
	}

	// Set terminal condition when cluster is in create-failed state
	if *latestStatus == StatusCreateFailed {
		return nil, ackerr.NewTerminalError(fmt.Errorf("cluster is in '%s' state, cannot be updated", StatusCreateFailed))
	}

	annotations := desired.ko.ObjectMeta.GetAnnotations()

	// Handle asynchronous rollback of NodeType. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNodeType]; ok && desired.ko.Spec.NodeType != nil {
		if val == *desired.ko.Spec.NodeType && delta.DifferentAt("Spec.NodeType") {
			return nil, ackerr.NewTerminalError(errors.New("cannot update NodeType"))
		}
	}

	// Handle asynchronous rollback of NumShards. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNumShards]; ok && desired.ko.Spec.NumShards != nil {
		numShards, err := strconv.ParseInt(val, 10, 64)
		if err == nil && numShards == *desired.ko.Spec.NumShards && delta.DifferentAt("Spec.NumShards") {
			return nil, ackerr.NewTerminalError(errors.New("cannot update NumShards"))
		}
	}

	// Handle asynchronous rollback of NumShards. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNumReplicasPerShard]; ok && desired.ko.Spec.NumReplicasPerShard != nil {
		numReplicasPerShard, err := strconv.ParseInt(val, 10, 64)
		if err == nil && numReplicasPerShard == *desired.ko.Spec.NumReplicasPerShard && delta.DifferentAt("Spec.NumReplicasPerShard") {
			return nil, ackerr.NewTerminalError(errors.New("cannot update NumReplicasPerShard"))
		}
	}

	return nil, nil
}

func (rm *resourceManager) newMemoryDBClusterUploadPayload(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) *svcsdk.UpdateClusterInput {
	res := &svcsdk.UpdateClusterInput{}

	if delta.DifferentAt("Spec.ACLName") && desired.ko.Spec.ACLName != nil {
		res.ACLName = desired.ko.Spec.ACLName
	}
	if desired.ko.Spec.Name != nil {
		res.ClusterName = desired.ko.Spec.Name
	}
	if delta.DifferentAt("Spec.Description") && desired.ko.Spec.Description != nil {
		res.Description = desired.ko.Spec.Description
	}
	if delta.DifferentAt("Spec.MaintenanceWindow") && desired.ko.Spec.MaintenanceWindow != nil {
		res.MaintenanceWindow = desired.ko.Spec.MaintenanceWindow
	}
	if delta.DifferentAt("Spec.ParameterGroupName") && desired.ko.Spec.ParameterGroupName != nil {
		res.ParameterGroupName = desired.ko.Spec.ParameterGroupName
	}
	if delta.DifferentAt("Spec.SecurityGroupIDs") && desired.ko.Spec.SecurityGroupIDs != nil {
		f8 := []string{}
		for _, f8iter := range desired.ko.Spec.SecurityGroupIDs {
			f8 = append(f8, *f8iter)
		}
		res.SecurityGroupIds = f8
	}
	if delta.DifferentAt("Spec.SnapshotRetentionLimit") && desired.ko.Spec.SnapshotRetentionLimit != nil {
		res.SnapshotRetentionLimit = aws.Int32(int32(*desired.ko.Spec.SnapshotRetentionLimit))
	}
	if delta.DifferentAt("Spec.SnapshotWindow") && desired.ko.Spec.SnapshotWindow != nil {
		res.SnapshotWindow = desired.ko.Spec.SnapshotWindow
	}
	if delta.DifferentAt("Spec.SNSTopicARN") && desired.ko.Spec.SNSTopicARN != nil {
		res.SnsTopicArn = desired.ko.Spec.SNSTopicARN
	}
	if delta.DifferentAt("Spec.SNSTopicStatus") && desired.ko.Status.SNSTopicStatus != nil {
		res.SnsTopicStatus = desired.ko.Status.SNSTopicStatus
	}

	if delta.DifferentAt("Spec.EngineVersion") && desired.ko.Spec.EngineVersion != nil {
		res.EngineVersion = desired.ko.Spec.EngineVersion
	}

	// Determine if we are trying to scale up an instance
	scaleUpUpdate := false
	if delta.DifferentAt("Spec.NodeType") && desired.ko.Spec.NodeType != nil {
		for _, instance := range desired.ko.Status.AllowedScaleUpNodeTypes {
			if *instance == *desired.ko.Spec.NodeType {
				scaleUpUpdate = true
				break
			}
		}
	}

	// Determine if we are doing scaling up/down along with resharding
	if delta.DifferentAt("Spec.NodeType") && desired.ko.Spec.NodeType != nil &&
		delta.DifferentAt("Spec.NumShards") {
		if latest.ko.Spec.NumShards != nil && desired.ko.Spec.NumShards != nil {
			// If we are scaling in, then perform scale up or down it does not matter.
			if *latest.ko.Spec.NumShards > *desired.ko.Spec.NumShards {
				scaleUpUpdate = true
			}
		}
	}

	// This means we are not scaling out, so we can perform scale up/down. Reason we give preference to scale down
	// instead of scale in is we perform scale down and update engine version together.
	if scaleUpUpdate {
		res.NodeType = desired.ko.Spec.NodeType
	}

	engineUpgradeOrScaling := delta.DifferentAt("Spec.EngineVersion") || scaleUpUpdate

	if !engineUpgradeOrScaling && delta.DifferentAt("Spec.NumShards") && desired.ko.Spec.NumShards != nil {
		shardConfig := &svcsdktypes.ShardConfigurationRequest{}
		shardConfig.ShardCount = int32(*desired.ko.Spec.NumShards)
		res.ShardConfiguration = shardConfig
	}

	reSharding := delta.DifferentAt("Spec.NumShards")

	// Ensure no resharding would be done in API call
	if !reSharding && delta.DifferentAt("Spec.NodeType") && desired.ko.Spec.NodeType != nil {
		res.NodeType = desired.ko.Spec.NodeType
	}

	// If no scaling or engine upgrade then perform replica scaling.
	engineUpgradeOrScaling = delta.DifferentAt("Spec.EngineVersion") || delta.DifferentAt("Spec.NodeType")

	if !engineUpgradeOrScaling && !reSharding &&
		delta.DifferentAt("Spec.NumReplicasPerShard") && desired.ko.Spec.NumReplicasPerShard != nil {
		replicaConfig := &svcsdktypes.ReplicaConfigurationRequest{}
		replicaConfig.ReplicaCount = int32(*desired.ko.Spec.NumReplicasPerShard)
		res.ReplicaConfiguration = replicaConfig
	}

	return res
}

// getEvents gets events from Cluster in service.
func (rm *resourceManager) getEvents(
	ctx context.Context,
	r *resource,
) ([]*svcapitypes.Event, error) {
	maxResults := int32(svcutil.MaxEvents)
	input := &svcsdk.DescribeEventsInput{
		SourceName: r.ko.Spec.Name,
		SourceType: svcsdktypes.SourceTypeCluster,
		MaxResults: &maxResults,
	}
	resp, err := rm.sdkapi.DescribeEvents(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeEvents", err)
	if err != nil {
		rm.log.V(1).Info("Error during DescribeEvents-Cluster", "error", err)
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

// isClusterAvailable returns true when the status of the given Cluster is set to `available`
func (rm *resourceManager) isClusterAvailable(
	latest *resource,
) bool {
	latestStatus := latest.ko.Status.Status
	return latestStatus != nil && *latestStatus == resourceStatusAvailable
}

// getTags gets tags from given Cluster.
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
	tags := resourceTagsFromSDKTags(resp.TagList)
	return tags, nil
}

// updateTags updates tags of given Cluster to desired tags.
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
		rlog.Debug("removing tags from cluster", "tags", toDelete)
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
		rlog.Debug("adding tags to cluster", "tags", toAdd)
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
