package mysql

import (
	"errors"
	"fmt"
	"github.com/klwxsrx/expense-tracker/pkg/common/app/logger"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const (
	migrationLock      = "perform_migration_lock"
	migrationTableName = "migration"

	migrationFileRegexString = "^([0-9]+).sql$"
	migrationIDVariable      = "%migration_id%"
	migrationFileNamePattern = migrationIDVariable + ".sql"
)

type Migration struct {
	client                 TransactionalClient
	logger                 logger.Logger
	migrationDirectoryPath string
}

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

	return m.performMigrations(m.migrationDirectoryPath)
}

func (m *Migration) createMigrationTableIfNotExists() error {
	_, err := m.client.Exec(
		"CREATE TABLE IF NOT EXISTS `" + migrationTableName + "` (" +
			"id int," +
			"PRIMARY KEY (id)" +
			") ENGINE=InnoDB CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	return err
}

func (m *Migration) performMigrations(migrationDirectoryPath string) error {
	performedMigrationIDs, err := getPerformedMigrationIDs(m.client)
	if err != nil {
		return err
	}

	fileMigrationIDs, err := getFileMigrationIDs(migrationDirectoryPath)
	if err != nil {
		return err
	}

	for _, migrationID := range fileMigrationIDs {
		if !performedMigrationIDs[migrationID] {
			m.logger.Info("execute migration #", migrationID)
			migrationSql, err := getMigrationSql(migrationDirectoryPath, migrationID)
			if err != nil {
				m.logger.WithError(err).Error(fmt.Sprintf("failed to obtain migration #%d sql", migrationID))
				return err
			}
			tx, err := m.client.Begin()
			if err != nil {
				return err
			}
			err = performMigration(tx, migrationSql, migrationID)
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

func getFileMigrationIDs(migrationDirectoryPath string) ([]int, error) {
	files, err := ioutil.ReadDir(migrationDirectoryPath)
	if err != nil {
		return nil, err
	}
	var result []int
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		migrationID, err := getMigrationIDFromFileName(file.Name())
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

func getMigrationSql(migrationDirectoryPath string, migrationID int) (string, error) {
	migrationFilePath := filepath.Join(migrationDirectoryPath, getMigrationFileNameFromMigrationID(migrationID))
	content, err := ioutil.ReadFile(migrationFilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getMigrationFileNameFromMigrationID(migrationID int) string {
	return strings.Replace(migrationFileNamePattern, migrationIDVariable, strconv.Itoa(migrationID), 1)
}

func performMigration(client Client, sql string, migrationID int) error {
	if sql == "" {
		return errors.New("empty migration")
	}
	err := createMigrationRecord(client, migrationID)
	if err != nil {
		return err
	}
	_, err = client.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createMigrationRecord(client Client, migrationID int) error {
	_, err := client.Exec("INSERT INTO `"+migrationTableName+"` VALUES (?)", migrationID)
	return err
}

func NewMigration(client TransactionalClient, logger logger.Logger, migrationDirectoryPath string) (*Migration, error) {
	if migrationFileRegexError != nil {
		return nil, migrationFileRegexError
	}
	return &Migration{client, logger, migrationDirectoryPath}, nil
}
