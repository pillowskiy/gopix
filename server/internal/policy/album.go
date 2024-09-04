package policy

import "github.com/pillowskiy/gopix/internal/domain"

type albumAccessPolicy struct{}

func NewAlbumAccessPolicy() *albumAccessPolicy {
	return &albumAccessPolicy{}
}

func (p *albumAccessPolicy) CanModify(user *domain.User, album *domain.Album) bool {
	if user == nil {
		return false
	}

	isOwner := user.ID == album.AuthorID
	isAdmin := user.HasPermission(domain.PermissionsAdmin)
	return isOwner || isAdmin
}
