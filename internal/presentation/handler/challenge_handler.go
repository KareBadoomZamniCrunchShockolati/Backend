package handler

import (
	"net/http"
	"strconv"

	"challenge-app/internal/application/dto"
	"challenge-app/internal/domain/exception"
	serviceinterface "challenge-app/internal/application/service/interface"

	"github.com/gin-gonic/gin"
)

type ChallengeHandler struct {
	Challengeservice serviceinterface.ChallengeServicer
}

func NewChallengeHandler(Challengeservice serviceinterface.ChallengeServicer) *ChallengeHandler {
	return &ChallengeHandler{Challengeservice: Challengeservice}
}

func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	var input dto.CreateChallengeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	challenge, err := h.Challengeservice.CreateChallenge(&input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusCreated, challenge)
}

func (h *ChallengeHandler) GetChallengeByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	challenge, err := h.Challengeservice.GetChallengeByID(uint(id))
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	var input dto.UpdateChallengeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	updated, err := h.Challengeservice.UpdateChallenge(uint(id), getCurrentUserID(c), &input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	c.JSON(http.StatusOK, updated)
}

func (h *ChallengeHandler) DeleteChallenge(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	if err := h.Challengeservice.DeleteChallenge(uint(id), getCurrentUserID(c)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}


// --- Helper functions ---

func getPagination(c *gin.Context) (offset, limit int) {
	offsetStr := c.Query("offset")
	limitStr := c.Query("limit")

	offset = 0
	limit = 10

	if offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil {
			offset = o
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil {
			limit = l
		}
	}
	return
}

func getCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("currentUserID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func handleServiceError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *exception.NotFoundException:
		c.JSON(http.StatusNotFound, gin.H{"error": e.Error()})
	case *exception.BadRequestException:
		c.JSON(http.StatusBadRequest, gin.H{"error": e.Error()})
	case *exception.UnauthorizedException:
		c.JSON(http.StatusUnauthorized, gin.H{"error": e.Error()})
	case *exception.ConflictException:
		c.JSON(http.StatusConflict, gin.H{"error": e.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}
