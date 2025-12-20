package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(adminAddresses []string) Params {
	return Params{AdminAddresses: adminAddresses}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	// Use an explicit empty slice (not nil) to make equality checks stable in tests and JSON.
	return NewParams([]string{})
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate validates the set of params
func (p Params) Validate() error {
	seen := make(map[string]struct{}, len(p.AdminAddresses))
	for _, a := range p.AdminAddresses {
		if _, err := sdk.AccAddressFromBech32(a); err != nil {
			return fmt.Errorf("invalid admin address %q: %w", a, err)
		}
		if _, ok := seen[a]; ok {
			return fmt.Errorf("duplicate admin address %q", a)
		}
		seen[a] = struct{}{}
	}
	return nil
}
