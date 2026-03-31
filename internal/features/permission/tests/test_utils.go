package permission

import "github.com/google/uuid"

func TestUserID() uuid.UUID {
	return uuid.MustParse("aaaa1111-1111-1111-1111-111111111111")
}

func TestCodeNames() []string {
	return []string{"view_experiment", "manage_experiment"}
}

func TestDuplicateCodeNames() []string {
	return []string{"READ", "READ"}
}
