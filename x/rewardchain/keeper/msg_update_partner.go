package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"reward-chain/x/rewardchain/types"
)

func (k msgServer) UpdatePartner(goCtx context.Context, req *types.MsgUpdatePartner) (*types.MsgUpdatePartnerResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "empty request")
	}
	if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid creator address")
	}
	if req.Id == 0 {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "id must be > 0")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := requireAdmin(k.Keeper, ctx, req.Creator); err != nil {
		return nil, err
	}

	p, found := k.GetPartner(ctx, req.Id)
	if !found {
		return nil, types.ErrPartnerNotFound
	}
	if p.Disabled {
		return nil, types.ErrPartnerDisabled
	}

	// Update provided fields (treat empty as "leave unchanged")
	if v := strings.TrimSpace(req.Name); v != "" {
		p.Name = v
	}
	if v := strings.TrimSpace(req.Category); v != "" {
		p.Category = v
	}
	if v := strings.TrimSpace(req.Location); v != "" {
		p.Location = v
	}
	if v := strings.TrimSpace(req.Country); v != "" {
		p.Country = v
	}

	if strings.TrimSpace(p.Name) == "" || strings.TrimSpace(p.Country) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "name and country are required")
	}

	if err := k.SetPartner(ctx, p); err != nil {
		return nil, err
	}

	return &types.MsgUpdatePartnerResponse{}, nil
}


