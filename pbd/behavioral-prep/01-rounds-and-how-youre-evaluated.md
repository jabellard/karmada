# 01. The Rounds and How You're Evaluated

Behavioral evaluation doesn't happen in one "behavioral round." It starts with the recruiter's first call and ends when the last reference is checked. Every person you talk to writes something about you, and those notes decide both **whether** you get an offer and **at what level**. This module walks the whole journey, round by round, so you know who's in the room, what they're writing down, and what a Staff answer looks like at each stage.

---

## At a glance

| Stage | Who | What they're really deciding | Your #1 goal | Main material |
|---|---|---|---|---|
| Recruiter screen | Recruiter / sourcer | Fit, level, motivation, logistics | Get slotted at **Staff** | 60s pitch, why leaving, level line (`10-motivation-career-and-about-you.md`) |
| Hiring manager / EM screen | The manager you'd report to | "Can this person own a hard area on my team?" | Talk about *their* problems | Stories #1, #2, #6 |
| Behavioral / leadership round(s) | EM, senior IC, or trained interviewer | Competencies at Staff scope | One clear, quantified story per question | Story bank (`04-story-bank.md`) |
| Past-project deep dive | Senior/Staff IC, sometimes EM | Ownership, judgment, depth, "what did *you* do" | Show your decisions and rejected alternatives | Story #2 (architecture evolution) |
| Cross-functional / bar raiser | Someone outside the team | "Does this person raise the bar for the company?" | Org-level scope, high standards | Stories #1, #4, #9, #13 |
| Skip-level / director / VP | Your manager's manager | Strategy, judgment, long-term fit | Talk like a partner on direction | Business framing of #1, #6 |
| Team matching | Several EMs | Mutual fit with a specific team | Find the right scope; be easy to say yes to | Pitch + targeted stories |
| Founder / exec chat | CEO, CTO, cofounder | Mission, ownership, pace, ego | Show conviction and low ego | Why this company, Story #10 |
| References | Former managers, peers, mentees | "Is the loop's picture true?" | Brief your references well | 3 people who can speak to Staff scope |
| Debrief / leveling | Interviewers, HM, hiring committee | Hire? At which level? | (Already done) Consistent signal across rounds | Everything above |

**The key idea:** coding and design rounds mostly decide *hire or no hire*. Behavioral rounds, the project deep dive, and the HM's read are usually what decide *Staff or Senior*. Treat them as leveling rounds, not soft rounds.

---

## 1. Recruiter screen (20–30 min)

**Who's in the room:** a recruiter, sometimes a sourcer. Not technical, but often influential on level and timeline. They are your advocate inside the process if you make their job easy.

**What they care about:**
- Do your background and the role match well enough to spend interviewer time on you?
- What level should the loop be calibrated for?
- Are you motivated by *this* company, or just shopping?
- Logistics: location, timeline, other processes, comp expectations.

**What they write in their notes:** level they'll propose, a one-line summary of your background, motivation, red flags (vague, negative about employer, unrealistic expectations), timeline and competing processes.

**Typical questions:**
- "Tell me about yourself."
- "Why are you looking?" / "Why us?"
- "What level are you targeting?" / "What's your current title?"
- "What are your comp expectations?"
- "What's your timeline? Are you interviewing elsewhere?"

**Your goal:** leave with the loop calibrated at Staff, a clear picture of the process, and no number given first.

**Staff signals here:**
- A crisp 60-second pitch with scope and numbers (40+ fleets, weeks → minutes, CNCF maintainer).
- You state your level calmly and early, with evidence:
> Bloomberg doesn't have formal levels above Senior, but my scope maps to Staff: I set the firm-wide standard for multi-cluster DR, own the platform end to end, sit on the Architecture Review Board, and maintain a CNCF graduated project. I'd like to be considered for Staff.
- You know what you want: platform / AI infrastructure, a hard area to own.

**Pitfalls:**
- Rambling pitch (over 2 minutes).
- Letting "Senior" stick because you didn't say anything. Changing level after the loop is much harder.
- Naming a number first, or answering salary history (in NYC, employers can't ask; see `11-curveballs-hypotheticals-and-recovery.md` §8).
- Sounding like you're running *from* Bloomberg rather than *toward* something.

**Questions to ask them:**
- "What does the loop look like, and which rounds carry the most weight for leveling?"
- "Is leveling decided before or after the loop? Can I be considered for Staff?"
- "What's the range for Staff on this role?" (NYC postings must include ranges.)
- "Is there a team-matching phase?"
- "What problem is the team hiring to solve?"

**Follow-up:** same-day thank-you email with one line restating level interest and one line on what excites you about the team.

---

## 2. Hiring manager / EM screen (30–60 min)

**Who's in the room:** the engineering manager you'd likely report to. Sometimes a senior IC joins. This is the most important single conversation for your level, because the HM is often the strongest voice in the debrief.

**What they care about:**
- Can I hand this person an ambiguous area and stop worrying about it?
- Do their skills map to *my* team's problems?
- Will they make my team better (mentoring, standards, direction)?
- Will they partner with me, or create management overhead?

**What they write:** problem-space fit, scope of past work, ownership, communication, "would I want this person on my team," proposed level, and which rounds to probe harder ("verify depth on X in the deep dive").

**Typical questions:**
- "Walk me through your current role / your platform."
- "What's the hardest technical problem you've worked on?"
- "How do you work with your manager? With product? With other teams?"
- "Tell me about a time you set direction for a team."
- "Why this team?" / "What do you want to do next?"
- "What would you do in your first 90 days?"

**Your goal:** make the EM picture you owning their hardest problem.

**Staff signals here:**
- Talk about **their** problems. Research the team's talks, blogs, and open source first.
> From your talk on X, it sounds like Y is a pain point. At Bloomberg we hit a version of that, and here's what worked and what didn't.
- Show a 90-day approach: learn the system and the people → find leverage points → deliver an early win → set direction.
- Show you think about the EM's goals: team velocity, on-call load, hiring, roadmap.

**Pitfalls:**
- Treating it like a screen you have to pass, instead of a two-way conversation between future partners.
- Over-explaining Karmada internals. To them, Karmada is evidence you can do the work, not the topic.
- Bluffing about GPU/HPC depth. Name the gap calmly and show how you'd close it (`11-curveballs-hypotheticals-and-recovery.md` §9).

**Questions to ask them:**
- "What's the hardest open problem you'd want this person to own?"
- "What would make someone exceptional in this role in the first year?"
- "Where is the team's technical direction unclear right now?"
- "How do Staff engineers here influence the roadmap? How do you and your Staff engineers split responsibilities?"

**Follow-up:** thank-you within 24 hours naming the specific problem they described and one sentence on how you'd approach it.

See §11 below for how EMs evaluate Staff ICs in more detail.

---

## 3. Behavioral / leadership rounds in the loop (45–60 min, often 1–2 of them)

**Who's in the room:** an EM from a neighboring team, a senior/Staff IC, or a trained behavioral interviewer. At some companies, a specific person owns the "leadership" or "culture" signal.

**What they care about:** consistent evidence on a fixed set of competencies, at the scope of the level you're interviewing for. They usually have a question list mapped to competencies and need enough evidence to rate each one.

**What they write:** per competency, a rating plus evidence ("Candidate described X; their role was Y; result Z"). Weak notes say "seemed good." Strong notes quote your decisions and numbers. Your job is to give them quotable material.

**Common competencies (Staff level):**

| Competency | Sample question | Primary story |
|---|---|---|
| Technical direction | "Tell me about the biggest technical decision you drove." | #1, #2 |
| Influence without authority | "Convince people you didn't manage?" | #3, #1 |
| Ambiguity | "A problem with no clear owner or solution?" | #1, #10 |
| Cross-team execution | "A project spanning several teams?" | #4, #7 |
| Conflict | "Disagreed with a peer or manager?" | #12, #13 |
| Growing others | "How have you grown other engineers?" | #5, #14 |
| Failure / learning | "Tell me about a failure." | #11 |
| Prioritization / saying no | "Pushed back on a stakeholder?" | #9 |
| Operational excellence / customers | "Improved reliability or reduced toil?" | #6, #8 |

Full question banks per competency live in `05` through `10`.

**Your goal:** 4–6 clean stories in 45 minutes, each 2–3 minutes, each with *your* decisions and a number.

**Staff signals here:** scope beyond your team, decisions with rejected alternatives, multiplier effect, and a real Learning at the end (`02-staff-signals.md`, `03-storytelling-framework.md`).

**Pitfalls:**
- Five-minute monologues that leave no time for follow-ups (fewer stories = fewer data points = weaker ratings).
- Team-scope stories when the level needs org scope.
- "We" for everything, so the interviewer can't credit you.
- Polished stories with no tension, disagreement, or tradeoff. They sound rehearsed and shallow.

**Questions to ask them:** "What distinguishes the Staff engineers you've worked with here?" "How do Staff ICs and EMs partner on this team?"

---

## 4. Past-project deep dive (45–60 min)

Covered technically in `../interview-playbook.md` §5.3. This section is the behavioral side, which is often what decides the level.

**Who's in the room:** a senior/Staff IC, sometimes with the EM. They're probing whether the project was as big as it sounds and whether *you* drove it.

**What they care about:**
- What did *you* do versus the team?
- Did you make the important calls, and were they good calls?
- Do you understand the system end to end, including the parts you didn't build?
- Can you explain tradeoffs and what you'd change?

**What they write:** your role (driver, contributor, or observer), depth, quality of decisions, honesty about mistakes, scope of impact.

**Typical probes:**
- "Which part did you personally build or decide?"
- "Who disagreed, and how was it resolved?"
- "Why stretched control planes instead of active/passive?"
- "What happens if both management clusters fail?"
- "What would you do differently?"
- "How did you get the other teams to adopt it?"

**Your goal:** make your ownership unmistakable while being generous to the team.

**Staff signals:**
- "I decided… because… I rejected… because…" for every major decision (Story #2: standalone → stretched → source of truth).
- You bring up what went wrong and what you learned before they ask.
- You describe the human side: aligning 10+ platform teams, the ARB, tenants.
- You can say what's next: conformance testing, automated DR testing.

**Pitfalls:** claiming the whole team's work (they'll drill until it breaks), or the opposite: "we did it" so often you get rated as a contributor. Internal names or metrics beyond what's public (say "an internal managed etcd service").

---

## 5. Cross-functional / bar raiser round (45–60 min)

**Who's in the room:** someone outside the hiring team, often senior, sometimes with veto power (Amazon's bar raiser is the best-known version). They're protecting the company's bar, not filling the team's headcount.

**What they care about:** Would this person raise the average? Do they hold high standards? Are they someone the whole org benefits from, not just this team?

**What they write:** a hire/no-hire with heavy weight on judgment, standards, and integrity; comparisons against people already at that level.

**Typical questions:** mixed behavioral and judgment: "Tell me about a time you raised the bar for your org." "A time you had to make a decision with incomplete data." "When did you disagree and commit?" Often more follow-ups and less warmth by design.

**Your goal:** raise the scope of every answer to org level.

**Staff signals:**
- Org-wide standards: Karmada as the firm's DR standard (#1), ARB reviews (#13), shared-responsibility model for tenants (#9).
- External impact: CNCF graduation (#4), upstream features (#3).
- High standards you hold others to, and how you did it kindly.

**Pitfalls:** treating the neutral or skeptical tone as a sign it's going badly. It's often deliberate. Stay calm, stay specific (`11-curveballs-hypotheticals-and-recovery.md` §5).

---

## 6. EM, skip-level, director, or VP chats (30–45 min)

**Who's in the room:** your manager's manager or above. Usually late in the loop or after it, sometimes as a "sell" call, sometimes as a real evaluation. Assume it's evaluative.

**What they care about:** strategic thinking, judgment, long-term fit, whether you'll operate at the level they're paying for, and whether you'll stay.

**What they write:** a short read on scope ("thinks at org level" or "very hands-on, team-level"), culture/values fit, any concerns.

**Typical questions:**
- "Where do you see infrastructure going in the next few years?"
- "What would you want to be doing in 3 years?"
- "What's the biggest bet you've made?"
- "What do you need from leadership to be successful?"

**Your goal:** talk like a partner in the org's direction, not a candidate asking for a job.

**Staff signals:** connect technical work to business outcomes (DR mandate, tenant velocity, on-call load, cost). Have a point of view on their space (e.g., multi-cluster and multi-region for AI workloads, GPU fleet reliability).

**Questions to ask them:**
- "What will the infrastructure org look like in two years, and what's the biggest bet?"
- "Where do you want Staff engineers to have the most influence?"
- "What's a decision you're wrestling with right now?"

---

## 7. Team matching (common at Google- and Meta-style companies)

**How it works:** you pass a general loop, then talk with several EMs whose teams have headcount. Each chat is usually 30–45 minutes. Either side can say no. At some companies the offer isn't final until a match happens.

**What EMs care about:** fit with their specific problems, and whether you'd be excited (they've seen your packet, so they're mostly selling and checking chemistry).

**Your goal:** find the team where you'd own something hard at Staff scope, and be easy to say yes to.

**Do:**
- Prepare a version of your pitch per team.
- Ask about scope and Staff engineers on the team.
- Ask the questions that matter to you: fleet/AI-infra content, upstream involvement (can you keep maintaining Karmada?), on-call, promotion path.
- Be honest about preferences, kindly. A bad match is worse than a slower one.

**Questions to ask:** "What would I own end to end?" "Who are the Staff engineers here and what do they work on?" "How does promotion to Senior Staff work on this team?"

---

## 8. Founder / exec chats at startups and AI labs (30–60 min)

**Who's in the room:** CEO, CTO, a cofounder, or head of infrastructure. Common at startups (Modal, Upbound, Spectro Cloud, Rafay) and some AI labs.

**What they care about:** mission alignment, ownership ("will they do whatever it takes?"), pace, low ego, and whether you'd raise the talent bar.

**Typical questions:** "Why us, specifically?" "What's the most impressive thing you've built?" "What would you do in your first month?" "What do you think we're getting wrong?"

**Your goal:** show conviction about their problem and that you'll own outcomes, not tasks.

**Staff signals:** you've done real homework, you have a thoughtful (not arrogant) take on their product, and you've shown initiative outside your job description (Story #10 Hoppscotch, Story #3 upstreaming).

**Pitfalls:** big-company vocabulary ("I'd need a team to…"), comp questions too early, or flattery without substance.

---

## 9. Reference checks

**Who's involved:** the references you give, and sometimes back-channel references (people they know who've worked with you, especially in a small community like CNCF).

**What they ask your references:** scope of your work, how you handled conflict, strengths, growth areas, "would you hire them again," and "at what level would you place them."

**Your goal:** references that confirm Staff scope with specific examples.

**Pick three (from `../interview-playbook.md` §6):**
- A former manager (scope, judgment, delivery).
- A Karmada co-maintainer (influence without authority, technical leadership).
- A tenant lead or mentee (multiplier, customer focus). A mentee who now drives proposals or scheduler work is a strong Staff signal.

**Brief each one:** the role, the level, and the 2–3 themes to emphasize, with the story they can tell for each. Thank them afterward and tell them the outcome.

**Discretion:** don't list current Bloomberg colleagues unless you're comfortable with them knowing you're looking.

---

## 10. Post-loop debrief and leveling

You aren't in this room, but knowing how it works changes how you interview.

### How hiring decisions are usually made

1. **Written feedback first.** Each interviewer writes notes and a rating, often before seeing anyone else's, to avoid anchoring. Vivid, quotable evidence ("described rejecting active/passive because…") carries the most weight.
2. **Debrief.** Interviewers and the HM discuss. Mixed signals get argued over. A single strong advocate helps a lot; a single strong concern can sink you.
3. **Hiring committee** (at some companies, such as Google). A group that didn't interview you reviews the written packet only. If it isn't in the notes, it didn't happen. That's why clear, structured answers matter.
4. **Leveling calibration.** Someone compares your packet against the bar for each level. This is where "strong hire at Senior" vs "hire at Staff" is decided.
5. **Offer and comp** follow the level.

### The level-down risk

The most common bad outcome for Staff candidates isn't a rejection. It's an offer at Senior. It usually happens because:

| Signal in the packet | How it reads |
|---|---|
| Stories scoped to your own team or your own code | Senior |
| "We" everywhere, unclear personal role | Senior (or weaker) |
| Strong coding, solid design, thin behavioral evidence | Senior |
| You solved problems you were given | Senior |
| You found the problem, set direction, aligned other teams | Staff |
| You grew others who now operate independently | Staff |
| You made tradeoffs with business context and can say what you'd change | Staff |

**Implication:** behavioral rounds are where you prove scope. Lead with your biggest-scope stories (#1, #4, #5, #13), not your most technically interesting ones. If you're leveled Senior anyway, see `../interview-playbook.md` §7.6.

---

## 11. How EMs evaluate Staff ICs

An EM hiring a Staff engineer is hiring a partner. They're asking six questions:

| EM's question | What convinces them | Your evidence |
|---|---|---|
| **Scope:** has this person operated beyond one team? | Firm-wide or cross-org outcomes | #1 firm standard, #4 CNCF graduation, #13 ARB |
| **Ambiguity:** can I hand them a fuzzy area? | You defined the problem, not just solved it | #1 (DR mandate → platform), #10 (problem nobody owned) |
| **Partnership:** will they work *with* me? | You think about roadmap, staffing, priorities, and talk tradeoffs openly | Product-owner role, office hours, roadmap decisions |
| **Multiplier:** will they grow my team? | Others got better because of you | #5 (<5 → 20+), #14 mentoring |
| **Low maintenance:** will they create work for me? | You resolve conflict yourself and escalate well | #12, #9 |
| **Judgment:** will they make the right call alone? | Decisions with alternatives, risk awareness, knowing when to say no | #2, #7, #9 |

### Talking to an EM as a peer partner

- **Ask about their goals, not just the role.** "What are you trying to achieve with the team this year? Where's the gap?"
- **Describe how you'd split the work.** "I'd expect to own technical direction for X and partner with you on prioritization and staffing. How do you like to work with Staff engineers?"
- **Show you make managers' lives easier.** "I tend to surface risks early with a recommendation, not just the problem."
- **Show you care about people too.** Mention mentoring, onboarding docs, and how you've helped people grow into bigger work.
- **Avoid:** asking only about tech stack and perks, or positioning yourself as someone who needs heavy direction.

> In my current role I'm both tech lead and product owner, so I'm used to partnering on the roadmap: I bring a recommendation with tradeoffs, we agree on priorities, and I own getting it done across teams.

---

## 12. Company-style differences

Processes change often, so confirm with your recruiter. These are common patterns, not guarantees.

| Company type | What's commonly emphasized | How to adapt |
|---|---|---|
| **Amazon / AWS** | Leadership Principles (Ownership, Dive Deep, Have Backbone; Disagree and Commit, Earn Trust, Deliver Results…). Questions map to principles. A bar raiser from outside the team often has strong influence | Map 2 stories per principle. Expect many follow-ups and "dive deep" drilling. Use data. Prepare more stories than usual, since repeats across interviewers are discouraged |
| **Google** | "Googleyness and leadership" round; written packet reviewed by a hiring committee; team match after the loop | Make answers structured and quotable (the committee only sees notes). Be ready for hypotheticals. Prepare for team-match chats |
| **Meta** | A behavioral round probing signal areas such as conflict resolution, growth/feedback, handling ambiguity, driving results, and working with others | Cover each area with a distinct story. Show results with numbers and what you learned |
| **Stripe, Datadog, and similar infra companies** | Ownership, rigor, clear writing, customer focus; some include writing exercises or doc reviews | Show how you write designs and docs (reliability guide, reference architecture). Emphasize customer (tenant) outcomes |
| **CNCF-heavy companies** (Red Hat, Microsoft, Upbound, SUSE, Grafana) | Open-source collaboration, community influence | Lead with #3, #4, #5. Upstream work is directly relevant |
| **AI labs and AI-infra companies** (CoreWeave, Anthropic, OpenAI, Databricks, Nebius, Modal) | Mission alignment, low ego, speed, judgment under uncertainty, comfort with very high scale and fast change | Have a sincere answer for why AI infrastructure. Show you ramp fast (adopter → maintainer). Be honest about the GPU gap. Show pragmatic speed, not process for its own sake |
| **Startups** | Ownership, breadth, pace, founder fit | Show initiative (#10), willingness to do unglamorous work, and fast iteration |

---

## 13. Round-by-round prep checklist

| Round | Prepare beforehand | Have ready in the room | Done |
|---|---|---|---|
| Recruiter | 60s pitch, why leaving, why them, level line, comp deflection | Questions about loop and leveling, Staff range | [ ] |
| HM / EM screen | Team research (talks, blogs, repos), 3 tailored stories, 90-day plan | 3–4 questions about their hardest problem | [ ] |
| Behavioral | All 14 stories drafted in STAR-L, 2-min and 5-min versions | Story-to-question map (`04-story-bank.md`) | [ ] |
| Deep dive | Public-safe diagrams, decisions + alternatives, failure modes, "what I'd change" | Clear split of "I" vs "we" | [ ] |
| Bar raiser | Org-scope versions of #1, #4, #9, #13; a disagree-and-commit story | Calm under skeptical follow-ups | [ ] |
| Skip-level / VP | Point of view on their space and where infra is going | Strategic questions | [ ] |
| Team match | Per-team pitch variant, your must-haves | Scope and growth questions | [ ] |
| Founder / exec | Product homework, a thoughtful take | Why this mission, now | [ ] |
| References | 3 references chosen and briefed | Themes and stories for each | [ ] |

---

## 14. What each interviewer is writing about you

| Interviewer | Headline they're writing | Dimensions | What makes the note strong |
|---|---|---|---|
| Recruiter | "Strong fit for Staff, motivated, timeline X" | Fit, level, motivation, logistics | Crisp pitch, clear level ask, positive reasons |
| Hiring manager | "Could own [area]; would raise the team" | Problem fit, ownership, partnership, level | Talked about their problems; strong 90-day thinking |
| Behavioral interviewer | "Strong evidence on leadership, conflict, growth" | One rating per competency, with evidence | Quotable decisions, numbers, real Learnings |
| Deep-dive interviewer | "Drove the architecture; strong judgment" | Role, depth, decision quality, honesty | Clear "I decided / rejected," failure modes, what's next |
| Bar raiser | "Raises the bar at Staff" | Standards, judgment, integrity, scope | Org-level scope, disagree-and-commit, high standards held kindly |
| Skip-level / VP | "Thinks at org level, good long-term fit" | Strategy, business judgment, retention risk | Business outcomes, point of view on the space |
| Team-match EM | "Want them on my team" | Fit, excitement, chemistry | Specific interest in their problems |
| Founder / exec | "Owner, mission-driven, low ego" | Ownership, pace, mission, talent bar | Initiative, conviction, humility |
| Reference | "Would hire again at Staff" | Scope, conflict, strengths, growth areas | Specific examples matching your stories |
| Hiring committee / calibration | "Hire at Staff" vs "Hire at Senior" | Consistency across all notes vs. level rubric | Same Staff picture in every round |

**Use this table the night before:** for each round on your schedule, picture the headline you want them to write, and pick the 2–3 stories that earn it.

---

## Self-check

1. Which rounds mostly decide your level, and why?
2. What exact sentence do you say when a recruiter asks your level?
3. Name the six questions an EM asks about a Staff IC. Which story answers each?
4. What's the level-down risk, and which three behaviors in a packet cause it?
5. For your next interview, what headline do you want each interviewer to write?

Next: `02-staff-signals.md` for the signals themselves, then `03-storytelling-framework.md` for how to package them.
