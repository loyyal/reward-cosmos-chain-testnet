package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePartner{}

func NewMsgCreatePartner(creator string, name string, category string, country string, currency string, earnCostPerPoint string, burnCostPerPoint string, totalLiquidity string) *MsgCreatePartner {
	return &MsgCreatePartner{
		Creator:          creator,
		Name:             name,
		Category:         category,
		Country:          country,
		Currency:         currency,
		EarnCostPerPoint: earnCostPerPoint,
		BurnCostPerPoint: burnCostPerPoint,
		TotalLiquidity:   totalLiquidity,
	}
}

func (msg *MsgCreatePartner) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
