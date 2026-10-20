# 07. Binary Search

Binary search isn't "find a number in a sorted array" but **find the boundary of a monotonic yes/no predicate in O(log range) evaluations.** That covers sorted-array lookups, rotated arrays, peaks, and optimization problems ("minimum capacity", "minimum speed", "minimize the maximum") that search over the answer itself. Most bugs are boundary bugs, so this module uses one convention throughout.

## At a glance

- **Use when:** sorted (or rotated) data with O(log n) required; "minimum X such that" / "minimize the maximum" with a monotonic feasibility check; huge answer range with a cheap check.
- **Templates:**
  - Template 1: first True (the universal one), for any monotonic predicate over indices (also peaks, partitions)
  - Template 2: the `bisect` module, for lower/upper bound lookups in a sorted array
  - Template 3: binary search on the answer, for min/max X with a greedy feasibility check
  - Template 4: floating-point binary search, for real-valued answers within a tolerance
- **Typical complexity:** O(log n) per lookup; O(check · log(range)) on the answer.
- **Top trap:** mixing interval conventions (`lo = mid` with floor mid loops forever).

## Core idea

Find where a monotonic predicate first flips from False to True:

```
 x:     lo                               hi
 pred:  F  F  F  F  F  T  T  T  T  T  T
                       ^
                       first_true
```

**Convention for the whole module: half-open `[lo, hi)`, find first True.**

```python
while lo < hi:
    mid = (lo + hi) // 2
    if pred(mid): hi = mid        # mid may be the answer, keep it in range
    else:         lo = mid + 1    # mid is not the answer, exclude it
return lo                         # lo == hi
```

**Loop invariant:** every index `< lo` is known False; every index `>= hi` is known True (or the "none found" sentinel). `[lo, hi)` is unknown; when it's empty (`lo == hi`), `lo` is the first True.

**Why it terminates:** when `lo < hi`, `mid = (lo + hi) // 2` satisfies `lo <= mid < hi`, so both `hi = mid` and `lo = mid + 1` shrink the interval. Infinite loops come from mixing conventions: with floor mid, `lo = mid` stalls forever when `hi = lo + 1`. Rule: **the branch that keeps `mid` must be the side floor-mid never equals**: the `hi` side.

**Derived searches, all from first-True:**

| Want | Predicate on index i | Result |
|---|---|---|
| lower_bound: first `a[i] >= t` | `a[i] >= t` | `bisect_left` |
| upper_bound: first `a[i] > t` | `a[i] > t` | `bisect_right` |
| last `a[i] <= t` | first `a[i] > t`, then `- 1` | `bisect_right - 1` |
| last `a[i] < t` | first `a[i] >= t`, then `- 1` | `bisect_left - 1` |
| exact match | lower_bound, then check `a[i] == t` | index or -1 |
| count of `t` | `upper_bound - lower_bound` | |
| last x where `feasible(x)` (T T T F F) | first `not feasible(x)`, then `- 1` | max feasible |

**Binary search on the answer.** For "minimum X such that something is achievable", with monotonic achievability (X works ⇒ X+1 works), write `feasible(X)` (usually a greedy O(n) check) and binary search X over `[lowest possible, highest possible]`. **"What is the optimal X?" → "can I do it with X?" → monotonic yes/no → first True.** Cost O(check · log(range)).

## Recognition signals

- Sorted (or rotated sorted) input, n ≥ 10^5, many queries: O(log n) per query.
- "Find first / last occurrence", "insert position", "how many elements ≤ x": lower/upper bound.
- "Minimum capacity / speed / time / size such that ... within D / k / h": binary search on answer.
- "Minimize the maximum" / "maximize the minimum" over a partition or placement: binary search on answer + greedy check.
- Huge answer range (10^9 or 10^18), cheap check per candidate: binary search on answer.
- An explicit O(log n) requirement.
- "Find any peak / local minimum": binary search on slope, no sort needed.
- Real-valued answer with tolerance ("within 10^-6"): floating-point binary search.
- Implicitly sorted structure: matrix with sorted rows/cols, monotonic function f(x), "k-th smallest" in a multiplication table or pair distances (count elements ≤ mid).

## Disguises

1. **Autoscaler sizing.** "Jobs with given CPU demands must be processed in order within D batches. What's the smallest per-batch CPU capacity that finishes on time?"
   Decode: smallest capacity finishing within D, order preserved. Feasibility: greedy batch count ≤ D, monotonic in capacity (LeetCode: Capacity To Ship Packages Within D Days).

2. **Rate limit tuning.** "Each client has a backlog; a worker drains k requests per minute from one client at a time. What's the lowest k that clears all backlogs in h minutes?"
   Decode: minutes needed = sum of ceil(backlog / k), non-increasing in k. First k with minutes ≤ h (Koko Eating Bananas).

3. **Shard balancing.** "Split an ordered list of tables into k contiguous shards to minimize the largest shard size."
   Decode: minimize the maximum over contiguous groups. Guess the max, greedily count shards (Split Array Largest Sum).

4. **Router placement.** "Place k routers at given positions along a hallway to maximize the minimum distance between any two."
   Decode: maximize the minimum. Feasibility: greedily place at the first position ≥ d away; d works if k fit. Last True (LeetCode: Magnetic Force Between Two Balls; also Aggressive Cows).

5. **Rotated log index.** "A ring buffer of sorted timestamps was read from an arbitrary offset. Find a given timestamp in O(log n)."
   Decode: sorted array, rotated by an unknown offset; one half around mid is always sorted (Search in Rotated Sorted Array).

6. **Version bisect.** "Builds 1..n; once a build breaks, all later builds are broken. Find the first broken build with the fewest checks."
   Decode: literally first True over `isBad(version)` (LeetCode: First Bad Version), like `git bisect`.

## Templates and worked examples

All templates are in `code/07_binary_search.py`.

### Template 1: first True (the universal one)

> - **Example problem:** [Find Peak Element (162)](https://leetcode.com/problems/find-peak-element/)
> - **Why this template:** "Return any peak" in O(log n): the first index where "downhill at i" (`nums[i] > nums[i+1]`) holds is a peak.
> - **Worked below:** Example 1.2
> - **Also fits:** [Find Minimum in Rotated Sorted Array (153)](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/)

```python
def first_true(lo: int, hi: int, pred: Callable[[int], bool]) -> int:
    while lo < hi:
        mid = (lo + hi) // 2
        if pred(mid):       # TODO: your monotonic predicate
            hi = mid        # mid might be the answer: keep it
        else:
            lo = mid + 1    # mid is definitely not the answer
    return lo


def lower_bound(a: list[int], target: int) -> int:
    """First index i with a[i] >= target (== bisect.bisect_left)."""
    return first_true(0, len(a), lambda i: a[i] >= target)


def upper_bound(a: list[int], target: int) -> int:
    """First index i with a[i] > target (== bisect.bisect_right)."""
    return first_true(0, len(a), lambda i: a[i] > target)


def binary_search_exact(a: list[int], target: int) -> int:
    i = lower_bound(a, target)
    return i if i < len(a) and a[i] == target else -1
```

Knobs: `lo`, `hi` (half-open; `hi` is the "none" sentinel), and the predicate. `pred(mid)` never sees `mid == hi`, so `hi = len(a)` is safe. In interviews, write the loop inline.

For reference, the closed-interval exact-match form (fine for "found or not", but generalizes less cleanly to boundaries):

```python
def binary_search_exact_closed(a: list[int], target: int) -> int:
    lo, hi = 0, len(a) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if a[mid] == target:
            return mid
        if a[mid] < target:
            lo = mid + 1
        else:
            hi = mid - 1
    return -1
```

#### Example 1.1: Search in a rotated sorted array [Medium]

**Problems:** [Search in Rotated Sorted Array (33)](https://leetcode.com/problems/search-in-rotated-sorted-array/), [Find Minimum in Rotated Sorted Array (153)](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/) · **Template:** Template 1 (first True (the universal one)); closed-interval exact-match form for the search, first-True for the minimum

**Problem.** A sorted array of distinct values was rotated at an unknown pivot. Return a target's index in O(log n), or -1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Sorted but rotated, O(log n). At any `mid`, `[lo, mid]` or `[mid, hi]` is sorted.

**Approach:** Closed interval, exact match (clearest, since the predicate isn't monotonic on index). Compare `nums[lo]` and `nums[mid]` to find the sorted half; search it if the target is in its value range, else the other half.

```python
def search_rotated(nums: list[int], target: int) -> int:
    lo, hi = 0, len(nums) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if nums[mid] == target:
            return mid
        if nums[lo] <= nums[mid]:                    # left half [lo, mid] sorted
            if nums[lo] <= target < nums[mid]:
                hi = mid - 1
            else:
                lo = mid + 1
        else:                                        # right half [mid, hi] sorted
            if nums[mid] < target <= nums[hi]:
                lo = mid + 1
            else:
                hi = mid - 1
    return -1
```

Related, and clean with first-True: a rotated array's minimum is the first index with `nums[i] <= nums[-1]`.

```python
def find_min_rotated(nums: list[int]) -> int:
    last = nums[-1]
    return nums[first_true(0, len(nums), lambda i: nums[i] <= last)]
```

**Complexity:** O(log n) time, O(1) space.

**Edge cases:** `<=` in `nums[lo] <= nums[mid]` matters when `lo == mid` (two elements left). With duplicates (Search in Rotated Sorted Array II), `nums[lo] == nums[mid] == nums[hi]` is uninformative; shrink `lo += 1, hi -= 1`, worst case O(n).

</details>

#### Example 1.2: Find a peak [Medium]

**Problems:** [Find Peak Element (162)](https://leetcode.com/problems/find-peak-element/), [Search a 2D Matrix (74)](https://leetcode.com/problems/search-a-2d-matrix/) · **Template:** Template 1 (first True (the universal one))

**Problem.** Adjacent elements differ. In O(log n), return the index of any element strictly greater than its neighbors (out-of-range neighbors are −∞).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** O(log n) on unsorted data means binary search on a local property: the slope at `i` shows which side must hold a peak.

**Approach:** Downhill at `i` (`nums[i] > nums[i+1]`) means a peak at `i` or to its left; otherwise to the right. Run first-True of "downhill at i" over `[0, n-1)`; if never downhill, the last index is a peak. The predicate isn't globally monotonic, but each step preserves "a peak exists in `[lo, hi]`", which is all binary search needs.

```python
def find_peak_element(nums: list[int]) -> int:
    n = len(nums)
    return first_true(0, n - 1, lambda i: nums[i] > nums[i + 1]) if n > 1 else 0
```

**Complexity:** O(log n) time, O(1) space.

**Edge cases:** Range `[0, n-1)` keeps `nums[i + 1]` valid; the sentinel `n-1` is a valid answer. Plateaus break the slope argument (hence adjacent elements differ).

Related: a 2D sorted matrix (LeetCode: Search a 2D Matrix) is a flattened sorted array, with index `k` at `(k // n, k % n)`.

```python
def search_matrix(matrix: list[list[int]], target: int) -> bool:
    if not matrix or not matrix[0]:
        return False
    m, n = len(matrix), len(matrix[0])
    k = first_true(0, m * n, lambda k: matrix[k // n][k % n] >= target)
    return k < m * n and matrix[k // n][k % n] == target
```



</details>

#### Example 1.3: Median of two sorted arrays [Hard]

**Problem:** [Median of Two Sorted Arrays (4)](https://leetcode.com/problems/median-of-two-sorted-arrays/) · **Template:** Template 1 (first True (the universal one)), applied to the partition index

**Problem.** Given sorted arrays `a` (length m) and `b` (length n), return the median of their union in O(log(min(m, n))).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** O(log) over two sorted arrays: binary search a **partition**, not values.

**Approach:** Cut `a` after `i` elements and `b` after `j = half - i`, where `half = (m + n + 1) // 2` is the left side's size. The cut is correct when every left element ≤ every right one: `a[i-1] <= b[j]` and `b[j-1] <= a[i]`. As `i` grows, `a[i]` rises and `b[j-1]` falls, so `a[i] >= b[j-1]` is monotonic in `i`; the first `i` where it holds satisfies the other condition too (since `i-1` failed). Search the shorter array so `j` stays in range.

```
 a:  a[0] .. a[i-1] | a[i] .. a[m-1]
 b:  b[0] .. b[j-1] | b[j] .. b[n-1]
     left side has half elements; median uses max(left) and min(right)
```

```python
def find_median_sorted_arrays(a: list[int], b: list[int]) -> float:
    if len(a) > len(b):
        a, b = b, a
    m, n = len(a), len(b)
    half = (m + n + 1) // 2
    NEG, POS = float("-inf"), float("inf")

    # first i in [0, m] with a[i] >= b[half - i - 1]  (i == m counts as True)
    def ok(i: int) -> bool:
        j = half - i
        return i == m or j == 0 or a[i] >= b[j - 1]

    i = first_true(0, m + 1, ok)
    j = half - i
    left_max = max(a[i - 1] if i > 0 else NEG, b[j - 1] if j > 0 else NEG)
    if (m + n) % 2:
        return float(left_max)
    right_min = min(a[i] if i < m else POS, b[j] if j < n else POS)
    return (left_max + right_min) / 2
```

**Complexity:** O(log min(m, n)) time, O(1) space.

**Edge cases:** m = 0: range `[0, 1)`, `i = 0`, works. Use ±∞ sentinels for missing neighbors. For odd totals the extra element is on the left, so the median is `left_max`. Generalization: k-th smallest of two sorted arrays (discard k/2 per step).

</details>

### Template 2: the `bisect` module

> - **Example problem:** [Find First and Last Position of Element in Sorted Array (34)](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/)
> - **Why this template:** "First and last position" of a target in a sorted array is exactly `bisect_left` and `bisect_right - 1`.
> - **Worked below:** Example 2.1
> - **Also fits:** [Time Based Key-Value Store (981)](https://leetcode.com/problems/time-based-key-value-store/)

```python
import bisect
bisect.bisect_left(a, x)        # first i with a[i] >= x   (lower_bound)
bisect.bisect_right(a, x)       # first i with a[i] > x    (upper_bound)
bisect.insort(a, x)             # insert keeping sorted: O(n) because of list shift
bisect.bisect_left(a, x, lo, hi)          # restrict range
bisect.bisect_left(a, x, key=lambda r: r[0])  # Python 3.10+: key applies to elements, not x
```

Use `bisect` for array lookups (correct, C-fast); hand-write the loop for predicate searches (rotated arrays, answer search).

#### Example 2.1: First and last position of a value [Medium]

**Problem:** [Find First and Last Position of Element in Sorted Array (34)](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/) · **Template:** Template 2 (the `bisect` module)

**Problem.** In O(log n), return the first and last index of a target in a sorted array with duplicates, or `[-1, -1]`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** First and last in sorted data: lower_bound and upper_bound − 1.

**Approach:** `bisect_left` gives the first index ≥ target; if that slot isn't the target, it's absent. Otherwise the last is `bisect_right - 1`.

```python
def search_range(nums: list[int], target: int) -> list[int]:
    left = bisect.bisect_left(nums, target)
    if left == len(nums) or nums[left] != target:
        return [-1, -1]
    right = bisect.bisect_right(nums, target) - 1
    return [left, right]
```

**Complexity:** O(log n) time, O(1) space.

**Edge cases:** Check `left == len(nums)` before indexing. "Find any match, then walk left/right" degrades to O(n) on all duplicates.

</details>

### Template 3: binary search on the answer

> - **Example problem:** [Koko Eating Bananas (875)](https://leetcode.com/problems/koko-eating-bananas/)
> - **Why this template:** "Minimum integer k" to eat all piles "within h hours", and faster never takes longer: binary search the speed with a feasibility check.
> - **Worked below:** Example 3.1, Example 3.2
> - **Also fits:** [Capacity To Ship Packages Within D Days (1011)](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/)

```python
def search_answer_min(lo: int, hi: int, feasible: Callable[[int], bool]) -> int:
    """Smallest x in [lo, hi] with feasible(x). Caller guarantees feasible(hi)."""
    return first_true(lo, hi + 1, feasible)  # TODO: choose lo/hi so the answer is inside


def search_answer_max(lo: int, hi: int, feasible: Callable[[int], bool]) -> int:
    """Largest x in [lo, hi] with feasible(x), feasible is T...T F...F. Caller guarantees feasible(lo)."""
    return first_true(lo, hi + 1, lambda x: not feasible(x)) - 1
```

Knobs:
1. **Bounds.** `lo` = smallest answer that could work (often `max(items)`: one item must fit). `hi` = an answer that surely works (often `sum(items)` or `max * count`). Bounds must be correct, not tight.
2. **Feasibility check.** Usually greedy ("pack as much as possible before a new group", "place at the earliest legal spot"). It must be monotonic; prove it in one sentence: "a larger capacity can always mimic a smaller one's packing".
3. **Direction.** Minimize X with "X works ⇒ X+1 works": `search_answer_min`. Maximize X with "X works ⇒ X−1 works": `search_answer_max`.

#### Example 3.1: Minimum eating speed [Medium]

**Problem:** [Koko Eating Bananas (875)](https://leetcode.com/problems/koko-eating-bananas/) · **Template:** Template 3 (binary search on the answer)

**Problem.** Each hour you eat up to `k` bananas from one pile. Find the smallest integer `k` that finishes all piles within `h` hours.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Smallest k finishing within h, and faster never takes longer: binary search on the answer.

**Approach:** `hours(k) = sum(ceil(p / k))` is non-increasing in k. Bounds: `k = 1` to `k = max(piles)` (finishes in `len(piles) ≤ h` hours, always feasible).

```python
def min_eating_speed(piles: list[int], h: int) -> int:
    def hours(speed: int) -> int:
        return sum((p + speed - 1) // speed for p in piles)   # ceil division
    return search_answer_min(1, max(piles), lambda s: hours(s) <= h)
```

**Complexity:** O(n log M) time where M = max pile, O(1) space.

**Edge cases:** Integer ceil `(p + k - 1) // k`, not `math.ceil(p / k)` (float error on large values). Lower bound 1, not 0 (division by zero).

</details>

#### Example 3.2: Smallest ship capacity for a deadline [Medium]

**Problem:** [Capacity To Ship Packages Within D Days (1011)](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/) · **Template:** Template 3 (binary search on the answer)

**Problem.** Packages on a conveyor must ship in order. Each day the ship loads from the front until the next package would exceed capacity. Find the smallest capacity that ships everything within `days` days. Example: weights `[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]`, days = 5 → `15` (`[1..5] [6, 7] [8] [9] [10]`).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Smallest capacity finishing within D days, order preserved. A bigger ship never needs more days (monotonic feasibility): binary search on the answer with a greedy check.

**Approach:** For capacity `cap`, load each day until the next package would overflow. Greedy minimizes days (delaying a package never reduces the total). `days_needed(cap)` is non-increasing, so `days_needed(cap) <= days` is F...F T...T and `search_answer_min` finds the first True. Bounds: `max(weights)` (the heaviest must fit) to `sum(weights)` (one day, feasible since `days >= 1`).

```python
def ship_within_days(weights: list[int], days: int) -> int:
    def days_needed(cap: int) -> int:
        d, load = 1, 0
        for w in weights:
            if load + w > cap:
                d, load = d + 1, 0
            load += w
        return d
    return search_answer_min(max(weights), sum(weights), lambda c: days_needed(c) <= days)
```

**Complexity:** O(n log S) time for S = `sum(weights)` (O(log S) probes × O(n) simulation), O(1) extra space.

**Edge cases:** `lo` must be `max(weights)`, not 1 or 0: with a smaller cap, a package that never fits silently counts as its own overweight day. Start the day counter at 1. Same code as Split Array Largest Sum (Example 3.3), with days as parts.

</details>

#### Example 3.3: Minimize the largest part [Hard]

**Problems:** [Split Array Largest Sum (410)](https://leetcode.com/problems/split-array-largest-sum/), [Minimum Time to Complete Trips (2187)](https://leetcode.com/problems/minimum-time-to-complete-trips/) · **Template:** Template 3 (binary search on the answer)

**Problem.** Split a non-negative integer array into `k` non-empty contiguous parts to minimize the largest part sum.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Minimize the maximum over contiguous partitions: guess the cap, check if k parts suffice.

**Approach:** For cap C, extend the current part until the next element would exceed C. Greedy minimizes parts (any partition can shift elements into earlier parts without adding parts). Parts needed is non-increasing in C, so "parts_needed(C) ≤ k" is monotonic. Fewer than k parts is fine: split further without raising the max (k ≤ n). Bounds: `max(nums)` (the largest element must fit) to `sum(nums)` (one part).

```python
def split_array(nums: list[int], k: int) -> int:
    def parts_needed(cap: int) -> int:
        parts, cur = 1, 0
        for x in nums:
            if cur + x > cap:          # greedy: start a new part only when forced
                parts, cur = parts + 1, 0
            cur += x
        return parts
    return search_answer_min(max(nums), sum(nums), lambda cap: parts_needed(cap) <= k)
```

**Complexity:** O(n log S) time where S = sum, O(1) space. (DP is O(k n^2).)

**Edge cases:** `lo = max(nums)`, not 0, or the check must handle an element larger than the cap. Same code as shipping packages and "painter's partition": one family.

**Variants of "minimum time to finish jobs on k machines":**
- Parallel machines, each repeating jobs at its own rate (LeetCode: Minimum Time to Complete Trips): feasibility is `sum(T // t) >= total`, O(n) per check.

  ```python
  def minimum_time_trips(time: list[int], total_trips: int) -> int:
      hi = min(time) * total_trips                       # fastest machine alone suffices
      return search_answer_min(1, hi, lambda T: sum(T // t for t in time) >= total_trips)
  ```

- Whole jobs, any order, k workers (LeetCode: Find Minimum Time to Finish All Jobs): "do all jobs fit with max load L?" is bin packing (NP-hard), so feasibility is backtracking with pruning (`min_time_assign_jobs` in the code file; see `15-backtracking.md`), inside a binary search on L.
- Contiguous jobs in order: identical to Split Array Largest Sum.

</details>

### Template 4: floating-point binary search

> - **Example problem:** [Separate Squares I (3453)](https://leetcode.com/problems/separate-squares-i/)
> - **Why this template:** Minimum real y where area above equals area below (within 10^-5); area below grows with y, so bisect over reals.
> - **Worked below:** Example 4.1

```python
def search_float(lo: float, hi: float, pred: Callable[[float], bool], iters: int = 100) -> float:
    for _ in range(iters):
        mid = (lo + hi) / 2
        if pred(mid):
            hi = mid
        else:
            lo = mid
    return hi
```

Use a fixed iteration count, not `while hi - lo > eps`: at large magnitudes `eps` can be below the spacing between doubles, so the loop never ends. 100 iterations shrink any interval by 2^100. `lo = mid` is fine here (no integer rounding).

#### Example 4.1: Horizontal line splitting square areas in half [Medium]

**Problem:** [Separate Squares I (3453)](https://leetcode.com/problems/separate-squares-i/) · **Template:** Template 4 (floating-point binary search)

**Problem.** Each square is `[x, y, l]`: bottom-left corner `(x, y)`, side `l`. Find the minimum real y-coordinate of a horizontal line such that the total square area above it equals the total area below it (an overlapping region counts once for each square covering it). Answers within 10^-5 are accepted. Example: `[[0, 0, 1], [2, 2, 1]]` → 1.0 (any line in [1, 2] splits 1 / 1; the minimum is 1).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "Minimum real y" + "within 10^-5" + a quantity that only grows as the line rises → floating-point binary search on y.

**Approach:** A square contributes `l * clamp(Y - y, 0, l)` below the line, so `area_below(Y)` is continuous and non-decreasing. Area above = total - below, so "above equals below" is `2 * area_below(Y) == total`. Search for the first Y where `2 * area_below(Y) >= total`: first-true gives the minimum when the area is flat over a gap. Bounds: lowest bottom edge (below = 0) to highest top edge (below = total). Use a fixed iteration count, not `while hi - lo > eps`: with coordinates near 2·10^9, adjacent doubles are about 2·10^-7 apart, so a tight epsilon like 10^-9 can never be reached and the loop spins forever. 100 halvings of a range up to ~2·10^9 go far past double precision, and the cost is predictable.

```python
def separate_squares(squares: list[list[int]]) -> float:
    total = sum(l * l for _, _, l in squares)

    def area_below(y: float) -> float:
        return sum(l * min(max(y - b, 0), l) for _, b, l in squares)

    lo = min(b for _, b, _ in squares)
    hi = max(b + l for _, b, l in squares)
    return search_float(lo, hi, lambda y: 2 * area_below(y) >= total)
```

**Complexity:** O(n · I) time with I = 100 iterations, O(1) extra space. (Exact alternative: sort the 2n edges and sweep the piecewise-linear area, O(n log n).)

**Edge cases:** Gaps where the area is flat: return the left end (the minimum). One square → its midpoint. Compare `2 * below >= total` to avoid halving. x is irrelevant. Separate Squares II (union area, overlaps once) needs a sweep line instead.

</details>

## Template cheat sheet

Every template from this module in one place, for warm-ups and review.

### Template 1: first True (the universal one)

```python
def first_true(lo: int, hi: int, pred: Callable[[int], bool]) -> int:
    while lo < hi:
        mid = (lo + hi) // 2
        if pred(mid):       # TODO: your monotonic predicate
            hi = mid        # mid might be the answer: keep it
        else:
            lo = mid + 1    # mid is definitely not the answer
    return lo


def lower_bound(a: list[int], target: int) -> int:
    """First index i with a[i] >= target (== bisect.bisect_left)."""
    return first_true(0, len(a), lambda i: a[i] >= target)


def upper_bound(a: list[int], target: int) -> int:
    """First index i with a[i] > target (== bisect.bisect_right)."""
    return first_true(0, len(a), lambda i: a[i] > target)


def binary_search_exact(a: list[int], target: int) -> int:
    i = lower_bound(a, target)
    return i if i < len(a) and a[i] == target else -1
```

```python
def binary_search_exact_closed(a: list[int], target: int) -> int:
    lo, hi = 0, len(a) - 1
    while lo <= hi:
        mid = (lo + hi) // 2
        if a[mid] == target:
            return mid
        if a[mid] < target:
            lo = mid + 1
        else:
            hi = mid - 1
    return -1
```

### Template 2: the `bisect` module

```python
import bisect
bisect.bisect_left(a, x)        # first i with a[i] >= x   (lower_bound)
bisect.bisect_right(a, x)       # first i with a[i] > x    (upper_bound)
bisect.insort(a, x)             # insert keeping sorted: O(n) because of list shift
bisect.bisect_left(a, x, lo, hi)          # restrict range
bisect.bisect_left(a, x, key=lambda r: r[0])  # Python 3.10+: key applies to elements, not x
```

### Template 3: binary search on the answer

```python
def search_answer_min(lo: int, hi: int, feasible: Callable[[int], bool]) -> int:
    """Smallest x in [lo, hi] with feasible(x). Caller guarantees feasible(hi)."""
    return first_true(lo, hi + 1, feasible)  # TODO: choose lo/hi so the answer is inside


def search_answer_max(lo: int, hi: int, feasible: Callable[[int], bool]) -> int:
    """Largest x in [lo, hi] with feasible(x), feasible is T...T F...F. Caller guarantees feasible(lo)."""
    return first_true(lo, hi + 1, lambda x: not feasible(x)) - 1
```

### Template 4: floating-point binary search

```python
def search_float(lo: float, hi: float, pred: Callable[[float], bool], iters: int = 100) -> float:
    for _ in range(iters):
        mid = (lo + hi) / 2
        if pred(mid):
            hi = mid
        else:
            lo = mid
    return hi
```

## Variations and extensions

- **K-th smallest via counting** (sorted matrix, multiplication table, pair distances): binary search on value v; `count(≤ v)` is monotonic, and the first v with count ≥ k is the answer (always an actual element).
- **Maximize the minimum**: placement problems (Magnetic Force Between Two Balls), "maximum minimum sweetness" (Divide Chocolate). Use `search_answer_max`.
- **Binary search + sliding window**: longest L satisfying X, where checking one L is a fixed window (`06-sliding-window.md`).
- **Binary search + BFS/DSU**: "minimum effort path", "swim in rising water": binary search the threshold, check connectivity (`13-graphs-traversal-and-topological-sort.md`, `14-union-find-shortest-paths-mst.md`).
- **Binary search inside DP**: Longest Increasing Subsequence in O(n log n) uses `bisect_left` on a tails array (`17-dynamic-programming-advanced.md`).
- **Ternary search** for unimodal real functions; for integers, prefer binary search on the slope `f(x) < f(x+1)` (same as peak finding).
- **Unknown size** (Search in a Sorted Array of Unknown Size): double `hi` until `a[hi] >= target` (exponential search), then binary search.
- **Follow-ups**: duplicates? (rotated: O(n) worst case); huge array on disk? (B-tree / block-wise search); many queries? (sort once, O(log n) each).

## Common pitfalls

- Mixing conventions: half-open init (`hi = len(a)`) with `while lo <= hi`, or `hi = mid - 1` in a first-True search. Pick one.
- `lo = mid` with floor mid (infinite loop when `hi = lo + 1`). Fix with ceil mid `(lo + hi + 1) // 2` or, more cleanly, first-True of the negated predicate minus 1.
- Non-monotonic predicate. Ask first: "if X works, does X+1 always work?" If not, binary search on the answer is invalid.
- Wrong bounds: `lo` too high skips the answer; an infeasible `hi` comes back as a false answer. Ensure `feasible(hi)` or handle "no answer".
- Overflow isn't a Python issue; mention `lo + (hi - lo) // 2` for other languages.
- Forgetting to check the found index (`i < len(a) and a[i] == target`) for exact matches.
- Float search with `while hi - lo > 1e-9` on large values may never terminate. Use fixed iterations.
- A greedy feasibility check that isn't actually optimal (e.g., greedily assigning arbitrary-order jobs to k machines).

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| Binary search on answer | DP | Feasibility not monotonic, or the check must optimize over choices with no greedy |
| Binary search on sorted array | Two pointers | Pair/triple sums in a sorted array: two pointers is O(n) vs O(n log n) (`05-two-pointers.md`) |
| Binary search on answer | Heap / greedy directly | "Minimum number of meeting rooms", "k-th largest element": a direct structure is simpler |
| Binary search in 2D | Staircase walk | Rows and columns sorted independently but rows don't chain (Search a 2D Matrix II): O(m + n) from the top-right corner |
| Searching unsorted data | Hash map | Exact lookups, order irrelevant: O(1) hash |
| Rotated array search | Linear scan | With many duplicates the worst case is O(n) anyway; say so |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Search in Rotated Sorted Array (33)](https://leetcode.com/problems/search-in-rotated-sorted-array/) · Medium · Rotated array: one half always sorted.
2. [Find Minimum in Rotated Sorted Array (153)](https://leetcode.com/problems/find-minimum-in-rotated-sorted-array/) · Medium · Rotated array, boundary: first `nums[i] <= nums[-1]`.
3. [Find First and Last Position of Element in Sorted Array (34)](https://leetcode.com/problems/find-first-and-last-position-of-element-in-sorted-array/) · Medium · lower/upper bound: two boundary searches.
4. [Koko Eating Bananas (875)](https://leetcode.com/problems/koko-eating-bananas/) · Medium · Binary search on the answer: ceil-sum feasibility.
5. [Capacity To Ship Packages Within D Days (1011)](https://leetcode.com/problems/capacity-to-ship-packages-within-d-days/) · Medium · Binary search on the answer + greedy check: greedy day count.

### Optional

**Classic / bounds**
- [Binary Search (704)](https://leetcode.com/problems/binary-search/) · Easy · Exact match.
- [Search Insert Position (35)](https://leetcode.com/problems/search-insert-position/) · Easy · `bisect_left`.
- [First Bad Version (278)](https://leetcode.com/problems/first-bad-version/) · Easy · First-True over an API.
- [Sqrt(x) (69)](https://leetcode.com/problems/sqrtx/) · Easy · Last r with r*r ≤ x.
- [Time Based Key-Value Store (981)](https://leetcode.com/problems/time-based-key-value-store/) · Medium · `bisect_right - 1` on timestamps per key.

**Rotated / peak / 2D**
- [Find Peak Element (162)](https://leetcode.com/problems/find-peak-element/) · Medium · Slope predicate.
- [Search a 2D Matrix (74)](https://leetcode.com/problems/search-a-2d-matrix/) · Medium · Flatten index.
- [Search a 2D Matrix II (240)](https://leetcode.com/problems/search-a-2d-matrix-ii/) · Medium · Staircase, not binary search.

**Binary search on the answer**
- [Minimum Time to Complete Trips (2187)](https://leetcode.com/problems/minimum-time-to-complete-trips/) · Medium · `sum(T // t)`.
- [Magnetic Force Between Two Balls (1552)](https://leetcode.com/problems/magnetic-force-between-two-balls/) · Medium · Maximize the minimum gap.
- [Kth Smallest Element in a Sorted Matrix (378)](https://leetcode.com/problems/kth-smallest-element-in-a-sorted-matrix/) · Medium · Binary search on value with count ≤ mid.
- [Split Array Largest Sum (410)](https://leetcode.com/problems/split-array-largest-sum/) · Hard · Minimize the maximum part.
- [Find Minimum Time to Finish All Jobs (1723)](https://leetcode.com/problems/find-minimum-time-to-finish-all-jobs/) · Hard · Binary search + backtracking feasibility.

**Partition**
- [Median of Two Sorted Arrays (4)](https://leetcode.com/problems/median-of-two-sorted-arrays/) · Hard · Binary search the cut in the shorter array.

## Self-check

1. State the invariant of the half-open first-True loop. Why does it always terminate?
2. Express "last index with `a[i] <= t`" using only first-True.
3. In Split Array Largest Sum, why is `parts_needed(cap) <= k` monotonic in `cap`, and why are the bounds `max(nums)` and `sum(nums)`?
4. Why does `lo = mid` with floor mid loop forever, and what are the two fixes?
5. In Find Peak Element, the predicate isn't monotonic. Why does binary search still work?
6. In Median of Two Sorted Arrays, why search the shorter array, and what makes the partition condition monotonic in `i`?
7. Why use a fixed iteration count for floating-point binary search?
