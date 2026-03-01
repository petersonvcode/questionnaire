package database

import (
	"embed"
	"sort"
	"strconv"
	"strings"
	"time"

	"log/slog"
)

//go:embed migrations/*.sql
var migrationFolder embed.FS

type MigrationRun struct {
	Name         string `json:"name"`
	RunAt        string `json:"runAt"`
	Success      bool   `json:"success"`
	ErrorMessage string `json:"errorMessage"`
}

type MigrationFile struct {
	Name    string `json:"name"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

func (s *DatabaseService) Migrate() error {
	slog.Info("Migrating database ...")
	err := createMigrationRunsTable(s)
	if err != nil {
		slog.Error("failed to create migration runs table", "err", err.Error())
		return err
	}

	files, err := migrationFolder.ReadDir("migrations")
	if err != nil {
		slog.Error("failed to read migration folder", "err", err.Error())
		return err
	}

	sqlFiles := make([]string, 0)
	for _, file := range files {
		name := file.Name()
		if !file.IsDir() && strings.HasSuffix(name, ".sql") {
			sqlFiles = append(sqlFiles, name)
		}
	}

	migrationFiles := make([]MigrationFile, 0, len(sqlFiles))

	for _, file := range sqlFiles {
		slog.Debug("Reading migration file: " + file)
		sqlBytes, err := migrationFolder.ReadFile("migrations/" + file)
		if err != nil {
			slog.Error("failed to read migration file", "err", err.Error())
			return err
		}
		sql := string(sqlBytes)

		migrationFiles = append(migrationFiles, MigrationFile{
			Name:    file,
			Path:    "migrations/" + file,
			Content: sql,
		})
	}
	sort.Slice(migrationFiles, func(i, j int) bool {
		iPrefix := strings.Split(migrationFiles[i].Name, "-")[0]
		jPrefix := strings.Split(migrationFiles[j].Name, "-")[0]
		iNumber, err := strconv.Atoi(iPrefix)
		if err != nil {
			return false
		}
		jNumber, err := strconv.Atoi(jPrefix)
		if err != nil {
			return false
		}
		return iNumber < jNumber
	})
	slog.Debug("Migration files", "migrationFiles", getMigrationNames(migrationFiles), "count", len(migrationFiles))

	alreadyRun := getMigrationAlreadyRan(s)
	slog.Debug("Already ran migrations", "alreadyRunNames", getAlreadyRunNames(alreadyRun), "count", len(alreadyRun))
	for _, migrationFile := range migrationFiles {
		alreadyRun, migrationRun := hasAlreadyRun(alreadyRun, migrationFile)
		if alreadyRun {
			slog.Debug("Migration file already ran", "name", migrationFile.Name, "when", migrationRun.RunAt)
			continue
		}
		err = runMigration(s, migrationFile)
		if err != nil {
			slog.Error("failed to run migration file", "err", err.Error(), "name", migrationFile.Name, "sql", migrationFile.Content)
			return err
		}
		slog.Debug("Migration file ran", "name", migrationFile.Name, "when", migrationRun.RunAt)
	}

	slog.Info("Done Migrating database !!")
	return nil
}

func createMigrationRunsTable(s *DatabaseService) error {
	query := `
CREATE TABLE IF NOT EXISTS migration_runs (
  name TEXT PRIMARY KEY,
  run_at TEXT NOT NULL,
  success BOOLEAN NOT NULL,
  error_message TEXT
);`

	_, err := s.db.Exec(query)
	if err != nil {
		slog.Error("failed to create migration runs table", "err", err.Error(), "query", query)
		return err
	}
	return nil
}

func hasAlreadyRun(alreadyRun []MigrationRun, file MigrationFile) (bool, MigrationRun) {
	for _, migrationRun := range alreadyRun {
		if migrationRun.Name == file.Name && migrationRun.Success {
			return true, migrationRun
		}
	}
	return false, MigrationRun{}
}

func getMigrationNames(migrationFiles []MigrationFile) string {
	migrationNames := ""
	for _, file := range migrationFiles {
		migrationNames += file.Name + ","
	}

	if len(migrationNames) > 0 {
		migrationNames = migrationNames[0 : len(migrationNames)-1]
	}
	return migrationNames
}

func getAlreadyRunNames(alreadyRun []MigrationRun) string {
	migrationNames := ""
	for _, migrationRun := range alreadyRun {
		migrationNames += migrationRun.Name + ","
	}
	if len(migrationNames) > 0 {
		migrationNames = migrationNames[0 : len(migrationNames)-1]
	}
	return migrationNames
}

func runMigration(s *DatabaseService, migrationFile MigrationFile) error {
	run := MigrationRun{
		Name:         migrationFile.Name,
		RunAt:        time.Now().Format(time.RFC3339),
		Success:      true,
		ErrorMessage: "",
	}

	_, err := s.db.Exec(migrationFile.Content)
	if err != nil {
		slog.Error("failed to execute migration file", "err", err.Error(), "name", migrationFile.Name, "sql", migrationFile.Content)
		run.Success = false
		run.ErrorMessage = err.Error()
	}

	insertOrUpdateQuery := `
		INSERT INTO migration_runs (name, run_at, success, error_message)
		VALUES (?, ?, ?, ?)
		ON CONFLICT(name) DO UPDATE SET
			run_at=excluded.run_at,
			success=excluded.success,
			error_message=excluded.error_message
	`
	_, migrationErr := s.db.Exec(insertOrUpdateQuery, run.Name, run.RunAt, run.Success, run.ErrorMessage)
	if migrationErr != nil {
		slog.Error("failed to insert or update migration run", "err", migrationErr.Error(), "name", migrationFile.Name, "sql", migrationFile.Content)
	}

	if err != nil || migrationErr != nil {
		return err
	}

	return nil
}

func getMigrationAlreadyRan(s *DatabaseService) []MigrationRun {
	rows, err := s.db.Query("SELECT name, run_at, success, error_message FROM migration_runs")
	if err != nil {
		slog.Error("failed to get migration already run", "err", err.Error())
		return nil
	}
	defer rows.Close()

	migrationRuns := make([]MigrationRun, 0)
	for rows.Next() {
		migrationRun := MigrationRun{}
		err = rows.Scan(&migrationRun.Name, &migrationRun.RunAt, &migrationRun.Success, &migrationRun.ErrorMessage)
		if err != nil {
			slog.Error("failed to scan migration run", "err", err.Error())
			return nil
		}
		migrationRuns = append(migrationRuns, migrationRun)
	}
	return migrationRuns
}