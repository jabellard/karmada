# 23. Practice Plan, Review System, and Interview Day

Knowing the patterns isn't enough; you need **fast, reliable recall under pressure**. This module is the training system: schedule, problem protocol, miss review, mock interviews, and the final week.

---

## 1. Your target

By the end you should:
- **Recognize** an unseen medium's pattern within 3–5 minutes.
- **Write any pattern's core template** from a blank editor in under 5 minutes.
- **Solve an unseen medium** cleanly in 20–25 minutes, with tests and complexity.
- **Make solid progress on an unseen hard** in 40 minutes (correct approach, mostly coded).
- **Narrate** throughout without going silent.

Track these weekly as readiness metrics.

---

## 2. Schedule (6 weeks, ~1–1.5 hours/day; compress to 4 if needed)

**Doing the full course in 4 weeks?** Follow the day-by-day [`study-plan.md`](study-plan.md) instead of the table below.

| Week | Modules | Daily focus |
|---|---|---|
| 1 | 01, 02, 03, 04, 05 | Framework + recognition; arrays/hashing/prefix sums, two pointers. 3 problems/day |
| 2 | 06, 07, 08, 09 | Sliding window, binary search (including on the answer), monotonic stack, linked lists |
| 3 | 10, 11, 12 | Trees, heaps, intervals. Add 2 mixed problems/day from earlier weeks |
| 4 | 13, 14, 15 | Graphs (BFS/DFS/topo), union-find/Dijkstra/MST, backtracking |
| 5 | 16, 17, 18 | DP foundations + advanced, greedy. Heavier review of weak patterns |
| 6 | 19, 20, 21, 22 | Tries/strings, bits/math, practical coding, **disguised drills**. 2 mock interviews |

**Fast track (2 weeks):** 01, 02, 03 → 04, 06, 07, 13, 16 → 05, 08, 10, 11, 12 → 21, 22. Solve only the **★ top 35** in `practice-tracker.md`. Skip advanced DP and string algorithms unless targeting algorithm-heavy companies.

**Priority if time is short** (by interview frequency):
1. Arrays/hashing, two pointers, sliding window
2. Trees (DFS/BFS)
3. Graphs (BFS/DFS, topo sort, grid)
4. Binary search (including on the answer)
5. Heaps / top-k / intervals
6. DP foundations
7. Monotonic stack, linked lists, backtracking
8. Practical coding (essential for infra / AI-infra companies)
9. Union-find, Dijkstra, tries, greedy
10. Advanced DP, string algorithms, bits/math

---

## 3. Daily routine (60–90 min)

1. **Warm-up (10 min):** write one template from memory (rotating patterns) and check it against the module's **Template cheat sheet**.
2. **Learn (15–20 min, new-module days):** read the module's At a glance, Core idea, Signals, and Disguises. Type each template (don't copy-paste) and attempt its worked examples before opening the solutions.
3. **Solve (40–50 min):** 2–3 problems with the protocol below, **must-dos first** from `practice-tracker.md` (5 per pattern). Optional problems come only after all must-dos, focused on weak patterns.
4. **Review (10 min):** update the error log and schedule re-attempts.

---

## 4. Problem protocol (every problem)

1. **Read and classify (≤ 5 min):** write the one-line recognition note:
   `Signal: ___ → Pattern: ___ → Template: ___`
2. **Plan (≤ 5 min):** invariant, data structures, complexity.
3. **Code (≤ 15 min):** from the template. Talk aloud, even alone, to build the habit.
4. **Test:** trace one example by hand, then the edge cases.
5. **Timebox:** stuck for 20–25 min → read *one* hint (the module's disguise/decoding notes first, then an editorial). After reading a solution, **close it and re-implement from scratch**.
6. **Post-mortem (2 min):** if you missed the pattern, write down the missed signal.

**Don't** grind volume without review: 15 problems with full post-mortems beat 50 rushed ones.

---

## 5. Error log (the most valuable artifact you'll build)

Keep a simple table (spreadsheet or markdown):

| Date | Problem | Pattern | Result | Missed signal / bug | Re-attempt dates |
|---|---|---|---|---|---|
| 10/12 | Fruit into baskets | Sliding window (≤ 2 distinct) | Missed pattern | "pick from consecutive trees" = contiguous | 10/15, 10/19 |
| 10/13 | Course schedule II | Topo sort | Bug | no cycle check (result length < n) | 10/16 |

**Categories of misses:**
- **Recognition:** didn't see the pattern → add the signal to your personal trigger list.
- **Template:** knew the pattern, fumbled the code → rewrite the template from memory daily for 3 days.
- **Edge case / bug:** add it to that module's pitfall list in your notes.
- **Speed:** right answer, too slow → re-solve timed.

---

## 6. Spaced repetition

- Re-attempt every missed problem blind after **3 days**, then **7 days**, then **21 days**.
- One flashcard per pattern. Front: recognition signals + a disguised example. Back: template + invariant + pitfalls.
- Weekly 30-minute "pattern sprint": 10 random problem statements from mixed modules, 2 minutes each, **classify only** (no coding). This trains recognition directly.

---

## 7. Mock interviews

Start in week 4; do at least 4–6 before real loops.
- **Format:** 45 minutes, one medium + a follow-up, or one hard. A partner plays interviewer (peers, fellow maintainers, or a paid mock platform).
- **Conditions:** shared plain editor or whiteboard tool, no autocomplete, camera on, talking throughout.
- **Record yourself** at least once; watch for silent stretches, skipped clarifications, and untested code.
- Do at least one multi-part **practical-coding mock** (module 21) that extends code.

### Self-evaluation rubric (score 1–4 after each mock)

| Dimension | 4 = strong hire signal |
|---|---|
| Problem understanding | Clarified ambiguities, made examples, found edge cases early |
| Approach | Stated brute force, found the waste, picked the right pattern, explained tradeoffs |
| Communication | Narrated continuously, checked in, took hints well |
| Code quality | Clean structure, good names, helpers, idiomatic Python |
| Correctness and testing | Traced examples, handled edge cases, fixed bugs calmly |
| Complexity | Justified time/space; discussed optimizations |
| Speed | Finished with time for follow-ups |

Anything scoring ≤ 2 becomes next week's focus.

---

## 8. Python fluency drills

Remove friction for timed rounds:
- Write each from memory in under 2 minutes: grid BFS, Dijkstra, the DSU class, lower_bound binary search, the sliding-window skeleton, a backtracking skeleton, a memoized DP with `@cache`, Kahn's topological sort, heap top-k.
- Know `collections` (`Counter`, `defaultdict`, `deque`, `OrderedDict`), `heapq` (`heappush`, `heappop`, `heapify`, `nlargest`), `bisect` (`bisect_left`, `bisect_right`), `itertools`, `functools.cache`, `math.inf`. See module 03.
- Practice in a plain editor without autocomplete at least half the time.

---

## 9. Final week taper

- **No new patterns.** Review templates, the error log, and the disguised drills (module 22).
- One mixed timed session daily (two unseen mediums in 50 minutes).
- Re-read `02-pattern-recognition.md` and your personal trigger list.
- Rest the day before. Sleep beats one more problem.

---

## 10. Interview day checklist

- [ ] Editor or whiteboard tool tested; Python version known; camera, mic, and light set up.
- [ ] Water and scratch paper.
- [ ] Open with: restate → clarify → examples (module 01).
- [ ] State the brute force and its complexity before optimizing.
- [ ] Announce the pattern and invariant; get a nod before coding.
- [ ] Narrate while coding; use helpers.
- [ ] Trace an example and the edge cases before the interviewer asks.
- [ ] State time and space complexity.
- [ ] Offer follow-ups (streaming, scale, concurrency), a staff-level signal; see `../interview-playbook.md`.
- [ ] If stuck: brute force → smaller example → relax a constraint → think aloud.

---

## 11. Course completion checklist

- [ ] Can explain the 8-step framework without notes (01)
- [ ] Recognition quiz in 02: 15/15
- [ ] Every template in modules 04–20 written from memory at least once
- [ ] All 86 must-dos in `practice-tracker.md` solved blind in under 25 min (misses redone after 7 days)
- [ ] Optional problems: 3+ extra in each of your weakest patterns
- [ ] All 36 disguised drills (22) attempted blind; misses redone after 3 and 7 days
- [ ] All practical-coding exercises (21) implemented from scratch, timed, at least once
- [ ] 4+ mock interviews scoring ≥ 3 in every rubric dimension
- [ ] Unseen medium solve rate ≥ 80% in under 25 minutes over the last 10 attempts
