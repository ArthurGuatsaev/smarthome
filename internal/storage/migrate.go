package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"io/fs"
	"sort"
	"strconv"
	"strings"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func Migrate(ctx context.Context, db *sql.DB) error {
	// Гарантируем таблицу версий (на всякий случай, если миграция 0001 не успела)
	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS schema_migrations (version INTEGER PRIMARY KEY)`); err != nil {
		return err
	}

	cur, err := currentVersion(ctx, db)
	if err != nil {
		return err
	}

	files, err := fs.Glob(migrationsFS, "migrations/*.sql")
	if err != nil {
		return err
	}
	sort.Strings(files)

	for _, f := range files {
		v, err := parseVersion(f) // migrations/0001_init.sql -> 1
		if err != nil {
			return err
		}
		if v <= cur {
			continue
		}

		sqlBytes, err := migrationsFS.ReadFile(f)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		if _, err := tx.ExecContext(ctx, string(sqlBytes)); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("migration %s failed: %w", f, err)
		}

		if _, err := tx.ExecContext(ctx, `INSERT INTO schema_migrations(version) VALUES (?)`, v); err != nil {
			_ = tx.Rollback()
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}

		cur = v
	}

	return nil
}

func currentVersion(ctx context.Context, db *sql.DB) (int, error) {
	row := db.QueryRowContext(ctx, `SELECT COALESCE(MAX(version), 0) FROM schema_migrations`)
	var v int
	if err := row.Scan(&v); err != nil {
		return 0, err
	}
	return v, nil
}

func parseVersion(path string) (int, error) {
	// "migrations/0001_init.sql" -> "0001"
	base := path[strings.LastIndex(path, "/")+1:]
	if len(base) < 4 {
		return 0, fmt.Errorf("bad migration filename: %s", path)
	}
	n, err := strconv.Atoi(base[:4])
	if err != nil {
		return 0, fmt.Errorf("bad migration version in %s: %w", path, err)
	}
	return n, nil
}
