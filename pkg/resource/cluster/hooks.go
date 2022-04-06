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
	"errors"

	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

var (
	condMsgCurrentlyDeleting     = "cluster currently being deleted"
	condMsgNoDeleteWhileUpdating = "cluster is being updated. cannot delete"
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("delete is in progress"),
		ackrequeue.DefaultRequeueAfterDuration,
	)
	requeueWaitWhileUpdating = ackrequeue.NeededAfter(
		errors.New("update is in progress"),
		ackrequeue.DefaultRequeueAfterDuration,
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
