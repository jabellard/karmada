# 04. Your Story Bank

Behavioral prep comes down to this: **12–15 real stories, each rehearsed until you can tell it at three lengths and survive four levels of follow-up.** Every question in modules 05–10 maps back to these stories. This module drafts each story from your real work (`../interview-playbook.md`, `../profile/`, `../office-hours-log.md`) in STAR-L shape.

**Full spoken drafts of every story, including the missing ones, are in `13-story-drafts.md`.**

**Everything in `[brackets]` is yours to fill in.** I haven't invented any details. A story with brackets left in isn't ready to tell. Interviewers drill for specifics, and a made-up detail falls apart on the second follow-up.

---

## At a glance

| # | Story | Primary signals | Readiness |
|---|---|---|---|
| 1 | Karmada adoption as Bloomberg's DR/HA standard | Direction, ambiguity, influence, business impact | Mostly drafted; fill the skeptic and the alternatives you evaluated |
| 2 | Architecture evolution: standalone → stretched → source of truth | Judgment, tradeoffs, iteration | Drafted from the reference architecture |
| 3 | Upstreaming external etcd, custom CA, and token rotation into Karmada | Influence without authority, open source | Fill reviewer pushback |
| 4 | CNCF graduation: first North American adopter, adopter interviews | Cross-org leadership, representing the company | Fill your specific coordination steps |
| 5 | Contributors: fewer than 5 → 20+ | Multiplier, mentoring | Fill one named-by-role mentee's trajectory |
| 6 | Noisy alerts → SLO-driven operations | Operational excellence, customer focus | Fill the before/after metric if shareable |
| 7 | Zero-downtime ingress-nginx → Istio migration | Execution, risk management | **Needs details** |
| 8 | Scalability testing harness (CL2 + KWOK) | Engineering rigor, leverage | Drafted |
| 9 | Office hours and shared-responsibility lines; saying no | Stakeholder management, boundaries | Fill one concrete "no" |
| 10 | Hoppscotch: policy-compliant Postman replacement | Initiative, ownership beyond your role | **Needs details** |
| 11 | A real failure | Ownership, learning | **You must pick one** |
| 12 | A real conflict or disagreement | Conflict resolution, judgment | **You must pick one** |
| 13 | An Architecture Review Board case | Raising the bar beyond your team | **You must pick one** |
| 14 | Mentoring a new grad / ColorStack mentee | Growing people across differences | **You must pick one** |

**Priority order to finish:** 11, 12, 1, 5, 6, 9, 7, 13, then the rest. Failure and conflict are asked in almost every behavioral round, and they're the two you don't have yet.

---

## 1. What makes a story bank-worthy

| Test | Why |
|---|---|
| **You made decisions in it** | Interviewers credit decisions, not attendance |
| **There's tension** (skeptic, constraint, risk, or wrong turn) | No tension reads as low stakes |
| **There's a result** you can state in one sentence | Otherwise it's an anecdote |
| **It's recent** (ideally the last 2–3 years) | Old stories suggest you peaked |
| **It's public-safe** | You can tell it fully without internal names or numbers |
| **You can go four levels deep** | Follow-ups decide the level (`03-storytelling-framework.md` §5) |

**Coverage rule:** every competency in the coverage matrix (§4) should have **two** stories, so you never have to reuse your best one in back-to-back rounds.

---

## 2. The stories

Each story has: the **headline**, a **30-second version**, the **STAR-L beats**, **follow-up prep**, and the **gaps** to fill. Expand to 2–3 minutes by adding the Action decisions in full.

---

### Story #1: Leading Karmada adoption as Bloomberg's DR/HA standard

**Signals:** technical direction · ambiguity · influence across 10+ teams · business impact
**Use for:** biggest technical decision · setting direction · ambiguous problem · leading across teams · most impactful project · "tell me about a project you're proud of"

**Headline:** "I led Bloomberg's adoption of Karmada as the firm's standard way to meet its disaster-recovery mandate, and built the managed platform that took it from proof of concept to 40+ cluster fleets."

**30 seconds:**
> "Surviving a data-center loss is a firm-wide mandate at Bloomberg, and platform teams were each solving it on their own. I led the effort to solve it once: we standardized on Karmada and I designed a managed platform where tenants declare a control plane and inherit DR. I worked with 10+ core platform and infrastructure teams and 50+ engineers to get there. Setup went from weeks to minutes, and it now backs 40+ cluster fleets across AI/ML, data analytics, and general compute."

**Situation**
- DR (surviving the loss of an entire data center, including every cluster in it) is an organizational mandate; workloads sit close to the global financial markets.
- Teams were solving multi-cluster resilience on their own: inconsistent, expensive, and each needing specialists.
- `[What triggered the effort: a mandate deadline, an incident, a leadership ask, or you noticing the duplication?]`

**Task**
- As tech lead and product owner, own the answer to "how does Bloomberg do multi-cluster DR, once, for everyone?"

**Action (pick 3 decisions for a 2–3 minute telling)**
1. **Choosing the foundation.** Karmada, because its primitives (policy-based placement, automated cross-cluster failover) map cleanly onto the DR and placement outcomes needed. `[Alternatives you evaluated: e.g., Open Cluster Management, KubeFleet, Argo CD-based fan-out, building in-house. Why each lost.]`
2. **Making DR a platform guarantee, not a tenant project.** Principle: adopting the platform *is* becoming DR-compliant. Pushing that complexity onto tenants doesn't scale.
3. **Abstracting the control plane.** Tenants declare a `KarmadaControlPlane` as desired state; operators reconcile ingress, PKI, etcd, and member-cluster registration. Tenants never operate Karmada.
4. **Composing instead of rebuilding.** Reused the cluster inventory API, token exchange service and SPIFFE/SPIRE, the managed etcd service, and DNS-based service discovery. This let a small team deliver a large outcome, but it meant coordinating with each owning team.
5. **Bringing teams along.** `[How you got the first tenants: who went first, why, what you promised them.]` Office hours so the platform was built in the open (Story #9).

**Tension / skeptic:** `[Who doubted it? e.g., a platform team that preferred its own solution, someone worried about depending on a young CNCF project, leadership concerned about cost. What did they care about, and what changed their mind?]`

**Result**
- Standing up a DR-compliant control plane: **weeks → minutes**, self-service.
- **40+ cluster fleets**; tenants include the Streaming Platform (managed Flink), Managed Trino, the PaaS, and the Data Science Platform.
- Those platforms became DR-compliant by building on it instead of solving cross-data-center resilience themselves.
- It became the firm's standard. Bloomberg became the first North American adopter in Karmada's graduation (Story #4).

**Learning:** pick one that's true for you:
- "I'd have built DR verification in from day one. DR testing for management clusters is still a human-driven runbook, and *DR you never test is DR you don't have*. Now I treat verification as part of the first milestone, not a follow-up."
- "Composing internal services was the right call, but each dependency is a team with its own roadmap. I'd map those dependencies and their owners in the first week next time."

**Follow-up prep**

| Probe | Your answer (fill or confirm) |
|---|---|
| Why Karmada over the alternatives? | `[criteria and comparison]` |
| What was *your* part vs the team's? | Set direction and architecture, made build-vs-compose calls, owned tenant adoption. `[What you personally built]` |
| What was the hardest part? | `[technical (etcd topology?) or people (first adopters?)]` |
| What if Karmada had been abandoned upstream? | You became a maintainer, which reduces that risk; the patterns (control-plane abstraction, intent-based placement) outlive the tool |
| How did you get funding/headcount? | `[fill]` |
| What would you do differently? | Learning above |

---

### Story #2: Architecture evolution (standalone → stretched → source of truth + operators)

**Signals:** design judgment · tradeoffs · iterating on a design
**Use for:** design you're proud of · technical decision with tradeoffs · how a design evolved · "walk me through your architecture" · deep dive opener

**Headline:** "The platform went through three topologies, each fixing a failure mode the last one couldn't: a single point of failure, then manual bespoke deployments."

**30 seconds:**
> "We started with a standalone topology: one management cluster ran everything. It proved the model, but that cluster was a single point of failure, which is unacceptable under a DR mandate. So we stretched each control plane across management clusters in two data centers, behind one endpoint, on one multi-data-center etcd cluster. Then we put a declarative source of truth with operators on top, which turned a bespoke deployment into a self-service product."

**Action: the decisions and their tradeoffs**

| Stage | Decision | Tradeoff accepted | Why it was right |
|---|---|---|---|
| Standalone | One management cluster with the upstream operator plus managed-karmada-operator | Single point of failure | Fastest way to prove the model and the integrations |
| Stretched | Same API server CA (cert-manager) on both sites; one DNS-based unified endpoint; one multi-DC etcd cluster per control plane | Cross-DC etcd write latency; a two-site design (losing both sites takes an instance down) | Losing a site degrades an instance but doesn't take it down; control-plane write latency matters less than surviving a site loss |
| Source of truth | Minimal control plane (kube-apiserver, controller-manager, cert-manager); sot-operator, sync-operator, managed-karmada-operator | More moving parts to operate | Tenants declare intent and the platform converges; onboarding becomes two steps |

**Rejected alternative to name:** `[e.g., active/passive with a failover runbook; per-site independent control planes; a single global control plane. Why each lost.]`

**Result:** stretched topology is what runs in production; the abstraction "has held up": tenants really don't operate Karmada.

**Learning:** "The stretched design has explicit limits, and I'd rather state them than hide them: it's two-site, so losing both management clusters takes an instance down. But recovery is deterministic: bring the management clusters back and operators reconcile to ready, with registrations intact. What I'd improve is automated DR testing, which is still manual."

**Follow-up prep**

| Probe | Answer |
|---|---|
| What happens if both management clusters fail? | Instance is down; recovery is deterministic via reconcile. That's a stated limit of a two-site design |
| etcd quorum across data centers: split-brain? | `[How quorum is laid out across sites; what a partition does; the managed etcd service's role]` |
| Why not one global control plane? | `[blast radius, tenant isolation, upgrade independence: confirm your reasoning]` |
| How do tenants upgrade versions safely? | Versions are declared in the CR; `[your rollout and validation process]`; conformance testing is the gap you're closing upstream |
| What did *you* design? | `[which stages and components were yours]` |

---

### Story #3: Upstreaming external etcd, custom CA, and token rotation into Karmada

**Signals:** influence without authority · cross-org collaboration · open-source judgment
**Use for:** convince people you don't manage · work with an external team · build vs fork · "disagreement with a reviewer"

**Headline:** "The stretched topology needed features Karmada didn't have, and instead of forking, I got them accepted upstream, where I had no authority over the roadmap."

**30 seconds:**
> "Our design needed Karmada's operator to support external etcd, custom CAs, and token rotation. Forking would have been faster that quarter but meant maintaining a fork forever. I chose to upstream them. That meant writing the use case so it made sense for every adopter, not just Bloomberg, and working through maintainer review. They landed, and that work is a big part of how I became a maintainer and the operator sub-project's tech lead."

**Situation:** stretched control planes need a multi-DC etcd that isn't managed by the operator (external etcd), a shared API server CA across sites (custom CA), and short-lived credentials (token rotation).

**Action**
1. **Upstream over fork.** Rejected a long-lived internal fork: ongoing rebase cost, drift from the community, and the same reason you tell tenants "internal forks aren't a paved path" (Story #9).
2. **Framed it for the community.** `[How you generalized the use case: e.g., an issue describing the need for any HA/DR deployment before writing code.]`
3. **Worked the review.** `[Reviewer concerns: API shape, backward compatibility, scope. What you changed.]`

**Result:** features landed upstream; `[timeline]`; Bloomberg runs upstream rather than a fork; you're a maintainer and tech lead of the operator sub-project.

**Learning:** `[e.g., "Open the issue with the use case first. Agreeing on the design before the PR is faster than debating it in review."]`

**Follow-up prep:** what a reviewer pushed back on and whether they were right · what you'd have done if they'd said no · how you balanced Bloomberg's timeline with upstream's pace · `[did you run a patched build in the meantime?]`

---

### Story #4: CNCF graduation

**Signals:** cross-org leadership · coordination · representing the company externally
**Use for:** rallied multiple teams · represented your org · contributed beyond your job

**Headline:** "I helped drive Karmada's graduation in the CNCF. Bloomberg was the first North American adopter, and I led the adopter interviews across three Bloomberg teams."

**30 seconds:**
> "Graduation requires evidence of real production adoption. I partnered with Bloomberg's CTO office and open source program office to get Bloomberg on record as an adopter, and led the adopter interviews across three teams. Karmada graduated, and Bloomberg was its first North American adopter."

**Action**
- `[What graduation required from adopters and how you organized it]`
- Coordinated three internal teams for the interviews: `[how you prepared them, what you had to clear through the CTO office/OSPO]`
- `[Other graduation work: docs, security, governance items you drove]`

**Tension:** `[Getting approval to speak publicly about internal usage? Scheduling three busy teams? Confidentiality review?]`

**Result:** Karmada graduated; Bloomberg is publicly on record as the first North American adopter; a CNCF end-user case study / reference architecture followed.

**Learning:** `[fill]`

---

### Story #5: Growing contributors from fewer than 5 to 20+

**Signals:** multiplier · mentoring · building a community
**Use for:** grew other engineers · mentoring · built a team or community · "how do you develop others?"

**Headline:** "I grew Bloomberg's Karmada contributor base from fewer than 5 to over 20 engineers, and mentees went on to drive design proposals and scheduler enhancements."

**30 seconds:**
> "When I started, fewer than 5 Bloomberg engineers contributed to Karmada. I made contributing easier: mentoring people from their first PR to bigger work, and running a Contribufest with 20+ engineers. Now there are 20+ contributors, and some of my mentees drive design proposals and scheduler enhancements on their own."

**Action (the mechanism, not just the activity)**
1. **Lowered the entry cost:** `[good-first-issues curated from your own backlog? docs? pairing on first PRs?]`
2. **Contribufest** with 20+ engineers: `[how it was structured; what came out of it]`
3. **Graduated people to bigger scope:** first PR → `[feature]` → design proposal. `[One mentee, by role only: where they started, what you did, where they are now.]`
4. **Stepped back:** `[when you stopped reviewing their designs in detail]`

**Result:** fewer than 5 → **20+** contributors; mentees driving design proposals and Karmada scheduler enhancements.

**Learning:** `[e.g., "Mentoring scales when people get real ownership early. Handing someone a whole problem grew them faster than reviewing their PRs."]`

**Follow-up prep:** how you picked who to invest in · someone it didn't work for, and why · how you made time for it alongside running the platform · how you'd do this on a team you just joined.

---

### Story #6: Noisy alerts → SLO-driven operations

**Signals:** operational excellence · customer focus · reducing toil
**Use for:** improved reliability · reduced toil or on-call load · customer obsession · data-driven decision

**Headline:** "I moved our platform from noisy, infrastructure-centric alerting to SLOs scoped to tenant experience, and tenants can now fix their own issues."

**30 seconds:**
> "Our alerts were noisy and infrastructure-focused, so on-call chased things tenants never felt. And the existing observability guidance assumed a single cluster, which broke at cluster boundaries. I defined SLOs around tenant experience, built symptom-based alerting across clusters on a Thanos global view, and attached a runbook to every alert. Tenants now fix their own issues, and I upstreamed a reliability guide into Karmada."

**Action**
1. **Redefined what "healthy" means:** SLIs and SLOs scoped to tenant experience, not infrastructure health. Hard because they have to be fair across very different tenant workloads.
2. **Solved the cross-cluster gap:** per-cluster Prometheus + Thanos global view; symptom-based alerting across clusters.
3. **Burn-rate alerting:** alerts on SLO burn-rate violations, documented in a reliability guide. Calibrated SLO thresholds `[how: e.g., against production history]`.
4. **Runbooks with every alert** (e.g., `MemberClusterUnreachable`) so whoever is paged starts with a path to resolution.
5. **Gave tenants ownership:** `[how tenants now see and act on their own signals]`.

**Result:** `[before/after: alert volume, pages per week, MTTD/MTTR, if shareable]`; tenants fix their own issues; Karmada reliability engineering guide and observability docs upstream.

**Learning:** "SLOs are never done. Fairness across very different workloads needs continuous refinement, so I treat thresholds as something to recalibrate on a schedule, not set once."

---

### Story #7: Zero-downtime ingress-nginx → Istio migration

**Signals:** execution · risk management · delivering on live infrastructure
**Use for:** risky migration · complex delivery · change with no downtime allowed · project under pressure

**Headline:** "I led migrating ingress for our managed control planes from ingress-nginx to Istio with zero downtime for tenants."

**Needs from you:**
- `[Why migrate? (ingress-nginx's retirement/maintenance status, features needed for unified ingress, security?)]`
- `[Scale: how many control planes, how many tenants affected]`
- `[The plan: parallel run? per-control-plane cutover? DNS weighting? rollback path?]`
- `[What went wrong or almost went wrong]`
- `[How you communicated with tenants]`
- `[Result and learning]`

**Shape to aim for:** stakes (every tenant's API endpoint) → options (big bang in a window vs incremental cutover) → your plan and rollback → the scary moment → zero downtime → the learning (`[e.g., "the rollback path is the plan; the cutover is the easy part"]`).

You've given talks on zero-downtime platform migrations in multi-cluster environments, so mine your talk material for the details.

---

### Story #8: Scalability testing harness (ClusterLoader2 + KWOK)

**Signals:** engineering rigor · building leverage · evidence over intuition
**Use for:** built tools for others · proactive problem-solving · data-driven decisions · raising quality

**Headline:** "Before onboarding tenants with very different scale profiles, I wanted evidence instead of intuition, so we built a declarative scalability testing harness on ClusterLoader2 and KWOK."

**30 seconds:**
> "Tenants come with very different scale profiles, and we needed to know how the platform behaves before they commit. We built a harness on ClusterLoader2, the tool Kubernetes scalability SLOs are measured with, and KWOK to simulate hundreds of nodes per member cluster cheaply. Then we turned the tedious manual lifecycle into two CRDs: one `kubectl apply` builds an environment like 20 member clusters of 500 simulated nodes across two regions. We catch scaling regressions before production, and we're taking it upstream."

**Key decisions:** KWOK fidelity is right because you're measuring propagation across the fleet, not data-plane pod execution · CRD + operator instead of scripts, so the lifecycle (creation, readiness, execution, cleanup) is reproducible · upstream partnership (karmada-io/karmada#7385).

**Result:** repeatable load, head-to-head configuration comparisons, regressions caught before production. `[A specific finding it surfaced?]`

**Learning:** "It's not fully automated yet: an engineer still applies the CRDs by hand. Next is Argo Workflows for continuous test matrices. If I did it again I'd plan the automation layer from the start."

---

### Story #9: Office hours and shared-responsibility lines (saying no)

**Signals:** stakeholder management · setting boundaries · prioritization · product ownership
**Use for:** pushed back on a stakeholder · said no · prioritization · managing many customers · product thinking

**Headline:** "As product owner for a platform with many tenant teams, I set up office hours to build in the open, and drew explicit shared-responsibility lines, like tenants owning their own DR tests and internal forks not being a paved path."

**30 seconds:**
> "With 10+ teams depending on the platform, every team had requests. I ran regular office hours where we shared finished, ongoing, and upcoming work and took feedback, so the roadmap was visible and tenants could influence it. I also drew clear lines: tenants own their own DR tests, and internal forks aren't a paved path. Saying no clearly, with a supported alternative, kept the platform maintainable and kept tenants on it."

**Action**
- **Office hours:** a recurring forum (finished / ongoing / upcoming work), to communicate work and get constructive feedback. Platform built in the open.
- **Boundaries:** tenants own their DR tests; forks aren't supported; tenant-level limits (e.g., a control-plane limit per tenant).
- **One concrete "no":** `[the request, the requester's reasonable need, your reasoning, the alternative you offered, the outcome]` (sample telling in `03-storytelling-framework.md` §7).

**Result:** `[e.g., features shipped from tenant feedback (feature-gate customization, audit config, which appear in office-hours notes), fewer one-off requests, tenants staying on the paved path]`.

**Learning:** "A 'no' lands when it comes with the supported path and a date. When I hear the same request twice, I write the answer down as policy."

---

### Story #10: Hoppscotch, a policy-compliant Postman replacement

**Signals:** initiative · ownership beyond your role · security awareness
**Use for:** saw a problem nobody owned · went above and beyond · bias for action

**Headline:** "When engineers needed a policy-compliant replacement for Postman, I stood up Hoppscotch, and it's now used by 375+ engineers."

**Needs from you:** `[Why Postman became a problem (policy/data concerns)]` · `[Why nobody owned it]` · `[Why Hoppscotch; alternatives]` · `[What you did: deployment, security review, rollout]` · `[How adoption grew to 375+]` · `[Who owns it now, if you handed it off]` · `[Learning]`

**Staff angle:** you noticed a gap outside your remit that affected hundreds of engineers, solved it with an existing open-source tool instead of building, cleared it with security, and `[ideally]` set it up so it didn't depend on you.

---

### Story #11: A real failure (you must choose)

Every loop asks some version of "tell me about a failure." Pick a story that's **real, consequential, mostly your fault, fixed, and changed how you work**. Full guidance is in `08-failure-judgment-and-prioritization.md`.

**Candidates from your record (confirm what's true):**

| Candidate | Where it shows up | What you'd need to own |
|---|---|---|
| Race condition that caused member-cluster state flapping, fixed with a safeguard | Office hours, 2026-05-14 | `[Was it your design/code? Tenant impact? How was it detected? Why didn't testing catch it?]` |
| Two-phase deletion safety mechanism for control planes | Office hours, 2026-02-18 | `[What prompted it? A near-miss or incident where a control plane was deleted or nearly deleted?]` |
| Detector controller concurrent reconciliation triggering SLO violations | Upstream issue karmada-io/karmada#7120 | `[Did you ship the config that exposed it? Or did you catch it? A failure story needs your mistake]` |
| Starting with the standalone topology | Story #2 | `[Did you underestimate the SPOF and have to redo work? Only a failure if it cost something]` |
| Building DR verification late | Profile "needs improvement" | `[Did a gap actually bite you?]` |
| `[Something non-technical: an overcommitment, a stakeholder you lost, a hire or mentee you misjudged]` | | |

**Template:**
- **Headline:** "I `[did X]`, and it caused `[impact]`."
- **What happened:** `[fill]`
- **My part:** `[what you did or didn't do that led to it]`
- **Response:** `[detection, mitigation, communication]`
- **Fix:** `[immediate fix + systemic fix]`
- **Learning:** `[the mechanism you now use]`

---

### Story #12: A real conflict or disagreement (you must choose)

Pick one where the **stakes were real**, the **other party was reasonable**, and the **resolution came from your behavior**. Full guidance and a conflict-story finder are in `07-conflict-disagreement-and-pushback.md`.

**Where to look:** an ARB design review where you disagreed with a team · a tenant team pushing for something the platform shouldn't do · upstream reviewers questioning your design · disagreement on ingress migration timing · who owns alerts or SLO thresholds · a peer team whose roadmap you depended on (managed etcd, token exchange, inventory API) · a disagreement with your manager about priorities.

**Template:**
- **Headline:** "`[Role]` and I disagreed about `[X]`."
- **Their position and why it was reasonable:** `[fill]`
- **My position and why:** `[fill]`
- **What I did:** `[understood incentives, made it about criteria and data, prototype, compromise, escalation if any]`
- **Outcome:** `[who was right; did you disagree and commit?]`
- **Relationship after:** `[fill]`
- **Learning:** `[fill]`

**Have two:** one where you won the argument and one where you **lost or changed your mind**. "Disagree and commit" and "time you were wrong" are both common.

---

### Story #13: An Architecture Review Board case (you must choose)

**Signals:** raising the bar beyond your team · org-level judgment · tactful influence
**Use for:** influenced another team's design · set standards · gave hard feedback to a senior peer

**Needs from you:** `[a design you reviewed (public-safe description)]` · `[the issue you raised, e.g., a DR gap, a single point of failure]` · `[how you delivered the feedback]` · `[what changed]` · `[any standard or guidance that came out of it]` · `[advising teams on DR strategy: a specific instance]`

---

### Story #14: Mentoring a new grad / ColorStack mentee (you must choose)

**Signals:** growing people · mentoring across differences · investing in the pipeline
**Use for:** mentoring someone specific · helping someone find their footing · diversity and inclusion questions

**Needs from you:** `[one person, by role only]` · `[where they started and what they struggled with]` · `[what you did: concrete and repeated]` · `[where they are now]` · `[what you learned about mentoring]` · AfroTech feature on mentorship at Bloomberg (public, can reference).

---

## 3. Spare stories worth drafting

These are smaller, but useful as second examples or for rapid-fire rounds:

| Story seed | Good for |
|---|---|
| Observability docs and Grafana dashboards upstream; structured logging and metric standardization in Karmada | Raising quality across a project; multiplier |
| Conformance testing proposal upstream (karmada-io/karmada#7384) | Setting standards; long-term vision |
| KubeCon NA 2026 Program Committee, Platform Engineering track | Community leadership; judgment about quality |
| Talks at KubeCon and Observability Summit | Communication; representing the company |
| TAG Operational Resilience Tech Lead nomination | Recognition by peers; ecosystem influence |
| Composing internal services instead of rebuilding (part of #1) | Build vs buy judgment |
| Feature-gate and audit-config customization, shipped from tenant feedback | Customer-driven roadmap |

---

## 4. Coverage matrix

Each row is a common competency. **P** = primary story, **b** = backup. Make sure no column carries too much weight; if Story #1 is primary for six rows, find alternatives.

| Competency | 1 | 2 | 3 | 4 | 5 | 6 | 7 | 8 | 9 | 10 | 11 | 12 | 13 | 14 |
|---|---|---|---|---|---|---|---|---|---|---|---|---|---|---|
| Technical direction / biggest decision | P | b | | | | | | | | | | | | |
| Design with tradeoffs | b | P | | | | | | b | | | | | | |
| Influence without authority | b | | P | | | | | | | | | | b | |
| Leading across teams | P | | | b | | | | | | | | | | |
| Ambiguity / no clear owner | P | | | | | | | | | b | | | | |
| Growing others / mentoring | | | | | P | | | | | | | | | b |
| Conflict / disagreement | | | b | | | | | | | | | P | b | |
| Saying no / pushback | | | | | | | | | P | | | b | | |
| Failure / mistake | | | | | | | | | | | P | | | |
| Reliability / operations | | | | | | P | b | | | | | | | |
| Risky delivery / execution | | | | | | | P | | | | | | | |
| Customer / tenant focus | | | | | | b | | | P | | | | | |
| Building leverage / tools | | | | | b | | | P | | b | | | | |
| Initiative / above and beyond | | | | b | | | | | | P | | | | |
| Raising the bar beyond your team | | | | | | | | | | | | | P | |
| Representing the org externally | | | b | P | | | | | | | | | | |
| Data-driven decision | | | | | | b | | P | | | | | | |
| Prioritization | | | | | | | | | P | | | | | |

---

## 5. Choosing stories by company

| Company type | Lead with | Why |
|---|---|---|
| AI-infra / GPU cloud (e.g., CoreWeave) | #2, #6, #1, #8 | Control planes as a product, ops maturity, scale evidence |
| Large tech with leveling committees (Google, Meta, Amazon) | #1, #3, #5, #12, #11 | Broad competency coverage; influence and growth weigh heavily |
| AI labs | #1, #10, #11, #5 | Judgment, low ego, initiative, honesty about limits |
| Infra/devtools companies with open-source DNA | #3, #4, #5, #8 | Upstream influence and community leadership |
| Startups | #10, #1, #7, #9 | Ownership, speed, scrappiness, product sense |

Company-style differences are in `01-rounds-and-how-youre-evaluated.md`.

---

## 6. Worksheet (copy one per story)

```markdown
### Story #__: <title>

Headline (1 sentence):
30-second version:

S: where / stakes / why hard
T: what I owned
A: decision 1 (why, rejected alternative)
   decision 2 (why, rejected alternative)
   decision 3 / who I brought along and how
R: number + business meaning
L: concrete mechanism I now use

Tension: skeptic / constraint / risk / wrong turn
"I" check: three decisions that were mine:

Follow-ups:
- What did you specifically do?
- Hardest part?
- Who disagreed and how did you handle it?
- Alternatives rejected?
- How did you measure success?
- What would you do differently?
- What happened afterward?
- What would <skeptic/manager/tenant> say?

Confidentiality check: anything here I can't say publicly?
Questions this answers:
```

---

## 7. Self-check

1. Are Stories #11 and #12 filled in with real events? (If not, stop and do that first.)
2. Could you tell Stories #1, #2, #5, #6, #11, #12 **without notes** at 30 seconds and 2–3 minutes?
3. Does every story pass the six bank-worthy tests in §1?
4. Does every competency row in §4 have **two** stories?
5. Have you removed every bracket from the stories you plan to use this week?
