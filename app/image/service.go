package image

type Repository interface {
	GetImage()
}

type Service interface {
	GetImage()
}

type service struct {
	r *Repository
}

func (s *service) GetImage(){
	
}
