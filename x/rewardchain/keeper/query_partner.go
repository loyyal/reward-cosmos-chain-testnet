package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"reward-chain/x/rewardchain/types"
)

func (k Keeper) Partner(goCtx context.Context, req *types.QueryPartnerRequest) (*types.QueryPartnerResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	if req.Id == 0 {
		return nil, status.Error(codes.InvalidArgument, "id must be > 0")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	p, found := k.GetPartner(ctx, req.Id)
	if !found {
		return nil, status.Error(codes.NotFound, errorsmod.Wrap(types.ErrPartnerNotFound, "partner not found").Error())
	}

	return &types.QueryPartnerResponse{Partner: p}, nil
}

func (k Keeper) Partners(goCtx context.Context, req *types.QueryPartnersRequest) (*types.QueryPartnersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	partners, pageRes, err := k.PaginatePartners(ctx, req.Pagination, req.IncludeDisabled)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPartnersResponse{
		Partners:   partners,
		Pagination: pageRes,
	}, nil
}


