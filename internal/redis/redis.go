package redis

import (
	"context"
	"database/sql"

	"github.com/vald3mare/ML-workbench/internal/annotation"
	//_ "github.com/lib/pq"
)

type RedisAnnotationRepo struct {
	db *sql.DB // ну тут тип бд редис, но для примера оставим sql.DB
}

func NewRedisRepo(dsn string) (*RedisAnnotationRepo, error) {
	db, err := sql.Open("redis", dsn) // ну тут тип бд редис, но для примера оставим sql.Open
	if err != nil {
		return nil, err
	}

	return &RedisAnnotationRepo{db: db}, nil
}

func (r *RedisAnnotationRepo) Save(ctx context.Context, annotation *annotation.Annotation) error {
	// реализация сохранения аннотации в редисе
	return nil
}

func (r *RedisAnnotationRepo) Delete(ctx context.Context, annotationID string) error {
	// реализация удаления аннотации из редиса
	return nil
}

func (r *RedisAnnotationRepo) GetByImage(ctx context.Context, imageID string) ([]*annotation.Annotation, error) {
	// реализация получения аннотаций по id изображения из редиса
	return nil, nil
}

func (r *RedisAnnotationRepo) GetByID(ctx context.Context, annotationID string) (*annotation.Annotation, error) {
	// реализация получения аннотации по id из редиса
	return nil, nil
}
