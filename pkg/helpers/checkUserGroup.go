package helpers

func UserExists(user string) bool {
	return RunSilent("id", user) == nil
}

func GroupExists(group string) bool {
	return RunSilent("getent", "group", group) == nil
}
