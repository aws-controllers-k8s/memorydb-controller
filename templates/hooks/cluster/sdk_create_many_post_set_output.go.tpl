	ko, err = rm.setShardDetails(ctx, desired, ko)
	
	if err != nil {
		return nil, err
	}
	
	// Update the annotations to handle async rollback
	rm.setNodeTypeAnnotation(input.NodeType, ko)
	rm.setNumReplicasPerShardAnnotation(input.NumReplicasPerShard, ko)
	rm.setNumShardAnnotation(input.NumShards, ko)