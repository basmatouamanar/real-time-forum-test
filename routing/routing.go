package routing

import (
	"net/http"
	"time"

	"forum/handlers"
	"forum/middleware"
)

func Routing() {
	rateLimiterLogin := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterRegister := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterPost := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterComment := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterRefresh := middleware.NewRateLimiterManager(200, time.Minute)
	rateLimiterReaction := middleware.NewRateLimiterManager(20, time.Minute)
	rateLimiterCommentReaction := middleware.NewRateLimiterManager(20, time.Minute)

	http.HandleFunc("/api/posts", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.HanldlerShowHome))
	http.HandleFunc("/createcomment", middleware.RateLimitMiddleware(rateLimiterComment, handlers.CreateCommentHandler))
	
	http.HandleFunc("/loginAuth", middleware.RateLimitMiddleware(rateLimiterLogin, handlers.LoginHandler))
	http.HandleFunc("/registerAuth", middleware.RateLimitMiddleware(rateLimiterRegister, handlers.RegisterHandler))
	
	http.HandleFunc("/login", middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.CheckLogin(handlers.Showloginhandler)))
	http.HandleFunc("/register", middleware.RateLimitMiddleware(rateLimiterRefresh, middleware.CheckLogin(handlers.Showregister)))
	http.HandleFunc("/static/", handlers.StyleFunc)
	http.HandleFunc("/", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.HanldlerShowHome))
	http.HandleFunc("/createpost", middleware.RateLimitMiddleware(rateLimiterPost, middleware.Checksession(handlers.CreatePostHandler)))
	http.HandleFunc("/logout", middleware.RateLimitMiddleware(rateLimiterRefresh, handlers.LogOutHandler))
	http.HandleFunc("/reaction", middleware.RateLimitMiddleware(rateLimiterReaction, middleware.Checksession(handlers.ReactionHandler)))
	http.HandleFunc("/comment-reaction", middleware.RateLimitMiddleware(rateLimiterCommentReaction, middleware.Checksession(handlers.CommentReactionHandler)))
}