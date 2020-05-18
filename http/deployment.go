package http

import (
	"fmt"
	"net/http"
	"time"
)

func deployHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
	)
	defer httpErrorCommon(c, &httpCode, &err)
	data, err := c.GetRawData()
	if err != nil {
		return
	}
	deployment, err := svc.Deployer.Deploy(string(data))
	if err != nil {
		return
	}
	go svc.Deployer.WaitForPodContainersRunning(deployment.ClusterName, deployment.Namespace, deployment.Name,
		time.Second*100, time.Second*3, func(cluster, namespace, deploymentName, info string) {
			fmt.Println(info)
		})
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: deployment,
	})
}

func updateHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `json:"cluster" binding:"required"`
			Namespace      string `json:"namespace" binding:"required"`
			DeploymentName string `json:"deploymentName" binding:"required"`
			Image          string `json:"image" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	deployment, err := svc.Deployer.Update(req.Cluster, req.Namespace, req.DeploymentName, req.Image)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: deployment,
	})
}

func undoHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `json:"cluster" binding:"required"`
			Namespace      string `json:"namespace" binding:"required"`
			DeploymentName string `json:"deploymentName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	if err = svc.Deployer.Undo(req.Cluster, req.Namespace, req.DeploymentName); err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
	})
}

func rollBackHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `json:"cluster" binding:"required"`
			Namespace      string `json:"namespace" binding:"required"`
			DeploymentName string `json:"deploymentName" binding:"required"`
			Rs             string `json:"rs" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	deployment, err := svc.Deployer.RollBack(req.Cluster, req.Namespace, req.DeploymentName, req.Rs)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: deployment,
	})
}

func deploymentHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `form:"cluster" binding:"required"`
			Namespace      string `form:"namespace" binding:"required"`
			DeploymentName string `form:"deploymentName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	deployment, err := svc.Deployer.Deployment(req.Cluster, req.Namespace, req.DeploymentName)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: deployment,
	})
}

func deleteDeploymentHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `json:"cluster" binding:"required"`
			Namespace      string `json:"namespace" binding:"required"`
			DeploymentName string `json:"deploymentName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	if err = svc.Deployer.DeploymentDelete(req.Cluster, req.Namespace, req.DeploymentName); err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
	})
}

func replicaSetHandler(c *MyContext) {
	var (
		err      error
		httpCode = http.StatusInternalServerError
		req      struct {
			Cluster        string `form:"cluster" binding:"required"`
			Namespace      string `form:"namespace" binding:"required"`
			DeploymentName string `form:"deploymentName" binding:"required"`
		}
	)
	defer httpErrorCommon(c, &httpCode, &err)
	if err = c.Bind(&req); err != nil {
		httpCode = http.StatusBadRequest
		return
	}
	replicaSetList, err := svc.Deployer.ReplicaSetList(req.Cluster, req.Namespace, req.DeploymentName)
	if err != nil {
		return
	}
	c.JSON(http.StatusOK, &RespCommon{
		Code: HTTPErrorCodeSuccess,
		Data: replicaSetList,
	})
}
