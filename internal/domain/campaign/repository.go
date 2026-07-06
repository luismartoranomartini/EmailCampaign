package campaign

type Repository interface {
	Create(campaign *Campaign) error
	Get() ([]Campaign, error)
	GetBy(id string) (*Campaign, error)
	Update(campaign *Campaign) error
	Delete(campaign *Campaign) error
}
