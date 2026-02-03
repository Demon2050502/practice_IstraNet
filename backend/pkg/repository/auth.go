package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"

	dbmodel "practice_IstraNet/pkg/DB_model"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

var (
    ErrUserExists   = errors.New("user already exists")
    ErrUserNotFound = errors.New("user not found")
    ErrUserInactive = errors.New("user inactive")
    ErrRoleNotFound = errors.New("role not found")
    ErrNoUserRole   = errors.New("user role not found")
)

func (r *AuthPostgres) CreateUser(ctx context.Context, email, passwordHash, fullName, roleCode string) (userID int64, roleOut string, err error) {
    if roleCode == "" {
        roleCode = "user"
    }

    tx, err := r.db.BeginTxx(ctx, nil)
    if err != nil {
        return 0, "", err
    }
    defer func() {
        if err != nil {
            _ = tx.Rollback()
        }
    }()

    qUser := `INSERT INTO users (email, password_hash, full_name) VALUES ($1,$2,$3) RETURNING id`
    if err = tx.QueryRowContext(ctx, qUser, email, passwordHash, fullName).Scan(&userID); err != nil {
        var pqErr *pq.Error
        if errors.As(err, &pqErr) && pqErr.Code == "23505" {
            return 0, "", ErrUserExists
        }
        return 0, "", err
    }

    var roleID int64
    qRole := `SELECT id FROM roles WHERE code=$1`
    if err = tx.QueryRowContext(ctx, qRole, roleCode).Scan(&roleID); err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return 0, "", ErrRoleNotFound
        }
        return 0, "", err
    }

    qUR := `INSERT INTO user_roles (user_id, role_id) VALUES ($1,$2)`
    if _, err = tx.ExecContext(ctx, qUR, userID, roleID); err != nil {
        return 0, "", err
    }

    if err = tx.Commit(); err != nil {
        return 0, "", err
    }

    return userID, roleCode, nil
}

func (r *AuthPostgres) GetUserByEmail(ctx context.Context, email string) (dbmodel.UserDB, error) {
    var u dbmodel.UserDB
    q := `SELECT id, email, password_hash, full_name, is_active FROM users WHERE email=$1`
    err := r.db.GetContext(ctx, &u, q, email)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return dbmodel.UserDB{}, ErrUserNotFound
        }
        return dbmodel.UserDB{}, err
    }
    if !u.IsActive {
        return dbmodel.UserDB{}, ErrUserInactive
    }
    return u, nil
}

func (r *AuthPostgres) GetUserRoleCode(ctx context.Context, userID int64) (string, error) {
    var code string
    q := `
        SELECT r.code
        FROM roles r
        JOIN user_roles ur ON ur.role_id=r.id
        WHERE ur.user_id=$1
        ORDER BY r.id
        LIMIT 1
    `
    err := r.db.GetContext(ctx, &code, q, userID)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return "", ErrNoUserRole
        }
        return "", err
    }
    return code, nil
}
