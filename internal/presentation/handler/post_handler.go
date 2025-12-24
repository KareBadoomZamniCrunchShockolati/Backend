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
	userID := mustGetUserID(c)
	input := Validated[dto.CreatePostDTO](c)

	post, err := h.postService.CreatePost(c.Request.Context(), userID, &input)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusCreated, "", post)
}

func (h *PostHandler) GetPost(c *gin.Context) {
	var userID uint = 0
	if v, exists := c.Get("userID"); exists {
		if id, ok := v.(uint); ok {
			userID = id
		}
	}

	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	post, err := h.postService.GetPost(postID, userID)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", post)
}

func (h *PostHandler) UpdatePost(c *gin.Context) {
	userID := mustGetUserID(c)
	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	input := Validated[dto.UpdatePostDTO](c)

	post, err := h.postService.UpdatePost(postID, userID, &input)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", post)
}

func (h *PostHandler) DeletePost(c *gin.Context) {
	userID := mustGetUserID(c)
	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	if err := h.postService.DeletePost(postID, userID); err != nil {
		panic(err)
	}

	Response(c, http.StatusNoContent, "", nil)
}

func (h *PostHandler) GetUserPosts(c *gin.Context) {
	targetUserID := mustParseUintParam(c, "user_id", "INVALID_USER_ID")

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetUserPosts(targetUserID, offset, limit)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", posts)
}

func (h *PostHandler) GetFeedPosts(c *gin.Context) {
	userID := mustGetUserID(c)

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetFeedPosts(userID, offset, limit)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", posts)
}

func (h *PostHandler) GetPostsByChallenge(c *gin.Context) {
	var userID uint = 0
	if v, exists := c.Get("userID"); exists {
		if id, ok := v.(uint); ok {
			userID = id
		}
	}

	challengeID := mustParseUintParam(c, "challenge_id", "INVALID_CHALLENGE_ID")

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	posts, err := h.postService.GetPostsByChallenge(challengeID, userID, offset, limit)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", posts)
}

func (h *PostHandler) AddPostComment(c *gin.Context) {
	userID := mustGetUserID(c)
	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	type req struct {
		Content  string `json:"content" validate:"required,min=1,max=1000"`
		ParentID *uint  `json:"parent_id"`
	}
	input := Validated[req](c)

	commentDTO := &dto.CommentRequestDTO{
		EntityType: "post",
		EntityID:   postID,
		Content:    input.Content,
		ParentID:   input.ParentID,
	}

	comment, err := h.postService.AddComment(userID, commentDTO)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusCreated, "", comment)
}

func (h *PostHandler) GetPostComments(c *gin.Context) {
	var userID uint = 0
	if v, exists := c.Get("userID"); exists {
		if id, ok := v.(uint); ok {
			userID = id
		}
	}

	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	comments, err := h.postService.GetComments("post", postID, userID, offset, limit)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", comments)
}

func (h *PostHandler) LikePost(c *gin.Context) {
	userID := mustGetUserID(c)
	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	likeDTO := &dto.LikeRequestDTO{
		EntityType: "post",
		EntityID:   postID,
	}

	if err := h.postService.LikeEntity(userID, likeDTO); err != nil {
		panic(err)
	}

	Response(c, http.StatusNoContent, "", nil)
}

func (h *PostHandler) UnlikePost(c *gin.Context) {
	userID := mustGetUserID(c)
	postID := mustParseUintParam(c, "id", "INVALID_POST_ID")

	likeDTO := &dto.LikeRequestDTO{
		EntityType: "post",
		EntityID:   postID,
	}

	if err := h.postService.UnlikeEntity(userID, likeDTO); err != nil {
		panic(err)
	}

	Response(c, http.StatusNoContent, "", nil)
}

func (h *PostHandler) PresignPostImages(c *gin.Context) {
	userID := mustGetUserID(c)

	req := Validated[dto.PresignPostImagesRequest](c)

	resp, err := h.postService.PresignPostImages(c.Request.Context(), userID, req)
	if err != nil {
		panic(err)
	}

	Response(c, http.StatusOK, "", resp)
}

func mustGetUserID(c *gin.Context) uint {
	v, exists := c.Get("userID")
	if !exists {
		panic(exception.NewMissingUserIDException())
	}
	id, ok := v.(uint)
	if !ok || id == 0 {
		panic(exception.NewContextCastError(nil))
	}
	return id
}

func mustParseUintParam(c *gin.Context, paramName string, code string) uint {
	raw := c.Param(paramName)
	v, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		panic(exception.NewBadRequestException(
			code,
			map[string]any{"param": paramName, "value": raw},
		))
	}
	return uint(v)
}
