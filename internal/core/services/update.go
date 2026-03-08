package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/ksckaan1/crtui/internal/core/models"
	"github.com/samber/lo"
)

type UpdateCheckClient interface {
	GetLatestRelease(ctx context.Context, username string, repo string) (*models.Release, error)
}

type UpdateService struct {
	client         UpdateCheckClient
	isTest         bool
	currentVersion string
}

func NewUpdateService(
	client UpdateCheckClient,
	isTest bool,
	currentVersion string,
) *UpdateService {
	return &UpdateService{
		client: client,
		isTest: isTest,
	}
}

type updateCache struct {
	Release        *models.Release `json:"release"`
	LastCheckedAt  time.Time       `json:"last_checked_at"`
	CurrentVersion string          `json:"current_version"`
}

func (s *UpdateService) CheckForUpdate(ctx context.Context) (*models.Release, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, fmt.Errorf("os.UserCacheDir: %w", err)
	}

	fileName := lo.Ternary(s.isTest, "update_cache_test.json", "update_cache.json")
	updateCachePath := filepath.Join(cacheDir, "crtui", fileName)

	var updateCacheData updateCache

	data, err := os.ReadFile(updateCachePath)
	if err == nil {
		if err := json.Unmarshal(data, &updateCacheData); err != nil {
			updateCacheData = updateCache{}
		}
	}

	cacheExpired := time.Since(updateCacheData.LastCheckedAt) > 24*time.Hour

	if !cacheExpired && updateCacheData.Release != nil && updateCacheData.CurrentVersion == s.currentVersion {
		return updateCacheData.Release, nil
	}

	latestRelease, err := s.client.GetLatestRelease(ctx, "ksckaan1", "crtui")
	if err != nil {
		return nil, fmt.Errorf("client.GetLatestRelease: %w", err)
	}

	updateCacheData = updateCache{
		Release:        latestRelease,
		LastCheckedAt:  time.Now(),
		CurrentVersion: s.currentVersion,
	}

	if newData, err := json.Marshal(&updateCacheData); err == nil {
		if mkdirErr := os.MkdirAll(filepath.Dir(updateCachePath), 0755); mkdirErr == nil {
			_ = os.WriteFile(updateCachePath, newData, 0644)
		}
	}

	return latestRelease, nil
}
