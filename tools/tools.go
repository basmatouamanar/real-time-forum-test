package tools

type Post struct {
	ID           int      `json:"id"`
	Title        string   `json:"title"`
	Description  string   `json:"description"`
	ImageUrl     string   `json:"imageUrl"`
	UserName     string   `json:"userName"`
	CreationDate string   `json:"creationDate"`
	Categories   []string `json:"categories"`
}

type Category struct {
	ID       int    `json:"id"`
	Category string `json:"category"`
}

type ReactionStats struct {
	PostID        int `json:"postId"`
	LikesCount    int `json:"likesCount"`
	DislikesCount int `json:"dislikesCount"`
}

type PageData struct {
	Posts                  []Post                       `json:"posts"`
	Categories             []Category                   `json:"categories"`
	IsLogin                IsLogin                      `json:"isLogin"`
	ReactionStats          map[int]ReactionStats        `json:"reactionStats"`
	UserReactions          map[int]int                  `json:"userReactions"`
	Comment                map[int][]Comment            `json:"comments"`
	ConnectUserName        string                       `json:"connectUserName"`
	CommentReactionStats   map[int]CommentReactionStats `json:"commentReactionStats"`
	UserCommentReactions   map[int]int                  `json:"userCommentReactions"`
}

type IsLogin struct {
	LoggedIn bool `json:"loggedIn"`
	UserID   int  `json:"userId"`
}

type Comment struct {
	ID           int    `json:"id"`
	CommentText  string `json:"commentText"`
	PostID       int    `json:"postId"`
	UserID       int    `json:"userId"`
	UserName     string `json:"userName"`
	CreationDate string `json:"creationDate"`
}

type CommentReactionStats struct {
	CommentID     int `json:"commentId"`
	LikesCount    int `json:"likesCount"`
	DislikesCount int `json:"dislikesCount"`
}