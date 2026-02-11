package types

const (
	// ModuleName defines the module name
	ModuleName = "rewardchain"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_rewardchain"
)

var (
	ParamsKey = []byte("p_rewardchain")

	PartnerCountKey  = []byte("p_rewardchain_partner_count")
	PartnerKeyPrefix = []byte("p_rewardchain_partner/")
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}
