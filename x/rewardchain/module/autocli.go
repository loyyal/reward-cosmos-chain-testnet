package rewardchain

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "rewardchain/api/rewardchain/rewardchain"
)

// AutoCLIOptions implements the autocli.HasAutoCLIConfig interface.
func (am AppModule) AutoCLIOptions() *autocliv1.ModuleOptions {
	return &autocliv1.ModuleOptions{
		Query: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Query_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "Params",
					Use:       "params",
					Short:     "Shows the parameters of the module",
				},
				{
					RpcMethod: "Partner",
					Use:       "partner [id]",
					Short:     "Shows a partner by id",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "id"},
					},
				},
				{
					RpcMethod: "Partners",
					Use:       "partners",
					Short:     "Lists partners",
				},
				// this line is used by ignite scaffolding # autocli/query
			},
		},
		Tx: &autocliv1.ServiceCommandDescriptor{
			Service: modulev1.Msg_ServiceDesc.ServiceName,
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod:      "CreatePartner",
					Use:            "create-partner [name] [category] [country] [currency] [earn-cost-per-point] [burn-cost-per-point] [total-liquidity]",
					Short:          "Send a create-partner tx",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "name"},
						{ProtoField: "category"},
						{ProtoField: "country"},
						{ProtoField: "currency"},
						{ProtoField: "earnCostPerPoint"},
						{ProtoField: "burnCostPerPoint"},
						{ProtoField: "totalLiquidity"},
					},
				},
				{
					RpcMethod:      "AddPartnerLiquidity",
					Use:            "add-partner-liquidity [partner-id] [amount] [currency] [ext-wallet]",
					Short:          "Add external liquidity to a partner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "partnerId"},
						{ProtoField: "amount"},
						{ProtoField: "currency"},
						{ProtoField: "extWallet"},
					},
				},
				{
					RpcMethod:      "Swap",
					Use:            "swap [partner-id] [route] [points]",
					Short:          "Swap points and tokens for a partner",
					PositionalArgs: []*autocliv1.PositionalArgDescriptor{
						{ProtoField: "partnerId"},
						{ProtoField: "route"},
						{ProtoField: "points"},
					},
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
