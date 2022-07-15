package models

import (
	"errors"

	"github.com/spf13/viper"
)

type NodeData struct {
	Name     string `form:"name" json:"name" xml:"name" binding:"required"`
	Uid      string `form:"uid" json:"uid" xml:"uid" binding:"required"`
	Type     string `form:"type" json:"type" xml:"type" binding:"required"`
	Host     string `form:"host" json:"host" xml:"host" binding:"required"`
	Location string `form:"location" json:"location" xml:"location" binding:"required"`
	Custom   string `form:"custom" json:"custom" xml:"custom" `
}

func SetNode(nodeData *NodeData) error {
	name := viper.GetString("node.name")
	nodeType := viper.GetString("node.type")
	host := viper.GetString("node.host")
	location := viper.GetString("node.location")
	custom := viper.GetString("node.custom")
	disabled := viper.GetBool("node.disabled")
	if disabled == true {
		error := errors.New("Node is disabled")
		return error
	}
	if name == "" || host == "" {
		error := errors.New("Lost Node Info")
		return error
	}
	nodeData.Name = name
	nodeData.Type = nodeType
	nodeData.Host = host
	nodeData.Location = location
	nodeData.Custom = custom
	return nil
}
