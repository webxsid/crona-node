# CSV Export Spec

The CSV export uses a JSON spec instead of Handlebars.

Supported keys:
- `headers`: ordered CSV header labels
- `columns`: ordered field names pulled from each session row

Common row fields:
- `id`
- `issueId`
- `issueTitle`
- `repoName`
- `streamName`
- `startTime`
- `endTime`
- `durationSeconds`
- `summary`
