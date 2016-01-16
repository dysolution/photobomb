package main

import sdk "github.com/dysolution/espsdk"

// A Missile represents an action the API client performs whose URL isn't
// known until runtime, such as the retrieval or deletion of the most
// recently created Batch.
type Missile struct {
	client    *sdk.Client
	Name      string                     `json:"name"`
	Operation func() (sdk.Result, error) `json:"-"`
}

// Deploy fires the Missile.
func (m Missile) Deploy() (sdk.Result, error) {
	result, err := m.Operation()
	if err != nil {
		result.Log().Errorf("%s.Deploy %v: %v", m.Name, m.Operation, err)
		return sdk.Result{}, err
	}
	result.Log().Debugf("%s.Deploy", m.Name)
	return result, nil
}
