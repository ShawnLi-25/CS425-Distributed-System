package node

import (
	"log"
	"sort"
)

var membershipList []string

type Updater struct{}

type UpdateQuery struct {
	queryType int //0-GET, 1-ADD, 2-DELETE
	ID        string
}

//Open a go routine for this function, whenever needs update, build a channel; output will be 
func UpdateMembershipList() {

	for {
		
		select {
		case updateQuery:=<-upQryChan:
			if updateQuery.queryType == 0 {
				memListChan <- membershipList
			} else if updateQuery.queryType == 1 {
				membershipList = AddNewNode(updateQuery.ID, membershipList)
				SortMembershipList(membershipList)
				memListChan <- membershipList
			} else if updateQuery.queryType == 2 {
				membershipList = DeleteNode(updateQuery.ID, membershipList)
				memListChan <- membershipList
			}
		}
	}
}

func SortMembershipList(list []string) []string {
	sort.Strings(list)
	return list
}

func AddNewNode(newNodeID string, list []string) []string {
	log.Print("=== Current List is: ===")
	log.Print(list)
	newList := append(list, newNodeID)
	log.Print("=== Current List is: ===")
	log.Print(newList)
	return newList
}

func DeleteNode(nodeID string, list []string) []string {
	for i := 0; i < len(list); i++ {
		if list[i] == nodeID {

			if i != len(list)-1 {
				list = append(list[:i], list[i+1:]...)
			} else {
				list = list[:i]
			}
		}
	}
	return list
}

type MonitorList struct {
	localID   string
	monitorID []string
	Content   []string
}

func (*MonitorList) InitList() {
	var monitorList []string

}
