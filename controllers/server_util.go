package controller

import (
	"github.com/gravitl/netmaker/logger"
	"github.com/gravitl/netmaker/logic"
	"github.com/gravitl/netmaker/models"
	"github.com/gravitl/netmaker/mq"
	"github.com/gravitl/netmaker/servercfg"
)

func runServerPeerUpdate(node *models.Node, ifaceDelta bool) error {

	err := logic.TimerCheckpoint()
	if err != nil {
		logger.Log(3, "error occurred on timer,", err.Error())
	}

	if err := mq.PublishPeerUpdate(node); err != nil {
		logger.Log(0, "failed to inform peers of new node ", err.Error())
	}

	if servercfg.IsClientMode() != "on" {
		return nil
	}
	var currentServerNode, getErr = logic.GetNetworkServerLeader(node.Network)
	if err != nil {
		return getErr
	}
	if err = logic.ServerUpdate(&currentServerNode, ifaceDelta); err != nil {
		logger.Log(1, "server node:", currentServerNode.ID, "failed update")
		return err
	}
	return nil
}
