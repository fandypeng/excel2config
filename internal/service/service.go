package service

import (
	pb "excel2config/api"
	"excel2config/internal/dao"
	"excel2config/internal/def"
	"excel2config/internal/service/auth"
	"github.com/go-kratos/kratos/pkg/conf/paladin"
	"github.com/go-kratos/kratos/pkg/ecode"
	"strings"

	"github.com/google/wire"
)

var Provider = wire.NewSet(wire.Bind(new(pb.SheetBMServer), new(*Service)), New)

// Service service.
type Service struct {
	As  *auth.Service
	ac  *paladin.Map
	dao dao.Dao
}

// New new a service and return.
func New(d dao.Dao) (s *Service, cf func(), err error) {
	s = &Service{
		As:  auth.NewAuthService(d),
		ac:  &paladin.TOML{},
		dao: d,
	}
	cf = s.Close
	err = paladin.Watch("application.toml", s.ac)
	s.registerEcodes()
	return
}

// Close close the resource.
func (s *Service) Close() {
}

func (s Service) registerEcodes() {
	cms := map[int]string{}
	for i := def.ErrCodeStart; i > def.ErrCodeEnd; i-- {
		if !strings.HasPrefix(i.String(), "ErrorCode") {
			cms[int(i)] = i.String()
		}
	}
	ecode.Register(cms)
}
