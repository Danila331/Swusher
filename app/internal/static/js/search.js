// Обработка формы поиска
document.addEventListener('DOMContentLoaded', () => {
    const searchForm = document.querySelector('.search-form');
    if (!searchForm) return;

    searchForm.addEventListener('submit', (e) => {
        e.preventDefault();
        
        // Получаем значения из формы
        const searchQuery = document.getElementById('search-input').value.trim();
        const category = document.getElementById('category-select').value;
        const location = document.getElementById('location-select').value;
        const priceRange = document.getElementById('price-range').value;
        const sortBy = document.getElementById('sort-select').value;

        // Создаем объект с параметрами поиска
        const searchParams = new URLSearchParams();
        
        if (searchQuery) {
            searchParams.append('q', searchQuery);
        }
        if (category && category !== 'all') {
            searchParams.append('category', category);
        }
        if (location && location !== 'all') {
            searchParams.append('location', location);
        }
        if (priceRange) {
            searchParams.append('price', priceRange);
        }
        if (sortBy && sortBy !== 'default') {
            searchParams.append('sort', sortBy);
        }

        // Формируем URL для переадресации
        const catalogUrl = `/catalog.html${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
        
        // Переадресация на страницу каталога с параметрами
        window.location.href = catalogUrl;
    });

    // Обработка быстрого поиска (если есть)
    const quickSearchInput = document.querySelector('.quick-search-input');
    if (quickSearchInput) {
        quickSearchInput.addEventListener('keypress', (e) => {
            if (e.key === 'Enter') {
                e.preventDefault();
                const query = quickSearchInput.value.trim();
                if (query) {
                    window.location.href = `/catalog.html?q=${encodeURIComponent(query)}`;
                }
            }
        });
    }

    // Обработка фильтров в каталоге
    const filterForm = document.querySelector('.filter-form');
    if (filterForm) {
        filterForm.addEventListener('submit', (e) => {
            e.preventDefault();
            
            const formData = new FormData(filterForm);
            const searchParams = new URLSearchParams();

            // Собираем все активные фильтры
            for (let [key, value] of formData.entries()) {
                if (value && value !== 'all') {
                    searchParams.append(key, value);
                }
            }

            // Обновляем URL без перезагрузки страницы
            const newUrl = `${window.location.pathname}${searchParams.toString() ? '?' + searchParams.toString() : ''}`;
            window.history.pushState({}, '', newUrl);

            // Здесь можно добавить функцию для обновления результатов поиска
            updateSearchResults(searchParams);
        });
    }
});

// Функция обновления результатов поиска
async function updateSearchResults(params) {
    try {
        const response = await fetch(`/api/items/search?${params.toString()}`);
        if (!response.ok) {
            throw new Error('Ошибка при получении результатов поиска');
        }

        const data = await response.json();
        const resultsContainer = document.querySelector('.items-grid');
        
        if (resultsContainer) {
            // Очищаем текущие результаты
            resultsContainer.innerHTML = '';
            
            // Отображаем новые результаты
            data.items.forEach(item => {
                const itemCard = createItemCard(item);
                resultsContainer.appendChild(itemCard);
            });

            // Обновляем счетчик результатов
            const resultsCount = document.querySelector('.results-count');
            if (resultsCount) {
                resultsCount.textContent = `Найдено: ${data.total} ${declOfNum(data.total, ['предмет', 'предмета', 'предметов'])}`;
            }
        }
    } catch (error) {
        console.error('Ошибка при обновлении результатов:', error);
        // Показываем сообщение об ошибке пользователю
        const errorMessage = document.createElement('div');
        errorMessage.className = 'error-message';
        errorMessage.textContent = 'Произошла ошибка при загрузке результатов. Пожалуйста, попробуйте позже.';
        document.querySelector('.items-grid').appendChild(errorMessage);
    }
}

// Вспомогательная функция для склонения числительных
function declOfNum(number, titles) {
    const cases = [2, 0, 1, 1, 1, 2];
    return titles[
        (number % 100 > 4 && number % 100 < 20) ? 2 : cases[(number % 10 < 5) ? number % 10 : 5]
    ];
}

// Функция создания карточки товара
function createItemCard(item) {
    const card = document.createElement('div');
    card.className = 'item-card';
    card.innerHTML = `
        <img src="${item.image}" alt="${item.title}">
        <div class="item-info">
            <h3>${item.title}</h3>
            <div class="item-meta">
                <span class="price">${item.price} ₽/день</span>
                <span class="location">${item.location}</span>
            </div>
            <div class="item-details">
                <div class="item-rating">
                    <i class="fas fa-star"></i>
                    <span>${item.rating}</span>
                </div>
                <div class="item-deposit">
                    <i class="fas fa-shield-alt"></i>
                    <span>${item.deposit} ₽</span>
                </div>
            </div>
            <a href="item.html?id=${item.id}" class="btn btn-primary">Подробнее</a>
        </div>
    `;
    return card;
} 