package device

type DeviceCreationRequest struct {
	PushToken string `json:"push_token" validate:"required"`
	Brand     string `json:"brand,omitempty"`
	Model     string `json:"model,omitempty"`
	OS        string `json:"os,omitempty"`
	OSVersion string `json:"os_version,omitempty"`
}

type DeviceActivationRequest struct {
	PushToken string `json:"push_token" validate:"required"`
	IsActive  *bool  `json:"is_active" validate:"required"`
}
