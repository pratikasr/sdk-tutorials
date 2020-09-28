package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/sdk-tutorials/starport-nameservice/nameservice/x/nameservice/types"
)

// CreateWhois creates a whois
func (k Keeper) CreateWhois(ctx sdk.Context, whois types.Whois) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.WhoisPrefix + whois.Value)
	value := k.cdc.MustMarshalBinaryLengthPrefixed(whois)
	store.Set(key, value)
}

// GetWhois returns the whois information
func (k Keeper) GetWhois(ctx sdk.Context, key string) (types.Whois, error) {
	store := ctx.KVStore(k.storeKey)
	var whois types.Whois
	byteKey := []byte(types.WhoisPrefix + key)
	err := k.cdc.UnmarshalBinaryLengthPrefixed(store.Get(byteKey), &whois)
	if err != nil {
		return whois, err
	}
	return whois, nil
}

// SetWhois sets a whois
func (k Keeper) SetWhois(ctx sdk.Context, whois types.Whois) {
	whoisKey := whois.Value
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshalBinaryLengthPrefixed(whois)
	key := []byte(types.WhoisPrefix + whoisKey)
	store.Set(key, bz)
}

// DeleteWhois deletes a whois
func (k Keeper) DeleteWhois(ctx sdk.Context, key string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete([]byte(types.WhoisPrefix + key))
}

//
// Functions used by querier
//

func listWhois(ctx sdk.Context, k Keeper) ([]byte, error) {
	var whoisList []types.Whois
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, []byte(types.WhoisPrefix))
	for ; iterator.Valid(); iterator.Next() {
		var whois types.Whois
		k.cdc.MustUnmarshalBinaryLengthPrefixed(store.Get(iterator.Key()), &whois)
		whoisList = append(whoisList, whois)
	}
	res := codec.MustMarshalJSONIndent(k.cdc, whoisList)
	return res, nil
}

func getWhois(ctx sdk.Context, path []string, k Keeper) (res []byte, sdkError error) {
	key := path[0]
	whois, err := k.GetWhois(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err = codec.MarshalJSONIndent(k.cdc, whois)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// nolint: unparam
func resolveName(ctx sdk.Context, path []string, keeper Keeper) ([]byte, error) {
	value := keeper.ResolveName(ctx, path[0])

	if value == "" {
		return []byte{}, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, "could not resolve name")
	}

	res, err := codec.MarshalJSONIndent(keeper.cdc, types.QueryResResolve{Value: value})
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONMarshal, err.Error())
	}

	return res, nil
}

// Get creator of the item
func (k Keeper) GetOwner(ctx sdk.Context, key string) sdk.AccAddress {
	whois, _ := k.GetWhois(ctx, key)
	return whois.Owner
}

// Check if the key exists in the store
func (k Keeper) Exists(ctx sdk.Context, key string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(types.WhoisPrefix + key))
}

// ResolveName - returns the string that the name resolves to
func (k Keeper) ResolveName(ctx sdk.Context, name string) string {
	whois, _ := k.GetWhois(ctx, name)
	return whois.Value
}

// SetName - sets the value string that a name resolves to
func (k Keeper) SetName(ctx sdk.Context, name string, value string) {
	whois, _ := k.GetWhois(ctx, name)
	whois.Value = value
	k.SetWhois(ctx, whois)
}

// HasOwner - returns whether or not the name already has an owner
func (k Keeper) HasOwner(ctx sdk.Context, name string) bool {
	whois, _ := k.GetWhois(ctx, name)
	return !whois.Owner.Empty()
}

// SetOwner - sets the current owner of a name
func (k Keeper) SetOwner(ctx sdk.Context, name string, owner sdk.AccAddress) {
	whois, _ := k.GetWhois(ctx, name)
	whois.Owner = owner
	k.SetWhois(ctx, whois)
}

// GetPrice - gets the current price of a name
func (k Keeper) GetPrice(ctx sdk.Context, name string) sdk.Coins {
	whois, _ := k.GetWhois(ctx, name)
	return whois.Price
}

// SetPrice - sets the current price of a name
func (k Keeper) SetPrice(ctx sdk.Context, name string, price sdk.Coins) {
	whois, _ := k.GetWhois(ctx, name)
	whois.Price = price
	k.SetWhois(ctx, whois)
}

// Check if the name is present in the store or not
func (k Keeper) IsNamePresent(ctx sdk.Context, name string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has([]byte(name))
}

// Get an iterator over all names in which the keys are the names and the values are the whois
func (k Keeper) GetNamesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, []byte{})
}
