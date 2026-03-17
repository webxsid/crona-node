# Daily Report Template Variables

The daily report templates use Handlebars-style variables.

The same variable set is available to both the markdown and PDF templates.

Supported blocks in v1:
- `{{variable}}`
- `{{#if variable}}...{{/if}}`
- `{{#each items}}...{{/each}}`

Top-level variables:
- `{{date}}`
- `{{generatedAt}}`
- `{{summary.totalIssues}}`
- `{{summary.issueDoneCount}}`
- `{{summary.issueActiveCount}}`
- `{{summary.issueBlockedCount}}`
- `{{summary.issueAbandonedCount}}`
- `{{summary.issueCompletion}}`
- `{{summary.completedIssues}}`
- `{{summary.abandonedIssues}}`
- `{{summary.totalEstimatedMinutes}}`
- `{{summary.estimatedTime}}`
- `{{summary.workedSeconds}}`
- `{{summary.workedTime}}`
- `{{summary.workedEstimate}}`
- `{{summary.varianceTime}}`
- `{{summary.habitsDueCount}}`
- `{{summary.habitsCompletedCount}}`
- `{{summary.habitsPendingCount}}`
- `{{summary.habitCompletion}}`
- `{{dayHealth}}`
- `{{#each highlights}} ... {{/each}}`
- `{{#each risks}} ... {{/each}}`

Nested issue groups for the default template:
- `{{#each repos}} ... {{/each}}`
  - `{{name}}`
  - `{{#each streams}} ... {{/each}}`
    - `{{name}}`
    - `{{#each completedIssues}} ... {{/each}}`
    - `{{#each activeIssues}} ... {{/each}}`
    - `{{#each attentionIssues}} ... {{/each}}`
      - `{{id}}`
      - `{{title}}`
      - `{{status}}`
      - `{{estimateMinutes}}`
      - `{{estimateTime}}`
      - `{{workedSeconds}}`
      - `{{workedTime}}`
      - `{{workedEstimate}}`

Nested habit groups for the default template:
- `{{#each habitRepos}} ... {{/each}}`
  - `{{name}}`
  - `{{#each streams}} ... {{/each}}`
    - `{{name}}`
    - `{{#each completedHabits}} ... {{/each}}`
    - `{{#each pendingHabits}} ... {{/each}}`
      - `{{id}}`
      - `{{name}}`
      - `{{status}}`
      - `{{durationMinutes}}`
      - `{{durationTime}}`
      - `{{notes}}`

Optional objects:
- `{{checkIn.mood}}`
- `{{checkIn.energy}}`
- `{{checkIn.sleepHours}}`
- `{{checkIn.sleepScore}}`
- `{{checkIn.screenTimeMinutes}}`
- `{{checkIn.screenTime}}`
- `{{checkIn.notes}}`

Collections:
- `{{#each issues}} ... {{/each}}`
  - `{{id}}`
  - `{{title}}`
  - `{{repoName}}`
  - `{{streamName}}`
  - `{{status}}`
  - `{{estimateMinutes}}`
  - `{{estimateTime}}`
  - `{{workedSeconds}}`
  - `{{workedTime}}`
  - `{{workedEstimate}}`
- `{{#each sessions}} ... {{/each}}`
  - `{{id}}`
  - `{{issueId}}`
  - `{{issueTitle}}`
  - `{{repoName}}`
  - `{{streamName}}`
  - `{{startTime}}`
  - `{{endTime}}`
  - `{{durationSeconds}}`
  - `{{summary}}`
- `{{#each habits}} ... {{/each}}`
  - `{{id}}`
  - `{{name}}`
  - `{{repoName}}`
  - `{{streamName}}`
  - `{{status}}`
  - `{{durationMinutes}}`
  - `{{durationTime}}`
  - `{{notes}}`

Metrics:
- `{{metrics.sessionCount}}`
- `{{metrics.workedSeconds}}`
- `{{metrics.workedTime}}`
- `{{metrics.restSeconds}}`
- `{{metrics.restTime}}`
- `{{metrics.burnout.level}}`
- `{{metrics.burnout.score}}`
