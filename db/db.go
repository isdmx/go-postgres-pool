package db

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"database/sql"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	stdlib "github.com/jackc/pgx/v5/stdlib"

	"github.com/iliadmitriev/go-postgres-pool/config"
)

func New(ctx context.Context, cfg *config.DBConfig) (*sql.DB, error) {
	// Build the connection string (DSN)
	connString, err := buildConnectionString(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to build connection string: %w", err)
	}

	// Parse the connection string into a pgxpool.Config
	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	// Configure connection pool settings
	poolConfig.MaxConns = cfg.Pool.MaxConns
	poolConfig.MinConns = cfg.Pool.MinConns
	poolConfig.MaxConnLifetime = cfg.Pool.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.Pool.MaxConnIdleTime
	poolConfig.HealthCheckPeriod = cfg.Pool.HealthCheckPeriod

	// Configure TLS if enabled
	if cfg.TLS.Enabled {
		tlsConfig, errTLS := configureTLS(cfg.TLS)
		if errTLS != nil {
			return nil, fmt.Errorf("failed to configure TLS: %w", errTLS)
		}
		poolConfig.ConnConfig.TLSConfig = tlsConfig
	}

	// Create a pgxpool.Pool
	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	// Use pgxpool.OpenDBFromPool to create a *sql.DB object
	db := stdlib.OpenDBFromPool(pool)
	return db, nil
}

// buildConnectionString constructs a PostgreSQL DSN (connection string) from the configuration.
func buildConnectionString(cfg *config.DBConfig) (string, error) {
	// Create a URL object
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(cfg.User, cfg.Password),
		Host:   cfg.Host,
		Path:   cfg.Database,
	}

	// Add query parameters
	query := u.Query()
	for key, value := range cfg.Params {
		query.Set(key, value)
	}
	u.RawQuery = query.Encode()

	return u.String(), nil
}

// configureTLS initializes and returns a *tls.Config based on the provided TLS configuration.
func configureTLS(tlsConfig config.TLSConfig) (*tls.Config, error) {
	cfg := &tls.Config{
		Certificates: []tls.Certificate{},
		RootCAs:      x509.NewCertPool(),
	}

	// Load CA certificate
	if tlsConfig.CACert != "" {
		caCert, err := os.ReadFile(tlsConfig.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		if !cfg.RootCAs.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA certificate to pool")
		}
	}

	// Load client certificate and key
	if tlsConfig.ClientCert != "" && tlsConfig.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsConfig.ClientCert, tlsConfig.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate and key: %w", err)
		}
		cfg.Certificates = append(cfg.Certificates, cert)
	}

	// Set InsecureSkipVerify
	cfg.InsecureSkipVerify = tlsConfig.SkipVerify

	return cfg, nil
}
