# 19. Tries and String Algorithms

Most string problems reduce to four tools: a **trie** for prefix queries over many words, **expand-around-center** for palindromes, the **prefix function (KMP)** for exact matching and periodicity, and **rolling hashes** for O(1) substring comparisons. Infra analogs: routing tables and autocomplete (tries), log matching (KMP), deduplication and content-defined chunking (rolling hashes), wire formats (length-prefixed encoding).

Companion code: `code/19_tries_strings.py`.

## At a glance

- **Use when:** many words with prefix queries (trie); contiguous palindromes (centers); exact matching or periodicity (prefix function); many same-length substring comparisons (rolling hash).
- **Templates:**
  - Template 1: Trie (node class): several node fields or methods to add.
  - Template 2: Trie as dict-of-dicts: fast to write; grid word search, bit tries.
  - Template 3: expand around center: longest/count palindromic substrings.
  - Template 4: prefix function (KMP): linear matching, period, longest palindromic prefix.
  - Template 5: rolling hash: fixed-length windows, binary search on length.
- **Typical complexity:** trie O(L) per op; centers O(n^2) time, O(1) space; KMP O(n + m); rolling hash O(n) per length.
- **Top trap:** treating a prefix path as a word (forgetting `is_end`), and trusting an unverified hash match.

## Core idea

| Tool | Key insight | Invariant |
|---|---|---|
| Trie | Shared prefixes share a path; a prefix query costs O(len(prefix)), regardless of word count | The node on path p represents exactly prefix p; `is_end` marks complete words |
| Expand around center | Every palindrome has a center (char or gap): 2n-1 centers | While `s[lo] == s[hi]`, `s[lo..hi]` is a palindrome |
| Prefix function | On a mismatch, the pattern says how far to fall back | `pi[i]` = length of the longest proper prefix of `s[:i+1]` that is also its suffix |
| Rolling hash | A window's polynomial hash updates in O(1) per slide | `h = sum(s[j] * B^(L-1-(j-i)))  mod M` for window `[i, i+L)` |

```
Trie for {app, apple, apt, bat}:          Prefix function of "aabaaab":
root                                        a a b a a a b
 ├─ a ─ p ─ p*─ l ─ e*                      0 1 0 1 2 2 3
 │       └─ t*                              (pi[6]=3: "aab" is both prefix and suffix)
 └─ b ─ a ─ t*          (* = is_end)
```

## Recognition signals

- **Many words, prefix questions**: "starts with", "autocomplete", "shortest root", "wildcard dictionary". Trie.
- **Many words searched at once in a grid or text**: Word Search II, multi-pattern matching. The trie drives one DFS for all words.
- **Maximize XOR** of pairs, or constrained XOR queries: binary trie over bits, greedy from the high bit.
- **Palindromic substring** (contiguous): expand around center. Subsequence instead: DP (`16-dynamic-programming-foundations.md`).
- **Pattern in text in linear time**, **is the string a repetition**, **longest prefix that is also a suffix**, **fewest chars to prepend for a palindrome**: prefix function.
- **Many same-length substring comparisons**, "repeated substring of length L", "longest duplicate substring" (n up to 10^5): rolling hash, often with binary search on length.
- **Group strings by equivalence** (anagrams, shifts, patterns): canonical signature as dict key (`04-arrays-hashing-prefix-sums.md`).
- **Serialize arbitrary strings**: length-prefix encoding.
- **Evaluate an expression string**: tokenize, then stack parsing (`08-stacks-and-monotonic-structures.md`).

## Disguises

1. **Kubernetes label autocomplete.** As the user types each character, a CLI suggests up to three existing resource names, alphabetically.
   **Decode:** per-keystroke prefix query, bounded top-k. Trie with a precomputed `top` list per node (Search Suggestions System).

2. **Path normalization by mount points.** Given mount prefixes, rewrite every path token to the shortest mount prefix it starts with.
   **Decode:** shortest dictionary prefix per word. Trie; stop at the first `is_end` on the walk (Replace Words).

3. **Banned-term scanner over a character grid.** Find which of 10^4 banned terms can be spelled by adjacent cells in a board.
   **Decode:** many words over a grid. One backtracking DFS follows trie edges and prunes dead branches (Word Search II).

4. **Detecting replayed log segments.** Find every 10-character segment of an event log that appears more than once.
   **Decode:** fixed-length substrings compared many times. Rolling hash over a 10-wide window plus a seen set (Repeated DNA Sequences).

5. **Config file is a repeated stanza?** Decide whether a configuration string is one block repeated two or more times.
   **Decode:** periodicity. With `b = pi[-1]`, the string is periodic exactly when `b > 0` and `n % (n - b) == 0` (Repeated Substring Pattern).

6. **Wire protocol for string lists.** Design `encode(list[str]) -> str` and `decode` so the strings may contain any character, including your delimiter.
   **Decode:** delimiter escaping is fragile; a length prefix plus separator is unambiguous (Encode and Decode Strings).

## Templates and worked examples

### Template 1: Trie (node class)

> - **Example problem:** [Implement Trie (Prefix Tree) (208)](https://leetcode.com/problems/implement-trie-prefix-tree/)
> - **Why this template:** `insert`, `search`, and `startsWith` map directly onto the class's `insert`, `search` (checks `is_end`), and `starts_with` (walk only).
> - **Worked below:** Example 1.1
> - **Also fits:** [Design Add and Search Words Data Structure (211)](https://leetcode.com/problems/design-add-and-search-words-data-structure/)

```python
class TrieNode:
    __slots__ = ("children", "is_end", "count", "word")

    def __init__(self) -> None:
        self.children: dict[str, TrieNode] = {}
        self.is_end = False      # a word ends exactly here
        self.count = 0           # number of words passing through (prefix count)
        self.word: str | None = None   # optional: store the word at its end node


class Trie:
    def __init__(self) -> None:
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for ch in word:
            node = node.children.setdefault(ch, TrieNode())
            node.count += 1
        node.is_end = True
        node.word = word

    def _walk(self, s: str) -> TrieNode | None:
        node = self.root
        for ch in s:
            node = node.children.get(ch)
            if node is None:
                return None
        return node

    def search(self, word: str) -> bool:
        node = self._walk(word)
        return node is not None and node.is_end

    def starts_with(self, prefix: str) -> bool:
        return self._walk(prefix) is not None

    def count_prefix(self, prefix: str) -> int:
        if not prefix:
            return 0
        node = self._walk(prefix)
        return node.count if node else 0
```

Knobs (node payload, per problem):
- `is_end`: required for exact search.
- `count`: words with this prefix (Map Sum Pairs stores a sum instead).
- `word`: the full word at its end node, so DFS needn't rebuild strings (Word Search II).
- `top`: bounded list of best completions (autocomplete).
- Children as `dict` (any alphabet; the Python default) or `[None] * 26` (faster, fixed lowercase).

#### Example 1.1: Implement a trie with wildcard search [Medium]

**Problems:** [Implement Trie (Prefix Tree) (208)](https://leetcode.com/problems/implement-trie-prefix-tree/), [Design Add and Search Words Data Structure (211)](https://leetcode.com/problems/design-add-and-search-words-data-structure/) · **Template:** Template 1 (Trie, node class)

Problem: support `insert(word)`, `search(word)`, `startsWith(prefix)`. Variant: `search` patterns may contain `.` matching any one letter.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** explicit trie design; the wildcard variant adds a DFS branching on `.`.

**Part (a): Implement Trie.**

**Approach:** each node maps characters to children and records whether a word ends there. `insert` creates missing nodes, `search` walks and checks `is_end`, `starts_with` only checks the path exists. Same code as Template 1:

```python
class TrieNode:
    __slots__ = ("children", "is_end", "count", "word")

    def __init__(self) -> None:
        self.children: dict[str, TrieNode] = {}
        self.is_end = False      # a word ends exactly here
        self.count = 0           # number of words passing through (prefix count)
        self.word: str | None = None   # optional: store the word at its end node


class Trie:
    def __init__(self) -> None:
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for ch in word:
            node = node.children.setdefault(ch, TrieNode())
            node.count += 1
        node.is_end = True
        node.word = word

    def _walk(self, s: str) -> TrieNode | None:
        node = self.root
        for ch in s:
            node = node.children.get(ch)
            if node is None:
                return None
        return node

    def search(self, word: str) -> bool:
        node = self._walk(word)
        return node is not None and node.is_end

    def starts_with(self, prefix: str) -> bool:
        return self._walk(prefix) is not None

    def count_prefix(self, prefix: str) -> int:
        if not prefix:
            return 0
        node = self._walk(prefix)
        return node.count if node else 0
```

**Complexity:** O(L) per operation; O(total characters inserted) space.

**Part (b): wildcard search.**

**Approach:** recurse: a literal follows one edge; `.` tries every child. Exponential in the number of dots, but each level is bounded by the alphabet and words are short.

```python
class WordDictionary:
    def __init__(self) -> None:
        self.root = TrieNode()

    def add_word(self, word: str) -> None:
        node = self.root
        for ch in word:
            node = node.children.setdefault(ch, TrieNode())
        node.is_end = True

    def search(self, word: str) -> bool:
        def dfs(node: TrieNode, i: int) -> bool:
            if i == len(word):
                return node.is_end
            ch = word[i]
            if ch == ".":
                return any(dfs(child, i + 1) for child in node.children.values())
            nxt = node.children.get(ch)
            return nxt is not None and dfs(nxt, i + 1)
        return dfs(self.root, 0)
```

**Complexity:** insert O(L). Search O(L) without dots, O(26^d · L) worst case with d dots.

**Edge cases:** after inserting only `"apple"`, `search("ap")` is False (path exists, `is_end` false) while `startsWith("ap")` is True. `search("")` is True only if `""` was inserted.

</details>

#### Example 1.2: Autocomplete (Search Suggestions System) [Medium]

**Problem:** [Search Suggestions System (1268)](https://leetcode.com/problems/search-suggestions-system/) · **Template:** Template 1 (Trie, node class) with a `top` payload

Problem: given product names and a search word, after each typed character return up to three lexicographically smallest products with that prefix.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** a prefix lookup after every typed character against a fixed product list: a trie whose nodes carry a `top` payload of the best completions.

**Approach:** insert products in sorted order; each node keeps the first k words through it, which are the k smallest, so each query step is one edge walk. Trie-free alternative: sort, `bisect_left` the prefix (`07-binary-search.md`), check the next three entries.

```python
class SuggestNode:
    __slots__ = ("children", "top")

    def __init__(self) -> None:
        self.children: dict[str, SuggestNode] = {}
        self.top: list[str] = []          # up to k best words through this node


def suggested_products(products: list[str], search_word: str, k: int = 3) -> list[list[str]]:
    root = SuggestNode()
    for p in sorted(products):            # sorted insert => first k seen are the k smallest
        node = root
        for ch in p:
            node = node.children.setdefault(ch, SuggestNode())
            if len(node.top) < k:
                node.top.append(p)
    out, node = [], root
    for ch in search_word:
        node = node.children.get(ch) if node else None
        out.append(list(node.top) if node else [])
    return out
```

**Complexity:** builds in O(S log n + total chars), answers in O(len(search_word) · k).

**Edge cases:** once a prefix has no match, no longer prefix does (`node` stays `None`). For top-by-frequency autocomplete, store `(−count, word)` in a bounded sorted list or heap per node; updates must refresh every node on the path.

</details>

### Template 2: Trie as dict-of-dicts

> - **Example problem:** [Word Search II (212)](https://leetcode.com/problems/word-search-ii/)
> - **Why this template:** Many words in one grid: one DFS walks the trie, and storing the word under `END` lets it report and prune without rebuilding strings.
> - **Worked below:** Example 2.3
> - **Also fits:** [Replace Words (648)](https://leetcode.com/problems/replace-words/) (Example 2.2)

```python
END = "$"   # sentinel key; pick a char that cannot appear in words


def dict_trie_build(words: list[str]) -> dict:
    root: dict = {}
    for w in words:
        node = root
        for ch in w:
            node = node.setdefault(ch, {})
        node[END] = w            # store the word itself (handy for word search II)
    return root


def dict_trie_search(root: dict, word: str) -> bool:
    node = root
    for ch in word:
        if ch not in node:
            return False
        node = node[ch]
    return END in node
```

Dict-of-dicts is 8 lines and fast under pressure; the class is clearer when nodes carry several fields or methods. Be fluent in one.

#### Example 2.1: Maximum XOR of two numbers (brief) [Medium]

**Problem:** [Maximum XOR of Two Numbers in an Array (421)](https://leetcode.com/problems/maximum-xor-of-two-numbers-in-an-array/) · **Template:** Template 2 (Trie as dict-of-dicts), keyed by bits instead of characters

Problem: maximum of `a ^ b` over pairs in an array of non-negative ints.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "maximize XOR": decide bits high to low.

**Approach:** a higher bit outweighs all lower bits combined, so greedily prefer the opposite bit at each level. A binary trie answers "is there a number with this bit prefix?" in O(bits).

```python
def find_maximum_xor(nums: list[int], bits: int = 31) -> int:
    root: dict = {}
    for x in nums:
        node = root
        for b in range(bits - 1, -1, -1):
            node = node.setdefault((x >> b) & 1, {})
    best = 0
    for x in nums:
        node, cur = root, 0
        for b in range(bits - 1, -1, -1):
            bit = (x >> b) & 1
            want = 1 - bit                # greedy: prefer opposite bit at high positions
            if want in node:
                cur |= 1 << b
                node = node[want]
            else:
                node = node[bit]
        best = max(best, cur)
    return best
```

**Complexity:** O(n · B) time and space, B = 31.

**Edge cases:** a single number gives 0 (`x ^ x`); all-equal values give 0; `bits` must cover the largest value.

Extension: Maximum XOR With an Element From Array (sort queries by limit offline, insert numbers incrementally).

</details>

#### Example 2.2: Replace Words [Medium]

**Problem:** [Replace Words (648)](https://leetcode.com/problems/replace-words/) · **Template:** Template 2 (Trie as dict-of-dicts)

Problem: given dictionary roots and a sentence, replace each word with its shortest root prefix (keep words with no root).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** shortest dictionary prefix of each word: a prefix lookup against a fixed dictionary, so walk a trie along the word.

**Approach:** walk the trie along the word; the first `END` marker is the shortest root. Stop at a missing edge.

```python
def replace_words(roots: list[str], sentence: str) -> str:
    trie = dict_trie_build(roots)

    def shortest_root(word: str) -> str:
        node = trie
        for i, ch in enumerate(word):
            if ch not in node:
                break
            node = node[ch]
            if END in node:
                return word[: i + 1]      # first end marker = shortest root
        return word

    return " ".join(shortest_root(w) for w in sentence.split())
```

**Complexity:** O(total root chars + sentence length).

**Edge cases:** a word with no root is kept unchanged (the walk breaks at a missing edge); with roots `ca` and `cat`, the first `END` reached returns the shorter `ca`.

</details>

#### Example 2.3: Word Search II [Hard]

**Problem:** [Word Search II (212)](https://leetcode.com/problems/word-search-ii/) · **Template:** Template 2 (Trie as dict-of-dicts)

Problem: given a letter grid and a word list, return every word traceable through horizontally or vertically adjacent cells, each cell used at most once per word.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** many words, one grid, backtracking. Word Search I per word costs O(W · R · C · 4^L); a trie lets one DFS explore all words and abandon a path once no word has that prefix.

**Approach:** build a dict trie storing each full word at its end. DFS from each cell whose letter is a root child. At each step:
1. Move to this letter's child node.
2. If it holds a word, record it and **remove** the marker so it's reported once.
3. Mark the cell visited (`#`), recurse into neighbors whose letter is a child, restore.
4. **Prune**: if the child node is now empty, delete it from its parent so later calls skip it. This makes it fast on adversarial inputs.

```python
def find_words(board: list[list[str]], words: list[str]) -> list[str]:
    trie = dict_trie_build(words)
    rows, cols = len(board), len(board[0])
    found: list[str] = []

    def dfs(r: int, c: int, parent: dict) -> None:
        ch = board[r][c]
        node = parent[ch]
        if END in node:
            found.append(node.pop(END))  # pop: report each word once
        board[r][c] = "#"                 # mark visited
        for dr, dc in ((1, 0), (-1, 0), (0, 1), (0, -1)):
            nr, nc = r + dr, c + dc
            if 0 <= nr < rows and 0 <= nc < cols and board[nr][nc] in node:
                dfs(nr, nc, node)
        board[r][c] = ch                  # restore
        if not node:                      # prune: leaf with no remaining words
            parent.pop(ch)

    for r in range(rows):
        for c in range(cols):
            if board[r][c] in trie:
                dfs(r, c, trie)
    return found
```

**Complexity:** O(R · C · 4 · 3^(L-1)) worst case, L = max word length (after the first step at most 3 directions are unvisited); pruning makes typical runs far faster. Space O(total word chars) plus O(L) recursion.

**Edge cases:** the `#` sentinel must not be a trie letter. The trie deduplicates repeated input words. Mutating the board is fine if you restore it; mention it if input mutation matters. General choose/explore/unchoose shape: `15-backtracking.md`.

</details>

### Template 3: expand around center

> - **Example problem:** [Longest Palindromic Substring (5)](https://leetcode.com/problems/longest-palindromic-substring/)
> - **Why this template:** Every contiguous palindrome has a center at a character or gap, so expanding from 2n-1 centers finds the longest in O(n^2) time and O(1) space.
> - **Worked below:** Example 3.1
> - **Also fits:** [Palindromic Substrings (647)](https://leetcode.com/problems/palindromic-substrings/)

```python
def expand(s: str, lo: int, hi: int) -> tuple[int, int]:
    """Grow [lo, hi] while it is a palindrome; return the final inclusive bounds."""
    while lo >= 0 and hi < len(s) and s[lo] == s[hi]:
        lo -= 1
        hi += 1
    return lo + 1, hi - 1
```

Call with `(i, i)` for odd and `(i, i + 1)` for even centers. Knob: the per-step action (track longest, count, check a property).

#### Example 3.1: Palindromic substrings [Medium]

**Problems:** [Longest Palindromic Substring (5)](https://leetcode.com/problems/longest-palindromic-substring/), [Palindromic Substrings (647)](https://leetcode.com/problems/palindromic-substrings/) · **Template:** Template 3 (expand around center)

Problem A: return the longest contiguous palindromic substring. Problem B: count all palindromic substrings (by position).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "palindromic substring" (contiguous): expand around each of the 2n-1 centers.

**Approach:** expand while the ends match. Each successful step is a distinct palindrome, so B's count comes free. For A, keep the widest bounds.

```python
def longest_palindrome(s: str) -> str:
    if not s:
        return ""
    best_lo, best_hi = 0, 0
    for i in range(len(s)):
        for lo, hi in (expand(s, i, i), expand(s, i, i + 1)):   # odd, even centers
            if hi - lo > best_hi - best_lo:
                best_lo, best_hi = lo, hi
    return s[best_lo: best_hi + 1]


def count_substrings(s: str) -> int:
    """Count palindromic substrings: each successful expansion step is one palindrome."""
    total = 0
    for center in range(2 * len(s) - 1):
        lo, hi = center // 2, center // 2 + center % 2
        while lo >= 0 and hi < len(s) and s[lo] == s[hi]:
            total += 1
            lo -= 1
            hi += 1
    return total
```

**Complexity:** O(n^2) time, O(1) space. The DP table (`dp[i][j] = s[i] == s[j] and dp[i+1][j-1]`) is also O(n^2) time but O(n^2) space; prefer expansion.

**Manacher (mention only):** reuses mirror images inside the rightmost palindrome for O(n): "expand around center with memoized radii". Rarely required; the companion file's `manacher` is cross-checked against brute force.

**Edge cases:** an even center `(i, i + 1)` with `i + 1 == n` returns `(i + 1, i)`, negative width, never wins. Single-character strings. All-same characters (`"aaaa"` has 10 palindromic substrings).

</details>

### Template 4: prefix function (KMP)

> - **Example problem:** [Find the Index of the First Occurrence in a String (28)](https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string/)
> - **Why this template:** "First occurrence of needle in haystack" is exact matching; KMP does it in O(n + m) without backing up in the text.
> - **Worked below:** Example 4.1
> - **Also fits:** [Repeated Substring Pattern (459)](https://leetcode.com/problems/repeated-substring-pattern/)

```python
def prefix_function(s: str) -> list[int]:
    """pi[i] = length of longest proper prefix of s[:i+1] that is also its suffix."""
    pi = [0] * len(s)
    k = 0                                  # length of current matched border
    for i in range(1, len(s)):
        while k > 0 and s[i] != s[k]:
            k = pi[k - 1]                  # fall back to next shorter border
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    return pi


def kmp_search(text: str, pattern: str) -> list[int]:
    """All start indices of pattern in text, O(n + m)."""
    if not pattern:
        return list(range(len(text) + 1))
    pi = prefix_function(pattern)
    res, k = [], 0
    for i, ch in enumerate(text):
        while k > 0 and ch != pattern[k]:
            k = pi[k - 1]
        if ch == pattern[k]:
            k += 1
        if k == len(pattern):
            res.append(i - k + 1)
            k = pi[k - 1]                  # allow overlapping matches
    return res
```

Why linear: `k` rises at most 1 per character and every fallback lowers it, so total fallbacks are bounded by total increases. Alternative to `kmp_search`: `prefix_function(pattern + "#" + text)`, reporting positions where `pi == len(pattern)` (the separator must appear in neither string).

Three uses to know cold:
- **Matching**: as above.
- **Period**: the smallest period of s is `n - pi[-1]`; s is a full repetition exactly when that period divides n and is less than n.
- **Longest palindromic prefix**: `pi[-1]` of `s + "#" + reverse(s)`.

#### Example 4.1: Prefix function applications [Easy/Hard]

**Problems:** [Find the Index of the First Occurrence in a String (28)](https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string/), [Repeated Substring Pattern (459)](https://leetcode.com/problems/repeated-substring-pattern/), [Shortest Palindrome (214)](https://leetcode.com/problems/shortest-palindrome/) · **Template:** Template 4 (prefix function, KMP)

Problem A: first index of needle in haystack, or -1.
Problem B: can s be built by repeating one of its proper substrings?
Problem C: add the fewest characters to the **front** of s to make it a palindrome; return the result.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** A is matching. B is periodicity. C needs the longest palindromic prefix of s, and "prefix of s equal to a suffix of reverse(s)" is a border question.

**Approach:** B: with smallest period `p = n - pi[-1]`, if `pi[-1] > 0` and `p` divides `n`, s is `s[:p]` repeated `n / p` times (one-liner: `s in (s + s)[1:-1]`). C: in `t = s + "#" + reverse(s)`, `pi[-1]` is the longest prefix of s equal to a suffix of reverse(s), which is a palindromic prefix. Prepend the reverse of the rest.

```python
def str_str(haystack: str, needle: str) -> int:
    hits = kmp_search(haystack, needle)
    return hits[0] if hits else -1


def repeated_substring_pattern(s: str) -> bool:
    n = len(s)
    border = prefix_function(s)[-1] if s else 0
    period = n - border
    return border > 0 and n % period == 0


def shortest_palindrome(s: str) -> str:
    """Prepend the fewest chars to make s a palindrome.
    Longest palindromic prefix = border of s + '#' + reverse(s)."""
    rev = s[::-1]
    k = prefix_function(s + "#" + rev)[-1]
    return rev[: len(s) - k] + s
```

**Complexity:** all O(n + m) time and space.

**Edge cases:** the `#` separator stops a border spanning both halves; without it, `"aa"` gives `pi[-1] = 3` on `"aaaa"`. For A, `haystack.find(needle)` is fine unless the interviewer wants the algorithm; say so.

</details>

### Template 5: rolling hash (Rabin-Karp)

> - **Example problem:** [Repeated DNA Sequences (187)](https://leetcode.com/problems/repeated-dna-sequences/)
> - **Why this template:** "All the 10-letter-long sequences that occur more than once": fixed-length windows whose hash updates in O(1) per slide.
> - **Worked below:** Example 5.1
> - **Also fits:** [Longest Duplicate Substring (1044)](https://leetcode.com/problems/longest-duplicate-substring/)

```python
MOD = (1 << 61) - 1
BASE = 131


def rolling_hashes(s: str, L: int, base: int = BASE, mod: int = MOD):
    """Yield (start, hash) for every window of length L in O(n)."""
    if L == 0 or L > len(s):
        return
    h = 0
    for ch in s[:L]:
        h = (h * base + ord(ch)) % mod
    yield 0, h
    top = pow(base, L - 1, mod)            # weight of the char leaving the window
    for i in range(L, len(s)):
        h = (h - ord(s[i - L]) * top) % mod
        h = (h * base + ord(s[i])) % mod
        yield i - L + 1, h
```

Knobs:
- `MOD`: a large prime. `2^61 - 1` makes collisions astronomically rare; Python big ints keep the multiplication safe. Fixed-width languages use `10^9 + 7` or double hashing.
- `BASE`: larger than the alphabet; a random base defends against adversarial inputs.
- **Verify on match** (compare actual substrings) when correctness matters. If asked about collisions: "I verify on hash match, or use two moduli".
- For arbitrary substring hashes, precompute prefix hashes `H[i]` and powers `P[i]`: `hash(l, r) = (H[r] - H[l] * P[r - l]) mod M`.

#### Example 5.1: Rolling hash [Medium/Hard]

**Problems:** [Repeated DNA Sequences (187)](https://leetcode.com/problems/repeated-dna-sequences/), [Longest Duplicate Substring (1044)](https://leetcode.com/problems/longest-duplicate-substring/) · **Template:** Template 5 (rolling hash, Rabin-Karp)

Problem A: return all length-10 substrings that occur more than once.
Problem B: return any longest substring that occurs at least twice (occurrences may overlap); n up to 3 · 10^4.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** A compares fixed-length windows: rolling hash (or slicing into a set, since L = 10 is small). B asks for the longest L with a duplicate.

**Approach:** existence is **monotone** in L (a length-L duplicate contains duplicates of every shorter length), so binary search on L with an O(n) rolling-hash check (`07-binary-search.md`).

```python
def rabin_karp(text: str, pattern: str) -> list[int]:
    if not pattern or len(pattern) > len(text):
        return []
    _, target = next(rolling_hashes(pattern, len(pattern)))
    # verify on hash match to rule out collisions
    return [i for i, h in rolling_hashes(text, len(pattern))
            if h == target and text[i: i + len(pattern)] == pattern]


def find_repeated_dna_sequences(s: str, L: int = 10) -> list[str]:
    seen, out = set(), set()
    for i, h in rolling_hashes(s, L):
        if h in seen:
            out.add(s[i: i + L])
        seen.add(h)
    return sorted(out)


def longest_dup_substring(s: str) -> str:
    """Binary search on length + rolling hash. Monotone: a dup of length L
    implies a dup of length L-1, so the predicate is T...T F...F."""
    def dup_of_len(L: int) -> int:
        seen: dict[int, list[int]] = defaultdict(list)
        for i, h in rolling_hashes(s, L):
            for j in seen[h]:              # verify to be collision-safe
                if s[j: j + L] == s[i: i + L]:
                    return i
            seen[h].append(i)
        return -1

    lo, hi, best = 1, len(s) - 1, ""
    while lo <= hi:
        mid = (lo + hi) // 2
        i = dup_of_len(mid)
        if i != -1:
            best = s[i: i + mid]
            lo = mid + 1
        else:
            hi = mid - 1
    return best
```

**Complexity:** A is O(n). B is O(n log n) expected (each check O(n); verification only on rare collisions).

**Edge cases:** for A, packing each base into 2 bits gives an exact 20-bit key (`20-bit-manipulation-and-math.md`). In B, `hi = n - 1` because a duplicate can't be the whole string. Subtracting before `% mod` is safe in Python (`%` is non-negative); in C++/Java add `mod` first.

</details>

### Beyond the templates

#### Example B.1: Encoding, compression, signatures, tokenizing [Medium]

**Problems:** [Encode and Decode Strings (271)](https://leetcode.com/problems/encode-and-decode-strings/) (Premium), [String Compression (443)](https://leetcode.com/problems/string-compression/), [Group Anagrams (49)](https://leetcode.com/problems/group-anagrams/) · **Template:** none of the above; string-processing mechanics (length prefixes, read/write pointers, signatures)

**Why no template:** These are string-processing mechanics (length prefixes, read/write pointers, signatures, scanning), not trie or matching algorithms.

Problem A (LeetCode: Encode and Decode Strings): serialize a list of strings to one string and back; strings may contain any character.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** strings may contain any character, so no delimiter is safe alone; a length prefix means the decoder never interprets content.

**Approach:** write `len(s)`, `#`, then `s`. Decode reads digits up to the first `#` (lengths never contain `#`), then exactly `len` characters, whatever they are. Real wire formats do this (HTTP/2 frames, protobuf length-delimited fields, Redis RESP).

```python
def encode(strs: list[str]) -> str:
    return "".join(f"{len(s)}#{s}" for s in strs)


def decode(data: str) -> list[str]:
    out, i = [], 0
    while i < len(data):
        j = data.index("#", i)             # length digits never contain '#'
        n = int(data[i:j])
        out.append(data[j + 1: j + 1 + n])
        i = j + 1 + n
    return out
```

**Edge cases:** empty strings (`"0#"`), strings containing `#` or digits, and an empty list (encodes to `""`) all round-trip.

</details>

Problem B (LeetCode: String Compression): run-length encode a char array in place; return the new length.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** read/write two pointers (`05-two-pointers.md`); write never overtakes read because an encoded run is never longer than the run itself.

```python
def compress(chars: list[str]) -> int:
    """In-place run-length compression ("aaab" -> "a3b"); returns new length."""
    write = read = 0
    while read < len(chars):
        ch, start = chars[read], read
        while read < len(chars) and chars[read] == ch:
            read += 1
        chars[write] = ch
        write += 1
        if read - start > 1:
            for d in str(read - start):
                chars[write] = d
                write += 1
    return write
```

</details>

Problem C (LeetCode: Group Anagrams): group words that are anagrams.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** key = a 26-count tuple (O(L)) or `"".join(sorted(w))` (O(L log L)). More signature tricks (shift-normalization, pattern isomorphism): `04-arrays-hashing-prefix-sums.md`.

```python
def group_anagrams(strs: list[str]) -> list[list[str]]:
    groups: dict[tuple[int, ...], list[str]] = defaultdict(list)
    for w in strs:
        sig = [0] * 26
        for ch in w:
            sig[ord(ch) - 97] += 1
        groups[tuple(sig)].append(w)
    return list(groups.values())
```

</details>

Problem D: tokenizing.

<details>
<summary>Decode, approach, and solution</summary>

**Approach:** calculator and parser problems start by tokenizing; people slip on scanning multi-digit numbers with an inner `while`. Evaluation then uses a stack or recursive descent (`08-stacks-and-monotonic-structures.md`, Basic Calculator family).

```python
def tokenize(expr: str) -> list[str]:
    """Split an arithmetic expression into numbers, operators, parens."""
    tokens, i = [], 0
    while i < len(expr):
        ch = expr[i]
        if ch.isspace():
            i += 1
        elif ch.isdigit():
            j = i
            while j < len(expr) and expr[j].isdigit():
                j += 1
            tokens.append(expr[i:j])
            i = j
        elif ch in "+-*/()":
            tokens.append(ch)
            i += 1
        else:
            raise ValueError(f"bad char {ch!r} at {i}")
    return tokens
```

**Complexity:** all O(total length).

</details>

## Template cheat sheet

All templates in one place, for warm-ups.

### Template 1: Trie (node class)
```python
class TrieNode:
    __slots__ = ("children", "is_end", "count", "word")

    def __init__(self) -> None:
        self.children: dict[str, TrieNode] = {}
        self.is_end = False      # a word ends exactly here
        self.count = 0           # number of words passing through (prefix count)
        self.word: str | None = None   # optional: store the word at its end node


class Trie:
    def __init__(self) -> None:
        self.root = TrieNode()

    def insert(self, word: str) -> None:
        node = self.root
        for ch in word:
            node = node.children.setdefault(ch, TrieNode())
            node.count += 1
        node.is_end = True
        node.word = word

    def _walk(self, s: str) -> TrieNode | None:
        node = self.root
        for ch in s:
            node = node.children.get(ch)
            if node is None:
                return None
        return node

    def search(self, word: str) -> bool:
        node = self._walk(word)
        return node is not None and node.is_end

    def starts_with(self, prefix: str) -> bool:
        return self._walk(prefix) is not None

    def count_prefix(self, prefix: str) -> int:
        if not prefix:
            return 0
        node = self._walk(prefix)
        return node.count if node else 0
```

### Template 2: Trie as dict-of-dicts
```python
END = "$"   # sentinel key; pick a char that cannot appear in words


def dict_trie_build(words: list[str]) -> dict:
    root: dict = {}
    for w in words:
        node = root
        for ch in w:
            node = node.setdefault(ch, {})
        node[END] = w            # store the word itself (handy for word search II)
    return root


def dict_trie_search(root: dict, word: str) -> bool:
    node = root
    for ch in word:
        if ch not in node:
            return False
        node = node[ch]
    return END in node
```

### Template 3: expand around center
```python
def expand(s: str, lo: int, hi: int) -> tuple[int, int]:
    """Grow [lo, hi] while it is a palindrome; return the final inclusive bounds."""
    while lo >= 0 and hi < len(s) and s[lo] == s[hi]:
        lo -= 1
        hi += 1
    return lo + 1, hi - 1
```

### Template 4: prefix function (KMP)
```python
def prefix_function(s: str) -> list[int]:
    """pi[i] = length of longest proper prefix of s[:i+1] that is also its suffix."""
    pi = [0] * len(s)
    k = 0                                  # length of current matched border
    for i in range(1, len(s)):
        while k > 0 and s[i] != s[k]:
            k = pi[k - 1]                  # fall back to next shorter border
        if s[i] == s[k]:
            k += 1
        pi[i] = k
    return pi


def kmp_search(text: str, pattern: str) -> list[int]:
    """All start indices of pattern in text, O(n + m)."""
    if not pattern:
        return list(range(len(text) + 1))
    pi = prefix_function(pattern)
    res, k = [], 0
    for i, ch in enumerate(text):
        while k > 0 and ch != pattern[k]:
            k = pi[k - 1]
        if ch == pattern[k]:
            k += 1
        if k == len(pattern):
            res.append(i - k + 1)
            k = pi[k - 1]                  # allow overlapping matches
    return res
```

### Template 5: rolling hash (Rabin-Karp)
```python
MOD = (1 << 61) - 1
BASE = 131


def rolling_hashes(s: str, L: int, base: int = BASE, mod: int = MOD):
    """Yield (start, hash) for every window of length L in O(n)."""
    if L == 0 or L > len(s):
        return
    h = 0
    for ch in s[:L]:
        h = (h * base + ord(ch)) % mod
    yield 0, h
    top = pow(base, L - 1, mod)            # weight of the char leaving the window
    for i in range(L, len(s)):
        h = (h - ord(s[i - L]) * top) % mod
        h = (h * base + ord(s[i])) % mod
        yield i - L + 1, h
```

## Variations and extensions

- **Trie + DFS/backtracking**: Word Search II, Concatenated Words (DFS over the trie with memo), Palindrome Pairs (trie of reversed words).
- **Trie + heap / counts**: top-k autocomplete by frequency, Map Sum Pairs (prefix sums).
- **Trie with deletion**: decrement `count` along the path; delete children whose count hits 0.
- **Compressed trie (radix tree)**: collapse single-child chains; used by IP routing tables and Kubernetes/etcd key-prefix indexes. Talking point, rarely coded.
- **Aho-Corasick**: trie + KMP-style failure links, matching many patterns in one pass. Mention it for "scan logs for 10^4 banned terms".
- **Z-function**: `z[i]` = longest common prefix of s and s[i:]. Interchangeable with the prefix function for most problems.
- **Rolling hash + binary search**: longest duplicate substring, longest common substring of two strings, Longest Happy Prefix (also prefix function).
- **Content-defined chunking** (infra): rolling hash over a sliding window; cut where `hash % K == 0` (rsync, dedup backups).
- **Suffix arrays / automata**: out of scope for interviews.

Follow-ups:
- "Trie memory is tight" → 26-arrays only for small dense alphabets; otherwise a radix tree, or sorted list + binary search.
- "Words stream in, query top suggestions" → per-node bounded heaps; update O(L · log k).
- "Hash collisions?" → verify on match, 61-bit Mersenne prime modulus, or double hash.

## Common pitfalls

- Trie: confusing `search` with `startsWith` (forgetting `is_end`); accidentally marking the root `is_end` for the empty string.
- Word Search II: not removing found words (duplicate output), not restoring the cell, not pruning (TLE on adversarial inputs).
- Expand around center: forgetting even-length centers.
- Prefix function: starting at `i = 0` (pi[0] is always 0), or `if` instead of `while` for the fallback.
- KMP search: forgetting `k = pi[k - 1]` after a full match (misses overlaps or reads out of range).
- Rolling hash: removal weight `base^L` instead of `base^(L-1)`; negative values before `%` in other languages; calling built-in `hash()` on slices (O(L) per window, not O(1)).
- Encoding with a plain delimiter (`",".join`) breaks when strings contain it.
- Building strings with `+=` in loops is O(n^2) worst case; collect into a list and `"".join`.

## When NOT to use / look-alikes

| Situation | Use instead | Why |
|---|---|---|
| Exact lookups only, no prefix queries | `set` / `dict` | A trie pays off only for prefix structure |
| Palindromic **subsequence** | DP on intervals | Not contiguous; centers don't apply |
| One pattern, one text, just needs to work | `str.find` / `in` | Say "built-in is fine unless you want KMP" |
| Substring with at most k distinct / no repeats | Sliding window (`06-sliding-window.md`) | Constraint on window contents, not substring equality |
| Edit distance, LCS | DP (`16`, `17`) | Alignment, not exact matching |
| Anagram of pattern inside text | Fixed sliding window with counts | Order doesn't matter, so sequence hashes are wrong |
| Max XOR of a subarray | Prefix XOR + binary trie | Still a trie, but over prefix XORs (cross-ref `04`) |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Implement Trie (Prefix Tree) (208)](https://leetcode.com/problems/implement-trie-prefix-tree/) · Medium · Trie basics: node class, `is_end`.
2. [Design Add and Search Words Data Structure (211)](https://leetcode.com/problems/design-add-and-search-words-data-structure/) · Medium · Trie + wildcard DFS: DFS on `.`.
3. [Word Search II (212)](https://leetcode.com/problems/word-search-ii/) · Hard · Trie + grid backtracking: trie-driven DFS, pop found words, prune empty nodes.
4. [Longest Palindromic Substring (5)](https://leetcode.com/problems/longest-palindromic-substring/) · Medium · Expand around center: 2n-1 centers.
5. [Palindromic Substrings (647)](https://leetcode.com/problems/palindromic-substrings/) · Medium · Expand around center (count): count expansion steps.

### Optional

**Tries**
- [Replace Words (648)](https://leetcode.com/problems/replace-words/) · Medium · First `is_end` on the walk.
- [Search Suggestions System (1268)](https://leetcode.com/problems/search-suggestions-system/) · Medium · Sorted insert, top-3 per node.
- [Map Sum Pairs (677)](https://leetcode.com/problems/map-sum-pairs/) · Medium · Store sums along the path; handle key overwrite with a delta.
- [Maximum XOR of Two Numbers in an Array (421)](https://leetcode.com/problems/maximum-xor-of-two-numbers-in-an-array/) · Medium · Binary trie, greedy high bit.
- [Palindrome Pairs (336)](https://leetcode.com/problems/palindrome-pairs/) · Hard · Trie of reversed words or hash map of reversed words with split checks.

**Palindromes**
- [Valid Palindrome (125)](https://leetcode.com/problems/valid-palindrome/) · Easy · Two pointers skipping non-alphanumerics.
- [Valid Palindrome II (680)](https://leetcode.com/problems/valid-palindrome-ii/) · Easy · On first mismatch, try skipping either side once.

**Prefix function / matching**
- [Find the Index of the First Occurrence in a String (28)](https://leetcode.com/problems/find-the-index-of-the-first-occurrence-in-a-string/) · Easy · KMP or built-in.
- [Repeated Substring Pattern (459)](https://leetcode.com/problems/repeated-substring-pattern/) · Easy · Period `n - pi[-1]` divides n.
- [Longest Happy Prefix (1392)](https://leetcode.com/problems/longest-happy-prefix/) · Hard · Literally `pi[-1]`.
- [Shortest Palindrome (214)](https://leetcode.com/problems/shortest-palindrome/) · Hard · `pi` of `s + "#" + rev(s)`.

**Rolling hash**
- [Repeated DNA Sequences (187)](https://leetcode.com/problems/repeated-dna-sequences/) · Medium · Fixed window hash (or 2-bit encoding).
- [Longest Duplicate Substring (1044)](https://leetcode.com/problems/longest-duplicate-substring/) · Hard · Binary search on length + Rabin-Karp.

**Encoding / signatures / parsing**
- [Group Anagrams (49)](https://leetcode.com/problems/group-anagrams/) · Medium · Count-tuple key.
- [Encode and Decode Strings (271)](https://leetcode.com/problems/encode-and-decode-strings/) · Medium · Premium · Length prefix.
- [String Compression (443)](https://leetcode.com/problems/string-compression/) · Medium · Read/write pointers.
- [Decode String (394)](https://leetcode.com/problems/decode-string/) · Medium · Stack of (count, partial string).

## Self-check

1. What is the difference between `search` and `startsWith` in a trie, and which node field distinguishes them?
2. In Word Search II, what two modifications to the trie during DFS make it fast and correct?
3. Why are there 2n-1 centers for palindromes, and what does each successful expansion step correspond to?
4. Define `pi[i]`. Why is computing the prefix function O(n) despite the inner `while`?
5. How do you get the smallest period of a string from its prefix function? How does that solve Repeated Substring Pattern?
6. In a rolling hash, what is removed from the hash when the window slides, and why is the weight `base^(L-1)`?
7. Why is Longest Duplicate Substring binary-searchable on length, and what does each check cost?
