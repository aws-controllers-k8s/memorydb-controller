    created, err = rm.customTryCopySnapshot(ctx, desired)
    if created != nil || err != nil {
	    return created, err
    }