# 05. Leadership, Direction, and Influence

This is the module that decides your level. Every Staff loop asks some version of "how do you set direction and get people who don't report to you to follow it?" Senior candidates describe good work they did. Staff candidates describe a call they made on an ambiguous, cross-team problem, why it was the right call, and how they brought other people with them.

Your strongest material lives here: you made Karmada the firm's standard for DR/HA, designed the platform through three architectures, pushed features upstream into a project you didn't control, and coordinated across teams for a CNCF graduation. The goal of this module is to make sure each of those stories lands as *direction and influence*, not as "I built a thing."

---

## At a glance

| Question family | What they're really testing | Best stories |
|---|---|---|
| Setting technical direction / biggest technical decision | Do you make high-leverage calls with sound reasoning, and own them? | #1 (Karmada adoption), #2 (architecture evolution) |
| Influencing without authority | Can you move people you don't manage: peers, other teams, upstream, leadership? | #3 (upstreaming), #1, #9 (shared responsibility) |
| Leading across teams | Can you coordinate many groups toward one outcome? | #4 (CNCF graduation), #1 (10+ teams onboarded) |
| Ambiguity / no clear owner | Do you create structure and ownership where there is none? | #1, #10 (Hoppscotch) |
| Strategy and long-term vision | Can you see 1–3 years out and turn it into a roadmap? | #2, #1, #6 |
| Business impact | Do you connect technical work to outcomes leadership cares about? | #1 (DR mandate), #6 (on-call, tenant velocity) |
| Alignment, consensus, design docs | How do you get a decision made and stuck? | #3, #1, #13 (ARB) |
| Force multiplier via standards | Do you scale your judgment beyond your own code? | #13 (ARB), #1 (paved path), #6 (upstream guide) |
| Building trust with leadership | Would a director give you a hard problem and walk away? | #1, #4 |
| Changing someone's mind / being convinced | Do you persuade with evidence, and can you be persuaded? | #3, #12 |

**The one-line Staff test for this module:** *Would this person make the right call on an ambiguous, cross-team problem without supervision, and get the organization to actually follow it?*

---

## What Staff looks like here

| Dimension | Senior answer | Staff answer |
|---|---|---|
| Scope | "My team needed X, so I built it" | "Teams across the firm were each solving DR on their own; I made it a platform guarantee" |
| Decision | Describes what was chosen | Describes the alternatives, the criteria, and why the others lost |
| Influence | "I convinced my manager" | "I aligned 10+ teams, upstream maintainers, and leadership, each with different incentives" |
| Ownership | Executes a direction someone else set | Sets the direction, then hands execution to others |
| Time horizon | This quarter's deliverable | Where the platform needs to be in 2 years, and the sequence to get there |
| Outcome | "It shipped" | "Standing up DR-compliant control planes went from weeks to minutes; 40+ fleets now inherit DR" |
| Learning | "It went well" | "Here's what I'd change about how I drove the decision" |

**Anti-signals that cost you the level here:**
- Every story is about your own code.
- "We decided" with no clear "I" for the call itself.
- No one ever disagreed with you (it sounds either small or untrue).
- Escalating to a manager as the first move.
- Describing a direction without describing how you got people to adopt it.
- Vision with no sequencing ("we should be multi-region") and no tradeoffs.
- Claiming the whole team's or the whole community's work.

---

## 1. Setting technical direction / biggest technical decision

**Common phrasings:**
- "Tell me about the most important technical decision you've made."
- "Tell me about a time you set the technical direction for a team or org."
- "What's the biggest technical bet you've made? How did it play out?"
- "Tell me about a decision with long-term consequences."
- "Walk me through a time you chose between two good options."
- "Tell me about a technology you introduced to your organization."
- "What's a decision you made that others disagreed with at first?"
- "Tell me about a design you're most proud of."

**What they're really asking:** Did you choose the problem (not just the solution)? Did you weigh real alternatives with explicit criteria? Did the decision stick across the org? Do you own the downsides as well as the upsides?

**Use:** Story #1 (Karmada adoption as the DR/HA standard). Backup: Story #2 (architecture evolution), which works best for "design you're proud of."

**Answer skeleton (~2.5 min):**
1. **Stakes (20 s):** the DR mandate; teams solving it independently; the cost.
2. **The decision and the options (45 s):** what you were choosing between, your criteria.
3. **Why this one, and what you gave up (30 s):** the tradeoff stated out loud.
4. **How you got it adopted (30 s):** the influence part. Don't skip it.
5. **Result (15 s):** weeks → minutes, 40+ fleets.
6. **Learning (10 s):** what you'd do differently.

<details>
<summary>Sample answer: Story #1 (Karmada adoption)</summary>

> At Bloomberg, surviving the loss of a whole data center is an organizational mandate. When I started on this, platform teams were each solving multi-cluster DR on their own, which meant duplicated effort, inconsistent guarantees, and a lot of teams becoming accidental experts in failover.
>
> I saw the decision as two questions. First, should DR be something every team builds, or a guarantee a platform provides? I argued for a platform, because the mandate applied to everyone and the hard parts (control-plane HA, cross-data-center state, failover) are the same for everyone. Second, what do we build it on? The options were [roughly: build our own federation layer, adopt an Argo-style GitOps fan-out, or adopt a multi-cluster orchestrator like Karmada]. My criteria were: native policy-based placement and failover, a Kubernetes-native API so tenants didn't learn something new, an active upstream we could influence, and an operator model we could run as a service.
>
> I chose Karmada. What I gave up: it was a younger project than some alternatives, and we'd depend on an upstream we didn't control. I mitigated that by planning to become a contributor, not just a consumer. That's where the upstreamed features and my maintainer role came from.
>
> Adoption was the harder part. I worked with [10+] foundational platform and infrastructure teams and well over 50 engineers. I started with [the first tenant team] because their need was urgent and their success would be visible, and I made the platform compose existing internal services (DNS-based discovery, managed etcd, identity) instead of asking those teams to adopt something new.
>
> The result: Karmada became the firm's standard for meeting the DR/HA mandate. Standing up a DR-compliant control plane went from weeks to minutes, and the platform now backs 40+ cluster fleets across AI/ML, data and analytics, and general compute.
>
> If I did it again, [I'd invest in automated DR verification from day one. DR you never test is DR you don't have, and we added that later than I'd like].

</details>

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "Who disagreed, and how did you handle it?" | Name the concern (dependency risk, maturity, a team preferring their own approach) and the specific thing you did: data, a prototype, a mitigation. `[supply the real objection]` |
| "Why not build it yourselves?" | Undifferentiated work; upstream already had placement and failover; building would make you the only maintainers forever. |
| "What would have made you reverse the decision?" | Name a concrete tripwire: e.g., upstream stalling, or failover semantics that couldn't meet the mandate. Shows you held the decision with appropriate confidence. |
| "What part was you vs. the team?" | "I made the platform and technology call and designed the architecture; the team built and operates it with me." |
| "How did you measure success?" | Time to DR-compliant control plane (weeks → minutes), fleets onboarded (40+), tenants inheriting DR rather than building it. |
| "What did it cost?" | Upstream dependency, cross-DC etcd latency (Story #2), the time you invested upstream. |
| "How did leadership sign off?" | Frame it as aligning with the mandate they already owned; `[supply: who sponsored it and how]`. |

**Traps:**
- Spending 90 seconds explaining what Karmada is. One sentence: "a CNCF multi-cluster orchestrator."
- Presenting it as obvious. If it was obvious, it wasn't a Staff decision.
- Forgetting the adoption half.

---

## 2. Influencing without authority

**Common phrasings:**
- "Tell me about a time you influenced people who didn't report to you."
- "How do you get buy-in for an idea?"
- "Tell me about a time you convinced another team to change direction."
- "Tell me about getting something done that required a team outside your org."
- "How have you influenced an open-source project / external partners?"
- "Tell me about a time you influenced senior leadership."
- "Tell me about a time you had to sell an unpopular idea."
- "How do you lead when you don't have formal authority?"

**What they're really asking:** Staff ICs have no reports, so all their leverage is influence. They want to see that you understand *other people's incentives* and move them with evidence and empathy, not volume or escalation.

**Use:** Story #3 (upstreaming external etcd, custom CA, token rotation into Karmada). Backups: Story #1 (internal adoption), Story #9 (getting tenants to own their DR tests).

**Answer skeleton (~2.5 min):**
1. **What you needed and who held the decision (20 s).**
2. **Their incentive and their objection (30 s):** this is the beat that separates Staff from Senior.
3. **What you did to move them (60 s):** specific tactics (see §6).
4. **Result (20 s).**
5. **Learning (15 s).**

<details>
<summary>Sample answer: Story #3 (upstreaming into Karmada)</summary>

> Our platform needed things upstream Karmada didn't support: running control planes against an external etcd cluster (we use a managed, multi-data-center etcd), bringing a custom CA (stretched control planes behind one endpoint need a shared API server CA), and token rotation for short-lived credentials. We had two options: maintain an internal fork, or get the features accepted upstream.
>
> A fork would have been faster for us and expensive forever, and it would have contradicted something I was telling our own tenants: internal forks aren't a paved path. So I chose upstream, which meant convincing maintainers I had no authority over.
>
> Their incentive was different from ours. Maintainers care about the project's general usefulness and their own maintenance burden, not one company's architecture. So I framed each feature as a general need: external etcd is how most serious operators run etcd; custom CA matters to anyone with an enterprise PKI. I wrote [proposals / issues] that led with the use case and the operational problem, not our implementation. I kept changes small and reviewable, wrote the docs and tests myself so they weren't taking on burden, and I showed up consistently in reviews and community meetings.
>
> All three landed upstream. Beyond the features, that consistency is a big part of why I became a maintainer and later the tech lead of the operator sub-project.
>
> What I learned: in open source, influence is earned by reducing other people's burden. I now use the same approach internally: when I need another team to change something, I show up with the design, the tests, and an offer to do the work.

</details>

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "Was there a feature they pushed back on?" | `[supply a real review disagreement]`. Show you adapted the design to their constraints rather than forcing yours. |
| "What if they'd said no?" | You'd have kept the change as a thin, well-isolated extension and kept advocating, with a plan to retire it. Not a long-lived fork. |
| "How long did it take?" | `[supply]`. Admit upstream is slower; explain how you planned around it. |
| "How do you influence leadership specifically?" | Lead with their outcome (mandate, risk, cost), one page, a clear ask, options with your recommendation. |
| "How do you influence a peer team with different priorities?" | Find what they're measured on; make your ask cheap for them (you do the work, compose their service). Tie to Story #1. |
| "When did influence fail?" | Have one ready. `[supply]`. Pair with the learning. |

**Traps:**
- Influence = "I explained it well and they agreed." Show the *objection*.
- Making the other side sound foolish.
- Leading with escalation.

---

## 3. Leading across teams

**Common phrasings:**
- "Tell me about a project where you coordinated multiple teams."
- "Tell me about a time you rallied people around a goal."
- "How do you drive a cross-org initiative?"
- "Tell me about a time you represented your org externally."
- "Tell me about the most complex stakeholder landscape you've worked in."
- "How do you keep many teams aligned over a long project?"

**What they're really asking:** Can you run a multi-party effort with nobody reporting to you, keep it moving, and land it?

**Use:** Story #4 (CNCF graduation). Backup: Story #1 (onboarding 10+ teams).

**Answer skeleton (~2.5 min):** goal and why it mattered (20 s) → the parties and what each needed (30 s) → how you structured the work: owners, cadence, artifacts (50 s) → a moment it got stuck and what you did (30 s) → result (15 s) → learning (10 s).

<details>
<summary>Sample answer: Story #4 (CNCF graduation)</summary>

> Karmada was moving from incubation to graduation in the CNCF, and graduation requires evidence of real production adoption, including interviews with end users. Bloomberg was the first North American adopter, so our adoption story mattered for the project.
>
> The work spanned several groups with different incentives. Inside Bloomberg, the [3] teams that would speak in the adopter interviews were busy running production platforms; this wasn't their priority. Our CTO office and OSPO cared about how Bloomberg was represented externally. Upstream, the maintainers needed the evidence on the CNCF's timeline.
>
> I partnered with the CTO office and OSPO to get Bloomberg onboarded as a public adopter and agree on what we could say. Then I led the adopter interviews across the three teams. I made it cheap for them: I prepared each team with the questions in advance, framed what was public versus confidential, and scheduled around their on-call. On the upstream side, I helped drive the graduation work as a maintainer.
>
> [Where it got stuck: supply, e.g., approval for external representation took longer than expected; what you did.]
>
> Karmada graduated. Beyond the result, Bloomberg's presence as a named adopter and our CNCF reference architecture made it easier for other enterprises to adopt.
>
> What I learned: in cross-team work, the person coordinating has to absorb the overhead. If participating costs other teams an hour of prep, they'll defer it; if it costs them 20 minutes with everything ready, they'll show up.

</details>

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "What did you personally own vs. others?" | Bloomberg's adopter onboarding (with CTO/OSPO) and the adopter interviews; the graduation itself was a community effort you helped drive. |
| "How did you track progress?" | `[supply: doc, checklist, cadence]`. Show you made status visible. |
| "A team didn't deliver on time. What then?" | Find out why (priority, not willingness); reduce their cost; agree a new date; only then involve their manager. |
| "How did you handle confidentiality?" | Agreed public framing up front with OSPO; this is a Staff judgment signal. |
| "What's another example?" | Story #1: onboarding 10+ teams onto the platform. |

**Traps:** listing parties without explaining their incentives; taking credit for graduation itself.

---

## 4. Ambiguity and problems with no clear owner

**Common phrasings:**
- "Tell me about a time you worked on an ambiguous problem."
- "Tell me about a time requirements were unclear or kept changing."
- "Tell me about a problem nobody owned that you took on."
- "Tell me about a time you had to make a decision without enough information."
- "How do you approach a problem when you don't know where to start?"
- "Tell me about something you did that wasn't in your job description."

**What they're really asking:** Do you wait for clarity, or create it? Do you scope sensibly, make reversible decisions quickly, and set up ownership so it survives you?

**Use:** Story #10 (Hoppscotch: nobody owned the Postman problem). Backup: Story #1 (DR solved team by team, no platform owner).

**Answer skeleton (~2 min):** the gap and who was affected (20 s) → why nobody owned it (15 s) → how you created structure: framing the problem, constraints, first reversible step (45 s) → how you made ownership durable (20 s) → result (15 s) → learning (10 s).

<details>
<summary>Sample answer: Story #10 (Hoppscotch)</summary>

> Engineers across the firm used Postman to work with APIs, and it no longer fit our policies [supply the specific policy reason in public-safe terms, e.g., data leaving the environment]. That affected hundreds of engineers, but it wasn't any team's job: security owned the policy, not the replacement, and every team was improvising on its own.
>
> I treated it as a product problem. First I pinned down the constraints: it had to be policy-compliant, self-hosted, and close enough to Postman that people would actually switch. I evaluated [options] and chose Hoppscotch, an open-source API client we could run internally. I [stood up a working instance / partnered with security for approval / supply] before asking anyone to commit, because a working tool is a far stronger argument than a proposal.
>
> I made it durable by [supply: documenting it, finding an owning team, setting up support]. Today 375+ engineers use it.
>
> What I learned: when nobody owns a problem, the first job is to make the problem legible (who's affected, what the constraints are) and the second is to make sure it has an owner after you.

</details>

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "Why did you take it on instead of your actual work?" | Show you sized it: small cost to you, large benefit to many; your manager agreed `[confirm]`. |
| "How did you decide with incomplete information?" | Separate reversible from irreversible decisions; move fast on reversible ones. |
| "Who owns it now?" | Have a clear answer. Staff work survives the person. |
| "How did you know it was the right solution?" | Adoption (375+) and policy compliance. |
| "What about Story #1 as ambiguity?" | Ambiguity was "should DR be a platform at all, and whose?" You framed it before solving it. |

**Traps:** confusing "ambiguous" with "hard"; heroics with no ownership plan.

---

## 5. Strategy, long-term vision, and roadmaps

**Common phrasings:**
- "Where should your platform be in two years? How did you decide?"
- "How do you build a technical roadmap?"
- "Tell me about a time you planned for scale you didn't have yet."
- "How do you balance short-term delivery against long-term investment?"
- "What's your vision for [our area]?" (common with directors)
- "How did your architecture evolve, and did you anticipate it?"

**What they're really asking:** Do you see the system as something that evolves? Can you sequence work so each step delivers value and unlocks the next?

**Use:** Story #2 (standalone → stretched → source of truth + operators). Backups: Story #6, Story #8 (investing in scale confidence early).

**Answer skeleton (~2.5 min):** where you started and why it was right then (20 s) → what forced each evolution (60 s) → how you sequenced it so tenants kept working (30 s) → what's next and why (25 s) → learning (15 s).

<details>
<summary>Sample answer: Story #2 (architecture evolution as a roadmap)</summary>

> I think of the platform's roadmap as three stages, each fixing the biggest risk of the last.
>
> **Standalone:** one management cluster ran the Karmada operator plus our own managed operator. It proved the model fast, but the management cluster was a single point of failure, which is unacceptable for a DR mandate. I knew going in it was a stepping stone.
>
> **Stretched:** each control plane runs on management clusters in multiple data centers behind one DNS-based endpoint, with a shared API server CA from cert-manager and one multi-data-center etcd cluster behind both. The tradeoff was cross-data-center etcd latency in exchange for surviving a site loss. For control-plane writes, that's the right trade.
>
> **Source of truth plus operators:** tenants declare a control plane as a custom resource against one API, and operators reconcile ingress, PKI, etcd, and member-cluster registration. That's what made it self-service and took setup from weeks to minutes.
>
> On sequencing: each step had to deliver tenant value on its own, because I couldn't ask tenants to wait a year. [Supply: how you migrated existing control planes between stages.]
>
> What's next: [from public material and office hours: conformance and scalability testing upstream, automated DR verification, better failover insights]. My principle is that the next investment goes where the next class of failure is.
>
> Learning: [supply, e.g., I'd have designed for the source-of-truth model earlier, because retrofitting self-service was harder than building it in].

</details>

**How you build a roadmap (have this ready as a short answer):**
> I start from the outcome we owe tenants and the biggest risk to it, not a feature list. I keep a short list of bets ranked by risk reduced and tenant value, I make each step independently useful, and I keep a feedback loop open: for us that was monthly tenant office hours, where we share finished, ongoing, and upcoming work and take feedback. I reserve capacity for operational work and upstream, because both compound.

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "How did you balance roadmap work against tenant requests?" | Office hours plus shared-responsibility boundaries (Story #9); say no with a reason and an alternative. |
| "What would you do differently at 10× scale?" | Name the bottleneck honestly (etcd per control plane, management-cluster capacity, registration churn) and what you'd change. |
| "How do you know you're investing early enough?" | Story #8: the scalability harness exists so you find limits before tenants do. |
| "Vision for our area?" | Use target-companies.md; tie their problem to your patterns; admit what you'd need to learn. |

**Traps:** a vision with no tradeoffs; treating the current architecture as final.

---

## 6. Business impact: connecting technology to outcomes

**Common phrasings:**
- "Tell me about the business impact of your work."
- "How do you decide what's worth building?"
- "Tell me about a time you made a technical tradeoff for business reasons."
- "How do you explain technical work to non-technical stakeholders?"
- "What's the most impactful thing you've done?"

**What they're really asking:** Do you know *why* your work matters to the company, and can you say it in their terms?

**Use:** Story #1 (DR mandate compliance as a default). Backups: Story #6 (on-call load, tenant self-service), Story #10 (policy compliance for 375+ engineers).

**The translation table (memorize):**
| Technical fact | Business framing |
|---|---|
| Stretched control planes, multi-DC etcd | "Teams meet the firm's DR mandate by adopting the platform, not by building it" |
| Weeks → minutes control-plane setup | "Platform teams ship DR-compliant services faster" |
| 40+ fleets across AI/ML, analytics, general compute | "One investment, reused across the firm's core platforms" |
| SLO-driven, symptom-based alerting | "Less on-call toil; tenants fix their own issues; reliability measured in tenant experience" |
| Upstreaming instead of forking | "No permanent maintenance tax from a private fork" |
| Hoppscotch | "Removed a policy risk for 375+ engineers" |

<details>
<summary>Sample answer: "What's the most impactful thing you've done?"</summary>

> Making DR a platform guarantee instead of every team's project. Bloomberg's workloads sit close to the global financial markets, so surviving the loss of a data center is mandatory. Before the platform, each team paid that cost separately. Now a team becomes DR-compliant by adopting the platform: standing up a compliant control plane takes minutes instead of weeks, and 40+ cluster fleets inherit it, including our streaming, managed Trino, PaaS, and data science platforms. The impact is that the firm pays for multi-cluster DR expertise once.

</details>

**Likely follow-ups:** "How did you quantify it?" (weeks → minutes, fleets, `[team-time saved if shareable]`) · "Did leadership see it that way?" · "What did you not build because the business case was weak?" (Story #9: internal forks aren't a paved path).

**Traps:** only technical metrics; inventing dollar figures you can't back up.

---

## 7. Driving alignment, consensus, and design docs

**Common phrasings:**
- "How do you drive a technical decision to closure?"
- "Tell me about a design doc or RFC you wrote. How did it go through review?"
- "How do you handle a decision where people can't agree?"
- "Tell me about a time you had to build consensus."
- "Disagree and commit: tell me about a time you did it."

**What they're really asking:** Can you turn a contested question into a decision that gets made and stays made, without unanimous agreement?

**Use:** Story #3 (upstream proposals and design review). Backups: Story #13 (ARB), Story #1.

**Your process (say it as a short framework, then a story):**
1. **Write it down.** Problem, constraints, options, recommendation, open questions. Writing exposes weak reasoning (yours included).
2. **Pre-wire.** Talk to the most affected and most skeptical people one-on-one before the big review. No one should be surprised in the meeting.
3. **Name the decider and the deadline.** Consensus isn't unanimity. Know who decides and by when.
4. **Separate facts from preferences.** Settle facts with data or a prototype; for preferences, pick and move.
5. **Record the decision and the dissent.** Note what would make you revisit it.
6. **Disagree and commit, visibly.** When it goes against you, support it fully.

<details>
<summary>Sample answer: "How do you drive a decision to closure?"</summary>

> I write it down first: problem, constraints, two or three real options, my recommendation, and what would change my mind. Then I pre-wire: before a design goes to review, I talk to the people most affected and most skeptical, so the review is about refinement, not surprise.
>
> Upstream in Karmada that looks like [a proposal / issue] for features like [external etcd support], where I led with the general use case, took review feedback into the design, and [supply a concrete change you made based on a reviewer's input]. Inside Bloomberg, I sit on the Architecture Review Board, which reviews system designs and shapes engineering-wide standards, so I see this from the reviewer side too: the designs that move fastest are the ones that state their alternatives and tradeoffs honestly.
>
> When people still disagree, I make sure we know who decides and by when, settle the factual parts with data or a quick prototype, and record the decision and the conditions for revisiting it.

</details>

**Likely follow-ups:** "Tell me about a time you disagreed and committed" (`[supply]`, see module 07) · "What makes a good design doc?" · "How do you review someone else's design?" (Story #13) · "What if the decider is wrong?"

**Traps:** "we discussed until everyone agreed" (slow and unrealistic); no mention of writing.

---

## 8. Force multiplier through standards and paved paths

**Common phrasings:**
- "How do you scale your impact beyond your own work?"
- "Tell me about a standard or best practice you established."
- "How do you raise the technical bar across an org?"
- "Tell me about a time you made other teams more effective."

**What they're really asking:** Does your judgment reach engineers who never talk to you?

**Use:** Story #13 (ARB: engineering-wide design review and DR advice), Story #1 (paved path), Story #6 (upstream reliability guide: SLIs, SLOs, runbooks), Story #8 (harness others reuse).

<details>
<summary>Sample answer</summary>

> I scale through three things: paved paths, standards, and reusable artifacts.
>
> The platform itself is the biggest one. Teams get DR by adopting it, so the right architecture is the default, not something each team has to get right.
>
> Second, I sit on Bloomberg's Architecture Review Board, which reviews system designs and shapes engineering-wide standards, and I advise platform teams on DR strategy. [Supply one specific review: the design, the risk you flagged, what changed.]
>
> Third, artifacts. When we moved to SLO-driven operations, I found existing guidance assumed a single cluster and broke at cluster boundaries. I wrote a reliability engineering guide upstream in Karmada (SLIs, SLOs, runbooks), plus metrics reference docs and Grafana dashboards, so every Karmada operator benefits, not just us.

</details>

**Likely follow-ups:** "How do you get teams to follow a standard?" (make the paved path the easy path; mandates only for real risk) · "When should a team go off the paved path?" (Story #9) · "How do you know a standard is working?"

---

## 9. Building trust with leadership

**Common phrasings:**
- "How do you work with your manager / director?"
- "Tell me about a time leadership trusted you with something big."
- "How do you keep leadership informed?"
- "Tell me about delivering bad news upward."

**What they're really asking:** Can a director give you a problem and stop worrying about it?

**What to say:** you bring problems with options and a recommendation; you surface risk early ("no surprises"); you report in outcomes; you've been trusted as product owner end to end. Being given the platform's tech lead *and* product owner role, a seat on the ARB, and the external representation work for graduation (with the CTO office and OSPO) are the evidence.

<details>
<summary>Sample answer: delivering bad news</summary>

> [Supply a real instance, e.g., a timeline slip or an incident.] My rule is that leadership should hear bad news from me first, early, with the impact, what I'm doing about it, and what I need from them. [Situation → what you said → outcome.] The trust that builds is worth more than the awkward conversation.

</details>

---

## 10. Changing someone's mind, and being convinced

**Common phrasings:**
- "Tell me about a time you changed someone's mind."
- "Tell me about a time someone changed yours."
- "Tell me about a time you were wrong about a technical decision."
- "How do you react when your idea is challenged?"

**What they're really asking:** Are you persuadable by evidence? Staff engineers who can't be convinced are dangerous.

**Use:** for persuading: Story #3 or Story #9. For being convinced: `[supply a real one: an upstream reviewer, a mentee, or a tenant who changed your design]`. This is a gap to fill; it's asked often.

**Answer skeleton for "someone changed your mind" (~2 min):** your original position and why it was reasonable (30 s) → what they showed you (30 s) → how you checked it (20 s) → what changed (20 s) → what it taught you about your own blind spots (20 s).

> Good shape: "I'd assumed X. [Name/role] showed me Y with [data/a case I hadn't considered]. I validated it by Z, changed the design to W, and it was the better call. Now I [habit you adopted]."

**Traps:** a trivial example; sounding grudging.

---

## 11. Influence tactics (your toolkit, mapped to stories)

Use these words in answers. Naming the tactic shows it was deliberate.

| Tactic | What it means | Where you used it |
|---|---|---|
| **Lead with their incentive** | Frame the ask in terms of what they're measured on | Upstream maintainers: general usefulness, low burden (Story #3); tenants: DR compliance for free (Story #1) |
| **Data over opinion** | Settle factual disputes with measurement | Scalability harness to find limits before debating them (Story #8); SLOs scoped to tenant experience (Story #6) |
| **Prototype it** | A working thing beats a proposal | `[confirm: PoC before platform]` (Story #1); Hoppscotch instance (Story #10) |
| **Pre-wire stakeholders** | No surprises in the big meeting | Adopter interviews prep (Story #4); `[design reviews]` |
| **Write it down** | Clear docs travel without you | Upstream proposals (Story #3); reliability guide (Story #6); ARB designs (Story #13) |
| **Reduce their cost** | Do the work, tests, docs | Upstream PRs with tests and docs (Story #3); composing existing internal services (Story #1) |
| **Create a forum** | A recurring place for feedback and alignment | Monthly tenant office hours (Story #9) |
| **Clear boundaries** | Shared-responsibility model | Tenants own their DR tests; forks aren't a paved path (Story #9) |
| **Earn credibility first** | Consistent contribution before big asks | Contributor → maintainer → operator sub-project TL |
| **Escalate last, together** | Joint escalation with a shared write-up of both views | `[supply if you have one]` |

---

## Rapid-fire variants

| If they ask… | Cluster | Story |
|---|---|---|
| "Tell me about a time you took initiative." | 4 Ambiguity | #10, #1 |
| "Tell me about your biggest professional accomplishment." | 1 Direction / 6 Impact | #1 |
| "How do you make decisions?" | 7 Alignment | Framework + #1 |
| "Tell me about a time you led without being the manager." | 2 Influence | #3, #4 |
| "Tell me about a time you had to get buy-in from a skeptical team." | 2 Influence | #1, #9 |
| "How do you prioritize across many stakeholders?" | 5 Strategy | #9 (see 08) |
| "Tell me about a time you saw something broken across the org and fixed it." | 4 Ambiguity / 8 Multiplier | #10, #6 |
| "How do you stay technical as a leader?" | Direction | Upstream PRs, operator sub-project TL |
| "Tell me about representing your company externally." | 3 Cross-team | #4, KubeCon talks, Program Committee |
| "Tell me about a time you challenged the status quo." | 1 / 4 | #6 (noisy alerts → SLOs), #1 |
| "Tell me about a decision you made quickly." | 4 Ambiguity | `[supply]` |
| "Tell me about a time you simplified something complex." | 1 / 6 | #2 (one CR, operators reconcile the rest) |
| "How do you get engineers to adopt your platform?" | 8 Multiplier | #1, #9 |
| "What's your leadership style?" | All | "Context over control: set direction, explain why, hand off ownership, stay close on risk." |

---

## Self-check

1. For Story #1, can you name the alternatives you rejected and one concrete reason each lost?
2. For every story in this module, can you name the person or group who disagreed and what you did about it?
3. Can you state the business framing of your platform in one sentence with no Kubernetes terms?
4. Do you have a real "someone changed my mind" story? If not, find one before your next loop.
5. Can you describe your roadmap process in under 45 seconds?
6. For Story #10, who owns Hoppscotch now?
7. Which three influence tactics do you use most, and which story shows each?
8. Have you filled every `[bracket]` in the sample answers you plan to use?

**Drill (30 min):** set a timer, answer clusters 1, 2, and 4 out loud, 2.5 minutes each, recording yourself. Then score each on: (a) clear "I" for the decision, (b) alternatives named, (c) objection named, (d) one number, (e) a learning. Redo any that miss two or more. Rubric details are in module 12.
