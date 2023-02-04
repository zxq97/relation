package biz

import "github.com/pkg/errors"

var (
	ErrBlacked = errors.New("relation: blacked")
)

type BlackRepo interface {
}

type BlackUseCase struct {
	repo BlackRepo
}

func NewBlackUseCase(repo BlackRepo) *BlackUseCase {
	return &BlackUseCase{repo: repo}
}
