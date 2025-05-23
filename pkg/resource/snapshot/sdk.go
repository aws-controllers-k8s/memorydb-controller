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

package snapshot

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackcompare "github.com/aws-controllers-k8s/runtime/pkg/compare"
	ackcondition "github.com/aws-controllers-k8s/runtime/pkg/condition"
	ackerr "github.com/aws-controllers-k8s/runtime/pkg/errors"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	ackrtlog "github.com/aws-controllers-k8s/runtime/pkg/runtime/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	svcsdk "github.com/aws/aws-sdk-go-v2/service/memorydb"
	svcsdktypes "github.com/aws/aws-sdk-go-v2/service/memorydb/types"
	smithy "github.com/aws/smithy-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

// Hack to avoid import errors during build...
var (
	_ = &metav1.Time{}
	_ = strings.ToLower("")
	_ = &svcsdk.Client{}
	_ = &svcapitypes.Snapshot{}
	_ = ackv1alpha1.AWSAccountID("")
	_ = &ackerr.NotFound
	_ = &ackcondition.NotManagedMessage
	_ = &reflect.Value{}
	_ = fmt.Sprintf("")
	_ = &ackrequeue.NoRequeue{}
	_ = &aws.Config{}
)

// sdkFind returns SDK-specific information about a supplied resource
func (rm *resourceManager) sdkFind(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkFind")
	defer func() {
		exit(err)
	}()
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
	var resp *svcsdk.DescribeSnapshotsOutput
	resp, err = rm.sdkapi.DescribeSnapshots(ctx, input)
	rm.metrics.RecordAPICall("READ_MANY", "DescribeSnapshots", err)
	if err != nil {
		var awsErr smithy.APIError
		if errors.As(err, &awsErr) && awsErr.ErrorCode() == "SnapshotNotFoundFault" {
			return nil, ackerr.NotFound
		}
		return nil, err
	}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

	found := false
	for _, elem := range resp.Snapshots {
		if elem.ARN != nil {
			if ko.Status.ACKResourceMetadata == nil {
				ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
			}
			tmpARN := ackv1alpha1.AWSResourceName(*elem.ARN)
			ko.Status.ACKResourceMetadata.ARN = &tmpARN
		}
		if elem.ClusterConfiguration != nil {
			f1 := &svcapitypes.ClusterConfiguration{}
			if elem.ClusterConfiguration.Description != nil {
				f1.Description = elem.ClusterConfiguration.Description
			}
			if elem.ClusterConfiguration.EngineVersion != nil {
				f1.EngineVersion = elem.ClusterConfiguration.EngineVersion
			}
			if elem.ClusterConfiguration.MaintenanceWindow != nil {
				f1.MaintenanceWindow = elem.ClusterConfiguration.MaintenanceWindow
			}
			if elem.ClusterConfiguration.Name != nil {
				f1.Name = elem.ClusterConfiguration.Name
			}
			if elem.ClusterConfiguration.NodeType != nil {
				f1.NodeType = elem.ClusterConfiguration.NodeType
			}
			if elem.ClusterConfiguration.NumShards != nil {
				numShardsCopy := int64(*elem.ClusterConfiguration.NumShards)
				f1.NumShards = &numShardsCopy
			}
			if elem.ClusterConfiguration.ParameterGroupName != nil {
				f1.ParameterGroupName = elem.ClusterConfiguration.ParameterGroupName
			}
			if elem.ClusterConfiguration.Port != nil {
				portCopy := int64(*elem.ClusterConfiguration.Port)
				f1.Port = &portCopy
			}
			if elem.ClusterConfiguration.Shards != nil {
				f1f8 := []*svcapitypes.ShardDetail{}
				for _, f1f8iter := range elem.ClusterConfiguration.Shards {
					f1f8elem := &svcapitypes.ShardDetail{}
					if f1f8iter.Configuration != nil {
						f1f8elemf0 := &svcapitypes.ShardConfiguration{}
						if f1f8iter.Configuration.ReplicaCount != nil {
							replicaCountCopy := int64(*f1f8iter.Configuration.ReplicaCount)
							f1f8elemf0.ReplicaCount = &replicaCountCopy
						}
						if f1f8iter.Configuration.Slots != nil {
							f1f8elemf0.Slots = f1f8iter.Configuration.Slots
						}
						f1f8elem.Configuration = f1f8elemf0
					}
					if f1f8iter.Name != nil {
						f1f8elem.Name = f1f8iter.Name
					}
					if f1f8iter.Size != nil {
						f1f8elem.Size = f1f8iter.Size
					}
					if f1f8iter.SnapshotCreationTime != nil {
						f1f8elem.SnapshotCreationTime = &metav1.Time{*f1f8iter.SnapshotCreationTime}
					}
					f1f8 = append(f1f8, f1f8elem)
				}
				f1.Shards = f1f8
			}
			if elem.ClusterConfiguration.SnapshotRetentionLimit != nil {
				snapshotRetentionLimitCopy := int64(*elem.ClusterConfiguration.SnapshotRetentionLimit)
				f1.SnapshotRetentionLimit = &snapshotRetentionLimitCopy
			}
			if elem.ClusterConfiguration.SnapshotWindow != nil {
				f1.SnapshotWindow = elem.ClusterConfiguration.SnapshotWindow
			}
			if elem.ClusterConfiguration.SubnetGroupName != nil {
				f1.SubnetGroupName = elem.ClusterConfiguration.SubnetGroupName
			}
			if elem.ClusterConfiguration.TopicArn != nil {
				f1.TopicARN = elem.ClusterConfiguration.TopicArn
			}
			if elem.ClusterConfiguration.VpcId != nil {
				f1.VPCID = elem.ClusterConfiguration.VpcId
			}
			ko.Status.ClusterConfiguration = f1
		} else {
			ko.Status.ClusterConfiguration = nil
		}
		if elem.KmsKeyId != nil {
			ko.Spec.KMSKeyID = elem.KmsKeyId
		} else {
			ko.Spec.KMSKeyID = nil
		}
		if elem.Name != nil {
			ko.Spec.Name = elem.Name
		} else {
			ko.Spec.Name = nil
		}
		if elem.Source != nil {
			ko.Status.Source = elem.Source
		} else {
			ko.Status.Source = nil
		}
		if elem.Status != nil {
			ko.Status.Status = elem.Status
		} else {
			ko.Status.Status = nil
		}
		found = true
		break
	}
	if !found {
		return nil, ackerr.NotFound
	}

	rm.setStatusDefaults(ko)
	// custom set output from response
	ko, err = rm.customDescribeSnapshotSetOutput(resp, ko)
	if err != nil {
		return nil, err
	}

	if rm.isSnapshotAvailable(&resource{ko}) {
		resourceARN := (*string)(ko.Status.ACKResourceMetadata.ARN)
		tags, err := rm.getTags(ctx, *resourceARN)
		if err != nil {
			return nil, err
		}
		ko.Spec.Tags = tags
	}

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
) (*svcsdk.DescribeSnapshotsInput, error) {
	res := &svcsdk.DescribeSnapshotsInput{}

	if r.ko.Spec.Name != nil {
		res.SnapshotName = r.ko.Spec.Name
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
	defer func() {
		exit(err)
	}()
	created, err = rm.customTryCopySnapshot(ctx, desired)
	if created != nil || err != nil {
		return created, err
	}
	input, err := rm.newCreateRequestPayload(ctx, desired)
	if err != nil {
		return nil, err
	}

	var resp *svcsdk.CreateSnapshotOutput
	_ = resp
	resp, err = rm.sdkapi.CreateSnapshot(ctx, input)
	rm.metrics.RecordAPICall("CREATE", "CreateSnapshot", err)
	if err != nil {
		return nil, err
	}
	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := desired.ko.DeepCopy()

	if ko.Status.ACKResourceMetadata == nil {
		ko.Status.ACKResourceMetadata = &ackv1alpha1.ResourceMetadata{}
	}
	if resp.Snapshot.ARN != nil {
		arn := ackv1alpha1.AWSResourceName(*resp.Snapshot.ARN)
		ko.Status.ACKResourceMetadata.ARN = &arn
	}
	if resp.Snapshot.ClusterConfiguration != nil {
		f1 := &svcapitypes.ClusterConfiguration{}
		if resp.Snapshot.ClusterConfiguration.Description != nil {
			f1.Description = resp.Snapshot.ClusterConfiguration.Description
		}
		if resp.Snapshot.ClusterConfiguration.EngineVersion != nil {
			f1.EngineVersion = resp.Snapshot.ClusterConfiguration.EngineVersion
		}
		if resp.Snapshot.ClusterConfiguration.MaintenanceWindow != nil {
			f1.MaintenanceWindow = resp.Snapshot.ClusterConfiguration.MaintenanceWindow
		}
		if resp.Snapshot.ClusterConfiguration.Name != nil {
			f1.Name = resp.Snapshot.ClusterConfiguration.Name
		}
		if resp.Snapshot.ClusterConfiguration.NodeType != nil {
			f1.NodeType = resp.Snapshot.ClusterConfiguration.NodeType
		}
		if resp.Snapshot.ClusterConfiguration.NumShards != nil {
			numShardsCopy := int64(*resp.Snapshot.ClusterConfiguration.NumShards)
			f1.NumShards = &numShardsCopy
		}
		if resp.Snapshot.ClusterConfiguration.ParameterGroupName != nil {
			f1.ParameterGroupName = resp.Snapshot.ClusterConfiguration.ParameterGroupName
		}
		if resp.Snapshot.ClusterConfiguration.Port != nil {
			portCopy := int64(*resp.Snapshot.ClusterConfiguration.Port)
			f1.Port = &portCopy
		}
		if resp.Snapshot.ClusterConfiguration.Shards != nil {
			f1f8 := []*svcapitypes.ShardDetail{}
			for _, f1f8iter := range resp.Snapshot.ClusterConfiguration.Shards {
				f1f8elem := &svcapitypes.ShardDetail{}
				if f1f8iter.Configuration != nil {
					f1f8elemf0 := &svcapitypes.ShardConfiguration{}
					if f1f8iter.Configuration.ReplicaCount != nil {
						replicaCountCopy := int64(*f1f8iter.Configuration.ReplicaCount)
						f1f8elemf0.ReplicaCount = &replicaCountCopy
					}
					if f1f8iter.Configuration.Slots != nil {
						f1f8elemf0.Slots = f1f8iter.Configuration.Slots
					}
					f1f8elem.Configuration = f1f8elemf0
				}
				if f1f8iter.Name != nil {
					f1f8elem.Name = f1f8iter.Name
				}
				if f1f8iter.Size != nil {
					f1f8elem.Size = f1f8iter.Size
				}
				if f1f8iter.SnapshotCreationTime != nil {
					f1f8elem.SnapshotCreationTime = &metav1.Time{*f1f8iter.SnapshotCreationTime}
				}
				f1f8 = append(f1f8, f1f8elem)
			}
			f1.Shards = f1f8
		}
		if resp.Snapshot.ClusterConfiguration.SnapshotRetentionLimit != nil {
			snapshotRetentionLimitCopy := int64(*resp.Snapshot.ClusterConfiguration.SnapshotRetentionLimit)
			f1.SnapshotRetentionLimit = &snapshotRetentionLimitCopy
		}
		if resp.Snapshot.ClusterConfiguration.SnapshotWindow != nil {
			f1.SnapshotWindow = resp.Snapshot.ClusterConfiguration.SnapshotWindow
		}
		if resp.Snapshot.ClusterConfiguration.SubnetGroupName != nil {
			f1.SubnetGroupName = resp.Snapshot.ClusterConfiguration.SubnetGroupName
		}
		if resp.Snapshot.ClusterConfiguration.TopicArn != nil {
			f1.TopicARN = resp.Snapshot.ClusterConfiguration.TopicArn
		}
		if resp.Snapshot.ClusterConfiguration.VpcId != nil {
			f1.VPCID = resp.Snapshot.ClusterConfiguration.VpcId
		}
		ko.Status.ClusterConfiguration = f1
	} else {
		ko.Status.ClusterConfiguration = nil
	}
	if resp.Snapshot.KmsKeyId != nil {
		ko.Spec.KMSKeyID = resp.Snapshot.KmsKeyId
	} else {
		ko.Spec.KMSKeyID = nil
	}
	if resp.Snapshot.Name != nil {
		ko.Spec.Name = resp.Snapshot.Name
	} else {
		ko.Spec.Name = nil
	}
	if resp.Snapshot.Source != nil {
		ko.Status.Source = resp.Snapshot.Source
	} else {
		ko.Status.Source = nil
	}
	if resp.Snapshot.Status != nil {
		ko.Status.Status = resp.Snapshot.Status
	} else {
		ko.Status.Status = nil
	}

	rm.setStatusDefaults(ko)
	// custom set output from response
	ko, err = rm.customCreateSnapshotSetOutput(resp, ko)
	if err != nil {
		return nil, err
	}
	return &resource{ko}, nil
}

// newCreateRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Create API call for the resource
func (rm *resourceManager) newCreateRequestPayload(
	ctx context.Context,
	r *resource,
) (*svcsdk.CreateSnapshotInput, error) {
	res := &svcsdk.CreateSnapshotInput{}

	if r.ko.Spec.ClusterName != nil {
		res.ClusterName = r.ko.Spec.ClusterName
	}
	if r.ko.Spec.KMSKeyID != nil {
		res.KmsKeyId = r.ko.Spec.KMSKeyID
	}
	if r.ko.Spec.Name != nil {
		res.SnapshotName = r.ko.Spec.Name
	}
	if r.ko.Spec.Tags != nil {
		f3 := []svcsdktypes.Tag{}
		for _, f3iter := range r.ko.Spec.Tags {
			f3elem := &svcsdktypes.Tag{}
			if f3iter.Key != nil {
				f3elem.Key = f3iter.Key
			}
			if f3iter.Value != nil {
				f3elem.Value = f3iter.Value
			}
			f3 = append(f3, *f3elem)
		}
		res.Tags = f3
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
) (*resource, error) {
	return rm.customUpdate(ctx, desired, latest, delta)
}

// sdkDelete deletes the supplied resource in the backend AWS service API
func (rm *resourceManager) sdkDelete(
	ctx context.Context,
	r *resource,
) (latest *resource, err error) {
	rlog := ackrtlog.FromContext(ctx)
	exit := rlog.Trace("rm.sdkDelete")
	defer func() {
		exit(err)
	}()
	if isDeleting(r) {
		// Setting resource synced condition to false will trigger a requeue of
		// the resource.
		ackcondition.SetSynced(
			r,
			corev1.ConditionFalse,
			&condMsgCurrentlyDeleting,
			nil,
		)
		// Need to return a requeue error here, otherwise:
		// - reconciler.deleteResource() marks the resource unmanaged
		// - reconciler.HandleReconcileError() does not update status for unmanaged resource
		// - reconciler.handleRequeues() is not invoked for delete code path.
		// TODO: return err as nil when reconciler is updated.
		return r, requeueWaitWhileDeleting
	}

	if isNotReadyForDeleting(r) {
		// return a requeue if snapshot is not ready to be deleted
		return r, requeueWaitSnapshotIsReadyForDeleting
	}
	input, err := rm.newDeleteRequestPayload(r)
	if err != nil {
		return nil, err
	}
	var resp *svcsdk.DeleteSnapshotOutput
	_ = resp
	resp, err = rm.sdkapi.DeleteSnapshot(ctx, input)
	rm.metrics.RecordAPICall("DELETE", "DeleteSnapshot", err)
	if err == nil {
		rp, _ := rm.setSnapshotOutput(r, resp.Snapshot)
		// Setting resource synced condition to false will trigger a requeue of
		// the resource.
		ackcondition.SetSynced(
			rp,
			corev1.ConditionFalse,
			&condMsgCurrentlyDeleting,
			nil,
		)
		// Need to return a requeue error here, otherwise:
		// - reconciler.deleteResource() marks the resource unmanaged
		// - reconciler.HandleReconcileError() does not update status for unmanaged resource
		// - reconciler.handleRequeues() is not invoked for delete code path.
		// TODO: return err as nil when reconciler is updated.
		return rp, requeueWaitWhileDeleting
	}
	return nil, err
}

// newDeleteRequestPayload returns an SDK-specific struct for the HTTP request
// payload of the Delete API call for the resource
func (rm *resourceManager) newDeleteRequestPayload(
	r *resource,
) (*svcsdk.DeleteSnapshotInput, error) {
	res := &svcsdk.DeleteSnapshotInput{}

	if r.ko.Spec.Name != nil {
		res.SnapshotName = r.ko.Spec.Name
	}

	return res, nil
}

// setStatusDefaults sets default properties into supplied custom resource
func (rm *resourceManager) setStatusDefaults(
	ko *svcapitypes.Snapshot,
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
	var termError *ackerr.TerminalError
	if rm.terminalAWSError(err) || err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
		if terminalCondition == nil {
			terminalCondition = &ackv1alpha1.Condition{
				Type: ackv1alpha1.ConditionTypeTerminal,
			}
			ko.Status.Conditions = append(ko.Status.Conditions, terminalCondition)
		}
		var errorMessage = ""
		if err == ackerr.SecretTypeNotSupported || err == ackerr.SecretNotFound || errors.As(err, &termError) {
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

	var terminalErr smithy.APIError
	if !errors.As(err, &terminalErr) {
		return false
	}
	switch terminalErr.ErrorCode() {
	case "InvalidParameterCombinationException",
		"InvalidParameterValueException",
		"InvalidParameter",
		"SnapshotAlreadyExistsFault":
		return true
	default:
		return false
	}
}
