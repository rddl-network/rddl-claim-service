package service

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	planetmint "github.com/planetmint/planetmint-go/lib"
	daotypes "github.com/planetmint/planetmint-go/x/dao/types"
	"github.com/rddl-network/rddl-claim-service/config"
)

type IPlanetmintClient interface {
	SendConfirmation(claimID int, beneficiary string) (txResponse sdk.TxResponse, err error)
}

type PlanetmintClient struct{}

func NewPlanetmintClient() *PlanetmintClient {
	return &PlanetmintClient{}
}

func (pc *PlanetmintClient) SendConfirmation(claimID int, beneficiary string) (txResponse sdk.TxResponse, err error) {
	cfg := config.GetConfig()
	addressString := cfg.PlanetmintAddress
	addr := sdk.MustAccAddressFromBech32(addressString)
	msg := daotypes.NewMsgConfirmRedeemClaim(addressString, uint64(claimID), beneficiary)

	out, err := planetmint.BroadcastTxWithFileLock(addr, msg)
	if err != nil {
		return
	}

	return planetmint.GetTxResponseFromOut(out)
}
