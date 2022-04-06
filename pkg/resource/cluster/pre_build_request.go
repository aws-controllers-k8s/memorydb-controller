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
	"fmt"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
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
	if latestStatus == nil || *latestStatus != StatusAvailable {
		return nil, requeue.NeededAfter(
			fmt.Errorf("cluster cannot be updated as its status is not '%s'", StatusAvailable),
			requeue.DefaultRequeueAfterDuration)
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
