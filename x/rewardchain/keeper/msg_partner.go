package keeper

import (
	"context"
	"strings"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"reward-chain/x/rewardchain/types"
)

func (k msgServer) CreatePartner(goCtx context.Context, req *types.MsgCreatePartner) (*types.MsgCreatePartnerResponse, error) {
	if req == nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "empty request")
	}
	if _, err := sdk.AccAddressFromBech32(req.Creator); err != nil {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "invalid creator address")
	}
	if strings.TrimSpace(req.Name) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "name is required")
	}
	if strings.TrimSpace(req.Country) == "" {
		return nil, errorsmod.Wrap(types.ErrInvalidPartner, "country is required")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	if err := requireAdmin(k.Keeper, ctx, req.Creator); err != nil {
		return nil, err
	}

	id, err := k.GetNextPartnerID(ctx)
	if err != nil {
		return nil, err
	}

	p := types.Partner{
		Id:       id,
		Name:     strings.TrimSpace(req.Name),
		Category: strings.TrimSpace(req.Category),
		Location: strings.TrimSpace(req.Location),
		Country:  strings.TrimSpace(req.Country),
		Disabled: false,
	}

	if err := k.SetPartner(ctx, p); err != nil {
		return nil, err
	}
	k.SetPartnerCounter(ctx, id)

	return &types.MsgCreatePartnerResponse{Id: id}, nil
}

func (k msgServer) DisablePartner(goCtx context.Context, req *types.MsgDisablePartner) (*types.MsgDisablePartnerResponse, error) {
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
		return &types.MsgDisablePartnerResponse{}, nil
	}
	p.Disabled = true
	if err := k.SetPartner(ctx, p); err != nil {
		return nil, err
	}

	return &types.MsgDisablePartnerResponse{}, nil
}


