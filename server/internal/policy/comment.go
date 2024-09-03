package policy

import "github.com/pillowskiy/gopix/internal/domain"

type commentAccessPolicy struct{}

func NewCommentAccessPolicy() *commentAccessPolicy {
	return &commentAccessPolicy{}
}

func (p *commentAccessPolicy) CanModify(user *domain.User, comment *domain.Comment) bool {
	if user == nil {
		return false
	}

	isOwner := user.ID == comment.AuthorID
	isAdmin := user.HasPermission(domain.PermissionsAdmin)
	return isOwner || isAdmin
}
