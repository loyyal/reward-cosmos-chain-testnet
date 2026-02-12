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

// Swap handles MsgSwap messages.
// For route == "points_to_token", it removes that many points from the partner's liquidity.
// For route == "token_to_points", it currently does nothing (to be defined later).
func (k msgServer) Swap(goCtx context.Context, msg *types.MsgSwap) (*types.MsgSwapResponse, error) {
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

	p, found := k.GetPartner(ctx, msg.PartnerId)
	if !found {
		return nil, types.ErrPartnerNotFound
	}
	if p.Disabled {
		return nil, types.ErrPartnerDisabled
	}

	route := strings.TrimSpace(strings.ToLower(msg.Route))
	pointsStr := strings.TrimSpace(msg.Points)
	if pointsStr == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "points is required")
	}

	pointsDec, err := math.LegacyNewDecFromStr(pointsStr)
	if err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid points")
	}
	if pointsDec.IsNegative() {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "points must be >= 0")
	}

	switch route {
	case "points_to_token":
		// Remove points from partner's liquidity.
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

		if pointsDec.GT(totalDec) || pointsDec.GT(availDec) {
			return nil, errorsmod.Wrap(types.ErrInvalidPartner, "insufficient liquidity for points_to_token swap")
		}

		newTotal := totalDec.Sub(pointsDec)
		newAvail := availDec.Sub(pointsDec)

		p.TotalLiquidity = newTotal.String()
		p.AvailableLiquidity = newAvail.String()

		if err := k.SetPartner(ctx, p); err != nil {
			return nil, err
		}

		// Emit event for the swap.
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"partner_swap",
				sdk.NewAttribute("partner_id", strconv.FormatUint(msg.PartnerId, 10)),
				sdk.NewAttribute("route", msg.Route),
				sdk.NewAttribute("points", msg.Points),
			),
		)

	case "token_to_points":
		// TODO: Define behavior for token_to_points later.
		// For now, just emit an event without changing state.
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				"partner_swap",
				sdk.NewAttribute("partner_id", strconv.FormatUint(msg.PartnerId, 10)),
				sdk.NewAttribute("route", msg.Route),
				sdk.NewAttribute("points", msg.Points),
			),
		)

	default:
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid route (must be points_to_token or token_to_points)")
	}

	return &types.MsgSwapResponse{}, nil
}

