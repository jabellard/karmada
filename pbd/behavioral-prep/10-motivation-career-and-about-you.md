# 10. Motivation, Career, and "About You" Questions

These questions show up in every round, from the recruiter screen to the director chat, and they're often the first and last thing an interviewer hears from you. They sound soft, but they decide your level ("does this person think at Staff scope?"), your fit ("why us?"), and your risk ("will they leave in a year?"). Prepare them as carefully as your stories; most candidates wing them and it shows.

This module aligns with `../interview-playbook.md` §1.1–1.3 and the CoreWeave intro in `../coreweave-cks-em-intro.md` §3. If you change your pitch, change it in all three places.

---

## At a glance

| Question family | What they're really testing | Use |
|---|---|---|
| Tell me about yourself | Can you frame your own career crisply, at the right level, for this audience? | Pitch (30s / 60s / 2 min) |
| Walk me through your resume | Is there a coherent arc and growing scope? | Arc script |
| Why leaving / why now | Are you running toward something, or away from something? | Positive "next scale" answer |
| Why this company / role | Did you do homework, and does your experience map to their problem? | Formula + worksheet |
| What you're looking for | Will this role actually satisfy you? | Three must-haves |
| IC vs. management | Do you understand the Staff path, and are you a flight risk? | IC-by-choice answer |
| Career goals | Ambition with realism; does it fit here? | 3–5 year answer |
| Strengths / weaknesses | Self-awareness, honesty, growth mechanism | Evidence-backed strengths; real weakness with mechanism |
| What would others say / feedback | Self-awareness, coachability | Real feedback + change |
| Motivation / culture / manager | Fit and how to manage you | Energizers, drains, working style |
| Staying current | Learning habits; depth in the field | Maintainer, program committee, talks |
| Proudest accomplishment | Values: what you think "impact" means | Story #1 or #5 |
| Level conversation | Can you advocate for yourself calmly? | Staff scope mapping |
| Questions for them | Curiosity, judgment, seriousness | Question bank by interviewer |
| Closing | Last impression, interest | Closing scripts |

---

## What Staff looks like here

| Dimension | Senior answer | Staff answer |
|---|---|---|
| Pitch | Job title and tech stack | Scope, the problem you own, outcome, and where you want to go |
| Why leaving | Vague or negative | "I took X to Y; the next scale of this problem is here" |
| Why this company | "Great culture, cool tech" | "Your problem is Z; I've solved a version of Z; here's where it'd stretch me" |
| Goals | "Grow technically" | "Own a hard area end to end and set direction across teams" |
| Weakness | Humblebrag ("I care too much") | Real gap, how you found it, what you're doing about it, evidence it's working |
| Questions | Perks, process | The team's hardest problem, how Staff engineers shape direction |

**Anti-signals**
- Reciting your resume chronologically for 4 minutes.
- Any criticism of Bloomberg, your manager, or your team.
- "I'm open to anything." Staff candidates know what they want.
- Sounding desperate or sounding bored.
- A weakness that's secretly a strength, or one that's disqualifying for the role.
- No questions at the end.

---

## 1. Tell me about yourself

**Common phrasings**
- "Tell me about yourself."
- "Walk me through your background."
- "Give me the quick version of your story."
- "What should I know about you?"
- "Introduce yourself."

**What they're really asking:** In 60–90 seconds, who are you professionally, what level do you operate at, and why are you talking to me?

**Formula:** **Present** (role + scope + one outcome) → **Proof** (one or two signature achievements) → **Beyond the day job** (upstream, community) → **Future** (why this conversation). End with a hook, not "...so yeah."

### 30-second version (recruiter opener, networking)

> I'm a platform engineer at Bloomberg, where I'm the tech lead and product owner for the firm's managed multi-cluster Kubernetes platform. It runs 40+ cluster fleets for AI/ML, data and analytics, and general compute, and it took standing up a DR-compliant control plane from weeks to minutes. I'm also a maintainer of Karmada, the CNCF multi-cluster project. I'm looking for the next scale of fleet and infrastructure problems, especially under AI workloads.

### 60-second version (default)

> I'm a platform engineer at Bloomberg, where I'm the tech lead and product owner for the firm's managed multi-cluster Kubernetes platform. I led Bloomberg's adoption of Karmada as the firm's standard for meeting its DR and HA mandate, designed the platform, and took it from proof of concept to production, working with over 10 core platform and infrastructure teams. It now runs 40+ Kubernetes cluster fleets across AI/ML, data and analytics, and general compute, and standing up a DR-compliant control plane went from weeks to minutes.
>
> I'm also a maintainer of Karmada and tech lead of its operator sub-project. I helped drive its CNCF graduation and grew Bloomberg's contributor base from fewer than 5 to over 20.
>
> I'm looking for the next scale of fleet and infrastructure problems, especially the infrastructure under AI workloads, where scheduling, capacity, and failover are the hard problems. Happy to go deeper on any of that.

### 2-minute version (HM screen, director chat)

Add one sentence to the 60-second version for each of these:
1. **Architecture:** "Each control plane is stretched across data centers on a multi-data-center etcd cluster behind one endpoint, and tenants declare control planes as custom resources that our operators reconcile."
2. **Operations:** "We moved from noisy infrastructure alerts to SLOs scoped to tenant experience, so tenants now fix their own issues."
3. **People:** "Upstream, my mentees now drive design proposals and scheduler enhancements."
4. **Community:** "I'm on the KubeCon NA 2026 program committee for platform engineering, and I speak at KubeCon on DR, observability, and migrations."

### Tailoring by audience

| Audience | Lead with | Close with |
|---|---|---|
| **Platform engineering team** | Paved path; tenants inherit DR; self-service; SLOs | "I want to work on platform problems at a bigger scale" |
| **AI-infra team** | Platform underneath the data science platform (distributed training across heterogeneous clusters); capacity-aware placement; failover | "The hardest fleet problems now sit under AI workloads: scheduling, capacity, failover" |
| **Kubernetes-as-a-product (CKS, GKE, AKS)** | "We run Kubernetes control planes as a service": CR → stretched control plane, etcd, PKI, credentials, registration | "Here it's the core product, at bigger scale and on harder hardware" |
| **EM** | Ownership, outcomes, people (contributors <5 → 20+), product owner role | "I want to own a hard area end to end and set direction" |
| **Director / VP** | Business outcome (DR mandate met by adoption; firm standard) | Where you want to go and why their org |
| **Peer engineer** | Shorter; get to the technical hook fast | "Happy to go deep on the stretched topology" |

**Traps**
- Explaining Karmada internals. To them, Karmada is evidence, not the topic.
- Going over 90 seconds without being asked.
- Omitting the "why I'm here" ending; it sets up "why this company."

---

## 2. Walk me through your resume

**Common phrasings:** "Walk me through your resume." · "How did you get to where you are?" · "Tell me about your career path." · "Why did you make each move?"

**What they're really asking:** Is there an arc of growing scope, and were your moves deliberate?

**Answer skeleton (2–3 min):** one line per earlier stage, focused on *what changed in your scope*, then 60% of the time on the current role.

> I started `[at Bloomberg, or wherever you started]` as `[earlier role / team]`, working on `[area]`. That's where I learned `[one thing]`. Then I moved to `[next role]`, where my scope grew to `[x]`. `[Earlier tech: e.g., C# work, if relevant]`. For the last `[N]` years I've been the tech lead and product owner for the managed multi-cluster platform, `[60-second pitch body]`. Along the way I became a Karmada maintainer. The thread through all of it is `[e.g., taking infrastructure problems teams solve badly on their own and turning them into paved paths]`.

**Fill in:** `[start year, earlier roles, education]` (see `../linkedin-guide.md` checklist for start dates).

**Likely follow-ups:** Why stay at one company so long? (scope kept growing: from `[x]` to firm-wide standard; plus external growth through upstream) · Which move taught you most? · What would you change?

**Trap:** long tenure framed as inertia. Frame it as compounding scope.

---

## 3. Why are you leaving? Why now?

**Common phrasings:** "Why are you looking?" · "Why leave Bloomberg?" · "Why now?" · "What's missing in your current role?" · "What would make you stay?"

**What they're really asking:** Are you running toward something specific? Is there a problem (performance, conflict) behind this?

**Formula:** **Accomplished** (what you built) + **stable point** (why now is a natural time) + **pull** (the problem you want next) + **why it needs a new company** (this is the core product there).

> I took the platform from an idea to the firm's standard for multi-cluster DR, and it's in a stable, mature place. Karmada just graduated, which closes a chapter I'd been driving. I want to take what I've learned to larger-scale fleet problems, especially the infrastructure under AI workloads, at a company where this work is the core product rather than an internal service.

**Follow-ups**
| Follow-up | Guide |
|---|---|
| "What would make you stay?" | "Honestly, the problem I want next isn't the core business at Bloomberg. I have a lot of respect for the team." |
| "What don't you like about your current role?" | Reframe to what you want more of: scale, AI infra, product-as-core. Never criticize. |
| "Are you interviewing elsewhere?" | "Yes, a few places in platform and AI infra. This is one of my top choices because `[specific]`." Mention timelines if it helps speed things up. |
| "Would you take a counteroffer?" | "I'm looking because of the problem, not comp, so a counteroffer doesn't change that." (Decide your real walk-away terms first: playbook §7.) |

**Traps:** any hint of frustration; "I've been there too long"; "I'm bored."

---

## 4. Why this company? Why this role?

**Common phrasings:** "Why us?" · "Why this team?" · "What do you know about what we do?" · "Why this role over others you're considering?" · "What excites you about this position?"

**What they're really asking:** Did you do homework, and is the fit real or generic?

**Formula (3 beats, 45–60s):**
1. **Their problem** (specific, from their talks, blog, product, or job post).
2. **Your mapping** (where you've solved a version of it, with one proof point).
3. **The stretch** (what's new and hard for you there; shows honesty and real motivation).

**Example (CoreWeave CKS):**
> Three reasons. Kubernetes *is* the product here, and the customers are the most demanding AI teams, so the bar for control-plane reliability is as high as it gets. The problems map directly to what I do: hosted control planes, etcd, credential lifecycle, multi-tenant observability, capacity-aware scheduling. And bare metal plus GPUs makes the failure modes much harder: node health, burn-in, topology. That's the part I want to learn.

**Worksheet (fill per company; pull candidates from `../target-companies.md`)**

| Field | Your notes |
|---|---|
| Their core infra problem (1 sentence) | |
| Evidence (talk, blog, product page, job post line) | |
| Your matching story # and proof point | |
| What would stretch you there | |
| CNCF/upstream angle (★ rating in target-companies.md) | |
| One sharp question that proves homework | |
| Why this one over your other options (honest) | |

**Traps:** praising "culture" or "mission" without specifics; reading their About page back to them; claiming expertise in their stack you don't have.

---

## 5. What are you looking for in your next role?

**Common phrasings:** "What are you looking for?" · "What's your ideal role?" · "What matters most to you in your next job?" · "How will you decide between offers?"

**Answer (three must-haves, then one nice-to-have):**
> Three things. First, a hard infrastructure problem I can own end to end, ideally fleet-scale platform or the infrastructure under AI workloads. Second, Staff-level scope: setting technical direction across teams, not just within one. Third, a team where infrastructure is the core product, so the bar is high. A bonus is room to keep contributing upstream, since Karmada is a big part of how I've grown.

**Follow-ups:** "What would make you turn down an offer?" (wrong level, no real ownership, `[location/remote constraints: NYC-based, open to remote; relocation only for exceptional]`) · "Remote preferences?" `[decide beforehand]`.

---

## 6. IC vs. management

**Common phrasings:** "Why IC and not management?" · "Would you consider managing?" · "Do you want to be a manager someday?" · "You lead a lot. Why not be an EM?"

**What they're really asking:** Do you understand the Staff role as leadership through technical direction? Will you get bored and leave for an EM role?

> I've chosen the IC path deliberately. What I enjoy most is the technical direction itself: making the architecture calls, setting standards, and growing people through the work, like taking contributors from first PR to driving design proposals. Staff lets me do that across teams without giving up depth. I've done plenty of people leadership as a tech lead and product owner, so I understand the job, and I'd never say never, but it's not what I'm optimizing for now.

`[Adjust if your real answer is different. Don't fake certainty either way.]`

**Trap:** sounding like you avoid people work. Staff is people work.

---

## 7. Career goals: where do you see yourself in 3–5 years?

**Common phrasings:** "Where do you see yourself in five years?" · "What are your long-term goals?" · "What's next after this role?" · "What do you want to learn?"

> In three to five years I want to be the person a company trusts with its hardest infrastructure direction in fleet management and AI infrastructure: a Senior Staff or Principal-scope engineer whose designs other teams build on, and someone known in the community for it, through upstream work and talks. Concretely, in this role, the first year is about learning your systems and people deeply, then owning `[their hardest problem]`.

**Trap:** a goal that obviously requires leaving them (founding a startup next year, switching fields).

---

## 8. Strengths

**Common phrasings:** "What are your greatest strengths?" · "What are you best at?" · "What's your superpower?" · "Why should we hire you?"

**Formula:** strength → evidence (story #) → why it matters for *this* role. Pick two or three, not seven.

| Strength | Evidence | Why it matters |
|---|---|---|
| Turning hard distributed-systems problems into paved paths | #1, #2: DR as a guarantee tenants inherit; weeks → minutes | Platform teams live or die on adoption |
| Influence beyond my team and org | #3, #4: upstreamed features; CNCF graduation; maintainer | Staff work crosses boundaries |
| Growing people | #5: contributors <5 → 20+; mentees drive designs | Multiplier effect |
| Operational judgment | #6, #7: SLOs scoped to tenants; zero-downtime migration | Reliability is the product for infra |

> My strongest skill is taking a hard infrastructure problem that every team is solving badly on its own and turning it into a paved path people actually adopt. Multi-cluster DR at Bloomberg is the clearest case: teams went from building DR themselves to inheriting it, and setup went from weeks to minutes. The second is influence without authority. I got the features we needed into Karmada upstream rather than forking, and ended up a maintainer.

---

## 9. Weaknesses and growth areas

**Common phrasings:** "What's your biggest weakness?" · "What are you working on improving?" · "Where do you need to grow?" · "What would your manager say you need to work on?" · "What's a skill you lack for this role?"

**What they're really asking:** Self-awareness and honesty, plus a working mechanism for growth.

**Formula:** **Real weakness** (relevant but not disqualifying) → **how you noticed** → **what you're doing** (a mechanism) → **evidence it's working**.

**Candidate A: hands-on GPU / AI-infra depth (use for AI-infra roles; it's the obvious gap, so naming it first is a strength)**
> The honest gap for this role is hands-on GPU infrastructure. I've been the platform underneath our data science platform, so I know what training teams need from the control plane and scheduling layer, but I haven't operated GPU fleets at the node level: InfiniBand, NCCL behavior, burn-in, topology-aware scheduling. I'm closing it deliberately: `[studying DRA, Kueue, JobSet; a mini-project like GPU-capacity-aware placement with Karmada + Kueue; writing about it]`. And I have a track record of ramping fast: I went from Karmada adopter to maintainer.

**Candidate B: `[pick a real behavioral one]`**. Options that are honest and safe for Staff, if true for you:
- *Doing it myself instead of delegating* → mechanism: `[e.g., mentoring others to own components; mentees now drive scheduler designs]`.
- *Saying yes to too many tenant asks early on* → mechanism: shared-responsibility lines, office hours, a roadmap (#9).
- *Over-investing in written detail before getting buy-in* → mechanism: `[short proposals first]`.

Pick one that's real. Interviewers probe ("give me an example"), so have the story.

**Traps:** humblebrags; a weakness that's core to the role ("I don't like ambiguity"); no evidence of improvement.

---

## 10. What would your manager or peers say about you?

**Common phrasings:** "How would your manager describe you?" · "What would your peers say?" · "What would your mentees say?" · "If I called your references, what would I hear?" · "What's a criticism you've gotten?"

**Formula:** two positives with a concrete reason, one growth area that matches section 9.

> My manager would say I own things end to end: the platform, the roadmap, and the relationships with tenant teams, and they don't have to chase me. Peers and tenants would say I'm the person who explains the "why" and writes things down: the reliability guide, the office-hours notes. The constructive feedback I've gotten is `[real feedback, e.g., to delegate earlier]`, and `[what you changed]`.

Brief your references to match this (playbook §6).

---

## 11. Best feedback you've received / how you've grown

**Common phrasings:** "What's the best feedback you've gotten?" · "Tell me about a time you received critical feedback." · "How have you changed as an engineer in the last few years?"

**Skeleton:** the feedback (verbatim if you can) → your first reaction (honest) → what you changed → evidence.

> `[Real example. If none comes to mind, a growth arc works: "The biggest shift was from building the right thing to getting the right thing adopted. Early on I'd design something sound and expect it to win on merit. Driving Karmada adoption across 10+ teams taught me that adoption is its own engineering problem: onboarding automation, office hours, clear contracts."]`

---

## 12. Motivation, energizers, and drains

**Common phrasings:** "What motivates you?" · "What gets you excited about work?" · "What kind of work drains you?" · "Describe your ideal day." · "When were you most engaged at work?"

> I'm most energized when a hard technical problem has a real user at the end of it, like turning DR from a project every team dreads into something they inherit. I also get a lot out of watching someone I've mentored take on something I'd have done myself. What drains me is churn without decisions: long debates where nobody owns the call. I usually respond by writing up the options and a recommendation so we can decide.

**Trap:** naming a drain that's a core part of the job (meetings, reviews, on-call). Pair any drain with how you handle it.

---

## 13. Manager style and culture

**Common phrasings:** "What kind of manager do you work best with?" · "Describe your ideal team culture." · "How do you like to receive feedback?" · "What kind of environment do you thrive in?"

> I work best with a manager who gives me a hard problem and context, then trusts me to own it, and who'll tell me directly when I'm off. I like teams that write things down, make decisions explicitly, and own their production systems. Feedback: direct and early beats polished and late.

**Trap:** describing a manager who sounds like "leaves me alone completely."

---

## 14. How do you stay current?

**Common phrasings:** "How do you keep up with the field?" · "What have you learned recently?" · "What tech trend are you excited about?" · "What are you reading?"

> Mostly by being in the work upstream. As a Karmada maintainer I'm in design reviews and community meetings every week. I'm on the KubeCon NA 2026 program committee for platform engineering, which means reading a lot of submissions about what practitioners are actually doing, and I speak at KubeCon and Observability Summit. Recently I've been going deep on AI infrastructure: `[DRA, Kueue, JobSet, LeaderWorkerSet, topology-aware scheduling: pick what you actually studied]`.

**Have ready:** one recent thing you learned and an opinion on it `[fill in]`.

---

## 15. Proudest accomplishment

**Common phrasings:** "What are you most proud of?" · "What's the achievement you'd want on your tombstone (career-wise)?" · "What's the best thing you've built?"

**Choose by audience:** technical audience → Story #1/#2; EM or values-heavy → Story #5 (people). Your own words in the CoreWeave notes: proud of the community influence more than the code.

> Technically, the platform: DR became something teams inherit instead of build. But honestly I'm prouder of the people side upstream: Bloomberg's Karmada contributor base went from fewer than 5 to over 20, and people I mentored now drive design proposals and scheduler work. That outlasts any feature I write.

---

## 16. The level conversation

**Common phrasings:** "What level are you targeting?" · "Why should we level you at Staff?" · "Bloomberg doesn't have Staff. What's your level there?" · "What's the difference between Senior and Staff, in your view?" · "Would you consider Senior?"

**What they're really asking:** Do you understand Staff scope, and can you back the claim with evidence without sounding entitled?

**Script (say once, calmly, early: ideally the recruiter screen)**
> Bloomberg doesn't have formal levels above Senior, but my scope maps to Staff: I set the firm-wide standard for multi-cluster DR, own the platform end to end as tech lead and product owner, sit on the Architecture Review Board, and maintain a CNCF graduated project. I'd like to be considered for Staff.

**"What's the difference between Senior and Staff?"**
> Senior engineers deliver hard projects well within a team. Staff engineers decide which projects matter, set direction across teams, and make the people and systems around them better. Concretely for me, that's choosing Karmada as the firm standard, getting 10+ teams onto it, and growing contributors.

**"Would you consider Senior?"**
> My goal is Staff, and I think the loop will show that scope. If you level me Senior, I'd want to understand the specific gaps and have a clear path and timeline to Staff. `[Your real answer; see playbook §7.6.]`

**Trap:** arguing; over-justifying; bringing it up in every round (once with the recruiter, once with the HM if natural).

---

## 17. Gaps, transitions, and tenure

**Common phrasings:** "Why have you been at one company so long?" · "Explain this gap." · "Why the switch from `[X]` to infrastructure?"

**Long tenure:**
> My scope kept growing, from `[earlier role]` to leading a firm-wide platform, and I grew externally too, as an upstream maintainer working with engineers from many companies. That gave me a range of perspectives beyond Bloomberg.

`[Add any gap or field-switch explanation you need: short, factual, forward-looking.]`

---

## 18. Outside of work

**Common phrasings:** "What do you do for fun?" · "Tell me something not on your resume." · "What are you passionate about outside work?"

Keep it to 20 seconds, real, and human. `[Fill in: e.g., training/fitness, mentoring with ColorStack, something else]`. Mentoring through Bloomberg × ColorStack (and the AfroTech feature on mentorship) works as a bridge if you want one.

---

## 19. Questions to ask them

Always have 2–3 per round, and never ask what a search answers. A good question signals how you think; pick the ones that fit the interviewer.

**Recruiter**
| Question | Signals |
|---|---|
| "What does the loop look like, and which rounds weigh most for leveling?" | You care about level and process |
| "Is leveling decided before or after the loop? Can I be considered for Staff?" | Calm self-advocacy |
| "What problem is this team hiring to solve?" | Outcome focus |
| "Is there a team-matching phase?" | You know how big-company hiring works |
| "What's the timeline? I have other processes moving." | Urgency without pressure |

**Peer engineer**
| Question | Signals |
|---|---|
| "What's the most painful part of operating your system today?" | Operator mindset |
| "How do design decisions get made here: docs, reviews, one person?" | You care about decision process |
| "What did the last incident teach the team?" | Blameless, learning culture |
| "What would you change about the codebase if you had a quarter?" | Tech-debt awareness |

**Staff / principal engineer**
| Question | Signals |
|---|---|
| "What are you working on that you think is most important, and why?" | Respect, curiosity |
| "How do Staff engineers here influence the roadmap: proposals, staffing, both?" | You understand the role |
| "Where is technical direction unclear or contested right now?" | You want the hard problems |
| "What separates Staff from Senior Staff here?" | Long-term thinking |

**Engineering manager**
| Question | Signals |
|---|---|
| "What's the hardest problem on the roadmap you'd want this person to own?" | Ownership |
| "What would make someone exceptional in this role after 6 and 12 months?" | Bar awareness |
| "How is the team organized, and where do you see someone with my background?" | Fit thinking |
| "How do you support Staff engineers' growth? What does promotion beyond Staff look like?" | Ambition |
| "Does the team contribute upstream? Would maintaining Karmada fit?" | Community value (ask at ★★★ companies) |

**Director / VP**
| Question | Signals |
|---|---|
| "What will the infrastructure org look like in two years, and what's the biggest bet?" | Strategy |
| "What's the biggest risk to that plan?" | Judgment |
| "Where do you need senior ICs to lead that you don't have today?" | You want to fill a real gap |
| "How do you measure whether the platform is succeeding?" | Business outcomes |

**Don't ask:** comp or PTO before the offer stage (leave comp to the recruiter), "what does the company do," or anything on the website.

---

## 20. Closing statements

The last 60 seconds of a round stick. End with interest, one specific link, and next steps.

**After an HM / EM round**
> I've really enjoyed this. What you said about `[specific problem]` is exactly the kind of work I want to do next, and it maps closely to `[your story]`. What are the next steps from here?

**After a technical or behavioral loop round**
> Thanks, this was fun, especially `[specific part]`. Is there anything about my background you'd like me to clarify before we wrap?

(That question surfaces doubts while you can still answer them. Use it when the round felt mixed.)

**After a director / VP chat**
> Thank you. The `[two-year bet]` you described is why I'm excited about this team. I'd love to help make it happen.

**Thank-you note (within 24h):** two or three sentences, one specific thing from the conversation, one line restating interest (and, after the recruiter screen, your level interest).

---

## Rapid-fire variants

| Question | Answer in one line |
|---|---|
| "Describe yourself in three words." | "Owner, multiplier, operator." `[pick yours]` |
| "Why should we hire you?" | Strengths answer + their problem: "I've built and run the version of this problem you're hiring for." |
| "What's your dream job?" | Section 5, aimed at this role |
| "What's your management style?" (IC) | "Lead through direction and context; write it down; grow people through the work." |
| "Are you a generalist or specialist?" | "Deep in fleet/control planes, broad across ops, DR, observability." |
| "What would you do in your first 90 days?" | Learn systems and people → find leverage → early win → set direction |
| "What's something you changed your mind about?" | `[real one: e.g., alerting on component health vs. tenant SLOs]` |
| "What's your biggest professional regret?" | Real, small, with learning (module 08) |
| "Salary expectations?" | "I'd like to understand the range for this level first." (playbook §2) |
| "When can you start?" | `[notice period]`, flexible within reason |
| "Do you have other offers?" | Honest; mention timelines if helpful |
| "Anything else we should know?" | One sentence on Staff scope + interest |

---

## Self-check

1. Can you deliver the 30s, 60s, and 2-minute pitches cleanly, ending with a hook?
2. Do you have a filled worksheet (section 4) for each company in your active pipeline?
3. Is your weakness real, with a mechanism and evidence, and can you give an example on demand?
4. Can you say the level script in under 20 seconds without sounding defensive?
5. Do you have 3 questions ready for each interviewer type?
6. Are all `[brackets]` in your pitch and "why leaving" filled in and consistent with `../interview-playbook.md` and `../coreweave-cks-em-intro.md`?

## Drill (30 min)

1. Record the 60-second pitch three times. Keep the best; trim filler words.
2. Do the platform, AI-infra, and EM variants back to back. Notice which sentence changes.
3. Fill the worksheet for your top two companies.
4. Have a friend ask "why leaving," "weakness," and "why Staff?" in random order with follow-ups. Keep each answer under 60 seconds.
5. Write your thank-you note template now so it takes two minutes after each round.
