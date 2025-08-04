// Package storage proporciona la funcionalidad para interactuar con la base de datos MongoDB.
package storage

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"iot-device-simulator/internal/sensor"
)

// MongoDB representa un cliente de base de datos para almacenar lecturas y configuraciones.
type MongoDB struct {
	client   *mongo.Client
	database *mongo.Database
}

// NewMongoDB crea y devuelve una nueva instancia de MongoDB conectada a la base de datos.
// Devuelve un error si la conexión falla.
func NewMongoDB(uri, dbName string) (*MongoDB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

// SaveReading guarda una única lectura de sensor en la colección 'readings'.
func (m *MongoDB) SaveReading(reading sensor.Reading) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.database.Collection("readings").InsertOne(ctx, reading)
	if err != nil {
		log.Printf("Error saving reading: %v", err)
	}
	return err
}

// SaveConfig guarda la configuración completa de un dispositivo en la colección 'configurations'.
// Utiliza una operación de upsert para crear o reemplazar la configuración existente.
func (m *MongoDB) SaveConfig(deviceID string, configs map[string]any) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc := map[string]any{
		"device_id":  deviceID,
		"configs":    configs,
		"timestamp":  time.Now(),
	}

	opts := options.Replace().SetUpsert(true)
	_, err := m.database.Collection("configurations").ReplaceOne(
		ctx,
		map[string]any{"device_id": deviceID},
		doc,
		opts,
	)
	
	if err != nil {
		log.Printf("Error saving config: %v", err)
	}
	return err
}

// GetLatestReadings recupera las últimas 'limit' lecturas para un sensorID específico,
// ordenadas por timestamp descendente.
func (m *MongoDB) GetLatestReadings(sensorID string, limit int) ([]sensor.Reading, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"sensor_id": sensorID}
	opts := options.Find().SetSort(bson.D{{Key: "timestamp", Value: -1}}).SetLimit(int64(limit))

	cursor, err := m.database.Collection("readings").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var readings []sensor.Reading
	if err := cursor.All(ctx, &readings); err != nil {
		return nil, err
	}

	return readings, nil
}

// Close cierra la conexión con la base de datos MongoDB.
func (m *MongoDB) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return m.client.Disconnect(ctx)
}