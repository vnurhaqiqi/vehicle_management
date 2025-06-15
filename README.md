# Vehicle Management App

## Depedencies:
- PostgreSQL
- RabbitMQ
- Eclipse Mosquitto (MQTT broker)

## Requirements

- Docker
- Docker Compose

## Running the App

1. **Clone the repository**

```bash
git clone hhttps://github.com/vnurhaqiqi/vehicle_management.git
cd vehicle-management
```

2. **Create `.env`, the example is on `.env.example`**

3. **Migrate database schema on `init.sql` file**

4. **Build and Run**

```bash
docker-compose up --build
```

```
This command will:

Build and run the Go application

Start PostgreSQL on port 5432

Start RabbitMQ (with UI on http://localhost:15672, user/pass: guest/guest)

Start MQTT broker on port 1883
```

5. Run the publish vehicle location
```bash
cd mqttpublishermock/
go run main.go
```


