package forum​


type Usecase interface {
	CreateForum(forum *models.Forum) error
}
