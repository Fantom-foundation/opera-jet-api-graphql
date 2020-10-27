/*
Package rpc implements bridge to Lachesis full node API interface.

We recommend using local IPC for fast and the most efficient inter-process communication between the API server
and an Opera/Lachesis node. Any remote RPC connection will work, but the performance may be significantly degraded
by extra networking overhead of remote RPC calls.

You should also consider security implications of opening Lachesis RPC interface for a remote access.
If you considering it as your deployment strategy, you should establish encrypted channel between the API server
and Lachesis RPC interface with connection limited to specified endpoints.

We strongly discourage opening Lachesis RPC interface for unrestricted Internet access.
*/
package rpc

//go:generate abigen --abi ./contracts/sfc-2.0.2-rc2.abi --pkg rpc --type SfcContract --out ./smc_sfc.go

import (
	"fantom-api-graphql/internal/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"math/big"
)

// sfcContractAddress represents the address on which the Sfc contract is deployed.
var sfcContractAddress = common.HexToAddress("0xfc00face00000000000000000000000000000000")

// SfcVersion returns current version of the SFC contract as a single number.
func (ftm *FtmBridge) SfcVersion() (hexutil.Uint64, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return 0, err
	}

	// get the version information from the contract
	var ver [3]byte
	ver, err = contract.Version(nil)
	if err != nil {
		ftm.log.Criticalf("failed to get the SFC version; %v", err)
		return 0, err
	}

	return hexutil.Uint64((uint64(ver[0]) << 16) | (uint64(ver[1]) << 8) | uint64(ver[2])), nil
}

// CurrentEpoch extract the current epoch id from SFC smart contract.
func (ftm *FtmBridge) CurrentEpoch() (hexutil.Uint64, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return 0, err
	}

	// get the value from the contract
	epoch, err := contract.CurrentEpoch(nil)
	if err != nil {
		ftm.log.Errorf("failed to get the current epoch: %v", err)
		return 0, err
	}

	// get the value
	return hexutil.Uint64(epoch.Uint64()), nil
}

// CurrentSealedEpoch extract the current sealed epoch id from SFC smart contract.
func (ftm *FtmBridge) CurrentSealedEpoch() (hexutil.Uint64, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return 0, err
	}

	// get the value from the contract
	epoch, err := contract.CurrentSealedEpoch(nil)
	if err != nil {
		ftm.log.Errorf("failed to get the current sealed epoch: %v", err)
		return 0, err
	}

	// get the value
	return hexutil.Uint64(epoch.Uint64()), nil
}

// LastStakerId returns the last staker id in Opera blockchain.
func (ftm *FtmBridge) LastStakerId() (hexutil.Uint64, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return 0, err
	}

	// get the value from the contract
	sl, err := contract.StakersLastID(nil)
	if err != nil {
		ftm.log.Errorf("failed to get the last staker ID: %v", err)
		return 0, err
	}

	// get the value
	return hexutil.Uint64(sl.Uint64()), nil
}

// StakersNum returns the number of stakers in Opera blockchain.
func (ftm *FtmBridge) StakersNum() (hexutil.Uint64, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return 0, err
	}

	// get the value from the contract
	sn, err := contract.StakersNum(nil)
	if err != nil {
		ftm.log.Errorf("failed to get the current number of stakers: %v", err)
		return 0, err
	}

	// get the value
	return hexutil.Uint64(sn.Uint64()), nil
}

// stakerStatusFromSfc updates staker information using SFC binding.
func (ftm *FtmBridge) stakerStatusFromSfc(contract *SfcContract, staker *types.Staker) error {
	// log action
	ftm.log.Debug("updating staker info from SFC")

	// get the value from the contract
	si, err := contract.Stakers(nil, big.NewInt(int64(staker.Id)))
	if err != nil {
		ftm.log.Errorf("failed to get the staker information from SFC: %v", err)
		return err
	}

	// do we have a valid record?
	if si.Status != nil {
		// update some invalid information
		staker.DelegatedMe = (*hexutil.Big)(si.DelegatedMe)
		staker.Stake = (*hexutil.Big)(si.StakeAmount)
		staker.Status = hexutil.Uint64(si.Status.Uint64())

		if staker.Stake != nil && staker.DelegatedMe != nil {
			// recalculate the total stake
			staker.TotalStake = (*hexutil.Big)(big.NewInt(0).Add(si.DelegatedMe, si.StakeAmount))

			// calculate delegation limit
			staker.TotalDelegatedLimit = ftm.maxDelegatedLimit(staker.Stake, contract)

			// calculate available limit for staking
			val := new(big.Int).Sub((*big.Int)(&staker.TotalDelegatedLimit), (*big.Int)(staker.DelegatedMe))
			staker.DelegatedLimit = (hexutil.Big)(*val)
		}
	} else {
		// log issue
		ftm.log.Debug("staker info update from SFC failed, no data received")
	}

	// get the value
	return nil
}

// stakerLockFromSfc updates staker lock details using SFC binding.
func (ftm *FtmBridge) stakerLockFromSfc(contract *SfcContract, staker *types.Staker) error {
	// log action
	ftm.log.Debug("updating staker locking details from SFC")

	// get staker locking detail
	lock, err := contract.LockedStakes(nil, big.NewInt(int64(staker.Id)))
	if err != nil {
		ftm.log.Errorf("stake lock query failed; %v", err)
		return nil
	}

	// are lock timers available?
	if lock.FromEpoch == nil || lock.EndTime == nil {
		ftm.log.Errorf("stake lock details not available")
		return nil
	}

	// apply the lock values
	staker.LockedFromEpoch = hexutil.Uint64(lock.FromEpoch.Uint64())
	staker.LockedUntil = hexutil.Uint64(lock.EndTime.Uint64())

	// get the value
	return nil
}

// maxDelegatedLimit calculate maximum amount of tokens allowed to be delegated to a staker.
func (ftm *FtmBridge) maxDelegatedLimit(staked *hexutil.Big, contract *SfcContract) hexutil.Big {
	// if we don't know the staked amount, return zero
	if staked == nil {
		return (hexutil.Big)(*hexutil.MustDecodeBig("0x0"))
	}

	// ratio unit is used to calculate the value (1.000.000)
	// please note this formula is taken from SFC contract and can change
	ratioUnit := hexutil.MustDecodeBig("0xF4240")

	// get delegation ration
	ratio, err := contract.MaxDelegatedRatio(nil)
	if err != nil {
		ftm.log.Errorf("can not get the delegation ratio; %s", err.Error())
		return (hexutil.Big)(*hexutil.MustDecodeBig("0x0"))
	}

	// calculate the delegation limit temp value
	temp := new(big.Int).Mul((*big.Int)(staked), ratio)

	// adjust to percent
	value := new(big.Int).Div(temp, ratioUnit)
	return (hexutil.Big)(*value)
}

// extendStaker extends staker information using SFC contract binding.
func (ftm *FtmBridge) extendStaker(staker *types.Staker) (*types.Staker, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return nil, err
	}

	// update status detail
	err = ftm.stakerStatusFromSfc(contract, staker)
	if err != nil {
		ftm.log.Critical("staker status could not be updated from SFC")
	}

	// update locking detail
	err = ftm.stakerLockFromSfc(contract, staker)
	if err != nil {
		ftm.log.Critical("staker locking could not be updated from SFC")
	}

	return staker, nil
}

// Staker extract a staker information by numeric id.
func (ftm *FtmBridge) Staker(id hexutil.Uint64) (*types.Staker, error) {
	// keep track of the operation
	ftm.log.Debugf("loading staker #%d", id)

	// call for data
	var st types.Staker
	err := ftm.rpc.Call(&st, "sfc_getStaker", id, "0x2")
	if err != nil {
		ftm.log.Error("staker information could not be extracted")
		return nil, err
	}

	// keep track of the operation
	ftm.log.Debugf("staker #%d loaded", id)
	return ftm.extendStaker(&st)
}

// StakerByAddress extracts a staker information by address.
func (ftm *FtmBridge) StakerByAddress(addr common.Address) (*types.Staker, error) {
	// keep track of the operation
	ftm.log.Debugf("loading staker %s", addr.String())

	// call for data
	var st types.Staker
	err := ftm.rpc.Call(&st, "sfc_getStakerByAddress", addr, "0x2")
	if err != nil {
		ftm.log.Error("staker information could not be extracted")
		return nil, err
	}

	// keep track of the operation
	ftm.log.Debugf("staker %s loaded", addr.String())
	return ftm.extendStaker(&st)
}

// Epoch extract information about an epoch from SFC smart contract.
func (ftm *FtmBridge) Epoch(id hexutil.Uint64) (types.Epoch, error) {
	// instantiate the contract and display its name
	contract, err := NewSfcContract(sfcContractAddress, ftm.eth)
	if err != nil {
		ftm.log.Criticalf("failed to instantiate SFC contract: %v", err)
		return types.Epoch{}, err
	}

	// extract epoch snapshot
	epo, err := contract.EpochSnapshots(nil, big.NewInt(int64(id)))
	if err != nil {
		ftm.log.Errorf("failed to extract epoch information: %v", err)
		return types.Epoch{}, err
	}

	return types.Epoch{
		Id:                     id,
		EndTime:                (hexutil.Big)(*epo.EndTime),
		Duration:               (hexutil.Big)(*epo.Duration),
		EpochFee:               (hexutil.Big)(*epo.EpochFee),
		TotalBaseRewardWeight:  (hexutil.Big)(*epo.TotalBaseRewardWeight),
		TotalTxRewardWeight:    (hexutil.Big)(*epo.TotalTxRewardWeight),
		BaseRewardPerSecond:    (hexutil.Big)(*epo.BaseRewardPerSecond),
		StakeTotalAmount:       (hexutil.Big)(*epo.StakeTotalAmount),
		DelegationsTotalAmount: (hexutil.Big)(*epo.DelegationsTotalAmount),
		TotalSupply:            (hexutil.Big)(*epo.TotalSupply),
	}, nil
}


