package mysql

import (
	"errors"
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
	performedMigrationIds, err := getPerformedMigrationIDs(m.client)
	if err != nil {
		return err
	}

	fileMigrationIds, err := getFileMigrationIds(migrationDirectoryPath)
	if err != nil {
		return err
	}

	for _, migrationId := range fileMigrationIds {
		if !performedMigrationIds[migrationId] {
			migrationSql, err := getMigrationSql(migrationDirectoryPath, migrationId)
			if err != nil {
				return err
			}
			tx, err := m.client.Begin()
			if err != nil {
				return err
			}
			err = performMigration(tx, migrationSql, migrationId)
			if err != nil {
				_ = tx.Rollback()
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

func getFileMigrationIds(migrationDirectoryPath string) ([]int, error) {
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

func getMigrationSql(migrationDirectoryPath string, migrationId int) (string, error) {
	migrationFilePath := filepath.Join(migrationDirectoryPath, getMigrationFileNameFromMigrationID(migrationId))
	content, err := ioutil.ReadFile(migrationFilePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func getMigrationFileNameFromMigrationID(migrationID int) string {
	return strings.Replace(migrationFileNamePattern, migrationIDVariable, strconv.Itoa(migrationID), 1)
}

func performMigration(client Client, sql string, migrationId int) error {
	if sql == "" {
		return errors.New("empty migration")
	}
	err := createMigrationRecord(client, migrationId)
	if err != nil {
		return err
	}
	_, err = client.Exec(sql)
	if err != nil {
		return err
	}
	return nil
}

func createMigrationRecord(client Client, migrationId int) error {
	_, err := client.Exec("INSERT INTO `"+migrationTableName+"` VALUES (?)", migrationId)
	return err
}

func NewMigration(client TransactionalClient, migrationDirectoryPath string) (*Migration, error) {
	if migrationFileRegexError != nil {
		return nil, migrationFileRegexError
	}
	return &Migration{client, migrationDirectoryPath}, nil
}
