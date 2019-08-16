package rsrc

import "github.com/pkg/errors"

func identifier(i interface{}) (*int, *string, error) {
	if id, ok := i.(int); ok {
		return &id, nil, nil
	} else if name, ok := i.(string); ok {
		return nil, &name, nil
	}
	return nil, nil, errors.Errorf("wrong identifier %v(%T) for resource identifier", i, i)
}
