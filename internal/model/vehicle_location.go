package model

import (
	"fmt"
	"strings"
)

type VehicleLocation struct {
	ID        int64   `db:"id"`
	VehicleID string  `db:"vehicle_id"`
	Latitude  float64 `db:"latitude"`
	Longitude float64 `db:"longitude"`
	Timestamp int64   `db:"timestamp"` // Unix timestamp in seconds
}

type VehicleLocationFilter struct {
	VechileID string
	Start     int64
	End       int64
	OrderBy   string
	Sort      string
}

func (f *VehicleLocationFilter) SetOrderBy(orderBy string) {
	f.OrderBy = orderBy
}

func (f *VehicleLocationFilter) SetSortBy(sortBy string) {
	f.Sort = sortBy
}

func (f VehicleLocationFilter) ComposeFilterClause() (string, []interface{}, error) {
	args := make([]interface{}, 0)
	clause := make([]string, 0)

	if f.VechileID != "" {
		clause = append(clause, "vehicle_id = $1")
		args = append(args, f.VechileID)
	}

	if f.Start != 0 {
		clause = append(clause, "timestamp >= $2")
		args = append(args, f.Start)
	}

	if f.End != 0 {
		clause = append(clause, "timestamp <= $3")
		args = append(args, f.End)
	}

	whereClause := strings.Join(clause, " AND ")

	if f.OrderBy != "" && f.Sort != "" {
		whereClause += fmt.Sprintf(" ORDER BY %s %s ", f.OrderBy, f.Sort)
	} else {
		whereClause += " ORDER BY timestamp ASC "
	}

	return whereClause, args, nil

}
