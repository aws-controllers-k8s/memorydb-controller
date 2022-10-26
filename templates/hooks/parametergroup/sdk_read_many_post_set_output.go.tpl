resourceARN := (*string)(ko.Status.ACKResourceMetadata.ARN)
tags, err := rm.getTags(ctx, *resourceARN)
if err != nil {
	return nil, err
}
ko.Spec.Tags = tags


ko, err = rm.setParameters(ctx, ko)

if err != nil {
    return nil, err
}