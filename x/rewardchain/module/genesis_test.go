package rewardchain_test

import (
	"testing"

	keepertest "rewardchain/testutil/keeper"
	"rewardchain/testutil/nullify"
	rewardchain "rewardchain/x/rewardchain/module"
	"rewardchain/x/rewardchain/types"

	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.RewardchainKeeper(t)
	rewardchain.InitGenesis(ctx, k, genesisState)
	got := rewardchain.ExportGenesis(ctx, k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	// this line is used by starport scaffolding # genesis/test/assert
}
