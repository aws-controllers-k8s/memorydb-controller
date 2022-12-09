    var subnets []*string
    for _, subnet := range resp.SubnetGroups[0].Subnets {
	    if  subnet.Identifier != nil{
		    subnets = append(subnets, subnet.Identifier)
	    }
    }

    ko.Spec.SubnetIDs = subnets