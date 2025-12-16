package dto

type PresignPostImagesRequest struct {
	Count       int    `json:"count"`
	ContentType string `json:"content_type"` 
}

type PresignPostImagesResponse struct {
	Uploads []PresignedUploadDTO `json:"uploads"`
}

type PresignedUploadDTO struct {
	Key          string            `json:"key"`
	UploadURL    string            `json:"upload_url"`
	Headers      map[string]string `json:"headers"`
	TempPublicURL string           `json:"temp_public_url"`
	ExpiresAt    string            `json:"expires_at"`
}
