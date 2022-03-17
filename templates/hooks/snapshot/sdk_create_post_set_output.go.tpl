    // custom set output from response
    ko, err = rm.customCreateSnapshotSetOutput(resp, ko)
    if err != nil {
	    return nil, err
    }