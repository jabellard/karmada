# 06. Growing and Developing Others

At Staff, "how do you grow others?" isn't a soft question. It's a scope question. A Staff engineer's output is measured partly by what the people around them can do that they couldn't before. Engineering managers in particular will probe this hard, because they want a Staff IC who makes their team stronger without trying to manage it.

You have real material: Bloomberg's Karmada contributor base went from fewer than 5 to 20+, mentees went on to drive design proposals and scheduler enhancements, you ran a Contribufest with 20+ engineers, you mentor new-grad hires and ColorStack students, and you built forums (office hours, docs, guides) that let other people work without you. The job of this module is to frame all of that as **trajectories you created**, not activities you did.

**The Staff framing rule:** "I mentored engineers" is Senior. "I grew contributors from fewer than 5 to 20+, and people I mentored now drive designs" is Staff. Always describe where the person or group ended up.

---

## At a glance

| Question family | What they're really testing | Best stories |
|---|---|---|
| How do you grow others? (philosophy) | Do you have a repeatable mechanism, or just good intentions? | #5 (contributors <5 → 20+) |
| Mentoring a specific person | Can you describe a real trajectory with your specific interventions? | #5 (a mentee), #14 (new-grad / ColorStack) |
| Giving hard feedback | Can you be direct and kind, and does it change behavior? | `[supply]`, #5 (code/design review) |
| Receiving feedback | Are you coachable? | `[supply]` |
| Struggling teammate / underperformer | Do you diagnose before judging, and know your lane vs. the EM's? | `[supply]` |
| Delegation and letting go | Do you hand off work you could do faster? | #5, #3 (operator sub-project) |
| Raising the bar | Do you scale quality through review, docs, standards? | #6 (reliability guide), #13 (ARB), #8 |
| Growing a community | Can you build a pipeline, not just help individuals? | #5 (Contribufest, <5 → 20+) |
| Hiring and onboarding | Do you build the team's future? | `[supply: interviewing at Bloomberg]`, #5 (contributor onboarding) |
| Culture | What norms do you create? | #9 (office hours, build in the open), #6 |
| Sponsorship vs. mentorship | Do you spend your own credibility on others? | #5 (mentees driving scheduler designs) |
| Successor / not being the bottleneck | Would the work survive if you left? | #5, #2 (operators, self-service) |
| Mentoring across differences | Can you meet people where they are? | #14 (ColorStack, new-grads) |

---

## What Staff looks like here

| Dimension | Senior answer | Staff answer |
|---|---|---|
| Unit of growth | One person I helped | A pipeline: first PR → bigger scope → design ownership, repeated |
| Evidence | "They said it helped" | "They now drive scheduler enhancements on their own" |
| Method | Answering questions | Giving context and ownership, then stepping back |
| Delegation | Hands off tasks | Hands off problems, including ones they'd do faster themselves |
| Bar raising | Careful code reviews | Standards, guides, and forums that raise quality without them in the room |
| Feedback | Avoids or softens it | Specific, timely, about behavior and impact, followed up |
| Relationship to EM | Doesn't mention them, or acts like a shadow manager | Partners with the EM on growth; stays in the technical lane on performance |
| Success test | "I helped a lot" | "The work and the people would be fine if I left" |

**Anti-signals:**
- "I'm always available for questions." (You're a bottleneck, not a multiplier.)
- Only ever mentoring people who were already doing well.
- Taking over when a mentee struggles.
- Describing performance management as if you were their manager.
- No real example of feedback you received and acted on.
- Measuring mentoring in hours spent rather than outcomes.

---

## 1. How do you grow and develop others? (philosophy + mechanism)

**Common phrasings:**
- "How do you help engineers grow?"
- "Tell me about your approach to mentoring."
- "How do you make the people around you better?"
- "How do you develop junior engineers?"
- "What does being a multiplier mean to you?"
- "How do you level up a team?"
- "Tell me about a time you grew the capabilities of a team."

**What they're really asking:** Do you have a *mechanism* that works repeatedly, and evidence it produced results?

**Use:** Story #5 (contributors <5 → 20+). Backup: Story #14 (new-grad mentoring).

**Answer skeleton (~2.5 min):**
1. **Philosophy in one line (10 s):** context over tasks; ownership over help.
2. **The mechanism (50 s):** the ladder you move people along.
3. **The evidence (40 s):** <5 → 20+, mentees driving designs, Contribufest.
4. **One specific person (30 s):** a short trajectory.
5. **What you learned about doing it well (15 s).**

**Your growth ladder (name it in the answer):**
| Stage | What you give them | What they own | Your role |
|---|---|---|---|
| 1. First contribution | A well-scoped issue, the "why," pointers to code | A small PR | Fast, kind, detailed review |
| 2. Real feature | The problem and the constraints, not the design | A feature end to end | Review the approach early, then step back |
| 3. Design ownership | A problem area with ambiguity | A design proposal through community review | Sponsor: pre-wire reviewers, back them in review |
| 4. Teaching others | Visibility and credit | Reviewing and mentoring newer contributors | Get out of the way |

<details>
<summary>Sample answer: Story #5</summary>

> My approach is to give people context and ownership rather than tasks, and to move them up a deliberate ladder: first PR, then a real feature, then owning a design, then mentoring the next person.
>
> The clearest example is Bloomberg's contributions to Karmada. When I started, fewer than five Bloomberg engineers contributed upstream. I made the first step easy: I picked well-scoped issues that mattered to our platform, explained *why* they mattered, and reviewed quickly and in detail so the first experience was good. Once someone had a couple of PRs in, I'd give them a problem rather than a solution, a feature where they had to make the design calls, and I'd review the approach early and then step back. Once they'd done that, I'd encourage them to take on a design proposal and I'd help them through community review.
>
> To widen the funnel, I ran a Contribufest with 20+ engineers, a hands-on session to get people through their first contribution.
>
> Today Bloomberg has 20+ contributors, and engineers I mentored went on to drive design proposals and Karmada scheduler enhancements on their own. [Supply one name-free example: "One engineer started with a [docs/bug fix] PR and within [N months] was leading [a scheduler enhancement proposal]."]
>
> What I learned: the biggest lever is the first review. A slow or harsh first review loses a contributor; a fast, generous one usually keeps them.

</details>

**Likely follow-ups:**
| Follow-up | How to handle |
|---|---|
| "Walk me through one person in detail." | Cluster 2. Have one trajectory fully prepared. |
| "How do you pick what to give someone?" | Slightly beyond what they've done, important enough to matter, reversible enough that a mistake is survivable. |
| "How do you know it's working?" | They make good calls without checking with you; they start mentoring others. |
| "What about people who didn't grow?" | Be honest: some contributors drift off. You learned to check motivation early and match work to what they care about. |
| "How do you find time for this?" | It's part of the job, not extra: every hour that makes someone independent saves many later. |
| "Did you do this inside your team, not just open source?" | `[supply: new-grads on your team, tenant engineers through office hours]`. Have an internal example. |

**Traps:** a generic philosophy with no numbers; only open-source examples if the interviewer is an EM asking about teams.

---

## 2. Mentoring a specific person (the trajectory story)

**Common phrasings:**
- "Tell me about someone you mentored. How did they grow?"
- "Tell me about a time you helped someone get promoted / take on more scope."
- "Tell me about the person you're proudest of having helped."
- "Tell me about a time you coached someone through a hard problem."

**What they're really asking:** Can you describe a person's before and after, and the specific moves *you* made?

**Use:** a mentee from Story #5 who now drives scheduler work (`[supply details, no name needed]`). Backup: Story #14 (`[a new-grad or ColorStack mentee]`).

**Answer skeleton (~2 min):** where they started and what they wanted (20 s) → what held them back (20 s) → your 2–3 specific interventions (50 s) → where they are now (20 s) → what you learned (10 s).

<details>
<summary>Sample answer template (fill in the brackets)</summary>

> One engineer I mentored started as [their starting point: e.g., a new contributor with strong Go but no Kubernetes controller experience]. They wanted [their goal: e.g., to own something visible upstream].
>
> What was holding them back was [diagnosis: e.g., they'd wait for a complete spec before starting, and they were nervous about public review].
>
> I did three things. First, I [intervention 1: e.g., gave them the problem and the constraints instead of a design, and we talked through options together for one session, then I stepped back]. Second, I [intervention 2: e.g., reviewed their first design draft privately before it went public, so their first public review went well]. Third, I [intervention 3: e.g., pointed maintainers to their work and had them present it in the community meeting].
>
> Within [timeframe], they [outcome: e.g., were leading a scheduler enhancement proposal through review and reviewing newer contributors' PRs].
>
> What I learned: [e.g., confidence was the bottleneck, not skill. The private-review-first step mattered more than any technical advice].

</details>

**Likely follow-ups:** "What did you do when they got stuck?" (ask questions before giving answers) · "Did you ever disagree with their design?" (let them make reversible mistakes; step in on irreversible ones) · "How did you give them visibility?" (cluster 10, sponsorship) · "Are you still in touch?"

**Traps:** a hero story where you did the work; a mentee with no clear "after."

---

## 3. Giving hard or constructive feedback

**Common phrasings:**
- "Tell me about a time you gave someone difficult feedback."
- "Tell me about feedback you gave a peer (or someone more senior)."
- "How do you give feedback in code review / design review?"
- "Tell me about a time you had to tell someone their work wasn't good enough."
- "How do you handle a teammate who keeps making the same mistake?"

**What they're really asking:** Are you direct enough to help, kind enough to be heard, and do you follow up?

**Use:** `[supply a real instance: a contributor's design, a peer's approach, a tenant team's architecture via ARB]`. Backup: Story #13 (ARB review that flagged a DR risk).

**Your method (say it):**
1. **Private, soon, specific.** The behavior, its impact, and what good looks like.
2. **Ask before telling.** "How do you think that went?" Often they already know.
3. **Separate the work from the person.**
4. **Agree on a next step.** Then follow up.
5. **In reviews:** distinguish blocking issues from preferences; explain *why*; praise what's good, specifically.

<details>
<summary>Sample answer template</summary>

> [Situation: e.g., a contributor's PR / a peer's design / an ARB submission] had [the issue: e.g., a failover design that assumed both data centers would never be degraded at the same time].
>
> I raised it [privately / in a 1:1 / in review with a direct but specific comment] rather than [in a big meeting]. I said what I saw, why it mattered (the impact on [tenants / the DR guarantee]), and asked how they'd thought about it. [Their reaction.]
>
> We agreed [next step], and I [followed up: e.g., reviewed the revised design]. The result was [outcome].
>
> What I learned: [e.g., the earlier the feedback, the cheaper it is. Now I review designs at the outline stage].

</details>

**Likely follow-ups:** "How did they react?" · "What if they'd disagreed?" · "Feedback to someone senior to you?" (same method, more emphasis on impact and asking questions) · "Feedback that didn't land?" (be honest; what you changed).

**Traps:** feedback so softened it's not feedback; feedback delivered in public.

---

## 4. Receiving feedback

**Common phrasings:**
- "Tell me about the toughest feedback you've received."
- "Tell me about a time you got feedback you disagreed with."
- "What's a piece of feedback that changed how you work?"
- "What would your manager say you need to improve?"

**What they're really asking:** Are you coachable, and do you actually change?

**Use:** `[supply a real one]`. This is commonly asked and you don't have it written yet. Good candidates: feedback about delegating more, about communicating status upward, about over-investing in perfection, or about how you come across in review.

**Answer skeleton (~1.5–2 min):** the feedback, quoted as directly as you can (15 s) → your honest first reaction (15 s) → what you did to check it (20 s) → the specific change (30 s) → evidence the change stuck (20 s).

> Good shape: "My manager told me [feedback]. My first reaction was [honest]. I asked [two people] whether they saw it too, and they did. I changed [specific practice]. Since then [evidence]."

**Traps:** a disguised strength ("I care too much"); no evidence of change.

---

## 5. Helping a struggling teammate or underperformer

**Common phrasings:**
- "Tell me about a time you helped a struggling teammate."
- "How have you dealt with an underperformer on your team?"
- "Tell me about a teammate who wasn't pulling their weight."
- "What do you do when a project is at risk because of one person?"

**What they're really asking:** Do you diagnose before judging? Do you know where your lane ends and the manager's begins?

**Use:** `[supply]`. Backup framing from Story #5: a contributor who stalled.

**Your approach:**
1. **Diagnose first.** Skill, clarity, motivation, or something outside work? Each has a different fix.
2. **Fix what's in your control.** Clearer scope, pairing, smaller milestones, unblocking.
3. **Make expectations explicit.** Ambiguity is the most common cause of "underperformance."
4. **Bring in the EM early, as a partner.** Performance management is their job; your job is technical support and honest input.
5. **Protect the project and the person.** Re-scope the work rather than silently doing it yourself.

<details>
<summary>Sample answer template</summary>

> [Context: a teammate / contributor was missing milestones on X.] Before assuming anything, I talked with them one-on-one. It turned out [diagnosis: e.g., the scope was ambiguous and they didn't feel able to ask, or they were new to controllers].
>
> I [fix: e.g., broke the work into smaller milestones with clear done criteria and paired with them on the first one]. I also let our EM know what I was seeing and what I was trying, so they had context and could support from their side. I didn't take the work over.
>
> [Outcome: e.g., they delivered the next milestone on time, and within [timeframe] were working independently.]
>
> Learning: [e.g., what looks like a performance problem is often a clarity problem].

</details>

**Traps:** quietly redoing their work; talking like their manager; never involving the EM.

---

## 6. Delegation and letting go

**Common phrasings:**
- "Tell me about a time you delegated something you could have done faster yourself."
- "How do you decide what to delegate?"
- "Tell me about a time you let someone make a mistake."
- "How do you avoid being a bottleneck?"

**What they're really asking:** Do you scale yourself, or are you the hero everything routes through?

**Use:** Story #5 (handing design work to mentees). Backup: platform work you've handed to team members `[supply]`.

**How you decide (short framework):**
| Delegate it if… | Keep it (for now) if… |
|---|---|
| It's a growth step for someone | It's irreversible and high-risk, and no one is ready |
| It's reversible | It needs context only you have (then write that context down) |
| Someone else will own it long term | It's a short, urgent incident fix |

> Good line: "I delegate the problem, not the task, and I agree up front which decisions they make alone and which we check together."

<details>
<summary>Sample answer template</summary>

> [Feature or area, e.g., a scheduler-related enhancement upstream / an operator feature on the platform] was something I could have done in [a week]. I gave it to [a mentee / teammate] because [they were ready for design ownership, and I wanted it to have an owner other than me].
>
> I gave them the problem and constraints, agreed which calls were theirs and which we'd check together, and reviewed their approach early. It took [longer] than if I'd done it. [They made a call I wouldn't have: describe; it was reversible, so I let it go / it turned out fine / they learned from it.]
>
> Outcome: [they own that area now / they went on to drive design proposals]. The time cost paid back many times over.

</details>

**Traps:** delegating only grunt work; delegating and then disappearing.

---

## 7. Raising the team's bar

**Common phrasings:**
- "How do you raise the engineering bar on your team?"
- "How do you approach code review?"
- "Tell me about a standard or practice you introduced."
- "How do you improve quality without slowing people down?"
- "How do you make good engineering practices stick?"

**What they're really asking:** Does quality improve because of systems you built, or only when you review personally?

**Use:** Story #6 (reliability guide: SLIs, SLOs, runbooks; metrics reference; Grafana dashboards; structured logging and metric standardization upstream), Story #13 (ARB), Story #8 (scalability harness), conformance testing upstream.

<details>
<summary>Sample answer</summary>

> I raise the bar by making the right thing the default and making expectations explicit.
>
> On our platform, we moved from noisy, infrastructure-centric alerts to SLOs scoped to tenant experience and symptom-based alerting. I wrote a reliability guide documenting the SLO foundations and the burn-rate alerting framework, so anyone on the team, or a tenant, can understand why an alert fires and what to do. Upstream in Karmada I did the same: a reliability engineering guide, metrics reference docs, Grafana dashboards, and metric standardization and structured-logging adoption aligned with SIG Instrumentation conventions.
>
> For correctness at scale, we built a scalability harness with ClusterLoader2 and KWOK so performance claims are measured, not argued. Upstream I'm driving conformance and scalability testing efforts modeled on Kubernetes.
>
> And in review, at both code and design level (including the Architecture Review Board), I separate blocking issues from preferences and always explain the why, because the why is what people carry into the next design.

</details>

**Likely follow-ups:** "How do you keep review from becoming a bottleneck?" (grow other reviewers; checklists; automate what you can) · "What do you look for in a design review?" (failure modes, blast radius, migration, operability, alternatives) · "How did the team react to new standards?"

---

## 8. Growing a community or pipeline

**Common phrasings:**
- "Tell me about building a community or a group of practitioners."
- "How did you grow contributors / adoption?"
- "Tell me about a program you started."

**Use:** Story #5. Also tenant office hours (Story #9) as a community of users, and Slack support upstream.

**Points to hit:** a funnel, not a favor. Lower the entry cost (curated issues, Contribufest with 20+ engineers), make the first experience good (fast reviews), create a path upward (feature → design → reviewer), make it visible (credit in community meetings). Result: fewer than 5 → 20+ contributors; mentees driving designs. Upstream, you also help users troubleshoot on Slack and turn recurring questions into docs, which grows the user community as well.

---

## 9. Hiring, interviewing, and onboarding

**Common phrasings:**
- "What's your experience interviewing and hiring?"
- "What do you look for in a candidate?"
- "How do you onboard new engineers?"
- "How do you get a new hire productive quickly?"

**Use:** `[supply: your interviewing experience at Bloomberg, e.g., number of loops, rounds you run]`. Onboarding: contributor onboarding (Story #5), new-grad mentoring (Story #14), docs and guides that let people self-serve.

**What you look for (have a crisp answer):**
> I look for three things: how they reason when they don't know the answer, whether they think about operating what they build (failure modes, observability), and whether they make people around them better. I care more about the trajectory than the starting point.

**Onboarding (have a crisp answer):**
> A well-scoped first change in week one so they ship early, written context (architecture docs, runbooks) so they're not dependent on me, a named buddy, and an increasing-scope plan for the first 90 days. The test is how quickly they can make a good decision without asking.

---

## 10. Sponsorship vs. mentorship

**Common phrasings:**
- "What's the difference between mentoring and sponsoring?"
- "Tell me about a time you advocated for someone."
- "How have you helped someone get visibility?"

**What they're really asking:** Do you spend your own credibility on other people's careers?

**Answer:**
> Mentorship is advice; sponsorship is spending your credibility on someone. Mentoring tells them how to write a good design proposal. Sponsoring is putting their name forward to own it, pre-wiring the maintainers so they get a fair hearing, and making sure they get the credit when it lands. [Supply a real instance: e.g., recommending a mentee to lead a scheduler enhancement, nominating someone to present, writing input for their performance review.]

**Traps:** saying you've sponsored without a concrete instance.

---

## 11. Growing a successor / not being the bottleneck

**Common phrasings:**
- "If you left tomorrow, what would happen to your project?"
- "How do you make sure knowledge isn't concentrated in you?"
- "Tell me about a time you made yourself unnecessary."

**What they're really asking:** Did you build something that lasts?

<details>
<summary>Sample answer</summary>

> It would keep running, and that's by design. The platform is self-service: tenants declare control planes as custom resources and operators reconcile the rest, so day-to-day operation doesn't go through me. Operations are documented: the reliability guide, runbooks, and alerts tied to tenant-facing SLOs mean on-call doesn't need me. Tenants own their DR tests under our shared-responsibility model. And upstream, there are now 20+ Bloomberg contributors, several of whom drive designs, so our upstream influence doesn't depend on me either.
>
> [Supply: who on the team could step into the tech lead role, and what you've done to prepare them, e.g., having them lead office hours or design reviews.]

</details>

---

## 12. Mentoring across differences

**Common phrasings:**
- "How do you mentor people with different backgrounds from you?"
- "Tell me about supporting someone early in their career."
- "How do you contribute to an inclusive team?"

**Use:** Story #14 (Bloomberg × ColorStack mentoring, new-grad mentoring, AfroTech feature on mentorship at Bloomberg).

**Points to hit:** start by asking what they want rather than assuming; make the unwritten rules explicit (how decisions really get made, how to ask for help, how to be visible); connect them to people and opportunities, not just advice. `[Supply one concrete example and outcome.]`

> Good line: "A lot of what helps early-career engineers isn't technical. It's the unwritten rules. I try to make those explicit."

---

## The EM lens: what managers want to hear from a Staff IC

Engineering managers interviewing you are asking: *will this person make my team better, and will they make my job easier or harder?*

| They want to hear | Say it like |
|---|---|
| You partner with the EM on growth | "I share what I'm seeing about people's technical growth with the EM and ask where they want me to invest." |
| You stay in your lane on performance | "Performance and career decisions are the EM's. I give honest technical input and support the plan." |
| You create stretch opportunities | "I look for problems that are the right size for someone's next step and help scope them." |
| You give the EM signal for reviews and promotions | "I'm happy to write specific, evidence-based input for reviews and promo cases." |
| You don't hoard the interesting work | "I take the riskiest, most ambiguous piece and hand off the rest as growth opportunities." |
| You make the team less dependent on you | "My test is whether things go well when I'm on vacation." |
| You raise problems early | "If I see someone struggling or a team dynamic going wrong, the EM hears it from me early, privately." |

**Questions to ask an EM (they signal you care about people):**
- "How do Staff engineers and managers split responsibility for people growth here?"
- "What does the team's seniority mix look like, and where do you most need a multiplier?"
- "Who on the team is ready for more scope, and how could a Staff engineer help?"

**Avoid:** "I've basically been managing the team" (sounds like you want the EM's job), or "I prefer to focus on code" (sounds Senior).

---

## Rapid-fire variants

| If they ask… | Cluster | Story |
|---|---|---|
| "What's your leadership style with junior engineers?" | 1 | #5 ladder |
| "How do you balance mentoring with your own deliverables?" | 1 / 6 | #5: mentoring is how deliverables scale |
| "Tell me about a time you taught a complex topic." | 7 | Reliability guide, Contribufest |
| "How do you handle someone who rejects your feedback?" | 3 | Ask why, look for the valid part, escalate only on real risk |
| "Tell me about building psychological safety." | Culture | Office hours "built in the open"; blameless reviews `[confirm]` |
| "How do you run a design review?" | 7 | ARB (#13) |
| "How do you help tenants / customers become self-sufficient?" | 7 / 11 | #6 (tenants fix their own issues), #9 |
| "Who has mentored you, and what did you learn?" | 4 | `[supply]` |
| "Tell me about recognizing someone's work." | 10 | Credit in community meetings `[supply]` |
| "What's one thing you'd want new engineers to learn first?" | 9 | "How to operate what they build" |
| "How do you handle a senior engineer who resists change?" | 3 / see 07 | Find their incentive; data; pilot |
| "Have you ever had to let someone go / recommend it?" | 5 | Your lane: honest input to the EM; `[supply if real]` |
| "How do you build a team culture?" | Culture | Building in the open (office hours), reliability as a shared value, fast generous reviews |

---

## Self-check

1. Can you describe one mentee's before and after in under 90 seconds, with three specific things you did?
2. Can you name your growth ladder's four stages without notes?
3. Do you have a real "feedback I received and acted on" story with evidence it stuck?
4. Do you have a real "struggling teammate" story, and does it include the EM?
5. Do you have at least one internal (non–open source) growth example for EM interviewers?
6. What's one concrete instance of sponsorship, not mentorship?
7. If you left tomorrow, who would step in, and what have you done to prepare them?
8. Is every number you'd say (fewer than 5 → 20+, Contribufest 20+) something you can back up?

**Drill (25 min):** answer clusters 1, 2, and 4 out loud (2 minutes each), then have a friend ask three follow-ups from the tables. Score each answer: (a) a trajectory, not an activity, (b) your specific interventions, (c) an outcome, (d) a learning. Then fill every `[bracket]` you hit and log new stories in 04-story-bank.md.
