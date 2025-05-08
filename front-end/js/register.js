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
document.querySelector('.auth-form').addEventListener('submit', (e) => {
    e.preventDefault();
    
    const name = document.getElementById('reg-name').value;
    const email = document.getElementById('reg-email').value;
    const password = document.getElementById('reg-password').value;
    const confirmPassword = document.getElementById('reg-password2').value;
    
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