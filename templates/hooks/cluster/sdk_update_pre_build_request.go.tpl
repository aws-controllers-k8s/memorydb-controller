	res, err := rm.validateClusterNeedsUpdate(desired, latest, delta)
	
	if err != nil || res!= nil{
		return res, err
	}

	if delta.DifferentAt("Spec.Tags") {
        err = rm.updateTags(ctx, desired, latest)
        if err != nil {
            return nil, err
        }
    }

    if !delta.DifferentExcept("Spec.Tags") {
    	return desired, nil
    }
