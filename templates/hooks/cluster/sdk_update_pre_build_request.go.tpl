	res, err := rm.validateClusterNeedsUpdate(desired, latest, delta)
	
	if err != nil || res!= nil{
		return res, err
	}