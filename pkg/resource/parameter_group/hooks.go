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
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"

	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
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
	tags := make([]*svcapitypes.Tag, 0, len(resp.TagList))
	for _, tag := range resp.TagList {
		tags = append(tags, &svcapitypes.Tag{
			Key:   tag.Key,
			Value: tag.Value,
		})
	}
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

	toAdd, toDelete := computeTagsDelta(
		desired.ko.Spec.Tags, latest.ko.Spec.Tags,
	)

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

func computeTagsDelta(
	desired []*svcapitypes.Tag,
	latest []*svcapitypes.Tag,
) (added []*svcapitypes.Tag, removed []*string) {
	toDelete := []*string{}
	toAdd := []*svcapitypes.Tag{}

	desiredTags := map[string]string{}
	key := ""
	value := ""
	for _, tag := range desired {
		if tag.Key != nil {
			key = *tag.Key
			value = ""
			if tag.Value != nil {
				value = *tag.Value
			}
			desiredTags[key] = value
		}
	}

	for _, tag := range desired {
		toAdd = append(toAdd, tag)
	}
	for _, tag := range latest {
		_, ok := desiredTags[*tag.Key]
		if !ok {
			toDelete = append(toDelete, tag.Key)
		}
	}
	return toAdd, toDelete
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
