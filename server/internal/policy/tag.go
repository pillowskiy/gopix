package policy

import "github.com/pillowskiy/gopix/internal/domain"

type tagAccessPolicy struct{}

func NewTagAccessPolicy() *tagAccessPolicy {
	return &tagAccessPolicy{}
}

func (p *tagAccessPolicy) CanModifyImageTags(user *domain.User, image *domain.Image) bool {
	if user == nil || image == nil {
		return false
	}

	isAuthor := user.ID == image.AuthorID
	isAdmin := user.HasPermission(domain.PermissionsAdmin)
	return isAuthor || isAdmin
}
