# 20. Bit Manipulation and Math

Why this matters: some interview questions just check whether you know the trick: XOR cancels pairs, `x & (x - 1)` drops the lowest set bit, a bitmask is a set, `pow(a, b, m)` is O(log b), reservoir sampling handles unbounded streams. Easy once seen, nearly impossible to derive under pressure. This module collects them with the reasoning that makes them stick, plus the matrix and geometry mechanics of "implementation" rounds.

Companion code: `code/20_bit_math.py`.

## At a glance

- **Use when:** "twice except one" in O(1) space, n ≤ 20 subsets, "mod 10^9 + 7", 32-bit overflow, random picks, in-place matrix, integer geometry.
- **Templates:**
  - Template 1: bit helpers and 32-bit emulation: XOR tricks, bit counts, add without `+`.
  - Template 2: subset enumeration with bitmasks: subsets, submasks, DP over masks.
  - Template 3: fast modular exponentiation: huge exponents, modular inverses.
  - Template 4: reservoir sampling, Fisher-Yates, weighted pick: randomized selection.
- **Typical complexity:** O(n) or O(32n) bit tricks; O(n · 2^n) bitmask DP; O(log e) powers.
- **Top trap:** unbounded Python ints: bit loops on negatives never end without `& 0xFFFFFFFF`.

## Core idea

Bits:

| Identity | Meaning |
|---|---|
| `x ^ x == 0`, `x ^ 0 == x`, XOR is commutative/associative | Pairs cancel in any order |
| `x & (x - 1)` | Clear the lowest set bit |
| `x & -x` | Isolate the lowest set bit (two's complement: `-x == ~x + 1`) |
| `(x >> i) & 1`, `x \| (1 << i)`, `x & ~(1 << i)`, `x ^ (1 << i)` | Get, set, clear, toggle bit i |
| `n > 0 and n & (n - 1) == 0` | n is a power of two |
| `mask >> i & 1` over `range(1 << n)` | Enumerate all subsets of n items |
| `sub = (sub - 1) & mask` | Step to the next smaller submask |

**Python ints are unbounded.** No overflow, no fixed sign bit: `-1` acts like infinitely many 1s, so `-1 >> 1 == -1` and `bin(-5)` is `'-0b101'`. To emulate 32 bits, mask with `& 0xFFFFFFFF` after every operation, then convert back to signed:

```
x & 0xFFFFFFFF            -> low 32 bits as an unsigned value in [0, 2^32)
v - 2^32 if v >= 2^31     -> reinterpret as signed 32-bit
```

Math:
- **Euclid**: `gcd(a, b) = gcd(b, a mod b)`. O(log min(a, b)).
- **Fast exponentiation**: over the exponent's bits, square the base each step and multiply it in when the bit is 1. O(log e).
- **Modular inverse** for prime p: `a^(p-2) mod p` (Fermat's little theorem). Enables division under a modulus, which nCr mod p needs.
- **Sieve**: each prime p crosses out its multiples starting at p*p. O(n log log n).

Randomization: an algorithm is uniform if **every outcome is equally likely**, usually shown by a telescoping product (reservoir: `k/(i+1) · (i+1)/(i+2) · ... = k/n`).

## Recognition signals

- "Every element appears twice/three times except one", "find the missing/duplicate number" in O(1) space: XOR or bit counting.
- "Without using + / * / division": bit-level arithmetic.
- n ≤ 20 (sometimes ≤ 16) and "choose a subset": bitmask enumeration or bitmask DP (2^20 ≈ 10^6).
- "Count 1 bits", "Hamming distance", "power of two/four": bit identities.
- "Return the answer modulo 10^9 + 7": modular arithmetic; division needs a modular inverse.
- "Count primes below n", n up to 10^6 or 10^7: sieve.
- "Number of paths/ways in a grid without obstacles": combinatorics nCr instead of DP.
- "32-bit signed integer", "return 0 if it overflows": explicit bound checks.
- "Stream of unknown length, pick uniformly", "shuffle", "pick with probability proportional to weight": reservoir, Fisher-Yates, prefix sums + bisect.
- "In place" with a 2D matrix: layer-by-layer, transpose+reverse, or the first row/column as storage.
- "Points on a line", "orientation", "is it convex": cross product, normalized integer slopes.

## Disguises

1. **Corrupted replica.** Every config blob ID appears on exactly two replicas except one that lost its pair. Find the orphan in one pass with O(1) memory.
   **Decode:** pairs cancel; XOR everything (Single Number).

2. **Feature flag combos.** 15 independent flags; score every combination, and for each combination also sum over all its sub-combinations.
   **Decode:** n ≤ 20 subsets: iterate `range(1 << n)`; sub-combinations via submask iteration, O(3^n) total.

3. **Sampling a log firehose.** Keep a uniform random sample of 1000 log lines from a stream of unknown length, in O(1000) memory.
   **Decode:** reservoir sampling.

4. **Weighted load balancer.** Backends have integer weights; route each request to backend i with probability `w[i] / sum(w)`.
   **Decode:** Random Pick with Weight: prefix sums + `bisect` on a uniform random number.

5. **Canary rotation for a heatmap.** Rotate an n x n dashboard bitmap tile 90 degrees clockwise without allocating another tile.
   **Decode:** Rotate Image: transpose, then reverse each row.

6. **Collinear sensor detection.** Given integer sensor coordinates, what is the most sensors on one straight line?
   **Decode:** Max Points on a Line: per anchor, group others by reduced `(dy, dx)` with a canonical sign. Never float slopes.

## Templates and worked examples

### Template 1: bit helpers and 32-bit emulation

> - **Example problem:** [Single Number (136)](https://leetcode.com/problems/single-number/)
> - **Why this template:** "Every element appears twice except for one" with constant extra space: XOR cancels pairs (and Single Number III splits on `lowest_set_bit`).
> - **Worked below:** Example 1.2
> - **Also fits:** [Sum of Two Integers (371)](https://leetcode.com/problems/sum-of-two-integers/) (the template's `get_sum`)

```python
MASK32 = 0xFFFFFFFF


def to_signed32(x: int) -> int:
    """Interpret the low 32 bits of x as a signed two's-complement int."""
    x &= MASK32
    return x - (1 << 32) if x & 0x80000000 else x


def get_bit(x: int, i: int) -> int:
    return (x >> i) & 1


def set_bit(x: int, i: int) -> int:
    return x | (1 << i)


def clear_bit(x: int, i: int) -> int:
    return x & ~(1 << i)


def toggle_bit(x: int, i: int) -> int:
    return x ^ (1 << i)


def lowest_set_bit(x: int) -> int:
    return x & -x
```

32-bit emulation (Sum of Two Integers, add without `+`): XOR gives the carry-free sum, `(a & b) << 1` the carries; loop until no carry. Without masking, a negative operand loops forever, since carries propagate into infinitely many 1 bits.

```python
def get_sum(a: int, b: int) -> int:
    """Add without + using xor (sum bits) and and<<1 (carries), 32-bit semantics."""
    a, b = a & MASK32, b & MASK32
    while b:
        a, b = (a ^ b) & MASK32, ((a & b) << 1) & MASK32
    return to_signed32(a)


def reverse_bits(n: int) -> int:
    out = 0
    for _ in range(32):
        out = (out << 1) | (n & 1)
        n >>= 1
    return out
```

#### Example 1.1: Counting bits and powers of two [Easy]

**Problems:** [Number of 1 Bits (191)](https://leetcode.com/problems/number-of-1-bits/), [Counting Bits (338)](https://leetcode.com/problems/counting-bits/), [Power of Two (231)](https://leetcode.com/problems/power-of-two/) · **Template:** Template 1 (bit helpers and 32-bit emulation)

**Problems.** (a) Count the set bits of a 32-bit unsigned integer. (b) Return the set-bit count of every i in 0..n, O(n) total. (c) Is an integer a power of two?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "set bits" and "power of two" concern the binary representation; the common tool is `x & (x - 1)` (clear the lowest set bit).

**Approach:**
- Kernighan: `x &= x - 1` removes one set bit per iteration, so the loop runs popcount(x) times, not 32.
- Counting Bits: `bits[i] = bits[i >> 1] + (i & 1)` (i's bits are i//2's plus its last bit). Alternative: `bits[i] = bits[i & (i - 1)] + 1`.
- Power of two: exactly one set bit, so clearing it gives zero. `n > 0` excludes 0 and negatives.

```python
def popcount_kernighan(x: int) -> int:
    """x & (x-1) clears the lowest set bit; loop runs popcount times."""
    count = 0
    while x:
        x &= x - 1
        count += 1
    return count


def count_bits(n: int) -> list[int]:
    """bits[i] = bits[i >> 1] + (i & 1)."""
    bits = [0] * (n + 1)
    for i in range(1, n + 1):
        bits[i] = bits[i >> 1] + (i & 1)
    return bits


def is_power_of_two(n: int) -> bool:
    return n > 0 and n & (n - 1) == 0
```

**Complexity:** O(popcount), O(n), O(1).

**Edge cases:** Kernighan never terminates on a negative Python int; mask first (`x & MASK32`). Built-ins: `bin(x).count("1")`, `int.bit_count()` (3.10+), `int.bit_length()`. Power of four: a power of two whose bit is in an even position, `n & 0x55555555 != 0`.

</details>

#### Example 1.2: XOR family [Easy/Medium]

**Problems:** [Single Number (136)](https://leetcode.com/problems/single-number/), [Single Number II (137)](https://leetcode.com/problems/single-number-ii/), [Single Number III (260)](https://leetcode.com/problems/single-number-iii/), [Missing Number (268)](https://leetcode.com/problems/missing-number/) · **Template:** Template 1 (bit helpers and 32-bit emulation)

Problem: (I) all values appear twice except one; find it. (II) all appear three times except one. (III) all appear twice except two; find both. (Missing) n distinct values from 0..n with one missing.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "appears k times except", "O(1) extra space". XOR cancels pairs; bit counting mod k handles any k.

**Approach:**
- I: XOR everything.
- II: for each of 32 bits, count the numbers with it set; the single has that bit exactly when the count isn't divisible by 3. Convert to signed at the end: we rebuilt only 32 bits, while Python negatives have infinitely many.
- III: XOR everything to get `a ^ b` (nonzero since a ≠ b). Any set bit, e.g. `diff = x & -x`, is where a and b differ. Partition by that bit; each group XORs down to one single.
- Missing: XOR all indices 0..n and all values; everything present cancels.

```python
def single_number(nums: list[int]) -> int:
    """Every value appears twice except one: pairs cancel under xor."""
    acc = 0
    for x in nums:
        acc ^= x
    return acc


def single_number_ii(nums: list[int]) -> int:
    """Every value appears three times except one: count each bit mod 3."""
    res = 0
    for b in range(32):
        cnt = sum((x >> b) & 1 for x in nums)
        if cnt % 3:
            res |= 1 << b
    return to_signed32(res)


def single_number_iii(nums: list[int]) -> list[int]:
    """Two values appear once, rest twice: split by any differing bit."""
    xor_all = 0
    for x in nums:
        xor_all ^= x
    diff = xor_all & -xor_all          # a bit where the two singles differ
    a = 0
    for x in nums:
        if x & diff:
            a ^= x
    return sorted([a, xor_all ^ a])


def missing_number(nums: list[int]) -> int:
    """nums has n distinct values from 0..n with one missing."""
    acc = len(nums)
    for i, x in enumerate(nums):
        acc ^= i ^ x
    return acc
```

**Complexity:** O(n) time (II is O(32n)), O(1) space.

**Edge cases:** in II, `(x >> b) & 1` reads a negative x's two's-complement bit correctly in Python, so rebuilding 32 bits plus `to_signed32` works. Missing Number alternative: `n*(n+1)//2 - sum(nums)` (overflow-prone outside Python). XOR swap (`a ^= b; b ^= a; a ^= b`) is trivia; it breaks when both refer to the same memory location.

</details>

### Template 2: subset enumeration with bitmasks

> - **Example problem:** [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/)
> - **Why this template:** with n <= 16 the set of used items fits in a bitmask, and `dp[mask]` builds from masks with one fewer bit.
> - **Worked below:** Example 2.1

```python
def all_subsets(items: list[T]) -> list[list[T]]:
    n = len(items)
    return [[items[i] for i in range(n) if mask >> i & 1] for mask in range(1 << n)]


def iter_submasks(mask: int) -> Iterator[int]:
    """All submasks of mask, descending, including mask and 0. Total over all
    masks of n bits is O(3^n)."""
    sub = mask
    while True:
        yield sub
        if sub == 0:
            return
        sub = (sub - 1) & mask
```

Why `(sub - 1) & mask` works: subtracting 1 flips `sub`'s lowest set bit and all bits below it; ANDing with `mask` drops bits outside the mask, giving the next smaller subset. Why O(3^n) overall: each bit is outside the mask, in the mask only, or in both mask and submask.

Knobs: what you compute per mask (sum, validity), and whether you need DP over masks (`dp[mask]` from `dp[mask without bit i]`): Partition to K Equal Sum Subsets, TSP-style Hamiltonian paths, assignment.

#### Example 2.1: Bitmask DP and Gray code [Medium]

**Problems:** [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/), [Gray Code (89)](https://leetcode.com/problems/gray-code/) · **Template:** Template 2 (subset enumeration with bitmasks)

Problem: can the array (n ≤ 16) be split into k groups of equal sum?

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** n ≤ 16 + "assign each item to a group": state = set of used items.

**Approach:** fill buckets one at a time. `dp[mask]` = current bucket's fill level after using exactly the items in `mask`, or -1 if unreachable. Item i is valid if it doesn't overflow the bucket; `% target` rolls over to a new bucket when one fills. Sorting lets us `break` once an item is too large. Answer: `dp[full] == 0`. One value per mask loses nothing, because the mask determines the fill level (its sum mod target).

```python
def can_partition_k_subsets(nums: list[int], k: int) -> bool:
    """Bitmask DP: dp[mask] = fill level of the current bucket using items in mask, or -1."""
    total = sum(nums)
    if total % k:
        return False
    target = total // k
    nums = sorted(nums)
    if nums and nums[-1] > target:
        return False
    n = len(nums)
    dp = [-1] * (1 << n)
    dp[0] = 0
    for mask in range(1 << n):
        if dp[mask] < 0:
            continue
        for i in range(n):
            if mask >> i & 1:
                continue
            if dp[mask] + nums[i] > target:
                break                       # sorted: later items also too big
            nxt = mask | (1 << i)
            if dp[nxt] < 0:
                dp[nxt] = (dp[mask] + nums[i]) % target
    return dp[(1 << n) - 1] == 0
```

**Complexity:** O(n · 2^n) time, O(2^n) space. Backtracking alternative: `15-backtracking.md`.

Gray code: `i ^ (i >> 1)` for i in `0 .. 2^n - 1` changes exactly one bit between consecutive values (the same sequence as reflect-and-prefix).

```python
def gray_code(n: int) -> list[int]:
    return [i ^ (i >> 1) for i in range(1 << n)]
```

**Edge cases:** a total not divisible by k, or an item larger than the target, returns False before the DP; k = 1 is always True. Gray code with n = 0 is `[0]`.

</details>

### Template 3: fast modular exponentiation

> - **Example problem:** [Pow(x, n) (50)](https://leetcode.com/problems/powx-n/)
> - **Why this template:** n up to 2^31 rules out n multiplications; square-and-multiply over n's bits gives O(log n).
> - **Worked below:** Example 3.1

```python
def mod_pow(base: int, exp: int, mod: int) -> int:
    """Binary exponentiation: square the base, multiply in when the bit is set."""
    result, base = 1 % mod, base % mod
    while exp:
        if exp & 1:
            result = result * base % mod
        base = base * base % mod
        exp >>= 1
    return result
```

In real solutions use built-in `pow(a, b, m)`; write it by hand when asked (or for `Pow(x, n)` on floats). `pow(a, -1, m)` gives a modular inverse (Python 3.8+).

#### Example 3.1: Number theory toolkit: primes, fast power, nCr mod p, Pascal's triangle [Easy/Medium]

**Problems:** [Count Primes (204)](https://leetcode.com/problems/count-primes/), [Pow(x, n) (50)](https://leetcode.com/problems/powx-n/), [Pascal's Triangle (118)](https://leetcode.com/problems/pascals-triangle/) · **Template:** Template 3 (fast modular exponentiation) for Pow(x, n) and modular inverses; the sieve and gcd need no template

**Problems.** (a) Count the primes strictly less than n. (b) Compute x raised to an integer power n (n may be negative) in O(log n) multiplications. (c) Generate the first numRows rows of Pascal's triangle. Also: nCr modulo a prime.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** primes, gcd, powers, "mod 10^9 + 7", and combinations are standard number-theory pieces; recognize which one is needed.

**Approach:** a toolkit, one piece at a time.

GCD/LCM: compute lcm as `a // gcd * b` to keep intermediates small (matters in fixed-width languages). Python has `math.gcd` and `math.lcm`.

```python
def gcd(a: int, b: int) -> int:
    while b:
        a, b = b, a % b
    return abs(a)


def lcm(a: int, b: int) -> int:
    return a // gcd(a, b) * b if a and b else 0
```

Pow(x, n) on floats, including negative n:

```python
def my_pow(x: float, n: int) -> float:
    if n < 0:
        x, n = 1 / x, -n
    result = 1.0
    while n:
        if n & 1:
            result *= x
        x *= x
        n >>= 1
    return result
```

Modular inverse (prime p only; for composite moduli use extended Euclid, or `pow(a, -1, m)` when `gcd(a, m) == 1`):

```python
def mod_inverse(a: int, p: int) -> int:
    """Fermat: a^(p-2) mod p, valid when p is prime and p does not divide a."""
    return pow(a, p - 2, p)
```

Sieve of Eratosthenes: start crossing at `i * i`, since smaller primes already crossed smaller multiples. Slice assignment on a `bytearray` keeps it fast.

```python
def sieve(n: int) -> list[int]:
    """All primes <= n in O(n log log n)."""
    if n < 2:
        return []
    is_p = bytearray([1]) * (n + 1)
    is_p[0] = is_p[1] = 0
    for i in range(2, math.isqrt(n) + 1):
        if is_p[i]:
            is_p[i * i:: i] = bytearray(len(range(i * i, n + 1, i)))
    return [i for i in range(n + 1) if is_p[i]]


def count_primes(n: int) -> int:
    """Primes strictly less than n."""
    return len(sieve(n - 1))


def prime_factors(n: int) -> dict[int, int]:
    """Trial division up to sqrt(n): O(sqrt n)."""
    out: dict[int, int] = {}
    d = 2
    while d * d <= n:
        while n % d == 0:
            out[d] = out.get(d, 0) + 1
            n //= d
        d += 1
    if n > 1:
        out[n] = out.get(n, 0) + 1
    return out
```

A leftover `n > 1` after the loop is the (at most one) prime factor larger than sqrt of the original.

nCr mod p: `C(n, r) = n! / (r! (n-r)!)`; division mod p becomes multiplication by inverses. Compute one inverse with Fermat, then fill the rest backward with `inv_fact[i-1] = inv_fact[i] * i`.

```python
MOD = 10**9 + 7


class Combinatorics:
    """Precompute factorials mod p for O(1) nCr queries."""

    def __init__(self, n: int, mod: int = MOD) -> None:
        self.mod = mod
        self.fact = [1] * (n + 1)
        for i in range(1, n + 1):
            self.fact[i] = self.fact[i - 1] * i % mod
        self.inv_fact = [1] * (n + 1)
        self.inv_fact[n] = pow(self.fact[n], mod - 2, mod)
        for i in range(n, 0, -1):             # inv_fact[i-1] = inv_fact[i] * i
            self.inv_fact[i - 1] = self.inv_fact[i] * i % mod

    def ncr(self, n: int, r: int) -> int:
        if r < 0 or r > n:
            return 0
        return self.fact[n] * self.inv_fact[r] % self.mod * self.inv_fact[n - r] % self.mod
```

Pascal's triangle is the DP form `C(n, r) = C(n-1, r-1) + C(n-1, r)`; use it when n is small (≤ ~1000) or the modulus isn't prime. Unique Paths on an m x n grid is `C(m + n - 2, m - 1)` (`math.comb`).

```python
def pascal(num_rows: int) -> list[list[int]]:
    rows: list[list[int]] = []
    for r in range(num_rows):
        row = [1] * (r + 1)
        for c in range(1, r):
            row[c] = rows[r - 1][c - 1] + rows[r - 1][c]
        rows.append(row)
    return rows
```

**Complexity:** gcd O(log min), pow O(log n), sieve O(n log log n), factorization O(sqrt n), nCr O(n) build and O(1) query.

**Edge cases:** `count_primes(n)` is 0 for n ≤ 2; `my_pow` with x = 0 and negative n divides by zero; the Fermat inverse needs p prime and p not dividing a; `ncr` returns 0 when r < 0 or r > n; `pascal(0)` is `[]`.

</details>

### Template 4: reservoir sampling, Fisher-Yates, weighted pick

> - **Example problem:** [Random Pick with Weight (528)](https://leetcode.com/problems/random-pick-with-weight/)
> - **Why this template:** "Probability of picking index i is w[i] / sum(w)" is prefix sums plus bisect: exactly `WeightedPicker`.
> - **Worked below:** Example 4.1
> - **Also fits:** [Linked List Random Node (382)](https://leetcode.com/problems/linked-list-random-node/) (reservoir, k = 1), [Shuffle an Array (384)](https://leetcode.com/problems/shuffle-an-array/) (Fisher-Yates)

```python
def reservoir_sample(stream: Iterable[T], k: int, rng: random.Random = random) -> list[T]:
    """Uniform k-sample from a stream of unknown length in O(k) memory."""
    res: list[T] = []
    for i, x in enumerate(stream):
        if i < k:
            res.append(x)
        else:
            j = rng.randint(0, i)          # inclusive: i+1 choices
            if j < k:
                res[j] = x                 # keep x with probability k/(i+1)
    return res


def fisher_yates(arr: list[T], rng: random.Random = random) -> list[T]:
    a = list(arr)
    for i in range(len(a) - 1, 0, -1):
        j = rng.randint(0, i)              # j in [0, i], NOT [0, n-1]
        a[i], a[j] = a[j], a[i]
    return a


class WeightedPicker:
    """Random Pick with Weight: prefix sums + bisect, O(n) build, O(log n) pick."""

    def __init__(self, w: list[int], rng: random.Random = random) -> None:
        self.prefix = list(itertools.accumulate(w))
        self.rng = rng

    def pick_index(self) -> int:
        target = self.rng.random() * self.prefix[-1]          # in [0, total)
        return bisect.bisect_right(self.prefix, target)        # first prefix > target
```

Correctness sketches:
- **Reservoir**: item i (0-indexed, i ≥ k) enters with probability `k/(i+1)`. A reservoir item survives step j with probability `1 - (k/(j+1)) · (1/k) = j/(j+1)`. The product from i+1 to n-1 telescopes to `(i+1)/n`, so every item ends with probability `k/n`.
- **Fisher-Yates**: position i gets each of the i+1 remaining elements with probability `1/(i+1)`: n! equally likely paths, one per permutation. Picking `j` from the whole array gives n^n paths, generally not divisible by n!, so it is biased.
- **Weighted pick**: index i owns `[prefix[i-1], prefix[i])`, of length `w[i]`. `bisect_right` returns the first prefix strictly greater than the target: the owner. Zero-weight items own empty intervals and are never picked.

Testing randomized code: seed a `random.Random` and assert **properties** (outputs are permutations, frequencies within tolerance), not exact outputs, as the companion file does.

#### Example 4.1: Random Pick with Weight [Medium]

**Problem:** [Random Pick with Weight (528)](https://leetcode.com/problems/random-pick-with-weight/) · **Template:** Template 4 (reservoir sampling, Fisher-Yates, weighted pick)

**Problem.** Given positive weights `w`, build an object whose `pickIndex()` returns index i with probability `w[i] / sum(w)`; it is called many times. `w = [1, 3]` → 0 a quarter of the time, 1 three quarters.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "probability proportional to weight", many picks over a fixed array → prefix sums once, binary search per pick.

**Approach:** lay the weights end to end on `[0, total)`; index i owns `[prefix[i-1], prefix[i])`, of length `w[i]`. Draw a uniform `target` in `[0, total)` and return its owner: the first prefix strictly greater than `target` (`bisect_right`). The probability of landing in a segment equals its length over `total`. Passing a seeded `rng` makes it testable: the companion file asserts frequencies within ±0.02 of `w[i] / total` over 20k–40k picks, and that zero-weight indices never appear.

```python
class WeightedPicker:
    """Random Pick with Weight: prefix sums + bisect, O(n) build, O(log n) pick."""

    def __init__(self, w: list[int], rng: random.Random = random) -> None:
        self.prefix = list(itertools.accumulate(w))
        self.rng = rng

    def pick_index(self) -> int:
        target = self.rng.random() * self.prefix[-1]          # in [0, total)
        return bisect.bisect_right(self.prefix, target)        # first prefix > target
```

**Complexity:** O(n) build, O(log n) per pick, O(n) space.

**Edge cases:** a single weight always returns 0; zero weights (not in 528, but common in variants) own empty segments, which `bisect_right` skips; with integer targets use `randint(1, total)` and `bisect_left` instead (don't mix the conventions). Follow-up: if items arrive as a stream of unknown length, you can't build prefix sums; use reservoir sampling ([Linked List Random Node (382)](https://leetcode.com/problems/linked-list-random-node/) is k = 1), or its weighted form (key `u ** (1 / w)`, keep the largest).

</details>

### Beyond the templates

#### Example B.1: Overflow and base conversion [Easy/Medium]

**Problems:** [Reverse Integer (7)](https://leetcode.com/problems/reverse-integer/), [String to Integer (atoi) (8)](https://leetcode.com/problems/string-to-integer-atoi/), [Excel Sheet Column Title (168)](https://leetcode.com/problems/excel-sheet-column-title/) · **Template:** none directly; 32-bit overflow checks related to Template 1 (bit helpers and 32-bit emulation)

**Why no template:** Digit-by-digit arithmetic with overflow guards and bijective base-26 are number mechanics, not bit tricks, bitmask enumeration, modular power, or sampling.

**Problems.** (a) Reverse the digits of a signed 32-bit integer; return 0 on overflow. (b) Parse a string into a signed 32-bit integer: skip leading spaces, read an optional sign and digits, clamp to the 32-bit range. (c) Convert a positive integer into its spreadsheet column title (1 → A, 26 → Z, 27 → AA).

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "signed 32-bit", "clamp", "return 0 on overflow".

**Approach:** Python won't overflow, so check as a fixed-width language would: **before** the multiply-add.

```python
INT_MAX, INT_MIN = 2**31 - 1, -(2**31)


def reverse_integer(x: int) -> int:
    """Reverse digits; return 0 if result leaves signed 32-bit range."""
    sign = -1 if x < 0 else 1
    x, rev = abs(x), 0
    while x:
        x, d = divmod(x, 10)
        # In a fixed-width language, check BEFORE multiplying:
        if rev > (INT_MAX - d) // 10:
            return 0
        rev = rev * 10 + d
    return sign * rev


def my_atoi(s: str) -> int:
    i, n = 0, len(s)
    while i < n and s[i] == " ":
        i += 1
    sign = 1
    if i < n and s[i] in "+-":
        sign = -1 if s[i] == "-" else 1
        i += 1
    val = 0
    while i < n and s[i].isdigit():
        val = val * 10 + int(s[i])
        i += 1
    return max(INT_MIN, min(INT_MAX, sign * val))
```

**Edge cases:** `reverse_integer` checks negatives against `INT_MAX`, slightly conservative (INT_MIN's magnitude is one larger), but no valid 32-bit input reverses to exactly 2^31 (that needs input 8463847412), so it is correct. `my_atoi` clamps at the end, fine in Python; in Java/C++ clamp inside the loop. Python gotcha: `-7 // 2 == -4` and `-7 % 2 == 1` (floor semantics); use `abs` and reapply the sign, or `int(a / b)` to truncate toward zero on small values.

Base conversion: repeated `divmod` by the base, least significant digit first. Excel columns are **bijective** base 26 (no zero digit: A=1 … Z=26), so subtract 1 before each `divmod`.

```python
DIGITS = "0123456789abcdefghijklmnopqrstuvwxyz"


def to_base(n: int, b: int) -> str:
    if n == 0:
        return "0"
    sign, n = ("-", -n) if n < 0 else ("", n)
    out = []
    while n:
        n, r = divmod(n, b)
        out.append(DIGITS[r])
    return sign + "".join(reversed(out))


def from_base(s: str, b: int) -> int:
    return int(s, b)   # manual: val = val * b + DIGITS.index(ch)


def excel_column_title(n: int) -> str:
    """Bijective base 26: there is no zero digit, so shift by one each step."""
    out = []
    while n:
        n, r = divmod(n - 1, 26)
        out.append(chr(ord("A") + r))
    return "".join(reversed(out))
```

**Complexity:** O(number of digits): O(log x) for reverse, O(len(s)) for atoi, O(log_b n) for base conversion.

</details>

#### Example B.2: Matrix manipulation [Medium]

**Problems:** [Rotate Image (48)](https://leetcode.com/problems/rotate-image/), [Spiral Matrix (54)](https://leetcode.com/problems/spiral-matrix/), [Set Matrix Zeroes (73)](https://leetcode.com/problems/set-matrix-zeroes/) · **Template:** none; in-place matrix mechanics

**Why no template:** In-place index manipulation (transpose plus reverse, boundary shrinking, marker rows) uses none of this module's bit, math, or sampling templates.

**Problems.** (a) Rotate an n×n matrix 90° clockwise in place. (b) Return all elements of an m×n matrix in spiral order. (c) If a cell is 0, zero its whole row and column, in place, with O(1) extra space.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** in-place matrix transforms: an index mapping (rotate), layer peeling with four boundaries (spiral), or the matrix itself as scratch space (zeroes).

**Approach:** rotate clockwise: transpose (swap across the main diagonal), then reverse each row. Counter-clockwise: transpose, then reverse each column (or reverse rows, then transpose). Loop only `j > i` in the transpose, or each pair swaps twice and undoes itself.

```python
def rotate(matrix: list[list[int]]) -> None:
    """Rotate n x n 90 degrees clockwise in place: transpose, then reverse rows."""
    n = len(matrix)
    for i in range(n):
        for j in range(i + 1, n):
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
    for row in matrix:
        row.reverse()
```

Spiral order: keep four boundaries and peel one layer per iteration (top row, right column, bottom row reversed, left column upward). The guards prevent double-reading a last remaining row or column.

```python
def spiral_order(matrix: list[list[int]]) -> list[int]:
    if not matrix or not matrix[0]:
        return []
    top, bottom, left, right = 0, len(matrix) - 1, 0, len(matrix[0]) - 1
    out = []
    while top <= bottom and left <= right:
        out.extend(matrix[top][c] for c in range(left, right + 1))
        out.extend(matrix[r][right] for r in range(top + 1, bottom + 1))
        if top < bottom:                    # guard single remaining row
            out.extend(matrix[bottom][c] for c in range(right - 1, left - 1, -1))
        if left < right:                    # guard single remaining column
            out.extend(matrix[r][left] for r in range(bottom - 1, top, -1))
        top, bottom, left, right = top + 1, bottom - 1, left + 1, right - 1
    return out
```

Set Matrix Zeroes in O(1) space: use row 0 and column 0 as markers. Save their original state in two booleans first, and zero them last so markers survive until read.

```python
def set_zeroes(matrix: list[list[int]]) -> None:
    """O(1) extra space: use row 0 / col 0 as marker arrays."""
    rows, cols = len(matrix), len(matrix[0])
    first_row_zero = any(matrix[0][c] == 0 for c in range(cols))
    first_col_zero = any(matrix[r][0] == 0 for r in range(rows))
    for r in range(1, rows):
        for c in range(1, cols):
            if matrix[r][c] == 0:
                matrix[r][0] = matrix[0][c] = 0
    for r in range(1, rows):
        for c in range(1, cols):
            if matrix[r][0] == 0 or matrix[0][c] == 0:
                matrix[r][c] = 0
    if first_row_zero:
        for c in range(cols):
            matrix[0][c] = 0
    if first_col_zero:
        for r in range(rows):
            matrix[r][0] = 0
```

**Complexity:** all O(R · C) time; O(1) extra space (spiral's output aside).

**Edge cases:** 1×1 and empty matrices; a single remaining row or column in spiral (the guards); zeros in row 0 or column 0 for Set Matrix Zeroes (the two booleans).

</details>

#### Example B.3: Geometry (cross product) [Hard]

**Problem:** [Max Points on a Line (149)](https://leetcode.com/problems/max-points-on-a-line/) · **Template:** none; cross product and slope normalization

**Why no template:** It's computational geometry (slopes normalized by gcd, cross products), which none of the templates cover.

**Problem.** Given points on a 2D plane with integer coordinates, return the maximum number of points that lie on a single straight line.

<details>
<summary>Decode, approach, and solution</summary>

**Decode:** "on a single straight line" with integer coordinates: test collinearity exactly (cross product or reduced slope), never with floats.

**Approach:** the cross product of `OA` and `OB` gives orientation: positive is counter-clockwise, negative clockwise, zero collinear. Exact on integers, it underlies convex hull, segment intersection, and polygon area (shoelace).

```python
def cross(o: tuple[int, int], a: tuple[int, int], b: tuple[int, int]) -> int:
    """>0: o->a->b turns counter-clockwise; <0: clockwise; 0: collinear."""
    return (a[0] - o[0]) * (b[1] - o[1]) - (a[1] - o[1]) * (b[0] - o[0])
```

Max points: for each anchor, every other point defines a direction, shared by points on the same line through the anchor. Represent it exactly as `(dy / g, dx / g)` with `g = gcd(dx, dy)` and a canonical sign (dx > 0, or dx == 0 and dy > 0), so `(1, -2)` and `(-1, 2)` hash the same. Float slopes fail on precision and vertical lines.

```python
def max_points(points: list[list[int]]) -> int:
    """Max points on one line. Slope as reduced (dy, dx) with a canonical sign."""
    if len(points) <= 2:
        return len(points)
    best = 0
    for i, (x1, y1) in enumerate(points):
        slopes: Counter[tuple[int, int]] = Counter()
        for x2, y2 in points[i + 1:]:
            dx, dy = x2 - x1, y2 - y1
            g = math.gcd(dx, dy)
            dx, dy = dx // g, dy // g
            if dx < 0 or (dx == 0 and dy < 0):   # canonical sign
                dx, dy = -dx, -dy
            slopes[(dy, dx)] += 1
        best = max(best, 1 + max(slopes.values(), default=0))
    return best
```

**Complexity:** O(n^2 log C) time (gcd per pair), O(n) space per anchor.

**Edge cases:** duplicate points (if allowed) give `g = 0`; count them separately and add them to every slope bucket for that anchor. Only pairs `j > i` are needed, since each line is counted from its first point.

</details>

## Template cheat sheet

### Template 1: bit helpers and 32-bit emulation
```python
MASK32 = 0xFFFFFFFF


def to_signed32(x: int) -> int:
    """Interpret the low 32 bits of x as a signed two's-complement int."""
    x &= MASK32
    return x - (1 << 32) if x & 0x80000000 else x


def get_bit(x: int, i: int) -> int:
    return (x >> i) & 1


def set_bit(x: int, i: int) -> int:
    return x | (1 << i)


def clear_bit(x: int, i: int) -> int:
    return x & ~(1 << i)


def toggle_bit(x: int, i: int) -> int:
    return x ^ (1 << i)


def lowest_set_bit(x: int) -> int:
    return x & -x
```
```python
def get_sum(a: int, b: int) -> int:
    """Add without + using xor (sum bits) and and<<1 (carries), 32-bit semantics."""
    a, b = a & MASK32, b & MASK32
    while b:
        a, b = (a ^ b) & MASK32, ((a & b) << 1) & MASK32
    return to_signed32(a)


def reverse_bits(n: int) -> int:
    out = 0
    for _ in range(32):
        out = (out << 1) | (n & 1)
        n >>= 1
    return out
```

### Template 2: subset enumeration with bitmasks
```python
def all_subsets(items: list[T]) -> list[list[T]]:
    n = len(items)
    return [[items[i] for i in range(n) if mask >> i & 1] for mask in range(1 << n)]


def iter_submasks(mask: int) -> Iterator[int]:
    """All submasks of mask, descending, including mask and 0. Total over all
    masks of n bits is O(3^n)."""
    sub = mask
    while True:
        yield sub
        if sub == 0:
            return
        sub = (sub - 1) & mask
```

### Template 3: fast modular exponentiation
```python
def mod_pow(base: int, exp: int, mod: int) -> int:
    """Binary exponentiation: square the base, multiply in when the bit is set."""
    result, base = 1 % mod, base % mod
    while exp:
        if exp & 1:
            result = result * base % mod
        base = base * base % mod
        exp >>= 1
    return result
```

### Template 4: reservoir sampling, Fisher-Yates, weighted pick
```python
def reservoir_sample(stream: Iterable[T], k: int, rng: random.Random = random) -> list[T]:
    """Uniform k-sample from a stream of unknown length in O(k) memory."""
    res: list[T] = []
    for i, x in enumerate(stream):
        if i < k:
            res.append(x)
        else:
            j = rng.randint(0, i)          # inclusive: i+1 choices
            if j < k:
                res[j] = x                 # keep x with probability k/(i+1)
    return res


def fisher_yates(arr: list[T], rng: random.Random = random) -> list[T]:
    a = list(arr)
    for i in range(len(a) - 1, 0, -1):
        j = rng.randint(0, i)              # j in [0, i], NOT [0, n-1]
        a[i], a[j] = a[j], a[i]
    return a


class WeightedPicker:
    """Random Pick with Weight: prefix sums + bisect, O(n) build, O(log n) pick."""

    def __init__(self, w: list[int], rng: random.Random = random) -> None:
        self.prefix = list(itertools.accumulate(w))
        self.rng = rng

    def pick_index(self) -> int:
        target = self.rng.random() * self.prefix[-1]          # in [0, total)
        return bisect.bisect_right(self.prefix, target)        # first prefix > target
```

## Variations and extensions

- **Bits + hashing**: DNA as 2 bits per base gives exact integer keys (Repeated DNA Sequences); a letter set as a 26-bit mask (Maximum Product of Word Lengths: no shared letters exactly when `mask1 & mask2 == 0`).
- **Bits + prefix sums**: prefix XOR answers subarray XOR queries; vowel-parity masks (Find the Longest Substring Containing Vowels in Even Counts) pair with "first index seen" maps (`04-arrays-hashing-prefix-sums.md`).
- **Bits + trie**: maximum XOR pair (`19-tries-and-string-algorithms.md`).
- **Bitmask DP**: TSP / shortest path visiting all nodes (BFS over `(node, mask)`), assignment, partition (`17-dynamic-programming-advanced.md`).
- **Bitwise AND of a range**: the common binary prefix of the endpoints; shift both right until equal.
- **Randomization**: Random Pick Index (reservoir, k = 1, over matching indices), Linked List Random Node, Insert Delete GetRandom O(1) (array + index map, swap-with-last on delete), random-pivot quickselect (`11-heaps-top-k-and-merging.md`).
- **Weighted pick with updates**: Fenwick tree over weights, O(log n) updates and picks.
- **Matrix**: rotate by layers with 4-way swaps, Game of Life in place (old and new state in two bits per cell), Spiral Matrix II (fill instead of read).
- **Geometry**: convex hull via monotone chain (Erect the Fence), Valid Square (four equal sides via squared distances, two equal diagonals), rectangle overlap (intersect intervals per axis, `12-intervals-and-sweep-line.md`).

Follow-ups:
- "Your language has 32-bit ints, what changes?" Overflow checks before multiplication, masking, and `>>>` vs `>>`.
- "The stream is distributed across machines; sample k uniformly." Each machine keeps a reservoir and its count; merge by sampling from each in proportion to its count, or give each item a random key and keep the global k smallest.

## Common pitfalls

- Forgetting Python has no fixed width: negatives have infinitely many leading 1s, so bit loops on them never end without `& MASK32`.
- Operator precedence: in Python, comparisons bind looser than `&`, `|`, `^`, so `n & (n - 1) == 0` is `(n & (n - 1)) == 0` and works. In C, Java, and JS `==` binds tighter, so the same text is a bug. Parenthesize anyway. In Python `+` and `<<` bind tighter than `&`: `x & 1 << i` is `x & (1 << i)`.
- Prefer `1 << n` to `2 ** n`: same value, clearer intent.
- Fisher-Yates with `randint(0, n - 1)` instead of `randint(0, i)` is biased.
- Reservoir with `randint(0, i - 1)` (off by one) is not uniform.
- Weighted pick: mixing `bisect_left` on an integer target in `[1, total]` with `bisect_right` on a float in `[0, total)`. Pick one convention and verify a zero-weight item is never chosen.
- Set Matrix Zeroes: zeroing row 0 / column 0 before using them as markers.
- Spiral: missing the single-row / single-column guards duplicates elements.
- Max Points: float slopes, unnormalized signs, forgetting vertical lines.
- Modular arithmetic: `% MOD` only at the end (slow with huge ints in Python, overflow elsewhere); dividing under a modulus without an inverse.
- Fermat inverse needs a prime modulus and `a % p != 0`.

## When NOT to use / look-alikes

| Situation | Use instead | Why |
|---|---|---|
| "Single number" with O(n) memory allowed and arbitrary repeat counts | `Counter` | XOR only cancels even multiplicities |
| Find the duplicate in [1..n] with n+1 values, read-only array | Floyd cycle detection or binary search on value (`07`, `09`) | XOR fails when the duplicate repeats more than twice |
| n > 25 and "subsets" | DP over sums / meet in the middle / greedy | 2^n too big |
| Count paths in a grid **with obstacles** | DP (`16`) | Closed-form nCr no longer applies |
| nCr where p is small or not prime | Pascal DP or Lucas' theorem | Factorials become 0 mod p |
| Pick uniformly from an array you already have | `random.choice` | Reservoir is for unknown-length streams |
| Rotating a non-square matrix | Build a new matrix `zip(*matrix[::-1])` | In-place transpose needs a square |

## Practice ladder

Track your progress in [`practice-tracker.md`](practice-tracker.md).

### Must-do

1. [Single Number (136)](https://leetcode.com/problems/single-number/) · Easy · XOR: XOR all.
2. [Counting Bits (338)](https://leetcode.com/problems/counting-bits/) · Easy · Bit DP: `bits[i >> 1] + (i & 1)`.
3. [Missing Number (268)](https://leetcode.com/problems/missing-number/) · Easy · XOR / math: XOR indices with values.
4. [Rotate Image (48)](https://leetcode.com/problems/rotate-image/) · Medium · Matrix in place: transpose then reverse rows.
5. [Spiral Matrix (54)](https://leetcode.com/problems/spiral-matrix/) · Medium · Matrix traversal: four boundaries with guards.

### Optional

**Bits and XOR**
- [Number of 1 Bits (191)](https://leetcode.com/problems/number-of-1-bits/) · Easy · Kernighan.
- [Power of Two (231)](https://leetcode.com/problems/power-of-two/) · Easy · `n & (n - 1) == 0`.
- [Reverse Bits (190)](https://leetcode.com/problems/reverse-bits/) · Easy · Shift out, shift in, 32 times.
- [Single Number II (137)](https://leetcode.com/problems/single-number-ii/) · Medium · Bit counts mod 3.
- [Single Number III (260)](https://leetcode.com/problems/single-number-iii/) · Medium · Split by lowest differing bit.
- [Sum of Two Integers (371)](https://leetcode.com/problems/sum-of-two-integers/) · Medium · XOR + carry with 32-bit masking.
- [Bitwise AND of Numbers Range (201)](https://leetcode.com/problems/bitwise-and-of-numbers-range/) · Medium · Common prefix.

**Bitmasks**
- [Subsets (78)](https://leetcode.com/problems/subsets/) · Medium · `range(1 << n)`.
- [Gray Code (89)](https://leetcode.com/problems/gray-code/) · Medium · `i ^ (i >> 1)`.
- [Partition to K Equal Sum Subsets (698)](https://leetcode.com/problems/partition-to-k-equal-sum-subsets/) · Medium · DP over masks.

**Math**
- [Pascal's Triangle (118)](https://leetcode.com/problems/pascals-triangle/) · Easy · Row from previous row.
- [Excel Sheet Column Title (168)](https://leetcode.com/problems/excel-sheet-column-title/) · Easy · Bijective base 26.
- [Pow(x, n) (50)](https://leetcode.com/problems/powx-n/) · Medium · Binary exponentiation, negative n.
- [Count Primes (204)](https://leetcode.com/problems/count-primes/) · Medium · Sieve from i*i.
- [Unique Paths (62)](https://leetcode.com/problems/unique-paths/) · Medium · `C(m+n-2, m-1)`.
- [Reverse Integer (7)](https://leetcode.com/problems/reverse-integer/) · Medium · Check before multiply.
- [String to Integer (atoi) (8)](https://leetcode.com/problems/string-to-integer-atoi/) · Medium · Whitespace, sign, digits, clamp.

**Randomization**
- [Linked List Random Node (382)](https://leetcode.com/problems/linked-list-random-node/) · Medium · Reservoir, k = 1.
- [Shuffle an Array (384)](https://leetcode.com/problems/shuffle-an-array/) · Medium · Fisher-Yates.
- [Random Pick with Weight (528)](https://leetcode.com/problems/random-pick-with-weight/) · Medium · Prefix sums + bisect.

**Matrix and geometry**
- [Set Matrix Zeroes (73)](https://leetcode.com/problems/set-matrix-zeroes/) · Medium · First row/column as markers.
- [Max Points on a Line (149)](https://leetcode.com/problems/max-points-on-a-line/) · Hard · Normalized `(dy, dx)` per anchor.

## Self-check

1. Why does `x & (x - 1)` clear the lowest set bit, and what does `x & -x` return?
2. How do you emulate 32-bit signed addition in Python, and why does the naive carry loop hang on negative inputs?
3. In Single Number III, why is any set bit of `a ^ b` enough to separate a and b?
4. Explain why `sub = (sub - 1) & mask` enumerates all submasks and why iterating submasks of all masks costs O(3^n).
5. Derive fast exponentiation from the binary representation of the exponent. How do you compute nCr mod 10^9+7 in O(1) per query?
6. Prove reservoir sampling keeps each item with probability k/n. Why is "swap with any index" shuffling biased?
7. Why must Max Points on a Line use reduced integer pairs with a canonical sign instead of float slopes?
