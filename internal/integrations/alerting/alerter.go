package alerting

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/go-multierror"

	"github.com/hatchet-dev/hatchet/internal/encryption"
	"github.com/hatchet-dev/hatchet/internal/integrations/alerting/alerttypes"
	"github.com/hatchet-dev/hatchet/internal/integrations/email"
	"github.com/hatchet-dev/hatchet/internal/repository"
	"github.com/hatchet-dev/hatchet/internal/repository/prisma/db"
	"github.com/hatchet-dev/hatchet/internal/repository/prisma/sqlchelpers"

	"github.com/hatchet-dev/timediff"
)

type TenantAlertManager struct {
	repo      repository.EngineRepository
	enc       encryption.EncryptionService
	serverURL string
	email     email.EmailService
}

func New(repo repository.EngineRepository, e encryption.EncryptionService, serverURL string, email email.EmailService) *TenantAlertManager {
	return &TenantAlertManager{repo, e, serverURL, email}
}

func (t *TenantAlertManager) HandleAlert(tenantId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// read in the tenant alerting settings and determine if we should alert
	tenantAlerting, err := t.repo.TenantAlertingSettings().GetTenantAlertingSettings(ctx, tenantId)

	if err != nil {
		return err
	}

	lastAlertedAt := tenantAlerting.Settings.LastAlertedAt.Time.UTC()
	maxFrequency, err := time.ParseDuration(tenantAlerting.Settings.MaxFrequency)

	if err != nil {
		return err
	}

	isZero := lastAlertedAt.IsZero()

	if isZero || time.Since(lastAlertedAt) > maxFrequency {
		// update the lastAlertedAt
		now := time.Now().UTC()

		// if we're in the zero state, we don't want to alert since the very beginning of the interval
		if isZero {
			lastAlertedAt = now.Add(-1 * maxFrequency)
		}

		err = t.repo.TenantAlertingSettings().UpdateTenantAlertingSettings(ctx, tenantId, &repository.UpdateTenantAlertingSettingsOpts{
			LastAlertedAt: &now,
		})

		if err != nil {
			return err
		}

		return t.sendAlert(ctx, tenantAlerting, lastAlertedAt)
	}

	return nil
}

func (t *TenantAlertManager) SendAlert(tenantId string, prevLastAlertedAt time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// read in the tenant alerting settings and determine if we should alert
	tenantAlerting, err := t.repo.TenantAlertingSettings().GetTenantAlertingSettings(ctx, tenantId)

	if err != nil {
		return err
	}

	return t.sendAlert(ctx, tenantAlerting, prevLastAlertedAt)
}

func (t *TenantAlertManager) sendAlert(ctx context.Context, tenantAlerting *repository.GetTenantAlertingSettingsResponse, prevLastAlertedAt time.Time) error {
	// read in all failed workflow runs since the last alerted time, ordered by the most recent runs first
	statuses := []db.WorkflowRunStatus{
		db.WorkflowRunStatusFailed,
	}

	limit := 5

	tenantId := sqlchelpers.UUIDToStr(tenantAlerting.Settings.TenantId)

	failedWorkflowRuns, err := t.repo.WorkflowRun().ListWorkflowRuns(
		ctx,
		tenantId,
		&repository.ListWorkflowRunsOpts{
			Statuses:       &statuses,
			Limit:          &limit,
			OrderBy:        repository.StringPtr("createdAt"),
			OrderDirection: repository.StringPtr("DESC"),
			FinishedAfter:  &prevLastAlertedAt,
		},
	)

	if err != nil {
		return err
	}

	if failedWorkflowRuns.Count == 0 {
		return nil
	}

	failedItems := t.getFailedItems(failedWorkflowRuns)

	// iterate through possible alerters
	for _, slackWebhook := range tenantAlerting.SlackWebhooks {
		if innerErr := t.sendSlackAlert(slackWebhook, failedWorkflowRuns.Count, failedItems); innerErr != nil {
			err = multierror.Append(err, innerErr)
		}
	}

	for _, emailGroup := range tenantAlerting.EmailGroups {
		if innerErr := t.sendEmailAlert(tenantAlerting.Tenant, emailGroup, failedWorkflowRuns.Count, failedItems); innerErr != nil {
			err = multierror.Append(err, innerErr)
		}
	}

	return nil
}

func (t *TenantAlertManager) getFailedItems(failedWorkflowRuns *repository.ListWorkflowRunsResult) []alerttypes.WorkflowRunFailedItem {
	res := make([]alerttypes.WorkflowRunFailedItem, 0)

	for _, workflowRun := range failedWorkflowRuns.Rows {
		workflowRunId := sqlchelpers.UUIDToStr(workflowRun.WorkflowRun.ID)
		tenantId := sqlchelpers.UUIDToStr(workflowRun.WorkflowRun.TenantId)

		readableId := workflowRun.WorkflowRun.DisplayName.String

		if readableId == "" {
			readableId = workflowRun.Workflow.Name
		}

		res = append(res, alerttypes.WorkflowRunFailedItem{
			Link:                  fmt.Sprintf("%s/workflow-runs/%s?tenant=%s", t.serverURL, workflowRunId, tenantId),
			WorkflowName:          workflowRun.Workflow.Name,
			WorkflowRunReadableId: readableId,
			RelativeDate:          timediff.TimeDiff(workflowRun.WorkflowRun.FinishedAt.Time),
			AbsoluteDate:          workflowRun.WorkflowRun.FinishedAt.Time.Format("2006-01-02 15:04:05"),
		})
	}

	return res
}
