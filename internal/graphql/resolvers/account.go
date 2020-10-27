// Package resolvers implements GraphQL resolvers to incoming API requests.
package resolvers

import (
	"fantom-api-graphql/internal/repository"
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// accMaxTransactionsPerRequest maximal number of transaction end-client can request in one query.
const accMaxTransactionsPerRequest = 50

// Account represents resolvable blockchain account structure.
type Account struct {
	repo          repository.Repository
	rfDelegations []types.Delegation
	rfStaker      *types.Staker
	rfBalance     *hexutil.Big

	/* extended delegated amounts pre-loaded */
	dlExtendedAmount           *big.Int
	dlExtendedAmountInWithdraw *big.Int

	types.Account
}

// NewAccount builds new resolvable account structure.
func NewAccount(acc *types.Account, repo repository.Repository) *Account {
	return &Account{
		repo:    repo,
		Account: *acc,
	}
}

// Account resolves blockchain account by address.
func (rs *rootResolver) Account(args struct{ Address common.Address }) (*Account, error) {
	// simply pull the block by hash
	acc, err := rs.repo.Account(&args.Address)
	if err != nil {
		rs.log.Errorf("could not get the specified account")
		return nil, err
	}

	return NewAccount(acc, rs.repo), nil
}

// Resolves total number of active accounts on the blockchain.
func (rs *rootResolver) AccountsActive() (hexutil.Uint64, error) {
	return rs.repo.AccountsActive()
}

// Balance resolves total balance of the account.
func (acc *Account) Balance() (hexutil.Big, error) {
	if acc.rfBalance == nil {
		// get the sender by address
		bal, err := acc.repo.AccountBalance(&acc.Account)
		if err != nil {
			return hexutil.Big{}, err
		}

		acc.rfBalance = bal
	}

	return *acc.rfBalance, nil
}

// TxCount resolves the number of transaction sent by the account, also known as nonce.
func (acc *Account) TxCount() (hexutil.Uint64, error) {
	// get the sender by address
	bal, err := acc.repo.AccountNonce(&acc.Account)
	if err != nil {
		return hexutil.Uint64(0), err
	}

	return *bal, nil
}

// TxList resolves list of transaction associated with the account.
func (acc *Account) TxList(args struct {
	Cursor *Cursor
	Count  int32
}) (*TransactionList, error) {
	// limit query size; the count can be either positive or negative
	// this controls the loading direction
	args.Count = listLimitCount(args.Count, accMaxTransactionsPerRequest)

	// get the transaction hash list from repository
	bl, err := acc.repo.AccountTransactions(&acc.Account, (*string)(args.Cursor), args.Count)
	if err != nil {
		return nil, err
	}

	return NewTransactionList(bl, acc.repo), nil
}

// Contract resolves the account smart contract detail,
// if the account is a smart contract address.
func (acc *Account) Contract() (*Contract, error) {
	// is this actually a contract account?
	if acc.ContractTx == nil {
		return nil, nil
	}

	// get new contract
	con, err := acc.repo.Contract(&acc.Address)
	if err != nil {
		return nil, err
	}

	return NewContract(con, acc.repo), nil
}