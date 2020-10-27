package gqlschema

// Auto generated GraphQL schema bundle; created 2020-10-09 12:18
const schema = `
# TransactionList is a list of transaction edges provided by sequential access request.
type TransactionList {
    # Edges contains provided edges of the sequential list.
    edges: [TransactionListEdge!]!

    # TotalCount is the maximum number of transactions available for sequential access.
    totalCount: BigInt!

    # PageInfo is an information about the current page of transaction edges.
    pageInfo: ListPageInfo!
}

# TransactionListEdge is a single edge in a sequential list of transactions.
type TransactionListEdge {
    cursor: Cursor!
    transaction: Transaction!
}


# BlockList is a list of block edges provided by sequential access request.
type BlockList {
    # Edges contains provided edges of the sequential list.
    edges: [BlockListEdge!]!

    # TotalCount is the maximum number of blocks available for sequential access.
    totalCount: BigInt!

    # PageInfo is an information about the current page of block edges.
    pageInfo: ListPageInfo!
}

# BlockListEdge is a single edge in a sequential list of blocks.
type BlockListEdge {
    cursor: Cursor!
    block: Block!
}

# StakerInfo represents extended staker information from smart contract.
type StakerInfo {
    "Name represents the name of the staker."
    name: String

    "LogoUrl represents staker logo URL."
    logoUrl: String

    "Website represents a link to stakers website."
    website: String

    "Contact represents a link to contact to the staker."
    contact: String
}

# ListPageInfo contains information about a sequential access list page.
type ListPageInfo {
    # First is the cursor of the first edge of the edges list. null for empty list.
    first: Cursor

    # Last if the cursor of the last edge of the edges list. null for empty list.
    last: Cursor

    # HasNext specifies if there is another edge after the last one.
    hasNext: Boolean!

    # HasNext specifies if there is another edge before the first one.
    hasPrevious: Boolean!
}

# Transaction is an Opera block chain transaction.
type Transaction {
    # Hash is the unique hash of this transaction.
    hash: Hash!

    # Nonce is the number of transactions sent by the account prior to this transaction.
    nonce: Long!

    # Index is the index of this transaction in the block. This will
    # be null if the transaction is in a pending pool.
    index: Long

    # From is the address of the account that sent this transaction
    from: Address!

    # Sender is the account that sent this transaction
    sender: Account!

    # To is the account the transaction was sent to.
    # This is null for contract creating transactions.
    to: Address

    # contractAddress represents the address of smart contract deployed by this transaction;
    # null if the transaction is not contract creation
    contractAddress: Address

    # Recipient is the account that received this transaction.
    # Null for contract creating transaction.
    recipient: Account

    # Value is the value sent along with this transaction in WEI.
    value: BigInt!

    # GasUsed is the amount of gas that was used on processing this transaction.
    # If the transaction is pending, this field will be null.
    gasUsed: Long

    # InputData is the data supplied to the target of the transaction.
    # Contains smart contract byte code if this is contract creation.
    # Contains encoded contract state mutating function call if recipient
    # is a contract address.
    inputData: Bytes!

    # BlockHash is the hash of the block this transaction was assigned to.
    # Null if the transaction is pending.
    blockHash: Hash

    # BlockHash is the hash of the block this transaction was assigned to.
    # Null if the transaction is pending.
    blockNumber: Long

    # Block is the block this transaction was assigned to. This will be null if
    # the transaction is pending.
    block: Block

    # Status is the return status of the transaction. This will be 1 if the
    # transaction succeeded, or 0 if it failed (due to a revert, or due to
    # running out of gas). If the transaction has not yet been processed, this
    # field will be null.
    status: Long
}

# Block is an Opera block chain block.
type Block {
    # Number is the number of this block, starting at 0 for the genesis block.
    number: Long!

    # Hash is the unique block hash of this block.
    hash: Hash!

    # Parent is the parent block of this block.
    parent: Block

    # TransactionCount is the number of transactions in this block.
    transactionCount: Int

    # Timestamp is the unix timestamp at which this block was mined.
    timestamp: Long!

    # GasLimit represents the maximum gas allowed in this block.
    gasLimit: Long!

    # GasUsed represents the actual total used gas by all transactions in this block.
    gasUsed: Long!

    # txHashList is the list of unique hash values of transaction
    # assigned to the block.
    txHashList: [Hash!]!

    # txList is a list of transactions assigned to the block.
    txList: [Transaction!]!
}

# Represents epoch information.
type Epoch {
    "Number the epoch end."
    id: Long!

    "Timestamp of the epoch end."
    endTime: BigInt!

    "Epoch duration in seconds."
    duration: BigInt!

    "Fee at the epoch."
    epochFee: BigInt!

    "Total base reward weight on epoch."
    totalBaseRewardWeight: BigInt!

    "Total transaction reward weight on epoch."
    totalTxRewardWeight: BigInt!

    "Base reward per second of epoch."
    baseRewardPerSecond: BigInt!

    "Total amount staked."
    stakeTotalAmount: BigInt!

    "Total amount delegated."
    delegationsTotalAmount: BigInt!

    "Total supply amount."
    totalSupply: BigInt!
}

# Contract defines block-chain smart contract information container
type Contract {
    "Address represents the contract address."
    address: Address!

    "DeployedBy represents the smart contract deployment transaction reference."
    deployedBy: Transaction!

    "transactionHash represents the smart contract deployment transaction hash."
    transactionHash: Hash!

    "Smart contract name. Empty if not available."
    name: String!

    "Smart contract version identifier. Empty if not available."
    version: String!

    """
    License specifies an open source license the contract was published with.
    Empty if not specified.
    """
    license: String!

    "Smart contract author contact. Empty if not available."
    supportContact: String!

    "Smart contract compiler identifier. Empty if not available."
    compiler: String!

    "Smart contract source code. Empty if not available."
    sourceCode: String!

    "Smart contract ABI definition. Empty if not available."
    abi: String!

    """
    Validated is the unix timestamp at which the source code was validated
    against the deployed byte code. Null if not validated yet.
    """
    validated: Long

    "Timestamp is the unix timestamp at which this smart contract was deployed."
    timestamp: Long!
}

# ContractValidationInput represents a set of data sent from client
# to validate deployed contract with the provided source code.
input ContractValidationInput {
    "Address of the contract being validated."
    address: Address!

    "Optional smart contract name. Maximum allowed length is 64 characters."
    name: String

    "Optional smart contract version identifier. Maximum allowed length is 14 characters."
    version: String

    "Optional smart contract author contact. Maximum allowed length is 64 characters."
    supportContact: String

    """
    License specifies an open source license the contract was published with.
    Empty if not specified.
    """
    license: String

    "Optimized specifies if the compiler was set to optimize the byte code."
    optimized: Boolean = true

    """
    OptimizeRuns specifies number of optimization runs the compiler was set
    to execute during the byte code optimizing.
    """
    optimizeRuns: Int = 200

    "Smart contract source code."
    sourceCode: String!
}

# ContractList is a list of smart contract edges provided by sequential access request.
type ContractList {
    # Edges contains provided edges of the sequential list.
    edges: [ContractListEdge!]!

    # TotalCount is the maximum number of contracts available for sequential access.
    totalCount: BigInt!

    # PageInfo is an information about the current page of contract edges.
    pageInfo: ListPageInfo!
}

# TransactionListEdge is a single edge in a sequential list of transactions.
type ContractListEdge {
    cursor: Cursor!
    contract: Contract!
}

# Hash is a 32 byte binary string, represented by 0x prefixed hexadecimal.
scalar Hash

# Address is a 20 byte Opera address, represented as 0x prefixed hexadecimal number.
scalar Address

# BigInt is a large integer value. Input is accepted as either a JSON number,
# or a hexadecimal string alternatively prefixed with 0x. Output is 0x prefixed hexadecimal.
scalar BigInt

# Long is a 64 bit unsigned integer value.
scalar Long

# Bytes is an arbitrary length binary string, represented as 0x-prefixed hexadecimal.
# An empty byte string is represented as '0x'.
scalar Bytes

# Cursor is a string representing position in a sequential list of edges.
scalar Cursor

# CurrentState represents the current active state
# of the chain information condensed on one place.
type CurrentState {
    # epoch is the last sealed Epoch structure
    sealedEpoch: Epoch!

    # blocks represents number of blocks in the chain.
    blocks: Long!

    # transactions represents number of transactions in the chain.
    transactions: Long!

    # validators represents number of validators in the network.
    validators: Long!

    # accounts represents number of accounts participating on transactions.
    accounts: Long!

    # sfcVersion indicates the current version of the SFC contract.
    # The version is encoded into 3 bytes representing ASCII version numbers
    # with the most significant byte first [<8bit major><8bit minor><8bit revision>].
    # I.e. Version 1.0.2 = "102" = 0x313032
    sfcVersion: Long!
}

# Represents staker information.
type Staker {
    "Id number the staker."
    id: Long!

    "Staker address."
    stakerAddress: Address!

    "Amount of total staked tokens in WEI."
    totalStake: BigInt

    "Amount of own staked tokens in WEI."
    stake: BigInt

    "Amount of tokens delegated to the staker in WEI."
    delegatedMe: BigInt

    """
    Maximum total amount of tokens allowed to be delegated
    to the staker in WEI.
    This value depends on the amount of self staked tokens.
    """
    totalDelegatedLimit: BigInt!

    """
    Maximum amount of tokens allowed to be delegated to the staker
    on a new delegation in WEI.
    This value depends on the amount of self staked tokens.
    """
    delegatedLimit: BigInt!

    "Is this a validator record."
    isValidator: Boolean!

    "Is the staker active."
    isActive: Boolean!

    "Is the staker considered to be cheater."
    isCheater: Boolean!

    "Is the staker offline."
    isOffline: Boolean!

    "isStakeLocked signals if the staker locked the stake."
    isStakeLocked: Boolean!

    "Epoch in which the staker was created."
    createdEpoch: Long!

    "Timestamp of the staker creation."
    createdTime: Long!

    "lockedFromEpoch is the identifier of the epoch the stake lock was created."
    lockedFromEpoch: Long!

    "lockedUntil is the timestamp up to which the stake is locked, zero if not locked."
    lockedUntil: Long!

    "Epoch in which the staker was deactivated."
    deactivatedEpoch: Long!

    "Timestamp of the staker deactivation."
    deactivatedTime: Long!

    "How many blocks the staker missed."
    missedBlocks: Long!

    "Number of seconds the staker is offline."
    downtime: Long!

    "Proof of importance score."
    poi: BigInt

    "Base weight for rewards distribution."
    baseRewardWeight: BigInt

    "Weight for transaction rewards distribution."
    txRewardWeight: BigInt

    "Validation score."
    validationScore: BigInt

    "Origination score."
    originationScore: BigInt

    "Amount of rewards claimed in WEI."
    claimedRewards: BigInt

    "Amount of rewards claimed by delegators in WEI."
    delegationClaimedRewards: BigInt

    """
    Status is a binary encoded status of the staker.
    Ok = 0, bin 1 = Fork Detected, bin 256 = Validator Offline
    """
    status: Long!

    "StakerInfo represents extended staker information from smart contract."
    stakerInfo: StakerInfo

    """
    List of withdraw requests of the stake.
    Contains only withdrawal requests of the staking account,
    not the requests of the stake delegators.
    """
    withdrawRequests: [WithdrawRequest!]!
}


# Account defines block-chain account information container
type Account {
    "Address is the address of the account."
    address: Address!

    "Balance is the current balance of the Account in WEI."
    balance: BigInt!

    """
    TotalValue is the current total value fo the account in WEI.
    It includes available balance,
    delegated amount and pending staking rewards.
    """
    totalValue: BigInt!

    "txCount represents number of transaction sent from the account."
    txCount: Long!

    """
    txList represents list of transactions of the account
    in form of TransactionList.
    """
    txList (cursor:Cursor, count:Int!): TransactionList!

    "Details of a staker, if the account is a staker."
    staker: Staker

    "Details about smart contract, if the account is a smart contract."
    contract: Contract
}

# Root schema definition
schema {
    query: Query
    mutation: Mutation
    subscription: Subscription
}

# Entry points for querying the API
type Query {
    "State represents the current state of the blockchain and network."
    state: CurrentState!

    "Total number of accounts active on the Opera blockchain."
    accountsActive:Long!

    "Get an Account information by hash address."
    account(address:Address!):Account!

    """
    Get list of Contracts with at most <count> edges.
    If <count> is positive, return edges after the cursor,
    if negative, return edges before the cursor.
    For undefined cursor, positive <count> starts the list from top,
    negative <count> starts the list from bottom.
    ValidatedOnly specifies if the list should contain all the Contracts,
    or just contracts with validated byte code and available source/ABI.
    """
    contracts(validatedOnly: Boolean = false, cursor:Cursor, count:Int!):ContractList!

    """
    Get block information by number or by hash.
    If neither is provided, the most recent block is given.
    """
    block(number:Long, hash: Hash):Block

    "Get transaction information for given transaction hash."
    transaction(hash:Hash!):Transaction

    """
    Get list of Blocks with at most <count> edges.
    If <count> is positive, return edges after the cursor,
    if negative, return edges before the cursor.
    For undefined cursor, positive <count> starts the list from top,
    negative <count> starts the list from bottom.
    """
    blocks(cursor:Cursor, count:Int!):BlockList!

    """
    Get list of Transactions with at most <count> edges.
    If <count> is positive, return edges after the cursor,
    if negative, return edges before the cursor.
    For undefined cursor, positive <count> starts the list from top,
    negative <count> starts the list from bottom.
    """
    transactions(cursor:Cursor, count:Int!):TransactionList!

    "Get the id of the current epoch of the Opera blockchain."
    currentEpoch:Long!

    """
    Get information about specified epoch. Returns current epoch information
    if id is not provided.
    """
    epoch(id: Long!): Epoch!

    "The last staker id in Opera blockchain."
    lastStakerId: Long!

    "The number of stakers in Opera blockchain."
    stakersNum: Long!

    """
    Staker information. The staker is loaded either by numeric ID,
    or by address. null if none is provided.
    """
    staker(id: Long, address: Address): Staker

    "List of staker information from SFC smart contract."
    stakers: [Staker!]!

    "Returns the current price per gas in WEI units."
    gasPrice: Long!

    "Get price details of the Opera blockchain token for the given target symbols."
    price(to:String!):Price!
}

# Mutation endpoints for modifying the data
type Mutation {
    """
    SendTransaction submits a raw signed transaction into the block chain.
    The tx parameter represents raw signed and RLP encoded transaction data.
    """
    sendTransaction(tx: Bytes!):Transaction

    """
    Validate a deployed contract byte code with the provided source code
    so potential users can check the contract source code, access contract ABI
    to be able to interact with the contract and get the right metadata.
    Returns updated contract information. If the contract can not be validated,
    it raises a GraphQL error.
    """
    validateContract(contract: ContractValidationInput!): Contract!
}

# Subscriptions to live events broadcasting
type Subscription {
    "Subscribe to receive information about new blocks in the blockchain."
    onBlock: Block!

    "Subscribe to receive information about new transactions in the blockchain."
    onTransaction: Transaction!
}

`
