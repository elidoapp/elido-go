package elido

import (
	"context"
	"encoding/json"
	"net/url"
	"strconv"
)

// SmsCampaigns — /v1/workspaces/{id}/sms-campaigns.
//
// Per-recipient short-link tracking + Twilio dispatch, with STOP/HELP
// keyword handling on the public webhook surface. The CRUD shape here
// mirrors the dashboard editor: create draft → list → schedule/pause/
// resume/cancel.

type SmsCampaigns struct{ c *Client }

type SmsCampaign struct {
	ID          int64   `json:"id"`
	WorkspaceID int64   `json:"workspace_id"`
	Name        string  `json:"name"`
	Body        string  `json:"body"`
	SenderID    string  `json:"sender_id"`
	Status      string  `json:"status"`
	SendAt      *string `json:"send_at,omitempty"`
	StartedAt   *string `json:"started_at,omitempty"`
	CompletedAt *string `json:"completed_at,omitempty"`
	TargetURL   string  `json:"target_url"`
	CreatedBy   *int64  `json:"created_by,omitempty"`
	CreatedAt   string  `json:"created_at"`
}

type SmsCampaignStats struct {
	Total         int64 `json:"total"`
	Queued        int64 `json:"queued"`
	Sent          int64 `json:"sent"`
	Delivered     int64 `json:"delivered"`
	Failed        int64 `json:"failed"`
	OptedOut      int64 `json:"opted_out"`
	Undeliverable int64 `json:"undeliverable"`
}

type SmsCampaignDetail struct {
	Campaign SmsCampaign      `json:"campaign"`
	Stats    SmsCampaignStats `json:"stats"`
}

type SmsMessage struct {
	ID            int64           `json:"id"`
	CampaignID    int64           `json:"campaign_id"`
	RecipientE164 string          `json:"recipient_e164"`
	MergeData     json.RawMessage `json:"merge_data"`
	LinkID        *int64          `json:"link_id,omitempty"`
	TwilioSID     *string         `json:"twilio_sid,omitempty"`
	Status        string          `json:"status"`
	SentAt        *string         `json:"sent_at,omitempty"`
	DeliveredAt   *string         `json:"delivered_at,omitempty"`
	FailedReason  *string         `json:"failed_reason,omitempty"`
	CostMicrosEUR *int32          `json:"cost_micros_eur,omitempty"`
}

type SmsRecipientInput struct {
	E164      string          `json:"e164"`
	MergeData json.RawMessage `json:"merge_data,omitempty"`
}

type CreateSmsCampaignInput struct {
	Name       string              `json:"name"`
	Body       string              `json:"body"`
	SenderID   string              `json:"sender_id"`
	TargetURL  string              `json:"target_url"`
	SendAt     *string             `json:"send_at,omitempty"`
	Recipients []SmsRecipientInput `json:"recipients"`
}

type UpdateSmsCampaignInput struct {
	Name      *string `json:"name,omitempty"`
	Body      *string `json:"body,omitempty"`
	SenderID  *string `json:"sender_id,omitempty"`
	TargetURL *string `json:"target_url,omitempty"`
	SendAt    *string `json:"send_at,omitempty"`
}

type SmsCampaignsListResponse struct {
	Items  []SmsCampaign `json:"items"`
	Limit  int           `json:"limit"`
	Offset int           `json:"offset"`
}

type SmsMessagesListResponse struct {
	Items  []SmsMessage `json:"items"`
	Limit  int          `json:"limit"`
	Offset int          `json:"offset"`
}

func (s *SmsCampaigns) List(ctx context.Context, workspaceID int64, limit, offset int) (*SmsCampaignsListResponse, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	out := &SmsCampaignsListResponse{}
	return out, s.c.request(ctx, s.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns", nil, out, q)
}

func (s *SmsCampaigns) Create(ctx context.Context, workspaceID int64, in CreateSmsCampaignInput) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns", in, out, nil)
}

func (s *SmsCampaigns) Get(ctx context.Context, workspaceID, campaignID int64) (*SmsCampaignDetail, error) {
	out := &SmsCampaignDetail{}
	return out, s.c.request(ctx, s.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID), nil, out, nil)
}

func (s *SmsCampaigns) ListMessages(ctx context.Context, workspaceID, campaignID int64, limit, offset int) (*SmsMessagesListResponse, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	out := &SmsMessagesListResponse{}
	return out, s.c.request(ctx, s.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/messages", nil, out, q)
}

func (s *SmsCampaigns) Update(ctx context.Context, workspaceID, campaignID int64, in UpdateSmsCampaignInput) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "PATCH",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID), in, out, nil)
}

type scheduleSmsCampaignBody struct {
	SendAt *string `json:"send_at,omitempty"`
}

func (s *SmsCampaigns) Schedule(ctx context.Context, workspaceID, campaignID int64, sendAt *string) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/schedule",
		scheduleSmsCampaignBody{SendAt: sendAt}, out, nil)
}

func (s *SmsCampaigns) Pause(ctx context.Context, workspaceID, campaignID int64) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/pause", nil, out, nil)
}

func (s *SmsCampaigns) Resume(ctx context.Context, workspaceID, campaignID int64) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/resume", nil, out, nil)
}

func (s *SmsCampaigns) Cancel(ctx context.Context, workspaceID, campaignID int64) (*SmsCampaign, error) {
	out := &SmsCampaign{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/cancel", nil, out, nil)
}

type testSmsCampaignBody struct {
	ToE164 string `json:"to_e164"`
}

type SmsTestResponse struct {
	OK        bool   `json:"ok"`
	TwilioSID string `json:"twilio_sid"`
}

func (s *SmsCampaigns) Test(ctx context.Context, workspaceID, campaignID int64, toE164 string) (*SmsTestResponse, error) {
	out := &SmsTestResponse{}
	return out, s.c.request(ctx, s.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-campaigns/"+itoa(campaignID)+"/test",
		testSmsCampaignBody{ToE164: toE164}, out, nil)
}

// --- Opt-outs --------------------------------------------------------------

type SmsOptOuts struct{ c *Client }

type SmsOptOut struct {
	ID            int64  `json:"id"`
	WorkspaceID   int64  `json:"workspace_id"`
	RecipientE164 string `json:"recipient_e164"`
	Reason        string `json:"reason"`
	OptedOutAt    string `json:"opted_out_at"`
}

type SmsOptOutsListResponse struct {
	Items  []SmsOptOut `json:"items"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

type AddSmsOptOutInput struct {
	RecipientE164 string `json:"recipient_e164"`
	Reason        string `json:"reason,omitempty"`
}

func (o *SmsOptOuts) List(ctx context.Context, workspaceID int64, limit, offset int) (*SmsOptOutsListResponse, error) {
	q := url.Values{}
	if limit > 0 {
		q.Set("limit", strconv.Itoa(limit))
	}
	if offset > 0 {
		q.Set("offset", strconv.Itoa(offset))
	}
	out := &SmsOptOutsListResponse{}
	return out, o.c.request(ctx, o.c.base, "GET",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-opt-outs", nil, out, q)
}

func (o *SmsOptOuts) Add(ctx context.Context, workspaceID int64, in AddSmsOptOutInput) (*SmsOptOut, error) {
	out := &SmsOptOut{}
	return out, o.c.request(ctx, o.c.base, "POST",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-opt-outs", in, out, nil)
}

func (o *SmsOptOuts) Remove(ctx context.Context, workspaceID int64, e164 string) error {
	return o.c.request(ctx, o.c.base, "DELETE",
		"/v1/workspaces/"+itoa(workspaceID)+"/sms-opt-outs/"+url.PathEscape(e164),
		nil, nil, nil)
}
