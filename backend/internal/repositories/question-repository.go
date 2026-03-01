package repositories

import (
	"database/sql"
	"log/slog"
	"slices"
	"strings"

	"github.com/petersonvcode/questionnaire/backend/internal/domain"
)

type QuestionRepository struct {
	db DatabaseService
}

func NewQuestionRepository(db DatabaseService) *QuestionRepository {
	return &QuestionRepository{db: db}
}

func (r *QuestionRepository) InsertQuestions(questions []domain.QuestionCreationItemRequest) ([]*domain.Question, error) {
	tx, err := r.db.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	results := make([]*domain.Question, 0, len(questions))
	for _, question := range questions {
		dbQuestion, err := r.insertQuestion(tx, question)
		if err != nil {
			return nil, err
		}
		options, err := r.insertOptions(tx, dbQuestion.ID, question.Options)
		if err != nil {
			return nil, err
		}
		tags, err := r.insertTags(tx, dbQuestion.ID, question.Tags)
		if err != nil {
			return nil, err
		}
		dbQuestion.Options = options
		dbQuestion.Tags = tags
		results = append(results, dbQuestion)
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return results, nil
}

func (r *QuestionRepository) insertQuestion(tx *sql.Tx, question domain.QuestionCreationItemRequest) (*domain.Question, error) {
	query := `INSERT INTO questions (text, type) VALUES (?, ?) RETURNING id;`
	row := tx.QueryRow(query, question.Text, question.Type)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	return &domain.Question{ID: id, Text: question.Text, Type: question.Type}, nil
}

func (r *QuestionRepository) insertOptions(tx *sql.Tx, questionID int64, options []domain.QuestionOptionCreationRequest) ([]domain.QuestionOption, error) {
	if len(options) == 0 {
		return nil, nil
	}
	var query strings.Builder
	query.WriteString(`INSERT INTO question_options (question_id, text, correct) VALUES `)
	params := make([]any, 0, len(options)*3)
	for i, option := range options {
		if i > 0 {
			query.WriteString(", ")
		}
		query.WriteString("(?, ?, ?)")
		params = append(params, questionID, option.Text, option.Correct)
	}
	query.WriteString(" RETURNING id;")

	rows, err := tx.Query(query.String(), params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.QuestionOption, 0, len(options))
	index := 0
	var id int64
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		results = append(results, domain.QuestionOption{ID: id, QuestionID: questionID, Text: options[index].Text, Correct: options[index].Correct})
		index++
	}
	return results, nil
}

func (r *QuestionRepository) insertTags(tx *sql.Tx, questionID int64, tags []string) ([]domain.QuestionTag, error) {
	exists, notExists, err := r.checkTagExistence(tx, tags)
	if err != nil {
		return nil, err
	}

	// Insert new tags into tags table
	tagIDs := make(map[string]int64)
	for _, qt := range exists {
		tagIDs[qt.Tag] = qt.ID
	}
	for _, tag := range notExists {
		row := tx.QueryRow(`INSERT INTO tags (tag, "group") VALUES (?, ?) RETURNING id;`, tag, "")
		var id int64
		err := row.Scan(&id)
		if err != nil {
			return nil, err
		}
		tagIDs[tag] = id
	}

	// Insert into question_tags junction table
	if len(tagIDs) > 0 {
		query := `INSERT INTO question_tags (question_id, tag_id) VALUES `
		params := make([]any, 0, len(tags)*2)
		for i, tag := range tags {
			if i > 0 {
				query += ", "
			}
			query += "(?, ?)"
			params = append(params, questionID, tagIDs[tag])
		}
		_, err = tx.Exec(query, params...)
		if err != nil {
			return nil, err
		}
	}

	// Build result in order
	results := make([]domain.QuestionTag, 0, len(tags))
	for _, tag := range tags {
		results = append(results, domain.QuestionTag{ID: tagIDs[tag], Tag: tag, Group: ""})
	}
	return results, nil
}

func (r *QuestionRepository) checkTagExistence(tx *sql.Tx, tags []string) (exists []domain.QuestionTag, notExists []string, err error) {
	if len(tags) == 0 {
		return nil, nil, nil
	}
	query := `SELECT id, tag, "group" FROM tags WHERE tag IN (`
	args := make([]any, len(tags))
	for i := range tags {
		if i > 0 {
			query += ", "
		}
		query += "?"
		args[i] = tags[i]
	}
	query += ")"

	rows, err := tx.Query(query, args...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	exists = make([]domain.QuestionTag, 0)
	notExists = make([]string, 0)
	seenTags := make(map[string]bool)
	for _, t := range tags {
		seenTags[t] = false
	}
	var id int64
	var tag string
	var group string
	for rows.Next() {
		err = rows.Scan(&id, &tag, &group)
		if err != nil {
			slog.Error("Failed to scan tag", "error", err)
			return nil, nil, err
		}
		if slices.Contains(tags, tag) {
			exists = append(exists, domain.QuestionTag{ID: id, Tag: tag, Group: group})
			seenTags[tag] = true
		}
	}
	for _, t := range tags {
		if !seenTags[t] {
			notExists = append(notExists, t)
		}
	}

	return exists, notExists, nil
}

func (r *QuestionRepository) GetQuestionByID(id int64) (*domain.Question, error) {
	tx, err := r.db.BeginTransaction()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	query := `SELECT id, text, type FROM questions WHERE id = ?;`
	row := tx.QueryRow(query, id)

	question := domain.Question{}
	err = row.Scan(&question.ID, &question.Text, &question.Type)
	if err != nil {
		if err == sql.ErrNoRows {
			slog.Error("Question not found", "error", err)
			return nil, nil
		}

		slog.Error("Failed to scan question", "error", err)
		return nil, err
	}
	question.Options, err = r.getQuestionOptionsByQuestionID(tx, id)
	if err != nil {
		return nil, err
	}
	question.Tags, err = r.getQuestionTagsByQuestionID(tx, id)
	if err != nil {
		return nil, err
	}

	return &question, nil
}

func (r *QuestionRepository) getQuestionOptionsByQuestionID(tx *sql.Tx, questionID int64) ([]domain.QuestionOption, error) {
	query := `SELECT id, question_id, text, correct FROM question_options WHERE question_id = ?;`
	rows, err := tx.Query(query, questionID)
	if err != nil {
		slog.Error("Failed to query question options", "error", err)
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.QuestionOption, 0)
	for rows.Next() {
		opt := domain.QuestionOption{}
		err = rows.Scan(&opt.ID, &opt.QuestionID, &opt.Text, &opt.Correct)
		if err != nil {
			slog.Error("Failed to scan question option", "error", err)
			return nil, err
		}
		results = append(results, opt)
	}

	return results, nil
}

func (r *QuestionRepository) getQuestionTagsByQuestionID(tx *sql.Tx, questionID int64) ([]domain.QuestionTag, error) {
	query := `SELECT tag_id FROM question_tags WHERE question_id = ?`
	rows, err := tx.Query(query, questionID)
	if err != nil {
		slog.Error("Failed to query question tags IDs", "error", err)
		return nil, err
	}
	defer rows.Close()

	ids := make([]int64, 0)
	for rows.Next() {
		var id int64
		err = rows.Scan(&id)
		if err != nil {
			slog.Error("Failed to scan question tag ID", "error", err)
			return nil, err
		}
		ids = append(ids, id)
	}

	query = `SELECT id, tag, "group" FROM tags WHERE id IN (`
	params := make([]any, len(ids))
	for i := range ids {
		query += "?, "
		params[i] = ids[i]
	}
	query = strings.TrimSuffix(query, ", ")
	query += ");"

	rows, err = tx.Query(query, params...)
	if err != nil {
		slog.Error("Failed to query question tags", "error", err, "query", query)
		return nil, err
	}
	defer rows.Close()

	results := make([]domain.QuestionTag, 0)
	for rows.Next() {
		qt := domain.QuestionTag{}
		err = rows.Scan(&qt.ID, &qt.Tag, &qt.Group)
		if err != nil {
			slog.Error("Failed to scan question tag", "error", err)
			return nil, err
		}
		results = append(results, qt)
	}

	return results, nil
}

func (r *QuestionRepository) DeleteQuestionByID(id int64) error {
	tx, err := r.db.BeginTransaction()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	question, err := r.GetQuestionByID(id)
	if err != nil {
		return err
	}
	if question == nil {
		return nil
	}

	_, err = tx.Exec(`DELETE FROM question_options WHERE question_id = ?;`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM question_tags WHERE question_id = ?;`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM tags WHERE id IN (SELECT tag_id FROM question_tags WHERE question_id = ?);`, id)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM questions WHERE id = ?;`, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *QuestionRepository) GetRandomQuestion() (*domain.Question, error) {
	query := `SELECT id FROM questions ORDER BY RANDOM() LIMIT 1;`
	row := r.db.QueryRow(query)
	var id int64
	err := row.Scan(&id)
	if err != nil {
		return nil, err
	}
	return r.GetQuestionByID(id)
}
