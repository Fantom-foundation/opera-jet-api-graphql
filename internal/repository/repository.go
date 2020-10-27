/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Internally it utilizes RPC to access Opera/Lachesis full node for blockchain interaction. Mongo database
for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"fantom-api-graphql/internal/config"
	"fantom-api-graphql/internal/logger"
	"fantom-api-graphql/internal/repository/cache"
	"fantom-api-graphql/internal/repository/db"
	"fantom-api-graphql/internal/repository/rpc"
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/compiler"
	"github.com/ethereum/go-ethereum/common/hexutil"
	ftm "github.com/ethereum/go-ethereum/rpc"
)

// Repository interface defines functions the underlying implementation provides to API resolvers.
type Repository interface {
	// FtmConnection returns open connection to Opera/Lachesis full node.
	FtmConnection() *ftm.Client

	// Account returns account at Opera blockchain for an address, nil if not found.
	Account(*common.Address) (*types.Account, error)

	// AccountBalance returns the current balance of an account at Opera blockchain.
	AccountBalance(*types.Account) (*hexutil.Big, error)

	// AccountNonce returns the current number of sent transactions of an account at Opera blockchain.
	AccountNonce(*types.Account) (*hexutil.Uint64, error)

	// AccountTransactions returns list of transaction hashes for account at Opera blockchain.
	//
	// String cursor represents cursor based on which the list is loaded. If null,
	// it loads either from top, or bottom of the list, based on the value
	// of the integer count. The integer represents the number of transaction loaded at most.
	//
	// For positive number, the list starts right after the cursor
	// (or on top without one) and loads at most defined number of transactions older than that.
	//
	// For negative number, the list starts right before the cursor
	// (or at the bottom without one) and loads at most defined number
	// of transactions newer than that.
	//
	// Transactions are always sorted from newer to older.
	AccountTransactions(*types.Account, *string, int32) (*types.TransactionHashList, error)

	// Returns total number of accounts known to repository.
	AccountsActive() (hexutil.Uint64, error)

	// Block returns a block at Opera blockchain represented by a number.
	// Top block is returned if the number is not provided.
	// If the block is not found, ErrBlockNotFound error is returned.
	BlockByNumber(*hexutil.Uint64) (*types.Block, error)

	// BlockHeight returns the current height of the Opera blockchain in blocks.
	BlockHeight() (*hexutil.Big, error)

	// LastKnownBlock returns number of the last block known to the repository.
	LastKnownBlock() (uint64, error)

	// CurrentEpoch returns the id of the current epoch.
	CurrentEpoch() (hexutil.Uint64, error)

	// CurrentSealedEpoch returns the data of the latest sealed epoch.
	CurrentSealedEpoch() (*types.Epoch, error)

	// Epoch returns the id of the current epoch.
	Epoch(hexutil.Uint64) (types.Epoch, error)

	// Block returns a block at Opera blockchain represented by a hash.
	// Top block is returned if the hash is not provided.
	// If the block is not found, ErrBlockNotFound error is returned.
	BlockByHash(*types.Hash) (*types.Block, error)

	// AddTransaction notifies a new incoming transaction from blockchain to the repository.
	AddTransaction(*types.Block, *types.Transaction) error

	// Transaction returns a transaction at Opera blockchain by a hash, nil if not found.
	Transaction(*types.Hash) (*types.Transaction, error)

	// TransactionsCount returns total number of transactions in the block chain.
	TransactionsCount() (hexutil.Uint64, error)

	// Transactions returns list of transaction hashes at Opera blockchain.
	Transactions(*string, int32) (*types.TransactionHashList, error)

	// Collection pulls list of blocks starting on the specified block number
	// and going up, or down based on count number.
	Blocks(*uint64, int32) (*types.BlockList, error)

	// LastStakerId returns the last staker id in Opera blockchain.
	LastStakerId() (hexutil.Uint64, error)

	// StakersNum returns the number of stakers in Opera blockchain.
	StakersNum() (hexutil.Uint64, error)

	// SfcVersion returns current version of the SFC contract.
	SfcVersion() (hexutil.Uint64, error)

	// SendTransaction sends raw signed and RLP encoded transaction to the block chain.
	SendTransaction(hexutil.Bytes) (*types.Transaction, error)

	// SetBlockChannel registers a channel for notifying new block events.
	SetBlockChannel(chan *types.Block)

	// SetTrxChannel registers a channel for notifying new transaction events.
	SetTrxChannel(chan *types.Transaction)

	// Contract extract a smart contract information by address if available.
	Contract(*common.Address) (*types.Contract, error)

	// Contracts returns list of smart contracts at Opera blockchain.
	Contracts(bool, *string, int32) (*types.ContractList, error)

	// ValidateContract tries to validate contract byte code using
	// provided source code. If successful, the contract information
	// is updated the the repository.
	ValidateContract(*types.Contract) error

	// Close and cleanup the repository.
	Close()
}

// Proxy represents Repository interface implementation and controls access to data
// trough several low level bridges.
type proxy struct {
	cache *cache.MemBridge
	db    *db.MongoDbBridge
	rpc   *rpc.FtmBridge
	log   logger.Logger

	// smart contract compilers
	solCompiler string

	// official ballot source addresses
	ballotSources []string

	// service orchestrator reference
	orc *orchestrator
}

// New creates new instance of Repository implementation, namely proxy structure.
func New(cfg *config.Config, log logger.Logger) (Repository, error) {
	// create new in-memory cache bridge
	caBridge, err := cache.New(cfg, log)
	if err != nil {
		log.Criticalf("can not create in-memory cache bridge, %s", err.Error())
		return nil, err
	}

	// create new database connection bridge
	dbBridge, err := db.New(cfg, log)
	if err != nil {
		log.Criticalf("can not connect backend persistent storage, %s", err.Error())
		return nil, err
	}

	// create new Lachesis RPC bridge
	rpcBridge, err := rpc.New(cfg, log)
	if err != nil {
		log.Criticalf("can not connect Lachesis RPC interface, %s", err.Error())
		return nil, err
	}

	// try to validate the solidity compiler by asking for it's version
	if _, err := compiler.SolidityVersion(cfg.SolCompilerPath); err != nil {
		log.Criticalf("can not invoke the Solidity compiler, %s", err.Error())
		return nil, err
	}

	// construct the proxy instance
	p := proxy{
		cache: caBridge,
		db:    dbBridge,
		rpc:   rpcBridge,
		log:   log,

		// keep reference to the SOL compiler
		solCompiler: cfg.SolCompilerPath,

		// keep the ballot sources ref
		ballotSources: cfg.VotingSources,
	}

	// inform about voting sources
	log.Infof("voting ballots accepted from %s", cfg.VotingSources)

	// propagate callbacks
	dbBridge.SetBalance(p.AccountBalance)

	// make the service orchestrator
	p.orc = newOrchestrator(&p, log)

	// return the proxy
	return &p, nil
}

// Close with close all connections and clean up the pending work for graceful termination.
func (p *proxy) Close() {
	// inform about actions
	p.log.Notice("repository is closing")

	// initiate orchestrator closing process
	p.orc.close()

	// close connections
	p.db.Close()
	p.rpc.Close()

	// inform about actions
	p.log.Notice("repository done")
}

// FtmClient returns open connection to Opera/Lachesis full node.
func (p *proxy) FtmConnection() *ftm.Client {
	return p.rpc.Connection()
}

// SetBlockChannel registers a channel for notifying new block events.
func (p *proxy) SetBlockChannel(ch chan *types.Block) {
	p.orc.setBlockChannel(ch)
}

// SetTrxChannel registers a channel for notifying new transactions events.
func (p *proxy) SetTrxChannel(ch chan *types.Transaction) {
	p.orc.setTrxChannel(ch)
}
