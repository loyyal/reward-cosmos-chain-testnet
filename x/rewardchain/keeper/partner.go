package keeper

import (
	"context"
	"encoding/binary"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	"github.com/cosmos/cosmos-sdk/types/query"

	"rewardchain/x/rewardchain/types"
)

func (k Keeper) getPartnerStore(ctx context.Context) prefix.Store {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return prefix.NewStore(store, types.PartnerKeyPrefix)
}

func (k Keeper) GetNextPartnerID(ctx context.Context) (uint64, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.PartnerCountKey)
	if bz == nil {
		// IDs start at 1
		return 1, nil
	}
	if len(bz) != 8 {
		return 0, errorsmod.Wrap(types.ErrInvalidPartner, "invalid partner counter")
	}
	return binary.BigEndian.Uint64(bz) + 1, nil
}

func (k Keeper) SetPartnerCounter(ctx context.Context, lastID uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, lastID)
	store.Set(types.PartnerCountKey, bz)
}

func (k Keeper) SetPartner(ctx context.Context, p types.Partner) error {
	ps := k.getPartnerStore(ctx)
	bz, err := k.cdc.Marshal(&p)
	if err != nil {
		return err
	}
	ps.Set(types.PartnerKey(p.Id), bz)
	return nil
}

func (k Keeper) GetPartner(ctx context.Context, id uint64) (types.Partner, bool) {
	ps := k.getPartnerStore(ctx)
	bz := ps.Get(types.PartnerKey(id))
	if bz == nil {
		return types.Partner{}, false
	}
	var p types.Partner
	k.cdc.MustUnmarshal(bz, &p)
	return p, true
}

func (k Keeper) GetAllPartners(ctx context.Context) ([]types.Partner, error) {
	ps := k.getPartnerStore(ctx)
	iter := ps.Iterator(nil, nil)
	defer iter.Close()

	out := make([]types.Partner, 0)
	for ; iter.Valid(); iter.Next() {
		var p types.Partner
		k.cdc.MustUnmarshal(iter.Value(), &p)
		out = append(out, p)
	}
	return out, nil
}

func (k Keeper) PaginatePartners(
	ctx context.Context,
	pageReq *query.PageRequest,
	includeDisabled bool,
) ([]types.Partner, *query.PageResponse, error) {
	ps := k.getPartnerStore(ctx)
	partners := make([]types.Partner, 0)

	// Collect all partners first, then filter if needed
	iter := ps.Iterator(nil, nil)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		var p types.Partner
		k.cdc.MustUnmarshal(iter.Value(), &p)
		if includeDisabled || !p.Disabled {
			partners = append(partners, p)
		}
	}

	// Apply pagination
	start := uint64(0)
	limit := uint64(100) // default limit
	if pageReq != nil {
		if pageReq.Offset > 0 {
			start = pageReq.Offset
		}
		if pageReq.Limit > 0 {
			limit = pageReq.Limit
		}
	}

	total := uint64(len(partners))
	end := start + limit
	if start >= total {
		return []types.Partner{}, &query.PageResponse{
			NextKey: nil,
			Total:   total,
		}, nil
	}
	if end > total {
		end = total
	}

	paginatedPartners := partners[start:end]
	var nextKey []byte
	if end < total && len(paginatedPartners) > 0 {
		// Set next key based on the last partner's ID
		nextKey = types.PartnerKey(paginatedPartners[len(paginatedPartners)-1].Id + 1)
	}

	return paginatedPartners, &query.PageResponse{
		NextKey: nextKey,
		Total:   total,
	}, nil
}
