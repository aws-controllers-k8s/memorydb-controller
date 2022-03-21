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

package acl

import (
	"github.com/pkg/errors"

	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

// validateACLNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable
func (rm *resourceManager) validateACLNeedsUpdate(
	latest *resource,
) error {

	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != "active" {
		return requeue.NeededAfter(
			errors.New("ACL cannot be modified as its status is not 'active'."),
			requeue.DefaultRequeueAfterDuration)
	}

	return nil
}
