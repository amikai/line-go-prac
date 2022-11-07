package linebot

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/linebot"

	"github.com/amikai/line-go-prac/internal/dao"
	"github.com/amikai/line-go-prac/pkg/linebotkit"
	"github.com/amikai/line-go-prac/pkg/logkit"
)

type Handler struct {
	linebotClient *linebotkit.LinebotClient
	svc           Service
	logger        *logkit.Logger
}

func NewGinHandler(
	linebotClient *linebotkit.LinebotClient,
	svc Service,
	logger *logkit.Logger,
) *Handler {
	handler := &Handler{
		linebotClient: linebotClient,
		svc:           svc,
		logger:        logger,
	}

	return handler
}

func (g *Handler) Webhook(ctx *gin.Context) {
	events, err := g.linebotClient.ParseRequest(ctx.Request)
	if err == linebot.ErrInvalidSignature {
		abortWithErrMsg(ctx, http.StatusBadRequest, err.Error())
		return
	}
	if err != nil {
		abortWithErrMsg(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	for _, event := range events {
		err = g.svc.HandleEvent(ctx.Request.Context(), event)
		if err != nil {
			ctx.abortWithErrMsg(http.StatusInternalServerError, err)
			return
		}
	}
}

func (g *Handler) GetUserByID(ctx *gin.Context) {
	var req GetUserByIDRequest
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		abortWithErrMsg(ctx, http.StatusBadRequest, "bad request")
	}

	userID := ctx.Param("userID")
	if userID == "" {
		abortWithErrMsg(ctx, http.StatusBadRequest, "missing user id")
	}

	daoMessages, err := g.svc.GetUserByID(ctx, userID, dao.MessagePagination{
		Count: req.Count,
		After: time.Unix(req.After, 0),
	})
	if err != nil {
		return
	}

	messages := make([]*Message, 0, len(daoMessages))
	for _, daoMessage := range daoMessages {
		var message Message
		message.FromMessageDAO(daoMessage)
		messages = append(messages, &message)
	}

	res := &GetUserByIDResponse{
		Messages: messages,
	}
	ctx.JSON(http.StatusOK, res)
}

func abortWithErrMsg(c *gin.Context, code int, message string) {
	c.AbortWithStatusJSON(code, gin.H{
		"code":    code,
		"message": message,
	})
}
