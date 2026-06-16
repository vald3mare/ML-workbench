package handlers

import (
	"context"
	"net/http"

	"github.com/vald3mare/ML-workbench/internal/annotation"
	"github.com/vald3mare/ML-workbench/internal/postgres"
)

// это транспортный слой приложения, который отвечает за обработку HTTP запросов и передачу данных между клиентом и сервисом,
// он будет использовать сервис для работы с бизнес логикой, а не с данными напрямую
// определен там где нужен потребителю, так как он может использоваться, а не в месте реализации
type AnnotationService interface {
	Save(ctx context.Context, annotation *annotation.Annotation) error
	Delete(ctx context.Context, annotationID string) error
	GetByID(ctx context.Context, annotationID string) (*annotation.Annotation, error)
}

type AnnotationHandler struct {
	service AnnotationService
}

func (h *AnnotationHandler) Save(w http.ResponseWriter, r *http.Request) {
	// здесь будет логика для сохранения аннотации, например чтение данных из запроса, валидация и вызов сервиса для сохранения данных
	postgresService, ok := h.service.(*postgres.PostgresService)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// пример использования postgresService для сохранения аннотации, в реальной реализации здесь будет логика для чтения данных из запроса и валидации
	err := postgresService.Save(context.Background(), &annotation.Annotation{
		ID:      "test-id",
		ImageID: "test-image-id",
	})
	if err != nil {
		http.Error(w, "Failed to save annotation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *AnnotationHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// здесь будет логика для удаления аннотации, например чтение данных из запроса, валидация и вызов сервиса для удаления данных
	postgresService, ok := h.service.(*postgres.PostgresService)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// пример использования postgresService для удаления аннотации, в реальной реализации здесь будет логика для чтения данных из запроса и валидации
	err := postgresService.Delete(context.Background(), "test-id")
	if err != nil {
		http.Error(w, "Failed to delete annotation", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *AnnotationHandler) GetByImage(w http.ResponseWriter, r *http.Request) {
	// здесь будет логика для получения аннотаций по id изображения, например чтение данных из запроса, валидация и вызов сервиса для получения данных
	postgresService, ok := h.service.(*postgres.PostgresService)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// пример использования postgresService для получения аннотаций по id изображения, в реальной реализации здесь будет логика для чтения данных из запроса и валидации
	_, err := postgresService.GetByImage(context.Background(), "test-image-id")
	if err != nil {
		http.Error(w, "Failed to get annotations", http.StatusInternalServerError)
		return
	}
	// здесь будет логика для сериализации данных и отправки их в ответе, например с помощью json.Marshal
	w.WriteHeader(http.StatusOK)
	// w.Write(jsonData)
}

func (h *AnnotationHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// здесь будет логика для получения аннотации по id, например чтение данных из запроса, валидация и вызов сервиса для получения данных
	postgresService, ok := h.service.(*postgres.PostgresService)
	if !ok {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// пример использования postgresService для получения аннотации по id, в реальной реализации здесь будет логика для чтения данных из запроса и валидации
	_, err := postgresService.GetByID(context.Background(), "test-id")
	if err != nil {
		http.Error(w, "Failed to get annotation", http.StatusInternalServerError)
		return
	}
	// здесь будет логика для сериализации данных и отправки их в ответе, например с помощью json.Marshal
	w.WriteHeader(http.StatusOK)
	// w.Write(jsonData)
}
