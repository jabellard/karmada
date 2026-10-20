# 13. Graphs: Traversal and Topological Sort

Dependency resolution, service call graphs, network reachability, build systems, puzzles, and grid maps are graph problems in disguise. BFS, DFS, and topological sort cover most of them; the hard part is seeing the graph, deciding what a node is, and what goes into `visited`.

## At a glance

- **Use when:** edge lists, prerequisites, grids, or states with moves; "fewest steps", "how many regions", "valid order".
- **Templates:**
  - Template 1: BFS shortest path: fewest edges, one source, unweighted.
  - Template 2: multi-source BFS: simultaneous spread, distance to nearest X.
  - Template 3: grid DFS / flood fill: count or measure regions.
  - Template 4: Kahn's topological sort: dependency order, cycle check via count.
  - Template 5: 3-color DFS: directed cycle detection, safe states, postorder topo.
  - Template 6: state-space BFS: fewest moves over generated states.
- **Typical complexity:** O(V + E), or O(R*C) on a grid.
- **Top trap:** marking visited on dequeue, or a state missing something that affects future moves.

## Core idea

A graph is nodes plus a neighbor function. Traversal visits each reachable node once using a `visited` set.

- **BFS** (queue) explores in layers. Invariant: a node's distance when first discovered is the fewest edges from the source(s), so BFS finds shortest paths in **unweighted** graphs.
- **DFS** (stack or recursion) goes deep first: reachability, components, flood fill, cycle detection, postorder (topological order).
- **Topological sort** orders a DAG so every edge `u -> v` has `u` first. Kahn's repeatedly removes indegree-0 nodes; nodes never reaching 0 are on or behind a cycle.
- **Implicit graphs**: nodes are states (a lock combination, a board, `(cell, keys_held)`), edges are legal moves, neighbors are generated on demand.

```
BFS layers from S:            Kahn on a -> b, a -> c, b -> d, c -> d
    S                          indeg: a0 b1 c1 d1
   / \      layer 0: S         pop a -> b0 c0   order: a
  A   B     layer 1: A B       pop b -> d0      order: a b
  |   |     layer 2: C D       pop c            order: a b c
  C   D                        pop d            order: a b c d
```

In BFS, mark visited on **enqueue**. Marking on dequeue lets a node be enqueued many times (blowups on dense grids). Dijkstra is the opposite: finalize on pop (`14-union-find-shortest-paths-mst.md`).

## Recognition signals

- Input: edge list, adjacency list/matrix, `prerequisites`, `dependencies`, `tickets`, `(a, b)` pairs, or a 2D grid.
- "Fewest steps / moves / transformations / minutes", uniform cost → BFS.
- "Number of groups / regions / islands / provinces", "all connected", "can reach" → DFS/BFS components (or DSU, module 14).
- "X before Y", "prerequisites", "build order", "possible to finish" → topological sort / cycle detection.
- "Spreads from all X simultaneously", "distance to nearest X" → multi-source BFS.
- "Split into two groups with no conflicts inside a group" → bipartite check.
- Small state (a few numbers or a short string), "minimum moves" → state-space BFS, if the states fit in memory (roughly ≤ 10^6 to 10^7).
- Grid up to 10^3 x 10^3 → O(rows*cols); iterative DFS or BFS avoids recursion limits.

## Disguises

1. *A package manager has manifests saying which packages depend on which. Produce an install order, or report the circular dependency.*
   Decode: "X needs Y first" = edge Y -> X. Install order = topological sort; "circular" = cycle, so Kahn's output is shorter than n.

2. *A datacenter floor plan marks failed racks. Heat spreads to adjacent racks each minute. When does the last healthy rack overheat?*
   Decode: spread from many origins at once at uniform speed → multi-source BFS. Answer = number of layers; -1 if a cell is unreachable.

3. *You are given a log of messages "server A pinged server B". How many isolated clusters are there?*
   Decode: undirected edges, count components → DFS/BFS from each unvisited node, or DSU.

4. *A config value can change one character at a time, and every intermediate value must pass validation (be in an allowed set). Minimum edits from current to desired?*
   Decode: strings are nodes, one-character edits within the allowed set are edges → BFS over states (Word Ladder / Genetic Mutation).

5. *Assign each of n workers to shift A or B. Some pairs refuse to share a shift. Is it possible?*
   Decode: refusals are edges, shifts are two colors → bipartite check.

6. *A robot on a grid can smash through at most k walls. Fewest moves to the exit?*
   Decode: BFS, but position alone isn't the state. State = `(r, c, smashes_left)`.

## Templates and worked examples

Imports (companion file: `code/13_graphs_traversal_topo.py`):

```python
from collections import defaultdict, deque
from typing import Callable, Hashable, Iterable
```

Build the graph first: edge lists become adjacency lists; grids use a direction helper.

```python
def build_adj(n: int, edges: list[list[int]], directed: bool = False) -> list[list[int]]:
    """Adjacency list for nodes 0..n-1 from an edge list."""
    adj: list[list[int]] = [[] for _ in range(n)]
    for u, v in edges:
        adj[u].append(v)
        if not directed:
            adj[v].append(u)
    return adj


DIRS4 = ((1, 0), (-1, 0), (0, 1), (0, -1))
DIRS8 = DIRS4 + ((1, 1), (1, -1), (-1, 1), (-1, -1))


def grid_neighbors(r: int, c: int, rows: int, cols: int, dirs=DIRS4):
    """Yield in-bounds neighbors of (r, c)."""
    for dr, dc in dirs:
        nr, nc = r + dr, c + dc
        if 0 <= nr < rows and 0 <= nc < cols:
            yield nr, nc
```

For string or tuple labels, use `defaultdict(list)`.

### Template 1: BFS shortest path (unweighted)

> - **Example problem:** [Shortest Path in Binary Matrix (1091)](https://leetcode.com/problems/shortest-path-in-binary-matrix/)
> - **Why this template:** "shortest clear path" with unit-cost moves = fewest edges in an unweighted graph: plain BFS.
> - **Worked below:** Example 1.1

```python
def bfs_shortest_path(adj: list[list[int]], src: int, dst: int) -> int:
    """Fewest edges from src to dst in an unweighted graph, -1 if unreachable.

    Mark visited when ENQUEUING: each node enters the queue once, and the first
    time it is reached is via a shortest path (BFS explores by layers).
    """
    if src == dst:
        return 0
    dist = {src: 0}
    q = deque([src])
    while q:
        u = q.popleft()
        for v in adj[u]:                      # TODO: neighbor generator
            if v not in dist:                 # TODO: extra filters (walls, constraints)
                dist[v] = dist[u] + 1
                if v == dst:                  # early exit is safe on enqueue for BFS
                    return dist[v]
                q.append(v)
    return -1
```

Knobs: `dist` doubles as `visited`. For layer-by-layer processing (minutes, levels), loop `for _ in range(len(q))` inside the `while`. A `parent` map reconstructs the path.

#### Example 1.1: Shortest clear path across a grid, diagonals allowed [Medium]

**Problem:** [Shortest Path in Binary Matrix (1091)](https://leetcode.com/problems/shortest-path-in-binary-matrix/) · **Template:** Template 1 (BFS shortest path (unweighted))

Problem: an n x n grid of 0 (open) and 1 (blocked). Moving in any of 8 directions through open cells, return the number of cells on the shortest path from top-left to bottom-right (both included), or -1. Example: `[[0, 1], [1, 0]]` → 2 (one diagonal step).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "shortest" + every move costs 1 = BFS on the grid graph; diagonals just mean 8 neighbors (`DIRS8`).

**Approach:** BFS from `(0, 0)` with `dist[0][0] = 1`, since the answer counts cells, not moves. Mark a cell visited when it's enqueued (setting `dist`), so each cell enters the queue once at its shortest distance. Return when the target is popped.

```python
def shortest_path_binary_matrix(grid: list[list[int]]) -> int:
    """Shortest Path in Binary Matrix: cells on the shortest 8-directional clear path, or -1."""
    n = len(grid)
    if grid[0][0] or grid[n - 1][n - 1]:
        return -1
    dist = [[0] * n for _ in range(n)]        # 0 = unvisited, else cells on the path so far
    dist[0][0] = 1
    q = deque([(0, 0)])
    while q:
        r, c = q.popleft()
        if r == c == n - 1:
            return dist[r][c]
        for nr, nc in grid_neighbors(r, c, n, n, DIRS8):
            if grid[nr][nc] == 0 and not dist[nr][nc]:
                dist[nr][nc] = dist[r][c] + 1  # mark on enqueue
                q.append((nr, nc))
    return -1
```

**Complexity:** O(n^2) time and space: each cell is enqueued once and checks 8 neighbors.

**Edge cases:** blocked start or end → -1 immediately; `[[0]]` → 1 (start is the target); marking on dequeue instead lets a cell be enqueued many times (still correct, but slower).

</details>

### Template 2: multi-source BFS

> - **Example problem:** [Rotting Oranges (994)](https://leetcode.com/problems/rotting-oranges/)
> - **Why this template:** all rotten oranges spread each minute at once, so every rotten cell is a source at distance 0.
> - **Worked below:** Example 2.1
> - **Also fits:** [Walls and Gates (286)](https://leetcode.com/problems/walls-and-gates/) (Premium)

```python
def multi_source_bfs(grid: list[list[int]], is_source: Callable[[int], bool],
                     passable: Callable[[int], bool]) -> list[list[int]]:
    """Distance from every cell to its nearest source (-1 if unreachable).

    Equivalent to adding a virtual super-source with 0-weight edges to all sources.
    """
    rows, cols = len(grid), len(grid[0])
    dist = [[-1] * cols for _ in range(rows)]
    q: deque[tuple[int, int]] = deque()
    for r in range(rows):
        for c in range(cols):
            if is_source(grid[r][c]):         # TODO: what counts as a source
                dist[r][c] = 0
                q.append((r, c))
    while q:
        r, c = q.popleft()
        for nr, nc in grid_neighbors(r, c, rows, cols):
            if dist[nr][nc] == -1 and passable(grid[nr][nc]):   # TODO: passability
                dist[nr][nc] = dist[r][c] + 1
                q.append((nr, nc))
    return dist
```

Knobs: sources and passability. Seeding all sources at distance 0 gives each cell its nearest source in O(V), versus O(k * V) for one BFS per source. Think in reverse: "distance from each empty cell to the nearest gate" = BFS *from all gates*.

#### Example 2.1: Simultaneous spread [Medium]

**Problems:** [Rotting Oranges (994)](https://leetcode.com/problems/rotting-oranges/), [Walls and Gates (286)](https://leetcode.com/problems/walls-and-gates/) (Premium) · **Template:** Template 2 (multi-source BFS)

Problem: a grid of empty cells, fresh oranges, and rotten oranges. Each minute, rot spreads to 4-adjacent fresh oranges. Return minutes until none are fresh, or -1. Variant: fill each empty room with its distance to the nearest gate.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** all rotten oranges spread together = multi-source BFS; minutes = layers.

**Approach:** enqueue every rotten orange, count fresh ones, and process one layer (minute) at a time until fresh hits 0. Fresh left when the queue empties is unreachable.

```python
def oranges_rotting(grid: list[list[int]]) -> int:
    """Rotting Oranges: 0 empty, 1 fresh, 2 rotten. Minutes until no fresh, or -1."""
    rows, cols = len(grid), len(grid[0])
    q: deque[tuple[int, int]] = deque()
    fresh = 0
    for r in range(rows):
        for c in range(cols):
            if grid[r][c] == 2:
                q.append((r, c))
            elif grid[r][c] == 1:
                fresh += 1
    minutes = 0
    while q and fresh:
        for _ in range(len(q)):               # process one full layer = one minute
            r, c = q.popleft()
            for nr, nc in grid_neighbors(r, c, rows, cols):
                if grid[nr][nc] == 1:
                    grid[nr][nc] = 2
                    fresh -= 1
                    q.append((nr, nc))
        minutes += 1
    return -1 if fresh else minutes
```

Walls and Gates is Template 2 with gates as sources and only empty rooms passable:

```python
def walls_and_gates(rooms: list[list[int]]) -> list[list[int]]:
    """Walls and Gates: -1 wall, 0 gate, INF empty -> fill with distance to nearest gate."""
    INF = 2**31 - 1
    dist = multi_source_bfs(rooms, is_source=lambda v: v == 0, passable=lambda v: v == INF)
    return [[rooms[r][c] if rooms[r][c] == -1 or dist[r][c] == -1 else dist[r][c]
             for c in range(len(rooms[0]))] for r in range(len(rooms))]
```

**Complexity:** O(R*C) time and space; each cell is enqueued at most once.

**Edge cases:** no fresh oranges → 0 (the `while q and fresh` guard prevents an extra minute); fresh but no rotten, or an isolated fresh orange → -1. Same template: **01 Matrix** (sources = zeros), **As Far from Land as Possible** (sources = land, answer = max distance).

</details>

### Template 3: grid DFS / flood fill

> - **Example problem:** [Number of Islands (200)](https://leetcode.com/problems/number-of-islands/)
> - **Why this template:** islands are 4-connected land groups, so each flood fill from unvisited land counts one component.
> - **Worked below:** Example 3.1
> - **Also fits:** [Max Area of Island (695)](https://leetcode.com/problems/max-area-of-island/)

```python
def grid_dfs_iterative(grid: list[list[str]], r: int, c: int, target: str, mark: str) -> int:
    """Flood fill from (r, c) over cells equal to `target`, overwriting with `mark`.

    Returns component size. Iterative to avoid Python's recursion limit on big grids.
    Marking on push guarantees each cell is pushed once.
    """
    rows, cols = len(grid), len(grid[0])
    if grid[r][c] != target:
        return 0
    grid[r][c] = mark                         # in-place visited marker (TODO: or a visited set)
    stack = [(r, c)]
    size = 0
    while stack:
        cr, cc = stack.pop()
        size += 1
        for nr, nc in grid_neighbors(cr, cc, rows, cols):   # TODO: DIRS8 for diagonal adjacency
            if grid[nr][nc] == target:
                grid[nr][nc] = mark
                stack.append((nr, nc))
    return size
```

Knobs: `DIRS4` vs `DIRS8`; in-place marking (ask if mutation is allowed) vs a `visited` set; return size, perimeter, or a touches-border flag. `deque.popleft()` makes it BFS. Recursive DFS is shorter but hits Python's ~1000 recursion limit on a 300x300 all-land grid.

#### Example 3.1: Count connected regions in a grid [Medium]

**Problems:** [Number of Islands (200)](https://leetcode.com/problems/number-of-islands/), [Max Area of Island (695)](https://leetcode.com/problems/max-area-of-island/) · **Template:** Template 3 (grid DFS / flood fill)

Problem: a grid of land `"1"` and water `"0"`. Count islands (4-directionally connected land). Variant: return the largest island's size.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** connected regions in a grid → flood fill (Template 3).

**Approach:** scan every cell; at each unvisited land cell, flood-fill it to water and count one. Each cell is visited a constant number of times.

```python
def num_islands(grid: list[list[str]]) -> int:
    """Number of Islands: count connected components of '1' (mutates grid)."""
    count = 0
    for r in range(len(grid)):
        for c in range(len(grid[0])):
            if grid[r][c] == "1":
                grid_dfs_iterative(grid, r, c, "1", "0")
                count += 1
    return count
```

The recursive version reads cleanly and suits small grids (≤ 50x50):

```python
def max_area_of_island(grid: list[list[int]]) -> int:
    """Max Area of Island, recursive DFS variant (fine for small grids)."""
    rows, cols = len(grid), len(grid[0])

    def dfs(r: int, c: int) -> int:
        if not (0 <= r < rows and 0 <= c < cols) or grid[r][c] != 1:
            return 0
        grid[r][c] = 0
        return 1 + sum(dfs(r + dr, c + dc) for dr, dc in DIRS4)

    return max((dfs(r, c) for r in range(rows) for c in range(cols)), default=0)
```

**Complexity:** O(R*C) time; O(R*C) worst-case stack.

**Edge cases:** strings `"1"` vs ints `1`; mutating input (offer a `visited` set); recursion depth; "8 directions" means diagonals.

Island variants: **closed islands / enclaves** (flood from the border first, count what remains); **surrounded regions** (same border trick, flip the rest); **Pacific Atlantic** (multi-source from each ocean's border going *uphill*, intersect the reachable sets); **distinct island shapes** (record the DFS path signature relative to the start).

</details>

### Template 4: Kahn's topological sort

> - **Example problem:** [Course Schedule II (210)](https://leetcode.com/problems/course-schedule-ii/)
> - **Why this template:** "b before a" is a directed edge; "return an order" is a topological order (empty if a cycle blocks it).
> - **Worked below:** Example 4.1
> - **Also fits:** [Alien Dictionary (269)](https://leetcode.com/problems/alien-dictionary/) (Premium)

```python
def topo_sort_kahn(n: int, edges: list[list[int]]) -> list[int]:
    """Kahn's algorithm. Edge [u, v] means u must come before v.

    Returns an order, or [] if a cycle exists (some nodes never reach indegree 0).
    """
    adj: list[list[int]] = [[] for _ in range(n)]
    indeg = [0] * n
    for u, v in edges:
        adj[u].append(v)
        indeg[v] += 1
    q = deque(i for i in range(n) if indeg[i] == 0)   # TODO: heap here for lexicographically smallest order
    order: list[int] = []
    while q:
        u = q.popleft()
        order.append(u)
        for v in adj[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                q.append(v)
    return order if len(order) == n else []
```

Knobs: **edge direction** (`[a, b]` meaning "b before a" needs edge `b -> a`); queue vs heap (lexicographically smallest order); processing by layers gives the minimum semesters / parallel rounds (longest DAG path). The order is unique exactly when the queue never holds more than one node.

#### Example 4.1: Prerequisite ordering [Medium]

**Problems:** [Course Schedule (207)](https://leetcode.com/problems/course-schedule/), [Course Schedule II (210)](https://leetcode.com/problems/course-schedule-ii/) · **Template:** Template 4 (Kahn's topological sort)

Problem: n courses; pair `[a, b]` means take b before a. (I) Can you finish all? (II) Return any valid order, or empty.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "before" constraints → directed graph; "possible" = no cycle; "order" = topological sort.

**Approach:** flip each pair into edge `b -> a` and run Kahn's. All n output means a DAG: a cycle node keeps indegree ≥ 1, so it never enters the queue.

```python
def can_finish(num_courses: int, prerequisites: list[list[int]]) -> bool:
    """Course Schedule: pair [a, b] means b before a. Possible iff DAG."""
    return len(topo_sort_kahn(num_courses, [[b, a] for a, b in prerequisites])) == num_courses


def find_order(num_courses: int, prerequisites: list[list[int]]) -> list[int]:
    """Course Schedule II: any valid order, or [] if impossible."""
    return topo_sort_kahn(num_courses, [[b, a] for a, b in prerequisites])
```

**Complexity:** O(V + E) time and space.

**Edge cases:** initialize all n nodes so courses without prerequisites appear; duplicate pairs are fine if counted consistently; self-loop `[a, a]` is a cycle. Follow-ups: minimum semesters with unlimited parallelism = number of Kahn layers; at most k per semester makes it bitmask DP.

</details>

#### Example 4.2: Derive an order from sorted data [Hard]

**Problem:** [Alien Dictionary (269)](https://leetcode.com/problems/alien-dictionary/) (Premium) · **Template:** Template 4 (Kahn's topological sort)

Problem: words are sorted by an unknown alphabet. Return any consistent letter order, or `""` if the input is contradictory.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "sorted by unknown rules" → comparisons give "x before y" constraints → topological sort on a graph you extract.

**Approach:** in each adjacent pair, the first differing position gives edge `a -> b`; later positions say nothing. No difference with the first word longer (`"abc"` before `"ab"`) is invalid. Make every letter a node and run Kahn's; a cycle means contradiction.

```python
def alien_order(words: list[str]) -> str:
    """Alien Dictionary: derive letter order from a sorted word list; '' if invalid."""
    adj: dict[str, set[str]] = {ch: set() for w in words for ch in w}
    indeg = {ch: 0 for ch in adj}
    for w1, w2 in zip(words, words[1:]):
        for a, b in zip(w1, w2):
            if a != b:                        # first difference is the only info this pair gives
                if b not in adj[a]:
                    adj[a].add(b)
                    indeg[b] += 1
                break
        else:
            if len(w1) > len(w2):             # "abc" before "ab" is impossible
                return ""
    q = deque(ch for ch in adj if indeg[ch] == 0)
    out: list[str] = []
    while q:
        ch = q.popleft()
        out.append(ch)
        for nxt in adj[ch]:
            indeg[nxt] -= 1
            if indeg[nxt] == 0:
                q.append(nxt)
    return "".join(out) if len(out) == len(adj) else ""
```

**Complexity:** O(C) to build edges (C = total characters) plus O(U + E) for Kahn's over U unique letters.

**Edge cases:** the prefix rule; duplicate edges inflating indegree (the `set` prevents it); letters appearing only in uninformative positions must still be output; adjacent pairs suffice because sortedness is transitive.

</details>

### Template 5: 3-color DFS (directed cycle detection, DFS topo sort)

> - **Example problem:** [Find Eventual Safe States (802)](https://leetcode.com/problems/find-eventual-safe-states/)
> - **Why this template:** a node is "safe" exactly when no path from it reaches a cycle, which GRAY/BLACK track.
> - **Worked below:** Example 5.1
> - **Also fits:** [Course Schedule (207)](https://leetcode.com/problems/course-schedule/)

```python
def has_cycle_directed(n: int, adj: list[list[int]]) -> bool:
    """3-color DFS. WHITE=unvisited, GRAY=on current path, BLACK=finished.

    A back edge to a GRAY node closes a cycle. BLACK nodes are safe to skip.
    """
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * n

    def dfs(u: int) -> bool:
        color[u] = GRAY
        for v in adj[u]:
            if color[v] == GRAY:
                return True
            if color[v] == WHITE and dfs(v):
                return True
        color[u] = BLACK                      # TODO: append u here for DFS-postorder topo sort
        return False

    return any(color[u] == WHITE and dfs(u) for u in range(n))
```

Two colors fail on directed graphs: a cross edge (`a -> b`, `a -> c -> b`) reaches a visited node without a cycle. Only an edge into the **current path** (GRAY) closes one. Topological order is reversed postorder:

```python
def topo_sort_dfs(n: int, edges: list[list[int]]) -> list[int]:
    """Topological order via reversed DFS postorder; [] on cycle."""
    adj = build_adj(n, edges, directed=True)
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * n
    post: list[int] = []

    def dfs(u: int) -> bool:
        color[u] = GRAY
        for v in adj[u]:
            if color[v] == GRAY or (color[v] == WHITE and not dfs(v)):
                return False
        color[u] = BLACK
        post.append(u)                        # u finishes after all its descendants
        return True

    for u in range(n):
        if color[u] == WHITE and not dfs(u):
            return []
    return post[::-1]
```

**Undirected** cycle detection differs: each edge appears twice, so skip the edge to the parent once; any other visited neighbor is a cycle. (Or DSU: an edge joining already-connected endpoints closes a cycle.)

```python
def has_cycle_undirected(n: int, edges: list[list[int]]) -> bool:
    """Iterative DFS tracking the parent; seeing a visited non-parent means a cycle."""
    adj = build_adj(n, edges)
    visited = [False] * n
    for s in range(n):
        if visited[s]:
            continue
        visited[s] = True
        stack = [(s, -1)]
        while stack:
            u, parent = stack.pop()
            for v in adj[u]:
                if v == parent:
                    parent = -2               # skip the tree edge once (handles parallel edges)
                    continue
                if visited[v]:
                    return True
                visited[v] = True
                stack.append((v, u))
    return False
```

#### Example 5.1: Nodes that can't reach a cycle [Medium]

**Problem:** [Find Eventual Safe States (802)](https://leetcode.com/problems/find-eventual-safe-states/) · **Template:** Template 5 (3-color DFS (directed cycle detection, DFS topo sort))

Problem: a directed graph as an adjacency list `graph[u]`. A node is safe if every path from it ends at a terminal node (no outgoing edges). Return all safe nodes in increasing order. Example: `[[1, 2], [2, 3], [5], [0], [5], [], []]` → `[2, 4, 5, 6]` (0 → 1 → 3 → 0 is a cycle).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "every path ends" = no path reaches a cycle = DFS from the node never hits a GRAY (on-stack) node.

**Approach:** 3-color DFS where colors double as a memo. BLACK = proven safe (all neighbors safe). A node that reaches a GRAY node, or any unsafe node, returns `False` and stays GRAY forever, so later visits see GRAY and correctly treat it as unsafe. WHITE nodes are explored once, so the whole graph costs one pass.

```python
def eventual_safe_nodes(graph: list[list[int]]) -> list[int]:
    """Find Eventual Safe States: sorted nodes from which every path reaches a terminal node."""
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * len(graph)

    def safe(u: int) -> bool:
        if color[u] != WHITE:
            return color[u] == BLACK          # GRAY: on the current path, or proven unsafe
        color[u] = GRAY
        if any(not safe(v) for v in graph[u]):
            return False                      # stays GRAY forever = unsafe
        color[u] = BLACK
        return True

    return [u for u in range(len(graph)) if safe(u)]
```

**Complexity:** O(V + E) time: each node is colored once and each edge scanned at most once. O(V) space.

**Edge cases:** self-loop (`[[0]]`) → unsafe; terminal nodes are safe (`any` over no neighbors is `False`); recursion depth up to V (raise the limit, or use the alternative). Alternative: Kahn's on the **reversed** graph. Terminal nodes have out-degree 0; peel them and decrement the out-degree of their predecessors. Nodes that reach out-degree 0 are safe; cycle nodes and anything that leads to them never do.

</details>

### Template 6: state-space BFS

> - **Example problem:** [Open the Lock (752)](https://leetcode.com/problems/open-the-lock/)
> - **Why this template:** "minimum turns" over combinations generated on the fly (8 neighbors each, dead ends blocked) is BFS over states, not a given graph.
> - **Worked below:** Example 6.1
> - **Also fits:** [Word Ladder (127)](https://leetcode.com/problems/word-ladder/)

```python
def state_space_bfs(start: Hashable, is_goal: Callable[[Hashable], bool],
                    neighbors: Callable[[Hashable], Iterable[Hashable]],
                    blocked: set | frozenset = frozenset()) -> int:
    """Fewest moves from start to a goal state, where states are generated on the fly."""
    if start in blocked:
        return -1
    if is_goal(start):
        return 0
    seen = {start} | set(blocked)             # treat blocked states as pre-visited
    q = deque([(start, 0)])
    while q:
        state, d = q.popleft()
        for nxt in neighbors(state):          # TODO: the move generator is the whole problem
            if nxt not in seen:
                if is_goal(nxt):
                    return d + 1
                seen.add(nxt)
                q.append((nxt, d + 1))
    return -1
```

Knobs: **state encoding** (hashable: string, tuple, int bitmask); the move generator; what goes into `seen`. The state must include *everything affecting future moves*: if two visits to a cell can have different futures (keys held, budget left), the cell alone isn't the state.

#### Example 6.1: Shortest path over generated states [Medium/Hard]

**Problems:** [Open the Lock (752)](https://leetcode.com/problems/open-the-lock/), [Word Ladder (127)](https://leetcode.com/problems/word-ladder/) · **Template:** Template 6 (state-space BFS)

Problem: a 4-wheel lock (digits 0-9) starts at `"0000"`; a move turns one wheel up or down by one, wrapping. Never land on a dead end. Fewest moves to the target, or -1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** fewest moves, unit cost, small string states → BFS on an implicit graph of 10^4 states, 8 neighbors each.

**Approach:** Template 6; only the neighbor generator is problem-specific. Dead ends go into `seen` up front so they're never expanded (if `"0000"` is one, the answer is -1).

```python
def open_lock(deadends: list[str], target: str) -> int:
    """Open the Lock: 4 wheels, each move turns one wheel by +/-1. -1 if impossible."""
    def neighbors(s: str):
        for i in range(4):
            d = int(s[i])
            for step in (1, -1):
                yield s[:i] + str((d + step) % 10) + s[i + 1:]

    return state_space_bfs("0000", lambda s: s == target, neighbors, set(deadends))
```

**Word Ladder** has the same shape: words differing by one letter are neighbors. All-pairs comparison is O(n^2 * L); instead bucket words by wildcard pattern (`"h*t"` holds `hot`, `hit`, `hat`). Clear a used bucket so it's never rescanned.

```python
def ladder_length(begin: str, end: str, word_list: list[str]) -> int:
    """Word Ladder: number of WORDS in the shortest transformation sequence, 0 if none.

    Wildcard buckets ("h*t") connect words differing in one position without O(n^2) comparisons.
    """
    words = set(word_list)
    if end not in words:
        return 0
    L = len(begin)
    buckets: dict[str, list[str]] = defaultdict(list)
    for w in words | {begin}:
        for i in range(L):
            buckets[w[:i] + "*" + w[i + 1:]].append(w)
    seen = {begin}
    q = deque([(begin, 1)])
    while q:
        w, d = q.popleft()
        for i in range(L):
            key = w[:i] + "*" + w[i + 1:]
            for nxt in buckets[key]:
                if nxt == end:
                    return d + 1
                if nxt not in seen:
                    seen.add(nxt)
                    q.append((nxt, d + 1))
            buckets[key] = []                 # bucket fully expanded; never scan it again
    return 0
```

**Complexity:** Open the Lock is at most O(10^4 * 8 * 4). Word Ladder is O(n * L^2) to build keys (slicing costs L).

**Edge cases:** start is a dead end; start equals target (0); Word Ladder counts **words**, not edges (start at 1); end word missing from the list.

</details>

#### Example 6.2: BFS where position is not enough [Hard]

**Problem:** [Shortest Path in a Grid with Obstacles Elimination (1293)](https://leetcode.com/problems/shortest-path-in-a-grid-with-obstacles-elimination/) · **Template:** Template 6 (state-space BFS)

Problem: grid of 0 (open) and 1 (obstacle). Move 4-directionally from top-left to bottom-right, removing up to k obstacles. Fewest steps, or -1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** BFS for fewest steps, but arrivals at a cell with different remaining budgets have different futures → state = `(r, c, remaining)`.

**Approach:** a `seen` set of triples works (O(R*C*k) states). Tighter: keep the **best remaining budget** per cell; BFS arrives in nondecreasing distance, so a later arrival helps only with more budget. If `k >= rows + cols - 2`, bulldoze a straight Manhattan path.

```python
def shortest_path_with_eliminations(grid: list[list[int]], k: int) -> int:
    """Shortest Path in a Grid with Obstacles Elimination.

    State = (r, c, eliminations_left). Visiting a cell with more budget left
    is strictly better, so we track the best remaining budget per cell.
    """
    rows, cols = len(grid), len(grid[0])
    if k >= rows + cols - 2:                  # enough budget to walk any Manhattan path
        return rows + cols - 2
    best = [[-1] * cols for _ in range(rows)]  # max eliminations left seen at each cell
    best[0][0] = k
    q = deque([(0, 0, k, 0)])
    while q:
        r, c, left, d = q.popleft()
        if (r, c) == (rows - 1, cols - 1):
            return d
        for nr, nc in grid_neighbors(r, c, rows, cols):
            nleft = left - grid[nr][nc]
            if nleft > best[nr][nc]:          # dominates any earlier visit of this cell
                best[nr][nc] = nleft
                q.append((nr, nc, nleft, d + 1))
    return -1
```

**Complexity:** O(R*C*k) time and space worst case (each cell's best budget improves at most k+1 times).

**Edge cases:** 1x1 grid → 0; negative `nleft` is filtered because `best` starts at -1; the shortcut handles huge k.

Same idea with a bitmask: **Shortest Path to Get All Keys**, state = `(r, c, keys_mask)`; a door is passable only if its key bit is set.

```python
def shortest_path_all_keys(grid: list[str]) -> int:
    """Shortest Path to Get All Keys: state = (r, c, keys_bitmask)."""
    rows, cols = len(grid), len(grid[0])
    all_keys = 0
    for r in range(rows):
        for c in range(cols):
            ch = grid[r][c]
            if ch == "@":
                sr, sc = r, c
            elif ch.islower():
                all_keys |= 1 << (ord(ch) - ord("a"))
    seen = {(sr, sc, 0)}
    q = deque([(sr, sc, 0, 0)])
    while q:
        r, c, keys, d = q.popleft()
        if keys == all_keys:
            return d
        for nr, nc in grid_neighbors(r, c, rows, cols):
            ch = grid[nr][nc]
            if ch == "#":
                continue
            if ch.isupper() and not keys >> (ord(ch.lower()) - ord("a")) & 1:
                continue                      # locked door, no key yet
            nkeys = keys | (1 << (ord(ch) - ord("a"))) if ch.islower() else keys
            if (nr, nc, nkeys) not in seen:
                seen.add((nr, nc, nkeys))
                q.append((nr, nc, nkeys, d + 1))
    return -1
```

</details>

## Template cheat sheet

All templates, for warm-ups and review.

### Template 1: BFS shortest path (unweighted)

```python
def bfs_shortest_path(adj: list[list[int]], src: int, dst: int) -> int:
    """Fewest edges from src to dst in an unweighted graph, -1 if unreachable.

    Mark visited when ENQUEUING: each node enters the queue once, and the first
    time it is reached is via a shortest path (BFS explores by layers).
    """
    if src == dst:
        return 0
    dist = {src: 0}
    q = deque([src])
    while q:
        u = q.popleft()
        for v in adj[u]:                      # TODO: neighbor generator
            if v not in dist:                 # TODO: extra filters (walls, constraints)
                dist[v] = dist[u] + 1
                if v == dst:                  # early exit is safe on enqueue for BFS
                    return dist[v]
                q.append(v)
    return -1
```

### Template 2: multi-source BFS

```python
def multi_source_bfs(grid: list[list[int]], is_source: Callable[[int], bool],
                     passable: Callable[[int], bool]) -> list[list[int]]:
    """Distance from every cell to its nearest source (-1 if unreachable).

    Equivalent to adding a virtual super-source with 0-weight edges to all sources.
    """
    rows, cols = len(grid), len(grid[0])
    dist = [[-1] * cols for _ in range(rows)]
    q: deque[tuple[int, int]] = deque()
    for r in range(rows):
        for c in range(cols):
            if is_source(grid[r][c]):         # TODO: what counts as a source
                dist[r][c] = 0
                q.append((r, c))
    while q:
        r, c = q.popleft()
        for nr, nc in grid_neighbors(r, c, rows, cols):
            if dist[nr][nc] == -1 and passable(grid[nr][nc]):   # TODO: passability
                dist[nr][nc] = dist[r][c] + 1
                q.append((nr, nc))
    return dist
```

### Template 3: grid DFS / flood fill

```python
def grid_dfs_iterative(grid: list[list[str]], r: int, c: int, target: str, mark: str) -> int:
    """Flood fill from (r, c) over cells equal to `target`, overwriting with `mark`.

    Returns component size. Iterative to avoid Python's recursion limit on big grids.
    Marking on push guarantees each cell is pushed once.
    """
    rows, cols = len(grid), len(grid[0])
    if grid[r][c] != target:
        return 0
    grid[r][c] = mark                         # in-place visited marker (TODO: or a visited set)
    stack = [(r, c)]
    size = 0
    while stack:
        cr, cc = stack.pop()
        size += 1
        for nr, nc in grid_neighbors(cr, cc, rows, cols):   # TODO: DIRS8 for diagonal adjacency
            if grid[nr][nc] == target:
                grid[nr][nc] = mark
                stack.append((nr, nc))
    return size
```

### Template 4: Kahn's topological sort

```python
def topo_sort_kahn(n: int, edges: list[list[int]]) -> list[int]:
    """Kahn's algorithm. Edge [u, v] means u must come before v.

    Returns an order, or [] if a cycle exists (some nodes never reach indegree 0).
    """
    adj: list[list[int]] = [[] for _ in range(n)]
    indeg = [0] * n
    for u, v in edges:
        adj[u].append(v)
        indeg[v] += 1
    q = deque(i for i in range(n) if indeg[i] == 0)   # TODO: heap here for lexicographically smallest order
    order: list[int] = []
    while q:
        u = q.popleft()
        order.append(u)
        for v in adj[u]:
            indeg[v] -= 1
            if indeg[v] == 0:
                q.append(v)
    return order if len(order) == n else []
```

### Template 5: 3-color DFS (directed cycle detection, DFS topo sort)

```python
def has_cycle_directed(n: int, adj: list[list[int]]) -> bool:
    """3-color DFS. WHITE=unvisited, GRAY=on current path, BLACK=finished.

    A back edge to a GRAY node closes a cycle. BLACK nodes are safe to skip.
    """
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * n

    def dfs(u: int) -> bool:
        color[u] = GRAY
        for v in adj[u]:
            if color[v] == GRAY:
                return True
            if color[v] == WHITE and dfs(v):
                return True
        color[u] = BLACK                      # TODO: append u here for DFS-postorder topo sort
        return False

    return any(color[u] == WHITE and dfs(u) for u in range(n))
```

```python
def topo_sort_dfs(n: int, edges: list[list[int]]) -> list[int]:
    """Topological order via reversed DFS postorder; [] on cycle."""
    adj = build_adj(n, edges, directed=True)
    WHITE, GRAY, BLACK = 0, 1, 2
    color = [WHITE] * n
    post: list[int] = []

    def dfs(u: int) -> bool:
        color[u] = GRAY
        for v in adj[u]:
            if color[v] == GRAY or (color[v] == WHITE and not dfs(v)):
                return False
        color[u] = BLACK
        post.append(u)                        # u finishes after all its descendants
        return True

    for u in range(n):
        if color[u] == WHITE and not dfs(u):
            return []
    return post[::-1]
```

```python
def has_cycle_undirected(n: int, edges: list[list[int]]) -> bool:
    """Iterative DFS tracking the parent; seeing a visited non-parent means a cycle."""
    adj = build_adj(n, edges)
    visited = [False] * n
    for s in range(n):
        if visited[s]:
            continue
        visited[s] = True
        stack = [(s, -1)]
        while stack:
            u, parent = stack.pop()
            for v in adj[u]:
                if v == parent:
                    parent = -2               # skip the tree edge once (handles parallel edges)
                    continue
                if visited[v]:
                    return True
                visited[v] = True
                stack.append((v, u))
    return False
```

### Template 6: state-space BFS

```python
def state_space_bfs(start: Hashable, is_goal: Callable[[Hashable], bool],
                    neighbors: Callable[[Hashable], Iterable[Hashable]],
                    blocked: set | frozenset = frozenset()) -> int:
    """Fewest moves from start to a goal state, where states are generated on the fly."""
    if start in blocked:
        return -1
    if is_goal(start):
        return 0
    seen = {start} | set(blocked)             # treat blocked states as pre-visited
    q = deque([(start, 0)])
    while q:
        state, d = q.popleft()
        for nxt in neighbors(state):          # TODO: the move generator is the whole problem
            if nxt not in seen:
                if is_goal(nxt):
                    return d + 1
                seen.add(nxt)
                q.append((nxt, d + 1))
    return -1
```

## Variations and extensions

**Bipartite / two-coloring (LeetCode: Is Graph Bipartite?).** BFS or DFS with alternating colors; a conflict means an odd cycle. Loop over all nodes, since the graph may be disconnected.

```python
def is_bipartite(graph: list[list[int]]) -> bool:
    """Is Graph Bipartite: BFS two-coloring over every component."""
    color = [-1] * len(graph)
    for s in range(len(graph)):
        if color[s] != -1:
            continue
        color[s] = 0
        q = deque([s])
        while q:
            u = q.popleft()
            for v in graph[u]:
                if color[v] == -1:
                    color[v] = 1 - color[u]
                    q.append(v)
                elif color[v] == color[u]:    # edge inside one side -> odd cycle
                    return False
    return True
```

**More state-space graphs.** Minimum Genetic Mutation is Word Ladder over `ACGT` with a bank. Sliding Puzzle flattens the 2x3 board to a string and swaps `0` with precomputed neighbors; only 6!/2 = 360 states are reachable from any start.

```python
def min_mutation(start: str, end: str, bank: list[str]) -> int:
    """Minimum Genetic Mutation: one-letter changes, every intermediate must be in bank."""
    bank_set = set(bank)

    def neighbors(s: str):
        for i in range(len(s)):
            for ch in "ACGT":
                if ch != s[i]:
                    cand = s[:i] + ch + s[i + 1:]
                    if cand in bank_set:
                        yield cand

    return state_space_bfs(start, lambda s: s == end, neighbors)


def sliding_puzzle(board: list[list[int]]) -> int:
    """Sliding Puzzle (2x3): state = board flattened to a string."""
    start = "".join(str(x) for row in board for x in row)
    # Precomputed neighbors of each index in a 2x3 grid.
    moves = {0: (1, 3), 1: (0, 2, 4), 2: (1, 5), 3: (0, 4), 4: (1, 3, 5), 5: (2, 4)}

    def neighbors(s: str):
        z = s.index("0")
        for j in moves[z]:
            lst = list(s)
            lst[z], lst[j] = lst[j], lst[z]
            yield "".join(lst)

    return state_space_bfs(start, lambda s: s == "123450", neighbors)
```

**Bidirectional BFS.** With large branching factor b and a known goal, search from both ends, expanding the smaller frontier, until they meet: about O(b^(d/2)) instead of O(b^d). Mention it for Word Ladder; implement only if asked.

```python
def bidirectional_bfs(start: str, goal: str, neighbors: Callable[[str], Iterable[str]]) -> int:
    """Fewest moves from start to goal, expanding the smaller frontier each round."""
    if start == goal:
        return 0
    front, back = {start}, {goal}
    seen = {start, goal}
    d = 0
    while front and back:
        if len(front) > len(back):
            front, back = back, front         # always expand the smaller side
        d += 1
        nxt_front: set[str] = set()
        for s in front:
            for nxt in neighbors(s):
                if nxt in back:
                    return d
                if nxt not in seen:
                    seen.add(nxt)
                    nxt_front.add(nxt)
        front = nxt_front
    return -1
```

**Clone Graph.** Traverse while building an `old -> new` map; it doubles as the visited set, so cycles terminate. Same pattern as copying a linked list with random pointers (`09-linked-lists.md`).

```python
class Node:
    def __init__(self, val: int = 0, neighbors: list["Node"] | None = None):
        self.val = val
        self.neighbors = neighbors if neighbors is not None else []


def clone_graph(node: Node | None) -> Node | None:
    """Clone Graph: old->new map doubles as the visited set (handles cycles)."""
    if node is None:
        return None
    clones = {node: Node(node.val)}
    q = deque([node])
    while q:
        u = q.popleft()
        for v in u.neighbors:
            if v not in clones:
                clones[v] = Node(v.val)
                q.append(v)
            clones[u].neighbors.append(clones[v])
    return clones[node]
```

Other combinations and follow-ups:

- **Weighted edges** → Dijkstra; **0/1 weights** → 0-1 BFS with a deque (both in `14-union-find-shortest-paths-mst.md`).
- **Eventual safe states**: 3-color DFS or reverse-graph Kahn's (Example 5.1).
- **Minimum height trees**: peel leaves layer by layer (Kahn's on an undirected tree, degree 1 = ready). The last 1-2 nodes are the centers.
- **Graph valid tree**: n-1 edges and connected (or acyclic and connected).
- "Graph too big for memory?" Keep only the BFS frontier; for state spaces, A* with an admissible heuristic or bidirectional search.
- "Return the path, not the length": set `parent[v] = u` on enqueue and walk back from the target.
- "All shortest paths" (Word Ladder II): BFS builds a layered DAG of parents, then backtracking enumerates paths.

## Common pitfalls

- Marking visited on dequeue in BFS. Mark on enqueue.
- Reversed prerequisite edges (the most common topo-sort bug). Write "edge u -> v means u first" before coding.
- Forgetting isolated nodes: loop over all n nodes for components, topo sort, and bipartite checks.
- Recursive DFS on a 10^5-node path or a large grid: `RecursionError`. Use an explicit stack (`sys.setrecursionlimit` still risks a C stack overflow).
- 2-color visited for directed cycles: cross edges look like cycles. Use 3 colors.
- Undirected cycle detection that counts the edge back to the parent as a cycle.
- State-space BFS: an incomplete state (keys/budget missing from `seen`) gives wrong answers; an over-detailed state (including step count) never terminates.
- Building all O(n^2) pairwise edges for Word Ladder. Generate neighbors (or wildcard buckets).
- Layered BFS: `range(len(q))` must snapshot the layer size before the inner loop mutates the queue (Python evaluates it once, as desired).

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| BFS shortest path | Dijkstra / 0-1 BFS (`14`) | edges have different costs |
| DFS components | Union-Find (`14`) | edges arrive over time with interleaved queries ("after each addition, how many groups?") |
| Grid traversal | DP (`16`) | moves only go right/down (no cycles); you need min cost or path count |
| State-space BFS | Backtracking (`15`) | you need *all* solutions or any valid configuration, not fewest moves |
| Topological sort | Interval scheduling (`12`) | constraints are time overlap, not "A before B" |
| Tree traversal (`10`) | Graph traversal | trees have no cycles and one parent; graphs need `visited` |
| Bipartite | General k-coloring | k ≥ 3 is NP-hard; only 2-coloring is a BFS |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Number of Islands (200)](https://leetcode.com/problems/number-of-islands/) · Medium · Grid DFS/BFS: count flood-fill launches.
2. [Rotting Oranges (994)](https://leetcode.com/problems/rotting-oranges/) · Medium · Multi-source BFS: layered BFS from all rotten oranges.
3. [Course Schedule (207)](https://leetcode.com/problems/course-schedule/) · Medium · Cycle detection / topo sort: Kahn's processed count == n.
4. [Course Schedule II (210)](https://leetcode.com/problems/course-schedule-ii/) · Medium · Topological order (Kahn): return Kahn's order.
5. [Word Ladder (127)](https://leetcode.com/problems/word-ladder/) · Hard · State-space BFS: wildcard buckets; count words, not edges.

### Optional

**Grid traversal / components**
- [Flood Fill (733)](https://leetcode.com/problems/flood-fill/) · Easy · Return early if the new color equals the old.
- [Max Area of Island (695)](https://leetcode.com/problems/max-area-of-island/) · Medium · Return size from DFS.
- [Surrounded Regions (130)](https://leetcode.com/problems/surrounded-regions/) · Medium · Flood from the border first.
- [Number of Enclaves (1020)](https://leetcode.com/problems/number-of-enclaves/) · Medium · Border flood, count remaining land.
- [Pacific Atlantic Water Flow (417)](https://leetcode.com/problems/pacific-atlantic-water-flow/) · Medium · Two multi-source searches uphill, intersect.
- [Number of Connected Components in an Undirected Graph (323)](https://leetcode.com/problems/number-of-connected-components-in-an-undirected-graph/) · Medium · Premium · Adjacency list + DFS (or DSU).

**BFS / multi-source**
- [Shortest Path in Binary Matrix (1091)](https://leetcode.com/problems/shortest-path-in-binary-matrix/) · Medium · 8 directions, BFS from (0,0).
- [Walls and Gates (286)](https://leetcode.com/problems/walls-and-gates/) · Medium · Premium · BFS from all gates.
- [01 Matrix (542)](https://leetcode.com/problems/01-matrix/) · Medium · Sources are zeros.
- [As Far from Land as Possible (1162)](https://leetcode.com/problems/as-far-from-land-as-possible/) · Medium · Sources are land; answer is the max layer.

**Topological sort / cycles**
- [Find Eventual Safe States (802)](https://leetcode.com/problems/find-eventual-safe-states/) · Medium · 3-color DFS or reverse-graph Kahn's.
- [Minimum Height Trees (310)](https://leetcode.com/problems/minimum-height-trees/) · Medium · Leaf peeling.
- [Alien Dictionary (269)](https://leetcode.com/problems/alien-dictionary/) · Hard · Premium · Extract edges from adjacent words; prefix rule.

**Coloring / cloning**
- [Is Graph Bipartite? (785)](https://leetcode.com/problems/is-graph-bipartite/) · Medium · BFS two-coloring over all components.
- [Possible Bipartition (886)](https://leetcode.com/problems/possible-bipartition/) · Medium · Build graph from dislikes, then 2-color.
- [Clone Graph (133)](https://leetcode.com/problems/clone-graph/) · Medium · old->new map is the visited set.

**State-space BFS**
- [Open the Lock (752)](https://leetcode.com/problems/open-the-lock/) · Medium · 8 neighbors per state; dead ends pre-visited.
- [Minimum Genetic Mutation (433)](https://leetcode.com/problems/minimum-genetic-mutation/) · Medium · Word Ladder on ACGT.
- [Sliding Puzzle (773)](https://leetcode.com/problems/sliding-puzzle/) · Hard · Board as string, swap the zero.
- [Shortest Path in a Grid with Obstacles Elimination (1293)](https://leetcode.com/problems/shortest-path-in-a-grid-with-obstacles-elimination/) · Hard · State includes remaining budget.
- [Shortest Path to Get All Keys (864)](https://leetcode.com/problems/shortest-path-to-get-all-keys/) · Hard · State includes key bitmask.

## Self-check

1. Why must BFS mark visited on enqueue rather than dequeue? What goes wrong otherwise?
2. Why is multi-source BFS correct for "distance to the nearest source", and why is it faster than one BFS per source?
3. In a directed graph, why doesn't reaching a visited node imply a cycle? How do three colors fix it?
4. In Kahn's algorithm, why are cycle nodes never output? How do you detect a unique topological order?
5. In Alien Dictionary, why does only the first differing character of an adjacent pair carry information, and what input is invalid without a cycle?
6. Given a grid problem with keys, doors, and a teleporter, how do you decide the BFS state?
7. When would you choose iterative over recursive DFS in Python, and when doesn't it matter?
