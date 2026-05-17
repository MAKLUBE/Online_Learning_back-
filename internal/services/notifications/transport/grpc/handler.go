package grpc

import (
	"context"
	"online-learning-platform/internal/shared/grpcerrors"
	"online-learning-platform/internal/shared/validate"
	"time"

	"online-learning-platform/internal/services/notifications/model"
	"online-learning-platform/internal/services/notifications/usecase"
	notificationpb "online-learning-platform/pkg/gen/notifications"
)

type Handler struct {
	notificationpb.UnimplementedNotificationServiceServer
	useCase *usecase.UseCase
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) SendEmail(ctx context.Context, req *notificationpb.SendEmailRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetSubject()) || !validate.Required(req.GetBody()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email, subject and body are required")
	}
	notification, err := h.useCase.SendEmail(ctx, req.GetUserId(), req.GetEmail(), req.GetSubject(), req.GetBody(), "direct")
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) SendWelcomeEmail(ctx context.Context, req *notificationpb.SendWelcomeEmailRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetName()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email and name are required")
	}
	notification, err := h.useCase.SendWelcomeEmail(ctx, req.GetUserId(), req.GetEmail(), req.GetName())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) SendPasswordResetEmail(ctx context.Context, req *notificationpb.SendPasswordResetEmailRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetResetToken()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email and reset_token are required")
	}
	notification, err := h.useCase.SendPasswordResetEmail(ctx, req.GetUserId(), req.GetEmail(), req.GetResetToken())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) SendCourseEnrollmentEmail(ctx context.Context, req *notificationpb.SendCourseEnrollmentEmailRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetCourseId()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email and course_id are required")
	}
	notification, err := h.useCase.SendCourseEnrollmentEmail(ctx, req.GetUserId(), req.GetEmail(), req.GetCourseId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) SendAssignmentGradedEmail(ctx context.Context, req *notificationpb.SendAssignmentGradedEmailRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetAssignmentId()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email and assignment_id are required")
	}
	notification, err := h.useCase.SendAssignmentGradedEmail(ctx, req.GetUserId(), req.GetEmail(), req.GetAssignmentId(), req.GetGrade())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) SendDeadlineReminder(ctx context.Context, req *notificationpb.SendDeadlineReminderRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetAssignmentId()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email and assignment_id are required")
	}
	notification, err := h.useCase.SendDeadlineReminder(ctx, req.GetUserId(), req.GetEmail(), req.GetAssignmentId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func (h *Handler) GetNotificationHistory(ctx context.Context, req *notificationpb.GetNotificationHistoryRequest) (*notificationpb.NotificationHistoryResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	notifications, err := h.useCase.GetNotificationHistory(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	items := make([]*notificationpb.NotificationResponse, 0, len(notifications))
	for _, notification := range notifications {
		items = append(items, toNotification(notification))
	}
	return &notificationpb.NotificationHistoryResponse{Notifications: items}, nil
}

func (h *Handler) GetNotificationById(ctx context.Context, req *notificationpb.GetNotificationByIdRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	notification, err := h.useCase.GetNotificationByID(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("notification not found")
	}
	return toNotification(notification), nil
}

func (h *Handler) MarkNotificationRead(ctx context.Context, req *notificationpb.MarkNotificationReadRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	notification, err := h.useCase.MarkNotificationRead(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("notification not found")
	}
	return toNotification(notification), nil
}

func (h *Handler) GetUnreadCount(ctx context.Context, req *notificationpb.GetUnreadCountRequest) (*notificationpb.UnreadCountResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	count, err := h.useCase.GetUnreadCount(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return &notificationpb.UnreadCountResponse{Count: int32(count)}, nil
}

func (h *Handler) UpdateNotificationSettings(ctx context.Context, req *notificationpb.UpdateNotificationSettingsRequest) (*notificationpb.NotificationSettingsResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	settings, err := h.useCase.UpdateNotificationSettings(ctx, req.GetUserId(), req.GetEmailEnabled())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return &notificationpb.NotificationSettingsResponse{UserId: settings.UserID, EmailEnabled: settings.EmailEnabled}, nil
}

func (h *Handler) DeleteNotification(ctx context.Context, req *notificationpb.DeleteNotificationRequest) (*notificationpb.Empty, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	return &notificationpb.Empty{}, grpcerrors.Internal(h.useCase.DeleteNotification(ctx, req.GetId()))
}

func (h *Handler) ScheduleNotification(ctx context.Context, req *notificationpb.ScheduleNotificationRequest) (*notificationpb.NotificationResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Email(req.GetEmail()) || !validate.Required(req.GetSubject()) || !validate.Required(req.GetBody()) || !validate.Required(req.GetSendAt()) {
		return nil, grpcerrors.InvalidArgument("user_id, valid email, subject, body and send_at are required")
	}
	sendAt, err := time.Parse(time.RFC3339, req.GetSendAt())
	if err != nil {
		return nil, grpcerrors.InvalidArgument("send_at must be RFC3339")
	}
	notification, err := h.useCase.ScheduleNotification(ctx, req.GetUserId(), req.GetEmail(), req.GetSubject(), req.GetBody(), sendAt)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toNotification(notification), nil
}

func toNotification(notification model.Notification) *notificationpb.NotificationResponse {
	return &notificationpb.NotificationResponse{
		Id: notification.ID, UserId: notification.UserID, Email: notification.RecipientEmail,
		Subject: notification.Subject, Body: notification.Body, Type: notification.Type,
		Status: notification.Status, Read: notification.Read,
	}
}
