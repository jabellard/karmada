# 16. Dynamic Programming Foundations

Dynamic programming is recursion that never solves a subproblem twice, turning exponential search over choices into polynomial work over states. It covers most "number of ways", "minimum cost", "maximum value", and "is it possible" questions. This module teaches the **method**; module 17 applies it to harder states.

## At a glance

- **Use when:** counting, optimizing, or deciding feasibility over interacting choices, with small brute-force recursion arguments (an index, a remaining amount, two indices).
- **Templates:**
  - Template 1: 1D linear: one sequence, answer over prefixes.
  - Template 2: 2D grid: right/down moves only.
  - Template 3: Two sequences: compare, align, or transform two strings.
  - Template 4: 0/1 knapsack: each item at most once, modest capacity.
  - Template 5: Unbounded knapsack: reusable items, modest capacity.
- **Typical complexity:** states x work per transition; usually O(n), O(n^2), O(m·n), or O(n·cap).
- **Top trap:** a vague state; write "dp[i] = ..." as a full sentence before coding.

## Core idea

A problem is DP when it has:

- **Optimal substructure**: the best answer is built from best answers to smaller problems of the same kind.
- **Overlapping subproblems**: naive recursion revisits the same smaller problem many times.

```
rob(0)                         naive recursion for House Robber
├── rob(1)
│   ├── rob(2)  <──┐
│   │   ├── rob(3)  │  rob(2) and rob(3) are computed again and again:
│   │   └── rob(4)  │  the tree has ~1.6^n nodes but only n distinct labels
│   └── rob(3)      │
└── rob(2)  ────────┘
```

DP computes each distinct label once. Time = states x work per transition.

**Invariant**: when computing `dp[state]`, every state it depends on is final. Evaluation order exists to keep this true.

### The 7-step method

Say these steps out loud in an interview.

| Step | Question | House Robber answer |
|---|---|---|
| 1. Recognize | Choices + optimize/count/feasible + repeats? | Choose houses, maximize, no two adjacent |
| 2. State in words | "dp[i] = the answer for ..." | dp[i] = max money from the first i houses |
| 3. Transition | What is the **last decision**? | House i-1 robbed or not: `max(dp[i-1], dp[i-2] + nums[i-1])` |
| 4. Base cases | Smallest states with obvious answers | dp[0] = 0, dp[1] = nums[0] |
| 5. Order | Which states must be ready first? | i increasing |
| 6. Answer location | Which cell(s) hold the answer? | dp[n] |
| 7. Space | Which earlier cells does the transition read? | Only i-1, i-2: two scalars |

**Step 2 is where most people fail.** Two common shapes:

- **"up to i"** (prefix): dp[i] = best answer over the first i items. Answer is `dp[n]`.
- **"ending at i"**: dp[i] = best structure **ending exactly at** i, needed when the next step must know where the last piece ended (LIS, max subarray). Answer is `max(dp)`.

**Step 3 trick**: think about the *last* decision. "How did I get here?" gives the predecessors.

### The progression on one problem (House Robber)

Every DP solution has four stages; Example 1.2 walks through all four.

| Stage | Time | Space | When to use |
|---|---|---|---|
| Brute recursion | O(2^n) | O(n) | Only to discover the recurrence |
| Memoization | O(states x transition) | O(states) + recursion stack | Fastest to write; sparse state spaces; awkward orders |
| Bottom-up table | same | O(states) | Deep recursion risk; need the table for reconstruction |
| Rolling | same | O(one or two rows) | Transition reads a fixed window of earlier rows |

**Top-down vs bottom-up**: same complexity. Top-down touches only reachable states and needs no ordering. Bottom-up avoids Python's recursion limit (~1000) and allows rolling space; prefer it if n can be 10^5.

## Recognition signals

- **"Number of ways"**, "how many distinct", "count the paths", often "modulo 10^9 + 7" (too many to enumerate).
- **"Minimum cost / maximum value / longest / shortest"** where greedy isn't obviously safe.
- **"Is it possible / can you reach / can it be split"** with interacting choices.
- Input is a **sequence**, **two sequences** (compare/align/transform), a **grid** with right/down moves, or a modest **capacity/target** (<= ~10^4).
- **Constraints**: n <= 10^3 suggests O(n^2); n <= 10^5 suggests O(n) or O(n log n); target <= 10^4 with n in the 100s suggests knapsack O(n * target).
- A backtracking solution with **small** recursion arguments is a memoization waiting to happen.

## Disguises

1. **Retry budget.** "A client retries a failed request after 1 or 2 ticks. How many distinct retry schedules exactly fill a 30-tick window?"
   **Decode:** sequences of 1s and 2s summing to n = climbing stairs: dp[i] = dp[i-1] + dp[i-2].

2. **Non-adjacent maintenance windows.** "Running batch jobs on a night yields a throughput gain, but two nights in a row overheats the cluster. Maximize total gain."
   **Decode:** "no two adjacent" + maximize = House Robber.

3. **Ambiguous encoding.** "A legacy protocol maps 1..26 to tokens and dropped the separators. Count possible decodings of a digit string."
   **Decode:** count splits whose last chunk is 1 or 2 digits = Decode Ways.

4. **Config diff.** "Fewest single-line insertions, deletions, or replacements turning one config file into another."
   **Decode:** two sequences, transform, min operations = edit distance over lines.

5. **Shard balancing.** "Can these shards be split across two nodes so both store exactly the same bytes?"
   **Decode:** subset summing to total/2 = 0/1 knapsack feasibility (Partition Equal Subset Sum).

6. **Denominations.** "Credits come in packs of 1, 5, 10, 25. How many pack combinations (order ignored) buy exactly 100 credits? What is the fewest packs?"
   **Decode:** unlimited reuse = unbounded knapsack; "order ignored" = items in the outer loop; "fewest" = min instead of count.

## Templates and worked examples

### Loop order for counting: combinations vs permutations

Counting ways to reach a total with reusable parts:

| Outer loop | Inner loop | Counts | Why |
|---|---|---|---|
| items (coins) | amount ascending | **combinations** (order ignored) | All uses of coin c come before any later coin, so each multiset is built in one order |
| amount | items | **permutations** (sequences) | For each total, any item can be last, so `1+2` and `2+1` count separately |

### Template 1: 1D linear

> - **Example problem:** [Decode Ways (91)](https://leetcode.com/problems/decode-ways/)
> - **Why this template:** Counting decodings is counting over prefixes; the last decision (1- or 2-digit chunk) reads only dp[i-1] and dp[i-2].
> - **Worked below:** Example 1.1 (House Robber in Example 1.2)
> - **Also fits:** [House Robber (198)](https://leetcode.com/problems/house-robber/) (solved stage by stage in Example 1.2)

```python
def dp_1d_template(items: list[int]) -> int:
    """1D linear DP over a prefix. Customize: meaning of dp[i], base cases,
    which earlier indices the transition reads, combine op (max/min/sum/or)."""
    n = len(items)
    dp = [0] * (n + 1)           # dp[i] = answer for items[:i]
    dp[0] = 0                    # TODO: base case for the empty prefix
    for i in range(1, n + 1):
        # TODO: transition from the last decision about items[i-1]
        dp[i] = dp[i - 1] + items[i - 1]
    return dp[n]                 # TODO: or max(dp) if the answer can end anywhere
```

Knobs: window read (`i-1, i-2`, or all `j < i` for O(n^2)); combine op (`max`, `min`, `+` count, `or` feasibility); answer `dp[n]` (prefix) or `max(dp)` (ending-at). Size `n + 1` with `dp[0]` = empty prefix avoids most off-by-one cases.

#### Example 1.1: Climbing stairs and decode ways (1D counting) [Easy/Medium]

**Problems:** [Climbing Stairs (70)](https://leetcode.com/problems/climbing-stairs/), [Decode Ways (91)](https://leetcode.com/problems/decode-ways/) · **Template:** Template 1 (1D linear)

**Climbing stairs** (LeetCode: Climbing Stairs): count ways to reach step n taking 1 or 2 steps at a time.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** the last step came from i-1 or i-2: ways[i] = ways[i-1] + ways[i-2], ways[0] = ways[1] = 1. Fibonacci.

```python
def climb_stairs(n: int) -> int:
    """Ways to reach step n with steps of 1 or 2. ways[i] = ways[i-1] + ways[i-2]."""
    a, b = 1, 1                  # ways[0], ways[1]
    for _ in range(n - 1):
        a, b = b, a + b
    return b if n >= 1 else 1
```

</details>

**Decode ways** (LeetCode: Decode Ways): digits map to letters 1..26; count decodings of a digit string.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** count + split a string = 1D DP over prefixes.

**Approach:** dp[i] = decodings of `s[:i]`. Sum over the valid last tokens: `s[i-1]` (not `'0'`) and `s[i-2:i]` (10..26). dp[0] = 1 so the first token counts correctly.

```python
def num_decodings(s: str) -> int:
    """Decode Ways. dp[i] = ways to decode s[:i]. Last decision: the final
    chunk is one digit (1-9) or two digits (10-26)."""
    n = len(s)
    dp = [0] * (n + 1)
    dp[0] = 1                                    # empty prefix: one way
    for i in range(1, n + 1):
        if s[i - 1] != "0":
            dp[i] += dp[i - 1]
        if i >= 2 and 10 <= int(s[i - 2:i]) <= 26:
            dp[i] += dp[i - 2]
    return dp[n]
```

**Complexity:** O(n) time, O(n) space (O(1) rolling).

**Edge cases:** leading `'0'` gives 0. `"10"` is 1 (only "J"). `"06"` is 0 (not 10..26). `"2101"` is 1.

</details>

#### Example 1.2: House robber, brute force to O(1) space [Medium]

**Problem:** [House Robber (198)](https://leetcode.com/problems/house-robber/) · **Template:** Template 1 (1D linear)

Problem: houses along a street each hold a non-negative amount of cash; robbing two adjacent houses trips the alarm. Return the most cash you can take. For `[2, 7, 9, 3, 1]` the answer is 12 (houses 0, 2, 4: 2 + 9 + 1).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** maximize + take/skip per element + "no two adjacent" = 1D prefix DP reading i-1 and i-2. Greedy fails: `[2, 1, 1, 2]` wants houses 0 and 3, which "take the bigger neighbor" misses.

**Approach:** run the 7-step method (table above), then the four stages. Submit stage 2 or 4; stages 1 and 3 get you there.

**Stage 1: brute-force recursion.** At house i, skip (to i+1) or take (to i+2). O(2^n), because `best(i)` is recomputed along many paths.

```python
def rob_brute(nums: list[int]) -> int:
    """Step 1: brute-force recursion straight from the last decision.
    At house i: skip it (go to i+1) or take it (go to i+2). O(2^n)."""
    def best(i: int) -> int:
        if i >= len(nums):
            return 0
        return max(best(i + 1), nums[i] + best(i + 2))
    return best(0)
```

**Stage 2: memoization (top-down).** The argument `i` is the state: n + 1 states, each O(1) once, so O(n).

```python
def rob_memo(nums: list[int]) -> int:
    """Step 2: same recursion, cached. n distinct states x O(1) work = O(n)."""
    @cache
    def best(i: int) -> int:
        if i >= len(nums):
            return 0
        return max(best(i + 1), nums[i] + best(i + 2))
    return best(0)
```

**Stage 3: bottom-up table.** Use the prefix form so the table fills left to right: dp[i] = best from the first i houses, skipping house i-1 (`dp[i-1]`) or robbing it (`dp[i-2] + nums[i-1]`). Each dp[i] reads only smaller, final indices. No recursion limit or cache overhead.

```python
def rob_table(nums: list[int]) -> int:
    """Step 3: bottom-up. dp[i] = best using the first i houses.
    dp[i] = max(dp[i-1], dp[i-2] + nums[i-1]). Answer dp[n]."""
    n = len(nums)
    dp = [0] * (n + 1)
    if n:
        dp[1] = nums[0]
    for i in range(2, n + 1):
        dp[i] = max(dp[i - 1], dp[i - 2] + nums[i - 1])
    return dp[n]
```

**Stage 4: rolling variables.** Only the last two entries are read, so keep two scalars: O(1) space.

```python
def rob_rolling(nums: list[int]) -> int:
    """Step 4: dp[i] reads only dp[i-1], dp[i-2] -> keep two scalars. O(1) space."""
    prev2, prev1 = 0, 0          # dp[i-2], dp[i-1]
    for x in nums:
        prev2, prev1 = prev1, max(prev1, prev2 + x)
    return prev1
```

**Complexity:** stage 1: O(2^n) time (about 1.6^n calls), O(n) stack. Stage 2: O(n) time, O(n) cache and recursion depth (hits Python's limit near n = 1000). Stage 3: O(n) time and space. Stage 4: O(n) time, O(1) space.

**Edge cases:** empty list returns 0 (stage 3 guards `dp[1]` with `if n`); one house returns its value. Circular version (House Robber II): see Variations.

</details>

#### Example 1.3: Word break (1D with a variable window) [Medium]

**Problem:** [Word Break (139)](https://leetcode.com/problems/word-break/) · **Template:** Template 1 (1D linear)

Problem (LeetCode: Word Break): can a string be split into a sequence of dictionary words?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "can it be split" = feasibility over prefixes.

**Approach:** dp[i] = `s[:i]` can be segmented = any(dp[j] and s[j:i] in dict), with j bounded by the longest word length.

```python
def word_break(s: str, words: list[str]) -> bool:
    """dp[i] = s[:i] can be segmented. Last decision: the last word s[j:i]."""
    word_set = set(words)
    max_len = max(map(len, word_set), default=0)
    n = len(s)
    dp = [False] * (n + 1)
    dp[0] = True
    for i in range(1, n + 1):
        for j in range(max(0, i - max_len), i):
            if dp[j] and s[j:i] in word_set:
                dp[i] = True
                break
    return dp[n]
```

**Complexity:** O(n * L * L), L = max word length (L choices of j, O(L) to slice and hash). Space O(n).

**Edge cases:** empty string is True. Word Break II (module 15) lists all splits: backtracking + memo.

</details>

### Template 2: 2D grid

> - **Example problem:** [Minimum Path Sum (64)](https://leetcode.com/problems/minimum-path-sum/)
> - **Why this template:** Right/down moves + "minimize the sum" is exactly `grid[r][c] + min(up, left)`.
> - **Worked below:** Example 2.1
> - **Also fits:** [Unique Paths II (63)](https://leetcode.com/problems/unique-paths-ii/)

```python
def dp_grid_template(grid: list[list[int]]) -> int:
    """2D grid DP, moves right/down. dp[r][c] = answer for reaching (r, c)."""
    R, C = len(grid), len(grid[0])
    INF = float("inf")
    dp = [[INF] * C for _ in range(R)]
    for r in range(R):
        for c in range(C):
            if r == 0 and c == 0:
                dp[r][c] = grid[0][0]           # TODO: base
                continue
            up = dp[r - 1][c] if r > 0 else INF
            left = dp[r][c - 1] if c > 0 else INF
            dp[r][c] = grid[r][c] + min(up, left)   # TODO: combine
    return dp[R - 1][C - 1]
```

Knobs: combine (`min` + cost, or `+` to count), blocked cells (0 or INF). Row-major order works because up and left come earlier. One row suffices: `dp[c]` holds "up" until overwritten.

#### Example 2.1: Unique paths with obstacles and minimum path sum (2D grid) [Medium]

**Problems:** [Unique Paths II (63)](https://leetcode.com/problems/unique-paths-ii/), [Minimum Path Sum (64)](https://leetcode.com/problems/minimum-path-sum/) · **Template:** Template 2 (2D grid)

Problem (LeetCode: Unique Paths II): count right/down paths from top-left to bottom-right avoiding blocked cells. Problem (LeetCode: Minimum Path Sum): min-cost right/down path.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** grid + right/down moves only = acyclic, so row-major order is valid.

**Approach:** paths[r][c] = above + left, 0 if blocked. In a rolling row, `dp[c]` before update is the row above; `dp[c-1]` is already the current row's left neighbor.

```python
def unique_paths_with_obstacles(grid: list[list[int]]) -> int:
    """Paths from top-left to bottom-right moving right/down, 1 = blocked.
    Rolling single row: dp[c] holds the row above until overwritten."""
    if not grid or grid[0][0] == 1:
        return 0
    C = len(grid[0])
    dp = [0] * C
    dp[0] = 1
    for row in grid:
        for c in range(C):
            if row[c] == 1:
                dp[c] = 0
            elif c > 0:
                dp[c] += dp[c - 1]               # dp[c] (from above) + dp[c-1] (from left)
    return dp[-1]


def min_path_sum(grid: list[list[int]]) -> int:
    R, C = len(grid), len(grid[0])
    dp = [float("inf")] * C
    dp[0] = 0
    for r in range(R):
        dp[0] += grid[r][0]
        for c in range(1, C):
            dp[c] = grid[r][c] + min(dp[c], dp[c - 1])
    return int(dp[-1])
```

**Complexity:** O(R * C) time, O(C) space.

**Edge cases:** blocked start or end gives 0. A blocked first-row cell blocks everything to its right (rolling handles this: `dp[c]` stays 0, adding only from the left). Up/left moves create cycles that break DP order: use BFS/Dijkstra (modules 13, 14).

</details>

### Template 3: Two sequences

> - **Example problem:** [Longest Common Subsequence (1143)](https://leetcode.com/problems/longest-common-subsequence/)
> - **Why this template:** A common subsequence of two strings means dp[i][j] over prefixes `a[:i]`, `b[:j]`; the template code is LCS.
> - **Worked below:** Example 3.2
> - **Also fits:** [Edit Distance (72)](https://leetcode.com/problems/edit-distance/)

```python
def dp_two_seq_template(a: str, b: str) -> int:
    """Two-sequence DP. dp[i][j] = answer for prefixes a[:i], b[:j].
    Row 0 / col 0 encode 'one string empty'. This instance is LCS."""
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    # TODO: fill dp[i][0] and dp[0][j] base cases (0 for LCS; i / j for edit distance)
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1          # TODO: match case
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])  # TODO: mismatch case
    return dp[m][n]
```

Knobs: base row/column, match case, mismatch case. Diagonal `(i-1, j-1)` consumes a char of each; up `(i-1, j)` one of `a`; left `(i, j-1)` one of `b`. Space: two rows, or one row plus a saved diagonal (O(min(m, n)) with the shorter string as columns).

#### Example 3.1: Longest palindromic subsequence [Medium]

**Problem:** [Longest Palindromic Subsequence (516)](https://leetcode.com/problems/longest-palindromic-subsequence/) · **Template:** Template 3 (Two sequences), as LCS(s, reversed s); the range form is a preview of interval DP

Problem (LeetCode: Longest Palindromic Subsequence): length of the longest subsequence that reads the same both ways.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** one string, palindrome, subsequence (not substring) = DP over ranges `[i, j]`.

**Approach:** dp[i][j] = LPS within `s[i..j]`. Matching ends both join (`dp[i+1][j-1] + 2`); otherwise drop one end. Row i reads row i+1, so iterate i downward. This is the gateway to interval DP (module 17); it equals LCS(s, reversed s), which the tests verify.

```python
def longest_palindrome_subseq(s: str) -> int:
    """dp[i][j] = LPS length within s[i..j]. Fill by i descending so that
    dp[i+1][*] is ready. Equivalent to LCS(s, reversed(s))."""
    n = len(s)
    if n == 0:
        return 0
    dp = [[0] * n for _ in range(n)]
    for i in range(n - 1, -1, -1):
        dp[i][i] = 1
        for j in range(i + 1, n):
            if s[i] == s[j]:
                dp[i][j] = dp[i + 1][j - 1] + 2
            else:
                dp[i][j] = max(dp[i + 1][j], dp[i][j - 1])
    return dp[0][n - 1]
```

**Complexity:** O(n^2) time and space (O(n) with a rolling row).

**Edge cases:** when `j = i + 1` and chars match, `dp[i+1][j-1]` is `dp[i+1][i]`, an empty range that is 0: correct. Longest palindromic **substring** is contiguous; expand around center is simpler there.

</details>

#### Example 3.2: Two-sequence DP: LCS, edit distance, distinct subsequences [Medium/Hard]

**Problems:** [Longest Common Subsequence (1143)](https://leetcode.com/problems/longest-common-subsequence/), [Edit Distance (72)](https://leetcode.com/problems/edit-distance/), [Distinct Subsequences (115)](https://leetcode.com/problems/distinct-subsequences/) · **Template:** Template 3 (Two sequences)

**LCS** (LeetCode: Longest Common Subsequence): given two strings, return the length of their longest common subsequence (characters in order, not necessarily contiguous). Example: `"abcde"`, `"ace"` → `3` (`"ace"`).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** two strings + "subsequence" + "longest" = two-sequence DP over prefixes.

**Approach:** dp[i][j] = LCS of `a[:i]` and `b[:j]`. Matching last characters extend the LCS: dp[i-1][j-1] + 1. Otherwise one is unused: max(dp[i-1][j], dp[i][j-1]). Row/column 0 are 0; answer dp[m][n].

```python
def lcs(a: str, b: str) -> int:
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]   # row 0 / col 0: one string empty -> 0
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1          # both last chars join the LCS
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])  # drop one of them
    return dp[m][n]
```

**Complexity:** O(m·n) time and space; O(min(m, n)) space with rolling rows.

**Edge cases:** an empty string gives 0. On a mismatch the LCS uses at most one of the last characters, not neither. To recover the subsequence, walk the full table back from (m, n).

</details>

**Edit distance** (LeetCode: Edit Distance): min insert/delete/replace operations turning `a` into `b`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** two strings + transform + min operations.

**Approach:** dp[i][j] = min edits turning `a[:i]` into `b[:j]`. Matching last chars: dp[i-1][j-1]. Otherwise 1 + min of replace (diagonal), delete `a[i-1]` (up), insert `b[j-1]` (left). Base: dp[i][0] = i, dp[0][j] = j.

```python
def edit_distance(a: str, b: str) -> int:
    """dp[i][j] = min edits turning a[:i] into b[:j]. Last operation on the
    last chars: match (free), replace, delete a[i-1], insert b[j-1]."""
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(m + 1):
        dp[i][0] = i                             # delete everything
    for j in range(n + 1):
        dp[0][j] = j                             # insert everything
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = 1 + min(dp[i - 1][j - 1],   # replace
                                   dp[i - 1][j],       # delete from a
                                   dp[i][j - 1])       # insert into a
    return dp[m][n]
```

</details>

**Distinct subsequences** (LeetCode: Distinct Subsequences): count the ways `t` appears as a subsequence of `s`.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** the counting version of the same grid: skip `s[i-1]` (`dp[i-1][j]`) or, when equal, match it with `t[j-1]` (`+ dp[i-1][j-1]`). Base: empty `t` appears once in any prefix.

```python
def num_distinct(s: str, t: str) -> int:
    """Count subsequences of s equal to t. dp[i][j] = ways t[:j] appears in s[:i].
    s[i-1] either is skipped, or (if equal) matched to t[j-1]."""
    m, n = len(s), len(t)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(m + 1):
        dp[i][0] = 1                             # empty t: one way (pick nothing)
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            dp[i][j] = dp[i - 1][j]
            if s[i - 1] == t[j - 1]:
                dp[i][j] += dp[i - 1][j - 1]
    return dp[m][n]
```

**Complexity:** O(m * n) time and space for edit distance and distinct subsequences; O(min(m, n)) space with rolling rows.

**Edge cases:** the base row/column handle empty strings. Rolling edit distance must save the old diagonal before overwriting `dp[j-1]`.

</details>

### Template 4: 0/1 knapsack (each item at most once)

> - **Example problem:** [Partition Equal Subset Sum (416)](https://leetcode.com/problems/partition-equal-subset-sum/)
> - **Why this template:** Each number is used at most once; the modest target total/2 is the capacity.
> - **Worked below:** Example 4.1
> - **Also fits:** [Target Sum (494)](https://leetcode.com/problems/target-sum/)

```python
def knapsack_01_template(weights: list[int], values: list[int], cap: int) -> int:
    """0/1 knapsack, 1D. dp[w] = best value with capacity w using items seen
    so far. Iterate w DOWNWARD so dp[w - wt] still holds the previous item
    row, i.e. this item has not been used yet."""
    dp = [0] * (cap + 1)
    for wt, val in zip(weights, values):
        for w in range(cap, wt - 1, -1):            # backward: each item once
            dp[w] = max(dp[w], dp[w - wt] + val)
    return dp[cap]
```

**Why backward?** In 2D, `dp[i][w] = max(dp[i-1][w], dp[i-1][w-wt] + val)` reads only the **previous item's row**, which 1D overwrites in place. Going high to low, `dp[w - wt]` isn't overwritten yet, so it is still the previous row. Going forward, it would already include the current item, allowing reuse.

```
cap = 4, item wt=1 val=15, starting dp = [0,0,0,0,0]
forward : dp = [0,15,30,45,60]   item used 4 times (wrong for 0/1)
backward: dp = [0,15,15,15,15]   item used at most once
```

Knobs: combine op (`max` value, `or` feasibility, `+` count); initial values (`dp[0] = 0/True/1`; others `0/False/0` for "at most cap", or `-inf`/`inf` if capacity must be filled exactly).

#### Example 4.1: 0/1 knapsack: partition equal subset sum and target sum [Medium]

**Problems:** [Partition Equal Subset Sum (416)](https://leetcode.com/problems/partition-equal-subset-sum/), [Target Sum (494)](https://leetcode.com/problems/target-sum/) · **Template:** Template 4 (0/1 knapsack)

**Partition equal subset sum** (LeetCode: Partition Equal Subset Sum): can the array be split into two equal-sum groups?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** two equal groups = one subset summing to total/2; each element at most once = boolean 0/1 knapsack.

**Approach:** `dp[s]` = some subset so far sums to `s`. Iterate `s` from `half` down to `x` so each item is used once.

```python
def can_partition(nums: list[int]) -> bool:
    """Split into two equal-sum halves <=> some subset sums to total/2.
    0/1 knapsack over booleans, backward iteration."""
    total = sum(nums)
    if total % 2:
        return False
    half = total // 2
    dp = [False] * (half + 1)
    dp[0] = True
    for x in nums:
        for s in range(half, x - 1, -1):
            dp[s] = dp[s] or dp[s - x]
    return dp[half]
```

</details>

**Target sum** (LeetCode: Target Sum): put `+` or `-` before each number; count assignments that evaluate to a target.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** signs look like 2^n backtracking, but with P = sum of the `+` group, `P - (total - P) = target`, so `P = (total + target) / 2`: count subsets summing to P.

```python
def target_sum_ways(nums: list[int], target: int) -> int:
    """Assign + or - to each number to hit target. With P = sum of '+' set:
    P - (total - P) = target  ->  P = (total + target) / 2. Count subsets
    summing to P (0/1 knapsack, counting)."""
    total = sum(nums)
    if abs(target) > total or (total + target) % 2:
        return 0
    p = (total + target) // 2
    dp = [1] + [0] * p
    for x in nums:
        for s in range(p, x - 1, -1):
            dp[s] += dp[s - x]
    return dp[p]
```

**Complexity:** O(n * S) time, S = subset target, O(S) space. Pseudo-polynomial: fine for small sums (e.g., <= 2 * 10^4), useless for values up to 10^9 (see look-alikes).

**Edge cases:** odd total means no partition. Zeros in Target Sum double the count (`+0` and `-0` differ); a zero item does `dp[s] += dp[s]`, handling it. The `abs(target)` check and the formula handle a negative target.

</details>

### Template 5: Unbounded knapsack (reuse allowed)

> - **Example problem:** [Coin Change (322)](https://leetcode.com/problems/coin-change/)
> - **Why this template:** "An infinite number of each kind of coin" means reuse; the amount (up to 10^4) is the capacity; combine with min.
> - **Worked below:** Example 5.1
> - **Also fits:** [Coin Change II (518)](https://leetcode.com/problems/coin-change-ii/)

```python
def knapsack_unbounded_template(weights: list[int], values: list[int], cap: int) -> int:
    """Unbounded knapsack. Iterate w UPWARD so dp[w - wt] may already include
    this item: that is exactly 'reuse allowed'."""
    dp = [0] * (cap + 1)
    for wt, val in zip(weights, values):
        for w in range(wt, cap + 1):                # forward: unlimited copies
            dp[w] = max(dp[w], dp[w - wt] + val)
    return dp[cap]
```

The only difference from 0/1 is the inner loop's direction. Memorize the reason, not the direction.

#### Example 5.1: Coin change: min coins, combinations, permutations (unbounded knapsack) [Medium]

**Problems:** [Coin Change (322)](https://leetcode.com/problems/coin-change/), [Coin Change II (518)](https://leetcode.com/problems/coin-change-ii/), [Combination Sum IV (377)](https://leetcode.com/problems/combination-sum-iv/) · **Template:** Template 5 (Unbounded knapsack)

Problem A (LeetCode: Coin Change): fewest coins to make an amount, or -1.
Problem B (LeetCode: Coin Change II): number of coin combinations (order ignored).
Problem C (LeetCode: Combination Sum IV): number of ordered sequences.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** unlimited denominations + target = unbounded knapsack. "Fewest" = min; combinations vs sequences sets loop order.

**Approach:** (A) dp[a] = fewest coins summing to a; the last coin is some c: dp[a] = 1 + min(dp[a - c]). `amount + 1` works as infinity (all 1-coins need only `amount`).

```python
def coin_change_min(coins: list[int], amount: int) -> int:
    """Fewest coins summing to amount, or -1. dp[a] = min coins for a.
    Unbounded knapsack with min instead of max."""
    INF = amount + 1
    dp = [0] + [INF] * amount
    for a in range(1, amount + 1):
        for c in coins:
            if c <= a and dp[a - c] + 1 < dp[a]:
                dp[a] = dp[a - c] + 1
    return dp[amount] if dp[amount] != INF else -1
```

For min/max either loop order works. For counting (B, C) it matters (see the loop-order table):

```python
def coin_change_ways(coins: list[int], amount: int) -> int:
    """Number of COMBINATIONS (order ignored). Coins in the OUTER loop: each
    coin is introduced once, so a multiset is counted in one canonical order."""
    dp = [1] + [0] * amount
    for c in coins:
        for a in range(c, amount + 1):
            dp[a] += dp[a - c]
    return dp[amount]


def count_ordered_ways(nums: list[int], target: int) -> int:
    """Number of SEQUENCES (order matters; Combination Sum IV). Amount in the
    OUTER loop: for each total, any number may be the last one placed."""
    dp = [1] + [0] * target
    for a in range(1, target + 1):
        for x in nums:
            if x <= a:
                dp[a] += dp[a - x]
    return dp[target]
```

**Complexity:** O(len(coins) * amount) time, O(amount) space for all three.

**Edge cases:** amount 0 returns 0 coins and 1 way. Greedy fails for A: coins `[1, 3, 4]`, amount 6: greedy 4+1+1 (3 coins) vs 3+3 (2). With `coins = [1, 2]`, amount 3: B returns 2 (`{1,1,1}`, `{1,2}`), C returns 3 (`111`, `12`, `21`).

</details>

## Template cheat sheet

### Template 1: 1D linear
```python
def dp_1d_template(items: list[int]) -> int:
    """1D linear DP over a prefix. Customize: meaning of dp[i], base cases,
    which earlier indices the transition reads, combine op (max/min/sum/or)."""
    n = len(items)
    dp = [0] * (n + 1)           # dp[i] = answer for items[:i]
    dp[0] = 0                    # TODO: base case for the empty prefix
    for i in range(1, n + 1):
        # TODO: transition from the last decision about items[i-1]
        dp[i] = dp[i - 1] + items[i - 1]
    return dp[n]                 # TODO: or max(dp) if the answer can end anywhere
```

### Template 2: 2D grid
```python
def dp_grid_template(grid: list[list[int]]) -> int:
    """2D grid DP, moves right/down. dp[r][c] = answer for reaching (r, c)."""
    R, C = len(grid), len(grid[0])
    INF = float("inf")
    dp = [[INF] * C for _ in range(R)]
    for r in range(R):
        for c in range(C):
            if r == 0 and c == 0:
                dp[r][c] = grid[0][0]           # TODO: base
                continue
            up = dp[r - 1][c] if r > 0 else INF
            left = dp[r][c - 1] if c > 0 else INF
            dp[r][c] = grid[r][c] + min(up, left)   # TODO: combine
    return dp[R - 1][C - 1]
```

### Template 3: Two sequences
```python
def dp_two_seq_template(a: str, b: str) -> int:
    """Two-sequence DP. dp[i][j] = answer for prefixes a[:i], b[:j].
    Row 0 / col 0 encode 'one string empty'. This instance is LCS."""
    m, n = len(a), len(b)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    # TODO: fill dp[i][0] and dp[0][j] base cases (0 for LCS; i / j for edit distance)
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if a[i - 1] == b[j - 1]:
                dp[i][j] = dp[i - 1][j - 1] + 1          # TODO: match case
            else:
                dp[i][j] = max(dp[i - 1][j], dp[i][j - 1])  # TODO: mismatch case
    return dp[m][n]
```

### Template 4: 0/1 knapsack (each item at most once)
```python
def knapsack_01_template(weights: list[int], values: list[int], cap: int) -> int:
    """0/1 knapsack, 1D. dp[w] = best value with capacity w using items seen
    so far. Iterate w DOWNWARD so dp[w - wt] still holds the previous item
    row, i.e. this item has not been used yet."""
    dp = [0] * (cap + 1)
    for wt, val in zip(weights, values):
        for w in range(cap, wt - 1, -1):            # backward: each item once
            dp[w] = max(dp[w], dp[w - wt] + val)
    return dp[cap]
```
```
cap = 4, item wt=1 val=15, starting dp = [0,0,0,0,0]
forward : dp = [0,15,30,45,60]   item used 4 times (wrong for 0/1)
backward: dp = [0,15,15,15,15]   item used at most once
```

### Template 5: Unbounded knapsack (reuse allowed)
```python
def knapsack_unbounded_template(weights: list[int], values: list[int], cap: int) -> int:
    """Unbounded knapsack. Iterate w UPWARD so dp[w - wt] may already include
    this item: that is exactly 'reuse allowed'."""
    dp = [0] * (cap + 1)
    for wt, val in zip(weights, values):
        for w in range(wt, cap + 1):                # forward: unlimited copies
            dp[w] = max(dp[w], dp[w - wt] + val)
    return dp[cap]
```

## Variations and extensions

- **Reconstruct the solution, not just its value.** Keep the full table (or parent pointers) and walk back from the answer cell, checking which transition produced each value. Rolling arrays lose this; mention the trade-off.
- **Circular input** (House Robber II): run the linear DP twice, excluding the first element, then the last.
- **Bounded knapsack** (item i has `c_i` copies): 0/1 over binary-split bundles (1, 2, 4, ...), or unbounded with a count limit.
- **Multiple capacities** (Ones and Zeroes): 0/1 knapsack with 2D capacity; iterate both dimensions backward.
- **Follow-ups**: "O(1)/O(n) space?" (rolling). "n is 10^9?" (matrix exponentiation for linear recurrences like Fibonacci). "Modulo 10^9 + 7" (reduce at every addition).
- **DP + other patterns**: prefix sums for range costs in transitions; binary search for a predecessor (module 17, weighted job scheduling); monotonic deque for a sliding-window max over dp (Jump Game VI / Constrained Subsequence Sum).

## Common pitfalls

- **Vague state.** "dp[i] is the answer at i" is not a definition; say "for the prefix of length i" or "ending exactly at index i".
- **Off-by-one between dp and input index.** With size `n + 1`, `dp[i]` refers to `s[i-1]`. Pick a convention and comment it.
- **Wrong counting base case.** dp[0] = 1 (one way to make nothing); 0 zeroes everything.
- **Forward iteration in 0/1 knapsack** (silently allows reuse) or backward in unbounded (silently forbids it).
- **Wrong loop order when counting with reusable parts**: combinations need items outer; sequences need amount outer.
- **Returning `float('inf')`** where an int is expected; `INF + 1` overflows in fixed-int languages (not Python, but say it).
- **Shared rows**: `[[0] * n] * m` aliases every row. Use a comprehension.
- **Recursion depth** with `@cache` at n = 10^4+. Go bottom-up or raise the limit (deep Python recursion can also crash the interpreter's C stack).
- **Greedy temptation**: coin change with arbitrary denominations, or "take the bigger neighbor" in House Robber, fail on small counterexamples. Spend 30 seconds hunting one before choosing greedy.

## When NOT to use / look-alikes

| Looks like DP | But is | Tell |
|---|---|---|
| List all solutions | Backtracking (module 15) | Output itself is exponential; DP only helps count/optimize |
| Shortest path with arbitrary moves | BFS / Dijkstra (modules 13, 14) | Moves in all directions (cycles), so no topological order |
| Locally best choice provably safe | Greedy (module 18) | Exchange argument works (e.g., interval scheduling by end time) |
| Max subarray sum, contiguous | Kadane (a tiny DP) or prefix sums | One pass, "ending at i" state, O(1) memory |
| Sum target with negative numbers, contiguous subarray | Prefix sums + hash map (module 04) | Contiguous + count, not a subset |
| Huge values (10^9) in a subset-sum | Meet in the middle, or other structure | Pseudo-polynomial table too large |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [House Robber (198)](https://leetcode.com/problems/house-robber/) · Medium · 1D take/skip: do all four stages of the progression.
2. [Coin Change (322)](https://leetcode.com/problems/coin-change/) · Medium · Unbounded knapsack (min): min over last coin.
3. [Word Break (139)](https://leetcode.com/problems/word-break/) · Medium · 1D over prefixes: variable window; bound by max word length.
4. [Longest Common Subsequence (1143)](https://leetcode.com/problems/longest-common-subsequence/) · Medium · Two-sequence: the template.
5. [Edit Distance (72)](https://leetcode.com/problems/edit-distance/) · Medium · Two-sequence: three operations map to three neighbors.
6. [Partition Equal Subset Sum (416)](https://leetcode.com/problems/partition-equal-subset-sum/) · Medium · 0/1 knapsack: backward boolean knapsack.

### Optional

**1D linear**
- [Climbing Stairs (70)](https://leetcode.com/problems/climbing-stairs/) · Easy · Fibonacci from the last step.
- [Min Cost Climbing Stairs (746)](https://leetcode.com/problems/min-cost-climbing-stairs/) · Easy · Same shape with costs; mind where the "top" is.
- [House Robber II (213)](https://leetcode.com/problems/house-robber-ii/) · Medium · Circle = two linear runs.
- [Decode Ways (91)](https://leetcode.com/problems/decode-ways/) · Medium · Two possible last tokens, validity checks.
- [Maximum Subarray (53)](https://leetcode.com/problems/maximum-subarray/) · Medium · Kadane: best sum ending at i = max(x, best ending at i-1 + x).
- [Maximum Product Subarray (152)](https://leetcode.com/problems/maximum-product-subarray/) · Medium · Track both max and min ending at i.

**Unbounded knapsack**
- [Coin Change II (518)](https://leetcode.com/problems/coin-change-ii/) · Medium · Items outer loop for combinations.
- [Combination Sum IV (377)](https://leetcode.com/problems/combination-sum-iv/) · Medium · Amount outer loop for sequences.
- [Perfect Squares (279)](https://leetcode.com/problems/perfect-squares/) · Medium · Coins are the squares.

**2D grid**
- [Unique Paths (62)](https://leetcode.com/problems/unique-paths/) · Medium · Additive counting from top and left.
- [Unique Paths II (63)](https://leetcode.com/problems/unique-paths-ii/) · Medium · Same as Unique Paths; obstacles force 0.
- [Minimum Path Sum (64)](https://leetcode.com/problems/minimum-path-sum/) · Medium · Rolling row.
- [Triangle (120)](https://leetcode.com/problems/triangle/) · Medium · Bottom-up from the last row makes the answer a single cell.
- [Maximal Square (221)](https://leetcode.com/problems/maximal-square/) · Medium · dp = side of the largest square ending at (r, c): min of three neighbors + 1.

**Two sequences**
- [Delete Operation for Two Strings (583)](https://leetcode.com/problems/delete-operation-for-two-strings/) · Medium · Derive from LCS.
- [Interleaving String (97)](https://leetcode.com/problems/interleaving-string/) · Medium · dp[i][j] = can a[:i] and b[:j] form c[:i+j].
- [Distinct Subsequences (115)](https://leetcode.com/problems/distinct-subsequences/) · Hard · Counting version; empty target base row.

**0/1 knapsack and palindromes**
- [Target Sum (494)](https://leetcode.com/problems/target-sum/) · Medium · Algebra converts signs to a subset count.
- [Last Stone Weight II (1049)](https://leetcode.com/problems/last-stone-weight-ii/) · Medium · Minimize difference of two groups = subset closest to total/2.
- [Ones and Zeroes (474)](https://leetcode.com/problems/ones-and-zeroes/) · Medium · Two-dimensional capacity.
- [Longest Palindromic Subsequence (516)](https://leetcode.com/problems/longest-palindromic-subsequence/) · Medium · Range DP, i descending.

## Self-check

1. Apply the 7 steps to "minimum cost to climb stairs where each step has a cost."
2. Why does iterating capacity backward in 1D 0/1 knapsack prevent reusing an item? Show it on a two-cell example.
3. Coin Change II and Combination Sum IV have nearly identical code. What single change separates them, and why does it change what is counted?
4. When is "dp[i] = answer ending at i" required instead of "answer for prefix up to i"? Give an example of each.
5. In edit distance, which neighbor is insertion and which is deletion, and why?
6. What does dp[0] = 1 mean in a counting DP, and what goes wrong if it is 0?
7. Your memoized solution hits RecursionError at n = 50,000. Name two fixes, and which you'd choose.
