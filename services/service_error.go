// Convert AppError to ServiceError for compatibility
package services

import "github.com/AlsoShantanuBorkar/budget_max/errors"

type ServiceError struct {
	Code    int
	Message string
}

func (e *ServiceError) Error() string {
	return e.Message
}

func NewServiceError(code int, message string) *ServiceError {
	return &ServiceError{
		Code:    code,
		Message: message,
	}

}

func ServiceErrorFromAppError(appErr *errors.AppError) *ServiceError {
       if appErr == nil {
              return nil
       }
       return &ServiceError{
              Code:    appErr.Code,
              Message: appErr.Message,
       }
}
