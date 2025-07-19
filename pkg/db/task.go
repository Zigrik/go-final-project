package db

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func AddTask(task *Task) (int64, error) {
	if Database == nil {
		return 0, fmt.Errorf("database not initialized")
	}

	tx, err := Database.Begin()
	if err != nil {
		return 0, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback() // Безопасный откат если не будет Commit

	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat)`
	res, err := tx.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
	)
	if err != nil {
		return 0, fmt.Errorf("failed to insert task: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("failed to commit transaction: %w", err)
	}
	return id, nil
}

func Tasks(limit int, search, dateForm string) ([]*Task, error) {
	if Database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	var rows *sql.Rows
	var err error
	var date time.Time

	if date, err = time.Parse("02.01.2006", search); err == nil {
		query := `SELECT id, date, title, comment, repeat 
                 FROM scheduler 
				 WHERE date = ? 
                 ORDER BY date ASC, id ASC LIMIT ?`
		rows, err = Database.Query(query, date.Format(dateForm), limit)
	} else if search != "" {
		query := `SELECT id, date, title, comment, repeat 
				  FROM scheduler 
				  WHERE title 
				  LIKE ? OR comment LIKE ? 
				  ORDER BY date ASC, id ASC LIMIT ?`
		search = "%" + search + "%"
		rows, err = Database.Query(query, search, search, limit)
	} else {
		query := `SELECT id, date, title, comment, repeat 
                 FROM scheduler 
                 ORDER BY date ASC, id ASC LIMIT ?`
		rows, err = Database.Query(query, limit)
	}

	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat); err != nil {
			return nil, fmt.Errorf("scan db failed: %w", err)
		}
		tasks = append(tasks, &task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while going through the records: %w", err)
	}

	if tasks == nil {
		tasks = make([]*Task, 0)
	}

	return tasks, nil
}

func GetTask(id string) (*Task, error) {
	if Database == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	if id == "" {
		return nil, fmt.Errorf("task ID cannot be empty")
	}

	_, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("invalid task ID format")
	}

	query := `SELECT id, date, title, comment, repeat 
			  FROM scheduler 
			  WHERE id = ?`

	var task Task
	err = Database.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)

	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func UpdateTask(task *Task) error {
	if Database == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id `
	res, err := tx.Exec(query,
		sql.Named("date", task.Date),
		sql.Named("title", task.Title),
		sql.Named("comment", task.Comment),
		sql.Named("repeat", task.Repeat),
		sql.Named("id", task.ID),
	)
	if err != nil {
		return fmt.Errorf("failed to execute update: %w", err)
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf(`incorrect id for updating task`)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func DeleteTask(id string) error {
	if Database == nil {
		return fmt.Errorf("database not initialized")
	}

	tx, err := Database.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	query := `DELETE FROM scheduler WHERE id = :id`
	result, err := tx.Exec(query, sql.Named("id", id))
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %s not found", id)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}
