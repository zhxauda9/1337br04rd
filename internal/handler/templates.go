package handler

import (
	"html/template"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"1337b04rd/pkg/middleware"
)

type TemplateHandler struct {
	templates *template.Template
	logger    *slog.Logger
}

func NewTemplateHandler(logger *slog.Logger) *TemplateHandler {
	tmpl := template.Must(template.ParseGlob("internal/adapters/frontend/templates/*.html"))
	return &TemplateHandler{templates: tmpl, logger: logger}
}

func (t *TemplateHandler) RenderHomePage(sessionHandler *SessionHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := r.Cookie("session_id")
		if err != nil {
			session, err := sessionHandler.sessionService.CreateSession(r.Context(), "")
			if err != nil {
				t.logger.Error("Failed to create session", "error", err, "method", r.Method)
				handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to create session.")
				handler.ServeHTTP(w, r)
				return
			}

			http.SetCookie(w, &http.Cookie{
				Name:     "session_id",
				Value:    session.ID,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
				Expires:  session.ExpiresAt,
			})
			http.SetCookie(w, &http.Cookie{
				Name:     "session_name",
				Value:    session.Name,
				Path:     "/",
				HttpOnly: true,
				Secure:   false,
				SameSite: http.SameSiteStrictMode,
				Expires:  session.ExpiresAt,
			})
		}

		// client := externalapi.NewRickAndMortyClient()
		// characters, err := client.FetchAllCharacters()
		// if err != nil {
		// 	t.logger.Error("Failed to fetch characters from the API.", "error", err, "method", r.Method)
		// 	handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to fetch characters from the API.")
		// 	handler.ServeHTTP(w, r)
		// 	return
		// }

		// if len(characters) > 21 {
		// 	characters = characters[:21]
		// }

		data := map[string]interface{}{
			"Title": "1337b04rd",
		}

		err = t.templates.ExecuteTemplate(w, "home-content.html", data)
		if err != nil {
			t.logger.Error("Failed to render the home page.", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render the home page.")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered home page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderCreatePostPage() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Title": "Create Post",
		}

		err := t.templates.ExecuteTemplate(w, "create-post.html", data)
		if err != nil {
			t.logger.Error("Failed to render created post page.", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render created post page")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered created post page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderCatalogPage(PostHandler *PostHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts, err := PostHandler.postService.GetAllPosts(r.Context())
		if err != nil {
			t.logger.Error("Failed to fetch posts", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to fetch posts")
			handler.ServeHTTP(w, r)
			return
		}

		session, err := middleware.GetSession(r.Context())
		if err != nil {
			t.logger.Error("Error retrieving session", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusUnauthorized, "Error retrieving session.")
			handler.ServeHTTP(w, r)
			return
		}

		data := map[string]interface{}{
			"Title": "Catalog",
			"Name":  session.Name,
			"Posts": posts,
		}

		err = t.templates.ExecuteTemplate(w, "catalog.html", data)
		if err != nil {
			t.logger.Error("Failed to render catalog page", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render catalog page")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered catalog page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderArchivePage(ArchiveHandler *ArchiveHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		posts, err := ArchiveHandler.service.GetAllArchivedPosts(r.Context())
		if err != nil {
			t.logger.Error("Failed to fetch archived posts", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to fetch archived posts")
			handler.ServeHTTP(w, r)
			return
		}

		session, err := middleware.GetSession(r.Context())
		if err != nil {
			t.logger.Error("Failed to fetch characters from the API", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to fetch characters from the API.")
			handler.ServeHTTP(w, r)
			return
		}

		data := map[string]interface{}{
			"Title": "Archived Threads",
			"Name":  session.Name,
			"Posts": posts,
		}
		err = t.templates.ExecuteTemplate(w, "archive.html", data)
		if err != nil {
			t.logger.Error("Failed to render archive page", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render archive page")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered archive page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderPostPage(PostHandler *PostHandler, CommentHandler *CommentHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.URL.Path, "/")
		postID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			t.logger.Error("Invalid Post ID", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusBadRequest, "Invalid Post ID")
			handler.ServeHTTP(w, r)
			return
		}

		post, err := PostHandler.postService.GetPost(r.Context(), postID)
		if err != nil {
			t.logger.Error("Post not found", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusNotFound, "Post not found")
			handler.ServeHTTP(w, r)
			return
		}

		comments, err := CommentHandler.commentService.GetAllCommentsOfPost(r.Context(), postID)
		if err != nil {
			t.logger.Error("Failed to fetch comments", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to fetch comments")
			handler.ServeHTTP(w, r)
			return
		}

		data := map[string]interface{}{
			"Title":    "Post - " + post.Title,
			"Post":     post,
			"Comments": comments,
		}

		err = t.templates.ExecuteTemplate(w, "post.html", data)
		if err != nil {
			t.logger.Error("Failed to render page", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render page")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered post page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderArchivedPostPage(ArchiveHandler *ArchiveHandler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
		if len(parts) < 2 {
			t.logger.Error("Invalid URL format", "method", r.Method)
			handler := t.RenderErrorPage(http.StatusBadRequest, "Invalid URL format")
			handler.ServeHTTP(w, r)
			return
		}
		postID, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			t.logger.Error("Invalid post ID", "method", r.Method)
			handler := t.RenderErrorPage(http.StatusBadRequest, "Invalid post ID")
			handler.ServeHTTP(w, r)
			return
		}

		post := ArchiveHandler.service.GetArchivedPostByID(r.Context(), postID)
		if err != nil {
			t.logger.Error("Post not found", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusNotFound, "Post not found")
			handler.ServeHTTP(w, r)
			return
		}

		data := map[string]interface{}{
			"Title":    "Archived Post",
			"Post":     post,
			"Comments": post.Comments,
		}

		err = t.templates.ExecuteTemplate(w, "archive-post.html", data)
		if err != nil {
			t.logger.Error("Failed to render page", "error", err, "method", r.Method)
			handler := t.RenderErrorPage(http.StatusInternalServerError, "Failed to render page")
			handler.ServeHTTP(w, r)
			return
		}
		t.logger.Info("Successfully rendered archived post page", "method", r.Method)
	})
}

func (t *TemplateHandler) RenderErrorPage(code int, message string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]interface{}{
			"Code":    code,
			"Message": message,
		}

		err := t.templates.ExecuteTemplate(w, "error.html", data)
		if err != nil {
			t.logger.Error("Failed to render error page", "error", err, "method", r.Method)
			http.Error(w, "Failed to render error page", http.StatusInternalServerError)
			return
		}
		t.logger.Info("Successfully rendered error page", "method", r.Method)
	})
}
