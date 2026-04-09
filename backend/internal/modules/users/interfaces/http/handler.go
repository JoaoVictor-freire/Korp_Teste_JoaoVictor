package http

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"korp_backend/internal/modules/users/application"
	"korp_backend/internal/platform/auth"
	"korp_backend/internal/platform/httpx"
)

type Handler struct {
	register application.RegisterUseCase
	login    application.LoginUseCase
	signer   auth.TokenSigner
}

func NewHandler(
	register application.RegisterUseCase,
	login application.LoginUseCase,
	signer auth.TokenSigner,
) Handler {
	return Handler{
		register: register,
		login:    login,
		signer:   signer,
	}
}

func (h Handler) Register(c *gin.Context) {
	var request registerRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.register.Execute(c.Request.Context(), application.RegisterInput{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		switch {
		case errors.Is(err, application.ErrNameRequired),
			errors.Is(err, application.ErrEmailRequired),
			errors.Is(err, application.ErrPasswordRequired):
			httpx.Error(c, http.StatusBadRequest, err.Error())
			return
		case errors.Is(err, application.ErrUserExists):
			httpx.Error(c, http.StatusConflict, err.Error())
			return
		default:
			httpx.Error(c, http.StatusInternalServerError, "failed to register user")
			return
		}
	}

	token, err := h.signer.Sign(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "failed to issue token")
		return
	}

	httpx.JSON(c, http.StatusCreated, gin.H{
		"user":  user,
		"token": token,
	})
}

func (h Handler) Login(c *gin.Context) {
	var request loginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		httpx.Error(c, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.login.Execute(c.Request.Context(), application.LoginInput{
		Email:    request.Email,
		Password: request.Password,
	})
	if err != nil {
		if errors.Is(err, application.ErrInvalidCredentials) {
			httpx.Error(c, http.StatusUnauthorized, err.Error())
			return
		}

		httpx.Error(c, http.StatusInternalServerError, "failed to login")
		return
	}

	token, err := h.signer.Sign(user.ID, user.Email, 24*time.Hour)
	if err != nil {
		httpx.Error(c, http.StatusInternalServerError, "failed to issue token")
		return
	}

	httpx.JSON(c, http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}
