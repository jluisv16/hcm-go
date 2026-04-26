package httpapi

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/jluisv16/hcm-go/internal/employees/application"
	"github.com/jluisv16/hcm-go/internal/employees/domain"
)

const hireDateLayout = "2006-01-02"

type Handler struct {
	service *application.Service
}

type upsertEmployeeRequest struct {
	FirstName  string  `json:"first_name" binding:"required"`
	LastName   string  `json:"last_name" binding:"required"`
	Email      string  `json:"email" binding:"required"`
	Department string  `json:"department" binding:"required"`
	Role       string  `json:"role" binding:"required"`
	Salary     float64 `json:"salary" binding:"required"`
	HireDate   string  `json:"hire_date" binding:"required"`
}

type employeeResponse struct {
	ID         string  `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	Email      string  `json:"email"`
	Department string  `json:"department"`
	Role       string  `json:"role"`
	Salary     float64 `json:"salary"`
	HireDate   string  `json:"hire_date"`
}

func NewHandler(service *application.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Register(v1 *gin.RouterGroup) {
	employeeGroup := v1.Group("/employees")
	employeeGroup.GET("", h.List)
	employeeGroup.GET("/:id", h.GetByID)
	employeeGroup.POST("", h.Create)
	employeeGroup.PUT("/:id", h.Update)
	employeeGroup.DELETE("/:id", h.Delete)
}

func (h *Handler) List(ctx *gin.Context) {
	employees, err := h.service.List(ctx.Request.Context())
	if err != nil {
		writeError(ctx, http.StatusInternalServerError, "could not list employees")
		return
	}

	response := make([]employeeResponse, 0, len(employees))
	for _, employee := range employees {
		response = append(response, toEmployeeResponse(employee))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"employees": response,
	})
}

func (h *Handler) GetByID(ctx *gin.Context) {
	employee, err := h.service.GetByID(ctx.Request.Context(), ctx.Param("id"))
	if err != nil {
		writeDomainError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toEmployeeResponse(employee))
}

func (h *Handler) Create(ctx *gin.Context) {
	var request upsertEmployeeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeError(ctx, http.StatusBadRequest, "invalid payload")
		return
	}

	hireDate, err := time.Parse(hireDateLayout, strings.TrimSpace(request.HireDate))
	if err != nil {
		writeError(ctx, http.StatusBadRequest, "hire_date must use YYYY-MM-DD format")
		return
	}

	input := application.UpsertEmployeeInput{
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Role:       request.Role,
		Salary:     request.Salary,
		HireDate:   hireDate,
	}

	employee, err := h.service.Create(ctx.Request.Context(), input)
	if err != nil {
		writeDomainError(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, toEmployeeResponse(employee))
}

func (h *Handler) Update(ctx *gin.Context) {
	var request upsertEmployeeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		writeError(ctx, http.StatusBadRequest, "invalid payload")
		return
	}

	hireDate, err := time.Parse(hireDateLayout, strings.TrimSpace(request.HireDate))
	if err != nil {
		writeError(ctx, http.StatusBadRequest, "hire_date must use YYYY-MM-DD format")
		return
	}

	input := application.UpsertEmployeeInput{
		ID:         ctx.Param("id"),
		FirstName:  request.FirstName,
		LastName:   request.LastName,
		Email:      request.Email,
		Department: request.Department,
		Role:       request.Role,
		Salary:     request.Salary,
		HireDate:   hireDate,
	}

	employee, err := h.service.Update(ctx.Request.Context(), input)
	if err != nil {
		writeDomainError(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, toEmployeeResponse(employee))
}

func (h *Handler) Delete(ctx *gin.Context) {
	if err := h.service.Delete(ctx.Request.Context(), ctx.Param("id")); err != nil {
		writeDomainError(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func writeDomainError(ctx *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmployeeNotFound):
		writeError(ctx, http.StatusNotFound, err.Error())
	case errors.Is(err, domain.ErrEmailAlreadyInUse):
		writeError(ctx, http.StatusConflict, err.Error())
	default:
		writeError(ctx, http.StatusBadRequest, err.Error())
	}
}

func writeError(ctx *gin.Context, statusCode int, message string) {
	ctx.JSON(statusCode, gin.H{
		"error": message,
	})
}

func toEmployeeResponse(employee domain.Employee) employeeResponse {
	return employeeResponse{
		ID:         employee.ID,
		FirstName:  employee.FirstName,
		LastName:   employee.LastName,
		Email:      employee.Email,
		Department: employee.Department,
		Role:       employee.Role,
		Salary:     employee.Salary,
		HireDate:   employee.HireDate.Format(hireDateLayout),
	}
}
