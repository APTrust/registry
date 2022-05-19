package webui

import (
	"net/http"

	"github.com/APTrust/registry/common"
	"github.com/gin-gonic/gin"
)

func NsqShow(c *gin.Context) {
	req := NewRequest(c)
	stats, err := common.Context().NSQClient.GetStats()
	if AbortIfError(c, err) {
		return
	}
	req.TemplateData["stats"] = stats
	c.HTML(http.StatusOK, "nsq/show.html", req.TemplateData)
}

func NsqTopicPause(c *gin.Context) {

}

func NsqTopicUnpause(c *gin.Context) {

}

func NsqTopicEmpty(c *gin.Context) {

}

func NsqTopicDelete(c *gin.Context) {

}

func NsqChannelPause(c *gin.Context) {

}

func NsqChannelUnpause(c *gin.Context) {

}

func NsqChannelEmpty(c *gin.Context) {

}

func NsqChannelDelete(c *gin.Context) {

}
