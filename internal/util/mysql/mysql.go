package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// MySQLClient encapsulates MySQL database operations
type MySQLClient struct {
	db      *sql.DB
	timeout time.Duration
}

var (
	MYSQL_CLIENT *MySQLClient
	once         sync.Once
)

// Initialize creates MySQL clients with the given configurations
func InitializeMySQL(host string, port int, db string, user string, password string, timeout int) error {
	var err error
	once.Do(func() {
		tmout := 10
		if timeout > 0 {
			tmout = timeout
		}

		db, dbErr := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, db))
		if dbErr != nil {
			err = dbErr
			return
		}

		// Set connection pool settings
		db.SetMaxOpenConns(25)
		db.SetMaxIdleConns(5)
		db.SetConnMaxLifetime(1 * time.Hour)

		pingCtx, cancel := context.WithTimeout(context.Background(), time.Duration(tmout)*time.Second)
		defer cancel()

		if pingErr := db.PingContext(pingCtx); pingErr != nil {
			err = pingErr
			return
		}

		client := &MySQLClient{
			db:      db,
			timeout: time.Duration(timeout) * time.Second,
		}

		MYSQL_CLIENT = client
	})

	return err
}

// Close closes the MySQL database connection
func (c *MySQLClient) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

// Query executes a query that returns rows
func (c *MySQLClient) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if c.db == nil {
		return nil, fmt.Errorf("MySQL client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.db.QueryContext(ctx, query, args...)
}

// QueryRow executes a query that returns a single row
func (c *MySQLClient) QueryRow(query string, args ...interface{}) *sql.Row {
	if c.db == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.db.QueryRowContext(ctx, query, args...)
}

// Exec executes a query that doesn't return rows
func (c *MySQLClient) Exec(query string, args ...interface{}) (sql.Result, error) {
	if c.db == nil {
		return nil, fmt.Errorf("MySQL client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.db.ExecContext(ctx, query, args...)
}

// QueryToMap executes a query and returns the results as a slice of maps
func (c *MySQLClient) QueryToMap(query string, args ...interface{}) ([]map[string]interface{}, error) {
	if c.db == nil {
		return nil, fmt.Errorf("MySQL client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	rows, err := c.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Get column names
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// Create a slice of interface{} to hold each row's values
	values := make([]interface{}, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	// Create the result slice
	var result []map[string]interface{}

	// Iterate through the rows
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		// Create a map for this row
		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]

			// Handle nil values
			if val == nil {
				row[col] = nil
				continue
			}

			// Handle different types
			switch v := val.(type) {
			case []byte:
				// Convert []byte to string for text data
				row[col] = string(v)
			default:
				row[col] = v
			}
		}

		result = append(result, row)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}

// Begin starts a new transaction
func (c *MySQLClient) Begin() (*sql.Tx, error) {
	if c.db == nil {
		return nil, fmt.Errorf("MySQL client not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	return c.db.BeginTx(ctx, nil)
}
