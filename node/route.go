package node

import (
	"strings"

	"github.com/Dliv3/Venom/utils"
)

// NetworkTopology 网络拓扑
// RouteTable 路由表, 路由表的Key为目标节点, Value为下一跳节点，注意在该多级代理的应用场景中，暂不支持节点间形成环路的情况
// 从管理节点到其他节点，有且仅有一条道路，所以不涉及路由选路的问题，即仅支持树形拓扑
// NetworkMap 网络拓扑, key为节点id，value为该节点下直接连接的节点id
type NetworkTopology struct {
	RouteTable map[string]string
	NetworkMap map[string]([]string)
}

func (nt *NetworkTopology) recursiveUpdateRouteTable(root string, key string) {
	if value, ok := nt.NetworkMap[key]; ok {
		for _, v := range value {
			// avoid adding the current node to the route table
			if v != CurrentNode.HashID {
				nt.RouteTable[v] = root
				nt.recursiveUpdateRouteTable(v, v)
			}
		}
	}
}

// UpdateRouteTable 通过NeworkMap中的数据生成路由表
func (nt *NetworkTopology) UpdateRouteTable() {
	// 清空现有路由表，这样当有节点断开时网络拓扑会实时改变
	nt.RouteTable = make(map[string]string)

	if value, ok := nt.NetworkMap[CurrentNode.HashID]; ok {
		for _, v := range value {
			nt.RouteTable[v] = v
			nt.recursiveUpdateRouteTable(v, v)
		}
	}
}

// AddRoute 在路由表中添加一条路由表
func (nt *NetworkTopology) AddRoute(targetNode string, nextNode string) {
	nt.RouteTable[targetNode] = nextNode
}

// AddNetworkMap 向网络拓扑中添加节点，key为父节点，nodeId为子节点
func (nt *NetworkTopology) AddNetworkMap(parent string, chlid string) {
	nt.NetworkMap[parent] = append(nt.NetworkMap[parent], chlid)
}

// InitNetworkMap 初始化网络拓扑, 初始网络拓扑仅包含与本节点直接相连的节点
func (nt *NetworkTopology) InitNetworkMap() {
	nt.NetworkMap = make(map[string]([]string))
	for i := range Nodes {
		if Nodes[i].DirectConnection {
			nt.AddNetworkMap(CurrentNode.HashID, Nodes[i].HashID)
		}
	}
}

// ResolveNetworkMapData 解析SyncPacket中包含的NetworkMap数据
func (nt *NetworkTopology) ResolveNetworkMapData(data []byte) {
	var networkMap = strings.Split(string(data), "$")
	for i := range networkMap {
		each := strings.Split(networkMap[i], "#")
		key := each[0]
		if _, ok := nt.NetworkMap[key]; ok {
			nt.NetworkMap[key] = append(nt.NetworkMap[key], strings.Split(each[1], "|")...)
		} else {
			tempNodes := strings.Split(each[1], "|")
			var nodes []string
			for j := range tempNodes {
				// 删除当前节点
				if tempNodes[j] == CurrentNode.HashID {
					nodes = append(tempNodes[:j], tempNodes[j+1:]...)
					break
				}
			}
			nt.NetworkMap[key] = nodes
		}
		nt.NetworkMap[key] = utils.RemoveDuplicateElement(nt.NetworkMap[key])
	}
}

// GenerateNetworkMapData 更具网络拓扑生成SyncPacket中使用的NetworkMap数据
func (nt *NetworkTopology) GenerateNetworkMapData() []byte {
	// fmt.Println(CurrentNode.HashID)
	var networkMap []string

	for key := range nt.NetworkMap {
		networkMap = append(networkMap, key+"#"+strings.Join(nt.NetworkMap[key], "|"))
	}
	var networkMapBytes = []byte(strings.Join(networkMap, "$"))
	return networkMapBytes
}

func (nt *NetworkTopology) DeleteNode(node *Node) {
	// nodeNumber := nt.NodeUUID2Number[node.HashID]

	// // 删除节点标号
	// delete(nt.NodeNumber2UUID, nodeNumber)
	// delete(nt.NodeUUID2Number, node.HashID)
	// // 删除路由表
	// delete(nt.RouteTable, node.HashID)
	// // 删除网络拓扑
	// delete(nt.NetworkMap, node.HashID)

	// for key, _ := range nt.NetworkMap {
	// 	for i, value := range nt.NetworkMap[key] {
	// 		if value == node.HashID {
	// 			nt.NetworkMap[key] = append(nt.NetworkMap[key][:i], nt.NetworkMap[key][i+1:]...)
	// 		}
	// 	}
	// }
	// // 删除节点描述
	// delete(nt.NodeDescription, nodeNumber)
}
