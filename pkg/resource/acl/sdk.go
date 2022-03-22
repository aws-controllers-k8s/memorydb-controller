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

package acl

import (
	"context"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go/aws"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &aws.JSONValue{}
	_ = &svcsdk.MemoryDB{}
	_ = &svcapitypes.ACL{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer exit(err)
	// If any required fields in the input shape are missing, AWS resource is
	// not created yet. Return NotFound here to indicate to callers that the
	// resource isn't yet created.
	if rm.requiredFieldsMissingFromReadManyInput(r) {
		return nil, ackerr.NotFound
	}

	input, err := rm.newListRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DescribeACLsOutput
	resp, err = rm.sdkapi.DescribeACLsWithContext(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeACLs", err)
	if err != nil {
		if awsErr, ok := ackerr.AWSError(err); ok && awsErr.Code() == "ACLNotFoundFault" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	found := false
	for _, elem := range resp.ACLs {
		if elem.ARN != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.ARN)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.Clusters != nil {
			f1 := []*string{}
			for _, f1iter := range elem.Clusters {
				var f1elem string
				f1elem = *f1iter
				f1 = append(f1, &f1elem)
			}
			ko.Status.Clusters = f1
		} else {
			ko.Status.Clusters = nil
		}
		if elem.MinimumEngineVersion != nil {
			ko.Status.MinimumEngineVersion = elem.MinimumEngineVersion
		} else {
			ko.Status.MinimumEngineVersion = nil
		}
		if elem.Name != nil {
			ko.Spec.Name = elem.Name
		} else {
			ko.Spec.Name = nil
		}
		if elem.PendingChanges != nil {
			f4 := &svcapitypes.ACLPendingChanges{}
			if elem.PendingChanges.UserNamesToAdd != nil {
				f4f0 := []*string{}
				for _, f4f0iter := range elem.PendingChanges.UserNamesToAdd {
					var f4f0elem string
					f4f0elem = *f4f0iter
					f4f0 = append(f4f0, &f4f0elem)
				}
				f4.UserNamesToAdd = f4f0
			}
			if elem.PendingChanges.UserNamesToRemove != nil {
				f4f1 := []*string{}
				for _, f4f1iter := range elem.PendingChanges.UserNamesToRemove {
					var f4f1elem string
					f4f1elem = *f4f1iter
					f4f1 = append(f4f1, &f4f1elem)
				}
				f4.UserNamesToRemove = f4f1
			}
			ko.Status.PendingChanges = f4
		} else {
			ko.Status.PendingChanges = nil
		}
		if elem.Status != nil {
			ko.Status.Status = elem.Status
		} else {
			ko.Status.Status = nil
		}
		if elem.UserNames != nil {
			f6 := []*string{}
			for _, f6iter := range elem.UserNames {
				var f6elem string
				f6elem = *f6iter
				f6 = append(f6, &f6elem)
			}
			ko.Spec.UserNames = f6
		} else {
			ko.Spec.UserNames = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// requiredFieldsMissingFromReadManyInput returns true if there are any fields
// for the ReadMany Input shape that are required but not present in the
// resource's Spec or Status
func (rm *resourceManager) requiredFieldsMissingFromReadManyInput(
	r *resource,
) bool {
	return r.ko.Spec.Name == nil

}

// newListRequestPayload returns SDK-specific struct for the HTTP request
// payload of the List API call for the resource
func (rm *resourceManager) newListRequestPayload(
	r *resource,
) (*svcsdk.DescribeACLsInput, error) {
	res := &svcsdk.DescribeACLsInput{}

	if r.ko.Spec.Name != nil {
		res.SetACLName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkCreate creates the supplied resource in the backend AWS service API and
// returns a copy of the resource with resource fields (in both Spec and
// Status) filled in with values from the CREATE API operation's Output shape.
func (rm *resourceManager) sdkCreate(
	ctx context.Context,
	desired *resource,
) (created *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkCreate")
	defer exit(err)
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateACLOutput
	_ = resp
	resp, err = rm.sdkapi.CreateACLWithContext(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateACL", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ACL.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ACL.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ACL.Clusters != nil {
		f1 := []*string{}
		for _, f1iter := range resp.ACL.Clusters {
			var f1elem string
			f1elem = *f1iter
			f1 = append(f1, &f1elem)
		}
		ko.Status.Clusters = f1
	} else {
		ko.Status.Clusters = nil
	}
	if resp.ACL.MinimumEngineVersion != nil {
		ko.Status.MinimumEngineVersion = resp.ACL.MinimumEngineVersion
	} else {
		ko.Status.MinimumEngineVersion = nil
	}
	if resp.ACL.Name != nil {
		ko.Spec.Name = resp.ACL.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ACL.PendingChanges != nil {
		f4 := &svcapitypes.ACLPendingChanges{}
		if resp.ACL.PendingChanges.UserNamesToAdd != nil {
			f4f0 := []*string{}
			for _, f4f0iter := range resp.ACL.PendingChanges.UserNamesToAdd {
				var f4f0elem string
				f4f0elem = *f4f0iter
				f4f0 = append(f4f0, &f4f0elem)
			}
			f4.UserNamesToAdd = f4f0
		}
		if resp.ACL.PendingChanges.UserNamesToRemove != nil {
			f4f1 := []*string{}
			for _, f4f1iter := range resp.ACL.PendingChanges.UserNamesToRemove {
				var f4f1elem string
				f4f1elem = *f4f1iter
				f4f1 = append(f4f1, &f4f1elem)
			}
			f4.UserNamesToRemove = f4f1
		}
		ko.Status.PendingChanges = f4
	} else {
		ko.Status.PendingChanges = nil
	}
	if resp.ACL.Status != nil {
		ko.Status.Status = resp.ACL.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.ACL.UserNames != nil {
		f6 := []*string{}
		for _, f6iter := range resp.ACL.UserNames {
			var f6elem string
			f6elem = *f6iter
			f6 = append(f6, &f6elem)
		}
		ko.Spec.UserNames = f6
	} else {
		ko.Spec.UserNames = nil
	}

	rm.setStatusDefaults(ko)
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateACLInput, error) {
	res := &svcsdk.CreateACLInput{}

	if r.ko.Spec.Name != nil {
		res.SetACLName(*r.ko.Spec.Name)
	}
	if r.ko.Spec.Tags != nil {
		f1 := []*svcsdk.Tag{}
		for _, f1iter := range r.ko.Spec.Tags {
			f1elem := &svcsdk.Tag{}
			if f1iter.Key != nil {
				f1elem.SetKey(*f1iter.Key)
			}
			if f1iter.Value != nil {
				f1elem.SetValue(*f1iter.Value)
			}
			f1 = append(f1, f1elem)
		}
		res.SetTags(f1)
	}
	if r.ko.Spec.UserNames != nil {
		f2 := []*string{}
		for _, f2iter := range r.ko.Spec.UserNames {
			var f2elem string
			f2elem = *f2iter
			f2 = append(f2, &f2elem)
		}
		res.SetUserNames(f2)
	}

	return res, nil
}

// sdkUpdate patches the supplied resource in the backend AWS service API and
// returns a new resource with updated fields.
func (rm *resourceManager) sdkUpdate(
	ctx context.Context,
	desired *resource,
	latest *resource,
	delta *ackcompare.Delta,
) (updated *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkUpdate")
	defer exit(err)
	validationErr := rm.validateACLNeedsUpdate(latest)

	if validationErr != nil {
		return nil, err
	}
	input, err := rm.newUpdateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}
	createMapForUserNames := func(userIds []*string) map[string]bool {
		userIdsMap := make(map[string]bool)

		for _, userId := range userIds {
			userIdsMap[*userId] = true
		}

		return userIdsMap
	}

	for _, diff := range delta.Differences {
		if diff.Path.Contains("Spec.UserNames") {
			existingUserNamesMap := createMapForUserNames(diff.B.([]*string))
			requiredUserNamesMap := createMapForUserNames(diff.A.([]*string))

			// If a user ID is not required to be deleted or added set its value as false
			for userName, _ := range existingUserNamesMap {
				if _, ok := requiredUserNamesMap[userName]; ok {
					requiredUserNamesMap[userName] = false
					existingUserNamesMap[userName] = false
				}
			}

			if err != nil {
				return nil, err
			}

			// User Ids to add
			{
				var userNamesToAdd []*string

				for userName, include := range requiredUserNamesMap {
					if include {
						userNamesToAdd = append(userNamesToAdd, &userName)
					}
				}

				input.SetUserNamesToAdd(userNamesToAdd)
			}

			// User Ids to remove
			{
				var userNamesToRemove []*string

				for userName, include := range existingUserNamesMap {
					if include {
						userNamesToRemove = append(userNamesToRemove, &userName)
					}
				}

				input.SetUserNamesToRemove(userNamesToRemove)
			}
		}
	}

	var resp *svcsdk.UpdateACLOutput
	_ = resp
	resp, err = rm.sdkapi.UpdateACLWithContext(ctx, input)
	rm.metrics.RecordAPICall("UPDATE", "UpdateACL", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.ACL.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.ACL.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.ACL.Clusters != nil {
		f1 := []*string{}
		for _, f1iter := range resp.ACL.Clusters {
			var f1elem string
			f1elem = *f1iter
			f1 = append(f1, &f1elem)
		}
		ko.Status.Clusters = f1
	} else {
		ko.Status.Clusters = nil
	}
	if resp.ACL.MinimumEngineVersion != nil {
		ko.Status.MinimumEngineVersion = resp.ACL.MinimumEngineVersion
	} else {
		ko.Status.MinimumEngineVersion = nil
	}
	if resp.ACL.Name != nil {
		ko.Spec.Name = resp.ACL.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.ACL.PendingChanges != nil {
		f4 := &svcapitypes.ACLPendingChanges{}
		if resp.ACL.PendingChanges.UserNamesToAdd != nil {
			f4f0 := []*string{}
			for _, f4f0iter := range resp.ACL.PendingChanges.UserNamesToAdd {
				var f4f0elem string
				f4f0elem = *f4f0iter
				f4f0 = append(f4f0, &f4f0elem)
			}
			f4.UserNamesToAdd = f4f0
		}
		if resp.ACL.PendingChanges.UserNamesToRemove != nil {
			f4f1 := []*string{}
			for _, f4f1iter := range resp.ACL.PendingChanges.UserNamesToRemove {
				var f4f1elem string
				f4f1elem = *f4f1iter
				f4f1 = append(f4f1, &f4f1elem)
			}
			f4.UserNamesToRemove = f4f1
		}
		ko.Status.PendingChanges = f4
	} else {
		ko.Status.PendingChanges = nil
	}
	if resp.ACL.Status != nil {
		ko.Status.Status = resp.ACL.Status
	} else {
		ko.Status.Status = nil
	}
	if resp.ACL.UserNames != nil {
		f6 := []*string{}
		for _, f6iter := range resp.ACL.UserNames {
			var f6elem string
			f6elem = *f6iter
			f6 = append(f6, &f6elem)
		}
		ko.Spec.UserNames = f6
	} else {
		ko.Spec.UserNames = nil
	}

	rm.setStatusDefaults(ko)
	ko.Spec.UserNames = desired.ko.Spec.UserNames
	return &resource{ko}, nil
}

// newUpdateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Update API call for the resource
func (rm *resourceManager) newUpdateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.UpdateACLInput, error) {
	res := &svcsdk.UpdateACLInput{}

	if r.ko.Spec.Name != nil {
		res.SetACLName(*r.ko.Spec.Name)
	}

	return res, nil
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer exit(err)
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteACLOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteACLWithContext(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteACL", err)
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteACLInput, error) {
	res := &svcsdk.DeleteACLInput{}

	if r.ko.Spec.Name != nil {
		res.SetACLName(*r.ko.Spec.Name)
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.ACL,
) {
	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if ko.Status.ACKResourceMetadata.Region == nil {
		ko.Status.ACKResourceMetadata.Region = &rm.awsRegion
	}
	if ko.Status.ACKResourceMetadata.OwnerAccountID == nil {
		ko.Status.ACKResourceMetadata.OwnerAccountID = &rm.awsAccountID
	}
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
}

// updateConditions returns updated resource, true; if conditions were updated
// else it returns nil, false
func (rm *resourceManager) updateConditions(
	r *resource,
	onSuccess bool,
	err error,
) (*resource, bool) {
	ko := r.ko.DeepCopy()
	rm.setStatusDefaults(ko)

	// Terminal condition
	var terminalCondition *ackv1alpha1.Condition = nil
	var recoverableCondition *ackv1alpha1.Condition = nil
	var syncCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeTerminal {
			terminalCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeRecoverable {
			recoverableCondition = condition
		}
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			syncCondition = condition
		}
	}

	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound {
			errorMessage = err.Error()
		} else {
			awsErr, _ := ackerr.AWSError(err)
			errorMessage = awsErr.Error()
		}
		terminalCondition.Status = corev1.ConditionTrue
		terminalCondition.Message = &errorMessage
	} else {
		// Clear the terminal condition if no longer present
		if terminalCondition != nil {
			terminalCondition.Status = corev1.ConditionFalse
			terminalCondition.Message = nil
		}
		// Handling Recoverable Conditions
		if err != nil {
			if recoverableCondition == nil {
				// Add a new Condition containing a non-terminal error
				recoverableCondition = &ackv1alpha1.Condition{
					Type: ackv1alpha1.ConditionTypeRecoverable,
				}
				ko.Status.Conditions = append(ko.Status.Conditions, recoverableCondition)
			}
			recoverableCondition.Status = corev1.ConditionTrue
			awsErr, _ := ackerr.AWSError(err)
			errorMessage := err.Error()
			if awsErr != nil {
				errorMessage = awsErr.Error()
			}
			recoverableCondition.Message = &errorMessage
		} else if recoverableCondition != nil {
			recoverableCondition.Status = corev1.ConditionFalse
			recoverableCondition.Message = nil
		}
	}
	// Required to avoid the "declared but not used" error in the default case
	_ = syncCondition
	if terminalCondition != nil || recoverableCondition != nil || syncCondition != nil {
		return &resource{ko}, true // updated
	}
	return nil, false // not updated
}

// terminalAWSError returns awserr, true; if the supplied error is an aws Error type
// and if the exception indicates that it is a Terminal exception
// 'Terminal' exception are specified in generator configuration
func (rm *resourceManager) terminalAWSError(err error) bool {
	if err == nil {
		return false
	}
	awsErr, ok := ackerr.AWSError(err)
	if !ok {
		return false
	}
	switch awsErr.Code() {
	case "ACLAlreadyExistsFault",
		"DefaultUserRequired",
		"UserNotFoundFault",
		"DuplicateUserNameFault",
		"ACLQuotaExceededFault",
		"InvalidParameterValueException",
		"InvalidParameterCombinationException",
		"TagQuotaPerResourceExceeded":
		return true
	default:
		return false
	}
}