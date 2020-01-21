package users

import "github.com/pkg/errors"

func HandleGetUSer(repo UserRepo, userID int) (User, error) {
	user, err := repo.Get(userID)
	if err != nil {
		return User{}, errors.Wrap(err, "error getting user by ID")
	}
	return *user, nil
}

func HandleListUSer(repo UserRepo) ([]User, error) {
	allUsers, err := repo.List()
	if err != nil {
		return nil, errors.Wrap(err, "error getting list of users")
	}
	return allUsers, nil
}
