package rewardchain

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"rewardchain/x/rewardchain/keeper"
	"rewardchain/x/rewardchain/types"
)

// InitGenesis initializes the module's state from a provided genesis state.
func InitGenesis(ctx sdk.Context, k keeper.Keeper, genState types.GenesisState) {
	// this line is used by starport scaffolding # genesis/module/init
	if err := k.SetParams(ctx, genState.Params); err != nil {
		panic(err)
	}

	var maxID uint64
	for _, p := range genState.Partners {
		if err := k.SetPartner(ctx, p); err != nil {
			panic(err)
		}
		if p.Id > maxID {
			maxID = p.Id
		}
	}
	if maxID > 0 {
		k.SetPartnerCounter(ctx, maxID)
	}
}

// ExportGenesis returns the module's exported genesis.
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	genesis := types.DefaultGenesis()
	genesis.Params = k.GetParams(ctx)

	// this line is used by starport scaffolding # genesis/module/export
	partners, err := k.GetAllPartners(ctx)
	if err != nil {
		panic(err)
	}
	genesis.Partners = partners

	return genesis
}
