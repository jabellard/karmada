# 02. Staff Signals: What Interviewers Listen For

Behavioral rounds decide **level** more than any other part of the loop. Coding mostly answers "can this person build?" and system design answers "can they design?" The behavioral rounds answer "at what **scope** do they operate?" Two candidates can tell the same story and get leveled differently because of *how* they tell it. This module teaches you what Staff sounds like so every answer sends it.

---

## At a glance

| Dimension | The interviewer's question | Senior sounds like | Staff sounds like |
|---|---|---|---|
| **Scope** | How big is the problem they own? | "My team's service" | "The firm's approach to multi-cluster DR, across 10+ teams" |
| **Ambiguity** | Do they need a defined problem? | "I was asked to build X" | "Nobody owned this; I framed the problem and got it funded" |
| **Direction** | Do they decide *what* to build, or just *how*? | Executes the roadmap | Shapes the roadmap; says what *not* to build |
| **Influence** | Can they move people they don't manage? | Convinces their own team | Moves other teams, leadership, and upstream maintainers |
| **Judgment** | Do they make the right call with tradeoffs? | Picks a good solution | Names alternatives, picks one for stated reasons, revisits when facts change |
| **Multiplier** | Do others get better because of them? | Mentors when asked | Builds mechanisms (standards, docs, paved paths) and grows people into owners |
| **Ownership** | Do they own outcomes, not tasks? | "I shipped it" | "I owned whether it worked for tenants: SLOs, on-call, adoption" |
| **Business sense** | Do they connect work to why it matters? | Technical outcome | Outcome for the org: DR mandate met, velocity, risk, cost |
| **Self-awareness** | Do they learn? | Successes only | Clear about mistakes and what they changed because of them |

**The Staff test in one line** (same as `../interview-playbook.md` §0): *Would this person make the right call on an ambiguous, cross-team problem without supervision, and make the people around them better?*

---

## 1. What "Staff" means to the people judging you

Companies differ, but most use some version of these **Staff archetypes** (popularized by Will Larson's *Staff Engineer*). Knowing which one a role wants lets you weight your stories.

| Archetype | What they do | Your evidence | Lead with these stories |
|---|---|---|---|
| **Tech lead** | Guides the approach and execution of a team or area | TL and product owner of the managed multi-cluster platform | #1, #2, #7, #9 |
| **Architect** | Owns direction and quality across a critical domain | Set the firm-wide multi-cluster DR standard; Architecture Review Board | #1, #2, #13 |
| **Solver** | Goes deep on the hardest problems, wherever they are | Stretched control planes on multi-DC etcd; upstream etcd/CA/token work | #2, #3, #8 |
| **Right hand** | Extends a leader's reach across an org | CNCF graduation coordination; adopter interviews across 3 teams | #4, #1 |

**Your profile reads as Tech lead + Architect.** That's the most common Staff shape at infra companies. Don't claim Solver depth on GPU internals you haven't touched (see `11-curveballs-hypotheticals-and-recovery.md`).

### What changes at Staff (and what doesn't)

- **Doesn't change:** you still need to be technically credible. Behavioral answers that are all process and no substance read as "manager in IC clothing."
- **Changes:** the *unit* of your work. Senior engineers deliver projects. Staff engineers deliver **outcomes across boundaries**: teams, orgs, codebases, companies.
- **Changes:** where your time goes. Interviewers listen for writing, aligning, reviewing, deciding, and growing others, *alongside* building.

---

## 2. The seven Staff signals, with how to send each

### Signal 1: You chose the problem

Staff engineers find and frame problems. Senior engineers are handed them.

| ❌ Senior framing | ✅ Staff framing |
|---|---|
| "My manager asked me to evaluate multi-cluster tools." | "Teams were each solving the DR mandate on their own, which was expensive and inconsistent. I made the case that the platform should solve it once, and I led the evaluation that landed on Karmada." |

**Phrases that send it:** "I noticed…", "Nobody owned…", "I made the case that…", "I reframed the problem as…"

**Only say this if it's true.** If the problem was handed to you, the Staff move is how you *reshaped* it: "I was asked for X. I pushed to make it Y because…"

### Signal 2: You decided, with alternatives

The single strongest judgment signal is **naming the option you rejected and why**.

> "We could have gone active/passive, which is simpler to reason about. I chose stretched control planes on one multi-data-center etcd cluster because tenants needed one endpoint that survives a site loss without a failover runbook. The cost is cross-DC write latency, which is fine for control-plane traffic."

**Formula:** *Options → criteria → choice → cost you accepted → how you'd know if you were wrong.*

### Signal 3: You moved people you don't manage

Count the boundaries you crossed. Each one is a level signal.

| Boundary | Your evidence |
|---|---|
| Other teams inside the company | 10+ core platform and infra teams, 50+ engineers, adopting the platform |
| Leadership | Firm-wide standard for the DR mandate; partnered with the CTO office and OSPO |
| Outside the company | Upstreamed external etcd, custom CA, and token rotation into Karmada; CNCF graduation |
| Peers' designs | Architecture Review Board |

**How to tell it:** influence stories need a **skeptic**. "Everyone agreed" isn't influence. Name who was unconvinced, what they cared about, and what changed their mind (data, prototype, a tradeoff you gave up). See `05-leadership-direction-and-influence.md`.

### Signal 4: You multiplied others

Staff framing rule (from `../linkedin-guide.md` §4.6): **describe the trajectory you created, not the act of mentoring.**

| ❌ Activity | ✅ Trajectory |
|---|---|
| "I mentor engineers and do code reviews." | "I grew Bloomberg's Karmada contributors from fewer than 5 to 20+. Mentees went from first PRs to driving design proposals and scheduler enhancements." |
| "I wrote documentation." | "I wrote the reliability guide that defines our SLOs and burn-rate alerting, then upstreamed a version into Karmada, so other adopters start from a baseline." |

Mechanisms beat heroics: office hours, paved paths, guides, standards, test harnesses (#8), shared-responsibility lines (#9). See `06-growing-and-developing-others.md`.

### Signal 5: You owned the outcome, not the task

Ownership means you cared whether it **worked for the people using it**, after launch.

> "Shipping the control plane wasn't the goal. Tenants inheriting DR was. So I owned the SLOs, scoped them to tenant experience instead of infrastructure health, and changed alerting so tenants could fix their own issues."

**Tell-tale phrases:** "after launch…", "the metric I watched was…", "on-call went from… to…", "adoption", "tenants".

### Signal 6: You connected it to the business

Every major story should have a "so what" a director would care about.

| Technical result | Business translation |
|---|---|
| Stretched control planes | Teams meet the firm's DR mandate by adopting the platform, not by building it |
| Weeks → minutes control plane setup | Tenant platform teams ship faster; no specialist time per onboarding |
| SLO-driven alerting | Less on-call load; incidents caught by symptom, not by noise |
| Hoppscotch | 375+ engineers on a policy-compliant tool; compliance risk removed |
| Upstreaming instead of forking | No long-lived internal fork to maintain; Bloomberg benefits from community work |

### Signal 7: You're honest and you learn

Staff candidates are trusted with ambiguity, which requires trusting their self-reports. That means:
- Admit what you got wrong, concretely (`08-failure-judgment-and-prioritization.md`).
- Say "I don't know" about GPU internals, then reason.
- Give the Learning beat in every story: what you'd do differently, or the **mechanism** you changed.

---

## 3. "I" vs "we": the most common leveling mistake

Interviewers can only credit what *you* did. "We" hides you. But claiming everything sounds dishonest.

**Rule:** **"I" for decisions, framing, and influence. "We" for team execution.**

> "**I** proposed the stretched topology and wrote the design. **We** built the operators as a team. **I** owned the etcd and PKI decisions, and **I** worked with the managed etcd team to get multi-DC clusters per control plane."

**Test yourself:** after a story, could a listener list three decisions that were yours? If not, retell it.

**When they ask "what was your specific role?"**, answer in one breath: *"I was the tech lead and product owner: I set the architecture, made the build-vs-compose calls, and owned tenant adoption. The team built most of the operators; I wrote [the parts you wrote]."*

---

## 4. Calibrating scope words

Vague size words read as small. Specific scope reads as Staff. Memorize the numbers sheet (`../interview-playbook.md` §1.5).

| Instead of | Say |
|---|---|
| "a lot of clusters" | "40+ Kubernetes cluster fleets" |
| "much faster" | "from weeks to minutes, self-service" |
| "several teams" | "10+ core platform and infrastructure teams, 50+ engineers" |
| "grew the community" | "fewer than 5 to 20+ contributors; a Contribufest with 20+ engineers" |
| "widely used" | "375+ engineers" |
| "important project" | "the firm's standard for its DR and HA mandate" |

**If you can't share a number** (confidentiality), give the shape: "I can't share the internal figure, but it was most of the tenants running stateful services." That itself signals judgment.

---

## 5. Anti-signals: what gets you leveled down

| Anti-signal | What it sounds like | What the interviewer writes |
|---|---|---|
| Waiting for direction | "Once the requirements were finalized, I…" | "Needs well-defined problems" |
| Hero coding | "I stayed up all weekend and fixed it myself" | "Doesn't scale through others" |
| Hiding behind "we" | "We decided… we built… we launched…" | "Unclear individual contribution" |
| No tension | "Everyone agreed it was a great idea" | "Story lacks substance" or "low-stakes examples" |
| Villains | "The other team was just being difficult" | "Low empathy; may be hard to work with" |
| Tool tunnel vision | "Karmada is the answer" to everything | "Over-indexes on familiar tech" |
| Humblebrag failure | "My weakness is I care too much" | "Low self-awareness" |
| No business link | Pure technical narration | "Doesn't connect work to impact" |
| Monologue | 6-minute answer, no check-in | "Unstructured communicator" |
| Bluffing | Confident claims about NCCL or InfiniBand you haven't used | "Credibility concern" (this one can sink a loop) |
| Badmouthing | "Bloomberg is too slow" | "Possible culture risk" |

---

## 6. Phrase bank: language that signals scope

Use these naturally. Don't stack them all in one answer.

**Framing and direction**
- "The real problem wasn't X, it was Y."
- "I wrote a one-pager that laid out three options…"
- "I set the principle that… and that resolved most of the later debates."
- "We explicitly decided not to…"

**Tradeoffs**
- "That bought us X at the cost of Y. Given Z, it was the right trade."
- "I'd make the opposite call if [condition changed]."
- "The reversible decisions I made fast. The irreversible one, the etcd topology, I took slowly."

**Influence**
- "Before the review, I met the skeptics one-on-one to understand their concerns."
- "I built a small prototype so the debate was about data, not opinions."
- "I found what they cared about, which was [their incentive], and showed how this helped it."

**Multiplier**
- "I made it so the team didn't need me for that anymore."
- "I turned it into a guide/standard/harness so every team got it for free."
- "Their next design, they drove themselves; I just reviewed."

**Ownership and learning**
- "In hindsight, I'd have [concrete change]. Now I [mechanism] by default."
- "The number I watched after launch was…"

---

## 7. Sending Staff signals in each round type

| Round | Primary signal to send | How |
|---|---|---|
| Recruiter screen | Your scope maps to Staff | Ask for Staff consideration early; one sentence of evidence (`10-motivation-career-and-about-you.md`) |
| HM / EM screen | You can own *their* problem | Talk about their problems; frame your stories as transferable patterns |
| Behavioral round | All seven signals | Story bank, tight structure, "I" for decisions |
| Project deep dive | Judgment + ownership | Alternatives rejected, what you'd redo, what you personally built |
| Bar raiser | Org-level scope and bar | Raise the altitude: "for the org…", "for every tenant…" |
| Director / VP | Strategy and business sense | Talk direction and bets, not implementation |
| Team match | Mutual fit at Staff scope | Ask what you'd own end to end and how Staff engineers shape the roadmap |

Round mechanics are in `01-rounds-and-how-youre-evaluated.md`.

---

## 8. Self-check

1. For each of your top 6 stories, can you name the **alternative you rejected**?
2. Which stories include a **skeptic** and what changed their mind?
3. Which stories show **trajectory** for someone else, not just activity?
4. Can you say each story's **business result** in one sentence a VP would care about?
5. If an interviewer wrote "unclear individual contribution" about one of your stories, which would it be? Fix it.
6. What's the honest one-sentence boundary of your GPU/AI-infra experience?

**Drill:** take Story #1 and tell it twice, first as a Senior engineer would ("I was asked to…") and then as a Staff engineer. Record both. The difference you hear is this module.
