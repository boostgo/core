package mongox

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/boostgo/core/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const migrationsCollection = "schema_migrations"

type Migration struct {
	Version   int       `bson:"version"`
	Name      string    `bson:"name"`
	AppliedAt time.Time `bson:"applied_at"`
	Checksum  string    `bson:"checksum"`
}

type MigrationFile struct {
	Version  int
	Name     string
	Content  string
	Checksum string
}

type Migrator struct {
	client     Client
	collection *mongo.Collection
}

func NewMigrator(client Client) *Migrator {
	return &Migrator{
		client:     client,
		collection: client.Collection(migrationsCollection),
	}
}

func (m *Migrator) Init(ctx context.Context) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.D{{Key: "version", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	_, err := m.collection.Indexes().CreateOne(ctx, indexModel)
	return err
}

func (m *Migrator) LastAppliedVersion(ctx context.Context) (int, error) {
	opts := options.FindOne().SetSort(bson.D{{Key: "version", Value: -1}})

	var migration Migration
	err := m.collection.FindOne(ctx, bson.M{}, opts).Decode(&migration)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, nil
		}

		return 0, err
	}

	return migration.Version, nil
}

func (m *Migrator) RunMigrations(ctx context.Context, migrationsPath string) error {
	if err := m.Init(ctx); err != nil {
		return fmt.Errorf("failed to init migrations collection: %w", err)
	}

	// Получаем последнюю примененную версию
	lastVersion, err := m.LastAppliedVersion(ctx)
	if err != nil {
		return fmt.Errorf("failed to get last applied version: %w", err)
	}

	// Читаем файлы миграций
	migrations, err := m.loadMigrations(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Применяем новые миграции
	for _, migration := range migrations {
		if migration.Version <= lastVersion {
			continue // Пропускаем уже примененные
		}

		log.Info().
			Int("version", migration.Version).
			Str("name", migration.Name).
			Msg("Applying migration")

		if err := m.applyMigration(ctx, migration); err != nil {
			return fmt.Errorf("failed to apply migration %d: %w", migration.Version, err)
		}

		// Записываем информацию о примененной миграции
		migrationRecord := Migration{
			Version:   migration.Version,
			Name:      migration.Name,
			AppliedAt: time.Now(),
			Checksum:  migration.Checksum,
		}

		if _, err = m.collection.InsertOne(ctx, migrationRecord); err != nil {
			return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
		}

		log.Info().
			Int("version", migration.Version).
			Str("name", migration.Name).
			Msg("Migration applied successfully")
	}

	return nil
}

func (m *Migrator) loadMigrations(migrationsPath string) ([]MigrationFile, error) {
	var migrations []MigrationFile

	err := filepath.WalkDir(migrationsPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || !strings.HasSuffix(path, ".js") {
			return nil
		}

		filename := filepath.Base(path)
		version, name, err := parseMigrationFilename(filename)
		if err != nil {
			return fmt.Errorf("invalid migration filename %s: %w", filename, err)
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", path, err)
		}

		checksum := fmt.Sprintf("%x", sha256.Sum256(content))

		migrations = append(migrations, MigrationFile{
			Version:  version,
			Name:     name,
			Content:  string(content),
			Checksum: checksum,
		})

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Сортируем по версии
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func parseMigrationFilename(filename string) (int, string, error) {
	// Ожидаем формат: 001_migration_name.js
	parts := strings.SplitN(filename, "_", 2)
	if len(parts) != 2 {
		return 0, "", fmt.Errorf("filename should be in format: NNN_name.js")
	}

	version, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", fmt.Errorf("invalid version number: %w", err)
	}

	name := strings.TrimSuffix(parts[1], ".js")
	return version, name, nil
}

func (m *Migrator) applyMigration(ctx context.Context, migration MigrationFile) error {
	// Выполняем JavaScript код миграции
	// Это упрощенная версия - в реальности нужно использовать JavaScript engine
	// или выполнять команды через mongo shell

	return m.executeJavaScript(ctx, migration.Content)
}

func (m *Migrator) executeJavaScript(ctx context.Context, script string) error {
	// В реальной реализации здесь должен быть JavaScript engine
	// Для простоты показываем как можно выполнять базовые операции

	// Парсим и выполняем команды
	// Это очень упрощенная реализация - в продакшене лучше использовать
	// готовые решения или выполнять через mongo shell

	return m.client.Database().RunCommand(ctx, bson.D{
		{Key: "eval", Value: script},
	}).Err()
}
