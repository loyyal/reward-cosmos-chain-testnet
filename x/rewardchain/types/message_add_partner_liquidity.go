package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddPartnerLiquidity{}

func NewMsgAddPartnerLiquidity(creator string, partnerId uint64, amount, currency, extWallet string) *MsgAddPartnerLiquidity {
	return &MsgAddPartnerLiquidity{
		Creator:   creator,
		PartnerId: partnerId,
		Amount:    amount,
		Currency:  currency,
		ExtWallet: extWallet,
	}
}

func (msg *MsgAddPartnerLiquidity) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	if msg.PartnerId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "partner_id must be > 0")
	}
	if len(msg.Amount) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "amount is required")
	}
	return nil
}

