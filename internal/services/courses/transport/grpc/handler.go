package grpc

import (
	"context"
	"online-learning-platform/internal/shared/authz"
	"online-learning-platform/internal/shared/grpcerrors"
	"online-learning-platform/internal/shared/validate"

	"online-learning-platform/internal/services/courses/model"
	"online-learning-platform/internal/services/courses/usecase"
	coursepb "online-learning-platform/pkg/gen/courses"
)

type Handler struct {
	coursepb.UnimplementedCourseServiceServer
	useCase *usecase.CourseUseCase
}

func NewHandler(useCase *usecase.CourseUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (h *Handler) CreateCourse(ctx context.Context, req *coursepb.CreateCourseRequest) (*coursepb.CourseResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("title is required")
	}
	if !validate.Required(req.GetInstructorId()) {
		return nil, grpcerrors.InvalidArgument("instructor_id is required")
	}
	course, err := h.useCase.Create(ctx, req.GetTitle(), req.GetDescription(), req.GetInstructorId())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toCourse(course), nil
}

func (h *Handler) GetCourse(ctx context.Context, req *coursepb.GetCourseRequest) (*coursepb.CourseResponse, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	course, err := h.useCase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("course not found")
	}
	return toCourse(course), nil
}

func (h *Handler) UpdateCourse(ctx context.Context, req *coursepb.UpdateCourseRequest) (*coursepb.CourseResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("id and title are required")
	}
	course, err := h.useCase.GetByID(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("course not found")
	}
	course.Title = req.GetTitle()
	course.Description = req.GetDescription()
	updated, err := h.useCase.Update(ctx, course)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toCourse(updated), nil
}

func (h *Handler) DeleteCourse(ctx context.Context, req *coursepb.DeleteCourseRequest) (*coursepb.Empty, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	return &coursepb.Empty{}, grpcerrors.Internal(h.useCase.Delete(ctx, req.GetId()))
}

func (h *Handler) ListCourses(ctx context.Context, _ *coursepb.ListCoursesRequest) (*coursepb.ListCoursesResponse, error) {
	courses, err := h.useCase.List(ctx)
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toCourseList(courses), nil
}

func (h *Handler) SearchCourses(ctx context.Context, req *coursepb.SearchCoursesRequest) (*coursepb.ListCoursesResponse, error) {
	courses, err := h.useCase.Search(ctx, req.GetQuery())
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toCourseList(courses), nil
}

func (h *Handler) PublishCourse(ctx context.Context, req *coursepb.PublishCourseRequest) (*coursepb.CourseResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	course, err := h.useCase.Publish(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("course not found")
	}
	return toCourse(course), nil
}

func (h *Handler) CreateModule(ctx context.Context, req *coursepb.CreateModuleRequest) (*coursepb.ModuleResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetCourseId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("course_id and title are required")
	}
	module, err := h.useCase.CreateModule(ctx, req.GetCourseId(), req.GetTitle(), int(req.GetPosition()))
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toModule(module), nil
}

func (h *Handler) UpdateModule(ctx context.Context, req *coursepb.UpdateModuleRequest) (*coursepb.ModuleResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("id and title are required")
	}
	module, err := h.useCase.UpdateModule(ctx, model.Module{ID: req.GetId(), Title: req.GetTitle(), Position: int(req.GetPosition())})
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toModule(module), nil
}

func (h *Handler) DeleteModule(ctx context.Context, req *coursepb.DeleteModuleRequest) (*coursepb.Empty, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	return &coursepb.Empty{}, grpcerrors.Internal(h.useCase.DeleteModule(ctx, req.GetId()))
}

func (h *Handler) CreateLesson(ctx context.Context, req *coursepb.CreateLessonRequest) (*coursepb.LessonResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetCourseId()) || !validate.Required(req.GetModuleId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("course_id, module_id and title are required")
	}
	lesson, err := h.useCase.CreateLesson(ctx, req.GetCourseId(), req.GetModuleId(), req.GetTitle(), req.GetContent(), int(req.GetPosition()))
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toLesson(lesson), nil
}

func (h *Handler) UpdateLesson(ctx context.Context, req *coursepb.UpdateLessonRequest) (*coursepb.LessonResponse, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) || !validate.Required(req.GetTitle()) {
		return nil, grpcerrors.InvalidArgument("id and title are required")
	}
	lesson, err := h.useCase.UpdateLesson(ctx, model.Lesson{ID: req.GetId(), Title: req.GetTitle(), Content: req.GetContent(), Position: int(req.GetPosition())})
	if err != nil {
		return nil, grpcerrors.Internal(err)
	}
	return toLesson(lesson), nil
}

func (h *Handler) DeleteLesson(ctx context.Context, req *coursepb.DeleteLessonRequest) (*coursepb.Empty, error) {
	if err := authz.RequireRole(ctx, "instructor", "admin"); err != nil {
		return nil, err
	}
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	return &coursepb.Empty{}, grpcerrors.Internal(h.useCase.DeleteLesson(ctx, req.GetId()))
}

func (h *Handler) GetLesson(ctx context.Context, req *coursepb.GetLessonRequest) (*coursepb.LessonResponse, error) {
	if !validate.Required(req.GetId()) {
		return nil, grpcerrors.InvalidArgument("id is required")
	}
	lesson, err := h.useCase.GetLesson(ctx, req.GetId())
	if err != nil {
		return nil, grpcerrors.NotFound("lesson not found")
	}
	return toLesson(lesson), nil
}

func toCourse(course model.Course) *coursepb.CourseResponse {
	return &coursepb.CourseResponse{Id: course.ID, Title: course.Title, Description: course.Description, InstructorId: course.InstructorID, Published: course.Published}
}

func toCourseList(courses []model.Course) *coursepb.ListCoursesResponse {
	items := make([]*coursepb.CourseResponse, 0, len(courses))
	for _, course := range courses {
		items = append(items, toCourse(course))
	}
	return &coursepb.ListCoursesResponse{Courses: items}
}

func toModule(module model.Module) *coursepb.ModuleResponse {
	return &coursepb.ModuleResponse{Id: module.ID, CourseId: module.CourseID, Title: module.Title, Position: int32(module.Position)}
}

func toLesson(lesson model.Lesson) *coursepb.LessonResponse {
	return &coursepb.LessonResponse{Id: lesson.ID, CourseId: lesson.CourseID, ModuleId: lesson.ModuleID, Title: lesson.Title, Content: lesson.Content, Position: int32(lesson.Position)}
}
