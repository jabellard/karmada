# 12. Intervals and Sweep Line

Scheduling, reservations, resource allocation, log time ranges, and IP/port ranges are all intervals on a line. Pick the right sort key and sweep left to right with a little state. The same sweep powers difference arrays for bulk range updates.

## At a glance

- **Use when:** a list of `[start, end]` pairs; overlap, conflict, booking, rooms, free time, coverage.
- **Templates:**
  - Template 1: merge (sort by start): union, insert, gaps.
  - Template 2: event sweep: peak overlap, minimum rooms, capacity.
  - Template 3: sort-by-end greedy: max non-overlapping, min removals, min arrows.
  - Template 4: difference array: many range additions on a small integer axis.
- **Typical complexity:** O(n log n) sort; O(n + k) difference array.
- **Top trap:** endpoint semantics. Decide half-open vs closed first; it fixes `<` vs `<=` and tie order.

## Core idea

Sort, then sweep with a small summary of everything seen.

- **Merge (sort by start).** A new interval can only overlap the block being built. Invariant: output is sorted and disjoint; only its last element can grow.
- **Selection (sort by end).** To keep the most non-overlapping intervals, take the one finishing first. Exchange argument: any optimal solution can swap its first interval for the earliest-ending one at no loss.
- **Counting (event sweep).** `+1` at each start, `-1` at each end. Sorted by time, the running sum counts active intervals; its maximum is the peak overlap.
- **Range update (difference array).** `diff[l] += v; diff[r+1] -= v` adds v to [l, r]; a prefix sum rebuilds the array. It's the event sweep on an indexed axis.

```
time:     0    5    10   15   20   25   30
A [0,30)  |==============================|
B [5,10)       |====|
C [15,20)                |====|
events:  +1   +1   -1   +1   -1        -1
active:   1    2    1    2    1         0     peak = 2
```

### Choosing the sort key and tie-breakers

| Goal | Sort key |
|---|---|
| Merge / union / insert / free time | start |
| Max non-overlapping, min removals, min arrows, min points to stab | end |
| Peak overlap / rooms / capacity | events by time |
| Remove covered intervals | start asc, end desc (the widest of equal starts covers the rest) |

Most bugs are ties at equal time. Decide: **half-open** `[s, e)` (a meeting ending at 10 doesn't conflict with one starting at 10) or **closed** `[s, e]` (balloons touching at 10 share an arrow)?

- Half-open counting: **end before start** at equal times. Python sorts `(t, -1)` before `(t, +1)` already.
- Closed counting: **start before end**, or shift integer ends to `e + 1`.
- Merge: `s <= last_end` merges touching intervals; `s < last_end` doesn't. The problem's examples decide.

## Recognition signals

- Pairs like `[start, end]`, `[from, to]`, `[open, close]`, `(ts_begin, ts_end)`.
- Words: overlap, conflict, schedule, booking, meeting, coverage, gap, free time, busy, concurrent.
- "Minimum rooms/servers/machines" → peak overlap (event sweep or heap of ends).
- "Maximum non-overlapping", "minimum removals", "minimum arrows to hit all" → sort-by-end greedy.
- Many range increments, coordinates up to ~10^6 → difference array.
- Online "can I book this?" → sorted structure with `bisect`.
- n up to 10^5 → O(n log n); O(n^2) pairwise checks time out.

## Disguises

1. *A CI cluster logs each job as (queued_at, finished_at). What is the smallest number of runners that would have let every job start the moment it was queued?*
   Decode: each job holds a runner for a range, so "fewest runners" is peak overlap. Event sweep or heap of ends.

2. *An on-call rotation file has shifts from several people, possibly overlapping. List the time windows when nobody was on call.*
   Decode: merge all shifts, emit the gaps (Employee Free Time).

3. *A shuttle has 40 seats. Given group bookings (size, pickup_km, dropoff_km), does any stretch of road exceed capacity?*
   Decode: range additions plus a capacity check → difference array or event sweep. Drop-off frees seats, so half-open.

4. *A conference hall can host one talk at a time. Talks have fixed times. Pick as many talks as possible.*
   Decode: max non-overlapping → sort-by-end greedy.

5. *A firewall has a rule list of port ranges. Compress it into the fewest equivalent ranges, treating adjacent ranges like 80-89 and 90-99 as one.*
   Decode: merge, with `s <= last_end + 1` for integer adjacency.

6. *Two engineers each publish their sorted list of free slots. Find every slot when both are free.*
   Decode: intersect two sorted disjoint lists → two pointers, advance the one ending first.

## Templates and worked examples

Imports (companion file: `code/12_intervals_sweep.py`):

```python
import bisect
import heapq
from itertools import accumulate
```

### Template 1: merge (sort by start)

> - **Example problem:** [Merge Intervals (56)](https://leetcode.com/problems/merge-intervals/)
> - **Why this template:** "merge overlapping intervals" = sort by start, extend the last block while the next starts inside it.
> - **Worked below:** Example 1.1
> - **Also fits:** [Insert Interval (57)](https://leetcode.com/problems/insert-interval/)

```python
def merge_template(intervals: list[list[int]]) -> list[list[int]]:
    """Sort by start, then extend or close the current block.

    Invariant: `out` holds disjoint, sorted intervals covering everything seen;
    only out[-1] can still grow, because later starts are >= its start.
    """
    out: list[list[int]] = []
    for s, e in sorted(intervals):            # TODO: sort key (start)
        if out and s <= out[-1][1]:           # TODO: '<=' merges touching [1,2],[2,3]; '<' does not
            out[-1][1] = max(out[-1][1], e)   # TODO: combine payload (max end, sum weight, ...)
        else:
            out.append([s, e])
    return out
```

Knobs: overlap test (`<=`, `<`, or `<= last + 1`); merge action (max end, count, sum); output (blocks or gaps).

#### Example 1.1: Merge overlapping intervals [Easy/Medium]

**Problem:** [Merge Intervals (56)](https://leetcode.com/problems/merge-intervals/) · **Template:** Template 1 (merge, sort by start)

Problem: given unsorted closed intervals, return their union as sorted disjoint intervals. Touching intervals (`[1,4]`, `[4,5]`) merge.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "union", "combine overlapping" → Template 1.

**Approach:** sort by start; extend the last block while the next start is `<=` its end. `max` on the end handles containment (`[1,4]` swallowing `[2,3]`).

```python
def merge_intervals(intervals: list[list[int]]) -> list[list[int]]:
    """Merge Intervals."""
    return merge_template(intervals)
```

**Complexity:** O(n log n) time, O(n) output.

**Edge cases:** empty input; containment; touching intervals; unsorted input.

</details>

#### Example 1.2: Insert into a sorted disjoint list [Medium]

**Problem:** [Insert Interval (57)](https://leetcode.com/problems/insert-interval/) · **Template:** Template 1 (merge, sort by start), on input that is already sorted

Problem: insert a new interval into disjoint intervals sorted by start, merging as needed.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** already sorted, so don't re-sort. O(n) three-phase scan.

**Approach:** copy intervals ending before the new one; absorb those starting `<=` its end; copy the rest. Each phase is contiguous because the list is sorted and disjoint.

```python
def insert_interval(intervals: list[list[int]], new: list[int]) -> list[list[int]]:
    """Insert Interval: input sorted and disjoint. Three phases, O(n)."""
    out: list[list[int]] = []
    i, n = 0, len(intervals)
    s, e = new
    while i < n and intervals[i][1] < s:      # entirely left of new
        out.append(intervals[i])
        i += 1
    while i < n and intervals[i][0] <= e:     # overlaps new: absorb
        s = min(s, intervals[i][0])
        e = max(e, intervals[i][1])
        i += 1
    out.append([s, e])
    out.extend(intervals[i:])                 # entirely right of new
    return out
```

**Complexity:** O(n) time and space. `bisect` finds the phase boundaries in O(log n), but building the output is O(n).

**Edge cases:** empty list; new interval before or after all; new interval swallowing everything.

</details>

### Template 2: event sweep (peak overlap)

> - **Example problem:** [Meeting Rooms II (253)](https://leetcode.com/problems/meeting-rooms-ii/) (Premium)
> - **Why this template:** "minimum rooms" = most meetings open at once = +1/-1 event sweep.
> - **Worked below:** Example 2.1
> - **Also fits:** [Car Pooling (1094)](https://leetcode.com/problems/car-pooling/)

```python
def event_sweep_template(intervals: list[tuple[int, int]]) -> int:
    """Max number of simultaneously active half-open intervals [s, e).

    Events: (time, delta). Sorting tuples puts (t, -1) before (t, +1), so an
    interval ending at t frees its slot before one starting at t takes it.
    """
    events: list[tuple[int, int]] = []
    for s, e in intervals:
        events.append((s, +1))
        events.append((e, -1))                # TODO: closed intervals [s, e] -> use (e + 1, -1) or flip tie order
    events.sort()                             # TODO: tie-breaker encoded in the delta sign
    active = best = 0
    for _, d in events:
        active += d
        best = max(best, active)              # TODO: record time, check capacity, etc.
    return best
```

Knobs: weighted deltas (passengers, CPU); tie order; what to track (peak, first time over capacity, total busy time). To know *which* resource is reused, use a min-heap of end times (Example 2.1).

#### Example 2.1: Overlap detection and peak overlap [Easy, Medium]

**Problems:** [Meeting Rooms (252)](https://leetcode.com/problems/meeting-rooms/) (Premium), [Meeting Rooms II (253)](https://leetcode.com/problems/meeting-rooms-ii/) (Premium) · **Template:** Template 2 (event sweep, peak overlap)

Problem: meetings are half-open `[start, end)`. (a) Can one person attend all? (b) What is the minimum number of rooms?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "attend all" = no overlap. "Minimum rooms" = most meetings active at once.

**Approach:** (a) sort by start; any overlap shows up between neighbors.

```python
def can_attend_all(intervals: list[list[int]]) -> bool:
    """Meeting Rooms: any overlap among half-open intervals?"""
    intervals = sorted(intervals)
    return all(intervals[k - 1][1] <= intervals[k][0] for k in range(1, len(intervals)))
```

(b) Three equivalent ways:

1. Event sweep (Template 2). Ends sort first at equal times, so back-to-back meetings share a room.
2. Two sorted arrays (Template 2 without tuples): before opening a room at `s`, release every room ending `<= s`.
3. Min-heap of end times: reuse the earliest-freeing room if free by `s`; the final heap size is the answer. Use it when asked which room each meeting gets (`11-heaps-top-k-and-merging.md`).

```python
def min_meeting_rooms(intervals: list[list[int]]) -> int:
    """Meeting Rooms II via two sorted arrays (starts, ends)."""
    starts = sorted(s for s, _ in intervals)
    ends = sorted(e for _, e in intervals)
    rooms = best = 0
    j = 0
    for s in starts:
        # Free every room whose meeting ended at or before s (half-open).
        while j < len(ends) and ends[j] <= s:
            rooms -= 1
            j += 1
        rooms += 1
        best = max(best, rooms)
    return best


def min_meeting_rooms_heap(intervals: list[list[int]]) -> int:
    """Meeting Rooms II via min-heap of end times (rooms currently in use)."""
    heap: list[int] = []
    for s, e in sorted(intervals):
        if heap and heap[0] <= s:
            heapq.heapreplace(heap, e)        # reuse the room that frees earliest
        else:
            heapq.heappush(heap, e)
    return len(heap)
```

Decoupling starts from ends is valid: the peak depends only on how many starts and ends precede each time, not which end pairs with which start.

**Complexity:** O(n log n) time, O(n) space for all three.

**Edge cases:** n identical meetings need n rooms; back-to-back `[1,5],[5,10]` needs 1 (`<` instead of `<=` gives 2).

</details>

### Template 3: sort-by-end greedy (activity selection)

> - **Example problem:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/)
> - **Why this template:** "minimum removals" = n minus the most non-overlapping intervals kept, which sort-by-end solves.
> - **Worked below:** Example 3.1
> - **Also fits:** [Minimum Number of Arrows to Burst Balloons (452)](https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/)

```python
def sort_by_end_greedy_template(intervals: list[list[int]]) -> int:
    """Max number of pairwise non-overlapping intervals (activity selection).

    Picking the interval that ends earliest leaves the most room for the rest
    (exchange argument). Touching endpoints are allowed here.
    """
    count, last_end = 0, float("-inf")
    for s, e in sorted(intervals, key=lambda iv: iv[1]):   # sort by END
        if s >= last_end:                     # TODO: '>' if touching counts as overlap
            count += 1
            last_end = e
    return count
```

Knobs: `>=` vs `>`; return `count` (max kept), `n - count` (min removed), or `count` with `>` on closed intervals (min stabbing points). Proof and variants: `18-greedy.md`.

#### Example 3.1: Sort by end, min removals and min arrows [Medium]

**Problems:** [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/), [Minimum Number of Arrows to Burst Balloons (452)](https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/) · **Template:** Template 3 (sort-by-end greedy)

Problem (a): remove the fewest intervals so the rest don't overlap; touching is fine. (b): balloons span closed x-ranges; an arrow at x bursts those with `s <= x <= e`. Find the fewest arrows.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** (a) fewest removed = n minus most kept → activity selection. (b) Its dual: fewest stabbing points = most pairwise disjoint balloons.

**Approach:** sort by end. (a) Keep an interval starting at or after the last kept end. (b) Shoot at the earliest unburst end; it hits every balloon starting at or before it, so only one starting strictly after needs a new arrow.

```python
def erase_overlap_intervals(intervals: list[list[int]]) -> int:
    """Non-overlapping Intervals: min removals = n - max kept (touching ok)."""
    return len(intervals) - sort_by_end_greedy_template(intervals)


def min_arrows(points: list[list[int]]) -> int:
    """Minimum Number of Arrows to Burst Balloons: closed intervals."""
    arrows, last = 0, float("-inf")
    for s, e in sorted(points, key=lambda p: p[1]):
        if s > last:                          # strict: touching balloons share an arrow
            arrows += 1
            last = e                          # shoot at the earliest end
    return arrows
```

Why not sort by start? "Keep the first" fails on `[1,100],[2,3],[4,5]`: it keeps the long one and loses two. Sorting by start works only if overlaps keep the smaller end; sort by end is cleaner.

**Complexity:** O(n log n) time, O(1) extra.

**Edge cases:** `>=` vs `>` is the only half-open/closed difference. Coordinates reach `-2^31`; a `float("-inf")` sentinel avoids overflow-style bugs.

</details>

### Template 4: difference array (range additions)

> - **Example problem:** [Car Pooling (1094)](https://leetcode.com/problems/car-pooling/)
> - **Why this template:** many "add p passengers from pickup to dropoff" updates on a small axis, then one scan.
> - **Worked below:** Example 4.1
> - **Also fits:** [Corporate Flight Bookings (1109)](https://leetcode.com/problems/corporate-flight-bookings/)

```python
def difference_array_template(n: int, updates: list[tuple[int, int, int]]) -> list[int]:
    """Apply many range additions [l, r] += v on an array of length n in O(n + k)."""
    diff = [0] * (n + 1)
    for l, r, v in updates:
        diff[l] += v
        diff[r + 1] -= v                      # r inclusive; slot n absorbs overflow
    return list(accumulate(diff[:n]))
```

Knobs: inclusive vs exclusive right end (`r + 1` vs `r`); 1-indexed input; huge coordinates → dict of deltas, sorted keys (Template 2). 2D: four corner updates, then a 2D prefix sum.

#### Example 4.1: Range additions with a difference array [Medium]

**Problems:** [Car Pooling (1094)](https://leetcode.com/problems/car-pooling/), [Corporate Flight Bookings (1109)](https://leetcode.com/problems/corporate-flight-bookings/) · **Template:** Template 4 (difference array)

Problem (a): a vehicle drives one way; trips are `[passengers, pickup, dropoff]`, and riders leave at `dropoff` before others board. Is capacity ever exceeded? (b): flights 1..n; booking `[first, last, seats]` reserves seats on each flight in the inclusive range. Return seats per flight.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** range increments on a bounded axis, then a scan → difference array. "Leave before board" = half-open `[pickup, dropoff)`.

**Approach:** `+p` at pickup, `-p` at dropoff, prefix-sum, check each point. Locations are bounded (~1000), so an array beats sorting events. For (b), convert to 0-indexed and reuse Template 4.

```python
def car_pooling(trips: list[list[int]], capacity: int) -> bool:
    """Car Pooling: trips[i] = [passengers, from, to]; drop-off at `to` frees seats."""
    max_loc = max((t for _, _, t in trips), default=0)
    diff = [0] * (max_loc + 1)
    for p, f, t in trips:
        diff[f] += p
        diff[t] -= p                          # half-open [from, to)
    load = 0
    for d in diff:
        load += d
        if load > capacity:
            return False
    return True


def corp_flight_bookings(bookings: list[list[int]], n: int) -> list[int]:
    """Corporate Flight Bookings: 1-indexed inclusive ranges."""
    return difference_array_template(n, [(f - 1, l - 1, s) for f, l, s in bookings])
```

**Complexity:** O(n + k) for k updates on an axis of size n. For a huge axis (timestamps to 10^9), sort the 2k events: O(k log k).

**Edge cases:** inclusive vs exclusive end; 1-indexing; one extra `diff` slot.

</details>

### Beyond the templates

#### Example B.1: Intersect two sorted lists [Medium]

**Problem:** [Interval List Intersections (986)](https://leetcode.com/problems/interval-list-intersections/) · **Template:** none of this module's templates; two pointers over two sorted lists (`05-two-pointers.md`)

**Why no template:** It's two pointers over two sorted lists, from module 05.

Problem: two lists of closed intervals, each sorted and disjoint. Return all pairwise intersections.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** two sorted sequences walked together → two pointers.

**Approach:** the intersection is `[max(starts), min(ends)]` if non-empty. Advance whichever ends first: the other interval extends past it and later ones start later, so it can't intersect anything else.

```python
def interval_intersections(a: list[list[int]], b: list[list[int]]) -> list[list[int]]:
    """Interval List Intersections: two sorted disjoint lists, closed intervals."""
    i = j = 0
    out: list[list[int]] = []
    while i < len(a) and j < len(b):
        lo = max(a[i][0], b[j][0])
        hi = min(a[i][1], b[j][1])
        if lo <= hi:
            out.append([lo, hi])
        # The one ending first cannot intersect anything else in the other list.
        if a[i][1] < b[j][1]:
            i += 1
        else:
            j += 1
    return out
```

**Complexity:** O(m + n) time and output.

**Edge cases:** single-point intersections like `[5,5]` (closed, so `lo <= hi`); an empty list.

</details>

## Template cheat sheet

All templates, for warm-ups and review.

### Template 1: merge (sort by start)

```python
def merge_template(intervals: list[list[int]]) -> list[list[int]]:
    """Sort by start, then extend or close the current block.

    Invariant: `out` holds disjoint, sorted intervals covering everything seen;
    only out[-1] can still grow, because later starts are >= its start.
    """
    out: list[list[int]] = []
    for s, e in sorted(intervals):            # TODO: sort key (start)
        if out and s <= out[-1][1]:           # TODO: '<=' merges touching [1,2],[2,3]; '<' does not
            out[-1][1] = max(out[-1][1], e)   # TODO: combine payload (max end, sum weight, ...)
        else:
            out.append([s, e])
    return out
```

### Template 2: event sweep (peak overlap)

```python
def event_sweep_template(intervals: list[tuple[int, int]]) -> int:
    """Max number of simultaneously active half-open intervals [s, e).

    Events: (time, delta). Sorting tuples puts (t, -1) before (t, +1), so an
    interval ending at t frees its slot before one starting at t takes it.
    """
    events: list[tuple[int, int]] = []
    for s, e in intervals:
        events.append((s, +1))
        events.append((e, -1))                # TODO: closed intervals [s, e] -> use (e + 1, -1) or flip tie order
    events.sort()                             # TODO: tie-breaker encoded in the delta sign
    active = best = 0
    for _, d in events:
        active += d
        best = max(best, active)              # TODO: record time, check capacity, etc.
    return best
```

### Template 3: sort-by-end greedy (activity selection)

```python
def sort_by_end_greedy_template(intervals: list[list[int]]) -> int:
    """Max number of pairwise non-overlapping intervals (activity selection).

    Picking the interval that ends earliest leaves the most room for the rest
    (exchange argument). Touching endpoints are allowed here.
    """
    count, last_end = 0, float("-inf")
    for s, e in sorted(intervals, key=lambda iv: iv[1]):   # sort by END
        if s >= last_end:                     # TODO: '>' if touching counts as overlap
            count += 1
            last_end = e
    return count
```

### Template 4: difference array (range additions)

```python
def difference_array_template(n: int, updates: list[tuple[int, int, int]]) -> list[int]:
    """Apply many range additions [l, r] += v on an array of length n in O(n + k)."""
    diff = [0] * (n + 1)
    for l, r, v in updates:
        diff[l] += v
        diff[r + 1] -= v                      # r inclusive; slot n absorbs overflow
    return list(accumulate(diff[:n]))
```

## Variations and extensions

**Employee free time (LeetCode: Employee Free Time).** Merge all busy intervals and emit the gaps, or k-way merge the sorted lists with a heap (`11-heaps-top-k-and-merging.md`) in O(N log k).

```python
def employee_free_time(schedules: list[list[list[int]]]) -> list[list[int]]:
    """Employee Free Time: gaps in the union of everyone's busy intervals."""
    all_iv = [iv for person in schedules for iv in person]
    merged = merge_template(all_iv)
    return [[merged[k - 1][1], merged[k][0]] for k in range(1, len(merged))
            if merged[k - 1][1] < merged[k][0]]
```

**Online booking (LeetCode: My Calendar I).** Interleaved inserts and queries rule out sorting once. Keep starts sorted; `bisect` finds the two neighbors that could conflict. List insertion is O(n); a balanced BST or `sortedcontainers.SortedList` is O(log n).

```python
class MyCalendar:
    """My Calendar I: reject a booking that overlaps any existing one (half-open)."""

    def __init__(self) -> None:
        self.starts: list[int] = []
        self.ends: list[int] = []

    def book(self, start: int, end: int) -> bool:
        i = bisect.bisect_right(self.starts, start)
        # Neighbor on the left must end by `start`; neighbor on the right must start at/after `end`.
        if i > 0 and self.ends[i - 1] > start:
            return False
        if i < len(self.starts) and self.starts[i] < end:
            return False
        self.starts.insert(i, start)          # O(n) insert; a balanced BST/SortedList makes it O(log n)
        self.ends.insert(i, end)
        return True
```

**Peak overlap after each booking (LeetCode: My Calendar III).** Re-sweep a time → delta map. Fine for few calls; a lazy segment tree scales.

```python
class MyCalendarThree:
    """My Calendar III: max concurrent bookings after each insert (event map sweep)."""

    def __init__(self) -> None:
        self.delta: dict[int, int] = {}

    def book(self, start: int, end: int) -> int:
        self.delta[start] = self.delta.get(start, 0) + 1
        self.delta[end] = self.delta.get(end, 0) - 1
        active = best = 0
        for t in sorted(self.delta):          # O(n log n) per call; fine for ~400 calls
            active += self.delta[t]
            best = max(best, active)
        return best
```

Other combinations and follow-ups:

- **Double but not triple booking (My Calendar II):** store bookings and pairwise overlaps; reject a booking hitting any overlap.
- **Sweep + heap (Skyline):** events by x, max-heap of active heights with lazy deletion; emit a point when the max changes.
- **Interval cover (Video Stitching, Minimum Number of Taps to Open to Water a Garden):** sort by start, greedily extend the farthest reach (jump game, `18-greedy.md`).
- **Weighted interval scheduling:** greedy fails with values. Sort by end; DP with `bisect` for the last compatible interval (`16-dynamic-programming-foundations.md`).
- **Max events attended, one day each from its range:** sort by start, sweep days, min-heap of end days; attend the one expiring soonest.
- **Streaming intervals (Data Stream as Disjoint Intervals):** sorted starts with `bisect`; merge with neighbors on insert.
- "Input doesn't fit in memory?" External sort by start, then one streaming merge with O(1) state.
- "Streaming intervals; current peak?" A heap of end times gives the active count; its running max is the peak.

## Common pitfalls

- Not asking if endpoints are inclusive. `[1,5]` and `[5,8]` overlap for balloons, not meetings.
- Merging without `max(out[-1][1], e)`: contained intervals shrink the block.
- Sorting by start for max non-overlapping. Sort by end.
- Event ties: `(t, +1)` before `(t, -1)` overcounts back-to-back meetings. Use `-1` for ends, or an explicit key.
- Mutating the caller's input. Fine in interviews if you say so.
- Re-sorting sorted input (Insert Interval): passes, but misses the O(n) solution.
- Difference array: off by one on `r + 1`, no extra slot, or mixed 1- and 0-indexing.
- Difference array with coordinates to 10^9. Sort events or coordinate-compress.
- `heapq.heapreplace` pops then pushes; `heappushpop` pushes then pops. They differ when the new item is smaller than the top.

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| Interval merge | Sliding window (`06-sliding-window.md`) | contiguous subarray of one sequence, not a set of ranges |
| Sort-by-end greedy | Weighted interval DP | intervals carry different values |
| Event sweep | Segment tree / Fenwick | online range updates **and** queries interleaved at scale |
| Meeting rooms | Task scheduler with cooldown | identical tasks, gap constraint, no fixed times (`18-greedy.md`, heaps) |
| Merge | Union-Find (`14-union-find-shortest-paths-mst.md`) | "overlap" is transitive across non-adjacent items in a graph, not a line |
| Difference array | Prefix sums (`04-arrays-hashing-prefix-sums.md`) | prefix sums answer range *queries*; difference arrays do range *updates*; they're inverses |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Merge Intervals (56)](https://leetcode.com/problems/merge-intervals/) · Medium · Sort by start + merge: `max` on the end when merging.
2. [Insert Interval (57)](https://leetcode.com/problems/insert-interval/) · Medium · Merge with a new interval: three phases, no sort needed.
3. [Non-overlapping Intervals (435)](https://leetcode.com/problems/non-overlapping-intervals/) · Medium · Sort by end (greedy): answer is n minus max kept.
4. [Meeting Rooms II (253)](https://leetcode.com/problems/meeting-rooms-ii/) · Medium · Premium · Max overlap (sweep / heap): event sweep or heap of ends; end before start on ties.
5. [Interval List Intersections (986)](https://leetcode.com/problems/interval-list-intersections/) · Medium · Two pointers on intervals: advance the one that ends first.

### Optional

**Merge / insert / intersect**
- [Meeting Rooms (252)](https://leetcode.com/problems/meeting-rooms/) · Easy · Premium · Sort by start, check neighbors.
- [Remove Covered Intervals (1288)](https://leetcode.com/problems/remove-covered-intervals/) · Medium · Sort start asc, end desc; track max end.
- [Employee Free Time (759)](https://leetcode.com/problems/employee-free-time/) · Hard · Premium · Merge all, emit gaps (or heap k-way merge).
- [Data Stream as Disjoint Intervals (352)](https://leetcode.com/problems/data-stream-as-disjoint-intervals/) · Hard · Sorted structure, merge with neighbors.

**Peak overlap / event sweep**
- [Car Pooling (1094)](https://leetcode.com/problems/car-pooling/) · Medium · Weighted events or difference array.
- [Corporate Flight Bookings (1109)](https://leetcode.com/problems/corporate-flight-bookings/) · Medium · Difference array, inclusive end.
- [Range Addition (370)](https://leetcode.com/problems/range-addition/) · Medium · Premium · Pure difference array.
- [My Calendar III (732)](https://leetcode.com/problems/my-calendar-iii/) · Hard · Event map; segment tree at scale.
- [The Skyline Problem (218)](https://leetcode.com/problems/the-skyline-problem/) · Hard · Sweep + max-heap with lazy deletion.

**Sort by end greedy / cover**
- [Minimum Number of Arrows to Burst Balloons (452)](https://leetcode.com/problems/minimum-number-of-arrows-to-burst-balloons/) · Medium · Closed intervals, strict `>`.
- [Video Stitching (1024)](https://leetcode.com/problems/video-stitching/) · Medium · Sort by start, extend farthest reach.
- [Maximum Number of Events That Can Be Attended (1353)](https://leetcode.com/problems/maximum-number-of-events-that-can-be-attended/) · Medium · Day sweep + min-heap of ends.

**Online booking**
- [My Calendar I (729)](https://leetcode.com/problems/my-calendar-i/) · Medium · `bisect` to find the two neighbors.
- [My Calendar II (731)](https://leetcode.com/problems/my-calendar-ii/) · Medium · Track double-booked overlaps.

## Self-check

1. Why does sorting by start mean a new interval can only overlap the last merged block?
2. Give three intervals where "sort by start, keep first" fails for max non-overlapping.
3. In a meeting-rooms sweep, does a start or an end come first at equal time, and why? What changes for closed intervals?
4. Why can Meeting Rooms II sort starts and ends independently?
5. Minimum arrows uses `s > last`; non-overlapping intervals uses `s >= last_end`. Explain via interval closure.
6. When would you choose a difference array over sorting events, and vice versa?
7. In interval intersection, why is discarding the interval that ends first safe?
