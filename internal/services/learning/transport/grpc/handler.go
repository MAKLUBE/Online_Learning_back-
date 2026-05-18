package grpc

import (
	"context"
	"online-learning-platform/internal/shared/authz"
	"online-learning-platform/internal/shared/grpcerrors"
	"online-learning-platform/internal/shared/validate"
	"time"

	"online-learning-platform/internal/services/learning/model"
	"online-learning-platform/internal/services/learning/usecase"
	learningpb "online-learning-platform/pkg/gen/learning"
)

type Handler struct {
	learningpb.UnimplementedLearningServiceServer
	useCase *usecase.UseCase
}

func NewHandler(useCase *usecase.UseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) EnrollCourse(ctx context.Context, req *learningpb.EnrollCourseRequest) (*learningpb.EnrollmentResponse, error) {
	if err := authz.RequireRole(ctx, "student", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetUserId()) || !validate.Required(req.GetCourseId()) {
		return nil, grpcerrors.InvalidArgument("user_id and course_id are required")
	}
	enrollment, err := h.useCase.EnrollCourse(ctx, req.GetUserId(), req.GetCourseId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toEnrollment(enrollment), nil
}

func (h *Handler) UnenrollCourse(ctx context.Context, req *learningpb.UnenrollCourseRequest) (*learningpb.Empty, error) {
	if err := authz.RequireRole(ctx, "student", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetUserId()) || !validate.Required(req.GetCourseId()) {
		return nil, grpcerrors.InvalidArgument("user_id and course_id are required")
	}
	return &learningpb.Empty{}, grpcerrors.Internal(h.useCase.UnenrollCourse(ctx, req.GetUserId(), req.GetCourseId()))
}

func (h *Handler) ListMyCourses(ctx context.Context, req *learningpb.ListMyCoursesRequest) (*learningpb.ListEnrollmentsResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	enrollments, err := h.useCase.ListMyCourses(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	items := make([]*learningpb.EnrollmentResponse, 0, len(enrollments))
	for _, enrollment := range enrollments {
		items = append(items, toEnrollment(enrollment))
	}
	return &learningpb.ListEnrollmentsResponse{Enrollments: items}, nil
}

func (h *Handler) GetEnrollment(ctx context.Context, req *learningpb.GetEnrollmentRequest) (*learningpb.EnrollmentResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Required(req.GetCourseId()) {
		return nil, grpcerrors.InvalidArgument("user_id and course_id are required")
	}
	enrollment, err := h.useCase.GetEnrollment(ctx, req.GetUserId(), req.GetCourseId())
	if err != nil {
		return nil, grpcerrors.NotFound("enrollment not found")
	}
	return toEnrollment(enrollment), nil
}

func (h *Handler) MarkLessonCompleted(ctx context.Context, req *learningpb.MarkLessonCompletedRequest) (*learningpb.ProgressResponse, error) {
	if err := authz.RequireRole(ctx, "student", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetUserId()) || !validate.Required(req.GetCourseId()) || !validate.Required(req.GetLessonId()) {
		return nil, grpcerrors.InvalidArgument("user_id, course_id and lesson_id are required")
	}
	progress, err := h.useCase.MarkLessonCompleted(ctx, req.GetUserId(), req.GetCourseId(), req.GetLessonId(), int(req.GetTotalLessons()))
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toProgress(progress), nil
}

func (h *Handler) GetCourseProgress(ctx context.Context, req *learningpb.GetCourseProgressRequest) (*learningpb.ProgressResponse, error) {
	if !validate.Required(req.GetUserId()) || !validate.Required(req.GetCourseId()) {
		return nil, grpcerrors.InvalidArgument("user_id and course_id are required")
	}
	progress, err := h.useCase.GetCourseProgress(ctx, req.GetUserId(), req.GetCourseId())
	if err != nil {
		return nil, grpcerrors.NotFound("progress not found")
	}
	return toProgress(progress), nil
}

func (h *Handler) GetUserProgress(ctx context.Context, req *learningpb.GetUserProgressRequest) (*learningpb.ListProgressResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	progressItems, err := h.useCase.GetUserProgress(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	items := make([]*learningpb.ProgressResponse, 0, len(progressItems))
	for _, progress := range progressItems {
		items = append(items, toProgress(progress))
	}
	return &learningpb.ListProgressResponse{Progress: items}, nil
}

func (h *Handler) CreateAssignment(ctx context.Context, req *learningpb.CreateAssignmentRequest) (*learningpb.AssignmentResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetCourseId()) || !validate.Required(req.GetLessonId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("course_id, lesson_id and title are required")
	}
	dueAt, _ := time.Parse(time.RFC3339, req.GetDueAt())
	assignment, err := h.useCase.CreateAssignment(ctx, req.GetCourseId(), req.GetLessonId(), req.GetTitle(), dueAt)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toAssignment(assignment), nil
}

func (h *Handler) UpdateAssignment(ctx context.Context, req *learningpb.UpdateAssignmentRequest) (*learningpb.AssignmentResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("id and title are required")
	}
	dueAt, _ := time.Parse(time.RFC3339, req.GetDueAt())
	assignment, err := h.useCase.UpdateAssignment(ctx, model.Assignment{ID: req.GetId(), Title: req.GetTitle(), DueAt: dueAt})
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toAssignment(assignment), nil
}

func (h *Handler) SubmitAssignment(ctx context.Context, req *learningpb.SubmitAssignmentRequest) (*learningpb.SubmissionResponse, error) {
	if !validate.Required(req.GetAssignmentId()) || !validate.Required(req.GetUserId()) || !validate.Required(req.GetAnswer()) {
		return nil, grpcerrors.InvalidArgument("assignment_id, user_id and answer are required")
	}
	submission, err := h.useCase.SubmitAssignment(ctx, req.GetAssignmentId(), req.GetUserId(), req.GetAnswer())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toSubmission(submission), nil
}

func (h *Handler) GetSubmission(ctx context.Context, req *learningpb.GetSubmissionRequest) (*learningpb.SubmissionResponse, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	submission, err := h.useCase.GetSubmission(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("submission not found")
	}
	return toSubmission(submission), nil
}

func (h *Handler) ListSubmissions(ctx context.Context, req *learningpb.ListSubmissionsRequest) (*learningpb.ListSubmissionsResponse, error) {
	if !validate.Required(req.GetAssignmentId()) {
		return nil, grpcerrors.InvalidArgument("assignment_id is required")
	}
	submissions, err := h.useCase.ListSubmissions(ctx, req.GetAssignmentId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toSubmissionList(submissions), nil
}

func (h *Handler) GradeAssignment(ctx context.Context, req *learningpb.GradeAssignmentRequest) (*learningpb.SubmissionResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetSubmissionId()) {
		return nil, grpcerrors.InvalidArgument("submission_id is required")
	}
	if req.GetGrade() < 0 || req.GetGrade() > 100 {
		return nil, grpcerrors.InvalidArgument("grade must be between 0 and 100")
	}
	submission, err := h.useCase.GradeAssignment(ctx, req.GetSubmissionId(), req.GetGrade())
	if err != nil {
		return nil, grpcerrors.NotFound("submission not found")
	}
	return toSubmission(submission), nil
}

func (h *Handler) GetGrades(ctx context.Context, req *learningpb.GetGradesRequest) (*learningpb.ListSubmissionsResponse, error) {
	if !validate.Required(req.GetUserId()) {
		return nil, grpcerrors.InvalidArgument("user_id is required")
	}
	submissions, err := h.useCase.GetGrades(ctx, req.GetUserId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toSubmissionList(submissions), nil
}

func toEnrollment(enrollment model.Enrollment) *learningpb.EnrollmentResponse {
	return &learningpb.EnrollmentResponse{Id: enrollment.ID, UserId: enrollment.UserID, CourseId: enrollment.CourseID, Status: enrollment.Status}
}

func toProgress(progress model.Progress) *learningpb.ProgressResponse {
	return &learningpb.ProgressResponse{UserId: progress.UserID, CourseId: progress.CourseID, CompletedLessons: int32(progress.CompletedLessons), TotalLessons: int32(progress.TotalLessons), Percentage: progress.Percentage}
}

func toAssignment(assignment model.Assignment) *learningpb.AssignmentResponse {
	return &learningpb.AssignmentResponse{Id: assignment.ID, CourseId: assignment.CourseID, LessonId: assignment.LessonID, Title: assignment.Title, DueAt: assignment.DueAt.Format(time.RFC3339)}
}

func toSubmission(submission model.Submission) *learningpb.SubmissionResponse {
	return &learningpb.SubmissionResponse{Id: submission.ID, AssignmentId: submission.AssignmentID, UserId: submission.UserID, Answer: submission.Answer, Grade: submission.Grade, Status: submission.Status}
}

func toSubmissionList(submissions []model.Submission) *learningpb.ListSubmissionsResponse {
	items := make([]*learningpb.SubmissionResponse, 0, len(submissions))
	for _, submission := range submissions {
		items = append(items, toSubmission(submission))
	}
	return &learningpb.ListSubmissionsResponse{Submissions: items}
}
