package types

// Organization API (Group uses valid identifiers; User has hyphens).

// Group is the group payload for Client.CreateGroup and
// Client.UpdateGroup.
type Group struct {
	// Name is the unique group name. Required for create.
	Name string `json:"name,omitempty"`
	// Comments is a free-text description of the group.
	Comments string `json:"comments,omitempty"`
	// User lists the member usernames.
	User []string `json:"user,omitempty"`
}

// User is the user payload for Client.CreateUser and
// Client.UpdateUser.
type User struct {
	// UserID is the unique username. Required for create.
	UserID string `json:"user-id,omitempty"`
	// FullName is the display name.
	FullName string `json:"full-name,omitempty"`
	// UserComments is a free-text comment.
	UserComments string `json:"user-comments,omitempty"`
	// Email is the email address.
	Email string `json:"email,omitempty"`
	// Password is the plain-text password (pass hashPassword=true
	// to Client.CreateUser to have OpenNMS hash it on receipt).
	Password string `json:"password,omitempty"`
	// PasswordSalt is true when Password is already salted.
	PasswordSalt *bool `json:"passwordSalt,omitempty"`
	// DutySchedule holds duty schedule strings,
	// e.g. ["MoTuWeThFrSaSu800-2300"].
	DutySchedule []string `json:"duty-schedule,omitempty"`
	// Role holds security roles, e.g. ["ROLE_ADMIN"].
	Role []string `json:"role,omitempty"`
}

// Category is the category payload for Client.CreateCategory and
// Client.UpdateCategory.
type Category struct {
	// Name is the unique category name. Required for create.
	Name string `json:"name,omitempty"`
	// AuthorizedGroups lists the groups authorized to see this
	// category.
	AuthorizedGroups []string `json:"authorizedGroups,omitempty"`
}
