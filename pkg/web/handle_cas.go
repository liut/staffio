package web

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/liut/staffio/pkg/models/cas"
)

const (
	ticketCKey = "cTGT"
)

func NewTGC(c *gin.Context, ticket *cas.Ticket) {
	tgt := cas.NewTicket("TGT", ticket.Service, ticket.UID, false)
	session := ginSession(c)
	session.Set(ticketCKey, tgt)
}

func GetTGC(c *gin.Context) *cas.Ticket {
	session := ginSession(c)
	v := session.Get(ticketCKey)
	if t, ok := v.(*cas.Ticket); ok {
		return t
	}
	return nil
}

func DeleteTGC(c *gin.Context) {
	session := ginSession(c)
	session.Set(ticketCKey, nil)
}

func casLogout(c *gin.Context) {
	tgc := GetTGC(c)
	if tgc != nil {
		DeleteTGC(c)
		fmt.Fprint(c.Writer, "User has been logged out")
	} else {
		fmt.Fprint(c.Writer, "User is not logged in")
	}
}

func (s *server) casValidateV1(c *gin.Context) {
	service := c.Request.FormValue("service")
	ticket := c.Request.FormValue("ticket")

	if ticket == "" {
		fmt.Fprint(c.Writer, "no\n") //nolint
	} else {
		t, err := s.service.GetTicket(ticket)
		if err != nil {
			logger().Infow("load tick fail", "val", ticket, "err", err)
			fmt.Fprint(c.Writer, "no\n") //nolint
		} else {
			if t.Service != service {
				fmt.Fprint(c.Writer, "no\n") //nolint
			} else {
				_ = s.service.DeleteTicket(ticket)  //TODO: fix it
				fmt.Fprint(c.Writer, "yes\n"+t.UID) //nolint
			}
		}
	}
}

func (s *server) casValidateV2(c *gin.Context) {
	service := c.Request.FormValue("service")
	ticket := c.Request.FormValue("ticket")
	format := c.Request.FormValue("format")
	if format != "JSON" {
		format = "XML"
	}

	if format == "XML" {
		c.Writer.Header().Set("Content-Type", "application/xml;charset=UTF-8")
	} else {
		c.Writer.Header().Set("Content-Type", "application/json;charset=UTF-8")
	}
	st, err := s.service.GetTicket(ticket)
	if err != nil {
		logger().Infow("load tick fail", "val", ticket, "err", err)
		//nolint
		fmt.Fprintf(c.Writer, v2ResponseFailure(cas.NewCasError(
			"ticket is invalid or expired", cas.ERROR_CODE_INVALID_TICKET), format))
		return
	}

	casErr := st.Check()
	if casErr == nil {
		casErr = cas.ValidateService(service)
		if casErr == nil {
			if "" == service || st.Service != service {
				casErr = cas.NewCasError("mismatch service", cas.ERROR_CODE_INVALID_SERVICE)
			}
		}
	}

	if casErr != nil {
		fmt.Fprintf(c.Writer, v2ResponseFailure(casErr, format)) //nolint
		logger().Infow("casValidateV2 fail", "uri", c.Request.URL, "err", casErr)
		return
	}

	fmt.Fprintf(c.Writer, v2ResponseSuccess(st, format)) //nolint
}

// v2ResponseFailure produces XML string for failure
func v2ResponseFailure(casError *cas.CasError, format string) string {
	if format == "XML" {
		return fmt.Sprintf(v2ValidationFailureXML, casError.Code, casError.InnerError)
	}

	return fmt.Sprintf(v2ValidationFailureJSON,
		casError.Code, casError.InnerError)
}

// v2ResponseSuccess produces XML string for success
func v2ResponseSuccess(ticket *cas.Ticket, format string) string {
	if format == "XML" {
		return fmt.Sprintf(v2ValidationSuccessXML, ticket.UID, ticket.Value)
	}

	return fmt.Sprintf(v2ValidationSuccessJSON, ticket.UID, ticket.Value)
}

const (
	v2ValidationSuccessXML = `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
 <cas:authenticationSuccess>
  <cas:user>%s</cas:user>
  <cas:proxyGrantingTicket>%s</cas:proxyGrantingTicket>
 </cas:authenticationSuccess>
</cas:serviceResponse>`
	v2ValidationSuccessJSON = `{
  "serviceResponse" : {
    "authenticationSuccess" : {
      "user" : "%s",
      "proxyGrantingTicket" : "%s"
    }
  }
}`
	v2ValidationFailureXML = `<cas:serviceResponse xmlns:cas="http://www.yale.edu/tp/cas">
 <cas:authenticationFailure code="%s">
    %s
  </cas:authenticationFailure>
</cas:serviceResponse>
`
	v2ValidationFailureJSON = `{
  "serviceResponse" : {
    "authenticationFailure" : {
      "code" : "%s",
      "description" : "%s"
    }
  }
}`
)
