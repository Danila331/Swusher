

// Sample items data
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
        }
    },
    {
        id: 2,
        title: 'Велосипед Trek',
        price: '800 ₽/день',
        location: 'Санкт-Петербург',
        image: 'images/items/bicycle.jpg',
        rating: 4.9,
        category: 'sports',
        deposit: '15000 ₽',
        owner: {
            id: 2,
            name: 'Мария К.',
            rating: 4.7,
            verified: true
        }
    },
    {
        id: 3,
        title: 'Дрель Bosch',
        price: '500 ₽/день',
        location: 'Екатеринбург',
        image: 'images/items/drill.jpg',
        rating: 4.7,
        category: 'tools',
        deposit: '5000 ₽',
        owner: {
            id: 3,
            name: 'Иван С.',
            rating: 4.8,
            verified: true
        }
    },
    {
        id: 4,
        title: 'DJI Mavic Air 2',
        price: '2000 ₽/день',
        location: 'Москва',
        image: 'images/items/drone.jpg',
        rating: 4.9,
        category: 'electronics',
        deposit: '40000 ₽'
    },
    {
        id: 5,
        title: 'Палатка Coleman',
        price: '600 ₽/день',
        location: 'Казань',
        image: 'images/items/tent.jpg',
        rating: 4.6,
        category: 'sports',
        deposit: '8000 ₽'
    },
    {
        id: 6,
        title: 'Проектор Epson',
        price: '1200 ₽/день',
        location: 'Новосибирск',
        image: 'images/items/projector.jpg',
        rating: 4.5,
        category: 'electronics',
        deposit: '20000 ₽'
    }
];

// Function to create item card
function createItemCard(item) {
    return `
        <div class="item-card" data-category="${item.category}">
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
                        <span>Залог: ${item.deposit}</span>
                    </div>
                </div>
                <button class="btn btn-primary" onclick="window.location.href='item.html?id=${item.id}'">Подробнее</button>
            </div>
        </div>
    `;
}

// Function to filter items
function filterItems() {
    const selectedCategories = Array.from(document.querySelectorAll('input[name="category"]:checked'))
        .map(checkbox => checkbox.value);
    
    const minPrice = document.querySelector('.price-min').value;
    const maxPrice = document.querySelector('.price-max').value;
    const location = document.querySelector('.location-input').value.toLowerCase();
    const date = document.querySelector('.date-input').value;
    const minRating = document.querySelector('.rating-filter .stars').children.length;

    const filteredItems = items.filter(item => {
        const matchesCategory = selectedCategories.length === 0 || selectedCategories.includes(item.category);
        const matchesPrice = (!minPrice || parseInt(item.price) >= parseInt(minPrice)) &&
                           (!maxPrice || parseInt(item.price) <= parseInt(maxPrice));
        const matchesLocation = !location || item.location.toLowerCase().includes(location);
        const matchesRating = item.rating >= minRating;

        return matchesCategory && matchesPrice && matchesLocation && matchesRating;
    });

    displayItems(filteredItems);
}

// Function to display items
function displayItems(itemsToShow) {
    const itemsGrid = document.querySelector('.items-grid');
    itemsGrid.innerHTML = itemsToShow.map(item => createItemCard(item)).join('');
    
    // Update items count
    document.querySelector('.items-header h2').textContent = `Найдено ${itemsToShow.length} вещей`;
}

// Function to sort items
function sortItems(sortBy) {
    const itemsGrid = document.querySelector('.items-grid');
    const items = Array.from(itemsGrid.children);

    items.sort((a, b) => {
        const priceA = parseInt(a.querySelector('.price').textContent);
        const priceB = parseInt(b.querySelector('.price').textContent);
        const ratingA = parseFloat(a.querySelector('.item-rating span').textContent);
        const ratingB = parseFloat(b.querySelector('.item-rating span').textContent);

        switch (sortBy) {
            case 'price-asc':
                return priceA - priceB;
            case 'price-desc':
                return priceB - priceA;
            case 'rating':
                return ratingB - ratingA;
            default:
                return 0;
        }
    });

    items.forEach(item => itemsGrid.appendChild(item));
}

// Event Listeners
document.addEventListener('DOMContentLoaded', () => {
    // Initial display
    displayItems(items);

    // Filter events
    document.querySelectorAll('input[name="category"]').forEach(checkbox => {
        checkbox.addEventListener('change', filterItems);
    });

    document.querySelector('.price-slider').addEventListener('input', (e) => {
        document.querySelector('.price-max').value = e.target.value;
        filterItems();
    });

    document.querySelector('.location-input').addEventListener('input', filterItems);
    document.querySelector('.date-input').addEventListener('change', filterItems);

    // Rating filter
    const stars = document.querySelectorAll('.rating-filter .stars i');
    stars.forEach((star, index) => {
        star.addEventListener('click', () => {
            stars.forEach((s, i) => {
                s.style.color = i <= index ? '#F59E0B' : '#E5E7EB';
            });
            filterItems();
        });
    });

    // Sort events
    document.querySelector('.sort-select').addEventListener('change', (e) => {
        sortItems(e.target.value);
    });

    // Pagination
    document.querySelectorAll('.page-numbers span').forEach(page => {
        page.addEventListener('click', () => {
            document.querySelector('.page-numbers span.active').classList.remove('active');
            page.classList.add('active');
            // Here you would typically load the corresponding page of items
        });
    });
}); 