package service

import (
	"authn_authz/db"

	"github.com/google/uuid"
)

type UserService struct {
	Data []*db.File
}

// Retrieve a file
func (u *UserService) RetrieveFiles(userId uuid.UUID) ([]*db.File, error) {
	// retrieve files for the logged in user
	var files []*db.File
	for _, f := range u.Data {
		if f.UserId == userId {
			files = append(files, f)
		}
	}
	return files, nil
}

func (u *UserService) UploadFile(userId uuid.UUID, fileName string) (*db.File, error) {
	// create a new file
	newFile := &db.File{
		ID:     uuid.New(),
		Name:   fileName,
		UserId: userId,
	}
	u.Data = append(u.Data, newFile)
	return newFile, nil
}
