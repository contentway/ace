package ace

import "math"

const (
	// 2xx status codes
	StatusOK       = 200
	StatusCreated  = 201
	StatusAccepted = 202

	// 4xx status codes
	StatusBadRequest                  = 400
	StatusUnauthorized                = 401
	StatusPaymentRequired             = 402
	StatusForbidden                   = 403
	StatusNotFound                    = 404
	StatusMethodNotAllowed            = 405
	StatusNotAcceptable               = 406
	StatusProxyAuthenticationRequired = 407
	StatusRequestTimeout              = 408
	StatusConflict                    = 409
	StatusGone                        = 410
	StatusLengthRequired              = 411
	StatusPreconditionFailed          = 412
	StatusRequestEntityTooLarge       = 413
	StatusURITooLong                  = 414
	StatusUnsupportedMediaType        = 415
	StatusRangeNotSatisfiable         = 416
	StatusExpectationFailed           = 417
	StatusImATeapot                   = 418
	StatusUnprocessableEntity         = 422
	StatusLocked                      = 423

	// 5xx status codes
	StatusInternalServerError = 500
	StatusNotImplemented      = 501
	StatusBadGateway          = 502
	StatusServiceUnavailable  = 503
	StatusGatewayTimeout      = 504

	HeaderAuthorization      = "Authorization"
	HeaderContentType        = "Content-Type"
	HeaderContentDisposition = "ContentDisposition"

	AbortMiddlewareIndex = math.MaxInt8 / 2
)
