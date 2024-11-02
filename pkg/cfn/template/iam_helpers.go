package template

// AttachAllowPolicy constructs a role with allow policy for given resources and actions
func (t *Template) AttachAllowPolicy(name string, refRole *Value, resources interface{}, actions []string) {
	t.NewResource(name, &IAMPolicy{
		PolicyName: MakeName(name),
		Roles:      MakeSlice(refRole),
		PolicyDocument: MakePolicyDocument(MapOfInterfaces{
			"Effect":   "Allow",
			"Resource": resources,
			"Action":   actions,
		}),
	})
}

// MakePolicyDocument constructs a policy with given statements
func MakePolicyDocument(statements ...MapOfInterfaces) MapOfInterfaces {
	return MapOfInterfaces{
		"Version":   "2012-10-17",
		"Statement": statements,
	}
}

// MakeAssumeRolePolicyDocumentForServices constructs a trust policy for given services
func MakeAssumeRolePolicyDocumentForServices(services ...string) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": []string{"sts:AssumeRole"},
		"Principal": map[string][]string{
			"Service": services,
		},
	})
}

// MakeAssumeRoleWithWebIdentityPolicyDocument constructs a trust policy for given a web identity priovider with given conditions
func MakeAssumeRoleWithWebIdentityPolicyDocument(providerARN string, condition MapOfInterfaces) MapOfInterfaces {
	return MakePolicyDocument(MapOfInterfaces{
		"Effect": "Allow",
		"Action": []string{"sts:AssumeRoleWithWebIdentity"},
		"Principal": map[string]string{
			"Federated": providerARN,
		},
		"Condition": condition,
	})
}
