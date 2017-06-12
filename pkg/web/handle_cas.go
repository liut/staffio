package web

import (
	"fmt"
	"log"

	"github.com/gorilla/sessions"

	"lcgc/platform/staffio/pkg/backends"
	"lcgc/platform/staffio/pkg/models/cas"
)

const (
	ticketCKey = "cTGT"
)

func NewTGC(sess *sessions.Session, ticket *cas.Ticket) {
	tgt := cas.NewTicket("TGT", ticket.Service, ticket.Uid, false)
	sess.Values[ticketCKey] = tgt
}

func GetTGC(sess *sessions.Session) *cas.Ticket {
	if v, ok := sess.Values[ticketCKey]; ok {
		return v.(*cas.Ticket)
	}
	return nil
}

func DeleteTGC(sess *sessions.Session) {
	delete(sess.Values, ticketCKey)
}

func casLogout(c *Context) error {
	tgc := GetTGC(c.Session)
	if tgc != nil {
		DeleteTGC(c.Session)
		fmt.Fprint(c.Writer, "User has been logged out")
	} else {
		fmt.Fprint(c.Writer, "User is not logged in")
	}
	return nil
}

func casValidateV1(c *Context) error {
	service := c.Request.FormValue("service")
	ticket := c.Request.FormValue("ticket")

	if ticket == "" {
		fmt.Fprint(c.Writer, "no\n")
	} else {
		t, err := backends.GetTicket(ticket)
		if err != nil {
			log.Printf("load ticket %s ERR: %s", ticket, err)
			fmt.Fprint(c.Writer, "no\n")
		} else {
			if t.Service != service {
				fmt.Fprint(c.Writer, "no\n")
			} else {
				backends.DeleteTicket(ticket)
				fmt.Fprint(c.Writer, "yes\n"+t.Uid)
			}
		}
	}
	return nil
}

func casValidateV2(c *Context) error {
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
	st, err := backends.GetTicket(ticket)
	if err != nil {
		fmt.Fprintf(c.Writer, v2ResponseFailure(cas.NewCasError(
			"ticket is invalid or expired", cas.ERROR_CODE_INVALID_TICKET), format))
		log.Printf("casValidateV2 %s ERR: %s", c.Request.URL, err)
		return nil
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
		return nil
	}

	fmt.Fprintf(c.Writer, v2ResponseSuccess(st, format))
	return nil
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
