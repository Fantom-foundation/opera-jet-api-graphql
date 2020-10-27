/*
Package repository implements repository for handling fast and efficient access to data required
by the resolvers of the API server.

Internally it utilizes RPC to access Opera/Lachesis full node for blockchain interaction. Mongo database
for fast, robust and scalable off-chain data storage, especially for aggregated and pre-calculated data mining
results. BigCache for in-memory object storage to speed up loading of frequently accessed entities.
*/
package repository

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// SfcVersion returns current version of the SFC contract.
func (p *proxy) SfcVersion() (hexutil.Uint64, error) {
	return p.rpc.SfcVersion()
}

// CurrentEpoch returns the id of the current epoch.
func (p *proxy) CurrentEpoch() (hexutil.Uint64, error) {
	return p.rpc.CurrentEpoch()
}

// Epoch returns the id of the current epoch.
func (p *proxy) Epoch(id hexutil.Uint64) (types.Epoch, error) {
	return p.rpc.Epoch(id)
}

// LastStakerId returns the last staker id in Opera blockchain.
func (p *proxy) LastStakerId() (hexutil.Uint64, error) {
	return p.rpc.LastStakerId()
}

// StakersNum returns the number of stakers in Opera blockchain.
func (p *proxy) StakersNum() (hexutil.Uint64, error) {
	return p.rpc.StakersNum()
}

// CurrentSealedEpoch returns the data of the latest sealed epoch.
// This is used for reward estimation calculation and we don't need
// real time data, but rather faster response time.
// So, we use cache for handling the response.
// It will not be updated in sync with the SFC contract.
// If you need real time response, please use the Epoch(id) function instead.
func (p *proxy) CurrentSealedEpoch() (*types.Epoch, error) {
	// inform what we do
	p.log.Debug("latest sealed epoch requested")

	// try to use the in-memory cache
	if ep := p.cache.PullLastEpoch(); ep != nil {
		// inform what we do
		p.log.Debug("latest sealed epoch loaded from cache")

		// return the block
		return ep, nil
	}

	// we need to go the slow path
	id, err := p.rpc.CurrentSealedEpoch()
	if err != nil {
		// inform what we do
		p.log.Errorf("can not get the id of the last sealed epoch; %s", err.Error())
		return nil, err
	}

	// get the epoch from SFC
	ep, err := p.rpc.Epoch(id)
	if err != nil {
		// inform what we do
		p.log.Errorf("can not get data of the last sealed epoch; %s", err.Error())
		return nil, err
	}

	// try to store the block in cache for future use
	err = p.cache.PushLastEpoch(&ep)
	if err != nil {
		p.log.Error(err)
	}

	// inform what we do
	p.log.Debugf("epoch [%s] loaded from sfc", id.String())
	return &ep, nil
}
