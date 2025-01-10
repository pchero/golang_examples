package main

import (
	"fmt"
	"net/http"
	customer_manager "testoapi/gens/models/customer_manager"
	"testoapi/gens/server"

	"github.com/gin-gonic/gin"
)

func (h *Server) GetAccesskeys(c *gin.Context, params server.GetAccesskeysParams) {
	pageSize := 0
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}
	fmt.Printf("request detail. page_size: %d, page_token: %s, request: %v\n", pageSize, pageToken, params)

	result := []customer_manager.Accesskey{
		{
			Id: strPtr("619cfcd0-ce2e-11ef-8115-fba7cc388145"),
		},
	}

	c.JSON(http.StatusOK, server.GetAccesskeys200JSONResponse{
		NextPageToken: strPtr("next_page_namge"),
		Result:        &result,
	})
}

func (h *Server) PostAccesskeys(c *gin.Context) {

	req := server.PostAccesskeysJSONBody{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	ID := "c75884ee-ce2a-11ef-9b36-93c4817072e0"
	tmp := customer_manager.Accesskey{
		Id: &ID,
	}

	res := server.PostAccesskeys201JSONResponse(tmp)
	c.JSON(http.StatusOK, res)
}

func (h *Server) DeleteAccesskeysId(c *gin.Context, id string) {
	res := server.DeleteAccesskeysId204Response{}
	c.JSON(http.StatusNoContent, res)
}

func (h *Server) GetAccesskeysId(c *gin.Context, id string) {
	tmp := customer_manager.Accesskey{
		Id: &id,
	}

	res := server.GetAccesskeysId200JSONResponse(tmp)
	c.JSON(http.StatusOK, res)
}

func (h *Server) PutAccesskeysId(c *gin.Context, id string) {
	tmp := customer_manager.Accesskey{
		Id: &id,
	}

	res := server.PutAccesskeysId200JSONResponse(tmp)
	c.JSON(http.StatusOK, res)
}
