package groups

import (
	"application/core/model"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/rokwire/logging-library-go/v2/logs"
)

// Adapter groups adapter
type Adapter struct {
	baseURL        string
	internalAPIKey string

	//TODO: Switch to service account once Group BB has been updated
	// serviceAccountManager *authservice.ServiceAccountManager

	logger *logs.Logger
}

// NewGroupsAdapter creates a new Groups BB adapter instance
func NewGroupsAdapter(notificationHost string, internalAPIKey string, logger *logs.Logger) (*Adapter, error) {
	return &Adapter{baseURL: notificationHost, internalAPIKey: internalAPIKey, logger: logger}, nil
}

// GetGroupsMembership retrieves all groups that a user is a member
func (a *Adapter) GetGroupsMembership(userToken string) (*model.GroupMembership, error) {
	if len(userToken) == 0 {
		return nil, nil
	}

	url := fmt.Sprintf("%s/api/user/group-memberships", a.baseURL)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userToken))
	if err != nil {
		log.Printf("error GetGroupsMembership: request - %s", err)
		return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error GetGroupsMembership: request - %s", err)
		return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error GetGroupsMembership: request - %s", err)
			return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
		}

		log.Printf("error GetGroupsMembership: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
		return nil, fmt.Errorf("error GetGroupsMembership: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error GetGroupsMembership: request - %s", err)
		return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
	}

	var groups []model.UserGroup
	err = json.Unmarshal(data, &groups)
	if err != nil {
		log.Printf("error GetGroupsMembership: request - %s", err)
		return nil, fmt.Errorf("error GetGroupsMembership: request - %s", err)
	}

	membership := model.GroupMembership{}
	if len(groups) > 0 {
		for _, group := range groups {
			if group.MembershipStatus == "member" {
				membership.GroupIDsAsMember = append(membership.GroupIDsAsMember, group.ID)
			} else if group.MembershipStatus == "admin" {
				membership.GroupIDsAsAdmin = append(membership.GroupIDsAsAdmin, group.ID)
			}
		}
	}

	return &membership, nil
}

// GetGroupMembers retrieves all group members from a group
func (a *Adapter) GetGroupMembers(userToken string, groupID string) (*[]model.GroupMember, error) {
	if len(groupID) == 0 {
		return nil, nil
	}

	url := fmt.Sprintf("%s/api/group/%s/members", a.baseURL, groupID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userToken))
	if err != nil {
		log.Printf("error GetGroupMembers: request - %s", err)
		return nil, fmt.Errorf("error GetGroupMembers: request - %s", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error GetGroupMembers: request - %s", err)
		return nil, fmt.Errorf("error GetGroupMembers: request - %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error GetGroupMembers: request - %s", err)
			return nil, fmt.Errorf("error GetGroupMembers: request - %s", err)
		}

		log.Printf("error GetGroupMembers: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
		return nil, fmt.Errorf("error GetGroupMembers: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error GetGroupMembers: request - %s", err)
		return nil, fmt.Errorf("error GetGroupMembers: request - %s", err)
	}

	var groups []model.GroupMember
	err = json.Unmarshal(data, &groups)
	if err != nil {
		log.Printf("error GetGroupMembers: request - %s", err)
		return nil, fmt.Errorf("error GetGroupMembers: request - %s", err)
	}

	return &groups, nil
}

// GetGroupDetails retrieves group details
func (a *Adapter) GetGroupDetails(userToken string, groupID string) (*model.Group, error) {
	if len(groupID) == 0 {
		return nil, nil
	}

	url := fmt.Sprintf("%s/api/v2/groups/%s", a.baseURL, groupID)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Printf("error GetGroupDetails: request - %s", err)
		return nil, fmt.Errorf("error GetGroupDetails: request - %s", err)
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", userToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error GetGroupDetails: request - %s", err)
		return nil, fmt.Errorf("error GetGroupDetails: request - %s", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		errorBody, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("error GetGroupDetails: request - %s", err)
			return nil, fmt.Errorf("error GetGroupDetails: request - %s", err)
		}

		log.Printf("error GetGroupDetails: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
		return nil, fmt.Errorf("error GetGroupDetails: request - %d. Error: %s, Body: %s", resp.StatusCode, err, string(errorBody))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error GetGroupDetails: request - %s", err)
		return nil, fmt.Errorf("error GetGroupDetails: request - %s", err)
	}

	var group model.Group
	err = json.Unmarshal(data, &group)
	if err != nil {
		log.Printf("error GetGroupDetails: request - %s", err)
		return nil, fmt.Errorf("error GetGroupDetails: request - %s", err)
	}

	return &group, nil
}

// SendGroupNotification Sends a notification to members of a group
func (a *Adapter) SendGroupNotification(groupID string, notification model.GroupNotification) {
	go a.sendGroupNotification(groupID, notification)
}

// SendGroupNotification Sends a group notification
func (a *Adapter) sendGroupNotification(groupID string, notification model.GroupNotification) {
	if len(groupID) == 0 || len(notification.Subject) == 0 || len(notification.Body) == 0 {
		return
	}

	bodyBytes, err := json.Marshal(notification)
	if err != nil {
		log.Printf("error creating group notification request body - %s", err)
		return
	}

	url := fmt.Sprintf("%s/api/int/group/%s/notification", a.baseURL, groupID)
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	req.Header.Set("INTERNAL-API-KEY", a.internalAPIKey)
	if err != nil {
		log.Printf("error SendGroupNotification: request - %s", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Printf("error SendGroupNotification: request - %s", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("error SendGroupNotification: request - %d. Error: %s", resp.StatusCode, err)
	}
}
