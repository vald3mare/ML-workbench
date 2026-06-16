package annotation

import (
	"context"
	"fmt"
)

// Annotation это доменный слой приложения, который описывает структуру данных и интерфейсы для работы с аннотациями (разметкой).
type Annotation struct {
	ID      string
	ImageID string
	Shapes  []Shape
}

type Shape struct {
	Type   string // "bbox", "polygon", "point"
	Points []Point
	Label  string
}

type Point struct {
	X float64
	Y float64
}

// интерфейс для работы с аннотациями(разметкой)
// объявляем его тут потому что он может использоваться в разных местах, например в сервисе и в контроллере и там будет его реализация, а не в одном месте
// это основной интерфейс для работы с аннотациями, он будет использоваться в сервисе и в контроллере, а реализация будет в другом месте, например в postgres или в памяти
type AnnotationRepository interface {
	Save(ctx context.Context, annotation *Annotation) error
	Delete(ctx context.Context, annotationID string) error
	GetByImage(ctx context.Context, imageID string) ([]*Annotation, error)
	GetByID(ctx context.Context, annotationID string) (*Annotation, error)
}

// это основной сервис для работы с аннотациями, он будет использовать репозиторий для работы с данными,
// а контроллер будет использовать сервис для работы с бизнес логикой, а не с данными напрямую
type AnnotationService struct {
	repo AnnotationRepository
}

// конструктор для сервиса, он будет принимать репозиторий и возвращать сервис, который будет использовать этот репозиторий для работы с данными
func NewAnnotationService(repo AnnotationRepository) *AnnotationService {
	// бизнес логика
	if repo == nil {
		panic("AnnotationRepository cannot be nil")
	}

	if _, ok := repo.(*AnnotationService); ok {
		panic("repo cannot be of type AnnotationService to avoid circular dependency")
	}

	if annotation, err := repo.GetByID(context.Background(), "test"); err != nil {
		panic("repo GetByID method is not working: " + err.Error())
	} else if annotation != nil {
		panic("repo GetByID method should return nil for non-existent annotation")
	}

	return &AnnotationService{repo: repo}
}

func (s *AnnotationService) Save(ctx context.Context, annotation *Annotation) error {
	if annotation == nil {
		return fmt.Errorf("annotation cannot be nil")
	}

	if annotation.ID == "" {
		return fmt.Errorf("annotation ID cannot be empty")
	}

	if annotation.ImageID == "" {
		return fmt.Errorf("annotation ImageID cannot be empty")
	}

	return s.repo.Save(ctx, annotation)
}

func (s *AnnotationService) Delete(ctx context.Context, annotationID string) error {
	return s.repo.Delete(ctx, annotationID)
}

func (s *AnnotationService) GetByImage(ctx context.Context, imageID string) ([]*Annotation, error) {
	return s.repo.GetByImage(ctx, imageID)
}

func (s *AnnotationService) GetByID(ctx context.Context, annotationID string) (*Annotation, error) {
	return s.repo.GetByID(ctx, annotationID)
}
