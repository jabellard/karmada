# Karmada Push-Mode Cluster Secret Rotation Analysis

## Overview

This document analyzes how Karmada handles secret/token rotation for clusters registered in **push mode**. Specifically, it examines whether the Karmada controllers observe changes to the Secret referenced by a Cluster's `Spec.SecretRef` and whether they rebuild their kubeconfig/client connections with the new credentials.

---

## TL;DR — What Works and What Breaks

After a token rotation (assuming normal in-place update of the Secret's `.data.token`):

| Capability | Works after rotation? | Why | Evidence |
|---|---|---|---|
| **Health checks** (`/readyz`, `/healthz`) | **Yes** | Fresh client built every 10s via `BuildClusterConfig` → `secretGetter` | `cluster_status_controller.go:191,197` |
| **Pushing resources** (Create/Update/Delete to member cluster) | **Yes** | Fresh client built per operation via `ClusterClientSetFunc` | `objectwatcher.go:89,158,229` |
| **Cluster ready status** | **Yes** | Health check passes → `ClusterReady=True` maintained | `cluster_status_controller.go:197-199,209` |
| **Node/pod resource modeling** (NodeSummary, ResourceSummary) | **No** | TypedInformerManager holds stale client, cache freezes | `typedmanager/multi-cluster-manager.go:110-113` |
| **Work status collection** (status of propagated resources) | **No** | GenericInformerManager holds stale client, watches break | `work_status_controller.go:535`, `genericmanager/multi-cluster-manager.go:101-103` |
| **Scheduler capacity-based decisions** | **No** | Downstream of frozen ResourceSummary | `cluster_status_controller.go:278-279` |
| **Drift detection** (detecting external changes on member cluster) | **No** | GenericInformerManager watches broken, cache stale | `execution_controller.go:318`, `objectwatcher.go:106` |
| **Self-healing / automatic recovery** | **No** | No code detects auth failures on informer watches | `cluster_status_controller.go:134-135` (Stop only on cluster deletion) |

### The core problem

Karmada becomes **"write-only"** with respect to the member cluster. It can still push resources but is blind to their actual state.

- **Karmada can send changes to the member cluster** (push works — fresh credentials every time)
- **Karmada can no longer see what's happening on the member cluster** (observation broken — informers hold expired token)

### Concrete consequences

1. **You deploy something and it fails on the member cluster** — Karmada doesn't know. The `Work` object still shows the old successful status. Operators and dashboards see stale data.

2. **A node goes down on the member cluster** — Karmada doesn't know. `Cluster.Status.NodeSummary` and `ResourceSummary` still report the old capacity. The scheduler may keep sending workloads to a cluster that's actually overloaded.

3. **Someone manually modifies or deletes a resource on the member cluster** — Karmada doesn't know. It cannot detect drift or reconcile the difference. If a resource was deleted externally, Karmada's stale cache still thinks it exists, so it tries to Update instead of Create, getting repeated NotFound errors.

4. **The situation never recovers on its own** — the broken informers retry with the expired token forever. Recovery requires restarting the `karmada-controller-manager` or removing and re-joining the cluster.

### Root cause

The `TypedInformerManager` and `GenericInformerManager` are created once per cluster (`typedmanager/multi-cluster-manager.go:110-113`, `genericmanager/multi-cluster-manager.go:101-103`). Their `ForCluster` methods return the existing manager if one exists, never replacing the embedded client. There is no code that detects authentication failures on these watches and triggers recreation.

**Evidence:** `typedmanager/single-cluster-manager.go:91` — `informers.NewSharedInformerFactory(client, defaultResync)` captures the client at creation time. `cluster_status_controller.go:341-342` — `initializeGenericInformerManagerForCluster` short-circuits with `IsManagerExist`. `cluster_status_controller.go:134-135` — `Stop` is only called when the cluster is deleted, never on auth failure.

---

## How Credentials Are Used

### Secret Structure

Push-mode clusters store their credentials in a Kubernetes Secret referenced by `Cluster.Spec.SecretRef`. The secret contains:
- `token` key (`SecretTokenKey`) — the bearer token
- `caBundle` key (`SecretCADataKey`) — the CA certificate

### Building a Cluster Client

The function `BuildClusterConfig()` (`pkg/util/membercluster_client.go:192-252`) is responsible for constructing a `rest.Config` for a member cluster. It reads the secret fresh from the API server every time it is called:

```go
// pkg/util/membercluster_client.go:209
secret, err := secretGetter(cluster.Spec.SecretRef.Namespace, cluster.Spec.SecretRef.Name)
```

The `secretGetter` (`pkg/util/membercluster_client.go:260-266`) uses the controller-runtime `client.Client`:

```go
func secretGetter(client client.Client) func(string, string) (*corev1.Secret, error) {
    return func(namespace string, name string) (*corev1.Secret, error) {
        secret := &corev1.Secret{}
        err := client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, secret)
        return secret, err
    }
}
```

**Note:** The controller-runtime `client.Client` in karmada-controller-manager uses an informer-backed cache by default (see `cmd/controller-manager/app/controllermanager.go:174-207` — no cache bypass is configured for Secrets). This means `client.Get()` reads from the cache, not directly from the API server. However, the cache is kept current via watches on the Karmada control plane API server, so the latest secret data is available within milliseconds of an update. For practical purposes, this behaves the same as a direct API call — when the secret is updated in-place, the new token is available immediately on the next `BuildClusterConfig` call.

---

## No Secret Watcher Exists

Karmada controllers do **not** set up an informer or watch on the Secret objects that contain cluster credentials. There is no mechanism that triggers a reconciliation when a secret's data changes.

- `pkg/controllers/cluster/cluster_controller.go:196` — the cluster controller only watches `Cluster` objects
- `pkg/controllers/status/cluster_status_controller.go:175` — the cluster-status-controller watches only `Cluster` objects with `predicate.GenerationChangedPredicate{}`

---

## Per-Component Behavior After Token Rotation

### 1. Health Checks (cluster-status-controller) — Uses Fresh Credentials

Every reconcile cycle (default: every 10 seconds per `cmd/controller-manager/app/options/options.go:171`), the cluster-status-controller creates a brand-new `ClusterClient`:

- `cluster_status_controller.go:191` — `c.ClusterClientSetFunc(cluster.Name, c.Client, c.ClusterClientOption)`
- `controllermanager.go:362` — `ClusterClientSetFunc` is wired to `util.NewClusterClientSet`
- `membercluster_client.go:112` — `NewClusterClientSet` calls `BuildClusterConfig` which reads the secret from the controller-runtime cache (kept current via watches)
- `cluster_status_controller.go:197` — health check uses this fresh client: `getClusterHealthStatus(clusterClient)` hits `/readyz` or `/healthz`

**Result:** After a normal in-place token rotation, the health check picks up the new token within ~10 seconds and the cluster stays healthy.

#### Caveat: BuildClusterConfig Failure Bypasses Threshold

If `BuildClusterConfig` itself fails, the code at `cluster_status_controller.go:192-194` goes directly to `setStatusCollectionFailedCondition`, which at line 288-290 sets `ClusterConditionReady = False` with reason `StatusCollectionFailed`:

```go
readyCondition := util.NewCondition(clusterv1alpha1.ClusterConditionReady, statusCollectionFailed, message, metav1.ConditionFalse)
return updateStatusCondition(ctx, c.Client, cluster, readyCondition)
```

This path **bypasses** the `thresholdAdjustedReadyCondition` at line 199. The failure threshold (default 30s per `options.go:195`) only applies to health-check results, not to client construction failures.

`BuildClusterConfig` fails when:
- The Secret referenced by `Cluster.Spec.SecretRef` does not exist (the `secretGetter` returns an error)
- The Secret exists but is missing the `token` key or the `caBundle` key
- The `Cluster.Spec.APIEndpoint` is empty
- The `Cluster.Spec.SecretRef` is nil

**Important nuance:** Since the `secretGetter` reads from the controller-runtime cache, which is kept current via watches, the timing matters:
- **In-place update** (update the Secret's `.data.token` field): the cache receives the update event and serves the new data. `BuildClusterConfig` does NOT fail. This is the safe rotation method.
- **Delete-then-create** (delete the old Secret, create a new one): There is a window between when the delete event reaches the cache and when the create event does. During this window, `secretGetter` returns NotFound, `BuildClusterConfig` fails, and the cluster is immediately marked unhealthy with no grace period. The window is typically sub-second but depends on API server and controller-runtime watch latency.

### 2. Resource Push/Sync (execution-controller / ObjectWatcher) — Uses Fresh Credentials

The `objectWatcherImpl` creates a fresh dynamic client for every operation:

- `objectwatcher.go:89` — `Create` calls `o.ClusterClientSetFunc(clusterName, ...)` per invocation
- `objectwatcher.go:158` — `Update` calls `o.ClusterClientSetFunc(clusterName, ...)` per invocation
- `objectwatcher.go:229` — `Delete` calls `o.ClusterClientSetFunc(clusterName, ...)` per invocation

Each call goes through `NewClusterDynamicClientSet` → `BuildClusterConfig` → `secretGetter`, reading the current token from the controller-runtime cache. No client caching — a new `rest.Config` with the current bearer token is constructed each time.

**Result:** The actual API calls to the member cluster (Create/Update/Delete) always use fresh credentials.

**Caveat:** Both the execution controller (`execution_controller.go:318`) and the objectWatcher (`objectwatcher.go:106`) use `helper.GetObjectFromCache()` with the `GenericInformerManager` to check whether an object already exists in the member cluster before deciding to Create or Update. Since the GenericInformerManager's cache goes stale after token rotation (see Section 4), this can cause subtle issues:

- If a resource was **externally deleted** from the member cluster after the cache went stale, the stale cache still reports it as existing. The execution controller at `execution_controller.go:318-332` would call `Update` instead of `Create`. The `Update` at `objectwatcher.go:170` sets the ResourceVersion from the stale cached object. The API call at `objectwatcher.go:189` would then fail with a NotFound error (because the object no longer exists). The controller retries, but each retry hits the same stale cache and fails the same way.
- For **new resources** that were never in the cache, the cache correctly returns NotFound and the Create path proceeds normally with fresh credentials.

In practice, this edge case requires a resource to have existed in the member cluster before the token expired AND to have been externally deleted from that cluster afterward. Under normal Karmada operation (where Karmada is the sole manager of propagated resources), this scenario is uncommon.

### 3. TypedInformerManager (Node/Pod Watches) — Holds Stale Credentials

The `TypedInformerManager` is created once per cluster and never refreshed:

- `cluster_status_controller.go:356-358` — `buildInformerForCluster` first checks `GetSingleClusterManager`. Only if it returns `nil` does it create a new one
- `typedmanager/multi-cluster-manager.go:110-113` — `ForCluster` returns the existing manager if one exists:
  ```go
  // If informer manager already exist, just return
  if manager, exist := m.managers[cluster]; exist {
      return manager
  }
  ```
- `typedmanager/single-cluster-manager.go:88-101` — `NewSingleClusterInformerManager` captures the `client` at creation time and creates `informers.NewSharedInformerFactory(client, defaultResync)` at line 91. This factory uses that client for all list/watch operations forever.

Additionally, `IsInformerSynced` returns `true` forever once initially synced:

- `typedmanager/single-cluster-manager.go:159-164` — checks `syncedInformers` map
- `typedmanager/single-cluster-manager.go:261-262` — entries are added to `syncedInformers` but **never removed**
- `cluster_status_controller.go:364-374` — because `IsInformerSynced` returns `true`, the function returns the stale manager immediately without re-syncing

**Result:** After token rotation, the informer's watches break. The `listNodes()` and `listPods()` calls (`cluster_status_controller.go:269-274`) return stale cached data. Errors from listing are logged but not returned (lines 270-271, 274-275), so the stale data is silently written into `Cluster.Status.NodeSummary` and `Cluster.Status.ResourceSummary`.

### 4. GenericInformerManager (Resource Watches) — Holds Stale Credentials

Same pattern as TypedInformerManager:

- `cluster_status_controller.go:340-343` — `initializeGenericInformerManagerForCluster` short-circuits if manager exists:
  ```go
  if c.GenericInformerManager.IsManagerExist(clusterClient.ClusterName) {
      return
  }
  ```
- `genericmanager/multi-cluster-manager.go:101-103` — `ForCluster` returns existing manager
- `work_status_controller.go:532-544` — work-status-controller follows the same pattern:
  ```go
  singleClusterInformerManager := c.InformerManager.GetSingleClusterManager(cluster.Name)
  if singleClusterInformerManager == nil {
      // only creates a new one if none exists
  }
  return singleClusterInformerManager, nil
  ```

**Result:** Watches on member cluster resources break. Work status in Karmada stops reflecting changes to propagated resources.

### 5. No Self-Healing Mechanism

The only calls to `GenericInformerManager.Stop()` / `TypedInformerManager.Stop()` in the cluster-status-controller are at lines 134-135, inside the `apierrors.IsNotFound(err)` branch — meaning the cluster object no longer exists:

```go
if apierrors.IsNotFound(err) {
    c.GenericInformerManager.Stop(req.NamespacedName.Name)
    c.TypedInformerManager.Stop(req.NamespacedName.Name)
```

There is **no code** in `cluster_status_controller.go` that calls `Stop` when a secret changes, when authentication fails, or when informer watches break.

### 6. ClusterAccessCredentialChanged Does Not Detect Secret Changes

`pkg/util/cluster.go:242-250`:

```go
func ClusterAccessCredentialChanged(newSpec, oldSpec clusterv1alpha1.ClusterSpec) bool {
    if oldSpec.APIEndpoint == newSpec.APIEndpoint &&
        oldSpec.InsecureSkipTLSVerification == newSpec.InsecureSkipTLSVerification &&
        oldSpec.ProxyURL == newSpec.ProxyURL &&
        equality.Semantic.DeepEqual(oldSpec.ProxyHeader, newSpec.ProxyHeader) {
        return false
    }
    return true
}
```

It checks exactly four fields: `APIEndpoint`, `InsecureSkipTLSVerification`, `ProxyURL`, `ProxyHeader`. It does **not** check `SecretRef` (the name/namespace of the secret), nor can it check the secret's `.data` (since it only receives `ClusterSpec`, not the `Secret` object).

The metrics adapter (`metricsadapter/controller.go:177-181`) calls `stopInformerManager` when `ClusterAccessCredentialChanged` returns `true`, but this is irrelevant because the function doesn't detect secret data changes.

---

## Summary Table

| Component | Picks Up Rotated Token? | Mechanism | Evidence |
|---|---|---|---|
| Health checks | **Yes** | Fresh `BuildClusterConfig` every 10s | `cluster_status_controller.go:191→197` |
| Resource push/sync (API calls) | **Yes** | Fresh `ClusterClientSetFunc` per operation | `objectwatcher.go:89, 158, 229` |
| Resource push/sync (cache checks) | **No** | `GetObjectFromCache` uses stale GenericInformerManager | `execution_controller.go:318`, `objectwatcher.go:106` |
| TypedInformerManager (node/pod watches) | **No** | Created once, never refreshed | `typedmanager/multi-cluster-manager.go:110-113` |
| GenericInformerManager (resource watches) | **No** | Created once, never refreshed | `genericmanager/multi-cluster-manager.go:101-103` |
| Work status collection | **No** | Uses cached GenericInformerManager | `work_status_controller.go:535` |
| Metrics adapter (cached clients) | **No** | Only invalidated on `Cluster.Spec` changes | `metricsadapter/controller.go:177` |

---

## Side Effects

### Can Karmada Mark a Cluster Unhealthy by Mistake?

| Rotation Scenario | Cluster Marked Unhealthy? | Details |
|---|---|---|
| Secret `.data.token` updated in-place (atomic) | **No** | Controller-runtime cache receives the update event; next `BuildClusterConfig` call uses new token; health check passes |
| Secret deleted and recreated (even briefly) | **Yes** | During the window between delete and create events reaching the cache, `secretGetter` returns NotFound, `BuildClusterConfig` fails, and `setStatusCollectionFailedCondition` sets `ClusterConditionReady=False` bypassing the failure threshold. Window is typically sub-second. |
| Token expires before new one is written | **Yes** | Health check uses fresh client with expired token, gets 401, `getClusterHealthStatus` reports `online=false`. The `thresholdAdjustedReadyCondition` at line 199 provides a 30s grace period (`options.go:195`) before transitioning to unhealthy. |

### Silent Side Effects (Regardless of Rotation Method)

Once the old token expires, the `TypedInformerManager` and `GenericInformerManager` informers lose their watch connections to the member cluster. The Kubernetes `Reflector` inside each informer retries list/watch with exponential backoff, but it reuses the **same client** that was created with the old bearer token (`typedmanager/single-cluster-manager.go:91` — `informers.NewSharedInformerFactory(client, defaultResync)`). Every retry fails with 401 Unauthorized. This continues **indefinitely** — there is no code that detects the auth failure and recreates the informer with a fresh client.

These failures are **permanent until the karmada-controller-manager is restarted** (or the cluster is removed and re-joined). They are not transient or self-healing.

#### 1. Stale Resource Modeling — Scheduler Makes Decisions on Frozen Capacity Data

**What breaks:** `Cluster.Status.NodeSummary` and `Cluster.Status.ResourceSummary` freeze at whatever values were last cached before the token expired.

**Why it's serious:** The karmada scheduler reads `Cluster.Status.ResourceSummary` to make replica assignment decisions in dynamic scheduling (e.g., spreading replicas based on available cluster capacity). With frozen data, the scheduler may:
- Over-schedule a cluster whose real capacity has decreased (nodes removed, resources consumed)
- Under-schedule a cluster whose capacity has increased (nodes added, pods freed)

**Why it's silent:** Errors from `listNodes` and `listPods` are logged but not returned (`cluster_status_controller.go:270-275`). The function continues with stale cached data and writes it into the Cluster status alongside a `ClusterReady=True` condition. Nothing in the Cluster status indicates the data is stale.

**Evidence chain:**
- `cluster_status_controller.go:262` — calls `buildInformerForCluster(clusterClient)` which returns the stale TypedInformerManager
- `cluster_status_controller.go:269-274` — `listNodes`/`listPods` read from the broken informer's cache, errors are logged but not returned
- `cluster_status_controller.go:278-279` — stale data written to `currentClusterStatus.NodeSummary` and `ResourceSummary`
- `cluster_status_controller.go:227-228` — persisted with `ClusterReady=True`

#### 2. Work Status Stops Updating — Loss of Visibility into Propagated Resources

**What breaks:** The `work-status-controller` can no longer observe the actual state of resources it has pushed to the member cluster. The `Work` objects in the Karmada control plane stop reflecting reality.

**Why it's serious:** Operators and higher-level controllers rely on `Work.Status` to know whether propagated resources (Deployments, Services, ConfigMaps, etc.) are healthy. With frozen work status:
- A Deployment rollout that fails on the member cluster will still appear successful in Karmada
- A deleted or evicted Pod will still appear running in Karmada
- Rollback decisions, health-based rescheduling, and operational dashboards all operate on stale data

**Why it's silent:** The `GenericInformerManager`'s watches silently fail and the cache goes stale. No condition or event is surfaced on the Work or Cluster object to indicate the status is outdated.

**Evidence chain:**
- `work_status_controller.go:535` — `GetSingleClusterManager(cluster.Name)` returns the existing manager with the old token
- `genericmanager/multi-cluster-manager.go:101-103` — `ForCluster` returns cached manager without refreshing the client
- The informer's Reflector retries with the same stale client indefinitely

#### 3. No Recovery Without Manual Intervention

**What breaks:** There is no self-healing mechanism. The informer managers are never stopped and recreated in response to authentication failures.

**Why it's serious:** The failures in (1) and (2) persist until an operator either:
- Restarts the `karmada-controller-manager` pod (which destroys all informer managers and forces them to be recreated with fresh clients on startup)
- Removes the Cluster object and re-joins the cluster (the only code path that calls `Stop` is cluster deletion at `cluster_status_controller.go:134-135`)

**Why it's silent:** The cluster remains marked as `ClusterReady=True` (health checks use a fresh client and succeed). There is no alert, condition, or event indicating that the informer watches are broken. The only indicator is error-level log lines from the Reflector's failed list/watch retries, which are easy to miss in a busy cluster.

#### Severity Assessment

| Failure | Severity | Impact Scope | Duration | Visibility |
|---|---|---|---|---|
| Frozen resource modeling | **High** | Scheduling decisions for all workloads targeting this cluster | Permanent until restart | None — no condition or event surfaced |
| Frozen work status | **High** | All propagated resources on this cluster | Permanent until restart | None — Work objects show last-known-good status |
| Misleading cluster health | **Medium** | Operators trust the cluster is fully functional | Permanent until restart | Actively misleading — `ClusterReady=True` with stale data |
| Stale cache for workload sync decisions | **Low** | Only affects resources externally deleted from member cluster after cache went stale | Permanent until restart | Error logs from failed Update API calls |

---

## Key File References

| File | Relevant Lines | Purpose |
|---|---|---|
| `pkg/util/membercluster_client.go` | 191-252, 260-266 | `BuildClusterConfig`, `secretGetter` |
| `pkg/util/cluster.go` | 241-250 | `ClusterAccessCredentialChanged` |
| `pkg/controllers/status/cluster_status_controller.go` | 130-147, 181-232, 288-291, 340-397 | Cluster status reconciliation, informer initialization |
| `pkg/controllers/status/cluster_condition_cache.go` | 44-84 | Threshold-adjusted ready condition |
| `pkg/controllers/status/work_status_controller.go` | 530-545 | Work status informer manager caching |
| `pkg/util/objectwatcher/objectwatcher.go` | 88-125, 150-205, 207-245 | ObjectWatcher Create/Update/Delete |
| `pkg/controllers/execution/execution_controller.go` | 311-338 | Execution controller tryCreateOrUpdateWorkload |
| `pkg/util/fedinformer/typedmanager/multi-cluster-manager.go` | 106-118 | TypedInformerManager `ForCluster` caching |
| `pkg/util/fedinformer/typedmanager/single-cluster-manager.go` | 86-101, 159-164, 253-266 | SingleClusterInformerManager, `IsInformerSynced`, `syncedInformers` |
| `pkg/util/fedinformer/genericmanager/multi-cluster-manager.go` | 97-108 | GenericInformerManager `ForCluster` caching |
| `pkg/metricsadapter/controller.go` | 164-183, 277-281 | Metrics adapter credential change handling |
| `cmd/controller-manager/app/options/options.go` | 171, 194-195 | Default `ClusterStatusUpdateFrequency` (10s), thresholds (30s) |
| `cmd/controller-manager/app/controllermanager.go` | 360-369 | Controller wiring |
