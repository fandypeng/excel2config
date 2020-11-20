package service

import (
	pb "excel2config/api"
	"excel2config/internal/dao"
	"github.com/go-kratos/kratos/pkg/conf/paladin"

	"github.com/google/wire"
)

var Provider = wire.NewSet(New, wire.Bind(new(pb.SheetBMServer), new(*Service)))

// Service service.
type Service struct {
	ac  *paladin.Map
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	return
}

// Close close the resource.
func (s *Service) Close() {
}
