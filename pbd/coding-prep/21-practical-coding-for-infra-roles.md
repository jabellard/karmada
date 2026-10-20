# 21. Practical Coding for Infra Roles

Infra and AI-infra companies often swap a LeetCode round for a **practical round**: build a realistic component (cache, rate limiter, scheduler, controller loop) in 2–4 cumulative parts, then field follow-ups on concurrency, failures, and scale. The algorithms are rarely hard. The grade is code a teammate would merge: clean API, correct edge cases, tests, production judgment. For a Staff platform engineer this is home turf, if you deliver incrementally and don't over-engineer part 1.

All code lives and is tested in `code/21_practical_coding.py`. Time-dependent components take an injectable `clock` or `sleep`, so tests are deterministic.

## What practical rounds evaluate

| Dimension | What they watch for | What "strong" looks like |
|---|---|---|
| Correctness | Spec and edge cases work | Ran tests; explicit boundaries (TTL exactly at expiry, empty input) |
| API design | Names, signatures, return types, error contract | Small surface; raise vs return chosen on purpose; injectable dependencies |
| Incremental delivery | Part 1 works before part 2 starts | Simplest correct version, then extend without rewriting |
| Testing | Do you verify your own code? | Early asserts; table-driven cases; deterministic time |
| Ambiguity | Handling under-specified prompts | 2–4 sharp questions, assumptions stated aloud, moves on |
| Production thinking | Follow-ups: threads, crashes, memory, scale | Concrete: "one lock around get+put because get mutates order" |
| Code quality | Readability under time pressure | Clearly named helpers, no dead code, comments on invariants only |

**How they differ from LeetCode rounds**

| LeetCode round | Practical round |
|---|---|
| One function, fixed signature | You design the class/API |
| Optimal complexity is the bar | Correct, readable, extensible is the bar; complexity is usually obvious |
| Clean input | Messy input (malformed log lines, unknown commands) |
| Single part | Multi-part; part N reuses part N-1 |
| Code rarely run | Code and tests often run in a real editor |
| Algorithmic follow-ups | Systems follow-ups: locking, persistence, distribution, backpressure |

## Approach playbook

**1. Clarify requirements and the API (3–5 min).** Restate the task in one sentence; ask only design-changing questions:
- Interface: method names and signatures? Return value or exception on missing/invalid?
- Boundaries: is a 10s TTL expired at exactly t=10? Does `get` refresh recency?
- Scale: single-threaded first? In-memory only? Roughly how many keys?
- Time: may I inject a clock for testing? (Always ask; it signals testing maturity.)

Write assumptions as a comment atop the file. If told "your call", pick, say why, move on.

**2. Write the interface and 2–3 tests before the body.** Stub signatures and docstrings, then assert the happy path and one edge case. This forces API decisions early and gives a fast check after each change.

**3. Implement the simplest correct version.** `OrderedDict` before a hand-rolled linked list; a coarse lock before striping; sorting before a streaming sketch. Say: "This is O(n log n); if it matters I can switch to X."

**4. Run the tests, then extend.** Each part should be additive. If part 2 forces a rewrite, name the refactor and why ("I'll pull expiry into a helper so get and cleanup share it").

**5. Narrate each tradeoff in one sentence:** "Lazy expiry keeps reads O(1) but dead keys accumulate, so I'll add a bounded cleanup."

**6. Handle errors on purpose.** Pick one contract (raise `KeyError`, return `None`, return a sentinel) and keep it. Validate constructor args (`capacity <= 0`). Parsers count and skip malformed input, never crash.

**7. Budget time.** 60 minutes, 3 parts: 5 min clarify, 15 per part, 5–10 min follow-ups. If part 1 runs long, stop polishing: there are usually more parts than time, and plain working code that gets further beats a perfect part 1.

**8. Extending an unfamiliar codebase.** Some rounds hand you a repo.
- Read the README; run the tests before changing anything to learn the baseline (some may already fail).
- Find entry points (`main`, CLI parsing, HTTP handlers, the class named in the prompt); trace one request end to end.
- Copy the naming, error handling, and test style of the most similar existing feature.
- Make the smallest working change, add a test beside the existing ones, run the suite.
- Narrate what you read: "This looks like the store interface; `MemoryStore` and `RedisStore` both implement it, so I'll add TTL to the interface first."
- Don't refactor untouched code unless it blocks you; mention it as a follow-up.

## Shared test helper

Time-dependent components take `clock: Callable[[], float] = time.monotonic`; tests pass a `FakeClock` and advance it by hand. Use `time.monotonic` in production, not `time.time`: wall-clock time can jump backward on NTP adjustments.

```python
class FakeClock:
    """Callable clock. Pass `clock=fake` wherever a component takes `clock`."""

    def __init__(self, start: float = 0.0):
        self.now = start

    def __call__(self) -> float:
        return self.now

    def advance(self, seconds: float) -> None:
        self.now += seconds
```

## Exercise 1: LRU cache

**Spec (as given).** "Implement a cache with a fixed capacity. `get(key)` returns the value or a default. `put(key, value)` inserts or updates. When full, evict the least recently used entry. Both operations O(1)."

**Clarifying questions.**
- Does `get` count as a use? (Normally yes.) Does updating a key refresh recency? (Yes.)
- On a miss, does `get` return `None`, `-1`, or raise?
- Capacity 0 allowed? Thread safety now or later?

**Design notes.** O(1) lookup needs a hash map; O(1) "move to most recent" and "evict oldest" need O(1) middle removal from an ordered structure: a doubly linked list. `OrderedDict` is exactly that (in C), so start there and offer the manual version if asked. Sentinel head/tail nodes remove all `None` checks in link/unlink.

**Version A: OrderedDict.**

```python
class LRUCacheOD:
    """LRU via OrderedDict: end = most recently used, front = eviction victim."""

    def __init__(self, capacity: int):
        if capacity <= 0:
            raise ValueError("capacity must be positive")
        self.capacity = capacity
        self.data: OrderedDict[Any, Any] = OrderedDict()

    def get(self, key: Any, default: Any = None) -> Any:
        if key not in self.data:
            return default
        self.data.move_to_end(key)
        return self.data[key]

    def put(self, key: Any, value: Any) -> None:
        if key in self.data:
            self.data.move_to_end(key)
        self.data[key] = value
        if len(self.data) > self.capacity:
            self.data.popitem(last=False)

    def __len__(self) -> int:
        return len(self.data)
```

**Version B: hash map + doubly linked list.**

```python
class _Node:
    __slots__ = ("key", "val", "prev", "next")

    def __init__(self, key: Any = None, val: Any = None):
        self.key, self.val = key, val
        self.prev: _Node | None = None
        self.next: _Node | None = None


class LRUCache:
    """LRU via hash map + doubly linked list with sentinels.

    Invariant: head.next is MRU, tail.prev is LRU; map[key] is the node in the list.
    """

    def __init__(self, capacity: int):
        if capacity <= 0:
            raise ValueError("capacity must be positive")
        self.capacity = capacity
        self.map: dict[Any, _Node] = {}
        self.head, self.tail = _Node(), _Node()  # sentinels: no None checks
        self.head.next, self.tail.prev = self.tail, self.head

    def _unlink(self, node: _Node) -> None:
        node.prev.next, node.next.prev = node.next, node.prev

    def _push_front(self, node: _Node) -> None:
        node.prev, node.next = self.head, self.head.next
        self.head.next.prev = node
        self.head.next = node

    def get(self, key: Any, default: Any = None) -> Any:
        node = self.map.get(key)
        if node is None:
            return default
        self._unlink(node)
        self._push_front(node)
        return node.val

    def put(self, key: Any, value: Any) -> None:
        node = self.map.get(key)
        if node is not None:
            node.val = value
            self._unlink(node)
            self._push_front(node)
            return
        node = _Node(key, value)
        self.map[key] = node
        self._push_front(node)
        if len(self.map) > self.capacity:
            lru = self.tail.prev
            self._unlink(lru)
            del self.map[lru.key]

    def keys_mru_to_lru(self) -> list:
        out, cur = [], self.head.next
        while cur is not self.tail:
            out.append(cur.key)
            cur = cur.next
        return out

    def __len__(self) -> int:
        return len(self.map)
```

**Tests (excerpt).**

```python
c = LRUCache(2)
c.put("a", 1); c.put("b", 2)
assert c.get("a") == 1          # a is now MRU
c.put("c", 3)                   # evicts b
assert c.get("b") is None
c.put("a", 10)                  # update refreshes recency
c.put("d", 4)                   # evicts c
assert c.get("c") is None and c.get("a") == 10
```

**Complexity.** O(1) per operation; O(capacity) space.

**Follow-ups.**
- *Thread safety?* → One lock around `get` and `put`. `get` mutates recency order, so every op is a write; a read-write lock doesn't help.

```python
class ThreadSafeLRU(LRUCacheOD):
    """Follow-up answer: one coarse lock. get() mutates order, so it needs the lock too."""

    def __init__(self, capacity: int):
        super().__init__(capacity)
        self._lock = threading.Lock()

    def get(self, key: Any, default: Any = None) -> Any:
        with self._lock:
            return super().get(key, default)

    def put(self, key: Any, value: Any) -> None:
        with self._lock:
            super().put(key, value)
```

- *Lock contention at high QPS?* → Shard into N LRUs by `hash(key) % N` (globally approximate LRU), or use approximate LRU (Redis-style sampled eviction, CLOCK), which skips reordering on reads.
- *Evict by bytes, not count?* → Track `total_bytes`; evict from the tail while over budget. Reject items larger than the budget.
- *TTL too?* → Combine with Exercise 2: store `expires_at` in the node, check on `get`.
- *LFU instead?* → Map key → (value, freq), freq → OrderedDict of keys, plus `min_freq`. All O(1).
- *Distributed?* → Consistent hashing (Exercise 8) routes keys to nodes, each running a local LRU.

## Exercise 2: Key-value store with TTL

**Spec.** "Build an in-memory KV store. `set(key, value, ttl=None)`, `get(key)`. Keys with a TTL disappear after `ttl` seconds. Part 2: memory keeps growing with expired keys nobody reads; fix that."

**Clarifying questions.**
- Expired *at* `expires_at` or only after? (Pick `now >= expires_at` and test it.)
- Does overwriting reset the TTL? Does `set` without TTL clear one? (Redis: yes.)
- Background cleanup thread, or call-driven cleanup?
- May I inject the clock?

**Design notes.**
- **Lazy expiry**: check `expires_at` on read, delete if past. O(1), but never-read dead keys leak memory.
- **Active expiry**: min-heap of `(expires_at, version, key)`; `cleanup()` pops while the top is due. Overwrites leave stale heap entries; a per-write `version` lets cleanup skip them without searching the heap.
- `max_items` bounds work per cleanup so an expiration burst can't stall the caller (Redis uses random sampling for the same reason).
- A `_MISSING` sentinel distinguishes "absent" from a stored `None`.

```python
_MISSING = object()  # sentinel: distinguishes 'absent' from a stored None


class TTLStore:
    """KV store with per-key TTL.

    Lazy expiry on read; cleanup() does active expiry using a min-heap of
    (expires_at, version, key). Stale heap entries (key overwritten/deleted)
    are skipped by comparing the version.
    """

    def __init__(self, clock: Callable[[], float] = time.monotonic):
        self.clock = clock
        self.data: dict[str, tuple[Any, float | None, int]] = {}  # key -> (value, expires_at, version)
        self.heap: list[tuple[float, int, str]] = []
        self._version = 0

    def set(self, key: str, value: Any, ttl: float | None = None) -> None:
        if ttl is not None and ttl <= 0:
            raise ValueError("ttl must be positive")
        self._version += 1
        expires = None if ttl is None else self.clock() + ttl
        self.data[key] = (value, expires, self._version)
        if expires is not None:
            heapq.heappush(self.heap, (expires, self._version, key))

    def _expired(self, expires: float | None) -> bool:
        return expires is not None and self.clock() >= expires

    def get(self, key: str, default: Any = None) -> Any:
        entry = self.data.get(key)
        if entry is None:
            return default
        value, expires, _ = entry
        if self._expired(expires):
            del self.data[key]  # lazy expiry
            return default
        return value

    def delete(self, key: str) -> bool:
        return self.data.pop(key, None) is not None

    def ttl(self, key: str) -> float | None:
        """Remaining seconds, None if no TTL, raises KeyError if absent/expired."""
        if self.get(key, _MISSING) is _MISSING:
            raise KeyError(key)
        expires = self.data[key][1]
        return None if expires is None else expires - self.clock()

    def cleanup(self, max_items: int | None = None) -> int:
        """Remove expired keys. Bounded work per call if max_items is set."""
        removed, now = 0, self.clock()
        while self.heap and self.heap[0][0] <= now:
            if max_items is not None and removed >= max_items:
                break
            _, version, key = heapq.heappop(self.heap)
            entry = self.data.get(key)
            if entry is not None and entry[2] == version:  # not overwritten since
                del self.data[key]
                removed += 1
        return removed

    def __len__(self) -> int:  # includes expired-but-not-yet-cleaned keys
        return len(self.data)
```

**Tests (excerpt).** The exact-boundary and stale-heap-entry cases are common bugs.

```python
clk = FakeClock()
s = TTLStore(clock=clk)
s.set("a", 1, ttl=10)
clk.advance(9.999); assert s.get("a") == 1
clk.advance(0.001); assert s.get("a") is None   # expired exactly at t=10

s.set("c", 3, ttl=5)        # heap entry expires t=15
clk.advance(3)
s.set("c", 33, ttl=5)       # new entry expires t=18; old heap entry is now stale
clk.advance(3)              # t=16
assert s.cleanup() == 0 and s.get("c") == 33
```

**Complexity.** `set` O(log n) (heap push), `get` O(1), `cleanup` O(k log n) for k popped entries. Space O(n + stale heap entries).

**Follow-ups.**
- *Thread safety?* → One lock around every method. Background `cleanup` uses bounded `max_items` so it releases the lock often.
- *Background cleaner thread?* → `threading.Thread(daemon=True)` looping `cleanup(); stop_event.wait(interval)`. `Event.wait` (not `sleep`) makes shutdown immediate.
- *Stale heap entries pile up under heavy overwrites?* → Rebuild the heap from `data` when `len(heap) > 2 * len(data)`, or use a timing wheel (per-second buckets, O(1) insert).
- *Persistence?* → Append-only `set/delete` log with absolute `expires_at` (relative TTLs would extend every key on restart). Replay on start, skipping expired; periodic snapshot + log truncation.
- *Clock skew across replicas?* → Store absolute expiry from the leader's clock and make expiry a leader decision (Redis replicas wait for the primary's DEL).

## Exercise 3: Key-value store with transactions

**Spec.** "Implement a command interpreter: `SET k v`, `GET k` (prints value or NULL), `DELETE k`, `COUNT v` (how many keys have value v). Then add `BEGIN`, `ROLLBACK`, `COMMIT`. Transactions nest. `ROLLBACK` undoes the innermost open transaction. `COMMIT` commits all open transactions. Print `NO TRANSACTION` when there's nothing to roll back or commit."

**Clarifying questions.**
- Does `COMMIT` close all levels or just the innermost? (Classic: all.)
- Does `GET` in a transaction see uncommitted writes? (Yes, single client.)
- Must `COUNT` be fast? (Usually O(1) or O(log n).)
- Single client, or concurrent clients with isolation? (Single first.)

**Design notes.** Two standard designs:

| Design | Write | Read | Rollback | Commit |
|---|---|---|---|---|
| **Undo log** (apply writes now, record old values per txn) | O(1) | O(1) | O(writes in txn) | O(1) |
| Overlay stack (each txn holds a diff; reads search top-down) | O(1) | O(depth) | O(1) | O(total writes) to merge |

Undo log wins for the classic spec: O(1) reads and exact `COUNT`. Record a key's old value only on its **first** write in a txn, so rollback restores the value as of `BEGIN`. Rollback goes through `_raw_set` to keep `value_counts` consistent.

```python
class TxnStore:
    """SET/GET/DELETE/COUNT with BEGIN/ROLLBACK/COMMIT, nested.

    Each open transaction keeps an undo log: key -> previous value (or _MISSING).
    Only the first write to a key inside a txn records the old value.
    ROLLBACK replays the top undo log; COMMIT discards all logs (changes are
    already applied). value_counts supports COUNT in O(1).
    """

    def __init__(self):
        self.data: dict[str, Any] = {}
        self.value_counts: Counter = Counter()
        self.undo_stack: list[dict[str, Any]] = []

    def _raw_set(self, key: str, value: Any) -> None:
        old = self.data.get(key, _MISSING)
        if old is not _MISSING:
            self.value_counts[old] -= 1
        if value is _MISSING:
            self.data.pop(key, None)
        else:
            self.data[key] = value
            self.value_counts[value] += 1

    def _record(self, key: str) -> None:
        if self.undo_stack and key not in self.undo_stack[-1]:
            self.undo_stack[-1][key] = self.data.get(key, _MISSING)

    def set(self, key: str, value: Any) -> None:
        self._record(key)
        self._raw_set(key, value)

    def get(self, key: str) -> Any:
        return self.data.get(key)

    def delete(self, key: str) -> None:
        self._record(key)
        self._raw_set(key, _MISSING)

    def count(self, value: Any) -> int:
        return self.value_counts[value]

    def begin(self) -> None:
        self.undo_stack.append({})

    def rollback(self) -> bool:
        if not self.undo_stack:
            return False  # "NO TRANSACTION"
        for key, old in self.undo_stack.pop().items():
            self._raw_set(key, old)
        return True

    def commit(self) -> bool:
        if not self.undo_stack:
            return False
        self.undo_stack.clear()  # commits all open transactions
        return True
```

Driver in the usual interview shape:

```python
def run_txn_commands(lines: list[str]) -> list[str]:
    """Command-line driver, the shape many interviewers give (one command per line)."""
    store, out = TxnStore(), []
    for line in lines:
        cmd, *args = line.split()
        cmd = cmd.upper()
        if cmd == "SET":
            store.set(args[0], args[1])
        elif cmd == "GET":
            v = store.get(args[0])
            out.append("NULL" if v is None else v)
        elif cmd == "DELETE":
            store.delete(args[0])
        elif cmd == "COUNT":
            out.append(str(store.count(args[0])))
        elif cmd == "BEGIN":
            store.begin()
        elif cmd == "ROLLBACK":
            if not store.rollback():
                out.append("NO TRANSACTION")
        elif cmd == "COMMIT":
            if not store.commit():
                out.append("NO TRANSACTION")
        else:
            out.append(f"UNKNOWN COMMAND {cmd}")
    return out
```

**Tests.**

```python
def test_txn_store() -> None:
    s = TxnStore()
    s.set("a", "10")
    s.begin()
    s.set("a", "20")
    s.begin()
    s.set("a", "30")
    s.delete("a")
    assert s.get("a") is None
    assert s.rollback() and s.get("a") == "20"
    assert s.rollback() and s.get("a") == "10"
    assert not s.rollback()
    # Nested commit commits everything.
    s.begin()
    s.set("b", "1")
    s.begin()
    s.set("c", "1")
    assert s.count("1") == 2
    assert s.commit() and not s.rollback()
    assert s.get("b") == "1" and s.get("c") == "1"
    # Rollback restores a key created inside the txn as absent.
    s.begin()
    s.set("new", "x")
    s.rollback()
    assert s.get("new") is None and s.count("x") == 0

    out = run_txn_commands([
        "SET x 5", "GET x", "BEGIN", "SET x 6", "SET y 6", "COUNT 6",
        "BEGIN", "DELETE x", "GET x", "ROLLBACK", "GET x", "ROLLBACK",
        "GET x", "GET y", "COUNT 6", "COMMIT", "FOO",
    ])
    assert out == ["5", "2", "NULL", "6", "5", "NULL", "0", "NO TRANSACTION", "UNKNOWN COMMAND FOO"], out
```

**Complexity.** SET/GET/DELETE/COUNT O(1). BEGIN/COMMIT O(1). ROLLBACK O(distinct keys written in that txn). Space O(n + total undo entries).

**Follow-ups.**
- *COMMIT only the innermost level?* → Merge the top undo log into its parent, keeping the parent's (older) entry on conflict: `for k, v in top.items(): parent.setdefault(k, v)`.
- *Concurrent clients with isolation?* → Per-client write set, validated and applied under a global lock at commit (optimistic concurrency control); or MVCC, where readers see a snapshot of versioned values at their start version.
- *Durability?* → Write-ahead log: append the commit's write set as one record, fsync, then apply. On restart, replay only committed records.
- *A txn writes millions of keys?* → The undo log grows with keys touched. Spill to disk, or cap txn size and reject.

## Exercise 4: Rate limiters

**Spec.** "Implement `allow(key) -> bool`: each client key may make at most N requests per window. Part 2: allow short bursts but enforce a sustained rate. Part 3: make it thread safe."

**Clarifying questions.**
- Hard limit over any sliding window, or approximate OK? (Log vs counter.)
- Does a client idle for an hour get a burst? How big?
- Per-key memory budget? Millions of keys?
- Single process or fleet-wide?

**Design notes.**

| Algorithm | Memory per key | Accuracy | Bursts |
|---|---|---|---|
| Fixed window counter | O(1) | Up to 2N at a window boundary | Yes, at boundary |
| Sliding window log | O(N) timestamps | Exact | No |
| Sliding window counter | O(1) | Approximate (assumes uniform arrival last window) | Smoothed |
| Token bucket | O(1) | Exact for its model | Up to `capacity`, then `rate`/s |

Token bucket refills **lazily** from elapsed time on each call, with no timer thread. All three hold one lock per call around a tiny critical section.

```python
class TokenBucketLimiter:
    """Per-key token bucket. Refill lazily on each call: no background thread.

    Allows bursts up to `capacity`, sustained rate `rate` tokens/sec.
    """

    def __init__(self, capacity: float, rate: float, clock: Callable[[], float] = time.monotonic):
        self.capacity, self.rate, self.clock = capacity, rate, clock
        self.buckets: dict[str, tuple[float, float]] = {}  # key -> (tokens, last_refill)
        self._lock = threading.Lock()

    def allow(self, key: str, cost: float = 1.0) -> bool:
        with self._lock:
            now = self.clock()
            tokens, last = self.buckets.get(key, (self.capacity, now))
            tokens = min(self.capacity, tokens + (now - last) * self.rate)
            allowed = tokens >= cost
            if allowed:
                tokens -= cost
            self.buckets[key] = (tokens, now)
            return allowed
```

```python
class SlidingWindowLogLimiter:
    """Exact: keep timestamps of accepted requests within the last `window` seconds.

    Memory O(limit) per key.
    """

    def __init__(self, limit: int, window: float, clock: Callable[[], float] = time.monotonic):
        self.limit, self.window, self.clock = limit, window, clock
        self.logs: dict[str, deque[float]] = defaultdict(deque)
        self._lock = threading.Lock()

    def allow(self, key: str) -> bool:
        with self._lock:
            now = self.clock()
            log = self.logs[key]
            while log and log[0] <= now - self.window:  # evict outside (now-window, now]
                log.popleft()
            if len(log) < self.limit:
                log.append(now)
                return True
            return False
```

```python
class SlidingWindowCounterLimiter:
    """Approximate: two fixed-window counters, weight the previous by overlap.

    estimate = prev_count * (1 - elapsed_in_current/window) + curr_count. O(1) memory per key.
    """

    def __init__(self, limit: int, window: float, clock: Callable[[], float] = time.monotonic):
        self.limit, self.window, self.clock = limit, window, clock
        self.state: dict[str, list] = {}  # key -> [window_index, curr_count, prev_count]
        self._lock = threading.Lock()

    def allow(self, key: str) -> bool:
        with self._lock:
            now = self.clock()
            idx = int(now // self.window)
            st = self.state.setdefault(key, [idx, 0, 0])
            if idx != st[0]:
                # Rolled forward: if exactly one window passed, current becomes previous.
                st[2] = st[1] if idx == st[0] + 1 else 0
                st[1], st[0] = 0, idx
            elapsed_frac = (now - idx * self.window) / self.window
            estimate = st[2] * (1 - elapsed_frac) + st[1]
            if estimate < self.limit:
                st[1] += 1
                return True
            return False
```

**Tests (excerpt).**

```python
clk = FakeClock()
tb = TokenBucketLimiter(capacity=3, rate=1.0, clock=clk)
assert [tb.allow("u") for _ in range(4)] == [True, True, True, False]  # burst of 3
assert tb.allow("other")             # keys are independent
clk.advance(1)
assert tb.allow("u") and not tb.allow("u")
clk.advance(100)                     # refill caps at capacity
assert [tb.allow("u") for _ in range(4)] == [True, True, True, False]

clk = FakeClock()
sc = SlidingWindowCounterLimiter(limit=10, window=10, clock=clk)
assert sum(sc.allow("u") for _ in range(15)) == 10
clk.advance(15)                      # halfway into next window: prev weighted 0.5 -> 5 used
assert sum(sc.allow("u") for _ in range(15)) == 5
```

**Complexity.** Token bucket and counter: O(1) time and space per key. Log: amortized O(1) time (each timestamp popped once), O(limit) space per key.

**Follow-ups.**
- *Millions of idle keys?* → Keep state in an LRU or TTL store (Exercises 1–2). A bucket idle for `capacity/rate` seconds is full, so dropping it loses nothing.
- *Lock contention?* → Stripe locks by `hash(key) % N`; per-key state is independent.
- *Fleet-wide?* → Central store with atomic ops: Redis `INCR` + `EXPIRE` for fixed window, a Lua script for token bucket (the read-modify-write must be atomic). Alternatives: local limiters at `limit / num_instances` (no coordination, inaccurate under skew), or periodically synced local counts.
- *What does the client get?* → HTTP 429 with `Retry-After`, which token bucket computes as `(cost - tokens) / rate`.
- *Clock goes backward?* → Impossible in-process with `time.monotonic`. Across machines, clamp negative elapsed to 0.

## Exercise 5: Worker pool with retries and backoff

**Spec.** "Build a pool of N worker threads that executes submitted tasks. Tasks can fail transiently; retry them with exponential backoff, up to M retries. Collect results. Part 2: same thing with asyncio."

**Clarifying questions.**
- Which errors are retryable? (Transient: timeouts, 503. Permanent, never retried: validation errors.)
- Max retries, base delay, cap? Jitter? (Say yes, and why.)
- Should `submit` block on a full queue (backpressure)?
- Shutdown: drain the queue or abandon pending work?

**Design notes.**
- Workers loop on `get()` from the thread-safe `queue.Queue`. `task_done()` sits in a `finally` so `join()` can't hang on an exception.
- `wait()` uses `Queue.join()`. Retries run inside the worker, so a task stays "not done" until its final attempt.
- Shutdown: one **poison pill** per worker, then `join()` each thread.
- **Full jitter** backoff: `uniform(0, min(cap, base * 2**attempt))`. Without jitter, clients that failed together retry together and recreate the spike (thundering herd).
- `sleep` and `rng` are injectable; the test's fake sleep records delays, so it runs instantly and deterministically.

```python
def backoff_delay(attempt: int, base: float = 0.1, cap: float = 10.0,
                  rng: random.Random | None = None) -> float:
    """'Full jitter': uniform(0, min(cap, base * 2**attempt)). attempt starts at 0."""
    rng = rng or random
    return rng.uniform(0, min(cap, base * (2 ** attempt)))


class RetryableError(Exception):
    """Raise for transient failures. Anything else is treated as permanent."""


@dataclass
class TaskResult:
    task_id: Any
    ok: bool
    value: Any = None
    error: BaseException | None = None
    attempts: int = 0
```

```python
class WorkerPool:
    """Fixed number of threads pulling from a queue.Queue.

    Retries RetryableError up to max_retries with backoff; permanent errors fail fast.
    `sleep` is injectable so tests don't actually wait.
    """

    _STOP = object()

    def __init__(self, num_workers: int, max_retries: int = 3, base_delay: float = 0.1,
                 sleep: Callable[[float], None] = time.sleep, rng: random.Random | None = None):
        self.tasks: queue.Queue = queue.Queue()
        self.results: dict[Any, TaskResult] = {}
        self.max_retries, self.base_delay, self.sleep = max_retries, base_delay, sleep
        self.rng = rng or random.Random()
        self._lock = threading.Lock()  # guards results and rng
        self.threads = [threading.Thread(target=self._worker, daemon=True) for _ in range(num_workers)]
        for t in self.threads:
            t.start()

    def submit(self, task_id: Any, fn: Callable[[], Any]) -> None:
        self.tasks.put((task_id, fn))

    def _run_with_retries(self, task_id: Any, fn: Callable[[], Any]) -> TaskResult:
        for attempt in range(self.max_retries + 1):
            try:
                return TaskResult(task_id, True, value=fn(), attempts=attempt + 1)
            except RetryableError as e:
                if attempt == self.max_retries:
                    return TaskResult(task_id, False, error=e, attempts=attempt + 1)
                with self._lock:
                    delay = backoff_delay(attempt, self.base_delay, rng=self.rng)
                self.sleep(delay)
            except Exception as e:  # permanent: don't retry
                return TaskResult(task_id, False, error=e, attempts=attempt + 1)
        raise AssertionError("unreachable")

    def _worker(self) -> None:
        while True:
            item = self.tasks.get()
            try:
                if item is self._STOP:
                    return
                task_id, fn = item
                res = self._run_with_retries(task_id, fn)
                with self._lock:
                    self.results[task_id] = res
            finally:
                self.tasks.task_done()

    def wait(self) -> dict[Any, TaskResult]:
        """Block until every submitted task has finished (including retries)."""
        self.tasks.join()
        return self.results

    def shutdown(self) -> None:
        for _ in self.threads:
            self.tasks.put(self._STOP)  # one poison pill per worker
        for t in self.threads:
            t.join()
```

**Asyncio version.** `asyncio.Semaphore` bounds concurrency; `gather` collects results. Take a coroutine **factory**, because a coroutine object can be awaited only once. Release the semaphore during backoff so a sleeping task doesn't hold a slot.

```python
async def run_async_with_limit(coros_fns: dict[Any, Callable[[], Any]], limit: int,
                               max_retries: int = 3, base_delay: float = 0.01,
                               sleep: Callable[[float], Any] = asyncio.sleep,
                               rng: random.Random | None = None) -> dict[Any, TaskResult]:
    """asyncio version: Semaphore bounds concurrency, gather collects results.

    coros_fns maps task_id -> zero-arg function returning a fresh coroutine
    (a coroutine object can only be awaited once, so retries need a factory).
    """
    sem = asyncio.Semaphore(limit)
    rng = rng or random.Random()

    async def one(task_id: Any, factory: Callable[[], Any]) -> TaskResult:
        for attempt in range(max_retries + 1):
            try:
                async with sem:  # hold the slot only while doing work, not while backing off
                    value = await factory()
                return TaskResult(task_id, True, value=value, attempts=attempt + 1)
            except RetryableError as e:
                if attempt == max_retries:
                    return TaskResult(task_id, False, error=e, attempts=attempt + 1)
                await sleep(backoff_delay(attempt, base_delay, rng=rng))
            except Exception as e:
                return TaskResult(task_id, False, error=e, attempts=attempt + 1)
        raise AssertionError("unreachable")

    results = await asyncio.gather(*(one(tid, f) for tid, f in coros_fns.items()))
    return {r.task_id: r for r in results}
```

**Tests.**

```python
def test_worker_pool() -> None:
    sleeps: list[float] = []
    sleeps_lock = threading.Lock()

    def fake_sleep(d: float) -> None:
        with sleeps_lock:
            sleeps.append(d)

    attempts: Counter = Counter()
    lock = threading.Lock()

    def flaky(task_id: str, fail_n: int) -> Callable[[], str]:
        def fn() -> str:
            with lock:
                attempts[task_id] += 1
                n = attempts[task_id]
            if n <= fail_n:
                raise RetryableError(f"{task_id} attempt {n}")
            return task_id.upper()
        return fn

    def permanent() -> None:
        raise ValueError("bad input")

    pool = WorkerPool(num_workers=4, max_retries=3, base_delay=0.1, sleep=fake_sleep, rng=random.Random(1))
    for i in range(20):
        pool.submit(f"ok{i}", flaky(f"ok{i}", 0))
    pool.submit("retry2", flaky("retry2", 2))
    pool.submit("never", flaky("never", 99))
    pool.submit("perm", permanent)
    res = pool.wait()
    pool.shutdown()
    assert all(res[f"ok{i}"].ok and res[f"ok{i}"].value == f"OK{i}" for i in range(20))
    assert res["retry2"].ok and res["retry2"].attempts == 3
    assert not res["never"].ok and res["never"].attempts == 4 and isinstance(res["never"].error, RetryableError)
    assert not res["perm"].ok and res["perm"].attempts == 1 and isinstance(res["perm"].error, ValueError)
    assert len(sleeps) == 2 + 3                    # retry2 slept twice, never slept 3 times
    assert all(t.is_alive() is False for t in pool.threads)
```

The async test asserts peak concurrency never exceeds the semaphore limit, using `await asyncio.sleep(0)` in each task to force interleaving.

**Complexity.** O(1) per submit/dispatch; task work plus backoff dominates total time. Space O(queued tasks + results).

**Follow-ups.**
- *Backpressure?* → `queue.Queue(maxsize=K)`: `submit` blocks when full (or `put(timeout=...)` and raise).
- *Why not `ThreadPoolExecutor`?* → Use it in production. The manual version shows queue, shutdown, and `task_done` handling. Know both.
- *Threads vs asyncio vs processes?* → Threads for blocking I/O with non-async libraries; asyncio for high-concurrency I/O (thousands of connections); processes for CPU-bound work (the GIL serializes Python bytecode).
- *Per-task timeout?* → asyncio: `asyncio.wait_for(factory(), timeout)` (or `asyncio.timeout` in 3.11+). Threads can't be killed; use cooperative cancellation or a subprocess.
- *Don't retries amplify outages?* → Yes. Add a retry budget (e.g., retries ≤ 10% of requests) or a circuit breaker (Exercise 11). Retry only idempotent operations, or use idempotency keys.
- *Durability across crashes?* → Use a durable queue (SQS, Kafka, DB table with `status` + lease timeout). At-least-once delivery requires idempotent handlers.

## Exercise 6: Task scheduler with dependencies (DAG)

**Spec.** "Given tasks and their dependencies, return a valid execution order, or report a cycle. Part 2: actually run them, with independent tasks in parallel, at most K at a time. Part 3: if a task fails, don't run anything that depends on it."

**Clarifying questions.**
- Input: `task -> prerequisites`, or edges? Can a dependency be unknown?
- On cycle: raise, or return its members?
- On failure: stop all, or continue independent branches?
- Deterministic ties? (Easier tests; use a heap.)

**Design notes.** Kahn's algorithm (`13-graphs-traversal-and-topological-sort.md`): start from indegree-0 nodes; decrement children as parents finish. Leftover nodes with indegree > 0 are on or behind a cycle. For parallel runs, the same bookkeeping drives a `ThreadPoolExecutor`: submit a task when its indegree hits 0; `wait(..., FIRST_COMPLETED)` reacts to whichever finishes first. Only the coordinating thread touches the bookkeeping, so no locks.

```python
class CycleError(Exception):
    pass


def topo_order(deps: dict[str, list[str]]) -> list[str]:
    """Kahn's algorithm. deps[t] = tasks that must finish before t.

    Raises CycleError listing tasks stuck in/behind a cycle.
    """
    nodes = set(deps) | {d for ds in deps.values() for d in ds}
    indegree = {n: 0 for n in nodes}
    children: dict[str, list[str]] = defaultdict(list)
    for t, ds in deps.items():
        for d in ds:
            children[d].append(t)
            indegree[t] += 1
    ready = sorted(n for n in nodes if indegree[n] == 0)  # sorted -> deterministic
    heapq.heapify(ready)
    order = []
    while ready:
        n = heapq.heappop(ready)
        order.append(n)
        for c in children[n]:
            indegree[c] -= 1
            if indegree[c] == 0:
                heapq.heappush(ready, c)
    if len(order) != len(nodes):
        raise CycleError(sorted(n for n in nodes if indegree[n] > 0))
    return order
```

```python
def run_dag(tasks: dict[str, Callable[[], Any]], deps: dict[str, list[str]],
            max_workers: int = 4) -> dict[str, Any]:
    """Run tasks respecting deps, independent ones concurrently (<= max_workers).

    On a task failure: stop scheduling new tasks, let in-flight ones finish,
    mark descendants as skipped. Returns {task: value | exception | 'SKIPPED'}.
    """
    topo_order({t: deps.get(t, []) for t in tasks})  # validate up front: fail before running anything
    indegree = {t: len(deps.get(t, [])) for t in tasks}
    children: dict[str, list[str]] = defaultdict(list)
    for t in tasks:
        for d in deps.get(t, []):
            if d not in tasks:
                raise KeyError(f"{t} depends on unknown task {d}")
            children[d].append(t)

    results: dict[str, Any] = {}
    failed = False
    with ThreadPoolExecutor(max_workers=max_workers) as ex:
        in_flight = {ex.submit(tasks[t]): t for t in sorted(t for t in tasks if indegree[t] == 0)}
        while in_flight:
            done, _ = wait(in_flight, return_when=FIRST_COMPLETED)
            for fut in done:
                t = in_flight.pop(fut)
                try:
                    results[t] = fut.result()
                except Exception as e:
                    results[t] = e
                    failed = True
                    continue
                if failed:
                    continue
                for c in children[t]:
                    indegree[c] -= 1
                    if indegree[c] == 0:
                        in_flight[ex.submit(tasks[c])] = c
    for t in tasks:
        results.setdefault(t, "SKIPPED")
    return results
```

**Tests (excerpt).**

```python
deps = {"deploy": ["build", "test"], "test": ["build"], "build": ["fetch"], "lint": ["fetch"]}
assert topo_order(deps) == ["fetch", "build", "lint", "test", "deploy"]
expect_raises(CycleError, topo_order, {"a": ["b"], "b": ["c"], "c": ["a"], "d": []})

def boom(): raise RuntimeError("x")
res = run_dag({"a": lambda: 1, "b": boom, "c": lambda: 3, "d": lambda: 4},
              {"c": ["b"], "d": ["c"]}, max_workers=1)
assert res["a"] == 1 and isinstance(res["b"], RuntimeError)
assert res["c"] == "SKIPPED" and res["d"] == "SKIPPED"
```

The companion file also runs a fan-out/fan-in DAG with `max_workers=2`, logs start/end events under a lock, and asserts peak concurrency ≤ 2 and that each task started after all its dependencies ended.

**Complexity.** Ordering O(V + E) (O(V log V + E) with the tie-breaking heap). Execution uses the same bookkeeping; with enough workers, wall time ≈ critical path length.

**Follow-ups.**
- *Continue independent branches after a failure?* → Drop the global `failed` flag; mark only the failed task's descendants (BFS) skipped and keep scheduling the rest.
- *Per-task retries?* → Wrap each callable in Exercise 5's retry loop.
- *Priorities?* → Hold ready tasks in a priority heap instead of submitting immediately (longest remaining path first heuristically minimizes makespan).
- *Survive a coordinator crash?* → Persist task states (pending/running/done/failed) in a DB; workers claim tasks with leases; the scheduler reconciles over that table (Airflow/Argo's model).
- *Report the cycle path?* → Three-color DFS: an edge to a gray node is a back edge; rebuild the path from the DFS stack.

## Exercise 7: Reconciliation loop

**Spec.** "You have desired state (a map of name → spec) and an API to list and mutate actual state. Write a controller that converges actual to desired. Operations can fail transiently; retry with backoff. Part 2: make repeated runs safe."

A miniature Kubernetes controller, popular at infra companies because it opens real-systems talk: level- vs edge-triggered, idempotency, work queues, resync.

**Clarifying questions.**
- Reconcile per object (keyed work queue) or whole world? (Per object, like controller-runtime.)
- Are mutations idempotent? If not, can I make them so (create-or-update)?
- Delete what's in actual but not desired? (Real controllers scope this with labels/owner references to avoid deleting what they don't own.)
- Max backoff? Periodic resync?

**Design notes.**
- **Pure diff**: `diff_state(desired, actual) -> actions`, no I/O, table-testable.
- **Level-triggered**: each reconcile reads current state and decides from it, never trusting "what changed" events to be complete, so missed or duplicate events are harmless.
- **Idempotent actions**: create on existing = update; delete on missing = no-op. Replay after a crash is safe.
- **Deduplicating work queue**: a key is queued at most once (`queued` set), like client-go's workqueue, so many events for one object collapse into one reconcile.
- **Per-key exponential backoff**: on error, requeue after `base * 2**failures`, capped; success resets. A failing key doesn't block others.
- Injected clock: the queue is a heap of `(ready_at, name)`; `process_ready` handles everything due at `clock()`.

```python
@dataclass(frozen=True)
class Action:
    kind: str  # "create" | "update" | "delete"
    name: str
    spec: Any = None


def diff_state(desired: dict[str, Any], actual: dict[str, Any]) -> list[Action]:
    """Level-triggered diff. Pure function: same inputs -> same actions."""
    actions = []
    for name in sorted(desired):
        if name not in actual:
            actions.append(Action("create", name, desired[name]))
        elif actual[name] != desired[name]:
            actions.append(Action("update", name, desired[name]))
    for name in sorted(actual):
        if name not in desired:
            actions.append(Action("delete", name))
    return actions
```

```python
class FakeCluster:
    """Actual-state backend. Operations are idempotent (create on existing = update,
    delete on missing = no-op), so replaying an action after a crash is safe."""

    def __init__(self, fail_times: dict[str, int] | None = None):
        self.objects: dict[str, Any] = {}
        self.fail_times = dict(fail_times or {})  # name -> remaining injected failures
        self.calls = 0

    def list(self) -> dict[str, Any]:
        return dict(self.objects)

    def apply(self, action: Action) -> None:
        self.calls += 1
        if self.fail_times.get(action.name, 0) > 0:
            self.fail_times[action.name] -= 1
            raise RuntimeError(f"transient error on {action.name}")
        if action.kind in ("create", "update"):
            self.objects[action.name] = action.spec
        elif action.kind == "delete":
            self.objects.pop(action.name, None)
```

```python
class Reconciler:
    """Work queue of object names, per-key exponential backoff on failure.

    Like controller-runtime: reconcile(name) reads desired+actual for ONE key,
    converges it, and on error requeues with backoff. Dedup: a key is in the
    queue at most once.
    """

    def __init__(self, desired: dict[str, Any], cluster: FakeCluster,
                 clock: Callable[[], float] = time.monotonic,
                 base_delay: float = 1.0, max_delay: float = 60.0):
        self.desired, self.cluster, self.clock = desired, cluster, clock
        self.base_delay, self.max_delay = base_delay, max_delay
        self.queue: list[tuple[float, str]] = []  # (ready_at, name) min-heap
        self.queued: set[str] = set()
        self.failures: dict[str, int] = defaultdict(int)

    def enqueue(self, name: str, delay: float = 0.0) -> None:
        if name in self.queued:
            return
        self.queued.add(name)
        heapq.heappush(self.queue, (self.clock() + delay, name))

    def enqueue_all(self) -> None:
        """Periodic resync: enqueue every key in desired OR actual."""
        for name in sorted(set(self.desired) | set(self.cluster.list())):
            self.enqueue(name)

    def reconcile(self, name: str) -> None:
        actual = self.cluster.list()
        want = {name: self.desired[name]} if name in self.desired else {}
        have = {name: actual[name]} if name in actual else {}
        for action in diff_state(want, have):
            self.cluster.apply(action)

    def process_ready(self) -> int:
        """Process every item whose ready_at <= now. Returns items processed."""
        n = 0
        while self.queue and self.queue[0][0] <= self.clock():
            _, name = heapq.heappop(self.queue)
            self.queued.discard(name)
            n += 1
            try:
                self.reconcile(name)
                self.failures.pop(name, None)  # success resets backoff
            except Exception:
                delay = min(self.max_delay, self.base_delay * 2 ** self.failures[name])
                self.failures[name] += 1
                self.enqueue(name, delay)
        return n
```

**Tests.**

```python
def test_reconciler() -> None:
    assert diff_state({"a": 1, "b": 2}, {"b": 3, "c": 4}) == [
        Action("create", "a", 1), Action("update", "b", 2), Action("delete", "c")]
    assert diff_state({"a": 1}, {"a": 1}) == []
    cases = [
        ({"a": 1}, {}, [Action("create", "a", 1)]),
        ({"a": 1}, {"a": 2}, [Action("update", "a", 1)]),
        ({}, {"a": 1}, [Action("delete", "a")]),
        ({"a": 1}, {"a": 1}, []),
    ]
    for desired, actual, want in cases:
        assert diff_state(desired, actual) == want, (desired, actual)

    clk = FakeClock()
    cluster = FakeCluster(fail_times={"web": 2})
    cluster.objects = {"web": {"replicas": 1}, "old": {"replicas": 1}}
    desired = {"web": {"replicas": 3}, "api": {"replicas": 2}}
    r = Reconciler(desired, cluster, clock=clk, base_delay=1.0)
    r.enqueue_all()
    r.enqueue("api")                   # dedup: already queued
    assert len(r.queue) == 3
    assert r.process_ready() == 3      # web fails once -> requeued at t+1
    assert cluster.objects == {"web": {"replicas": 1}, "api": {"replicas": 2}}
    assert r.process_ready() == 0      # not ready yet
    clk.advance(1)
    assert r.process_ready() == 1      # fails again -> requeue at t+2
    clk.advance(1)
    assert r.process_ready() == 0      # backoff doubled
    clk.advance(1)
    assert r.process_ready() == 1      # succeeds
    assert cluster.objects == desired and "web" not in r.failures

    # Idempotence: a full resync on a converged system issues no writes.
    calls = cluster.calls
    r.enqueue_all()
    r.process_ready()
    assert cluster.calls == calls

    # Drift: someone edits actual state out-of-band; resync repairs it.
    cluster.objects["api"] = {"replicas": 99}
    r.enqueue_all()
    r.process_ready()
    assert cluster.objects == desired

    # Backoff caps at max_delay.
    r2 = Reconciler({"x": 1}, FakeCluster(fail_times={"x": 100}), clock=clk, base_delay=1, max_delay=4)
    r2.enqueue("x")
    delays = []
    for _ in range(5):
        r2.process_ready()
        ready_at = r2.queue[0][0]
        delays.append(ready_at - clk())
        clk.now = ready_at
    assert delays == [1, 2, 4, 4, 4]
```

**Complexity.** `diff_state` O(n log n) (sorted for deterministic output). Per-key reconcile O(1) plus the API call; queue ops O(log q).

**Mapping to Kubernetes controllers.**

| Here | controller-runtime / client-go |
|---|---|
| `desired[name]` | The object's `.spec` |
| `cluster.list()` | Informer cache (local, watch-fed), not live API reads |
| `enqueue(name)` + `queued` set | `workqueue` dedup: a key is never in the queue twice |
| `failures` + `base * 2**n` | `RateLimitingQueue` / `ItemExponentialFailureRateLimiter` |
| `enqueue_all()` | Periodic resync / `SyncPeriod` |
| `reconcile(name)` | `Reconcile(ctx, req)` returning `Result{RequeueAfter}` or error |

**Follow-ups.**
- *Multiple workers?* → The work queue gives each key to one worker at a time (client-go tracks "processing" and "dirty" sets); different keys run in parallel.
- *Two controller replicas?* → Leader election via a Lease object, or shard keys across replicas.
- *Conflicting writes?* → Optimistic concurrency on `resourceVersion`: updating a changed object fails with 409 Conflict; requeue and retry from a fresh read.
- *Deleting safely?* → Finalizers: mark for deletion, clean external resources, remove the finalizer. Owner references + garbage collection handle children.
- *Status reporting?* → Write observed state to `.status` with `observedGeneration` so users can tell if the latest spec was seen.
- *Thundering requeue after an outage?* → Capped backoff with jitter, plus a global queue rate limiter (client-go combines per-item exponential with an overall token bucket).

## Exercise 8: Consistent hashing with virtual nodes

**Spec.** "Distribute keys across cache nodes. Adding or removing a node should move as few keys as possible. Implement `add_node`, `remove_node`, `get_node(key)`. Then measure how many keys move when a node is added."

**Clarifying questions.**
- Weighted nodes (different capacities)?
- Replication (N distinct nodes per key)?
- Stable across processes? (Yes, so avoid Python's `hash()`, which salts strings per process.)

**Design notes.** Hash each node to many points (virtual nodes) on a ring of 64-bit integers. A key belongs to the first point clockwise from its hash: `bisect_right` on a sorted list, wrapping to index 0. Changing N under `hash(key) % N` remaps about (N-1)/N of keys; adding node N+1 to a ring moves about 1/(N+1). Virtual nodes smooth load, since one point per node gives wildly uneven arcs.

```python
def _hash(s: str) -> int:
    # Stable across processes (unlike built-in hash(), which is salted for str).
    return int.from_bytes(hashlib.md5(s.encode()).digest()[:8], "big")


class HashRing:
    """Sorted list of (hash, node) points; key goes to first point clockwise."""

    def __init__(self, nodes: Iterable[str] = (), vnodes: int = 100):
        self.vnodes = vnodes
        self.points: list[int] = []
        self.owner: dict[int, str] = {}
        self.nodes: set[str] = set()
        for n in nodes:
            self.add_node(n)

    def add_node(self, node: str) -> None:
        if node in self.nodes:
            return
        self.nodes.add(node)
        for i in range(self.vnodes):
            h = _hash(f"{node}#{i}")
            if h in self.owner:  # astronomically rare collision: skip
                continue
            self.owner[h] = node
            bisect.insort(self.points, h)

    def remove_node(self, node: str) -> None:
        if node not in self.nodes:
            return
        self.nodes.discard(node)
        for i in range(self.vnodes):
            h = _hash(f"{node}#{i}")
            if self.owner.get(h) == node:
                del self.owner[h]
                self.points.pop(bisect.bisect_left(self.points, h))

    def get_node(self, key: str) -> str:
        if not self.points:
            raise LookupError("empty ring")
        i = bisect.bisect_right(self.points, _hash(key)) % len(self.points)  # wrap around
        return self.owner[self.points[i]]

    def get_nodes(self, key: str, n: int) -> list[str]:
        """Replication follow-up: n distinct physical nodes walking clockwise."""
        if n > len(self.nodes):
            raise ValueError("not enough nodes")
        out: list[str] = []
        i = bisect.bisect_right(self.points, _hash(key))
        while len(out) < n:
            node = self.owner[self.points[i % len(self.points)]]
            if node not in out:
                out.append(node)
            i += 1
        return out


def moved_fraction(before: HashRing, after: HashRing, keys: list[str]) -> float:
    return sum(before.get_node(k) != after.get_node(k) for k in keys) / len(keys)
```

**Tests (excerpt).**

```python
keys = [f"key-{i}" for i in range(10000)]
ring = HashRing(["n1", "n2", "n3", "n4"], vnodes=100)
load = Counter(ring.get_node(k) for k in keys)
assert max(load.values()) / min(load.values()) < 1.6    # vnodes smooth the load

bigger = HashRing(["n1", "n2", "n3", "n4", "n5"], vnodes=100)
frac = moved_fraction(ring, bigger, keys)
assert 0.1 < frac < 0.3                                 # ideal is 1/5; mod-N would move ~80%
# Every moved key moved TO the new node.
assert all(bigger.get_node(k) == "n5" for k in keys if ring.get_node(k) != bigger.get_node(k))
```

**Complexity.** `get_node` O(log(N·V)). `add_node`/`remove_node` O(V · N·V) because `insort` shifts the list; fine for hundreds of nodes. For frequent membership changes, collect all points and sort once: O(N·V log(N·V)). Space O(N·V).

**Follow-ups.**
- *Weighted nodes?* → Vnode count proportional to capacity.
- *Replication?* → `get_nodes(key, n)`: walk clockwise collecting distinct physical nodes (Dynamo's preference list); for rack/zone awareness, skip nodes in already-used zones.
- *How many vnodes?* → Load standard deviation shrinks roughly like 1/sqrt(V). 100–200 per node is typical; more costs memory and slows membership changes.
- *Alternatives?* → Rendezvous (highest random weight) hashing scores every node per key and picks the max: O(N) lookup, no ring, minimal movement. Jump consistent hash: O(log N), no memory, but nodes are added/removed only at the end (fine for shards, not arbitrary failures).
- *Hot keys?* → Hashing can't spread one huge key. Replicate it and read from any copy, or add a local cache tier.
- *Data migration on add?* → Only keys in arcs the new vnodes take over move; stream them from their old owners before switching reads.

## Exercise 9: Log processing

**Spec.** "Given an access log file, report: total requests, the top-k endpoints by count, the error rate (5xx), and p50/p95/p99 latency. Some lines are malformed."

**Clarifying questions.**
- Exact line format? (Get a sample line; write the regex against it.)
- Are `/users?id=1` and `/users?id=2` one endpoint? (Usually: strip the query string. Templating `/users/123` → `/users/:id` is a good follow-up.)
- Is 4xx an error? (Usually "error rate" means 5xx; ask.)
- Fits in memory? Single file or a stream?
- Nearest-rank or interpolated percentiles?

**Design notes.**
- Compiled regex; return `None` on mismatch and count malformed lines. A parser that throws on line 10 million is a production incident.
- Iterate lazily (`for line in f`) so the file never loads whole. Counters are O(distinct endpoints).
- Exact percentiles need all latencies (O(n) memory) and a sort. Nearest-rank: index `ceil(p/100 * n) - 1`.
- Top-k: `heapq.nsmallest(k, ..., key=(-count, path))` is O(m log k) with deterministic ties; `Counter.most_common(k)` breaks ties by insertion order.

```python
LOG_RE = re.compile(
    r'^(?P<ip>\S+) \S+ \S+ \[(?P<ts>[^\]]+)\] "(?P<method>[A-Z]+) (?P<path>\S+) [^"]*" '
    r'(?P<status>\d{3}) (?P<size>\d+|-) (?P<latency>[\d.]+)$'
)


@dataclass
class LogEntry:
    method: str
    path: str
    status: int
    latency: float


def parse_line(line: str) -> LogEntry | None:
    m = LOG_RE.match(line.strip())
    if not m:
        return None  # malformed: caller counts it, never crash on bad input
    path = m["path"].split("?", 1)[0]  # group /a?x=1 and /a?x=2 together
    return LogEntry(m["method"], path, int(m["status"]), float(m["latency"]))


def percentile(sorted_vals: list[float], p: float) -> float:
    """Nearest-rank percentile on pre-sorted data. p in (0, 100]."""
    if not sorted_vals:
        raise ValueError("no data")
    rank = max(1, math.ceil(p / 100 * len(sorted_vals)))
    return sorted_vals[rank - 1]


def analyze_logs(lines: Iterable[str], k: int = 3) -> dict[str, Any]:
    counts: Counter = Counter()
    latencies: list[float] = []
    errors = malformed = total = 0
    for line in lines:  # streaming over lines; only latencies are kept
        e = parse_line(line)
        if e is None:
            malformed += 1
            continue
        total += 1
        counts[e.path] += 1
        latencies.append(e.latency)
        if e.status >= 500:
            errors += 1
    latencies.sort()
    pct = {f"p{p}": percentile(latencies, p) for p in (50, 95, 99)} if latencies else {}
    # Tie-break by path for deterministic output; most_common's ties are insertion-ordered.
    top = heapq.nsmallest(k, counts.items(), key=lambda kv: (-kv[1], kv[0]))
    return {"total": total, "malformed": malformed, "top_k": top,
            "error_rate": errors / total if total else 0.0, **pct}
```

**Streaming percentile approximation.** For latencies that don't fit in memory or arrive live, use fixed log-spaced buckets. Each bucket's upper bound is `growth` times the previous, so relative error is at most `growth - 1` (1% here), the core of Prometheus histograms and DDSketch. Histograms **merge** by adding counts, so per-host results combine exactly; exact percentiles can't be averaged across hosts.

```python
class StreamingHistogram:
    """Bounded-memory percentile approximation with fixed log-spaced buckets.

    Relative error bounded by bucket growth factor (like Prometheus histograms / DDSketch idea).
    """

    def __init__(self, min_val: float = 1e-4, max_val: float = 1e4, growth: float = 1.05):
        self.min_val, self.growth = min_val, growth
        self.nbuckets = int(math.log(max_val / min_val, growth)) + 2
        self.counts = [0] * self.nbuckets
        self.n = 0
        self._log_g = math.log(growth)

    def add(self, v: float) -> None:
        i = 0 if v <= self.min_val else int(math.log(v / self.min_val) / self._log_g) + 1
        self.counts[min(i, self.nbuckets - 1)] += 1
        self.n += 1

    def quantile(self, p: float) -> float:
        target = max(1, math.ceil(p / 100 * self.n))
        running = 0
        for i, c in enumerate(self.counts):
            running += c
            if running >= target:
                return self.min_val * self.growth ** i  # upper bound of bucket i
        raise ValueError("no data")
```

**Tests (excerpt).**

```python
r = analyze_logs(lines, k=2)
assert r["total"] == 5 and r["malformed"] == 2
assert r["top_k"] == [("/api/orders", 2), ("/api/users", 2)]   # tie broken by path
assert r["error_rate"] == 2 / 5
assert r["p50"] == 0.2 and r["p95"] == 1.5 and r["p99"] == 1.5

h = StreamingHistogram(growth=1.01)        # 20k exponential samples
for p in (50, 95, 99):
    assert abs(h.quantile(p) - percentile(data, p)) / percentile(data, p) < 0.02
```

**Complexity.** Parse O(total bytes). Exact percentiles O(n log n) time, O(n) space (`quickselect` gives O(n) per percentile). Top-k O(m log k) for m endpoints. Histogram: O(1) add, O(buckets) query, O(buckets) space independent of n.

**Streaming approximations to name.**

| Need | Technique | Memory |
|---|---|---|
| Percentiles | Log-bucket histogram, t-digest, DDSketch | O(buckets) / O(compression) |
| Top-k heavy hitters | Count-Min Sketch + min-heap of size k; Misra-Gries / Space-Saving | O(k) or O(width × depth) |
| Distinct count | HyperLogLog | ~KB for ~1% error |
| Percentiles over a sliding time window | Ring of per-minute histograms; merge the last N | O(N × buckets) |

**Follow-ups.**
- *1 TB file?* → Stream it; shard by file or byte range across processes (`multiprocessing`), each producing a `Counter` + histogram, then merge (histograms merge; exact percentiles don't).
- *Live log (tail -f)?* → Loop reading from the last offset; buffer a partial last line until its newline; detect rotation by inode change.
- *Per-endpoint p99?* → One histogram per endpoint. Template paths first, or cardinality explodes memory.
- *Out-of-order timestamps across sources?* → Merge sorted per-source streams (Exercise 10) with a bounded reorder buffer.

## Exercise 10: Iterators and generators

**Spec (three small parts, often a warm-up).** (a) "Write an iterator that flattens arbitrarily nested lists of integers." (b) "Wrap an iterator to add `peek()`." (c) "Merge k sorted streams lazily; the streams may be huge or infinite."

**Clarifying questions.**
- Can nested lists be empty? (Yes; this breaks naive `has_next`.)
- `has_next()` (Java-style) or Python's iterator protocol?
- Can values be `None` or falsy? (Use a sentinel, not truthiness.)
- Streams sorted ascending? Stable merge required?

**Design notes.**
- **Flatten**: a stack of iterators (not lists) copies nothing and stays lazy. `_advance()` pre-fetches the next integer, making `has_next()` a flag and empty sublists like `[[], [[]]]` correct.
- **Peeking**: buffer one item, with a private "exhausted" sentinel because `None` may be real.
- **K-way merge**: heap of `(value, stream_index)`, one entry per live stream (`11-heaps-top-k-and-merging.md`). The index breaks ties so iterators are never compared. In production, use `heapq.merge(*streams)`.

```python
class NestedIterator:
    """Flatten arbitrarily nested lists lazily using a stack of iterators.

    Invariant after _advance(): self._next holds the next int, or stack is empty.
    """

    def __init__(self, nested: list):
        self.stack: list[Iterator] = [iter(nested)]
        self._has = False
        self._next: Any = None
        self._advance()

    def _advance(self) -> None:
        self._has = False
        while self.stack:
            try:
                x = next(self.stack[-1])
            except StopIteration:
                self.stack.pop()
                continue
            if isinstance(x, list):
                self.stack.append(iter(x))
            else:
                self._next, self._has = x, True
                return

    def __iter__(self):
        return self

    def __next__(self) -> Any:
        if not self._has:
            raise StopIteration
        val = self._next
        self._advance()
        return val

    def has_next(self) -> bool:
        return self._has


def flatten_gen(nested: list) -> Iterator:
    """Generator version. Recursion depth = nesting depth."""
    for x in nested:
        if isinstance(x, list):
            yield from flatten_gen(x)
        else:
            yield x
```

```python
class PeekingIterator:
    """Wrap any iterator; peek() returns the next item without consuming it."""

    _END = object()

    def __init__(self, it: Iterable):
        self.it = iter(it)
        self._buf = next(self.it, self._END)

    def has_next(self) -> bool:
        return self._buf is not self._END

    def peek(self) -> Any:
        if self._buf is self._END:
            raise StopIteration
        return self._buf

    def __iter__(self):
        return self

    def __next__(self) -> Any:
        if self._buf is self._END:
            raise StopIteration
        val, self._buf = self._buf, next(self.it, self._END)
        return val
```

```python
def merge_sorted_streams(streams: list[Iterable]) -> Iterator:
    """Manual k-way merge. Heap holds one (value, stream_idx) per non-empty stream.

    stream_idx breaks ties so values need not be comparable beyond <, and
    iterators themselves are never compared.
    """
    iters = [iter(s) for s in streams]
    heap = []
    for i, it in enumerate(iters):
        first = next(it, _MISSING)
        if first is not _MISSING:
            heap.append((first, i))
    heapq.heapify(heap)
    while heap:
        val, i = heap[0]
        yield val
        nxt = next(iters[i], _MISSING)
        if nxt is _MISSING:
            heapq.heappop(heap)
        else:
            heapq.heapreplace(heap, (nxt, i))  # pop + push in one sift
```

**Tests (excerpt).**

```python
assert list(NestedIterator([1, [2, [3, []], 4], [], [[5]], 6])) == [1, 2, 3, 4, 5, 6]
assert list(NestedIterator([[], [[]]])) == []
p = PeekingIterator(iter([1, 2, 3]))
assert p.peek() == 1 and next(p) == 1 and p.peek() == 2
assert list(PeekingIterator([None, 0])) == [None, 0]     # falsy values are real items
streams = [[1, 4, 7], [2, 5, 8], [], [0, 9], [3, 3]]
assert list(merge_sorted_streams(streams)) == list(heapq.merge(*streams))
m = merge_sorted_streams([evens(), odds()])               # infinite generators
assert [next(m) for _ in range(6)] == [0, 1, 2, 3, 4, 5]
```

**Complexity.** Flatten: O(1) amortized per element, O(depth) space. Peek: O(1). Merge: O(log k) per element, O(k) space.

**Follow-ups.**
- *Recursion limit on deep nesting?* → The stack version doesn't recurse; `flatten_gen` recurses per level and hits Python's default limit (~1000).
- *Merge by key (e.g., timestamp)?* → Push `(key(item), index, item)`, or `heapq.merge(*streams, key=...)`.
- *Network streams that stall?* → Lazy merge blocks on the slowest stream. Use per-stream buffers with timeouts, or a watermark: emit up to the minimum timestamp seen across streams.
- *Resettable or thread-safe iterator?* → Generators aren't thread safe; lock `__next__`, or give each consumer its own iterator (`itertools.tee` for fan-out, at a memory cost).

## Exercise 11: Circuit breaker

**Spec.** "Wrap calls to a flaky dependency. After N consecutive failures, stop calling it and fail fast for T seconds. After T, allow a single trial call: success closes the circuit, failure opens it again."

**Clarifying questions.**
- Consecutive failures, or windowed failure rate? (Start with consecutive.)
- Half-open: one trial or several? What do concurrent callers get meanwhile?
- Which exceptions count? (A 404 probably isn't a dependency failure.)
- Fail fast by raising, or return a fallback?

**Design notes.** Three states:

```
         N consecutive failures
CLOSED ---------------------------> OPEN
   ^                                 |  reset_timeout elapsed
   |  trial succeeds                 v
   +-------------------------- HALF_OPEN --trial fails--> OPEN (timer restarts)
```

- OPEN → HALF_OPEN is time-driven; compute it on read (`_current_state`) instead of using a timer thread.
- `trial_in_flight` admits exactly one half-open probe; concurrent callers fail fast.
- **Never hold the lock while calling `fn`**, or a slow dependency serializes every caller.

```python
class CircuitOpenError(Exception):
    pass


class CircuitBreaker:
    """CLOSED --(failure_threshold consecutive failures)--> OPEN
    OPEN --(reset_timeout elapsed)--> HALF_OPEN (allow one trial call)
    HALF_OPEN --success--> CLOSED, --failure--> OPEN (timer restarts)
    """

    CLOSED, OPEN, HALF_OPEN = "closed", "open", "half_open"

    def __init__(self, failure_threshold: int = 3, reset_timeout: float = 30.0,
                 clock: Callable[[], float] = time.monotonic):
        self.failure_threshold, self.reset_timeout, self.clock = failure_threshold, reset_timeout, clock
        self._state = self.CLOSED
        self.failures = 0
        self.opened_at = 0.0
        self.trial_in_flight = False
        self._lock = threading.Lock()

    @property
    def state(self) -> str:
        with self._lock:
            return self._current_state()

    def _current_state(self) -> str:
        # OPEN -> HALF_OPEN is time-driven, so compute it on read.
        if self._state == self.OPEN and self.clock() - self.opened_at >= self.reset_timeout:
            self._state = self.HALF_OPEN
            self.trial_in_flight = False
        return self._state

    def _trip(self) -> None:
        self._state, self.opened_at, self.failures = self.OPEN, self.clock(), 0

    def call(self, fn: Callable[[], Any]) -> Any:
        with self._lock:
            st = self._current_state()
            if st == self.OPEN:
                raise CircuitOpenError("circuit open")
            if st == self.HALF_OPEN:
                if self.trial_in_flight:
                    raise CircuitOpenError("trial in progress")
                self.trial_in_flight = True
        try:
            result = fn()  # never hold the lock while calling out
        except Exception:
            with self._lock:
                if self._state == self.HALF_OPEN:
                    self._trip()
                else:
                    self.failures += 1
                    if self.failures >= self.failure_threshold:
                        self._trip()
            raise
        with self._lock:
            self._state, self.failures, self.trial_in_flight = self.CLOSED, 0, False
        return result
```

**Tests (excerpt).**

```python
clk = FakeClock()
cb = CircuitBreaker(failure_threshold=2, reset_timeout=10, clock=clk)
expect_raises(IOError, cb.call, fail)
assert cb.call(lambda: 1) == 1                    # success resets the consecutive count
expect_raises(IOError, cb.call, fail)
expect_raises(IOError, cb.call, fail)
assert cb.state == "open"
expect_raises(CircuitOpenError, cb.call, lambda: 1)
clk.advance(10)
assert cb.state == "half_open"
expect_raises(IOError, cb.call, fail)             # trial fails -> open again
clk.advance(10)
assert cb.call(lambda: "recovered") == "recovered" and cb.state == "closed"
```

The companion file also tests that half-open admits only one concurrent trial: one thread's trial blocks on an `Event` while the main thread asserts a second call gets `CircuitOpenError`. Events, not sleeps, make the interleaving deterministic.

**Complexity.** O(1) per call.

**Follow-ups.**
- *Failure rate instead of consecutive failures?* → Ring buffer of the last N outcomes (or per-second buckets, like Hystrix/resilience4j); trip when failures/total ≥ X% and total meets a minimum volume.
- *Per-host breakers?* → A dict of breakers keyed by host, so one bad replica doesn't cut off healthy ones.
- *Slow calls?* → Count calls over a latency threshold as failures; time out `fn`.
- *Interaction with retries?* → Retries go *inside* the breaker's view, so a retry storm against an open circuit fails fast instead of piling on.
- *Metrics?* → Emit state transitions, the most useful alert signal.
- *Distributed?* → Usually local per client, since shared state adds a dependency to protect a dependency. Service meshes (Envoy outlier detection) do it per proxy.

## Exercise 12: In-memory file system (brief)

**Spec.** "Implement `mkdir(path)` (creates parents like `mkdir -p`), `ls(path)` (sorted names; for a file, just its name), `write(path, content)`, `read(path)`."

**Clarifying questions.** Absolute paths only? Does `write` create missing parents? Append or overwrite? Error type for reading a directory?

**Design notes.** A trie keyed by path components, with one node type (`content is None` means directory). One `_walk` helper handles traversal, optional creation, and "component is a file" errors, keeping each public method to a few lines. Built-in exceptions (`FileNotFoundError`, `IsADirectoryError`) make the contract self-documenting.

```python
@dataclass
class _FSNode:
    children: dict[str, "_FSNode"] = field(default_factory=dict)
    content: str | None = None  # None -> directory

    @property
    def is_file(self) -> bool:
        return self.content is not None


class FileSystem:
    def __init__(self):
        self.root = _FSNode()

    @staticmethod
    def _parts(path: str) -> list[str]:
        if not path.startswith("/"):
            raise ValueError("absolute paths only")
        return [p for p in path.split("/") if p]

    def _walk(self, parts: list[str], create: bool = False) -> _FSNode:
        node = self.root
        for p in parts:
            if node.is_file:
                raise NotADirectoryError(p)
            if p not in node.children:
                if not create:
                    raise FileNotFoundError("/" + "/".join(parts))
                node.children[p] = _FSNode()
            node = node.children[p]
        return node

    def mkdir(self, path: str) -> None:
        """mkdir -p semantics."""
        node = self._walk(self._parts(path), create=True)
        if node.is_file:
            raise FileExistsError(path)

    def ls(self, path: str) -> list[str]:
        parts = self._parts(path)
        node = self._walk(parts)
        return [parts[-1]] if node.is_file else sorted(node.children)

    def write(self, path: str, content: str, append: bool = False) -> None:
        parts = self._parts(path)
        if not parts:
            raise IsADirectoryError("/")
        parent = self._walk(parts[:-1], create=True)
        if parent.is_file:
            raise NotADirectoryError(path)
        node = parent.children.get(parts[-1])
        if node is None:
            node = parent.children[parts[-1]] = _FSNode(content="")
        elif not node.is_file:
            raise IsADirectoryError(path)
        node.content = node.content + content if append else content

    def read(self, path: str) -> str:
        node = self._walk(self._parts(path))
        if not node.is_file:
            raise IsADirectoryError(path)
        return node.content
```

**Complexity.** O(path depth) per op, plus O(c log c) for sorting c children in `ls`.

**Follow-ups.**
- *`rm -r`?* → Delete the subtree from its parent's dict.
- *`mv`?* → Detach and reattach the node; reject moving into its own subtree.
- *Size accounting?* → Store sizes; update ancestors on write.
- *Thread safety?* → One lock, or per-node locks acquired root to leaf.
- *Globbing?* → DFS with `fnmatch` on each component.

## Testing strategy

**Inject time; never sleep for it.** Components with timeouts, TTLs, rates, or backoff take `clock` (plus `sleep` if they wait); tests call `FakeClock.advance()`. `time.sleep(1)` in a test is slow and flaky; an injected clock takes microseconds and hits exact boundaries (`advance(9.999)` then `advance(0.001)`).

**Inject randomness.** Pass `rng: random.Random`; seed it in tests, or assert bounds (`0 <= delay <= cap`) instead of exact values.

**Pure core, thin I/O shell.** Put logic in pure, table-testable functions (`diff_state`, `parse_line`) wrapped by I/O (`Reconciler`, `analyze_logs`).

**Table-driven asserts.** Quick to write, easy to extend:

```python
cases = [
    ({"a": 1}, {}, [Action("create", "a", 1)]),
    ({"a": 1}, {"a": 2}, [Action("update", "a", 1)]),
    ({}, {"a": 1}, [Action("delete", "a")]),
    ({"a": 1}, {"a": 1}, []),
]
for desired, actual, want in cases:
    assert diff_state(desired, actual) == want, (desired, actual)
```

Inputs in the assert message make a failure name its case.

**Edge cases checklist.** Empty input; single element; capacity 1; exact boundary times; duplicate keys / overwrites; falsy values (`0`, `None`, `""`); malformed input; ops on missing keys; methods called in the wrong state (rollback with no txn).

**Testing concurrency.**
- **Synchronize, don't sleep**: `threading.Event` forces an interleaving (half-open breaker test); `Queue.join()` and `Thread.join()` wait for completion.
- **Invariant hammer tests**: several threads run seeded random ops, then assert invariants (`len(cache) <= capacity`, no exceptions, consistent values). They catch races probabilistically but can't prove their absence.
- **Record events under a lock**, then assert ordering (DAG test: each task's `start` follows its dependencies' `end`).
- **Measure peak concurrency** with a locked active counter to verify a semaphore/pool limit.
- Inject `sleep` into retry code so tests skip real backoff.
- Keep real sleeps to milliseconds, only where true overlap is needed; never assert on wall-clock duration.
- Worker-thread exceptions don't fail a test on their own: capture them into results and assert in the main thread.

## Self-check

1. Why does `ThreadSafeLRU.get` need the same exclusive lock as `put`, and what are two ways to reduce contention?
2. In the TTL store, what bug appears if `cleanup` deletes by key without checking the version? Write the 4-line test that catches it.
3. In the transactional store, why does the undo log record a key's old value only on the first write within a transaction? What changes if `COMMIT` should close only the innermost level?
4. Compare token bucket, sliding window log, and sliding window counter on memory, accuracy, and burst behavior. Which would you put in front of a public API with millions of client keys?
5. Why full jitter? What goes wrong if 10,000 clients retry with pure exponential backoff after a shared outage?
6. Why does the asyncio worker take coroutine factories instead of coroutines, and why release the semaphore during backoff?
7. What makes a reconciler level-triggered, and why does that make missed or duplicate events harmless? Name three client-go mechanisms that correspond to parts of `Reconciler`.
8. Adding a 5th node to a 4-node ring moves about what fraction of keys? With `hash % N`? Why are virtual nodes needed?
9. Why can't you average per-host p99s, and what data structure lets you combine them correctly?
10. Why must the circuit breaker release its lock before calling the wrapped function?
11. Given an unfamiliar repo and 45 minutes to add a feature, what do you do in the first 5 minutes?
