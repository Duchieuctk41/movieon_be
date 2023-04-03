package conf

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	// dbDefault presents current active database,
	// should be initialized on starting app by calling MustOpenDefault or OpenDefault
	dbDefault *gorm.DB
)

// Config presents configuration that's necessary to work with database
type Config struct {
	Driver          string `env:"DB_DRIVER" envDefault:"mysql"`
	DSN             string `env:"DB_DSN"`
	MaxOpenConns    int    `env:"DB_MAX_OPEN_CONNS" envDefault:"25"`
	MaxIdleConns    int    `env:"DB_MAX_IDLE_CONNS" envDefault:"25"`
	ConnMaxLifetime int    `env:"DB_CONN_MAX_LIFETIME" envDefault:"600"`

	Host   string `env:"DB_HOST" envDefault:"localhost"`
	Port   string `env:"DB_PORT" envDefault:"3306"`
	User   string `env:"DB_USER" envDefault:"doan"`
	Pass   string `env:"DB_PASS" envDefault:"doan"`
	Name   string `env:"DB_NAME" envDefault:"doan"`
	Schema string `env:"DB_SCHEMA" envDefault:"public"`
}

// GetDSN returns a dsn that is read from ENV or built from separated env DB_*
func (c Config) GetDSN() string {
	if c.DSN != "" {
		return c.DSN
	}
	c.DSN = fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable connect_timeout=5",
		c.Host,
		c.Port,
		c.User,
		c.Name,
		c.Pass,
	)

	return c.DSN
}

// Open open a DB connection
//
//	dbDefault, err := Open(config)
func Open(config *Config) (*gorm.DB, error) {
	naming := &schema.NamingStrategy{
		SingularTable: true,
	}
	cfg := &gorm.Config{
		NamingStrategy: naming,
		Logger:         logger.Default.LogMode(logger.Silent),
	}

	var dialector gorm.Dialector
	switch config.Driver {
	case "sqlite", "sqlite3":
		dialector = sqlite.Open(config.GetDSN())
	case "postgres":
		dialector = postgres.Open(config.GetDSN())
		naming.TablePrefix = config.Schema + "."
	case "mysql":
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", config.User, config.Pass, config.Host, config.Port, config.Name)
		dialector = mysql.Open(dsn)
	default:
		return nil, fmt.Errorf("unsupported driver %s", config.Driver)
	}

	db, err := gorm.Open(dialector, cfg)
	if err != nil {
		return nil, err
	}

	theDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	if config.MaxIdleConns > 0 {
		theDB.SetMaxIdleConns(config.MaxIdleConns)
	}
	if config.MaxOpenConns > 0 {
		theDB.SetMaxOpenConns(config.MaxOpenConns)
	}
	if config.ConnMaxLifetime > 0 {
		theDB.SetConnMaxLifetime(time.Duration(config.ConnMaxLifetime) * time.Second)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	if err = theDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error while ping DB: %v", err)
	}

	return db, nil
}

// Close release a DB instance
func Close(db *gorm.DB) {
	if db != nil {
		if dbInstance, err := db.DB(); err == nil {
			_ = dbInstance.Close()
		}
	}
}

// inMemorySqliteCfg presents configuration for quick testing
// this is lightweight database, should consider to user a real DB
// in more advanced testing like concurrency writing
var inMemorySqliteCfg = &Config{
	Driver:          "sqlite3",
	DSN:             ":memory:",
	MaxOpenConns:    1, // should be 1 cuz sqlite doesn't support concurrency writing operation.
	MaxIdleConns:    1,
	ConnMaxLifetime: 600,
}

// MustSetupTest setups an in-memory DB for testing and set to default
// it'll panic if errors occur
func MustSetupTest() {
	db, err := Open(inMemorySqliteCfg)
	if err != nil {
		panic(err)
	}

	dbDefault = db
}

// GetDB gets default database connection
func GetDB() *gorm.DB {
	if dbDefault == nil {
		panic(errors.New("uninitialized database. Please connect first"))
	}
	return dbDefault
}

// OpenDefault opens default database connection and assign to default
func OpenDefault(config *Config) error {
	db, err := Open(config)
	if err != nil {
		return err
	}
	dbDefault = db

	return nil
}

// MustOpenDefault open connection & assign to dbDefault, this will panic application if failed
func MustOpenDefault(config *Config) {
	if err := OpenDefault(config); err != nil {
		panic(err)
	}
}

// CloseDB closes default database
func CloseDB() {
	Close(dbDefault)
	dbDefault = nil
}
