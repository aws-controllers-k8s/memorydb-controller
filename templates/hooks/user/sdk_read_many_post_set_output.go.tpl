    ko.Status.Events, err = rm.getEvents(ctx, r)
    if err != nil {
        return nil, err
    }

    if rm.isUserActive(&resource{ko}) {
		resourceARN := (*string)(ko.Status.ACKResourceMetadata.ARN)
		tags, err := rm.getTags(ctx, *resourceARN)
		if err != nil {
			return nil, err
		}
		ko.Spec.Tags = tags
	}
