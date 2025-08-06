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

    const name = document.getElementById('reg-name').value;
    const lastName = document.getElementById('reg-lastname').value;
    const email = document.getElementById('reg-email').value;
    const phone = document.getElementById('reg-phone').value;
    const password = document.getElementById('reg-password').value;
    const confirmPassword = document.getElementById('reg-password2').value;

    if (password !== confirmPassword) {
        alert('Пароли не совпадают');
        return;
    }

    try {
        const response = await fetch('/register/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                name,
                lastName,
                email,
                phone,
                password
            })
        });

        let data;
        try {
            data = await response.json();
        } catch (jsonErr) {
            // Если ответ не JSON
            console.error('Ошибка парсинга ответа:', jsonErr);
            alert('Неизвестная ошибка сервера. Попробуйте позже.');
            return;
        }
        if (response.status === 409) {
            // Пользователь уже существует
            alert('Пользователь с таким email уже зарегистрирован.');
            return;
        }
        if (!response.ok) {
            // Показываем сообщение от сервера, если оно есть
            const errorMsg = data.error || data.message || 'Ошибка при регистрации';
            alert(errorMsg);
            console.error('Ошибка регистрации:', errorMsg, data);
            return;
        }

        sessionStorage.setItem('currentUser', JSON.stringify(data));
        window.location.href = '/profile-saler/';
    } catch (error) {
        // Ошибка сети или fetch
        alert('Ошибка соединения с сервером. Проверьте интернет или попробуйте позже.');
        console.error('Ошибка fetch:', error);
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
    const clientId = '7326f484a76e4cc58292b13494195c18'; // Замените на ваш ID приложения Яндекс
    const redirectUri = 'http://localhost:8080/register/';
    const responseType = 'token';
    
    const yandexAuthUrl = `https://oauth.yandex.ru/authorize?response_type=${responseType}&client_id=${clientId}&redirect_uri=${redirectUri}`;
    
    window.location.href = yandexAuthUrl;
}

// Обработка ответа от Яндекс OAuth
window.addEventListener('load', async function() {
    if (window.location.pathname !== '/register/') return;
    const hash = window.location.hash;
    if (hash) {
        const params = new URLSearchParams(hash.substring(1));
        const accessToken = params.get('access_token');
        if (accessToken) {
            try {
                const response = await fetch('/register/yandex', {
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
                window.location.href = '/profile-saler/';
            } catch (error) {
                alert(error.message);
            }
        }
    }
});

// Обработчик входа через Google
async function handleGoogleSignIn(response) {
    try {
        const result = await fetch('/register/google', {
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

function togglePassword(inputId, btn) {
    const input = document.getElementById(inputId);
    const icon = btn.querySelector('i');
    if (input.type === 'password') {
        input.type = 'text';
        icon.className = 'fa-solid fa-eye-slash';
    } else {
        input.type = 'password';
        icon.className = 'fa-solid fa-eye';
    }
}