package utils

var RoleAccepted = []string{
	"11111111-1111-1111-1111-111111111111",
	"22222222-2222-2222-2222-222222222222",
	"33333333-3333-3333-3333-333333333333",
}

func CheckRoleAccepted(role string) bool {

	for _, v := range RoleAccepted {
		if v == role {
			return true
		}
	}
	return false
}
