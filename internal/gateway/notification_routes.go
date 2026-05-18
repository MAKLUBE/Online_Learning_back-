package gateway

import (
	"net/http"

	notificationpb "online-learning-platform/pkg/gen/notifications"
)

func registerNotificationRoutes(mux *http.ServeMux, clients *Clients) {
	mux.HandleFunc("/api/notifications/email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendEmailRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendEmail(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/welcome-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendWelcomeEmailRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendWelcomeEmail(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/password-reset-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendPasswordResetEmailRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendPasswordResetEmail(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/course-enrollment-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendCourseEnrollmentEmailRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendCourseEnrollmentEmail(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/assignment-graded-email", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendAssignmentGradedEmailRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendAssignmentGradedEmail(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/deadline-reminder", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.SendDeadlineReminderRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.SendDeadlineReminder(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/schedule", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		var req notificationpb.ScheduleNotificationRequest
		if !decode(w, r, &req) {
			return
		}
		resp, err := clients.Notifications.ScheduleNotification(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/users/", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/notifications/users/")
		if len(parts) != 2 || r.Method != http.MethodGet {
			http.NotFound(w, r)
			return
		}
		userID := parts[0]
		switch parts[1] {
		case "history":
			resp, err := clients.Notifications.GetNotificationHistory(requestContext(r), &notificationpb.GetNotificationHistoryRequest{UserId: userID})
			writeGRPC(w, resp, err)
		case "unread-count":
			resp, err := clients.Notifications.GetUnreadCount(requestContext(r), &notificationpb.GetUnreadCountRequest{UserId: userID})
			writeGRPC(w, resp, err)
		default:
			http.NotFound(w, r)
		}
	})

	mux.HandleFunc("/api/notifications/settings/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPatch {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/notifications/settings/")
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		var req notificationpb.UpdateNotificationSettingsRequest
		if !decode(w, r, &req) {
			return
		}
		req.UserId = parts[0]
		resp, err := clients.Notifications.UpdateNotificationSettings(requestContext(r), &req)
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/notifications/", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/notifications/")
		if len(parts) == 0 {
			http.NotFound(w, r)
			return
		}
		notificationID := parts[0]
		if len(parts) == 2 && parts[1] == "read" {
			if r.Method != http.MethodPatch {
				methodAllowed(w)
				return
			}
			resp, err := clients.Notifications.MarkNotificationRead(requestContext(r), &notificationpb.MarkNotificationReadRequest{Id: notificationID})
			writeGRPC(w, resp, err)
			return
		}
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			resp, err := clients.Notifications.GetNotificationById(requestContext(r), &notificationpb.GetNotificationByIdRequest{Id: notificationID})
			writeGRPC(w, resp, err)
		case http.MethodDelete:
			resp, err := clients.Notifications.DeleteNotification(requestContext(r), &notificationpb.DeleteNotificationRequest{Id: notificationID})
			writeGRPC(w, resp, err)
		default:
			methodAllowed(w)
		}
	})
}
