package admin

import (
	"github.com/go-sonic/sonic/handler/web"
	"github.com/go-sonic/sonic/model/param"
	"github.com/go-sonic/sonic/service"
	"github.com/go-sonic/sonic/util/xerr"
)

type EmailHandler struct {
	EmailService service.EmailService
}

func NewEmailHandler(emailService service.EmailService) *EmailHandler {
	return &EmailHandler{
		EmailService: emailService,
	}
}

func (e *EmailHandler) Test(ctx web.Context) (interface{}, error) {
	p := &param.TestEmail{}
	if err := ctx.BindJSON(p); err != nil {
		return nil, xerr.WithStatus(err, xerr.StatusBadRequest).WithMsg("param error ")
	}
	return nil, e.EmailService.SendTextEmail(ctx.RequestContext(), p.To, p.Subject, p.Content)
}
