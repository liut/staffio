package web

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

func (s *server) wechatCheckinList(c *gin.Context) {
	days, err := strconv.Atoi(c.DefaultQuery("days", "7"))
	if err != nil {
		apiError(c, 400, err)
		return
	}
	data, err := s.checkin.ListCheckin(days, c.QueryArray("uid")...)
	if err != nil {
		logger().Infow("ListCheckin fail", "err", err)
		apiError(c, 400, err)
		return
	}

	apiOk(c, data.CheckInData, len(data.CheckInData))
}
