	if desired.ko.Spec.AuthenticationMode.Passwords != nil {
		f1f0 := []string{}
		for _, f1f0iter := range desired.ko.Spec.AuthenticationMode.Passwords {
			var f1f0elem string
			if f1f0iter != nil {
				tmpSecret, err := rm.rr.SecretValueFromReference(ctx, f1f0iter)
				if err != nil {
					return nil, ackrequeue.Needed(err)
				}
				if tmpSecret != "" {
					f1f0elem = tmpSecret
				}
			}
			f1f0 = append(f1f0, f1f0elem)
		}
		input.AuthenticationMode.Passwords = f1f0
	}