package policy

import "github.com/pillowskiy/gopix/internal/domain"

type imageAccessPolicy struct{}

func NewImageAccessPolicy() *imageAccessPolicy {
	return &imageAccessPolicy{}
}

func (p *imageAccessPolicy) CanModify(user *domain.User, image *domain.Image) bool {
	if user == nil {
		return false
	}

	isOwner := user.ID == image.AuthorID
	isAdmin := user.HasPermission(domain.PermissionsAdmin)
	return isOwner || isAdmin
}
