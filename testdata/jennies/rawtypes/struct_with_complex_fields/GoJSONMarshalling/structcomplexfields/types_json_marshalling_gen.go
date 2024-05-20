package structcomplexfields
func (resource StringOrBool) MarshalJSON() ([]byte, error) {
	if resource.String != nil {
		return json.Marshal(resource.String)
	}

	if resource.Bool != nil {
		return json.Marshal(resource.Bool)
	}

	return nil, fmt.Errorf("no value for disjunction of scalars")
}

func (resource *StringOrBool) UnmarshalJSON(raw []byte) error {
	if raw == nil {
		return nil
	}

	var errList []error

	// String
	var stringArg string
	if err := json.Unmarshal(raw, &stringArg); err != nil {
		errList = append(errList, err)
		resource.String = nil
	} else {
		resource.String = &stringArg
		return nil
	}

	// Bool
	var boolArg bool
	if err := json.Unmarshal(raw, &boolArg); err != nil {
		errList = append(errList, err)
		resource.Bool = nil
	} else {
		resource.Bool = &boolArg
		return nil
	}

	return errors.Join(errList...)
}

