package util

func Int32OrNil(val *int64) *int32 {
	if val == nil {
		return nil
	}
	ret := int32(*val)
	return &ret
}

func Int64OrNil(val *int32) *int64 {
	if val == nil {
		return nil
	}
	ret := int64(*val)
	return &ret
}
