# 07. Conflict, Disagreement, and Pushback

Every Staff loop has at least one conflict question, and many have three in disguise ("disagreed with your manager," "said no to a customer," "a decision you didn't agree with"). Interviewers aren't checking whether you've had conflict. Everyone has. They're checking **how you behave when smart people want different things**: do you make it about the problem, use data, protect the relationship, and get to a decision that sticks?

This module gives you one framework you can narrate, the question clusters with sample answers, and a worksheet to find your real conflict stories (Story #12 is still a gap; fill it here).

---

## At a glance

| Question family | What they're really testing | Best stories |
|---|---|---|
| Disagreed with a peer on a technical decision | Do you argue with criteria and data, and can you be persuaded? | #12, #2, #13 |
| Disagreed with your manager / leadership | Courage plus judgment: do you push up respectfully, then commit? | #12, #1 |
| Disagree and commit | Can you lose an argument gracefully and still execute fully? | #12 `[need a real one]` |
| Said no to a stakeholder / tenant / customer | Customer focus without being a pushover; protecting the platform | #9, #6 |
| Conflict with another team (ownership, priorities, dependencies) | Cross-team influence when nobody reports to you | #1, #9, #6 |
| Disagreement in open source | Influence without any authority, across companies | #3, #5 |
| Difficult coworker / interpersonal friction | Empathy, maturity, not making it personal | #12 `[supply]` |
| You were wrong and changed your mind | Intellectual honesty, ego in check | #12, #2 |
| Mediating conflict between others | TL-level leadership: decision process, not picking favorites | `[supply]`, #5 |
| Escalation | Judgment about when, how, and with what framing | #12, #9 |

---

## What Staff looks like here

| Senior answer | Staff answer |
|---|---|
| "I convinced them I was right." | "We disagreed on criteria. Once we agreed what mattered, the answer was obvious, and in one area they were right." |
| Conflict is between two people | Conflict is between two reasonable goals (e.g., tenant velocity vs. platform reliability) |
| Resolves it by debating harder | Resolves it with a prototype, a measurement, a written tradeoff, or a time-boxed experiment |
| Escalates when stuck, or never escalates | Escalates deliberately, jointly, with both options written up fairly |
| Wins the argument, loses the relationship | Repairs the relationship and follows up after the decision |
| Ends at "they agreed" | Ends with a mechanism: a standard, a doc, a policy that prevents the same fight next time |

**Anti-signals**
- The other person is a villain or "didn't get it."
- You "won" and that's the whole story.
- No stakes: a fight over tabs vs. spaces, or a disagreement resolved in one meeting with no cost.
- You avoided the conflict, or went around the person.
- You escalated first and talked second.
- You can't name a single thing the other side was right about.
- Venting about Bloomberg, a manager, or a tenant team.

---

## The conflict framework (narrate these beats)

Use this as the spine of your Action section. You don't need all eight in every answer; pick the 4–5 that actually happened and name them out loud.

| # | Beat | What you say in the story |
|---|---|---|
| 1 | **Separate the person from the problem** | "I assumed good intent. They were optimizing for something real." |
| 2 | **Understand their incentive or constraint** | "I met them 1:1 to understand what was driving it: their deadline, their on-call load, their customer." |
| 3 | **Find the shared goal** | "We both wanted tenants to be DR-compliant without extra toil." |
| 4 | **Make the disagreement about criteria and data** | "I proposed we agree on how to decide first: blast radius, operational cost, time to deliver." |
| 5 | **Propose an experiment or prototype** | "Instead of debating, I suggested a time-boxed spike / load test with the scalability harness." |
| 6 | **Decide and record** | "We wrote the decision and the rejected option down, with who decided and why." |
| 7 | **Disagree and commit** | "When it went the other way, I executed it as if it were my idea." |
| 8 | **Repair and follow up** | "I followed up afterward, and we revisited it at the checkpoint we'd agreed on." |

**One-line version** you can say in any answer:
> "I try to move a disagreement from *who's right* to *what are we optimizing for*, then let data or a small experiment decide, write it down, and fully commit to whatever we pick."

---

## Choosing a conflict story that's real but safe

A good conflict story passes all five checks:

| Check | Why | Fails if… |
|---|---|---|
| **Stakes were real** | Shows Staff scope | Nothing would have broken either way |
| **You weren't a jerk** | They're checking you'd be safe to work with | You raised your voice, went around them, or gloated |
| **The other party isn't a villain** | Shows empathy and judgment | "They just didn't understand Kubernetes" |
| **Resolution came from your behavior** | Shows agency | A manager stepped in and decided for you |
| **You can name a learning** | Shows growth | "I'd do it the same way" |

**Avoid:** conflicts about compensation or promotion, anything involving HR, stories where the other person is identifiable and still at Bloomberg in an embarrassing way, and anything that requires confidential details to make sense.

**Have two:** one where you **won through process** (data, prototype) and one where you **lost or were wrong** and committed anyway. Many loops ask both.

---

## Conflict story finder (worksheet)

Story #12 is the biggest gap in your bank. Work through these prompts, write one line for each that triggers a memory, then pick the two strongest and draft them with `03-storytelling-framework.md`.

| Where conflict likely happened | Prompt | Your notes |
|---|---|---|
| **Architecture Review Board** (Story #13) | A design you reviewed where you pushed back, e.g., on a team's DR approach. What did they want, what did you require, how did it end? | `[ ]` |
| **Tenant requests** (Story #9) | A tenant who wanted a fork, a custom feature, or the platform to own something (e.g., their DR tests). How did you say no and keep them? | `[ ]` |
| **Upstream PR reviews / design proposals** (Story #3) | External etcd, custom CA, token rotation: did a reviewer or maintainer push back on scope, API shape, or approach? | `[ ]` |
| **Migration timing** (Story #7) | ingress-nginx → Istio: did anyone want to go faster, slower, or not at all? Who carried the risk? | `[ ]` |
| **SLO thresholds and alerting ownership** (Story #6) | Moving from infra-centric alerts to SLOs scoped to tenant experience: who disagreed on thresholds, paging, or who gets paged? | `[ ]` |
| **Platform vs. tenant scope** (Story #9) | Shared-responsibility lines: did a tenant team or a dependency team disagree about who owns what? | `[ ]` |
| **Karmada adoption** (Story #1) | Did anyone favor a different approach (build in-house, another multi-cluster tool, active/passive)? | `[ ]` |
| **Architecture evolution** (Story #2) | Standalone → stretched: did anyone argue against cross-DC etcd latency, or for active/passive? | `[ ]` |
| **Dependency teams** | Managed etcd, DNS-based discovery, Kubernetes-as-a-service, identity/secrets: a priority or ownership clash with an infrastructure services team? | `[ ]` |
| **Your manager** | A time you disagreed on priority (e.g., upstream work vs. internal roadmap), staffing, or a deadline | `[ ]` |
| **Mentees / your team** | Two engineers who disagreed on a design and you had to help them land it | `[ ]` |

**For each candidate, fill in:** who (role only), what they wanted, what you wanted, why both were reasonable, what you did, what was decided, what happened after, what you'd do differently.

---

## Cluster 1: Disagreed with a peer on a technical decision

**Common phrasings**
- "Tell me about a time you disagreed with a coworker."
- "Tell me about a technical disagreement with a peer. How did you resolve it?"
- "Describe a time someone on your team proposed an approach you thought was wrong."
- "Tell me about a design review where you pushed back."
- "Have you ever had to convince a senior engineer they were wrong?"
- "Tell me about a time you had a strong opinion that others didn't share."

**What they're really asking:** Do you resolve disagreements with criteria and evidence, and can you be persuaded?

**Use:** Story #12 `[your real one]`. Backup: Story #13 (an ARB design you pushed back on) or Story #2 (stretched vs. alternatives, if someone argued for a different topology).

**Answer skeleton (~2.5 min)**
1. **Context (20 s):** the decision, why it mattered, who the other person was (role, not name).
2. **The disagreement (20 s):** both positions, stated fairly. Say what was reasonable about theirs.
3. **What you did (75 s):** framework beats 2, 4, 5, 6 (their constraint, agree on criteria, experiment, record).
4. **Result (20 s):** the decision, a number if you have one, the relationship after.
5. **Learning (15 s):** a changed habit or mechanism.

<details><summary>Sample answer (template; fill the brackets with your real story)</summary>

> On the managed multi-cluster platform, `[a peer / a tenant platform's tech lead]` and I disagreed about `[the decision, e.g., how member-cluster registration should handle credential rotation]`. They wanted `[their option]` because `[their real constraint: delivery date, simplicity, on-call cost]`. I wanted `[your option]` because `[your reason: blast radius, DR guarantee, operability]`.
>
> We'd gone back and forth in the PR comments without converging, so I set up a 30-minute call. I started by asking what they were most worried about, and it turned out to be `[their constraint]`, which I hadn't weighted enough. We agreed the decision should come down to three things: `[e.g., what happens when a data center is lost, operational cost for my team, and time to ship]`.
>
> Against those, my option won on `[criterion]` and theirs won on `[criterion]`. So I proposed a `[time-boxed spike / load test with our CL2 + KWOK scalability harness]` to settle the one we couldn't agree on. The data showed `[result]`. We went with `[decision]`, kept `[the part of their idea that was right]`, and I wrote it up as a short decision record so we wouldn't relitigate it.
>
> `[Result: shipped by X, no incidents from it, they later reused the approach.]` What I took away is to agree on decision criteria before arguing about solutions. I now open contentious design reviews with "here's how I think we should decide," and it shortens most of them.

</details>

**Likely follow-ups**
- *"What if the data had gone against you?"* Then you'd have gone with their option and said so. Show you meant it.
- *"How did the relationship end up?"* Concrete: you worked together on `[next thing]`, they asked you to review `[X]`.
- *"What was the other person right about?"* Have a real answer ready. This is the most-asked follow-up.
- *"Did you consider escalating?"* "Not yet. We hadn't tried to agree on criteria. Escalation is for when we've done that and still disagree."
- *"How long did this take?"* Show it was proportionate: a spike of days, not a month-long war.
- *"Would you make the same call today?"* Yes with nuance, or "no, because we later learned X."

**Traps**
- Picking a disagreement with zero stakes.
- Telling it as a debate you won.
- Not stating their position fairly.

---

## Cluster 2: Disagreed with your manager or leadership

**Common phrasings**
- "Tell me about a time you disagreed with your manager."
- "Have you ever pushed back on a decision from leadership?"
- "Tell me about a time you thought your org's direction was wrong."
- "Describe a time you challenged a deadline or priority given to you."
- "When have you told your manager something they didn't want to hear?"

**What they're really asking:** Will you raise concerns upward with courage and tact, and still be a good partner once the call is made?

**Use:** Story #12 `[a real manager/leadership disagreement]`. Candidates to mine: upstream Karmada work vs. internal roadmap priority; scope of the platform vs. tenant responsibilities; a timeline for GA or a migration. Backup: Story #1 if you had to push leadership toward a firm-wide standard rather than per-team DR solutions `[confirm there was real pushback]`.

**Answer skeleton (~2.5 min)**
1. **Context (20 s):** what was decided or proposed, and the stakes.
2. **Your concern (20 s):** specific, tied to business outcome (DR mandate, tenant velocity, on-call).
3. **How you raised it (60 s):** privately first, with data and an alternative, framed around *their* goals. Offered a recommendation, not just a complaint.
4. **Outcome (30 s):** they changed course, compromised, or held. If they held, you committed.
5. **Learning (15 s).**

<details><summary>Sample answer (template)</summary>

> `[Context: e.g., my manager wanted to deprioritize upstream Karmada contributions to hit an internal milestone.]` I understood why: `[the milestone / the commitment to tenants]` was real.
>
> My concern was that `[e.g., carrying internal patches for external etcd and custom CA would make every Karmada upgrade more expensive, and our tenants would pay that in slower upgrades]`. Rather than argue in a team meeting, I brought it to our 1:1 with `[data: e.g., how much time the last upgrade took with patches vs. without, or a list of patches we'd have to carry]` and an alternative: `[e.g., upstream only the two changes that blocked us, and time-box them]`.
>
> I framed it around what they cared about, `[delivery date / tenant trust]`, not around what I preferred. `[Outcome: they agreed to the time-boxed version / they held the original call.]` `[If held:]` I said I still had the concern, asked that we revisit at `[checkpoint]`, and then put my full effort into the milestone.
>
> `[Result.]` The learning: disagreement upward works best as a written tradeoff with a recommendation, delivered privately, plus an explicit checkpoint. That way neither of us has to win in the moment.

</details>

**Likely follow-ups**
- *"What would you have done if they'd said no?"* Commit, track the risk, revisit at the agreed checkpoint.
- *"Did you ever go over their head?"* Almost always no; if you did, it was transparent and agreed (see Cluster 10).
- *"How did you know you were right?"* You didn't know; you had evidence and a risk. That's why you proposed a checkpoint.
- *"How do you disagree with someone more senior than your manager?"* Same pattern, more written, routed through your manager.
- *"Was there a cost to you for pushing back?"* Honest answer. Good managers value it; mention that it strengthened trust if true.

**Traps**
- Implying your manager was incompetent.
- Pushing back publicly first.
- No alternative offered, only objection.

---

## Cluster 3: Disagree and commit

**Common phrasings**
- "Tell me about a time you disagreed with a decision but went along with it."
- "Tell me about a time you lost an argument."
- "Describe a time a decision went against your recommendation. What did you do?"
- "Have you ever implemented something you didn't believe in?"
- "Tell me about a time you had to support a decision you disagreed with."

**What they're really asking:** Can you lose gracefully, execute fully, and not say "I told you so"?

**Use:** Story #12 `[must be one where you genuinely lost]`. Look in: ARB outcomes, upstream reviewers choosing a different API shape for one of your features (Story #3), a tenant or leadership call on scope.

**Answer skeleton (~2 min)**
1. **The decision and your position (30 s).**
2. **How you made your case (30 s):** once, clearly, with data.
3. **Why you committed (20 s):** their reasons were legitimate, the decision-maker owned it, delay was worse than either option.
4. **How you executed (30 s):** concrete evidence of full commitment (you built it, defended it to others, made it succeed).
5. **What happened, and learning (20 s):** honest, including if they were right.

<details><summary>Sample answer (template)</summary>

> When I proposed `[feature / approach]` upstream in Karmada, `[a maintainer / reviewers]` preferred `[their approach]` because `[their reason: consistency with existing APIs, maintenance burden, other adopters' needs]`. I thought `[your approach]` fit our use case better because `[reason]`.
>
> I made the case once in the design discussion with `[evidence]`. The other maintainers weighed in for `[their approach]`. Their reason was legitimate: Karmada serves many adopters, not just Bloomberg, and `[their approach]` was more general.
>
> So I committed. I implemented `[their approach]`, adapted our platform to it, and `[e.g., wrote the docs for it]`. `[Outcome: it worked well / it needed a small follow-up that we did together.]`
>
> In hindsight, they were right about `[X]`. The learning I took: in a project with many stakeholders, "best for my use case" and "best for the project" are different questions, and as a maintainer my job is the second one.

</details>

**Likely follow-ups**
- *"Were they right in the end?"* Answer honestly either way. If you were right, don't gloat: "It needed a follow-up, which we did together."
- *"How did your team feel?"* You owned the message to them; you didn't blame the decision-maker.
- *"When would you not commit?"* Ethics, security, legal, or real customer harm. Then escalate (Cluster 10).
- *"How do you show commitment rather than just compliance?"* You defend it to others, invest in making it work, don't sandbag.

**Traps**
- "I disagreed but it turned out I was right" as the whole point.
- Commitment that was really quiet resentment.

---

## Cluster 4: Saying no to a stakeholder, tenant, or customer

**Common phrasings**
- "Tell me about a time you said no to a customer or stakeholder."
- "Tell me about a time you pushed back on a feature request."
- "How do you handle a stakeholder who wants something that isn't good for the system?"
- "Tell me about a time you had to set boundaries with another team."
- "Describe a time you balanced customer requests against platform health."
- "Tell me about a time a customer wanted something you couldn't give them."

**What they're really asking:** Do you protect the platform and other tenants while still solving the requester's real problem?

**Use:** Story #9 (tenant office hours + shared responsibility: internal forks aren't a paved path; tenants own their DR tests). Backup: Story #6 (SLO-driven ops: tenants fix their own issues, which meant saying no to "the platform team should page on everything").

**Answer skeleton (~2.5 min)**
1. **The request and requester (20 s):** what they wanted, why it was reasonable from their seat.
2. **Why it was a no (25 s):** cost to the platform and other tenants, not your preference.
3. **How you said no (60 s):** found the underlying need, offered an alternative, made the boundary a written policy, kept a forum open (office hours).
4. **Result (25 s):** they got what they needed; the policy scaled to other tenants.
5. **Learning (15 s).**

<details><summary>Sample answer: internal forks aren't a paved path (Story #9)</summary>

> On our managed multi-cluster platform, a tenant team asked to run `[their own fork of a Karmada component / a custom build]` because they needed `[specific behavior]` sooner than upstream would ship it. From their seat that was reasonable: they had a deadline.
>
> But we run control planes as a managed service. Every fork is a version we'd have to test, upgrade, and support on call, and it breaks the guarantee that every tenant inherits the same DR behavior. If I said yes once, I'd be saying yes to every tenant.
>
> So I said no to the fork, but not to the need. I asked what they actually needed, which was `[behavior]`. We found `[a supported path: a feature gate we could expose / an upstream change / a configuration]`. `[If true: We later rolled out tenant-configurable feature gates and API server audit config, which covered several requests like this.]` I also made the boundary explicit in our docs: forks aren't a paved path, here's the supported way to request changes, and here's how upstream contributions get prioritized. And I used our monthly tenant office hours as the forum for these requests, so nobody felt stonewalled.
>
> `[Result: they shipped on X using the supported path; we haven't had a fork request since / we handle them in office hours.]` The learning: a "no" scales only if it comes with a paved alternative and a written policy. Otherwise you're negotiating the same thing with every tenant.

</details>

<details><summary>Sample answer: tenants own their DR tests (Story #9)</summary>

> Tenants inherit DR from the platform: every control plane is stretched across data centers on a multi-data-center etcd cluster. Some tenants took that to mean the platform also owned testing that *their* applications survived a site loss. `[A tenant asked us to run DR tests for their workloads.]`
>
> I pushed back because only they know whether their app is correct after failover: their state, their dependencies, their clients. If we owned it, we'd be signing off on things we couldn't verify. I wrote a shared-responsibility split: the platform guarantees the control plane survives and fails workloads over; tenants own verifying their workloads behave correctly, and we give them the tooling and guidance to do it.
>
> `[Result: tenants run their own DR tests; we publish guidance.]` The learning: draw responsibility lines at the boundary of what each side can actually verify. "DR you never test is DR you don't have," but the right owner has to test it.

</details>

**Likely follow-ups**
- *"What if a senior leader had insisted?"* Make the cost visible in writing, offer the alternative, and if they still insist, it's their call with eyes open.
- *"How did you keep the relationship?"* Office hours, fast turnaround on the supported path, credit them publicly when it works.
- *"How do you decide what's a yes?"* Benefits many tenants, fits the paved path, supportable on call, aligned with roadmap.
- *"Did you ever reverse a no?"* Good chance to show flexibility: e.g., feature-gate customization became a supported feature once multiple tenants needed it.
- *"How do you prioritize tenant requests?"* See `08-failure-judgment-and-prioritization.md`, Cluster 5.

**Traps**
- Sounding bureaucratic ("it's against policy").
- No alternative offered.
- Implying the tenant was unreasonable.

---

## Cluster 5: Conflict with another team (ownership, priorities, dependencies)

**Common phrasings**
- "Tell me about a time you depended on another team that wasn't prioritizing your work."
- "Describe a conflict between teams over ownership."
- "Tell me about a time two teams had conflicting goals."
- "How do you get another team to do something they don't want to do?"
- "Tell me about a cross-team project where priorities clashed."

**What they're really asking:** Can you get cross-team outcomes without authority, by aligning incentives instead of applying pressure?

**Use:** Story #1 (Karmada adoption, working with 10+ core platform and infrastructure teams and 50+ engineers). Backup: the platform's dependency teams (managed etcd, DNS-based service discovery, Kubernetes-as-a-service, identity/secrets) `[a specific clash]`, or Story #9 for ownership boundaries.

**Answer skeleton (~2.5 min)**
1. **The dependency or overlap (20 s).**
2. **Their priorities, stated fairly (20 s).**
3. **What you did (75 s):** found what they got out of it, made the ask small and concrete, offered to do the work yourself (PR into their repo, design doc), made ownership explicit.
4. **Result (25 s).**
5. **Learning (10 s).**

<details><summary>Sample answer (template built on the platform's dependencies)</summary>

> Our platform composes existing internal services rather than building new ones: a managed etcd service, DNS-based service discovery, a Kubernetes-as-a-service for the clusters underneath, and identity and secret management. That was a deliberate choice, but it meant our roadmap depended on teams with their own priorities.
>
> `[Specific: e.g., we needed X from the managed etcd team to support multi-data-center clusters per control plane, and it wasn't on their roadmap.]` From their side, `[their reason: other customers, staffing]`.
>
> I didn't escalate. I set up time with their lead to understand their roadmap and looked for overlap: `[what they got: a big new consumer, a validated multi-DC pattern, fewer bespoke requests from individual teams]`. I made the ask as small as possible and offered `[my team to write the design / contribute the change / take on testing]`. We wrote down who owned what, so on-call boundaries were clear.
>
> `[Result.]` The learning: cross-team asks land when they're small, concrete, come with your own effort attached, and are framed as a win for the other team.

</details>

**Likely follow-ups**
- *"What if they still said no?"* Find a workaround you can own, document the risk, escalate jointly with both options written up.
- *"How did you decide who should own it?"* Whoever can verify it and operate it; same principle as Story #9.
- *"How do you keep dependencies from slipping?"* Regular syncs, written milestones, early warning when slippage threatens your dates.
- *"How many teams were involved?"* Story #1: 10+ core platform and infrastructure teams, 50+ engineers.

**Traps**
- "I escalated to their director" as the main move.
- Portraying the other team as lazy.

---

## Cluster 6: Disagreement in open source

**Common phrasings**
- "Tell me about a time you had to convince people you don't work with."
- "How do you influence a project you don't control?"
- "Tell me about a time a reviewer pushed back on your design."
- "Describe getting a contentious change accepted upstream."
- "How do you handle disagreement in a community with competing company interests?"

**What they're really asking:** Influence with zero authority, across company boundaries, where the only currency is trust and good engineering.

**Use:** Story #3 (upstreaming external etcd, custom CA, token rotation into Karmada). Backup: Story #5 (growing contributors) or proposals you've driven (conformance testing, scalability testing).

**Answer skeleton (~2.5 min)**
1. **What Bloomberg needed and why upstream (20 s):** avoid carrying a fork; benefit other adopters.
2. **The pushback (20 s):** what reviewers worried about (scope, generality, maintenance).
3. **What you did (75 s):** reframed around other adopters' needs, wrote a proposal, split into smaller PRs, addressed review thoroughly, offered to maintain it.
4. **Result (25 s):** merged; you became a maintainer and lead the operator sub-project.
5. **Learning (10 s).**

<details><summary>Sample answer (Story #3)</summary>

> Our platform needed Karmada control planes to use an external, multi-data-center etcd cluster, a custom CA so stretched control planes could share an API server CA, and token rotation for short-lived credentials. We could have patched these internally, but every internal patch makes upgrades more expensive, so I wanted them upstream.
>
> `[The pushback: e.g., reviewers worried external etcd support would add configuration surface and maintenance burden for a use case they thought was niche.]` That was a fair concern for a project serving many adopters.
>
> So I reframed it from "Bloomberg needs this" to "any adopter who needs HA across failure domains needs this." I wrote `[a proposal / issue]` describing the general use case, kept the API minimal, `[split the work into smaller PRs]`, and worked through review feedback. I also committed to maintaining it.
>
> All three landed upstream. The broader result is that I became a maintainer and now lead Karmada's operator sub-project. The learning: upstream, you're not asking for a feature, you're proposing to share a maintenance burden. Proving you'll carry it is half the argument.

</details>

**Likely follow-ups**
- *"Was there a feature that didn't get accepted?"* Have an honest answer `[supply]`.
- *"How do you handle competing company interests?"* Argue for the general use case, keep company specifics out, respect maintainer consensus.
- *"How did you earn trust as a newer contributor?"* Small fixes first, docs, reviewing others' PRs, showing up consistently.
- *"How long did it take?"* `[supply rough timeline]`

**Traps**
- Making it sound like Bloomberg pushed features through by weight.
- Getting lost in technical detail.

---

## Cluster 7: Difficult coworker or interpersonal friction

**Common phrasings**
- "Tell me about a time you worked with someone difficult."
- "Describe a coworker whose style clashed with yours."
- "Tell me about a time a teammate wasn't pulling their weight."
- "How have you handled someone who was dismissive or abrasive in reviews?"
- "Tell me about a time you had to give someone hard feedback."

**What they're really asking:** Empathy and maturity. Do you assume good intent, address it directly and privately, and focus on behavior and impact?

**Use:** Story #12 `[a real interpersonal one, anonymized]`. Search: a reviewer (internal or upstream) whose tone was harsh; a tenant contact who escalated everything; a collaborator who missed commitments.

**Answer skeleton (~2 min)**
1. **The situation and its impact on the work (25 s):** behavior, not personality.
2. **Understanding them (25 s):** what you learned about their pressures.
3. **The direct conversation (40 s):** private, specific (situation, behavior, impact), a request.
4. **Result (20 s):** what changed, including in you.
5. **Learning (10 s).**

<details><summary>Sample answer (template)</summary>

> I worked with `[role]` who `[specific behavior: e.g., left blunt, one-line rejections on PRs from newer contributors]`. The impact was real: `[e.g., two newer contributors stopped submitting]`, which mattered because I was trying to grow the contributor base.
>
> Before reacting, I tried to understand it. `[What you learned: they were reviewing far more PRs than anyone else and were short on time.]` So the behavior made sense, even if the effect wasn't okay.
>
> I talked to them privately. I described a specific review, the effect it had, and asked if we could try something: `[e.g., a short review checklist and me taking first-pass reviews on newcomer PRs to lighten their load]`. I didn't frame it as their tone being bad. I framed it as both of us wanting more people able to contribute.
>
> `[Result.]` The learning: most "difficult people" are overloaded people. Fixing the load often fixes the behavior.

</details>

**Likely follow-ups**
- *"What if it hadn't changed?"* Raise it with their manager or a lead, framed around impact, after you'd addressed it directly.
- *"Did you involve your manager?"* Only after a direct conversation, unless it was harassment or a safety issue.
- *"What did you learn about yourself?"* Something real about your own style.

**Traps**
- Character assassination.
- Going to their manager first.
- Stories with harassment or HR implications (avoid).

---

## Cluster 8: You were wrong and changed your mind

**Common phrasings**
- "Tell me about a time you changed your mind."
- "Tell me about a time someone convinced you that you were wrong."
- "Describe a strong opinion you held that you later reversed."
- "When did a junior engineer teach you something?"

**What they're really asking:** Intellectual honesty. Staff engineers update fast on evidence.

**Use:** Story #12 `[a real reversal]`. Story #2 can work if *you* initially favored the standalone topology or another approach and changed course `[confirm]`. A mentee from Story #5 who proposed a better scheduler design is a strong option `[supply]`.

**Answer skeleton (~2 min):** your original position and why it was reasonable (25 s) → what changed your mind: data, an argument, a prototype (40 s) → how you acknowledged it publicly (20 s) → result (20 s) → learning about how you form opinions (15 s).

<details><summary>Sample answer (template)</summary>

> I was initially convinced that `[position]`, because `[reason]`. `[Who]` disagreed and argued `[their view]`. I pushed back at first, but they showed me `[evidence / prototype / production data]`, and it was clear that `[what you'd missed]`.
>
> I said so directly, in the same forum where I'd argued the other side, and we went with their approach. `[Result.]` The learning: I now try to write down what evidence would change my mind *before* a debate. If I can't name any, I'm holding a belief, not a technical position.

</details>

**Likely follow-ups**
- *"How did you acknowledge it?"* Publicly, giving credit.
- *"How do you know when to hold firm vs. update?"* Safety, security, and DR guarantees you hold firm on; implementation choices you update on evidence.

**Traps:** a fake reversal ("I was wrong to work so hard"); a reversal with no cost.

---

## Cluster 9: Mediating conflict between others

**Common phrasings**
- "Tell me about a time you resolved a conflict between two team members."
- "Two senior engineers disagree on a design and come to you. What do you do?"
- "How do you build consensus on a team that's split?"
- "Tell me about a time you had to break a tie."

**What they're really asking:** As a TL/Staff, can you run a fair decision process instead of picking a favorite?

**Use:** `[a real instance from your team, mentees, or Karmada design discussions]`. Backup: a hypothetical answer using your process (see `11-curveballs-hypotheticals-and-recovery.md`).

**Answer skeleton (~2 min):** the two positions, stated fairly (25 s) → you met each 1:1 to understand constraints (20 s) → got them to agree on criteria, often in writing (30 s) → decided or helped them decide; documented, including the rejected option (25 s) → follow-up with the person whose option lost (10 s) → learning (10 s).

<details><summary>Sample answer (template)</summary>

> Two engineers `[on my team / in the Karmada community]` disagreed on `[design question]`. One wanted `[A]` for `[reason]`, the other `[B]` for `[reason]`. It had stalled `[the PR / the proposal]` for `[time]`.
>
> I talked to each separately first, so neither had to defend a position in front of the other. Then I brought them together with one question: what are we optimizing for? We wrote down three criteria, `[e.g., upgrade safety, operational simplicity, time to ship]`. Scored against those, `[A]` won, but `[B]` had a good point about `[X]`, which we folded in.
>
> I wrote it up with the rejected option and why, and I checked in with `[the engineer whose option lost]` afterward to make sure they felt heard, and I gave them `[ownership of the part they cared about]`. `[Result.]` The learning: most design deadlocks are really disagreements about priorities. Making the criteria explicit unsticks them.

</details>

**Likely follow-ups**
- *"What if they still disagreed after criteria?"* As the owner, you make the call and own it, with a revisit date.
- *"What if one of them was more senior?"* Same process; seniority isn't a criterion.
- *"What if it's personal, not technical?"* Address separately and privately; see Cluster 7.

**Traps:** deciding by seniority; splitting the difference to keep the peace.

---

## Cluster 10: Escalation

**Common phrasings**
- "When do you escalate?"
- "Tell me about a time you escalated an issue."
- "Tell me about a time you should have escalated sooner."
- "How do you escalate without damaging the relationship?"

**What they're really asking:** Judgment about timing, and whether escalation is a tool you use together or a weapon you use against someone.

**Use:** Story #12 or a dependency-team blocker `[supply]`. Story #9 if a tenant escalated to your leadership and you handled it.

**Your escalation rule (say this):**
> "I escalate when we've agreed on the facts and criteria and still disagree, when the decision is above both of us, or when waiting costs more than the disagreement, like a safety, security, or DR risk. And I do it jointly: 'We disagree, here are both options written fairly, we need a decision by Friday.'"

| Escalate | Don't escalate (yet) |
|---|---|
| Real risk to reliability, security, compliance, or DR guarantee | You haven't talked to the person directly |
| Decision belongs to someone above both of you | You haven't agreed on criteria |
| Deadlock after a real attempt, and time matters | It's a preference, not a risk |
| A dependency slip threatens a commitment you made | You want a manager to "win" for you |

**Answer skeleton (~2 min):** situation and why it couldn't be resolved at your level (30 s) → what you tried first (30 s) → how you escalated: joint, written, both options, a deadline (30 s) → outcome (20 s) → learning, ideally about timing (10 s).

<details><summary>Sample answer (template)</summary>

> We needed `[X]` from `[team]` to `[hit a commitment to tenants]`, and it kept slipping. I'd met their lead twice, we agreed on the facts, but their priorities came from `[their leadership]`, so neither of us could resolve it.
>
> I told their lead I thought we needed our managers to make the call, and asked if they'd co-write a one-pager with me. We wrote both options, the cost of each, and the date we needed a decision by. It went to `[both managers]`, who decided `[outcome]` in one meeting.
>
> The learning: I'd waited `[two weeks]` longer than I should have. Now, when a dependency threatens a commitment, I set an escalation date at the start, so escalating is the plan, not a betrayal.

</details>

**Likely follow-ups**
- *"How did the other person react?"* Because it was joint, it didn't damage trust.
- *"Have you ever escalated and regretted it?"* Honest answer `[supply]`.
- *"Have you been escalated on?"* How you responded: no defensiveness, got the facts, fixed or explained.

**Traps:** escalating unilaterally; escalation as the first move; never escalating anything.

---

## Rapid-fire variants

| Phrasing | Cluster | Story |
|---|---|---|
| "Tell me about a time you had to deliver bad news to a stakeholder." | 4 | #9, #7 |
| "How do you handle a tenant who escalates every issue to leadership?" | 4, 10 | #9 |
| "Tell me about a time your team disagreed on priorities." | 9 | `[supply]` |
| "Describe a time you influenced a decision without authority." | 5, 6 | #1, #3 (see also `05-leadership-direction-and-influence.md`) |
| "Tell me about a time you pushed back in a design review." | 1 | #13 |
| "Tell me about a time you challenged the status quo." | 2 | #6, #1 |
| "Tell me about a time you had to compromise." | 1, 3 | #12, #3 |
| "Describe a disagreement about on-call or alerting." | 1, 5 | #6 |
| "Tell me about a time someone took credit for your work." | 7 | `[only if real; keep it generous]` |
| "Tell me about a time you disagreed with a customer's technical approach." | 4 | #9, #13 |
| "How do you handle a stakeholder with unrealistic expectations?" | 4 | #9 |
| "Tell me about a time a migration's timing caused conflict." | 1, 5 | #7 `[confirm]` |
| "How do you say no to your manager?" | 2 | #12 |
| "What do you do when two stakeholders want opposite things?" | 9, 4 | #9 |
| "Tell me about feedback you disagreed with." | 8 | `[supply]` |

---

## Self-check

1. Can you state the other side's position in your Story #12 so fairly that they'd agree with your summary?
2. Do you have one conflict story you won through process and one you lost or were wrong about?
3. For each conflict story, what was the other person right about?
4. What's your escalation rule, in one sentence?
5. Does each story end with a mechanism (a policy, doc, standard, decision record) rather than just "they agreed"?
6. Have you checked that no story needs confidential details or makes a Bloomberg colleague look bad?
7. Can you tell your Story #9 "no" so that the tenant sounds reasonable and you sound helpful?

## Drill (30 min)

1. **Fill the worksheet** above: 10 minutes, one line per prompt.
2. **Pick two** candidates and draft them in STAR-L (`03-storytelling-framework.md`). Mark which framework beats you used.
3. **Say each out loud** in under 2.5 minutes. Record it.
4. **Play it back** and check: did you say "I" for your decisions? Did you name what they were right about? Is there a number? Is the learning a mechanism?
5. **Answer three follow-ups** for each, picked at random from the lists above.
6. Add both to `04-story-bank.md` as Story #12a and #12b.
