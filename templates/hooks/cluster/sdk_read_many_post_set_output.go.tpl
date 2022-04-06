	cluster := resp.Clusters[0]
	if cluster.NumberOfShards != nil {
		ko.Spec.NumShards = cluster.NumberOfShards
	} else {
		ko.Spec.NumShards = nil
	}

	if cluster.Shards != nil && cluster.Shards[0].NumberOfNodes != nil {
		replicas := *cluster.Shards[0].NumberOfNodes - 1
		ko.Spec.NumReplicasPerShard = &replicas
	} else {
		ko.Spec.NumReplicasPerShard = nil
	}

	if cluster.SecurityGroups != nil {
		var securityGroupIds []*string
		for _, securityGroup := range cluster.SecurityGroups {
			if securityGroup.SecurityGroupId != nil {
				securityGroupIds = append(securityGroupIds, securityGroup.SecurityGroupId)
			}
		}
		ko.Spec.SecurityGroupIDs = securityGroupIds
	} else {
		ko.Spec.SecurityGroupIDs = nil
	}

	respErr := rm.setAllowedNodeTypeUpdates(ctx, ko)
	if respErr != nil {
		return nil, respErr
	}