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
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err != nil {
		ctx.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	for _, event := range events {
		err = g.svc.HandleEvent(ctx.Request.Context(), event)
		if err != nil {
			ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}
}

func (g *Handler) GetUserByID(ctx *gin.Context) {
	var req GetUserByIDRequest

	err := ctx.ShouldBindQuery(&req)
	userID := ctx.Param("userID")
	if err != nil || userID == "" {
		ctx.AbortWithError(http.StatusBadRequest, err)
	}
	daoMessages, err := g.svc.GetUserByID(ctx.Request.Context(), userID, dao.MessagePagination{
		Count: req.Count,
		After: time.Unix(req.After, 0),
	})

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