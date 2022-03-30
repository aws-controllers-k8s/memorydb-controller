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
	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	"strconv"
)

const (
	// AnnotationLastRequestedNumShards is an annotation whose value is the value of NumShards
	// passed in as input to either the create or modify API called most recently
	AnnotationLastRequestedNumShards = svcapitypes.AnnotationPrefix + "last-requested-num-shards"
	// AnnotationLastRequestedNumReplicasPerShard is an annotation whose value is the value of NumReplicasPerShard
	// passed in as input to either the create or modify API called most recently
	AnnotationLastRequestedNumReplicasPerShard = svcapitypes.AnnotationPrefix + "last-requested-num-replicas-per-shard"
	// AnnotationLastRequestedNodeType is an annotation whose value is the value of NodeType
	// passed in as input to either the create or modify API called most recently
	AnnotationLastRequestedNodeType = svcapitypes.AnnotationPrefix + "last-requested-node-type"
)

// setNumShardAnnotation sets the AnnotationLastRequestedNumShards annotation for cluster resource
// This should only be called upon a successful create or modify call.
func (rm *resourceManager) setNumShardAnnotation(
	numShards *int64,
	ko *svcapitypes.Cluster,
) {
	if ko.ObjectMeta.Annotations == nil {
		ko.ObjectMeta.Annotations = make(map[string]string)
	}

	annotations := ko.ObjectMeta.Annotations

	if numShards != nil {
		annotations[AnnotationLastRequestedNumShards] = strconv.FormatInt(*numShards, 10)
	}
}

// setNumReplicasPerShardAnnotation sets the AnnotationLastRequestedNumReplicasPerShard annotation for cluster resource
// This should only be called upon a successful create or modify call.
func (rm *resourceManager) setNumReplicasPerShardAnnotation(
	numReplicasPerShard *int64,
	ko *svcapitypes.Cluster,
) {
	if ko.ObjectMeta.Annotations == nil {
		ko.ObjectMeta.Annotations = make(map[string]string)
	}

	annotations := ko.ObjectMeta.Annotations

	if numReplicasPerShard != nil {
		annotations[AnnotationLastRequestedNumReplicasPerShard] = strconv.FormatInt(*numReplicasPerShard, 10)
	}
}

// setNodeTypeAnnotation sets the AnnotationLastRequestedNodeType annotation for cluster resource
// This should only be called upon a successful create or modify call.
func (rm *resourceManager) setNodeTypeAnnotation(
	nodeType *string,
	ko *svcapitypes.Cluster,
) {
	if ko.ObjectMeta.Annotations == nil {
		ko.ObjectMeta.Annotations = make(map[string]string)
	}

	annotations := ko.ObjectMeta.Annotations

	if nodeType != nil {
		annotations[AnnotationLastRequestedNodeType] = *nodeType
	}
}
