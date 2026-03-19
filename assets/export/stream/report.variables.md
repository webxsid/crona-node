# Stream Report Template Variables

Top-level:
- `{{startDate}}`
- `{{endDate}}`
- `{{generatedAt}}`
- `{{repo.name}}`
- `{{stream.name}}`
- `{{stream.description}}`

Summary:
- `{{summary.issueCount}}`
- `{{summary.habitCount}}`
- `{{summary.sessionCount}}`

Issues:
- `{{#each issues}}`
  - `{{id}}`
  - `{{title}}`
  - `{{status}}`
  - `{{estimateTime}}`
  - `{{description}}`
  - `{{notes}}`
  - `{{#each sessions}} -> {{summary}}, {{commit}}, {{context}}, {{work}}, {{notes}}`

Habits:
- `{{#each habits}} -> {{name}}, {{scheduleType}}`
