# 08. Stacks and Monotonic Structures

Use a stack when the newest unresolved item resolves first: nested brackets, expressions, undo, explicit recursion. The **monotonic stack** finds every element's next/previous greater or smaller in O(n), powering histogram rectangles, stock spans, and contribution counting.

## At a glance

- **Use when:** nesting/matching ("valid", `k[...]`, calculators), or each element's nearest larger/smaller one (waits, spans, rectangles, subarray-min sums).
- **Templates:**
  - Template 1: matching / nesting: brackets, parsing, collisions
  - Template 2: monotonic stack "next": the popped element gets the answer (deep dive: `08a-monotonic-stacks.md`)
  - Template 3: monotonic stack "previous": the survivor answers the new element
  - Template 4: monotonic deque: window elements expire
  - Template 5: augmented stack: O(1) min/max with pop
  - Template 6: queue from two stacks: FIFO from LIFO
- **Typical complexity:** O(n) time (each index pushed and popped once), O(n) space.
- **Top trap:** wrong strictness on duplicates (double-counting contributions); no sentinel leaves trailing elements unprocessed.

## Core idea

**Plain stack (LIFO):** push waiting items; pop when their match arrives. Nesting guarantees the newest waiter resolves first.

**Monotonic stack:** indices whose values stay sorted (increasing or decreasing) bottom to top. A new element pops every element it "beats"; it is each popped element's *next* greater/smaller. The remaining top is the new element's *previous* greater/smaller.

```
 next greater, decreasing stack (values shown), process 2 1 2 4 3

 i=0 x=2  stack: [2]
 i=1 x=1  stack: [2 1]
 i=2 x=2  pop 1 (ans=2)          stack: [2 2]
 i=3 x=4  pop 2, pop 2 (ans=4)   stack: [4]
 i=4 x=3  stack: [4 3]           leftovers have no next greater -> -1
```

**Invariant:** the stack holds exactly the elements not yet beaten by anything to their right. A popped element is never a future answer: its popper is closer and at least as good.

**Why O(n):** each index is pushed once and popped at most once, so the inner `while` runs at most n times total.

### Direction and strictness cheat sheet

| Goal | Stack order (bottom → top) | Pop while | Answer assigned |
|---|---|---|---|
| Next greater | decreasing | `nums[top] < x` | to popped, = i |
| Next greater or equal | strictly decreasing | `nums[top] <= x` | to popped, = i |
| Next smaller | increasing | `nums[top] > x` | to popped, = i |
| Next smaller or equal | strictly increasing | `nums[top] >= x` | to popped, = i |
| Previous greater | decreasing | `nums[top] <= x` | to x, = top after popping |
| Previous smaller | increasing | `nums[top] >= x` | to x, = top after popping |

Mnemonic: next **greater** needs a **decreasing** stack (the opposite of what you seek). One pass gives both directions: when `x` pops `j`, `x` is j's next and the element under `j` is j's previous.

**Duplicates:**
- Pop on `<` (strict): equals don't pop, so "next greater" skips them.
- Pop on `<=` (non-strict): equals pop each other, so "next greater or equal" finds the next twin.
- Contribution/counting: **strict on one side, non-strict on the other**, so a subarray with tied minimums is credited once (Sum of Subarray Minimums).

## Recognition signals

- Brackets, tags, "valid", "balanced", "well-formed": matching stack.
- Calculators, `k[...]` encodings, path simplification (`/a/./b/../c`): parsing stack.
- "Next/previous larger/smaller/warmer/taller" per element: monotonic stack.
- "Days until...", "how far can you see", "span of days with price ≤ today": next/previous greater.
- "Largest rectangle", "maximal area", "water trapped between bars": nearest smaller/greater on both sides.
- "Sum over all subarrays of min/max": contribution technique.
- "Remove k digits for the smallest number", "smallest subsequence of distinct characters": greedy monotonic stack.
- "Max/min of each sliding window": monotonic deque (`06-sliding-window.md`).
- Adjacent survivors colliding or cancelling (asteroids, "remove adjacent duplicates"): simulation stack.
- n up to 10^5 and an obvious O(n^2) "scan right for a bigger one" brute force.

## Disguises

1. **Alert backlog.** "You record CPU load each minute. For each minute, how many minutes until load is strictly higher? 0 if never."
   Decode: next strictly greater, as a distance (Daily Temperatures).

2. **Skyline billboard.** "Buildings of given heights stand in a row. Find the largest rectangular billboard painted across contiguous faces, as tall as the shortest building it covers."
   Decode: min height × width over a contiguous range = Largest Rectangle in Histogram.

3. **Config templating.** "A config DSL writes repeated blocks as `2[host-a,3[x]]`, nested arbitrarily. Expand it."
   Decode: nested `k[...]`: save outer context on `[`, splice on `]` (Decode String).

4. **Cost of cheapest node per window.** "For every contiguous range of servers in a rack, take the cheapest price. Sum over all ranges."
   Decode: sum of subarray mins = contribution via previous/next smaller (Sum of Subarray Minimums).

5. **Message cancellation.** "Messages move along a bus at equal speed: positive ids right, negative left. When two meet, the smaller absolute id is dropped (both on a tie). Which survive?"
   Decode: only right-then-left movers meet, latest right-mover first: stack simulation (Asteroid Collision).

6. **Shortest version string.** "Remove k digits from a version number to minimize it."
   Decode: greedy increasing stack; pop a larger previous digit while removals remain (Remove K Digits).

## Templates and worked examples

All templates are in `code/08_stacks_monotonic.py`.

### Template 1: matching / nesting

> - **Example problem:** [Valid Parentheses (20)](https://leetcode.com/problems/valid-parentheses/)
> - **Why this template:** brackets closing in nesting order: the most recent unmatched opener must match next.
> - **Worked below:** Example 1.1

```python
def is_balanced(s: str) -> bool:
    pairs = {")": "(", "]": "[", "}": "{"}
    stack: list[str] = []
    for ch in s:
        if ch in pairs:
            if not stack or stack[-1] != pairs[ch]:
                return False
            stack.pop()
        elif ch in "([{":
            stack.append(ch)
    return not stack
```

Knobs: what you push (char, index for removals, context for parsing), what "match" means, and what's left (unmatched openers).

#### Example 1.1: Balanced brackets [Easy]

**Problem:** [Valid Parentheses (20)](https://leetcode.com/problems/valid-parentheses/) · **Template:** Template 1 (matching / nesting)

**Problem.** A string contains only `()[]{}`. Decide whether every opener is closed by the same type, in correct nesting order.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** nesting + "valid" = matching stack.

**Approach:** push openers; each closer must match the top; end empty.

```python
def is_valid(s: str) -> bool:
    pairs = {")": "(", "]": "[", "}": "{"}
    stack: list[str] = []
    for ch in s:
        if ch in pairs:
            if not stack or stack.pop() != pairs[ch]:
                return False
        else:
            stack.append(ch)
    return not stack
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** closer on an empty stack (`"]"`); leftover openers (`"(("`). Per-type counts fail: `"([)]"` balances but nests badly.

Related repairs (code file):
- `min_add_to_make_valid` (Minimum Add to Make Parentheses Valid): one type, so a counter suffices; answer = unmatched `)` + unmatched `(`.
- `min_remove_to_make_valid` (Minimum Remove to Make Valid Parentheses): push **indices** of `(`; blank unmatched `)` immediately, leftover `(` at the end.

</details>

#### Example 1.2: Expand nested repeats [Medium]

**Problem:** [Decode String (394)](https://leetcode.com/problems/decode-string/) · **Template:** Template 1 (matching / nesting)

**Problem.** Expand an encoding where `k[s]` repeats `s` `k` times, nested arbitrarily: `3[a2[c]]` → `accaccacc`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** inner groups finish first = stack of saved contexts.

**Approach:** track the current level's string (`cur`) and pending count (`num`). On `[`, push `(cur, num)` and reset. On `]`, pop `(prev, k)`; `cur = prev + cur * k`.

```python
def decode_string(s: str) -> str:
    stack: list[tuple[str, int]] = []
    cur: list[str] = []
    num = 0
    for ch in s:
        if ch.isdigit():
            num = num * 10 + int(ch)         # multi-digit counts
        elif ch == "[":
            stack.append(("".join(cur), num))
            cur, num = [], 0
        elif ch == "]":
            prev, k = stack.pop()
            cur = [prev + "".join(cur) * k]
        else:
            cur.append(ch)
    return "".join(cur)
```

**Complexity:** O(output length) time and space; worst case O(output × depth), since each nesting level rebuilds its characters.

**Edge cases:** multi-digit counts (`10[a]`); letters outside brackets; reset `num` after pushing. Recursive descent (parse to `]`, return string and index) also works.

Related (code file): `eval_rpn` (Evaluate Reverse Polish Notation): push numbers; on an operator pop `b` then `a`, push `a op b`; truncate division with `int(a / b)`.

</details>

#### Example 1.3: Asteroid collision (simulation) [Medium]

**Problem:** [Asteroid Collision (735)](https://leetcode.com/problems/asteroid-collision/) · **Template:** Template 1 (matching / nesting), simulation variant

**Problem.** Asteroids move along a line: absolute value is size, sign is direction (positive = right, negative = left). When a right-mover meets a left-mover, the smaller explodes; equal sizes both explode. Return the survivors in order.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** only a stacked right-mover and a new left-mover collide, most recent right-mover first.

**Approach:** loop while the new asteroid is alive and the top moves right.

```python
def asteroid_collision(asteroids: list[int]) -> list[int]:
    stack: list[int] = []
    for a in asteroids:
        alive = True
        while alive and a < 0 and stack and stack[-1] > 0:
            if stack[-1] < -a:
                stack.pop()          # top explodes, keep checking
            elif stack[-1] == -a:
                stack.pop()          # both explode
                alive = False
            else:
                alive = False        # new one explodes
        if alive:
            stack.append(a)
    return stack
```

**Complexity:** O(n) amortized.

**Edge cases:** `[-2, 2]` move apart and never collide; stacked left-movers never collide with later asteroids.

</details>

#### Example 1.4: Calculator with + - and parentheses [Hard]

**Problem:** [Basic Calculator (224)](https://leetcode.com/problems/basic-calculator/) · **Template:** Template 1 (matching / nesting)

**Problem.** Evaluate a string with non-negative integers, `+`, `-`, parentheses, spaces, and unary minus (e.g., `-(2+3)`). No `eval`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** nested parentheses = stack of saved contexts. With only `+` and `-`, the expression is a signed sum; each number needs only its effective sign.

**Approach:** track `result` (this level's sum), `sign` (for the next number), and `num` (digits read). On `+`/`-`, flush `num` and set `sign`. On `(`, push `result` and `sign` (the group's sign), then reset. On `)`, flush; `result = saved_result + saved_sign * result`. Unary minus is free: `-(` sets `sign = -1` with `result = 0`.

```python
def calculate(s: str) -> int:
    stack: list[int] = []
    result, sign, num = 0, 1, 0
    for ch in s:
        if ch.isdigit():
            num = num * 10 + int(ch)
        elif ch in "+-":
            result += sign * num
            num = 0
            sign = 1 if ch == "+" else -1
        elif ch == "(":
            stack.append(result)
            stack.append(sign)               # sign applied to the whole group
            result, sign = 0, 1
        elif ch == ")":
            result += sign * num
            num = 0
            group_sign = stack.pop()
            result = stack.pop() + group_sign * result
    return result + sign * num
```

**Complexity:** O(n) time, O(depth) space.

**Edge cases:** flush the last number after the loop; spaces are skipped implicitly; pop in reverse push order (sign on top).

Follow-ups: `*` and `/` without parentheses (Basic Calculator II, `calculate_ii`): push signed terms, apply `*`/`/` to the top immediately, sum the stack. All four operators plus parentheses (Basic Calculator III): recurse on `(`, or two-stack shunting-yard.

</details>

### Template 2: monotonic stack, "next" form

> - **Example problem:** [Daily Temperatures (739)](https://leetcode.com/problems/daily-temperatures/)
> - **Why this template:** "days until a warmer temperature" is each element's next strictly greater element.
> - **Worked below:** Example 2.1

```python
def next_greater_index(nums: list[int]) -> list[int]:
    res = [-1] * len(nums)
    stack: list[int] = []
    for i, x in enumerate(nums):
        while stack and nums[stack[-1]] < x:   # TODO: comparison picks greater/smaller, strict/non-strict
            res[stack.pop()] = i               # x is the first element to beat the popped one
        stack.append(i)
    return res
```

**Deep dive:** [`08a-monotonic-stacks.md`](08a-monotonic-stacks.md) covers monotonic stacks step by step (who answers whom, the four variants, duplicates, sentinels, the contribution technique), with seven worked examples.

#### Example 2.1: Days until warmer [Medium]

**Problems:** [Daily Temperatures (739)](https://leetcode.com/problems/daily-temperatures/), [Online Stock Span (901)](https://leetcode.com/problems/online-stock-span/) · **Template:** Template 2 (monotonic stack, "next" form); Template 3 (monotonic stack, "previous" form) for the stock span variant

**Problem.** For each day, return how many days until a strictly warmer temperature, or 0.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "wait until strictly higher" = next strictly greater, as a distance.

**Approach:** Template 2 with a decreasing stack of indices; when today beats top `j`, j's wait is `i - j`.

```python
def daily_temperatures(temps: list[int]) -> list[int]:
    res = [0] * len(temps)
    stack: list[int] = []                    # indices, temps non-increasing
    for i, t in enumerate(temps):
        while stack and temps[stack[-1]] < t:
            j = stack.pop()
            res[j] = i - j
        stack.append(i)
    return res
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** strict `<`, since equal isn't warmer; default 0, not -1.

Previous-greater variant (Online Stock Span): span = consecutive days ending today with price ≤ today. Keep `(price, span)` pairs with strictly decreasing prices, merging popped spans; this works online.

```python
class StockSpanner:
    def __init__(self) -> None:
        self.stack: list[tuple[int, int]] = []

    def next(self, price: int) -> int:
        span = 1
        while self.stack and self.stack[-1][0] <= price:
            span += self.stack.pop()[1]
        self.stack.append((price, span))
        return span
```

</details>

### Template 3: monotonic stack, "previous" form

> - **Example problem:** [Online Stock Span (901)](https://leetcode.com/problems/online-stock-span/)
> - **Why this template:** "consecutive days ending today with price at most today's" is the distance back to the previous strictly greater price.
> - **Worked in:** Example 2.1 (the Online Stock Span variant) above, under Template 2.
> - **Also fits:** [Largest Rectangle in Histogram (84)](https://leetcode.com/problems/largest-rectangle-in-histogram/)

```python
def previous_smaller_index(nums: list[int]) -> list[int]:
    res = [-1] * len(nums)
    stack: list[int] = []
    for i, x in enumerate(nums):
        while stack and nums[stack[-1]] >= x:  # pop non-smaller; survivors are strictly smaller
            stack.pop()
        res[i] = stack[-1] if stack else -1
        stack.append(i)
    return res
```

Knobs for Templates 2 and 3:
1. **Comparison** (cheat sheet above): greater vs smaller, strict vs non-strict.
2. **What to record**: index (distances/widths), value, or merged aggregate (span counts).
3. **Leftovers** have no next greater/smaller: default them (`-1`, `n`, `0`) or append a **sentinel** (`0` or `-inf`) to flush them.
4. **Circular arrays** (Next Greater Element II): iterate `range(2n)` with `nums[i % n]`, pushing only in the first pass.
5. **Right-to-left** iteration turns "next" into "previous"; use either.

### Template 4: monotonic deque (window max/min)

> - **Example problem:** [Sliding Window Maximum (239)](https://leetcode.com/problems/sliding-window-maximum/)
> - **Why this template:** "maximum of each size-k window" needs a monotonic structure whose oldest entries expire from the front.
> - **Worked below:** Example 4.1

A monotonic stack that also expires from the front (`06-sliding-window.md`).

```python
def sliding_window_max(nums: list[int], k: int) -> list[int]:
    dq: deque[int] = deque()
    out = []
    for i, x in enumerate(nums):
        while dq and nums[dq[-1]] <= x:
            dq.pop()
        dq.append(i)
        if dq[0] <= i - k:
            dq.popleft()
        if i >= k - 1:
            out.append(nums[dq[0]])
    return out
```

#### Example 4.1: Shortest subarray with sum at least k, negatives allowed [Hard]

**Problem:** [Shortest Subarray with Sum at Least K (862)](https://leetcode.com/problems/shortest-subarray-with-sum-at-least-k/) · **Template:** Template 4 (monotonic deque (window max/min))

**Problem.** Given integers (possibly negative) and `k`, return the length of the shortest non-empty contiguous subarray with sum ≥ `k`, or -1. Example: `[2, -1, 2]`, `k = 3` → 3.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "shortest subarray, sum ≥ k" suggests a sliding window, but negatives break it: shrinking from the left can *raise* the sum (`06-sliding-window.md`). Rewrite as: for each end `i`, find the latest start `j < i` with `pre[i] - pre[j] >= k`. Many candidate starts, some dominated, some expiring: monotonic deque of prefix indices.

**Approach:** keep prefix indices with strictly increasing `pre` values. For each `i`:
- **Front:** while `pre[i] - pre[front] >= k`, record `i - front` and pop. Any later end would only make this start longer, so it's finished.
- **Back:** while `pre[back] >= pre[i]`, pop. `i` is later (shorter subarrays) and no larger (sum at least as big), so `back` can never beat it.

Each index is pushed and popped once.

```python
def shortest_subarray(nums: list[int], k: int) -> int:
    pre = [0]
    for x in nums:
        pre.append(pre[-1] + x)
    best = len(nums) + 1
    dq: deque[int] = deque()                  # prefix indices, pre values strictly increasing
    for i, p in enumerate(pre):
        while dq and p - pre[dq[0]] >= k:
            best = min(best, i - dq.popleft())    # later ends only make this start longer
        while dq and pre[dq[-1]] >= p:
            dq.pop()                          # i is later and no larger: a better start
        dq.append(i)
    return best if best <= len(nums) else -1
```

**Complexity:** O(n) time (each prefix index enters and leaves the deque once), O(n) space.

**Edge cases:** iterate over `pre` (length n+1) so subarrays starting at index 0 count; nothing qualifies → -1; all negatives (`[-5, -2]`, `k = 1`) → -1. With only non-negative values, the plain shrinking window suffices.

</details>

### Template 5: augmented stack (min-stack)

> - **Example problem:** [Min Stack (155)](https://leetcode.com/problems/min-stack/)
> - **Why this template:** O(1) minimum alongside push/pop: store the running min with each entry.
> - **Worked below:** Example 5.1

Store the aggregate with each element; pop restores it automatically.

```python
class MinStack:
    def __init__(self) -> None:
        self.stack: list[tuple[int, int]] = []

    def push(self, val: int) -> None:
        cur_min = min(val, self.stack[-1][1]) if self.stack else val
        self.stack.append((val, cur_min))

    def pop(self) -> None:
        self.stack.pop()

    def top(self) -> int:
        return self.stack[-1][0]

    def get_min(self) -> int:
        return self.stack[-1][1]
```

Works for any aggregate that never needs "un-applying" (min, max, gcd). Two of these as a queue give a sliding-window aggregate for any associative operation.

#### Example 5.1: Stack with constant-time minimum [Medium]

**Problem:** [Min Stack (155)](https://leetcode.com/problems/min-stack/) · **Template:** Template 5 (augmented stack (min-stack))

**Problem.** Design a stack supporting `push(val)`, `pop()`, `top()`, and `get_min()` (current smallest value), each O(1). `pop`, `top`, and `get_min` are only called on a non-empty stack. Example: push −2, push 0, push −3, `get_min()` → −3, `pop()`, `top()` → 0, `get_min()` → −2.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a single `min` variable can't recover the previous minimum on pop, and a heap is O(log n). Removals are LIFO, so the minimum "as of" each entry never changes: store it with the entry.

**Approach:** push `(val, min(val, previous min))`. Invariant: `stack[i][1]` is the minimum of `stack[0..i]`, so after a pop the new top already holds the correct minimum.

```python
class MinStack:
    def __init__(self) -> None:
        self.stack: list[tuple[int, int]] = []

    def push(self, val: int) -> None:
        cur_min = min(val, self.stack[-1][1]) if self.stack else val
        self.stack.append((val, cur_min))

    def pop(self) -> None:
        self.stack.pop()

    def top(self) -> int:
        return self.stack[-1][0]

    def get_min(self) -> int:
        return self.stack[-1][1]
```

**Complexity:** O(1) per operation; O(n) space (one extra integer per element).

**Edge cases:** duplicate minimums (push 0, 0, pop) work since each entry stores its own min. A separate stack of minimums must push on `val <= current min`, not `<`, or popping one duplicate loses the min. LeetCode names it `getMin`; ask whether empty-stack calls need handling (here: `IndexError`).

</details>

### Template 6: queue from two stacks

> - **Example problem:** [Implement Queue using Stacks (232)](https://leetcode.com/problems/implement-queue-using-stacks/)
> - **Why this template:** a FIFO queue from two stacks is exactly the inbox/outbox shift.
> - **Worked below:** Example 6.1

```python
class MyQueue:
    def __init__(self) -> None:
        self.inbox: list[int] = []
        self.outbox: list[int] = []

    def push(self, x: int) -> None:
        self.inbox.append(x)

    def _shift(self) -> None:
        if not self.outbox:
            while self.inbox:
                self.outbox.append(self.inbox.pop())

    def pop(self) -> int:
        self._shift()
        return self.outbox.pop()

    def peek(self) -> int:
        self._shift()
        return self.outbox[-1]

    def empty(self) -> bool:
        return not self.inbox and not self.outbox
```

Amortized O(1) per operation (why: Example 6.1). In real Python, use `collections.deque`.

#### Example 6.1: FIFO queue from two stacks [Easy]

**Problem:** [Implement Queue using Stacks (232)](https://leetcode.com/problems/implement-queue-using-stacks/) · **Template:** Template 6 (queue from two stacks)

**Problem.** Build a FIFO queue with `push(x)`, `pop()`, `peek()`, and `empty()`, using only stack operations (push to top, pop/peek top, size, is-empty). `pop` and `peek` are only called on a non-empty queue. Example: push 1, push 2, `peek()` → 1, `pop()` → 1, `empty()` → false.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** pouring one stack into another reverses it, turning LIFO into FIFO.

**Approach:** `push` goes to `inbox`. `pop`/`peek` read `outbox`; only when `outbox` is empty, pour all of `inbox` into it. Invariant: `outbox` top-to-bottom, then `inbox` bottom-to-top, is the queue front-to-back. Pouring into a non-empty `outbox` would bury older elements under newer ones.

```python
class MyQueue:
    def __init__(self) -> None:
        self.inbox: list[int] = []
        self.outbox: list[int] = []

    def push(self, x: int) -> None:
        self.inbox.append(x)

    def _shift(self) -> None:
        if not self.outbox:
            while self.inbox:
                self.outbox.append(self.inbox.pop())

    def pop(self) -> int:
        self._shift()
        return self.outbox.pop()

    def peek(self) -> int:
        self._shift()
        return self.outbox[-1]

    def empty(self) -> bool:
        return not self.inbox and not self.outbox
```

**Complexity:** amortized O(1) per operation: each element is pushed to `inbox`, moved to `outbox` at most once, and popped once, so n operations cost O(n) total. A single `pop` can be O(n) (a full shift). O(n) space.

**Edge cases:** interleaved pushes and pops (pushes land in `inbox` while `outbox` still holds older items); `empty()` checks both stacks. LeetCode guarantees valid calls; otherwise `pop` on empty raises `IndexError`.

</details>

### Combined examples

#### Example C.1: Sum of minimums over all subarrays [Medium]

**Problem:** [Sum of Subarray Minimums (907)](https://leetcode.com/problems/sum-of-subarray-minimums/) · **Template:** Template 2 (monotonic stack, "next" form) and Template 3 (monotonic stack, "previous" form)

**Problem.** Sum `min(sub)` over every contiguous subarray, modulo 10^9+7.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** n ~3·10^4+ rules out O(n^2) enumeration. Flip the sum: count the subarrays where each element is the minimum (**contribution technique**).

**Approach:** `arr[j]` is the min of the `(j - left) * (right - j)` subarrays starting after `left` (previous smaller) and ending before `right` (next smaller). To avoid counting `[2, 2]` twice, use `left` = previous **strictly** smaller and `right` = next smaller **or equal**, crediting each subarray to its rightmost tied minimum. Popping while `arr[top] >= cur` gives both in one pass: when `i` pops `j`, `i` is j's right and the element under `j` is j's left. A `-inf` sentinel flushes the stack.

```
 arr = [3, 1, 2, 4]     contribution of 1 (index 1): left=-1, right=4
                        (1 - (-1)) * (4 - 1) = 6 subarrays, sum += 1 * 6
```

```python
def sum_subarray_mins(arr: list[int], mod: int = 10**9 + 7) -> int:
    total = 0
    stack: list[int] = []                    # indices, values strictly increasing
    for i in range(len(arr) + 1):
        cur = arr[i] if i < len(arr) else float("-inf")   # sentinel flushes the stack
        while stack and arr[stack[-1]] >= cur:
            j = stack.pop()
            left = stack[-1] if stack else -1
            total += arr[j] * (j - left) * (i - j)
        stack.append(i)
    return total % mod
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** duplicates (the strict/non-strict asymmetry). Python ints don't overflow, so mod at the end (per step in other languages). Follow-up: sum of subarray *ranges* = sum of maxes − sum of mins, two passes.

</details>

#### Example C.2: Largest rectangle in a histogram [Hard]

**Problems:** [Largest Rectangle in Histogram (84)](https://leetcode.com/problems/largest-rectangle-in-histogram/), [Maximal Rectangle (85)](https://leetcode.com/problems/maximal-rectangle/) · **Template:** Template 2 (monotonic stack, "next" form) and Template 3 (monotonic stack, "previous" form)

**Problem.** Given bar heights of width 1, find the largest rectangle area in the histogram.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** the shortest bar limits the area: treat each bar as the limiting height and extend to its nearest shorter bars on both sides.

**Approach:** increasing stack of indices. When bar `i` is ≤ top `j`, pop `j`: right boundary `i`, left boundary the new top (previous strictly shorter), width `i - left - 1`. A height-0 sentinel flushes the rest. Popping equal heights on `>=` undercounts the earlier bar's width, but the last tied bar gets the full width, so the max is right.

```
 heights: 2 1 5 6 2 3
 at i=4 (h=2): pop 6 -> width 4-2-1=1, area 6
               pop 5 -> width 4-1-1=2, area 10   <- best
```

```python
def largest_rectangle_area(heights: list[int]) -> int:
    best = 0
    stack: list[int] = []
    for i in range(len(heights) + 1):
        h = heights[i] if i < len(heights) else 0   # sentinel 0 flushes the stack
        while stack and heights[stack[-1]] >= h:
            j = stack.pop()
            left = stack[-1] if stack else -1
            best = max(best, heights[j] * (i - left - 1))
        stack.append(i)
    return best
```

**Complexity:** O(n) time, O(n) space.

**Edge cases:** without the sentinel, a trailing increasing run is never evaluated (`[2, 4]`). An empty stack after popping means the bar reaches the left edge (`left = -1`).

Extension (Maximal Rectangle): each matrix row is a histogram base where `heights[c]` counts consecutive `1`s ending there; run the histogram per row. O(rows × cols).

```python
def maximal_rectangle(matrix: list[list[str]]) -> int:
    if not matrix or not matrix[0]:
        return 0
    heights = [0] * len(matrix[0])
    best = 0
    for row in matrix:
        for c, v in enumerate(row):
            heights[c] = heights[c] + 1 if v == "1" else 0
        best = max(best, largest_rectangle_area(heights))
    return best
```

</details>

## Template cheat sheet

All templates, for review.

### Template 1: matching / nesting
```python
def is_balanced(s: str) -> bool:
    pairs = {")": "(", "]": "[", "}": "{"}
    stack: list[str] = []
    for ch in s:
        if ch in pairs:
            if not stack or stack[-1] != pairs[ch]:
                return False
            stack.pop()
        elif ch in "([{":
            stack.append(ch)
    return not stack
```
### Template 2: monotonic stack, "next" form
```python
def next_greater_index(nums: list[int]) -> list[int]:
    res = [-1] * len(nums)
    stack: list[int] = []
    for i, x in enumerate(nums):
        while stack and nums[stack[-1]] < x:   # TODO: comparison picks greater/smaller, strict/non-strict
            res[stack.pop()] = i               # x is the first element to beat the popped one
        stack.append(i)
    return res
```
### Template 3: monotonic stack, "previous" form
```python
def previous_smaller_index(nums: list[int]) -> list[int]:
    res = [-1] * len(nums)
    stack: list[int] = []
    for i, x in enumerate(nums):
        while stack and nums[stack[-1]] >= x:  # pop non-smaller; survivors are strictly smaller
            stack.pop()
        res[i] = stack[-1] if stack else -1
        stack.append(i)
    return res
```
### Template 4: monotonic deque (window max/min)
```python
def sliding_window_max(nums: list[int], k: int) -> list[int]:
    dq: deque[int] = deque()
    out = []
    for i, x in enumerate(nums):
        while dq and nums[dq[-1]] <= x:
            dq.pop()
        dq.append(i)
        if dq[0] <= i - k:
            dq.popleft()
        if i >= k - 1:
            out.append(nums[dq[0]])
    return out
```
### Template 5: augmented stack (min-stack)
```python
class MinStack:
    def __init__(self) -> None:
        self.stack: list[tuple[int, int]] = []

    def push(self, val: int) -> None:
        cur_min = min(val, self.stack[-1][1]) if self.stack else val
        self.stack.append((val, cur_min))

    def pop(self) -> None:
        self.stack.pop()

    def top(self) -> int:
        return self.stack[-1][0]

    def get_min(self) -> int:
        return self.stack[-1][1]
```
### Template 6: queue from two stacks
```python
class MyQueue:
    def __init__(self) -> None:
        self.inbox: list[int] = []
        self.outbox: list[int] = []

    def push(self, x: int) -> None:
        self.inbox.append(x)

    def _shift(self) -> None:
        if not self.outbox:
            while self.inbox:
                self.outbox.append(self.inbox.pop())

    def pop(self) -> int:
        self._shift()
        return self.outbox.pop()

    def peek(self) -> int:
        self._shift()
        return self.outbox[-1]

    def empty(self) -> bool:
        return not self.inbox and not self.outbox
```

## Variations and extensions

- **Greedy monotonic stack** (Remove K Digits, Remove Duplicate Letters, Smallest Subsequence of Distinct Characters): pop the top while it's worse than the incoming element *and* droppable (removals left, or it reappears).
- **Trapping Rain Water**: the stack version pops a bottom and fills water bounded by the new top and incoming bar; two pointers use O(1) space (`05-two-pointers.md`).
- **Next greater in another array** (Next Greater Element I): compute once for the reference array; store in a dict.
- **Count visible people / buildings with a view**: count pops, plus one if the stack is non-empty.
- **Contribution generalizations**: sum of subarray ranges, subarrays where an element is the max, "maximum min-product" (prefix sums + previous/next smaller).
- **Monotonic deque + DP** (Jump Game VI, Constrained Subsequence Sum): each state takes the max over a sliding window of earlier states.
- **Stack for tree/graph**: iterative DFS and in-order traversal (`10-trees-and-bst.md`), validate preorder serialization, "132 pattern" (decreasing stack from the right).
- **Follow-ups**: streaming → monotonic stacks are naturally online for "previous" queries; O(1) extra space → usually impossible for next-greater (trapping water is the exception); min-stack with O(1) extra total → store differences from the min.

## Common pitfalls

- Storing values instead of indices, then needing distances or widths.
- Wrong strictness: double-counting contributions or wrong answers on duplicates (equal temperatures).
- No sentinel or leftover handling, so trailing elements go unprocessed.
- Mixing up who gets the answer: the popped element in "next", the new element in "previous".
- RPN / calculator operand order: `b = pop(); a = pop(); a - b`.
- Python `//` floors toward negative infinity; problems usually want truncation toward zero (`int(a / b)`).
- Two-stack queue: shifting into a non-empty outbox breaks FIFO order.

## When NOT to use / look-alikes

| Looks like | Actually | Tell |
|---|---|---|
| Next greater element | Sorting / BST | "Next greater *value* anywhere" (not by position): sort, or balanced BST / `bisect` |
| Window max with stack | Monotonic deque | Elements expire from the old end: need popleft (`06-sliding-window.md`) |
| Bracket validity | Counter | One bracket type, validity only: a counter suffices |
| Remove invalid parentheses (all answers) | BFS / backtracking | Must enumerate all minimal results |
| Largest rectangle | Two pointers | Container With Most Water: only two lines matter, no "min over range" |
| Longest valid parentheses | Stack of indices or DP | Keep the last unmatched index as the base; plain matching isn't enough |
| k-th largest / top-k | Heap | Global value order, not nearest by position (`11-heaps-top-k-and-merging.md`) |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Valid Parentheses (20)](https://leetcode.com/problems/valid-parentheses/) · Easy · Matching: push openers, pop on matching closer.
2. [Min Stack (155)](https://leetcode.com/problems/min-stack/) · Medium · Augmented stack: store running min per entry.
3. [Daily Temperatures (739)](https://leetcode.com/problems/daily-temperatures/) · Medium · Next greater (monotonic stack): next greater distance.
4. [Decode String (394)](https://leetcode.com/problems/decode-string/) · Medium · Parsing with a stack: push context on `[`.
5. [Largest Rectangle in Histogram (84)](https://leetcode.com/problems/largest-rectangle-in-histogram/) · Hard · Previous/next smaller: sentinel, width from new top.

### Optional

**Matching / parsing**
- [Minimum Add to Make Parentheses Valid (921)](https://leetcode.com/problems/minimum-add-to-make-parentheses-valid/) · Medium · Counter.
- [Minimum Remove to Make Valid Parentheses (1249)](https://leetcode.com/problems/minimum-remove-to-make-valid-parentheses/) · Medium · Stack of indices.
- [Evaluate Reverse Polish Notation (150)](https://leetcode.com/problems/evaluate-reverse-polish-notation/) · Medium · Operand order, truncation.
- [Simplify Path (71)](https://leetcode.com/problems/simplify-path/) · Medium · Split on `/`, stack of directory names.
- [Basic Calculator II (227)](https://leetcode.com/problems/basic-calculator-ii/) · Medium · Precedence via immediate `*` `/`.
- [Basic Calculator (224)](https://leetcode.com/problems/basic-calculator/) · Hard · Push result and sign on `(`.

**Monotonic stack**
- [Online Stock Span (901)](https://leetcode.com/problems/online-stock-span/) · Medium · Previous greater with merged spans.
- [Next Greater Element II (503)](https://leetcode.com/problems/next-greater-element-ii/) · Medium · Circular, iterate 2n.
- [Remove K Digits (402)](https://leetcode.com/problems/remove-k-digits/) · Medium · Greedy increasing stack.
- [Sum of Subarray Minimums (907)](https://leetcode.com/problems/sum-of-subarray-minimums/) · Medium · Contribution, strict/non-strict asymmetry.
- [Maximal Rectangle (85)](https://leetcode.com/problems/maximal-rectangle/) · Hard · Histogram per row.
- [Trapping Rain Water (42)](https://leetcode.com/problems/trapping-rain-water/) · Hard · Stack or two pointers.

**Monotonic deque**
- [Sliding Window Maximum (239)](https://leetcode.com/problems/sliding-window-maximum/) · Hard · Decreasing deque of indices.

**Design / simulation**
- [Implement Queue using Stacks (232)](https://leetcode.com/problems/implement-queue-using-stacks/) · Easy · Amortized shift.
- [Asteroid Collision (735)](https://leetcode.com/problems/asteroid-collision/) · Medium · Collide only right-then-left.

## Self-check

1. Why is a monotonic stack O(n) even though it has a nested `while` loop?
2. For next *smaller*, is the stack increasing or decreasing? What changes for "smaller or equal"?
3. When `x` pops `j`, what two facts do you learn about `j`, and which one about `x`?
4. In Sum of Subarray Minimums, what goes wrong if both sides use strict comparisons? Both non-strict?
5. In Largest Rectangle in Histogram, why is the left boundary the new top after popping, and why is the sentinel needed?
6. In Basic Calculator, what do you push on `(`, and how does unary minus fall out?
7. Why is the two-stack queue amortized O(1), and what breaks if you shift on every pop?
