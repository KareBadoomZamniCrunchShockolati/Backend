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
	ChallengeServicer serviceinterface.ChallengeServicer
}

func NewChallengeHandler(Challengeservice serviceinterface.ChallengeServicer) *ChallengeHandler {
	return &ChallengeHandler{ChallengeServicer: Challengeservice}
}


func (h *ChallengeHandler) CreateChallenge(c *gin.Context) {
	var input dto.CreateChallengeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	challenge, err := h.ChallengeServicer.CreateChallenge(&input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	responseDTO, err := h.ChallengeServicer.ToChallengeResponseDTO(challenge, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build response"})
		return
	}

	c.JSON(http.StatusCreated, responseDTO)
}

func (h *ChallengeHandler) GetChallengeByID(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	challenge, err := h.ChallengeServicer.GetChallengeByID(id)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	responseDTO, err := h.ChallengeServicer.ToChallengeResponseDTO(challenge, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build response"})
		return
	}

	c.JSON(http.StatusOK, responseDTO)
}

func (h *ChallengeHandler) GetAllChallenges(c *gin.Context) {
	challenges, err := h.ChallengeServicer.GetAllChallenges()
	if err != nil {
		handleServiceError(c, err)
		return
	}

	responseDTOs, err := h.ChallengeServicer.ToChallengeResponseDTOs(challenges, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build response"})
		return
	}

	c.JSON(http.StatusOK, responseDTOs)
}

func (h *ChallengeHandler) UpdateChallenge(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	var input dto.UpdateChallengeDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	updated, err := h.ChallengeServicer.UpdateChallenge(id, getCurrentUserID(c), &input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	responseDTO, err := h.ChallengeServicer.ToChallengeResponseDTO(updated, getCurrentUserID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to build response"})
		return
	}

	c.JSON(http.StatusOK, responseDTO)
}

func (h *ChallengeHandler) DeleteChallenge(c *gin.Context) {
	id, err := parseIDParam(c, "id")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid challenge id"})
		return
	}

	if err := h.ChallengeServicer.DeleteChallenge(id, getCurrentUserID(c)); err != nil {
		handleServiceError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}


func getCurrentUserID(c *gin.Context) uint {
	userID, exists := c.Get("currentUserID")
	if !exists {
		return 0
	}
	return userID.(uint)
}

func parseIDParam(c *gin.Context, param string) (uint, error) {
	idStr := c.Param(param)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
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
