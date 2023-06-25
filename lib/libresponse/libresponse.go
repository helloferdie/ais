package libresponse

import (
	"encoding/json"
	"math"

	"github.com/helloferdie/golib/libtime"
)

// Response -
type Response struct {
	Success  bool                   `json:"success"`
	Code     int                    `json:"code"`
	Message  string                 `json:"message"`
	Error    string                 `json:"error"`
	ErrorVar map[string]interface{} `json:"error_var"`
	Data     interface{}            `json:"data"`
}

// Error -
type Error struct {
	Error    string                 `json:"error"`
	ErrorVar map[string]interface{} `json:"error_var"`
}

// Pagination -
type Pagination struct {
	Items      []interface{} `json:"items"`
	TotalItems int64         `json:"total_items"`
	TotalPages int64         `json:"total_pages"`
}

// GetDefault - Return default response
func GetDefault() *Response {
	resp := new(Response)
	resp.Success = false
	resp.Code = 500
	resp.Message = "common.error.general"
	return resp
}

// GetFormatOutput - Get format for return response
func GetFormatOutput() map[string]interface{} {
	return map[string]interface{}{
		"timezone": "Asia/Jakarta",
	}
}

// MapOutput - Map output based on provided format
func MapOutput(obj interface{}, stdTimestamp bool, format map[string]interface{}) map[string]interface{} {
	tz, ok := format["timezone"].(string)
	if !ok {
		tz = "UTC"
	}
	databytes, _ := json.Marshal(obj)
	m := map[string]interface{}{}
	json.Unmarshal(databytes, &m)
	if stdTimestamp {
		m["created_at"] = libtime.NullFormat(m["created_at"], tz)
		m["updated_at"] = libtime.NullFormat(m["updated_at"], tz)
		m["deleted_at"] = libtime.NullFormat(m["deleted_at"], tz)
	}
	return m
}

// TotalPages - Calculate total pages given total items & total items per page
func TotalPages(itemPerPage int64, totalItems int64) float64 {
	return math.Ceil(float64(totalItems) / float64(itemPerPage))
}

// SuccessDefault -
func (r *Response) SuccessDefault() *Response {
	r.Success = true
	r.Code = 200
	r.Message = "common.success.default"
	return r
}

// SuccessList -
func (r *Response) SuccessList() *Response {
	r.Success = true
	r.Code = 200
	r.Message = "common.success.list"
	return r
}

// SuccessCreate -
func (r *Response) SuccessCreate() *Response {
	r.Success = true
	r.Code = 200
	r.Message = "common.success.create"
	return r
}

// SuccessUpdate -
func (r *Response) SuccessUpdate() *Response {
	r.Success = true
	r.Code = 200
	r.Message = "common.success.update"
	return r
}

// SuccessDelete -
func (r *Response) SuccessDelete() *Response {
	r.Success = true
	r.Code = 200
	r.Message = "common.success.delete"
	return r
}

// ErrorValidation -
func (r *Response) ErrorValidation() *Response {
	r.Code = 422
	r.Message = "validation.error.default"
	return r
}

// ErrorList -
func (r *Response) ErrorList() *Response {
	r.Code = 500
	r.Message = "common.error.service.list"
	return r
}

// ErrorCreate -
func (r *Response) ErrorCreate() *Response {
	r.Code = 500
	r.Message = "common.error.service.create"
	return r
}

// ErrorUpdate -
func (r *Response) ErrorUpdate() *Response {
	r.Code = 500
	r.Message = "common.error.service.update"
	return r
}

// ErrorDelete -
func (r *Response) ErrorDelete() *Response {
	r.Code = 500
	r.Message = "common.error.service.delete"
	return r
}

// ErrorDataNotFound -
func (r *Response) ErrorDataNotFound() *Response {
	r.Code = 404
	r.Error = "common.error.service.data.not_found"
	return r
}
