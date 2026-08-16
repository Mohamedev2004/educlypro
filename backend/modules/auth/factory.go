package auth

func NewCenter(name, slug string) Center {
	return Center{
		Name: name,
		Slug: slug,
	}
}
