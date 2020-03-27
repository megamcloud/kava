package keeper

import (
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/kava-labs/kava/x/incentive/types"
)

// Keeper keeper for the incentive module
type Keeper struct {
	key           sdk.StoreKey
	cdc           *codec.Codec
	paramSubspace subspace.Subspace
	supplyKeeper  types.SupplyKeeper
	cdpKeeper     types.CdpKeeper
	accountKeeper types.AccountKeeper
	codespace     sdk.CodespaceType
}

// NewKeeper creates a new keeper
func NewKeeper(cdc *codec.Codec, key sdk.StoreKey, paramstore subspace.Subspace, sk types.SupplyKeeper, cdpk types.CdpKeeper, ak types.AccountKeeper, codespace sdk.CodespaceType) Keeper {

	return Keeper{
		key:           key,
		cdc:           cdc,
		paramSubspace: paramstore.WithKeyTable(types.ParamKeyTable()),
		supplyKeeper:  sk,
		cdpKeeper:     cdpk,
		accountKeeper: ak,
		codespace:     codespace,
	}
}

func (k Keeper) Codespace() sdk.CodespaceType {
	return k.codespace
}

// GetRewardPeriod returns the reward period from the store for the input denom and a boolean for if it was found
func (k Keeper) GetRewardPeriod(ctx sdk.Context, denom string) (types.RewardPeriod, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.RewardPeriodKeyPrefix)
	bz := store.Get(types.GetDenomBytes(denom))
	if bz == nil {
		return types.RewardPeriod{}, false
	}
	var rp types.RewardPeriod
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &rp)
	return rp, true
}

// SetRewardPeriod sets the reward period in the store for the input deno,
func (k Keeper) SetRewardPeriod(ctx sdk.Context, rp types.RewardPeriod) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.RewardPeriodKeyPrefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(rp)
	store.Set(types.GetDenomBytes(rp.Denom), bz)
}

// DeleteRewardPeriod deletes the reward period in the store for the input denom,
func (k Keeper) DeleteRewardPeriod(ctx sdk.Context, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.RewardPeriodKeyPrefix)
	store.Delete(types.GetDenomBytes(denom))
}

// IterateRewardPeriods iterates over all reward period objects in the store and preforms a callback function
func (k Keeper) IterateRewardPeriods(ctx sdk.Context, cb func(rp types.RewardPeriod) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.RewardPeriodKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var rp types.RewardPeriod
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &rp)
		if cb(rp) {
			break
		}
	}
}

// IterateClaimPeriodIDKeysAndValues iterates over the claim period id (value) and denom (key) of each claim period id in the store and performs a callback function
func (k Keeper) IterateClaimPeriodIDKeysAndValues(ctx sdk.Context, cb func(denom string, id uint64) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.NextClaimPeriodIDPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		id := types.BytesToUint64(iterator.Value())
		denom := types.GetDenomFromBytes(iterator.Key())
		if cb(denom, id) {
			break
		}
	}
}

// GetNextClaimPeriodID returns the highest claim period id in the store for the input denom
func (k Keeper) GetNextClaimPeriodID(ctx sdk.Context, denom string) uint64 {
	store := prefix.NewStore(ctx.KVStore(k.key), types.NextClaimPeriodIDPrefix)
	bz := store.Get(types.GetDenomBytes(denom))
	id := types.BytesToUint64(bz)
	return id
}

// SetNextClaimPeriodID sets the highest claim period id in the store for the input denom
func (k Keeper) SetNextClaimPeriodID(ctx sdk.Context, denom string, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.NextClaimPeriodIDPrefix)
	store.Set(types.GetDenomBytes(denom), types.GetIDBytes(id))
}

// GetClaimPeriod returns claim period in the store for the input ID and denom and a boolean for if it was found
func (k Keeper) GetClaimPeriod(ctx sdk.Context, id uint64, denom string) (types.ClaimPeriod, bool) {
	var cp types.ClaimPeriod
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimPeriodKeyPrefix)
	bz := store.Get(types.GetClaimPeriodPrefix(denom, id))
	if bz == nil {
		return types.ClaimPeriod{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &cp)
	return cp, true
}

// SetClaimPeriod sets the claim period in the store for the input ID and denom
func (k Keeper) SetClaimPeriod(ctx sdk.Context, cp types.ClaimPeriod) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimPeriodKeyPrefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(cp)
	store.Set(types.GetClaimPeriodPrefix(cp.Denom, cp.ID), bz)
}

// DeleteClaimPeriod deletes the claim period in the store for the input ID and denom
func (k Keeper) DeleteClaimPeriod(ctx sdk.Context, id uint64, denom string) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimPeriodKeyPrefix)
	store.Delete(types.GetClaimPeriodPrefix(denom, id))
}

// IterateClaimPeriods iterates over all claim period objects in the store and preforms a callback function
func (k Keeper) IterateClaimPeriods(ctx sdk.Context, cb func(cp types.ClaimPeriod) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimPeriodKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var cp types.ClaimPeriod
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &cp)
		if cb(cp) {
			break
		}
	}
}

// GetClaim returns the claim in the store corresponding the the input address denom and id and a boolean for if the claim was found
func (k Keeper) GetClaim(ctx sdk.Context, addr sdk.AccAddress, denom string, id uint64) (types.Claim, bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimKeyPrefix)
	bz := store.Get(types.GetClaimPrefix(addr, denom, id))
	if bz == nil {
		return types.Claim{}, false
	}
	var c types.Claim
	k.cdc.MustUnmarshalBinaryLengthPrefixed(bz, &c)
	return c, true
}

// SetClaim sets the claim in the store corresponding to the input address, denom, and id
func (k Keeper) SetClaim(ctx sdk.Context, c types.Claim) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimKeyPrefix)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(c)
	store.Set(types.GetClaimPrefix(c.Owner, c.Denom, c.ID), bz)

}

// DeleteClaim deletes the claim in the store corresponding to the input address, denom, and id
func (k Keeper) DeleteClaim(ctx sdk.Context, owner sdk.AccAddress, denom string, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimKeyPrefix)
	store.Delete(types.GetClaimPrefix(owner, denom, id))
}

// IterateClaims iterates over all claim  objects in the store and preforms a callback function
func (k Keeper) IterateClaims(ctx sdk.Context, cb func(c types.Claim) (stop bool)) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.ClaimKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var c types.Claim
		k.cdc.MustUnmarshalBinaryLengthPrefixed(iterator.Value(), &c)
		if cb(c) {
			break
		}
	}
}

// GetPreviousBlockTime get the blocktime for the previous block
func (k Keeper) GetPreviousBlockTime(ctx sdk.Context) (blockTime time.Time, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousBlockTimeKey)
	b := store.Get([]byte{})
	if b == nil {
		return time.Time{}, false
	}
	k.cdc.MustUnmarshalBinaryLengthPrefixed(b, &blockTime)
	return blockTime, true
}

// SetPreviousBlockTime set the time of the previous block
func (k Keeper) SetPreviousBlockTime(ctx sdk.Context, blockTime time.Time) {
	store := prefix.NewStore(ctx.KVStore(k.key), types.PreviousBlockTimeKey)
	store.Set([]byte{}, k.cdc.MustMarshalBinaryLengthPrefixed(blockTime))
}