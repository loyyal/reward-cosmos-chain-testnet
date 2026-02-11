package keeper

import (
	"context"

	"rewardchain/x/rewardchain/types"
)

func (k Keeper) IsAdmin(ctx context.Context, addr string) bool {
	// NOTE: keeper.GetParams expects a context.Context; sdk.Context satisfies it.
	params := k.GetParams(ctx)
	for _, a := range params.AdminAddresses {
		if a == addr {
			return true
		}
	}
	return false
}

func requireAdmin(k Keeper, ctx context.Context, addr string) error {
	if k.IsAdmin(ctx, addr) {
		return nil
	}
	return types.ErrUnauthorized
}
