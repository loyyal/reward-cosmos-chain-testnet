package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSwap{}

func NewMsgSwap(creator string, partnerId uint64, route, points string) *MsgSwap {
	return &MsgSwap{
		Creator:   creator,
		PartnerId: partnerId,
		Route:     route,
		Points:    points,
	}
}

func (msg *MsgSwap) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.PartnerId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "partner_id must be > 0")
	}
	if msg.Route == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "route is required")
	}
	if msg.Points == "" {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "points is required")
	}
	return nil
}

