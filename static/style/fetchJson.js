document.addEventListener('DOMContentLoaded', () => {
    
    fetchapi();
});


async function fetchapi() {
    const container = document.getElementById('posts-container');

    try {
        const response = await fetch("/api/posts");

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const pageData = await response.json();
        displayPosts(pageData);
    } catch (error) {
        console.error('Error fetching posts:', error);
        container.innerHTML = '<div class="error">❌ Failed to load posts. Please try again later.</div>';
    }
}


function displayPosts(pageData) {
    const container = document.getElementById('posts-container');
    container.innerHTML = '';

    if (!pageData.posts || pageData.posts.length === 0) {
        container.innerHTML = '<div class="loading">No posts yet. Be the first to post!</div>';
        return;
    }

    pageData.posts.forEach(post => {
        const reactionStats = pageData.reactionStats[post.id] || { likesCount: 0, dislikesCount: 0 };
        const userReaction = pageData.userReactions[post.id]; // 1 = like, -1 = dislike
        const comments = pageData.comments[post.id] || [];
        
        const postElement = createPostElement(post, reactionStats, userReaction, comments);
        container.appendChild(postElement);
    });
}


function createPostElement(post, reactionStats, userReaction, comments) {
    const postDiv = document.createElement('div');
    postDiv.className = 'post';
    postDiv.id = `post-${post.id}`;

    postDiv.innerHTML = `
        <div class="top-section">
            <div class="full-profile">
                <img src="/static/images/user.png" alt="User Profile" class="profile-image">
                <h2 class="post-username">${escapeHtml(post.userName)}</h2>
            </div>
            <span class="post-time">${new Date(post.creationDate).toLocaleDateString()}</span>
        </div>
        
        <div class="body-section">
            <div class="category-section">
                ${post.categories ? post.categories.map(cat => `<span class="post-category">${escapeHtml(cat)}</span>`).join('') : ''}
            </div>
            <h1 class="post-title">${escapeHtml(post.title)}</h1>
            ${post.imageUrl ? `
                <div class="body-image">
                    <img src="${escapeHtml(post.imageUrl)}" alt="post image" class="post-img">
                </div>
            ` : ''}
            <p class="post-description">${escapeHtml(post.description)}</p>
        </div>
        
        <div class="reactions-section">
            <form class="reaction-form" data-post-id="${post.id}">
                <input type="hidden" name="postId" value="${post.id}">
                
                <div class="likes-section">
                    <button type="submit" class="btn-reaction ${userReaction === 1 ? 'active' : ''}" name="reaction" value="1">
                        <i class="fa-solid fa-thumbs-up icon-reaction"></i>
                    </button>
                    <span class="total-reactions">${reactionStats.likesCount}</span>
                </div>
                
                <div class="deslikes-section">
                    <button type="submit" class="btn-reaction ${userReaction === -1 ? 'active' : ''}" name="reaction" value="-1">
                        <i class="fa-solid fa-thumbs-down icon-reaction"></i>
                    </button>
                    <span class="total-reactions">${reactionStats.dislikesCount}</span>
                </div>
            </form>
            
            <button class="comment-toggle">
                <i class="fa-solid fa-comment"></i> ${comments.length}
            </button>
        </div>
        
        <div class="comments-box" id="comments-box-${post.id}" style="display: none;">
            <div class="comments-list"></div>
            
            <form class="comment-form" data-post-id="${post.id}">
                <input type="hidden" name="postId" value="${post.id}">
                <input type="text" name="comment" placeholder="Add a comment..." class="comment-input">
                <button type="submit" class="btn-comment">Post</button>
            </form>
        </div>
    `;

    // Attach event listeners
    attachPostEvents(postDiv, comments);
    
    return postDiv;
}


function attachPostEvents(postElement, comments) {
    const postId = postElement.id.replace('post-', '');
    
    // REACTION FORM
    const reactionForm = postElement.querySelector('.reaction-form');
    reactionForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        const button = e.submitter;
        if (!button) return;

        const reactionValue = button.value;
        const postId = reactionForm.querySelector("input[name='postId']").value;

        const params = new URLSearchParams();
        params.append("postId", postId);
        params.append("reaction", reactionValue);

        try {
            const res = await fetch('/reaction', {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            });

            const data = await res.json();

            if (!res.ok) {
                alert(data.error || "Erreur lors de la réaction");
                return;
            }

            updateReactionUI(reactionForm, data);

        } catch (err) {
            console.error("Erreur fetch:", err);
        }
    });

    // COMMENT TOGGLE
    const commentToggle = postElement.querySelector('.comment-toggle');
    const commentsBox = postElement.querySelector('.comments-box');
    const commentsList = postElement.querySelector('.comments-list');
    
    commentToggle.addEventListener('click', () => {
        if (commentsBox.style.display === 'none') {
            commentsBox.style.display = 'block';
            renderComments(commentsList, comments);
        } else {
            commentsBox.style.display = 'none';
        }
    });

    // COMMENT FORM
    const commentForm = postElement.querySelector('.comment-form');
    commentForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const postId = commentForm.querySelector("input[name='postId']").value;
        const commentText = commentForm.querySelector("input[name='comment']").value.trim();

        if (!commentText) return;

        const params = new URLSearchParams();
        params.append("postId", postId);
        params.append("comment", commentText);

        try {
            const res = await fetch("/createcomment", {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            });

            const data = await res.json();

            if (!res.ok) {
                alert(data.error || "Erreur lors de l'ajout du commentaire");
                return;
            }

            addCommentToDOM(data, commentForm);
            commentForm.querySelector("input[name='comment']").value = "";
            
            // Update comment count
            const commentToggle = postElement.querySelector('.comment-toggle');
            const currentCount = parseInt(commentToggle.textContent.match(/\d+/)[0]);
            commentToggle.innerHTML = `<i class="fa-solid fa-comment"></i> ${currentCount + 1}`;

        } catch (err) {
            console.error("Erreur fetch:", err);
        }
    });
}


function updateReactionUI(form, reactionData) {
    const likesSection = form.querySelector('.likes-section');
    const dislikesSection = form.querySelector('.deslikes-section');

    if (!likesSection || !dislikesSection) return;

    const likesSpan = likesSection.querySelector('.total-reactions');
    const dislikesSpan = dislikesSection.querySelector('.total-reactions');

    if (likesSpan) likesSpan.textContent = reactionData.likesCount || 0;
    if (dislikesSpan) dislikesSpan.textContent = reactionData.dislikesCount || 0;

    const likeButton = likesSection.querySelector('.btn-reaction');
    const dislikeButton = dislikesSection.querySelector('.btn-reaction');

    likeButton?.classList.remove('active');
    dislikeButton?.classList.remove('active');

    if (reactionData.userReaction === 1) {
        likeButton?.classList.add('active');
    } else if (reactionData.userReaction === -1) {
        dislikeButton?.classList.add('active');
    }
}


function renderComments(container, comments) {
    container.innerHTML = '';
    
    if (comments.length === 0) {
        container.innerHTML = '<p class="noComments">No comments yet</p>';
        return;
    }

    comments.forEach(comment => {
        const commentDiv = document.createElement('div');
        commentDiv.className = 'comments';
        commentDiv.innerHTML = `
            <img src="/static/images/user.png" alt="User Profile" class="profile-comment">
            <div class="comment-content-wrapper">
                <span class="username">${escapeHtml(comment.userName)}</span>
                <p class="comment-content">${escapeHtml(comment.commentText)}</p>
                <span class="post-time">${new Date(comment.creationDate).toLocaleDateString()}</span>
            </div>
        `;
        container.appendChild(commentDiv);
    });
}


function addCommentToDOM(commentData, form) {
    const postId = commentData.postId;
    const commentsBox = document.getElementById(`comments-box-${postId}`);
    if (!commentsBox) return;

    const commentsList = commentsBox.querySelector('.comments-list');
    const noComments = commentsList.querySelector('.noComments');
    if (noComments) noComments.remove();

    const commentDiv = document.createElement('div');
    commentDiv.className = 'comments';
    commentDiv.innerHTML = `
        <img src="/static/images/user.png" alt="User Profile" class="profile-comment">
        <div class="comment-content-wrapper">
            <span class="username">${escapeHtml(commentData.userName)}</span>
            <p class="comment-content">${escapeHtml(commentData.comment)}</p>
            <span class="post-time">${new Date().toLocaleDateString()}</span>
        </div>
    `;

    commentsList.appendChild(commentDiv);
}


function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}