package rewardchain

import (
	autocliv1 "cosmossdk.io/api/cosmos/autocli/v1"

	modulev1 "reward-chain/api/rewardchain/rewardchain"
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
			Service:              modulev1.Msg_ServiceDesc.ServiceName,
			EnhanceCustomCommand: true, // only required if you want to use the custom command
			RpcCommandOptions: []*autocliv1.RpcCommandOptions{
				{
					RpcMethod: "UpdateParams",
					Skip:      true, // skipped because authority gated
				},
				{
					RpcMethod: "CreatePartner",
					Use:       "create-partner [name] [category] [location] [country]",
					Short:     "Create a partner",
				},
				{
					RpcMethod: "DisablePartner",
					Use:       "disable-partner [id]",
					Short:     "Disable a partner by id",
				},
				{
					RpcMethod: "UpdatePartner",
					Use:       "update-partner [id] [name] [category] [location] [country]",
					Short:     "Update a partner",
				},
				// this line is used by ignite scaffolding # autocli/tx
			},
		},
	}
}
