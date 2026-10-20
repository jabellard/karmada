# 18. Greedy

Why this matters: greedy makes one locally best, irrevocable choice at a time. When it works it is the shortest, fastest solution (usually a sort plus one pass). But it looks right far more often than it is, so the skill is **justifying** the choice, or quickly breaking it and switching to DP.

Companion code: `code/18_greedy.py`.

## At a glance

- **Use when:** min/max count or cost, input sortable by one key, future depends on one scalar (last end, reach, capacity).
- **Templates:**
  - Template 1: sort by a justified key, then sweep: interval selection, matching sorted lists.
  - Template 2: reachability frontier: reach / min jumps; reachable set is a prefix.
  - Template 3: greedy with a heap: regret (evict worst) or unlock-and-pick-best.
  - Template 4: counterexample search harness: attack a rule with brute force.
  - Template 5: two-pass constraints: left sweep, then right sweep, combine per element.
- **Typical complexity:** O(n log n) (sort) or O(n).
- **Top trap:** trusting a rule because it passes the prompt's examples; prove it or break it on 3–4 items, then switch to DP.

## Core idea

A greedy algorithm is correct when:

1. **Greedy-choice property**: some optimal solution begins with the greedy choice.
2. **Optimal substructure**: what remains after that choice is a smaller instance of the same problem.

Prove property 1 with one of two arguments.

**Exchange argument.** Where an optimal solution OPT first differs from greedy, swap in greedy's choice and show the result is still feasible and no worse. Repeated swaps turn OPT into greedy's solution without losing value.

```
OPT:     [ o1 ][ o2 ][ o3 ] ...        greedy picks g1 with end(g1) <= end(o1)
swap:    [ g1 ][ o2 ][ o3 ] ...        still feasible (g1 ends no later), same count
```

**"Greedy stays ahead" argument.** Define a progress measure after k steps (end of the k-th chosen interval, farthest reach after k jumps) and prove by induction that greedy's is always at least as good as any other algorithm's. Never behind means never worse.

```
step k:     1    2    3    4
greedy:    e1 <= e2 <= e3 <= e4      (greedy's k-th end time)
any alg:   f1    f2    f3    f4      invariant: e_k <= f_k for all k
```

**Invariant to state out loud.** Interval scheduling: "after processing intervals sorted by end, `end` is the earliest finish of any max-size compatible set of the intervals seen so far". Jump Game: "`far` is the farthest index reachable using positions 0..i". If you can't phrase it, you probably don't have a correct greedy yet.

**Test before you trust.** Spend 60 seconds trying to break a proposed rule on 2–4 items (at home, automate with Template 4). Wrong rules usually fail at size 3 or 4.

## Recognition signals

- **Max count / min count / min cost** over input **sortable** by a key (end, size, ratio, deadline).
- n up to 10^5 or 10^6 rules out O(n^2) DP; expect O(n log n) or O(n).
- Choices are **independent once ordered**: each affects the future only through one scalar (end, capacity, reach).
- "Minimum number of X to cover / hit / pierce Y": arrows, platforms, rooms, taps, jumps.
- "Can you reach / complete": jump game, gas station, refueling.
- Matching two sorted lists: people to boats, cookies to children, workers to jobs.
- 1-D neighbor constraints ("higher rating than a neighbor gets more"): two passes.
- "Arrange to form the largest/smallest": custom sort comparator.
- "Choose k items; each choice unlocks others": heap of available options.
- Red flag: "exactly", "number of ways", or arbitrary weights/values (0/1 knapsack, arbitrary coins) mean DP.

## Disguises

1. **Conference talks.** A single hall gets requests for talks, each with a fixed start and end time. Accept as many as possible.
   **Decode:** single resource + fixed start/end + "as many as possible" = interval scheduling. Sort by end, take whenever start >= last end.

2. **Firmware rollout waves.** Each server can be patched during a window [earliest, latest]. A "wave" patches every server whose window contains its timestamp. What is the fewest waves?
   **Decode:** points that stab all intervals = Minimum Arrows. Sort by end; place a point at the end of the first unstabbed interval.

3. **Hop-limited relays.** Relay i forwards a packet up to `r[i]` positions ahead. What is the fewest forwards from relay 0 to the last relay?
   **Decode:** Jump Game II. BFS layers collapse into a frontier sweep because the reachable indices form a prefix.

4. **Electric delivery loop.** Charging stations sit on a circular route; station i gives `g[i]` kWh and the leg to the next costs `c[i]` kWh. Find a station from which the van can complete the loop.
   **Decode:** Gas Station. If total supply >= total demand, start just after the running balance's minimum.

5. **Pod co-scheduling.** Each node hosts at most two pods within a memory limit L; each pod needs <= L. Place all pods on the fewest nodes.
   **Decode:** Boats to Save People. Sort, pair heaviest with lightest when they fit (two pointers, `05-two-pointers.md`).

6. **CPU job cooldown.** A single core runs jobs of several types; same-type jobs need at least n idle-or-other ticks between them. Minimum total ticks?
   **Decode:** Task Scheduler. The most frequent type dictates frames of width n+1.

## Templates and worked examples

### Template 1: sort by a justified key, then sweep

> - **Example problem:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/)
> - **Why this template:** "Minimum intervals to remove" = keep a maximum compatible set; the exchange argument justifies sorting by end, and the only state is the last kept end.
> - **Worked below:** Example 1.2
> - **Also fits:** [Assign Cookies (455)](https://leetcode.com/problems/assign-cookies/)

```python
def greedy_sort_and_sweep(items: list, key: Callable, can_take: Callable, take: Callable, state):
    """Generic sort-then-sweep greedy.

    key:      the ordering that the exchange argument justifies (TODO per problem)
    can_take: (state, item) -> bool, is the item compatible with choices so far
    take:     (state, item) -> new state
    """
    chosen = []
    for item in sorted(items, key=key):
        if can_take(state, item):
            state = take(state, item)
            chosen.append(item)
    return chosen, state
```

Knobs:

| Problem | key | state | can_take |
|---|---|---|---|
| Max non-overlapping intervals | end | last end | start >= end |
| Min arrows | end | last arrow x | start > x (else covered) |
| Assign cookies | size (both lists) | pointer into other list | cookie >= greed |
| Course Schedule III | deadline | total time + heap | add then evict if over |

In practice you inline this; the table is the point. The hard part is choosing `key`, which the exchange argument decides.

#### Example 1.1: Matching by sorting [Easy/Medium]

**Problems:** [Assign Cookies (455)](https://leetcode.com/problems/assign-cookies/), [Boats to Save People (881)](https://leetcode.com/problems/boats-to-save-people/) · **Template:** Template 1 (sort by a justified key, then sweep)

Problem A: children have greed factors, cookies have sizes; a child is content with a cookie at least their greed, one cookie each. Maximize content children.
Problem B: boats carry at most two people and at most `limit` weight. Minimize boats.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** two collections paired under a size constraint, maximizing matches or minimizing groups. Sort both and sweep.

**Approach:**
- A: give each cookie, smallest first, to the least greedy unmet child it satisfies. Exchange: if OPT gives the smallest usable cookie to a greedier child, swapping with the least greedy child keeps both satisfied or frees a bigger cookie.
- B: pair the heaviest person with the lightest. If the lightest doesn't fit, nobody does, so the heaviest goes alone. If they fit, any OPT pairing the heaviest with someone else can swap partners with the lightest without breaking either boat.

```python
def find_content_children(greed: list[int], cookies: list[int]) -> int:
    greed, cookies = sorted(greed), sorted(cookies)
    child = 0
    for c in cookies:                       # smallest cookie first
        if child < len(greed) and c >= greed[child]:
            child += 1                      # least greedy unmet child gets it
    return child


def num_rescue_boats(people: list[int], limit: int) -> int:
    people = sorted(people)
    lo, hi, boats = 0, len(people) - 1, 0
    while lo <= hi:
        if people[lo] + people[hi] <= limit:
            lo += 1                         # lightest rides with heaviest
        hi -= 1                             # heaviest always leaves
        boats += 1
    return boats
```

**Complexity:** O(n log n) for sorting, O(1) extra beyond the sorted copies.

**Edge cases:** `lo <= hi` (not `<`) so a lone middle person gets a boat. More than two per boat is bin packing (NP-hard), where greedy fails.

</details>

#### Example 1.2: Non-overlapping intervals and arrows [Medium]

**Problems:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/), [Minimum Number of Arrows to Burst Balloons (452)](https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/) · **Template:** Template 1 (sort by a justified key, then sweep)

Problem: delete the fewest intervals so the rest are pairwise non-overlapping (touching endpoints is fine). Companion: find the fewest points such that every closed interval contains one.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "minimum to delete so none overlap" = n minus the maximum compatible set (classic interval scheduling).

**Approach:** sort by **end** and keep the earliest finisher; it leaves the most room.
Exchange: let OPT's first interval be o1 and greedy's g1, with end(g1) <= end(o1). Replacing o1 with g1 keeps OPT feasible (everything after o1 starts at or after end(o1) >= end(g1)) and keeps the count. Induct on the rest.
Why not by start? A long early interval blocks many short ones: `[0,10], [1,2], [3,4]`. Why not shortest first? In `[0,4], [4,8], [3,5]`, the short middle one blocks both others.

```python
def max_non_overlapping(intervals: list[list[int]]) -> int:
    """Max number of pairwise non-overlapping intervals; touching endpoints OK."""
    count, end = 0, float("-inf")
    for s, e in sorted(intervals, key=lambda iv: iv[1]):   # earliest end first
        if s >= end:
            count += 1
            end = e
    return count


def erase_overlap_intervals(intervals: list[list[int]]) -> int:
    """Min removals so the rest do not overlap = n - max kept."""
    return len(intervals) - max_non_overlapping(intervals)


def min_arrows(points: list[list[int]]) -> int:
    """Min arrows to burst balloons [s, e] (closed; touching shares an arrow)."""
    if not points:
        return 0
    arrows, end = 0, float("-inf")
    for s, e in sorted(points, key=lambda p: p[1]):
        if s > end:          # strictly after: needs a new arrow at e
            arrows += 1
            end = e
    return arrows
```

Arrows (stays ahead): an arrow at the earliest end is as far right as possible while still hitting that balloon, so it hits a superset of the remaining balloons any other valid position would.

**Complexity:** O(n log n) time for the sort, O(1) extra (O(n) if sort copies).

**Edge cases:** `>=` vs `>` encodes whether touching counts as overlap; read the statement. Empty input. `float("-inf")` is safe for huge coordinates in Python.

</details>

#### Example 1.3: Largest Number [Medium]

**Problem:** [Largest Number (179)](https://leetcode.com/problems/largest-number/) · **Template:** Template 1 (sort by a justified key, then sweep): the key is a custom comparator

Problem: arrange non-negative integers to form the largest concatenated number, returned as a string.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** the output is an ordering of the inputs, so look for a pairwise rule justified by an exchange argument: a custom sort comparator.

**Approach:** put a before b when `a + b > b + a` as strings. This is a valid total order (transitive; the proof treats concatenation as `a * 10^len(b) + b`). Exchange: swapping adjacent a, b that violate the order makes the number larger, so the optimum is sorted under it.

```python
def largest_number(nums: list[int]) -> str:
    strs = [str(x) for x in nums]

    def cmp(a: str, b: str) -> int:
        # a before b iff a+b > b+a
        return -1 if a + b > b + a else (1 if a + b < b + a else 0)

    strs.sort(key=cmp_to_key(cmp))
    out = "".join(strs)
    return "0" if out[0] == "0" else out
```

**Complexity:** O(n log n · L) for string compares of length L.

**Edge cases:** all zeros gives `"00"`, so normalize. Plain descending string sort fails on `[3, 30]` (`"303"` < `"330"`).

</details>

### Template 2: reachability frontier

> - **Example problem:** [Jump Game (55)](https://leetcode.com/problems/jump-game/)
> - **Why this template:** "Maximum jump length at that position" + "can you reach the last index": reachable indices form a prefix, so one scalar `far` suffices.
> - **Worked below:** Example 2.1
> - **Also fits:** [Jump Game II (45)](https://leetcode.com/problems/jump-game-ii/)

```python
def frontier_sweep(nums: list[int]) -> int:
    """Return the farthest index reachable from 0 (stops early if stuck).

    Invariant: after processing i, `far` = max index reachable using
    positions 0..i. If i > far, i is unreachable and so is everything after it.
    """
    far = 0
    for i, step in enumerate(nums):
        if i > far:          # TODO: problem-specific "stuck" handling
            break
        far = max(far, i + step)
    return min(far, len(nums) - 1)
```

Knobs: return a boolean (Jump Game), count layer boundaries (Jump Game II), or treat intervals `[i - r, i + r]` as reach (Minimum Taps, Video Stitching).

#### Example 2.1: Jump Game I and II [Medium]

**Problems:** [Jump Game (55)](https://leetcode.com/problems/jump-game/), [Jump Game II (45)](https://leetcode.com/problems/jump-game-ii/) · **Template:** Template 2 (reachability frontier)

Problem: each cell holds a maximum forward jump. (I) Can you reach the last index from 0? (II) If so, in how few jumps?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** reachability on a line where the reachable set is a prefix, so BFS becomes one scalar frontier.

**Approach:**
- I: track `far`, the farthest reachable index. An index beyond `far` means you are stuck.
- II: implicit BFS. Layer k (indices reachable in exactly k jumps) is contiguous. Sweep i through the current layer, computing `far` (the next layer's end); when i hits the layer's end, jump. Stays ahead: by induction, after k jumps greedy's frontier is the maximum reach of any k-jump path.

```python
def can_jump(nums: list[int]) -> bool:
    far = 0
    for i, step in enumerate(nums):
        if i > far:
            return False
        far = max(far, i + step)
    return True


def min_jumps(nums: list[int]) -> int:
    """BFS by layers without a queue: [start, cur_end] is the current layer."""
    jumps, cur_end, far = 0, 0, 0
    for i in range(len(nums) - 1):          # never need to jump from the last index
        far = max(far, i + nums[i])
        if i == cur_end:                    # finished this layer
            jumps += 1
            cur_end = far
    return jumps
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** a single element needs 0 jumps (the loop runs zero times). Loop to `len(nums) - 1`, not `len(nums)`, or you count an extra jump when `cur_end` lands exactly on the last index. If II doesn't guarantee reachability, return -1 when `far <= i` at a boundary.

</details>

#### Example 2.2: Gas Station [Medium]

**Problem:** [Gas Station (134)](https://leetcode.com/problems/gas-station/) · **Template:** closest is Template 2 (reachability frontier): a one-pass scalar (running balance) that restarts when it goes negative

Problem: stations on a circle; station i provides `gas[i]`, and driving to i+1 costs `cost[i]`. Return the unique start index that allows a full loop, or -1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** circular route with a running balance that must never go negative.

**Approach:** two facts.
1. If `sum(gas) < sum(cost)` no start works; otherwise one does.
2. If starting at s first goes negative after station i, **no start in [s, i]** works: you reached any intermediate k with a non-negative balance, so starting there with 0 fuel is no better. Skip to i+1.

So the answer is the station after the prefix balance's global minimum.

```python
def can_complete_circuit(gas: list[int], cost: list[int]) -> int:
    if sum(gas) < sum(cost):
        return -1
    start, tank = 0, 0
    for i in range(len(gas)):
        tank += gas[i] - cost[i]
        if tank < 0:            # no start in [start, i] works; restart after i
            start, tank = i + 1, 0
    return start
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** no wrap-around check needed: total surplus >= 0 and [start, n) has a non-negative running sum, so it covers the deficit before `start`. The companion tests check this against an O(n^2) brute force.

</details>

#### Example 2.3: Partition Labels [Medium]

**Problem:** [Partition Labels (763)](https://leetcode.com/problems/partition-labels/) · **Template:** Template 2 (reachability frontier): `end` grows like `far`, and a cut happens when `i` reaches it

Problem: split a string into as many parts as possible so each letter appears in at most one part. Return the part sizes.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a part containing a letter must extend to its last occurrence: interval merging in disguise, with each letter an interval [first, last].

**Approach:** precompute last occurrences. Sweep, extending `end` to cover every letter seen. At `i == end`, nothing in the current part appears later, so cut. Cutting as early as possible maximizes the part count (stays ahead).

```python
def partition_labels(s: str) -> list[int]:
    last = {ch: i for i, ch in enumerate(s)}
    sizes, start, end = [], 0, 0
    for i, ch in enumerate(s):
        end = max(end, last[ch])            # this part must reach ch's last spot
        if i == end:
            sizes.append(end - start + 1)
            start = i + 1
    return sizes
```

**Complexity:** O(n) time, O(alphabet) space.

**Edge cases:** a single character gives `[1]`; all-distinct letters give all 1s; a letter whose last occurrence is at the end forces one big part.

</details>

### Template 3: greedy with a heap (regret / best available)

> - **Example problem:** [IPO (502)](https://leetcode.com/problems/ipo/)
> - **Why this template:** "At most k projects" where capital unlocks projects: push newly affordable profits into a max-heap and pop the best (unlock and pick best).
> - **Worked below:** Example 3.2
> - **Also fits:** [Course Schedule III (630)](https://leetcode.com/problems/course-schedule-iii/) (the regret shape, which is the template's exact code)

```python
def greedy_with_heap(events: list[tuple[int, int]]) -> int:
    """Process items in deadline order, keep a heap of chosen costs; when a
    constraint breaks, evict the worst chosen item (the "regret" move).

    events = (cost, deadline). Max number of items such that every prefix of
    chosen items finishes by its deadline (Course Schedule III shape).
    """
    total = 0
    heap: list[int] = []                     # max-heap of chosen costs (negated)
    for cost, deadline in sorted(events, key=lambda e: e[1]):
        heapq.heappush(heap, -cost)
        total += cost
        if total > deadline:                 # TODO: constraint check
            total += heapq.heappop(heap)     # evict the most expensive (adds -cost)
    return len(heap)
```

Two shapes of heap-greedy (`11-heaps-top-k-and-merging.md`):

- **Regret**: tentatively take every item; when a constraint breaks, undo the worst choice so far. Course Schedule III, Minimum Number of Refueling Stops (take the largest passed station only when you run dry), Furthest Building You Can Reach.
- **Unlock and pick best**: sort by unlock threshold; each step, push everything now affordable and pop the best. IPO, Single-Threaded CPU.

#### Example 3.1: Task Scheduler [Medium]

**Problem:** [Task Scheduler (621)](https://leetcode.com/problems/task-scheduler/) · **Template:** Template 3 (greedy with a heap) for the simulation; the formula needs no template

Problem: letter-labeled tasks take one tick each; identical labels need at least n ticks between them. Minimize total ticks, idles included.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** cooldown between identical items, minimize length. The most frequent label is the bottleneck.

**Approach:** formula. Let `max_c` be the highest frequency and `num_max` the number of labels with it. Lay out `max_c - 1` frames of width `n + 1`, plus a last frame holding the `num_max` top labels:

```
n = 2, tasks A A A B B B C
A B C | A B _ | A B
frame   frame   tail      (3-1)*(2+1) + 2 = 8
```

Other labels fill idle slots. With more tasks than slots there are no idles, giving `len(tasks)`. Hence `max(len(tasks), (max_c - 1) * (n + 1) + num_max)`.

Simulation: each tick, run the most frequent available label (max-heap) and park it in a FIFO cooldown queue until `time + n`. Most frequent first because that label is most likely to force idles later. The simulation generalizes (return the actual order, Reorganize String, Rearrange String k Distance Apart).

```python
def least_interval(tasks: list[str], n: int) -> int:
    counts = Counter(tasks).values()
    max_c = max(counts)
    num_max = sum(1 for c in counts if c == max_c)
    # (max_c - 1) full frames of width n+1, then a final partial frame.
    return max(len(tasks), (max_c - 1) * (n + 1) + num_max)


def least_interval_sim(tasks: list[str], n: int) -> int:
    """Simulation: each tick run the most frequent available task."""
    heap = [-c for c in Counter(tasks).values()]
    heapq.heapify(heap)
    cooling: deque[tuple[int, int]] = deque()  # (time_available, -remaining), FIFO
    time = 0
    while heap or cooling:
        time += 1
        if heap:
            c = heapq.heappop(heap) + 1     # run one unit (counts are negative)
            if c:
                cooling.append((time + n, c))
        if cooling and cooling[0][0] == time:
            heapq.heappush(heap, cooling.popleft()[1])
    return time
```

**Complexity:** formula O(T) to count; simulation O(T log 26) = O(T).

**Edge cases:** n = 0 gives `len(tasks)`. Cooldown release times are increasing, so a FIFO deque suffices.

</details>

#### Example 3.2: IPO [Hard]

**Problem:** [IPO (502)](https://leetcode.com/problems/ipo/) · **Template:** Template 3 (greedy with a heap): unlock and pick best

Problem: start with capital w; each project has a capital requirement and a profit. Do at most k projects to maximize final capital.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** at most k picks, where each pick's gain unlocks more options: keep a max-heap of everything currently affordable and take the best.

**Approach:** unlock-and-pick-best. Sort projects by requirement. Each round, push newly affordable profits into a max-heap and take the largest. Exchange: the most profitable affordable project gives the most capital, which only unlocks more options.

```python
def find_maximized_capital(k: int, w: int, profits: list[int], capital: list[int]) -> int:
    projects = sorted(zip(capital, profits))
    avail: list[int] = []                   # max-heap of profits we can afford
    i = 0
    for _ in range(k):
        while i < len(projects) and projects[i][0] <= w:
            heapq.heappush(avail, -projects[i][1])
            i += 1
        if not avail:
            break
        w -= heapq.heappop(avail)
    return w
```

**Complexity:** O(n log n + k log n).

**Edge cases:** if nothing is affordable, stop early (`break`); k larger than the number of projects just exhausts the heap.

</details>

### Template 4: counterexample search harness

> - **Example problem:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/)
> - **Why this template:** Tempting rules (shortest first, earliest start first) pass the prompt's examples; a subset brute force breaks them in seconds and supports sort-by-end.
> - **Worked below:** Example 4.1

```python
def find_counterexample(greedy: Callable, brute: Callable, gen: Callable, trials: int = 2000, seed: int = 0):
    """Return the first input where greedy != brute, or None.

    gen(rng) must produce a small input (n <= 8 or so) so brute force is cheap.
    """
    rng = random.Random(seed)
    for _ in range(trials):
        x = gen(rng)
        if greedy(x) != brute(x):
            return x
    return None
```

While practicing, write a 5-line exponential brute force (all subsets / permutations) and let the harness attack your greedy (Example 4.1; the companion file also catches largest-coin-first for coins `[1, 3, 4]`). In a live interview, do this by hand on 3-element inputs.

#### Example 4.1: Test a greedy before trusting it [Medium]

**Problem:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/) · **Template:** Template 4 (counterexample search harness)

**Problem.** Remove the fewest intervals so the rest don't overlap (touching endpoints is fine); equivalently, keep the most. Two candidate rules: keep the shortest interval first, or sort by end. Which one is right? `[[1,2],[2,3],[3,4],[1,3]]` → remove 1.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** several plausible sort keys, and the examples don't separate them → attack each with brute force before writing a proof.

**Approach:** the brute force tries all 2^n subsets and keeps the largest pairwise-compatible one (sorted by start, each must end by the next one's start). `gen_intervals` makes n ≤ 7 inputs, so a trial costs ~128 subsets. The harness compares each greedy with the brute force on 2000 seeded inputs:

- Shortest first fails. Seed 0's first failure is `[[2,4],[9,12],[5,9],[8,10]]`: `[8,10]` (length 2) is taken early and blocks both `[5,9]` and `[9,12]`, so greedy keeps 2 vs. optimal 3. Hand-minimized: `[0,4],[4,8],[3,5]`.
- Sort by end (`max_non_overlapping`, Example 1.2) survives every trial. That's evidence, not proof; the exchange argument in Example 1.2 is the proof.

```python
def max_non_overlapping_brute(intervals: list[list[int]]) -> int:
    best = 0
    n = len(intervals)
    for mask in range(1 << n):
        pick = sorted(intervals[i] for i in range(n) if mask >> i & 1)
        if all(pick[i][1] <= pick[i + 1][0] for i in range(len(pick) - 1)):
            best = max(best, len(pick))
    return best


def shortest_first(intervals: list[list[int]]) -> int:
    """Tempting but WRONG greedy: keep the shortest compatible intervals first."""
    taken = []
    for s, e in sorted(intervals, key=lambda iv: iv[1] - iv[0]):
        if all(e <= a or s >= b for a, b in taken):
            taken.append((s, e))
    return len(taken)


def gen_intervals(rng: random.Random) -> list[list[int]]:
    """Small random intervals (n <= 7) so the 2^n brute force stays cheap."""
    out = []
    for _ in range(rng.randint(0, 7)):
        s = rng.randint(0, 10)
        out.append([s, s + rng.randint(1, 5)])
    return out


def attack_interval_greedies(trials: int = 2000):
    """Run the harness on both rules. Returns (counterexample or None) per rule."""
    bad = find_counterexample(shortest_first, max_non_overlapping_brute, gen_intervals, trials)
    good = find_counterexample(max_non_overlapping, max_non_overlapping_brute, gen_intervals, trials)
    return bad, good
```

**Complexity:** O(trials · 2^n · n log n) for the brute force, about 2000 · 128 · 20 operations here: seconds. Each greedy is O(n log n) or O(n^2).

**Edge cases:** include empty and single-interval inputs in the generator (n can be 0); keep coordinates small so ties and touching endpoints actually occur; fix the seed so a failure reproduces.

</details>

### Template 5: two-pass constraints (left sweep, then right sweep)

> - **Example problem:** [Candy (135)](https://leetcode.com/problems/candy/)
> - **Why this template:** "Higher rating than a neighbor gets more" constrains each child from both sides; a forward pass handles left neighbors, a backward pass handles right neighbors, and `max` combines them.
> - **Worked below:** Example 5.1
> - **Also fits:** [Product of Array Except Self (238)](https://leetcode.com/problems/product-of-array-except-self/)

```python
def two_pass(nums: list[int]) -> list[int]:
    """Satisfy left-neighbor constraints forward, right-neighbor ones backward.

    Invariant: after the forward pass, left[i] is the least value meeting every
    constraint from the left; the backward pass does the same from the right.
    """
    n = len(nums)
    left = [1] * n                          # TODO: base value per problem
    for i in range(1, n):                   # forward: depends on i - 1 only
        if nums[i] > nums[i - 1]:           # TODO: left-neighbor constraint
            left[i] = left[i - 1] + 1
    right = [1] * n
    for i in range(n - 2, -1, -1):          # backward: depends on i + 1 only
        if nums[i] > nums[i + 1]:           # TODO: right-neighbor constraint
            right[i] = right[i + 1] + 1
    return [max(a, b) for a, b in zip(left, right)]   # TODO: combine (max/min/product)
```

Each element is constrained by both neighbors. The forward pass computes the tightest bound coming from the left, the backward pass the tightest bound from the right, and the combine step merges them per element: `max` for lower bounds (Candy), `min` for upper bounds (Trapping Rain Water's water level), product for "everything except me" (Product of Array Except Self).

Why it's correct: each pass enforces one side's constraints and ignores the other. Combining with `max` never breaks either side: if `nums[i] > nums[i - 1]`, the backward pass leaves `i - 1` at its base value, so `left[i] = left[i - 1] + 1` already exceeds both of `i - 1`'s bounds (and symmetrically for right neighbors). Each bound is necessary, so their max is also minimal.

#### Example 5.1: Candy [Hard]

**Problem:** [Candy (135)](https://leetcode.com/problems/candy/) · **Template:** Template 5 (two-pass constraints)

Problem: children in a row have ratings. Each gets at least one candy; a child rated strictly higher than a neighbor gets more than that neighbor. Minimize total candy.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** constraints only between neighbors, in both directions. Satisfy each direction in its own pass.

**Approach:** the left pass enforces "higher than left neighbor gets more" minimally. The right pass does the same for right neighbors, taking `max` so left constraints hold. Each value is the max of two necessary lower bounds (the increasing run ending here from each side), so the total is minimal.

```python
def candy(ratings: list[int]) -> int:
    n = len(ratings)
    give = [1] * n
    for i in range(1, n):                   # satisfy left-neighbor constraints
        if ratings[i] > ratings[i - 1]:
            give[i] = give[i - 1] + 1
    for i in range(n - 2, -1, -1):          # satisfy right-neighbor, keep left
        if ratings[i] > ratings[i + 1]:
            give[i] = max(give[i], give[i + 1] + 1)
    return sum(give)
```

**Complexity:** O(n) time, O(n) space (an O(1) slope-counting version exists; not worth it in interviews).

**Edge cases:** equal neighbors impose no constraint (`[1, 2, 2]` gives 4). Same two-pass idea: Product of Array Except Self, Trapping Rain Water (max from each side).

</details>

## Template cheat sheet

### Template 1: sort by a justified key, then sweep
```python
def greedy_sort_and_sweep(items: list, key: Callable, can_take: Callable, take: Callable, state):
    """Generic sort-then-sweep greedy.

    key:      the ordering that the exchange argument justifies (TODO per problem)
    can_take: (state, item) -> bool, is the item compatible with choices so far
    take:     (state, item) -> new state
    """
    chosen = []
    for item in sorted(items, key=key):
        if can_take(state, item):
            state = take(state, item)
            chosen.append(item)
    return chosen, state
```

### Template 2: reachability frontier
```python
def frontier_sweep(nums: list[int]) -> int:
    """Return the farthest index reachable from 0 (stops early if stuck).

    Invariant: after processing i, `far` = max index reachable using
    positions 0..i. If i > far, i is unreachable and so is everything after it.
    """
    far = 0
    for i, step in enumerate(nums):
        if i > far:          # TODO: problem-specific "stuck" handling
            break
        far = max(far, i + step)
    return min(far, len(nums) - 1)
```

### Template 3: greedy with a heap (regret / best available)
```python
def greedy_with_heap(events: list[tuple[int, int]]) -> int:
    """Process items in deadline order, keep a heap of chosen costs; when a
    constraint breaks, evict the worst chosen item (the "regret" move).

    events = (cost, deadline). Max number of items such that every prefix of
    chosen items finishes by its deadline (Course Schedule III shape).
    """
    total = 0
    heap: list[int] = []                     # max-heap of chosen costs (negated)
    for cost, deadline in sorted(events, key=lambda e: e[1]):
        heapq.heappush(heap, -cost)
        total += cost
        if total > deadline:                 # TODO: constraint check
            total += heapq.heappop(heap)     # evict the most expensive (adds -cost)
    return len(heap)
```

### Template 4: counterexample search harness
```python
def find_counterexample(greedy: Callable, brute: Callable, gen: Callable, trials: int = 2000, seed: int = 0):
    """Return the first input where greedy != brute, or None.

    gen(rng) must produce a small input (n <= 8 or so) so brute force is cheap.
    """
    rng = random.Random(seed)
    for _ in range(trials):
        x = gen(rng)
        if greedy(x) != brute(x):
            return x
    return None
```

### Template 5: two-pass constraints (left sweep, then right sweep)
```python
def two_pass(nums: list[int]) -> list[int]:
    """Satisfy left-neighbor constraints forward, right-neighbor ones backward.

    Invariant: after the forward pass, left[i] is the least value meeting every
    constraint from the left; the backward pass does the same from the right.
    """
    n = len(nums)
    left = [1] * n                          # TODO: base value per problem
    for i in range(1, n):                   # forward: depends on i - 1 only
        if nums[i] > nums[i - 1]:           # TODO: left-neighbor constraint
            left[i] = left[i - 1] + 1
    right = [1] * n
    for i in range(n - 2, -1, -1):          # backward: depends on i + 1 only
        if nums[i] > nums[i + 1]:           # TODO: right-neighbor constraint
            right[i] = right[i + 1] + 1
    return [max(a, b) for a, b in zip(left, right)]   # TODO: combine (max/min/product)
```

## When greedy fails

Greedy fails when an early cheap choice destroys a better later combination. Counterexamples to memorize:

| Problem | Tempting greedy | Counterexample | Correct tool |
|---|---|---|---|
| Coin Change, coins [1, 3, 4] | largest coin first | amount 6: greedy 4+1+1 (3 coins), optimal 3+3 (2) | DP over amount (`16-dynamic-programming-foundations.md`) |
| Coin Change, coins [2, 5] | largest coin first | amount 6: greedy takes 5, stuck; optimal 2+2+2 | DP |
| 0/1 Knapsack | best value/weight ratio | cap 10, items (w6,v7), (w5,v5), (w5,v5): ratio takes 7, optimal 10 | DP (fractional knapsack is greedy) |
| Interval scheduling | shortest first | [0,4], [4,8], [3,5] | sort by end |
| Interval scheduling | earliest start | [0,10], [1,2], [3,4] | sort by end |
| Weighted interval scheduling | earliest end | (0,3,w1), (0,5,w100) | DP + binary search (`17-dynamic-programming-advanced.md`) |
| Min path sum in a grid | step to the cheaper neighbor | cheap first step leading into an expensive region | DP |
| Boats with capacity > 2 | first-fit decreasing | bin packing is NP-hard | backtracking / DP over subsets |

"Canonical" coin systems like US coins [1, 5, 10, 25] happen to be greedy-optimal; never assume this for arbitrary inputs.

```python
def coin_change_greedy(coins: list[int], amount: int) -> int:
    count = 0
    for c in sorted(coins, reverse=True):
        count += amount // c
        amount %= c
    return count if amount == 0 else -1


def coin_change_dp(coins: list[int], amount: int) -> int:
    INF = float("inf")
    dp = [0] + [INF] * amount
    for a in range(1, amount + 1):
        for c in coins:
            if c <= a and dp[a - c] + 1 < dp[a]:
                dp[a] = dp[a - c] + 1
    return dp[amount] if dp[amount] != INF else -1
```

### Decision procedure (use it every time)

1. **Propose** a rule: what to sort by, what to pick first.
2. **State the invariant** you believe it maintains.
3. **Attack it** on 3–4 adversarial items: one long vs. many short, one big vs. many small, ties, equal ratios, items that only pay off together.
4. **Try to prove it** (exchange or stays ahead). A failing swap step usually is the counterexample.
5. **If it breaks, fall back to DP.** The counterexample names the state: if a choice depends on "capacity left" or "which items remain", that is your DP dimension.
6. **Say it out loud:** "Greedy by end time; exchange argument: swapping OPT's first interval for the earliest-ending one keeps feasibility." That sentence often turns a pass into a strong pass.

## Variations and extensions

- **+ two pointers**: boats, Two City Scheduling (sort by cost difference), Bag of Tokens (`05-two-pointers.md`).
- **+ heap**: IPO, Course Schedule III, Refueling Stops, Furthest Building You Can Reach, Reorganize String (`11-heaps-top-k-and-merging.md`).
- **+ monotonic stack**: Remove K Digits, Remove Duplicate Letters: pop larger digits while removals remain (`08-stacks-and-monotonic-structures.md`).
- **+ binary search on the answer**: "minimize the maximum" (Split Array Largest Sum, Capacity to Ship Packages) with a greedy pack-until-full feasibility check (`07-binary-search.md`).
- **In graphs**: Dijkstra, Prim, and Kruskal are greedy with exchange-argument (cut property) proofs (`14-union-find-shortest-paths-mst.md`).
- **Range greedy**: Valid Parenthesis String tracks the [min_open, max_open] range of possible balances.

Follow-ups:
- "Prove it." → The exchange argument in two sentences.
- "What if items arrive as a stream?" → Sort-by-end needs all items; online, use a heap (keep the k cheapest, evict on conflict) or accept approximation.
- "Return the actual schedule." → Record choices in the sweep; heap simulations produce the order naturally.

## Common pitfalls

- Wrong sort key (start instead of end, value instead of ratio). Test it against a counterexample shape before coding.
- `>` vs `>=` for touching intervals: touching balloons share an arrow; touching non-overlapping intervals are allowed.
- Jump Game II: iterating to `len(nums)` instead of `len(nums) - 1` adds a spurious jump.
- Gas Station: skipping the `sum(gas) < sum(cost)` check, or simulating from every start (O(n^2)).
- Candy: forgetting `max` in the second pass breaks left constraints.
- Largest Number: not normalizing `"000"` to `"0"`. Python 3 `sort` has no `cmp=`; use `functools.cmp_to_key`.
- Task Scheduler: forgetting `max(len(tasks), ...)`, which matters with many distinct labels.
- Trusting a greedy because it passes the prompt's examples; they are chosen to be easy.

## When NOT to use / look-alikes

| Looks like greedy | Actually | Tell |
|---|---|---|
| Coin change, knapsack, "number of ways" | DP | Arbitrary weights/values; choices interact through a capacity |
| Max concurrent meetings / rooms needed | Sweep line or min-heap of end times | Asks for a peak count, not a selection |
| Jump Game with costs (min cost to reach end) | DP or Dijkstra | Costs break "farther is always better" |
| Max sum with adjacency bans (House Robber) | DP | Local best excludes neighbors that might be better together |
| Scheduling with weights and deadlines | Greedy for unit-time jobs (sort by profit, latest free slot); else usually DP | Unit lengths often keep greedy valid; general lengths usually don't |
| Shortest path with negative edges | Bellman-Ford, not Dijkstra | Greedy "closest first" breaks with negatives |

Rule of thumb: greedy is right when the state after a choice is **one scalar you want as small (or large) as possible**. If the future needs more, use DP.

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Jump Game (55)](https://leetcode.com/problems/jump-game/) · Medium · Reachability frontier: track `far`.
2. [Jump Game II (45)](https://leetcode.com/problems/jump-game-ii/) · Medium · BFS-like levels / frontier: count BFS layers with `cur_end`.
3. [Gas Station (134)](https://leetcode.com/problems/gas-station/) · Medium · Running deficit reset: reset after the balance goes negative.
4. [Partition Labels (763)](https://leetcode.com/problems/partition-labels/) · Medium · Last occurrence: extends the current cut.
5. [Candy (135)](https://leetcode.com/problems/candy/) · Hard · Two-pass constraints: left pass, right pass with max.

### Optional

**Interval selection**
- [Meeting Rooms (252)](https://leetcode.com/problems/meeting-rooms/) · Easy · Premium · Sort by start, check adjacent overlap.
- [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/) · Medium · Sort by end, count keepers.
- [Minimum Number of Arrows to Burst Balloons (452)](https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/) · Medium · Sort by end, new arrow when start > last.

**Reachability frontier**
- [Video Stitching (1024)](https://leetcode.com/problems/video-stitching/) · Medium · Intervals as reach; same as Jump Game II after sorting by start.
- [Minimum Number of Taps to Open to Water a Garden (1326)](https://leetcode.com/problems/minimum-number-of-taps-to-open-to-water-a-garden/) · Hard · Convert each tap to an interval, then Jump Game II.

**Running balance**
- [Valid Parenthesis String (678)](https://leetcode.com/problems/valid-parenthesis-string/) · Medium · Track a range of possible open counts.

**Matching by sorting**
- [Assign Cookies (455)](https://leetcode.com/problems/assign-cookies/) · Easy · Sort both, smallest cookie to least greedy child.
- [Boats to Save People (881)](https://leetcode.com/problems/boats-to-save-people/) · Medium · Heaviest with lightest.
- [Two City Scheduling (1029)](https://leetcode.com/problems/two-city-scheduling/) · Medium · Sort by cost difference.
- [Hand of Straights (846)](https://leetcode.com/problems/hand-of-straights/) · Medium · Smallest card must start a run.

**Counting / heap**
- [Task Scheduler (621)](https://leetcode.com/problems/task-scheduler/) · Medium · Frame formula or heap simulation.
- [Reorganize String (767)](https://leetcode.com/problems/reorganize-string/) · Medium · Possible iff max count <= (n + 1) // 2; place most frequent first.
- [IPO (502)](https://leetcode.com/problems/ipo/) · Hard · Unlock by capital, pop max profit.
- [Course Schedule III (630)](https://leetcode.com/problems/course-schedule-iii/) · Hard · Sort by deadline, evict the longest course on overflow.
- [Minimum Number of Refueling Stops (871)](https://leetcode.com/problems/minimum-number-of-refueling-stops/) · Hard · Regret heap of passed stations.

**Two-pass and comparators**
- [Largest Number (179)](https://leetcode.com/problems/largest-number/) · Medium · Comparator `a+b > b+a`.
- [Queue Reconstruction by Height (406)](https://leetcode.com/problems/queue-reconstruction-by-height/) · Medium · Sort tall-first, insert at index k.
- [Remove K Digits (402)](https://leetcode.com/problems/remove-k-digits/) · Medium · Monotonic stack greedy.

## Self-check

1. State the exchange argument for "sort intervals by end" in two sentences. Why does sorting by start fail?
2. In Jump Game II, what exactly do `cur_end` and `far` represent, and why is the reachable set always a contiguous prefix?
3. In Gas Station, why can you skip every start between the old start and the index where the tank went negative?
4. Derive the Task Scheduler formula. When does `len(tasks)` win the `max`?
5. Give a coin system and amount where largest-coin-first fails. What DP state fixes it?
6. Why does Candy need two passes, and why `max` in the second?
7. You have proposed a greedy rule in an interview. What are the next three things you do before writing code?
