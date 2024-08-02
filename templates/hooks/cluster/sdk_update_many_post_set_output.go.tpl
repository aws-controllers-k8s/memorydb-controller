	ko, err = rm.setShardDetails(ctx, desired, ko)
	
	if err != nil {
		return nil, err
	}
	
	// Do not perform spec patching as these fields eventually get updated
	ko.Spec.NumShards = desired.ko.Spec.NumShards
	ko.Spec.NumReplicasPerShard = desired.ko.Spec.NumReplicasPerShard
	ko.Spec.ACLName = desired.ko.Spec.ACLName
	ko.Spec.NodeType = desired.ko.Spec.NodeType
	ko.Spec.EngineVersion = desired.ko.Spec.EngineVersion
	
	// Update the annotations to handle async rollback
	rm.setNodeTypeAnnotation(input.NodeType, ko)
	if input.ReplicaConfiguration != nil && input.ReplicaConfiguration.ReplicaCount != nil {
		rm.setNumReplicasPerShardAnnotation(*input.ReplicaConfiguration.ReplicaCount, ko)
	}
	if input.ShardConfiguration != nil && input.ShardConfiguration.ShardCount != nil {
		rm.setNumShardAnnotation(*input.ShardConfiguration.ShardCount, ko)
	}
	return &resource{ko}, requeueWaitWhileUpdating
