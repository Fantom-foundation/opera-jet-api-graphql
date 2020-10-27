// Package resolvers implements GraphQL resolvers to incoming API requests.
package resolvers

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// CurrentEpoch resolves the id of the current epoch of the Opera blockchain.
func (rs *rootResolver) CurrentEpoch() (hexutil.Uint64, error) {
	return rs.repo.CurrentEpoch()
}

// Epoch resolves information about epoch of the given id.
func (rs *rootResolver) Epoch(args *struct{ Id hexutil.Uint64 }) (types.Epoch, error) {
	return rs.repo.Epoch(args.Id)
}

// Resolves the last staker id in Opera blockchain.
func (rs *rootResolver) LastStakerId() (hexutil.Uint64, error) {
	return rs.repo.LastStakerId()
}

// Resolves the number of stakers in Opera blockchain.
func (rs *rootResolver) StakersNum() (hexutil.Uint64, error) {
	return rs.repo.StakersNum()
}
