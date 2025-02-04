	cluster := resp.Clusters[0]
	if cluster.NumberOfShards != nil {
		ko.Spec.NumShards = aws.Int64(int64(*cluster.NumberOfShards))
	} else {
		ko.Spec.NumShards = nil
	}

	if cluster.Shards != nil && cluster.Shards[0].NumberOfNodes != nil {
		replicas := *cluster.Shards[0].NumberOfNodes - 1
		ko.Spec.NumReplicasPerShard = aws.Int64(int64(replicas))
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

    ko.Status.Events, err = rm.getEvents(ctx, r)
    if err != nil {
        return nil, err
    }

    if rm.isClusterAvailable(&resource{ko}) {
		resourceARN := (*string)(ko.Status.ACKResourceMetadata.ARN)
		tags, err := rm.getTags(ctx, *resourceARN)
		if err != nil {
			return nil, err
		}
		ko.Spec.Tags = tags
	}

