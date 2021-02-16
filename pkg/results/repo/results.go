package repo

import "github.com/IoanStoianov/Open-func/pkg/types"

// Results repository
type Results interface {
	AddRecord(result *types.FuncResult) error
	GetRecords(name string, count int64) ([]*types.FuncResult, error)
}
