// Auth state management
const authState = {
    isAuthenticated: false,
    currentUser: null
};

// Sample items data (same as in catalog.js)
const items = [
    {
        id: 1,
        title: 'Камера Sony A7III',
        price: '1500 ₽/день',
        location: 'Москва',
        image: 'images/items/camera.jpg',
        rating: 4.8,
        category: 'electronics',
        deposit: '30000 ₽',
        owner: {
            id: 1,
            name: 'Александр П.',
            rating: 4.9,
            verified: true
        },
        description: 'Полнокадровая беззеркальная камера Sony A7III с матрицей 24.2 Мп, системой автофокусировки с 693 точками и 5-осевой стабилизацией изображения. Идеально подходит для профессиональной фото- и видеосъемки.',
        specs: [
            'Матрица: 24.2 Мп, полный кадр',
            'Стабилизация: 5-осевая',
            'Автофокус: 693 точки',
            'Видео: 4K 30p',
            'Батарея: NP-FZ100'
        ],
        rentalTerms: [
            'Минимальный срок аренды: 1 день',
            'Максимальный срок аренды: 30 дней',
            'Требуется залог: 30000 ₽',
            'Страховка: опционально'
        ],
        reviews: [
            {
                id: 1,
                user: {
                    name: 'Мария К.',
                    avatar: 'images/users/user1.jpg'
                },
                rating: 5.0,
                text: 'Отличная камера, всё работает идеально. Владелец очень ответственный и приятный в общении. Рекомендую!',
                date: '12 марта 2024'
            }
        ]
    }
];

// Check authentication on page load
document.addEventListener('DOMContentLoaded', () => {
    const currentUser = JSON.parse(localStorage.getItem('currentUser')) || 
                       JSON.parse(sessionStorage.getItem('currentUser'));
    
    if (currentUser) {
        authState.isAuthenticated = true;
        authState.currentUser = currentUser;
        updateAuthUI();
    } else {
        window.location.href = 'auth.html';
    }

    // Get item ID from URL
    const urlParams = new URLSearchParams(window.location.search);
    const itemId = parseInt(urlParams.get('id'));
    
    // Find and display item
    const item = items.find(i => i.id === itemId);
    if (item) {
        displayItem(item);
    } else {
        window.location.href = 'catalog.html';
    }
});

// Update UI based on auth state
function updateAuthUI() {
    const authButtons = document.querySelector('.auth-buttons');
    if (authState.isAuthenticated) {
        authButtons.innerHTML = `
            <div class="user-menu">
                <img src="images/users/${authState.currentUser.id}.jpg" alt="User Avatar" class="user-avatar">
                <span class="user-name">${authState.currentUser.name}</span>
                <button class="btn btn-outline" id="logout">Выйти</button>
            </div>
        `;
        
        document.getElementById('logout').addEventListener('click', () => {
            localStorage.removeItem('currentUser');
            sessionStorage.removeItem('currentUser');
            window.location.href = 'auth.html';
        });
    }
}

// Display item details
function displayItem(item) {
    // Update main image and thumbnails
    document.querySelector('.main-image img').src = item.image;
    document.querySelector('.main-image img').alt = item.title;
    
    // Update item info
    document.querySelector('.item-info h1').textContent = item.title;
    document.querySelector('.rating span').textContent = item.rating;
    document.querySelector('.location span').textContent = item.location;
    document.querySelector('.daily-price .price').textContent = item.price;
    document.querySelector('.deposit span').textContent = `Залог: ${item.deposit}`;
    
    // Update description
    document.querySelector('.item-description p').textContent = item.description;
    
    // Update specs
    const specsList = document.querySelector('.specs-list');
    specsList.innerHTML = item.specs.map(spec => `<li>${spec}</li>`).join('');
    
    // Update rental terms
    const rentalTermsList = document.querySelector('.rental-terms');
    rentalTermsList.innerHTML = item.rentalTerms.map(term => `<li>${term}</li>`).join('');
    
    // Update owner info
    document.querySelector('.owner-details h3').textContent = item.owner.name;
    document.querySelector('.owner-rating span').textContent = item.owner.rating;
    
    // Update reviews
    const reviewsList = document.querySelector('.reviews-list');
    reviewsList.innerHTML = item.reviews.map(review => `
        <div class="review-card">
            <div class="reviewer-info">
                <img src="${review.user.avatar}" alt="${review.user.name}" class="reviewer-avatar">
                <div class="reviewer-details">
                    <h4>${review.user.name}</h4>
                    <div class="rating">
                        <i class="fas fa-star"></i>
                        <span>${review.rating}</span>
                    </div>
                </div>
            </div>
            <p class="review-text">${review.text}</p>
            <span class="review-date">${review.date}</span>
        </div>
    `).join('');
}

// Handle booking form
document.querySelector('.booking-form').addEventListener('submit', (e) => {
    e.preventDefault();
    
    const startDate = document.querySelector('.start-date').value;
    const endDate = document.querySelector('.end-date').value;
    
    if (!startDate || !endDate) {
        alert('Пожалуйста, выберите даты аренды');
        return;
    }
    
    // Calculate total price
    const pricePerDay = parseInt(document.querySelector('.daily-price .price').textContent);
    const days = Math.ceil((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24));
    const totalPrice = pricePerDay * days;
    
    // Update total price display
    document.querySelector('.total-price .price').textContent = `${totalPrice} ₽`;
    
    // Here you would typically send the booking request to the server
    alert('Бронирование успешно создано!');
});

// Handle date range changes
document.querySelectorAll('.date-range input').forEach(input => {
    input.addEventListener('change', () => {
        const startDate = document.querySelector('.start-date').value;
        const endDate = document.querySelector('.end-date').value;
        
        if (startDate && endDate) {
            const pricePerDay = parseInt(document.querySelector('.daily-price .price').textContent);
            const days = Math.ceil((new Date(endDate) - new Date(startDate)) / (1000 * 60 * 60 * 24));
            const totalPrice = pricePerDay * days;
            
            document.querySelector('.total-price .price').textContent = `${totalPrice} ₽`;
        }
    });
}); 