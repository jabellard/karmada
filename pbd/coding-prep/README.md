# Coding Interview Prep: Patterns, Not Problems

A course for passing coding interviews at top tech and AI-infrastructure companies by **mastering patterns** instead of memorizing problems, and by learning to decode unseen problems, including ones worded to hide the pattern.

All code is Python 3.11+, stdlib only. Every template and worked solution lives in `code/`, covered by tests.

---

## Philosophy

1. **There are about 20 patterns.** Nearly every interview problem is one of them, or two combined. Learning each pattern's *signal*, *invariant*, and *template* solves whole classes of problems.
2. **Recognition is the real skill.** Most failing candidates know the technique but don't see that it applies. This course trains recognition directly: signals, disguises, look-alikes, and blind drills.
3. **A procedure beats inspiration.** A repeatable framework (understand → constraints → brute force → find the waste → map to a pattern → plan → code → verify) gets you unstuck on unfamiliar problems.
4. **Templates are starting points, not answers.** Each has marked hooks to adapt, not recite.
5. **Review misses by signal, not by problem.** Every miss becomes a reusable rule.

---

## Course map

| # | Module | What you'll master |
|---|---|---|
| 01 | [Problem-Solving Framework](01-problem-solving-framework.md) | The 8-step loop, finding the waste, breaking down unknown problems, getting unstuck |
| 02 | [Pattern Recognition](02-pattern-recognition.md) | Signals → patterns, decision flow, trigger phrases, **disguise techniques**, look-alikes |
| 03 | [Python Toolkit & Complexity](03-python-toolkit-and-complexity.md) | Constraints → complexity, Python data structures and idioms, pitfalls |
| 04 | [Arrays, Hashing & Prefix Sums](04-arrays-hashing-prefix-sums.md) | Counting, complement lookup, prefix sums + hash map, difference arrays, cyclic sort |
| 05 | [Two Pointers](05-two-pointers.md) | Opposite ends, read/write pointers, partitioning, k-sum |
| 06 | [Sliding Window](06-sliding-window.md) | Fixed/variable windows, at-most-k counting, monotonic deque |
| 06a | [Monotonic Deques (companion)](06a-monotonic-deques.md) | How max and min deques work: adding, removing, order of operations, duplicates; window max/min, two-deque windows, deques over DP and transformed values |
| 06b | [Beyond the Sliding Window (companion)](06b-beyond-sliding-window.md) | When the window breaks (negatives, exact/modular targets, products): prefix sums + hash map, Kadane, prefix deque/stack, sorted prefixes, Fenwick counting, binary search on the answer |
| 07 | [Binary Search](07-binary-search.md) | One universal template, rotated arrays, **binary search on the answer** |
| 08 | [Stacks & Monotonic Structures](08-stacks-and-monotonic-structures.md) | Parsing, next greater element, histogram, contribution technique |
| 08a | [Monotonic Stacks (companion)](08a-monotonic-stacks.md) | Who answers whom, the four next/previous greater/smaller variants, duplicates and strictness, sentinels, the contribution technique; histogram, spans, circular arrays, greedy digit removal |
| 09 | [Linked Lists](09-linked-lists.md) | Dummy head, fast/slow, reversal, merging |
| 10 | [Trees & BST](10-trees-and-bst.md) | Top-down vs bottom-up DFS, global-answer pattern, BFS levels, BST ops |
| 11 | [Heaps, Top-K & Merging](11-heaps-top-k-and-merging.md) | Top-k, two heaps, k-way merge, heap-driven scheduling |
| 12 | [Intervals & Sweep Line](12-intervals-and-sweep-line.md) | Merge, overlap counting, event sweep, interval greedy |
| 13 | [Graphs: Traversal & Topological Sort](13-graphs-traversal-and-topological-sort.md) | BFS/DFS, grids, multi-source BFS, state-space BFS, topo sort, cycles |
| 14 | [Union-Find, Shortest Paths & MST](14-union-find-shortest-paths-mst.md) | DSU, Dijkstra, Bellman-Ford with k stops, Kruskal/Prim |
| 15 | [Backtracking](15-backtracking.md) | Choose/explore/unchoose, subsets/permutations/combinations, pruning, constraint satisfaction |
| 16 | [Dynamic Programming: Foundations](16-dynamic-programming-foundations.md) | DP method (state → transition → base → order); 1D, grid, two-sequence, knapsack |
| 17 | [Dynamic Programming: Advanced](17-dynamic-programming-advanced.md) | LIS, interval DP, state machines, bitmask DP, tree DP, digit DP |
| 18 | [Greedy](18-greedy.md) | Exchange arguments, interval/jump/gas/scheduling greedy, when greedy fails |
| 19 | [Tries & String Algorithms](19-tries-and-string-algorithms.md) | Trie, palindromes, KMP, rolling hash |
| 20 | [Bit Manipulation & Math](20-bit-manipulation-and-math.md) | XOR tricks, bitmasks, modular arithmetic, sampling, matrix manipulation |
| 21 | [Practical Coding for Infra Roles](21-practical-coding-for-infra-roles.md) | LRU, TTL/transactional KV, rate limiters, worker pools, DAG schedulers, reconcile loops, consistent hashing |
| 22 | [Disguised Problem Drills](22-disguised-problem-drills.md) | 36 mixed, story-wrapped problems with hint ladders and full solutions |
| 23 | [Practice Plan & Review](23-practice-plan-and-review.md) | 6-week schedule, problem protocol, error log, mocks, rubric, interview-day checklist |
| — | [**Practice Tracker**](practice-tracker.md) | **5 must-dos per pattern** (86 total; ★ top 35 for a 2-week sprint), optional extras, progress tables |
| — | [**4-Week Study Plan**](study-plan.md) | **Day-by-day plan for the full course** (~120 h: 3 h weekdays, 7.5 h weekend days), with checkpoints and a readiness check |

---

## How each pattern module is structured (04–20)

1. **At a glance:** when to use it, one line per template on when to pick it, typical complexity, top trap.
2. **Core idea and invariant:** why the pattern works.
3. **Recognition signals:** what in the statement, input, and constraints points here.
4. **Disguises:** story-wrapped statements, each with a short decode of the giveaway.
5. **Templates and worked examples:** each template opens with a fitting LeetCode problem, then the template code, then **worked examples numbered by template** (Example 2.1 = Template 2's first example). The problem is shown first; the **Decode**, **Approach**, code, **Complexity**, and **Edge cases** sit in a collapsible **"Decode, approach, and solution"** block so you can try it first. Some modules end this part with **Combined examples** (C.k: problems that need two templates together) or **Beyond the templates** (B.k: useful techniques outside the module's templates, each with a one-line reason).
6. **Template cheat sheet:** every template's code on one screen, for warm-ups and review.
7. **Variations, pitfalls, look-alikes.**
8. **Practice ladder:** must-dos (matching `practice-tracker.md`) and optional problems, each linked to LeetCode.
9. **Self-check questions.**

---

## How to use this course

1. **Start with 01, 02, 03** (framework, recognition, toolkit); they speed up every later module.
2. **For each pattern module:** read At a glance → Core idea → Signals → Disguises. **Type** each template and attempt its worked examples before opening the solutions. Then do the pattern's **must-dos** from [`practice-tracker.md`](practice-tracker.md) with the module 23 protocol; optional problems come after all must-dos.
3. **Keep an error log** of missed *signals* (module 23).
4. **Finish with 21 and 22:** practical coding (essential for infra / AI-infra companies) and blind disguised drills.
5. **Follow the schedule** in module 23 (6 weeks, or the 2-week fast track).

---

## Running the code

```bash
cd coding-prep
python3 run_all_tests.py          # runs every module's tests
python3 code/06_sliding_window.py # run one module
```

Each `code/` file holds the module's templates and worked solutions, with asserts for normal and edge cases.

**Navigating a code file:**
- The **index** at the top lists every function and where it's taught, e.g. `length_of_longest_substring  Example 2.1 · LeetCode 3: Longest Substring Without Repeating Characters · uses Template 2`.
- A **`# ▸` label** above each function repeats that information.
- "Supporting code" functions are helpers, test utilities, or extra variations the module text mentions.
- After editing a module, regenerate the labels with `python3 tools/label_code.py`.

---

## Related

- `../interview-playbook.md`: the full interview process (recruiter screen → offer) and Staff-level signals.
