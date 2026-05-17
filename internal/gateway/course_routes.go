package gateway

import (
	"net/http"

	coursepb "online-learning-platform/pkg/gen/courses"
)

func registerCourseRoutes(mux *http.ServeMux, clients *Clients) {
	mux.HandleFunc("/api/courses", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		switch r.Method {
		case http.MethodGet:
			resp, err := clients.Courses.ListCourses(requestContext(r), &coursepb.ListCoursesRequest{})
			writeGRPC(w, resp, err)
		case http.MethodPost:
			var req coursepb.CreateCourseRequest
			if !decode(w, r, &req) {
				return
			}
			resp, err := clients.Courses.CreateCourse(grpcContext(r), &req)
			writeGRPC(w, resp, err)
		default:
			methodAllowed(w)
		}
	})

	mux.HandleFunc("/api/courses/search", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			methodAllowed(w)
			return
		}
		if !clientsAvailable(w, clients) {
			return
		}
		resp, err := clients.Courses.SearchCourses(requestContext(r), &coursepb.SearchCoursesRequest{Query: r.URL.Query().Get("query")})
		writeGRPC(w, resp, err)
	})

	mux.HandleFunc("/api/courses/modules/", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/courses/modules/")
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		moduleID := parts[0]
		switch r.Method {
		case http.MethodPatch:
			var req coursepb.UpdateModuleRequest
			if !decode(w, r, &req) {
				return
			}
			req.Id = moduleID
			resp, err := clients.Courses.UpdateModule(grpcContext(r), &req)
			writeGRPC(w, resp, err)
		case http.MethodDelete:
			resp, err := clients.Courses.DeleteModule(grpcContext(r), &coursepb.DeleteModuleRequest{Id: moduleID})
			writeGRPC(w, resp, err)
		default:
			methodAllowed(w)
		}
	})

	mux.HandleFunc("/api/courses/lessons/", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/courses/lessons/")
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		lessonID := parts[0]
		switch r.Method {
		case http.MethodGet:
			resp, err := clients.Courses.GetLesson(requestContext(r), &coursepb.GetLessonRequest{Id: lessonID})
			writeGRPC(w, resp, err)
		case http.MethodPatch:
			var req coursepb.UpdateLessonRequest
			if !decode(w, r, &req) {
				return
			}
			req.Id = lessonID
			resp, err := clients.Courses.UpdateLesson(grpcContext(r), &req)
			writeGRPC(w, resp, err)
		case http.MethodDelete:
			resp, err := clients.Courses.DeleteLesson(grpcContext(r), &coursepb.DeleteLessonRequest{Id: lessonID})
			writeGRPC(w, resp, err)
		default:
			methodAllowed(w)
		}
	})

	mux.HandleFunc("/api/courses/", func(w http.ResponseWriter, r *http.Request) {
		if !clientsAvailable(w, clients) {
			return
		}
		parts := pathSegments(r.URL.Path, "/api/courses/")
		if len(parts) == 0 {
			http.NotFound(w, r)
			return
		}
		courseID := parts[0]
		if len(parts) == 2 && parts[1] == "publish" {
			if r.Method != http.MethodPost {
				methodAllowed(w)
				return
			}
			resp, err := clients.Courses.PublishCourse(grpcContext(r), &coursepb.PublishCourseRequest{Id: courseID})
			writeGRPC(w, resp, err)
			return
		}
		if len(parts) == 2 && parts[1] == "modules" {
			if r.Method != http.MethodPost {
				methodAllowed(w)
				return
			}
			var req coursepb.CreateModuleRequest
			if !decode(w, r, &req) {
				return
			}
			req.CourseId = courseID
			resp, err := clients.Courses.CreateModule(grpcContext(r), &req)
			writeGRPC(w, resp, err)
			return
		}
		if len(parts) == 4 && parts[1] == "modules" && parts[3] == "lessons" {
			if r.Method != http.MethodPost {
				methodAllowed(w)
				return
			}
			var req coursepb.CreateLessonRequest
			if !decode(w, r, &req) {
				return
			}
			req.CourseId = courseID
			req.ModuleId = parts[2]
			resp, err := clients.Courses.CreateLesson(grpcContext(r), &req)
			writeGRPC(w, resp, err)
			return
		}
		if len(parts) != 1 {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			resp, err := clients.Courses.GetCourse(requestContext(r), &coursepb.GetCourseRequest{Id: courseID})
			writeGRPC(w, resp, err)
		case http.MethodPatch:
			var req coursepb.UpdateCourseRequest
			if !decode(w, r, &req) {
				return
			}
			req.Id = courseID
			resp, err := clients.Courses.UpdateCourse(grpcContext(r), &req)
			writeGRPC(w, resp, err)
		case http.MethodDelete:
			resp, err := clients.Courses.DeleteCourse(grpcContext(r), &coursepb.DeleteCourseRequest{Id: courseID})
			writeGRPC(w, resp, err)
		default:
			methodAllowed(w)
		}
	})
}
