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

// SubnetGroupSpec defines the desired state of SubnetGroup.
//
// Represents the output of one of the following operations:
//
//   - CreateSubnetGroup
//
//   - UpdateSubnetGroup
//
// A subnet group is a collection of subnets (typically private) that you can
// designate for your clusters running in an Amazon Virtual Private Cloud (VPC)
// environment.
type SubnetGroupSpec struct {

	// A description for the subnet group.
	Description *string `json:"description,omitempty"`
	// The name of the subnet group.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// A list of VPC subnet IDs for the subnet group.
	SubnetIDs  []*string                                  `json:"subnetIDs,omitempty"`
	SubnetRefs []*ackv1alpha1.AWSResourceReferenceWrapper `json:"subnetRefs,omitempty"`
	// A list of tags to be added to this resource. A tag is a key-value pair. A
	// tag key must be accompanied by a tag value, although null is accepted.
	Tags []*Tag `json:"tags,omitempty"`
}

// SubnetGroupStatus defines the observed state of SubnetGroup
type SubnetGroupStatus struct {
	// All CRs managed by ACK have a common `Status.ACKResourceMetadata` member
	// that is used to contain resource sync state, account ownership,
	// constructed ARN for the resource
	// +kubebuilder:validation:Optional
	ACKResourceMetadata *ackv1alpha1.ResourceMetadata `json:"ackResourceMetadata"`
	// All CRs managed by ACK have a common `Status.Conditions` member that
	// contains a collection of `ackv1alpha1.Condition` objects that describe
	// the various terminal states of the CR and its backend AWS service API
	// resource
	// +kubebuilder:validation:Optional
	Conditions []*ackv1alpha1.Condition `json:"conditions"`
	// A list of subnets associated with the subnet group.
	// +kubebuilder:validation:Optional
	Subnets []*Subnet `json:"subnets,omitempty"`
	// The Amazon Virtual Private Cloud identifier (VPC ID) of the subnet group.
	// +kubebuilder:validation:Optional
	VPCID *string `json:"vpcID,omitempty"`
}

// SubnetGroup is the Schema for the SubnetGroups API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type SubnetGroup struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SubnetGroupSpec   `json:"spec,omitempty"`
	Status            SubnetGroupStatus `json:"status,omitempty"`
}

// SubnetGroupList contains a list of SubnetGroup
// +kubebuilder:object:root=true
type SubnetGroupList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SubnetGroup `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SubnetGroup{}, &SubnetGroupList{})
}
