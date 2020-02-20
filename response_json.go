package osin

import (
	"github.com/labstack/echo"
)

// OutputJSON encodes the Response to JSON and writes to the http.ResponseWriter
func OutputJSON(rs *Response, c echo.Context) error {
	// Add headers

	for i, k := range rs.Headers {
		for _, v := range k {
			c.Response().Header().Add(i, v)
		}
	}

	if rs.Type == REDIRECT {
		// Output redirect with parameters
		u, err := rs.GetRedirectUrl()
		if err != nil {
			return err
		}
		return c.Redirect(302, u)
	}

	// set content type if the response doesn't already have one associated with it
	if c.Response().Header().Get("Content-Type") == "" {
		c.Response().Header().Set("Content-Type", "application/json")
	}
	return c.JSON(rs.StatusCode, rs.Output)
}
