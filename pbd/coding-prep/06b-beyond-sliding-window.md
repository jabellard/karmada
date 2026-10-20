# 06b. Beyond the Sliding Window (companion to module 06)

The sliding window solves "best contiguous subarray under a constraint" in O(n), but only when validity is **monotonic**: growing a window only pushes it one way. Negative numbers, exact targets, divisibility, and sign-flipping products break that, and the window gives silently wrong answers. This companion covers what to use instead. The core move is to **rewrite every subarray as a difference of two prefix sums** and then pick the data structure that finds the best earlier prefix. Code: `code/06b_beyond_sliding_window.py`.

## At a glance

- **Use when:** contiguous subarray + sum/count/longest/shortest, but values can be **negative** (or zero, for products), the target is **exact** or **modular**, or the problem asks for a max/min sum with no bound.
- **Templates:**
  - Template 1: prefix sum + hash map, for an **exact** target (sum = k, sum % m = r, balance = 0): count, longest, shortest
  - Template 2: Kadane, for **max/min subarray sum** (and product, circular, one deletion) with no target
  - Template 3: prefix sums + increasing monotonic deque, for **shortest with sum ≥ k**, or max sum with **length in [lo, hi]**
  - Template 4: prefix sums + decreasing monotonic stack, for **longest with sum ≥ k**
  - Template 5: prefix sums + sorted prefixes, for **max sum ≤ k** / closest to k
  - Template 6: prefix sums + Fenwick tree (or merge sort), to **count** subarrays with sum ≥ k or in a range [lo, hi]
  - Template 7: binary search on the answer + transformed prefix check, for **max average / median** with a length constraint
- **Typical complexity:** O(n) for Templates 1–4, O(n log n) for 5–7 (7 is O(n log(range/ε))).
- **Top trap:** reaching for a window because the problem says "contiguous" without checking that every value is non-negative.

## Core idea

### Why the window breaks

A window drops a left endpoint **forever** when the window turns invalid. That's safe only if no later right endpoint could make that left endpoint valid again. With non-negative values, the sum only grows as `right` moves, so "too big now" means "too big forever". With a negative value, it doesn't:

```
nums = [5, -4, 3],  "longest subarray with sum <= 4"

right=0: [5] sum 5 > 4  -> window drops left=0 forever
right=2: window [-4, 3] -> answers 2
truth:   [5, -4, 3] sum 4 <= 4 -> answer 3   (left=0 was needed again)
```

Shortest with sum ≥ k breaks the other way: dropping a **negative** from the left *raises* the sum, so shrinking no longer moves toward invalid.

### The rewrite that always works

Let `P[0] = 0` and `P[j] = nums[0] + ... + nums[j-1]`. Then

```
sum(nums[i..j-1]) = P[j] - P[i]        for 0 <= i < j <= n
```

A subarray is a **pair of prefix indices `i < j`**. Scan `j` left to right; each problem becomes: *which earlier `P[i]` do I need, and what structure finds it fast?*

| Goal for end `j` | Need from earlier prefixes | Structure | Template |
|---|---|---|---|
| sum = k (count / longest / shortest) | `P[i] == P[j] - k` | hash map: count / first index / last index | 1 |
| max sum | smallest `P[i]` | one running minimum (Kadane) | 2 |
| shortest with sum ≥ k | the latest `i` with `P[i] ≤ P[j] - k` | increasing deque | 3 |
| max sum, length in [lo, hi] | min `P[i]` over a sliding range of `i` | increasing deque (sliding min) | 3 |
| longest with sum ≥ k | the earliest `i` with `P[i] ≤ P[j] - k` | decreasing stack of prefix minima | 4 |
| max sum ≤ k | smallest `P[i] ≥ P[j] - k` | sorted list + bisect | 5 |
| count with sum in [lo, hi] | how many `P[i]` lie in `[P[j] - hi, P[j] - lo]` | Fenwick tree over compressed values | 6 |
| max average ≥ x? | max sum of `(a - x)` ≥ 0 | binary search on x + running min | 7 |

### Can the window be rescued?

Check these before abandoning the window:

- **All values ≥ 0** (zeros are fine for sums): the window works. Exact targets on non-negative values: `atMost(k) − atMost(k−1)` (module 06, Template 4).
- **Products of positive values:** the window works (product only grows). With **zeros or negatives**, split at zeros and track sign parity (Example B.1).
- **Adding a constant to make everything positive does not work.** It adds `c × length` to every sum, so the target changes per subarray.

## Recognition signals

- "Contiguous subarray" + "sum" and the constraints say `-10^4 <= nums[i]` (read the constraints line, not just the story).
- "Exactly k", "equal to k", "divisible by k", "multiple of k", "equal number of X and Y": exact or modular target → Template 1.
- "Maximum sum", "largest total", "best stretch" with no target → Template 2.
- "Shortest ... at least k" with negatives → Template 3. "Longest ... at least k" (or "more X than Y") → Template 4.
- "Length at least L and at most R" with negatives → Template 3 (sliding min of prefixes).
- "No larger than k", "closest to k" → Template 5.
- "How many subarrays have sum in [lo, hi]" or "more X than Y" counts → Template 6.
- "Maximum average with length ≥ k", "median ≥ x" → Template 7.
- A 2D grid asking for a rectangle with sum = / ≤ target: fix a pair of rows, compress columns into a 1D array, then use 1 or 5 (Combined examples).
- Binary or category arrays ("more ones than zeros", "hours > 8"): map to ±1 first, then it's a sum problem with negatives.

## Disguises

1. **P&L streaks.** "A trading desk logs daily profit or loss. How many contiguous periods netted exactly $D?"
   Decode: losses are negatives, "exactly" → prefix sum + count map (Template 1).

2. **Balanced traffic.** "Each second a load balancer logs +1 for a request routed east and −1 for west. Find the longest period where traffic was perfectly balanced."
   Decode: longest subarray with sum 0 → first-index map (Template 1, Example 1.2's shape).

3. **Batch packing.** "Job sizes are in a list; you can remove one contiguous run of jobs so the remaining total divides evenly into p workers. Minimize the run you remove."
   Decode: modular target, shortest → last-index map on `prefix % p` (Example 1.4).

4. **Best trading window.** "Given daily price changes (positive or negative), find the best contiguous holding period."
   Decode: max subarray sum, no target → Kadane (Template 2). It's also [Best Time to Buy and Sell Stock (121)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock/) in prefix form.

5. **Ring buffer revenue.** "Shops are arranged in a circle; pick a contiguous arc maximizing net revenue."
   Decode: circular max subarray → max(Kadane max, total − Kadane min) (Example 2.2).

6. **Minimum burst.** "Per-minute net memory change (allocations minus frees). Find the shortest stretch whose net growth is at least K MB."
   Decode: frees are negatives, shortest with sum ≥ K → increasing deque (Template 3).

7. **SLA health.** "Each hour is good (+1) or bad (−1). Find the longest stretch with strictly more good hours than bad."
   Decode: ±1 transform, longest with sum ≥ 1 → Template 4 (Example 4.2).

8. **Budget cap.** "Find the most profitable contiguous run of quarters whose total doesn't exceed the cap B."
   Decode: max sum ≤ B with negatives → sorted prefixes (Template 5).

9. **Tolerance audit.** "Count the contiguous windows of sensor drift whose total stays within [lo, hi]."
   Decode: count pairs of prefixes with difference in a range → Fenwick (Template 6).

## Templates and worked examples

All templates are in `code/06b_beyond_sliding_window.py`. Knobs shared by all: the **prefix key** (raw sum, transformed sum, remainder), the **query** for end `j`, the **structure** holding earlier prefixes, and the **order**: query before inserting `P[j]`, so a subarray is never empty.

### Choosing a template

| Question | Template | Structure stores |
|---|---|---|
| Count subarrays with sum = k | 1 | prefix → count, seeded `{0: 1}` |
| Longest with sum = k (or balance 0) | 1 | prefix → first index, seeded `{0: -1}` |
| Shortest with sum = k (or remainder r) | 1 | prefix → last index, seeded `{0: -1}` |
| Max / min sum, no target | 2 | running best-ending-here |
| Shortest with sum ≥ k | 3 | deque of starts, prefixes increasing |
| Max sum with lo ≤ length ≤ hi | 3 | deque of starts in `[j-hi, j-lo]` |
| Longest with sum ≥ k | 4 | stack of prefix minima, scan ends right to left |
| Max sum ≤ k / closest to k | 5 | sorted list of prefixes |
| Count with sum ≥ k or in [lo, hi] | 6 | Fenwick over compressed prefixes |
| Max average / median with length ≥ k | 7 | binary search + running min |

### Template 1: prefix sum + hash map

> - **Example problem:** [Subarray Sum Equals K (560)](https://leetcode.com/problems/subarray-sum-equals-k/)
> - **Why this template:** an exact target turns "find a start" into "look up one value": `P[i] = P[j] - k`. No monotonicity needed.
> - **Worked below:** Examples 1.1–1.5
> - **Also fits:** [Subarray Sums Divisible by K (974)](https://leetcode.com/problems/subarray-sums-divisible-by-k/), [Binary Subarrays With Sum (930)](https://leetcode.com/problems/binary-subarrays-with-sum/)

```python
def count_sum_k(nums: list[int], k: int) -> int:
    seen: dict[int, int] = defaultdict(int)
    seen[0] = 1                         # TODO: seed the empty prefix ({0: 1} count, {0: -1} index)
    prefix = 0
    total = 0
    for x in nums:
        prefix += x                     # TODO: transform x first if needed (0 -> -1, x % m, ...)
        total += seen[prefix - k]       # TODO: look up the complement BEFORE inserting
        seen[prefix] += 1               # TODO: store count / first index / last index
    return total
```

Knobs: **what the map stores** decides the question. Counts answer "how many". The **first** index answers "longest" (the earliest start is the farthest away). The **last** index answers "shortest". The seed stands for the empty prefix `P[0]` so subarrays starting at index 0 are found.

#### Example 1.1: Count subarrays summing to k [Medium]

**Problem:** [Subarray Sum Equals K (560)](https://leetcode.com/problems/subarray-sum-equals-k/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** Given integers (possibly negative) and `k`, count the contiguous subarrays whose sum is exactly `k`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Exact sum + negatives → not a window. Count → map stores counts.

**Approach:** Template 1 as written. Each earlier prefix equal to `prefix - k` is one subarray ending here.

```python
def subarray_sum(nums: list[int], k: int) -> int:
    seen: dict[int, int] = defaultdict(int)
    seen[0] = 1
    prefix = total = 0
    for x in nums:
        prefix += x
        total += seen[prefix - k]
        seen[prefix] += 1
    return total
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** `k = 0` with zeros: `[0, 0]` has 3 subarrays; looking up before inserting keeps a prefix from pairing with itself (an empty subarray). Missing `{0: 1}` misses every subarray that starts at index 0.

</details>

#### Example 1.2: Longest balanced binary subarray [Medium]

**Problem:** [Contiguous Array (525)](https://leetcode.com/problems/contiguous-array/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** Given a binary array, return the length of the longest contiguous subarray with an equal number of 0s and 1s.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "Equal number of X and Y" → map 0 to −1, then it's longest subarray with sum 0. Longest → keep the first index of each balance.

**Approach:** If the balance at `i` was seen before at `first[b]`, the subarray between them sums to 0. Never overwrite the first index.

```python
def find_max_length(nums: list[int]) -> int:
    first = {0: -1}                     # balance -> earliest index it was seen
    balance = best = 0
    for i, x in enumerate(nums):
        balance += 1 if x == 1 else -1
        if balance in first:
            best = max(best, i - first[balance])
        else:
            first[balance] = i          # keep the FIRST index only: longest
    return best
```

**Complexity:** O(n) time, O(n) space (an array of size 2n+1 works too, since the balance is in [−n, n]).

**Edge cases:** All one value → 0. The `{0: -1}` seed handles a balanced prefix (`[0, 1]` → 2).

</details>

#### Example 1.3: Subarray sum divisible by k, length at least 2 [Medium]

**Problem:** [Continuous Subarray Sum (523)](https://leetcode.com/problems/continuous-subarray-sum/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** Given non-negative integers and `k ≥ 1`, return whether some subarray of length ≥ 2 has a sum that's a multiple of `k`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "Multiple of k" is a modular target: `P[j] ≡ P[i] (mod k)`. The values are non-negative, but divisibility isn't monotonic (adding more can make it divisible again), so a window still fails.

**Approach:** Key the map on `prefix % k` and store the first index. A repeat remainder at distance ≥ 2 is the answer.

```python
def check_subarray_sum(nums: list[int], k: int) -> bool:
    first = {0: -1}                     # remainder -> earliest index
    rem = 0
    for i, x in enumerate(nums):
        rem = (rem + x) % k
        if rem in first:
            if i - first[rem] >= 2:
                return True
        else:
            first[rem] = i
    return False
```

**Why "a repeat remainder at distance ≥ 2" is the answer:**

1. **A repeated remainder means a valid subarray.** `rem` at index `i` is the running total through `nums[i]`, mod k. Suppose the same remainder shows up at an earlier index `a` and at the current index `i`. Then the two running totals differ by a multiple of k (same clock position). That difference is the sum of the elements **after** `a` up to and including `i`:

   ```
   sum(nums[a+1 .. i]) = P(i) − P(a)  → a multiple of k
   ```

2. **The distance is the subarray's length.** `nums[a+1 .. i]` contains `i − a` elements, so "distance" just means "how long the subarray is":

   ```
   index:      -1     0     1     2
   nums:            23     2     4
   rem:         0     5     1     5
                      ^───────────^    same remainder 5 at a = 0 and i = 2
                                       subarray nums[1..2] = [2, 4], length 2 − 0 = 2, sum 6 ✓
   ```

3. **Why "≥ 2".** The problem requires at least two elements. A distance of 1 means a single element that's a multiple of k on its own, like `[6]` with `k = 6`, which doesn't count.

**Why keep the first index:** when a remainder repeats at distance 1, the code doesn't overwrite the stored index (the `else`). The earlier index might still give distance 2 one step later. `nums = [5, 0, 0]`, `k = 3` has remainders 2, 2, 2:

| i | rem | first[2] | distance | result |
|---|---|---|---|---|
| 0 | 2 | record 0 | | |
| 1 | 2 | stays 0 | 1 − 0 = 1 → too short | keep going |
| 2 | 2 | stays 0 | 2 − 0 = **2** | **True**: `[0, 0]` sums to 0 |

Overwriting `first[2] = 1` at i = 1 would make the distance at i = 2 equal to 1, and the function would wrongly return False. The earliest index gives the longest subarray, so it's the best chance to reach length 2.

**Why the `{0: -1}` seed:** before the array starts, the running total is 0, at "index −1". That lets a subarray starting at index 0 be found: for `[2, 4]` with `k = 6`, remainder 0 at i = 1 matches the seed, distance `1 − (−1) = 2` → True.

**Complexity:** O(n) time, O(min(n, k)) space.

**Edge cases:** `[0, 0]` → True (0 is a multiple of every k). In Java/C++, normalize negative remainders with `((r % k) + k) % k`; Python's `%` already returns a non-negative result for positive k.

</details>

#### Example 1.4: Shortest removal to make the sum divisible [Medium]

**Problem:** [Make Sum Divisible by P (1590)](https://leetcode.com/problems/make-sum-divisible-by-p/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** Remove the shortest contiguous subarray (not the whole array) so the remaining sum is divisible by `p`. Return its length, or −1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** The removed subarray must have sum ≡ `total % p` (call it `need`). Exact modular target, shortest → store the **last** index of each remainder.

**Approach:** For end `i` with remainder `rem`, a start works if its remainder is `(rem - need) % p`. The latest such start gives the shortest subarray.

```python
def min_subarray_to_remove(nums: list[int], p: int) -> int:
    need = sum(nums) % p
    if need == 0:
        return 0
    last = {0: -1}                      # remainder -> LATEST index: shortest
    rem = 0
    best = len(nums)
    for i, x in enumerate(nums):
        rem = (rem + x) % p
        want = (rem - need) % p
        if want in last:
            best = min(best, i - last[want])
        last[rem] = i
    return best if best < len(nums) else -1
```

**Complexity:** O(n) time, O(min(n, p)) space.

**Edge cases:** Already divisible → 0. Removing everything isn't allowed, so a best of `n` means −1. Overwriting `last[rem]` on every step is the whole point (contrast with Example 1.3).

</details>

#### Example 1.5: Count subarrays with median k [Hard]

**Problem:** [Count Subarrays With Median K (2488)](https://leetcode.com/problems/count-subarrays-with-median-k/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** `nums` is a permutation of `1..n`. Count subarrays whose median (lower middle for even length) equals `k`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** The subarray must contain `k`. Map each value to +1 (> k), −1 (< k), 0 (= k). `k` is the median exactly when the balance is 0 (odd length) or 1 (even length). That's an exact target with negatives, anchored at `k`'s position.

**Approach:** Count balances of every left part `nums[i..pos]`. For each right part `nums[pos..j]` with balance `b`, add left parts with balance `-b` or `1 - b`. It's Template 1 with two complements and the map built from one side.

```python
def count_subarrays_median_k(nums: list[int], k: int) -> int:
    pos = nums.index(k)
    left_counts: dict[int, int] = defaultdict(int)
    balance = 0
    for i in range(pos, -1, -1):        # balances of nums[i..pos]
        if nums[i] > k:
            balance += 1
        elif nums[i] < k:
            balance -= 1
        left_counts[balance] += 1
    total = balance = 0
    for j in range(pos, len(nums)):     # balances of nums[pos..j]
        if nums[j] > k:
            balance += 1
        elif nums[j] < k:
            balance -= 1
        total += left_counts[-balance] + left_counts[1 - balance]
    return total
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** `[k]` alone counts (both balances 0). `k`'s own 0 contributes to both sides, but that's fine because each side starts at `pos`, so it's counted once per pair.

</details>

### Template 2: Kadane

> - **Example problem:** [Maximum Subarray (53)](https://leetcode.com/problems/maximum-subarray/)
> - **Why this template:** with no target, the best start for end `j` is just the smallest earlier prefix. Kadane keeps that implicitly as "best sum ending here".
> - **Worked below:** Examples 2.1–2.4
> - **Also fits:** [Best Time to Buy and Sell Stock (121)](https://leetcode.com/problems/best-time-to-buy-and-sell-stock/) (Kadane on daily differences)

```python
def kadane(nums: list[int]) -> int:
    best = cur = nums[0]
    for x in nums[1:]:
        cur = max(x, cur + x)           # TODO: extend or restart (add more states for variants)
        best = max(best, cur)           # TODO: record
    return best
```

Two equivalent views: `cur = max(x, cur + x)` ("a negative running sum only hurts, so restart"), or `P[j] - min(P[0..j-1])`. Variants add states: a running min for circular, a min product for signs, a "deleted already" state for one deletion.

#### Example 2.1: Maximum subarray sum [Medium]

**Problem:** [Maximum Subarray (53)](https://leetcode.com/problems/maximum-subarray/) · **Template:** Template 2 (Kadane)

**Problem.** Return the largest sum of a non-empty contiguous subarray.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Max sum, no target, negatives → Kadane.

**Approach:** Template 2 as written.

```python
def max_sub_array(nums: list[int]) -> int:
    best = cur = nums[0]
    for x in nums[1:]:
        cur = max(x, cur + x)
        best = max(best, cur)
    return best
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** All negative → the largest single element. Initializing `best = 0` returns 0 there, which is wrong for a non-empty subarray. Follow-up "return the indices": record `start` when restarting.

</details>

#### Example 2.2: Maximum circular subarray [Medium]

**Problem:** [Maximum Sum Circular Subarray (918)](https://leetcode.com/problems/maximum-sum-circular-subarray/) · **Template:** Template 2 (Kadane)

**Problem.** The array is circular. Return the max sum of a non-empty subarray, where each element is used at most once.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** A wrapping subarray is the whole array minus a non-wrapping middle. So the best wrap = total − (min subarray sum). Run Kadane for max and min together.

**Approach:** Answer = max(best_max, total − best_min). If every value is negative, `total − best_min` is the empty subarray (0), so return best_max.

```python
def max_subarray_sum_circular(nums: list[int]) -> int:
    best_max = cur_max = nums[0]
    best_min = cur_min = nums[0]
    for x in nums[1:]:
        cur_max = max(x, cur_max + x)
        best_max = max(best_max, cur_max)
        cur_min = min(x, cur_min + x)
        best_min = min(best_min, cur_min)
    if best_max < 0:                    # all negative: wrap would be the empty subarray
        return best_max
    return max(best_max, sum(nums) - best_min)
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** All negative (`[-3, -2, -3]` → −2). Doubling the array and running a window doesn't work (negatives again). Doubling plus Example 3.2's deque (lengths 1..n) does work, at O(n) with more code.

</details>

#### Example 2.3: Maximum product subarray [Medium]

**Problem:** [Maximum Product Subarray (152)](https://leetcode.com/problems/maximum-product-subarray/) · **Template:** Template 2 (Kadane)

**Problem.** Return the largest product of a non-empty contiguous subarray (values include negatives and zeros).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Products with negatives aren't monotonic, and logs fail on zeros and negatives. Kadane with two states: the largest and smallest product ending here, because a negative turns the smallest into the largest.

**Approach:** On a negative, swap `hi` and `lo`, then extend or restart each.

```python
def max_product(nums: list[int]) -> int:
    best = hi = lo = nums[0]
    for x in nums[1:]:
        if x < 0:
            hi, lo = lo, hi             # a negative swaps the roles of max and min
        hi = max(x, hi * x)
        lo = min(x, lo * x)
        best = max(best, hi)
    return best
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** Zeros reset both states (restart wins). `[-2, 3, -4]` → 24 needs the `lo` state. A single negative element returns itself.

</details>

#### Example 2.4: Max subarray sum with at most one deletion [Medium]

**Problem:** [Maximum Subarray Sum with One Deletion (1186)](https://leetcode.com/problems/maximum-subarray-sum-with-one-deletion/) · **Template:** Template 2 (Kadane)

**Problem.** Return the max sum of a non-empty subarray after optionally deleting one element (the result must stay non-empty).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Kadane plus one bit of state: "have I deleted yet?"

**Approach:** `keep` = best ending here with no deletion. `deleted` = best ending here with one deletion: either extend an earlier deleted run with `x`, or delete `x` itself (giving `keep` from the previous step, which is non-empty).

```python
def maximum_sum_one_deletion(arr: list[int]) -> int:
    keep = arr[0]                       # best ending here, no deletion
    deleted = float("-inf")             # best ending here, exactly one deletion used
    best = arr[0]
    for x in arr[1:]:
        deleted = max(deleted + x, keep)    # extend a deleted run, or delete x now
        keep = max(keep + x, x)
        best = max(best, keep, deleted)
    return best
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** Update `deleted` **before** `keep` (it needs the previous `keep`). All negative → the max single element; `deleted` starts at −inf so an empty result never wins.

</details>

### Template 3: prefix sums + monotonic deque

> - **Example problem:** [Shortest Subarray with Sum at Least K (862)](https://leetcode.com/problems/shortest-subarray-with-sum-at-least-k/)
> - **Why this template:** a "≥ k" target with negatives needs the *latest* start with a small enough prefix. An increasing deque of start indices keeps only useful starts and drops each one once.
> - **Worked below:** Examples 3.1, 3.2
> - **Also fits:** [Jump Game VI (1696)](https://leetcode.com/problems/jump-game-vi/) (same deque over DP values; see module 06a)

```python
def shortest_sum_at_least(nums: list[int], k: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    dq: deque[int] = deque()            # start indices, prefix values increasing
    best = len(nums) + 1
    for j, pj in enumerate(prefix):
        while dq and pj - prefix[dq[0]] >= k:       # TODO: "valid with this start" condition
            best = min(best, j - dq.popleft())
        while dq and prefix[dq[-1]] >= pj:          # dominated starts
            dq.pop()
        dq.append(j)
    return best if best <= len(nums) else -1
```

Why each pop is safe:

- **Back pop** (`P[back] ≥ P[j]`): `j` is a later start with a smaller prefix, so it gives a larger sum *and* a shorter subarray for every future end. `back` can never win.
- **Front pop** (`P[j] − P[front] ≥ k`): record the length. Any later end using the same start is longer, so drop it.

This restores the monotonicity the window lost: the window's "sum grows with length" becomes "the deque's prefixes increase".

#### Example 3.1: Shortest subarray with sum at least k, negatives allowed [Hard]

**Problem:** [Shortest Subarray with Sum at Least K (862)](https://leetcode.com/problems/shortest-subarray-with-sum-at-least-k/) · **Template:** Template 3 (prefix sums + monotonic deque)

**Problem.** Return the length of the shortest non-empty subarray with sum ≥ `k`, or −1. Values can be negative.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Shortest + "at least" + negatives → Template 3. (With positives only, it's [Minimum Size Subarray Sum (209)](https://leetcode.com/problems/minimum-size-subarray-sum/), a plain window.)

**Approach:** Template 3 as written.

```python
def shortest_subarray(nums: list[int], k: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    dq: deque[int] = deque()
    best = len(nums) + 1
    for j, pj in enumerate(prefix):
        while dq and pj - prefix[dq[0]] >= k:
            best = min(best, j - dq.popleft())
        while dq and prefix[dq[-1]] >= pj:
            dq.pop()
        dq.append(j)
    return best if best <= len(nums) else -1
```

**Complexity:** O(n) time (each index is pushed and popped once), O(n) space.

**Edge cases:** Front pops come before back pops and the append, so `j` never pairs with itself. `[2, -1, 2], k = 3` → 3: a window would shrink past the 2 and miss it. No answer → −1.

</details>

#### Example 3.2: Max sum with length between lo and hi [Medium]

**Problem:** follow-up to [Maximum Subarray (53)](https://leetcode.com/problems/maximum-subarray/), not a LeetCode problem · **Template:** Template 3 (prefix sums + monotonic deque)

**Problem.** Given integers (any signs) and `1 ≤ lo ≤ hi`, return the max sum of a subarray whose length is in `[lo, hi]`, or `None` if `n < lo`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** For end `j`, the valid starts are `i ∈ [j − hi, j − lo]`, a range that slides right. Max `P[j] − P[i]` = `P[j] − min(P[i])` over a sliding range → sliding-window **min** deque over prefixes.

**Approach:** At each `j`, push start `j − lo` (it just became allowed), expire starts below `j − hi`, read the front.

```python
def max_sum_length_between(nums: list[int], lo: int, hi: int) -> int | None:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    dq: deque[int] = deque()            # start indices in [j - hi, j - lo], prefix increasing
    best = None
    for j in range(lo, len(prefix)):
        i = j - lo                      # newest start allowed for this end
        while dq and prefix[dq[-1]] >= prefix[i]:
            dq.pop()
        dq.append(i)
        if dq[0] < j - hi:              # oldest start is now too far back
            dq.popleft()
        cand = prefix[j] - prefix[dq[0]]
        best = cand if best is None else max(best, cand)
    return best
```

**Complexity:** O(n) time, O(hi − lo) space.

**Edge cases:** One expiry per step is enough, because exactly one start (`j − hi − 1`) leaves each step. `lo = 1, hi = n` reduces to Kadane. Only `lo` ("length ≥ k") → keep a running min instead of a deque, as in Template 7's check.

</details>

### Template 4: prefix sums + monotonic stack (longest)

> - **Example problem:** [Longest Well-Performing Interval (1124)](https://leetcode.com/problems/longest-well-performing-interval/)
> - **Why this template:** for "longest with sum ≥ k" you want the *earliest* start with a small enough prefix. Only running prefix minima can be that start; scanning ends right to left lets each start be popped once.
> - **Worked below:** Examples 4.1, 4.2
> - **Also fits:** [Maximum Width Ramp (962)](https://leetcode.com/problems/maximum-width-ramp/)

```python
def longest_sum_at_least(nums: list[int], k: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    stack: list[int] = []               # start candidates, prefix values strictly decreasing
    for i, p in enumerate(prefix):
        if not stack or p < prefix[stack[-1]]:
            stack.append(i)
    best = 0
    for j in range(len(prefix) - 1, -1, -1):
        while stack and prefix[j] - prefix[stack[-1]] >= k:   # TODO: validity condition
            best = max(best, j - stack.pop())
    return best
```

Why it works:

- **Pass 1:** a start `i` with some earlier `i' < i` and `P[i'] ≤ P[i]` is never best: `i'` is longer and at least as valid. So only strict prefix minima (left to right) are candidates, and they're decreasing.
- **Pass 2:** for the largest `j` first, pop every start it validates. A popped start's best partner is this `j` (smaller `j` are shorter), so it's done. A start that fails stays for smaller `j`.

Compare with Template 3: shortest wants the **latest** valid start (deque, scan forward); longest wants the **earliest** (stack of minima, scan backward).

#### Example 4.1: Widest ramp [Medium]

**Problem:** [Maximum Width Ramp (962)](https://leetcode.com/problems/maximum-width-ramp/) · **Template:** Template 4 (prefix sums + monotonic stack (longest))

**Problem.** A ramp is a pair `i < j` with `nums[i] ≤ nums[j]`. Return the max `j − i`, or 0.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "Farthest pair with left ≤ right" is Template 4's pair search on the raw values instead of prefixes. It's the bare technique.

**Approach:** Stack of strictly decreasing values from the left; scan `j` from the right, popping every start with `nums[start] ≤ nums[j]`.

```python
def max_width_ramp(nums: list[int]) -> int:
    stack: list[int] = []
    for i, x in enumerate(nums):
        if not stack or x < nums[stack[-1]]:
            stack.append(i)
    best = 0
    for j in range(len(nums) - 1, -1, -1):
        while stack and nums[stack[-1]] <= nums[j]:
            best = max(best, j - stack.pop())
    return best
```

**Complexity:** O(n) time, O(n) space. (Sorting indices by value is an O(n log n) alternative.)

**Edge cases:** Strictly decreasing → 0. `j == i` pairs give width 0, harmless.

</details>

#### Example 4.2: Longest stretch with more good hours than bad [Medium]

**Problem:** [Longest Well-Performing Interval (1124)](https://leetcode.com/problems/longest-well-performing-interval/) · **Template:** Template 4 (prefix sums + monotonic stack (longest))

**Problem.** An hour is tiring if `hours[i] > 8`. Return the length of the longest interval with strictly more tiring days than non-tiring days.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Map to +1 / −1 → longest subarray with sum ≥ 1, negatives present → Template 4 with `k = 1`.

**Approach:** Template 4 on the transformed prefixes.

```python
def longest_wpi(hours: list[int]) -> int:
    prefix = [0]
    for h in hours:
        prefix.append(prefix[-1] + (1 if h > 8 else -1))
    stack: list[int] = []
    for i, p in enumerate(prefix):
        if not stack or p < prefix[stack[-1]]:
            stack.append(i)
    best = 0
    for j in range(len(prefix) - 1, -1, -1):
        while stack and prefix[j] - prefix[stack[-1]] >= 1:
            best = max(best, j - stack.pop())
    return best
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** No tiring hours → 0. **±1 shortcut:** because the prefix moves by 1, the earliest prefix below `score` is the first occurrence of `score − 1`, giving a one-pass Template 1 solution:

```python
def longest_wpi_hash(hours: list[int]) -> int:
    first: dict[int, int] = {}
    score = best = 0
    for i, h in enumerate(hours):
        score += 1 if h > 8 else -1
        if score > 0:
            best = i + 1                # whole prefix is well-performing
        elif score - 1 in first:
            best = max(best, i - first[score - 1])
        first.setdefault(score, i)
    return best
```

The shortcut needs ±1 steps; Template 4 works for any values.

</details>

### Template 5: prefix sums + sorted prefixes

> - **Example problem:** [Max Sum of Rectangle No Larger Than K (363)](https://leetcode.com/problems/max-sum-of-rectangle-no-larger-than-k/) (1D core)
> - **Why this template:** "max sum ≤ k" needs the smallest earlier prefix that's still ≥ `P[j] − k`. That's a successor query, which needs sorted order, not a hash.
> - **Worked below:** Example 5.1; Combined example C.1
> - **Also fits:** "closest subarray sum to k" (query both neighbors)

```python
def max_sum_at_most_k(nums: list[int], k: int) -> int | None:
    seen = [0]                          # sorted earlier prefixes (empty prefix included)
    prefix = 0
    best = None
    for x in nums:
        prefix += x
        i = bisect.bisect_left(seen, prefix - k)    # TODO: which neighbor (>= or <=) you need
        if i < len(seen):
            cand = prefix - seen[i]
            best = cand if best is None else max(best, cand)
        bisect.insort(seen, prefix)
    return best
```

Knobs: `bisect_left(seen, prefix - k)` finds the smallest `P[i] ≥ prefix − k`, the start giving the largest sum still ≤ k. For "closest to k", check `seen[i]` and `seen[i − 1]`. `insort` is O(n) per insert but a fast memmove (fine to ~10^5); say "a balanced BST / `SortedList` makes it O(log n)" in the interview.

#### Example 5.1: Max subarray sum no larger than k [Medium]

**Problem:** follow-up to [Max Sum of Rectangle No Larger Than K (363)](https://leetcode.com/problems/max-sum-of-rectangle-no-larger-than-k/), not a LeetCode problem · **Template:** Template 5 (prefix sums + sorted prefixes)

**Problem.** Given integers (any signs) and `k`, return the max subarray sum that's ≤ `k` (or `None` if none).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** A cap on the sum + negatives. Kadane can overshoot the cap, and a window can't shrink reliably → Template 5.

**Approach:** Template 5 as written.

```python
def max_subarray_sum_at_most_k(nums: list[int], k: int) -> int | None:
    seen = [0]
    prefix = 0
    best = None
    for x in nums:
        prefix += x
        i = bisect.bisect_left(seen, prefix - k)
        if i < len(seen):
            cand = prefix - seen[i]
            best = cand if best is None else max(best, cand)
        bisect.insort(seen, prefix)
    return best
```

**Complexity:** O(n log n) comparisons; O(n^2) worst-case element moves from `insort` (O(n log n) with a balanced BST). O(n) space.

**Edge cases:** Every element > k → `None`. Negative `k` works (`[-5, -2], k = −3` → −5). Early exit when `best == k`.

</details>

### Template 6: prefix sums + Fenwick tree (counting)

> - **Example problem:** [Count of Range Sum (327)](https://leetcode.com/problems/count-of-range-sum/)
> - **Why this template:** counting subarrays with sum in a range means counting earlier prefixes in a **value range**. A Fenwick tree over compressed prefix values counts that in O(log n).
> - **Worked below:** Examples 6.1, 6.2
> - **Also fits:** count subarrays with sum ≥ k (set `hi` to infinity), [Reverse Pairs (493)](https://leetcode.com/problems/reverse-pairs/) (same "count earlier values in a range" core)

```python
class Fenwick:

    def __init__(self, n: int):
        self.tree = [0] * (n + 1)

    def add(self, i: int, delta: int = 1) -> None:      # i is 0-indexed
        i += 1
        while i < len(self.tree):
            self.tree[i] += delta
            i += i & -i

    def prefix_count(self, i: int) -> int:              # count at positions [0, i)
        total = 0
        while i > 0:
            total += self.tree[i]
            i -= i & -i
        return total


def count_sum_in_range(nums: list[int], lo: int, hi: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    values = sorted(set(prefix))        # coordinate compression
    tree = Fenwick(len(values))
    total = 0
    for p in prefix:
        left = bisect.bisect_left(values, p - hi)       # first value >= p - hi
        right = bisect.bisect_right(values, p - lo)     # past last value <= p - lo
        total += tree.prefix_count(right) - tree.prefix_count(left)
        tree.add(bisect.bisect_left(values, p))         # TODO: insert AFTER querying
    return total
```

Knobs: the query range `[p − hi, p − lo]`, found by bisecting the sorted distinct prefixes (the bounds themselves don't need to be in `values`). If the prefixes are small integers (±1 steps), skip compression and offset by `n` (Example 6.2). **Merge-sort alternative:** count cross pairs between sorted halves with two forward pointers, also O(n log n) (shown in Example 6.1).

#### Example 6.1: Count subarrays with sum in [lower, upper] [Hard]

**Problem:** [Count of Range Sum (327)](https://leetcode.com/problems/count-of-range-sum/) · **Template:** Template 6 (prefix sums + Fenwick tree (counting))

**Problem.** Count the subarrays whose sum lies in `[lower, upper]`. Values can be negative; n up to 10^5.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Count + a sum range + negatives. Template 1 handles one exact value, but a range would need one lookup per value in it → Template 6.

**Approach:** Template 6 as written.

```python
def count_range_sum(nums: list[int], lower: int, upper: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    values = sorted(set(prefix))
    tree = Fenwick(len(values))
    total = 0
    for p in prefix:
        left = bisect.bisect_left(values, p - upper)
        right = bisect.bisect_right(values, p - lower)
        total += tree.prefix_count(right) - tree.prefix_count(left)
        tree.add(bisect.bisect_left(values, p))
    return total
```

Merge-sort version: within a merge, each left-half prefix needs right-half prefixes in `[left + lower, left + upper]`. Both halves are sorted, so the two bounds only move forward.

```python
def count_range_sum_merge(nums: list[int], lower: int, upper: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)

    def sort_count(lo: int, hi: int) -> int:        # sorts prefix[lo:hi], returns pair count
        if hi - lo <= 1:
            return 0
        mid = (lo + hi) // 2
        count = sort_count(lo, mid) + sort_count(mid, hi)
        a = b = mid
        for left in prefix[lo:mid]:     # both halves sorted: windows only move right
            while a < hi and prefix[a] - left < lower:
                a += 1
            while b < hi and prefix[b] - left <= upper:
                b += 1
            count += b - a
        prefix[lo:hi] = sorted(prefix[lo:hi])
        return count

    return sort_count(0, len(prefix))
```

**Complexity:** O(n log n) time, O(n) space for both. (`sorted` on two already-sorted runs is a linear merge in Timsort, so the merge version stays O(n log n).)

**Edge cases:** Query before insert so `i < j`. Python ints don't overflow; in Java, prefixes need `long`.

</details>

#### Example 6.2: Count subarrays with more ones than zeros [Medium]

**Problem:** [Count Subarrays With More Ones Than Zeros (2031)](https://leetcode.com/problems/count-subarrays-with-more-ones-than-zeros/) · **Template:** Template 6 (prefix sums + Fenwick tree (counting))

**Problem.** Given a binary array, count subarrays with more 1s than 0s, modulo 10^9 + 7.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** ±1 transform → count subarrays with sum ≥ 1 → count earlier prefixes strictly below `P[j]`. The prefixes lie in `[−n, n]`, so offset by `n` instead of compressing.

**Approach:** Template 6 with `hi = ∞` and direct indexing.

```python
def subarrays_with_more_ones(nums: list[int]) -> int:
    mod = 10**9 + 7
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + (1 if x == 1 else -1))
    n = len(nums)
    tree = Fenwick(2 * n + 1)           # prefixes lie in [-n, n]: shift by n, no compression
    total = 0
    for p in prefix:
        total += tree.prefix_count(p + n)   # earlier prefixes strictly below p
        tree.add(p + n)
    return total % mod
```

**Complexity:** O(n log n) time, O(n) space. An O(n) version exists: since `P` moves by ±1, keep a count array of prefixes and a running "count below current" that changes by one bucket per step.

**Edge cases:** `[0]` → 0, `[1]` → 1. "At least as many" → use `prefix_count(p + n + 1)` (≤ instead of <).

</details>

### Template 7: binary search on the answer

> - **Example problem:** [Maximum Average Subarray II (644)](https://leetcode.com/problems/maximum-average-subarray-ii/)
> - **Why this template:** averages and medians aren't sums, so no prefix trick applies directly. Guessing the answer `x` and subtracting it turns "average ≥ x" into "sum ≥ 0", which prefixes can check in O(n).
> - **Worked below:** Example 7.1
> - **Also fits:** "max median of a subarray with length ≥ k" (guess x, map values to +1 if ≥ x else −1, check for a sum > 0)

```python
def max_average_at_least_len(nums: list[int], k: int, eps: float = 1e-6) -> float:
    def feasible(x: float) -> bool:                     # TODO: the transformed check
        prefix = 0.0
        prev = 0.0                      # B[j - k + 1]: latest start allowed for end j
        min_prev = float("inf")
        for j, a in enumerate(nums):
            prefix += a - x
            if j >= k - 1:
                if j >= k:
                    prev += nums[j - k] - x
                min_prev = min(min_prev, prev)
                if prefix - min_prev >= 0:
                    return True
        return False

    lo, hi = min(nums), max(nums)
    while hi - lo > eps:
        mid = (lo + hi) / 2
        if feasible(mid):
            lo = mid
        else:
            hi = mid
    return lo
```

Knobs: the **transform** (`a − x` for averages, ±1 for medians), the **check** (running min of prefixes at least `k` back, the length-≥-k case of Example 3.2), and the **search space** (real numbers with `eps`, or integers with module 07's template). Monotonic: if some subarray averages ≥ x, it also averages ≥ any smaller x.

#### Example 7.1: Max average with length at least k [Hard]

**Problem:** [Maximum Average Subarray II (644)](https://leetcode.com/problems/maximum-average-subarray-ii/) · **Template:** Template 7 (binary search on the answer)

**Problem.** Return the max average of a contiguous subarray with length ≥ `k` (error < 10^-5 accepted).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Length "at least" k and negatives → no fixed window (that's [Maximum Average Subarray I (643)](https://leetcode.com/problems/maximum-average-subarray-i/)). Average → binary search on the answer, with a prefix-min check.

**Approach:** Template 7 as written. The answer is between `min(nums)` and `max(nums)`.

```python
def find_max_average(nums: list[int], k: int) -> float:
    def feasible(x: float) -> bool:
        prefix = prev = 0.0
        min_prev = float("inf")
        for j, a in enumerate(nums):
            prefix += a - x
            if j >= k - 1:
                if j >= k:
                    prev += nums[j - k] - x
                min_prev = min(min_prev, prev)
                if prefix - min_prev >= 0:
                    return True
        return False

    lo, hi = min(nums), max(nums)
    while hi - lo > 1e-6:
        mid = (lo + hi) / 2
        if feasible(mid):
            lo = mid
        else:
            hi = mid
    return lo
```

**Complexity:** O(n log((max − min)/ε)) time, O(1) space.

**Edge cases:** `k = n` → the whole-array average. A fixed number of iterations (e.g. 60) avoids float-loop worries. An exact O(n) solution exists (convex hull of prefix points) but is beyond interview scope.

</details>

### Combined examples

#### Example C.1: Max rectangle sum no larger than k [Hard]

**Problem:** [Max Sum of Rectangle No Larger Than K (363)](https://leetcode.com/problems/max-sum-of-rectangle-no-larger-than-k/) · **Template:** Template 5 (prefix sums + sorted prefixes)

**Problem.** Given a matrix (any signs) and `k`, return the max sum of a rectangle that's ≤ `k`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** 2D rectangle → fix a pair of rows (top, bottom), collapse each column to its sum between them, and the rectangle becomes a 1D subarray. Cap + negatives → Template 5 on each collapsed row.

**Approach:** Iterate row pairs along the shorter dimension (transpose if needed), accumulate column sums, call Example 5.1.

```python
def max_sum_submatrix(matrix: list[list[int]], k: int) -> int:
    rows, cols = len(matrix), len(matrix[0])
    if rows > cols:                     # loop over pairs along the SHORTER side
        matrix = [list(col) for col in zip(*matrix)]
        rows, cols = cols, rows
    best = None
    for top in range(rows):
        col_sums = [0] * cols
        for bottom in range(top, rows):
            for c in range(cols):
                col_sums[c] += matrix[bottom][c]
            cand = max_subarray_sum_at_most_k(col_sums, k)
            if cand is not None and (best is None or cand > best):
                best = cand
                if best == k:
                    return k
    return best
```

**Complexity:** O(min(R,C)^2 · max(R,C) log max(R,C)) time, O(max(R,C)) space.

**Edge cases:** The problem guarantees an answer exists; the code returns `None` otherwise. Return `k` early on an exact hit.

</details>

#### Example C.2: Count submatrices summing to target [Hard]

**Problem:** [Number of Submatrices That Sum to Target (1074)](https://leetcode.com/problems/number-of-submatrices-that-sum-to-target/) · **Template:** Template 1 (prefix sum + hash map)

**Problem.** Count the non-empty submatrices whose sum equals `target`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Same row-pair collapse as C.1; exact target → Template 1 on each collapsed row.

**Approach:** For each (top, bottom), accumulate column sums and count subarrays summing to `target` with Example 1.1.

```python
def num_submatrix_sum_target(matrix: list[list[int]], target: int) -> int:
    rows, cols = len(matrix), len(matrix[0])
    total = 0
    for top in range(rows):
        col_sums = [0] * cols
        for bottom in range(top, rows):
            for c in range(cols):
                col_sums[c] += matrix[bottom][c]
            total += subarray_sum(col_sums, target)
    return total
```

**Complexity:** O(R^2 · C) time, O(C) space. Transpose if R > C.

**Edge cases:** `target = 0` with zeros counts many overlapping submatrices; the count map handles it like 1D.

</details>

### Beyond the templates

#### Example B.1: Longest subarray with positive product [Medium]

**Problem:** [Maximum Length of Subarray With Positive Product (1567)](https://leetcode.com/problems/maximum-length-of-subarray-with-positive-product/) · **Template:** none (prefix parity + reset at zeros)

**Problem.** Return the length of the longest subarray whose product is positive.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Only the sign matters, and the sign is the **parity of negatives**. A zero kills any subarray through it, so split at zeros. Within a segment, the product of `(i, j]` is positive when both prefixes have the same parity → longest → first index per parity. It's a two-key Template 1 with a reset; it's worth knowing because it's the standard rescue for "window over products with zeros and negatives".

**Approach:** Track the parity of negatives since the last zero and the first index of each parity; reset at every zero.

```python
def get_max_len(nums: list[int]) -> int:
    best = 0
    first = {0: -1}                     # parity of negatives -> earliest index (since last zero)
    parity = 0
    for i, x in enumerate(nums):
        if x == 0:                      # zero splits the array: restart
            first = {0: i}
            parity = 0
            continue
        if x < 0:
            parity ^= 1
        if parity in first:
            best = max(best, i - first[parity])
        else:
            first[parity] = i
    return best
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** All zeros → 0. The reset seeds `{0: i}` so a segment's own start counts as the empty prefix.

</details>

## Template cheat sheet

Every template from this module in one place, for warm-ups and review.

### Template 1: prefix sum + hash map

```python
def count_sum_k(nums: list[int], k: int) -> int:
    seen: dict[int, int] = defaultdict(int)
    seen[0] = 1                         # TODO: seed the empty prefix ({0: 1} count, {0: -1} index)
    prefix = 0
    total = 0
    for x in nums:
        prefix += x                     # TODO: transform x first if needed (0 -> -1, x % m, ...)
        total += seen[prefix - k]       # TODO: look up the complement BEFORE inserting
        seen[prefix] += 1               # TODO: store count / first index / last index
    return total
```

### Template 2: Kadane

```python
def kadane(nums: list[int]) -> int:
    best = cur = nums[0]
    for x in nums[1:]:
        cur = max(x, cur + x)           # TODO: extend or restart (add more states for variants)
        best = max(best, cur)           # TODO: record
    return best
```

### Template 3: prefix sums + monotonic deque

```python
def shortest_sum_at_least(nums: list[int], k: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    dq: deque[int] = deque()            # start indices, prefix values increasing
    best = len(nums) + 1
    for j, pj in enumerate(prefix):
        while dq and pj - prefix[dq[0]] >= k:       # TODO: "valid with this start" condition
            best = min(best, j - dq.popleft())
        while dq and prefix[dq[-1]] >= pj:          # dominated starts
            dq.pop()
        dq.append(j)
    return best if best <= len(nums) else -1
```

### Template 4: prefix sums + monotonic stack (longest)

```python
def longest_sum_at_least(nums: list[int], k: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    stack: list[int] = []               # start candidates, prefix values strictly decreasing
    for i, p in enumerate(prefix):
        if not stack or p < prefix[stack[-1]]:
            stack.append(i)
    best = 0
    for j in range(len(prefix) - 1, -1, -1):
        while stack and prefix[j] - prefix[stack[-1]] >= k:   # TODO: validity condition
            best = max(best, j - stack.pop())
    return best
```

### Template 5: prefix sums + sorted prefixes

```python
def max_sum_at_most_k(nums: list[int], k: int) -> int | None:
    seen = [0]                          # sorted earlier prefixes (empty prefix included)
    prefix = 0
    best = None
    for x in nums:
        prefix += x
        i = bisect.bisect_left(seen, prefix - k)    # TODO: which neighbor (>= or <=) you need
        if i < len(seen):
            cand = prefix - seen[i]
            best = cand if best is None else max(best, cand)
        bisect.insort(seen, prefix)
    return best
```

### Template 6: prefix sums + Fenwick tree (counting)

```python
class Fenwick:

    def __init__(self, n: int):
        self.tree = [0] * (n + 1)

    def add(self, i: int, delta: int = 1) -> None:      # i is 0-indexed
        i += 1
        while i < len(self.tree):
            self.tree[i] += delta
            i += i & -i

    def prefix_count(self, i: int) -> int:              # count at positions [0, i)
        total = 0
        while i > 0:
            total += self.tree[i]
            i -= i & -i
        return total


def count_sum_in_range(nums: list[int], lo: int, hi: int) -> int:
    prefix = [0]
    for x in nums:
        prefix.append(prefix[-1] + x)
    values = sorted(set(prefix))        # coordinate compression
    tree = Fenwick(len(values))
    total = 0
    for p in prefix:
        left = bisect.bisect_left(values, p - hi)       # first value >= p - hi
        right = bisect.bisect_right(values, p - lo)     # past last value <= p - lo
        total += tree.prefix_count(right) - tree.prefix_count(left)
        tree.add(bisect.bisect_left(values, p))         # TODO: insert AFTER querying
    return total
```

### Template 7: binary search on the answer

```python
def max_average_at_least_len(nums: list[int], k: int, eps: float = 1e-6) -> float:
    def feasible(x: float) -> bool:                     # TODO: the transformed check
        prefix = 0.0
        prev = 0.0                      # B[j - k + 1]: latest start allowed for end j
        min_prev = float("inf")
        for j, a in enumerate(nums):
            prefix += a - x
            if j >= k - 1:
                if j >= k:
                    prev += nums[j - k] - x
                min_prev = min(min_prev, prev)
                if prefix - min_prev >= 0:
                    return True
        return False

    lo, hi = min(nums), max(nums)
    while hi - lo > eps:
        mid = (lo + hi) / 2
        if feasible(mid):
            lo = mid
        else:
            hi = mid
    return lo
```

## Variations and extensions

- **Exact target on non-negative values:** both work. The window version is `atMost(k) − atMost(k−1)` (module 06); the prefix map is simpler to get right and handles negatives if the constraints change.
- **Modular keys:** "divisible by k" → key `prefix % k`. Count version: [Subarray Sums Divisible by K (974)](https://leetcode.com/problems/subarray-sums-divisible-by-k/), seeded `{0: 1}`.
- **Balance transforms:** "equal X and Y", "more X than Y", "median = k" → map categories to +1/−1/0, then a sum problem.
- **XOR instead of sum:** prefix XOR + hash map counts subarrays with XOR = k (`P[i] = P[j] ^ k`), e.g. [Count Triplets That Can Form Two Arrays of Equal XOR (1442)](https://leetcode.com/problems/count-triplets-that-can-form-two-arrays-of-equal-xor/).
- **Parity / bitmask keys:** "every vowel appears an even number of times" → key = bitmask of parities, first-index map ([Find the Longest Substring Containing Vowels in Even Counts (1371)](https://leetcode.com/problems/find-the-longest-substring-containing-vowels-in-even-counts/)).
- **2D:** fix a row pair, collapse columns (Combined examples). Use 2D prefix sums for repeated rectangle queries (module 04).
- **Streaming:** Templates 1, 2, and 3 are one-pass and work online; Template 4 needs the whole array (backward pass); Template 6 with compression needs all values up front (offline), or a dynamic structure online.
- **k-length constraint plus exact target:** add the length bound to the map (store first index, check distance) or expire old prefixes from a map keyed by value.

## Common pitfalls

- Using a window because the story says "contiguous" while the constraints allow negatives.
- Forgetting the empty-prefix seed (`{0: 1}` or `{0: -1}`): misses every subarray starting at index 0.
- Inserting `P[j]` before querying: pairs a prefix with itself (an empty subarray), wrong when `k = 0` or the length must be ≥ 1.
- Storing the wrong index: first index for **longest**, last index for **shortest**, counts for **how many**.
- Kadane initialized with `best = 0`: wrong when every value is negative.
- Circular Kadane without the all-negative check: returns 0 (the empty subarray).
- Product Kadane without the min state, or forgetting to swap on a negative.
- Template 3: appending `j` before the front pops lets `j` pair with itself (an empty subarray, wrong when `k ≤ 0`).
- Template 4: running the end scan left to right (pops starts too early).
- Fenwick compression that bisects the wrong side (`bisect_left` for the lower bound, `bisect_right` for the upper).
- `%` on negative numbers in Java/C++ (Python is fine for positive moduli).
- Float binary search with `while lo < hi` on reals: loop on `hi - lo > eps` or a fixed iteration count.

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| These templates | Plain sliding window (06) | Every value ≥ 0 and the constraint is "at most / at least" |
| Template 1 (exact count) | `atMost(k) − atMost(k−1)` window (06) | Non-negative values or counts; both work, the window is O(1) space |
| Template 2 (Kadane) | DP over subsequences (16) | "Subsequence", not "subarray": elements can be skipped |
| Template 3 (shortest ≥ k) | Window (06, Template 3) | Positives only: [Minimum Size Subarray Sum (209)](https://leetcode.com/problems/minimum-size-subarray-sum/) |
| Template 5 (≤ k) | Kadane | No cap: just the max sum |
| Template 6 (range count) | Template 1 | Range is a single value (`lo == hi`): a hash map is O(n) |
| Template 7 | Fixed window (06, Template 1) | Length is exactly k, not "at least k": [Maximum Average Subarray I (643)](https://leetcode.com/problems/maximum-average-subarray-i/) |
| Prefix sums | Two pointers (05) | Sorted array, pair sums, not contiguous |

## Practice ladder

### Must-do

1. [Subarray Sum Equals K (560)](https://leetcode.com/problems/subarray-sum-equals-k/) · Medium · Prefix sum + count map, seeded `{0: 1}`.
2. [Contiguous Array (525)](https://leetcode.com/problems/contiguous-array/) · Medium · ±1 transform + first-index map.
3. [Maximum Sum Circular Subarray (918)](https://leetcode.com/problems/maximum-sum-circular-subarray/) · Medium · Kadane max and min; all-negative check.
4. [Longest Well-Performing Interval (1124)](https://leetcode.com/problems/longest-well-performing-interval/) · Medium · ±1 transform + Template 4 (or the ±1 hash shortcut).
5. [Shortest Subarray with Sum at Least K (862)](https://leetcode.com/problems/shortest-subarray-with-sum-at-least-k/) · Hard · Increasing deque over prefix sums.

### Optional

**Prefix sum + hash map**
- [Continuous Subarray Sum (523)](https://leetcode.com/problems/continuous-subarray-sum/) · Medium · Remainder key, first index, length ≥ 2.
- [Subarray Sums Divisible by K (974)](https://leetcode.com/problems/subarray-sums-divisible-by-k/) · Medium · Remainder key, counts.
- [Make Sum Divisible by P (1590)](https://leetcode.com/problems/make-sum-divisible-by-p/) · Medium · Remainder key, last index.
- [Maximum Size Subarray Sum Equals k (325)](https://leetcode.com/problems/maximum-size-subarray-sum-equals-k/) · Medium · First-index map (premium).
- [Find the Longest Substring Containing Vowels in Even Counts (1371)](https://leetcode.com/problems/find-the-longest-substring-containing-vowels-in-even-counts/) · Medium · Bitmask-of-parities key.
- [Count Subarrays With Median K (2488)](https://leetcode.com/problems/count-subarrays-with-median-k/) · Hard · ±1 transform anchored at k, two complements.

**Kadane**
- [Maximum Subarray (53)](https://leetcode.com/problems/maximum-subarray/) · Medium · The base case.
- [Maximum Product Subarray (152)](https://leetcode.com/problems/maximum-product-subarray/) · Medium · Max and min states, swap on negatives.
- [Maximum Subarray Sum with One Deletion (1186)](https://leetcode.com/problems/maximum-subarray-sum-with-one-deletion/) · Medium · Extra "deleted" state.
- [Maximum Length of Subarray With Positive Product (1567)](https://leetcode.com/problems/maximum-length-of-subarray-with-positive-product/) · Medium · Parity + reset at zeros.

**Monotonic deque / stack over prefixes**
- [Maximum Width Ramp (962)](https://leetcode.com/problems/maximum-width-ramp/) · Medium · Template 4's core on raw values.
- [Jump Game VI (1696)](https://leetcode.com/problems/jump-game-vi/) · Medium · Sliding max deque over DP values.

**Sorted prefixes and counting**
- [Max Sum of Rectangle No Larger Than K (363)](https://leetcode.com/problems/max-sum-of-rectangle-no-larger-than-k/) · Hard · Row-pair collapse + sorted prefixes.
- [Count of Range Sum (327)](https://leetcode.com/problems/count-of-range-sum/) · Hard · Fenwick or merge sort over prefixes.
- [Count Subarrays With More Ones Than Zeros (2031)](https://leetcode.com/problems/count-subarrays-with-more-ones-than-zeros/) · Medium · ±1 + Fenwick (premium).
- [Number of Submatrices That Sum to Target (1074)](https://leetcode.com/problems/number-of-submatrices-that-sum-to-target/) · Hard · Row-pair collapse + count map.

**Binary search on the answer**
- [Maximum Average Subarray II (644)](https://leetcode.com/problems/maximum-average-subarray-ii/) · Hard · Subtract the guess, prefix-min check (premium).

## Self-check

1. Give a 3-element array where "longest subarray with sum ≤ k" breaks the sliding window. Which left endpoint did the window wrongly discard?
2. Why doesn't adding a constant to every element (to make them positive) rescue the window?
3. In Template 1, when do you store counts, the first index, or the last index? Why the seed `{0: 1}` vs `{0: -1}`?
4. Why is Kadane the same as "prefix minus the smallest earlier prefix"?
5. In circular max subarray, why is the wrapping answer `total − min subarray`, and when is it invalid?
6. In Template 3, justify both the back pop and the front pop. Why does shortest use a deque and longest use a stack scanned backward?
7. Why does Template 5 need sorted order while Template 1 only needs a hash map?
8. How do you count subarrays with sum ≥ k when values can be negative? What's the complexity, and when can you drop the coordinate compression?
9. Why is "some subarray with length ≥ k averages ≥ x" monotonic in x, and how does subtracting x turn it into a prefix problem?
10. How do you turn a rectangle-sum problem into a 1D subarray problem, and what does it cost?
