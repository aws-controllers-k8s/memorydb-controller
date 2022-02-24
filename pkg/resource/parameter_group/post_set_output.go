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

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
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
