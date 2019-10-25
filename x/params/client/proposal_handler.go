package client

import (
	govclient "github.com/gracenoah/cosmos-sdk/x/gov/client"
	"github.com/gracenoah/cosmos-sdk/x/params/client/cli"
	"github.com/gracenoah/cosmos-sdk/x/params/client/rest"
)

// param change proposal handler
var ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
