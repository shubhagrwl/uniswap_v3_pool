//go:generate mockgen -package=mock -destination=../../service/util/testutils/mocks/repository/pool/pool_mock.go uniswapper/internal/app/db/repository/pool IPoolLogsRepository
package pool

import (
	"context"
	"fmt"
	"uniswapper/internal/app/constants"
	"uniswapper/internal/app/db"

	pool_DBModels "uniswapper/internal/app/db/dto/pool"
)

type IPoolLogsRepository interface {
	StorePoolLogs(ctx context.Context, logs pool_DBModels.Logs) error
	GetPoolLogs(ctx context.Context, poolID, block string) (pool_DBModels.Logs, error)
	GetPoolLogsHistory(ctx context.Context, poolID string) ([]pool_DBModels.Logs, error)
}

type PoolLogsRepository struct {
	DBService *db.DBService
}

func NewPoolLogsRepository(dbService *db.DBService) IPoolLogsRepository {
	return &PoolLogsRepository{
		DBService: dbService,
	}
}

func (u *PoolLogsRepository) StorePoolLogs(ctx context.Context, logs pool_DBModels.Logs) error {
	tx := u.DBService.GetDB().Begin()
	defer tx.Rollback()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	err := tx.Table(pool_DBModels.TABLE_NAME).Create(&logs).Error
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (u *PoolLogsRepository) GetPoolLogs(ctx context.Context, poolID, block string) (pool_DBModels.Logs, error) {
	tx := u.DBService.GetDB()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var poolLogs pool_DBModels.Logs

	whr := fmt.Sprintf("%s = '%s'", pool_DBModels.COLUMN_POOL_ADDRESS, poolID)
	order := fmt.Sprintf("ABS(%s - %s)", pool_DBModels.COLUMN_BLOCK_NUMBER, block)

	switch block {
	case "latest":
		if err := tx.Table(pool_DBModels.TABLE_NAME).
			Where(whr).Order(fmt.Sprintf("%s DESC", pool_DBModels.COLUMN_CREATED_AT)).Limit(1).Scan(&poolLogs).Error; err != nil {
			return poolLogs, err
		}
	default:
		if err := tx.Table(pool_DBModels.TABLE_NAME).
			Where(whr).Order(order).Limit(1).Scan(&poolLogs).Error; err != nil {
			return poolLogs, err
		}
	}

	return poolLogs, nil
}

func (u *PoolLogsRepository) GetPoolLogsHistory(ctx context.Context, poolID string) ([]pool_DBModels.Logs, error) {
	tx := u.DBService.GetDB()
	tx.LogMode(constants.Config.DatabaseConfig.DB_LOG_MODE)

	var poolLogs []pool_DBModels.Logs
	whr := fmt.Sprintf("%s = '%s'", pool_DBModels.COLUMN_POOL_ADDRESS, poolID)

	if err := tx.Table(pool_DBModels.TABLE_NAME).
		Where(whr).Order(fmt.Sprintf("%s DESC", pool_DBModels.COLUMN_CREATED_AT)).Scan(&poolLogs).Error; err != nil {
		return poolLogs, err
	}

	return poolLogs, nil
}
