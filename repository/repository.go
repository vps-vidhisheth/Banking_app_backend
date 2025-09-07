package repository

import (
	"context"
	"errors"
	"reflect"

	"gorm.io/gorm"
)

type Repository[T any] struct {
	db *gorm.DB
}

func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}
	return r.db.WithContext(ctx).Create(entity).Error
}

func (r *Repository[T]) GetOne(ctx context.Context, query string, args ...interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).Where(query, args...).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) GetByID(ctx context.Context, id interface{}) (*T, error) {
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}
	return r.db.WithContext(ctx).Save(entity).Error
}

func (r *Repository[T]) Delete(ctx context.Context, query string, args ...interface{}) error {
	var entity T

	updates := map[string]interface{}{}
	t := reflect.TypeOf(entity)
	if _, ok := t.FieldByName("IsActive"); ok {
		updates["is_active"] = false
	}

	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&entity).Where(query, args...).Updates(updates).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).Where(query, args...).Delete(&entity).Error
}

func (r *Repository[T]) DeleteByID(ctx context.Context, id interface{}) error {
	var entity T

	updates := map[string]interface{}{}
	t := reflect.TypeOf(entity)
	if _, ok := t.FieldByName("IsActive"); ok {
		updates["is_active"] = false
	}

	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&entity).Where("id = ?", id).Updates(updates).Error; err != nil {
			return err
		}
	}

	return r.db.WithContext(ctx).Where("id = ?", id).Delete(&entity).Error
}

func (r *Repository[T]) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]T, error) {
	var entities []T

	db := r.db.Session(&gorm.Session{}).WithContext(ctx)

	query := db.Model(new(T))

	for k, v := range filters {
		query = query.Where(k, v)
	}

	var entity T
	t := reflect.TypeOf(entity)
	orderColumn := "id"

	if _, ok := t.FieldByName("CustomerID"); ok {
		orderColumn = "customer_id"
	} else if _, ok := t.FieldByName("BankID"); ok {
		orderColumn = "bank_id"
	} else if _, ok := t.FieldByName("AccountID"); ok {
		orderColumn = "account_id"
	}

	query = query.Order(orderColumn + " ASC")

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}

	return entities, nil
}

func (r *Repository[T]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	db := r.db.Session(&gorm.Session{}).WithContext(ctx)

	query := db.Model(new(T))

	for k, v := range filters {
		query = query.Where(k, v)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *Repository[T]) WithTransaction(tx *gorm.DB) *Repository[T] {
	return &Repository[T]{db: tx}
}

func (r *Repository[T]) FindWithQuery(ctx context.Context, queryFn func(*gorm.DB) *gorm.DB) ([]T, error) {
	var entities []T
	db := queryFn(r.db.WithContext(ctx))
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

func (r *Repository[T]) Aggregate(ctx context.Context, queryFn func(*gorm.DB) *gorm.DB, dest interface{}) error {
	return queryFn(r.db.WithContext(ctx).Model(new(T))).Scan(dest).Error
}

type UnitOfWork struct {
	tx        *gorm.DB
	committed bool
}

func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		tx: db.Begin(),
	}
}

func (u *UnitOfWork) Commit() error {
	if u.committed {
		return errors.New("transaction already committed")
	}
	if err := u.tx.Commit().Error; err != nil {
		return err
	}
	u.committed = true
	return nil
}

func (u *UnitOfWork) Rollback() error {
	if u.committed {
		return errors.New("cannot rollback a committed transaction")
	}
	return u.tx.Rollback().Error
}

func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}

func (r *Repository[T]) GetDB() *gorm.DB {
	return r.db
}
