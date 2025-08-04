// Item page functionality
document.addEventListener('DOMContentLoaded', function() {
    // Initialize mobile menu
    initializeMobileMenu();
    
    // Initialize image gallery
    initializeImageGallery();
    
    // Initialize booking form
    initializeBookingForm();
    
    // Initialize date validation
    initializeDateValidation();
    
    // Initialize lazy loading
    initializeLazyLoading();
    
    // Add smooth scrolling for better UX
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
});

// Mobile menu functionality
function initializeMobileMenu() {
    const mobileMenuBtn = document.querySelector('.mobile-menu-btn');
    const mobileNav = document.getElementById('mobileNav');

    if (mobileMenuBtn && mobileNav) {
        mobileMenuBtn.addEventListener('click', function() {
            mobileNav.classList.toggle('active');
            mobileMenuBtn.classList.toggle('active');
        });

        // Close mobile menu when clicking outside
        document.addEventListener('click', function(e) {
            if (!mobileMenuBtn.contains(e.target) && !mobileNav.contains(e.target)) {
                mobileNav.classList.remove('active');
                mobileMenuBtn.classList.remove('active');
            }
        });
    }
}

// Image gallery functionality
function initializeImageGallery() {
    const thumbnails = document.querySelectorAll('.thumbnail');
    const mainImage = document.getElementById('mainImage');

    thumbnails.forEach(thumbnail => {
        thumbnail.addEventListener('click', function() {
            // Remove active class from all thumbnails
            thumbnails.forEach(t => t.classList.remove('active'));
            
            // Add active class to clicked thumbnail
            this.classList.add('active');
            
            // Update main image
            if (mainImage) {
                mainImage.src = this.src;
                mainImage.alt = this.alt;
            }
        });
    });

    // Set first thumbnail as active by default
    if (thumbnails.length > 0) {
        thumbnails[0].classList.add('active');
    }
}

// Booking form functionality
function initializeBookingForm() {
    const bookingForm = document.getElementById('bookingForm');
    const startDate = document.getElementById('startDate');
    const endDate = document.getElementById('endDate');
    const totalPrice = document.getElementById('totalPrice');

    if (bookingForm && startDate && endDate && totalPrice) {
        // Get price per day from the page
        const priceElement = document.querySelector('.daily-price .price');
        const pricePerDay = priceElement ? parseInt(priceElement.textContent) : 0;

        // Calculate total when dates change
        function calculateTotal() {
            if (startDate.value && endDate.value) {
                const start = new Date(startDate.value);
                const end = new Date(endDate.value);
                const days = Math.ceil((end - start) / (1000 * 60 * 60 * 24));
                
                if (days > 0) {
                    totalPrice.textContent = (days * pricePerDay) + ' ₽';
                } else {
                    totalPrice.textContent = '0 ₽';
                }
            }
        }

        startDate.addEventListener('change', calculateTotal);
        endDate.addEventListener('change', calculateTotal);
        
        // Make calculateTotal function global for calendar integration
        window.calculateTotal = calculateTotal;

        // Handle form submission
        bookingForm.addEventListener('submit', function(e) {
            e.preventDefault();
            
            if (!startDate.value || !endDate.value) {
                showNotification('Пожалуйста, выберите даты аренды', 'error');
                return;
            }

            const start = new Date(startDate.value);
            const end = new Date(endDate.value);
            const days = Math.ceil((end - start) / (1000 * 60 * 60 * 24));
            
            if (days <= 0) {
                showNotification('Дата окончания должна быть позже даты начала', 'error');
                return;
            }

            // Show loading state
            const submitBtn = bookingForm.querySelector('button[type="submit"]');
            const originalText = submitBtn.textContent;
            submitBtn.textContent = 'Обработка...';
            submitBtn.disabled = true;

            // Simulate API call
            setTimeout(() => {
                showNotification('Бронирование успешно создано!', 'success');
                submitBtn.textContent = originalText;
                submitBtn.disabled = false;
                
                // Reset form
                bookingForm.reset();
                totalPrice.textContent = '0 ₽';
            }, 2000);
        });
    }
}

// Date validation
function initializeDateValidation() {
    const startDate = document.getElementById('startDate');
    const endDate = document.getElementById('endDate');

    if (startDate && endDate) {
        // Set minimum date to today
        const today = new Date().toISOString().split('T')[0];
        startDate.min = today;
        endDate.min = today;

        // Update end date minimum when start date changes
        startDate.addEventListener('change', function() {
            endDate.min = this.value;
            if (endDate.value && endDate.value < this.value) {
                endDate.value = this.value;
            }
        });
    }
}

// Change main image function (for external calls)
function changeMainImage(src) {
    const mainImage = document.getElementById('mainImage');
    if (mainImage) {
        mainImage.src = src;
    }
}

// Open chat with owner
function openChat(ownerId) {
    // Check if user is authenticated
    const authButtons = document.querySelector('.auth-buttons');
    if (authButtons.querySelector('.user-menu')) {
        // User is authenticated, redirect to chat
        window.location.href = `/chat?user=${ownerId}`;
    } else {
        // User is not authenticated, redirect to login
        window.location.href = '/login/';
    }
}

// Show all reviews
function showAllReviews() {
    const itemId = getItemIdFromUrl();
    if (itemId) {
        window.location.href = `/reviews?item=${itemId}`;
    }
}

// Get item ID from URL
function getItemIdFromUrl() {
    const urlParams = new URLSearchParams(window.location.search);
    return urlParams.get('id');
}

// Notification system
function showNotification(message, type = 'info') {
    // Remove existing notifications
    const existingNotifications = document.querySelectorAll('.notification');
    existingNotifications.forEach(notification => notification.remove());

    // Create notification element
    const notification = document.createElement('div');
    notification.className = `notification notification-${type}`;
    notification.innerHTML = `
        <div class="notification-content">
            <span class="notification-message">${message}</span>
            <button class="notification-close" onclick="this.parentElement.parentElement.remove()">
                <i class="fas fa-times"></i>
            </button>
        </div>
    `;

    // Add styles
    notification.style.cssText = `
        position: fixed;
        top: 20px;
        right: 20px;
        background: ${type === 'success' ? '#4CAF50' : type === 'error' ? '#f44336' : '#2196F3'};
        color: white;
        padding: 15px 20px;
        border-radius: 5px;
        box-shadow: 0 4px 12px rgba(0,0,0,0.15);
        z-index: 10000;
        max-width: 400px;
        animation: slideInRight 0.3s ease;
    `;

    // Add to page
    document.body.appendChild(notification);

    // Auto remove after 5 seconds
    setTimeout(() => {
        if (notification.parentElement) {
            notification.remove();
        }
    }, 5000);
}

// Add CSS animation for notifications
const style = document.createElement('style');
style.textContent = `
    @keyframes slideInRight {
        from {
            transform: translateX(100%);
            opacity: 0;
        }
        to {
            transform: translateX(0);
            opacity: 1;
        }
    }
    
    .notification-content {
        display: flex;
        align-items: center;
        justify-content: space-between;
    }
    
    .notification-close {
        background: none;
        border: none;
        color: white;
        cursor: pointer;
        margin-left: 10px;
        padding: 0;
        font-size: 16px;
    }
    
    .notification-close:hover {
        opacity: 0.8;
    }
`;
document.head.appendChild(style);

// Lazy loading for images
function initializeLazyLoading() {
    const images = document.querySelectorAll('img[data-src]');
    
    const imageObserver = new IntersectionObserver((entries, observer) => {
        entries.forEach(entry => {
            if (entry.isIntersecting) {
                const img = entry.target;
                img.src = img.dataset.src;
                img.removeAttribute('data-src');
                observer.unobserve(img);
            }
        });
    });

    images.forEach(img => imageObserver.observe(img));
}

// Initialize lazy loading when DOM is ready
if (document.readyState === 'loading') {
    document.addEventListener('DOMContentLoaded', initializeLazyLoading);
} else {
    initializeLazyLoading();
} 