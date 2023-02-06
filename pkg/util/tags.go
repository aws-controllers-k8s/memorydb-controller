package util

import (
	ackutil "github.com/aws-controllers-k8s/runtime/pkg/util"
	svcsdk "github.com/aws/aws-sdk-go/service/memorydb"

	svcapitypes "github.com/aws-controllers-k8s/memorydb-controller/apis/v1alpha1"
)

func ComputeTagsDelta(
	desired []*svcapitypes.Tag,
	latest []*svcapitypes.Tag,
) (addedOrUpdated []*svcapitypes.Tag, removed []*string) {
	var visitedIndexes []string
	var hasSameKey bool

	for _, latestElement := range latest {
		hasSameKey = false
		visitedIndexes = append(visitedIndexes, *latestElement.Key)
		for _, desiredElement := range desired {
			if equalStrings(latestElement.Key, desiredElement.Key) {
				hasSameKey = true
				if !equalStrings(latestElement.Value, desiredElement.Value) {
					addedOrUpdated = append(addedOrUpdated, desiredElement)
				}
				break
			}
		}
		if hasSameKey {
			continue
		}
		removed = append(removed, latestElement.Key)
	}
	for _, desiredElement := range desired {
		if !ackutil.InStrings(*desiredElement.Key, visitedIndexes) {
			addedOrUpdated = append(addedOrUpdated, desiredElement)
		}
	}
	return addedOrUpdated, removed
}

func ResourceTagsFromSDKTags(
	sdkTags []*svcsdk.Tag,
) []*svcapitypes.Tag {
	tags := make([]*svcapitypes.Tag, len(sdkTags))
	for i := range sdkTags {
		tags[i] = &svcapitypes.Tag{
			Key:   sdkTags[i].Key,
			Value: sdkTags[i].Value,
		}
	}
	return tags
}

func SDKTagsFromResourceTags(
	rTags []*svcapitypes.Tag,
) []*svcsdk.Tag {
	tags := make([]*svcsdk.Tag, len(rTags))
	for i := range rTags {
		tags[i] = &svcsdk.Tag{
			Key:   rTags[i].Key,
			Value: rTags[i].Value,
		}
	}
	return tags
}

func equalStrings(a, b *string) bool {
	if a == nil {
		return b == nil || *b == ""
	}
	return (*a == "" && b == nil) || *a == *b
}
