# 11. Heaps, Top-K, and Merging

A heap cheaply and repeatedly answers "what is the current best item?": top-k over huge or streaming data, merging sorted sources, running medians, and simulations picking the next event or job. In infra terms: the scheduler queue, the log merger, the "top talkers" dashboard.

Companion code: `code/11_heaps_top_k.py` (templates and solutions, with asserts and randomized brute-force checks).

## At a glance

- **Use when:** "k-th / top k / k closest", a stream with small k, running median, "merge k sorted", or "which job runs next".
- **Templates:**
  - Template 1: size-k heap: top-k in O(k) memory
  - Template 2: two heaps: medians/percentiles
  - Template 3: k-way merge: one frontier element per sorted source
  - Template 4: event simulation: arrivals plus priority among the available
- **Typical complexity:** O(n log k) for top-k; O(N log k) for k-way merge; O(log n) per push/pop.
- **Top trap:** wrong heap direction (k largest needs a min-heap), or tuples comparing non-comparable payloads on ties.

## Core idea

A binary heap keeps the minimum at index 0: O(log n) push/pop, O(1) peek. It is **not** sorted, only parent <= children, which suffices to repeatedly extract the extreme. It beats sorting when you need only the top few or interleave inserts with extractions.

Python `heapq` facts to know cold:

| Operation | Call | Cost |
|---|---|---|
| Build from list | `heapq.heapify(a)` | O(n) |
| Push / pop min | `heappush(h, x)` / `heappop(h)` | O(log n) |
| Peek min | `h[0]` | O(1) |
| Pop then push (size stays k) | `heapreplace(h, x)` | O(log n) |
| Push then pop | `heappushpop(h, x)` | O(log n), faster if x would pop immediately |
| Top k | `nlargest(k, it, key=)` / `nsmallest` | O(n log k) |
| Max-heap | push `-x` (or `(-priority, ...)` tuples) | |

Four patterns and their invariants:

| Pattern | Invariant |
|---|---|
| Size-k heap | Holds the best k so far; the top is the **weakest** (the gatekeeper). For k largest, use a **min**-heap |
| Two heaps | `lo` (max-heap) holds the smaller half, `hi` (min-heap) the larger; `max(lo) <= min(hi)`, sizes differ by at most 1 |
| K-way merge | One frontier element per non-exhausted source; the top is the global next element |
| Event simulation | Holds what's **available now**, by priority; a sorted pointer admits items as the clock advances |

```
Size-k (k = 3, keep largest):     heap = [5, 11, 12]   top 5 = gatekeeper
  next 7:  7 > 5 -> replace       heap = [7, 11, 12]
  next 4:  4 <= 7 -> ignore

Two heaps:   lo (max) [1 2 3]  |  hi (min) [4 5]      median = top of lo
```

## Recognition signals

- "k-th largest/smallest", "top k", "k most frequent", "k closest": size-k heap (or quickselect, or bucket sort).
- Huge n or a **stream**, small k: size-k heap, O(n log k) time, O(k) memory.
- Stream or sliding-window "median"/"percentile": two heaps.
- "Merge k sorted ...", "smallest range covering one from each", "k-th smallest in a sorted matrix", "k pairs with smallest sums": k-way merge.
- "Minimum rooms/servers/machines", "how many concurrent at any time": heap of end times (or sweep, module 12).
- Jobs with arrival times and priorities, "which runs next", CPU scheduling: event simulation.
- "Up to k projects to maximize capital", "options unlock as resources grow": sorted unlock pointer + max-heap.
- "Repeatedly combine the two largest/smallest" (Last Stone Weight, connect ropes, Huffman): heap as greedy engine.
- "Remove an arbitrary element", "priority updates": lazy deletion.
- Dijkstra and Prim are heap-driven too (`14-union-find-shortest-paths-mst.md`).

## Disguises

1. **Top talkers.** "A router streams (src_ip, bytes) records, billions per hour. At any moment, report the 10 sources with the most bytes so far."
   Decode: aggregate in a hash map; size-k min-heap for the report. Too many keys for memory: Count-Min Sketch + heap (approximate).

2. **Latency p50 dashboard.** "Requests complete one by one. After each, show the median latency so far." Variant: "of the last 1,000 requests".
   Decode: running median = two heaps; "last 1,000" = sliding window median (two heaps + lazy deletion, or `bisect`).

3. **Log merge across shards.** "Each of 500 shards writes a log sorted by timestamp. Produce one time-ordered stream without loading everything into memory."
   Decode: k-way merge, one cursor per file, O(k) memory (`heapq.merge`, external sort's merge phase).

4. **Cluster sizing.** "Given batch jobs with start and end times, find the fewest machines so no job waits."
   Decode: maximum overlap = Meeting Rooms II: heap of end times; reuse the earliest-freed machine.

5. **Single-core job runner.** "Jobs arrive at known times with known durations. A free core picks the shortest waiting job (ties: lowest ID), or idles until the next arrival. Output the run order."
   Decode: event simulation (Single-Threaded CPU): sort by arrival, heap of available by `(duration, id)`, jump the clock when idle.

6. **Funding rounds.** "A startup can take at most k contracts. Each needs a minimum bankroll to start and adds its profit to the bankroll. Maximize the final bankroll."
   Decode: IPO: contracts unlock as capital grows (sorted pointer); greedily take the most profitable unlocked one (max-heap).

## Templates and worked examples

### Template 1: size-k heap (top-k)

> - **Example problem:** [Kth Largest Element in an Array (215)](https://leetcode.com/problems/kth-largest-element-in-an-array/)
> - **Why this template:** the kth largest needs only the k largest so far, the smallest of them atop a min-heap.
> - **Worked below:** Example 1.1
> - **Also fits:** [K Closest Points to Origin (973)](https://leetcode.com/problems/k-closest-points-to-origin/)

```python
def template_size_k_heap(items: Iterable[Any], k: int, key=lambda x: x) -> list[Any]:
    heap: list[tuple] = []
    for i, x in enumerate(items):
        entry = (key(x), i, x)  # i breaks ties so x itself is never compared
        if len(heap) < k:
            heapq.heappush(heap, entry)
        elif entry[0] > heap[0][0]:  # TODO: strictly better than the gatekeeper
            heapq.heapreplace(heap, entry)  # pop + push in one O(log k) step
    return [x for _, _, x in heap]  # unordered; sort if the caller needs order
```

Knobs:
- **Direction**: k largest uses a min-heap (evict the smallest of the top-k); k smallest/closest a max-heap (negate the key).
- **Answer**: k-th largest is `heap[0]`; top-k is the whole heap (sort if order matters).
- **Key**: distance, frequency, `(-count, word)` for tie-breaking.
- One-liner if allowed: `heapq.nlargest(k, items, key=...)`; show you know how it works.

#### Example 1.1: Kth largest element (heap vs quickselect vs sort) [Medium]

**Problem:** [Kth Largest Element in an Array (215)](https://leetcode.com/problems/kth-largest-element-in-an-array/) · **Template:** Template 1 (size-k heap)

**Problem.** Return the k-th largest value in an unsorted array (duplicates count separately).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "k-th largest", one-shot, array in memory.

**Approach:** three options; the interviewer wants the tradeoff.

| Method | Time | Space | When |
|---|---|---|---|
| Sort | O(n log n) | O(n) or O(1) | Simplest; fine for modest n |
| Size-k min-heap | O(n log k) | O(k) | Streaming, or k much smaller than n |
| Quickselect | O(n) average, O(n^2) worst | O(n) as written (O(1) in place) | One-shot, all in memory, best average |

```python
def kth_largest_heap(nums: list[int], k: int) -> int:
    heap = nums[:k]
    heapq.heapify(heap)
    for x in nums[k:]:
        if x > heap[0]:
            heapq.heapreplace(heap, x)
    return heap[0]
```

```python
def kth_largest_quickselect(nums: list[int], k: int) -> int:
    target = k  # k-th largest among the current candidate list
    arr = nums
    while True:
        pivot = random.choice(arr)
        bigger = [x for x in arr if x > pivot]
        equal_count = sum(1 for x in arr if x == pivot)
        if target <= len(bigger):
            arr = bigger
        elif target <= len(bigger) + equal_count:
            return pivot
        else:
            target -= len(bigger) + equal_count
            arr = [x for x in arr if x < pivot]
```

**Why quickselect works:** a three-way partition fixes the pivot's rank range, so recurse only into the side holding rank k: expected n + n/2 + n/4 + ... = O(n). Random pivots avoid adversarial O(n^2); three-way partitioning handles many duplicates. Mention in-place Lomuto/Hoare for O(1) space.

**Complexity:** heap O(n log k) time, O(k) space; quickselect O(n) average; sort O(n log n).

**Edge cases:** k = 1, k = n, all equal values, negatives.

Streaming variant (Kth Largest Element in a Stream): keep the size-k heap between calls; `heap[0]` answers after each add.

</details>

#### Example 1.2: Top k frequent elements (heap vs bucket sort) and k closest points [Medium]

**Problems:** [Top K Frequent Elements (347)](https://leetcode.com/problems/top-k-frequent-elements/), [K Closest Points to Origin (973)](https://leetcode.com/problems/k-closest-points-to-origin/) · **Template:** Template 1 (size-k heap)

**Problem.** Return the k values that occur most often.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** count, then top-k by count.

**Approach:** count with `Counter`, then pick the top k with a heap or buckets.

**Approach A, heap:** `nlargest(k, counts, key=counts.get)`. O(n + m log k) for m distinct values.

**Approach B, bucket sort:** frequencies are integers in `[1, n]`, so bucket values by frequency and scan high to low. O(n) time and space: the answer to "can you beat O(n log k)?"

```python
def top_k_frequent_bucket(nums: list[int], k: int) -> list[int]:
    counts = Counter(nums)
    buckets: list[list[int]] = [[] for _ in range(len(nums) + 1)]
    for v, c in counts.items():
        buckets[c].append(v)
    out: list[int] = []
    for c in range(len(buckets) - 1, 0, -1):
        for v in buckets[c]:
            out.append(v)
            if len(out) == k:
                return out
    return out
```

Tie-breaking variant (Top K Frequent Words): count descending, then alphabetical. Take the k **smallest** by `(-count, word)`.

```python
def top_k_frequent_words(words: list[str], k: int) -> list[str]:
    counts = Counter(words)
    return heapq.nsmallest(k, counts, key=lambda w: (-counts[w], w))
```

K Closest Points to Origin: k smallest by distance, so a size-k **max**-heap on negated squared distance. Skip `sqrt`: it's monotonic, slower, and adds float error.

```python
def k_closest(points: list[list[int]], k: int) -> list[list[int]]:
    heap: list[tuple[int, int]] = []
    for i, (x, y) in enumerate(points):
        d = x * x + y * y  # no sqrt: monotonic, exact
        if len(heap) < k:
            heapq.heappush(heap, (-d, i))
        elif d < -heap[0][0]:
            heapq.heapreplace(heap, (-d, i))
    return [points[i] for _, i in heap]
```

**Complexity:** O(n log k) time, O(k) space; quickselect on distance gives O(n) average.

**Edge cases:** `(-d, point)` works for a list `point` (lists compare), but unorderable objects fail on ties; use an index.

</details>

### Template 2: two heaps

> - **Example problem:** [Sliding Window Median (480)](https://leetcode.com/problems/sliding-window-median/)
> - **Why this template:** each window's median comes from the tops of a lower-half max-heap and an upper-half min-heap.
> - **Worked below:** Example 2.2 (Find Median from Data Stream (295) in Example 2.1)
> - **Also fits:** [Find Median from Data Stream (295)](https://leetcode.com/problems/find-median-from-data-stream/)

```python
class TemplateTwoHeaps:
    def __init__(self) -> None:
        self.lo: list[int] = []  # negated values
        self.hi: list[int] = []

    def add(self, x: int) -> None:
        heapq.heappush(self.lo, -x)                         # 1. into lower half
        heapq.heappush(self.hi, -heapq.heappop(self.lo))    # 2. move lower's max up (ordering)
        if len(self.hi) > len(self.lo):                     # 3. rebalance sizes
            heapq.heappush(self.lo, -heapq.heappop(self.hi))

    def median(self) -> float:
        if len(self.lo) > len(self.hi):
            return float(-self.lo[0])
        return (-self.lo[0] + self.hi[0]) / 2
```

Knobs:
- **Other percentiles**: keep `len(lo) == ceil(p * n)` instead of half.
- **Removals** (sliding window): lazy deletion plus valid-size counters (Example 2.2).
- Generally, two heaps split any set into "chosen" and "not chosen" around a boundary (IPO's sorted list + heap is the same idea).

#### Example 2.1: Running median of a stream (two heaps) [Hard]

**Problem:** [Find Median from Data Stream (295)](https://leetcode.com/problems/find-median-from-data-stream/) · **Template:** Template 2 (two heaps)

**Problem.** Design a class that receives integers one at a time (`add_num`) and reports the median of everything so far (`find_median`) at any moment. With an even count, the median is the mean of the two middle values. Example: add 1, add 2, median gives 1.5; add 3, median gives 2.0.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** repeated medians of a growing stream: two heaps. Re-sorting per query is O(n log n); an insertion-sorted list is O(n) per add (shifting).

**Approach:** lower half `lo` (max-heap, negated) and upper half `hi` (min-heap), with two invariants:
1. Ordering: every value in `lo` is `<=` every value in `hi`.
2. Sizes: `len(lo) == len(hi)` or `len(lo) == len(hi) + 1`.

The median is `-lo[0]` (odd count) or the mean of the tops (even). `add_num` pushes through: push `x` into `lo`, then move `lo`'s max into `hi`. The value crossing is the lower half's largest (including `x`), so ordering holds without comparing both tops. If `hi` is now bigger, move its min back to `lo` (it's `>=` all of `lo`, so ordering survives).

```python
class MedianFinder:
    def __init__(self) -> None:
        self.lo: list[int] = []  # max-heap of the smaller half (values negated)
        self.hi: list[int] = []  # min-heap of the larger half

    def add_num(self, x: int) -> None:
        heapq.heappush(self.lo, -x)                       # 1. tentatively into the lower half
        heapq.heappush(self.hi, -heapq.heappop(self.lo))  # 2. lower's max moves up: ordering holds
        if len(self.hi) > len(self.lo):                   # 3. sizes: lo may lead by at most one
            heapq.heappush(self.lo, -heapq.heappop(self.hi))

    def find_median(self) -> float:
        if len(self.lo) > len(self.hi):
            return float(-self.lo[0])
        return (-self.lo[0] + self.hi[0]) / 2
```

**Complexity:** `add_num` O(log n) (three heap ops), `find_median` O(1); O(n) space.

**Edge cases:** the first value (moved to `hi`, then back, so `lo` is never empty at query time); duplicates and negatives; forgetting to negate `lo[0]`; returning an int instead of a float for odd counts. Follow-ups: values in `[0, 100]` → a 101-slot count array, O(100) walk to the middle; median of the last k values → removals (Example 2.2).

</details>

#### Example 2.2: Sliding window median (two heaps + lazy deletion) [Hard]

**Problem:** [Sliding Window Median (480)](https://leetcode.com/problems/sliding-window-median/) · **Template:** Template 2 (two heaps)

**Problem.** Report the median of every size-k window.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** median = two heaps; the window removes arbitrary (non-top) elements, which `heapq` can't do directly.

**Approach:** keep the two-heap invariant over **valid** elements, deleting lazily:
- `delayed[x]` counts pending deletions of value x.
- `lo_size` / `hi_size` count valid elements only; `len(heap)` includes dead entries.
- To erase x: compare with `-lo[0]` to find its heap (every valid element `<= max(lo)` is in `lo`), decrement that size, and prune at once if x is the top.
- `prune(heap)` pops while the top is delayed; call it whenever a top may be stale (after erasing at the top, after moving between heaps).
- On reads, both tops are valid, so the median is correct.

```python
def sliding_window_median(nums: list[int], k: int) -> list[float]:
    lo: list[int] = []  # max-heap (negated)
    hi: list[int] = []  # min-heap
    delayed: dict[int, int] = defaultdict(int)
    lo_size = hi_size = 0

    def prune(heap: list[int], sign: int) -> None:
        while heap and delayed[sign * heap[0]]:
            delayed[sign * heap[0]] -= 1
            heapq.heappop(heap)

    def rebalance() -> None:
        nonlocal lo_size, hi_size
        if lo_size > hi_size + 1:
            heapq.heappush(hi, -heapq.heappop(lo))
            lo_size -= 1
            hi_size += 1
            prune(lo, -1)
        elif lo_size < hi_size:
            heapq.heappush(lo, -heapq.heappop(hi))
            hi_size -= 1
            lo_size += 1
            prune(hi, 1)

    def insert(x: int) -> None:
        nonlocal lo_size, hi_size
        if not lo or x <= -lo[0]:
            heapq.heappush(lo, -x)
            lo_size += 1
        else:
            heapq.heappush(hi, x)
            hi_size += 1
        rebalance()

    def erase(x: int) -> None:
        nonlocal lo_size, hi_size
        delayed[x] += 1
        if x <= -lo[0]:  # x logically lives in lo
            lo_size -= 1
            if x == -lo[0]:
                prune(lo, -1)
        else:
            hi_size -= 1
            if hi and x == hi[0]:
                prune(hi, 1)
        rebalance()

    def median() -> float:
        if k % 2:
            return float(-lo[0])
        return (-lo[0] + hi[0]) / 2

    for x in nums[:k]:
        insert(x)
    out = [median()]
    for i in range(k, len(nums)):
        insert(nums[i])
        erase(nums[i - k])
        out.append(median())
    return out
```

**Complexity:** O(n log n) time (bounded pushes/pops per element; dead entries can reach O(n)), O(n) space worst case.

**Edge cases:** k = 1; duplicates (pruning by value is safe, since equal values are interchangeable); even k averaging; `(a + b) / 2` overflow in other languages.

**Simpler alternative:** a sorted window with `bisect.insort` and `bisect` + `pop`. O(nk) worst case from shifting, but often fast for moderate k and far less error-prone. `sortedcontainers.SortedList` gives O(log k) if allowed.

</details>

### Template 3: k-way merge

> - **Example problem:** [Smallest Range Covering Elements from K Lists (632)](https://leetcode.com/problems/smallest-range-covering-elements-from-k-lists/)
> - **Why this template:** k sorted lists and one number from each: a heap frontier with one element per list.
> - **Worked below:** Example 3.2 (Merge k Sorted Lists (23) in Example 3.1)
> - **Also fits:** [Merge k Sorted Lists (23)](https://leetcode.com/problems/merge-k-sorted-lists/)

```python
def template_k_way_merge(sources: list[list[int]]) -> list[int]:
    heap = [(src[0], i, 0) for i, src in enumerate(sources) if src]
    heapq.heapify(heap)
    out = []
    while heap:
        val, i, j = heapq.heappop(heap)
        out.append(val)  # TODO: per-pop work (count to k, track range, ...)
        if j + 1 < len(sources[i]):  # advance the source we just consumed from
            heapq.heappush(heap, (sources[i][j + 1], i, j + 1))
    return out
```

Knobs:
- Entry `(value, source_id, position)`: `source_id` doubles as tie-breaker, so payloads (`ListNode`, file handles) are never compared.
- **Implicit sources**: sorted matrix rows, or `a[i] + b[j]` for fixed `i` (k smallest pairs). Seed only `min(k, rows)` sources.
- Stop after k pops for k-th smallest; track the frontier's running max for ranges.
- Stdlib: `heapq.merge(*iterables, key=...)` is a lazy k-way merge.

#### Example 3.1: Merge k sorted linked lists (heap frontier) [Hard]

**Problem:** [Merge k Sorted Lists (23)](https://leetcode.com/problems/merge-k-sorted-lists/) · **Template:** Template 3 (k-way merge)

**Problem.** Given an array of k ascending linked lists (some may be empty), combine them into one sorted list and return its head. Example: `[1 -> 4 -> 5, 1 -> 3 -> 4, 2 -> 6]` gives `1 -> 1 -> 2 -> 3 -> 4 -> 4 -> 5 -> 6`; `[]` or `[None]` gives an empty list.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** k-way merge frontier with `ListNode` payloads, so entries need a tie-breaker besides the node.

**Approach:** the heap holds each list's smallest unmerged node, so its top is the smallest remaining overall. Pop it, append via module 09's dummy + tail builder, and push its successor (the only new candidate). Entries are `(node.val, i, node)`; without `i`, equal values make Python compare `ListNode`s and raise `TypeError`. The last popped node ends an input list, so the output is terminated.

```python
def merge_k_lists(lists: list[Optional[ListNode]]) -> Optional[ListNode]:
    heap = [(node.val, i, node) for i, node in enumerate(lists) if node]
    heapq.heapify(heap)  # i breaks ties, so two ListNodes are never compared
    dummy = tail = ListNode()
    while heap:
        _, i, node = heapq.heappop(heap)
        tail.next = node
        tail = node
        if node.next:  # advance the list we just took from
            heapq.heappush(heap, (node.next.val, i, node.next))
    return dummy.next
```

**Complexity:** O(N log k) for N total nodes (each pushed and popped once; heap ≤ k). O(k) extra space; nodes are relinked.

**Edge cases:** `lists` empty or all `None` (filter when seeding, return `None`); a single list; many equal values (the index tie-breaker). Alternative: divide and conquer, merging pairs with `merge_two` (module 09) over log k rounds: O(N log k) time, O(1) extra space iteratively. Merging one at a time into an accumulator is O(N k); say why you avoid it.

</details>

#### Example 3.2: Smallest range covering k sorted lists (k-way merge) [Hard]

**Problem:** [Smallest Range Covering Elements from K Lists (632)](https://leetcode.com/problems/smallest-range-covering-elements-from-k-lists/) · **Template:** Template 3 (k-way merge)

**Problem.** Given k sorted lists, find the smallest range `[a, b]` including at least one number from each (ties: smaller `a`).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a k-way merge frontier always holds one element per list, so `[min(frontier), max(frontier)]` is always valid.

**Approach:** seed the heap with each list's head; track the frontier's `cur_max`. Each step, `[heap_min, cur_max]` is a candidate; then advance the list owning the min (the only move that can shrink the range, since the max only grows). Stop when that list is exhausted: any later range would skip it.

```python
def smallest_range(nums: list[list[int]]) -> list[int]:
    heap = [(row[0], i, 0) for i, row in enumerate(nums)]
    heapq.heapify(heap)
    cur_max = max(row[0] for row in nums)
    best = [heap[0][0], cur_max]
    while True:
        lo, i, j = heapq.heappop(heap)
        if cur_max - lo < best[1] - best[0]:
            best = [lo, cur_max]
        if j + 1 == len(nums[i]):
            return best  # this list cannot advance: no later window covers it
        nxt = nums[i][j + 1]
        cur_max = max(cur_max, nxt)
        heapq.heappush(heap, (nxt, i, j + 1))
```

**Complexity:** O(N log k) for N total elements, O(k) space.

**Edge cases:** a single list (range `[x, x]` for its first element); identical lists.

Same frontier, different stop (Kth Smallest Element in a Sorted Matrix): rows are sources; after k-1 pops, the top is the answer. Seed only `min(n, k)` rows.

```python
def kth_smallest_matrix(matrix: list[list[int]], k: int) -> int:
    n = len(matrix)
    heap = [(matrix[r][0], r, 0) for r in range(min(n, k))]  # rows beyond k cannot matter
    heapq.heapify(heap)
    for _ in range(k - 1):
        _, r, c = heapq.heappop(heap)
        if c + 1 < len(matrix[r]):
            heapq.heappush(heap, (matrix[r][c + 1], r, c + 1))
    return heap[0][0]
```

O(k log min(n, k)). When k nears n^2, binary search on value with a staircase count is O(n log(max - min)) (`07-binary-search.md`).

</details>

### Template 4: event simulation with a heap

> - **Example problem:** [Single-Threaded CPU (1834)](https://leetcode.com/problems/single-threaded-cpu/)
> - **Why this template:** tasks arrive over time and the idle CPU picks the shortest available: arrivals by time plus a priority heap.
> - **Worked below:** Example 4.2

```python
def template_event_simulation(tasks: list[tuple[int, int]]) -> list[int]:
    order = sorted(range(len(tasks)), key=lambda i: tasks[i][0])
    ready: list[tuple[int, int]] = []
    time, p, out = 0, 0, []
    while p < len(order) or ready:
        if not ready and time < tasks[order[p]][0]:
            time = tasks[order[p]][0]  # idle: jump the clock to the next arrival
        while p < len(order) and tasks[order[p]][0] <= time:  # admit everything that has arrived
            i = order[p]
            heapq.heappush(ready, (tasks[i][1], i))  # TODO: priority key
            p += 1
        dur, i = heapq.heappop(ready)  # TODO: pick and execute
        time += dur
        out.append(i)
    return out
```

Two structures: a **sorted list + pointer** for "not yet arrived" (by time) and a **heap** for "arrived, waiting" (by priority). Knobs: the priority key; what "execute" does (advance time, add profit, free a resource); multiple servers (a second heap of busy servers by free time). Jump the clock when idle; ticking is O(max_time).

#### Example 4.1: Meeting rooms and task scheduling (greedy with heaps) [Medium]

**Problems:** [Meeting Rooms II (253)](https://leetcode.com/problems/meeting-rooms-ii/) (Premium), [Task Scheduler (621)](https://leetcode.com/problems/task-scheduler/) · **Template:** Template 4 (event simulation with a heap)

**Meeting Rooms II** (LeetCode: Meeting Rooms II). **Problem.** Minimum rooms so no two overlapping meetings share one.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** minimum rooms = peak concurrent meetings; reusing a room needs the earliest end time among rooms in use: a min-heap.

**Approach:** process meetings by start; the heap holds in-use rooms' end times. If the earliest-ending room is free (`ends[0] <= start`), reuse it (`heapreplace`); otherwise open one. The heap never shrinks, so its final size is the peak.

```python
def min_meeting_rooms(intervals: list[list[int]]) -> int:
    ends: list[int] = []
    for start, end in sorted(intervals):
        if ends and ends[0] <= start:  # earliest-ending room is free: reuse it
            heapq.heapreplace(ends, end)
        else:
            heapq.heappush(ends, end)
    return len(ends)
```

**Complexity:** O(n log n).

**Edge cases:** `<=` because a meeting ending at 10 frees the room for one starting at 10 (confirm). Sweep-line alternative: `12-intervals-and-sweep-line.md`; the heap version also tells you **which** room (keep `(end, room_id)`).

</details>

**Task Scheduler** (LeetCode: Task Scheduler). **Problem.** Letter-labeled tasks take one unit each; runs of the same letter need at least n units between them; idling is allowed. Minimum total time.

<details>
<summary>Decode, approach, and solution</summary>

**Approach A, simulation:** each tick, run the task with the most remaining (max-heap of counts) and park it in a FIFO cooldown queue with its ready time. Jump the clock when only cooldowns remain.

```python
def least_interval_sim(tasks: list[str], n: int) -> int:
    heap = [-c for c in Counter(tasks).values()]
    heapq.heapify(heap)
    cooldown: deque[tuple[int, int]] = deque()  # (ready_time, negated remaining count)
    time = 0
    while heap or cooldown:
        time += 1
        if heap:
            c = heapq.heappop(heap) + 1  # run one unit of the most frequent task
            if c:
                cooldown.append((time + n, c))
        else:
            time = cooldown[0][0]  # idle: jump to the next ready time
        if cooldown and cooldown[0][0] == time:
            heapq.heappush(heap, cooldown.popleft()[1])
    return time
```

**Approach B, counting:** the most frequent task (count `mx`) forces `mx - 1` frames of length `n + 1`, plus a final frame of every task tied at `mx`. If other tasks fill all gaps, no idling: the answer is `len(tasks)`.

```python
def least_interval_formula(tasks: list[str], n: int) -> int:
    counts = Counter(tasks).values()
    mx = max(counts)
    num_max = sum(1 for c in counts if c == mx)
    return max(len(tasks), (mx - 1) * (n + 1) + num_max)
```

Lead with the simulation if you can't derive B on the spot; it generalizes (output the schedule, Reorganize String, Rearrange String k Distance Apart). The code file cross-checks both on random inputs.

</details>

#### Example 4.2: Single-threaded CPU and IPO (event simulation) [Hard]

**Problems:** [Single-Threaded CPU (1834)](https://leetcode.com/problems/single-threaded-cpu/), [IPO (502)](https://leetcode.com/problems/ipo/) · **Template:** Template 4 (event simulation with a heap)

**Single-Threaded CPU** (LeetCode: Single-Threaded CPU). **Problem.** `tasks[i] = [enqueue_time, processing_time]`. When idle, the CPU runs the available task with the shortest processing time (ties: smaller index) to completion, or idles if none is available. Return the order.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** arrivals + priority among the available + a clock: Template 4 verbatim.

```python
def get_order(tasks: list[list[int]]) -> list[int]:
    order = sorted(range(len(tasks)), key=lambda i: tasks[i][0])
    ready: list[tuple[int, int]] = []
    time, p, out = 0, 0, []
    while len(out) < len(tasks):
        if not ready and time < tasks[order[p]][0]:
            time = tasks[order[p]][0]
        while p < len(order) and tasks[order[p]][0] <= time:
            i = order[p]
            heapq.heappush(ready, (tasks[i][1], i))
            p += 1
        proc, i = heapq.heappop(ready)
        time += proc
        out.append(i)
    return out
```

**Complexity:** O(n log n).

**Edge cases:** sort indices, not tasks (original indices are needed); jump the clock across idle gaps; `(proc, i)` encodes the tie-break.

</details>

**IPO** (LeetCode: IPO). **Problem.** At most k projects; project i needs `capital[i]` to start and adds `profits[i]`. Start with w. Maximize final capital.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** capital is the "clock". Sort projects by required capital; each round, admit every affordable project into a profit max-heap and take the best. Greedy is optimal: profit only grows capital, so the best affordable pick never closes options (exchange argument, `18-greedy.md`).

```python
def find_maximized_capital(k: int, w: int, profits: list[int], capital: list[int]) -> int:
    projects = sorted(zip(capital, profits))
    avail: list[int] = []
    p = 0
    for _ in range(k):
        while p < len(projects) and projects[p][0] <= w:
            heapq.heappush(avail, -projects[p][1])
            p += 1
        if not avail:
            break  # nothing affordable, and capital cannot grow
        w -= heapq.heappop(avail)
    return w
```

**Complexity:** O(n log n + k log n).

</details>

## Template cheat sheet

All templates, for review.

### Template 1: size-k heap (top-k)
```python
def template_size_k_heap(items: Iterable[Any], k: int, key=lambda x: x) -> list[Any]:
    heap: list[tuple] = []
    for i, x in enumerate(items):
        entry = (key(x), i, x)  # i breaks ties so x itself is never compared
        if len(heap) < k:
            heapq.heappush(heap, entry)
        elif entry[0] > heap[0][0]:  # TODO: strictly better than the gatekeeper
            heapq.heapreplace(heap, entry)  # pop + push in one O(log k) step
    return [x for _, _, x in heap]  # unordered; sort if the caller needs order
```
### Template 2: two heaps
```python
class TemplateTwoHeaps:
    def __init__(self) -> None:
        self.lo: list[int] = []  # negated values
        self.hi: list[int] = []

    def add(self, x: int) -> None:
        heapq.heappush(self.lo, -x)                         # 1. into lower half
        heapq.heappush(self.hi, -heapq.heappop(self.lo))    # 2. move lower's max up (ordering)
        if len(self.hi) > len(self.lo):                     # 3. rebalance sizes
            heapq.heappush(self.lo, -heapq.heappop(self.hi))

    def median(self) -> float:
        if len(self.lo) > len(self.hi):
            return float(-self.lo[0])
        return (-self.lo[0] + self.hi[0]) / 2
```
### Template 3: k-way merge
```python
def template_k_way_merge(sources: list[list[int]]) -> list[int]:
    heap = [(src[0], i, 0) for i, src in enumerate(sources) if src]
    heapq.heapify(heap)
    out = []
    while heap:
        val, i, j = heapq.heappop(heap)
        out.append(val)  # TODO: per-pop work (count to k, track range, ...)
        if j + 1 < len(sources[i]):  # advance the source we just consumed from
            heapq.heappush(heap, (sources[i][j + 1], i, j + 1))
    return out
```
### Template 4: event simulation with a heap
```python
def template_event_simulation(tasks: list[tuple[int, int]]) -> list[int]:
    order = sorted(range(len(tasks)), key=lambda i: tasks[i][0])
    ready: list[tuple[int, int]] = []
    time, p, out = 0, 0, []
    while p < len(order) or ready:
        if not ready and time < tasks[order[p]][0]:
            time = tasks[order[p]][0]  # idle: jump the clock to the next arrival
        while p < len(order) and tasks[order[p]][0] <= time:  # admit everything that has arrived
            i = order[p]
            heapq.heappush(ready, (tasks[i][1], i))  # TODO: priority key
            p += 1
        dur, i = heapq.heappop(ready)  # TODO: pick and execute
        time += dur
        out.append(i)
    return out
```

## Lazy deletion and tie-breakers

**Lazy deletion.** `heapq` has no `remove(x)` (finding x is O(n)). Mark x dead; discard it when it reaches the top. Rules: prune before every top read; track size separately from `len(heap)`. Used in sliding window median, skyline (buildings ending), and stock-price-fluctuation problems.

```python
class LazyHeap:
    def __init__(self) -> None:
        self.heap: list[Any] = []
        self.dead: Counter = Counter()
        self.size = 0  # number of live items

    def push(self, x: Any) -> None:
        heapq.heappush(self.heap, x)
        self.size += 1

    def remove(self, x: Any) -> None:
        """Caller guarantees x is currently live."""
        self.dead[x] += 1
        self.size -= 1

    def _prune(self) -> None:
        while self.heap and self.dead[self.heap[0]]:
            self.dead[self.heap[0]] -= 1
            heapq.heappop(self.heap)

    def top(self) -> Any:
        self._prune()
        return self.heap[0]

    def pop(self) -> Any:
        self._prune()
        self.size -= 1
        return heapq.heappop(self.heap)

    def __len__(self) -> int:
        return self.size
```

**Tie-breakers.** Tuples compare left to right, so equal priorities compare the next field: `TypeError` for objects like `ListNode` or dicts, arbitrary (non-FIFO) order otherwise. An increasing counter, `(priority, seq, item)`, gives a total order and FIFO among equals (fair schedulers). Entry-marking also supports priority **updates** (mark the old entry dead, push a new one): decrease-key for Dijkstra with `heapq`.

```python
class PriorityQueue:
    def __init__(self) -> None:
        self.heap: list[list] = []
        self.entries: dict[Hashable, list] = {}
        self.counter = itertools.count()

    def push(self, item: Hashable, priority: int) -> None:
        if item in self.entries:
            self.cancel(item)
        entry = [priority, next(self.counter), item, True]
        self.entries[item] = entry
        heapq.heappush(self.heap, entry)

    def cancel(self, item: Hashable) -> None:
        self.entries.pop(item)[3] = False

    def pop(self) -> Hashable:
        while self.heap:
            _, _, item, alive = heapq.heappop(self.heap)
            if alive:
                del self.entries[item]
                return item
        raise KeyError("pop from empty priority queue")
```

## Variations and extensions

- **Find K Pairs with Smallest Sums**: implicit matrix `a[i] + b[j]`; seed `(a[i] + b[0])` for the first k rows, advance `j` (`k_smallest_pairs`).
- **Ugly numbers / super ugly numbers**: k-way merge of multiplied sequences, or a min-heap with a seen set.
- **Reorganize String**: max-heap of counts; alternate the top two, or hold back the last letter for one round.
- **Multi-server simulation** (Process Tasks Using Servers, Meeting Rooms III): `free` heap by `(weight, id)`, `busy` by `(free_time, id)`. Before each assignment, return servers with `free_time <= now` to `free`; if none are free, jump the clock to the earliest busy one.
- **Minimum cost to hire k workers**: sort by ratio, size-k max-heap of quality (sorted pointer + heap, like IPO).
- **Furthest building you can reach**: size-k min-heap of climbs given ladders; overflow uses bricks.
- **Streaming top-k with bounded memory**: Count-Min Sketch or Space-Saving + heap (approximate).
- **External sort**: sort memory-sized chunks, then k-way merge the runs from disk.
- **Heap vs balanced BST**: for both min and max, arbitrary deletion, or rank queries, `SortedList` may beat heaps with lazy deletion.

Follow-ups:
- Doesn't fit in memory → size-k heap streams in O(k); k-way merge of sorted runs.
- k close to n → size n - k heap for the opposite end, or quickselect.
- Distributed top-k → local top-k per shard, then merge. Exact only if each key is fully aggregated on one shard; otherwise aggregate by key first (or approximate).
- Thread-safe → `queue.PriorityQueue` wraps `heapq` with a lock (`21-practical-coding-for-infra-roles.md`).

## Common pitfalls

- Wrong heap direction: k largest needs a **min**-heap of size k, not a max-heap of everything.
- Forgetting to negate on both push and pop for max-heaps; wrap in helpers if needed.
- Non-comparable payloads: `(priority, node)` crashes on ties. Add an index or counter.
- Using `len(heap)` as the logical size under lazy deletion.
- Peeking `heap[0]` without pruning dead entries first.
- Ticking the clock one unit at a time; jump to the next event.
- Sorting tasks and losing their original indices.
- Heapifying all n stream items when only k are needed: O(n) memory.
- `heappop` on an empty heap raises `IndexError`; guard with `if heap`.
- Assuming the heap is sorted. Only `heap[0]` is guaranteed; use `sorted(heap)` for order.

## When NOT to use / look-alikes

| Looks like | Actually | How to tell |
|---|---|---|
| "Maximum in each sliding window" | Monotonic deque (module 08), O(n) | Only max/min, window slides by one. Heap + lazy deletion works but is O(n log n) |
| "k-th smallest in a sorted matrix" with large k | Binary search on value (module 07) | k near n^2; counting is O(n) per probe |
| "Top k" with k = 1 | Single pass with a variable | No heap needed |
| Meeting rooms count | Sweep line (module 12) | Only the count matters; sweep is simpler, equally fast |
| "Merge two sorted arrays" | Two pointers (module 05) | k = 2; heap is overkill |
| Frequency top-k with small value range | Bucket / counting sort | Frequencies ≤ n give O(n) |
| Repeated "next best" where choices depend on future ones | DP (module 16), not greedy heap | Exchange argument fails; a locally best pick can block a better total |
| Shortest path with unit weights | Plain BFS (module 13) | Heap needed only for varying non-negative weights (Dijkstra) |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Kth Largest Element in an Array (215)](https://leetcode.com/problems/kth-largest-element-in-an-array/) · Medium · Size-k heap / quickselect: heap vs quickselect tradeoff.
2. [K Closest Points to Origin (973)](https://leetcode.com/problems/k-closest-points-to-origin/) · Medium · Top-k with a key: max-heap of size k; no sqrt.
3. [Task Scheduler (621)](https://leetcode.com/problems/task-scheduler/) · Medium · Heap simulation / counting: simulation vs counting formula.
4. [Merge k Sorted Lists (23)](https://leetcode.com/problems/merge-k-sorted-lists/) · Hard · k-way merge: tie-breaker index.
5. [Find Median from Data Stream (295)](https://leetcode.com/problems/find-median-from-data-stream/) · Hard · Two heaps: push-through rebalancing.

### Optional

**Size-k / top-k**
- [Kth Largest Element in a Stream (703)](https://leetcode.com/problems/kth-largest-element-in-a-stream/) · Easy · Persistent size-k min-heap.
- [Last Stone Weight (1046)](https://leetcode.com/problems/last-stone-weight/) · Easy · Max-heap as a greedy engine.
- [Top K Frequent Elements (347)](https://leetcode.com/problems/top-k-frequent-elements/) · Medium · Heap vs bucket sort.
- [Top K Frequent Words (692)](https://leetcode.com/problems/top-k-frequent-words/) · Medium · Tuple key for tie-break.

**Two heaps**
- [Sliding Window Median (480)](https://leetcode.com/problems/sliding-window-median/) · Hard · Lazy deletion with valid-size counters.
- [IPO (502)](https://leetcode.com/problems/ipo/) · Hard · Sorted unlock pointer + max-heap.

**K-way merge**
- [Kth Smallest Element in a Sorted Matrix (378)](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/) · Medium · Rows as sources; or binary search.
- [Find K Pairs with Smallest Sums (373)](https://leetcode.com/problems/find-k-pairs-with-smallest-sums/) · Medium · Implicit matrix, seed k rows.
- [Smallest Range Covering Elements from K Lists (632)](https://leetcode.com/problems/smallest-range-covering-elements-from-k-lists/) · Hard · Frontier min and running max.

**Scheduling / simulation**
- [Meeting Rooms II (253)](https://leetcode.com/problems/meeting-rooms-ii/) · Medium · Premium · Heap of end times.
- [Reorganize String (767)](https://leetcode.com/problems/reorganize-string/) · Medium · Hold back the last letter.
- [Single-Threaded CPU (1834)](https://leetcode.com/problems/single-threaded-cpu/) · Medium · Sorted arrivals + available heap + clock jumps.
- [Process Tasks Using Servers (1882)](https://leetcode.com/problems/process-tasks-using-servers/) · Medium · Free/busy heaps.
- [Meeting Rooms III (2402)](https://leetcode.com/problems/meeting-rooms-iii/) · Hard · Free rooms by id, busy rooms by end time, delays.
- [Minimum Cost to Hire K Workers (857)](https://leetcode.com/problems/minimum-cost-to-hire-k-workers/) · Hard · Sort by ratio, size-k max-heap.

## Self-check

1. Why does "k largest" use a min-heap? What does its top represent at every step?
2. State the running median's two invariants. Why does the push-through step in `add` keep ordering without comparing both tops?
3. In sliding window median, why can't you rebalance with `len(lo)`, and when must you prune?
4. In smallest range, why is advancing the minimum's list the only useful move, and why stop once one list is exhausted?
5. Give the time and space for top-k frequent via heap, bucket sort, and quickselect. When would you pick each?
6. Why add a counter to heap entries? Name two problems it solves.
7. In an event simulation, what lives in the sorted list and what in the heap? Why jump the clock instead of ticking?
