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
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
)

func (rm *resourceManager) newMemoryDBClusterUploadPayload(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) *svcsdk.UpdateClusterInput {
	res := &svcsdk.UpdateClusterInput{}

	if delta.DifferentAt("Spec.AclName") && desired.ko.Spec.ACLName != nil {
		res.SetACLName(*desired.ko.Spec.ACLName)
	}
	if desired.ko.Spec.Name != nil {
		res.SetClusterName(*desired.ko.Spec.Name)
	}
	if delta.DifferentAt("Spec.Description") && desired.ko.Spec.Description != nil {
		res.SetDescription(*desired.ko.Spec.Description)
	}
	if delta.DifferentAt("Spec.MaintenanceWindow") && desired.ko.Spec.MaintenanceWindow != nil {
		res.SetMaintenanceWindow(*desired.ko.Spec.MaintenanceWindow)
	}
	if delta.DifferentAt("Spec.ParameterGroupName") && desired.ko.Spec.ParameterGroupName != nil {
		res.SetParameterGroupName(*desired.ko.Spec.ParameterGroupName)
	}
	if delta.DifferentAt("Spec.SecurityGroupIDs") && desired.ko.Spec.SecurityGroupIDs != nil {
		f8 := []*string{}
		for _, f8iter := range desired.ko.Spec.SecurityGroupIDs {
			var f8elem string
			f8elem = *f8iter
			f8 = append(f8, &f8elem)
		}
		res.SetSecurityGroupIds(f8)
	}
	if delta.DifferentAt("Spec.SnapshotRetentionLimit") && desired.ko.Spec.SnapshotRetentionLimit != nil {
		res.SetSnapshotRetentionLimit(*desired.ko.Spec.SnapshotRetentionLimit)
	}
	if delta.DifferentAt("Spec.SnapshotWindow") && desired.ko.Spec.SnapshotWindow != nil {
		res.SetSnapshotWindow(*desired.ko.Spec.SnapshotWindow)
	}
	if delta.DifferentAt("Spec.SnsTopicARN") && desired.ko.Spec.SnsTopicARN != nil {
		res.SetSnsTopicArn(*desired.ko.Spec.SnsTopicARN)
	}
	if delta.DifferentAt("Spec.SnsTopicStatus") && desired.ko.Status.SnsTopicStatus != nil {
		res.SetSnsTopicStatus(*desired.ko.Status.SnsTopicStatus)
	}

	if delta.DifferentAt("Spec.EngineVersion") && desired.ko.Spec.EngineVersion != nil {
		res.SetEngineVersion(*desired.ko.Spec.EngineVersion)
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
		res.SetNodeType(*desired.ko.Spec.NodeType)
	}

	engineUpgradeOrScaling := delta.DifferentAt("Spec.EngineVersion") || scaleUpUpdate

	if !engineUpgradeOrScaling && delta.DifferentAt("Spec.NumShards") && desired.ko.Spec.NumShards != nil {
		shardConfig := &svcsdk.ShardConfigurationRequest{}
		shardConfig.SetShardCount(*desired.ko.Spec.NumShards)
		res.SetShardConfiguration(shardConfig)
	}

	reSharding := delta.DifferentAt("Spec.NumShards")

	// Ensure no resharding would be done in API call
	if !reSharding && delta.DifferentAt("Spec.NodeType") && desired.ko.Spec.NodeType != nil {
		res.SetNodeType(*desired.ko.Spec.NodeType)
	}

	// If no scaling or engine upgrade then perform replica scaling.
	engineUpgradeOrScaling = delta.DifferentAt("Spec.EngineVersion") || delta.DifferentAt("Spec.NodeType")

	if !engineUpgradeOrScaling && !reSharding &&
		delta.DifferentAt("Spec.NumReplicasPerShard") && desired.ko.Spec.NumReplicasPerShard != nil {
		replicaConfig := &svcsdk.ReplicaConfigurationRequest{}
		replicaConfig.SetReplicaCount(*desired.ko.Spec.NumReplicasPerShard)
		res.SetReplicaConfiguration(replicaConfig)
	}

	return res
}
