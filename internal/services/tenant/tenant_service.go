package tenant

import (
	"errors"
	"math"

	"github.com/akshit_tyagi/postgresql_project/internal/config"
	"github.com/akshit_tyagi/postgresql_project/internal/constants"
	tenantmodel "github.com/akshit_tyagi/postgresql_project/internal/models/tenant"
)

func GetAll(req tenantmodel.TenantListRequest) (*tenantmodel.TenantListResponse, error) {
	if req.CurrentPage <= 0 {
		req.CurrentPage = 1
	}
	if req.PerPage <= 0 {
		req.PerPage = 10
	}
	var total int64
	var roleList []tenantmodel.Tenant
	query := config.DB.Model(&tenantmodel.Tenant{})
	if req.Search != "" {
		query = query.Where("LOWER(name) LIKE LOWER(?)", "%"+req.Search+"%")
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}
	if req.SortBy != "name" && req.SortBy != "created_at" && req.SortBy != "updated_at" {
		req.SortBy = "id"
	}
	if req.SortOrder != "desc" && req.SortOrder != "DESC" {
		req.SortOrder = "ASC"
	}
	offset := (req.CurrentPage - 1) * req.PerPage
	if err := query.
		Order(req.SortBy + " " + req.SortOrder).
		Offset(offset).
		Limit(req.PerPage).
		Find(&roleList).Error; err != nil {
		return nil, err
	}
	lastPage := int(math.Ceil(float64(total) / float64(req.PerPage)))
	return &tenantmodel.TenantListResponse{
		Data:        roleList,
		CurrentPage: req.CurrentPage,
		PerPage:     req.PerPage,
		Total:       total,
		LastPage:    lastPage,
	}, nil
}

func Create(req tenantmodel.TenantRequest) (*tenantmodel.Tenant, error) {
	existing, err := FindByName(req.Name)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role := &tenantmodel.Tenant{
		Name:   req.Name,
		Status: true,
	}
	if err := createRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func GetByID(ID string) (*tenantmodel.Tenant, error) {
	role, err := FindByID(ID)
	if err != nil {
		if errors.Is(err, constants.RoleNotFound) {
			return nil, constants.RoleNotFound
		}
		return nil, err
	}
	return role, nil
}

func Update(id string, req tenantmodel.TenantRequest) (*tenantmodel.Tenant, error) {
	role, err := FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	existing, err := FindByName(req.Name, role.ID)
	if err == nil && existing != nil {
		return nil, constants.RoleAlreadyExists
	}
	role.Name = req.Name
	if err := updateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func Delete(id string) (*tenantmodel.Tenant, error) {
	role, err := FindByID(id)
	if err != nil {
		return nil, constants.RoleNotFound
	}
	if err := deleteRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func createRole(role *tenantmodel.Tenant) error {
	return config.DB.Create(role).Error
}

func updateRole(role *tenantmodel.Tenant) error {
	return config.DB.Save(role).Error
}

func deleteRole(role *tenantmodel.Tenant) error {
	return config.DB.Delete(role).Error
}

func FindByID(id string) (*tenantmodel.Tenant, error) {
	var role tenantmodel.Tenant
	err := config.DB.Where("id = ?", id).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func FindByName(name string, excludeID ...uint) (*tenantmodel.Tenant, error) {
	var role tenantmodel.Tenant
	query := config.DB.Where("LOWER(name) = LOWER(?)", name)
	if len(excludeID) > 0 {
		query = query.Where("id <> ?", excludeID[0])
	}
	err := query.First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
