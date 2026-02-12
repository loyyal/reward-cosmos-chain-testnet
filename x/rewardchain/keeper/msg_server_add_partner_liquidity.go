package keeper

import (
	"context"
	"strconv"
	"strings"

	errorsmod "cosmossdk.io/errors"
	math "cosmossdk.io/math"
	"rewardchain/x/rewardchain/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// AddPartnerLiquidity handles MsgAddPartnerLiquidity messages.
// It increases a partner's total and available liquidity based on the
// provided amount and the current RedeemCostPerPoint.
func (k msgServer) AddPartnerLiquidity(goCtx context.Context, msg *types.MsgAddPartnerLiquidity) (*types.MsgAddPartnerLiquidityResponse, error) {
	if msg == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "empty request")
	}

	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid creator address")
	}
	if msg.PartnerId == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "partner_id must be > 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	// NOTE: admin check is currently disabled in CreatePartner for testing.
	// Uncomment the following lines if you want to require admin privileges
	// for adding liquidity as well.
	//
	// if err := requireAdmin(k.Keeper, ctx, msg.Creator); err != nil {
	// 	return nil, err
	// }

	p, found := k.GetPartner(ctx, msg.PartnerId)
	if !found {
		return nil, types.ErrPartnerNotFound
	}
	if p.Disabled {
		return nil, types.ErrPartnerDisabled
	}

	amountStr := strings.TrimSpace(msg.Amount)
	if amountStr == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "amount is required")
	}
	redeemStr := strings.TrimSpace(p.RedeemCostPerPoint)
	if redeemStr == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "redeem_cost_per_point is required on partner")
	}

	amountDec, err := math.LegacyNewDecFromStr(amountStr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid amount")
	}
	redeemDec, err := math.LegacyNewDecFromStr(redeemStr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid redeem_cost_per_point on partner")
	}
	if redeemDec.IsZero() {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "redeem_cost_per_point must be > 0")
	}

	// points = amount / redeem_cost_per_point
	points := amountDec.Quo(redeemDec)

	totalStr := strings.TrimSpace(p.TotalLiquidity)
	if totalStr == "" {
		totalStr = "0"
	}
	totalDec, err := math.LegacyNewDecFromStr(totalStr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid total_liquidity on partner")
	}

	availStr := strings.TrimSpace(p.AvailableLiquidity)
	if availStr == "" {
		availStr = "0"
	}
	availDec, err := math.LegacyNewDecFromStr(availStr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid available_liquidity on partner")
	}

	newTotal := totalDec.Add(points)
	newAvail := availDec.Add(points)

	p.TotalLiquidity = newTotal.String()
	p.AvailableLiquidity = newAvail.String()

	if err := k.SetPartner(ctx, p); err != nil {
		return nil, err
	}

	// Emit an event for the added liquidity.
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"add_partner_liquidity",
			sdk.NewAttribute("partner_id", strconv.FormatUint(msg.PartnerId, 10)),
			sdk.NewAttribute("amount", msg.Amount),
			sdk.NewAttribute("currency", msg.Currency),
			sdk.NewAttribute("ext_wallet", msg.ExtWallet),
		),
	)

	return &types.MsgAddPartnerLiquidityResponse{}, nil
}

