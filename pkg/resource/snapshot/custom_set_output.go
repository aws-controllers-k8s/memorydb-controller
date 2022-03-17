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

package snapshot

import (
	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	"github.com/aws/aws-sdk-go/service/memorydb"
	corev1 "k8s.io/api/core/v1"
)

func (rm *resourceManager) customDescribeSnapshotSetOutput(
	resp *memorydb.DescribeSnapshotsOutput,
	ko *svcapitypes.Snapshot,
) (*svcapitypes.Snapshot, error) {
	if len(resp.Snapshots) == 0 {
		return ko, nil
	}
	elem := resp.Snapshots[0]
	rm.customSetOutput(elem, ko)
	return ko, nil
}

func (rm *resourceManager) customCreateSnapshotSetOutput(
	resp *memorydb.CreateSnapshotOutput,
	ko *svcapitypes.Snapshot,
) (*svcapitypes.Snapshot, error) {
	rm.customSetOutput(resp.Snapshot, ko)
	return ko, nil
}

func (rm *resourceManager) customCopySnapshotSetOutput(
	resp *memorydb.CopySnapshotOutput,
	ko *svcapitypes.Snapshot,
) *svcapitypes.Snapshot {
	rm.customSetOutput(resp.Snapshot, ko)
	return ko
}

func (rm *resourceManager) customSetOutput(
	respSnapshot *memorydb.Snapshot,
	ko *svcapitypes.Snapshot,
) {
	if ko.Status.Conditions == nil {
		ko.Status.Conditions = []*ackv1alpha1.Condition{}
	}
	snapshotStatus := respSnapshot.Status
	syncConditionStatus := corev1.ConditionUnknown
	if snapshotStatus != nil {
		if *snapshotStatus == "available" ||
			*snapshotStatus == "failed" {
			syncConditionStatus = corev1.ConditionTrue
		} else {
			// resource in "creating", "restoring","exporting"
			syncConditionStatus = corev1.ConditionFalse
		}
	}
	var resourceSyncedCondition *ackv1alpha1.Condition = nil
	for _, condition := range ko.Status.Conditions {
		if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
			resourceSyncedCondition = condition
			break
		}
	}
	if resourceSyncedCondition == nil {
		resourceSyncedCondition = &ackv1alpha1.Condition{
			Type:   ackv1alpha1.ConditionTypeResourceSynced,
			Status: syncConditionStatus,
		}
		ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
	} else {
		resourceSyncedCondition.Status = syncConditionStatus
	}
}
