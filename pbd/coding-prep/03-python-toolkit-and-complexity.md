# 03. Python Toolkit and Complexity

You'll lose more interview time to "what does `heapq` do with ties?" and "is this O(n) or O(n^2)?" than to a missing algorithm. This module gives the cost model and stdlib vocabulary to derive a target complexity from constraints, map it to a pattern, and code it without fighting the language.

Companion code: `code/03_python_toolkit.py` (asserts every idiom and pitfall below).

## 1. Big-O refresher for interviews

What interviewers expect:

- State **time and space** in the input's own names (`n`, `m`, `k`, `V`, `E`), not just "n".
- Count **auxiliary** space (beyond the input) and say so. Recursion depth counts.
- Give the **worst case** unless you say "amortized" or "expected" (hash maps are expected O(1)).
- Drop constants but keep independent variables: O(n + m) is not O(n); O(n log k) is not O(n log n) when k is small.

Common shapes and their sources:

| Complexity | Typical source |
|---|---|
| O(1) | Hash lookup, arithmetic, index access |
| O(log n) | Binary search, heap push/pop, halving loop, balanced tree op |
| O(n) | Single pass, two pointers, sliding window, prefix sums, tree BFS/DFS |
| O(n log n) | Sorting, n heap ops, n binary searches |
| O(n log k) | Top-k with a size-k heap, merging k sorted lists |
| O(V + E) | Graph BFS/DFS/topological sort |
| O(E log V) | Dijkstra with a binary heap |
| O(n^2) | All pairs, 2D DP over two strings/arrays, "for each i, expand" |
| O(n^3) | Interval DP, Floyd-Warshall, 3 nested loops |
| O(2^n) | All subsets, bitmask DP (often O(2^n * n)) |
| O(n!) | All permutations |

### Amortized analysis

Amortized cost = total cost of an operation sequence ÷ its length; rare expensive operations keep the average small.

**Dynamic array (`list.append`)**: a full backing array is replaced by a bigger copy (CPython grows by ~1.125x plus a constant; textbooks use 2x). With 2x growth, copies over n appends total 1 + 2 + 4 + ... + n < 2n, so append is **amortized O(1)**. Any constant factor > 1 works; additive growth (+10 each time) gives O(n^2) total.

Other amortized arguments:

- **Monotonic stack/deque**: each element is pushed once and popped at most once, so the inner `while` totals O(n) over the run, though one iteration may pop many.
- **Sliding window**: `left` only moves forward, so the shrink loop totals O(n).
- **Union-find** with path compression + union by rank: nearly O(1) per op (inverse Ackermann).

Interview phrasing: "The inner while loop looks quadratic, but each index is added and removed at most once, so it's O(n) total."

### Recursion cost

Time = **calls** × **work per call**; stack space = **max depth** × **frame size**.

| Recursion | Calls | Time | Stack |
|---|---|---|---|
| `f(n) = f(n-1) + O(1)` | n | O(n) | O(n) |
| `f(n) = f(n/2) + O(1)` (binary search) | log n | O(log n) | O(log n) |
| `f(n) = 2 f(n/2) + O(n)` (merge sort) | ~2n | O(n log n) | O(log n) + O(n) buffer |
| `f(n) = f(n-1) + f(n-2)` naive | ~1.6^n | O(1.6^n) | O(n) |
| Same, memoized | n distinct states | O(n) | O(n) |
| Branching factor b, depth d | b^d | O(b^d) | O(d) |
| Tree DFS | n nodes | O(n) | O(h), h = n if skewed |

For memoized recursion, **time = distinct states × work per state (excluding recursive calls)**, the most useful formula for DP analysis.

Hidden costs to call out: slicing inside recursion (see Pitfalls) and building strings with `+` in each frame.

## 2. Constraints → target complexity

Python does about **10^7 simple operations per second** and judges usually allow 1-2 seconds, so aim for 10^7 to 10^8 operations. Read the constraints first and work backwards:

| Max input size | Target complexity | Usual patterns |
|---|---|---|
| n ≤ 10 | O(n!), O(n^n) | Permutations, brute-force backtracking |
| n ≤ 20 | O(2^n), O(2^n · n) | Subsets, bitmask DP, meet-in-the-middle (n ≤ 40) |
| n ≤ 100 | O(n^4) maybe, O(n^3) safe | Small multi-dimensional DP |
| n ≤ 500 | O(n^3) | Interval DP, Floyd-Warshall, triple loop |
| n ≤ 5,000 | O(n^2) | 2D DP (LCS, edit distance), all pairs with pruning |
| n ≤ 10^5 to 10^6 | O(n log n) or O(n) | Sorting, heap, binary search, two pointers, sliding window, prefix sums + hash, monotonic stack, BFS/DFS, union-find |
| n ≤ 10^9 or more | O(log n), O(sqrt n), O(1) | Binary search on the answer, math formula, fast exponentiation, digit DP |

To guess the pattern **before** solving:

1. Find the largest variable. At `n ≤ 10^5`, O(n^2) = 10^10 is dead; say "sort, heap, binary search, or a linear scan with a hash map / window / stack."
2. Add the question shape. "Contiguous subarray" + O(n) → sliding window or prefix sums. "k-th / top k" + O(n log k) → heap. "Minimize the maximum" + values up to 10^9 → binary search on the answer.
3. Small n is a **hint**, not just permission: `n ≤ 15` almost always means bitmask/backtracking; `n ≤ 1000` with two strings often means O(n·m) DP.
4. Large values, smaller n (`n ≤ 10^5`, values ≤ 10^9): you can't index an array by value; use hashing, sorting, or coordinate compression.
5. Small values (`0 ≤ v ≤ 100`, or `1..n`): counting arrays, bucket sort, or index-as-hash work.

With no constraints given, **ask**: "What's the rough size of n?" is a senior question, not a weak one.

## 3. Python data structures and operation costs

### list (dynamic array)

| Operation | Cost |
|---|---|
| `a[i]`, `a[i] = x`, `len(a)` | O(1) |
| `a.append(x)`, `a.pop()` | O(1) amortized |
| `a.pop(0)`, `a.insert(0, x)`, `a.insert(i, x)` | O(n) |
| `x in a`, `a.index(x)`, `a.remove(x)`, `a.count(x)` | O(n) |
| `a[i:j]` | O(j - i), copies |
| `a.sort()`, `sorted(a)` | O(n log n), stable (Timsort) |
| `a.reverse()`, `a[::-1]` | O(n) |
| `a + b` | O(len(a) + len(b)), new list |
| `min(a)`, `max(a)`, `sum(a)` | O(n) |

A list is your stack: `append` / `pop`, and `a[-1]` to peek.

### dict and set (hash tables)

| Operation | Cost |
|---|---|
| `d[k]`, `d[k] = v`, `k in d`, `del d[k]`, `d.get(k, default)` | O(1) expected |
| `s.add(x)`, `x in s`, `s.discard(x)` | O(1) expected |
| Iteration | O(n) |
| `s & t` | O(min(len(s), len(t))) |
| `s | t` | O(len(s) + len(t)) |

- Dicts preserve **insertion order** (3.7+); sets promise no order.
- Keys must be hashable: int, str, tuple of hashables, frozenset. Not list, dict, set.
- `s.remove(x)` raises if missing; `s.discard(x)` doesn't.
- Adversarial collisions make ops O(n) worst case; say "expected O(1)" if pressed.

### deque

```python
from collections import deque
dq = deque([2, 3])
dq.appendleft(1); dq.append(4)   # O(1) both ends
dq.popleft(); dq.pop()           # O(1) both ends
dq[0], dq[-1]                    # O(1) peek; dq[i] in the middle is O(n)
deque(maxlen=k)                  # auto-evicts from the opposite end
```

Use it for BFS queues, sliding window max, and anywhere you'd write `list.pop(0)`.

### heapq (binary min-heap on a plain list)

| Operation | Cost |
|---|---|
| `heapq.heappush(h, x)`, `heapq.heappop(h)` | O(log n) |
| `h[0]` (peek min) | O(1) |
| `heapq.heapify(h)` | O(n), in place |
| `heapq.heappushpop(h, x)` (push then pop) | O(log n), faster than two calls |
| `heapq.heapreplace(h, x)` (pop then push) | O(log n) |
| `heapq.nlargest(k, it)`, `nsmallest` | O(n log k) |

Key techniques:

```python
# Min-heap of size k keeps the k LARGEST; root is the k-th largest.
for x in nums:
    heapq.heappush(heap, x)
    if len(heap) > k:
        heapq.heappop(heap)

# Max-heap: negate on the way in and out.
heapq.heappush(h, -x); largest = -heapq.heappop(h)

# Priority with tie-break: tuples compare element by element.
heapq.heappush(h, (priority, name))

# Non-comparable payloads: insert a counter so ties never reach the payload.
counter = itertools.count()
heapq.heappush(h, (priority, next(counter), job))
```

There's no decrease-key: push a new entry and skip stale ones on pop (lazy deletion), as Python Dijkstra does.

### bisect (binary search on a sorted list)

```python
import bisect
bisect.bisect_left(a, x)    # first index i with a[i] >= x
bisect.bisect_right(a, x)   # first index i with a[i] > x
bisect.insort(a, x)         # O(log n) search + O(n) insert
# count of values in [lo, hi]:
bisect.bisect_right(a, hi) - bisect.bisect_left(a, lo)
# key= is supported (3.10+): bisect.bisect_left(intervals, t, key=lambda iv: iv[0])
```

`insort` is O(n) because of the shift: fine up to ~10^4 inserts; beyond that, use a heap or another approach (`sortedcontainers` isn't stdlib and usually isn't allowed).

### sorting with key and cmp_to_key

```python
sorted(nums, key=lambda x: (-freq[x], x))    # freq desc, then value asc
intervals.sort(key=lambda iv: iv[0])
words.sort(key=len)                          # stable: equal lengths keep order
```

Use `functools.cmp_to_key` only when order is defined pairwise with no per-item key, e.g., "arrange numbers to form the largest concatenation": `a` goes before `b` when `a + b > b + a`.

Stability trick: sort by the secondary key, then the primary. A tuple key is usually clearer.

### Counter, defaultdict, OrderedDict

```python
from collections import Counter, defaultdict, OrderedDict
c = Counter("banana")         # {'a': 3, 'n': 2, 'b': 1}, O(n)
c["z"]                        # 0, no KeyError
c.most_common(k)              # O(n log k)
Counter(a) == Counter(b)      # anagram check
c1 - c2                       # drops non-positive counts

g = defaultdict(list); g[key].append(v)     # adjacency lists, grouping
cnt = defaultdict(int); cnt[x] += 1

od = OrderedDict()
od.move_to_end(k)             # O(1): mark as most recent
od.popitem(last=False)        # O(1): pop oldest  -> LRU cache in ~10 lines
```

### Strings

- Immutable: `s[i] = c` is a TypeError. Convert to `list(s)`, edit, then `"".join(...)`.
- `s += t` in a loop is O(n^2) worst case (CPython sometimes optimizes it in place; don't rely on it). Collect pieces in a list and `"".join(parts)` once: O(total length).
- `s[i:j]`, `s[::-1]`, `s.split()`, `s.replace(...)` are O(n); `t in s` is O(n·m) worst case.
- Handy: `ord(c) - ord('a')` for 26-slot arrays, `s.isalnum()`, `s.lower()`, `str.maketrans`.

### tuple and frozenset as keys

```python
seen = {(r, c)}                         # grid coordinates
memo[(i, j, k)]                         # multi-dimensional memo
groups[tuple(sorted(word))]             # canonical key
groups[tuple(counts_26)]                # count-signature key
edges = {frozenset((u, v))}             # unordered pair: {u,v} == {v,u}
```

Hashing a length-k tuple costs O(k): a 26-count key is O(26) per word, not free.

## 4. Idioms

```python
for i, x in enumerate(nums): ...            # index + value
for a, b in zip(xs, ys): ...                # stops at the shortest
for a, b in itertools.pairwise(nums): ...   # adjacent pairs (3.10+), no copy
first, *rest = nums                         # unpacking
a, b = b, a                                 # swap
r, c = divmod(idx, cols)                    # flat index -> 2D
squares = [x * x for x in nums if x > 0]    # comprehension
pos = {v: i for i, v in enumerate(nums)}    # value -> index
grid = [[0] * cols for _ in range(rows)]    # 2D grid, correct way
any(x < 0 for x in nums); all(...)          # short-circuit
max(words, key=len)                         # argmax by key
```

### Memoization

```python
from functools import cache, lru_cache

@cache                       # same as lru_cache(maxsize=None)
def dp(i: int, j: int) -> int: ...

def solve(nums):
    @cache                   # define inside so the cache does not leak across calls
    def go(i): ...
    return go(0)
```

Arguments must be hashable (convert lists to tuples). Deep recursion still hits the recursion limit (see Pitfalls); above a few thousand levels, prefer bottom-up DP. `go.cache_clear()` resets.

### itertools

| Function | Use |
|---|---|
| `accumulate(nums, initial=0)` | Length-n+1 prefix sums |
| `accumulate(nums, max)` | Running max |
| `combinations(items, k)` | C(n, k) k-subsets |
| `permutations(items)` | n! orderings |
| `product(range(3), repeat=2)` | Cartesian product; replaces nested loops |
| `groupby(sorted_items, key)` | Runs of **consecutive** equal keys |
| `pairwise(nums)` | Adjacent pairs |
| `count()` | Unique increasing tie-breakers |
| `chain.from_iterable(lists)` | Flatten one level |

### Numbers

```python
math.inf, -math.inf          # sentinels for min/max; compare fine with ints
float("inf")                 # same thing
7 // 2 == 3; -7 // 2 == -4   # // floors toward -infinity
int(-7 / 2) == -3            # truncation toward zero (what C/Java do)
-7 % 3 == 2                  # % takes the sign of the divisor
divmod(-17, 5) == (-4, 3)
-(-a // b)                   # ceil(a / b) for ints, no floats
math.isqrt(n)                # exact integer square root
pow(base, exp, mod)          # fast modular exponentiation
MOD = 10**9 + 7
```

Python ints never overflow, but under "modulo 10^9+7" still reduce as you go to keep numbers small and fast. For fixed-width languages, mention overflow in `mid = (lo + hi) // 2` and in products.

## 5. Pitfalls

| Pitfall | What happens | Fix |
|---|---|---|
| `[[0] * n] * m` | All m rows are the **same** list; one write changes every row | `[[0] * n for _ in range(m)]` |
| `def f(x, acc=[])` | Default evaluated once; state leaks across calls | `acc=None`, then `if acc is None: acc = []` |
| Slicing in loops/recursion | `a[1:]` copies O(k); `f(a[1:])` recursion becomes O(n^2) | Pass indices `lo, hi` |
| Recursion depth | Default limit ~1000; DFS on a 10^5-node path crashes | `sys.setrecursionlimit(10**6)` (may still overflow the C stack), or an explicit stack / bottom-up DP |
| String immutability | `s[i] = c` fails; `+=` in a loop is quadratic | `list(s)` + `"".join` |
| Mutating while iterating | Dict: `RuntimeError`. List: silently skips elements | Iterate over `list(d)` or build a new collection |
| Heap of non-comparable objects | Equal priorities make Python compare payloads → `TypeError` | `(priority, next(counter), obj)` |
| Float precision | `0.1 + 0.2 != 0.3`; `int(math.sqrt(10**18 - 1))` can be wrong | `math.isclose`, integer math, `math.isqrt`, cross-multiplied fractions |
| `is` vs `==` | `is` checks identity; small-int/string caching makes it "work" by accident | `==` for values, `is` only for `None` |
| defaultdict reads | `d[k]` inserts missing keys, changing `len(d)` and iteration | `k in d` or `d.get(k)` |
| `sort()` returns `None` | `x = a.sort()` gives `None` | `sorted(a)`, or sort then use `a` |
| `list.remove` / `in list` in a loop | Hidden O(n) inside an O(n) loop | Use a set or dict |
| Shadowing builtins | `list = ...`, `sum = 0`, `max = ...` break later calls | Pick other names |

Converting recursion to iteration (DFS):

```python
def dfs_iterative(graph: dict[int, list[int]], start: int) -> list[int]:
    order, seen, stack = [], set(), [start]
    while stack:
        u = stack.pop()
        if u in seen:
            continue
        seen.add(u)
        order.append(u)
        for v in reversed(graph.get(u, [])):    # reversed: match recursive visit order
            if v not in seen:
                stack.append(v)
    return order
```

For post-order work (children's results first), push `(node, state)` pairs, or push the node twice with a "children visited" flag.

## 6. Interview hygiene

**Structure the code.**
- Write the signature, then the core loop, then helpers. Name helpers for what they answer: `in_bounds(r, c)`, `neighbors(r, c)`, `can_finish_with(capacity)`, `expand(lo, hi)`.
- A 5-line helper isolating a check (e.g., a binary-search feasibility test) makes the main logic obviously correct and is testable alone.
- Names: `left/right` or `lo/hi`, `count`, `seen`, `freq`, `best`, `prefix`, `need`. Avoid `temp`, `arr2`, `flag`. Single letters are fine for loop indices and familiar math (`n`, `i`, `j`, `k`).
- Handle trivial cases at the top (`if not nums: return 0`) only when the general logic doesn't already, and say which.

**Test quickly.**
- Trace one small example aloud, tracking the key variables.
- List edge cases: empty, single element, all equal, all negative, sorted / reverse sorted, duplicates, k = 0 or k > n, max-size input (for complexity, not tracing).
- If you can run code, write asserts:

```python
assert solve([2, 7, 11, 15], 9) == [0, 1]
assert solve([], 0) == []
assert solve([1], 1) == []
```

- When a test fails, reason about the invariant to find the bug; don't make random edits.

**Talk through complexity.**
- Derive it, don't just label it: "Sorting is O(n log n). Then one pass with two pointers is O(n). Total O(n log n) time. Extra space O(1) apart from the sort, which in Python is O(n) for Timsort."
- Name what dominates and what you'd trade: "If memory were tight, I could drop the hash map and sort instead: O(n log n) time, O(1) extra space."
- For amortized or expected bounds, say so and give the one-sentence reason.
- For recursion: "n distinct states, O(1) work each, so O(n) time; recursion depth n, so O(n) stack."

## Self-check

1. Why is `list.append` amortized O(1) but `list.insert(0, x)` O(n)? What would make append O(n) amortized?
2. Constraints say `n ≤ 2 * 10^5`. Which complexities are acceptable, and which patterns does that suggest? What changes if `n ≤ 18`?
3. A memoized function has parameters `(i, j)` with `0 ≤ i, j < n` and an O(n) inner loop. What is the time complexity?
4. How do you build a max-heap of `(priority, task)` where tasks are custom objects and equal priorities must come out in insertion order?
5. What does `-7 // 2` return in Python, and how do you compute `ceil(a / b)` for positive ints without floats?
6. Why is `grid = [[0] * n] * m` wrong, but `[0] * n` alone fine?
7. Your recursive DFS crashes on a 10^5-node chain. Give two fixes and each one's trade-off.
8. Why does `f(arr[1:])` recursion turn an O(n) algorithm into O(n^2)?
9. When would you use `cmp_to_key` instead of `key=`? Give an example.
10. What does `bisect.insort` cost, and why?
