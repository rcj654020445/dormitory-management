// Package repository provides data access layer implementations.
// Layer 1: Depends on types (Layer 0) and pkg/database.
package repository

import (
	"context"

	"github.com/example/dormitory-management/internal/types"
	"github.com/example/dormitory-management/pkg/database"
)

// StudentRepository handles student data access.
type StudentRepository interface {
	Create(ctx context.Context, student *types.Student) error
	GetByID(ctx context.Context, id string) (*types.Student, error)
	GetByStudentID(ctx context.Context, studentID string) (*types.Student, error)
	List(ctx context.Context, page, pageSize int) ([]*types.Student, int, error)
	Update(ctx context.Context, student *types.Student) error
	Delete(ctx context.Context, id string) error
	AllocateRoom(ctx context.Context, studentID, roomID string) error
	VacateRoom(ctx context.Context, studentID string) error
}

type studentRepository struct {
	db *database.PostgresDB
}

// NewStudentRepository creates a new StudentRepository.
func NewStudentRepository(db *database.PostgresDB) StudentRepository {
	return &studentRepository{db: db}
}

func (r *studentRepository) Create(ctx context.Context, student *types.Student) error {
	query := `
		INSERT INTO students (id, student_no, name, gender, phone, email, major, grade, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, NOW(), NOW())
	`
	return r.db.Exec(ctx, query,
		student.ID,
		student.StudentID,
		student.Name,
		student.Gender,
		student.Phone,
		student.Email,
		student.Major,
		student.Grade,
		student.Status,
	)
}

func (r *studentRepository) GetByID(ctx context.Context, id string) (*types.Student, error) {
	query := `
		SELECT id, student_no, name, gender, phone, email, major, grade, room_id, status, created_at, updated_at
		FROM students WHERE id = $1
	`
	row := r.db.QueryRow(ctx, query, id)
	return r.scanStudent(row)
}

func (r *studentRepository) GetByStudentID(ctx context.Context, studentID string) (*types.Student, error) {
	query := `
		SELECT id, student_no, name, gender, phone, email, major, grade, room_id, status, created_at, updated_at
		FROM students WHERE student_no = $1
	`
	row := r.db.QueryRow(ctx, query, studentID)
	return r.scanStudent(row)
}

func (r *studentRepository) List(ctx context.Context, page, pageSize int) ([]*types.Student, int, error) {
	offset := (page - 1) * pageSize

	countQuery := `SELECT COUNT(*) FROM students`
	var total int
	if err := r.db.QueryRow(ctx, countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	query := `
		SELECT id, student_no, name, gender, phone, email, major, grade, room_id, status, created_at, updated_at
		FROM students
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	rows, err := r.db.Query(ctx, query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var students []*types.Student
	for rows.Next() {
		student, err := r.scanStudentRow(rows)
		if err != nil {
			return nil, 0, err
		}
		students = append(students, student)
	}

	return students, total, rows.Err()
}

func (r *studentRepository) Update(ctx context.Context, student *types.Student) error {
	query := `
		UPDATE students SET
			name = $2, gender = $3, phone = $4, email = $5,
			major = $6, grade = $7, room_id = $8,
			status = $9, updated_at = NOW()
		WHERE id = $1
	`
	return r.db.Exec(ctx, query,
		student.ID, student.Name, student.Gender, student.Phone,
		student.Email, student.Major, student.Grade, student.RoomID,
		student.Status,
	)
}

func (r *studentRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM students WHERE id = $1`
	return r.db.Exec(ctx, query, id)
}

func (r *studentRepository) AllocateRoom(ctx context.Context, studentID, roomID string) error {
	query := `UPDATE students SET room_id = $2, updated_at = NOW() WHERE id = $1`
	return r.db.Exec(ctx, query, studentID, roomID)
}

func (r *studentRepository) VacateRoom(ctx context.Context, studentID string) error {
	query := `UPDATE students SET room_id = NULL, updated_at = NOW() WHERE id = $1`
	return r.db.Exec(ctx, query, studentID)
}

func (r *studentRepository) scanStudent(row interface{ Scan(...any) error }) (*types.Student, error) {
	var student types.Student
	err := row.Scan(
		&student.ID, &student.StudentID, &student.Name, &student.Gender,
		&student.Phone, &student.Email, &student.Major, &student.Grade,
		&student.RoomID, &student.Status,
		&student.CreatedAt, &student.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}

func (r *studentRepository) scanStudentRow(rows interface{ Scan(...any) error }) (*types.Student, error) {
	var student types.Student
	err := rows.Scan(
		&student.ID, &student.StudentID, &student.Name, &student.Gender,
		&student.Phone, &student.Email, &student.Major, &student.Grade,
		&student.RoomID, &student.Status,
		&student.CreatedAt, &student.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &student, nil
}