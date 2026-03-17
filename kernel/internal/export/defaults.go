package export

const fallbackDailyReportTemplate = `# Daily Report - {{date}}

Generated at {{generatedAt}}

## Daily Snapshot

- Day health: {{dayHealth}}
- Issues: {{summary.totalIssues}} total, {{summary.issueDoneCount}} done, {{summary.issueActiveCount}} active, {{summary.issueBlockedCount}} blocked, {{summary.issueAbandonedCount}} abandoned
- Issue completion: {{summary.issueCompletion}}
- Worked / estimated: {{summary.workedEstimate}}
- Variance: {{summary.varianceTime}}
- Habits: {{summary.habitsCompletedCount}} completed / {{summary.habitsDueCount}} due
{{#if checkIn}}- Mood / energy: {{checkIn.mood}} / 5, {{checkIn.energy}} / 5{{/if}}
{{#if checkIn.sleepHours}}- Sleep: {{checkIn.sleepHours}}h{{/if}}
{{#if checkIn.screenTime}}- Screen time: {{checkIn.screenTime}}{{/if}}
{{#if metrics.burnout}}- Burnout: {{metrics.burnout.level}} ({{metrics.burnout.score}}){{/if}}

{{#if highlights}}
## Highlights

{{#each highlights}}
- {{this}}
{{/each}}
{{/if}}

{{#if risks}}
## Needs Attention

{{#each risks}}
- {{this}}
{{/each}}
{{/if}}

{{#if checkIn}}
## Wellbeing

- Mood: {{checkIn.mood}} / 5
- Energy: {{checkIn.energy}} / 5
{{#if checkIn.sleepHours}}- Sleep hours: {{checkIn.sleepHours}}{{/if}}
{{#if checkIn.sleepScore}}- Sleep score: {{checkIn.sleepScore}}{{/if}}
{{#if checkIn.screenTime}}- Screen time: {{checkIn.screenTime}}{{/if}}
{{#if metrics}}- Rest time: {{metrics.restTime}}{{/if}}
{{#if checkIn.notes}}

{{checkIn.notes}}
{{/if}}
{{/if}}

## Habits

{{#each habitRepos}}
### {{name}}

{{#each streams}}
#### {{name}}

{{#if completedHabits}}
##### Completed

{{#each completedHabits}}
###### {{name}}

- Status: {{status}}
{{#if durationTime}}- Time: {{durationTime}}{{/if}}
{{/each}}
{{/if}}

{{#if pendingHabits}}
##### Pending

{{#each pendingHabits}}
###### {{name}}

- Status: {{status}}
{{#if durationTime}}- Time: {{durationTime}}{{/if}}
{{/each}}
{{/if}}
{{/each}}
{{/each}}

## Issues

{{#each repos}}
### {{name}}

{{#each streams}}
#### {{name}}

{{#if completedIssues}}
##### Completed

{{#each completedIssues}}
###### #{{id}} {{title}}

- Status: {{status}}
- Worked / estimate: {{workedEstimate}}
{{/each}}
{{/if}}

{{#if activeIssues}}
##### Active

{{#each activeIssues}}
###### #{{id}} {{title}}

- Status: {{status}}
- Worked / estimate: {{workedEstimate}}
{{/each}}
{{/if}}

{{#if attentionIssues}}
##### Needs Attention

{{#each attentionIssues}}
###### #{{id}} {{title}}

- Status: {{status}}
- Worked / estimate: {{workedEstimate}}
{{/each}}
{{/if}}
{{/each}}
{{/each}}
`

const fallbackDailyReportPDFTemplate = `# Daily Report - {{date}}
Generated at {{generatedAt}}
## Snapshot
- Day health: {{dayHealth}}
- Issues: {{summary.issueDoneCount}} done / {{summary.totalIssues}} total
- Habits: {{summary.habitsCompletedCount}} completed / {{summary.habitsDueCount}} due
- Worked / estimated: {{summary.workedEstimate}}
- Variance: {{summary.varianceTime}}
{{#if checkIn}}- Mood / energy: {{checkIn.mood}} / 5, {{checkIn.energy}} / 5{{/if}}
{{#if checkIn.sleepHours}}- Sleep: {{checkIn.sleepHours}}h{{/if}}
{{#if metrics.burnout}}- Burnout: {{metrics.burnout.level}} ({{metrics.burnout.score}}){{/if}}
{{#if highlights}}
## Highlights
{{#each highlights}}
- {{this}}
{{/each}}
{{/if}}
{{#if risks}}
## Needs Attention
{{#each risks}}
- {{this}}
{{/each}}
{{/if}}
## Habits
{{#each habitRepos}}
### {{name}}
{{#each streams}}
#### {{name}}
{{#each completedHabits}}
- {{name}} | {{status}}{{#if durationTime}} | {{durationTime}}{{/if}}
{{/each}}
{{#each pendingHabits}}
- {{name}} | {{status}}{{#if durationTime}} | {{durationTime}}{{/if}}
{{/each}}
{{/each}}
{{/each}}
## Issues
{{#each repos}}
### {{name}}
{{#each streams}}
#### {{name}}
{{#each completedIssues}}
- #{{id}} {{title}} | {{status}} | {{workedEstimate}}
{{/each}}
{{#each activeIssues}}
- #{{id}} {{title}} | {{status}} | {{workedEstimate}}
{{/each}}
{{#each attentionIssues}}
- #{{id}} {{title}} | {{status}} | {{workedEstimate}}
{{/each}}
{{/each}}
{{/each}}
`

const fallbackDailyReportVariables = `# Daily Report Template Variables

The daily report template uses Handlebars-style variables.

The same variable set is available to both the markdown and PDF templates.

Supported blocks in v1:
- {{variable}}
- {{#if variable}}...{{/if}}
- {{#each items}}...{{/each}}

Top-level variables:
- {{date}}
- {{generatedAt}}
- {{summary.totalIssues}}
- {{summary.issueDoneCount}}
- {{summary.issueActiveCount}}
- {{summary.issueBlockedCount}}
- {{summary.issueAbandonedCount}}
- {{summary.issueCompletion}}
- {{summary.completedIssues}}
- {{summary.abandonedIssues}}
- {{summary.totalEstimatedMinutes}}
- {{summary.estimatedTime}}
- {{summary.workedSeconds}}
- {{summary.workedTime}}
- {{summary.workedEstimate}}
- {{summary.varianceTime}}
- {{summary.habitsDueCount}}
- {{summary.habitsCompletedCount}}
- {{summary.habitsPendingCount}}
- {{summary.habitCompletion}}
- {{dayHealth}}
- {{#each highlights}} ... {{/each}}
- {{#each risks}} ... {{/each}}

Nested issue groups for the default template:
- {{#each repos}} ... {{/each}}
  - {{name}}
  - {{#each streams}} ... {{/each}}
    - {{name}}
    - {{#each completedIssues}} ... {{/each}}
    - {{#each activeIssues}} ... {{/each}}
    - {{#each attentionIssues}} ... {{/each}}
      - {{id}}
      - {{title}}
      - {{status}}
      - {{estimateMinutes}}
      - {{estimateTime}}
      - {{workedSeconds}}
      - {{workedTime}}
      - {{workedEstimate}}

Nested habit groups for the default template:
- {{#each habitRepos}} ... {{/each}}
  - {{name}}
  - {{#each streams}} ... {{/each}}
    - {{name}}
    - {{#each completedHabits}} ... {{/each}}
    - {{#each pendingHabits}} ... {{/each}}
      - {{id}}
      - {{name}}
      - {{status}}
      - {{durationMinutes}}
      - {{durationTime}}
      - {{notes}}

Optional objects:
- {{checkIn.mood}}
- {{checkIn.energy}}
- {{checkIn.sleepHours}}
- {{checkIn.sleepScore}}
- {{checkIn.screenTimeMinutes}}
- {{checkIn.screenTime}}
- {{checkIn.notes}}

Collections:
- {{#each issues}} ... {{/each}}
  - {{id}}
  - {{title}}
  - {{repoName}}
  - {{streamName}}
  - {{status}}
  - {{estimateMinutes}}
  - {{estimateTime}}
  - {{workedSeconds}}
  - {{workedTime}}
  - {{workedEstimate}}
- {{#each sessions}} ... {{/each}}
  - {{id}}
  - {{issueId}}
  - {{issueTitle}}
  - {{repoName}}
  - {{streamName}}
  - {{startTime}}
  - {{endTime}}
  - {{durationSeconds}}
  - {{summary}}
- {{#each habits}} ... {{/each}}
  - {{id}}
  - {{name}}
  - {{repoName}}
  - {{streamName}}
  - {{status}}
  - {{durationMinutes}}
  - {{durationTime}}
  - {{notes}}

Metrics:
- {{metrics.sessionCount}}
- {{metrics.workedSeconds}}
- {{metrics.workedTime}}
- {{metrics.restSeconds}}
- {{metrics.restTime}}
- {{metrics.burnout.level}}
- {{metrics.burnout.score}}
`
