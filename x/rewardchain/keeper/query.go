package keeper

import (
	"reward-chain/x/rewardchain/types"
)

var _ types.QueryServer = Keeper{}
