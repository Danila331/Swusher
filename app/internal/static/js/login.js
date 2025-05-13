// Mobile Menu Toggle
const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
const navLinks = document.querySelector('.nav-links');

mobileMenuBtn.addEventListener('click', () => {
    navLinks.classList.toggle('active');
    authButtons.classList.toggle('active');
    mobileMenuBtn.classList.toggle('active');
});

// Mobile menu toggle
const mobileMenuBtn2 = document.querySelector('.mobile-menu-btn');
const mobileNav = document.getElementById('mobileNav');

if (mobileMenuBtn2 && mobileNav) {
    mobileMenuBtn2.addEventListener('click', () => {
        if (mobileNav.classList.contains('active')) {
            mobileNav.classList.add('closing');
            setTimeout(() => {
                mobileNav.classList.remove('active');
                mobileNav.classList.remove('closing');
                mobileMenuBtn2.classList.remove('active');
            }, 300);
        } else {
            mobileNav.classList.add('active');
            mobileMenuBtn2.classList.add('active');
        }
    });
    // Закрытие меню при клике на ссылку
    mobileNav.querySelectorAll('a').forEach(link => {
        link.addEventListener('click', () => {
            mobileNav.classList.add('closing');
            setTimeout(() => {
                mobileNav.classList.remove('active');
                mobileNav.classList.remove('closing');
                mobileMenuBtn2.classList.remove('active');
            }, 300);
        });
    });
}

// Handle register form submission
document.querySelector('.auth-form').addEventListener('submit', async (e) => {
    e.preventDefault();
    
    const email = document.getElementById('reg-email').value;
    const password = document.getElementById('reg-password').value;
    
    // Validate passwords match
    if (password !== confirmPassword) {
        alert('Пароли не совпадают');
        return;
    }

    try {
        const response = await fetch('http://localhost:8080/login/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email,
                password
            })
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.message || 'Ошибка при регистрации');
        }

        const userData = await response.json();
        
        // Save user data to session storage
        sessionStorage.setItem('currentUser', JSON.stringify(userData));
        
        // Show success message
        alert('Регистрация успешна! Добро пожаловать в ShareHub!');
        
        // Redirect to profile page
        window.location.href = 'profile.html';
    } catch (error) {
        alert(error.message);
    }
});

// Check if user is already logged in
document.addEventListener('DOMContentLoaded', () => {
    const currentUser = JSON.parse(localStorage.getItem('currentUser')) || 
                       JSON.parse(sessionStorage.getItem('currentUser'));
    
    if (currentUser) {
        authState.isAuthenticated = true;
        authState.currentUser = currentUser;
        window.location.href = 'catalog.html';
    }
});

// Toggle password visibility (глазик)
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOM Content Loaded');
    
    const toggleButtons = document.querySelectorAll('.toggle-password');
    console.log('Found toggle buttons:', toggleButtons.length);
    
    toggleButtons.forEach(function(toggle) {
        console.log('Adding click listener to:', toggle);
        
        toggle.addEventListener('click', function(e) {
            console.log('Toggle clicked');
            e.preventDefault();
            e.stopPropagation();
            
            const targetId = this.getAttribute('data-target');
            console.log('Target input ID:', targetId);
            
            const input = document.getElementById(targetId);
            console.log('Found input:', input);
            
            if (input) {
                if (input.type === 'password') {
                    console.log('Changing to text');
                    input.type = 'text';
                    this.querySelector('i').className = 'fa-solid fa-eye-slash';
                } else {
                    console.log('Changing to password');
                    input.type = 'password';
                    this.querySelector('i').className = 'fa-solid fa-eye';
                }
            }
        });
    });
});

// Обработчик входа через Яндекс
async function handleYandexLogin() {
    const clientId = 'YOUR_YANDEX_CLIENT_ID'; // Замените на ваш ID приложения Яндекс
    const redirectUri = encodeURIComponent(window.location.origin + '/register.html');
    const responseType = 'token';
    const display = 'popup';
    
    const yandexAuthUrl = `https://oauth.yandex.ru/authorize?response_type=${responseType}&client_id=${clientId}&redirect_uri=${redirectUri}&display=${display}`;
    
    window.open(yandexAuthUrl, 'yandexAuth', 'width=600,height=600');
}

// Обработчик входа через Google
async function handleGoogleSignIn(response) {
    try {
        const result = await fetch('http://localhost:3000/api/auth/google', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                credential: response.credential
            })
        });

        if (!result.ok) {
            throw new Error('Ошибка при авторизации через Google');
        }

        const userData = await result.json();
        sessionStorage.setItem('currentUser', JSON.stringify(userData));
        window.location.href = 'profile.html';
    } catch (error) {
        alert(error.message);
    }
}

// Обработка ответа от Яндекс OAuth
window.addEventListener('load', async function() {
    const hash = window.location.hash;
    if (hash) {
        const params = new URLSearchParams(hash.substring(1));
        const accessToken = params.get('access_token');
        if (accessToken) {
            try {
                const response = await fetch('http://localhost:3000/api/auth/yandex', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        access_token: accessToken
                    })
                });

                if (!response.ok) {
                    throw new Error('Ошибка при авторизации через Яндекс');
                }

                const userData = await response.json();
                sessionStorage.setItem('currentUser', JSON.stringify(userData));
                window.location.href = 'profile.html';
            } catch (error) {
                alert(error.message);
            }
        }
    }
});