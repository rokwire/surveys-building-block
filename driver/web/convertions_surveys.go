package web

import (
	"application/core/model"
	"time"

	"github.com/rokwire/core-auth-library-go/v3/tokenauth"
)

func surveyRequestToSurvey(claims *tokenauth.Claims, item model.SurveyRequest) model.Survey {
	item.Type = "user"
	//start
	var startValue *time.Time
	if item.StartDate != nil {
		startValueValue := time.Unix(int64(*item.StartDate), 0)
		startValue = &startValueValue
	}
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
		ResponseKeys: item.ResponseKeys, CalendarEventID: item.CalendarEventID, StartDate: startValue, EndDate: endValue,
		Public: item.Public, Archived: item.Archived}
}

func surveyToSurveyRequest(item model.Survey) model.SurveyRequest {
	var startDateUnixTimestamp int64
	if item.StartDate != nil {
		startDateUnixTimestamp = item.StartDate.Unix()
	}

	var endDateUnixTimestamp int64
	if item.EndDate != nil {
		endDateUnixTimestamp = item.EndDate.Unix()
	}

	return model.SurveyRequest{ID: item.ID, CreatorID: item.CreatorID, OrgID: item.OrgID, AppID: item.AppID, Title: item.Title,
		MoreInfo: item.MoreInfo, Data: item.Data, Scored: item.Scored, ResultRules: item.ResultRules, ResultJSON: item.ResultJSON,
		Type: item.Type, SurveyStats: item.SurveyStats, Sensitive: item.Sensitive, Anonymous: item.Anonymous, DefaultDataKey: item.DefaultDataKey,
		DefaultDataKeyRule: item.DefaultDataKeyRule, Constants: item.Constants, Strings: item.Strings, SubRules: item.SubRules,
		ResponseKeys: item.ResponseKeys, DateCreated: item.DateCreated, CalendarEventID: item.CalendarEventID, StartDate: &startDateUnixTimestamp,
		EndDate: &endDateUnixTimestamp, Public: item.Public, Archived: item.Archived}
}

func surveysToSurveyRequests(items []model.Survey) []model.SurveyRequest {
	list := make([]model.SurveyRequest, len(items))
	for index := range items {
		list[index] = surveyToSurveyRequest(items[index])
	}
	return list
}

func updateSurveyRequestToSurvey(claims *tokenauth.Claims, item model.SurveyRequest, id string) model.Survey {
	item.Type = "user"
	//start
	var startValue *time.Time
	if item.StartDate != nil {
		startValueTime := time.Unix(int64(*item.StartDate), 0)
		startValue = &startValueTime
	}
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
		ResponseKeys: item.ResponseKeys, CalendarEventID: item.CalendarEventID, StartDate: startValue, EndDate: endValue,
		Public: item.Public, Archived: item.Archived}
}

func surveyTimeFilter(item *model.SurveyTimeFilterRequest) *model.SurveyTimeFilter {

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
