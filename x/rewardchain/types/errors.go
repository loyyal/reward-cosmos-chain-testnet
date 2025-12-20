package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/rewardchain module sentinel errors
var (
	ErrInvalidSigner = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample        = sdkerrors.Register(ModuleName, 1101, "sample error")

	ErrInvalidPartner   = sdkerrors.Register(ModuleName, 1200, "invalid partner")
	ErrPartnerNotFound  = sdkerrors.Register(ModuleName, 1201, "partner not found")
	ErrPartnerDisabled  = sdkerrors.Register(ModuleName, 1202, "partner disabled")
	ErrUnauthorized     = sdkerrors.Register(ModuleName, 1203, "unauthorized")
)
