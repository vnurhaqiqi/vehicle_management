package repository

var (
	vehicleLocationQueries = struct {
		Select string
		Insert string
	}{
		Select: `
		SELECT 
			id,
			vehicle_id,
			latitude,
			longitude,
			timestamp
		FROM "vehicle"."vehicle_locations" 
		`,
		Insert: `
		INSERT INTO "vehicle"."vehicle_locations" (
			vehicle_id,
			latitude,
			longitude,
			timestamp
		) VALUES (
			:vehicle_id,
			:latitude,
			:longitude,
			:timestamp
		)
		`,
	}
)
