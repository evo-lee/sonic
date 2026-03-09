package api

import (
	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util"
)

type CommentHandler struct {
	BaseCommentService service.BaseCommentService
}

func NewCommentHandler(baseCommentService service.BaseCommentService) *CommentHandler {
	return &CommentHandler{
		BaseCommentService: baseCommentService,
	}
}

func (c *CommentHandler) Like(ctx web.Context) (interface{}, error) {
	commentID, err := util.ParamWebInt32(ctx, "commentID")
	if err != nil {
		return nil, err
	}
	return nil, c.BaseCommentService.IncreaseLike(ctx.RequestContext(), commentID)
}
