package http

import "net/http"

func ingressHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster   string `form:"cluster" binding:"required"`
			Namespace string `form:"namespace" binding:"required"`
			Name      string `form:"name" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	ingress, err := svc.Deployer.Ingress(req.Cluster, req.Namespace, req.Name)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: ingress,
	})
}
