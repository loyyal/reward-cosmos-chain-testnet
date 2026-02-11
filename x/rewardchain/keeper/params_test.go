package keeper_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	keepertest "rewardchain/testutil/keeper"
	"rewardchain/x/rewardchain/types"
)

func TestGetParams(t *testing.T) {
	k, ctx := keepertest.RewardchainKeeper(t)
	params := types.DefaultParams()

	require.NoError(t, k.SetParams(ctx, params))
	require.EqualValues(t, params, k.GetParams(ctx))
}
