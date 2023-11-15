package app

import (
	"encoding/json"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/syndtr/goleveldb/leveldb"
)

var PricesDB *leveldb.DB

// PricePair represents a structure to hold price information
type PricePair struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

// called before anything else
func init() {
	var err error
	PricesDB, err = leveldb.OpenFile("/Users/richard/workspace/osmosis-price-indexer/prices.db", nil)
	if err != nil {
		panic(err)
	}
}

func indexPrices(app *OsmosisApp, ctx sdk.Context) {
	// /RICHARD MODIFY End block ///
	DIVISOR := new(big.Float).SetFloat64(1000000000000000000)
	SOMMELIER_OSMO_POOL := 627
	ETH_OSMO_POOL := 704
	ctx.Logger().Error("RICHARD - INSIDE End blocker")
	ATOM_DENOM := "ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2"
	USDC_DENOM := "ibc/D189335C6E4A68B513C10AB227BF1C1D38C746766278BA3EEB4FB14124F1D858"
	ETH_DENOM := "ibc/EA1D43981D5C9A1C4AAEA9C23BB1D4FA126BA9BC7020A25E0AE4AA841EA25DC5"
	OSMO_DENOM := "uosmo"
	SOMMELIER_DENOM := "ibc/9BBA9A1C257E971E38C1422780CE6F0B0686F0A3085E2D61118D904BFE0F5F5E"
	atom_to_osmo, err := app.GAMMKeeper.CalculateSpotPrice(ctx, 1, OSMO_DENOM, ATOM_DENOM)
	a_to_o := 0.0
	int_a_to_o := atom_to_osmo.TruncateInt64()
	ctx.Logger().Info("int a to o", int_a_to_o)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		//a_to_o, err := atom_to_osmo.Float64()
		bigFloatVal := new(big.Float).SetInt(atom_to_osmo.BigInt())

		a_to_o, _ = new(big.Float).Quo(bigFloatVal, DIVISOR).Float64()
	}

	o_to_c := 0.0
	osmo_to_usdc, err := app.GAMMKeeper.CalculateSpotPrice(ctx, 678, USDC_DENOM, OSMO_DENOM)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		o_to_c, _ = new(big.Float).SetInt(osmo_to_usdc.BigInt()).Float64()
	}
	//o_to_c, err := osmo_to_usdc.Float64()

	s_to_o := 0.0
	somm_to_osmo, err := app.GAMMKeeper.CalculateSpotPrice(ctx, uint64(SOMMELIER_OSMO_POOL), SOMMELIER_DENOM, OSMO_DENOM)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		s_to_o, _ = new(big.Float).SetInt(somm_to_osmo.BigInt()).Float64()
	}
	//s_to_o, err := somm_to_osmo.Float64()

	e_to_o := 0.0
	eth_to_osmo, err := app.GAMMKeeper.CalculateSpotPrice(ctx, uint64(ETH_OSMO_POOL), ETH_DENOM, OSMO_DENOM)
	if err != nil {
		ctx.Logger().Error(err.Error())
	} else {
		e_to_o, _ = new(big.Float).SetInt(eth_to_osmo.BigInt()).Float64()
	}
	//e_to_o, err := eth_to_osmo.Float64()

	pricePairs := []PricePair{
		{Symbol: "OSMO/USDC", Price: o_to_c},
		{Symbol: "ETH/OSMO", Price: e_to_o},
		{Symbol: "ATOM/OSMO", Price: a_to_o},
		{Symbol: "SOMM/OSMO", Price: s_to_o},
	}
	// Convert the price pair to JSON
	data, err := json.Marshal(pricePairs)
	if err != nil {
		panic(err)
	}
	// Use current Unix timestamp as key
	timeKey := strconv.FormatInt(ctx.BlockTime().Unix(), 10)
	// Write the data to the database
	err = PricesDB.Put([]byte(timeKey), data, nil)
	if err != nil {
		panic(err)
	}
	ctx.Logger().Info("Price Data written successfully")

}
