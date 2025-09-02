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

// ✅ Create a new record
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}
	return r.db.WithContext(ctx).Create(entity).Error
}

// ✅ Get one record by condition
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

// ✅ Get by ID
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

// ✅ Update
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if entity == nil {
		return errors.New("entity cannot be nil")
	}
	return r.db.WithContext(ctx).Save(entity).Error
}

// Soft Delete by condition
func (r *Repository[T]) Delete(ctx context.Context, query string, args ...interface{}) error {
	var entity T

	updates := map[string]interface{}{}
	t := reflect.TypeOf(entity)
	if _, ok := t.FieldByName("IsActive"); ok {
		updates["is_active"] = false
	}

	// First, apply IsActive=false with the proper model
	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&entity).Where(query, args...).Updates(updates).Error; err != nil {
			return err
		}
	}

	// Then, soft delete using GORM
	return r.db.WithContext(ctx).Where(query, args...).Delete(&entity).Error
}

// Soft Delete by ID
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

// ✅ List with filters + pagination
func (r *Repository[T]) List(ctx context.Context, limit, offset int, filters map[string]interface{}) ([]T, error) {
	var entities []T
	query := r.db.WithContext(ctx).Model(new(T))

	for k, v := range filters {
		query = query.Where(k, v)
	}

	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}

	if err := query.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// ✅ Count with filters
func (r *Repository[T]) Count(ctx context.Context, filters map[string]interface{}) (int64, error) {
	var count int64
	query := r.db.WithContext(ctx).Model(new(T))

	for k, v := range filters {
		query = query.Where(k, v)
	}

	if err := query.Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// ✅ Transaction support
func (r *Repository[T]) WithTransaction(tx *gorm.DB) *Repository[T] {
	return &Repository[T]{db: tx}
}

// ✅ Extra helpers for advanced queries ----------------------

// Find with raw GORM query builder
func (r *Repository[T]) FindWithQuery(ctx context.Context, queryFn func(*gorm.DB) *gorm.DB) ([]T, error) {
	var entities []T
	db := queryFn(r.db.WithContext(ctx))
	if err := db.Find(&entities).Error; err != nil {
		return nil, err
	}
	return entities, nil
}

// Aggregate query (useful for SUM, AVG, etc.)
func (r *Repository[T]) Aggregate(ctx context.Context, queryFn func(*gorm.DB) *gorm.DB, dest interface{}) error {
	return queryFn(r.db.WithContext(ctx).Model(new(T))).Scan(dest).Error
}

// UnitOfWork wraps a GORM transaction
type UnitOfWork struct {
	tx        *gorm.DB
	committed bool
}

// NewUnitOfWork starts a new transaction.
// If `readonly` is true, we can skip committing (optional for read-only UOWs)
func NewUnitOfWork(db *gorm.DB) *UnitOfWork {
	return &UnitOfWork{
		tx: db.Begin(),
	}
}

// Commit commits the transaction
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

// Rollback rolls back the transaction
func (u *UnitOfWork) Rollback() error {
	if u.committed {
		return errors.New("cannot rollback a committed transaction")
	}
	return u.tx.Rollback().Error
}

// Tx returns the *gorm.DB instance for use in service methods
func (u *UnitOfWork) Tx() *gorm.DB {
	return u.tx
}
