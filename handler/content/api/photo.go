package api

import (
	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type PhotoHandler struct {
	PhotoService service.PhotoService
}

func NewPhotoHandler(photoService service.PhotoService) *PhotoHandler {
	return &PhotoHandler{
		PhotoService: photoService,
	}
}

func (p *PhotoHandler) Like(ctx web.Context) (interface{}, error) {
	id, err := util.ParamWebInt32(ctx, "photoID")
	if err != nil {
		return nil, err
	}
	return nil, p.PhotoService.IncreaseLike(ctx.RequestContext(), id)
}
