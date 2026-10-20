# 01. Problem-Solving Framework

Patterns help only once you know which one applies. Run this procedure on **every** problem, especially unfamiliar ones, until it's automatic; it keeps you calm, structured, and talking, which is half the evaluation.

---

## The 8-step loop

| Step | What you do | Time in a 45-min round |
|---|---|---|
| 1. Understand | Restate, clarify, examples, edge cases | 3–5 min |
| 2. Constraints → target | Estimate required complexity from input sizes | 1 min |
| 3. Brute force | State it and its cost; code it only if asked | 2–3 min |
| 4. Find the waste | What does brute force recompute or re-search? | 2–4 min |
| 5. Map to a pattern | Signals in `02-pattern-recognition.md` | 1–3 min |
| 6. Plan | Data structures, invariant, pseudocode; get a nod | 3–5 min |
| 7. Implement | Clean code from the template, with helpers | 12–18 min |
| 8. Verify | Trace, edge cases, complexity, follow-ups | 5–8 min |

Interviewers grade the process, not just the code: a clearly explained, correct plan that runs out of time often passes, while silently written perfect code often scores worse than expected.

---

## Step 1: Understand

**Restate the problem abstractly:** "So, given an array of integers and a target, return the number of contiguous subarrays whose sum equals the target."

**Clarifying questions** (only the ones that matter):
- **Input:** Sorted? Duplicates? Negatives or zeros? Empty? Size limits? Bounded values (1..n, lowercase letters)?
- **Output:** A value, indices, elements, or all solutions? Any valid answer or a specific one (lexicographically smallest)? What if none exists?
- **Mutation:** Can I modify the input? Must it be in place?
- **Scale:** Does it fit in memory? Is it a stream?
- **Ties:** How are ties broken?

**Make two examples:** a normal one and a tricky one (duplicates, all negative, single element). Solving them by hand often reveals the pattern.

---

## Step 2: Constraints → target complexity

Python manages about 10^7 simple operations per second, and limits are chosen so the intended complexity fits. Rule of thumb: n ≤ 10 → O(n!), ≤ 20 → O(2^n), ≤ 500 → O(n^3), ≤ 5,000 → O(n^2), ≤ 10^6 → O(n log n) or O(n), larger → O(log n) or O(1). At n = 10^5, O(n^2) is wrong; something smarter exists. The full table, with the usual pattern for each size, is in `03-python-toolkit-and-complexity.md` §2.

---

## Step 3: Brute force first

Say it out loud: "The naive approach checks every subarray: O(n^2) subarrays, O(n) to sum each, so O(n^3), or O(n^2) with running sums."

It proves understanding, gives you a fallback, and **exposes the repeated work** the pattern comes from.

---

## Step 4: Find the waste (the most important step)

Ask these of the brute force; each "yes" points to a technique family:

| Question | If yes → |
|---|---|
| **Recomputing sums/aggregates** over overlapping ranges? | Prefix sums, sliding window |
| Repeatedly **searching for a matching element**? | Hash map (complement lookup), sort + two pointers, binary search |
| **Re-solving the same subproblem** in recursion? | Memoization → DP |
| Repeatedly **finding the min/max** of a changing set? | Heap, monotonic deque, balanced structure |
| Scanning back for the **nearest greater/smaller** element? | Monotonic stack |
| Checking a candidate answer costs O(n), and **feasibility is monotonic** in it? | Binary search on the answer |
| Exploring **all orderings/choices**, most invalid early? | Backtracking with pruning |
| Would **sorting** expose the structure (input order irrelevant)? | Sort, then greedy / two pointers / sweep |
| Recomputing **connectivity** as edges arrive? | Union-Find |
| Matching **prefixes** of many strings? | Trie |

---

## Step 5: Map to a pattern

Combine three lenses (details in `02-pattern-recognition.md`):
1. **Input shape:** array, string, intervals, tree, graph, grid, stream.
2. **Question type:** longest/shortest, count, k-th, all solutions, min-cost, feasibility, ordering.
3. **Constraints:** n, value ranges, k small, sorted, non-negative.

If two patterns fit, pick the simplest one that meets the target complexity.

---

## Step 6: Plan before coding

Write a 4–8 line plan, in comments or out loud:
```
# pattern: prefix sums + hash map
# seen[prefix] = how many times this prefix sum occurred
# for each element: prefix += x; answer += seen[prefix - k]; seen[prefix] += 1
# init seen[0] = 1 (empty prefix)
# O(n) time, O(n) space
```
State the **invariant** ("the window always contains at most k distinct values"), then ask "Does this sound good before I code it?" A nod saves you from coding the wrong thing for 15 minutes.

---

## Step 7: Implement

- Start from the pattern's template (modules 04–20) and customize the hooks.
- Use clear names (`left`, `right`, `window_count`), not `i`, `j`, `k`.
- Extract helpers for anything non-trivial (`neighbors(r, c)`, `is_valid(state)`).
- Handle edge cases at the top (`if not nums: return 0`).
- Keep talking: "Now I expand the window... if it's invalid, I shrink from the left until it's valid again."

---

## Step 8: Verify

1. **Trace** a small example line by line, tracking variables.
2. **Edge cases:** empty, one element, all the same, all negative, already sorted, maximum size.
3. **Complexity:** time and space, justified.
4. **Follow-ups:** offer the next level unprompted: "If this were a stream, I'd keep only the last k elements in a deque."

---

## Techniques for breaking down an unfamiliar problem

When nothing clicks, these moves turn an unknown problem into a known one.

1. **Solve a small input by hand and watch yourself.** Your manual process often *is* the algorithm (scanning left to right remembering "the best so far" → greedy / running state).
2. **Simplify.** Drop a constraint (no negatives, k = 1, 1D instead of 2D), solve that, then add it back.
3. **Change the representation.**
   - Pairs/relations → a graph
   - Ranges → events on a line (sweep)
   - Running totals → prefix sums
   - Characters → counts (an array of 26)
   - States → nodes (BFS over states)
4. **Sort it.** If input order doesn't matter, sorting often reveals greedy, two-pointer, or sweep structure.
5. **Think in decisions.** What are the choices at each element (take/skip, left/right)? Choices + repeated subproblems → DP; no reuse → backtracking.
6. **Reverse the question.**
   - "Minimize the maximum" → "Is answer X feasible?" → binary search on X.
   - "Count subarrays with exactly K" → atMost(K) − atMost(K−1).
   - "At least one" → total − none.
7. **Work backwards** from the goal (DP from the end, reverse BFS from targets, multi-source BFS from all exits).
8. **Look at extremes.** All equal, sorted, n = 1? Extremes reveal invariants.
9. **Find the invariant.** What must hold for a partial solution to extend to a full one? That's your window condition, stack property, or DP state.
10. **Decompose** into independent parts (left/right halves, per row, per component).
11. **Precompute** for many fast queries (prefix sums, sparse tables, sorting, a trie).
12. **Bound the answer.** If it lies in a known range and candidates are testable, binary search it.

---

## Getting unstuck in the interview (a script)

1. "Let me go back to the brute force and see what work repeats." *(Step 4.)*
2. "Let me try a smaller example by hand." *(Technique 1.)*
3. "What if the array were sorted? / What if k were 1?" *(Techniques 2 and 4.)*
4. "I'm considering a sliding window, but negative numbers break the shrink condition, so prefix sums fit better." *(Thinking aloud lets the interviewer nudge you.)*
5. **Take hints gracefully** and run with them; using hints well is a positive signal.

Never go silent for more than ~30 seconds; narrate even "I'm checking whether a greedy choice works here."

---

## Worked demonstration 1: the framework end to end

**Problem (unfamiliar wording):** A service logs the change in active connections each minute (it can be negative). An incident review asks: how many stretches of consecutive minutes had a net change of exactly `target`?

1. **Understand:** count contiguous subarrays with sum == target in an int array with negatives. Example: `[1, -1, 1, 1]`, target 1 → `[1]`, `[1,-1,1]`, `[1]` (index 2), `[1]` (index 3), `[-1,1,1]` → 5.
2. **Constraints:** n up to 10^5 → O(n) or O(n log n).
3. **Brute force:** all (i, j) pairs with a running sum: O(n^2). Too slow.
4. **Waste:** for each end j, I search for starts i with `P[j+1] - P[i] == target`, i.e., earlier prefixes equal to `P[j+1] - target`: a **repeated lookup** → hash map.
5. **Pattern:** prefix sums + hash map (module 04). Not a sliding window: negatives break monotonic shrinking.
6. **Plan:** keep `count[prefix]`, initialized with `count[0] = 1`.
7. **Implement:**

```python
from collections import defaultdict

def count_windows_with_net_change(changes: list[int], target: int) -> int:
    seen = defaultdict(int)
    seen[0] = 1                    # empty prefix: lets subarrays start at index 0
    prefix = answer = 0
    for x in changes:
        prefix += x
        answer += seen[prefix - target]   # earlier prefixes that make sum == target
        seen[prefix] += 1
    return answer
```

8. **Verify:** trace `[1,-1,1,1]`, target 1 → 5. Edge cases: empty → 0; target 0 with `[0,0]` → 3. **Complexity:** O(n) time, O(n) space.

---

## Worked demonstration 2: reversing the question

**Problem:** You must run `jobs` (durations, in a fixed order) on `k` identical machines. Each machine takes a **contiguous** block of jobs. Minimize the finishing time of the busiest machine.

1. **Understand:** split the array into k contiguous parts, minimizing the maximum part sum.
2. **Constraints:** n up to 10^5, durations up to 10^4.
3. **Brute force:** all split points is exponential; DP is O(k · n^2), still too slow.
4. **Waste / reverse the question:** "minimize the maximum" → "Can max load ≤ X?" A greedy O(n) check answers it (fill each machine until it would exceed X), and feasibility is **monotonic**: if X works, any larger X does.
5. **Pattern:** binary search on the answer (module 07) + greedy check (module 18).
6. **Plan:** find the smallest feasible X in `[max(jobs), sum(jobs)]`.

```python
def min_max_load(jobs: list[int], k: int) -> int:
    def feasible(limit: int) -> bool:
        machines, load = 1, 0
        for d in jobs:
            if load + d > limit:       # start a new machine
                machines += 1
                load = 0
            load += d
        return machines <= k

    lo, hi = max(jobs), sum(jobs)       # answer lies in [lo, hi]
    while lo < hi:                      # find the first feasible limit
        mid = (lo + hi) // 2
        if feasible(mid):
            hi = mid
        else:
            lo = mid + 1
    return lo
```

**Complexity:** O(n log(sum)) time, O(1) space. *(Known as "Split Array Largest Sum.")*

---

## Communication template (use every time)

1. "Let me restate the problem…" → clarifying questions → examples.
2. "The constraints suggest we need about O(n log n)."
3. "Brute force would be ___, which costs ___."
4. "The repeated work is ___, so I'll use ___."
5. "Here's the plan… does that sound reasonable?"
6. *(code, narrating)*
7. "Let me trace through an example… edge cases… complexity is ___."
8. "If we needed ___, I'd change ___."

---

## Self-check

1. What are the 8 steps, and roughly how long does each take?
2. For n = 2 × 10^5, which complexities are acceptable?
3. Name four "find the waste" questions and the technique each points to.
4. How do you reframe "minimize the maximum ___"?
5. What do you do when stuck and the interviewer is silent?
6. Why state the brute force even if you know the optimal solution?

Companion code: `code/01_framework_examples.py`.
