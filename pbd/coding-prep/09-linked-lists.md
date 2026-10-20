# 09. Linked Lists

Linked-list problems test pointer handling, not algorithms: don't lose nodes, create cycles, or trip on empty/one-node inputs. A few moves (dummy head, fast/slow, reversal, merge) compose into nearly every list problem, and reappear in LRU caches, free lists, and intrusive queues.

Companion code: `code/09_linked_lists.py` (templates and solutions, with asserts).

## At a glance

- **Use when:** a `ListNode` head plus "O(1) extra space", or an array/function giving each element a single "next".
- **Templates:**
  - Template 1: dummy + tail builder: producing a new list
  - Template 2: dummy head for deletion: the head may be removed
  - Template 3: iterative reversal: reverse/reorder/palindrome
  - Template 4: fast/slow: middle, cycle, k-th from end
- **Typical complexity:** O(n) time, O(1) extra space.
- **Top trap:** overwriting `x.next` before saving it, or a relinked tail still pointing into the old list (cycle).

## Core idea

A singly linked list only reaches the next node. Everything else builds on four primitives:

| Primitive | What it buys you | Invariant |
|---|---|---|
| Dummy (sentinel) head | Head deletion/insertion isn't a special case | `dummy.next` is the result's true head |
| Dummy + tail builder | O(1) append while building a list | `tail` is the result's last node so far |
| Fast/slow pointers | Middle, cycle, k-th from end in one pass | `fast` has taken 2x (or k more) steps than `slow` |
| Iterative reversal | Reverse a list or segment in O(1) space | `prev` heads the reversed prefix; `cur` heads the untouched suffix |

Pointer discipline: **save what you still need from `x.next` before overwriting it.**

```
reversal step:   prev    cur -> nxt -> ...
                 1) nxt = cur.next      (save)
                 2) cur.next = prev     (flip)
                 3) prev, cur = cur, nxt (advance)

 None <- 1 <- 2    3 -> 4 -> None
         prev     cur
```

Draw 3 nodes as boxes and arrows before coding anything non-trivial; it catches most bugs.

## Recognition signals

- A `ListNode` head; an "O(1) extra space" follow-up forces pointer techniques over copying to an array.
- "Remove / insert / partition" where the head might change: dummy head.
- "Middle", "cycle", "k-th from the end", "intersection", "in one pass": two pointers with a speed or distance gap.
- "Reverse", "reorder", "palindrome", "alternate from both ends": find middle + reverse second half.
- "Merge sorted", "sort a list in O(n log n)": merge two, recursively (merge sort) or via a heap (k lists).
- "Deep copy" with extra pointers (random, child): old->new hash map, or interleaved copies.
- Numbers as digit lists: carry propagation, dummy+tail builder.
- n+1 values in [1, n], find the duplicate read-only: implicit list `i -> nums[i]`, Floyd's cycle start.
- O(1) get/put with recency order: hash map + doubly linked list (LRU, `21-practical-coding-for-infra-roles.md`).

## Disguises

1. **Playlist shuffle.** "A music app stores a playlist as a chain of song records. Rearrange it in place, without a new playlist, as first, last, second, second-to-last, and so on."
   Decode: Reorder List: middle + reverse second half + weave.

2. **Corrupted job chain.** "Each job record points to the next job. A bug may make some job point back to an earlier one, looping the scheduler forever. Detect this and report the first re-run job, in constant memory."
   Decode: "points back" = cycle, "first re-run job" = cycle entry, constant memory rules out a visited set: Floyd.

3. **Duplicate ID in a mapping table.** "An array of length n+1 holds IDs in 1..n, each slot pointing to another slot. Exactly one ID repeats. Find it read-only in O(1) space." (LeetCode: Find the Duplicate Number)
   Decode: treat `i -> nums[i]` as `next`. Two slots pointing at one index = two edges into one node = the cycle entry, which Floyd phase 2 returns.

4. **Batch reversal for packet framing.** "Packets arrive as a linked chain. Send every consecutive block of k reversed; a final partial block goes as-is."
   Decode: Reverse Nodes in k-Group.

5. **Snapshot a graph-ish config.** "Each config node has a `next` link and an optional `override` link to any node in the chain. Produce an independent deep copy."
   Decode: Copy List with Random Pointer: map old->new, or interleave for O(1) extra.

6. **Big-integer counters.** "Two odometers store readings as digit chains, least significant first. Produce the combined reading."
   Decode: Add Two Numbers; `or carry` in the loop condition emits the final overflow digit.

## Templates and worked examples

### Template 1: dummy + tail builder

> - **Example problem:** [Add Two Numbers (2)](https://leetcode.com/problems/add-two-numbers/)
> - **Why this template:** "return the sum as a linked list" builds a new list node by node: dummy + tail.
> - **Worked below:** Example 1.2 (Merge Two Sorted Lists (21) in Example 1.1)
> - **Also fits:** [Merge Two Sorted Lists (21)](https://leetcode.com/problems/merge-two-sorted-lists/)

Use whenever you produce a new list (merge, partition, filter, add numbers).

```python
def template_dummy_tail(head: Optional[ListNode]) -> Optional[ListNode]:
    dummy = tail = ListNode()
    cur = head
    while cur:
        nxt = cur.next
        if True:  # TODO: condition for keeping `cur`
            tail.next = cur
            tail = cur
        cur = nxt
    tail.next = None  # TODO: always terminate; the old tail may still point into the input
    return dummy.next
```

Knobs:
- Append an existing node (relinking, O(1) space) or a new `ListNode(v)` (building).
- Multiple builders (`small`/`big` in partition), then stitch.
- **Terminate the final tail** when relinking existing nodes, or you create cycles.

#### Example 1.1: Merge two sorted lists into one (dummy + tail relinking) [Easy]

**Problem:** [Merge Two Sorted Lists (21)](https://leetcode.com/problems/merge-two-sorted-lists/) · **Template:** Template 1 (dummy + tail builder)

**Problem.** Given the heads of two non-decreasing linked lists, splice their nodes into one sorted list and return its head. Example: `1 -> 2 -> 4` and `1 -> 3 -> 4` give `1 -> 1 -> 2 -> 3 -> 4 -> 4`; an empty list and `0` give `0`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a new list from existing nodes = dummy + tail. Sorted inputs: compare only the two current heads (the k = 2 case of module 11's k-way merge).

**Approach:** `tail` ends the sorted merged prefix of consumed nodes. While both lists are non-empty, the smaller head is the smallest remaining node: append it. When one runs out, the other's nodes are all ≥ `tail.val`: attach the remainder at once. No `tail.next = None` needed: the last attached node ends an input list.

```python
def merge_two(a: Optional[ListNode], b: Optional[ListNode]) -> Optional[ListNode]:
    dummy = tail = ListNode()
    while a and b:
        if a.val <= b.val:  # <= keeps it stable
            tail.next, a = a, a.next
        else:
            tail.next, b = b, b.next
        tail = tail.next
    tail.next = a or b
    return dummy.next
```

**Complexity:** O(n + m) time; O(1) extra space, since nodes are relinked.

**Edge cases:** one or both lists empty (`a or b` returns the other or `None`); one list entirely smaller; duplicates across lists (`<=` takes from `a` first, so the merge is stable). This is `sort_list`'s merge step.

</details>

#### Example 1.2: Add two numbers stored as lists (dummy + tail with carry) [Medium]

**Problem:** [Add Two Numbers (2)](https://leetcode.com/problems/add-two-numbers/) · **Template:** Template 1 (dummy + tail builder)

**Problem.** Two non-negative integers are stored as linked lists, one digit per node, least significant first. Return their sum in the same format.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** building a new list node by node: dummy + tail (Template 1). Least significant first lets you add in list order, like grade-school addition.

**Approach:** walk both lists, adding the two digits plus `carry`. Append `sum % 10`; carry `sum // 10`. Continue while either list has nodes *or* a carry remains.

```python
def add_two_numbers(a: Optional[ListNode], b: Optional[ListNode]) -> Optional[ListNode]:
    dummy = tail = ListNode()
    carry = 0
    while a or b or carry:  # `or carry` handles the final extra digit
        s = carry + (a.val if a else 0) + (b.val if b else 0)
        carry, digit = divmod(s, 10)
        tail.next = tail = ListNode(digit)
        a = a.next if a else None
        b = b.next if b else None
    return dummy.next
```

**Complexity:** O(max(m, n)) time; O(max(m, n)) for the output list, O(1) extra.

**Edge cases:**
- A final carry adds a digit (`5 + 5 → 0 → 1`); `or carry` in the loop condition handles it.
- Different lengths: treat a missing node as 0.
- `tail.next = tail = node` assigns left to right, so the old tail's `next` is set before `tail` is rebound. Use two lines if it feels too clever.
- Follow-up (Add Two Numbers II, most significant first): push digits onto two stacks, pop to add, and **prepend** each result digit (`head = ListNode(digit, head)`), avoiding a final reversal (`add_two_numbers_msd`).

</details>

### Template 2: dummy head for deletion

> - **Example problem:** [Remove Duplicates from Sorted List II (82)](https://leetcode.com/problems/remove-duplicates-from-sorted-list-ii/)
> - **Why this template:** deleting duplicated values may delete the head, so examine `prev.next` from a dummy.
> - **Worked below:** Example 2.1

```python
dummy = ListNode(0, head)
prev = dummy
while prev.next:
    if should_delete(prev.next):     # TODO
        prev.next = prev.next.next   # do NOT advance prev; the new prev.next is unchecked
    else:
        prev = prev.next
return dummy.next
```

Knob: examine `prev.next`, not `cur`; deletion needs the predecessor.

#### Example 2.1: Remove all duplicated values from a sorted list [Medium]

**Problem:** [Remove Duplicates from Sorted List II (82)](https://leetcode.com/problems/remove-duplicates-from-sorted-list-ii/) · **Template:** Template 2 (dummy head for deletion)

**Problem.** Given a sorted list, delete every node whose value appears more than once, keeping only unique values.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** deletion where the head may go (`1,1,2`): dummy head.

**Approach:** `prev` is the last known survivor. If `cur = prev.next` starts a run of equal values, skip `prev.next` past the run without advancing `prev` (another run may follow). Otherwise `cur` is unique: advance `prev`.

```python
def delete_duplicates_all(head: Optional[ListNode]) -> Optional[ListNode]:
    dummy = ListNode(0, head)
    prev = dummy  # last node known to be kept
    while prev.next:
        cur = prev.next
        if cur.next and cur.next.val == cur.val:
            v = cur.val
            while cur and cur.val == v:
                cur = cur.next
            prev.next = cur
        else:
            prev = cur
    return dummy.next
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** all duplicates (`1,1`) → empty; duplicates at the head; consecutive runs (`1,1,2,2,3`).

Related (code file): `partition_list` (Partition List) builds two lists and must set `b_tail.next = None`.

</details>

### Template 3: iterative reversal

> - **Example problem:** [Reverse Linked List II (92)](https://leetcode.com/problems/reverse-linked-list-ii/)
> - **Why this template:** reversing positions left to right is in-place reversal of a segment.
> - **Worked below:** Example 3.2 (Reverse Linked List (206) in Example 3.1)
> - **Also fits:** [Reverse Linked List (206)](https://leetcode.com/problems/reverse-linked-list/)

```python
def template_reverse(head: Optional[ListNode]) -> Optional[ListNode]:
    prev, cur = None, head
    while cur:
        nxt = cur.next      # 1. save
        cur.next = prev     # 2. flip
        prev, cur = cur, nxt  # 3. advance both
    return prev
```

Knobs:
- Initial `prev`: `None` for a whole list; for a segment, the node *after* it, so the reversed block is already linked forward (k-group).
- Segment stop condition: `cur is not stop_node`.
- **Head insertion** (Reverse Linked List II): keep `before` fixed and repeatedly move `start.next` right after it. Fewer pointers to reconnect.

#### Example 3.1: Reverse a whole list (iterative, plus recursive) [Easy]

**Problem:** [Reverse Linked List (206)](https://leetcode.com/problems/reverse-linked-list/) · **Template:** Template 3 (iterative reversal)

**Problem.** Reverse every link of a singly linked list and return the new head. Example: `1 -> 2 -> 3 -> 4` becomes `4 -> 3 -> 2 -> 1`; an empty list stays empty.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** Template 3 with `prev = None` over the whole list. Expect to give both iterative and recursive versions.

**Approach:** invariant: `prev` heads the reversed prefix, `cur` the untouched suffix. Each step saves `cur.next` (or the suffix is lost), points `cur` at `prev`, and advances both. When `cur` is `None`, `prev` (the old tail) is the new head.

```
None <- 1    2 -> 3 -> None     prev = 1, cur = 2
None <- 1 <- 2    3 -> None     prev = 2, cur = 3
None <- 1 <- 2 <- 3             prev = 3, cur = None  -> return 3
```

```python
def reverse_list(head: Optional[ListNode]) -> Optional[ListNode]:
    prev, cur = None, head
    while cur:
        nxt = cur.next        # save the rest before breaking the link
        cur.next = prev       # flip this node to point backward
        prev, cur = cur, nxt  # grow the reversed prefix by one
    return prev
```

Recursive version: reverse everything after `head`. `head.next` still points at the reversed sublist's last node, so `head.next.next = head` appends `head`; `head.next = None` makes it the tail.

```python
def reverse_list_recursive(head: Optional[ListNode]) -> Optional[ListNode]:
    if not head or not head.next:
        return head
    new_head = reverse_list_recursive(head.next)
    head.next.next = head
    head.next = None
    return new_head
```

**Complexity:** O(n) time both ways. Iterative O(1) space; recursive O(n) stack, which can hit Python's ~1000 recursion limit, so prefer iterative.

**Edge cases:** empty and single-node lists return unchanged; not saving `nxt` loses the rest; recursively, forgetting `head.next = None` leaves a two-node cycle.

</details>

#### Example 3.2: Reverse a segment of a list [Medium]

**Problem:** [Reverse Linked List II (92)](https://leetcode.com/problems/reverse-linked-list-ii/) · **Template:** Template 3 (iterative reversal), head-insertion knob

**Problem.** Reverse only the nodes at positions `left..right` (1-indexed) in one pass.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** partial reversal where `left` may be 1: dummy head.

**Approach:** walk `before` to position `left - 1`; `start = before.next` ends up last in the block. `right - left` times, unlink `start.next` and insert it after `before`. `start` stays linked to the rest.

```
before -> 2 -> 3 -> 4 -> 5     (start = 2)
before -> 3 -> 2 -> 4 -> 5     moved 3
before -> 4 -> 3 -> 2 -> 5     moved 4
```

```python
def reverse_between(head: Optional[ListNode], left: int, right: int) -> Optional[ListNode]:
    dummy = ListNode(0, head)
    before = dummy
    for _ in range(left - 1):
        before = before.next
    start = before.next  # will end up as the last node of the reversed block
    for _ in range(right - left):
        moved = start.next
        start.next = moved.next
        moved.next = before.next
        before.next = moved
    return dummy.next
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** `left == right` (loop runs 0 times); `left == 1` (dummy handles it); whole list.

</details>

#### Example 3.3: Reverse every block of k nodes (segment reversal in a loop) [Hard]

**Problem:** [Reverse Nodes in k-Group (25)](https://leetcode.com/problems/reverse-nodes-in-k-group/) · **Template:** Template 3 (iterative reversal), with `prev` seeded to the node after the block

**Problem.** Reverse each consecutive block of k nodes; a final shorter block stays in order. O(1) extra space.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** repeated segment reversal with a "k left?" check.

**Approach:** `group_prev` precedes the current group (starts at dummy). Walk k steps to `kth`; stop if you fall off. Reverse with `prev` seeded as `group_next = kth.next`, so the block's new tail already points forward. Set `group_prev.next = kth` and move `group_prev` to the old first node (now the block's tail).

```python
def reverse_k_group(head: Optional[ListNode], k: int) -> Optional[ListNode]:
    dummy = ListNode(0, head)
    group_prev = dummy
    while True:
        kth = group_prev
        for _ in range(k):
            kth = kth.next
            if not kth:
                return dummy.next
        group_next = kth.next
        prev, cur = group_next, group_prev.next
        while cur is not group_next:
            nxt = cur.next
            cur.next = prev
            prev, cur = cur, nxt
        first = group_prev.next  # old first, now last of the block
        group_prev.next = kth
        group_prev = first
```

**Complexity:** O(n) time (each node visited twice: count, reverse), O(1) space.

**Edge cases:** `k = 1` (no change); length a multiple of k; length < k.

</details>

### Template 4: fast/slow

> - **Example problem:** [Linked List Cycle II (142)](https://leetcode.com/problems/linked-list-cycle-ii/)
> - **Why this template:** the cycle's start node in O(1) memory is fast/slow (Floyd).
> - **Worked below:** Example 4.1
> - **Also fits:** [Middle of the Linked List (876)](https://leetcode.com/problems/middle-of-the-linked-list/)

```python
def template_fast_slow(head: Optional[ListNode]) -> Optional[ListNode]:
    slow = fast = head
    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next
    return slow
```

Knobs:

| Start | Odd length (1-2-3-4-5) | Even length (1-2-3-4) | Use when |
|---|---|---|---|
| `slow = fast = head` | 3 | 3 (second middle) | "return the middle node" |
| `slow, fast = head, head.next` | 3 | 2 (first middle) | splitting: left half must be non-empty and not everything (merge sort, reorder, palindrome) |

- Distance gap instead of speed gap: advance `lead` n steps, then move both by 1 (n-th from end).
- Cycle check: test `slow is fast` after moving.

### Floyd's cycle start: proof sketch

Let `a` = distance from head to cycle start, `c` = cycle length, and say they meet `b` steps into the cycle.

- slow walked `a + b`; fast walked `2(a + b)`.
- fast's extra `a + b` is whole laps: `a + b = k*c`.
- So `a = k*c - b`: `a` steps from the meeting point lands at `b + (k*c - b) = k*c`, the cycle start.
- A pointer from head also reaches the start after `a` steps, so moving both by 1, they meet there.

Why they meet at all: in the cycle the gap closes by 1 per step, so fast can't jump over slow, and slow is caught within its first lap: O(n).

#### Example 4.1: Find where a cycle begins (Floyd) [Medium]

**Problems:** [Linked List Cycle II (142)](https://leetcode.com/problems/linked-list-cycle-ii/), [Remove Nth Node From End of List (19)](https://leetcode.com/problems/remove-nth-node-from-end-of-list/) · **Template:** Template 4 (fast/slow); the n-th-from-end sibling uses its distance-gap knob

**Problem.** Return the node where a cycle begins, or `None` if none, in O(1) extra space.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "where the cycle begins" + O(1) space = Floyd phase 2.

**Approach:** phase 1: fast moves 2, slow 1, until they meet (or fast hits `None`: no cycle). Phase 2: reset one to head, move both by 1; they meet at the start (proof above).

```python
def detect_cycle(head: Optional[ListNode]) -> Optional[ListNode]:
    slow = fast = head
    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next
        if slow is fast:
            p = head
            while p is not slow:
                p = p.next
                slow = slow.next
            return p
    return None
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** cycle at the head (`pos = 0`; phase 2 returns immediately); single-node self-loop; compare with `is`, not `==` on values.

Distance-gap sibling (LeetCode: Remove Nth Node From End of List):

```python
def remove_nth_from_end(head: Optional[ListNode], n: int) -> Optional[ListNode]:
    dummy = ListNode(0, head)
    lead = trail = dummy
    for _ in range(n):
        lead = lead.next
    while lead.next:  # stop when lead is on the last node; trail is just before target
        lead = lead.next
        trail = trail.next
    trail.next = trail.next.next
    return dummy.next
```

Starting both at `dummy` makes removing the head (n = length) work without a special case.

</details>

### Combined examples

#### Example C.1: Reorder a list by weaving its ends (middle + reverse + weave) [Medium]

**Problems:** [Reorder List (143)](https://leetcode.com/problems/reorder-list/), [Palindrome Linked List (234)](https://leetcode.com/problems/palindrome-linked-list/) · **Template:** Template 4 (fast/slow) and Template 3 (iterative reversal)

**Problem.** Rearrange `L0, L1, ..., Ln` into `L0, Ln, L1, Ln-1, ...` in place.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** alternating from both ends, but you can't walk backward: reverse the back half so "from the end" becomes "from the front".

**Approach:**
1. Find the first middle (`fast = head.next`); on odd length the left half gets the extra node.
2. Cut after the middle; reverse the second half.
3. Weave; the second half is never longer, so loop on `second`.

```python
def reorder_list(head: Optional[ListNode]) -> None:
    if not head or not head.next:
        return
    slow, fast = head, head.next
    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next
    second = template_reverse(slow.next)
    slow.next = None
    first = head
    while second:
        n1, n2 = first.next, second.next
        first.next = second
        second.next = n1
        first, second = n1, n2
```

**Complexity:** O(n) time, O(1) space.

**Edge cases:** forgetting `slow.next = None` creates a cycle; lengths 0, 1, 2.

The same skeleton solves Palindrome Linked List in O(1) space: find middle, reverse second half, compare, then **reverse it back** (callers rarely expect mutated input; interviewers notice).

```python
def is_palindrome(head: Optional[ListNode]) -> bool:
    if not head or not head.next:
        return True
    slow, fast = head, head.next
    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next
    second = template_reverse(slow.next)
    p, q, ok = head, second, True
    while q:  # second half is never longer than the first
        if p.val != q.val:
            ok = False
            break
        p, q = p.next, q.next
    slow.next = template_reverse(second)  # restore
    return ok
```

</details>

### Beyond the templates

#### Example B.1: Deep-copy a list with random pointers (hash map vs interleaving) [Medium]

**Problem:** [Copy List with Random Pointer (138)](https://leetcode.com/problems/copy-list-with-random-pointer/) · **Template:** none of the four (hash map clone, or interleaving)

**Why no template:** It needs a node→copy hash map or interleaving, not one of the four list templates.

**Problem.** Each node has `next` and `random` (any node or `None`). Return a deep copy.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** deep copy where `random` may point to a node not yet copied.

**Approach:** both options build an old-to-new mapping.

**Approach A, hash map (say this first):** pass 1 maps each `old -> new`; pass 2 wires `new.next = map[old.next]`, `new.random = map[old.random]`. This clones any pointer graph (Clone Graph, `13-graphs-traversal-and-topological-sort.md`).

```python
def copy_random_list_map(head: Optional[RandomNode]) -> Optional[RandomNode]:
    if not head:
        return None
    clone: dict[RandomNode, RandomNode] = {}
    cur = head
    while cur:
        clone[cur] = RandomNode(cur.val)
        cur = cur.next
    cur = head
    while cur:
        clone[cur].next = clone.get(cur.next)
        clone[cur].random = clone.get(cur.random)
        cur = cur.next
    return clone[head]
```

**Approach B, interleaving (O(1) extra):** the list is the map. Insert each copy after its original (`A -> A' -> B -> B'`), so `X'.random = X.random.next`. Then unweave.

```python
def copy_random_list_interleave(head: Optional[RandomNode]) -> Optional[RandomNode]:
    if not head:
        return None
    cur = head
    while cur:  # 1. interleave copies
        cur.next = RandomNode(cur.val, cur.next)
        cur = cur.next.next
    cur = head
    while cur:  # 2. copy of X.random is X.random.next
        if cur.random:
            cur.next.random = cur.random.next
        cur = cur.next.next
    new_head = head.next
    cur = head
    while cur:  # 3. unweave, restoring the original
        copy = cur.next
        cur.next = copy.next
        copy.next = copy.next.next if copy.next else None
        cur = cur.next
    return new_head
```

**Complexity:** both O(n) time. A: O(n) extra; B: O(1) extra beyond the output.

**Edge cases:** self-pointing `random`; `random = None`; empty list. In B, set randoms in a separate pass *before* unweaving, which destroys `X.next = X'`.

</details>

## Template cheat sheet

All templates, for review.

### Template 1: dummy + tail builder
```python
def template_dummy_tail(head: Optional[ListNode]) -> Optional[ListNode]:
    dummy = tail = ListNode()
    cur = head
    while cur:
        nxt = cur.next
        if True:  # TODO: condition for keeping `cur`
            tail.next = cur
            tail = cur
        cur = nxt
    tail.next = None  # TODO: always terminate; the old tail may still point into the input
    return dummy.next
```
### Template 2: dummy head for deletion
```python
dummy = ListNode(0, head)
prev = dummy
while prev.next:
    if should_delete(prev.next):     # TODO
        prev.next = prev.next.next   # do NOT advance prev; the new prev.next is unchecked
    else:
        prev = prev.next
return dummy.next
```
### Template 3: iterative reversal
```python
def template_reverse(head: Optional[ListNode]) -> Optional[ListNode]:
    prev, cur = None, head
    while cur:
        nxt = cur.next      # 1. save
        cur.next = prev     # 2. flip
        prev, cur = cur, nxt  # 3. advance both
    return prev
```
### Template 4: fast/slow
```python
def template_fast_slow(head: Optional[ListNode]) -> Optional[ListNode]:
    slow = fast = head
    while fast and fast.next:
        slow = slow.next
        fast = fast.next.next
    return slow
```

## Variations and extensions

- **Merge k sorted lists**: heap of `(val, list_index, node)` (`11-heaps-top-k-and-merging.md`), or pairwise divide-and-conquer, O(N log k) without a heap.
- **Sort List in O(n log n)**: merge sort: split at the first middle, recurse, `merge_two`. O(log n) stack; bottom-up is O(1) but rarely required.

  ```python
  def sort_list(head: Optional[ListNode]) -> Optional[ListNode]:
      if not head or not head.next:
          return head
      slow, fast = head, head.next  # first middle so the left half is never the whole list
      while fast and fast.next:
          slow = slow.next
          fast = fast.next.next
      right = slow.next
      slow.next = None
      return merge_two(sort_list(head), sort_list(right))
  ```
- **Intersection of Two Linked Lists**: each pointer jumps to the other list's head at its end; both travel `a + c + b`, meeting at the intersection (or both at `None`).
- **Rotate List**: find length and tail, close the circle, cut at `length - k % length`.
- **LRU cache**: hash map from key to node in a **doubly** linked list with sentinel head and tail; move-to-front and evict-from-back are O(1) splices. Full code and variants (TTL, LFU) in `21-practical-coding-for-infra-roles.md`. `collections.OrderedDict` (`move_to_end`, `popitem(last=False)`) is fine if allowed; ask.
- **Flatten a multilevel doubly linked list**: DFS with an explicit stack, or recursion returning the flattened sublist's tail.
- **Implicit linked lists**: `i -> nums[i]` (Find the Duplicate Number), `n -> sum of squares of digits` (Happy Number). Any deterministic "next state" function on a finite set eventually cycles, so Floyd applies.

Follow-ups:
- O(1) extra space → swap the array copy or hash set for fast/slow or in-place reversal.
- Don't modify the input → reversal tricks need a restore step, or accept O(n) space.
- Doubly linked → maintain `prev` on every splice: four pointer writes per insert.
- Recursive reversal → fine, but mention that the O(n) stack hits Python's ~1000 limit on long lists.

## Common pitfalls

- Losing the rest of the list: overwriting `cur.next` before saving it.
- Creating a cycle: after relinking (partition, reorder, odd-even), the last node still points to its old successor. Set `tail.next = None`.
- Advancing `prev` after a deletion, leaving the new `prev.next` unchecked.
- Wrong middle: a second-middle split of 2 nodes leaves left = whole list, so merge sort recurses forever.
- `while fast.next and fast` instead of `while fast and fast.next`: the first check must guard the second.
- Comparing nodes by value (`==`) instead of identity (`is`); values repeat.
- Pushing `(val, node)` into `heapq`: ties compare `ListNode`s and raise `TypeError`. Add a counter or index.
- Silently mutating input (palindrome check). Restore it, or say why not.
- Not testing `head is None` and single-node lists first.

## When NOT to use / look-alikes

| Looks like | Actually | How to tell |
|---|---|---|
| "Find the duplicate" in an array | Implicit linked list cycle (Floyd) | Values in `[1, n]` index the array; O(1) space + read-only required |
| Cycle detection in a general graph | DFS with colors / topological sort | Multiple outgoing edges; Floyd needs out-degree 1 |
| Linked list with frequent random access | Convert to array | Repeated index access, space unconstrained: `vals = to_list(head)` in O(n) |
| "Reverse" a string/array | Two pointers swapping (module 05) | O(1) index access; no pointer surgery |
| LRU | Hash map + doubly linked list, not singly | O(1) removal of an arbitrary node needs its predecessor |
| Merge k lists | Heap (module 11) | k > 2 inputs; merging one by one is O(N k) |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Reverse Linked List (206)](https://leetcode.com/problems/reverse-linked-list/) · Easy · In-place reversal: iterative and recursive.
2. [Merge Two Sorted Lists (21)](https://leetcode.com/problems/merge-two-sorted-lists/) · Easy · Dummy head + tail: attach the leftover in one step.
3. [Linked List Cycle II (142)](https://leetcode.com/problems/linked-list-cycle-ii/) · Medium · Fast/slow, cycle start: Floyd phase 2; know the proof.
4. [Remove Nth Node From End of List (19)](https://leetcode.com/problems/remove-nth-node-from-end-of-list/) · Medium · Gap pointers + dummy: distance gap from a dummy.
5. [Reorder List (143)](https://leetcode.com/problems/reorder-list/) · Medium · Middle + reverse + weave: split at the middle, reverse the back half, interleave.

### Optional

**Dummy head / builder**
- [Remove Linked List Elements (203)](https://leetcode.com/problems/remove-linked-list-elements/) · Easy · Don't advance `prev` after deleting.
- [Partition List (86)](https://leetcode.com/problems/partition-list/) · Medium · Two builders; terminate the big list.
- [Remove Duplicates from Sorted List II (82)](https://leetcode.com/problems/remove-duplicates-from-sorted-list-ii/) · Medium · Skip whole runs; head can vanish.
- [Add Two Numbers (2)](https://leetcode.com/problems/add-two-numbers/) · Medium · Carry in the loop condition.
- [Odd Even Linked List (328)](https://leetcode.com/problems/odd-even-linked-list/) · Medium · Two interleaved builders; stitch odd tail to even head.

**Fast/slow**
- [Middle of the Linked List (876)](https://leetcode.com/problems/middle-of-the-linked-list/) · Easy · Which middle on even length?
- [Linked List Cycle (141)](https://leetcode.com/problems/linked-list-cycle/) · Easy · Speed gap, `is` comparison.
- [Intersection of Two Linked Lists (160)](https://leetcode.com/problems/intersection-of-two-linked-lists/) · Easy · Swap heads to equalize path lengths.
- [Find the Duplicate Number (287)](https://leetcode.com/problems/find-the-duplicate-number/) · Medium · Array as implicit list.

**Reversal**
- [Reverse Linked List II (92)](https://leetcode.com/problems/reverse-linked-list-ii/) · Medium · Head insertion.
- [Swap Nodes in Pairs (24)](https://leetcode.com/problems/swap-nodes-in-pairs/) · Medium · k-group with k = 2.
- [Reverse Nodes in k-Group (25)](https://leetcode.com/problems/reverse-nodes-in-k-group/) · Hard · Seed `prev` with the next group.

**Composite**
- [Palindrome Linked List (234)](https://leetcode.com/problems/palindrome-linked-list/) · Easy · Same skeleton as Reorder List, restore afterward.
- [Sort List (148)](https://leetcode.com/problems/sort-list/) · Medium · Merge sort with first middle.
- [Copy List with Random Pointer (138)](https://leetcode.com/problems/copy-list-with-random-pointer/) · Medium · Map vs interleave.
- [LRU Cache (146)](https://leetcode.com/problems/lru-cache/) · Medium · Map + doubly linked list with sentinels.
- [Merge k Sorted Lists (23)](https://leetcode.com/problems/merge-k-sorted-lists/) · Hard · Heap with tie-breaker index.

## Self-check

1. Why does a dummy head remove special cases, and what does it cost?
2. In Floyd's algorithm, why does resetting one pointer to head and moving both by 1 find the cycle start? Why can't fast skip over slow?
3. Why must list merge sort split at the first middle (`fast = head.next`)? What breaks on 2 nodes otherwise?
4. In k-group reversal, why seed `prev` with the node after the group instead of `None`?
5. In the interleaving copy, why must randoms be assigned in a separate pass before unweaving?
6. Which list problems create cycles if you forget one line, and which line?
7. Why does an LRU cache need a doubly linked list, and what do sentinels at both ends save?
