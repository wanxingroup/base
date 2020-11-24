package errors

type Error struct {
	ErrorCode    uint32 `json:"errorCode"`
	ErrorMessage string `json:"errorMessage"`
}
