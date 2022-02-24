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

package parameter_group

import (
	"context"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
)

func (rm *resourceManager) resetParameterGroup(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (updated *resource, err error) {
	if len(desired.ko.Spec.ParameterNameValues) == 0 {
		// Reset all the parameters to default
		return rm.resetAllParameters(ctx, desired)
	}

	return nil, nil
}

// resetAllParameters resets cache parameters for given ParameterGroup in desired custom resource.
func (rm *resourceManager) resetAllParameters(
	ctx context.Context,
	desired *resource,
) (updated *resource, err error) {
	input := &svcsdk.ResetParameterGroupInput{}
	input.SetParameterGroupName(*desired.ko.Spec.Name)
	input.SetAllParameters(true)

	resp, err := rm.sdkapi.ResetParameterGroupWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "ResetParameterGroup", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ParameterGroup.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ParameterGroup.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ParameterGroup.Description != nil {
		ko.Spec.Description = resp.ParameterGroup.Description
	} else {
		ko.Spec.Description = nil
	}
	if resp.ParameterGroup.Family != nil {
		ko.Spec.Family = resp.ParameterGroup.Family
	} else {
		ko.Spec.Family = nil
	}

	rm.setStatusDefaults(ko)
	ko, err = rm.setParameters(ctx, ko)

	if err != nil {
		return nil, err
	}
	return &resource{ko}, nil
}
