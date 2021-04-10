package tools

import (
	"errors"
	"net/http"
	"net/url"
)

type Category string

const (
	CategoryDeploy      Category = "deployment"
	CategoryPullRequest          = "pull-request"
)

const (
	TrackingId = "UA-42300869-1"
)

type Action string

const (
	ActionDeploymentRollout  Action = "rollout"
	ActionDeploymentRollback        = "rollback"
	ActionDeploymentDelete          = "delete"
	ActionPrCreated                 = "created"
	ActionPrMerged                  = "merged"
)

type DataSource string

const (
	Github DataSource = "github"
	Gitlab            = "gitlab"
)

type analytics struct {
	trackingId string
}

func NewAnalytics() *analytics {
	return &analytics{
		trackingId: TrackingId,
	}
}

func (a analytics) TrackDeploy(dataSource DataSource, customer string, user string) error {
	return a.trackEvent(dataSource, customer, user, CategoryDeploy, ActionDeploymentRollout)
}

func (a analytics) TrackRollout(dataSource DataSource, customer string, user string) error {
	return a.trackEvent(dataSource, customer, user, CategoryDeploy, ActionDeploymentRollback)
}

func (a analytics) TrackPullRequest(dataSource DataSource, customer string, user string, prAction Action) error {
	return a.trackEvent(dataSource, customer, user, CategoryPullRequest, prAction)
}

func (a analytics) trackEvent(dataSource DataSource, customer, user string, category Category, action Action) error {
	if category == "" || action == "" {
		return errors.New("analytics: category and action are required")
	}

	v := url.Values{
		"v":   {"1"},
		"tid": {a.trackingId},
		"cid": {customer},
		"uid": {user},
		"ds":  {string(dataSource)},
		"dr":  {string(dataSource)},
		"t":   {"event"},
		"ec":  {string(category)},
		"ea":  {string(action)},
		"ev":  {"1"},
	}

	_, err := http.PostForm("https://www.google-analytics.com/collect", v)
	return err
}
