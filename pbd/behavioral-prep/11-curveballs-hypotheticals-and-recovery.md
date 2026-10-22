# 11. Curveballs, Hypotheticals, and Recovery

Prepared stories cover most questions. This module covers the rest: the question no story fits, "what would you do if…" hypotheticals, deep drilling meant to test whether a story is real, skeptical interviewers, awkward questions (comp, levels, illegal ones), the AI/GPU gap, and how to recover when an answer goes badly. Staff candidates are judged partly on how they handle the unplanned moments, because that's what the job is.

---

## At a glance

| Situation | Your move | Section |
|---|---|---|
| You blank, or no story fits | Pause, clarify, adapt the nearest story, or answer honestly with how you'd handle it | §1 |
| "What would you do if…?" | Clarify → principles → concrete steps → tie to a real example | §2, §3 |
| "What exactly did *you* do?" and other drilling | Go one level deeper each time; specifics prove it's real | §4 |
| Skeptical or cold interviewer | Stay warm, specific, and brief; don't chase approval | §5 |
| Second failure question, or repeating stories | Have a backup; vary the angle, not the facts | §6 |
| Confidential details | Share the shape, not the specifics | §7 |
| Comp, salary history, levels, illegal questions | Polite deflection scripts | §8 |
| GPU / AI-infra gap | Name it calmly, show what transfers, show how you ramp | §9 |
| Running long or cut off | Headline first; yield gracefully | §10 |
| A bad answer earlier in the round | "Can I revisit something?" | §11 |
| Video or audio problems | Name it, fix it, continue | §12 |
| The last five minutes | Good questions, a short close | §13 |

---

## 1. When you blank or no story fits

Blanking happens to everyone. What the interviewer remembers is whether you handled it calmly.

### The pause script

> That's a good question. Let me take a few seconds to pick the best example.

Then actually take 5–10 seconds. Silence feels long to you and normal to them. Writing a word or two on paper helps you think.

### Four ways forward (in order of preference)

| Option | When | What it sounds like |
|---|---|---|
| **1. Clarify** | The question is broad or ambiguous | "Do you mean a disagreement with a peer or with my manager?" / "Are you more interested in a technical failure or a people one?" |
| **2. Adapt the nearest story** | A story fits 70% | "The closest example is when… It's more about X than Y, but the part that's relevant is…" |
| **3. Use a smaller real example** | No big story fits | "I don't have a dramatic example, but a smaller one is…" A small real story beats a big vague one |
| **4. Be honest, then reason** | You genuinely haven't faced it | "I haven't faced exactly that. Here's how I'd approach it, and the closest thing I've done is…" |

**Never invent a story.** Follow-ups expose fiction within two or three questions, and an integrity concern is much worse than a weak answer.

### Use the story-to-question map

Most "new" questions are old competencies in disguise. Map them back:

| Disguised question | Real competency | Try |
|---|---|---|
| "Tell me about a time you changed your mind." | Learning, openness | #2 (architecture evolution), #13 |
| "When did you go beyond your job description?" | Ownership, initiative | #10 Hoppscotch, #3 |
| "Tell me about working with a difficult person." | Conflict | #12 |
| "When did you have to learn something fast?" | Ramp-up, ambiguity | Karmada adopter → maintainer (#1, #3) |
| "Describe an unpopular decision you made." | Saying no, judgment | #9 (internal forks aren't a paved path) |
| "When did you simplify something complex?" | Leverage, customer focus | Source-of-truth API (#2), #6 |
| "When did you use data to make a decision?" | Judgment | #8 scalability harness, #6 SLOs |

Full map: `04-story-bank.md`.

---

## 2. Hypotheticals: a framework

"What would you do if…" questions test judgment, not memory. Weak answers jump to an action. Strong answers show how you think and then commit.

### The four steps

| Step | What you do | Time |
|---|---|---|
| **1. Clarify** | Ask 1–2 questions that change the answer: stakes, timeline, who's involved, what's been tried | 15–30s |
| **2. Principles** | Say what you're optimizing for: "My priorities would be safety first, then unblocking the team, then the long-term fix" | 15s |
| **3. Concrete steps** | 3–5 specific actions in order, including who you'd talk to and what you'd say | 60–90s |
| **4. Tie to a real example** | "I did a version of this when…" | 30s |

Step 4 is what separates a Staff answer from a textbook one. It turns a hypothetical into evidence.

**Common traps:**
- Escalating first. Staff engineers try to resolve it themselves, then escalate with a recommendation.
- Only the technical answer. Most hypotheticals are about people and alignment.
- Never committing. "It depends" is fine as a start, but finish with "Given that, I'd do X."

---

## 3. Worked hypotheticals

Each has a short answer outline. Expand into your own words and practice out loud.

<details>
<summary><strong>H1. A tenant team wants to bypass your platform and build their own.</strong></summary>

**Clarify:** Why? Missing feature, performance, trust, or timeline? Is their workload subject to the DR mandate?

**Principles:** Understand the need before defending the platform. A paved path has to earn adoption. Protect the org's guarantees (DR, security).

**Steps:**
1. Meet them to understand the real requirement. Often it's one missing capability.
2. Separate the need from their proposed solution. Can the platform meet it, and how fast?
3. Make the cost visible: they'd own DR compliance, upgrades, on-call, and security for their own stack.
4. If it's a real gap, put it on the roadmap with a date, or offer an extension point.
5. If they still need to diverge, agree on what they own and how they'd come back later. Document it.

**Real example:** Story #9. Tenants asked for internal forks; I held the line that forks aren't a paved path, and worked the needs into upstream features instead (#3: external etcd, custom CA, token rotation).
</details>

<details>
<summary><strong>H2. Your manager asks you to ship something you think is unsafe.</strong></summary>

**Clarify:** How unsafe: customer impact, data loss, security, or just risky? What's driving the date?

**Principles:** Make the risk concrete and the decision explicit. Offer options, not a veto. Disagree and commit unless it's a real safety or integrity issue.

**Steps:**
1. Write the risk down in one paragraph: what can break, the blast radius, likelihood, and how we'd detect and recover.
2. Bring options: ship behind a flag, a staged rollout, ship a smaller scope, or move the date by N days.
3. Recommend one, with the tradeoff.
4. If they decide to ship, commit fully and add the mitigations (rollback plan, alerts, on-call ready).
5. If it's a real safety, security, or compliance issue, escalate openly, and tell your manager you're doing it.

**Real example:** Story #7. The ingress-nginx → Istio migration had to be zero-downtime, so the plan was built around staged rollout and rollback rather than a single cutover. `[Add a real moment where you pushed back on a date or scope, if you have one.]`
</details>

<details>
<summary><strong>H3. Two senior engineers are deadlocked on a design.</strong></summary>

**Clarify:** How long has it been stuck? What's the cost of delay? Is it a one-way or two-way door?

**Principles:** Move the debate from opinions to criteria. Decide at the right level, at the right speed.

**Steps:**
1. Meet each one separately to understand the strongest version of their argument.
2. Get agreement on the decision criteria first (latency, operability, migration cost, team skills).
3. Write a short doc comparing options against those criteria. Have each person review the description of their own option.
4. If it's reversible, pick one, time-box it, and define what would make you switch. If it's irreversible, prototype or measure the key unknown.
5. Make or drive the call, credit both, and make sure the person whose option lost owns part of the result.

**Real example:** Story #13 (ARB) or a Karmada design review where contributors disagreed `[supply the specific case]`.
</details>

<details>
<summary><strong>H4. You join and the system is on fire.</strong></summary>

**Clarify:** On fire how: outages, on-call burnout, missed commitments? Is someone already leading the response?

**Principles:** Stabilize first, then find the root causes, then fix the system that produced the fire. Earn trust before changing things.

**Steps:**
1. Join on-call or the incident rotation right away to learn where it hurts.
2. Triage: top 3 sources of pain by impact. Quick, safe wins first (silence noisy alerts, add runbooks).
3. Define what "healthy" means: SLOs scoped to customer experience.
4. Build a short plan with the EM: what to stop, what to fix, what to defer.
5. Fix the process that caused it (alerting, release safety, ownership gaps).

**Real example:** Story #6. Noisy infra-centric alerts → SLOs scoped to tenant experience and symptom-based alerting across clusters; tenants now fix their own issues.
</details>

<details>
<summary><strong>H5. A tenant escalates to your VP about your platform.</strong></summary>

**Clarify:** What's the complaint? Is it valid? Did they try to reach the team first?

**Principles:** Fix the problem, not the politics. Don't get defensive. Close the loop with everyone.

**Steps:**
1. Reach out to the tenant directly, same day. Understand the impact.
2. Brief your manager and the VP: what happened, what we're doing, when it'll be resolved. No blame.
3. Fix or mitigate the immediate issue.
4. Ask why they felt they had to escalate. Usually it's a communication gap.
5. Fix that gap: office hours, a clear support channel, published roadmap.

**Real example:** Story #9. Tenant office hours exist so issues surface early and stakeholders shape the roadmap in the open.
</details>

<details>
<summary><strong>H6. You disagree with the roadmap.</strong></summary>

**Clarify:** Disagree on what: priorities, approach, timing? Has the decision been made already?

**Principles:** Disagree with data, in the right forum, early. Then commit.

**Steps:**
1. Make sure you understand the reasoning behind the roadmap. Ask the owner.
2. Write up the alternative: the cost of the current plan, the proposed change, the evidence.
3. Raise it with the decision maker, privately first.
4. If it's adopted, help drive it. If not, commit fully and say so. Revisit only if new data arrives.

**Real example:** `[A roadmap or prioritization disagreement from Bloomberg or Karmada.]` Or Story #1: making the case for Karmada as the firm standard.
</details>

<details>
<summary><strong>H7. An engineer on your team is struggling.</strong></summary>

**Clarify:** Struggling how: skills, scope, motivation, personal issues? Are you their manager or a peer?

**Principles:** Assume good intent. Find the cause before giving advice. Partner with the manager on performance; don't go around them.

**Steps:**
1. Have a private conversation: "How's it going? What's getting in the way?"
2. Diagnose: unclear expectations, a skill gap, a project that's too big, or something personal.
3. Act on the cause: break the work into smaller pieces, pair on it, give specific feedback, or connect them with resources.
4. Give them a visible win.
5. If it's performance, partner with their manager.

**Real example:** Story #5 / #14. Mentoring contributors from their first PR to design proposals by matching work to their level and increasing scope gradually.
</details>

<details>
<summary><strong>H8. You're asked to own an area you know nothing about (e.g., GPU scheduling).</strong></summary>

**Clarify:** What's the timeline? Who's the current expert? What decisions are pending?

**Principles:** Learn fast and in public. Lean on experts. Don't make irreversible calls until you understand the domain.

**Steps:**
1. Two weeks of fast learning: read the code and docs, shadow on-call, interview the 5 people who know it best and the main customers.
2. Write down what you've learned and the open questions; get it reviewed by the experts.
3. Find where your existing skills apply (for GPU scheduling: capacity-aware placement, failure domains, multi-tenant quotas).
4. Ship one small improvement to build credibility.
5. Then set direction, with the experts as partners.

**Real example:** Karmada: adopter → maintainer and operator sub-project tech lead. Also §9 below for how to say this about the GPU gap.
</details>

<details>
<summary><strong>H9. Your first 90 days as a Staff engineer.</strong></summary>

**Principles:** Learn before changing things. Earn trust with an early win. Then set direction.

| Days | Focus |
|---|---|
| 0–30 | Learn the system and the people. 1:1s with the team, partner teams, key customers. Join on-call. Read the design docs and postmortems |
| 30–60 | Find leverage points: the 2–3 problems that matter most and nobody owns. Deliver one early win |
| 60–90 | Write a short direction doc with the EM: where we should be in 12 months, and the first steps. Start mentoring and setting standards |

**Tie it back:** "That's roughly how I approached the platform at Bloomberg: understand why each team was solving DR on its own, find the common need, prove it with a POC, then make it the standard."
</details>

<details>
<summary><strong>H10. A critical project is slipping.</strong></summary>

**Clarify:** How far behind? Why: scope, unknowns, dependencies, staffing? What's the hard deadline driven by?

**Principles:** Surface it early with options. Protect the outcome, not the original plan.

**Steps:**
1. Get the real status: what's done, what's left, the biggest risk.
2. Find the cause. Don't just add people.
3. Bring options to stakeholders: cut scope, move the date, add focused help on the critical path, or change the approach.
4. Recommend one. Communicate the new plan widely.
5. Afterward, look at why the estimate was wrong.

**Real example:** `[A real Bloomberg or Karmada delivery story, e.g., the prod GA push or the Karmada v1.17 rollout, if one slipped.]`
</details>

<details>
<summary><strong>H11. A peer team ships a change that breaks your platform.</strong></summary>

**Steps:** mitigate first (roll back or work around), then a blameless joint review, then fix the contract: tests, a change-notification process, or conformance checks. **Real example:** conformance testing work in Karmada; the safeguard for the member-cluster state-flapping race condition.
</details>

<details>
<summary><strong>H12. Leadership wants a big new initiative, but on-call is drowning the team.</strong></summary>

**Steps:** quantify the toil (hours/week, pages), show the tradeoff explicitly, propose a split (e.g., one sprint of reliability work to free capacity, then the initiative), and get the decision made at the right level. **Real example:** Story #6.
</details>

<details>
<summary><strong>H13. You find out a decision you made months ago was wrong.</strong></summary>

**Steps:** say so early to the people affected, assess the impact, propose the fix and the migration path, and share what you learned. **Real example:** Story #2: the standalone topology made the management cluster a single point of failure, so the design evolved to stretched.
</details>

---

## 4. Deep follow-up drilling

Interviewers drill to find out whether a story is real and whether *you* drove it. A real story gets more specific with each question. A borrowed or inflated story gets vaguer.

### The "5 whys" pattern

A typical chain on Story #2:

1. "Why did you move away from the standalone topology?" → single point of failure, unacceptable for a DR mandate.
2. "Why not active/passive?" → `[your real reason: e.g., failover time, operational complexity, two endpoints for tenants]`.
3. "How did you handle etcd across data centers?" → one multi-DC etcd cluster from the internal managed etcd service; cross-DC latency was the tradeoff for surviving a site loss.
4. "What broke?" → `[a real issue]`.
5. "What would you do differently?" → `[your Learning]`.

**Prepare for this:** for each top story, write down 5 levels of "why" and "how" before the interview (`03-storytelling-framework.md`).

### Common drill questions

| Question | What they're testing | Your move |
|---|---|---|
| "What exactly did *you* do?" | Your role vs. the team's | Name 2–3 specific decisions or artifacts: "I wrote the design, chose X over Y, and drove the alignment with the etcd and DNS teams." Then credit the team for execution |
| "What would your teammates say about you?" | Self-awareness, consistency | One strength with evidence, one honest growth area: "They'd say I raise the bar on reliability; some would say I can over-invest in docs before shipping." `[make it true]` |
| "Why not X?" | Did you consider alternatives? | "We considered X. We rejected it because… If Z had been true, X would have been the right call." |
| "What was the hardest part?" | Depth, honesty | Usually people or alignment, not code. Be specific |
| "Who disagreed?" | Conflict handling, realism | Name the role (not the person), their concern, how you resolved it |
| "How did you measure success?" | Results orientation | The number: weeks → minutes, 40+ fleets, <5 → 20+ |
| "What would you do differently?" | Learning, maturity | A real, specific change, not "communicate more" |
| "What happened after you left / after launch?" | Durability, follow-through | What's still running, who owns it, what changed |
| "How big was the team? What was your role on it?" | Scope calibration | Honest numbers; describe your leadership role clearly |
| "How did you get buy-in?" | Influence | Who, what concern, what changed their mind |

### Signs you're answering well

- You give names of components, numbers, and dates without hesitation.
- You mention things that went wrong without being asked.
- Your answer to "why not X?" shows you thought about X *at the time*.

---

## 5. Skeptical or hostile interviewers

Some interviewers are neutral or cold by design (bar raisers especially). Some are just tired. Either way, it isn't a signal you're failing.

| What happens | What it probably means | Your move |
|---|---|---|
| Flat expression, no feedback | Their style, or they're taking notes | Keep your energy steady. Check in: "Want me to go deeper or move on?" |
| "That doesn't sound that hard." | Testing whether the scope was real | Don't get defensive. Add the specific constraint: "The hard part was doing it with zero downtime across N clusters while tenants kept deploying." |
| Interrupts constantly | Wants a specific signal, or time pressure | Let them. Answer what they asked, briefly |
| Challenges your decision | Testing conviction and openness | "That's a fair challenge. At the time we chose X because… If I were doing it today, I'd still pick X, but I'd add Y." |
| Openly dismissive | Rare, but happens | Stay professional and specific. Your written feedback still counts; a calm answer reads well in the debrief |

**Don't:** argue to win, chase approval, or let it shake the rest of your loop. Reset between rounds.

---

## 6. Repeating stories and the second failure question

### Can you repeat stories across rounds?

Usually yes, with care. Different interviewers often don't compare notes until the debrief, and a story used twice is fine if it's for **different competencies** and told from **different angles**. But:

- Amazon-style loops often discourage repeats, since each interviewer covers different principles. Prepare more stories.
- Don't use the same story for more than two rounds.
- In a single round, never repeat a story.

**Vary the angle, not the facts:** Story #1 can be about technical direction (why Karmada), influence (getting 10+ platform teams aligned), or ambiguity (turning a mandate into a platform). Same facts, different focus.

### "Tell me about another failure"

Have **two** failure stories and **two** conflict stories ready (`08-failure-judgment-and-prioritization.md`, `07-conflict-disagreement-and-pushback.md`). Candidates for the second failure story:

- The standalone topology's single point of failure (an early design that had to change, Story #2).
- The member-cluster state-flapping race condition (needed a safeguard after the fact).
- `[A people failure: a delegation, communication, or prioritization miss.]`

If you truly have nothing new: use a smaller one honestly. "A smaller example is…"

---

## 7. Confidentiality-limited answers

You'll often be asked for numbers or details you can't share. Saying so well is a positive signal.

> I can't share the specific numbers, but the shape of the problem was: dozens of control planes across multiple data centers, with tenants running mission-critical workloads. What I can tell you is that setup went from weeks to minutes and the platform now backs 40+ cluster fleets.

**Rules (from `../interview-playbook.md` §1.6):**
- Use public framing: talks, the CNCF reference architecture, Karmada GitHub.
- Say "an internal managed etcd service," not internal system names.
- Never share tenant details beyond what's public, internal metrics, or incident specifics.
- Offer the shape, the order of magnitude if public, and the outcome.

---

## 8. Comp, salary history, levels, and inappropriate questions

### Comp expectations

> I'd like to understand the range for this level first. Can you share the Staff range for this role?

If pushed:

> I'm focused on finding the right role and level. I'm confident we can land on something fair if there's a fit. I'd rather sort out numbers with the recruiter once we're further along.

With an interviewer (not the recruiter): "I'd rather leave that to the recruiter."

### Salary history

New York State and New York City law bars employers from asking about your salary history. You don't have to answer.

> I'd rather focus on the range for this role and the value I'd bring.

### "What's your level at Bloomberg?"

> Bloomberg doesn't have formal levels above Senior, but my scope maps to Staff: I set the firm-wide standard for multi-cluster DR, own the platform end to end, sit on the Architecture Review Board, and maintain a CNCF graduated project.

Say it once, calmly, and move on. Don't apologize for the title.

### Illegal or inappropriate questions

Questions about age, family plans, religion, national origin, disability, and similar topics are generally off-limits in US hiring. Options:

| Approach | Sounds like |
|---|---|
| Answer the job-relevant concern | "If you're asking about my availability, I can commit to the role's schedule." |
| Redirect | "I'd rather keep the focus on the role. Happy to talk more about X." |
| Decline politely | "I don't think that's relevant to the role." |

Afterward, consider telling the recruiter. It tells you something about the company.

---

## 9. The AI / GPU experience gap

This will come up at AI-infra companies. Answered calmly, it can be a strength: it shows honesty and how you ramp.

**The structure:** honest gap → what you do know → what transfers → how you'd ramp → evidence you ramp fast.

> Honestly, I haven't operated GPU fleets at the node level. My AI exposure is as the platform underneath our data science platform, which runs distributed training across heterogeneous clusters, so I know what those teams need from the control plane and scheduling layer. What transfers directly is control-plane lifecycle, fleet automation, failure-domain thinking, and capacity-aware placement. What I'd ramp on is the GPU-specific layer: node health and burn-in, the InfiniBand fabric, and topology-aware scheduling. I ramp fast on new domains. I went from Karmada adopter to maintainer.

**Don't:**
- Bluff about NCCL, InfiniBand, or training internals. One wrong detail costs more than the gap.
- Over-apologize or bring it up repeatedly.
- Say "I can learn anything." Show a specific plan.

**Do:** show you've done some homework (DRA, Kueue, JobSet, gang scheduling, checkpointing; see `../interview-playbook.md` §5.2), and ask a sharp question about how they handle it.

---

## 10. Running long or getting cut off

**Prevent it:** give the headline first, then the story.

> Short answer: I grew Karmada's contributor base from fewer than 5 to over 20. Here's how.

Aim for 2–3 minutes per story. If you see yourself past 3 minutes, wrap up:

> …so the result was X. Happy to go deeper on any part.

**If you're cut off:** that's fine. They likely got what they needed, or they're short on time. Don't fight to finish. "Sure." Then answer the new question.

**If you're unsure whether to keep going:** "Is this the level of detail you're looking for, or should I zoom out?"

---

## 11. Recovering from a bad answer

A weak answer early in a round can usually be fixed in the same round. Interviewers respect the self-awareness.

> Before we wrap up, can I revisit something I said earlier? When you asked about conflict, I think I underplayed my role. What I actually did was…

or

> I'd like to give a better example for the failure question. The one I gave was more of a setback. A real failure was…

**When:** at a natural transition, or when they ask "Any questions?" Keep it under a minute.

**If you realize after the round:** mention it briefly in your thank-you note, at most one line, only if it matters. Otherwise let it go and focus on the next round. One weak answer rarely sinks a loop.

**Between rounds:** take 2 minutes to note what went well and what to adjust. Then reset completely.

---

## 12. Virtual-interview hiccups

| Problem | Your move |
|---|---|
| Audio or video drops | "Sorry, I think I lost you. Could you repeat the last part?" Have a phone ready to dial in |
| Your connection fails | Rejoin quickly; email or message the recruiter right away if you can't |
| Interviewer is late or absent | Wait 10 minutes, then email the recruiter. Don't take it personally |
| Background noise or interruption | Name it briefly, mute, continue |
| You need a moment | "Give me a second to think about that." Totally fine |

**Setup checklist** (also `12-practice-plan-and-day-of.md`): stable internet, camera at eye level, good light, quiet room, phone backup, story map on paper nearby (not on screen where reading looks obvious).

---

## 13. Ending a round strong

The last five minutes shape the interviewer's memory of you.

1. **Ask 2–3 thoughtful questions** tailored to the interviewer's role (bank in `10-motivation-career-and-about-you.md`). For EMs: their hardest problem. For senior ICs: how Staff engineers influence direction. For execs: the biggest bet.
2. **Optionally add one missing signal:** "One thing I didn't get to mention that's relevant: I also sit on the Architecture Review Board, so I've done a lot of cross-team design review."
3. **Close briefly and warmly:** "I really enjoyed this. What you said about X is exactly the kind of problem I want to work on."
4. **Thank-you within 24 hours,** naming one specific thing from the conversation.

**Avoid:** "How did I do?" (puts them on the spot), running over time, or ending with "No questions."

---

## 14. Interviewer probes → what they're testing → your move

| Probe | Testing | Your move |
|---|---|---|
| "Tell me more about that." | Depth, realism | Go one level deeper: a specific decision, number, or person |
| "What was your role specifically?" | Ownership | Name your decisions; credit the team's execution |
| "Why not X?" | Alternatives considered | Explain the tradeoff at the time; say when X would win |
| "What would you do differently?" | Learning | A concrete change, not a platitude |
| "What did the other person think?" | Empathy, fairness | State their view fairly, even if you disagreed |
| "How did you know it worked?" | Results | The metric, before and after |
| "What if that hadn't worked?" | Contingency thinking | Your fallback plan at the time |
| "Was that your idea?" | Initiative | Honest answer; if not yours, what you added |
| "How long did that take?" | Pace, scope | Real timeline, and what drove it |
| "Who else was involved?" | Cross-team scope | The teams and roles you aligned |
| "What would you do if…?" | Judgment | Clarify → principles → steps → real example (§2) |
| "Anything else?" | Whether you have more | One crisp addition, or "That covers it" |
| Silence after your answer | Comfort, or waiting for more | Don't ramble. "Happy to go deeper on any part." |
| "That seems small." | Scope calibration | Add the constraint or stakes that made it hard |

---

## Self-check

1. Say your pause script out loud. What are the four options after it?
2. Walk through H1 and H8 using the four-step framework, ending with a real example.
3. For Story #2, write the 5-level "why" chain. Where does it get vague?
4. What's your second failure story? Your second conflict story?
5. Say your GPU-gap answer in under 60 seconds.
6. What do you say if asked for your current salary in New York?

Next: `12-practice-plan-and-day-of.md` to rehearse all of this under realistic conditions.
