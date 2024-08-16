// TABI-BOOKING - Microservice for Tabi project.
//
// API documents for Tabi-booking.
//
// ## Authentication
// Firstly, grab the **access_token** from the response of `/login`. Then include this header in all API calls:
// ```
// Authorization: Bearer ${access_token}
// ```
//
// For testing directly on this Swagger page, use the `Authorize` button right here bellow.
//
// Terms Of Service: N/A
//
//	Host: %{HOST}
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- bearer: []
//
//	SecurityDefinitions:
//	bearer:
//	     type: apiKey
//	     name: Authorization
//	     in: header
//
// swagger:meta
package main

import (
	httpcore "github.com/namhoai1109/tabi/core/http"
	"github.com/namhoai1109/tabi/core/server"
)

// swagger:parameters adminUserList PartnerBranchList PartnerRoomTypeListAll PublicBranchList PartnerBookingList UserBookingList
type ListRequest struct {
	httpcore.ListRequest
}

// Success empty response
// swagger:response ok
type SwaggOKResp struct{}

// Error empty response
// swagger:response err
type SwaggErrResp struct{}

// Error response with details
// swagger:response errDetails
type SwaggErrDetailsResp struct {
	//in: body
	Body server.ErrorResponse
}
