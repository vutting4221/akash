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

## Deployment Deposit

* Deposits are transfers from the source account to a module account.
* The unspent balance of deposits are transferred back to the source account when a deployment is closed.

Deposit:

|Field|Description|
|---|---|
|`DeploymentID`||
|`Balance`||
|`Transferred`|Denormalization for faster reporting|

Charge:

|Field|Description|
|---|---|
|`DeploymentID`||
|`LeaseID`||
|`Rate`|Tokens per block to transfer|
|`RebalancedAt`|Block that charges start|
|`Balance`|Balance currently reserved for owner|
|`Withdrawn`|Amount already withdrawn by owner|

#### Fast Path (naive)

If

```
sum([(c.CurrentHeight - c.RebalancedAt) * c.Rate + Balance + Withdrawn]) <= Deposit.Balance
```

Then for each Charge `c` of the Deposit `d`:

```
delta           = (CurrentHeight - c.RebalancedAt) * c.Rate
c.Balance      += delta
c.RebalancedAt  = CurrentHeight
d.Transferred  += delta
```

#### Accurate Path (naive)

1. Order charges by `RebalancedAt`, ascending.
1. Group charges by common open ranges, starting early.
1. For each group, if the total transfer will not overdraw the deposit, rebalance and continue.
1. If a group overdraws the deposit, distribute remaining balance to group entries weighted by their rate.

```
     |0|1|2|3|4|5|6|7|8|9|
c_a  | |*|*|*|*|*|*|*|*|*|
c_b  | | | |*|*|*|*|*|*|*|
c_c  | | | | | |*|*|*|*|*|
c_d  | | | | | | |*|*|*|*|
```

* when a new charge is added, always rebalance, removes need to group.
* when a withdraw happens, rebalance first.

===

Invariants:

* all existing charges have the same `RebalancedAt`.
* all balance has not been overdrawn.

```
Rebalance(d Deposit, height uint, charges []Charge)
  transfers []price
  for idx, c := range charges
    amount := (height - c.RebalancedAt) * c.Rate
    total += amount
    amount[idx] = amount

  if total <= (d.Balance - d.Transferred)
    d.Transferred += total
    for idx, c := range charges
      c.Balance += transfers[idx]
      c.RebalancedAt = height
    return

  remaining := (d.Balance - d.Transferred)

  // distribute, balance, weighted by transfer amount.

  for idx, c := range charges
    amount := (transfers[idx] / total) * remaining
    c.Balance += amount
    c.RebalancedAt = height

  d.Transferred += remaining
  d.State = OVERDRAWN
```

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

