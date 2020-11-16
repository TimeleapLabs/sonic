package gossip

import (
	"context"
	"crypto/ecdsa"
	"math/big"
	"time"

	"github.com/Fantom-foundation/lachesis-base/inter/idx"
	"github.com/Fantom-foundation/lachesis-base/inter/pos"
	"github.com/Fantom-foundation/lachesis-base/kvdb/flushable"
	"github.com/Fantom-foundation/lachesis-base/kvdb/memorydb"
	"github.com/Fantom-foundation/lachesis-base/lachesis"
	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"

	"github.com/Fantom-foundation/go-opera/evmcore"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/eventmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/evmmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/sealmodule"
	"github.com/Fantom-foundation/go-opera/gossip/blockproc/sfcmodule"
	"github.com/Fantom-foundation/go-opera/inter"
	"github.com/Fantom-foundation/go-opera/opera"
)

const (
	maxEpochDuration = time.Hour
	sameEpoch        = maxEpochDuration / 1000
	nextEpoch        = maxEpochDuration
)

type testEnv struct {
	network   opera.Config
	blockProc BlockProc
	store     *Store

	lastBlock     idx.Block
	lastBlockTime time.Time
}

func newTestEnv() *testEnv {
	network := opera.MainNetConfig()
	network.Dag.MaxEpochDuration = maxEpochDuration

	dbs := flushable.NewSyncedPool(
		memorydb.NewProducer(""))

	env := &testEnv{
		network: network,
		blockProc: BlockProc{
			SealerModule:        sealmodule.New(network),
			TxListenerModule:    sfcmodule.NewSfcTxListenerModule(network),
			GenesisTxTransactor: sfcmodule.NewSfcTxGenesisTransactor(network),
			PreTxTransactor:     sfcmodule.NewSfcTxPreTransactor(network),
			PostTxTransactor:    sfcmodule.NewSfcTxTransactor(network),
			EventsModule:        eventmodule.New(network),
			EVMModule:           evmmodule.New(network),
		},
		store: NewStore(dbs, LiteStoreConfig()),
	}
	_, _, err := env.store.ApplyGenesis(env.blockProc, &network)
	if err != nil {
		panic(err)
	}

	return env
}

func (env *testEnv) Close() {
}

func (env *testEnv) consensusCallbackBeginBlockFn() lachesis.BeginBlockFn {
	return consensusCallbackBeginBlockFn(
		env.network,
		env.store,
		env.blockProc,
		false, nil, nil,
	)
}

func (env *testEnv) ApplyBlock(spent time.Duration, txs ...*eth.Transaction) eth.Receipts {
	env.lastBlock++
	env.lastBlockTime = env.lastBlockTime.Add(spent)

	eBuilder := inter.MutableEventPayload{}
	eBuilder.SetTxs(txs)
	// TODO: fill the event

	event := eBuilder.Build()
	env.store.SetEvent(event)

	beginBlock := env.consensusCallbackBeginBlockFn()
	process := beginBlock(&lachesis.Block{
		Atropos: event.ID(),
	})

	process.ApplyEvent(event)
	_ = process.EndBlock()

	return nil
}

func (env *testEnv) Transfer(from int, to int, amount *big.Int) *eth.Transaction {
	nonce, _ := env.PendingNonceAt(nil, env.Address(from))
	key := env.privateKey(from)
	receiver := env.Address(to)
	gp := env.App.MinGasPrice()
	tx := eth.NewTransaction(nonce, receiver, amount, gasLimit, gp, nil)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) Contract(from int, amount *big.Int, hex string) *eth.Transaction {
	data := hexutil.MustDecode(hex)
	nonce, _ := env.PendingNonceAt(nil, env.Address(from))
	key := env.privateKey(from)
	gp := env.App.MinGasPrice()
	tx := eth.NewContractCreation(nonce, amount, gasLimit*10000, gp, data)
	tx, err := eth.SignTx(tx, env.signer, key)
	if err != nil {
		panic(err)
	}

	return tx
}

func (env *testEnv) privateKey(n int) *ecdsa.PrivateKey {
	key := crypto.FakeKey(n)
	return key
}

func (env *testEnv) Address(n int) common.Address {
	key := crypto.FakeKey(n)
	addr := crypto.PubkeyToAddress(key.PublicKey)
	return addr
}

func (env *testEnv) Payer(n int, amounts ...*big.Int) *bind.TransactOpts {
	key := env.privateKey(n)
	t := bind.NewKeyedTransactor(key)
	nonce, _ := env.PendingNonceAt(nil, env.Address(n))
	t.Nonce = big.NewInt(int64(nonce))
	t.Value = big.NewInt(0)
	for _, amount := range amounts {
		t.Value.Add(t.Value, amount)
	}
	return t
}

func (env *testEnv) ReadOnly() *bind.CallOpts {
	return &bind.CallOpts{}
}

func (env *testEnv) State() evmcore.StateDB {
	return env.Store.StateDB(env.lastState)
}

func (env *testEnv) incNonce(account common.Address) {
	env.nonces[account] += 1
}

/*
 * bind.ContractCaller interface
 */

var (
	errBlockNumberUnsupported = errors.New("simulatedBackend cannot access blocks other than the latest block")
)

// CodeAt returns the code of the given account. This is needed to differentiate
// between contract internal errors and the local chain being out of sync.
func (env *testEnv) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	if blockNumber != nil && idx.Block(blockNumber.Uint64()) != env.lastBlock {
		return nil, errBlockNumberUnsupported
	}

	code := env.State().GetCode(contract)
	return code, nil
}

// ContractCall executes an Ethereum contract call with the specified data as the
// input.
func (env *testEnv) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	if blockNumber != nil && idx.Block(blockNumber.Uint64()) != env.lastBlock {
		return nil, errBlockNumberUnsupported
	}

	h := env.App.GetEvmStateReader().GetHeader(common.Hash{}, uint64(env.lastBlock))
	block := &evmcore.EvmBlock{
		EvmHeader: *h,
	}

	rval, _, _, err := env.callContract(ctx, call, block, env.State())
	return rval, err
}

// callContract implements common code between normal and pending contract calls.
// state is modified during execution, make sure to copy it if necessary.
func (env *testEnv) callContract(
	ctx context.Context, call ethereum.CallMsg, block *evmcore.EvmBlock, statedb evmcore.StateDB,
) (
	ret []byte, usedGas uint64, failed bool, err error,
) {
	// Ensure message is initialized properly.
	if call.GasPrice == nil {
		call.GasPrice = big.NewInt(1)
	}
	if call.Gas == 0 {
		call.Gas = 50000000
	}
	if call.Value == nil {
		call.Value = new(big.Int)
	}
	// Set infinite balance to the fake caller account.
	from := statedb.MPT().GetOrNewStateObject(call.From)
	from.SetBalance(big.NewInt(math.MaxInt64))

	msg := callmsg{call}

	evmContext := evmcore.NewEVMContext(msg, block.Header(), env.App.GetEvmStateReader(), &call.From)
	// Create a new environment which holds all relevant information
	// about the transaction and calling mechanisms.
	vmenv := vm.NewEVM(evmContext, statedb, env.App.config.Net.EvmChainConfig(), vm.Config{})
	gaspool := new(evmcore.GasPool).AddGas(math.MaxUint64)
	res, err := evmcore.NewStateTransition(vmenv, msg, gaspool).TransitionDb()

	ret, usedGas, failed = res.Return(), res.UsedGas, res.Failed()
	return
}

/*
 * ContractTransactor interface
 */

// PendingCodeAt returns the code of the given account in the pending state.
func (env *testEnv) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	code := env.State().GetCode(account)
	return code, nil
}

// PendingNonceAt retrieves the current pending nonce associated with an account.
func (env *testEnv) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	nonce := env.nonces[account]
	return nonce, nil
}

// SuggestGasPrice retrieves the currently suggested gas price to allow a timely
// execution of a transaction.
func (env *testEnv) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	return env.App.MinGasPrice(), nil
}

// EstimateGas tries to estimate the gas needed to execute a specific
// transaction based on the current pending state of the backend blockchain.
// There is no guarantee that this is the true gas limit requirement as other
// transactions may be added or removed by miners, but it should provide a basis
// for setting a reasonable default.
func (env *testEnv) EstimateGas(ctx context.Context, call ethereum.CallMsg) (gas uint64, err error) {
	if call.To == nil {
		gas = gasLimit * 10000
	} else {
		gas = gasLimit * 10
	}
	return
}

// SendTransaction injects the transaction into the pending pool for execution.
func (env *testEnv) SendTransaction(ctx context.Context, tx *eth.Transaction) error {
	// do nothing to avoid executing by transactor, only generating needed
	return nil
}

/*
 *  bind.ContractFilterer interface
 */

// FilterLogs executes a log filter operation, blocking during execution and
// returning all the results in one batch.
func (env *testEnv) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]eth.Log, error) {
	panic("not implemented yet")
	return nil, nil
}

// SubscribeFilterLogs creates a background log filtering operation, returning
// a subscription immediately, which can be used to stream the found events.
func (env *testEnv) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- eth.Log) (ethereum.Subscription, error) {
	panic("not implemented yet")
	return nil, nil
}

// callmsg implements core.Message to allow passing it as a transaction simulator.
type callmsg struct {
	ethereum.CallMsg
}

func (m callmsg) From() common.Address { return m.CallMsg.From }
func (m callmsg) Nonce() uint64        { return 0 }
func (m callmsg) CheckNonce() bool     { return false }
func (m callmsg) To() *common.Address  { return m.CallMsg.To }
func (m callmsg) GasPrice() *big.Int   { return m.CallMsg.GasPrice }
func (m callmsg) Gas() uint64          { return m.CallMsg.Gas }
func (m callmsg) Value() *big.Int      { return m.CallMsg.Value }
func (m callmsg) Data() []byte         { return m.CallMsg.Data }

// fakeEngine implements Consensus interface
type fakeEngine struct {
	env *testEnv
}

// LastBlock returns current block.
func (f *fakeEngine) LastBlock() (idx.Block, hash.Event) {
	return f.env.lastBlock, hash.Event{}
}

// GetEpoch returns current epoch num.
func (f *fakeEngine) GetEpoch() idx.Epoch {
	return f.env.epoch
}

// GetValidators returns validators of current epoch.
func (f *fakeEngine) GetValidators() *pos.Validators {
	return f.env.validators.Build()
}

// GetEpochValidators atomically returns validators of current epoch, and the epoch.
func (f *fakeEngine) GetEpochValidators() (*pos.Validators, idx.Epoch) {
	return f.env.validators.Build(), f.env.epoch
}

// PushEvent takes event for processing.
func (f *fakeEngine) ProcessEvent(e *inter.Event) error {
	panic("Not implemented!")
	return nil
}

// GetGenesisHash returns hash of genesis poset works with.
func (f *fakeEngine) GetGenesisHash() common.Hash {
	panic("Not implemented!")
	return common.Hash{}
}

// GetVectorIndex returns internal vector clock if exists
func (f *fakeEngine) GetVectorIndex() *vector.Index {
	panic("Not implemented!")
	return nil
}

// Sets consensus fields. Returns nil if event should be dropped.
func (f *fakeEngine) Prepare(e *inter.Event) *inter.Event {
	panic("Not implemented!")
	return e
}

// GetConsensusTime calc consensus timestamp for given event.
func (f *fakeEngine) GetConsensusTime(id hash.Event) (inter.Timestamp, error) {
	panic("Not implemented!")
	return 0, nil
}

// Bootstrap must be called (once) before calling other methods
func (f *fakeEngine) Bootstrap(callbacks inter.ConsensusCallbacks) {
	panic("Not implemented!")
}
