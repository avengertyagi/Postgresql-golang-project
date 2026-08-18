package services

import (
	"context"
	"strings"

	roleconstants "github.com/akshit_tyagi/postgresql_project/internal/modules/role/constants"
	dto "github.com/akshit_tyagi/postgresql_project/internal/modules/role/dto"
	rolemodel "github.com/akshit_tyagi/postgresql_project/internal/modules/role/models"
	rolerepo "github.com/akshit_tyagi/postgresql_project/internal/modules/role/repositories"
)

type RoleService interface {
	Create(ctx context.Context, req dto.RoleRequest) (*rolemodel.Role, error)
	GetByID(ctx context.Context, id uint) (*rolemodel.Role, error)
	List(ctx context.Context, q dto.PaginationQuery) ([]rolemodel.Role, int64, error)
	Update(ctx context.Context, id uint, req dto.RoleRequest) (*rolemodel.Role, error)
	Delete(ctx context.Context, id uint) (*rolemodel.Role, error)
}

type roleService struct {
	repo rolerepo.RoleRepository
}

func NewRoleService(repo rolerepo.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) List(ctx context.Context, q dto.PaginationQuery) ([]rolemodel.Role, int64, error) {
	return s.repo.FindAll(ctx, q)
}

func (s *roleService) Create(ctx context.Context, req dto.RoleRequest) (*rolemodel.Role, error) {
	name := strings.TrimSpace(req.Name)
	exists, err := s.repo.ExistsByName(ctx, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, roleconstants.RoleAlreadyExists
	}
	role := &rolemodel.Role{
		Name: name,
	}
	if err := s.repo.Create(ctx, role); err != nil {
		return nil, err
	}
	if len(req.PermissionIDs) > 0 {
		if err := s.repo.SyncPermissions(ctx, role, req.PermissionIDs); err != nil {
			return nil, err
		}
	}
	return role, nil
}

func (s *roleService) GetByID(ctx context.Context, id uint) (*rolemodel.Role, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *roleService) Update(ctx context.Context, id uint, req dto.RoleRequest) (*rolemodel.Role, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(req.Name)
	if !strings.EqualFold(name, role.Name) {
		exists, err := s.repo.ExistsByName(ctx, name, id)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, roleconstants.RoleAlreadyExists
		}
		role.Name = name
	}
	if err := s.repo.Update(ctx, role); err != nil {
		return nil, err
	}
	if err := s.repo.SyncPermissions(ctx, role, req.PermissionIDs); err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) Delete(ctx context.Context, id uint) (*rolemodel.Role, error) {
	role, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Delete(ctx, id); err != nil {
		return nil, err
	}
	return role, nil
}
