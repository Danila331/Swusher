// Auth state management
const authState = {
    isAuthenticated: false,
    currentUser: null
};

// DOM Elements
const loginForm = document.getElementById('login-form');
const registerForm = document.getElementById('register-form');
const authButtons = document.querySelectorAll('.auth-buttons button');
const forms = document.querySelectorAll('.auth-form');

// Switch between login and register forms
authButtons.forEach(button => {
    button.addEventListener('click', () => {
        const formType = button.dataset.form;
        
        // Update button styles
        authButtons.forEach(btn => {
            btn.classList.remove('active', 'btn-primary');
            btn.classList.add('btn-outline');
        });
        button.classList.remove('btn-outline');
        button.classList.add('active', 'btn-primary');
        
        // Show selected form
        forms.forEach(form => {
            form.classList.remove('active');
            if (form.id === `${formType}-form`) {
                form.classList.add('active');
            }
        });
    });
});

// Handle login form submission
document.getElementById('login').addEventListener('submit', (e) => {
    e.preventDefault();
    
    const email = document.getElementById('login-email').value;
    const password = document.getElementById('login-password').value;
    const rememberMe = document.getElementById('remember-me').checked;
    
    // Get users from localStorage
    const users = JSON.parse(localStorage.getItem('users')) || [];
    const user = users.find(u => u.email === email && u.password === password);
    
    if (user) {
        // Update auth state
        authState.isAuthenticated = true;
        authState.currentUser = user;
        
        // Save to localStorage if remember me is checked
        if (rememberMe) {
            localStorage.setItem('currentUser', JSON.stringify(user));
        } else {
            sessionStorage.setItem('currentUser', JSON.stringify(user));
        }
        
        // Redirect to catalog
        window.location.href = 'catalog.html';
    } else {
        alert('Неверный email или пароль');
    }
});

// Handle register form submission
document.getElementById('register').addEventListener('submit', (e) => {
    e.preventDefault();
    
    const name = document.getElementById('register-name').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;
    const confirmPassword = document.getElementById('register-confirm-password').value;
    
    // Validate passwords match
    if (password !== confirmPassword) {
        alert('Пароли не совпадают');
        return;
    }
    
    // Get existing users
    const users = JSON.parse(localStorage.getItem('users')) || [];
    
    // Check if email already exists
    if (users.some(user => user.email === email)) {
        alert('Пользователь с таким email уже существует');
        return;
    }
    
    // Create new user
    const newUser = {
        id: Date.now(),
        name,
        email,
        password,
        rating: 0,
        reviews: 0,
        verified: false
    };
    
    // Save user
    users.push(newUser);
    localStorage.setItem('users', JSON.stringify(users));
    
    // Update auth state
    authState.isAuthenticated = true;
    authState.currentUser = newUser;
    
    // Save to session storage
    sessionStorage.setItem('currentUser', JSON.stringify(newUser));
    
    // Redirect to catalog
    window.location.href = 'catalog.html';
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

// Social login handlers
document.querySelector('.btn-social.vk').addEventListener('click', () => {
    alert('Вход через ВКонтакте будет реализован позже');
});

document.querySelector('.btn-social.google').addEventListener('click', () => {
    alert('Вход через Google будет реализован позже');
}); 