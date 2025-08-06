// Mobile Menu Toggle
const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
const navLinks = document.querySelector('.nav-links');
const authButtons = document.querySelector('.auth-buttons');

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

// Load featured items
function loadFeaturedItems() {
    const itemsGrid = document.querySelector('.items-grid');
    if (!itemsGrid) return;

    // Example items data (in real app, this would come from an API)
    const items = [
        {
            id: 1,
            title: 'Sony A7III',
            image: 'https://technovybor.ru/wa-data/public/shop/img/sony-a7-iv-01-pa191048-acr-820x547.jpg',
            price: '1500 ₽/день',
            location: 'Москва',
            rating: 4.8,
            deposit: '15000 ₽'
        },
        {
            id: 2,
            title: 'DJI Mavic Pro',
            image: 'https://ekipirovka.expert/upload/iblock/d06/nhcjnvsfp31hlnn7c14lfbyxv7upq61s.png',
            price: '2000 ₽/день',
            location: 'Санкт-Петербург',
            rating: 4.9,
            deposit: '20000 ₽'
        },
        {
            id: 3,
            title: 'MacBook Pro M1',
            image: 'https://images.unsplash.com/photo-1517336714731-489689fd1ca8?auto=format&fit=crop&w=400&q=80',
            price: '1000 ₽/день',
            location: 'Москва',
            rating: 4.7,
            deposit: '50000 ₽'
        }
    ];

    items.forEach(item => {
        const itemCard = document.createElement('div');
        itemCard.className = 'item-card';
        itemCard.innerHTML = `
            <img src="${item.image}" alt="${item.title}">
            <div class="item-info">
                <h3>${item.title}</h3>
                <div class="item-meta">
                    <span class="price">${item.price}</span>
                    <span class="location">${item.location}</span>
                </div>
                <div class="item-details">
                    <div class="item-rating">
                        <i class="fas fa-star"></i>
                        <span>${item.rating}</span>
                    </div>
                    <div class="item-deposit">
                        <i class="fas fa-shield-alt"></i>
                        <span>${item.deposit}</span>
                    </div>
                </div>
                <a href="item.html?id=${item.id}" class="btn btn-primary">Подробнее</a>
            </div>
        `;
        itemsGrid.appendChild(itemCard);
    });
}

// Mobile menu toggle
function initMobileMenu() {
    const menuToggle = document.querySelector('.menu-toggle');
    const navLinks = document.querySelector('.nav-links');
    const authButtons = document.querySelector('.auth-buttons');

    if (menuToggle) {
        menuToggle.addEventListener('click', () => {
            navLinks.classList.toggle('active');
            authButtons.classList.toggle('active');
            menuToggle.classList.toggle('active');
        });
    }
}

// Smooth scroll for navigation links
document.querySelectorAll('a[href^="#"]').forEach(anchor => {
    anchor.addEventListener('click', function (e) {
        e.preventDefault();
        const target = document.querySelector(this.getAttribute('href'));
        if (target) {
            target.scrollIntoView({
                behavior: 'smooth',
                block: 'start'
            });
        }
    });
});

// Scroll Reveal Animation
function revealOnScroll() {
    const elements = document.querySelectorAll('.reveal');
    
    elements.forEach(element => {
        const elementTop = element.getBoundingClientRect().top;
        const windowHeight = window.innerHeight;
        
        if (elementTop < windowHeight - 100) {
            element.classList.add('active');
        }
    });
}

// Animate numbers
function animateNumber(element, target, duration = 2000) {
    const start = 0;
    const increment = target / (duration / 16);
    let current = start;

    const updateNumber = () => {
        current += increment;
        if (current < target) {
            element.textContent = Math.floor(current);
            requestAnimationFrame(updateNumber);
        } else {
            element.textContent = target;
        }
    };

    updateNumber();
}

// Initialize number animations
function initNumberAnimations() {
    const numberElements = document.querySelectorAll('.stat-number[data-count]');
    
    const observer = new IntersectionObserver((entries) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const element = entry.target;
                const target = parseInt(element.getAttribute('data-count'));
                animateNumber(element, target);
                observer.unobserve(element);
            }
        });
    }, {
        threshold: 0.5
    });

    numberElements.forEach(element => {
        observer.observe(element);
    });
}

// Initialize all animations
document.addEventListener('DOMContentLoaded', () => {
    loadFeaturedItems();
    initMobileMenu();
    initNumberAnimations();
    
    // Add reveal class to elements that should animate
    const elementsToReveal = document.querySelectorAll('.category-card, .item-card, .feature-card, .step, .mission-content, .values-grid, .team-grid, .partner-card');
    elementsToReveal.forEach(element => {
        element.classList.add('reveal');
    });

    // Initial check for elements in view
    revealOnScroll();

    // Add scroll event listener
    window.addEventListener('scroll', revealOnScroll);

    // Add scroll event listener for header
    let lastScroll = 0;
    const header = document.querySelector('.header');
    
    window.addEventListener('scroll', () => {
        const currentScroll = window.pageYOffset;
        
        if (currentScroll <= 0) {
            header.classList.remove('scroll-up');
            return;
        }
        
        if (currentScroll > lastScroll && !header.classList.contains('scroll-down')) {
            header.classList.remove('scroll-up');
            header.classList.add('scroll-down');
        } else if (currentScroll < lastScroll && header.classList.contains('scroll-down')) {
            header.classList.remove('scroll-down');
            header.classList.add('scroll-up');
        }
        lastScroll = currentScroll;
    });

    // Enhanced hover effects for cards
    document.querySelectorAll('.category-card, .item-card, .feature-card, .partner-card').forEach(card => {
        card.addEventListener('mousemove', (e) => {
            const rect = card.getBoundingClientRect();
            const x = e.clientX - rect.left;
            const y = e.clientY - rect.top;
            
            card.style.setProperty('--mouse-x', `${x}px`);
            card.style.setProperty('--mouse-y', `${y}px`);
        });
    });

    // Обработка поискового бокса
    const searchBox = document.querySelector('.search-box');
    if (!searchBox) return;

    searchBox.addEventListener('submit', (e) => {
        e.preventDefault();
        
        const searchQuery = searchBox.querySelector('input[name="tovar"]').value.trim();
        const city = searchBox.querySelector('input[name="city"]').value.trim();
        
        // Формируем URL для переадресации
        const catalogUrl = `/catalog?tovar=${encodeURIComponent(searchQuery)}&city=${encodeURIComponent(city)}`;
        
        // Переадресация на страницу каталога
        window.location.href = catalogUrl;
    });
}); 