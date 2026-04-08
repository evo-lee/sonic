package admin

import (
	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/service/ai"
	"github.com/go-sonic/sonic/util/xerr"
)

type AIHandler struct {
	contentService ai.ContentService
}

func NewAIHandler(contentService ai.ContentService) *AIHandler {
	return &AIHandler{contentService: contentService}
}

type summarizeRequest struct {
	Content string `json:"content"`
}

type suggestTagsRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type polishRequest struct {
	Content string `json:"content"`
}

func (h *AIHandler) Summarize(ctx web.Context) (interface{}, error) {
	var req summarizeRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if req.Content == "" {
		return nil, xerr.BadParam.New("content is required")
	}
	summary, err := h.contentService.Summarize(ctx.RequestContext(), req.Content)
	if err != nil {
		return nil, err
	}
	return map[string]string{"summary": summary}, nil
}

func (h *AIHandler) SuggestTags(ctx web.Context) (interface{}, error) {
	var req suggestTagsRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if req.Content == "" {
		return nil, xerr.BadParam.New("content is required")
	}
	tags, err := h.contentService.SuggestTags(ctx.RequestContext(), req.Title, req.Content)
	if err != nil {
		return nil, err
	}
	return map[string]interface{}{"tags": tags}, nil
}

func (h *AIHandler) Polish(ctx web.Context) (interface{}, error) {
	var req polishRequest
	if err := ctx.BindJSON(&req); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("parameter error")
	}
	if req.Content == "" {
		return nil, xerr.BadParam.New("content is required")
	}
	polished, err := h.contentService.Polish(ctx.RequestContext(), req.Content)
	if err != nil {
		return nil, err
	}
	return map[string]string{"content": polished}, nil
}
