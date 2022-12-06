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

package user

import (
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	"github.com/pkg/errors"

	"github.com/aws-controllers-k8s/runtime/pkg/requeue"
)

// validateUserNeedsUpdate this function's purpose is to requeue if the resource is currently unavailable and
// to validate if resource update is required.
func (rm *resourceManager) validateUserNeedsUpdate(
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (*resource, error) {

	// requeue if necessary
	latestStatus := latest.ko.Status.Status
	if latestStatus == nil || *latestStatus != "active" {
		return nil, requeue.NeededAfter(
			errors.New("User cannot be modified as its status is not 'active'."),
			requeue.DefaultRequeueAfterDuration)
	}

	// AccessString might be transformed by the server to a different value when this happens
	// delta would be generated when old access string is supplied. We would need to re-patch the
	// Spec so it is updated with the transformed value.
	annotations := desired.ko.ObjectMeta.GetAnnotations()
	if val, ok := annotations[AnnotationLastRequestedAccessString]; ok && desired.ko.Spec.AccessString != nil {
		if val == *desired.ko.Spec.AccessString && delta.DifferentAt("Spec.AccessString") && len(delta.Differences) == 1 {
			return latest, nil
		}
	}

	return nil, nil
}
