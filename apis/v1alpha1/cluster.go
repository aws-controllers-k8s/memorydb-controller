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

// Code generated by ack-generate. DO NOT EDIT.

package v1alpha1

import (
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSpec defines the desired state of Cluster.
//
// Contains all of the attributes of a specific cluster.
type ClusterSpec struct {

	// The name of the Access Control List to associate with the cluster.
	ACLName *string                                  `json:"aclName,omitempty"`
	ACLRef  *ackv1alpha1.AWSResourceReferenceWrapper `json:"aclRef,omitempty"`
	// When set to true, the cluster will automatically receive minor engine version
	// upgrades after launch.
	AutoMinorVersionUpgrade *bool `json:"autoMinorVersionUpgrade,omitempty"`
	// An optional description of the cluster.
	Description *string `json:"description,omitempty"`
	// The name of the engine to be used for the cluster.
	Engine *string `json:"engine,omitempty"`
	// The version number of the Redis OSS engine to be used for the cluster.
	EngineVersion *string `json:"engineVersion,omitempty"`
	// The ID of the KMS key used to encrypt the cluster.
	KMSKeyID *string `json:"kmsKeyID,omitempty"`
	// Specifies the weekly time range during which maintenance on the cluster is
	// performed. It is specified as a range in the format ddd:hh24:mi-ddd:hh24:mi
	// (24H Clock UTC). The minimum maintenance window is a 60 minute period.
	//
	// Valid values for ddd are:
	//
	//   - sun
	//
	//   - mon
	//
	//   - tue
	//
	//   - wed
	//
	//   - thu
	//
	//   - fri
	//
	//   - sat
	//
	// Example: sun:23:00-mon:01:30
	MaintenanceWindow *string `json:"maintenanceWindow,omitempty"`
	// The name of the cluster. This value must be unique as it also serves as the
	// cluster identifier.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// The compute and memory capacity of the nodes in the cluster.
	// +kubebuilder:validation:Required
	NodeType *string `json:"nodeType"`
	// The number of replicas to apply to each shard. The default value is 1. The
	// maximum is 5.
	NumReplicasPerShard *int64 `json:"numReplicasPerShard,omitempty"`
	// The number of shards the cluster will contain. The default value is 1.
	NumShards *int64 `json:"numShards,omitempty"`
	// The name of the parameter group associated with the cluster.
	ParameterGroupName *string                                  `json:"parameterGroupName,omitempty"`
	ParameterGroupRef  *ackv1alpha1.AWSResourceReferenceWrapper `json:"parameterGroupRef,omitempty"`
	// The port number on which each of the nodes accepts connections.
	Port *int64 `json:"port,omitempty"`
	// A list of security group names to associate with this cluster.
	SecurityGroupIDs  []*string                                  `json:"securityGroupIDs,omitempty"`
	SecurityGroupRefs []*ackv1alpha1.AWSResourceReferenceWrapper `json:"securityGroupRefs,omitempty"`
	// A list of Amazon Resource Names (ARN) that uniquely identify the RDB snapshot
	// files stored in Amazon S3. The snapshot files are used to populate the new
	// cluster. The Amazon S3 object name in the ARN cannot contain any commas.
	SnapshotARNs []*string `json:"snapshotARNs,omitempty"`
	// The name of a snapshot from which to restore data into the new cluster. The
	// snapshot status changes to restoring while the new cluster is being created.
	SnapshotName *string                                  `json:"snapshotName,omitempty"`
	SnapshotRef  *ackv1alpha1.AWSResourceReferenceWrapper `json:"snapshotRef,omitempty"`
	// The number of days for which MemoryDB retains automatic snapshots before
	// deleting them. For example, if you set SnapshotRetentionLimit to 5, a snapshot
	// that was taken today is retained for 5 days before being deleted.
	SnapshotRetentionLimit *int64 `json:"snapshotRetentionLimit,omitempty"`
	// The daily time range (in UTC) during which MemoryDB begins taking a daily
	// snapshot of your shard.
	//
	// Example: 05:00-09:00
	//
	// If you do not specify this parameter, MemoryDB automatically chooses an appropriate
	// time range.
	SnapshotWindow *string `json:"snapshotWindow,omitempty"`
	// The Amazon Resource Name (ARN) of the Amazon Simple Notification Service
	// (SNS) topic to which notifications are sent.
	SNSTopicARN *string                                  `json:"snsTopicARN,omitempty"`
	SNSTopicRef *ackv1alpha1.AWSResourceReferenceWrapper `json:"snsTopicRef,omitempty"`
	// The name of the subnet group to be used for the cluster.
	SubnetGroupName *string                                  `json:"subnetGroupName,omitempty"`
	SubnetGroupRef  *ackv1alpha1.AWSResourceReferenceWrapper `json:"subnetGroupRef,omitempty"`
	// A flag to enable in-transit encryption on the cluster.
	TLSEnabled *bool `json:"tlsEnabled,omitempty"`
	// A list of tags to be added to this resource. Tags are comma-separated key,value
	// pairs (e.g. Key=myKey, Value=myKeyValue. You can include multiple tags as
	// shown following: Key=myKey, Value=myKeyValue Key=mySecondKey, Value=mySecondKeyValue.
	Tags []*Tag `json:"tags,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRS managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// A list node types which you can use to scale down your cluster.
	// +kubebuilder:validation:Optional
	AllowedScaleDownNodeTypes []*string `json:"allowedScaleDownNodeTypes,omitempty"`
	// A list node types which you can use to scale up your cluster.
	// +kubebuilder:validation:Optional
	AllowedScaleUpNodeTypes []*string `json:"allowedScaleUpNodeTypes,omitempty"`
	// Indicates if the cluster has a Multi-AZ configuration (multiaz) or not (singleaz).
	// +kubebuilder:validation:Optional
	AvailabilityMode *string `json:"availabilityMode,omitempty"`
	// The cluster's configuration endpoint
	// +kubebuilder:validation:Optional
	ClusterEndpoint *Endpoint `json:"clusterEndpoint,omitempty"`
	// Enables data tiering. Data tiering is only supported for clusters using the
	// r6gd node type. This parameter must be set when using r6gd nodes. For more
	// information, see Data tiering (https://docs.aws.amazon.com/memorydb/latest/devguide/data-tiering.html).
	// +kubebuilder:validation:Optional
	DataTiering *string `json:"dataTiering,omitempty"`
	// The Redis OSS engine patch version used by the cluster
	// +kubebuilder:validation:Optional
	EnginePatchVersion *string `json:"enginePatchVersion,omitempty"`
	// A list of events. Each element in the list contains detailed information
	// about one event.
	// +kubebuilder:validation:Optional
	Events []*Event `json:"events,omitempty"`
	// The name of the multi-Region cluster that this cluster belongs to.
	// +kubebuilder:validation:Optional
	MultiRegionClusterName *string `json:"multiRegionClusterName,omitempty"`
	// The number of shards in the cluster
	// +kubebuilder:validation:Optional
	NumberOfShards *int64 `json:"numberOfShards,omitempty"`
	// The status of the parameter group used by the cluster, for example 'active'
	// or 'applying'.
	// +kubebuilder:validation:Optional
	ParameterGroupStatus *string `json:"parameterGroupStatus,omitempty"`
	// A group of settings that are currently being applied.
	// +kubebuilder:validation:Optional
	PendingUpdates *ClusterPendingUpdates `json:"pendingUpdates,omitempty"`
	// A list of security groups used by the cluster
	// +kubebuilder:validation:Optional
	SecurityGroups []*SecurityGroupMembership `json:"securityGroups,omitempty"`
	// A list of shards that are members of the cluster.
	// +kubebuilder:validation:Optional
	Shards []*Shard `json:"shards,omitempty"`
	// The SNS topic must be in Active status to receive notifications
	// +kubebuilder:validation:Optional
	SNSTopicStatus *string `json:"snsTopicStatus,omitempty"`
	// The status of the cluster. For example, Available, Updating, Creating.
	// +kubebuilder:validation:Optional
	Status *string `json:"status,omitempty"`
}

// Cluster is the Schema for the Clusters API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ClusterSpec   `json:"spec,omitempty"`
	Status            ClusterStatus `json:"status,omitempty"`
}

// ClusterList contains a list of Cluster
// +kubebuilder:object:root=true
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
