package admin

import (
	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/service"
)

type StatisticHandler struct {
	StatisticService service.StatisticService
}

func NewStatisticHandler(l service.StatisticService) *StatisticHandler {
	return &StatisticHandler{
		StatisticService: l,
	}
}

func (s *StatisticHandler) Statistics(ctx web.Context) (interface{}, error) {
	return s.StatisticService.Statistic(ctx.RequestContext())
}

func (s *StatisticHandler) StatisticsWithUser(ctx web.Context) (interface{}, error) {
	return s.StatisticService.StatisticWithUser(ctx.RequestContext())
}
