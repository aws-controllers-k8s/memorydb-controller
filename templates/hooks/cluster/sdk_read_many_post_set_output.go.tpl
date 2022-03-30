	if resp.Clusters[0].NumberOfShards != nil {
		ko.Spec.NumShards = resp.Clusters[0].NumberOfShards
	} else {
		ko.Spec.NumShards = nil
	}

	if resp.Clusters[0].Shards != nil && resp.Clusters[0].Shards[0].NumberOfNodes != nil {
		replicas := *resp.Clusters[0].Shards[0].NumberOfNodes - 1
		ko.Spec.NumReplicasPerShard = &replicas
	} else {
		ko.Spec.NumReplicasPerShard = nil
	}

	if resp.Clusters[0].SecurityGroups != nil {
		var securityGroupIds []*string
		for _, securityGroup := range resp.Clusters[0].SecurityGroups {
			if securityGroup.SecurityGroupId != nil {
				securityGroupIds = append(securityGroupIds, securityGroup.SecurityGroupId)
			}
		}
		ko.Spec.SecurityGroupIDs = securityGroupIds
	} else {
		ko.Spec.SecurityGroupIDs = nil
	}

	rm.setAllowedNodeTypeUpdates(ctx, ko)