package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/vnurhaqiqi/vehicle_management/infras"
	"github.com/vnurhaqiqi/vehicle_management/internal/model"
)

type VechileLocationRepository interface {
	FindByFilter(ctx context.Context, filter model.VehicleLocationFilter) (vechileLocations []model.VehicleLocation, err error)
	Insert(ctx context.Context, vehicleLocation model.VehicleLocation) (err error)
}

type VechileLocationRepositoryImp struct {
	db *infras.PostgresConn
}

func ProvideVechileLocationRepository(db *infras.PostgresConn) VechileLocationRepository {
	return &VechileLocationRepositoryImp{
		db: db,
	}
}

func (r *VechileLocationRepositoryImp) FindByFilter(ctx context.Context, filter model.VehicleLocationFilter) (vechileLocations []model.VehicleLocation, err error) {
	clauses, args, err := filter.ComposeFilterClause()
	if err != nil {
		return
	}

	query := vehicleLocationQueries.Select
	if len(args) > 0 {
		query += " WHERE " + clauses
	} else {
		query += clauses
	}

	err = r.db.Conn.SelectContext(ctx, &vechileLocations, query, args...)
	if err != nil {
		return
	}

	return
}

func (r *VechileLocationRepositoryImp) Insert(ctx context.Context, vehicleLocation model.VehicleLocation) (err error) {
	err = r.db.WithTransaction(func(tx *sqlx.Tx, e chan error) {
		if _, err := tx.NamedExecContext(ctx, vehicleLocationQueries.Insert, vehicleLocation); err != nil {
			e <- err
		}
		e <- nil
	})
	if err != nil {
		return
	}

	return
}
