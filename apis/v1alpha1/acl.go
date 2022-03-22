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

// ACLSpec defines the desired state of ACL.
//
// An Access Control List. You can authenticate users with Access Contol Lists.
// ACLs enable you to control cluster access by grouping users. These Access
// control lists are designed as a way to organize access to clusters.
type ACLSpec struct {
	// The name of the Access Control List.
	// +kubebuilder:validation:Required
	Name *string `json:"name"`
	// A list of tags to be added to this resource. A tag is a key-value pair. A
	// tag key must be accompanied by a tag value, although null is accepted.
	Tags []*Tag `json:"tags,omitempty"`
	// The list of users that belong to the Access Control List.
	UserNames []*string `json:"userNames,omitempty"`
}

// ACLStatus defines the observed state of ACL
type ACLStatus struct {
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
	// A list of clusters associated with the ACL.
	// +kubebuilder:validation:Optional
	Clusters []*string `json:"clusters,omitempty"`
	// The minimum engine version supported for the ACL
	// +kubebuilder:validation:Optional
	MinimumEngineVersion *string `json:"minimumEngineVersion,omitempty"`
	// A list of updates being applied to the ACL.
	// +kubebuilder:validation:Optional
	PendingChanges *ACLPendingChanges `json:"pendingChanges,omitempty"`
	// Indicates ACL status. Can be "creating", "active", "modifying", "deleting".
	// +kubebuilder:validation:Optional
	Status *string `json:"status,omitempty"`
}

// ACL is the Schema for the ACLS API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ACL struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              ACLSpec   `json:"spec,omitempty"`
	Status            ACLStatus `json:"status,omitempty"`
}

// ACLList contains a list of ACL
// +kubebuilder:object:root=true
type ACLList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ACL `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ACL{}, &ACLList{})
}