# Issue Rollup Template Variables

Top-level:
- `{{startDate}}`
- `{{endDate}}`
- `{{generatedAt}}`

Issues:
- `{{#each issues}}`
  - `{{id}}`
  - `{{title}}`
  - `{{status}}`
  - `{{scope}}`
  - `{{sessionCount}}`
  - `{{workedTime}}`
  - `{{estimateTime}}`
  - `{{description}}`
  - `{{notes}}`
  - `{{#each sessions}} -> {{summary}}, {{commit}}, {{context}}, {{work}}, {{notes}}`
