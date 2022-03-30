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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/pkg/errors"
	"strconv"

	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

// validateClusterNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable and
// to validate if resource update is required.
func (rm *resourceManager) validateClusterNeedsUpdate(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {
	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != "available" {
		return nil, requeue.NeededAfter(
			errors.New("Cluster cannot be updated as its status is not 'active'."),
			requeue.DefaultRequeueAfterDuration)
	}

	// Set terminal condition when cluster is in create-failed state
	if latestStatus != nil || *latestStatus == "create-failed" {
		return nil, awserr.New("InvalidNodeStateFault", "Cluster is in create-failed state, cannot be updated.", nil)
	}

	annotations := desired.ko.ObjectMeta.GetAnnotations()

	// Handle asynchronous rollback of NodeType. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNodeType]; ok && desired.ko.Spec.NodeType != nil {
		if val == *desired.ko.Spec.NodeType && delta.DifferentAt("Spec.NodeType") {
			return nil, awserr.New("InvalidParameterCombinationException", "Cannot update NodeType.", nil)
		}
	}

	// Handle asynchronous rollback of NumShards. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNumShards]; ok && desired.ko.Spec.NumShards != nil {
		numShards, err := strconv.ParseInt(val, 10, 64)
		if err == nil && numShards == *desired.ko.Spec.NumShards && delta.DifferentAt("Spec.NumShards") {
			return nil, awserr.New("InvalidParameterCombinationException", "Cannot update NumShards.", nil)
		}
	}

	// Handle asynchronous rollback of NumShards. This can happen due to ICE, InsufficientMemoryException or other
	// errors TODO update the error message once we add describe events support
	if val, ok := annotations[AnnotationLastRequestedNumReplicasPerShard]; ok && desired.ko.Spec.NumReplicasPerShard != nil {
		numReplicasPerShard, err := strconv.ParseInt(val, 10, 64)
		if err == nil && numReplicasPerShard == *desired.ko.Spec.NumReplicasPerShard && delta.DifferentAt("Spec.NumReplicasPerShard") {
			return nil, awserr.New("InvalidParameterCombinationException", "Cannot update NumReplicasPerShard.", nil)
		}
	}

	return nil, nil
}
