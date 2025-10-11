package repositories

import (
	"context"
	"go-EdTech/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type RoleRepository struct {
	db *pgxpool.Pool
}

func NewRoleRepository(conn *pgxpool.Pool) *RoleRepository {
	return &RoleRepository{db: conn}
}

func (r *RoleRepository) GetRoleByID(c context.Context, roleID int) (*models.Role, error) {
	var role models.Role
	row := r.db.QueryRow(c, `select id, name from roles where id = $1`, roleID)

	if err := row.Scan(&role.Id, &role.Name); err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) GetRoleByName(c context.Context, name string) (*models.Role, error) {
	var role models.Role

	row := r.db.QueryRow(c, `select id, name from roles where name = $1`, name)

	if err := row.Scan(&role.Id, &role.Name); err != nil {
		return nil, err
	}

	return &role, nil
}

func (r *RoleRepository) Create(ctx context.Context, role *models.Role) error {
	_, err := r.db.Exec(ctx, `insert into roles (name) values ($1)`, role.Name)
	return err
}