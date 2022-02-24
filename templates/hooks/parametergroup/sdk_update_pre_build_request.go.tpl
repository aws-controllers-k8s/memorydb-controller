if delta.DifferentAt("Spec.ParameterNameValues") {
    ko, err := rm.resetParameterGroup(ctx, desired, latest)

    if ko != nil || err != nil {
        return ko, err
    }
}