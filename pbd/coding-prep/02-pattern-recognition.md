# 02. Pattern Recognition: Mapping Any Problem to a Pattern

You needn't have seen a problem; you need to see **what it is beneath the story**. This module maps signals → patterns, lists disguises, and separates look-alikes. Pair it with `01-problem-solving-framework.md` (the procedure) and `22-disguised-problem-drills.md` (the practice).

---

## 1. The three lenses

Every problem signals its pattern three ways. Check all three.

### Lens A: input shape

| Input | Default candidates |
|---|---|
| Unsorted array / string | Hash map, prefix sums, sliding window, sort first, monotonic stack |
| Sorted (or sortable without losing meaning) array | Two pointers, binary search, greedy |
| Contiguous subarray / substring | Sliding window, prefix sums, monotonic deque, Kadane-style DP |
| Linked list | Dummy head, fast/slow pointers, in-place reversal |
| Tree | DFS (top-down or bottom-up), BFS by level |
| Graph (edges, relations, dependencies) | BFS/DFS, topological sort, union-find, Dijkstra |
| Grid / matrix | Grid BFS/DFS, grid DP, multi-source BFS |
| Intervals / time ranges | Sort + merge, sweep line, heap of end times |
| Many strings / prefixes | Trie, hashing, sorting |
| Stream / online queries | Heap, two heaps, deque, hash map + running state |
| Small n (≤ 20) | Backtracking, bitmask DP |
| Values in 1..n | Index-as-hash (cyclic sort), counting array |

### Lens B: question type and trigger phrases

| The problem asks for… (typical phrasing) | Pattern |
|---|---|
| Longest / shortest **contiguous** segment satisfying X ("consecutive", "contiguous", "streak", "window", "in a row") | Sliding window if validity is monotonic; else prefix sums + hash map |
| The same, with a budget ("at most k replacements / flips / distinct") | Sliding window with a budget |
| **Count** subarrays with sum / property = k ("sum equals k", "balanced" 0s/1s, "exactly k") | Prefix sums + hash map; exactly k = atMost(k) − atMost(k−1) |
| Pair / triplet meeting a target; input "sorted" / "non-decreasing" | Hash map (unsorted) or two pointers (sorted); sorted input also suggests binary search |
| **k-th** largest / smallest; top k ("top", "most frequent", "closest k") | Size-k heap, quickselect, bucket, binary search on value |
| Running median / balance two halves ("median of a stream") | Two heaps |
| Merge k sorted sequences | Heap (k-way merge) |
| **Next / previous greater or smaller** ("next warmer", "next taller", "how many days until", "span") | Monotonic stack |
| Largest rectangle, water trapped | Monotonic stack / two pointers |
| Max/min of every size-k window | Monotonic deque |
| **Minimize the maximum** / maximize the minimum; smallest X satisfying a monotonic condition ("minimum capacity / speed / days / time such that") | Binary search on the answer |
| Search a "rotated sorted" array | Modified binary search |
| **All** combinations / permutations / arrangements ("generate all", "list every", "all possible") | Backtracking |
| **Number of ways**, min/max cost, or "is it possible" over choices ("how many distinct ways", "maximum profit", "can you reach") | DP (often mod 10^9+7); greedy if provable; BFS for reachability |
| A "subsequence" (not contiguous) | DP (LCS/LIS family); two pointers for is-subsequence |
| Palindromes | Expand around center / DP / two pointers |
| Fewest steps / shortest **unweighted** path ("minimum number of moves / steps / transformations") | BFS; state-space BFS when nodes are configurations |
| Simultaneous spread ("spreads to neighbors each minute") | Multi-source BFS |
| Regions in a grid ("islands", "regions", "connected areas") | Grid DFS/BFS or union-find |
| Shortest path, **weighted (non-negative)** ("cheapest", "fastest route", "latency") | Dijkstra |
| Cheapest path with at most k edges ("at most k stops / hops") | Bellman-Ford (k rounds) / BFS with (node, stops) state |
| All-pairs distances, small n ("every pair", "all distances") | O(n^2) / Floyd–Warshall |
| Valid **order** under dependencies ("prerequisites", "depends on", "build order", "must happen before") | Topological sort |
| **Connectivity** / grouping / merges over time ("friend circles", "linked accounts", "same group") | Union-find (BFS/DFS if static) |
| Connect everything at minimum cost | MST (Kruskal / Prim) |
| Merge / overlap / free time ("overlapping", "meetings", "bookings", "maintenance windows") | Intervals: sort + merge, sweep line |
| Max concurrent ("how many rooms / servers at once") | Sweep line / heap of end times |
| Max non-overlapping / min removals | Greedy (sort by end) |
| Prefix matching / many-word search ("prefix", "autocomplete", "starts with", "dictionary of words") | Trie |
| Pattern in text / repeated substring | KMP, rolling hash |
| Single / missing / duplicate number, toggles, subsets as masks ("appears once while others appear twice", "values from 1 to n") | XOR, cyclic sort / index marking |
| "Return any valid answer" / locally optimal choices | Greedy or constructive |
| Design / simulate a component ("implement a cache / limiter / scheduler") | Hash map + linked list / heap / deque (module 21) |

### Lens C: constraints

Input size sets the target complexity (n ≤ 20 → bitmask, ≤ 500 → O(n^3), ≤ 10^6 → O(n log n); full table with patterns in `03-python-toolkit-and-complexity.md` §2). Other constraints are just as telling:

| Constraint | What it suggests |
|---|---|
| Values up to 10^9, n small | Binary search on value; coordinate compression |
| All values non-negative | Sliding window works for sum constraints |
| Negatives allowed | Sliding window on sums breaks → prefix sums / DP |
| k small (≤ 10) | k-dimensional DP state, k-way heap, bitmask over k |
| "Answer modulo 10^9+7" | Counting DP / combinatorics |
| Sorted input | Two pointers / binary search (use the sortedness) |
| Grid up to 10^3 × 10^3 | O(R·C) BFS/DFS/DP only |

---

## 2. Master decision flow

```
Is the input a graph/grid/tree, or can it be modeled as one (relations, states, moves)?
├─ Tree → DFS (return values up, or pass state down) / BFS by level
├─ Need ordering under dependencies → Topological sort
├─ Connectivity / grouping (especially as edges arrive) → Union-Find
├─ Fewest steps, unweighted → BFS (multi-source if many starts; state-space BFS if nodes are configurations)
├─ Weighted shortest path → Dijkstra (non-negative) / Bellman-Ford (negative or ≤ k edges)
└─ Connect all at min cost → MST
Else: is it about intervals/time ranges?
├─ Merge/overlap → sort by start
├─ Max concurrency → sweep line / heap of ends
└─ Max non-overlapping → sort by end (greedy)
Else: does it ask for ALL solutions?  → Backtracking (n small)
Else: optimize/count over a sequence of CHOICES with repeated subproblems? → DP
Else: contiguous subarray/substring?
├─ Validity monotonic (shrinking helps) → Sliding window
├─ Sum = k with negatives / count → Prefix sums + hash map
└─ Window max/min → Monotonic deque
Else: "minimize the max" / "smallest X that works"? → Binary search on answer
Else: nearest greater/smaller, span, histogram? → Monotonic stack
Else: k-th / top-k / streaming extremes? → Heap (two heaps for median)
Else: sorted (or sortable) + pairs/targets? → Two pointers / binary search
Else: lookups, counts, grouping? → Hash map / counting
Else: prefixes across many strings? → Trie
Else: locally optimal choice provably safe? → Greedy (verify with exchange argument)
```

---

## 3. Disguise techniques (and how to see through them)

Interviewers rarely say "use a sliding window." Common ways they hide the pattern:

1. **Story wrapping.** "Servers," "trucks," "students," "GPUs." **Counter:** first rewrite it as abstract input → output in one sentence.
2. **Renamed data structures.**
   - "Flights between cities," "people who trust each other," "services that call services" → **graph**
   - "Map of land and water," "warehouse floor" → **grid**
   - "Bookings," "shifts," "maintenance windows" → **intervals**
   - "Org chart," "file system" → **tree**
3. **Hidden contiguity.** "Consecutive days," "an uninterrupted run of requests," "a stretch of road" → **subarray** → window or prefix sums.
4. **Hidden monotonicity.** "Find the minimum speed / capacity / time so that…": if a bigger value works whenever a smaller one does → **binary search on the answer**.
5. **Implicit graphs.** "Minimum moves to transform A into B," "lock combinations," "puzzle states" → **nodes are states**, edges are single moves → **BFS**.
6. **Inverted or complementary questions.** "Subarrays containing at least one X" = total − (subarrays with no X). "Exactly k" = atMost(k) − atMost(k−1). "Longest subarray after deleting one element" = longest window with at most one zero.
7. **Dimension reduction.** "Largest rectangle of 1s in a matrix" = per-row histogram + monotonic stack. "Submatrix sum = k" = fix a pair of rows → 1D prefix sums.
8. **Circular structures.** "Circular route / ring buffer" → double the array or use modulo indices; circular max-subarray = max(normal max-subarray, total − min-subarray).
9. **Order that doesn't matter.** "Pick any subset of jobs to run" → **sort first** → greedy / two pointers / DP over sorted items.
10. **Value-range leaks.** "IDs are 1..n" → index-as-hash. "Lowercase letters only" → counting array of 26. "k ≤ 16" → bitmask.
11. **"Minimum operations."** BFS (operations are steps in a state graph), DP (choices with overlapping subproblems), or greedy (provably optimal local choice); the state space and constraints decide.
12. **Game / turn-taking.** "Two players pick optimally" → DP on (state, whose turn) or minimax on intervals.
13. **Disguised sorting + sweep.** "How many trains need a platform at peak" → events (+1 arrive, −1 depart) sorted by time.
14. **Disguised knapsack.** "Split into two groups with equal total," "assign + or − signs to hit a target" → subset sum (0/1 knapsack).
15. **Disguised topological sort.** "Infer the alphabet order from sorted words," "build steps," "compile order" → edges from the constraints → topo sort.
16. **Disguised union-find.** "Equations like a == b, b != c: satisfiable?" → union the equalities, then check the inequalities.
17. **Red herrings.** Extra details (names, unneeded timestamps, colors). **Counter:** drop any detail whose removal doesn't change the answer.

---

## 4. Look-alikes: how to tell them apart

| Confusion | Decider |
|---|---|
| Sliding window vs prefix sums + hash map | Window needs monotonic validity (e.g., non-negative values, so growing only increases the sum). Negatives or an "exactly = k" count → prefix sums |
| Two pointers vs hash map (pair sum) | Sorted, or sortable (original indices not needed) → two pointers, O(1) space. Unsorted and original indices needed → hash map |
| BFS vs DFS | Shortest/fewest steps → BFS. Explore everything / components / paths with backtracking → DFS |
| BFS vs Dijkstra | Equal edge costs → BFS. Varied non-negative weights → Dijkstra. Weights 0/1 → 0-1 BFS |
| Dijkstra vs Bellman-Ford | Negative edges, or an edge-count limit → Bellman-Ford / BFS with state |
| Union-find vs BFS/DFS | Edges arrive over time / many connectivity queries → union-find. One-shot static traversal → either |
| Greedy vs DP | Try greedy and hunt for a counterexample on small inputs. If one exists, or choices interact later → DP |
| Backtracking vs DP | All solutions → backtracking. A count/optimum with repeating subproblems → DP (memoize the recursion) |
| Heap vs sort vs quickselect | Static data, all sorted → sort. Top-k / streaming → heap, O(n log k). One-off k-th → quickselect, average O(n) |
| Monotonic stack vs heap | Nearest greater/smaller in sequence order → stack. Global extremes of a changing set → heap |
| Binary search on index vs on answer | Sorted array → index. Numeric answer with a monotonic feasibility check → answer |
| Intervals sort by start vs by end | Merging/union → by start. Max non-overlapping / min removals → by end |
| Trie vs hash set | Prefix queries or many patterns sharing prefixes → trie. Exact-match lookups only → hash set |
| DP on subsequence vs two pointers | "Is s a subsequence of t" → two pointers. "Longest common / count of subsequences" → DP |

---

## 5. The recognition drill (do this on every practice problem)

Before solving, write one line:

> **Signal:** "longest + contiguous + at most k distinct" → **Pattern:** variable sliding window with a Counter → **Template:** longest-window.

When wrong, log the **missed signal**, not the problem; each miss becomes a reusable rule. See `23-practice-plan-and-review.md`.

---

## 6. Quick quiz (answers below)

Name the pattern for each:

1. Fewest minutes until every server is infected, if infection spreads to adjacent racks each minute.
2. Smallest per-node CPU limit so that all pods fit on 5 nodes, packed in their given order.
3. For each day's latency, how many days until a higher latency occurs?
4. Number of ways to climb to the top of a deployment pipeline of n stages, taking 1 or 2 stages at a time.
5. Longest run of requests containing at most 2 distinct error codes.
6. Count of stretches of consecutive minutes where net connections changed by exactly 0.
7. Order in which to install packages given "A requires B" rules.
8. Are all data centers reachable, given links that come online one at a time?
9. Minimum number of conference rooms for a list of meetings.
10. All ways to assign 8 workloads to 8 nodes with no two conflicting.
11. Cheapest way to send a packet across a network with weighted links.
12. Cheapest route with at most 2 layovers.
13. Real-time p50 latency as requests arrive.
14. Find the one request ID that appears an odd number of times.
15. Suggest commands as the user types a prefix.

<details><summary>Answers</summary>

1. Multi-source BFS (all infected servers start in the queue).
2. Binary search on the answer + greedy feasibility check.
3. Monotonic stack (next greater element).
4. 1D DP (Fibonacci-like).
5. Variable sliding window with a Counter (at most k distinct).
6. Prefix sums + hash map (negatives allowed, exact-sum count).
7. Topological sort.
8. Union-find (dynamic connectivity).
9. Sweep line or a heap of end times (max overlap).
10. Backtracking (constraint satisfaction, like N-Queens).
11. Dijkstra.
12. Bellman-Ford with k+1 rounds / BFS with (node, stops) state.
13. Two heaps.
14. XOR of all IDs.
15. Trie.

</details>

---

## Self-check

1. Which three lenses do you apply, and what does each tell you?
2. Name five phrases that signal binary search on the answer or a sliding window.
3. Why does a sliding window fail for "subarray sum = k" with negatives, and what replaces it?
4. How do you recognize an implicit graph?
5. Give three examples of disguised union-find problems.
6. When do you sort intervals by start, and when by end?
7. What do you do when two patterns both seem to fit?
