package rpg

func MapIdToResponse(id int) RpgIdResponse {
	return RpgIdResponse{
		Id: id,
	}
}

func MapSlugToResponse(slug string) RpgSlugResponse {
	return RpgSlugResponse{
		Slug: slug,
	}
}
