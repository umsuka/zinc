package auth

func VerifyCredentials(userID, password string) (SimpleUser, bool) {
	user, ok := ZINC_CACHED_USERS[userID]
	if !ok {
		return SimpleUser{}, false
	}

	incomingEncryptedPassword := GeneratePassword(password, user.Salt)
	if incomingEncryptedPassword == user.Password {
		return user, true
	}

	return SimpleUser{}, false
}
