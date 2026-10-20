# 15. Backtracking

Backtracking is careful exhaustive search: build a candidate one decision at a time, abandoning it once it can't lead to a valid answer. It handles "list all", "find any valid", and "count all configurations" on small inputs (n up to ~20) when nothing avoids enumeration, and is the brute-force start for most DP (module 16).

## At a glance

- **Use when:** "all" / "every" / "any valid" arrangement, tiny constraints (n ≤ ~20), each choice restricts later ones.
- **Templates:**
  - Template 1: universal template with explicit hooks: grids, partitions, constraint boards (Word Search, N-Queens, Sudoku).
  - Template 2: subsets and combinations: order doesn't matter; recurse from `i + 1` (or `i` with reuse).
  - Template 3: permutations: order matters; scan all indices with a `used` array.
- **Typical complexity:** O(n * 2^n) for subsets, O(n * n!) for permutations; pruning cuts real node counts.
- **Top trap:** appending `path` instead of `path[:]`, or a missing/mismatched undo.

## Core idea

Every backtracking problem walks an implicit **decision tree**: the root is the empty partial solution, each edge a choice, each leaf a complete solution or dead end.

```
                     []                       decision: take 1?
              /              \
           [1]                []              decision: take 2?
         /     \            /    \
     [1,2]     [1]       [2]     []           decision: take 3?
     /  \     /  \      /  \    /  \
 [1,2,3][1,2][1,3][1][2,3][2] [3]  []         leaves = all 8 subsets
```

You never build the tree. Walk **one mutable partial solution** (`path`) depth-first:

1. **Choose**: apply a move (append, mark used, place a queen).
2. **Explore**: recurse.
3. **Unchoose**: undo exactly that move.

**Invariant**: a returning call leaves shared state exactly as it found it, so siblings share one `path` and `used` array without copying. Break it and siblings see garbage.

**Pruning** separates backtracking from brute force: check constraints *before* descending to skip whole subtrees. Cost is roughly (nodes visited) x (work per node); without pruning, about `branching_factor ^ depth`.

| Problem shape | Depth | Branching | Leaves | Typical total |
|---|---|---|---|---|
| Subsets of n | n | 2 (in/out) | 2^n | O(n * 2^n) incl. copying |
| Permutations of n | n | n, n-1, ... | n! | O(n * n!) |
| k-combinations of n | k | <= n | C(n, k) | O(k * C(n, k)) |
| N-Queens | n rows | <= n cols | far fewer than n^n after pruning | ~O(n!) |
| Word search, word len L | L | 3 (after the first step) | | O(R * C * 3^L) |

The `* n` factor is copying `path` at each leaf.

## Recognition signals

- **"All"**: all subsets, permutations, arrangements, splits, paths. The output is exponential, so no polynomial algorithm exists.
- **"Find any valid configuration"** under constraints: Sudoku, N-Queens, scheduling with conflicts, graph coloring.
- **Tiny constraints**: n <= 10 (permutations), n <= 20 (subsets, 2^20 ~ 10^6), 9x9 board, word length <= 15.
- The answer is a **sequence of choices**, each restricting later ones (a used item, an occupied column, a cell already in the path).
- Strings **split** into valid pieces (palindromes, IP octets, dictionary words).
- Generating **well-formed structures**: balanced parentheses, valid expressions, letter combinations.
- Only a **count** or **min/max** with repeated subproblems → probably DP (see "When NOT to use").

## Disguises

1. **On-call rotation.** "Given 8 engineers and 8 weekly shifts, each engineer has blackout weeks. List every valid assignment where each engineer takes exactly one shift."
   Decode: one-to-one assignment = permutation; blackout weeks = prune; "list every" rules out DP. Permutation template with a `prune(engineer, week)` hook.

2. **Feature flag bundles.** "A release can enable any subset of 12 flags, but some flags are mutually exclusive. Enumerate all legal bundles."
   Decode: subsets of 12 = 4096 nodes; check mutual exclusion before including. Subset template with a conflict set.

3. **Budget-constrained instance mix.** "Instance types cost 2, 3, 6, 7 credits per hour. List every multiset of instances that costs exactly 7 credits/hour."
   Decode: "multiset" = order irrelevant, reuse allowed. Combination Sum: recurse with `i`, not `i + 1`; sort and break on overshoot.

4. **Rack placement.** "Place n replicas on an n x n grid of racks so that no two share a row, a column, or a diagonal power bus."
   Decode: N-Queens. One replica per row; sets for columns and both diagonals.

5. **Log line tokenizer.** "A log message lost its delimiters. Given a vocabulary, return every way to re-insert spaces so each token is a known word."
   Decode: partitioning at cut points. The suffix at index i is re-solved many times, so memoize on i (Word Break II).

6. **Config path lookup in a character grid.** "A 2D char grid of config shards; does a key appear as a path of adjacent cells without reusing a cell?"
   Decode: grid DFS with no reuse = mark visited going down, unmark coming up. Word Search.

## Templates and worked examples

### Knob summary

| Knob | Subsets | Combinations (size k) | Combination sum (reuse) | Combination sum (no reuse) | Permutations |
|---|---|---|---|---|---|
| Loop range | `start..n` | `start..n` | `start..n` | `start..n` | `0..n` with `used` |
| Recurse with | `i + 1` | `i + 1` | `i` | `i + 1` | (same level, `used` updated) |
| Record when | every node | `len(path) == k` | `remaining == 0` | `remaining == 0` | `len(path) == n` |
| Duplicates in input | sort, skip `i > start and a[i] == a[i-1]` | same | (usually distinct) | sort + same skip | sort, skip if `a[i] == a[i-1] and not used[i-1]` |
| Prune | | not enough items left | sorted, break on overshoot | sorted, break | problem constraints |

### Duplicate handling, precisely

Duplicates create identical subtrees; generate one canonical representative.

- **Subsets/combinations**: sort so equal values are adjacent. Within one level (one `for` loop), only a value's first copy starts a branch: skip `i > start and nums[i] == nums[i-1]`. Deeper levels can still use the second copy (giving `[2, 2]`) because there `i == start`.
- **Permutations**: after sorting, use copies left to right. Place copy `i` only if copy `i-1` is in the path: skip when `nums[i] == nums[i-1] and not used[i-1]`.

### Template 1: Universal template with explicit hooks

> - **Example problem:** [N-Queens (51)](https://leetcode.com/problems/n-queens/)
> - **Why this template:** one queen per row, prune attacked columns and diagonals, record full boards: every hook (choices, is_valid, make/undo move, snapshot) has a concrete job.
> - **Worked below:** Example 1.3
> - **Also fits:** [Word Search (79)](https://leetcode.com/problems/word-search/)

Every backtracking solution is this function with hooks filled in. Interviews inline them, but naming them helps design.

```python
def backtrack_universal(
    state: Any,
    is_solution: Callable[[Any], bool],
    candidates: Callable[[Any], Iterable[Any]],
    prune: Callable[[Any, Any], bool],
    make_move: Callable[[Any, Any], None],
    undo_move: Callable[[Any, Any], None],
    snapshot: Callable[[Any], Any],
    stop_at_solution: bool = True,
) -> list[Any]:
    """Universal choose -> explore -> unchoose skeleton with explicit hooks.

    Invariant: when backtrack(state) returns, `state` is exactly what it was
    on entry. That is what lets siblings in the decision tree share one
    mutable state object.

    stop_at_solution=False is for problems where every node is an answer
    (subsets): record, then keep extending.
    """
    results: list[Any] = []

    def backtrack() -> None:
        if is_solution(state):
            results.append(snapshot(state))   # copy! state keeps mutating
            if stop_at_solution:
                return
        for c in candidates(state):
            if prune(state, c):                # cut the branch before descending
                continue
            make_move(state, c)                # choose
            backtrack()                        # explore
            undo_move(state, c)                # unchoose

    backtrack()
    return results
```

| Hook | Question it answers | Examples |
|---|---|---|
| `is_solution` | Is the partial solution complete? | `len(path) == n`, `remaining == 0`, `row == n` |
| `candidates` | What choices exist here? | `range(start, n)`, `range(n)`, letters for a digit, 4 neighbors |
| `prune` | Can this choice never succeed? | used, column attacked, value > remaining, not a palindrome |
| `make_move` / `undo_move` | How to apply / revert the choice? | append/pop, add/remove from set, mark/unmark cell |
| `snapshot` | How to store an answer safely? | `path[:]`, `"".join(path)` |

The companion file wires permutations through this template (`permutations_via_universal`).

#### Example 1.1: Word search in a grid [Medium]

**Problem:** [Word Search (79)](https://leetcode.com/problems/word-search/) · **Template:** Template 1 (universal template with explicit hooks)

Problem: does a word appear in a character grid as a path of horizontally/vertically adjacent cells, using each cell at most once?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** grid + path + no reuse = DFS with undone visited state. A permanent visited set fails: a cell blocked on one path may be needed on another.

**Approach:** from every start cell, the cell at depth k must match `word[k]`. Mark it in place (`'#'`) so deeper calls can't reuse it, then restore it; this make/undo pair replaces a visited set.

```python
def word_exists(board: list[list[str]], word: str) -> bool:
    """Word Search. DFS from every cell; mark visited in place with '#'
    and restore on the way back (the unchoose step)."""
    if not board or not board[0]:
        return False
    R, C = len(board), len(board[0])

    def dfs(r: int, c: int, k: int) -> bool:
        if k == len(word):
            return True
        if not (0 <= r < R and 0 <= c < C) or board[r][c] != word[k]:
            return False
        saved, board[r][c] = board[r][c], "#"   # choose: mark visited
        found = (dfs(r + 1, c, k + 1) or dfs(r - 1, c, k + 1)
                 or dfs(r, c + 1, k + 1) or dfs(r, c - 1, k + 1))
        board[r][c] = saved                     # unchoose: restore
        return found

    return any(dfs(r, c, 0) for r in range(R) for c in range(C))
```

**Complexity:** O(R * C * 3^L) for word length L: after the first step each cell has at most 3 unvisited neighbors. O(L) recursion space.

**Edge cases:** the `or` chain short-circuits, but the restore still runs because it follows the expression; `return dfs(...) or dfs(...)` before restoring corrupts the board. Pruning: compare character counts of board and word first; search from the word's rarer end. For many words, use a trie (module 19, Word Search II).

</details>

#### Example 1.2: Palindrome partitioning [Medium]

**Problem:** [Palindrome Partitioning (131)](https://leetcode.com/problems/palindrome-partitioning/) · **Template:** Template 1 (universal template with explicit hooks)

Problem: split a string into palindromic pieces; return all such splits.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** all ways to split a string = partitioning backtracking; each node picks where the next piece ends.

**Approach:** state is the start of the unsplit suffix. Try every end; descend only if `s[start..end]` is a palindrome. A precomputed interval-DP `is_pal` table makes checks O(1), not O(n).

```python
def palindrome_partition(s: str) -> list[list[str]]:
    """Choose the next cut position; precompute palindrome table so each
    check is O(1)."""
    n = len(s)
    is_pal = [[False] * n for _ in range(n)]
    for i in range(n - 1, -1, -1):
        for j in range(i, n):
            is_pal[i][j] = s[i] == s[j] and (j - i < 2 or is_pal[i + 1][j - 1])

    res: list[list[str]] = []
    path: list[str] = []

    def backtrack(start: int) -> None:
        if start == n:
            res.append(path[:])
            return
        for end in range(start, n):
            if is_pal[start][end]:
                path.append(s[start:end + 1])
                backtrack(end + 1)
                path.pop()

    backtrack(0)
    return res
```

**Complexity:** O(n * 2^n): each of n - 1 gaps is cut or not, and copying a leaf is O(n). The table is O(n^2) time and space.

**Edge cases:** one character returns `[[c]]`. "Minimum number of cuts" is DP, not backtracking (module 17, Palindrome Partitioning II).

**Restore IP addresses** partitions into a fixed number of parts with a length prune: with `left` parts to fill, between `left` and `3 * left` characters must remain.

```python
def restore_ip_addresses(s: str) -> list[str]:
    """Exactly 4 parts, each 0..255, no leading zeros. Prune on remaining
    length: k parts left need between k and 3k characters."""
    res: list[str] = []
    parts: list[str] = []

    def backtrack(start: int) -> None:
        left = 4 - len(parts)
        remaining = len(s) - start
        if left == 0:
            if remaining == 0:
                res.append(".".join(parts))
            return
        if remaining < left or remaining > 3 * left:
            return                              # bound prune
        for length in (1, 2, 3):
            seg = s[start:start + length]
            if len(seg) < length:
                break
            if (seg[0] == "0" and length > 1) or int(seg) > 255:
                continue
            parts.append(seg)
            backtrack(start + length)
            parts.pop()

    backtrack(0)
    return res
```

**Complexity:** at most 3^4 = 81 leaves, so effectively O(1) (O(len(s)) per result to join).

</details>

#### Example 1.3: N-Queens [Hard]

**Problem:** [N-Queens (51)](https://leetcode.com/problems/n-queens/) · **Template:** Template 1 (universal template with explicit hooks)

Problem: place n non-attacking queens on an n x n board; return all boards.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** constraint satisfaction on a small board. Each row holds exactly one queen, so depth r picks row r's column.

**Approach:** row-by-row placement rules out row conflicts. Column and diagonal conflicts are O(1) set lookups: down-right diagonals share `r - c`, down-left share `r + c`. Prune attacked columns.

```python
def solve_n_queens(n: int) -> list[list[str]]:
    """One queen per row; O(1) conflict checks via column and diagonal sets.
    Cells on the same '\\' diagonal share r - c; on the same '/' share r + c."""
    res: list[list[str]] = []
    cols: set[int] = set()
    diag: set[int] = set()
    anti: set[int] = set()
    queen_col = [-1] * n

    def backtrack(r: int) -> None:
        if r == n:
            res.append(["." * c + "Q" + "." * (n - c - 1) for c in queen_col])
            return
        for c in range(n):
            if c in cols or (r - c) in diag or (r + c) in anti:
                continue
            cols.add(c); diag.add(r - c); anti.add(r + c)
            queen_col[r] = c
            backtrack(r + 1)
            cols.remove(c); diag.remove(r - c); anti.remove(r + c)

    backtrack(0)
    return res
```

**Complexity:** O(n!) nodes at most (row r has at most n - r free columns), O(n^2) to build each board. O(n) space.

**Edge cases:** n = 1 has one solution; n = 2 and 3 none. To count (N-Queens II), skip building boards. Bitmask sets (`cols | diag | anti` as ints) are a common speed follow-up.

**Sudoku solver** is the same pattern, but "find one": return `True` up the stack once the board is full so filled values aren't undone.

```python
def solve_sudoku(board: list[list[str]]) -> bool:
    """Fill '.' cells in place. Returns True once solved (the True propagates
    up and stops the search so the board keeps the solution)."""
    rows = [set() for _ in range(9)]
    cols = [set() for _ in range(9)]
    boxes = [set() for _ in range(9)]
    empties: list[tuple[int, int]] = []
    for r in range(9):
        for c in range(9):
            v = board[r][c]
            if v == ".":
                empties.append((r, c))
            else:
                rows[r].add(v); cols[c].add(v); boxes[(r // 3) * 3 + c // 3].add(v)

    def backtrack(k: int) -> bool:
        if k == len(empties):
            return True
        r, c = empties[k]
        b = (r // 3) * 3 + c // 3
        for d in "123456789":
            if d in rows[r] or d in cols[c] or d in boxes[b]:
                continue
            board[r][c] = d
            rows[r].add(d); cols[c].add(d); boxes[b].add(d)
            if backtrack(k + 1):
                return True                     # keep the filled board
            board[r][c] = "."
            rows[r].remove(d); cols[c].remove(d); boxes[b].remove(d)
        return False

    return backtrack(0)
```

Follow-up: fill the empty cell with the **fewest legal digits** next (most-constrained-variable heuristic). Same worst case, far fewer nodes in practice.

</details>

### Template 2: subsets and combinations (loop-start style)

> - **Example problem:** [Subsets II (90)](https://leetcode.com/problems/subsets-ii/)
> - **Why this template:** "all possible subsets", order irrelevant: each level picks only indices after the last pick.
> - **Worked below:** Example 2.2 (Subsets (78) in Example 2.1)
> - **Also fits:** [Subsets (78)](https://leetcode.com/problems/subsets/)

Order doesn't matter, so impose one: pick only indices after the last pick. Each tree node is a distinct subset.

```python
def subsets_template(nums: list[int]) -> list[list[int]]:
    """Subset/combination template (loop-start style).

    Knobs: record at every node (subsets) vs only at depth k (combinations);
    recurse with i+1 (each item once) vs i (reuse allowed); sort + skip dups.
    """
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int) -> None:
        res.append(path[:])                     # every node is a subset
        for i in range(start, len(nums)):
            path.append(nums[i])
            backtrack(i + 1)                    # i+1: no reuse; i: reuse
            path.pop()

    backtrack(0)
    return res
```

The **include/exclude** style makes the binary decision explicit, records only at leaves (`i == n`), and converts most directly to DP (index `i` is the state).

```python
def subsets_include_exclude(nums: list[int]) -> list[list[int]]:
    """Binary decision tree: at index i, either take nums[i] or skip it."""
    res: list[list[int]] = []
    path: list[int] = []

    def dfs(i: int) -> None:
        if i == len(nums):
            res.append(path[:])
            return
        path.append(nums[i])     # include
        dfs(i + 1)
        path.pop()
        dfs(i + 1)               # exclude

    dfs(0)
    return res
```

#### Example 2.1: All subsets of distinct values [Medium]

**Problem:** [Subsets (78)](https://leetcode.com/problems/subsets/) · **Template:** Template 2 (subsets and combinations)

Problem: return every subset (the power set) of distinct integers. For `[1, 2, 3]`, 8 lists: `[]`, `[1]`, `[1, 2]`, `[1, 2, 3]`, `[1, 3]`, `[2]`, `[2, 3]`, `[3]`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** every subset + distinct + n <= 10 = enumerate all 2^n choices, order irrelevant → loop-start subset template, no duplicate handling.

**Approach:** after `nums[i]`, consider only larger indices. Every node (including the empty root) is an answer, so record `path[:]` on entry. A subset `{a < b < c}` (by index) has exactly one root-to-node path (`a`, `b`, `c`), so nothing is missed or repeated. Include/exclude works equally.

```python
def subsets_template(nums: list[int]) -> list[list[int]]:
    """Subset/combination template (loop-start style).

    Knobs: record at every node (subsets) vs only at depth k (combinations);
    recurse with i+1 (each item once) vs i (reuse allowed); sort + skip dups.
    """
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int) -> None:
        res.append(path[:])                     # every node is a subset
        for i in range(start, len(nums)):
            path.append(nums[i])
            backtrack(i + 1)                    # i+1: no reuse; i: reuse
            path.pop()

    backtrack(0)
    return res
```

**Complexity:** O(n * 2^n) time (2^n nodes, each copied in O(n)); O(n) extra space for recursion and `path`, plus O(n * 2^n) output.

**Edge cases:** empty input returns `[[]]`, not `[]`. Repeated inputs make it Subsets II (Example 2.2).

</details>

#### Example 2.2: Subsets with duplicates [Medium]

**Problem:** [Subsets II (90)](https://leetcode.com/problems/subsets-ii/) · **Template:** Template 2 (subsets and combinations)

Problem: given integers that may repeat, return every distinct sub-multiset.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** every subset with n <= 10 = subset template. Repeats + distinct = sort and skip.

**Approach:** loop-start on sorted input. At one depth, a branch from the second `2` rebuilds the first `2`'s subsets, so skip it. `i > start` (not `i > 0`) still allows `[2, 2]`: one level deeper, the second `2` is the first candidate.

```python
def subsets_with_dup(nums: list[int]) -> list[list[int]]:
    """Subsets II. Sort; at one tree level, skip a value equal to the
    previous sibling (i > start), which would rebuild the same subtree."""
    nums = sorted(nums)
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int) -> None:
        res.append(path[:])
        for i in range(start, len(nums)):
            if i > start and nums[i] == nums[i - 1]:
                continue                        # same value, same level
            path.append(nums[i])
            backtrack(i + 1)
            path.pop()

    backtrack(0)
    return res
```

**Complexity:** O(n * 2^n) time (at most 2^n nodes, each copied in O(n)); O(n) extra space.

**Edge cases:** empty input returns `[[]]`. `[2, 2, 2]` returns 4 subsets, not 8.

</details>

#### Example 2.3: Combination sum, reuse vs no reuse [Medium]

**Problems:** [Combination Sum (39)](https://leetcode.com/problems/combination-sum/), [Combination Sum II (40)](https://leetcode.com/problems/combination-sum-ii/) · **Template:** Template 2 (subsets and combinations)

Problem A: distinct positive candidates, each reusable; list all multisets summing to a target.
Problem B: candidates may repeat, each element used at most once; list distinct combinations summing to the target.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** combinations (order irrelevant) + sum to target = loop-start template with a `remaining` counter. Reuse decides `i` vs `i + 1`.

**Approach:** sort, so once `cands[i] > remaining` all later candidates overshoot: `break`, not `continue`. A recurses with `i`: the same candidate may repeat, never an earlier one (that would yield permutations of one multiset). B recurses with `i + 1` plus the duplicate skip.

```python
def combination_sum(candidates: list[int], target: int) -> list[list[int]]:
    """Distinct positive candidates, each reusable. Recurse with i (reuse).
    Sorted order lets us break as soon as a candidate overshoots."""
    cands = sorted(candidates)
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int, remaining: int) -> None:
        if remaining == 0:
            res.append(path[:])
            return
        for i in range(start, len(cands)):
            if cands[i] > remaining:
                break                           # sorted: all later ones overshoot too
            path.append(cands[i])
            backtrack(i, remaining - cands[i])  # i, not i+1: reuse allowed
            path.pop()

    backtrack(0, target)
    return res


def combination_sum2(candidates: list[int], target: int) -> list[list[int]]:
    """Candidates may repeat; each used at most once; no duplicate combos."""
    cands = sorted(candidates)
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int, remaining: int) -> None:
        if remaining == 0:
            res.append(path[:])
            return
        for i in range(start, len(cands)):
            if i > start and cands[i] == cands[i - 1]:
                continue
            if cands[i] > remaining:
                break
            path.append(cands[i])
            backtrack(i + 1, remaining - cands[i])
            path.pop()

    backtrack(0, target)
    return res
```

**Complexity:** A is O(n^(T/m + 1)) worst case, T = target, m = smallest candidate (depth ≤ T/m). B is O(n * 2^n). Path space O(T/m) and O(n).

**Edge cases:** an unreachable target returns `[]`. Zero or negative candidates make A recurse forever; say so if constraints change. For only the **count**, use coin-change DP (module 16).

**k-combinations of 1..n** adds a bound prune: still needing `need` numbers, the next can be at most `n - need + 1`.

```python
def combine(n: int, k: int) -> list[list[int]]:
    """All k-subsets of 1..n, with the 'not enough items left' prune."""
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int) -> None:
        if len(path) == k:
            res.append(path[:])
            return
        need = k - len(path)
        # Upper bound: we need `need` numbers from start..n.
        for x in range(start, n - need + 2):
            path.append(x)
            backtrack(x + 1)
            path.pop()

    backtrack(1)
    return res
```

</details>

### Template 3: permutations

> - **Example problem:** [Permutations II (47)](https://leetcode.com/problems/permutations-ii/)
> - **Why this template:** "all unique permutations": order matters, so every level scans all unused elements with a `used` array.
> - **Worked below:** Example 3.2 (Permutations (46) in Example 3.1)
> - **Also fits:** [Permutations (46)](https://leetcode.com/problems/permutations/)

Order matters, so every level considers every unused element; a `used` array replaces the start index.

```python
def permutations_template(nums: list[int]) -> list[list[int]]:
    """Permutation template: order matters, so every level scans all items
    and a used[] array prevents picking one twice."""
    res: list[list[int]] = []
    path: list[int] = []
    used = [False] * len(nums)

    def backtrack() -> None:
        if len(path) == len(nums):
            res.append(path[:])
            return
        for i in range(len(nums)):
            if used[i]:
                continue
            used[i] = True
            path.append(nums[i])
            backtrack()
            path.pop()
            used[i] = False

    backtrack()
    return res
```

#### Example 3.1: All orderings of distinct values [Medium]

**Problem:** [Permutations (46)](https://leetcode.com/problems/permutations/) · **Template:** Template 3 (permutations)

Problem: given distinct integers, return all their orderings. For `[1, 2, 3]`, 6 lists: `[1, 2, 3]`, `[1, 3, 2]`, `[2, 1, 3]`, `[2, 3, 1]`, `[3, 1, 2]`, `[3, 2, 1]`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** every ordering = order matters = permutation template. Distinct = no duplicate handling. n <= 6 confirms n! enumeration.

**Approach:** depth = position being filled; each level tries every unused element; record at length n. Undo both `path` and `used[i]` so siblings see the same state. Each permutation is exactly one choice sequence, so no duplicates.

```python
def permutations_template(nums: list[int]) -> list[list[int]]:
    """Permutation template: order matters, so every level scans all items
    and a used[] array prevents picking one twice."""
    res: list[list[int]] = []
    path: list[int] = []
    used = [False] * len(nums)

    def backtrack() -> None:
        if len(path) == len(nums):
            res.append(path[:])
            return
        for i in range(len(nums)):
            if used[i]:
                continue
            used[i] = True
            path.append(nums[i])
            backtrack()
            path.pop()
            used[i] = False

    backtrack()
    return res
```

**Complexity:** O(n * n!) time (n! leaves, each copied in O(n); internal nodes add a constant factor); O(n) extra space.

**Edge cases:** empty input returns `[[]]` (one empty ordering). Forgetting `used[i] = False` yields far fewer results; a start index yields one ordering. Repeated values produce duplicates; see Example 3.2.

</details>

#### Example 3.2: Permutations with duplicates [Medium]

**Problem:** [Permutations II (47)](https://leetcode.com/problems/permutations-ii/) · **Template:** Template 3 (permutations)

Problem: return all distinct orderings of a list that may contain repeated values.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** all orderings = permutation template. Repeats = canonical order among copies.

**Approach:** sort and use a `used` array. Copy `i` of a value goes in only after copy `i-1`, so copies appear in index order and each distinct permutation appears once.

```python
def permute_unique(nums: list[int]) -> list[list[int]]:
    """Permutations II. Sort; among equal values, only allow using copy i
    if copy i-1 is already in the path (fixes a canonical order of copies)."""
    nums = sorted(nums)
    res: list[list[int]] = []
    path: list[int] = []
    used = [False] * len(nums)

    def backtrack() -> None:
        if len(path) == len(nums):
            res.append(path[:])
            return
        for i in range(len(nums)):
            if used[i]:
                continue
            if i > 0 and nums[i] == nums[i - 1] and not used[i - 1]:
                continue
            used[i] = True
            path.append(nums[i])
            backtrack()
            path.pop()
            used[i] = False

    backtrack()
    return res
```

**Complexity:** O(n * n!) worst case (all distinct); O(n) extra space.

**Edge cases:** `[1, 1]` returns one result. `not used[i-1]` and `used[i-1]` both work (different canonical orders), but `not used[i-1]` prunes far earlier.

</details>

### Generating structures

**Generate parentheses**: the rules "open if `open < n`" and "close if `close < open`" guarantee validity, so leaves need no check. There are Catalan C_n ~ 4^n / n^1.5 results.

```python
def generate_parentheses(n: int) -> list[str]:
    """Add '(' while open < n; add ')' while close < open. These two rules
    guarantee every completed string is valid, so no post-filter."""
    res: list[str] = []
    path: list[str] = []

    def backtrack(open_: int, close: int) -> None:
        if len(path) == 2 * n:
            res.append("".join(path))
            return
        if open_ < n:
            path.append("(")
            backtrack(open_ + 1, close)
            path.pop()
        if close < open_:
            path.append(")")
            backtrack(open_, close + 1)
            path.pop()

    backtrack(0, 0)
    return res
```

**Letter combinations of a phone number**: a Cartesian product; depth = digit index, branching = 3 or 4 letters, O(4^n * n). `itertools.product` does it too, but interviewers want the recursion.

```python
PHONE = {"2": "abc", "3": "def", "4": "ghi", "5": "jkl",
         "6": "mno", "7": "pqrs", "8": "tuv", "9": "wxyz"}


def letter_combinations(digits: str) -> list[str]:
    """Depth = digit index, branching = letters for that digit."""
    if not digits:
        return []
    res: list[str] = []
    path: list[str] = []

    def backtrack(i: int) -> None:
        if i == len(digits):
            res.append("".join(path))
            return
        for ch in PHONE[digits[i]]:
            path.append(ch)
            backtrack(i + 1)
            path.pop()

    backtrack(0)
    return res
```

Gotcha: empty input returns `[]`, not `[""]`.

## Template cheat sheet

All templates, for warm-ups and review.

### Template 1: Universal template with explicit hooks

```python
def backtrack_universal(
    state: Any,
    is_solution: Callable[[Any], bool],
    candidates: Callable[[Any], Iterable[Any]],
    prune: Callable[[Any, Any], bool],
    make_move: Callable[[Any, Any], None],
    undo_move: Callable[[Any, Any], None],
    snapshot: Callable[[Any], Any],
    stop_at_solution: bool = True,
) -> list[Any]:
    """Universal choose -> explore -> unchoose skeleton with explicit hooks.

    Invariant: when backtrack(state) returns, `state` is exactly what it was
    on entry. That is what lets siblings in the decision tree share one
    mutable state object.

    stop_at_solution=False is for problems where every node is an answer
    (subsets): record, then keep extending.
    """
    results: list[Any] = []

    def backtrack() -> None:
        if is_solution(state):
            results.append(snapshot(state))   # copy! state keeps mutating
            if stop_at_solution:
                return
        for c in candidates(state):
            if prune(state, c):                # cut the branch before descending
                continue
            make_move(state, c)                # choose
            backtrack()                        # explore
            undo_move(state, c)                # unchoose

    backtrack()
    return results
```

### Template 2: subsets and combinations (loop-start style)

```python
def subsets_template(nums: list[int]) -> list[list[int]]:
    """Subset/combination template (loop-start style).

    Knobs: record at every node (subsets) vs only at depth k (combinations);
    recurse with i+1 (each item once) vs i (reuse allowed); sort + skip dups.
    """
    res: list[list[int]] = []
    path: list[int] = []

    def backtrack(start: int) -> None:
        res.append(path[:])                     # every node is a subset
        for i in range(start, len(nums)):
            path.append(nums[i])
            backtrack(i + 1)                    # i+1: no reuse; i: reuse
            path.pop()

    backtrack(0)
    return res
```

```python
def subsets_include_exclude(nums: list[int]) -> list[list[int]]:
    """Binary decision tree: at index i, either take nums[i] or skip it."""
    res: list[list[int]] = []
    path: list[int] = []

    def dfs(i: int) -> None:
        if i == len(nums):
            res.append(path[:])
            return
        path.append(nums[i])     # include
        dfs(i + 1)
        path.pop()
        dfs(i + 1)               # exclude

    dfs(0)
    return res
```

### Template 3: permutations

```python
def permutations_template(nums: list[int]) -> list[list[int]]:
    """Permutation template: order matters, so every level scans all items
    and a used[] array prevents picking one twice."""
    res: list[list[int]] = []
    path: list[int] = []
    used = [False] * len(nums)

    def backtrack() -> None:
        if len(path) == len(nums):
            res.append(path[:])
            return
        for i in range(len(nums)):
            if used[i]:
                continue
            used[i] = True
            path.append(nums[i])
            backtrack()
            path.pop()
            used[i] = False

    backtrack()
    return res
```

## Variations and extensions

**When backtracking should become DP.** Does the future depend on the whole `path`, or only a small summary of it?

- In permutations, the future depends on *which* items are used, so memoizing needs a bitmask over n items (bitmask DP, module 17); feasible only for n <= ~20.
- In word segmentation, the future depends only on the start index; many paths reach the same index and re-solve the same suffix. Memoize on `start`:

```python
def word_break_all(s: str, words: list[str]) -> list[str]:
    """Word Break II. Plain backtracking re-solves the same suffix many times;
    memoize on the start index (the only thing the future depends on)."""
    word_set = set(words)
    max_len = max(map(len, word_set), default=0)

    @cache
    def sentences_from(start: int) -> tuple[str, ...]:
        if start == len(s):
            return ("",)
        out: list[str] = []
        for end in range(start + 1, min(len(s), start + max_len) + 1):
            w = s[start:end]
            if w in word_set:
                for rest in sentences_from(end):
                    out.append(w + (" " + rest if rest else ""))
        return tuple(out)

    return list(sentences_from(0))
```

Rule of thumb:

| Question asks for | Subproblems overlap? | Use |
|---|---|---|
| All solutions (list them) | no | backtracking |
| All solutions | yes | backtracking + memo on the summary state (output size still dominates) |
| Count / min / max / feasibility | yes | DP (module 16) |
| Count / min / max, state is a subset of <= 20 items | yes | bitmask DP (module 17) |
| Any one solution, huge space | | backtracking with strong pruning and ordering heuristics |

**Other extensions interviewers reach for:**
- "Return only the count": drop `snapshot`, return sums from recursion. Large counts suggest DP.
- "Return the k-th permutation": don't enumerate; use factorial number system arithmetic.
- "Next permutation in place": a two-pointer trick (module 05), not backtracking.
- "Make it faster": better pruning (sort to enable `break`), bitmasks for sets, most-constrained-first ordering, symmetry breaking (N-Queens: try half of row 0 and double).
- "Iterative version": build subsets iteratively (`res += [s + [x] for s in res]`) or with bitmasks `for mask in range(1 << n)`.
- Backtracking + trie (Word Search II), backtracking + bitmask (Partition to K Equal Sum Subsets, Matchsticks to Square).

## Common pitfalls

- **Appending `path` instead of `path[:]`.** Every result aliases one (eventually empty) list.
- **Missing or mismatched undo.** Every "choose" mutation needs an exact inverse in "unchoose", including sets, counters, and grid marks.
- **Returning before restoring state** in "find one" searches (`return dfs(...) or ...`, restore after). Compute into a variable first, or restore only on failure (Sudoku, which keeps the solved state on purpose).
- **`continue` vs `break` with sorting.** `break` is valid only when candidates are sorted and every later one is at least as bad.
- **Duplicate skip `i > 0` vs `i > start`.** `i > 0` silently drops legitimate deeper repeats (e.g., `[2, 2]`).
- **Confusing combinations and permutations.** `range(n)` with `used` when order doesn't matter produces each combination k! times.
- **Python recursion limit** (default 1000). Fine for typical depths; for long grid paths, use `sys.setrecursionlimit` or go iterative.
- **Backtracking a counting question with large n.** n = 100 and "how many ways" means DP.

## When NOT to use / look-alikes

| Looks like backtracking | But is | Tell |
|---|---|---|
| Count ways to make change / climb / decode | DP | Asks for a number, not the list; n large; state is an index or amount |
| Shortest path in a grid | BFS (module 13) | Wants minimum steps, not all paths; cells visited once globally |
| Is there a path / connected component | DFS/BFS with permanent visited | Reusing a cell across paths is irrelevant |
| Longest path with "no revisit" in a general grid | Backtracking (NP-hard in general) | Unless edges only go "uphill" (a DAG): then memoized DFS (module 17) |
| Pick items to maximize value under a budget | Knapsack DP | Optimization with small integer capacity |
| Choose locally best each step | Greedy (module 18) | Exchange argument shows the local choice is safe |
| All subsets but n = 40 | Meet in the middle | Two halves of 2^20 each, combined with sort/hash |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Subsets (78)](https://leetcode.com/problems/subsets/) · Medium · Include/exclude: both styles; record at every node.
2. [Permutations (46)](https://leetcode.com/problems/permutations/) · Medium · used[] array: mark, recurse, unmark.
3. [Combination Sum (39)](https://leetcode.com/problems/combination-sum/) · Medium · Start index, reuse allowed: recurse with `i` for reuse; sorted break.
4. [Word Search (79)](https://leetcode.com/problems/word-search/) · Medium · Grid backtracking: in-place mark and restore.
5. [N-Queens (51)](https://leetcode.com/problems/n-queens/) · Hard · Constraint satisfaction: diagonal keys `r - c`, `r + c`.

### Optional

**Subsets and combinations**
- [Subsets II (90)](https://leetcode.com/problems/subsets-ii/) · Medium · Sort + `i > start` skip.
- [Combinations (77)](https://leetcode.com/problems/combinations/) · Medium · Bound prune on remaining items.
- [Combination Sum II (40)](https://leetcode.com/problems/combination-sum-ii/) · Medium · `i + 1` plus duplicate skip.
- [Combination Sum III (216)](https://leetcode.com/problems/combination-sum-iii/) · Medium · Fixed size k and fixed sum; double prune.

**Permutations**
- [Permutations II (47)](https://leetcode.com/problems/permutations-ii/) · Medium · Canonical order of copies.
- [Letter Case Permutation (784)](https://leetcode.com/problems/letter-case-permutation/) · Medium · Binary branching only at letters.
- [Beautiful Arrangement (526)](https://leetcode.com/problems/beautiful-arrangement/) · Medium · Permutation with a positional constraint prune.

**Constraint satisfaction and grids**
- [N-Queens II (52)](https://leetcode.com/problems/n-queens-ii/) · Hard · Same diagonal keys as N-Queens; count instead of building boards.
- [Sudoku Solver (37)](https://leetcode.com/problems/sudoku-solver/) · Hard · Return True to keep the solution.
- [Word Search II (212)](https://leetcode.com/problems/word-search-ii/) · Hard · Trie to prune many words at once.
- [Unique Paths III (980)](https://leetcode.com/problems/unique-paths-iii/) · Hard · Visit every empty cell exactly once; count cells remaining.

**Partitioning and generation**
- [Palindrome Partitioning (131)](https://leetcode.com/problems/palindrome-partitioning/) · Medium · Precomputed palindrome table.
- [Restore IP Addresses (93)](https://leetcode.com/problems/restore-ip-addresses/) · Medium · Length bounds prune.
- [Generate Parentheses (22)](https://leetcode.com/problems/generate-parentheses/) · Medium · Validity built into the move rules.
- [Letter Combinations of a Phone Number (17)](https://leetcode.com/problems/letter-combinations-of-a-phone-number/) · Medium · Cartesian product.
- [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/) · Medium · Fill buckets; sort descending; skip identical empty buckets.
- [Word Break II (140)](https://leetcode.com/problems/word-break-ii/) · Hard · Memoize on start index.
- [Expression Add Operators (282)](https://leetcode.com/problems/expression-add-operators/) · Hard · Track last operand to undo multiplication precedence.

## Self-check

1. Draw the loop-start decision tree for subsets of `[1, 2, 3]`. How many nodes, and why is every node (not just leaves) an answer?
2. Why does `i > start` (not `i > 0`) in Subsets II still allow `[2, 2]`?
3. In Combination Sum, what goes wrong recursing with `0` instead of `i`? With `i + 1`?
4. Why do `r - c` and `r + c` suffice to identify diagonals in N-Queens?
5. In Word Search, why can't you use a global `visited` set that's never cleared?
6. Estimate the running time of permutations of 10 items and subsets of 25 items. Which fits an interview time limit?
7. Given a backtracking solution, how do you decide whether memoization helps, and what's the memo key?
