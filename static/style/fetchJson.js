document.addEventListener('DOMContentLoaded', () => {
    const registerContainer = document.getElementById('register-container');
    const registerview = document.getElementById('register-view');
    const logoutbtn = document.getElementById('btn-logout');
    const loginview = document.getElementById('login-view');
    const mainSection = document.querySelector('.main-section');
    const islogged = document.getElementById('app-container').dataset.loggedin === 'true';

    // ÉTAT INITIAL
    if (islogged) {
        registerview.style.display = 'none';
        registerContainer.style.display = 'none';
        loginview.style.display = 'none';
        mainSection.style.display = 'block';
        logoutbtn.style.display = 'flex';
        fetchapi();
    } else {
        registerview.style.display = 'none';
        registerContainer.style.display = 'none';
        loginview.style.display = 'flex';
        mainSection.style.display = 'none';
        logoutbtn.style.display = 'none';
    }

    // LOGIN FORM
    const logginform = document.getElementById('login-form');
    logginform.addEventListener('submit', async (e) => {
        e.preventDefault();
        const userName = document.getElementById('username').value;
        const password = document.getElementById('password').value;
        if (!userName || !password) return;

        const params = new URLSearchParams();
        params.append("username", userName);
        params.append("password", password);

        try {
            const res = await fetch("/loginAuth", {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            });

            const data = await res.json();

            if (!res.ok) {
                console.log("Erreur:", data);
                alert(data.error || "Erreur lors de login");
                return;
            }

            console.log("Login successful!");
            alert("Connexion réussie !");

            // Masquer login/register
            loginview.style.display = 'none';
            registerview.style.display = 'none';
            registerContainer.style.display = 'none';

            // Afficher la home page
            mainSection.style.display = 'block';
            logoutbtn.style.display = 'flex';

            // Charger les posts
            fetchapi();

        } catch (err) {
            console.error("Erreur fetch:", err);
            alert("Erreur réseau ou serveur");
        }
    });

    // REGISTER BUTTON
    const buttonregister = document.querySelector('.link-register');
    buttonregister.addEventListener('click', (e) => {
        e.preventDefault();
        console.log('Showing register form');
        registerview.style.display = 'flex';
        registerContainer.style.display = 'flex';
        loginview.style.display = 'none';
        mainSection.style.display = 'none';
    });

    // LINK TO LOGIN FROM REGISTER
    /*const linkLogin = document.querySelector('.link-login');
    linkLogin.addEventListener('click', (e) => {
        e.preventDefault();
        registerview.style.display = 'none';
        registerContainer.style.display = 'none';
        loginview.style.display = 'flex';
    });*/

    // REGISTER FORM
    const registerForm = document.getElementById('register-form');
    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault();

        const nickname = document.getElementById('nickname').value;
        const firstname = document.getElementById('firstname').value;
        const lastname = document.getElementById('lastname').value;
        const age = document.getElementById('age').value;
        const gender = document.getElementById('gender').value;
        const email = document.getElementById('email').value;
        const passwordre = document.getElementById('passwordre').value;

        const params = new URLSearchParams();
        params.append("nickname", nickname);
        params.append("firstname", firstname);
        params.append("lastname", lastname);
        params.append("age", age);
        params.append("gender", gender);
        params.append("email", email);
        params.append("passwordre", passwordre);

        try {
            const res = await fetch('/registerAuth', {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            });

            const data = await res.json();

            if (!res.ok) {
                console.log("Erreur:", data);
                alert(data.error || "Erreur lors de l'inscription");
                return;
            }

            console.log("Register successful!");
            alert("Inscription réussie !");

            // Retour au login
            registerview.style.display = 'none';
            registerContainer.style.display = 'none';
            loginview.style.display = 'flex';

        } catch (err) {
            console.error("Erreur fetch:", err);
            alert("Erreur réseau ou serveur");
        }
    });

    // LOGOUT
    logoutbtn.addEventListener('click', async (e) => {
        e.preventDefault();

        try {
            await fetch('/logout', { method: 'POST' });

            console.log('it done')
            loginview.style.display = 'flex';
            mainSection.style.display = 'none';
            logoutbtn.style.display = 'none';

            // ✅ Vider le container des posts
            const container = document.getElementById('posts-container');


        } catch (err) {
            console.error("Erreur logout:", err);
        }
    });
});


async function fetchapi() {
    const container = document.getElementById('posts-container');

    try {
        const response = await fetch("/api/posts");

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        const pageData = await response.json();
        displayCreationPost(pageData)
        console.log('here')
        displayPosts(pageData);
    } catch (error) {
        console.error('Error fetching posts:', error);
        container.innerHTML = '<div class="error">❌ Failed to load posts. Please try again later.</div>';
    }
}

/*
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
*/

function displayCreationPost(pageData) {
    const container = document.getElementById('post-creation')
    const containerDiv = document.createElement('div')
    const connectUsername = pageData.connectUserName
    container.innerHTML = ''
    // ✅ Générer les checkboxes de catégories
    let categoriesHTML = '<h3>Catégories :</h3>'
    pageData.categories.forEach(cat => {
        categoriesHTML += `
            <label class="category-label">
                <input type="checkbox" name="categories" value="${cat.id}"> 
                ${cat.category}
            </label>
        `
    })
    
    containerDiv.innerHTML = `
    <form class="create-post-form">
        <h1 class="title-create-post">Create New Post</h1>
        
        <input type="text" name="title" placeholder="Title" class="inputs title" maxlength="100" required>
        
        <textarea name="description" placeholder="Description" class="inputs description"
            maxlength="1000" required></textarea>
        
        <div class="categories">
            ${categoriesHTML}
        </div>
        
        <div class="upload-img">
            <label for="choose-file">Image (optionnel) :</label>
            <input type="file" name="choose-file" id="choose-file" class="choose-file" 
                   accept="image/jpeg,image/jpg,image/png,image/gif">
        </div>
        
        <button type="submit" class="btn btn-create-post">Create Post</button>
    </form>
    `
    
    container.appendChild(containerDiv)
    
    containerDiv.addEventListener('submit', async (e) => {
        e.preventDefault()
        
        const title = containerDiv.querySelector('[name="title"]').value.trim()
        const description = containerDiv.querySelector('[name="description"]').value.trim()
        const fileInput = containerDiv.querySelector('#choose-file')
        
        if (!title || !description) {
            alert("Le titre et la description sont requis")
            return
        }
        
        const selectedCategories = Array.from(
            containerDiv.querySelectorAll('[name="categories"]:checked')
        ).map(checkbox => checkbox.value)
        
        if (selectedCategories.length === 0) {
            alert("Veuillez sélectionner au moins une catégorie")
            return
        }
        
        const formData = new FormData()
        formData.append('title', title)
        formData.append('description', description)
        
        selectedCategories.forEach(catId => {
            formData.append('categories', catId)
        })
        
        if (fileInput.files.length > 0) {
            formData.append('choose-file', fileInput.files[0])
        }
        
        try {
            const res = await fetch('/createpost', {
                method: "POST",
                credentials: "same-origin",
                body: formData // Le navigateur ajoute automatiquement multipart/form-data
            });
            
            const data = await res.json()

            if (!res.ok) {
                console.log("Erreur:", data);
                alert(data.error || "Erreur lors de la création du post");
                return;
            }
            
            console.log("✅ Post créé:", data)
            alert("Post créé avec succès !")
            
            // ✅ Convertir les IDs de catégories en noms
            const categoryNames = data.categories.map(catId => {
                const cat = pageData.categories.find(c => c.id === parseInt(catId))
                return cat ? cat.category : "Unknown"
            })
            
            const newPost = {
                id: data.postId,
                title: data.title,
                description: data.description,
                imageUrl: data.imageUrl,
                userName: data.author,
                creationDate: data.date,
                categories: categoryNames
            }
            
            addNewPostToDOM(newPost)
            
           
            
        } catch (err) {
            console.error("Erreur fetch:", err);
            alert("Erreur réseau ou serveur");
        }
    })
}

function addNewPostToDOM(post) {
    const postContainer = document.getElementById('posts-container')
    
    const reactionStats = { likesCount: 0, dislikesCount: 0 }
    const userReaction = null
    const comments = []
    
    const postElement = createPostElement(post, reactionStats, userReaction, comments)
    
    postContainer.insertBefore(postElement, postContainer.firstChild)
}
//function 
/*
 const selectedCategories = Array.from(
            containerDiv.querySelectorAll('[name="categories"]:checked')
        ).map(checkbox => checkbox.value)
*/

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