package tags

import (
	sdk "github.com/gracenoah/cosmos-sdk/types"
)

// Slashing tags
var (
	ActionValidatorUnjailed = "validator-unjailed"

	Action    = sdk.TagAction
	Validator = "validator"
)
