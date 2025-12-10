package handler

import (
	"challenge-app/internal/application/dto"
	"challenge-app/internal/application/service"
	"challenge-app/internal/domain/exception"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{
		postService: postService,
	}
}

func (h *PostHandler) CreatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	var input dto.CreatePostDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid input", "INVALID_INPUT", nil))
		return
	}

	post, err := h.postService.CreatePost(userID.(uint), &input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	userID, _ := c.Get("userID")

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid post ID", "INVALID_POST_ID", nil))
		return
	}

	post, err := h.postService.GetPost(uint(postID), userID.(uint))
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid post ID", "INVALID_POST_ID", nil))
		return
	}

	var input dto.UpdatePostDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid input", "INVALID_INPUT", nil))
		return
	}

	post, err := h.postService.UpdatePost(uint(postID), userID.(uint), &input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	postID, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid post ID", "INVALID_POST_ID", nil))
		return
	}

	err = h.postService.DeletePost(uint(postID), userID.(uint))
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	targetUserID, err := strconv.ParseUint(c.Param("user_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid user ID", "INVALID_USER_ID", nil))
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetUserPosts(uint(targetUserID), offset, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetFeedPosts(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetFeedPosts(userID.(uint), offset, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) GetPostsByChallenge(c *gin.Context) {
	userID, _ := c.Get("userID")

	challengeID, err := strconv.ParseUint(c.Param("challenge_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid challenge ID", "INVALID_CHALLENGE_ID", nil))
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetPostsByChallenge(uint(challengeID), userID.(uint), offset, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, posts)
}

func (h *PostHandler) AddComment(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	var input dto.CommentRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid input", "INVALID_INPUT", nil))
		return
	}

	comment, err := h.postService.AddComment(userID.(uint), &input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, comment)
}

func (h *PostHandler) GetComments(c *gin.Context) {
	userID, _ := c.Get("userID")

	entityType := c.Query("entity_type")
	entityID, err := strconv.ParseUint(c.Query("entity_id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid entity ID", "INVALID_ENTITY_ID", nil))
		return
	}

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, err := h.postService.GetComments(entityType, uint(entityID), userID.(uint), offset, limit)
	if err != nil {
		handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, comments)
}

func (h *PostHandler) LikeEntity(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	var input dto.LikeRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid input", "INVALID_INPUT", nil))
		return
	}

	err := h.postService.LikeEntity(userID.(uint), &input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *PostHandler) UnlikeEntity(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, exception.NewUnauthorizedException("User not authenticated", "UNAUTHORIZED"))
		return
	}

	var input dto.LikeRequestDTO
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, exception.NewBadRequestException("Invalid input", "INVALID_INPUT", nil))
		return
	}

	err := h.postService.UnlikeEntity(userID.(uint), &input)
	if err != nil {
		handleError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

// Helper function to handle errors
func handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *exception.BadRequestException:
		c.JSON(http.StatusBadRequest, e)
	case *exception.UnauthorizedException:
		c.JSON(http.StatusUnauthorized, e)
	case *exception.ForbiddenException:
		c.JSON(http.StatusForbidden, e)
	case *exception.NotFoundException:
		c.JSON(http.StatusNotFound, e)
	case *exception.ConflictException:
		c.JSON(http.StatusConflict, e)
	default:
		c.JSON(http.StatusInternalServerError, exception.NewInternalServerException("Internal server error", "INTERNAL_SERVER_ERROR", nil))
	}
}
