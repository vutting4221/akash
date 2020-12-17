# Marketplace State Machine

## Lease Payments

Leases are paid from deployment owner (tenant) to the provider
through a deposit & withdraw mechanism.

Tenants are required to submit a deposit when creating a deployment.  Leases
will be paid passively from the balance of this deposit.  At any time,
a lease provider may withdraw the balance owed to them from this deposit.

If the available funds in the deposit ever reaches zero, a provider may
close the lease.

A tenant can add funds to their deposit at any time.

When a deployment is closed, the unspent portion of the balance will be returned
to the tenant.

## On-Chain Parameters

|Name|Initial Value|Description|
|---|---|---|
|`deployment_min_deposit`|`10akt`|Minimum deposit to make deployments.  Target: ~$10|
|`bid_min_deposit`|`100akt`|Deposit amount required to bid.  Target: ~$100|
|`bid_min_ttl`|20|Minimum lifetime of bids.|


## Transactions
### DeploymentCreate

#### Parameters

|Name|Description|
|---|---|
|`DeploymentID`| ID of Deployment. |
|`Deposit`| Deposit amount.  Must be greater than `deployment_min_deposit`|

### DeploymentFund

#### Parameters

|Name|Description|
|---|---|
|`DeploymentID`| ID of Deployment. |
|`Deposit`| Deposit amount.  Must be greater than `deployment_min_deposit`|

#### Actions

### BidCreate

#### Parameters

|name|description|
|---|---|
|`OrderID`| ID of Order |
|`TTL`| Number of blocks this bid is valid for |
|`Deposit`| Deposit amount.  `bid_min_deposit` if empty.|

#### Actions

### LeaseCreate

Creates a lease for the bid identified by `BidID`.

#### Parameters

|name|description|
|---|---|
|`BidID`|Bid to create a lease from|

#### Actions

1. Creates a `Lease` from the given `Bid`.
1. Returns deposits from all bids that lost.

### MarketWithdraw

This withdraws balances earned by providing for leases and deposits
of bids that have expired.

#### Parameters

|name|description|
|---|---|
|`Owner`|Provider ID to withdraw funds for.|

