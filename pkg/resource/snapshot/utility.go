package snapshot

import (
	"errors"
	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
	ackv1alpha1 "github.com/aws-controllers-k8s/runtime/apis/core/v1alpha1"
	ackrequeue "github.com/aws-controllers-k8s/runtime/pkg/requeue"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var (
	condMsgCurrentlyDeleting string = "snapshot currently being deleted"
	deleteStatus             string = "deleting"
)

var (
	requeueWaitWhileDeleting = ackrequeue.NeededAfter(
		errors.New("delete is in progress"),
		ackrequeue.DefaultRequeueAfterDuration,
	)
)

// isDeleting returns true if supplied snapshot resource state is 'deleting'
func isDeleting(r *resource) bool {
	if r == nil || r.ko.Status.Status == nil {
		return false
	}
	status := *r.ko.Status.Status
	return status == deleteStatus
}

func (rm *resourceManager) setSnapshotOutput(
	r *resource,
	obj *svcsdk.Snapshot,
) (*resource, error) {
	if obj == nil ||
		r == nil ||
		r.ko == nil {
		return nil, nil
	}
	resp := &svcsdk.DeleteSnapshotOutput{Snapshot: obj}

	// Merge in the information we read from the API call above to the copy of
	// the original Kubernetes object we passed to the function
	ko := r.ko.DeepCopy()

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
			f1.NumShards = resp.Snapshot.ClusterConfiguration.NumShards
		}
		if resp.Snapshot.ClusterConfiguration.ParameterGroupName != nil {
			f1.ParameterGroupName = resp.Snapshot.ClusterConfiguration.ParameterGroupName
		}
		if resp.Snapshot.ClusterConfiguration.Port != nil {
			f1.Port = resp.Snapshot.ClusterConfiguration.Port
		}
		if resp.Snapshot.ClusterConfiguration.Shards != nil {
			f1f8 := []*svcapitypes.ShardDetail{}
			for _, f1f8iter := range resp.Snapshot.ClusterConfiguration.Shards {
				f1f8elem := &svcapitypes.ShardDetail{}
				if f1f8iter.Configuration != nil {
					f1f8elemf0 := &svcapitypes.ShardConfiguration{}
					if f1f8iter.Configuration.ReplicaCount != nil {
						f1f8elemf0.ReplicaCount = f1f8iter.Configuration.ReplicaCount
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
			f1.SnapshotRetentionLimit = resp.Snapshot.ClusterConfiguration.SnapshotRetentionLimit
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
	rm.customSetOutput(obj, ko)
	return &resource{ko}, nil
}
