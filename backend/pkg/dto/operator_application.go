package dto

import "time"

type TakeApplicationRequest struct {
	ID int64 `json:"id" binding:"required"`
}

type ChangeStatusRequest struct {
	ID         int64   `json:"id" binding:"required"`
	StatusCode string  `json:"status_code" binding:"required"`
	Comment    *string `json:"comment,omitempty"`
}

type CloseApplicationRequest struct {
	ID      int64   `json:"id" binding:"required"`
	Comment *string `json:"comment,omitempty"`
}

type HistoryActor struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type ApplicationHistoryItem struct {
	ID        int64        `json:"id"`
	Action    string       `json:"action"`
	Field     *string      `json:"field,omitempty"`
	OldValue  *string      `json:"old_value,omitempty"`
	NewValue  *string      `json:"new_value,omitempty"`
	Actor     HistoryActor `json:"actor"`
	CreatedAt time.Time    `json:"created_at"`
}

type ApplicationHistoryResponse struct {
	Items []ApplicationHistoryItem `json:"items"`
}

