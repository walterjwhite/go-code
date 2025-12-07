package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/walterjwhite/go-code/lib/application/logging"
)

type Handler struct {
	userSvc  UserService
	validate *validator.Validate
	router   *gin.Engine
}

func NewHandler(userSvc UserService) *Handler {
	v := validator.New()
	r := gin.New()
	r.Use(gin.Recovery())
	r.NoMethod(func(c *gin.Context) { JSONError(c, http.StatusMethodNotAllowed, "method not allowed") })
	return &Handler{userSvc: userSvc, validate: v, router: r}
}

func (h *Handler) Router(dsn string) *gin.Engine {
	h.router.Use(LoggerMiddleware())


	db, err := NewSQLXDBSQLite(dsn)
	logging.Panic(err)

	ch := StartRequestLogWorker(db, 1000)
	h.router.Use(RequestLoggerMiddleware(ch))

	StartDailyRequestReportWorker(db)

	api := h.router.Group("/api")
	{
		users := api.Group("/users")
		users.POST("", h.createUser)
		users.GET("", h.listUsers)
		users.GET("/:id", h.getUser)
		users.PUT("/:id", h.updateUser)
		users.DELETE("/:id", h.deleteUser)
	}

	ServeStaticSPA(h.router, "", "./frontend/dist")

	h.router.NoRoute(NotFoundHandler)
	return h.router
}


func (h *Handler) parseUintParam(c *gin.Context, name string) (uint, bool) {
	idParam := c.Param(name)
	id, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		JSONError(c, http.StatusBadRequest, "invalid id")
		return 0, false
	}
	return uint(id), true
}

func (h *Handler) bindAndValidate(c *gin.Context, v interface{}) bool {
	if err := c.ShouldBindJSON(v); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return false
	}
	if err := h.validate.Struct(v); err != nil {
		JSONError(c, http.StatusBadRequest, err.Error())
		return false
	}
	return true
}

func (h *Handler) createUser(c *gin.Context) {
	var req CreateUserRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	dto, err := h.userSvc.Create(c.Request.Context(), req)
	if err != nil {
		if err == ErrAlreadyExists {
			JSONError(c, http.StatusConflict, err.Error())
			return
		}
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusCreated, dto)
}

func (h *Handler) listUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	sizeStr := c.DefaultQuery("size", "20")
	page, _ := strconv.Atoi(pageStr)
	size, _ := strconv.Atoi(sizeStr)
	list, total, err := h.userSvc.List(c.Request.Context(), page, size)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items": list,
		"meta": gin.H{
			"total": total,
			"page":  page,
			"size":  size,
		},
	})
}

func (h *Handler) getUser(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}

	u, err := h.userSvc.Get(c.Request.Context(), id)
	if err != nil {
		JSONError(c, http.StatusNotFound, err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) updateUser(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	var req UpdateUserRequest
	if !h.bindAndValidate(c, &req) {
		return
	}

	u, err := h.userSvc.Update(c.Request.Context(), id, req)
	if err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *Handler) deleteUser(c *gin.Context) {
	id, ok := h.parseUintParam(c, "id")
	if !ok {
		return
	}
	if err := h.userSvc.Delete(c.Request.Context(), id); err != nil {
		JSONError(c, http.StatusInternalServerError, err.Error())
		return
	}
	c.Status(http.StatusNoContent)
}
