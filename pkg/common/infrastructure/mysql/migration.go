package mysql

import (
	"errors"
	"fmt"
	"io/fs"
	"regexp"
	"strconv"
	"strings"

	"github.com/klwxsrx/budget-tracker/pkg/common/app/log"
)

const (
	migrationLock      = "perform_migration_lock"
	migrationTableName = "migration"

	migrationFileRegexString = `^(\d+).sql$`
	migrationFileNamePattern = "%v.sql"

	querySeparator = ";\n"
)

type Migration struct {
	client     TransactionalClient
	logger     log.Logger
	migrations fs.ReadDirFS
}

// nolint:gochecknoglobals
var migrationFileRegex, migrationFileRegexError = regexp.Compile(migrationFileRegexString)

func (m *Migration) Migrate() error {
	lock := NewLock(m.client, migrationLock)
	err := lock.Get()
	if err != nil {
		return err
	}
	defer lock.Release()

	err = m.createMigrationTableIfNotExists()
	if err != nil {
		return err
	}

	return m.performMigrations(m.migrations)
}

func (m *Migration) createMigrationTableIfNotExists() error {
	_, err := m.client.Exec(
		"CREATE TABLE IF NOT EXISTS `" + migrationTableName + "` (" +
			"id int," +
			"PRIMARY KEY (id)" +
			") ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	return err
}

func (m *Migration) performMigrations(migrations fs.ReadDirFS) error {
	performedMigrationIDs, err := getPerformedMigrationIDs(m.client)
	if err != nil {
		return err
	}

	fileMigrationIDs, err := getFileMigrationIDs(migrations)
	if err != nil {
		return err
	}

	for _, migrationID := range fileMigrationIDs {
		if performedMigrationIDs[migrationID] {
			continue
		}

		m.logger.Info(fmt.Sprintf("execute migration #%v", migrationID))
		migrationSQL, err := getMigrationSQL(migrations, migrationID)
		if err != nil {
			m.logger.WithError(err).Error(fmt.Sprintf("failed to obtain migration #%d sql", migrationID))
			return err
		}
		tx, err := m.client.Begin()
		if err != nil {
			return err
		}
		err = performMigration(tx, migrationSQL, migrationID)
		if err != nil {
			_ = tx.Rollback()
			m.logger.WithError(err).Error(fmt.Sprintf("migration #%d failed", migrationID))
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
	}
	return nil
}

func getPerformedMigrationIDs(client Client) (map[int]bool, error) {
	var ids []int
	err := client.Select(&ids, "SELECT id FROM `"+migrationTableName+"`")
	if err != nil {
		return nil, err
	}
	result := make(map[int]bool, len(ids))
	for _, id := range ids {
		result[id] = true
	}
	return result, nil
}

func getFileMigrationIDs(migrations fs.ReadDirFS) ([]int, error) {
	entries, err := migrations.ReadDir(".")
	if err != nil {
		return nil, err
	}
	result := make([]int, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		migrationID, err := getMigrationIDFromFileName(entry.Name())
		if err != nil {
			continue
		}
		result = append(result, migrationID)
	}
	return result, nil
}

func getMigrationIDFromFileName(fileName string) (int, error) {
	subMatch := migrationFileRegex.FindStringSubmatch(fileName)
	if subMatch == nil {
		return 0, errors.New("doesn't match the pattern")
	}
	migrationID, err := strconv.Atoi(subMatch[1])
	if err != nil {
		return 0, err
	}
	return migrationID, nil
}

func getMigrationSQL(migrations fs.ReadDirFS, migrationID int) (string, error) {
	content, err := fs.ReadFile(migrations, getMigrationFileNameFromMigrationID(migrationID))
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getMigrationFileNameFromMigrationID(migrationID int) string {
	return fmt.Sprintf(migrationFileNamePattern, strconv.Itoa(migrationID))
}

func performMigration(client Client, sql string, migrationID int) error {
	if sql == "" {
		return errors.New("empty migration")
	}
	err := createMigrationRecord(client, migrationID)
	if err != nil {
		return err
	}

	queries := splitToQueries(sql)
	for _, query := range queries {
		_, err = client.Exec(query)
		if err != nil {
			return err
		}
	}
	return nil
}

func splitToQueries(sql string) []string {
	return strings.Split(sql, querySeparator)
}

func createMigrationRecord(client Client, migrationID int) error {
	_, err := client.Exec("INSERT INTO `"+migrationTableName+"` VALUES (?)", migrationID)
	return err
}

func NewMigration(client TransactionalClient, logger log.Logger, migrations fs.ReadDirFS) (*Migration, error) {
	if migrationFileRegexError != nil {
		return nil, migrationFileRegexError
	}
	return &Migration{client, logger, migrations}, nil
}
