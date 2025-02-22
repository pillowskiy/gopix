package httprepo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/pillowskiy/gopix/internal/domain"
)

type vectorRepository struct {
	baseURL string
	client  *http.Client
}

type similarResponse struct {
	ImageID  domain.ID `json:"id"`
	Distance float64   `json:"distance"`
}

func NewVectorizationRepository(baseURL string) *vectorRepository {
	return &vectorRepository{baseURL: baseURL, client: &http.Client{}}
}

func (repo *vectorRepository) Features(
	ctx context.Context, imageID domain.ID, file *domain.FileNode,
) error {
	url := fmt.Sprintf("%s/features", repo.baseURL)

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)

	err := writer.WriteField("id", imageID.String())
	if err != nil {
		return fmt.Errorf("failed to add id field: %w", err)
	}

	part, err := writer.CreateFormFile("image", file.Name)
	if err != nil {
		return fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, file.Reader); err != nil {
		return fmt.Errorf("failed to copy reader into writer: %w", err)
	}
	writer.Close()

	req, err := http.NewRequestWithContext(ctx, "POST", url, &body)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", writer.FormDataContentType())

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		data := make(map[string]interface{})
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}
		return fmt.Errorf("server returned non-200 status: %d. message: %v", resp.StatusCode, data)
	}

	return nil
}

func (repo *vectorRepository) Similar(
	ctx context.Context, imageID domain.ID,
) ([]domain.ID, error) {
	url := fmt.Sprintf("%s/similar/%s?limit=%v", repo.baseURL, imageID.String(), 20)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := repo.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	if resp.StatusCode != http.StatusOK {
		data := make(map[string]interface{})
		if err := decoder.Decode(&data); err != nil {
			return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		return nil, fmt.Errorf(
			"server returned non-200 status: %d. %v", resp.StatusCode, data,
		)
	}

	var data []similarResponse
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	ids := make([]domain.ID, 0, len(data))
	for _, value := range data {
		ids = append(ids, value.ImageID)
	}

	return ids, nil
}

func (repo *vectorRepository) DeleteFeatures(ctx context.Context, imageID domain.ID) error {
	url := fmt.Sprintf("%s/features/%s", repo.baseURL, imageID.String())

	req, err := http.NewRequestWithContext(ctx, "DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := repo.client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data := make(map[string]interface{})
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			return fmt.Errorf("failed to unmarshal response body: %w", err)
		}

		return fmt.Errorf("server returned non-200 status: %d. message: %v", resp.StatusCode, data)
	}

	return nil
}
