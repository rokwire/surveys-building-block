package web

import (
	"application/core/model"
	"time"

	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
)

func surveyRequestToSurvey(claims *tokenauth.Claims, item model.SurveyRequest) model.Survey {
	item.Type = "user"
	//start
	startValueValue := time.Unix(int64(item.StartDate), 0)
	//end
	var endValue *time.Time
	if item.EndDate != nil {
		endValueTime := time.Unix(int64(*item.EndDate), 0)
		endValue = &endValueTime
	}

	return model.Survey{CreatorID: claims.Subject, OrgID: claims.OrgID, AppID: claims.AppID, Type: item.Type, Title: item.Title,
		MoreInfo: item.MoreInfo, Data: item.Data, Scored: item.Scored, ResultRules: item.ResultRules, ResultJSON: item.ResultJSON,
		SurveyStats: item.SurveyStats, Sensitive: item.Sensitive, Anonymous: item.Anonymous, DefaultDataKey: item.DefaultDataKey,
		DefaultDataKeyRule: item.DefaultDataKeyRule, Constants: item.Constants, Strings: item.Strings, SubRules: item.SubRules,
		ResponseKeys: item.ResponseKeys, CalendarEventID: item.CalendarEventID, StartDate: startValueValue, EndDate: endValue}
}

func surveyToSurveyRequest(item model.Survey) model.SurveyRequest {
	startDateUnixTimestamp := item.StartDate.Unix()
	var endDateUnixTimestamp int64
	if item.EndDate != nil {
		endDateUnixTimestamp = item.EndDate.Unix()
	}

	return model.SurveyRequest{ID: item.ID, CreatorID: item.CreatorID, OrgID: item.OrgID, AppID: item.AppID, Title: item.Title,
		MoreInfo: item.MoreInfo, Data: item.Data, Scored: item.Scored, ResultRules: item.ResultRules, ResultJSON: item.ResultJSON,
		Type: item.Type, SurveyStats: item.SurveyStats, Sensitive: item.Sensitive, Anonymous: item.Anonymous, DefaultDataKey: item.DefaultDataKey,
		DefaultDataKeyRule: item.DefaultDataKeyRule, Constants: item.Constants, Strings: item.Strings, SubRules: item.SubRules,
		ResponseKeys: item.ResponseKeys, DateCreated: item.DateCreated, CalendarEventID: item.CalendarEventID, StartDate: startDateUnixTimestamp,
		EndDate: &endDateUnixTimestamp}
}

func updateSurveyRequestToSurvey(claims *tokenauth.Claims, item model.SurveyRequest, id string) model.Survey {
	item.Type = "user"
	//start
	startValueValue := time.Unix(int64(item.StartDate), 0)
	//end
	var endValue *time.Time
	if item.EndDate != nil {
		endValueTime := time.Unix(int64(*item.EndDate), 0)
		endValue = &endValueTime
	}

	return model.Survey{ID: id, CreatorID: claims.Subject, OrgID: claims.OrgID, AppID: claims.AppID, Type: item.Type, Title: item.Title,
		MoreInfo: item.MoreInfo, Data: item.Data, Scored: item.Scored, ResultRules: item.ResultRules, ResultJSON: item.ResultJSON,
		SurveyStats: item.SurveyStats, Sensitive: item.Sensitive, Anonymous: item.Anonymous, DefaultDataKey: item.DefaultDataKey,
		DefaultDataKeyRule: item.DefaultDataKeyRule, Constants: item.Constants, Strings: item.Strings, SubRules: item.SubRules,
		ResponseKeys: item.ResponseKeys, CalendarEventID: item.CalendarEventID, StartDate: startValueValue, EndDate: endValue}
}

func surveyTimeFilter(item model.SurveyTimeFilterRequest) *model.SurveyTimeFilter {

	filter := model.SurveyTimeFilter{}

	if item.StartTimeBefore != nil {
		beforeStartTime := time.Unix(*item.StartTimeBefore, 0)
		filter.StartTimeBefore = &beforeStartTime
	}
	if item.StartTimeAfter != nil {
		afterStartTime := time.Unix(*item.StartTimeAfter, 0)
		filter.StartTimeAfter = &afterStartTime
	}

	if item.EndTimeBefore != nil {
		beforeEndTime := time.Unix(*item.EndTimeBefore, 0)
		filter.EndTimeBefore = &beforeEndTime
	}
	if item.EndTimeAfter != nil {
		afterEndTime := time.Unix(*item.EndTimeAfter, 0)
		filter.EndTimeAfter = &afterEndTime
	}

	return &model.SurveyTimeFilter{
		StartTimeAfter:  filter.StartTimeAfter,
		StartTimeBefore: filter.StartTimeBefore,
		EndTimeAfter:    filter.EndTimeAfter,
		EndTimeBefore:   filter.EndTimeBefore}
}
