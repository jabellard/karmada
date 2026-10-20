# 17. Dynamic Programming Advanced

Module 16's states were an index, an index pair, or a capacity. Hard problems use the same 7-step method with less obvious states: a range `[i, j]`, a "mode", a visited set, a subtree, a digit position. This module catalogs them.

## At a glance

- **Use when:** an index, index pair, or capacity can't capture what the next decision needs.
- **Templates:**
  - Template 1: LIS: chain with an ordering constraint
  - Template 2: Interval DP: merge/burst/cut on a contiguous range
  - Template 3: State machine DP: time steps with a few modes
  - Template 4: Bitmask DP: n <= 20 and which items are used matters
  - Template 5: Tree DP: parent-child constraint on a tree
  - Template 6: Memoized DFS on a DAG: strictly monotonic moves
  - Template 7: Digit DP: count integers <= N with a digit property
- **Typical complexity:** states × transitions per state: O(n log n) LIS, O(n^3) interval, O(2^n·n^2) bitmask, O(n) tree.
- **Top trap:** a state missing what the next decision needs, or memoizing on a graph with cycles.

## Core idea

Every DP state answers: **what is the minimum information about the past that the future needs?**

| State shape | What the future needs to know | Example |
|---|---|---|
| `i` (ending at i) | the last chosen element and its value | LIS |
| `[i, j]` range | what is left inside a contiguous block | burst balloons, merge stones |
| `(i, mode)` | which of a few situations you are in | stock buy/sell with cooldown |
| `(mask, last)` | which items are used, maybe where you are | TSP, assignment |
| node -> tuple | the subtree's best under each condition the parent cares about | house robber III, cameras |
| `(pos, tight, ...)` | digit position, whether still capped by N | digit DP |
| DAG node | best path from here (no cycles) | longest increasing path in a matrix |

**Invariant** (as in module 16): a state is computed only after its dependencies are final. Each sub-pattern ensures this its own way: increasing range length, time, mask value (a subset is numerically smaller than its superset), post-order, or memoized recursion on an acyclic graph.

## Recognition signals

- **Subsequence with an ordering constraint** ("strictly increasing", "each fits inside the previous"), n up to 10^5: LIS with binary search.
- **Merge/burst/cut on a contiguous array, cost from neighbors or the block's total**, n <= ~500: interval DP, O(n^3).
- **Time steps with a few modes** (holding, cooldown, k transactions left): state machine DP.
- **n <= 20 (sometimes 16) and "visit all", "assign each", "partition into groups"**: bitmask DP, O(2^n * n) or O(2^n * n^2).
- **Tree with a parent-child constraint** (adjacent nodes not both chosen; every node covered by a neighbor): tree DP returning a tuple.
- **Moves go strictly "uphill"** (values or time increase): a DAG, so memoized DFS needs no visited set.
- **"Count integers in [L, R] with a digit property"**, R up to 10^18: digit DP.
- **"Return modulo 10^9 + 7"**: counting DP with an astronomical count.

## Disguises

1. **Version compatibility chain.** "Each library release has an API level and a size. Find the longest upgrade chain where both strictly increase."
   Decode: 2D strict nesting = Russian Doll Envelopes: sort one dimension ascending (ties descending), LIS on the other.

2. **Compaction cost.** "An LSM tree has a row of sorted runs. A compaction merges exactly k adjacent runs into one, costing their total size. Minimum cost to reach a single run?"
   Decode: adjacent merges + cost = block total = interval DP (Minimum Cost to Merge Stones).

3. **Spot instance trading.** "Hold at most one reserved instance at a time; daily prices are known. After releasing one, wait a day before acquiring another. Maximize savings."
   Decode: days + modes (holding / just released / free) = state machine DP (stock with cooldown).

4. **Visit every datacenter.** "A mesh of 12 datacenters with link latencies. Find the minimum-latency route that starts anywhere and visits all of them."
   Decode: n = 12 + "visit all" = bitmask DP over (visited set, current node). With unit weights, BFS over the same state.

5. **Monitoring agents on a service tree.** "In a service dependency tree, an agent on a node monitors itself, its parent, and its children. Minimum agents to monitor every node?"
   Decode: tree + neighbor coverage = tree DP with 3 states (Binary Tree Cameras).

6. **Valid serial numbers.** "Count serial numbers up to N (N up to 10^15) with no repeated digit."
   Decode: count integers <= N with a digit property = digit DP with (pos, tight, started, used-digit mask).

## Templates and worked examples

Each template gives the dp meaning, transition, and order; fill in the TODOs.

### How to find the state when stuck

Work through in order; usually one item unlocks it.

1. **Write the brute-force recursion first.** Its parameters are a candidate state. If one is a whole list or path, ask which summary of it matters.
2. **What changes between subproblems?** Index, two indices, range, remaining budget, used items?
3. **What must I remember to choose the next move legally?** Last value (LIS), whether I hold a stock, which items are used, whether my parent has a camera. That goes into the state.
4. **What is the last decision?** Each option names a predecessor state. If your state can't name it, add the missing information.
5. **Prefix ("up to i") or ending-at-i?** Use ending-at if the next step needs to know where the last piece ended.
6. **Contiguous block that shrinks or splits?** Try `[i, j]` and "which element is processed **last** in this block" (often makes the halves independent).
7. **Is n tiny (<= 20)?** A bitmask of items is affordable.
8. **Bound the state space.** Multiply the parameter ranges. Above ~10^7, drop a parameter, find a monotonic structure, or decompose differently.
9. **Add a dimension, then try to remove it.** Start with `dp[i][j][k]`, then notice `k` is derivable (the worker index in assignment DP = popcount of mask).

### Template 1: LIS (ending-at-i, and patience sorting)

> - **Example problem:** [Longest Increasing Subsequence (300)](https://leetcode.com/problems/longest-increasing-subsequence/)
> - **Why this template:** Each element of a strictly increasing subsequence extends one ending at an earlier, smaller element: dp[i] ending at i, or `tails` with binary search.
> - **Worked below:** Example 1.1
> - **Also fits:** [Russian Doll Envelopes (354)](https://leetcode.com/problems/russian-doll-envelopes/)

```text
dp[i]   = length of best chain ENDING at i
dp[i]   = 1 + max(dp[j] for j < i if compatible(j, i))       O(n^2)
answer  = max(dp)
fast    : tails[k] = smallest tail of any chain of length k+1; binary search per element
```

#### Example 1.1: Longest increasing subsequence and Russian doll envelopes [Medium/Hard]

**Problems:** [Longest Increasing Subsequence (300)](https://leetcode.com/problems/longest-increasing-subsequence/), [Russian Doll Envelopes (354)](https://leetcode.com/problems/russian-doll-envelopes/) · **Template:** Template 1 (LIS)

Problem (LeetCode: Longest Increasing Subsequence): length of the longest strictly increasing subsequence. Problem (LeetCode: Russian Doll Envelopes): max envelopes that nest, where nesting needs both width and height strictly larger.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** subsequence + ordering constraint; n up to 10^5 rules out O(n^2), so patience sorting.

**Approach:** O(n^2): dp[i] = LIS ending at i, over j < i with `nums[j] < nums[i]`; answer `max(dp)`.

O(n log n): `tails[k]` = smallest tail of any increasing subsequence of length k+1. `tails` is strictly increasing (a longer chain's tail can't undercut a shorter chain's best tail). For each x, replace the first tail >= x (same length, smaller tail, more extendable), or append if x beats all. The answer is `len(tails)`; `tails` itself is **not** an LIS.

```python
def lis_quadratic(nums: list[int]) -> int:
    """dp[i] = length of the longest strictly increasing subsequence ENDING at i.
    Answer is max(dp), not dp[-1]."""
    if not nums:
        return 0
    dp = [1] * len(nums)
    for i in range(len(nums)):
        for j in range(i):
            if nums[j] < nums[i]:
                dp[i] = max(dp[i], dp[j] + 1)
    return max(dp)


def lis_patience(nums: list[int]) -> int:
    """tails[k] = smallest possible tail of an increasing subsequence of
    length k+1. tails is sorted, so each element does one binary search.
    bisect_left -> strictly increasing; bisect_right -> non-decreasing."""
    tails: list[int] = []
    for x in nums:
        k = bisect_left(tails, x)
        if k == len(tails):
            tails.append(x)          # extends the longest subsequence
        else:
            tails[k] = x             # same length, smaller tail: strictly better
    return len(tails)
```

**Envelopes**: sort by width ascending, ties by height **descending**, then LIS on heights. Equal widths now have decreasing heights, so two can never share an increasing run.

```python
def max_envelopes(envelopes: list[list[int]]) -> int:
    """Russian Doll. Sort by width asc, height DESC for equal widths (so two
    envelopes of the same width can never chain), then LIS on heights."""
    envelopes = sorted(envelopes, key=lambda e: (e[0], -e[1]))
    return lis_patience([h for _, h in envelopes])
```

**Complexity:** O(n^2) / O(n log n) time; O(n) space.

**Edge cases:** all equal values give 1 (with `bisect_left`). Non-decreasing variant: `bisect_right`. To reconstruct, store each element's predecessor (the element at `tails[k-1]` when placed).

</details>

### Template 2: Interval DP

> - **Example problem:** [Burst Balloons (312)](https://leetcode.com/problems/burst-balloons/)
> - **Why this template:** Bursting changes neighbors, so the subproblem is a contiguous range, split by its last-burst balloon.
> - **Worked below:** Example 2.1
> - **Also fits:** [Minimum Cost to Merge Stones (1000)](https://leetcode.com/problems/minimum-cost-to-merge-stones/)

```python
def interval_dp_template(n: int, cost) -> int:
    """dp[i][j] = best answer for the closed range [i, j]. Fill by increasing
    length so every sub-range is ready. Split point k = the 'last' or 'first'
    action inside the range. `cost(i, k, j)` is the price of that action."""
    INF = float("inf")
    dp = [[0] * n for _ in range(n)]
    for length in range(2, n + 1):
        for i in range(0, n - length + 1):
            j = i + length - 1
            dp[i][j] = INF
            for k in range(i, j):                       # TODO: split range
                dp[i][j] = min(dp[i][j], dp[i][k] + dp[k + 1][j] + cost(i, k, j))
    return dp[0][n - 1] if n else 0
```

Knobs: what k means (split point, or the **last** element processed), open vs closed ends, cost term (boundary product, block sum via prefix sums), min vs max. Increasing length ensures both halves are done. O(n^3) time, O(n^2) space.

#### Example 2.1: Interval DP: burst balloons, merge stones, palindrome partitioning II [Hard]

**Problems:** [Burst Balloons (312)](https://leetcode.com/problems/burst-balloons/), [Minimum Cost to Merge Stones (1000)](https://leetcode.com/problems/minimum-cost-to-merge-stones/), [Palindrome Partitioning II (132)](https://leetcode.com/problems/palindrome-partitioning-ii/) · **Template:** Template 2 (Interval DP)

**Burst balloons** (LeetCode: Burst Balloons): bursting balloon k earns `left * k * right` with its current neighbors; maximize total.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** cost depends on neighbors that change as you burst; n <= 300-500. Interval DP.

**Approach:** the **first** balloon burst doesn't split the problem (its neighbors merge and interact). The **last** balloon k in the open interval `(i, j)` does: its neighbors are then exactly `i` and `j`, and `(i, k)` and `(k, j)` never interacted. Pad both ends with 1s.

```python
def burst_balloons(nums: list[int]) -> int:
    """Pad with 1s. dp[i][j] = max coins from bursting everything strictly
    between i and j. Choose k = LAST balloon burst in (i, j): at that moment
    its neighbors are exactly i and j, so the two sides are independent."""
    a = [1] + nums + [1]
    n = len(a)
    dp = [[0] * n for _ in range(n)]
    for length in range(2, n):                  # gap j - i
        for i in range(0, n - length):
            j = i + length
            best = 0
            for k in range(i + 1, j):
                best = max(best, dp[i][k] + a[i] * a[k] * a[j] + dp[k][j])
            dp[i][j] = best
    return dp[0][n - 1]
```

**Matrix chain multiplication** is the template with split cost `dims[i] * dims[k+1] * dims[j+1]`:

```python
def matrix_chain(dims: list[int]) -> int:
    """Matrix i has shape dims[i] x dims[i+1]. Min scalar multiplications.
    Split k: (A_i..A_k)(A_{k+1}..A_j) costs dims[i]*dims[k+1]*dims[j+1]."""
    n = len(dims) - 1
    return interval_dp_template(n, lambda i, k, j: dims[i] * dims[k + 1] * dims[j + 1])
```

</details>

**Minimum cost to merge stones** (LeetCode: Minimum Cost to Merge Stones): merge exactly k adjacent piles at a time, cost = their total.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** feasibility first: each merge removes k-1 piles, so `(n-1) % (k-1)` must be 0. dp[i][j] = min cost to reduce `stones[i..j]` to as few piles as possible. Split so the left part `[i..m]` collapses to one pile (m steps by k-1). When `(len-1) % (k-1) == 0`, add the final merge cost (range sum via prefix sums).

```python
def merge_stones(stones: list[int], k: int) -> int:
    """Merge k consecutive piles at a time, cost = their total. Min total cost
    to end with one pile, or -1. dp[i][j] = min cost to reduce stones[i..j] to
    as few piles as possible ((len-1) % (k-1) + 1 piles)."""
    n = len(stones)
    if (n - 1) % (k - 1):
        return -1
    pre = [0]
    for s in stones:
        pre.append(pre[-1] + s)
    dp = [[0] * n for _ in range(n)]
    for length in range(k, n + 1):
        for i in range(0, n - length + 1):
            j = i + length - 1
            # left part [i..m] collapses to exactly one pile: step m by k-1
            dp[i][j] = min(dp[i][m] + dp[m + 1][j] for m in range(i, j, k - 1))
            if (length - 1) % (k - 1) == 0:     # can finish into one pile
                dp[i][j] += pre[j + 1] - pre[i]
    return dp[0][n - 1]
```

</details>

**Palindrome partitioning II** (LeetCode: Palindrome Partitioning II): minimum cuts so every piece is a palindrome.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** interval-flavored but O(n^2) with 1D `cuts[r]` = min cuts for `s[:r+1]`. Expand around every center; each palindrome `s[l..r]` offers `cuts[l-1] + 1` (0 if it starts at 0).

```python
def min_palindrome_cuts(s: str) -> int:
    """Palindrome Partitioning II. cuts[i] = min cuts for s[:i+1]. Expand
    around each center; every palindrome s[l..r] gives cuts[r] <= cuts[l-1]+1."""
    n = len(s)
    cuts = list(range(n))                       # worst case: cut between every char
    for center in range(n):
        for l, r in ((center, center), (center, center + 1)):
            while l >= 0 and r < n and s[l] == s[r]:
                cuts[r] = min(cuts[r], 0 if l == 0 else cuts[l - 1] + 1)
                l -= 1
                r += 1
    return cuts[-1] if n else 0
```

**Complexity:** burst balloons and matrix chain O(n^3) time, O(n^2) space. Merge stones O(n^3 / k). Palindrome cuts O(n^2) time, O(n) space.

**Edge cases:** empty balloon list is 0; a single stone pile costs 0. Skipping the merge-stones feasibility check returns a meaningless number instead of -1.

</details>

### Template 3: State machine DP

> - **Example problem:** [Best Time to Buy and Sell Stock with Cooldown (309)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-cooldown/)
> - **Why this template:** Prices arrive daily and you are always in one of a few modes (holding, just sold, resting); the cooldown rule is one arrow in the state diagram.
> - **Worked below:** Example 3.1
> - **Also fits:** [Best Time to Buy and Sell Stock with Transaction Fee (714)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/)

```text
states  = small set S (e.g. {hold, sold, rest} or {buy_t, sell_t for t in 1..k})
for each time step with input x:
    new[s] = best over transitions (s' -> s) of old[s'] + gain(s', s, x)
answer  = best over terminal states
```

Draw the state diagram first; each arrow is one term in a `max`. Use simultaneous assignment (`a, b = f(a, b), g(a, b)`) or copies so updates read the previous day's values.

#### Example 3.1: State machine DP: stock trading with cooldown, fee, and k transactions [Medium/Hard]

**Problems:** [Best Time to Buy and Sell Stock with Cooldown (309)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-cooldown/), [Best Time to Buy and Sell Stock with Transaction Fee (714)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/), [Best Time to Buy and Sell Stock IV (188)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-iv/) · **Template:** Template 3 (State machine DP)

Problems (LeetCode: Best Time to Buy and Sell Stock with Cooldown / with Transaction Fee / IV): maximize profit from daily prices under rules, holding at most one share.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** one pass over days + one of a few situations = state machine.

**Approach:** cooldown has three end-of-day states.

```
         buy (-p)
  rest ─────────────> hold ──┐ keep
   ^ │ keep            │ <───┘
   │ └──┘              │ sell (+p)
   │                   v
   └──── (next day) ── sold
```

`hold = max(hold, rest - p)`, `sold = hold + p`, `rest = max(rest, sold)`, assigned simultaneously.

```python
def max_profit_cooldown(prices: list[int]) -> int:
    """States after day i: hold, sold (just sold today), rest (free, not holding).
    hold -> hold | sold;  sold -> rest;  rest -> rest | hold."""
    hold, sold, rest = float("-inf"), float("-inf"), 0
    for p in prices:
        hold, sold, rest = max(hold, rest - p), hold + p, max(rest, sold)
    return int(max(sold, rest))


def max_profit_fee(prices: list[int], fee: int) -> int:
    cash, hold = 0, float("-inf")
    for p in prices:
        cash, hold = max(cash, hold + p - fee), max(hold, cash - p)
    return int(cash)
```

**At most k transactions**: add the transaction count to the state: `buy[t]` / `sell[t]` = best balance after the t-th buy / sell. If `2k >= n` the limit never binds: sum every positive daily difference.

```python
def max_profit_k(k: int, prices: list[int]) -> int:
    """At most k transactions. buy[t] / sell[t] = best balance after the t-th
    buy / sell. Updating t in increasing order within a day is safe: a same-day
    buy+sell nets zero profit, so it never inflates the answer."""
    if not prices or k == 0:
        return 0
    if 2 * k >= len(prices):                    # effectively unlimited
        return sum(max(0, b - a) for a, b in zip(prices, prices[1:]))
    buy = [float("-inf")] * (k + 1)
    sell = [0] * (k + 1)
    for p in prices:
        for t in range(1, k + 1):
            buy[t] = max(buy[t], sell[t - 1] - p)
            sell[t] = max(sell[t], buy[t] + p)
    return int(sell[k])
```

**Complexity:** O(n) time, O(1) space for cooldown and fee; O(n * k) time, O(k) space for k transactions.

**Edge cases:** empty or single price gives 0. Start `hold` at `-inf`, not 0 (you can't hold before buying). The table `dp[day][t][holding]` is equivalent; be able to derive the 1D form from it.

</details>

### Template 4: Bitmask DP

> - **Example problem:** [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/)
> - **Why this template:** n <= 16 + "divide into k subsets": the state is the set of used numbers, which fits in a bitmask.
> - **Worked below:** Example 4.1
> - **Also fits:** [Shortest Path Visiting All Nodes (847)](https://leetcode.com/problems/shortest-path-visiting-all-nodes/)

```text
dp[mask]        = best answer having used exactly the items in mask        (assignment, partition)
dp[mask][last]  = best answer having visited mask, currently at last       (TSP / paths)
order           : mask ascending (any subset of mask is numerically smaller)
transition      : for each item u not in mask: dp[mask | 1<<u] <- dp[mask] (+ cost)
complexity      : O(2^n * n) or O(2^n * n^2); n <= 20 or n <= 16 respectively
```

Bit idioms: `mask >> i & 1` (test), `mask | 1 << i` (add), `mask & ~(1 << i)` (remove), `bin(mask).count("1")` or `mask.bit_count()` (size), `(1 << n) - 1` (full set).

#### Example 4.1: Bitmask DP: TSP, assignment, k-equal-subsets, shortest path visiting all nodes [Medium/Hard]

**Problems:** [Shortest Path Visiting All Nodes (847)](https://leetcode.com/problems/shortest-path-visiting-all-nodes/), [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/) · **Template:** Template 4 (Bitmask DP)

**Problems.** (a) In an undirected connected graph with n ≤ 12 nodes, return the length of the shortest walk visiting every node (start and end anywhere; revisits allowed). (b) Decide whether n ≤ 16 positive integers can be split into k equal-sum groups. Also covers small-n TSP and the assignment problem.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** n <= 16-20, and the legal next move depends on **which** items are used, not just how many.

**Approach:** TSP (Held-Karp): dp[mask][v] = min cost of a path from node 0 visiting exactly `mask`, ending at v. Extend to an unvisited u; process masks ascending (adding a bit increases the integer).

```python
def tsp_min_cost(dist: list[list[int]]) -> int:
    """Held-Karp. dp[mask][v] = min cost to start at 0, visit exactly `mask`,
    and stand at v. O(2^n * n^2). Returns min tour cost back to 0."""
    n = len(dist)
    INF = float("inf")
    dp = [[INF] * n for _ in range(1 << n)]
    dp[1][0] = 0
    for mask in range(1 << n):
        if not mask & 1:
            continue
        for v in range(n):
            cur = dp[mask][v]
            if cur == INF or not (mask >> v) & 1:
                continue
            for u in range(n):
                if (mask >> u) & 1:
                    continue
                nxt = mask | (1 << u)
                if cur + dist[v][u] < dp[nxt][u]:
                    dp[nxt][u] = cur + dist[v][u]
    full = (1 << n) - 1
    return min(dp[full][v] + dist[v][0] for v in range(n)) if n > 1 else 0
```

**Shortest path visiting all nodes** (LeetCode: Shortest Path Visiting All Nodes): same state `(node, mask)`, but unit edges and revisits make it a shortest path on the state graph: BFS from every `(v, {v})` at once.

```python
def shortest_path_visiting_all(graph: list[list[int]]) -> int:
    """Unweighted graph, start anywhere, revisits allowed. State = (node, mask).
    Unit edges -> BFS over the state graph (multi-source from every node)."""
    n = len(graph)
    full = (1 << n) - 1
    q = deque((v, 1 << v, 0) for v in range(n))
    seen = {(v, 1 << v) for v in range(n)}
    while q:
        v, mask, d = q.popleft()
        if mask == full:
            return d
        for u in graph[v]:
            state = (u, mask | (1 << u))
            if state not in seen:
                seen.add(state)
                q.append((u, state[1], d + 1))
    return 0
```

**Assignment** (n workers, n jobs, min total cost): dp[mask] over jobs taken. Workers are assigned in order, so the worker index is `popcount(mask)`, not a dimension.

```python
def min_assignment_cost(cost: list[list[int]]) -> int:
    """n workers x n jobs. dp[mask] = min cost when the jobs in `mask` are
    taken by workers 0..popcount(mask)-1. Worker index is implied by popcount."""
    n = len(cost)
    INF = float("inf")
    dp = [INF] * (1 << n)
    dp[0] = 0
    for mask in range(1 << n):
        if dp[mask] == INF:
            continue
        w = bin(mask).count("1")                # next worker to assign
        if w == n:
            continue
        for j in range(n):
            if not (mask >> j) & 1:
                nxt = mask | (1 << j)
                dp[nxt] = min(dp[nxt], dp[mask] + cost[w][j])
    return dp[(1 << n) - 1]
```

**Partition to k equal sum subsets** (LeetCode: Partition to K Equal Sum Subsets): buckets fill one at a time, so the mask determines the current bucket's level, `sum(mask) % target`. Store it in dp[mask] (-1 if unreachable).

```python
def can_partition_k_subsets(nums: list[int], k: int) -> bool:
    """dp[mask] = fill level of the current (partially filled) bucket after
    using `mask`, or -1 if unreachable. Buckets are filled one at a time, so
    the fill level is just sum(mask) % target, and order of buckets is free."""
    total = sum(nums)
    if total % k:
        return False
    target = total // k
    nums = sorted(nums)
    if nums and nums[-1] > target:
        return False
    n = len(nums)
    dp = [-1] * (1 << n)
    dp[0] = 0
    for mask in range(1 << n):
        if dp[mask] < 0:
            continue
        for i in range(n):
            if (mask >> i) & 1:
                continue
            if dp[mask] + nums[i] > target:
                break                           # sorted: later ones overflow too
            nxt = mask | (1 << i)
            if dp[nxt] < 0:
                dp[nxt] = (dp[mask] + nums[i]) % target
    return dp[(1 << n) - 1] == 0
```

**Complexity:** TSP O(2^n * n^2) time, O(2^n * n) space (n = 16: ~16M operations, fine). Assignment and k-subsets O(2^n * n). Shortest path visiting all O(2^n * n^2) worst case (edges per state).

**Edge cases:** n = 1. In TSP, unreachable states stay INF. Assignment with n in the hundreds needs the Hungarian algorithm, not bitmask DP; say so.

</details>

### Template 5: Tree DP

> - **Example problem:** [House Robber III (337)](https://leetcode.com/problems/house-robber-iii/)
> - **Why this template:** A tree where two directly-linked houses can't both be robbed: each node returns (take, skip) to its parent.
> - **Worked below:** Example 5.1
> - **Also fits:** [Binary Tree Cameras (968)](https://leetcode.com/problems/binary-tree-cameras/)

```text
def dfs(node) -> tuple:                  # post-order
    if not node: return BASE_TUPLE
    L, R = dfs(node.left), dfs(node.right)
    # For each state s of node, combine children states compatible with s.
    return (state_0, state_1, ...)
answer = combine(dfs(root))
```

The tuple holds one value per situation the parent cares about (taken / not; camera / covered / uncovered). Post-order finishes children first. O(n) time, O(h) stack.

#### Example 5.1: Tree DP: house robber III and binary tree cameras [Medium/Hard]

**Problems:** [House Robber III (337)](https://leetcode.com/problems/house-robber-iii/), [Binary Tree Cameras (968)](https://leetcode.com/problems/binary-tree-cameras/) · **Template:** Template 5 (Tree DP)

**House robber III** (LeetCode: House Robber III): max sum of tree nodes with no parent-child pair both chosen (max-weight independent set on a tree).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** tree + adjacent-node constraint. The parent needs your best both if taken and if not.

**Approach:** post-order returns `(take, skip)`. Taking a node forces both children skipped; skipping lets each child use its better option.

```python
def rob_tree(root: Optional[TreeNode]) -> int:
    """House Robber III = max-weight independent set on a tree.
    Post-order returns (best if node taken, best if node skipped)."""
    def dfs(node: Optional[TreeNode]) -> tuple[int, int]:
        if not node:
            return 0, 0
        lt, ls = dfs(node.left)
        rt, rs = dfs(node.right)
        take = node.val + ls + rs               # children must be skipped
        skip = max(lt, ls) + max(rt, rs)        # children free
        return take, skip
    return max(dfs(root))
```

</details>

**Binary tree cameras** (LeetCode: Binary Tree Cameras): a camera covers itself, its parent, and its children. Minimum cameras to cover all nodes.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** three states per node, bottom-up: 0 = has camera, 1 = covered without camera, 2 = not covered (parent must have a camera). Any child uncovered: take a camera. Else any child has a camera: covered. Else uncovered, delegating to the parent. Null children count as covered, so leaves report "uncovered" and push cameras upward, where they cover more. If the root ends uncovered, add one.

```python
def min_camera_cover(root: Optional[TreeNode]) -> int:
    """Binary Tree Cameras. Three states per node, returned bottom-up:
      0 = has camera, 1 = covered (no camera), 2 = NOT covered (parent must fix).
    Greedy-DP: put cameras as high as possible (on parents of uncovered leaves)."""
    cameras = 0

    def dfs(node: Optional[TreeNode]) -> int:
        nonlocal cameras
        if not node:
            return 1                            # null counts as covered
        l, r = dfs(node.left), dfs(node.right)
        if l == 2 or r == 2:
            cameras += 1
            return 0
        if l == 0 or r == 0:
            return 1
        return 2

    if dfs(root) == 2:
        cameras += 1
    return cameras
```

This is the greedy collapse of a 3-state tree DP (min cost per state). If camera costs vary per node, return a tuple of three costs instead.

**Complexity:** O(n) time, O(h) recursion.

**Edge cases:** empty tree; a single node needs one camera. Skewed trees with 10^4+ nodes can hit the recursion limit: use iterative post-order or `sys.setrecursionlimit`.

</details>

### Template 6: Memoized DFS on a DAG

> - **Example problem:** [Longest Increasing Path in a Matrix (329)](https://leetcode.com/problems/longest-increasing-path-in-a-matrix/)
> - **Why this template:** Moves go four ways, but an increasing path never returns to a cell, so the move graph is a DAG and memo on (r, c) is valid.
> - **Worked below:** Example 6.1

```text
@cache
def best_from(v):  return base + max/min/sum(best_from(u) for u in next(v))
```

Valid only if `next` can never lead back to `v` (strictly increasing values or time, or a DAG). With cycles, memoization gives wrong answers; use BFS/Dijkstra or a visited-set search.

#### Example 6.1: DAG DP via memoized DFS: longest increasing path in a matrix [Hard]

**Problem:** [Longest Increasing Path in a Matrix (329)](https://leetcode.com/problems/longest-increasing-path-in-a-matrix/) · **Template:** Template 6 (Memoized DFS on a DAG)

Problem (LeetCode: Longest Increasing Path in a Matrix): longest path of strictly increasing values moving in four directions.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** four-direction moves normally create cycles, but "strictly increasing" never returns to a cell: the move graph is a DAG.

**Approach:** `longest_from(r, c)` = 1 + max over larger neighbors, memoized on (r, c), with no visited set.

```python
def longest_increasing_path(matrix: list[list[int]]) -> int:
    """Edges go from smaller to strictly larger neighbor -> a DAG, so no
    visited set is needed and memo on (r, c) is valid."""
    if not matrix or not matrix[0]:
        return 0
    R, C = len(matrix), len(matrix[0])

    @cache
    def longest_from(r: int, c: int) -> int:
        best = 1
        for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            nr, nc = r + dr, c + dc
            if 0 <= nr < R and 0 <= nc < C and matrix[nr][nc] > matrix[r][c]:
                best = max(best, 1 + longest_from(nr, nc))
        return best

    return max(longest_from(r, c) for r in range(R) for c in range(C))
```

**Complexity:** O(R * C) time (each cell once, 4 neighbors), O(R * C) space.

**Edge cases:** all equal values give 1. Recursion depth can reach R * C on a snake-shaped path; iteratively, Kahn's topological sort (module 13) peels cells layer by layer, and the layer count is the answer.

</details>

### Template 7: Digit DP

> - **Example problem:** [Number of Digit One (233)](https://leetcode.com/problems/number-of-digit-one/)
> - **Why this template:** "Count the total number of digit 1 appearing in all non-negative integers less than or equal to n" is a digit property over [0, n]: state (pos, tight, ones).
> - **Worked below:** Example 7.1

```text
digits = digits of N
@cache
def go(pos, tight, started, *extra):
    if pos == len(digits): return value_for_completed_number(started, extra)
    limit = digits[pos] if tight else 9
    return sum(go(pos+1, tight and d == limit, started or d != 0, update(extra, d))
               for d in range(limit + 1) if allowed(d, extra))
count(L..R) = count(R) - count(L-1)
```

`tight` = the prefix so far equals N's, so the next digit is capped. `started` = a nonzero digit has been placed (handles leading zeros when they matter). `extra` = property-specific memory (digit sum, last digit, used-digit mask, count of ones).

#### Example 7.1: Digit DP: total count of digit one up to n [Hard]

**Problem:** [Number of Digit One (233)](https://leetcode.com/problems/number-of-digit-one/) · **Template:** Template 7 (Digit DP)

**Problem.** Given n ≥ 0, count how many times the digit 1 is written across all integers 0, 1, ..., n. n = 13 → 6 (1, 10, 11 twice, 12, 13).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a digit property summed over every integer in [0, n], with n up to 10^9 (too many to enumerate) → digit DP.

**Approach:** build numbers digit by digit, left to right, as length-`len(n)` strings with leading zeros. State `(pos, tight, ones)`: `pos` = next digit index; `tight` = prefix equals n's prefix, so the next digit is capped at `digits[pos]` (else 9); `ones` = 1s placed so far. At the end, a completed number contributes `ones`. Each integer in [0, n] corresponds to exactly one digit string, so the sum is the answer. No `started` flag: leading zeros add no 1s, so padding is harmless.

```python
def count_digit_one(n: int) -> int:
    """Total count of digit '1' over all integers in [0, n].
    State: (pos, tight, ones_so_far). Template 7 shape."""
    if n <= 0:
        return 0
    digits = list(map(int, str(n)))

    @cache
    def go(pos: int, tight: bool, ones: int) -> int:
        if pos == len(digits):
            return ones
        limit = digits[pos] if tight else 9
        total = 0
        for d in range(limit + 1):
            total += go(pos + 1, tight and d == limit, ones + (d == 1))
        return total

    return go(0, True, 0)
```

**Complexity:** O(D^2 · 10) time for D = number of digits (≤ 10): D positions × 2 tight × ≤ D+1 ones states, 10 transitions each. O(D^2) memo space, recursion depth D.

**Edge cases:** n = 0 → 0; n = 1 → 1; powers of ten (n = 10 → 2). For a range [L, R], use `f(R) - f(L-1)`. Alternative without DP: positional counting, summing per power of ten how often that place shows a 1.

</details>

## Template cheat sheet

All templates in one place, for warm-ups and review.

### Template 1: LIS (ending-at-i, and patience sorting)
```text
dp[i]   = length of best chain ENDING at i
dp[i]   = 1 + max(dp[j] for j < i if compatible(j, i))       O(n^2)
answer  = max(dp)
fast    : tails[k] = smallest tail of any chain of length k+1; binary search per element
```

### Template 2: Interval DP
```python
def interval_dp_template(n: int, cost) -> int:
    """dp[i][j] = best answer for the closed range [i, j]. Fill by increasing
    length so every sub-range is ready. Split point k = the 'last' or 'first'
    action inside the range. `cost(i, k, j)` is the price of that action."""
    INF = float("inf")
    dp = [[0] * n for _ in range(n)]
    for length in range(2, n + 1):
        for i in range(0, n - length + 1):
            j = i + length - 1
            dp[i][j] = INF
            for k in range(i, j):                       # TODO: split range
                dp[i][j] = min(dp[i][j], dp[i][k] + dp[k + 1][j] + cost(i, k, j))
    return dp[0][n - 1] if n else 0
```

### Template 3: State machine DP
```text
states  = small set S (e.g. {hold, sold, rest} or {buy_t, sell_t for t in 1..k})
for each time step with input x:
    new[s] = best over transitions (s' -> s) of old[s'] + gain(s', s, x)
answer  = best over terminal states
```

### Template 4: Bitmask DP
```text
dp[mask]        = best answer having used exactly the items in mask        (assignment, partition)
dp[mask][last]  = best answer having visited mask, currently at last       (TSP / paths)
order           : mask ascending (any subset of mask is numerically smaller)
transition      : for each item u not in mask: dp[mask | 1<<u] <- dp[mask] (+ cost)
complexity      : O(2^n * n) or O(2^n * n^2); n <= 20 or n <= 16 respectively
```

### Template 5: Tree DP
```text
def dfs(node) -> tuple:                  # post-order
    if not node: return BASE_TUPLE
    L, R = dfs(node.left), dfs(node.right)
    # For each state s of node, combine children states compatible with s.
    return (state_0, state_1, ...)
answer = combine(dfs(root))
```

### Template 6: Memoized DFS on a DAG
```text
@cache
def best_from(v):  return base + max/min/sum(best_from(u) for u in next(v))
```

### Template 7: Digit DP
```text
digits = digits of N
@cache
def go(pos, tight, started, *extra):
    if pos == len(digits): return value_for_completed_number(started, extra)
    limit = digits[pos] if tight else 9
    return sum(go(pos+1, tight and d == limit, started or d != 0, update(extra, d))
               for d in range(limit + 1) if allowed(d, extra))
count(L..R) = count(R) - count(L-1)
```

## Variations and extensions

### Counting DP with modulo

For counts "modulo 10^9 + 7", reduce at every addition: Python won't overflow, but huge ints are slow, and it shows you understand the requirement. For subtraction, add MOD before reducing: `(a - b) % MOD` is safe in Python but not most languages.

```python
def num_rolls_to_target(n: int, faces: int, target: int) -> int:
    """dp[s] = ways for the dice so far to sum to s. Reduce mod at every add."""
    dp = [1] + [0] * target
    for _ in range(n):
        new = [0] * (target + 1)
        for s in range(target + 1):
            if dp[s]:
                for f in range(1, faces + 1):
                    if s + f > target:
                        break
                    new[s + f] = (new[s + f] + dp[s]) % MOD
        dp = new
    return dp[target]
```

O(n * target * faces); a prefix-sum over the previous row removes the `faces` factor.

### Digit DP

Total digit `1`s up to n: see Example 7.1. A second instance, integers up to n with all distinct digits (checked against brute force in the companion file), needs the `started` flag:

```python
def count_distinct_digit_numbers(n: int) -> int:
    """Count integers in [1, n] whose digits are all distinct (Count Special
    Integers). Digit DP template with tight + started + used-digit mask."""
    digits = list(map(int, str(n)))

    @cache
    def go(pos: int, tight: bool, started: bool, mask: int) -> int:
        if pos == len(digits):
            return 1 if started else 0          # exclude the number 0
        limit = digits[pos] if tight else 9
        total = 0
        for d in range(limit + 1):
            if started and (mask >> d) & 1:
                continue
            if not started and d == 0:          # still in leading zeros
                total += go(pos + 1, tight and d == limit, False, 0)
            else:
                total += go(pos + 1, tight and d == limit, True, mask | (1 << d))
        return total

    return go(0, True, False, 0)
```

States: digits (<= 19) x tight (2) x started (2) x mask (2^10). Leading zeros matter here (`007` must not count as using `0` twice), hence `started`.

### DP + binary search, and monotonic optimizations (brief)

- **Weighted interval scheduling** (LeetCode: Maximum Profit in Job Scheduling): sort by end; dp[i] = best from the first i jobs; taking job i jumps back to the last job ending by its start (binary search).

```python
def job_scheduling(start: list[int], end: list[int], profit: list[int]) -> int:
    """Weighted interval scheduling. Sort by end; dp[i] = best using the first
    i jobs. Take job i-1: add profit to dp[p] where p = number of jobs ending
    at or before its start (binary search on sorted ends)."""
    jobs = sorted(zip(end, start, profit))
    ends = [e for e, _, _ in jobs]
    dp = [0] * (len(jobs) + 1)
    for i, (e, s, p) in enumerate(jobs, 1):
        prev = bisect_right(ends, s, 0, i - 1)   # jobs ending <= s
        dp[i] = max(dp[i - 1], dp[prev] + p)
    return dp[-1]
```

- **Binary search on the answer + greedy/DP check** (Split Array Largest Sum): if "can we achieve X?" is monotonic in X, binary search X (module 07).
- **Monotonic deque over dp** (Constrained Subsequence Sum, Jump Game VI): `dp[i] = nums[i] + max(dp[i-k..i-1])` becomes O(n) with a sliding-window max deque (module 08).
- **Prefix sums in transitions**: when a transition sums a contiguous range of the previous row, prefix sums drop the window-size factor.
- Knuth, divide and conquer, and convex hull trick optimizations cut O(n^3) -> O(n^2) or O(n^2) -> O(n log n). Rare in interviews; mention them if you notice a monotonic optimal split point.

### Interviewer follow-ups

- "Reconstruct the choices" → store the argmax per state and walk back.
- "n is now 10^5" on an O(n^2) DP → a sorted structure (LIS-style binary search), a monotonic deque, or greedy.
- "k transactions, k huge" → collapses to the greedy sum of rises.
- "The tree is a general graph" → independent set on a general graph is NP-hard; the tree makes it linear.

## Common pitfalls

- **LIS: returning `dp[-1]`** instead of `max(dp)`; treating `tails` as the subsequence.
- **Envelopes: heights ascending on ties** lets same-width envelopes nest.
- **Interval DP: wrong loop order.** `i` and `j` both ascending reads unfinished cells. Loop by length, or `i` descending with `j` ascending.
- **Burst balloons: choosing the first balloon** instead of the last.
- **State machine: sequential updates without simultaneous assignment**, so a later update reads today's value, not yesterday's (sometimes harmless; know when).
- **Bitmask DP for n = 25+**: 2^25 * 25 is ~8 * 10^8, too slow in Python. Check n first.
- **Memoized DFS on a graph with cycles**: the memo caches a value computed mid-cycle. Only safe on DAGs.
- **Digit DP forgetting `tight` propagation**: `tight and d == limit`, not `d == limit`.
- **Modulo only at the end**: correct in Python but slow; in other languages, subtraction can go negative.

## When NOT to use / look-alikes

| Looks like | But is | Tell |
|---|---|---|
| Interval DP | Greedy with a heap (merging ropes/sticks) | Any two piles can merge, not just adjacent: Huffman-style greedy |
| Bitmask DP | Plain backtracking | Output is "list all", or pruning suffices |
| Bitmask DP for "visit all" on large n | Different problem class (MST, Euler path) | "Visit every edge" is Eulerian; "connect all" is MST (module 14) |
| Tree DP | Simple DFS aggregation | No parent-child interaction; one return value suffices (sum, height) |
| Memoized grid DFS | BFS / Dijkstra | Moves are not monotonic, so cycles exist |
| LIS | Longest consecutive sequence | Values must be consecutive integers; hash set, O(n) (module 04) |
| Stock state machine | Single pass min-so-far | One transaction: track min price, no states |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Longest Increasing Subsequence (300)](https://leetcode.com/problems/longest-increasing-subsequence/) · Medium · LIS, O(n log n): know both O(n^2) and patience.
2. [Best Time to Buy and Sell Stock with Cooldown (309)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-cooldown/) · Medium · State machine: three states.
3. [House Robber III (337)](https://leetcode.com/problems/house-robber-iii/) · Medium · DP on trees: (take, skip) tuple.
4. [Longest Increasing Path in a Matrix (329)](https://leetcode.com/problems/longest-increasing-path-in-a-matrix/) · Hard · Memoized DFS on a DAG: strictly increasing = DAG; memo DFS.
5. [Burst Balloons (312)](https://leetcode.com/problems/burst-balloons/) · Hard · Interval DP: last balloon in range.

### Optional

**LIS family**
- [Number of Longest Increasing Subsequence (673)](https://leetcode.com/problems/number-of-longest-increasing-subsequence/) · Medium · Carry (length, count) per index.
- [Russian Doll Envelopes (354)](https://leetcode.com/problems/russian-doll-envelopes/) · Hard · Sort with descending tie-break, then LIS.

**Interval DP**
- [Guess Number Higher or Lower II (375)](https://leetcode.com/problems/guess-number-higher-or-lower-ii/) · Medium · Minimax over the split.
- [Minimum Cost to Merge Stones (1000)](https://leetcode.com/problems/minimum-cost-to-merge-stones/) · Hard · Feasibility check and stepping by k-1.
- [Palindrome Partitioning II (132)](https://leetcode.com/problems/palindrome-partitioning-ii/) · Hard · Expand around centers for O(n^2).
- [Minimum Cost to Cut a Stick (1547)](https://leetcode.com/problems/minimum-cost-to-cut-a-stick/) · Hard · Add endpoints, sort cuts, interval over cut indices.

**State machine**
- [Best Time to Buy and Sell Stock with Transaction Fee (714)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-with-transaction-fee/) · Medium · Two states.
- [Best Time to Buy and Sell Stock III (123)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-iii/) · Hard · At most two transactions; transaction count in the state.
- [Best Time to Buy and Sell Stock IV (188)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock-iv/) · Hard · Generalize III to k transactions in the state.

**Bitmask DP**
- [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/) · Medium · dp[mask] = current bucket fill.
- [Shortest Path Visiting All Nodes (847)](https://leetcode.com/problems/shortest-path-visiting-all-nodes/) · Hard · BFS over (node, mask).
- [Find the Shortest Superstring (943)](https://leetcode.com/problems/find-the-shortest-superstring/) · Hard · TSP over overlap costs.
- [Find Minimum Time to Finish All Jobs (1723)](https://leetcode.com/problems/find-minimum-time-to-finish-all-jobs/) · Hard · Iterate submasks; or backtracking with pruning.

**Trees and DAGs**
- [Binary Tree Cameras (968)](https://leetcode.com/problems/binary-tree-cameras/) · Hard · Three coverage states.
- [Binary Tree Maximum Path Sum (124)](https://leetcode.com/problems/binary-tree-maximum-path-sum/) · Hard · Return the one-sided best, track the two-sided best globally.

**Counting, digits, binary search**
- [Number of Dice Rolls With Target Sum (1155)](https://leetcode.com/problems/number-of-dice-rolls-with-target-sum/) · Medium · Counting with modulo.
- [Knight Dialer (935)](https://leetcode.com/problems/knight-dialer/) · Medium · Small state machine over keypad digits, modulo.
- [Number of Digit One (233)](https://leetcode.com/problems/number-of-digit-one/) · Hard · Digit DP or positional math.
- [Count Special Integers (2376)](https://leetcode.com/problems/count-special-integers/) · Hard · Digit DP with used-digit mask and leading zeros.
- [Maximum Profit in Job Scheduling (1235)](https://leetcode.com/problems/maximum-profit-in-job-scheduling/) · Hard · Sort by end, binary search predecessor.

## Self-check

1. In patience sorting, why is `tails` always sorted, and why does replacing an element never shorten the answer?
2. Why must burst balloons choose the **last** balloon in a range, not the first? What goes wrong otherwise?
3. Draw the cooldown state machine from memory and derive the three update equations.
4. In assignment bitmask DP, why can the worker index be dropped from the state?
5. Why can longest increasing path in a matrix be memoized without a visited set, while "longest simple path in a grid" cannot?
6. In digit DP, what do `tight` and `started` each guard against? When can `started` be omitted?
7. Walk through the "find the state" checklist out loud on Russian Doll Envelopes.
