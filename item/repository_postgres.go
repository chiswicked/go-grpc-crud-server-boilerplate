package item

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/chiswicked/go-grpc-crud-server-boilerplate/model"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// PostgresRepository struct
type PostgresRepository struct {
	db *sql.DB
}

// NewPostgresRepository func
func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

// Create func
func (repo *PostgresRepository) Create(ctx context.Context, item *model.Item) (*string, error) {
	if len(item.Name) <= 0 {
		return nil, fmt.Errorf("Invalid Argument")
	}

	qry := `
		INSERT INTO itemtable (uuid, name)
		VALUES ($1, $2);
	`
	uid, err := uuid.NewV4()
	if err != nil {
		return nil, fmt.Errorf("Could not insert item into the database: %s", err)
	}

	out := uid.String()
	_, err = repo.db.ExecContext(ctx, qry, uid, item.Name)
	if err != nil {
		return nil, fmt.Errorf("Could not insert item into the database: %s", err)
	}

	return &out, nil
}

// GetByID func
func (repo *PostgresRepository) GetByID(ctx context.Context, id string) (*model.Item, error) {
	if _, err := uuid.FromString(id); err != nil {
		return nil, status.Errorf(codes.NotFound, "Not Found")
	}

	qry := `
		SELECT uuid, name
		FROM itemtable
		WHERE uuid = $1;
	`
	out := &model.Item{}
	err := repo.db.QueryRowContext(ctx, qry, id).Scan(&out.ID, &out.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "Not Found")
		}
		return nil, grpc.Errorf(codes.Internal, "Could not read item from the database: %s", err)
	}

	return out, nil
}

// Fetch func
func (repo *PostgresRepository) Fetch(ctx context.Context, num int64) ([]*model.Item, error) {
	qry := `
		SELECT uuid, name
		FROM itemtable
	`
	var outItems = []*model.Item{}

	// jozsi := sql.Rows
	rows, err := repo.db.QueryContext(ctx, qry)
	if err != nil {
		return nil, fmt.Errorf("Could not read items from the database: %s", err)
	}

	for rows.Next() {
		var id, name string
		err := rows.Scan(&id, &name)
		outItems = append(outItems, &model.Item{ID: id, Name: name})
		if err != nil {
			if err == sql.ErrNoRows {
				break
			}
			return nil, fmt.Errorf("Could not read items from the database: %s", err)
		}
	}

	return outItems, nil
}

// Update func
func (repo *PostgresRepository) Update(ctx context.Context, item *model.Item) (*model.Item, error) {
	qry := `
		UPDATE itemtable
		SET name = $2
		WHERE uuid = $1
	`

	res, err := repo.db.ExecContext(ctx, qry, item.ID, item.Name)
	if err != nil {
		return nil, fmt.Errorf("Could not update item in the database: %s", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("Error while updating item in the database: %s", err)
	}

	// TODO: check rows affected

	return item, nil
}

// Delete func
func (repo *PostgresRepository) Delete(ctx context.Context, id string) (bool, error) {
	qry := `
		DELETE FROM itemtable
		WHERE uuid = $1
	`

	res, err := repo.db.ExecContext(ctx, qry, id)
	if err != nil {
		return false, fmt.Errorf("Could not delete item from the database: %s", err)
	}

	_, err = res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("Error while deleting item from the database: %s", err)
	}

	// TODO: check rows affected

	return true, nil
}
