package platform

import (
	"context"
	"database/sql"
	"strings"

	"github.com/closeloopautomous/arms/internal/adapters/chromemstore"
	"github.com/closeloopautomous/arms/internal/adapters/sqlite"
	"github.com/closeloopautomous/arms/internal/config"
	"github.com/closeloopautomous/arms/internal/ports"
)

// newKnowledgeRepository selects the knowledge backend and whether Search should use FTS5 query sanitization (SQLite FTS only).
func newKnowledgeRepository(ctx context.Context, cfg config.Config, db *sql.DB) (ports.KnowledgeRepository, bool, error) {
	_ = ctx
	switch strings.ToLower(strings.TrimSpace(cfg.KnowledgeBackend)) {
	case "chromem":
		ks, err := chromemstore.NewKnowledgeStore(cfg)
		if err != nil {
			return nil, false, err
		}
		return ks, false, nil
	default:
		return sqlite.NewKnowledgeStore(db), true, nil
	}
}
