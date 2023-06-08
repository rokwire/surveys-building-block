package model

import (
	"time"
)

// Group struct wrapper
type Group struct {
	ID              string `json:"id"`
	ClientID        string `json:"client_id"`
	Category        string `json:"category"`
	Title           string `json:"title"`
	Privacy         string `json:"privacy"`
	HiddenForSearch bool   `json:"hidden_for_search"`
	Description     string `json:"description"`
	ImageURL        string `json:"image_url"`
	WebURL          string `json:"web_url"`
	Tags            string `json:"tags"`
	CreatorID       string `json:"creator_id"`
	CurrentMember   *struct {
		ID                       string `json:"id"`
		ClientID                 string `json:"client_id"`
		GroupID                  string `json:"group_id"`
		UserID                   string `json:"user_id"`
		ExternalID               string `json:"external_id"`
		Name                     string `json:"name"`
		NetID                    string `json:"net_id"`
		Email                    string `json:"email"`
		PhotoURL                 string `json:"photo_url"`
		Status                   string `json:"status"`
		Admin                    bool   `json:"admin"`
		RejectReason             string `json:"reject_reason"`
		NotificationsPreferences struct {
			OverridePreferences bool `json:"override_preferences"`
			AllMute             bool `json:"all_mute"`
			InvitationsMute     bool `json:"invitations_mute"`
			PostsMute           bool `json:"posts_mute"`
			EventsMute          bool `json:"events_mute"`
			PollsMute           bool `json:"polls_mute"`
		} `json:"notifications_preferences"`
		DateCreated  time.Time  `json:"date_created"`
		DateUpdated  *time.Time `json:"date_updated"`
		DateAttended *time.Time `json:"date_attended"`
	} `json:"current_member"`
	Stats struct {
		TotalCount      int `json:"total_count"`
		AdminsCount     int `json:"admins_count"`
		MemberCount     int `json:"member_count"`
		PendingCount    int `json:"pending_count"`
		RejectedCount   int `json:"rejected_count"`
		AttendanceCount int `json:"attendance_count"`
	} `json:"stats"`
	DateCreated                time.Time  `json:"date_created"`
	DateUpdated                *time.Time `json:"date_updated"`
	AuthmanEnabled             bool       `json:"authman_enabled"`
	AuthmanGroup               string     `json:"authman_group"`
	OnlyAdminsCanCreatePolls   bool       `json:"only_admins_can_create_polls"`
	CanJoinAutomatically       bool       `json:"can_join_automatically"`
	BlockNewMembershipRequests bool       `json:"block_new_membership_requests"`
	AttendanceGroup            bool       `json:"attendance_group"`
}

// IsCurrentUserAdmin checks if the user is a group admin
func (g *Group) IsCurrentUserAdmin(currentUserID string) bool {
	if g.CurrentMember != nil {
		if g.CurrentMember.UserID == currentUserID && g.CurrentMember.Status == "admin" {
			return true
		}
	}
	return false
}

type UserGroup struct {
	ID               string `json:"id"`
	Title            string `json:"title"`
	Privacy          string `json:"privacy"`
	MembershipStatus string `json:"membership_status"`
}

// GroupMembership mapping. Better to access map entry by key instead of iterating for check purpose.
type GroupMembership struct {
	GroupIDsAsAdmin  []string
	GroupIDsAsMember []string
}

// GroupMember represents the membership of a user to a given group
type GroupMember struct {
	ID         string `json:"id" bson:"_id"`
	ClientID   string `json:"client_id" bson:"client_id"`
	GroupID    string `json:"group_id" bson:"group_id"`
	UserID     string `json:"user_id" bson:"user_id"`
	ExternalID string `json:"external_id" bson:"external_id"`
	Name       string `json:"name" bson:"name"`
	NetID      string `json:"net_id" bson:"net_id"`
	Email      string `json:"email" bson:"email"`
	PhotoURL   string `json:"photo_url" bson:"photo_url"`

	Status string `json:"status" bson:"status"` //admin, pending, member, rejected

	RejectReason string `json:"reject_reason" bson:"reject_reason"`
	// MemberAnswers []MemberAnswer `json:"member_answers" bson:"member_answers"`
	SyncID string `json:"sync_id" bson:"sync_id"` //ID of sync that last updated this membership

	NotificationsPreferences NotificationsPreferences `json:"notifications_preferences" bson:"notifications_preferences"`

	DateCreated  time.Time  `json:"date_created" bson:"date_created"`
	DateUpdated  *time.Time `json:"date_updated" bson:"date_updated"`
	DateAttended *time.Time `json:"date_attended" bson:"date_attended"`
} //@name GroupMembership

type NotificationsPreferences struct {
	OverridePreferences bool `json:"override_preferences" bson:"override_preferences"`
	AllMute             bool `json:"all_mute" bson:"all_mute"`
	InvitationsMuted    bool `json:"invitations_mute" bson:"invitations_mute"`
	PostsMuted          bool `json:"posts_mute" bson:"posts_mute"`
	EventsMuted         bool `json:"events_mute" bson:"events_mute"`
	PollsMuted          bool `json:"polls_mute" bson:"polls_mute"`
} // @name NotificationsPreferences
