package web

import (
	"fmt"
	"log"

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
		fmt.Fprint(c.Writer, "no\n")
	} else {
		t, err := s.service.GetTicket(ticket)
		if err != nil {
			log.Printf("load ticket %s ERR: %s", ticket, err)
			fmt.Fprint(c.Writer, "no\n")
		} else {
			if t.Service != service {
				fmt.Fprint(c.Writer, "no\n")
			} else {
				s.service.DeleteTicket(ticket)
				fmt.Fprint(c.Writer, "yes\n"+t.Uid)
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
		fmt.Fprintf(c.Writer, v2ResponseFailure(cas.NewCasError(
			"ticket is invalid or expired", cas.ERROR_CODE_INVALID_TICKET), format))
		log.Printf("casValidateV2 %s ERR: %s", c.Request.URL, err)
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
		fmt.Fprintf(c.Writer, v2ResponseFailure(casErr, format))
		log.Printf("casValidateV2 %s ERR: %s", c.Request.URL, casErr)
		return
	}

	fmt.Fprintf(c.Writer, v2ResponseSuccess(st, format))
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
		return fmt.Sprintf(v2ValidationSuccessXML, ticket.Uid, ticket.Value)
	}

	return fmt.Sprintf(v2ValidationSuccessJSON, ticket.Uid, ticket.Value)
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
