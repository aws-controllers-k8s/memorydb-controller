	ko, err = rm.setShardDetails(ctx, desired, ko)
	
	if err != nil {
		return nil, err
	}
	
	// Update the annotations to handle async rollback
	rm.setNodeTypeAnnotation(input.NodeType, ko)
	if input.NumReplicasPerShard != nil {
		rm.setNumReplicasPerShardAnnotation(int64(*input.NumReplicasPerShard), ko)
	}

	if input.NumShards != nil {
		rm.setNumShardAnnotation(int64(*input.NumShards), ko)
	}