package role

const (
	SuperAdmin string = "SUPER_ADMIN"
	User       string = "USER"
)

// Role defines the role list
type Role struct {
	RoleId   string `json:"id"`
	RoleName string `json:"name"`
}

// GetRoleList generates role list
func GetRoleList() []string {
	var roleList []string
	roleList = append(roleList, SuperAdmin)
	roleList = append(roleList, User)

	return roleList
}

// GetRoleMap returns a boolean value to verify if the role valid or not
func GetRoleMap() map[string]bool {
	return map[string]bool{
		SuperAdmin: true,
		User:       true,
	}
}
