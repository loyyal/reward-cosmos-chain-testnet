package types

import (
	"fmt"
	"strings"
)

// this line is used by starport scaffolding # genesis/types/import

// DefaultIndex is the default global index
const DefaultIndex uint64 = 1

// DefaultGenesis returns the default genesis state
func DefaultGenesis() *GenesisState {
	return &GenesisState{
		// this line is used by starport scaffolding # genesis/types/default
		Params:   DefaultParams(),
		Partners: []Partner{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	// this line is used by starport scaffolding # genesis/types/validate
	seen := make(map[uint64]struct{}, len(gs.Partners))
	for _, p := range gs.Partners {
		if p.Id == 0 {
			return fmt.Errorf("partner id must be > 0")
		}
		if strings.TrimSpace(p.Name) == "" {
			return fmt.Errorf("partner name is required (id=%d)", p.Id)
		}
		if strings.TrimSpace(p.Country) == "" {
			return fmt.Errorf("partner country is required (id=%d)", p.Id)
		}
		if _, ok := seen[p.Id]; ok {
			return fmt.Errorf("duplicate partner id %d", p.Id)
		}
		seen[p.Id] = struct{}{}
	}

	return gs.Params.Validate()
}
