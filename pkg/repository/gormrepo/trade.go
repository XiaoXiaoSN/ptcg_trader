package gormrepo

import (
	"context"
	"ptcg_trader/pkg/model"

	"gorm.io/gorm"
)

func buildItemQuery(db *gorm.DB, query model.ItemQuery) *gorm.DB {
	if query.ID != nil {
		db = db.Where("id = ?", query.ID)
	}

	if query.PerPage <= 0 {
		query.PerPage = DefaultPerPage
	}
	if query.Page <= 0 {
		query.Page = 1
	}
	db = db.Limit(query.PerPage).
		Offset((query.Page - 1) * query.PerPage)

	return db
}

// GetItem get information of the targe item
func (repo *_repository) GetItem(ctx context.Context, query model.ItemQuery) (model.Item, error) {
	var item model.Item

	db := buildItemQuery(repo.DB(ctx), query).
		Model(&item)

	err := db.First(&item).Error
	if err != nil {
		return item, notFoundOrInternalError(err)
	}
	return item, nil
}

// CountItems counting total count of items
func (repo *_repository) CountItems(ctx context.Context, query model.ItemQuery) (int64, error) {
	var total int64

	db := buildItemQuery(repo.DB(ctx), query).
		Model(&model.Item{})

	err := db.Count(&total).Error
	if err != nil {
		return total, err
	}
	return total, nil
}

// ListItems list Items
func (repo *_repository) ListItems(ctx context.Context, query model.ItemQuery) ([]model.Item, error) {
	var itemList []model.Item

	db := buildItemQuery(repo.DB(ctx), query).
		Model(&model.Item{})

	err := db.Find(&itemList).Error
	if err != nil {
		return nil, err
	}
	return itemList, nil
}
