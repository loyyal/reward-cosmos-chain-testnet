package keeper

import (
	"context"
	"strconv"
	"strings"

	"rewardchain/x/rewardchain/types"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k msgServer) CreatePartner(goCtx context.Context, msg *types.MsgCreatePartner) (*types.MsgCreatePartnerResponse, error) {
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "empty request")
	}
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid creator address")
	}
	if strings.TrimSpace(msg.Name) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "name is required")
	}
	if strings.TrimSpace(msg.Country) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "country is required")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	// TODO: Re-enable admin check after adding addresses to genesis
	// Temporarily disabled for testing
	// if err := requireAdmin(k.Keeper, ctx, msg.Creator); err != nil {
	// 	return nil, err
	// }

	id, err := k.GetNextPartnerID(ctx)
	if err != nil {
		return nil, err
	}

	// Initialize available_liquidity to total_liquidity, on_hold_liquidity to "0"
	availableLiquidity := msg.TotalLiquidity
	if availableLiquidity == "" {
		availableLiquidity = "0"
	}
	onHoldLiquidity := "0"

	p := types.Partner{
		Id:                 id,
		Name:               strings.TrimSpace(msg.Name),
		Category:           strings.TrimSpace(msg.Category),
		Location:           "", // Not in message, set to empty
		Country:            strings.TrimSpace(msg.Country),
		Disabled:           false,
		TotalLiquidity:     strings.TrimSpace(msg.TotalLiquidity),
		AvailableLiquidity: availableLiquidity,
		OnHoldLiquidity:    onHoldLiquidity,
		EarnCostPerPoint:   strings.TrimSpace(msg.EarnCostPerPoint),
		RedeemCostPerPoint: strings.TrimSpace(msg.BurnCostPerPoint), // Map burnCostPerPoint to redeem_cost_per_point
		StartsFrom:         "",                                      // Not in message, set to empty
		EndsBefore:         "",                                      // Not in message, set to empty
	}

	if err := k.SetPartner(ctx, p); err != nil {
		return nil, err
	}
	k.SetPartnerCounter(ctx, id)

	return &types.MsgCreatePartnerResponse{Id: strconv.FormatUint(id, 10)}, nil
}
