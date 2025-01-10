package main

import (
	"fmt"
	"net/http"
	"testoapi/gens/models/flow_manager"
	"testoapi/gens/server"

	"github.com/gin-gonic/gin"
)

func (h *Server) PostActiveflows(c *gin.Context) {
	req := server.PostActiveflowsJSONBody{}
	if err := c.BindJSON(&req); err != nil {
		c.AbortWithStatus(400)
		return
	}

	id := "c75884ee-ce2a-11ef-9b36-93c4817072e0"
	tmp := flow_manager.Activeflow{
		Id: &id,
	}

	res := server.PostActiveflows201JSONResponse(tmp)
	c.JSON(http.StatusCreated, res)
}

func (h *Server) GetActiveflows(c *gin.Context, params server.GetActiveflowsParams) {
	pageSize := 0
	if params.PageSize != nil {
		pageSize = *params.PageSize
	}

	pageToken := ""
	if params.PageToken != nil {
		pageToken = *params.PageToken
	}
	fmt.Printf("request detail. page_size: %d, page_token: %s, request: %v\n", pageSize, pageToken, params)

	result := []flow_manager.Activeflow{
		{
			Id: strPtr("619cfcd0-ce2e-11ef-8115-fba7cc388145"),
		},
	}

	c.JSON(http.StatusCreated, server.GetActiveflows200JSONResponse{
		NextPageToken: strPtr("next_page_token"),
		Result:        &result,
	})
}

func (h *Server) DeleteActiveflowsId(c *gin.Context, id string) {
	res := server.DeleteActiveflowsId204Response{}
	c.JSON(http.StatusNoContent, res)
}

func (h *Server) GetActiveflowsId(c *gin.Context, id string) {
	tmp := flow_manager.Activeflow{
		Id: &id,
	}

	res := server.GetActiveflowsId200JSONResponse(tmp)
	c.JSON(http.StatusOK, res)
}

func (h *Server) PostActiveflowsIdStop(c *gin.Context, id string) {
	tmp := flow_manager.Activeflow{
		Id: &id,
	}

	res := server.PostActiveflowsIdStop200JSONResponse(tmp)
	c.JSON(http.StatusOK, res)
}
