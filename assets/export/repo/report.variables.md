# Repo Report Template Variables

Top-level:
- `{{startDate}}`
- `{{endDate}}`
- `{{generatedAt}}`
- `{{repo.name}}`
- `{{repo.description}}`

Summary:
- `{{summary.streamCount}}`
- `{{summary.issueCount}}`
- `{{summary.habitCount}}`
- `{{summary.sessionCount}}`

Streams:
- `{{#each streams}} -> {{name}}, {{description}}`

Issues:
- `{{#each issues}}`
  - `{{id}}`
  - `{{title}}`
  - `{{status}}`
  - `{{scope}}`
  - `{{estimateTime}}`
  - `{{description}}`
  - `{{notes}}`
  - `{{#each sessions}} -> {{summary}}, {{commit}}, {{context}}, {{work}}, {{notes}}`

Habits:
- `{{#each habits}} -> {{name}}, {{scheduleType}}`
