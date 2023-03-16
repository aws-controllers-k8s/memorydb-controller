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
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
	"github.com/samber/lo"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

const (
	maxNumberOfParametersUpdate = 20
)

func (rm *resourceManager) setParameters(
	ctx context.Context,
	ko *svcapitypes.ParameterGroup,
) (*svcapitypes.ParameterGroup, error) {
	if ko.Spec.Name == nil {
		return ko, nil
	}

	parameters, err := rm.describeParameters(ctx, ko.Spec.Name)
	if err != nil {
		return nil, err
	}
	parameterNameValues := []*svcapitypes.ParameterNameValue{}
	for _, p := range parameters {
		sp := svcapitypes.ParameterNameValue{
			ParameterName:  p.Name,
			ParameterValue: p.Value,
		}
		parameterNameValues = append(parameterNameValues, &sp)
	}
	ko.Spec.ParameterNameValues = parameterNameValues
	// Update parameter value to nil to not duplicate status and spec values.
	for _, p := range parameters {
		p.Value = nil
	}
	ko.Status.Parameters = parameters
	return ko, nil
}

func (rm *resourceManager) describeParameters(
	ctx context.Context,
	parameterGroupName *string,
) ([]*svcapitypes.Parameter, error) {

	parameters := []*svcapitypes.Parameter{}
	var nextToken *string = nil
	for {
		response, respErr := rm.sdkapi.DescribeParametersWithContext(ctx, &svcsdk.DescribeParametersInput{
			ParameterGroupName: parameterGroupName,
			NextToken:          nextToken,
		})
		rm.metrics.RecordAPICall("GET", "DescribeParameters", respErr)
		if respErr != nil {
			if awsErr, ok := ackerr.AWSError(respErr); ok && awsErr.Code() == "ParameterGroupNotFound" {
				return nil, ackerr.NotFound
			}
			return nil, respErr
		}
		if response.Parameters == nil {
			break
		}

		// Update the next token
		nextToken = response.NextToken

		for _, p := range response.Parameters {
			sp := svcapitypes.Parameter{
				Name:                 p.Name,
				Value:                p.Value,
				Description:          p.Description,
				DataType:             p.DataType,
				AllowedValues:        p.AllowedValues,
				MinimumEngineVersion: p.MinimumEngineVersion,
			}
			parameters = append(parameters, &sp)
		}

		if nextToken == nil || *nextToken == "" {
			break
		}
	}

	return parameters, nil
}

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

// customUpdate overrides sdkUpdate by custom logic
func (rm *resourceManager) customUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.customUpdate")
	defer func() {
		exit(err)
	}()
	if delta.DifferentAt("Spec.Tags") {
		err = rm.updateTags(ctx, desired, latest)
		if err != nil {
			return nil, err
		}
	}

	if delta.DifferentAt("Spec.ParameterNameValues") {
		ko, err := rm.resetParameterGroup(ctx, desired, latest)

		if ko != nil || err != nil {
			return ko, err
		}
	}

	if !delta.DifferentExcept("Spec.Tags") {
		return desired, nil
	}

	inputs := rm.updateRequestPayload(desired)

	var resp *svcsdk.UpdateParameterGroupOutput

	// Update 20 parameters each time
	for _, input := range inputs {
		resp, err = rm.sdkapi.UpdateParameterGroupWithContext(ctx, input)
		rm.metrics.RecordAPICall("UPDATE", "UpdateParameterGroup", err)
		if err != nil {
			return nil, err
		}
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}

	if resp != nil {
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
		if resp.ParameterGroup.Name != nil {
			ko.Spec.Name = resp.ParameterGroup.Name
		} else {
			ko.Spec.Name = nil
		}
	}

	rm.setStatusDefaults(ko)
	ko, err = rm.setParameters(ctx, ko)

	if err != nil {
		return nil, err
	}
	return &resource{ko}, nil
}

// updateRequestPayload returns an array of SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
// each element in return value has maximum 20 parameters
func (rm *resourceManager) updateRequestPayload(
	desired *resource,
) []*svcsdk.UpdateParameterGroupInput {
	parameterNameValues := lo.Chunk(desired.ko.Spec.ParameterNameValues, maxNumberOfParametersUpdate)
	return lo.Map(parameterNameValues, func(parameters []*svcapitypes.ParameterNameValue, index int) *svcsdk.UpdateParameterGroupInput {
		res := &svcsdk.UpdateParameterGroupInput{}
		if desired.ko.Spec.Name != nil {
			res.SetParameterGroupName(*desired.ko.Spec.Name)
		}
		f1 := []*svcsdk.ParameterNameValue{}
		for _, f1iter := range parameters {
			f1elem := &svcsdk.ParameterNameValue{}
			if f1iter.ParameterName != nil {
				f1elem.SetParameterName(*f1iter.ParameterName)
			}
			if f1iter.ParameterValue != nil {
				f1elem.SetParameterValue(*f1iter.ParameterValue)
			}
			f1 = append(f1, f1elem)
		}
		res.SetParameterNameValues(f1)
		return res
	})
}

// getTags gets tags from given ParameterGroup.
func (rm *resourceManager) getTags(
	ctx context.Context,
	resourceARN string,
) ([]*svcapitypes.Tag, error) {
	resp, err := rm.sdkapi.ListTagsWithContext(
		ctx,
		&svcsdk.ListTagsInput{
			ResourceArn: &resourceARN,
		},
	)
	rm.metrics.RecordAPICall("GET", "ListTags", err)
	if err != nil {
		return nil, err
	}
	tags := resourceTagsFromSDKTags(resp.TagList)
	return tags, nil
}

// updateTags updates tags of given ParameterGroup to desired tags.
func (rm *resourceManager) updateTags(
	ctx context.Context,
	desired *resource,
	latest *resource,
) (err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.syncTags")
	defer func() { exit(err) }()

	arn := (*string)(latest.ko.Status.ACKResourceMetadata.ARN)

	desiredTags := ToACKTags(desired.ko.Spec.Tags)
	latestTags := ToACKTags(latest.ko.Spec.Tags)

	added, _, removed := ackcompare.GetTagsDifference(latestTags, desiredTags)

	toAdd := FromACKTags(added)
	toRemove := FromACKTags(removed)

	var toDelete []*string
	for _, removedElement := range toRemove {
		toDelete = append(toDelete, removedElement.Key)
	}

	if len(toDelete) > 0 {
		rlog.Debug("removing tags from parameter group", "tags", toDelete)
		_, err = rm.sdkapi.UntagResourceWithContext(
			ctx,
			&svcsdk.UntagResourceInput{
				ResourceArn: arn,
				TagKeys:     toDelete,
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "UntagResource", err)
		if err != nil {
			return err
		}
	}

	if len(toAdd) > 0 {
		rlog.Debug("adding tags to parameter group", "tags", toAdd)
		_, err = rm.sdkapi.TagResourceWithContext(
			ctx,
			&svcsdk.TagResourceInput{
				ResourceArn: arn,
				Tags:        sdkTagsFromResourceTags(toAdd),
			},
		)
		rm.metrics.RecordAPICall("UPDATE", "TagResource", err)
		if err != nil {
			return err
		}
	}

	return nil
}

func sdkTagsFromResourceTags(
	rTags []*svcapitypes.Tag,
) []*svcsdk.Tag {
	tags := make([]*svcsdk.Tag, len(rTags))
	for i := range rTags {
		tags[i] = &svcsdk.Tag{
			Key:   rTags[i].Key,
			Value: rTags[i].Value,
		}
	}
	return tags
}

func resourceTagsFromSDKTags(
	sdkTags []*svcsdk.Tag,
) []*svcapitypes.Tag {
	tags := make([]*svcapitypes.Tag, len(sdkTags))
	for i := range sdkTags {
		tags[i] = &svcapitypes.Tag{
			Key:   sdkTags[i].Key,
			Value: sdkTags[i].Value,
		}
	}
	return tags
}
