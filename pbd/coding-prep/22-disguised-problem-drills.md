# 22. Disguised Problem Drills

Interviewers rarely say "use a sliding window." They describe a cluster, a log file, or a train schedule, and you find the pattern. This module trains that step: 36 mixed story problems, roughly easy to hard, each with graded hints and a full decode. Every solution is tested in `code/22_disguised_drills.py`.

## 1. The decoding method

Run this aloud on every problem until it's automatic.

1. **Strip the story.** Rewrite as "Given `<input type>`, return `<output>`." Delete nouns like node, pod, shop, train. A detail that changes neither the signature nor a constraint is noise.
2. **Name the objects.** Sequence (ordered? sorted?), intervals, graph (directed? weighted? DAG?), grid, tree, string, stream, weighted items.
3. **Name the question type.** Optimize (min/max), count, decide (yes/no), construct (an order, an assignment), k-th / top-k, all solutions.
4. **Read constraints as a complexity budget.** n ≤ 12-20 → exponential or bitmask. n ≤ 500 → O(n^3). n ≤ 5000 → O(n^2). n ≤ 10^5-10^6 → O(n log n) or O(n). Values up to 10^9 with a monotone check → binary search on the answer.
5. **Hunt trigger phrases** (table below). One rarely decides it; two or three that agree usually do.
6. **Brute force first, then find the repeated work.** State the brute force and its cost, then ask what it recomputes: overlapping windows (sliding window, prefix sums), subproblems (DP), "find min" (heap), "is it connected?" (union-find), "next greater" scans (monotonic stack). The fix is the pattern.

### Story words → abstract meaning

| Story says | Usually means | Module |
|---|---|---|
| "dependencies", "prerequisites", "must run before", "depends_on" | DAG, topological sort, cycle = impossible | `13-graphs-traversal-and-topological-sort.md` |
| "spreads to neighbors each minute", "all at once from several sources" | multi-source BFS, layers = time | `13-graphs-traversal-and-topological-sort.md` |
| "at most k changes / overrides / replacements" in a contiguous run | sliding window with a budget | `06-sliding-window.md` |
| "consecutive days / requests / readings", "contiguous span" | subarray: window or prefix sums | `06-sliding-window.md`, `04-arrays-hashing-prefix-sums.md` |
| "exactly sums to k" with negative values allowed | prefix sum + hash map | `04-arrays-hashing-prefix-sums.md` |
| "minimize the maximum", "smallest capacity such that it finishes in D" | binary search on the answer | `07-binary-search.md` |
| "first build / version where it broke", "monotone yes/no" | boundary binary search | `07-binary-search.md` |
| "next warmer / taller / higher", "until something bigger appears" | monotonic stack | `08-stacks-and-monotonic-structures.md` |
| "peak over every window of k" | monotonic deque | `08-stacks-and-monotonic-structures.md` |
| "merge overlapping maintenance windows", "how many at the same time" | intervals, sort + sweep | `12-intervals-and-sweep-line.md` |
| "groups that are connected / linked / same cluster", "extra link creates a loop" | union-find | `14-union-find-shortest-paths-mst.md` |
| "fastest / cheapest route", non-negative weights | Dijkstra | `14-union-find-shortest-paths-mst.md` |
| "cheapest route with at most k hops / relays" | Bellman-Ford limited to k+1 rounds, or BFS over (node, hops) | `14-union-find-shortest-paths-mst.md` |
| "connect everything at minimum total cost" | minimum spanning tree | `14-union-find-shortest-paths-mst.md` |
| "fewest steps", each step changes one thing | BFS on an implicit state graph | `13-graphs-traversal-and-topological-sort.md` |
| "list all", "every valid configuration" | backtracking | `15-backtracking.md` |
| "number of ways", "how many interpretations" | DP (count) | `16-dynamic-programming-foundations.md` |
| "split into two equal halves", "pick a subset hitting exactly X" | 0/1 knapsack | `16-dynamic-programming-foundations.md` |
| "cooldown", "hold / not hold", "after selling you must wait" | state-machine DP | `17-dynamic-programming-advanced.md` |
| "order of cuts / merges matters, cost depends on current segment" | interval DP | `17-dynamic-programming-advanced.md` |
| "n ≤ 12-20", "assign each to exactly one" | bitmask DP | `17-dynamic-programming-advanced.md` |
| "schedule to finish as early as possible", "always pick the most urgent" | greedy, often with a heap | `18-greedy.md`, `11-heaps-top-k-and-merging.md` |
| "k sorted feeds", "running median", "top k" | heap / two heaps | `11-heaps-top-k-and-merging.md` |
| "prefix typed so far", "autocomplete" | trie | `19-tries-and-string-algorithms.md` |
| "repeats itself", "find pattern in text in O(n)" | KMP prefix function / rolling hash | `19-tries-and-string-algorithms.md` |
| "O(1) extra memory", "everything appears twice except" | XOR / bit tricks, or Floyd cycle detection | `20-bit-manipulation-and-math.md`, `09-linked-lists.md` |
| "hierarchy", "reports to", "uplink / downlink" | tree | `10-trees-and-bst.md` |

**Noise:** prices or labels outside the objective, units (ms vs s), "sorted by priority" when priority is unused, whether it's a GPU or a train. **Never ignore:** "contiguous", "at most k", "exactly", "order is fixed", "each item used once", the size of n, and whether values can be negative.

## 2. Drills

Each drill is a story with collapsed hints and solution. Attempt it blind first (section 4).

### Drill 1: GPU quota pairing

Finance gives your team a weekly quota of GPU-hours and wants it spent exactly, with no leftover. The queue holds pending training jobs, listed in priority order, each with an estimated GPU-hour cost. Each job also has an owning team and a model name. Pick two different jobs whose estimates add up to exactly the quota. Return their indices `[i, j]` with `i < j`, choosing the pair whose second job appears earliest. Return `[]` if no pair works.

- `hours = [5, 12, 7, 3, 9], quota = 16` → `[2, 4]` (7 + 9)
- `hours = [4, 4], quota = 8` → `[0, 1]`
- Constraints: `n ≤ 10^5`, values up to `10^9`.

<details><summary>Hints</summary>

1. Team names, model names, and priority order never enter the answer.
2. Given an integer array and a target, find two indices whose values sum to the target.
3. Hash map of values seen so far (complement lookup).

</details>

<details><summary>Solution</summary>

**Decode:** "Two different jobs" + "add up to exactly" + `n ≤ 10^5` rules out all pairs (5·10^9). Each job's partner value is fixed (`quota - h`): have I seen it?

**Abstract problem:** Two Sum, return indices.

**Pattern:** Arrays and hashing (`04-arrays-hashing-prefix-sums.md`).

**Approach:** Scan left to right, checking each value's complement among earlier values before inserting it. The first hit has the earliest second element.

```python
def quota_pair(hours: list[int], quota: int) -> list[int]:
    seen: dict[int, int] = {}  # GPU-hours value -> index of an earlier job
    for j, h in enumerate(hours):
        if quota - h in seen:
            return [seen[quota - h], j]
        seen[h] = j
    return []
```

**Complexity:** O(n) time, O(n) space.

**Equivalent known problem:** Two Sum.

**Trap:** Sort + two pointers loses original indices and costs O(n log n). Inserting before checking pairs a job with itself when `h == quota / 2`.

</details>

### Drill 2: The missing shard

A storage cluster should contain shards numbered `0..n`. The inventory scan returns `n` distinct shard IDs in arbitrary order, so exactly one is missing. The scan runs in a sidecar with a hard memory cap, so use O(1) extra memory and do not reorder the input (it is shared with another reader). Return the missing ID.

- `ids = [3, 0, 1]` → `2`
- `ids = [0, 1, 2, 4, 5, 6, 3, 8]` → `7`
- Constraints: `1 ≤ n ≤ 10^4`.

<details><summary>Hints</summary>

1. "O(1) memory" and "do not reorder" rule out the set and the sort.
2. Numbers `0..n` with one absent; find it with a constant number of accumulators.
3. XOR cancellation (or the sum formula).

</details>

<details><summary>Solution</summary>

**Decode:** One value missing from a known range, plus the memory cap, signals cancellation: XOR every index and value; matched pairs cancel to 0.

**Abstract problem:** Missing number in `0..n`.

**Pattern:** Bit manipulation (`20-bit-manipulation-and-math.md`).

**Approach:** Seed `x = n`, then XOR in `i ^ ids[i]` for every `i`. Each present ID appears twice (as a value and as an index or `n`), so only the missing one survives.

```python
def missing_shard(ids: list[int]) -> int:
    x = len(ids)  # IDs are 0..n; n has no matching index, so seed with it
    for i, v in enumerate(ids):
        x ^= i ^ v  # every present ID cancels with its index; the missing one survives
    return x
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Missing Number.

**Trap:** Sorting violates "do not reorder" and costs O(n log n); a set uses O(n) memory. The sum formula `n(n+1)/2 - sum` works in Python but overflows fixed-width integers (a likely follow-up); XOR can't.

</details>

### Drill 3: The first regressed build

Nightly builds are numbered `1..n`. The latest build `n` fails the p99 latency benchmark. Once a regression lands, every later build also has it. Running the benchmark on one build takes 40 minutes on a dedicated GPU node, so minimize benchmark runs. You get `is_bad(build) -> bool`. Return the first bad build.

- `n = 10`, builds `7..10` are bad → `7`
- Constraints: `1 ≤ n ≤ 2^31 - 1`.

<details><summary>Hints</summary>

1. Latency and the GPU node are irrelevant. What shape does `is_bad` have over `1..n`?
2. A monotone boolean sequence `F F F ... T T T`; find the first `T`.
3. Boundary binary search.

</details>

<details><summary>Solution</summary>

**Decode:** "Every later build also has it": the predicate is monotone. "Minimize runs" with n up to 2^31: O(log n) calls.

**Abstract problem:** Find the first index where a monotone predicate becomes true.

**Pattern:** Binary search, boundary variant (`07-binary-search.md`).

**Approach:** Invariant: the answer is in `[lo, hi]`. Bad `mid` → `hi = mid` (it may be first); else `lo = mid + 1`.

```python
def first_bad_build(n: int, is_bad: Callable[[int], bool]) -> int:
    lo, hi = 1, n  # invariant: first bad build is in [lo, hi]; build n is known bad
    while lo < hi:
        mid = (lo + hi) // 2
        if is_bad(mid):
            hi = mid  # mid may be the first bad one
        else:
            lo = mid + 1
    return lo
```

**Complexity:** O(log n) predicate calls, O(1) space.

**Equivalent known problem:** First Bad Version.

**Trap:** Linear scan costs up to n benchmark runs. Mixing `while lo <= hi` with `hi = mid` loops forever.

</details>

### Drill 4: Minutes until a higher load

The autoscaler logs one CPU load reading per minute for a node. For each minute, it wants to know how many minutes it must wait until a strictly higher reading occurs, so it can decide whether a spike is "the peak". If no higher reading ever comes, report 0.

- `load = [40, 42, 38, 35, 41, 50, 45]` → `[1, 4, 2, 1, 1, 0, 0]`
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. For each position, find the next position to the right with a larger value.
2. "Next greater element" distance for every index.
3. Monotonic stack.

</details>

<details><summary>Solution</summary>

**Decode:** "Until a strictly higher reading" for every element is next-greater-element; `n ≤ 10^5` rules out the O(n^2) scan.

**Abstract problem:** For each `i`, the distance to the nearest `j > i` with `a[j] > a[i]`.

**Pattern:** Monotonic stack (`08-stacks-and-monotonic-structures.md`).

**Approach:** Stack of indices awaiting a higher value (loads non-increasing bottom to top). Each new reading resolves every waiting index with a smaller load.

```python
def minutes_until_higher(load: list[int]) -> list[int]:
    ans = [0] * len(load)
    stack: list[int] = []  # indices still waiting for a higher reading; loads non-increasing
    for i, x in enumerate(load):
        while stack and load[stack[-1]] < x:
            j = stack.pop()
            ans[j] = i - j
        stack.append(i)
    return ans
```

**Complexity:** O(n) time (each index pushed and popped once), O(n) space.

**Equivalent known problem:** Daily Temperatures.

**Trap:** Scanning right from each index looks fine on random data but is O(n^2) on a decreasing series.

</details>

### Drill 5: Failure incidents on the datacenter floor

The floor plan is a grid of strings: `'X'` is a rack in a failed state, `'.'` is healthy. Failed racks that share a wall (up, down, left, right) sit on the same power bus and count as one incident. Diagonal racks do not share a bus. How many incidents should on-call open?

```
XX..X
X...X      -> 4 incidents
..X..
....X
```

- Constraints: grid up to `300 × 300`.

<details><summary>Hints</summary>

1. Incidents are groups of failed cells glued by edges.
2. Count connected components of `'X'` cells under 4-neighborhood.
3. Grid DFS/BFS flood fill.

</details>

<details><summary>Solution</summary>

**Decode:** "Share a wall ... count as one" defines connectivity; "how many" counts connected components.

**Abstract problem:** Number of islands.

**Pattern:** Grid traversal (`13-graphs-traversal-and-topological-sort.md`).

**Approach:** Scan cells; on an unvisited `'X'`, count it and flood-fill its component with an explicit stack.

```python
def count_failure_zones(floor: list[str]) -> int:
    rows = len(floor)
    cols = len(floor[0]) if rows else 0
    seen: set[tuple[int, int]] = set()
    zones = 0
    for r in range(rows):
        for c in range(cols):
            if floor[r][c] != "X" or (r, c) in seen:
                continue
            zones += 1  # new unvisited failed rack: flood its whole zone
            seen.add((r, c))
            stack = [(r, c)]
            while stack:
                y, x = stack.pop()
                for ny, nx in ((y + 1, x), (y - 1, x), (y, x + 1), (y, x - 1)):
                    if (0 <= ny < rows and 0 <= nx < cols
                            and floor[ny][nx] == "X" and (ny, nx) not in seen):
                        seen.add((ny, nx))
                        stack.append((ny, nx))
    return zones
```

**Complexity:** O(R·C) time and space.

**Equivalent known problem:** Number of Islands.

**Trap:** Recursive DFS on a 300 × 300 all-`X` grid exceeds Python's default recursion limit. Including diagonals over-merges incidents.

</details>

### Drill 6: Two pods per node

A licensing limit allows at most two pods per node. Every node has memory capacity `limit`, and each pod's memory request is at most `limit`. Pods have names, owners, and priority classes. Find the minimum number of nodes that can host all pods.

- `mem = [3, 5, 3, 4], limit = 8` → `2` (3+5, 3+4)
- `mem = [3, 2, 2, 1], limit = 3` → `3`
- Constraints: `n ≤ 5·10^4`.

<details><summary>Hints</summary>

1. "At most two per node" turns general bin packing into a pairing problem.
2. Pair up numbers so each pair sums to at most `limit`; minimize pairs plus singles.
3. Sort, then two pointers from both ends (greedy).

</details>

<details><summary>Solution</summary>

**Decode:** Bin packing is NP-hard, but "at most two" makes each node a pair. Names and priorities don't matter.

**Abstract problem:** Minimum boats with capacity `limit` and at most two people per boat.

**Pattern:** Two pointers + greedy (`05-two-pointers.md`, `18-greedy.md`).

**Approach:** Sort. The heaviest pod needs a new node, and the lightest is its best partner: it rides along if it fits. If it doesn't, nothing else does.

```python
def min_nodes(mem: list[int], limit: int) -> int:
    mem = sorted(mem)
    i, j = 0, len(mem) - 1
    nodes = 0
    while i <= j:
        if mem[i] + mem[j] <= limit:  # lightest pod can share with the heaviest
            i += 1
        j -= 1  # heaviest remaining pod always gets this node
        nodes += 1
    return nodes
```

**Complexity:** O(n log n) for the sort, O(n) space for the sorted copy.

**Equivalent known problem:** Boats to Save People.

**Trap:** Reaching for first-fit-decreasing or DP because it "looks like bin packing". Pairing the two lightest first wastes capacity heavy pods need.

</details>

### Drill 7: Health-check streak with overrides

A service reports one health check per minute: `1` = pass, `0` = fail. For an SLA report, the SRE lead may mark at most `k` failed checks as "known flake" overrides. What is the longest run of consecutive minutes that can be reported as fully healthy?

- `checks = [1, 1, 0, 1, 0, 0, 1, 1, 1, 0], k = 2` → `6` (indices 3..8)
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. "Consecutive minutes" and "at most k" are the two signals.
2. Longest subarray containing at most `k` zeros.
3. Variable-size sliding window with a budget.

</details>

<details><summary>Solution</summary>

**Decode:** "Consecutive" = contiguous subarray; "at most k overrides" = a budget tracked incrementally. Shrinking a valid window keeps it valid (monotone).

**Abstract problem:** Longest subarray with at most `k` zeros.

**Pattern:** Sliding window, variable size (`06-sliding-window.md`).

**Approach:** Extend `right`; while failures exceed `k`, advance `left`. Each valid window is a candidate.

```python
def longest_healthy_streak(checks: list[int], k: int) -> int:
    best = left = fails = 0
    for right, ok in enumerate(checks):
        fails += ok == 0
        while fails > k:  # window invalid: too many failures to override
            fails -= checks[left] == 0
            left += 1
        best = max(best, right - left + 1)
    return best
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Max Consecutive Ones III.

**Trap:** Greedily spending overrides on the first `k` failures. A DP over "overrides used" is an unnecessary O(nk).

</details>

### Drill 8: Net replica change of exactly k

The autoscaler writes one entry per minute: the net replica change for that minute (positive for scale-up, negative for scale-down, zero for none). Capacity planning wants the number of contiguous time spans whose net change is exactly `k`.

- `deltas = [2, -1, 1, 2, -2, 2], k = 2` → `8`
- `deltas = [0, 0, 0], k = 0` → `6`
- Constraints: `n ≤ 2·10^4`, values in `[-1000, 1000]`.

<details><summary>Hints</summary>

1. Values can be negative. What does that do to a sliding window?
2. Count subarrays with sum exactly `k`.
3. Prefix sums + hash map of counts.

</details>

<details><summary>Solution</summary>

**Decode:** "Contiguous spans" + "count" + "exactly k". Negatives break the window (extending can lower the sum). Sum of `a[i..j]` = `P[j+1] - P[i]`, so count earlier prefixes equal to `P - k`.

**Abstract problem:** Number of subarrays with sum `k`.

**Pattern:** Prefix sums + hashing (`04-arrays-hashing-prefix-sums.md`).

**Approach:** Before recording each prefix, add the count of earlier prefixes equal to `prefix - k`. Seed `{0: 1}` for spans starting at index 0.

```python
def count_spans_with_net(deltas: list[int], k: int) -> int:
    count = prefix = 0
    seen = Counter({0: 1})  # prefix-sum value -> how many times seen (empty prefix counts)
    for d in deltas:
        prefix += d
        count += seen[prefix - k]  # earlier prefix p with prefix - p == k
        seen[prefix] += 1
    return count
```

**Complexity:** O(n) time, O(n) space.

**Equivalent known problem:** Subarray Sum Equals K.

**Trap:** Sliding window, since with negatives shrinking may not reduce the sum. Without the `{0: 1}` seed, spans starting at minute 0 are missed.

</details>

### Drill 9: Nearest common aggregation switch

The network is a binary tree of switches; each switch has at most two downlinks and a unique ID. IDs are assigned by procurement order, not by position. Two racks (identified by their switch IDs, both present) exchange heavy traffic. Find the lowest switch that all of their traffic must pass through. A switch counts as on its own path.

```
            1
          /   \
         2     3
        / \   / \
       4   5 6   7
          / \
         8   9
```

- `(8, 4)` → `2`; `(9, 7)` → `1`; `(5, 9)` → `5`
- Constraints: up to `10^5` switches.

<details><summary>Hints</summary>

1. "IDs not by position": assume no ordering.
2. Lowest node that has both targets in its subtree.
3. Lowest common ancestor in a binary tree (not a BST).

</details>

<details><summary>Solution</summary>

**Decode:** The lowest switch on the path between two tree nodes is their lowest common ancestor. The procurement-order remark rules out BST comparisons.

**Abstract problem:** LCA of two nodes in a binary tree.

**Pattern:** Trees, post-order recursion (`10-trees-and-bst.md`).

**Approach:** Return the node if it's a target, else recurse. If both subtrees return something, the targets split here; otherwise pass up whichever side found one.

```python
class TreeNode:
    def __init__(self, val: int, left: "TreeNode | None" = None,
                 right: "TreeNode | None" = None):
        self.val, self.left, self.right = val, left, right


def nearest_common_switch(root: Optional[TreeNode], a: int, b: int) -> Optional[TreeNode]:
    if root is None or root.val in (a, b):
        return root
    left = nearest_common_switch(root.left, a, b)
    right = nearest_common_switch(root.right, a, b)
    if left and right:  # a and b found in different subtrees: root is the split point
        return root
    return left or right
```

**Complexity:** O(n) time, O(h) recursion stack.

**Equivalent known problem:** Lowest Common Ancestor of a Binary Tree.

**Trap:** BST logic (`a < root.val` goes left) because IDs are integers. A 10^5-deep degenerate tree overflows Python recursion; offer an iterative parent-map version as a follow-up.

</details>

### Drill 10: How many CI runners?

CI jobs have fixed `[start, end)` times in minutes. A runner runs one job at a time and can start a new job at the same minute the previous one ends. Jobs belong to different repos and have different priorities. What is the minimum number of runners so that no job waits?

- `[(0, 30), (5, 10), (15, 20), (10, 25)]` → `3` (at minute 15: jobs 0, 2, 3)
- `[(1, 5), (5, 9)]` → `1`
- Constraints: `n ≤ 10^4`.

<details><summary>Hints</summary>

1. Repos and priorities are noise. When are the most jobs running at once?
2. Maximum number of overlapping half-open intervals.
3. Sweep line over start/end events (or a min-heap of end times).

</details>

<details><summary>Solution</summary>

**Decode:** Fixed start/end times + "minimum resources so none waits" = peak concurrency. Half-open: an end and a start at the same minute don't overlap.

**Abstract problem:** Minimum meeting rooms.

**Pattern:** Intervals, sweep line (`12-intervals-and-sweep-line.md`).

**Approach:** Sort events `(start, +1)` and `(end, -1)`; at equal times `-1` sorts first, encoding the half-open rule. Answer: the running sum's maximum.

```python
def min_runners(jobs: list[tuple[int, int]]) -> int:
    events = []
    for s, e in jobs:
        events.append((s, 1))
        events.append((e, -1))
    events.sort()  # at equal times -1 sorts first: a runner freed at t is reusable at t
    cur = best = 0
    for _, delta in events:
        cur += delta
        best = max(best, cur)
    return best
```

**Complexity:** O(n log n) time, O(n) space.

**Equivalent known problem:** Meeting Rooms II.

**Trap:** Merging intervals counts connected blocks, not peak concurrency. Sorting `+1` before `-1` at ties overcounts back-to-back jobs.

</details>

### Drill 11: Malware spreading across racks

A grid describes rack slots: `0` = empty slot, `1` = clean host, `2` = compromised host. Every minute, each compromised host infects every clean host that shares a side with it. How many minutes until no clean host remains? Return `-1` if some clean host can never be reached.

- `[[2,1,1],[1,1,0],[0,1,1]]` → `4`
- `[[2,1,1],[0,1,1],[1,0,1]]` → `-1` (bottom-left host is isolated)
- `[[0,2]]` → `0`
- Constraints: grid up to `100 × 100`.

<details><summary>Hints</summary>

1. Several hosts start infected and spread simultaneously.
2. Shortest distance from the nearest source to every cell; answer is the maximum.
3. Multi-source BFS.

</details>

<details><summary>Solution</summary>

**Decode:** "Every minute ... spreads to neighbors" = BFS layers measure time; simultaneous sources = seed the queue with all of them.

**Abstract problem:** Max over cells of distance to the nearest source; -1 if unreachable.

**Pattern:** Multi-source BFS (`13-graphs-traversal-and-topological-sort.md`).

**Approach:** Enqueue all `2`s and count `1`s. Each BFS layer is a minute; stop when no clean hosts remain.

```python
def minutes_to_full_compromise(racks: list[list[int]]) -> int:
    grid = [row[:] for row in racks]
    rows, cols = len(grid), len(grid[0]) if grid else 0
    q: deque[tuple[int, int]] = deque()
    clean = 0
    for r in range(rows):
        for c in range(cols):
            if grid[r][c] == 2:
                q.append((r, c))  # all infected hosts start in the queue together
            elif grid[r][c] == 1:
                clean += 1
    minutes = 0
    while q and clean:
        for _ in range(len(q)):  # one BFS layer == one minute
            r, c = q.popleft()
            for nr, nc in ((r + 1, c), (r - 1, c), (r, c + 1), (r, c - 1)):
                if 0 <= nr < rows and 0 <= nc < cols and grid[nr][nc] == 1:
                    grid[nr][nc] = 2
                    clean -= 1
                    q.append((nr, nc))
        minutes += 1
    return minutes if clean == 0 else -1
```

**Complexity:** O(R·C) time and space.

**Equivalent known problem:** Rotting Oranges.

**Trap:** A separate BFS per infected host costs O(k·R·C). Counting a minute for the final layer that infects nothing is off by one.

</details>

### Drill 12: Terraform apply order

Infrastructure modules declare `depends_on` lists. A module can be applied only after all of its dependencies. Produce an apply order. When several modules are ready at once, take them in the order they appear in the `modules` list. If the dependencies contain a loop, return `[]`.

- `modules = ["vpc", "subnets", "eks", "iam", "addons"]`
- `depends_on = {"subnets": ["vpc"], "eks": ["subnets", "iam"], "addons": ["eks"]}`
- → `["vpc", "iam", "subnets", "eks", "addons"]`
- Constraints: up to `2000` modules, `5000` edges.

<details><summary>Hints</summary>

1. "Must happen before" relations plus "loop means impossible".
2. Topological order of a directed graph; detect cycles.
3. Kahn's algorithm (BFS on in-degrees).

</details>

<details><summary>Solution</summary>

**Decode:** "Dependencies first" and "empty on a loop" are the two halves of topological sort: an order exists exactly when the graph is a DAG.

**Abstract problem:** Course Schedule II with deterministic tie-breaking.

**Pattern:** Topological sort (`13-graphs-traversal-and-topological-sort.md`).

**Approach:** Edges go `dependency → module`. Start with zero in-degree modules in list order; pop, emit, decrement its users. Emitting fewer than all means a cycle blocked the rest.

```python
def apply_order(modules: list[str], depends_on: dict[str, list[str]]) -> list[str]:
    indeg = {m: 0 for m in modules}
    users: dict[str, list[str]] = defaultdict(list)  # dependency -> modules needing it
    for m, deps in depends_on.items():
        for d in deps:
            users[d].append(m)
            indeg[m] += 1
    q = deque(m for m in modules if indeg[m] == 0)
    order = []
    while q:
        m = q.popleft()
        order.append(m)
        for u in users[m]:
            indeg[u] -= 1
            if indeg[u] == 0:
                q.append(u)
    return order if len(order) == len(modules) else []  # leftovers mean a cycle
```

**Complexity:** O(V + E) time and space.

**Equivalent known problem:** Course Schedule II.

**Trap:** DFS with only a visited set can't tell a back edge (cycle) from a cross edge; use three colors. Ordering by dependency count isn't a valid topological order.

</details>

### Drill 13: Merging per-node log files

Each of `k` nodes writes its own log file, already sorted by timestamp. Lines are `(timestamp, message)`. Produce a single merged timeline. On equal timestamps, lines from the lower-numbered file come first.

- files `[(1,a1),(4,a4),(9,a9)]`, `[(2,b2),(3,b3)]`, `[(5,c5)]` → timestamps `1, 2, 3, 4, 5, 9`
- Constraints: `k ≤ 10^4`, total lines `N ≤ 10^6`.

<details><summary>Hints</summary>

1. Inputs are sorted; you only need the smallest current head.
2. Merge `k` sorted lists.
3. Min-heap of list heads.

</details>

<details><summary>Solution</summary>

**Decode:** Sorted sources + one sorted output = k-way merge: the next line is the minimum of k heads, a heap gives it in O(log k).

**Abstract problem:** Merge k sorted lists.

**Pattern:** Heaps, k-way merge (`11-heaps-top-k-and-merging.md`).

**Approach:** Push `(timestamp, file_index, position)` per non-empty file. Pop, emit, push that file's next line. The file index breaks ties, so messages are never compared.

```python
def merge_log_streams(streams: list[list[tuple[int, str]]]) -> list[tuple[int, str]]:
    heap = [(s[0][0], i, 0) for i, s in enumerate(streams) if s]  # (ts, stream, pos)
    heapq.heapify(heap)
    out = []
    while heap:
        _, i, j = heapq.heappop(heap)
        out.append(streams[i][j])
        if j + 1 < len(streams[i]):
            heapq.heappush(heap, (streams[i][j + 1][0], i, j + 1))
    return out
```

**Complexity:** O(N log k) time, O(k) extra space besides the output.

**Equivalent known problem:** Merge k Sorted Lists.

**Trap:** Concatenate-and-sort is O(N log N), wastes the sortedness, and can't stream. Merging files one by one into an accumulator is O(N·k).

</details>

### Drill 14: Pager code interpretations

A legacy pager encodes alert names by mapping `A=1 ... Z=26` and concatenating digits with no separators. Given a received digit string, how many different letter strings could it represent?

- `"226"` → `3` (`BZ`, `VF`, `BBF`)
- `"12"` → `2`; `"06"` → `0`; `"10"` → `1`
- Constraints: length `≤ 100` (answers fit in a 64-bit int).

<details><summary>Hints</summary>

1. "How many" with length 100: enumeration is hopeless.
2. Ways to split a digit string into pieces from `1..26`, no leading zeros.
3. 1D DP over prefixes.

</details>

<details><summary>Solution</summary>

**Decode:** "How many ... could it represent" = count; the last piece is one or two digits and prefixes overlap, so DP.

**Abstract problem:** Decode ways.

**Pattern:** DP 1D (`16-dynamic-programming-foundations.md`).

**Approach:** `ways(i) = ways(i-1) [if s[i-1] != '0'] + ways(i-2) [if s[i-2:i] in 10..26]`. Keep two rolling variables.

```python
def count_alert_decodings(code: str) -> int:
    prev, cur = 1, (1 if code[0] != "0" else 0)  # ways(code[:0]), ways(code[:1])
    for i in range(2, len(code) + 1):
        nxt = 0
        if code[i - 1] != "0":
            nxt += cur  # last digit decoded alone
        if 10 <= int(code[i - 2:i]) <= 26:
            nxt += prev  # last two digits decoded together
        prev, cur = cur, nxt
    return cur
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Decode Ways.

**Trap:** Backtracking is exponential (`"111...1"` has Fibonacci-many decodings). Treating `"0"` as decodable, or allowing `"05"`.

</details>

### Drill 15: The redundant cable

A rack of `n` switches (IDs `1..n`) was cabled as a tree: connected, no loops. Later a technician added one extra cable, creating a loop that causes broadcast storms. The `n` cables are listed in installation order. Return the cable to remove so the network is a tree again; if several work, return the one installed last.

- `[(1,2), (1,3), (2,3)]` → `[2, 3]`
- `[(1,2), (2,3), (3,4), (1,4), (1,5)]` → `[1, 4]`
- Constraints: `3 ≤ n ≤ 1000`.

<details><summary>Hints</summary>

1. Process cables in order. When does a cable stop being useful?
2. The first edge whose endpoints are already connected closes the cycle.
3. Union-find.

</details>

<details><summary>Solution</summary>

**Decode:** "Already connected?" asked as edges arrive = union-find. The first edge joining two nodes of one component is the cycle's last-installed edge.

**Abstract problem:** Redundant connection in an undirected graph.

**Pattern:** Union-find (`14-union-find-shortest-paths-mst.md`).

**Approach:** `find` both ends of each cable: a shared root means it closes the loop; otherwise union.

```python
def redundant_cable(cables: list[tuple[int, int]]) -> list[int]:
    parent = list(range(len(cables) + 1))  # switches are 1..n and n == len(cables)

    def find(x: int) -> int:
        while parent[x] != x:
            parent[x] = parent[parent[x]]  # path halving
            x = parent[x]
        return x

    for a, b in cables:
        ra, rb = find(a), find(b)
        if ra == rb:  # already connected: this cable closes a loop
            return [a, b]
        parent[ra] = rb
    return []
```

**Complexity:** O(n α(n)) time, O(n) space.

**Equivalent known problem:** Redundant Connection.

**Trap:** A DFS reachability check per edge is O(n^2). Returning any cycle edge ignores the "installed last" tie-break.

</details>

### Drill 16: Config push reaching every node

A config change starts at node `source` in a gossip network of `n` nodes (IDs `1..n`). Directed link `(u, v, ms)` means `u` forwards to `v` after `ms` milliseconds. Nodes forward immediately on receipt. How long until every node has the change? Return `-1` if some node never gets it.

- `n = 4, links = [(2,1,1), (2,3,1), (3,4,1)], source = 2` → `2`
- `n = 3, links = [(1,2,10), (1,3,1), (3,2,2)], source = 1` → `3`
- Constraints: `n ≤ 100`, `links ≤ 6000`, `ms ≥ 0`.

<details><summary>Hints</summary>

1. Each node receives the change at its earliest arrival time.
2. Single-source shortest paths with non-negative weights; take the maximum.
3. Dijkstra.

</details>

<details><summary>Solution</summary>

**Decode:** "Until every node" = slowest node's earliest arrival = max shortest-path distance. Non-negative directed weights = Dijkstra.

**Abstract problem:** Network delay time.

**Pattern:** Dijkstra (`14-union-find-shortest-paths-mst.md`).

**Approach:** Lazy Dijkstra over `(dist, node)`, skipping stale entries. Return -1 if fewer than `n` nodes are finalized.

```python
def full_propagation_time(n: int, links: list[tuple[int, int, int]], source: int) -> int:
    graph: dict[int, list[tuple[int, int]]] = defaultdict(list)
    for u, v, w in links:
        graph[u].append((v, w))
    dist: dict[int, int] = {}
    heap = [(0, source)]
    while heap:
        d, u = heapq.heappop(heap)
        if u in dist:
            continue  # stale entry; u was finalized with a smaller distance
        dist[u] = d
        for v, w in graph[u]:
            if v not in dist:
                heapq.heappush(heap, (d + w, v))
    return max(dist.values()) if len(dist) == n else -1
```

**Complexity:** O(E log E) time, O(V + E) space.

**Equivalent known problem:** Network Delay Time.

**Trap:** BFS by hop count ignores latencies (second example: 2 hops beat 1). Summing the shortest-path tree's edges answers a different question.

</details>

### Drill 17: Instance mixes that hit an exact vCPU count

A cloud provider offers instance types with given vCPU counts, unlimited supply of each. The price per hour of each type is listed too. A capacity planner wants every distinct mix of instances whose vCPUs add up to exactly `target` (order inside a mix does not matter). List each mix sorted ascending; list mixes in lexicographic order.

- `sizes = [2, 3, 6, 7], target = 7` → `[[2, 2, 3], [7]]`
- `sizes = [2, 3, 5], target = 8` → `[[2, 2, 2, 2], [2, 3, 3], [3, 5]]`
- Constraints: `≤ 30` sizes, each `≥ 2`, `target ≤ 40`.

<details><summary>Hints</summary>

1. Prices are noise. The output is every solution, not a count.
2. All multisets from a set of numbers summing to `target`, with reuse.
3. Backtracking with a start index.

</details>

<details><summary>Solution</summary>

**Decode:** "Every distinct mix" = enumerate all = backtracking. "Unlimited supply" = reuse. "Order does not matter" = pick in non-decreasing order to avoid permutations.

**Abstract problem:** Combination Sum.

**Pattern:** Backtracking (`15-backtracking.md`).

**Approach:** Sort unique sizes. `backtrack(start, remaining)` tries sizes from `start`, recursing on the same index for reuse and breaking once a size exceeds what remains.

```python
def instance_combos(sizes: list[int], target: int) -> list[list[int]]:
    sizes = sorted(set(sizes))
    out: list[list[int]] = []
    combo: list[int] = []

    def backtrack(start: int, remaining: int) -> None:
        if remaining == 0:
            out.append(combo[:])
            return
        for i in range(start, len(sizes)):  # only sizes >= last pick: no permuted duplicates
            if sizes[i] > remaining:
                break  # sorted, so every later size is too big as well
            combo.append(sizes[i])
            backtrack(i, remaining - sizes[i])  # i, not i + 1: sizes are reusable
            combo.pop()

    backtrack(0, target)
    return out
```

**Complexity:** Exponential worst case, bounded by output size times depth (`target / min(size)`); fine here.

**Equivalent known problem:** Combination Sum.

**Trap:** Coin-change DP counts mixes but can't cheaply list them. Restarting at index 0 each level yields `[2,2,3]`, `[2,3,2]`, `[3,2,2]`.

</details>

### Drill 18: Follow-the-sun credit loop

A batch job hops around a ring of `n` regions to follow cheap off-peak power. In region `i` it earns `earn[i]` compute credits; moving from region `i` to `i+1` (wrapping around) costs `cost[i]` credits. It starts with zero credits and cannot go negative. From which region can it complete the full loop? The answer is unique if it exists; return `-1` otherwise.

- `earn = [1, 2, 3, 4, 5], cost = [3, 4, 5, 1, 2]` → `3`
- `earn = [2, 3, 4], cost = [3, 4, 3]` → `-1`
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. If total earn < total cost, no start works. Otherwise?
2. If the running balance from start `s` goes negative at `i`, no start in `[s, i]` works.
3. Greedy single pass (Gas Station).

</details>

<details><summary>Solution</summary>

**Decode:** Circular route, gain and cost per stop, "from which start", `n ≤ 10^5`: simulating every start (O(n^2)) is out. The skip-ahead lemma allows one pass.

**Abstract problem:** Gas station.

**Pattern:** Greedy (`18-greedy.md`).

**Approach:** If totals fail, return -1. Otherwise scan with a running tank; when it goes negative, the next index becomes the candidate (any start inside the failed stretch reaches the failure point with even less).

```python
def start_region(earn: list[int], cost: list[int]) -> int:
    if sum(earn) < sum(cost):
        return -1
    start = tank = 0
    for i in range(len(earn)):
        tank += earn[i] - cost[i]
        if tank < 0:  # no start in [start, i] can pass i; try i + 1
            start, tank = i + 1, 0
    return start
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Gas Station.

**Trap:** Trying every start is O(n^2). Skipping the total check returns a start even when no loop is possible.

</details>

### Drill 19: CLI autocomplete

You are building autocomplete for a kubectl-like CLI. Given the list of resource names and the string the user types one character at a time, after each keystroke show up to 3 resource names that start with the typed prefix, in lexicographic order.

- `names = ["deployments", "daemonsets", "deploy-keys", "configmaps", "deprecated", "devices"], typed = "dep"`
- → `[["daemonsets", "deploy-keys", "deployments"], ["deploy-keys", "deployments", "deprecated"], ["deploy-keys", "deployments", "deprecated"]]`
- Constraints: up to `1000` names, total characters `≤ 2·10^4`.

<details><summary>Hints</summary>

1. Every query is "names with this prefix", each extending the previous.
2. For each prefix of a query string, the 3 smallest words with that prefix.
3. Trie storing the top 3 at each node (or sort + binary search).

</details>

<details><summary>Solution</summary>

**Decode:** "Prefix typed so far" is the canonical trie signal; the query walks one trie path a character at a time.

**Abstract problem:** Search suggestions system.

**Pattern:** Trie (`19-tries-and-string-algorithms.md`).

**Approach:** Insert names in sorted order so each node's first 3 names are its smallest. Walk the typed string; once a character is missing, all later answers are empty.

```python
class _TrieNode:
    __slots__ = ("children", "top")

    def __init__(self) -> None:
        self.children: dict[str, _TrieNode] = {}
        self.top: list[str] = []  # up to 3 smallest names passing through this node


def autocomplete(names: list[str], typed: str) -> list[list[str]]:
    root = _TrieNode()
    for name in sorted(set(names)):  # sorted insert: first 3 seen are the smallest
        node = root
        for ch in name:
            node = node.children.setdefault(ch, _TrieNode())
            if len(node.top) < 3:
                node.top.append(name)
    out: list[list[str]] = []
    node: Optional[_TrieNode] = root
    for ch in typed:
        node = node.children.get(ch) if node else None
        out.append(node.top[:] if node else [])
    return out
```

**Complexity:** O(total chars · log n) to sort plus O(total chars) to build; O(len(typed)) per query walk.

**Equivalent known problem:** Search Suggestions System.

**Trap:** Rescanning all names per keystroke works at this size but doesn't scale; say so and offer the trie. Looking up children after the path breaks fails on `None`.

</details>

### Drill 20: The duplicated ticket

A ticketing service issued `n + 1` tickets, each labeled with an integer in `1..n`, so some label repeats. Exactly one label is duplicated, but it may appear more than twice. The ticket array lives in a read-only shared memory segment, and you may use only O(1) extra memory. Return the duplicated label.

- `[1, 3, 4, 2, 2]` → `2`
- `[3, 1, 3, 4, 2]` → `3`
- `[2, 2, 2, 2, 2]` → `2`
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. Read-only + O(1) memory rule out sets, sorting, and sign-marking.
2. Values are valid indices. Treat `i → tickets[i]` as a "next pointer".
3. Floyd cycle detection (fast/slow pointers).

</details>

<details><summary>Solution</summary>

**Decode:** Values in `1..n` across `n + 1` slots work as pointers. Two indices point at the duplicate, so following them from index 0 hits a cycle entering at the duplicate. "O(1) memory, no modification" signals linked-list cycle detection.

**Abstract problem:** Find the duplicate number (cycle entry of `i → a[i]`).

**Pattern:** Fast/slow pointers on an implicit linked list (`09-linked-lists.md`).

**Approach:** Phase 1: slow moves one step, fast two, until they meet in the cycle. Phase 2: restart slow at the head, move both one step; they meet at the cycle entry.

```python
def duplicate_ticket(tickets: list[int]) -> int:
    # Treat i -> tickets[i] as a linked list; the duplicate value is the cycle entry.
    slow = fast = tickets[0]
    while True:
        slow = tickets[slow]
        fast = tickets[tickets[fast]]
        if slow == fast:
            break
    slow = tickets[0]
    while slow != fast:
        slow = tickets[slow]
        fast = tickets[fast]
    return slow
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Find the Duplicate Number.

**Trap:** XOR works only if the duplicate appears exactly twice and nothing else is missing. Binary search on the value range (count values `≤ mid`) is an acceptable O(n log n) fallback.

</details>

### Drill 21: Migration bandwidth cap

Database shards must be migrated in fixed key-range order; sizes are in GB. Each day the migration proceeds in order until the next shard would exceed the daily bandwidth cap; shards cannot be split across days. What is the smallest daily cap that finishes all shards within `D` days?

- `shards = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10], D = 5` → `15`
- `shards = [3, 2, 2, 4, 1, 4], D = 3` → `6`
- Constraints: `n ≤ 5·10^4`, sizes `≤ 500`.

<details><summary>Hints</summary>

1. If a cap works, any bigger cap also works.
2. Smallest `cap` for which greedy packing into consecutive groups uses `≤ D` groups; equivalently, split into `D` contiguous parts minimizing the largest sum.
3. Binary search on the answer.

</details>

<details><summary>Solution</summary>

**Decode:** "Smallest cap such that it finishes within D" = minimize a parameter under a monotone feasibility check, and checking a cap is a greedy pass. The answer lies in `[max(shards), sum(shards)]`.

**Abstract problem:** Capacity to ship packages within D days (minimize the maximum contiguous group sum).

**Pattern:** Binary search on the answer (`07-binary-search.md`).

**Approach:** Binary search the smallest cap whose greedy `days_needed(cap) ≤ D`.

```python
def min_daily_capacity(shards: list[int], days: int) -> int:
    def days_needed(cap: int) -> int:
        d, load = 1, 0
        for s in shards:
            if load + s > cap:  # shard does not fit today: start a new day
                d += 1
                load = 0
            load += s
        return d

    lo, hi = max(shards), sum(shards)  # feasible(cap) is monotone in cap
    while lo < hi:
        mid = (lo + hi) // 2
        if days_needed(mid) <= days:
            hi = mid
        else:
            lo = mid + 1
    return lo
```

**Complexity:** O(n log(sum)) time, O(1) extra space.

**Equivalent known problem:** Capacity To Ship Packages Within D Days.

**Trap:** `ceil(sum / D)` ignores that shards are indivisible and ordered. Sorting breaks the fixed order. A partition DP is O(n^2·D).

</details>

### Drill 22: Rebuilding IPs from stripped logs

A log pipeline bug stripped the dots from IPv4 addresses. Given one damaged digit string, list every valid IPv4 address it could have been. Each octet is `0..255` and has no leading zero unless it is exactly `"0"`. All digits must be used.

- `"25525511135"` → `["255.255.11.135", "255.255.111.35"]`
- `"0000"` → `["0.0.0.0"]`
- `"101023"` → `["1.0.10.23", "1.0.102.3", "10.1.0.23", "10.10.2.3", "101.0.2.3"]`
- Constraints: length `≤ 20`.

<details><summary>Hints</summary>

1. "List every valid" = enumerate.
2. Place 3 dots so all 4 pieces are valid octets.
3. Backtracking with pruning.

</details>

<details><summary>Solution</summary>

**Decode:** "Every ... it could have been" = all solutions. Each choice (next octet length 1, 2, or 3) constrains the rest, and invalid prefixes are cut early.

**Abstract problem:** Restore IP addresses.

**Pattern:** Backtracking (`15-backtracking.md`).

**Approach:** Choose 4 segments left to right, skipping leading zeros and values above 255. Prune when leftover digits can't fill the remaining octets (1 to 3 digits each).

```python
def restore_ips(digits: str) -> list[str]:
    out: list[str] = []
    parts: list[str] = []

    def backtrack(i: int) -> None:
        if len(parts) == 4:
            if i == len(digits):
                out.append(".".join(parts))
            return
        remaining_parts = 4 - len(parts)
        if not remaining_parts <= len(digits) - i <= 3 * remaining_parts:
            return  # prune: leftover digits cannot fill the remaining octets
        for length in (1, 2, 3):
            seg = digits[i:i + length]
            if len(seg) < length:
                break
            if (seg[0] == "0" and length > 1) or int(seg) > 255:
                continue
            parts.append(seg)
            backtrack(i + length)
            parts.pop()

    backtrack(0)
    return out
```

**Complexity:** At most 3^4 = 81 branches, so O(1); the length prune rejects long inputs at once.

**Equivalent known problem:** Restore IP Addresses.

**Trap:** Missing the leading-zero rule (`"01"`) or not requiring all digits be used. Comparing octets as strings (`"99" > "255"`).

</details>

### Drill 23: Splitting jobs across two GPU nodes

You have two identical GPU nodes and a batch of training jobs with integer durations. Each job runs entirely on one node. Can you assign jobs so both nodes finish at exactly the same time?

- `[1, 5, 11, 5]` → `True` (`11` vs `1 + 5 + 5`)
- `[1, 2, 3, 5]` → `False`
- Constraints: `n ≤ 200`, durations `≤ 100`.

<details><summary>Hints</summary>

1. "Same finish time": each node gets half the total.
2. Is there a subset that sums to `total / 2`?
3. 0/1 knapsack (subset sum DP).

</details>

<details><summary>Solution</summary>

**Decode:** "Each job on one node" = each item used once; "exactly the same" = subset sum to half. The sum is at most 20,000, so DP over sums is cheap.

**Abstract problem:** Partition equal subset sum.

**Pattern:** 0/1 knapsack (`16-dynamic-programming-foundations.md`).

**Approach:** Odd total → False. Otherwise fill a boolean array over sums `0..target`, iterating sums downward per job so each is used at most once.

```python
def can_split_evenly(durations: list[int]) -> bool:
    total = sum(durations)
    if total % 2:
        return False
    target = total // 2
    reachable = [True] + [False] * target  # reachable[s]: some subset sums to s
    for d in durations:
        for s in range(target, d - 1, -1):  # descending: each job used at most once
            reachable[s] = reachable[s] or reachable[s - d]
    return reachable[target]
```

**Complexity:** O(n · total/2) time, O(total/2) space.

**Equivalent known problem:** Partition Equal Subset Sum.

**Trap:** Greedy "next largest job to the less loaded node" fails: `[3, 3, 2, 2, 2]` ends 7 vs 5, yet `{3, 3}` vs `{2, 2, 2}` works. Iterating sums upward reuses a job.

</details>

### Drill 24: Kernel launch cooldown

A GPU executes one kernel launch per time slot. Launches of the same kernel type need at least `gap` slots between them to let a cache cool down; a slot may stay idle. You may reorder launches freely. What is the minimum number of slots to run them all?

- `["A","A","A","B","B","B"], gap = 2` → `8` (`A B _ A B _ A B`)
- `["A","A","A","A","A","A","B","C","D","E","F","G"], gap = 2` → `16`
- Constraints: up to `10^4` launches, 26 kernel types.

<details><summary>Hints</summary>

1. The most frequent kernel dictates the skeleton.
2. Arrange counts so identical items are `gap + 1` apart; minimize length.
3. Greedy counting formula (or a max-heap simulation).

</details>

<details><summary>Solution</summary>

**Decode:** "Same type needs cooldown", "reorder freely", "minimum slots" = task scheduling. The busiest type forces `peak - 1` frames of size `gap + 1`, plus a last frame of every type tied for the peak. If other launches fill all idles, the answer is the count.

**Abstract problem:** Task scheduler.

**Pattern:** Greedy (`18-greedy.md`); heap simulation in `11-heaps-top-k-and-merging.md`.

**Approach:** `max(len(kernels), (peak - 1)(gap + 1) + n_peak)`.

```python
def min_slots(kernels: list[str], gap: int) -> int:
    counts = list(Counter(kernels).values())
    peak = max(counts)
    n_peak = counts.count(peak)
    # (peak - 1) full frames of size gap + 1, plus a final partial frame of the peak kernels
    return max(len(kernels), (peak - 1) * (gap + 1) + n_peak)
```

**Complexity:** O(n) time, O(1) space (26 types).

**Equivalent known problem:** Task Scheduler.

**Trap:** Simulating in input order (order is free). Dropping the `max` with `len(kernels)`, needed when many distinct types overflow the frames.

</details>

### Drill 25: Version-tag rollout path

Release tags are fixed-length strings. Policy: each rollout step changes exactly one character of the deployed tag, and every intermediate tag must be on the security team's approved list. What is the minimum number of steps from `start` to `target`? Return `-1` if impossible.

- `start = "v1.2.3", target = "v2.3.5"`, `approved = ["v1.2.4", "v1.3.4", "v2.3.4", "v2.3.5", "v9.9.9"]` → `4`
- same list, `target = "v9.9.9"` → `-1`
- Constraints: up to `5000` approved tags, length `≤ 10`.

<details><summary>Hints</summary>

1. Each tag is a state; one step moves to a neighbor state.
2. Shortest path in an unweighted implicit graph.
3. BFS on states.

</details>

<details><summary>Solution</summary>

**Decode:** "Minimum number of steps", equal step costs, one change per step = BFS on an implicit graph; generate neighbors on demand.

**Abstract problem:** Minimum genetic mutation / word ladder.

**Pattern:** BFS on an implicit graph (`13-graphs-traversal-and-topological-sort.md`).

**Approach:** BFS from `start` over approved, unseen one-character edits (alphabet = characters in approved tags). Return the level at which `target` is generated.

```python
def min_rollout_steps(start: str, target: str, allowed: list[str]) -> int:
    allowed_set = set(allowed)
    if start == target:
        return 0
    if target not in allowed_set:
        return -1
    alphabet = set("".join(allowed_set))
    q = deque([(start, 0)])
    seen = {start}
    while q:
        tag, steps = q.popleft()
        for i in range(len(tag)):
            for ch in alphabet:
                nxt = tag[:i] + ch + tag[i + 1:]
                if nxt in allowed_set and nxt not in seen:
                    if nxt == target:
                        return steps + 1
                    seen.add(nxt)
                    q.append((nxt, steps + 1))
    return -1
```

**Complexity:** O(N · L · |Σ|) neighbor checks, each O(L) to build the string; O(N) space.

**Equivalent known problem:** Minimum Genetic Mutation (Word Ladder is the same idea).

**Trap:** DFS finds a path, not the shortest. Greedily fixing characters toward `target` fails when the intermediate isn't approved (here the path must pass `v1.3.4`).

</details>

### Drill 26: Spot reservations with cooldown

A spot GPU market quotes one price per day. You hold at most one reservation at a time: buy it at the day's price and later release it back at that day's price. After releasing, compliance requires one full day before buying again. Maximize total profit.

- `[1, 2, 3, 0, 2]` → `3` (buy 1, sell 2, cooldown, buy 0, sell 2)
- `[5, 4, 3]` → `0`
- Constraints: `n ≤ 5000`.

<details><summary>Hints</summary>

1. Each day you're in one of a few modes; the rules say which follows which.
2. DP over days with states `hold`, `just sold`, `resting`.
3. State-machine DP.

</details>

<details><summary>Solution</summary>

**Decode:** "At most one at a time" + "after selling you must wait" = a small automaton: buy only from `rest`, sell only from `hold`, and `sold` must pass through `rest`. Track best profit per state per day.

**Abstract problem:** Best time to buy and sell stock with cooldown.

**Pattern:** State-machine DP (`17-dynamic-programming-advanced.md`).

**Approach:** `hold' = max(hold, rest - p)`, `sold' = hold + p`, `rest' = max(rest, sold)`. Update simultaneously.

```python
def max_spot_profit(prices: list[int]) -> int:
    hold, sold, rest = -inf, 0, 0  # best profit if: holding / sold today / free to buy
    for p in prices:
        hold, sold, rest = max(hold, rest - p), hold + p, max(rest, sold)
    return int(max(sold, rest))
```

**Complexity:** O(n) time, O(1) space.

**Equivalent known problem:** Best Time to Buy and Sell Stock with Cooldown.

**Trap:** Summing every positive daily rise (the no-cooldown answer) gives 4 here. Sequential instead of simultaneous updates allow buying on the day you sold.

</details>

### Drill 27: Cheapest transfer with at most k relays

Data moves between `n` regions over directed links `(u, v, cost_per_TB)`. Compliance allows data to pass through at most `k` intermediate regions. Find the cheapest cost from `src` to `dst`, or `-1`.

- `links = [(0,1,100), (1,2,100), (2,3,100), (0,2,500), (1,3,600)]`, `src = 0, dst = 3`
- `k = 1` → `600` (0→2→3); `k = 2` → `300` (0→1→2→3); `k = 0` → `-1`
- Constraints: `n ≤ 100`, `k < n`.

<details><summary>Hints</summary>

1. The cheapest path overall may use too many hops.
2. Shortest path using at most `k + 1` edges.
3. Bellman-Ford limited to `k + 1` rounds (or BFS/Dijkstra over `(node, hops)` states).

</details>

<details><summary>Solution</summary>

**Decode:** "Cheapest" + "at most k intermediate" = shortest path with an edge limit. Bellman-Ford round `i` finds the best cost with at most `i` edges, so stop after `k + 1` rounds.

**Abstract problem:** Cheapest flights within k stops.

**Pattern:** Bellman-Ford with a hop limit (`14-union-find-shortest-paths-mst.md`).

**Approach:** Each round relaxes every edge from the previous round's snapshot, so a round never chains two edges.

```python
def cheapest_transfer(n: int, links: list[tuple[int, int, int]],
                      src: int, dst: int, k: int) -> int:
    cost = [inf] * n
    cost[src] = 0
    for _ in range(k + 1):  # <= k relays means <= k + 1 edges
        nxt = cost[:]  # relax from last round's snapshot so one round adds one edge
        for u, v, w in links:
            if cost[u] + w < nxt[v]:
                nxt[v] = cost[u] + w
        cost = nxt
    return -1 if cost[dst] == inf else int(cost[dst])
```

**Complexity:** O(k · E) time, O(n) space.

**Equivalent known problem:** Cheapest Flights Within K Stops.

**Trap:** Plain Dijkstra on cost finalizes a node via a cheap long path, discarding the pricier short path you needed. Relaxing in place lets one round use several edges.

</details>

### Drill 28: Cheapest fiber layout

You must connect `n` sites with fiber. Any pair can be linked directly at a cost equal to the Manhattan distance between their coordinates. Every site must be reachable from every other, directly or through other sites. Minimize total cost.

- `[(0,0), (2,2), (3,10), (5,2), (7,0)]` → `20`
- `[(3,12), (-2,5), (-4,1)]` → `18`
- Constraints: `n ≤ 1000`.

<details><summary>Hints</summary>

1. The optimal network has no cycles.
2. Minimum spanning tree of a complete graph.
3. Prim's algorithm, O(n^2) array version for dense graphs.

</details>

<details><summary>Solution</summary>

**Decode:** "Connect all ... minimize total cost" with pairwise costs = MST. The graph is complete (~500k edges), so O(n^2) array Prim beats sorting all edges.

**Abstract problem:** Min cost to connect all points.

**Pattern:** Minimum spanning tree (`14-union-find-shortest-paths-mst.md`).

**Approach:** `best[v]` = cheapest link from the tree to `v`. Repeatedly add the closest outside site, then update the rest.

```python
def min_fiber_cost(sites: list[tuple[int, int]]) -> int:
    n = len(sites)
    best = [inf] * n  # cheapest link from the current tree to each site
    if n:
        best[0] = 0
    in_tree = [False] * n
    total = 0
    for _ in range(n):
        u = min((i for i in range(n) if not in_tree[i]), key=best.__getitem__)
        in_tree[u] = True
        total += best[u]
        x, y = sites[u]
        for v in range(n):
            if not in_tree[v]:
                d = abs(x - sites[v][0]) + abs(y - sites[v][1])
                if d < best[v]:
                    best[v] = d
    return int(total)
```

**Complexity:** O(n^2) time, O(n) space.

**Equivalent known problem:** Min Cost to Connect All Points.

**Trap:** A shortest-path tree minimizes distances from one site, not total cable. Kruskal works, but sorting n^2/2 edges is O(n^2 log n) and memory-heavy.

</details>

### Drill 29: Hostname rename cost

To estimate DNS churn, you count the minimum number of single-character edits (insert, delete, or replace one character) that turn an old hostname into a new one.

- `"horse"` → `"ros"`: `3`
- `"intention"` → `"execution"`: `5`
- `"db-east-1"` → `"db-west-2"`: `3`
- Constraints: lengths `≤ 500`.

<details><summary>Hints</summary>

1. Compare position by position? What does an insertion do?
2. Minimum insert/delete/replace operations between two strings.
3. 2D DP on prefixes.

</details>

<details><summary>Solution</summary>

**Decode:** Two strings, "minimum edits", three operations = edit distance; lengths ≤ 500 allow O(m·n).

**Abstract problem:** Edit distance.

**Pattern:** DP 2D (`16-dynamic-programming-foundations.md`).

**Approach:** `dp[i][j]` = cost to turn `a[:i]` into `b[:j]`: the diagonal if characters match, else 1 + min(delete, insert, replace). Keep only the previous row.

```python
def min_rename_ops(a: str, b: str) -> int:
    m, n = len(a), len(b)
    prev = list(range(n + 1))  # row i-1: cost to turn a[:i-1] into b[:j]
    for i in range(1, m + 1):
        cur = [i] + [0] * n
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                cur[j] = prev[j - 1]
            else:  # delete a[i-1] / insert b[j-1] / replace
                cur[j] = 1 + min(prev[j], cur[j - 1], prev[j - 1])
        prev = cur
    return prev[n]
```

**Complexity:** O(m·n) time, O(n) space.

**Equivalent known problem:** Edit Distance.

**Trap:** Hamming distance breaks when an insertion shifts everything. `m + n - 2·LCS` counts only inserts and deletes, overcounting replacements.

</details>

### Drill 30: Shortest heartbeat unit

A heartbeat log line is supposed to be one short unit repeated. Given the captured string, return the shortest unit that, repeated a whole number of times, produces exactly the string. If no shorter unit exists, return the string itself.

- `"ab.ab.ab."` → `"ab."`
- `"abcab"` → `"abcab"`
- `"abaaba"` → `"aba"`
- Constraints: length `≤ 10^6`.

<details><summary>Hints</summary>

1. For a repetition, the longest proper prefix that is also a suffix gives the shift.
2. Smallest period `p = n - border(n)`; it's a valid unit only if `p` divides `n`.
3. KMP prefix function.

</details>

<details><summary>Solution</summary>

**Decode:** "Repeats itself" at length 10^6 needs a linear string algorithm. The KMP prefix function gives the whole string's longest border; the smallest period is `n - border`.

**Abstract problem:** Repeated substring pattern, returning the unit.

**Pattern:** String algorithms, KMP prefix function (`19-tries-and-string-algorithms.md`).

**Approach:** Compute `pi`; `p = n - pi[n-1]`. If `n % p == 0` the unit is `s[:p]`; otherwise no proper unit tiles the string.

```python
def shortest_heartbeat_unit(s: str) -> str:
    n = len(s)
    pi = [0] * n  # pi[i]: length of longest proper border of s[:i+1]
    for i in range(1, n):
        k = pi[i - 1]
        while k and s[i] != s[k]:
            k = pi[k - 1]
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    period = n - pi[-1] if n else 0
    return s[:period] if n and n % period == 0 else s
```

**Complexity:** O(n) time and space.

**Equivalent known problem:** Repeated Substring Pattern.

**Trap:** Testing every divisor length by string comparison is O(n · d(n)), slow in Python at 10^6. `s in (s + s)[1:-1]` answers only yes/no, not the unit.

</details>

### Drill 31: Live p50 latency

Request latencies (ms) arrive one at a time. After each arrival the dashboard asks for the current median (p50) of everything seen so far. With an even count, report the mean of the two middle values.

- arrivals `120, 80, 95, 300, 101` → p50 after each: `120, 100, 95, 107.5, 101`
- Constraints: up to `5·10^4` records and queries.

<details><summary>Hints</summary>

1. You only need the middle, not a fully sorted list.
2. Keep lower and upper halves with fast access to their boundary values.
3. Two heaps.

</details>

<details><summary>Solution</summary>

**Decode:** "Stream" + "after each arrival" + "median" = online order statistic; re-sorting is O(n log n) per query.

**Abstract problem:** Find median from data stream.

**Pattern:** Two heaps (`11-heaps-top-k-and-merging.md`).

**Approach:** Max-heap `low`, min-heap `high`. Route each insert through `low` into `high` (keeping every low ≤ every high), then rebalance so `low` holds the extra element on odd counts.

```python
class P50Tracker:
    def __init__(self) -> None:
        self.low: list[int] = []   # max-heap (negated): smaller half
        self.high: list[int] = []  # min-heap: larger half
        # invariant: len(low) == len(high) or len(low) == len(high) + 1

    def record(self, ms: int) -> None:
        heapq.heappush(self.low, -ms)
        heapq.heappush(self.high, -heapq.heappop(self.low))  # keeps low <= high
        if len(self.high) > len(self.low):
            heapq.heappush(self.low, -heapq.heappop(self.high))

    def p50(self) -> float:
        if len(self.low) > len(self.high):
            return float(-self.low[0])
        return (-self.low[0] + self.high[0]) / 2
```

**Complexity:** O(log n) per record, O(1) per query, O(n) space.

**Equivalent known problem:** Find Median from Data Stream.

**Trap:** `bisect.insort` is O(n) per insert (quadratic overall). Choosing a heap by comparison without rebalancing breaks the size invariant.

</details>

### Drill 32: Peak load per window

You have per-second request rates for a service. To size burst capacity, report the peak rate in every window of `k` consecutive seconds, sliding by one second.

- `rates = [1, 3, 1, 0, 5, 3, 6, 7], k = 3` → `[3, 3, 5, 5, 6, 7]`
- Constraints: `n ≤ 10^5`, `1 ≤ k ≤ n`.

<details><summary>Hints</summary>

1. Once a newer, larger value enters, older smaller ones can never be the max.
2. Max of every fixed-size subarray.
3. Monotonic deque.

</details>

<details><summary>Solution</summary>

**Decode:** "Every window of k consecutive" = fixed sliding window; "peak" = max, which (unlike sum) can't be updated by subtraction. Discarding dominated elements gives a monotonic structure.

**Abstract problem:** Sliding window maximum.

**Pattern:** Monotonic deque (`08-stacks-and-monotonic-structures.md`, `06-sliding-window.md`).

**Approach:** Deque of indices with strictly decreasing values: on insert, pop smaller-or-equal values from the back; pop the front once it leaves the window. The front is the max.

```python
def window_peaks(load: list[int], k: int) -> list[int]:
    dq: deque[int] = deque()  # indices in window, loads strictly decreasing front->back
    out = []
    for i, x in enumerate(load):
        while dq and load[dq[-1]] <= x:
            dq.pop()  # dominated: older and not larger, can never be a max again
        dq.append(i)
        if dq[0] <= i - k:
            dq.popleft()  # front slid out of the window
        if i >= k - 1:
            out.append(load[dq[0]])
    return out
```

**Complexity:** O(n) time, O(k) space.

**Equivalent known problem:** Sliding Window Maximum.

**Trap:** Recomputing `max` per window is O(n·k). A heap needs lazy deletion of expired indices and costs O(n log n).

</details>

### Drill 33: Smallest trace slice covering services

Log lines are ordered and each is tagged with the service that wrote it. For a distributed trace you need the shortest contiguous slice containing every required service at least as many times as it appears in `required`. Return `(start, end)` inclusive, the earliest on ties, or `(-1, -1)`.

- `logs = ["auth","db","cache","auth","api","db","cache","api"]`
- `required = ["api","db","cache"]` → `(4, 6)`; `required = ["auth","api"]` → `(3, 4)`
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. "Contiguous" + "shortest" + "contains all".
2. Minimum window containing a multiset of tokens.
3. Variable sliding window with a `missing` counter.

</details>

<details><summary>Solution</summary>

**Decode:** "Shortest contiguous slice containing every required" = minimum window. Supersets of a valid window are valid, so expand right until valid, then shrink left while valid.

**Abstract problem:** Minimum window substring over tokens.

**Pattern:** Sliding window, shrink-while-valid (`06-sliding-window.md`).

**Approach:** `need` and `have` counts, plus `missing` = distinct services not yet satisfied, updated only when a count crosses its requirement, making the validity check O(1).

```python
def shortest_covering_span(logs: list[str], required: list[str]) -> tuple[int, int]:
    need = Counter(required)
    missing = len(need)  # distinct services whose required count is not met yet
    have: Counter[str] = Counter()
    best = (-1, -1)
    left = 0
    for right, svc in enumerate(logs):
        if svc in need:
            have[svc] += 1
            if have[svc] == need[svc]:
                missing -= 1
        while missing == 0:  # window valid: record, then shrink from the left
            if best[0] == -1 or right - left < best[1] - best[0]:
                best = (left, right)
            out = logs[left]
            if out in need:
                have[out] -= 1
                if have[out] < need[out]:
                    missing += 1
            left += 1
    return best
```

**Complexity:** O(n + m) time, O(m) space.

**Equivalent known problem:** Minimum Window Substring.

**Trap:** Comparing whole counters each step is O(n · m). A set instead of counts breaks when a service is required twice.

</details>

### Drill 34: Largest flat reservation

A capacity forecast gives free GPU count for each of the next `n` hours. You want to book one reservation of a constant size `c` over a contiguous range of hours, where `c` must not exceed the free capacity at any hour in the range. Maximize `c × hours` (GPU-hours booked).

- `free = [2, 1, 5, 6, 2, 3]` → `10` (`c = 5` for hours 2..3)
- `free = [3, 3, 3]` → `9`
- Constraints: `n ≤ 10^5`.

<details><summary>Hints</summary>

1. For a fixed range, the best `c` is the minimum in that range.
2. Maximize `min(range) × length` over all contiguous ranges.
3. Largest rectangle in a histogram: monotonic stack.

</details>

<details><summary>Solution</summary>

**Decode:** "Constant size under a curve over a contiguous range" = rectangle in a histogram. Flip it: treat each hour as the limiting height and extend to the nearest smaller bar on each side, found by a monotonic stack.

**Abstract problem:** Largest rectangle in histogram.

**Pattern:** Monotonic stack (`08-stacks-and-monotonic-structures.md`).

**Approach:** Non-decreasing stack of indices. When a lower bar arrives, pop: the popped bar's rectangle spans from the new top + 1 to the current index - 1. A trailing 0 sentinel flushes the stack.

```python
def largest_reservation(free: list[int]) -> int:
    stack: list[int] = []  # indices with non-decreasing free capacity
    best = 0
    for i, h in enumerate(free + [0]):  # sentinel 0 flushes the stack at the end
        while stack and free[stack[-1]] >= h:
            height = free[stack.pop()]
            left = stack[-1] + 1 if stack else 0  # first index the bar can extend to
            best = max(best, height * (i - left))
        stack.append(i)
    return best
```

**Complexity:** O(n) time, O(n) space.

**Equivalent known problem:** Largest Rectangle in Histogram.

**Trap:** Container With Most Water's two pointers use only endpoints, but here every hour in the range caps `c`. Expanding from each bar is O(n^2) on sorted input.

</details>

### Drill 35: Splitting a log archive

A log archive of `length` GB must be split at the given offsets before upload to cold storage. The splitting tool reads the entire chunk it is splitting, so each split costs the size of the current chunk. You choose the order of splits. Minimize total cost.

- `length = 7, offsets = [1, 3, 4, 5]` → `16`
- `length = 9, offsets = [5, 6, 1, 4, 2]` → `22`
- Constraints: `length ≤ 10^6`, up to `100` offsets.

<details><summary>Hints</summary>

1. The first split creates two independent subproblems.
2. For a segment between two cut points, choose which inner cut to make first.
3. Interval DP.

</details>

<details><summary>Solution</summary>

**Decode:** "Order of operations matters" + "cost depends on the current segment" + ~100 cut points (O(m^3) is fine) = interval DP. After a segment's first cut, the two sides never interact.

**Abstract problem:** Minimum cost to cut a stick.

**Pattern:** Interval DP (`17-dynamic-programming-advanced.md`).

**Approach:** Add `0` and `length` as boundaries and sort. `dp[i][j] = (cuts[j] - cuts[i]) + min over k in (i, j) of dp[i][k] + dp[k][j]`, filled by increasing gap.

```python
def min_split_cost(length: int, offsets: list[int]) -> int:
    cuts = [0] + sorted(offsets) + [length]
    m = len(cuts)
    # dp[i][j]: min cost to make every cut strictly between cuts[i] and cuts[j]
    dp = [[0] * m for _ in range(m)]
    for gap in range(2, m):
        for i in range(m - gap):
            j = i + gap
            dp[i][j] = cuts[j] - cuts[i] + min(dp[i][k] + dp[k][j] for k in range(i + 1, j))
    return dp[0][m - 1]
```

**Complexity:** O(m^3) time, O(m^2) space for `m` cut points.

**Equivalent known problem:** Minimum Cost to Cut a Stick.

**Trap:** Cutting nearest the middle isn't guaranteed optimal. Reversed, it resembles Huffman merging, but only adjacent pieces may merge, so Huffman's greedy fails. DP over `length` instead of cut indices blows up at 10^6.

</details>

### Drill 36: Certificate rotation assignment

`n` engineers must each take exactly one of `n` certificate rotation tasks (one per cluster). `cost[i][j]` is the hours engineer `i` needs for task `j`, based on their familiarity with that cluster. Minimize the total hours.

- `[[9, 2, 7], [6, 4, 3], [5, 8, 1]]` → `9` (engineer 0 → task 1, 1 → task 0, 2 → task 2)
- Constraints: `1 ≤ n ≤ 12`.

<details><summary>Hints</summary>

1. `n ≤ 12` is loud: `12!` is 479 million, but `2^12 · 12` is about 49k.
2. Assign engineers in order; the state is which tasks are already taken.
3. Bitmask DP.

</details>

<details><summary>Solution</summary>

**Decode:** "Each to exactly one" = a permutation; tiny `n` = exponential over subsets, not permutations. With engineers assigned in fixed order, the used-task set alone determines the next engineer (its popcount).

**Abstract problem:** Assignment problem (minimum-cost perfect matching on a small bipartite graph).

**Pattern:** Bitmask DP (`17-dynamic-programming-advanced.md`).

**Approach:** `dp[mask]` = min cost of giving tasks in `mask` to engineers `0..popcount(mask)-1`; extend each reachable mask with one free task for the next engineer.

```python
def min_rotation_cost(cost: list[list[int]]) -> int:
    n = len(cost)
    dp = [inf] * (1 << n)  # dp[mask]: min cost giving tasks in mask to engineers 0..popcount-1
    dp[0] = 0
    for mask in range(1 << n):
        eng = mask.bit_count()  # next engineer to assign
        if dp[mask] == inf or eng == n:
            continue
        for task in range(n):
            if not mask & (1 << task):
                nm = mask | (1 << task)
                dp[nm] = min(dp[nm], dp[mask] + cost[eng][task])
    return int(dp[-1])
```

**Complexity:** O(2^n · n) time, O(2^n) space.

**Equivalent known problem:** Minimum XOR Sum of Two Arrays (same DP with `cost = a[i] ^ b[j]`); also Maximum Compatibility Score Sum.

**Trap:** Greedy (each engineer takes their cheapest remaining task) gives engineer 0 task 1 and engineer 1 task 2, leaving engineer 2 cost 5: total 10, not 9. Enumerating permutations is O(n!). The Hungarian algorithm is O(n^3) but rarely expected.

</details>

## 3. Pattern coverage index

Use only after attempting drills blind.

| Drill | Story | Pattern(s) | Module |
|---|---|---|---|
| 1 | GPU quota pairing | Hash map (complement) | 04 |
| 2 | Missing shard | Bit manipulation (XOR) | 20 |
| 3 | First regressed build | Boundary binary search | 07 |
| 4 | Minutes until higher load | Monotonic stack | 08 |
| 5 | Failure incidents on floor | Grid DFS, connected components | 13 |
| 6 | Two pods per node | Two pointers + greedy | 05, 18 |
| 7 | Health-check streak | Sliding window with budget | 06 |
| 8 | Net replica change k | Prefix sums + hash map | 04 |
| 9 | Common aggregation switch | Tree LCA | 10 |
| 10 | CI runners | Intervals, sweep line | 12 |
| 11 | Malware spreading | Multi-source BFS | 13 |
| 12 | Terraform apply order | Topological sort (Kahn) | 13 |
| 13 | Merging node logs | Heap, k-way merge | 11 |
| 14 | Pager decodings | DP 1D (count ways) | 16 |
| 15 | Redundant cable | Union-find | 14 |
| 16 | Config push | Dijkstra | 14 |
| 17 | Instance mixes | Backtracking (with reuse) | 15 |
| 18 | Follow-the-sun loop | Greedy | 18 |
| 19 | CLI autocomplete | Trie | 19 |
| 20 | Duplicated ticket | Fast/slow pointers (Floyd) | 09 |
| 21 | Migration bandwidth cap | Binary search on answer | 07 |
| 22 | Rebuilding IPs | Backtracking with pruning | 15 |
| 23 | Two GPU nodes split | 0/1 knapsack (subset sum) | 16 |
| 24 | Kernel cooldown | Greedy counting (heap alternative) | 18, 11 |
| 25 | Version-tag rollout | BFS on implicit graph | 13 |
| 26 | Spot reservations | State-machine DP | 17 |
| 27 | Transfer with k relays | Bellman-Ford with hop limit | 14 |
| 28 | Fiber layout | MST (Prim) | 14 |
| 29 | Hostname rename | DP 2D (edit distance) | 16 |
| 30 | Heartbeat unit | KMP prefix function | 19 |
| 31 | Live p50 | Two heaps | 11 |
| 32 | Peak load per window | Monotonic deque | 08, 06 |
| 33 | Trace slice | Sliding window (min window) | 06 |
| 34 | Flat reservation | Monotonic stack (histogram) | 08 |
| 35 | Splitting archive | Interval DP | 17 |
| 36 | Certificate assignment | Bitmask DP | 17 |

By pattern: arrays/hashing 1; prefix sums 8; two pointers 6; sliding window 7, 32, 33; binary search 3, 21; monotonic stack/deque 4, 32, 34; linked list 20; trees 9; heaps 13, 24, 31; intervals 10; grid/BFS/DFS 5, 25; multi-source BFS 11; topo sort 12; union-find 15; Dijkstra 16; Bellman-Ford 27; MST 28; backtracking 17, 22; DP 1D 14; DP 2D 29; knapsack 23; interval DP 35; state-machine DP 26; bitmask DP 36; greedy 6, 18, 24; trie 19; string algorithms 30; bit manipulation 2.

## 4. How to use

1. **Blind first.** Read only the story. On a 10-minute timer, write the abstract statement, object types, question type, complexity budget, and pattern guess. No hints before the timer ends.
2. **Escalate one hint at a time.** Needing hint 1 means you missed a noticing step; needing hint 3 means you missed the mapping itself.
3. **Code it** without looking, compare with the solution, and run yours against the asserts in `code/22_disguised_drills.py`.
4. **Log each miss by the signal you missed**, not the problem name: "missed that negatives kill the window", "did not read n ≤ 12 as bitmask", "treated at-most-k-hops as plain Dijkstra". The log becomes your personal trigger list.
5. **Spaced redo.** Redo missed drills blind after 3 and 7 days. Retire a drill after decoding it in under 2 minutes twice in a row.
6. **Read the Trap line even when you got it right.** Follow-ups often push toward the trap to see if you can defend your choice.
