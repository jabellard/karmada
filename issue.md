**Title**: Push-mode cluster informers break permanently after credential Secret rotation

**What would you like to be added**:

Automatic recovery of member cluster informers when a push-mode cluster's credential Secret is rotated.

**Why is this needed**:

When a push-mode cluster's credential Secret is updated (token rotation), Karmada's informer-based functionality breaks permanently for that cluster. while health checks and resource pushing continue working. This creates a silent, partial failure that does not self-heal.

## Root Cause

The `ClusterStatusController` builds a fresh client from the Secret on every reconcile cycle. This fresh client is used for health checks, so health checks always pass after rotation. However, the `GenericInformerManager` and `TypedInformerManager` are lazily initialized once and then reused forever. Their internal clients have the old bearer token baked in as a static string. After rotation, the stale manager already exists, so the lazy init path skips recreation. The old watches retry with the expired token, get 401, and retry indefinitely.

There is no watch on Secrets. only on Cluster objects. A Secret update never triggers reconciliation, and even the periodic requeue hits the same "manager already exists" guard.

## What still works after rotation

- Cluster health checks (fresh client built each cycle from the current Secret)
- Resource pushing (also builds fresh clients)
- Cluster Ready condition remains `True`

## What breaks silently

- Status collection for propagated resources (work status, drift detection)
- Resource modeling and scheduler accuracy (stale node/pod data)
- Multi-cluster service endpoint collection
- FederatedHPA and metrics

The cluster appears healthy while silently losing all informer-driven observability.
