	ko, err = rm.setShardDetails(ctx, desired, ko)
	
	if err != nil {
		return nil, err
	}
	
	// Update the annotations to handle async rollback
	if input.NodeType != nil {
		rm.setNodeTypeAnnotation(input.NodeType, ko)
	}
	if input.NumReplicasPerShard != nil {
		rm.setNumReplicasPerShardAnnotation(*input.NumReplicasPerShard, ko)
	}

	if input.NumShards != nil {
		rm.setNumShardAnnotation(*input.NumShards, ko)
	}