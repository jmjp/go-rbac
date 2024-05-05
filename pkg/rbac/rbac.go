package rbac

import (
	"encoding/json"
	"os"
	"strings"
)

type RoleBasedAccessControl struct {
	Name        string
	Permissions []string
}

type Roles struct {
	Roles []RoleBasedAccessControl
}

type RBACBuilder struct {
	Roles Roles
}

func NewRBACBuilder() *RBACBuilder {
	return &RBACBuilder{Roles: Roles{}}
}

func (r *RBACBuilder) WithRolesFromFile(file string) *RBACBuilder {
	// parse json file
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		panic(err)
	}
	// unmarshal json
	var roles Roles
	if err := json.Unmarshal(fileBytes, &roles); err != nil {
		panic(err)
	}
	r.Roles = roles
	return r
}

func (r *RBACBuilder) WithRoles(roles []RoleBasedAccessControl) *RBACBuilder {
	r.Roles.Roles = roles
	return r
}

func (r *RBACBuilder) Build() Roles {
	return r.Roles
}

func (r *Roles) GetRolePermissions(role string) []string {
	for _, roleBasedAccessControl := range r.Roles {
		if roleBasedAccessControl.Name == role {
			return roleBasedAccessControl.Permissions
		}
	}
	return []string{}
}

func (r *Roles) HasPermission(permission, role string) bool {
	permissionChunks := strings.Split(permission, "::")
	if len(permissionChunks) != 3 {
		return false
	}
	permissions := r.GetRolePermissions(role)
	for _, p := range permissions {
		selfChunks := strings.Split(p, "::")
		if selfChunks[0] == permissionChunks[0] && (selfChunks[1] == permissionChunks[1] || selfChunks[1] == "*") && (selfChunks[2] == permissionChunks[2] || selfChunks[2] == "*") {
			return true
		}
	}
	return false
}
