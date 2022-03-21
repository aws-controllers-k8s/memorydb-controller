var resourceSyncedCondition *ackv1alpha1.Condition = nil
for _, condition := range ko.Status.Conditions {
	if condition.Type == ackv1alpha1.ConditionTypeResourceSynced {
		resourceSyncedCondition = condition
		break
	}
}
if resourceSyncedCondition == nil {
	resourceSyncedCondition = &ackv1alpha1.Condition{
		Type:   ackv1alpha1.ConditionTypeResourceSynced,
		Status: corev1.ConditionTrue,
	}
	ko.Status.Conditions = append(ko.Status.Conditions, resourceSyncedCondition)
} else {
	resourceSyncedCondition.Status = corev1.ConditionTrue
}