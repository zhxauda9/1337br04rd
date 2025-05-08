package app

import (
	"database/sql"
	"fmt"
	"log/slog"
	"net/http"

	"1337b04rd/internal/adapters/storage"
	"1337b04rd/internal/handler"
	"1337b04rd/internal/repository"
	"1337b04rd/internal/service"
	"1337b04rd/pkg/middleware"
)

func NewApp(db *sql.DB, logger *slog.Logger) http.Handler {
	// Initialize Minio storage
	s3Storage, err := storage.NewMinioClient(
		"minio:9000", // endpoint
		"minioadmin", // access key
		"minioadmin", // secret key
		false,        // useSSL
	)
	if err != nil {
		panic(fmt.Sprintf("failed to initialize Minio client: %v", err))
	}

	// Initialize repositories
	commentRepo := repository.NewCommentRepository(db)
	postRepo := repository.NewPostRepository(db)
	archiveRepo := repository.NewArchiveRepository(db)

	// Initialize services
	commentService := service.NewCommentService(commentRepo, postRepo)
	sessionService := service.NewSessionService(s3Storage, postRepo, commentRepo)
	postService := service.NewPostService(postRepo)
	archiveService := service.NewArchiveService(archiveRepo, postRepo)

	// Initialize handlers
	postHandler := handler.NewPostHandler(postService, s3Storage, logger)
	commentHandler := handler.NewCommentHandler(commentService, s3Storage, logger)
	sessionHandler := handler.NewSessionHandler(sessionService, logger)
	archiveHandler := handler.NewArchiveHandler(archiveService, logger)

	templateHandler := handler.NewTemplateHandler(logger)

	// Setup routes
	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("internal/adapters/frontend/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	s := NewCleanupService(postService)
	s.StartCleanupTask()

	// Post routes
	wrappedPostCreateHandler := middleware.InjectSessionMiddleware()(http.HandlerFunc(postHandler.CreatePost))
	mux.Handle("/", templateHandler.RenderHomePage(sessionHandler))
	mux.Handle("POST /posts/create", middleware.InjectSessionMiddleware()(wrappedPostCreateHandler))
	mux.HandleFunc("GET /posts/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware.InjectSessionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templateHandler.RenderPostPage(postHandler, commentHandler).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /posts/create", func(w http.ResponseWriter, r *http.Request) {
		middleware.InjectSessionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templateHandler.RenderCreatePostPage().ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /posts", func(w http.ResponseWriter, r *http.Request) {
		middleware.InjectSessionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templateHandler.RenderCatalogPage(postHandler).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /images/posts/{filename}", storage.ServePostImageHandler(s3Storage))
	mux.HandleFunc("GET /images/comments/{filename}", storage.ServeCommentImageHandler(s3Storage))

	// Comment routes
	wrappedCommentCreateHandler := middleware.InjectSessionMiddleware()(http.HandlerFunc(commentHandler.CreateComment))
	mux.Handle("POST /comments/create", wrappedCommentCreateHandler)
	mux.HandleFunc("GET /comments/post/{id}", commentHandler.GetCommentsOfPost)
	mux.HandleFunc("GET /comments/replies/{id}", commentHandler.GetRepliesToComment)

	// archive routes
	wrappedArchiveHandler := middleware.InjectSessionMiddleware()(http.HandlerFunc(archiveHandler.ArchivePost))
	mux.Handle("POST /archive/{id}", wrappedArchiveHandler)
	mux.HandleFunc("POST /archive-expired-posts", archiveHandler.ArchiveExpired)
	mux.HandleFunc("GET /archive", func(w http.ResponseWriter, r *http.Request) {
		middleware.InjectSessionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templateHandler.RenderArchivePage(archiveHandler).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	mux.HandleFunc("GET /archive/{id}", func(w http.ResponseWriter, r *http.Request) {
		middleware.InjectSessionMiddleware()(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			templateHandler.RenderArchivedPostPage(archiveHandler).ServeHTTP(w, r)
		})).ServeHTTP(w, r)
	})

	// Session routes
	mux.HandleFunc("POST /sessions/create", sessionHandler.CreateSession)
	mux.HandleFunc("GET /sessions/{id}", sessionHandler.GetSession)
	mux.HandleFunc("DELETE /sessions/{id}", sessionHandler.DeleteSession)
	mux.HandleFunc("GET /sessions", sessionHandler.GetAllSession)
	mux.HandleFunc("PUT /sessions/{id}", sessionHandler.UpdateSession)

	return mux
}
