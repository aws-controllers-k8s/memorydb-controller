    createMapForUserNames := func(userIds []*string) map[string]bool {
	    userIdsMap := make(map[string]bool)

	    for _, userId := range userIds {
		    userIdsMap[*userId] = true
	    }

	    return userIdsMap
    }

    for _, diff := range delta.Differences {
	    if diff.Path.Contains("Spec.UserNames") {
		    existingUserNamesMap := createMapForUserNames(diff.B.([]*string))
		    requiredUserNamesMap := createMapForUserNames(diff.A.([]*string))

		    // If a user ID is not required to be deleted or added set its value as false
		    for userName, _ := range existingUserNamesMap {
			    if _, ok := requiredUserNamesMap[userName]; ok {
				    requiredUserNamesMap[userName] = false
				    existingUserNamesMap[userName] = false
			    }
		    }

		    if err != nil {
			    return nil, err
		    }

		    // User Ids to add
		    {
			    var userNamesToAdd []*string

			    for userName, include := range requiredUserNamesMap {
				    if include {
					    userNamesToAdd = append(userNamesToAdd, &userName)
				    }
			    }

			    input.SetUserNamesToAdd(userNamesToAdd)
		    }

		    // User Ids to remove
		    {
			    var userNamesToRemove []*string

			    for userName, include := range existingUserNamesMap {
				    if include {
					    userNamesToRemove = append(userNamesToRemove, &userName)
				    }
			    }

			    input.SetUserNamesToRemove(userNamesToRemove)
		    }
	    }
    }
