if delta.DifferentAt("Spec.ParameterNameValues") {
    ko, err := rm.resetParameterGroup(ctx, desired, latest)
    if ko != nil || err != nil {
        return ko, err
    }
}

if delta.DifferentAt("Spec.Tags") {
    err = rm.updateTags(ctx, desired, latest)
    if err != nil {
        return nil, err
    }
}
