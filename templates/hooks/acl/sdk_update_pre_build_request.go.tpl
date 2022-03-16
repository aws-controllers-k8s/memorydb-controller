validationErr := rm.validateACLNeedsUpdate(latest)

if validationErr != nil {
	return nil, err
}