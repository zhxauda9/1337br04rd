package service

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"1337b04rd/internal/adapters/externalapi"
	"1337b04rd/internal/adapters/storage"
	"1337b04rd/internal/domain"
	"1337b04rd/internal/ports"
)

type SessionService struct {
	minioClient     *storage.MinioClient
	sessions        map[string]*domain.Session
	assignedAvatars map[string]bool
	allCharacters   []domain.Character
	postRepo        ports.PostRepository
	commentRepo     ports.CommentRepository
}

func NewSessionService(minioClient *storage.MinioClient, postRepo ports.PostRepository, commentRepo ports.CommentRepository) *SessionService {
	svc := &SessionService{
		minioClient:     minioClient,
		sessions:        make(map[string]*domain.Session),
		assignedAvatars: make(map[string]bool),
		allCharacters:   []domain.Character{},
		postRepo:        postRepo,
		commentRepo:     commentRepo,
	}
	svc.loadAllCharacters()
	return svc
}

func (s *SessionService) loadAllCharacters() {
	client := externalapi.NewRickAndMortyClient()

	characters, err := client.FetchAllCharacters()
	if err != nil {
		panic(fmt.Sprintf("Failed to load Rick & Morty characters: %v", err))
	}

	for _, c := range characters {
		if c.Image != "" {
			s.allCharacters = append(s.allCharacters, c)
		}
	}

	fmt.Printf("Loaded %d valid characters\n", len(s.allCharacters))
}

func (s *SessionService) pickRandomCharacter() domain.Character {
	rand.Seed(time.Now().UnixNano())
	idx := rand.Intn(len(s.allCharacters))
	return s.allCharacters[idx]
}

func (s *SessionService) resetAssignedAvatarsIfNeeded() {
	if len(s.assignedAvatars) >= 826 {
		fmt.Println("All avatars assigned. Resetting assigned avatars.")
		s.assignedAvatars = make(map[string]bool)
	}
}

func generateSessionID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%d-%d", time.Now().UnixNano(), rand.Intn(1000))
}

func (s *SessionService) CreateSession(ctx context.Context, name string) (*domain.Session, error) {
	s.resetAssignedAvatarsIfNeeded()

	var char domain.Character

	for attempt := 0; attempt < 10; attempt++ {
		char = s.pickRandomCharacter()

		if !s.assignedAvatars[char.Image] {
			break
		}

		fmt.Println("Avatar already assigned, retrying...")
	}

	avatarURL, err := s.minioClient.UploadAvatarFromURL(ctx, char.Image)
	if err != nil {
		return nil, fmt.Errorf("failed to upload avatar: %w", err)
	}

	if name == "" {
		name = char.Name
	}

	sessionID := generateSessionID()

	session, err := domain.NewSession(sessionID, name, avatarURL)
	if err != nil {
		return nil, err
	}

	s.sessions[session.ID] = session
	s.assignedAvatars[char.Image] = true

	fmt.Printf("Session created: ID=%s, Name=%s, AvatarURL=%s\n", session.ID, session.Name, session.AvatarURL)

	return session, nil
}

func (s *SessionService) UpdateSession(sessionID, newName string) (*domain.Session, error) {
	session, ok := s.sessions[sessionID]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}

	if newName == "" {
		return nil, fmt.Errorf("name cannot be empty")
	}

	// Update the session name
	session.Name = newName
	sessionM := &domain.Session{}
	sessionM.Name = newName

	// Now, update all posts and comments associated with this session
	err := s.updatePostsAndComments(sessionID, newName)
	if err != nil {
		return nil, fmt.Errorf("failed to update posts/comments: %w", err)
	}

	return session, nil
}

func (s *SessionService) updatePostsAndComments(sessionID, newName string) error {
	posts, err := s.getPostsByAuthor(sessionID)
	if err != nil {
		return fmt.Errorf("failed to fetch posts: %w", err)
	}

	comments, err := s.GetCommentsByAuthor(sessionID)
	if err != nil {
		return fmt.Errorf("failed to fetch comments: %w", err)
	}

	for _, post := range posts {
		post.AuthorName = newName
		if err := s.postRepo.Update(context.Background(), post); err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}
	}

	for _, comment := range comments {
		comment.AuthorName = newName
		if err := s.commentRepo.Update(context.Background(), comment); err != nil {
			return fmt.Errorf("failed to update comment: %w", err)
		}
	}

	return nil
}

func (s *SessionService) DeleteSession(w http.ResponseWriter, r *http.Request) error {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return fmt.Errorf("failed to get session cookie: %w", err)
	}

	delete(s.sessions, cookie.Value)

	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	return nil
}

// GetSession retrieves the session from the context (modified to get name from cookie)
func (s *SessionService) GetSession(ctx context.Context, sessionID string) (*domain.Session, error) {
	// Check if session exists in memory
	session, ok := s.sessions[sessionID]
	if ok {
		return session, nil
	}

	// If not, retrieve the session ID from the cookie and name from the cookie (if present)
	cookie, err := ctx.Value("request").(*http.Request).Cookie("session_id")
	if err != nil || cookie == nil {
		return nil, fmt.Errorf("session not found")
	}

	// Use the cookie value as the session ID and create a new session using cookie name
	session = &domain.Session{
		ID:   cookie.Value,
		Name: cookie.Name, // Assuming the name is stored in the cookie. Adjust as needed
	}

	return session, nil
}

func (s *SessionService) GetAllSessions(ctx context.Context) ([]*domain.Session, error) {
	sessions := make([]*domain.Session, 0, len(s.sessions))
	for _, session := range s.sessions {
		sessions = append(sessions, session)
	}
	return sessions, nil
}

// Retrieve posts linked to this session's author_id
func (s *SessionService) getPostsByAuthor(sessionID string) ([]*domain.Post, error) {
	posts, err := s.postRepo.FindByAuthorID(context.Background(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch posts: %w", err)
	}
	return posts, nil
}

// Retrieve comments linked to this session's author_id
func (s *SessionService) GetCommentsByAuthor(sessionID string) ([]*domain.Comment, error) {
	comments, err := s.commentRepo.FindByAuthorID(context.Background(), sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	return comments, nil
}
