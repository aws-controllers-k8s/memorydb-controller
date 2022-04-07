    // custom set output from response
    ko, err = rm.customDescribeSnapshotSetOutput(resp, ko)
    if err != nil {
    	return nil, err
    }