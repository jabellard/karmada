# 10. Trees and BST

Trees are the most common interview structure after arrays, and nearly every tree problem reduces to one decision: **what does the recursive function return for a subtree, and what state is passed down?** Then the code is ten lines. BSTs add ordering (O(n) searches become O(h)); "tree as graph" handles upward movement.

Companion code: `code/10_trees_bst.py` (TreeNode, `build_tree` from level-order lists, templates and solutions, with asserts).

## At a glance

- **Use when:** a `TreeNode` root or a hierarchy; depth/height/paths/ancestors; sorted-order queries on a BST.
- **Templates:**
  - Template 1: recursive DFS, return-value design: bottom-up summaries, top-down state, global + local paths
  - Template 2: BFS level order: per-level grouping, nearest-to-root
  - Template 3: iterative inorder: BST sorted order with early exit
- **Typical complexity:** O(n) time, O(h) stack (O(n) skewed); BST operations O(h).
- **Top trap:** validating a BST against the parent only (use bounds), or returning the bent path value to the parent in global + local.

## Core idea

A binary tree is a node plus two subtrees. Correctness is **structural induction**: if `dfs(child)` is right and the node combines correctly, `dfs(root)` is right.

Three DFS shapes cover nearly everything:

| Shape | Information flow | Signature | Examples |
|---|---|---|---|
| Top-down | Parent passes state down (depth, running sum, bounds, max-so-far) | `dfs(node, state) -> None` or `-> answer` | Path Sum, Count Good Nodes, Validate BST |
| Bottom-up | Children return summaries; parent combines | `dfs(node) -> summary` | Max Depth, Balanced, LCA, subtree sums |
| Global + local | Return what the parent can **extend**; record the path that **bends** here in a global | `dfs(node) -> extendable`, plus `best` | Diameter, Max Path Sum, Longest Univalue Path |

```
Global + local, for a path problem:

        node           returned to parent:  node + max(L, R)     (one arm, extendable)
       /    \          recorded globally:   L + node + R         (bent path, final)
     L-arm  R-arm
```

A path can't fork, so the parent extends only one arm of a child's path; hence the separate return value and global.

BST invariant: at every node, **all** left-subtree keys < node < **all** right-subtree keys. So inorder is sorted; search, insert, delete, floor/ceiling, and successor are O(h); validity means bounds, not parent comparison.

## Recognition signals

- A `TreeNode` root, or "hierarchy", "org chart", "file system", "nested categories", "parent pointers".
- "Depth", "height", "balanced", "diameter", "subtree sum", "is X a subtree": bottom-up.
- "Root-to-leaf", "max/min/sum so far from the root", "each node within a range": top-down.
- "Path anywhere" (not necessarily through the root): global + local.
- "Level", "row", "layer", "closest to the root", "visible from the side", "minimum depth": BFS level order.
- "Sorted", "k-th smallest", "successor", "range of keys", "closest value" on a BST: inorder or O(h) descent.
- "Common manager / common ancestor": LCA.
- "Count paths summing to target", any start: prefix sums on the root path (module 04 on a tree).
- "Distance from a node", "nodes k hops away", "burn/infect from a node": undirected graph (parent map) + BFS.
- "Rebuild from traversals", "send over the wire": construction and serialization.
- n up to 10^4..10^5, possibly skewed: recursion may hit Python's limit; mention iterative versions.

## Disguises

1. **Org chart escalation.** "Every employee has at most two direct reports. Given two employees, find the most junior manager with both in their reporting chain (an employee counts as their own chain)."
   Decode: lowest common ancestor; standard LCA already treats a node as its own ancestor.

2. **Network latency budget.** "A tree-shaped CDN has a latency on each link. Find the largest total latency between any two edge caches."
   Decode: weighted diameter (global + local): return the longest arm; record left + right.

3. **Config inheritance audit.** "Each config node may override a numeric limit. A node is 'compliant' if its limit is at least every ancestor's. Count compliant nodes."
   Decode: depends only on ancestors: top-down with `max_so_far` (Count Good Nodes).

4. **Budget rollups.** "Departments form a tree, each with a spend delta (possibly negative). Count top-to-bottom chains of command, starting and ending anywhere, whose deltas sum to exactly T."
   Decode: downward, any start/end, exact sum, negatives: prefix sums on the root path with undo (Path Sum III).

5. **Fire spreading in a hierarchy.** "An outage starts at service X in a dependency tree and spreads one link per minute in both directions. Which services go down at minute k?"
   Decode: upward spread needs parent links: undirected graph + k BFS layers (All Nodes Distance K).

6. **Leaderboard over a sorted index.** "Scores live in a binary search tree. Repeatedly find the next-higher score after a given one, in O(height) memory."
   Decode: in-order successor, or a BST iterator with a stack of left spines.

## Templates and worked examples

### Template 1: recursive DFS, return-value design

> - **Example problem:** [Binary Tree Maximum Path Sum (124)](https://leetcode.com/problems/binary-tree-maximum-path-sum/)
> - **Why this template:** the path needn't pass through the root, so each node returns a one-armed downward sum while a global records paths bending there.
> - **Worked below:** Example 1.6 (Diameter of Binary Tree (543) in Example 1.1; Serialize and Deserialize Binary Tree (297) in Example 1.7)
> - **Also fits:** [Lowest Common Ancestor of a Binary Tree (236)](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/)

Before coding, answer out loud:

1. **What does `dfs(node)` return for a subtree?** Precisely: "height in nodes", "max downward path sum starting at node", "(is_bst, min, max)".
2. **What does the empty subtree return?** A neutral value for the combine (0 for height, `-inf` for max, `True` for "valid").
3. **How do I combine children with the node?**
4. **Do I need ancestor state?** If yes, add parameters (top-down).
5. **Is the answer the root's return value, or a global updated at every node?** Use a global if the best answer can "bend" anywhere.

```python
def template_bottom_up(root: Optional[TreeNode]) -> int:
    def dfs(node: Optional[TreeNode]) -> int:
        if not node:
            return 0  # TODO: value for the empty subtree
        left = dfs(node.left)
        right = dfs(node.right)
        return 1 + max(left, right)  # TODO: combine children + node
    return dfs(root)
```

```python
def template_top_down(root: Optional[TreeNode]) -> int:
    best = 0

    def dfs(node: Optional[TreeNode], depth: int) -> None:  # TODO: state params
        nonlocal best
        if not node:
            return
        best = max(best, depth)  # TODO: use state at this node (often only at leaves)
        dfs(node.left, depth + 1)  # TODO: derive child state
        dfs(node.right, depth + 1)

    dfs(root, 1)
    return best
```

```python
def template_global_plus_local(root: Optional[TreeNode]) -> int:
    best = 0

    def dfs(node: Optional[TreeNode]) -> int:
        nonlocal best
        if not node:
            return 0
        left, right = dfs(node.left), dfs(node.right)
        best = max(best, left + right)  # TODO: answer if the path bends here
        return 1 + max(left, right)     # TODO: value the parent can extend
    dfs(root)
    return best
```

Knobs:
- Return a tuple when the parent needs several facts (`(height, balanced)`, `(rob, skip)` in House Robber III), or a sentinel like `-1` for "invalid" to short-circuit (`is_balanced`).
- Shared mutable state (a `path` list) needs **undo** on the way up (append before recursing, pop after): backtracking (`15-backtracking.md`).
- `nonlocal best` or returning it in a tuple both work; `nonlocal` is shorter.

#### Example 1.1: Longest path between any two nodes (global + local) [Easy]

**Problem:** [Diameter of Binary Tree (543)](https://leetcode.com/problems/diameter-of-binary-tree/) · **Template:** Template 1 (recursive DFS, return-value design), global + local form

**Problem.** Return the number of edges on the longest path between any two nodes of a binary tree, not necessarily through the root. Example: `[1, 2, 3, 4, 5]` gives 3 (`4 -> 2 -> 1 -> 3`, or `5 -> 2 -> 1 -> 3`); a single node gives 0.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** the path can bend anywhere, so the answer is a global (global + local).

**Approach:** `height(node)` returns subtree height in nodes (equivalently, edges on the longest downward path from the parent through `node`); empty is 0. Every path bends at a unique highest node, and the longest path bending at `node` has `left + right` edges. Record that globally; return `1 + max(left, right)` (the parent extends one arm). The max over all bend points is the diameter.

```python
def diameter(root: Optional[TreeNode]) -> int:
    best = 0

    def height(node: Optional[TreeNode]) -> int:
        nonlocal best
        if not node:
            return 0
        left, right = height(node.left), height(node.right)
        best = max(best, left + right)  # longest path that bends at node, in edges
        return 1 + max(left, right)     # height in nodes: what the parent can extend
    height(root)
    return best
```

**Complexity:** O(n) time; O(h) stack (O(n) skewed, O(log n) balanced).

**Edge cases:** edges vs nodes (height in nodes makes `left + right` the edge count; add 1 for nodes); the longest path may avoid the root (`[1, 2, None, 3, 4, 5, None, None, 6]` gives 4), hence the global; empty tree returns 0. A separate `height()` call per node is O(n^2) on a skewed tree.

</details>

#### Example 1.2: Count good nodes and validate a BST (top-down) [Medium]

**Problems:** [Count Good Nodes in Binary Tree (1448)](https://leetcode.com/problems/count-good-nodes-in-binary-tree/), [Validate Binary Search Tree (98)](https://leetcode.com/problems/validate-binary-search-tree/) · **Template:** Template 1 (recursive DFS, return-value design), top-down form

**Problem A.** A node is "good" if no node on its root path has a greater value. Count good nodes. (LeetCode: Count Good Nodes in Binary Tree)

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** depends only on ancestors: pass `max_so_far` down.

```python
def good_nodes(root: Optional[TreeNode]) -> int:
    def dfs(node: Optional[TreeNode], mx: int) -> int:
        if not node:
            return 0
        good = 1 if node.val >= mx else 0
        mx = max(mx, node.val)
        return good + dfs(node.left, mx) + dfs(node.right, mx)
    return dfs(root, -math.inf)
```

</details>

**Problem B.** Decide whether a tree is a valid BST (strict). (LeetCode: Validate Binary Search Tree)

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** each node must lie in an open interval `(lo, hi)` inherited from its ancestors; going left tightens `hi`, right tightens `lo`. Parent-only comparison fails: in `[5, 4, 6, null, null, 3, 7]`, node 3 is a valid left child of 6 but violates the root's bound.

```python
def is_valid_bst(root: Optional[TreeNode]) -> bool:
    def ok(node: Optional[TreeNode], lo: float, hi: float) -> bool:
        if not node:
            return True
        if not (lo < node.val < hi):
            return False
        return ok(node.left, lo, node.val) and ok(node.right, node.val, hi)
    return ok(root, -math.inf, math.inf)
```

**Complexity:** O(n) time, O(h) stack for both.

**Edge cases:** duplicates (strict inequality rejects `[2, 2, 2]`); values at integer extremes (use `math.inf`, not sentinel ints). Alternative for B: iterative inorder must be strictly increasing.

</details>

#### Example 1.3: Lowest common ancestor (bottom-up), binary tree and BST [Medium]

**Problems:** [Lowest Common Ancestor of a Binary Tree (236)](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/), [Lowest Common Ancestor of a Binary Search Tree (235)](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/) · **Template:** Template 1 (recursive DFS, return-value design), bottom-up form

**Problem.** Given two nodes `p` and `q` in the tree, return their lowest common ancestor (a node may be its own ancestor).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "lowest node with both below it".

**Approach:** `lca(node)` returns p or q if exactly one is in the subtree, the LCA if both are, else None. If `node` is `p` or `q`, return it (if the other is below, it's the LCA anyway). If both children return non-null, `p` and `q` split here: `node` is the LCA. Otherwise pass up the non-null side.

```python
def lca(root: Optional[TreeNode], p: TreeNode, q: TreeNode) -> Optional[TreeNode]:
    if not root or root is p or root is q:
        return root
    left = lca(root.left, p, q)
    right = lca(root.right, p, q)
    if left and right:
        return root
    return left or right
```

BST version: walk down; both smaller go left, both larger go right, else this is the split point. O(h), O(1) space.

```python
def lca_bst(root: TreeNode, p: TreeNode, q: TreeNode) -> TreeNode:
    lo, hi = min(p.val, q.val), max(p.val, q.val)
    while root:
        if hi < root.val:
            root = root.left
        elif lo > root.val:
            root = root.right
        else:
            return root  # split point
    raise ValueError("nodes not in tree")
```

**Complexity:** O(n) / O(h) time; O(h) / O(1) space.

**Edge cases:** the binary-tree version assumes both nodes exist; if not guaranteed, count finds and trust the result only if both were seen. With parent pointers, LCA becomes "intersection of two linked lists" (module 09).

</details>

#### Example 1.4: Count downward paths with a target sum (prefix sums on the root path) [Medium]

**Problem:** [Path Sum III (437)](https://leetcode.com/problems/path-sum-iii/) · **Template:** Template 1 (recursive DFS, return-value design), top-down form with undo

**Problem.** Count downward paths (from any node to itself or any descendant) summing to `target`. Values may be negative.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** on an array this is "count subarrays with sum k" (prefix sums + hash map, module 04). On a tree, the "array" is the current root-to-node path.

**Approach:** keep `seen[prefix] = count` for the current root path only. At running sum `cur`, `seen[cur - target]` paths end here. Add `cur` before recursing and **remove it after**, hiding it from sibling subtrees.

```python
def path_sum_count(root: Optional[TreeNode], target: int) -> int:
    seen: dict[int, int] = defaultdict(int)
    seen[0] = 1  # empty prefix: path starting at the root

    def dfs(node: Optional[TreeNode], cur: int) -> int:
        if not node:
            return 0
        cur += node.val
        res = seen[cur - target]
        seen[cur] += 1
        res += dfs(node.left, cur) + dfs(node.right, cur)
        seen[cur] -= 1  # this prefix is not an ancestor of sibling subtrees
        return res

    return dfs(root, 0)
```

**Complexity:** O(n) time; O(h) stack plus O(h) live map entries. Brute force (DFS from every node) is O(n·h), O(n^2) skewed.

**Edge cases:** `seen[0] = 1` counts paths from the root; read `seen[cur - target]` **before** inserting `cur`, or a zero target counts the empty path; zeros in the tree (`[0, 0, 0]`, target 0, gives 5).

Path Sum I and II are top-down: subtract along the way, check only at leaves. Path Sum II shares one list, popping on the way up and copying only on a match (`path_sum_all`).

</details>

#### Example 1.5: Rebuild a tree from preorder and inorder [Medium]

**Problem:** [Construct Binary Tree from Preorder and Inorder Traversal (105)](https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/) · **Template:** Template 1 (recursive DFS, return-value design)

**Problem.** Rebuild a tree with unique values from its preorder and inorder sequences.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** preorder gives the root; inorder splits left from right.

**Approach:** the next unused preorder element is the subtree's root; its inorder index splits the range. Recurse left **first**, since preorder lists the whole left subtree before the right. A value→index map makes splits O(1); a shared preorder cursor avoids slicing.

```python
def build_from_pre_in(preorder: list[int], inorder: list[int]) -> Optional[TreeNode]:
    idx = {v: i for i, v in enumerate(inorder)}
    pre_i = 0

    def build(lo: int, hi: int) -> Optional[TreeNode]:  # inorder[lo..hi]
        nonlocal pre_i
        if lo > hi:
            return None
        val = preorder[pre_i]
        pre_i += 1
        node = TreeNode(val)
        mid = idx[val]
        node.left = build(lo, mid - 1)   # left first: preorder consumes left subtree next
        node.right = build(mid + 1, hi)
        return node

    return build(0, len(inorder) - 1)
```

**Complexity:** O(n) time and space; slicing versions are O(n^2).

**Edge cases:** requires unique values. Postorder + inorder: consume postorder from the end, building the **right** subtree first. Preorder + postorder generally doesn't determine a unique tree.

</details>

#### Example 1.6: Maximum path sum anywhere in the tree (global + local) [Hard]

**Problem:** [Binary Tree Maximum Path Sum (124)](https://leetcode.com/problems/binary-tree-maximum-path-sum/) · **Template:** Template 1 (recursive DFS, return-value design), global + local form

**Problem.** A path is any sequence of connected nodes, each used once, not necessarily through the root. Find the maximum non-empty path sum. Values may be negative.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** any path, anywhere: global + local.

**Approach:** `gain(node)` returns the best downward sum **starting at** `node`, dropping negative children (`max(.., 0)`). Record `node.val + left + right` (bending here) globally; return `node.val + max(left, right)` (one arm).

```python
def max_path_sum(root: Optional[TreeNode]) -> int:
    best = -math.inf

    def gain(node: Optional[TreeNode]) -> int:
        nonlocal best
        if not node:
            return 0
        left = max(gain(node.left), 0)   # drop a negative branch
        right = max(gain(node.right), 0)
        best = max(best, node.val + left + right)  # path bends here
        return node.val + max(left, right)          # parent can extend one side only
    gain(root)
    return best
```

**Complexity:** O(n) time, O(h) space.

**Edge cases:** all negative (start `best = -inf`, not 0; the answer is the largest node); single node.

Same skeleton, different "arm" definitions:

| Problem | Arm returned | Bent value recorded |
|---|---|---|
| Diameter (edges) | `1 + max(L, R)` (height) | `L + R` |
| Max Path Sum | `val + max(L, R, 0)` | `val + max(L,0) + max(R,0)` |
| Longest Univalue Path | extends only through an equal-valued child, else 0 | `l + r` |

```python
def longest_univalue_path(root: Optional[TreeNode]) -> int:
    best = 0

    def arm(node: Optional[TreeNode]) -> int:
        nonlocal best
        if not node:
            return 0
        l, r = arm(node.left), arm(node.right)
        l = l + 1 if node.left and node.left.val == node.val else 0
        r = r + 1 if node.right and node.right.val == node.val else 0
        best = max(best, l + r)
        return max(l, r)
    arm(root)
    return best
```

Always call `arm(child)`, even when values differ: the best univalue path may lie entirely inside that child.

</details>

#### Example 1.7: Serialize and deserialize a binary tree (preorder with null markers) [Hard]

**Problem:** [Serialize and Deserialize Binary Tree (297)](https://leetcode.com/problems/serialize-and-deserialize-binary-tree/) · **Template:** Template 1 (recursive DFS, return-value design)

**Problem.** Design `serialize(root) -> str` and `deserialize(str) -> root` so any binary tree survives a round trip.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** values alone lose the shape, so encode structure too: a DFS that also writes a marker for every empty child pins the tree down.

**Approach:** preorder with null markers is self-delimiting: the deserializer reads tokens in written order, and `#` marks each subtree's end. No index math.

```python
class Codec:
    def serialize(self, root: Optional[TreeNode]) -> str:
        out: list[str] = []

        def dfs(node: Optional[TreeNode]) -> None:
            if not node:
                out.append("#")
                return
            out.append(str(node.val))
            dfs(node.left)
            dfs(node.right)

        dfs(root)
        return ",".join(out)

    def deserialize(self, data: str) -> Optional[TreeNode]:
        tokens = iter(data.split(","))

        def build() -> Optional[TreeNode]:
            t = next(tokens)
            if t == "#":
                return None
            node = TreeNode(int(t))
            node.left = build()
            node.right = build()
            return node

        return build()
```

**Complexity:** O(n) both ways. For a BST, drop the null markers: preorder alone suffices, rebuilt with bounds (Serialize and Deserialize BST).

**Edge cases:** an empty tree serializes to `"#"` and deserializes back to `None`; negative and multi-digit values are safe because tokens are comma-separated.

</details>

### Template 2: BFS level order

> - **Example problem:** [Binary Tree Level Order Traversal (102)](https://leetcode.com/problems/binary-tree-level-order-traversal/)
> - **Why this template:** "level by level, left to right" is BFS with a per-level queue-size snapshot.
> - **Worked below:** Example 2.1
> - **Also fits:** [Binary Tree Right Side View (199)](https://leetcode.com/problems/binary-tree-right-side-view/)

```python
def template_bfs_levels(root: Optional[TreeNode]) -> list[list[int]]:
    if not root:
        return []
    out = []
    q = deque([root])
    while q:
        level = []
        for _ in range(len(q)):  # snapshot: exactly the nodes of this level
            node = q.popleft()
            level.append(node.val)  # TODO: per-node work
            if node.left:
                q.append(node.left)
            if node.right:
                q.append(node.right)
        out.append(level)  # TODO: per-level work
    return out
```

Knobs: per-level aggregate (average, max, last for right view); reverse alternate levels (zigzag); stop at the first leaf (minimum depth); track `(node, index)` for widths (Maximum Width: children `2i`, `2i+1`).

#### Example 2.1: Group node values by depth (BFS with level snapshot) [Medium]

**Problem:** [Binary Tree Level Order Traversal (102)](https://leetcode.com/problems/binary-tree-level-order-traversal/) · **Template:** Template 2 (BFS level order)

**Problem.** Return a binary tree's values grouped by depth: one list per level, top first, each left to right. Example: `[3, 9, 20, null, null, 15, 7]` gives `[[3], [9, 20], [15, 7]]`; an empty tree gives `[]`.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "level by level, left to right" = BFS; grouping needs each level's end, hence the queue-size snapshot.

**Approach:** invariant: at the top of the outer loop, the queue holds exactly one level, left to right. Pop exactly `len(q)` nodes; their children (left first) queue behind as the next level. Without the snapshot, levels mix.

```python
def level_order(root: Optional[TreeNode]) -> list[list[int]]:
    if not root:
        return []
    out: list[list[int]] = []
    q = deque([root])
    while q:
        level = []
        for _ in range(len(q)):  # exactly the nodes that were queued for this level
            node = q.popleft()
            level.append(node.val)
            if node.left:
                q.append(node.left)
            if node.right:
                q.append(node.right)
        out.append(level)
    return out
```

**Complexity:** O(n) time; O(w) queue for max level width w (about n/2 in a complete tree); O(n) output.

**Edge cases:** empty tree returns `[]`, not `[[]]`; skewed tree (one node per level); `range(len(q))` is evaluated once, so appending inside the loop is safe. Use `deque.popleft()`, not `list.pop(0)` (O(n)). DFS passing `depth` and appending to `out[depth]` also works, with O(h) stack instead of O(w) queue.

</details>

### Template 3: iterative inorder

> - **Example problem:** [Kth Smallest Element in a BST (230)](https://leetcode.com/problems/kth-smallest-element-in-a-bst/)
> - **Why this template:** the BST's kth smallest is the k-th inorder visit, and the iterative form can stop early.
> - **Worked below:** Example 3.1

```python
def template_iterative_inorder(root: Optional[TreeNode]) -> list[int]:
    out, stack, cur = [], [], root
    while cur or stack:
        while cur:  # go as far left as possible
            stack.append(cur)
            cur = cur.left
        cur = stack.pop()
        out.append(cur.val)  # TODO: visit (k-th smallest: count and stop early)
        cur = cur.right
    return out
```

Invariant: the stack holds ancestors whose left subtree is in progress and visit pending; O(h) memory. Base for k-th smallest with early exit, the BST iterator, BST two-sum with two iterators, and recovering a BST with two swapped nodes.

Other iterative orders (code file): preorder with a stack (push right, then left); postorder as a reversed (node, right, left) preorder.

#### Example 3.1: k-th smallest value in a BST [Medium]

**Problem:** [Kth Smallest Element in a BST (230)](https://leetcode.com/problems/kth-smallest-element-in-a-bst/) · **Template:** Template 3 (iterative inorder)

**Problem.** Given a BST and `1 ≤ k ≤ n`, return the k-th smallest value (1-indexed). Example: `[3, 1, 4, null, 2]`, `k = 1` → 1; `[5, 3, 6, 2, 4, null, null, 1]`, `k = 3` → 3.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** BST + "k-th smallest" = k-th inorder visit, since inorder yields sorted order.

**Approach:** Template 3, counting pops. Each pop is the next value in sorted order; return at the k-th. Stopping early beats collecting the full list.

```python
def kth_smallest(root: Optional[TreeNode], k: int) -> int:
    stack, cur = [], root
    while cur or stack:
        while cur:
            stack.append(cur)
            cur = cur.left
        cur = stack.pop()
        k -= 1
        if k == 0:
            return cur.val
        cur = cur.right
    raise ValueError("k out of range")
```

**Complexity:** O(h + k) time: descending the left spine costs h, then each of the k pops is amortized O(1). O(h) space.

**Edge cases:** `k = 1` is the leftmost node; `k = n` traverses everything; a skewed tree makes h = n. Follow-up: with frequent inserts/deletes and repeated queries, store each node's subtree size; at a node with left size `L`, the answer is the node if `k == L + 1`, else go left (`k ≤ L`) or right with `k - L - 1`. O(h) per query, with sizes maintained on each update.

</details>

### Combined examples

#### Example C.1: All nodes at distance k from a target [Medium]

**Problem:** [All Nodes Distance K in Binary Tree (863)](https://leetcode.com/problems/all-nodes-distance-k-in-binary-tree/) · **Template:** Template 1 (recursive DFS, return-value design) to record parent links, then Template 2 (BFS level order) for k layers out from the target; both are needed because distance can run up through parents, which a downward traversal alone can't follow

**Problem.** Given a binary tree, a target node, and an integer k, return the values of all nodes exactly k edges from the target.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** distance can go up through parents, so treat the tree as an undirected graph.

**Approach:** one DFS records `parent[node]`; then BFS k layers from `target` over `{left, right, parent}`, with a `seen` set (edges go both ways now).

```python
def distance_k(root: Optional[TreeNode], target: TreeNode, k: int) -> list[int]:
    parent: dict[TreeNode, Optional[TreeNode]] = {root: None}
    stack = [root]
    while stack:
        node = stack.pop()
        for child in (node.left, node.right):
            if child:
                parent[child] = node
                stack.append(child)

    seen = {target}
    frontier = [target]
    for _ in range(k):
        nxt = []
        for node in frontier:
            for nb in (node.left, node.right, parent[node]):
                if nb and nb not in seen:
                    seen.add(nb)
                    nxt.append(nb)
        frontier = nxt
    return [n.val for n in frontier]
```

**Complexity:** O(n) time and space.

**Edge cases:** k = 0 returns `[target]`; k beyond reach returns `[]`. The same conversion solves "time for an infection to spread" (BFS until empty) and anything moving up (`to_adjacency` in the code file).

</details>

## Template cheat sheet

All templates, for review.

### Template 1: recursive DFS, return-value design
```python
def template_bottom_up(root: Optional[TreeNode]) -> int:
    def dfs(node: Optional[TreeNode]) -> int:
        if not node:
            return 0  # TODO: value for the empty subtree
        left = dfs(node.left)
        right = dfs(node.right)
        return 1 + max(left, right)  # TODO: combine children + node
    return dfs(root)
```
```python
def template_top_down(root: Optional[TreeNode]) -> int:
    best = 0

    def dfs(node: Optional[TreeNode], depth: int) -> None:  # TODO: state params
        nonlocal best
        if not node:
            return
        best = max(best, depth)  # TODO: use state at this node (often only at leaves)
        dfs(node.left, depth + 1)  # TODO: derive child state
        dfs(node.right, depth + 1)

    dfs(root, 1)
    return best
```
```python
def template_global_plus_local(root: Optional[TreeNode]) -> int:
    best = 0

    def dfs(node: Optional[TreeNode]) -> int:
        nonlocal best
        if not node:
            return 0
        left, right = dfs(node.left), dfs(node.right)
        best = max(best, left + right)  # TODO: answer if the path bends here
        return 1 + max(left, right)     # TODO: value the parent can extend
    dfs(root)
    return best
```
### Template 2: BFS level order
```python
def template_bfs_levels(root: Optional[TreeNode]) -> list[list[int]]:
    if not root:
        return []
    out = []
    q = deque([root])
    while q:
        level = []
        for _ in range(len(q)):  # snapshot: exactly the nodes of this level
            node = q.popleft()
            level.append(node.val)  # TODO: per-node work
            if node.left:
                q.append(node.left)
            if node.right:
                q.append(node.right)
        out.append(level)  # TODO: per-level work
    return out
```
### Template 3: iterative inorder
```python
def template_iterative_inorder(root: Optional[TreeNode]) -> list[int]:
    out, stack, cur = [], [], root
    while cur or stack:
        while cur:  # go as far left as possible
            stack.append(cur)
            cur = cur.left
        cur = stack.pop()
        out.append(cur.val)  # TODO: visit (k-th smallest: count and stop early)
        cur = cur.right
    return out
```

## BST operations reference

All O(h): O(log n) balanced, O(n) skewed.

| Operation | Idea |
|---|---|
| Search | Go left if smaller, right if larger |
| Insert | Recurse to the null slot; `root.left = insert(root.left, v)`; return root |
| Delete | 0/1 child: return the other child. 2 children: copy in the in-order successor (leftmost of right subtree), then delete it from the right subtree |
| k-th smallest | Iterative inorder, stop at k. Frequent inserts and queries: store subtree sizes, descend by size |
| In-order successor | Descend; if `p.val < node.val`, record node as candidate and go left; else go right |
| Floor / ceiling / closest | Same descent, tracking the best candidate |
| Range sum / trim | Prune subtrees entirely outside `[lo, hi]` |

```python
def delete_bst(root: Optional[TreeNode], key: int) -> Optional[TreeNode]:
    if not root:
        return None
    if key < root.val:
        root.left = delete_bst(root.left, key)
    elif key > root.val:
        root.right = delete_bst(root.right, key)
    else:
        if not root.left:
            return root.right
        if not root.right:
            return root.left
        succ = root.right  # two children: copy in-order successor, delete it from right
        while succ.left:
            succ = succ.left
        root.val = succ.val
        root.right = delete_bst(root.right, succ.val)
    return root
```

```python
def inorder_successor(root: Optional[TreeNode], p: TreeNode) -> Optional[TreeNode]:
    succ = None
    while root:
        if p.val < root.val:
            succ = root  # candidate; try to find a smaller one on the left
            root = root.left
        else:
            root = root.right
    return succ
```

BST iterator (Binary Search Tree Iterator): iterative inorder split across calls. `next()` pops a node and pushes its right child's left spine. Amortized O(1) (each node pushed and popped once); O(h) memory.

```python
class BSTIterator:
    def __init__(self, root: Optional[TreeNode]):
        self.stack: list[TreeNode] = []
        self._push_left(root)

    def _push_left(self, node: Optional[TreeNode]) -> None:
        while node:
            self.stack.append(node)
            node = node.left

    def next(self) -> int:
        node = self.stack.pop()
        self._push_left(node.right)
        return node.val

    def has_next(self) -> bool:
        return bool(self.stack)
```

## Variations and extensions

- **Tree DP**: a tuple of states per node: House Robber III returns `(best_if_robbed, best_if_skipped)`; Binary Tree Cameras, a 3-state code (`16-dynamic-programming-foundations.md`).
- **Subtree matching**: `same_tree` at each node (O(n·m)), or serialize both with null markers and substring-search (`19-tries-and-string-algorithms.md`).
- **Vertical order**: BFS carrying `(node, column)`, group by column, sort by `(row, value)` if required.
- **Flatten to linked list**: reverse postorder (right, left, node) with a `prev` pointer, or Morris-style rewiring.
- **Morris traversal**: O(1)-space inorder by temporarily threading predecessor right pointers. Mention for O(1) space; rarely coded.
- **N-ary trees and tries**: same templates over `children`. Tries: `19-tries-and-string-algorithms.md`.
- **Merging two BSTs into a sorted list**: merge two inorder generators (`iter_inorder`) like sorted lists. O(n + m), O(h1 + h2) memory.
- **Balanced BSTs in Python**: no stdlib TreeMap. Use `bisect` on a sorted list (O(n) insert, often fine), heaps with lazy deletion (`11-heaps-top-k-and-merging.md`), or `sortedcontainers.SortedList` if allowed (ask).

Follow-ups:
- 10^5-node chain → Python's recursion limit is ~1000: use an explicit stack, or `sys.setrecursionlimit` (risks a C stack overflow; say so).
- Repeated k-th smallest with inserts → subtree sizes in nodes.
- Huge tree on disk → B-tree; mention branching factor and page-sized nodes.

## Common pitfalls

- Validating a BST against the parent only. Use bounds or inorder.
- 0 as the neutral value for a max over possibly negative values. Use `-math.inf`.
- Loose leaf checks: in Path Sum, a one-child node isn't a leaf (`has_path_sum([1, 2], 1)` is `False`).
- Forgetting `path.pop()` on a shared list, or appending `path` instead of `path.copy()`.
- Skipping `seen[cur] -= 1` in Path Sum III, so siblings see non-ancestor prefixes.
- BFS without a level snapshot when grouping per level.
- Calling a recursive helper per node (`height()` inside `is_balanced()` is O(n^2)). Return combined info once.
- Global + local: returning the bent value to the parent, which can extend only one arm.
- Missing `nonlocal` for an integer accumulator: assignment creates a new local.
- Recursing on `preorder[1:mid+1]`-style slices: correct but O(n^2) time and memory.

## When NOT to use / look-alikes

| Looks like | Actually | How to tell |
|---|---|---|
| Tree path problem with upward moves | Graph BFS | Needs parent access or "distance from a node": build a parent map |
| "Hierarchy" with shared children | DAG / graph, not a tree | Multiple parents; need `visited` (module 13) |
| BST k-th smallest | Heap (module 11) | Input isn't a BST, or values stream in |
| "Level order" on a grid | Grid BFS (module 13) | No tree, just layers by distance |
| Sorted order queries with inserts in Python | `bisect` / heap / SortedList | Hand-built balanced BSTs are almost never expected |
| Path sum with only positive values on an array | Sliding window (module 06) | Trees have no window equivalent; prefix sums on the path work for any sign |
| Parentheses / nesting depth | Stack (module 08) | Implicitly a tree, but a counter or stack is simpler |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Binary Tree Level Order Traversal (102)](https://leetcode.com/problems/binary-tree-level-order-traversal/) · Medium · BFS by level: the level template.
2. [Validate Binary Search Tree (98)](https://leetcode.com/problems/validate-binary-search-tree/) · Medium · Top-down bounds: bounds, strict.
3. [Lowest Common Ancestor of a Binary Tree (236)](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-tree/) · Medium · Bottom-up return: return-value semantics.
4. [Diameter of Binary Tree (543)](https://leetcode.com/problems/diameter-of-binary-tree/) · Easy · Global answer + local return: edges vs nodes.
5. [Binary Tree Maximum Path Sum (124)](https://leetcode.com/problems/binary-tree-maximum-path-sum/) · Hard · Global answer + local return: clip negative arms.

### Optional

**DFS shapes**
- [Maximum Depth of Binary Tree (104)](https://leetcode.com/problems/maximum-depth-of-binary-tree/) · Easy · The canonical bottom-up.
- [Invert Binary Tree (226)](https://leetcode.com/problems/invert-binary-tree/) · Easy · Swap children, recurse; the simplest structural DFS.
- [Same Tree (100)](https://leetcode.com/problems/same-tree/) · Easy · Pairwise recursion on two nodes.
- [Symmetric Tree (101)](https://leetcode.com/problems/symmetric-tree/) · Easy · Pairwise recursion on mirrored children.
- [Subtree of Another Tree (572)](https://leetcode.com/problems/subtree-of-another-tree/) · Easy · Same Tree check launched from every node (or serialize both and substring-match).
- [Balanced Binary Tree (110)](https://leetcode.com/problems/balanced-binary-tree/) · Easy · Return height or -1 to avoid O(n^2).
- [Count Good Nodes in Binary Tree (1448)](https://leetcode.com/problems/count-good-nodes-in-binary-tree/) · Medium · Top-down max-so-far.
- [House Robber III (337)](https://leetcode.com/problems/house-robber-iii/) · Medium · Return a pair of states.

**BFS**
- [Binary Tree Right Side View (199)](https://leetcode.com/problems/binary-tree-right-side-view/) · Medium · Last per level.
- [Binary Tree Zigzag Level Order Traversal (103)](https://leetcode.com/problems/binary-tree-zigzag-level-order-traversal/) · Medium · Reverse alternate levels.

**Paths**
- [Path Sum (112)](https://leetcode.com/problems/path-sum/) · Easy · Check only at leaves.
- [Path Sum II (113)](https://leetcode.com/problems/path-sum-ii/) · Medium · Backtracking with copy-on-match.
- [Path Sum III (437)](https://leetcode.com/problems/path-sum-iii/) · Medium · Prefix sums on the root path, with undo.

**LCA and BST**
- [Lowest Common Ancestor of a Binary Search Tree (235)](https://leetcode.com/problems/lowest-common-ancestor-of-a-binary-search-tree/) · Medium · Split point.
- [Kth Smallest Element in a BST (230)](https://leetcode.com/problems/kth-smallest-element-in-a-bst/) · Medium · Iterative inorder, early exit.
- [Delete Node in a BST (450)](https://leetcode.com/problems/delete-node-in-a-bst/) · Medium · Successor swap.
- [Binary Search Tree Iterator (173)](https://leetcode.com/problems/binary-search-tree-iterator/) · Medium · Left-spine stack.

**Construction / graph**
- [Construct Binary Tree from Preorder and Inorder Traversal (105)](https://leetcode.com/problems/construct-binary-tree-from-preorder-and-inorder-traversal/) · Medium · Index map, left first.
- [All Nodes Distance K in Binary Tree (863)](https://leetcode.com/problems/all-nodes-distance-k-in-binary-tree/) · Medium · Parent map + BFS.
- [Serialize and Deserialize Binary Tree (297)](https://leetcode.com/problems/serialize-and-deserialize-binary-tree/) · Hard · Preorder with null markers.

## Self-check

1. For Diameter, what does the recursive function return, and why does it differ from what you record globally?
2. Why is comparing each node only to its parent insufficient for BST validation? Give a counterexample.
3. In Path Sum III, why must you decrement the prefix count after recursing, and why is `seen[0] = 1` needed?
4. In binary-tree LCA, what can `lca(node)` return, and what does each value mean? What breaks if one target is absent?
5. Why does preorder reconstruction recurse into the left subtree before the right? What changes for postorder + inorder?
6. Why is `BSTIterator.next()` amortized O(1) even though a single call can push h nodes?
7. When would you choose BFS over DFS on a tree? Name two problems where BFS stops early.
