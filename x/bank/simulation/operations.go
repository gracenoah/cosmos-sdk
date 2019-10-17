package simulation

import (
	"errors"
	"math/rand"

	"github.com/tendermint/tendermint/crypto"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/simapp/helpers"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank/internal/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/internal/types"
	"github.com/cosmos/cosmos-sdk/x/simulation"
)

// SimulateMsgSend tests and runs a single msg send where both
// accounts already exist.
// nolint: funlen
func SimulateMsgSend(ak types.AccountKeeper, bk keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		if !bk.GetSendEnabled(ctx) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		simAccount, toSimAcc, coins, skip, err := randomSendFields(r, ctx, accs, ak)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		if skip {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		msg := types.NewMsgSend(simAccount.Address, toSimAcc.Address, coins)

		err = sendMsgSend(r, app, ak, msg, ctx, chainID, []crypto.PrivKey{simAccount.PrivKey})
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// sendMsgSend sends a transaction with a MsgSend from a provided random account.
func sendMsgSend(
	r *rand.Rand, app *baseapp.BaseApp, ak types.AccountKeeper,
	msg types.MsgSend, ctx sdk.Context, chainID string, privkeys []crypto.PrivKey,
) (err error) {

	account := ak.GetAccount(ctx, msg.FromAddress)
	coins := account.SpendableCoins(ctx.BlockTime())

	var fees sdk.Coins
	coins, hasNeg := coins.SafeSub(msg.Amount)
	if !hasNeg {
		fees, err = simulation.RandomFees(r, ctx, coins)
		if err != nil {
			return err
		}
	}

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fees,
		chainID,
		[]uint64{account.GetAccountNumber()},
		[]uint64{account.GetSequence()},
		privkeys...,
	)

	res := app.Deliver(tx)
	if !res.IsOK() {
		return errors.New(res.Log)
	}

	return nil
}

// SimulateMsgMultiSend tests and runs a single msg multisend, with randomized, capped number of inputs/outputs.
// all accounts in msg fields exist in state
// nolint: funlen
func SimulateMsgMultiSend(ak types.AccountKeeper, bk keeper.Keeper) simulation.Operation {
	return func(
		r *rand.Rand, app *baseapp.BaseApp, ctx sdk.Context,
		accs []simulation.Account, chainID string,
	) (simulation.OperationMsg, []simulation.FutureOperation, error) {

		if !bk.GetSendEnabled(ctx) {
			return simulation.NoOpMsg(types.ModuleName), nil, nil
		}

		// random number of inputs/outputs between [1, 3]
		inputs := make([]types.Input, rand.Intn(3)+1)
		outputs := make([]types.Output, rand.Intn(3)+1)

		// collect signer privKeys
		privs := make([]crypto.PrivKey, len(inputs))

		// use map to check if address already exists as input
		usedAddrs := make(map[string]bool)

		var totalSentCoins sdk.Coins
		for i := range inputs {
			// generate random input fields, ignore to address
			simAccount, _, coins, skip, err := randomSendFields(r, ctx, accs, ak)

			// make sure account is fresh and not used in previous input
			for usedAddrs[simAccount.Address.String()] {
				simAccount, _, coins, skip, err = randomSendFields(r, ctx, accs, ak)
			}

			if err != nil {
				return simulation.NoOpMsg(types.ModuleName), nil, err
			}
			if skip {
				return simulation.NoOpMsg(types.ModuleName), nil, nil
			}

			// set input address in used address map
			usedAddrs[simAccount.Address.String()] = true

			// set signer privkey
			privs[i] = simAccount.PrivKey

			// set next input and accumulate total sent coins
			inputs[i] = types.NewInput(simAccount.Address, coins)
			totalSentCoins.Add(coins)
		}

		for o := range outputs {
			outAddr, _ := simulation.RandomAcc(r, accs)

			var outCoins sdk.Coins
			// split total sent coins into random subsets for output
			if o == len(outputs)-1 {
				outCoins = totalSentCoins
			} else {
				// take random subset of remaining coins for output
				// and update remaining coins
				outCoins = simulation.RandSubsetCoins(r, totalSentCoins)
				totalSentCoins = totalSentCoins.Sub(outCoins)
			}

			outputs[o] = types.NewOutput(outAddr.Address, outCoins)
		}

		msg := types.MsgMultiSend{
			Inputs:  inputs,
			Outputs: outputs,
		}

		err := sendMsgMultiSend(r, app, ak, msg, ctx, chainID, privs)
		if err != nil {
			return simulation.NoOpMsg(types.ModuleName), nil, err
		}

		return simulation.NewOperationMsg(msg, true, ""), nil, nil
	}
}

// sendMsgMultiSend sends a transaction with a MsgMultiSend from a provided random
// account.
func sendMsgMultiSend(
	r *rand.Rand, app *baseapp.BaseApp, ak types.AccountKeeper,
	msg types.MsgMultiSend, ctx sdk.Context, chainID string, privkeys []crypto.PrivKey,
) (err error) {

	accountNumbers := make([]uint64, len(msg.Inputs))
	sequenceNumbers := make([]uint64, len(msg.Inputs))

	for i := 0; i < len(msg.Inputs); i++ {
		acc := ak.GetAccount(ctx, msg.Inputs[i].Address)
		accountNumbers[i] = acc.GetAccountNumber()
		sequenceNumbers[i] = acc.GetSequence()
	}

	// feePayer is the first signer, i.e. first input address
	feePayer := ak.GetAccount(ctx, msg.Inputs[0].Address)
	coins := feePayer.SpendableCoins(ctx.BlockTime())

	var fees sdk.Coins
	coins, hasNeg := coins.SafeSub(msg.Inputs[0].Coins)
	if !hasNeg {
		fees, err = simulation.RandomFees(r, ctx, coins)
		if err != nil {
			return err
		}
	}

	tx := helpers.GenTx(
		[]sdk.Msg{msg},
		fees,
		chainID,
		accountNumbers,
		sequenceNumbers,
		privkeys...,
	)

	res := app.Deliver(tx)
	if !res.IsOK() {
		return errors.New(res.Log)
	}

	return nil
}

// randomSendFields returns the sender and recipient simulation accounts as well
// as the transferred amount.
func randomSendFields(
	r *rand.Rand, ctx sdk.Context, accs []simulation.Account, ak types.AccountKeeper,
) (simulation.Account, simulation.Account, sdk.Coins, bool, error) {

	simAccount, _ := simulation.RandomAcc(r, accs)
	toSimAcc, _ := simulation.RandomAcc(r, accs)

	// disallow sending money to yourself
	for simAccount.PubKey.Equals(toSimAcc.PubKey) {
		toSimAcc, _ = simulation.RandomAcc(r, accs)
	}

	acc := ak.GetAccount(ctx, simAccount.Address)
	if acc == nil {
		return simAccount, toSimAcc, nil, true, nil // skip error
	}

	coins := acc.SpendableCoins(ctx.BlockHeader().Time)
	if coins.Empty() {
		return simAccount, toSimAcc, nil, true, nil // skip error
	}

	sendCoins := simulation.RandSubsetCoins(r, coins)

	return simAccount, toSimAcc, sendCoins, false, nil
}
