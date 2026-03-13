package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/uptrace/bun"
)

func InitSchema(ctx context.Context, db *bun.DB) error {
	models := []any{
		(*RepoModel)(nil),
		(*StreamModel)(nil),
		(*IssueModel)(nil),
		(*SessionModel)(nil),
		(*StashModel)(nil),
		(*OpModel)(nil),
		(*CoreSettingsModel)(nil),
		(*SessionSegmentModel)(nil),
		(*ActiveContextModel)(nil),
		(*ScratchPadMetaModel)(nil),
	}

	for _, model := range models {
		if _, err := db.NewCreateTable().Model(model).IfNotExists().Exec(ctx); err != nil {
			return err
		}
	}

	if err := ensureIssueColumn(ctx, db, "completed_at"); err != nil {
		return err
	}
	if err := ensureIssueColumn(ctx, db, "abandoned_at"); err != nil {
		return err
	}
	for _, table := range []string{"repos", "streams", "issues"} {
		if err := ensurePublicIDColumn(ctx, db, table); err != nil {
			return err
		}
		if err := backfillPublicIDs(ctx, db, table); err != nil {
			return err
		}
	}

	indexes := []string{
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_repos_public_id ON repos (public_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_streams_public_id ON streams (public_id)",
		"CREATE UNIQUE INDEX IF NOT EXISTS idx_issues_public_id ON issues (public_id)",
		"CREATE INDEX IF NOT EXISTS idx_streams_repo_id ON streams (repo_id)",
		"CREATE INDEX IF NOT EXISTS idx_issues_stream_id ON issues (stream_id)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_issue_id ON sessions (issue_id)",
		"CREATE INDEX IF NOT EXISTS idx_stash_repo_id ON stash (repo_id)",
		"CREATE INDEX IF NOT EXISTS idx_stash_stream_id ON stash (stream_id)",
		"CREATE INDEX IF NOT EXISTS idx_stash_issue_id ON stash (issue_id)",
		"CREATE INDEX IF NOT EXISTS idx_ops_entity_entity_id ON ops (entity, entity_id)",
		"CREATE INDEX IF NOT EXISTS idx_repos_user_id ON repos (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_streams_user_id ON streams (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_issues_user_id ON issues (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_stash_user_id ON stash (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_ops_user_id ON ops (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_session_segments_session_id ON session_segments (session_id)",
		"CREATE INDEX IF NOT EXISTS idx_session_segments_user_id ON session_segments (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_active_context_user_id ON active_context (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_active_context_device_id ON active_context (device_id)",
		"CREATE INDEX IF NOT EXISTS idx_scratch_pad_meta_user_id ON scratch_pad_meta (user_id)",
		"CREATE INDEX IF NOT EXISTS idx_scratch_pad_meta_device_id ON scratch_pad_meta (device_id)",
		"CREATE INDEX IF NOT EXISTS idx_scratch_pad_meta_last_opened_at ON scratch_pad_meta (last_opened_at)",
	}

	for _, stmt := range indexes {
		if _, err := db.ExecContext(ctx, stmt); err != nil {
			return err
		}
	}

	return nil
}

func ensureIssueColumn(ctx context.Context, db *bun.DB, columnName string) error {
	rows, err := db.QueryContext(ctx, "PRAGMA table_info('issues')")
	if err != nil {
		return err
	}
	defer rows.Close()

	var (
		cid       int
		name      string
		typ       string
		notnull   int
		dfltValue sql.NullString
		pk        int
	)
	for rows.Next() {
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dfltValue, &pk); err != nil {
			return err
		}
		if name == columnName {
			return nil
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, fmt.Sprintf("ALTER TABLE issues ADD COLUMN %s text", columnName))
	return err
}
