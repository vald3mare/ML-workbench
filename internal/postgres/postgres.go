package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/vald3mare/ML-workbench/internal/annotation"
)

// PostgresService это реализация интерфейса AnnotationRepository для работы с PostgreSQL базой данных, по сути своей слой для доступа к данным,
// который будет использоваться в сервисе для работы с бизнес логикой, а не с данными напрямую
type PostgresService struct {
	db *sql.DB
}

// конструктор для PostgresService, он будет принимать строку подключения к базе данных
// и возвращать сервис, который будет использовать эту базу данных для работы с данными
func NewPostgresRepo(dsn string) (*PostgresService, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	// Проверяем соединение
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Создаем таблицу если она не существует
	if err := initDB(db); err != nil {
		return nil, err
	}

	return &PostgresService{db: db}, nil
}

// initDB инициализирует схему базы данных
func initDB(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS annotations (
		id VARCHAR(255) PRIMARY KEY,
		image_id VARCHAR(255) NOT NULL,
		title VARCHAR(255) NOT NULL,
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_annotations_image_id ON annotations(image_id);
	`

	_, err := db.Exec(schema)
	return err
}

func (p *PostgresService) Close() error {
	return p.db.Close()
}

// реализация метода Save для сохранения аннотации в базе данных
func (p *PostgresService) Save(ctx context.Context, annotation *annotation.Annotation) error {
	_, err := p.db.ExecContext(
		ctx,
		"INSERT INTO annotations (id, image_id, title, description) VALUES ($1, $2, $3, $4)",
		annotation.ID,
		annotation.ImageID,
		annotation.Shapes[0], // для примера сохраняем тип первой фигуры в описание, в реальной реализации будет логика для сохранения всех данных
	)
	if err != nil {
		fmt.Printf("Ошибка сохранения: %v", err)
		return err
	}
	return nil
}

// реализация метода Delete для удаления аннотации из базы данных, он будет использовать метод ExecContext для выполнения SQL запроса на удаление данных из таблицы annotations по id
func (p *PostgresService) Delete(ctx context.Context, annotationID string) error {
	_, err := p.db.ExecContext(ctx, "DELETE FROM annotations WHERE id = $1", annotationID)
	if err != nil {
		fmt.Printf("Ошибка удаления: %v", err)
		return err
	}
	return nil
}

func (p *PostgresService) GetByImage(ctx context.Context, imageID string) ([]*annotation.Annotation, error) {
	rows, err := p.db.QueryContext(ctx, "SELECT id, image_id, title, description FROM annotations WHERE image_id = $1", imageID)
	if err != nil {
		fmt.Printf("Ошибка получения по изображению: %v", err)
		return nil, err
	}
	defer rows.Close()

	var annotations []*annotation.Annotation
	for rows.Next() {
		var ann annotation.Annotation
		if err := rows.Scan(&ann.ID, &ann.ImageID, &ann.Shapes); err != nil {
			fmt.Printf("Ошибка сканирования: %v", err)
			return nil, err
		}
		annotations = append(annotations, &ann)
	}
	return annotations, nil
}

func (p *PostgresService) GetByID(ctx context.Context, annotationID string) (*annotation.Annotation, error) {
	row := p.db.QueryRowContext(ctx, "SELECT id, image_id, title, description FROM annotations WHERE id = $1", annotationID)

	var ann annotation.Annotation
	if err := row.Scan(&ann.ID, &ann.ImageID, &ann.Shapes); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // аннотация не найдена
		}
		fmt.Printf("Ошибка получения по ID: %v", err)
		return nil, err
	}
	return &ann, nil
}
