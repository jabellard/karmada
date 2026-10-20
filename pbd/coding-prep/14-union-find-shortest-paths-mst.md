# 14. Union-Find, Shortest Paths, and MST

When edges carry weights or connectivity changes over time, BFS/DFS isn't enough. Union-Find, weighted shortest paths, and minimum spanning trees handle routing, latency budgets, cluster merging, and cheapest wiring.

## At a glance

- **Use when:** transitive grouping or edges arriving over time (DSU); per-edge costs (shortest paths); "connect everything at minimum total cost" (MST).
- **Templates:**
  - Template 1: DSU: components, redundant edge, grouping by shared keys.
  - Template 2: Dijkstra: non-negative weights, one source (also bottleneck with `max`).
  - Template 3: Bellman-Ford with an edge limit: "at most k stops" or negative weights.
  - Template 4: 0-1 BFS: edge weights only 0 or 1.
  - Template 5: Floyd-Warshall: all pairs, n ≤ ~400.
  - Template 6: Kruskal: MST from an edge list (sparse).
  - Template 7: Prim: MST grown from one node; array form for dense/complete graphs.
- **Typical complexity:** DSU ~O(α(n)) per op; Dijkstra O((V + E) log V); the rest in the decision table.
- **Top trap:** Dijkstra with negative edges or a hop limit; use Bellman-Ford.

## Core idea

**Union-Find (DSU).** Each set is a tree whose root is its ID; `find` walks to the root, `union` links two roots. Two optimizations make operations effectively O(1) amortized (inverse Ackermann):
- *Path compression*: `find` points nodes closer to the root.
- *Union by size/rank*: the smaller tree goes under the larger, keeping trees shallow.

Invariant: two elements are connected exactly when they share a root. A failed `union` (returns False) means the edge closes a cycle.

**Dijkstra.** With non-negative weights, the closest unvisited node is final: any other path to it passes through a node at least as far, then adds non-negative edges. A min-heap yields it. `heapq` lacks decrease-key, so push duplicates and skip stale entries on pop.

**Bellman-Ford.** Relax every edge, k times. With snapshots, after round i `dist[v]` is the best cost using at most i edges. Handles negative weights; n-1 rounds suffice without negative cycles.

**0-1 BFS.** With weights in {0, 1}, a deque keeps nodes sorted by distance: 0-edges to the front, 1-edges to the back.

**Floyd-Warshall.** `d[i][j] = min(d[i][j], d[i][k] + d[k][j])` with k outermost. After processing k, `d[i][j]` is the shortest path using only intermediates `0..k`.

**MST.** Cut property: the cheapest edge crossing any cut is in some MST. Kruskal takes the cheapest edges that don't form a cycle (DSU check); Prim grows one tree by its cheapest outgoing edge (heap).

### Decision table

| Situation | Algorithm | Time |
|---|---|---|
| Unweighted (all edges cost 1) | BFS (`13`) | O(V + E) |
| Weights only 0 or 1 | 0-1 BFS with deque | O(V + E) |
| Non-negative weights, single source | Dijkstra with heap | O((V + E) log V) |
| Path cost = max/min edge (bottleneck), non-negative | Dijkstra with `max` combiner, binary search + BFS, or DSU in sorted edge order | O(E log V) |
| Negative weights, or "at most k edges/stops" | Bellman-Ford (k rounds with snapshot) | O(k * E) |
| Detect negative cycle | Bellman-Ford, one extra round | O(V * E) |
| All pairs, n ≤ ~400 | Floyd-Warshall | O(n^3) |
| DAG with weights (any sign) | Topological order + relax (`13`) | O(V + E) |
| Connectivity, edges arrive over time, union only | DSU | ~O(α(n)) per op |
| Connect everything at minimum total cost | MST: Kruskal (sparse) or Prim (dense) | O(E log E) / O(V^2) |

## Recognition signals

- "Groups", "merge if they share", "equivalent", "transitively connected", "after each addition, how many components?" → DSU.
- "Which edge to remove to make it a tree", "redundant link" → DSU, first edge whose `union` fails.
- "Minimum time/cost/latency from A to B" with per-edge costs → Dijkstra (any negative weights?).
- "Time for a signal to reach all nodes" → Dijkstra, then the max distance.
- "Minimize the largest step / steepest climb / max edge on the path" → bottleneck path (Dijkstra with `max`, or binary search on the answer).
- "At most k stops / hops / transfers" → Bellman-Ford with k+1 rounds or layered BFS, not Dijkstra on `node` alone.
- "Changing direction costs 1, following is free", "removing an obstacle costs 1" → 0-1 BFS.
- "Connect all servers / houses / points with minimum total cable" → MST.
- n ≤ 100 to 400 with queries between arbitrary pairs → Floyd-Warshall.

## Disguises

1. *A fleet of hosts reports pairs that share a storage volume. Volumes are shared transitively. How many independent storage domains exist, and which hosts are in each?*
   Decode: transitive grouping by shared resource = components. DSU over hosts (or volumes), then group by root (Accounts Merge shape).

2. *A service mesh has per-hop latencies. A config push starts at one gateway. How long until every sidecar has it, or report that some never get it?*
   Decode: one source, non-negative weights, "until all receive" = max shortest distance → Dijkstra (Network Delay Time).

3. *A hiking app wants a route between two points that minimizes the steepest single climb, not the total climb.*
   Decode: path cost is the max edge, not the sum. Dijkstra with a `max` combiner works because `max` is monotone (Path With Minimum Effort).

4. *Packets may traverse at most k intermediate routers due to TTL. Find the cheapest route under that limit.*
   Decode: an edge-count limit breaks Dijkstra's greedy choice → Bellman-Ford with k+1 snapshot rounds.

5. *You must lay fiber between n data centers. Any two can be linked at a cost equal to their Manhattan distance. Minimum total fiber so all are connected?*
   Decode: connect everything at minimum total weight, no source → MST. Complete graph, so O(n^2) Prim.

6. *Variables are constrained by statements like "x equals y" and "y differs from z". Can all statements hold?*
   Decode: equality is an equivalence relation → DSU over `==`, then check each `!=` pair is in different sets.

## Templates and worked examples

Imports (companion file: `code/14_union_find_shortest_paths_mst.py`):

```python
import heapq
from collections import defaultdict, deque
```

### Template 1: DSU (path compression + union by size)

> - **Example problem:** [Redundant Connection (684)](https://leetcode.com/problems/redundant-connection/)
> - **Why this template:** the extra edge that made a tree cyclic is the first edge whose endpoints already share a DSU set.
> - **Worked below:** Example 1.1
> - **Also fits:** [Number of Provinces (547)](https://leetcode.com/problems/number-of-provinces/)

```python
class DSU:
    """Disjoint Set Union with path compression (halving) and union by size.

    Amortized near-O(1) per operation (inverse Ackermann).
    """

    def __init__(self, n: int) -> None:
        self.parent = list(range(n))
        self.size = [1] * n
        self.count = n                        # number of disjoint sets

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]   # path halving
            x = self.parent[x]
        return x

    def union(self, a: int, b: int) -> bool:
        """Merge the sets of a and b. Returns False if already in the same set."""
        ra, rb = self.find(a), self.find(b)
        if ra == rb:
            return False
        if self.size[ra] < self.size[rb]:     # attach smaller tree under larger
            ra, rb = rb, ra
        self.parent[rb] = ra
        self.size[ra] += self.size[rb]
        self.count -= 1
        return True

    def connected(self, a: int, b: int) -> bool:
        return self.find(a) == self.find(b)
```

Knobs: IDs other than `0..n-1` (strings, coordinates) → map to indices, or dicts with `parent.setdefault(x, x)`; per-root data (size, min, sum) merged in `union`; `count` is the component count. Recursive `parent[x] = find(parent[x])` compresses fully but can hit the recursion limit on long chains; the iterative `find` can't.

**DSU vs BFS/DFS:**

| Use DSU when | Use BFS/DFS when |
|---|---|
| Edges arrive incrementally with queries in between | Graph is static and given up front |
| Only connectivity matters | You need distances, paths, or traversal order |
| Grouping by shared attributes (emails, values) | You need directed reachability (DSU is undirected only) |
| Kruskal's MST, cycle detection in an edge stream | Edges get deleted (DSU can't split; offline, process deletions in reverse as additions) |

#### Example 1.1: Components and the cycle-closing edge [Medium]

**Problems:** [Number of Provinces (547)](https://leetcode.com/problems/number-of-provinces/), [Redundant Connection (684)](https://leetcode.com/problems/redundant-connection/) · **Template:** Template 1 (DSU)

Problem (a): an n x n symmetric matrix says whether cities i and j are directly connected; count transitively connected groups. (b): a tree on nodes 1..n got one extra edge. Return an edge whose removal leaves a tree (the last in the input if several qualify).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** (a) components; (b) the extra cycle-creating edge = the first edge, in order, whose endpoints are already connected.

**Approach:** (a) union connected pairs; answer is `dsu.count`. (b) The first failing `union`, in order, is the answer, and the cycle's last input edge, since earlier cycle edges succeeded.

```python
def find_circle_num(is_connected: list[list[int]]) -> int:
    """Number of Provinces: components from an adjacency matrix."""
    n = len(is_connected)
    dsu = DSU(n)
    for i in range(n):
        for j in range(i + 1, n):
            if is_connected[i][j]:
                dsu.union(i, j)
    return dsu.count


def find_redundant_connection(edges: list[list[int]]) -> list[int]:
    """Redundant Connection: first edge (in input order) that closes a cycle. Nodes 1..n."""
    dsu = DSU(len(edges) + 1)
    for u, v in edges:
        if not dsu.union(u, v):
            return [u, v]
    return []
```

**Complexity:** (a) O(n^2 α(n)) for n^2 matrix entries. (b) O(n α(n)).

**Edge cases:** 1-indexed nodes (allocate n+1). The directed version (Redundant Connection II) is much harder: a node with two parents and/or a cycle needs case analysis.

</details>

#### Example 1.2: Grouping by shared keys [Medium]

**Problem:** [Accounts Merge (721)](https://leetcode.com/problems/accounts-merge/) · **Template:** Template 1 (DSU)

Problem: each account is a name plus emails. Accounts sharing any email (transitively) are one person. Merge them; output name plus sorted emails.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** transitive "share any ... belong together" → DSU. The trick is what to union: account indices, linked through an `email -> first account` map.

**Approach:** union each account with the first account that had each of its emails, then group emails by root. This avoids unioning every pair sharing an email (potentially quadratic).

```python
def accounts_merge(accounts: list[list[str]]) -> list[list[str]]:
    """Accounts Merge: union accounts that share any email."""
    dsu = DSU(len(accounts))
    owner: dict[str, int] = {}                # email -> first account index that had it
    for i, acc in enumerate(accounts):
        for email in acc[1:]:
            if email in owner:
                dsu.union(i, owner[email])
            else:
                owner[email] = i
    groups: dict[int, list[str]] = defaultdict(list)
    for email, i in owner.items():
        groups[dsu.find(i)].append(email)
    return [[accounts[root][0]] + sorted(emails) for root, emails in groups.items()]
```

**Complexity:** O(E log E) to sort emails per group plus O(E α(n)) for unions, E = total emails.

**Edge cases:** same name does **not** mean same person; an account repeating an email; accounts 0 and 2 merging transitively with nothing shared. BFS over an email graph also works; DSU is shorter.

</details>

### Template 2: Dijkstra with stale-entry skipping

> - **Example problem:** [Network Delay Time (743)](https://leetcode.com/problems/network-delay-time/)
> - **Why this template:** "time for a signal from k to reach all nodes" with non-negative times = single-source shortest paths, then the max.
> - **Worked below:** Example 2.1
> - **Also fits:** [Path With Minimum Effort (1631)](https://leetcode.com/problems/path-with-minimum-effort/)

```python
def dijkstra(n: int, adj: list[list[tuple[int, int]]], src: int) -> list[float]:
    """Single-source shortest paths, non-negative weights. adj[u] = [(v, w), ...].

    Invariant: when a node is popped with d == dist[node], that distance is final.
    Stale entries (d > dist[node]) are skipped instead of decrease-key.
    """
    dist = [float("inf")] * n
    dist[src] = 0
    heap = [(0, src)]
    while heap:
        d, u = heapq.heappop(heap)
        if d > dist[u]:
            continue                          # stale entry
        # TODO: early exit here if u is the target
        for v, w in adj[u]:
            nd = d + w                        # TODO: path cost combiner (sum, max, product...)
            if nd < dist[v]:
                dist[v] = nd
                heapq.heappush(heap, (nd, v))
    return dist
```

Knobs: the combiner (`d + w` for sums, `max(d, w)` for bottleneck, `-log(p)` or a max-heap of probabilities for most reliable); node identity (`(r, c)` or `(node, extra_state)`); early exit on popping the target. The heap holds up to E entries, so O(E log E) = O(E log V); the stale skip relaxes each node's edges once.

#### Example 2.1: Single-source weighted shortest paths [Medium]

**Problem:** [Network Delay Time (743)](https://leetcode.com/problems/network-delay-time/) · **Template:** Template 2 (Dijkstra with stale-entry skipping)

Problem: nodes 1..n, directed edges `(u, v, time)`, signal sent from node k. Return the time until all nodes receive it, or -1 if some never do.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** non-negative weights + one source + "all receive" → Dijkstra, then the max distance.

**Approach:** build the adjacency list, run Template 2, check for unreachable nodes.

```python
def network_delay_time(times: list[list[int]], n: int, k: int) -> int:
    """Network Delay Time: nodes 1..n, directed weighted edges; time for all to receive."""
    adj: list[list[tuple[int, int]]] = [[] for _ in range(n + 1)]
    for u, v, w in times:
        adj[u].append((v, w))
    dist = dijkstra(n + 1, adj, k)
    worst = max(dist[1:])
    return -1 if worst == float("inf") else int(worst)
```

**Complexity:** O((V + E) log V).

**Edge cases:** index 0 unused (slice it out of the max); unreachable nodes. Not BFS: a 2-hop path can beat a 1-hop path. Not DFS: it explores exponentially many paths without memoizing final distances.

</details>

#### Example 2.2: Bottleneck path cost [Medium]

**Problem:** [Path With Minimum Effort (1631)](https://leetcode.com/problems/path-with-minimum-effort/) · **Template:** Template 2 (Dijkstra with stale-entry skipping), with `max` as the path combiner

Problem: a grid of heights. A route's effort is its maximum absolute height difference between consecutive cells. Find the minimum effort from top-left to bottom-right.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "minimize the max step along a path" → bottleneck shortest path. Dijkstra only needs path cost never to decrease on extension; `max` satisfies this.

**Approach:** Dijkstra over cells with combiner `max(effort_so_far, step)`. The target's first pop is optimal.

```python
def minimum_effort_path(heights: list[list[int]]) -> int:
    """Path With Minimum Effort: path cost = MAX edge (abs height diff), minimize it."""
    rows, cols = len(heights), len(heights[0])
    best = [[float("inf")] * cols for _ in range(rows)]
    best[0][0] = 0
    heap = [(0, 0, 0)]
    while heap:
        e, r, c = heapq.heappop(heap)
        if e > best[r][c]:
            continue
        if (r, c) == (rows - 1, cols - 1):
            return e                          # first pop of target is optimal
        for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            nr, nc = r + dr, c + dc
            if 0 <= nr < rows and 0 <= nc < cols:
                ne = max(e, abs(heights[nr][nc] - heights[r][c]))   # bottleneck combiner
                if ne < best[nr][nc]:
                    best[nr][nc] = ne
                    heapq.heappush(heap, (ne, nr, nc))
    return 0
```

**Complexity:** O(RC log(RC)).

Alternatives: binary search the threshold with BFS over steps ≤ threshold (`07-binary-search.md`), or union edges in weight order until start and end connect (Kruskal-style). Both also solve Swim in Rising Water.

**Edge cases:** 1x1 grid → 0; grid DFS without memo times out.

</details>

### Template 3: Bellman-Ford with an edge limit

> - **Example problem:** [Cheapest Flights Within K Stops (787)](https://leetcode.com/problems/cheapest-flights-within-k-stops/)
> - **Why this template:** "at most k stops" caps paths at k+1 edges: exactly k+1 snapshot rounds of Bellman-Ford.
> - **Worked below:** Example 3.1

```python
def bellman_ford_k_edges(n: int, edges: list[tuple[int, int, int]], src: int, k: int) -> list[float]:
    """Shortest paths using AT MOST k edges. Handles negative weights.

    Each round relaxes from a snapshot of the previous round so a path cannot
    use more than one new edge per round. With k = n - 1 this is classic Bellman-Ford.
    """
    dist = [float("inf")] * n
    dist[src] = 0
    for _ in range(k):
        prev = dist[:]                        # snapshot: the crucial line for edge limits
        for u, v, w in edges:
            if prev[u] + w < dist[v]:
                dist[v] = prev[u] + w
    return dist
```

Knobs: k (`stops + 1` edges). Without the snapshot, a round can chain several edges and break the limit (harmless when unlimited). Negative cycle check: any improvement in a round after n-1 means a cycle.

```python
def has_negative_cycle(n: int, edges: list[tuple[int, int, int]]) -> bool:
    """After n-1 rounds, any further improvement means a reachable negative cycle."""
    dist = [0] * n                            # virtual source connected to all nodes
    for _ in range(n - 1):
        for u, v, w in edges:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
    return any(dist[u] + w < dist[v] for u, v, w in edges)
```

#### Example 3.1: Shortest path with a hop limit [Medium]

**Problem:** [Cheapest Flights Within K Stops (787)](https://leetcode.com/problems/cheapest-flights-within-k-stops/) · **Template:** Template 3 (Bellman-Ford with an edge limit)

Problem: n cities, directed priced flights. Find the cheapest price from src to dst with at most k stops, or -1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** at most k stops = at most k+1 edges. Dijkstra keyed on node alone may discard a pricier path with fewer stops, the only one within the limit.

**Approach:** k+1 Bellman-Ford rounds, each relaxing from a snapshot. After round i, `dist[v]` is the cheapest cost using at most i edges.

```python
def find_cheapest_price(n: int, flights: list[list[int]], src: int, dst: int, k: int) -> int:
    """Cheapest Flights Within K Stops: at most k stops = at most k + 1 edges."""
    dist = bellman_ford_k_edges(n, [tuple(f) for f in flights], src, k + 1)
    return -1 if dist[dst] == float("inf") else int(dist[dst])
```

Equivalent layered BFS (often easier to explain): one layer per edge, carrying forward only improvements.

```python
def find_cheapest_price_bfs(n: int, flights: list[list[int]], src: int, dst: int, k: int) -> int:
    """Same problem as layered BFS: one layer per edge, keep the best cost per node."""
    adj: list[list[tuple[int, int]]] = [[] for _ in range(n)]
    for u, v, w in flights:
        adj[u].append((v, w))
    best = [float("inf")] * n
    best[src] = 0
    frontier = [(src, 0)]
    for _ in range(k + 1):
        nxt: dict[int, int] = {}
        for u, cost in frontier:
            for v, w in adj[u]:
                nc = cost + w
                if nc < best[v] and nc < nxt.get(v, float("inf")):
                    nxt[v] = nc
        for v, c in nxt.items():
            best[v] = min(best[v], c)
        frontier = list(nxt.items())
    return -1 if best[dst] == float("inf") else int(best[dst])
```

Pruning with `best[v]` is safe: reaching v in a later layer at no better cost is dominated (more edges, not cheaper).

**Complexity:** O(k * E) time, O(n) space.

**Edge cases:** no snapshot (chain `0->1->2->3` with k=1 would use 3 edges in one round); `k + 1` vs `k`; dst unreachable. A third option is Dijkstra on state `(node, stops_used)`.

</details>

### Template 4: 0-1 BFS

> - **Example problem:** [Minimum Cost to Make at Least One Valid Path in a Grid (1368)](https://leetcode.com/problems/minimum-cost-to-make-at-least-one-valid-path-in-a-grid/)
> - **Why this template:** following a cell's arrow costs 0 and changing it costs 1, so every edge weighs 0 or 1.
> - **Worked below:** Example 4.1

```python
def zero_one_bfs(n: int, adj: list[list[tuple[int, int]]], src: int) -> list[float]:
    """Shortest paths when every weight is 0 or 1: deque instead of heap, O(V + E)."""
    dist = [float("inf")] * n
    dist[src] = 0
    dq = deque([src])
    while dq:
        u = dq.popleft()
        for v, w in adj[u]:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
                if w == 0:
                    dq.appendleft(v)          # same distance layer: front
                else:
                    dq.append(v)              # next layer: back
    return dist
```

The deque holds at most two distinct distances (d and d+1), in order, so it acts like Dijkstra's heap without the log factor.

#### Example 4.1: Cheapest arrow flips to reach the corner [Hard]

**Problem:** [Minimum Cost to Make at Least One Valid Path in a Grid (1368)](https://leetcode.com/problems/minimum-cost-to-make-at-least-one-valid-path-in-a-grid/) · **Template:** Template 4 (0-1 BFS)

**Problem.** Each grid cell holds an arrow (1 right, 2 left, 3 down, 4 up). From the top-left you walk along arrows; you may change any cell's arrow for cost 1. Return the minimum cost to reach the bottom-right. `[[1,1,1,1],[2,2,2,2],[1,1,1,1],[2,2,2,2]]` → 3.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** minimum cost, and each move costs either 0 (follow the arrow) or 1 (flip it) → 0-1 BFS.

**Approach:** cells are nodes with four outgoing edges: weight 0 along the cell's arrow, 1 otherwise. A flip only matters for the move leaving that cell, and an optimal path never revisits a cell, so per-move costs add up correctly. Zero-weight neighbors go to the deque's front, unit-weight ones to the back; the deque stays sorted by distance, so the first settled distance is final.

```python
def min_cost_valid_path(grid: list[list[int]]) -> int:
    """Minimum Cost to Make at Least One Valid Path in a Grid (0-1 BFS).

    Cell value 1..4 = arrow right, left, down, up. Following the arrow costs 0,
    any other direction costs 1 (changing the sign).
    """
    rows, cols = len(grid), len(grid[0])
    dirs = {1: (0, 1), 2: (0, -1), 3: (1, 0), 4: (-1, 0)}
    dist = [[float("inf")] * cols for _ in range(rows)]
    dist[0][0] = 0
    dq = deque([(0, 0)])
    while dq:
        r, c = dq.popleft()
        for sign, (dr, dc) in dirs.items():
            nr, nc = r + dr, c + dc
            if 0 <= nr < rows and 0 <= nc < cols:
                w = 0 if grid[r][c] == sign else 1
                if dist[r][c] + w < dist[nr][nc]:
                    dist[nr][nc] = dist[r][c] + w
                    if w == 0:
                        dq.appendleft((nr, nc))
                    else:
                        dq.append((nr, nc))
    return int(dist[rows - 1][cols - 1])
```

**Complexity:** O(rows · cols) time and space: each cell has 4 edges, and each is relaxed O(1) times.

**Edge cases:** 1x1 grid → 0; arrows pointing off the grid are just unusable 0-edges; a cell can be pushed twice (stale entries are harmless, the `<` check filters them).

</details>

### Template 5: Floyd-Warshall (small n, all pairs)

> - **Example problem:** [Find the City With the Smallest Number of Neighbors at a Threshold Distance (1334)](https://leetcode.com/problems/find-the-city-with-the-smallest-number-of-neighbors-at-a-threshold-distance/)
> - **Why this template:** n <= 100 and the answer needs every pairwise distance, so an O(n^3) table is simplest.
> - **Worked below:** Example 5.1

```python
def floyd_warshall(n: int, edges: list[tuple[int, int, int]]) -> list[list[float]]:
    """All-pairs shortest paths, O(n^3). Directed edges; add both ways for undirected."""
    INF = float("inf")
    d = [[0 if i == j else INF for j in range(n)] for i in range(n)]
    for u, v, w in edges:
        d[u][v] = min(d[u][v], w)
    for k in range(n):                        # k MUST be the outer loop
        dk = d[k]
        for i in range(n):
            dik = d[i][k]
            if dik == INF:
                continue
            di = d[i]
            for j in range(n):
                if dik + dk[j] < di[j]:
                    di[j] = dik + dk[j]
    return d
```

Fine for n ≤ ~400 in Python. A negative `d[i][i]` afterwards means a negative cycle through i. The same loops over booleans compute transitive closure (Course Schedule IV queries).

#### Example 5.1: City with the fewest reachable neighbors [Medium]

**Problem:** [Find the City With the Smallest Number of Neighbors at a Threshold Distance (1334)](https://leetcode.com/problems/find-the-city-with-the-smallest-number-of-neighbors-at-a-threshold-distance/) · **Template:** Template 5 (Floyd-Warshall)

**Problem.** n ≤ 100 cities, undirected weighted roads. For each city, count the other cities within shortest-path distance ≤ threshold. Return the city with the smallest count; ties go to the largest index. n=4, edges `[[0,1,3],[1,2,1],[1,3,4],[2,3,1]]`, threshold 4 → 3.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** every city's distances to every other city, with n ≤ 100 → all-pairs; O(n^3) = 10^6 is fine.

**Approach:** run Floyd-Warshall with each road added in both directions, then count `d[i][j] <= threshold` per row. Scanning i upward with `<=` lets a later (larger) city win ties. Dijkstra from every node also works (O(n · E log n)), but Floyd is shorter.

```python
def find_the_city(n: int, edges: list[list[int]], threshold: int) -> int:
    """City With the Smallest Number of Neighbors at a Threshold Distance (Floyd)."""
    d = floyd_warshall(n, [(u, v, w) for u, v, w in edges] + [(v, u, w) for u, v, w in edges])
    best_city, best_cnt = -1, n + 1
    for i in range(n):
        cnt = sum(1 for j in range(n) if j != i and d[i][j] <= threshold)
        if cnt <= best_cnt:                   # '<=' so ties go to the larger index
            best_city, best_cnt = i, cnt
    return best_city
```

**Complexity:** O(n^3) time, O(n^2) space.

**Edge cases:** disconnected cities (INF never ≤ threshold); all counts equal → n-1; parallel roads (Template 5 keeps the minimum).

</details>

### Template 6: Kruskal (sort edges + DSU)

> - **Example problem:** [Min Cost to Connect All Points (1584)](https://leetcode.com/problems/min-cost-to-connect-all-points/)
> - **Why this template:** "minimum cost to connect all points" is an MST; sort pairwise edges by cost and keep each that joins two DSU sets.
> - **Worked below:** Example 6.1

```python
def kruskal(n: int, edges: list[tuple[int, int, int]]) -> int:
    """MST total weight via Kruskal. edges = (w, u, v). Returns -1 if disconnected."""
    dsu = DSU(n)
    total = used = 0
    for w, u, v in sorted(edges):             # cheapest edge first
        if dsu.union(u, v):                   # skip edges that would form a cycle
            total += w
            used += 1
            if used == n - 1:
                break
    return total if used == n - 1 else -1
```

#### Example 6.1: MST of points via Kruskal [Medium]

**Problem:** [Min Cost to Connect All Points (1584)](https://leetcode.com/problems/min-cost-to-connect-all-points/) · **Template:** Template 6 (Kruskal)

**Problem.** Given points on a plane, connecting two costs their Manhattan distance. Return the minimum total cost to connect all of them. `[[0,0],[2,2],[3,10],[5,2],[7,0]]` → 20.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** connect everything at minimum total cost → MST. With a DSU at hand, Kruskal is the shortest correct code.

**Approach:** build all n(n-1)/2 edges `(dist, i, j)`, then run Template 6: sorted cheapest-first, keep each edge that joins two DSU components (cut property: the cheapest edge across any cut is in some MST). Stop at n-1 edges.

```python
def min_cost_connect_points_kruskal(points: list[list[int]]) -> int:
    """Same problem via Kruskal over all O(n^2) edges: O(n^2 log n)."""
    n = len(points)
    edges = [(abs(points[i][0] - points[j][0]) + abs(points[i][1] - points[j][1]), i, j)
             for i in range(n) for j in range(i + 1, n)]
    return kruskal(n, edges) if n > 1 else 0
```

**Complexity:** O(n^2 log n) time (sorting n^2/2 edges), O(n^2) space. On this complete graph, array Prim (Example 7.1) is O(n^2) time and O(n) space, so it wins at n = 1000. Kruskal is the better default on sparse edge lists, where E log E is small and the DSU does the cycle checks.

**Edge cases:** one point → 0 (no edges; guarded); duplicate points give zero-cost edges, which are fine.

</details>

### Template 7: Prim (lazy, with heap)

> - **Example problem:** [Min Cost to Connect All Points (1584)](https://leetcode.com/problems/min-cost-to-connect-all-points/)
> - **Why this template:** "connect all points at minimum total cost" is an MST; grow one tree by always taking the cheapest outgoing edge.
> - **Worked below:** Example 7.1

```python
def prim(n: int, adj: list[list[tuple[int, int]]]) -> int:
    """MST total weight via lazy Prim from node 0. adj[u] = [(v, w)]. -1 if disconnected."""
    seen = [False] * n
    heap = [(0, 0)]                           # (edge weight into tree, node)
    total = added = 0
    while heap and added < n:
        w, u = heapq.heappop(heap)
        if seen[u]:
            continue                          # stale: u already joined via a cheaper edge
        seen[u] = True
        total += w
        added += 1
        for v, wv in adj[u]:
            if not seen[v]:
                heapq.heappush(heap, (wv, v))
    return total if added == n else -1
```

Prim differs from Dijkstra only in the heap key: `w` (the edge into the tree) instead of `dist[u] + w`. Kruskal is simpler given an edge list (sparse graphs); heapless array Prim is O(V^2), best for dense or complete graphs.

#### Example 7.1: Minimum spanning tree on a complete graph [Medium]

**Problem:** [Min Cost to Connect All Points (1584)](https://leetcode.com/problems/min-cost-to-connect-all-points/) · **Template:** Template 7 (Prim), array-based O(n^2) form

Problem: points on a plane; connecting two costs their Manhattan distance. Find the minimum total cost to connect all points.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** connect all at minimum total cost, any structure → MST. Every pair is an edge, so E = n^2/2.

**Approach:** on a complete graph, array Prim is O(n^2) time and O(n) memory: each step adds the outside node with the cheapest `best[v]` (edge from the tree), then updates the rest.

```python
def min_cost_connect_points_prim(points: list[list[int]]) -> int:
    """Min Cost to Connect All Points: dense graph -> O(n^2) array Prim, no heap."""
    n = len(points)
    in_tree = [False] * n
    best = [float("inf")] * n                 # cheapest edge from tree to each node
    best[0] = 0
    total = 0
    for _ in range(n):
        u = min((i for i in range(n) if not in_tree[i]), key=best.__getitem__)
        in_tree[u] = True
        total += best[u]
        ux, uy = points[u]
        for v in range(n):
            if not in_tree[v]:
                d = abs(ux - points[v][0]) + abs(uy - points[v][1])
                if d < best[v]:
                    best[v] = d
    return int(total)
```

Kruskal also works, with O(n^2) memory for the edge list: see Example 6.1.

**Complexity:** Prim O(n^2) time, O(n) space. Kruskal (Example 6.1) O(n^2 log n) time, O(n^2) space.

**Edge cases:** one point → 0; duplicate points (zero-cost edges are fine).

</details>

## Template cheat sheet

All templates, for warm-ups and review.

### Template 1: DSU (path compression + union by size)

```python
class DSU:
    """Disjoint Set Union with path compression (halving) and union by size.

    Amortized near-O(1) per operation (inverse Ackermann).
    """

    def __init__(self, n: int) -> None:
        self.parent = list(range(n))
        self.size = [1] * n
        self.count = n                        # number of disjoint sets

    def find(self, x: int) -> int:
        while self.parent[x] != x:
            self.parent[x] = self.parent[self.parent[x]]   # path halving
            x = self.parent[x]
        return x

    def union(self, a: int, b: int) -> bool:
        """Merge the sets of a and b. Returns False if already in the same set."""
        ra, rb = self.find(a), self.find(b)
        if ra == rb:
            return False
        if self.size[ra] < self.size[rb]:     # attach smaller tree under larger
            ra, rb = rb, ra
        self.parent[rb] = ra
        self.size[ra] += self.size[rb]
        self.count -= 1
        return True

    def connected(self, a: int, b: int) -> bool:
        return self.find(a) == self.find(b)
```

### Template 2: Dijkstra with stale-entry skipping

```python
def dijkstra(n: int, adj: list[list[tuple[int, int]]], src: int) -> list[float]:
    """Single-source shortest paths, non-negative weights. adj[u] = [(v, w), ...].

    Invariant: when a node is popped with d == dist[node], that distance is final.
    Stale entries (d > dist[node]) are skipped instead of decrease-key.
    """
    dist = [float("inf")] * n
    dist[src] = 0
    heap = [(0, src)]
    while heap:
        d, u = heapq.heappop(heap)
        if d > dist[u]:
            continue                          # stale entry
        # TODO: early exit here if u is the target
        for v, w in adj[u]:
            nd = d + w                        # TODO: path cost combiner (sum, max, product...)
            if nd < dist[v]:
                dist[v] = nd
                heapq.heappush(heap, (nd, v))
    return dist
```

### Template 3: Bellman-Ford with an edge limit

```python
def bellman_ford_k_edges(n: int, edges: list[tuple[int, int, int]], src: int, k: int) -> list[float]:
    """Shortest paths using AT MOST k edges. Handles negative weights.

    Each round relaxes from a snapshot of the previous round so a path cannot
    use more than one new edge per round. With k = n - 1 this is classic Bellman-Ford.
    """
    dist = [float("inf")] * n
    dist[src] = 0
    for _ in range(k):
        prev = dist[:]                        # snapshot: the crucial line for edge limits
        for u, v, w in edges:
            if prev[u] + w < dist[v]:
                dist[v] = prev[u] + w
    return dist
```

```python
def has_negative_cycle(n: int, edges: list[tuple[int, int, int]]) -> bool:
    """After n-1 rounds, any further improvement means a reachable negative cycle."""
    dist = [0] * n                            # virtual source connected to all nodes
    for _ in range(n - 1):
        for u, v, w in edges:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
    return any(dist[u] + w < dist[v] for u, v, w in edges)
```

### Template 4: 0-1 BFS

```python
def zero_one_bfs(n: int, adj: list[list[tuple[int, int]]], src: int) -> list[float]:
    """Shortest paths when every weight is 0 or 1: deque instead of heap, O(V + E)."""
    dist = [float("inf")] * n
    dist[src] = 0
    dq = deque([src])
    while dq:
        u = dq.popleft()
        for v, w in adj[u]:
            if dist[u] + w < dist[v]:
                dist[v] = dist[u] + w
                if w == 0:
                    dq.appendleft(v)          # same distance layer: front
                else:
                    dq.append(v)              # next layer: back
    return dist
```

### Template 5: Floyd-Warshall (small n, all pairs)

```python
def floyd_warshall(n: int, edges: list[tuple[int, int, int]]) -> list[list[float]]:
    """All-pairs shortest paths, O(n^3). Directed edges; add both ways for undirected."""
    INF = float("inf")
    d = [[0 if i == j else INF for j in range(n)] for i in range(n)]
    for u, v, w in edges:
        d[u][v] = min(d[u][v], w)
    for k in range(n):                        # k MUST be the outer loop
        dk = d[k]
        for i in range(n):
            dik = d[i][k]
            if dik == INF:
                continue
            di = d[i]
            for j in range(n):
                if dik + dk[j] < di[j]:
                    di[j] = dik + dk[j]
    return d
```

### Template 6: Kruskal (sort edges + DSU)

```python
def kruskal(n: int, edges: list[tuple[int, int, int]]) -> int:
    """MST total weight via Kruskal. edges = (w, u, v). Returns -1 if disconnected."""
    dsu = DSU(n)
    total = used = 0
    for w, u, v in sorted(edges):             # cheapest edge first
        if dsu.union(u, v):                   # skip edges that would form a cycle
            total += w
            used += 1
            if used == n - 1:
                break
    return total if used == n - 1 else -1
```

### Template 7: Prim (lazy, with heap)

```python
def prim(n: int, adj: list[list[tuple[int, int]]]) -> int:
    """MST total weight via lazy Prim from node 0. adj[u] = [(v, w)]. -1 if disconnected."""
    seen = [False] * n
    heap = [(0, 0)]                           # (edge weight into tree, node)
    total = added = 0
    while heap and added < n:
        w, u = heapq.heappop(heap)
        if seen[u]:
            continue                          # stale: u already joined via a cheaper edge
        seen[u] = True
        total += w
        added += 1
        for v, wv in adj[u]:
            if not seen[v]:
                heapq.heappush(heap, (wv, v))
    return total if added == n else -1
```

## Variations and extensions

**Online connectivity (LeetCode: Number of Islands II).** Report the island count after each land addition. DFS per addition is O(k * mn). With DSU, a new cell adds an island and each successful union with a land neighbor removes one.

```python
def num_islands2(m: int, n: int, positions: list[list[int]]) -> list[int]:
    """Number of Islands II: island count after each land addition (online)."""
    dsu = DSU(m * n)
    land = [False] * (m * n)
    islands = 0
    out: list[int] = []
    for r, c in positions:
        idx = r * n + c
        if not land[idx]:                     # duplicate additions change nothing
            land[idx] = True
            islands += 1
            for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
                nr, nc = r + dr, c + dc
                if 0 <= nr < m and 0 <= nc < n and land[nr * n + nc]:
                    if dsu.union(idx, nr * n + nc):
                        islands -= 1
        out.append(islands)
    return out
```

**Equivalence constraints (LeetCode: Satisfiability of Equality Equations).** Union all `==`, then verify all `!=`. Checking `!=` in the first pass misses later equalities.

```python
def equations_possible(equations: list[str]) -> bool:
    """Satisfiability of Equality Equations: union all '==', then check every '!='."""
    dsu = DSU(26)
    for e in equations:
        if e[1] == "=":
            dsu.union(ord(e[0]) - 97, ord(e[3]) - 97)
    return all(not dsu.connected(ord(e[0]) - 97, ord(e[3]) - 97)
               for e in equations if e[1] == "!")
```

**Swaps generate permutations (LeetCode: Smallest String With Swaps).** With unlimited swaps, any index in a component can get any of its characters. Sort characters per component and place them at the sorted indices.

```python
def smallest_string_with_swaps(s: str, pairs: list[list[int]]) -> str:
    """Smallest String With Swaps: within a component, any permutation is reachable."""
    dsu = DSU(len(s))
    for a, b in pairs:
        dsu.union(a, b)
    comp: dict[int, list[int]] = defaultdict(list)
    for i in range(len(s)):
        comp[dsu.find(i)].append(i)           # indices appended in increasing order
    out = list(s)
    for idxs in comp.values():
        for i, ch in zip(idxs, sorted(s[i] for i in idxs)):
            out[i] = ch
    return "".join(out)
```

**0-1 BFS on a grid (LeetCode: Minimum Cost to Make at Least One Valid Path in a Grid).** See Example 4.1.

**All-pairs with a threshold (LeetCode: Find the City With the Smallest Number of Neighbors at a Threshold Distance).** See Example 5.1.

Other combinations and follow-ups:

- **Weighted DSU** (Evaluate Division): store the ratio to the parent; multiply along `find`.
- **DSU on values, not indices** (Longest Consecutive Sequence, Most Stones Removed): map values or rows/columns to nodes via a dict.
- **Offline deletions**: reverse time so deletions become additions, then DSU.
- **Virtual node trick** (Optimize Water Distribution in a Village): a well at cost c_i = edge of weight c_i from virtual node 0 to i; then plain MST.
- **Minimum operations to connect a network**: `components - 1` cables, needing at least n-1 edges total.
- **Dijkstra with state**: `(node, fuel)`, `(node, discounts_left)`, `(cell, k_left)` as heap entries and `dist` keys.
- **Max-probability path**: Dijkstra with a max-heap on products (push `-p`); valid because probabilities ≤ 1 never improve a path.
- **Path reconstruction**: set `parent[v] = u` on each relaxation.
- "Graph changes between queries?" DSU for additions; for weight changes, rerun Dijkstra (dynamic shortest paths are out of scope).

## Common pitfalls

- DSU without both optimizations on 10^5 elements: chains make `find` O(n).
- Grouping by `parent[x]` instead of `find(x)` at the end. Only roots are reliable.
- Recursive `find` hitting the recursion limit before compression happens.
- Dijkstra with negative edges: silently wrong. Use Bellman-Ford (or reweight; Johnson's algorithm is out of scope).
- Omitting the stale check `if d > dist[u]: continue`. Still correct, but redundant relaxations can blow up.
- Marking Dijkstra nodes visited on push (BFS habit). Finalize on pop.
- Bellman-Ford with a hop limit but no snapshot.
- Floyd-Warshall with k as an inner loop. k must be outermost.
- `float("inf")` leaking into an integer answer; check before returning.
- MST on a directed graph: Kruskal/Prim don't apply (that's the arborescence problem).
- Confusing MST with a shortest path tree: an MST minimizes total weight, not distance from a source; its path between two nodes usually isn't the shortest.

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| Dijkstra | BFS (`13`) | all edges cost the same |
| Dijkstra | 0-1 BFS | costs are only 0 or 1 |
| Dijkstra | Bellman-Ford | negative weights or a hop limit |
| Dijkstra | DP on DAG (`16`) | edges only go forward (grid moving right/down, time-ordered events) |
| DSU | BFS/DFS (`13`) | directed reachability, distances, or deletions |
| MST | Shortest path | a source and destination are named; one path matters, not total wiring |
| MST | Steiner tree | only a subset of nodes must connect (NP-hard, unlikely in interviews) |
| Floyd-Warshall | Repeated Dijkstra | n in the thousands; Floyd is O(n^3) |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Number of Provinces (547)](https://leetcode.com/problems/number-of-provinces/) · Medium · DSU basics: DSU or DFS on the matrix.
2. [Redundant Connection (684)](https://leetcode.com/problems/redundant-connection/) · Medium · DSU cycle detection: first failing union.
3. [Network Delay Time (743)](https://leetcode.com/problems/network-delay-time/) · Medium · Dijkstra: answer is the max of distances.
4. [Cheapest Flights Within K Stops (787)](https://leetcode.com/problems/cheapest-flights-within-k-stops/) · Medium · Bellman-Ford / BFS with stops: snapshot each round, k+1 rounds.
5. [Min Cost to Connect All Points (1584)](https://leetcode.com/problems/min-cost-to-connect-all-points/) · Medium · MST (Prim/Kruskal): dense graph, array Prim.

### Optional

**Union-Find**
- [Number of Connected Components in an Undirected Graph (323)](https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph/) · Medium · Premium · n minus successful unions.
- [Graph Valid Tree (261)](https://leetcode.com/problems/graph-valid-tree/) · Medium · Premium · Exactly n - 1 edges and no union joins an already-connected pair.
- [Satisfiability of Equality Equations (990)](https://leetcode.com/problems/satisfiability-of-equality-equations/) · Medium · Union `==` first, then check `!=`.
- [Accounts Merge (721)](https://leetcode.com/problems/accounts-merge/) · Medium · Union via an email -> account map.
- [Smallest String With Swaps (1202)](https://leetcode.com/problems/smallest-string-with-swaps/) · Medium · Sort chars within each component.
- [Number of Operations to Make Network Connected (1319)](https://leetcode.com/problems/number-of-operations-to-make-network-connected/) · Medium · Components - 1, if enough edges.
- [Most Stones Removed with Same Row or Column (947)](https://leetcode.com/problems/most-stones-removed-with-same-row-or-column/) · Medium · n - components; union rows with columns.
- [Evaluate Division (399)](https://leetcode.com/problems/evaluate-division/) · Medium · Weighted DSU or BFS on ratios.
- [Number of Islands II (305)](https://leetcode.com/problems/number-of-islands-ii/) · Hard · Premium · Online additions; DSU per new cell.

**Shortest paths**
- [Path With Minimum Effort (1631)](https://leetcode.com/problems/path-with-minimum-effort/) · Medium · Dijkstra with `max` combiner.
- [Path with Maximum Probability (1514)](https://leetcode.com/problems/path-with-maximum-probability/) · Medium · Max-heap on products.
- [Find the City With the Smallest Number of Neighbors at a Threshold Distance (1334)](https://leetcode.com/problems/find-the-city-with-the-smallest-number-of-neighbors-at-a-threshold-distance/) · Medium · Floyd-Warshall, n ≤ 100.
- [Swim in Rising Water (778)](https://leetcode.com/problems/swim-in-rising-water/) · Hard · Bottleneck: Dijkstra, binary search + BFS, or DSU.
- [Minimum Cost to Make at Least One Valid Path in a Grid (1368)](https://leetcode.com/problems/minimum-cost-to-make-at-least-one-valid-path-in-a-grid/) · Hard · 0-1 BFS.
- [Minimum Obstacle Removal to Reach Corner (2290)](https://leetcode.com/problems/minimum-obstacle-removal-to-reach-corner/) · Hard · 0-1 BFS with obstacles as weight 1.

**MST**
- [Connecting Cities With Minimum Cost (1135)](https://leetcode.com/problems/connecting-cities-with-minimum-cost/) · Medium · Premium · Kruskal, check connectivity.
- [Optimize Water Distribution in a Village (1168)](https://leetcode.com/problems/optimize-water-distribution-in-a-village/) · Hard · Premium · Virtual node + MST.

## Self-check

1. What do path compression and union by size each prevent? What's the complexity with both?
2. Why does Dijkstra fail with negative weights? Give a 3-node counterexample.
3. Why is skipping a popped entry with `d > dist[u]` safe, and what happens without that check?
4. In Cheapest Flights Within K Stops, why does Dijkstra on node alone fail, and why does Bellman-Ford need a snapshot of the previous round?
5. Why must k be the outermost loop in Floyd-Warshall?
6. State the cut property and how Kruskal and Prim each rely on it. When would you pick Prim over Kruskal?
7. Name three correct approaches to "minimize the maximum edge on a path from A to B".
