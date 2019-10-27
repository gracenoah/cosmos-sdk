package cli

import (
	"github.com/gracenoah/cosmos-sdk/client/context"
	"github.com/gracenoah/cosmos-sdk/client/utils"
	"github.com/gracenoah/cosmos-sdk/codec"
	sdk "github.com/gracenoah/cosmos-sdk/types"
	authtxb "github.com/gracenoah/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/gracenoah/cosmos-sdk/x/slashing"

	"github.com/spf13/cobra"
)

// GetCmdUnjail implements the create unjail validator command.
func GetCmdUnjail(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "unjail",
		Args:  cobra.NoArgs,
		Short: "unjail validator previously jailed for downtime",
		Long: `unjail a jailed validator:

$ gaiacli tx slashing unjail --from mykey
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(cdc)

			valAddr := cliCtx.GetFromAddress()

			msg := slashing.NewMsgUnjail(sdk.ValAddress(valAddr))
			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg}, false)
		},
	}
}
