package keeper

import (
	"rewardchain/x/rewardchain/types"
)

var _ types.QueryServer = Keeper{}
