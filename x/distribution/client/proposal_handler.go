package client

import (
	"github.com/gracenoah/cosmos-sdk/x/distribution/client/cli"
	"github.com/gracenoah/cosmos-sdk/x/distribution/client/rest"
	govclient "github.com/gracenoah/cosmos-sdk/x/gov/client"
)

// param change proposal handler
var (
	ProposalHandler = govclient.NewProposalHandler(cli.GetCmdSubmitProposal, rest.ProposalRESTHandler)
)
