res, err := rm.validateUserNeedsUpdate(desired, latest, delta)

if err != nil || res!= nil{
	return res, err
}