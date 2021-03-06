package app

import (
	"encoding/json"
	"errors"
	"github.com/QOSGroup/qbase/server"
	"github.com/QOSGroup/qbase/server/config"

	"github.com/QOSGroup/qos/account"
	"github.com/QOSGroup/qos/types"
	"github.com/spf13/pflag"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/ed25519"
	tmtypes "github.com/tendermint/tendermint/types"
)

// QOS初始状态
type GenesisState struct {
	CAPubKey   crypto.PubKey         `json:"ca_pub_key"`
	Accounts   []*account.QOSAccount `json:"accounts"`
	Validators []types.Validator     `json:"validators"`
}

func QOSAppInit() server.AppInit {
	fsAppGenState := pflag.NewFlagSet("", pflag.ContinueOnError)

	fsAppGenTx := pflag.NewFlagSet("", pflag.ContinueOnError)
	fsAppGenTx.String(server.FlagName, "", "validator moniker, required")
	fsAppGenTx.String(server.FlagClientHome, types.DefaultCLIHome,
		"home directory for the client, used for key generation")
	fsAppGenTx.Bool(server.FlagOWK, false, "overwrite the accounts created")

	return server.AppInit{
		FlagsAppGenState: fsAppGenState,
		FlagsAppGenTx:    fsAppGenTx,
		AppGenTx:         QOSAppGenTx,
		AppGenState:      QOSAppGenState,
	}
}

type QOSGenTx struct {
	Validator tmtypes.GenesisValidator `json:"validator"`
}

// Generate a genesis transaction
func QOSAppGenTx(cdc *amino.Codec, pk crypto.PubKey, genTxConfig config.GenTx) (
	appGenTx, cliPrint json.RawMessage, validator tmtypes.GenesisValidator, err error) {

	//JUST 占坑
	validator.PubKey = ed25519.PubKeyEd25519{}
	validator.Power = 1
	validator.Name = "Use app_state.validators Instead"

	simpleGenTx := QOSGenTx{tmtypes.GenesisValidator{
		PubKey: pk,
		Power:  10,
	}}
	bz, err := cdc.MarshalJSON(simpleGenTx)
	if err != nil {
		return
	}
	appGenTx = json.RawMessage(bz)
	return
}

// app_state初始配置项生成
func QOSAppGenState(cdc *amino.Codec, appGenTxs []json.RawMessage) (appState json.RawMessage, err error) {

	if len(appGenTxs) != 1 {
		err = errors.New("must provide a single genesis transaction")
		return
	}

	var genTx QOSGenTx
	err = cdc.UnmarshalJSON(appGenTxs[0], &genTx)
	if err != nil {
		return
	}

	appGenState := GenesisState{}

	appState, _ = cdc.MarshalJSONIndent(appGenState, "", " ")
	return
}
