import { fetchapi } from "./fetchJson";

document.addEventListener("DOMContentLoaded", () => {
    
    const logginform = document.getElementById('login-form')
    logginform.addEventListener('submit', async (e) => {
        e.preventDefault()
        const userName = document.getElementById('username').value
        const password = document.getElementById('password').value
        if (!userName || !password) return

        const params = new URLSearchParams();
        params.append("username", userName);
        params.append("password", password);

        try {
            const res = await fetch("/loginAuth", {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            })

            const data = await res.json()

            if (!res.ok) {
                console.log("Erreur:", data);
                alert(data.error || "Erreur lors de login");
                return;
            }

            // Afficher les données dans la console
            console.log("Login successful!");
            console.log("Data:", data);
            console.log("User ID:", data.userId);
            console.log("Username:", data.username);

            if (res.ok) {
                alert("Connexion réussie !");

                // Masquer login/register
                loginview.style.display = 'none';
                registerview.style.display = 'none';
                registerContainer.style.display = 'none';

                // Afficher la home page
                mainSection.style.display = 'block';
                logoutbtn.style.display = 'flex';

                // Appeler fetch pour charger les posts
                fetchapi();
            }


            // Attendre 1 seconde avant de recharger pour voir la console

        } catch (err) {
            console.error("Erreur fetch:", err);
            alert("Erreur réseau ou serveur");
        }
    })
    const registerContainer = document.getElementById('register-container')
    const registerview = document.getElementById('register-view')
    const logoutbtn = document.getElementById('btn-logout')
    const loginview = document.getElementById('login-view')
    const mainSection = document.querySelector('.main-section')
    const islogged =document.getElementById('app-container').dataset.loggedin === 'true'

    // ÉTAT INITIAL
    if (islogged) {
        registerview.style.display = 'none'
        registerContainer.style.display = 'none'
        loginview.style.display = 'none'
        mainSection.style.display = 'flex'
        mainSection.style.display = 'block'
        logoutbtn.style.display = 'flex'
        fetchapi()
    } else {
        registerview.style.display = 'none'
        registerContainer.style.display = 'none'
        loginview.style.display = 'flex'
        mainSection.style.display = 'none'
        logoutbtn.style.display = 'none'
    }

    const buttonregister = document.querySelector('.link-register')

    // When clicking the register link, just show the form
    buttonregister.addEventListener('click', (e) => {
        e.preventDefault()
        console.log('Showing register form')
        registerview.style.display = 'flex'
        registerContainer.style.display = 'flex'
        loginview.style.display = 'none'
        mainSection.style.display = 'none'
    })

    // Handle form submission separately
    const registerForm = document.getElementById('register-form') // Add ID to your form in HTML

    registerForm.addEventListener('submit', async (e) => {
        e.preventDefault()

        // Get values when form is submitted
        const nickname = document.getElementById('nickname').value
        const firstname = document.getElementById('firstname').value
        const lastname = document.getElementById('lastname').value
        const age = document.getElementById('age').value
        const gender = document.getElementById('gender').value
        const email = document.getElementById('email').value
        const passwordre = document.getElementById('passwordre').value

        console.log('Form values:', {
            nickname,
            firstname,
            lastname,
            age,
            gender,
            email,
            passwordre
        })

        const params = new URLSearchParams()
        params.append("nickname", nickname)
        params.append("firstname", firstname)
        params.append("lastname", lastname)
        params.append("age", age)
        params.append("gender", gender)
        params.append("email", email)
        params.append("passwordre", passwordre)

        try {
            const res = await fetch('/registerAuth', {
                method: "POST",
                credentials: "same-origin",
                headers: { "Content-Type": "application/x-www-form-urlencoded" },
                body: params.toString()
            })
            const data = await res.json()

            if (!res.ok) {
                console.log("Erreur:", data)
                alert(data.error || "Erreur lors de l'inscription")
                return
            }

            // Afficher les données dans la console
            console.log("Register successful!")
            console.log("Data:", data)

            alert("Inscription réussie !")

        } catch (err) {
            console.error("Erreur fetch:", err)
            alert("Erreur réseau ou serveur")
        }
    })

    // LOGOUT
    logoutbtn.addEventListener('click', async (e) => {
        e.preventDefault()

        // vrai logout backend
        await fetch('/logout', { method: 'POST' })

        // mise à jour SPA
        loginview.style.display = 'flex'
        mainSection.style.display = 'none'
        logoutbtn.style.display = 'none'
    })
});