# 08. Failure, Judgment, and Prioritization

These questions test the one thing a Staff engineer is paid for: **making good calls when the answer isn't obvious, and owning the bad ones**. "Tell me about a failure" is almost guaranteed. "How do you prioritize," "a decision with incomplete data," and "a tradeoff you made" show up in nearly every leadership round, often in disguise.

The bar is high on honesty. Interviewers have heard every humblebrag. What earns the Staff signal is a real mistake, owned in the first person, fixed, and followed by a **changed mechanism**, not a vague lesson.

---

## At a glance

| Question family | What they're really testing | Best stories |
|---|---|---|
| Failure / biggest mistake | Ownership, self-awareness, real learning | #11 `[pick]`, #2 |
| Outage or incident you owned | Calm under pressure, blameless rigor, prevention | #11, #6 |
| Wrong technical decision / design you'd redo | Design judgment and honesty about your own calls | #2, #11 |
| Decision with incomplete data | Comfort with ambiguity; reversible vs. irreversible thinking | #1, #2, #7 |
| Prioritization with competing demands | Business judgment; saying what you *won't* do | #9, office hours, upstream vs. roadmap |
| Speed vs. quality / tech debt | Pragmatism with guardrails | #7, #2, #8 |
| Calculated risk | Risk sizing, mitigation, reversibility | #7, #1 |
| Killing or descoping a project | Sunk-cost discipline | `[supply]` |
| Ethical / judgment calls | Integrity, security and compliance instincts | #10, confidentiality |
| "What would you do differently?" | Reflection on any story | Every story's Learning beat |

---

## What Staff looks like here

| Senior answer | Staff answer |
|---|---|
| Failure was small and someone else's fault | Failure was consequential, and "I made the call" |
| "I learned to communicate more" | "I added a safeguard / a review gate / a test, and it's caught two issues since" |
| Fixed the bug | Fixed the bug, the class of bug, and the process that let it ship |
| Prioritizes by who asked loudest | Prioritizes by impact, risk, and reversibility, and publishes the list of what's *not* happening |
| Decides only with full data, or decides on gut | Names what's unknown, sizes the cost of being wrong, picks a reversible step, sets a checkpoint |
| Tech debt is bad | Tech debt is a loan: take it deliberately, write down the interest, schedule repayment |

**Anti-signals**
- Humblebrags: "I care too much," "I took on too much," "I was too thorough."
- A failure with no consequence, or one that was really someone else's.
- Blaming tools, tenants, or leadership.
- "We failed" with no "I."
- A learning that's a platitude ("communication is key").
- Prioritization with no tradeoff: everything got done.
- Defensiveness when probed.

---

## The Learning beat: make it a mechanism

The "L" in STAR-L (`03-storytelling-framework.md`) is where this module's answers live or die. A Staff learning is **a change to how the system or the team works**, ideally one you can point to.

| Weak learning | Strong learning (mechanism) |
|---|---|
| "I learned to communicate better." | "I now send a written rollout plan with a rollback step to affected tenants before any control-plane change." |
| "I learned to test more." | "We added `[this scenario]` to the scalability harness, so every release runs it." |
| "I learned to be more careful." | "We added a two-phase deletion safeguard, so deleting a control plane requires an explicit second step." |
| "I learned to involve stakeholders earlier." | "Office hours now include an 'upcoming changes' section, so tenants see changes a month ahead." |
| "I learned to think about failure modes." | "Every design doc I review now has a 'what happens if a data center is lost' section." |

**Test:** could someone check whether you actually do it? If yes, it's a mechanism.

**Formula:**
> "What I changed was **[mechanism]**. Since then, **[evidence it works]**. The broader lesson is **[principle]**."

---

## Cluster 1: Failure or biggest mistake

**Common phrasings**
- "Tell me about a time you failed."
- "What's the biggest mistake you've made professionally?"
- "Tell me about a project that didn't go as planned."
- "Tell me about something you'd do differently."
- "Describe a time you let your team or a customer down."
- "What's a professional regret?"
- "Tell me about a time you missed a deadline or commitment."

**What they're really asking:** Do you own mistakes without flinching, and does your learning change what you do?

**How to pick your failure (Story #11)**

| Criterion | Ask yourself |
|---|---|
| **Real** | Would a colleague recognize it as a mistake? |
| **Consequential** | Did it cost something: tenant impact, time, trust, rework? |
| **Your fault** | Was it your decision or your miss, not the team's or a vendor's? |
| **Fixed** | Did you fix the immediate damage? |
| **Changed behavior** | Is there a mechanism you still use? |
| **Safe** | Not catastrophic, not confidential, not an integrity issue |

**Candidates from your work** (confirm and pick one):

| Candidate | Why it could work | Confirm |
|---|---|---|
| **Member-cluster state flapping race condition** | Real production issue; you added a safeguard (office hours, 05/2026) | `[What was your part? Did you design the code path, miss it in review, or not anticipate it? Tenant impact?]` |
| **Standalone topology's single point of failure** | Your design had a known SPOF; you evolved to stretched | `[Did you underestimate it, or was it a deliberate v1? If deliberate, it's a tradeoff story, not a failure.]` |
| **DR verification built in later than ideal** | Common, honest, fits "DR you never test is DR you don't have" | `[Is this true? When did automated DR testing land on the roadmap?]` |
| **Detector controller concurrent reconciliation → SLO violations** | Upstream issue you filed and tracked | `[Was this a failure on your side, or a bug you found? Only use if you own part of it.]` |
| **A missed estimate or commitment** | Universal, low-risk | `[supply]` |

**Use:** Story #11. Backup: Story #2 told as "the v1 I'd design differently."

**Answer skeleton (~2.5 min)**
1. **Context (20 s):** what you were doing and why it mattered.
2. **The mistake (25 s):** in one sentence, first person. "I decided / I missed / I assumed."
3. **Impact (15 s):** honest, sized.
4. **What you did (45 s):** immediate fix, then root cause, then prevention.
5. **Mechanism and evidence (30 s).**
6. **Principle (10 s).**

<details><summary>Sample answer (template using the state-flapping race condition; confirm every bracket)</summary>

> On our managed multi-cluster platform, member clusters register with Karmada control planes automatically as cluster inventory changes. `[Context: we'd shipped/changed X in the registration path.]`
>
> The mistake was mine: `[e.g., I designed the registration flow assuming updates to a member cluster's state would be serialized, and I didn't account for two components updating it concurrently]`. Under `[condition]`, that created a race that caused member cluster states to flap between ready and not ready.
>
> The impact: `[e.g., flapping could trigger unnecessary taint-based eviction / noisy alerts for N tenants over X days]`. `[Size honestly.]`
>
> I `[stabilized it immediately by Y]`, then dug into the root cause with `[who]`. We added a safeguard to prevent the race, and `[added a regression test / added it to the scalability harness]`. Because taint-based eviction was enabled by default around then, I also `[checked that flapping couldn't trigger eviction storms]`.
>
> What I changed was `[mechanism: e.g., any controller that writes cluster state now has to go through a single owner, and I added concurrency scenarios to our test suite]`. Since then, `[evidence]`. The broader lesson: in a reconciliation system, "who owns this field" is a design question, not an implementation detail.

</details>

**Likely follow-ups**
- *"Why didn't you catch it earlier?"* Honest answer about the gap (test coverage, assumptions), then the mechanism that closes it.
- *"How did tenants react?"* How you communicated: proactively, with a timeline and fix.
- *"What did your manager say?"* Keep it short; focus on what you did.
- *"What would you do differently on day one?"* Name the design choice.
- *"Has it happened again?"* Evidence of the mechanism working.
- *"Tell me about a different failure."* Have a second one ready (often a non-technical one: a missed commitment, a communication miss).

**Traps**
- Picking a failure so small it's a non-event.
- Blaming upstream, a tool, or a tenant.
- Over-explaining the technical detail and running out of time for the learning.

---

## Cluster 2: An outage or incident you owned

**Common phrasings**
- "Tell me about a production incident you handled."
- "Walk me through the worst outage you've been part of."
- "Tell me about a time you were on call and something broke."
- "How do you run a postmortem?"
- "Tell me about a time you prevented a class of incidents."

**What they're really asking:** Do you stay calm, communicate well, find root cause, and turn incidents into prevention, without blame?

**Use:** Story #11 (if incident-shaped). Backup: Story #6 (SLO-driven ops, symptom-based alerting) as the prevention story; the two-phase deletion safety mechanism for control planes as a preventive guardrail `[confirm what triggered it: a near-miss, an incident, or a risk review?]`.

**Incident answer structure (say the phases)**

| Phase | What to say |
|---|---|
| **Detect** | How you found out (alert, tenant report) and how fast |
| **Mitigate** | Stop the bleeding first: rollback, failover, disable |
| **Communicate** | Who you told, how often, what channel |
| **Root cause** | The technical cause *and* the reason it shipped |
| **Prevent** | The class-level fix and the process change |
| **Blameless review** | How you ran or contributed to the postmortem |

<details><summary>Sample answer (template using the two-phase deletion safeguard; confirm the trigger)</summary>

> Our tenants declare control planes as custom resources against a source-of-truth API, and operators reconcile the rest. That's great for self-service, but it means deleting one resource can tear down a stretched control plane and its etcd.
>
> `[Trigger: e.g., during a cleanup, a control plane was deleted that a tenant still depended on / we had a near-miss where a deletion almost went through / a risk review flagged it.]` `[Impact, honestly sized.]`
>
> `[If incident: We restored by X, communicated to the tenant within Y.]` In the review, the question wasn't who deleted it. It was why one step could cause that much damage. So we added a two-phase deletion safety mechanism: deleting a control plane takes an explicit second step, so a single mistaken change can't take it out.
>
> The broader change: for any irreversible operation on the platform, we now ask whether it needs a two-step guard. The lesson is that self-service raises the blast radius of every mistake, so guardrails have to scale with it.

</details>

<details><summary>Sample answer: preventing a class of incidents (Story #6)</summary>

> Early on, our alerting was infrastructure-centric and noisy: lots of pages about component health that didn't mean tenants were affected, and tenants often found real problems before we did.
>
> I led the move to SLOs scoped to tenant experience. We found that existing guidance assumed single clusters and broke down at cluster boundaries, so we built global metrics aggregation with Thanos and symptom-based alerting across clusters, then calibrated thresholds and alerted on SLO burn rate. I also wrote it up as a reliability guide, internally and upstream in Karmada.
>
> Result: alerts mean tenant impact, and tenants can fix their own issues with the dashboards and guides. The lesson: alert on what users feel, and make the signals good enough that tenants can self-serve.

</details>

**Likely follow-ups**
- *"How do you keep postmortems blameless?"* Ask "why did the system allow it," not "who did it"; name actions as system changes.
- *"How did you communicate during the incident?"* Cadence, channel, audience.
- *"How do you make sure action items actually get done?"* Owners, dates, tracked in the roadmap; review at office hours.
- *"What's your on-call philosophy?"* Page only on user impact; every page actionable; toil reduction is roadmap work.
- *"What metrics did you track?"* SLO burn rate; `[on-call metric if you can share]`.

**Traps:** heroics as the story (staying up all night) instead of the prevention; blame.

---

## Cluster 3: A wrong technical decision or design you'd redo

**Common phrasings**
- "Tell me about a technical decision you regret."
- "What's a design you'd do differently?"
- "Tell me about a time your architecture didn't hold up."
- "What would you change about the system you built?"
- "Tell me about tech debt you created."

**What they're really asking:** Can you critique your own designs with the same rigor as others'?

**Use:** Story #2 (standalone → stretched → source of truth + operators). Be precise about whether standalone was a *mistake* or a *deliberate v1*. Both are fine, but they're different stories.

**Two framings**

| If standalone was deliberate | If it was underestimated |
|---|---|
| "v1 had a known SPOF. I accepted it to prove the model, with a plan to stretch it. What I'd change: I'd have designed the control-plane spec so stretching didn't require `[migration work]`." | "I underestimated how unacceptable the management-cluster SPOF would be for a DR mandate. We had to re-architect to stretched, which cost `[time]`." |

<details><summary>Sample answer (Story #2; choose the framing that's true)</summary>

> Our first topology was standalone: one management cluster ran the Karmada operator and our managed operator, hosting tenants' control planes. It worked, but the management cluster was a single point of failure, which was a contradiction for a platform whose whole purpose was DR.
>
> `[Framing.]` We moved to a stretched topology: each control plane runs on management clusters in multiple data centers behind one DNS-based endpoint, with a shared API server CA via cert-manager and one multi-data-center etcd cluster. The tradeoff was cross-data-center etcd latency for surviving the loss of a site. For control-plane writes, that's the right trade. Then we added a source-of-truth control plane and operators so tenants declare intent and everything else is reconciled.
>
> What I'd do differently: `[e.g., design for stretched from day one, since the DR mandate made it inevitable / build automated DR verification into the first version instead of adding it later]`. The broader lesson: when the requirement is a guarantee, like surviving a site loss, the v1 has to honor the guarantee even if it's missing features.

</details>

**Likely follow-ups**
- *"Why didn't you start with stretched?"* Honest reason: time, unknowns about multi-DC etcd, proving the model.
- *"What did the migration cost?"* `[supply]`
- *"What's the current design's weakness?"* Have one: e.g., what happens if both management clusters fail; dependency on the managed etcd service.
- *"What's next for it?"* Conformance testing, automated DR testing, Argo-based test automation.

**Traps:** saying you'd change nothing; criticizing a design that was someone else's.

---

## Cluster 4: Decisions with incomplete data or under uncertainty

**Common phrasings**
- "Tell me about a time you made a decision without all the information."
- "How do you make decisions under ambiguity?"
- "Tell me about a time you had to move fast with uncertainty."
- "Describe a big bet you made."
- "How do you know when you have enough data?"

**What they're really asking:** Can you separate reversible from irreversible decisions, size the cost of being wrong, and move?

**Your decision framework (say it):**
> "I ask three things: how reversible is it, what does being wrong cost, and what's the cheapest way to reduce the biggest unknown. Reversible and cheap: decide now. Irreversible or expensive: buy information first with a spike, a load test, or a conversation, then set a checkpoint to revisit."

| Decision type | Your move | Example from your work |
|---|---|---|
| Reversible, low cost | Decide quickly, revisit | Feature gate defaults; alert thresholds (later calibrated) |
| Irreversible, high cost | Buy information first | Choosing Karmada as the firm standard; multi-DC etcd topology |
| High uncertainty about scale | Measure | Scalability harness (CL2 + KWOK) before committing |

**Use:** Story #1 (choosing Karmada as the standard). Backup: Story #8 (building the scalability harness to replace guesses with data), Story #7.

<details><summary>Sample answer (Story #1; fill brackets)</summary>

> Teams across Bloomberg were each solving the DR mandate on their own, which was expensive and inconsistent. I proposed a shared multi-cluster platform, and the big decision was the orchestrator. `[What was unknown: maturity at our scale, whether it could run stretched across data centers, community health.]`
>
> It wasn't easily reversible: tenants would build on it. So I bought information first: `[a proof of concept, evaluation against alternatives, conversations with maintainers]`. I couldn't eliminate every unknown, so I picked mitigations: we'd contribute upstream rather than fork, so we weren't dependent on others for fixes, and we'd measure scale ourselves (that later became the scalability harness).
>
> It went from proof of concept to production, and now backs 40+ cluster fleets. The bet on upstream paid off: the features we needed landed, I became a maintainer, and Karmada graduated in the CNCF. The lesson: for irreversible bets, reduce the risk you can and own the rest with a mitigation, not more analysis.

</details>

**Likely follow-ups**
- *"What alternatives did you consider?"* `[supply: other approaches, build in-house, active/passive]`
- *"What would have made you reverse?"* Name a concrete kill criterion.
- *"How did you get buy-in with uncertainty?"* Transparent about risks, mitigations, checkpoints. See `05-leadership-direction-and-influence.md`.

**Traps:** claiming you had certainty; analysis paralysis framed as rigor.

---

## Cluster 5: Prioritization with competing demands

**Common phrasings**
- "How do you prioritize?"
- "Tell me about a time you had more work than you could do."
- "How do you balance tenant requests, roadmap, operations, and tech debt?"
- "Tell me about a time you had to say what you *wouldn't* do."
- "How do you decide what to work on as a Staff engineer?"
- "Tell me about conflicting priorities from different stakeholders."

**What they're really asking:** Business judgment, and the discipline to drop things visibly.

**Your demand types**

| Demand | Example |
|---|---|
| Tenant requests | Feature gates, audit config, onboarding new teams (risk platform, data licence platform) |
| Platform roadmap | Version support (Karmada v1.17, Kubernetes v1.35), Kyverno support, structured logging |
| Operations / reliability | SLO calibration, safeguards, KaaS cluster alerts |
| Upstream | Conformance testing, scalability testing, operator features |
| Org work | ARB reviews, CNCF case study, mentoring |

**Your prioritization approach (say it):**
> "I rank by impact on tenants and risk to the platform's guarantees, then by leverage: does it help many tenants or one, does it remove toil permanently. Reliability and security issues jump the queue. Then I make the list public. At our monthly office hours, we show finished, ongoing, and upcoming work, so tenants see what's *not* happening this month and why."

**Use:** Office hours (Story #9) as the mechanism. Example of leverage: making feature gates and API server audit config tenant-configurable instead of handling each request by hand.

<details><summary>Sample answer</summary>

> As tech lead and product owner, I'm balancing four streams: tenant requests, our roadmap (version upgrades, new integrations), reliability work, and upstream Karmada contributions.
>
> My rule is: anything that threatens the platform's guarantees, like DR, security, or SLOs, goes first. Then I rank by leverage. A good example: several tenants needed different feature gates and audit settings. Handling each request by hand would have been endless, so we made both tenant-configurable. One piece of work closed a whole category of requests.
>
> Upstream work gets ranked the same way: we prioritized contributions that remove internal patches or unblock us, because they lower our upgrade cost. `[e.g., fixing Helm upgrades of the Karmada operator upstream.]`
>
> The important part is saying no visibly. Every month at tenant office hours, we publish finished, ongoing, and upcoming work. If something isn't on it, tenants can see that and argue for it in the room. `[Example of a request you deferred and how you communicated it.]`
>
> The lesson: prioritization only builds trust if it's visible. A private list just looks like you ignored people.

</details>

**Likely follow-ups**
- *"Tell me about something you deprioritized that came back to bite you."* Honest example `[supply]`.
- *"How do you handle a VP's pet request?"* Make the tradeoff visible: "If we do this, X slips. Is that the right trade?"
- *"How do you protect time for tech debt?"* A standing share of capacity, or tie debt to an SLO or upgrade cost.
- *"How do you prioritize your own time as Staff?"* Highest leverage: unblocking others, direction, reviews, then your own code.

**Traps:** "I work harder"; no example of something dropped.

---

## Cluster 6: Speed vs. quality, and tech debt

**Common phrasings**
- "Tell me about a time you had to cut corners."
- "How do you balance shipping fast with quality?"
- "Tell me about tech debt you chose to take on."
- "Tell me about a time you pushed back on shipping something."
- "How do you decide when something is good enough?"

**What they're really asking:** Pragmatism with guardrails. Do you know which corners are safe to cut?

**Your rule (say it):**
> "I'll cut scope, polish, and automation. I won't cut the things that protect the guarantee: DR behavior, security, rollback, observability. Debt I take on gets written down with the interest and a repayment trigger."

| Safe to cut | Not safe to cut |
|---|---|
| Manual steps you'll automate later | Rollback path |
| Feature breadth | Security and credential handling |
| UI and docs polish (for internal v1) | Observability for the new path |
| Generality | DR guarantee |

**Use:** Story #7 (zero-downtime ingress-nginx → Istio migration: what you didn't cut). Backup: Story #2 (standalone as deliberate debt), Story #8.

<details><summary>Sample answer (Story #7; fill brackets)</summary>

> We migrated ingress for the managed control planes from ingress-nginx to Istio, with zero downtime. `[Why: ingress-nginx retirement / platform standard / feature needs.]` There was pressure to move fast: `[the deadline or reason]`.
>
> I was willing to cut `[e.g., automation for later phases, migrating everything in one tool]`. I wasn't willing to cut a per-control-plane rollback and tenant-facing validation, because control-plane ingress is the front door for every tenant. `[Approach: e.g., ran both in parallel, shifted traffic gradually, validated per tenant, kept the old path until confirmed.]`
>
> `[Result: zero downtime across N control planes.]` The debt we took on was `[X]`, which we wrote down and `[paid back by Y]`. The lesson: speed comes from cutting the right things, not from skipping the safety net.

</details>

**Likely follow-ups**
- *"Did you ever cut a corner that hurt?"* Pairs with Cluster 1.
- *"How do you convince leadership to fund debt paydown?"* Tie it to upgrade cost, on-call load, or SLO risk with numbers.
- *"What's 'good enough' for an internal platform?"* The tenant can rely on it and you can operate it.

**Traps:** "I never cut corners" (unrealistic); cutting the safety net and getting lucky.

---

## Cluster 7: A calculated risk

**Common phrasings**
- "Tell me about a time you took a risk."
- "Tell me about a bold decision."
- "Describe a time you did something outside your comfort zone at work."
- "Tell me about a risky change to production."

**What they're really asking:** Do you size risk, mitigate it, and keep it reversible?

**Use:** Story #7 (risky migration on live infra). Backup: Story #1 (betting the firm standard on Karmada), Story #4 (Bloomberg going public as the first North American adopter in the graduation process `[confirm risk angle: e.g., approvals through CTO office and OSPO]`).

**Skeleton (~2 min):** the risk and the upside (25 s) → how you sized it: blast radius, likelihood, reversibility (25 s) → mitigations (40 s) → outcome (20 s) → what you'd do differently (10 s).

<details><summary>Sample answer (short form, Story #7)</summary>

> Changing ingress on live control planes risked cutting tenants off from their control plane, the worst-case outcome for us. I sized it as: blast radius is every tenant behind that ingress, likelihood moderate if we did it all at once, low if staged. So we staged it, kept the old path live, validated per tenant, and had a rollback for every step. It went through with zero downtime. What I'd do differently: `[supply]`.

</details>

**Follow-ups:** *"What was your rollback plan?"*; *"What would have made you abort?"* (name a concrete abort criterion); *"Did anyone think it was too risky?"* (links to `07-conflict-disagreement-and-pushback.md`).

---

## Cluster 8: Killing or descoping a project

**Common phrasings**
- "Tell me about a time you stopped a project."
- "Tell me about a time you cut scope to hit a date."
- "How do you know when to give up on an approach?"
- "Tell me about a time you recommended not building something."

**What they're really asking:** Sunk-cost discipline and the courage to say "stop."

**Use:** `[supply: an approach you abandoned, a feature you descoped, a tenant ask you declined to build]`. Candidates: the decision to compose existing internal services instead of building new ones (recommending *not* building), declining internal forks (Story #9), descoping for GA `[supply]`.

<details><summary>Sample answer (recommending not building, from the platform's design)</summary>

> When we designed the platform, there was a natural pull to build our own pieces: our own etcd management, our own service discovery, our own identity layer. I recommended against it. Bloomberg already had a managed etcd service, DNS-based service discovery, Kubernetes-as-a-service, and identity and secret management. Composing those meant we didn't own operating them, and we inherited their maturity.
>
> The cost was dependency on other teams' roadmaps, which we managed with clear ownership and early asks. The lesson: the cheapest code to operate is code someone else already operates well.

</details>

**Follow-ups:** *"How did the team feel?"*; *"What did you salvage?"*; *"What signal told you to stop?"*

**Traps:** a project killed by someone else that you just went along with.

---

## Cluster 9: Ethical and judgment calls

**Common phrasings**
- "Tell me about a time you had to make an ethical decision."
- "Tell me about a time you noticed a security or compliance problem nobody owned."
- "Tell me about a time you were asked to do something you weren't comfortable with."
- "How do you handle confidential information?"

**What they're really asking:** Integrity and good instincts on security, compliance, and confidentiality.

**Use:** Story #10 (Hoppscotch: a policy-compliant replacement for Postman, adopted by 375+ engineers). Also: how you handle confidentiality in this very interview (public framing only), which is itself a Staff judgment signal.

<details><summary>Sample answer (Story #10)</summary>

> `[Context: Postman usage conflicted with firm policy, e.g., around where request data and credentials could be stored.]` Engineers relied on it daily, and nobody owned a replacement.
>
> The easy path was to ignore it or just ban the tool. Neither was good: one left a compliance gap, the other blocked engineers. I took it on: `[evaluated options, set up Hoppscotch as a policy-compliant alternative]`. It's now used by 375+ engineers.
>
> The lesson: when a compliance problem has no owner, the fix that sticks gives people something better, not just a rule.

</details>

**On confidentiality (if asked, or when declining details):**
> "I can't share internal specifics, but the shape of the problem was…"

Interviewers read this as judgment, not evasion (see `interview-playbook.md` §1.6).

**Follow-ups:** *"What if your manager had told you to drop it?"*; *"What would you do if asked to ship something insecure for a deadline?"* (say no to that part, offer a safe scoped alternative, escalate if pressed).

---

## Cluster 10: "What would you do differently?"

This is asked after almost every story. Prepare one sentence for **every** story in `04-story-bank.md`.

| Story | Your "differently" (draft; confirm) |
|---|---|
| #1 Karmada adoption | `[e.g., involve dependency teams earlier]` |
| #2 Architecture evolution | `[e.g., stretched from day one / DR verification from v1]` |
| #3 Upstreaming | `[e.g., open the design discussion before writing code]` |
| #4 CNCF graduation | `[supply]` |
| #5 Growing contributors | `[supply]` |
| #6 SLO-driven ops | `[e.g., involve tenants in SLO definitions earlier]` |
| #7 Istio migration | `[supply]` |
| #8 Scalability harness | `[e.g., build it before the first scale question, not after]` |
| #9 Saying no | `[e.g., publish the shared-responsibility model at onboarding]` |
| #10 Hoppscotch | `[supply]` |
| #11 Failure | the mechanism from Cluster 1 |
| #12 Conflict | from `07-conflict-disagreement-and-pushback.md` |

**Rule:** never "nothing." Even a great outcome has a "faster," "earlier," or "with less toil" version.

---

## Rapid-fire variants

| Phrasing | Cluster | Story |
|---|---|---|
| "Tell me about feedback that was hard to hear." | 1 | `[supply]` |
| "Tell me about a time you underestimated something." | 1, 3 | #2, #11 |
| "Describe a time you didn't meet expectations." | 1 | #11 |
| "What's a hard tradeoff you made recently?" | 4, 6 | #2, #7 |
| "How do you decide build vs. buy?" | 8 | platform composition |
| "How do you handle being on call for something you didn't build?" | 2 | #6 |
| "Tell me about a time you had to choose between two good options." | 4 | #1, #2 |
| "Tell me about a time you were overloaded." | 5 | office hours, prioritization |
| "Tell me about a decision you made that was unpopular." | 5, 7 | #9 |
| "Tell me about a time a risk didn't pay off." | 7, 1 | `[supply]` |
| "How do you balance open source work with company priorities?" | 5 | #3, upstream vs. roadmap |
| "How do you know a project is in trouble?" | 8 | `[supply]` |
| "What's the most important thing you've learned in your career?" | Learning beat | pick your best mechanism |
| "Tell me about a near-miss." | 2 | two-phase deletion `[confirm]` |

---

## Self-check

1. Does your failure story start with "I" in the mistake sentence?
2. Is the impact sized honestly, with a number or a concrete consequence?
3. Is your learning a mechanism someone could verify?
4. Can you state your decision framework (reversibility, cost of being wrong, cheapest information) in one breath?
5. Can you name something you deliberately did *not* do, and how you communicated it?
6. For Story #2, have you decided whether standalone was a deliberate v1 or an underestimate?
7. Do you have a "what I'd do differently" line for every story?
8. Do you have a second failure ready, in case they ask for another?

## Drill (40 min)

1. **Pick Story #11** from the candidates table. Confirm every bracket with your own memory (10 min).
2. **Write the mistake sentence** in one line, first person, no softeners.
3. **Rewrite your learning** using the formula until it passes the "could someone check?" test.
4. **Say the failure story aloud** in under 2.5 minutes, recorded. Then answer: "Why didn't you catch it earlier?" and "Tell me about a different failure."
5. **Prioritization dry run:** take last month's office-hours list and explain, out loud in 90 seconds, why upcoming work was ordered the way it was and what you deferred.
6. **Fill the Cluster 10 table** for all stories.
