// Add Advertisement Form JavaScript
let uploadedPhotos = [];
let autoSaveTimer = null;
let yandexMap = null;
let mapPlacemark = null;

document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('advertisementForm');
    const photoInput = document.getElementById('photoInput');
    const photoUploadArea = document.getElementById('photoUploadArea');
    const photoPreviewGrid = document.getElementById('photoPreviewGrid');
    
    // Initialize form
    initializeForm();
    
    // Form submission
    form.addEventListener('submit', handleFormSubmit);
    
    // Photo upload functionality
    photoUploadArea.addEventListener('click', () => photoInput.click());
    photoInput.addEventListener('change', handlePhotoUpload);
    
    // Drag and drop for photos
    photoUploadArea.addEventListener('dragover', handleDragOver);
    photoUploadArea.addEventListener('dragleave', handleDragLeave);
    photoUploadArea.addEventListener('drop', handleDrop);
    
    // Auto-save functionality
    form.addEventListener('input', debounce(autoSave, 2000));
    
    // Price calculation helpers
    const costPerDayInput = document.getElementById('cost_per_day');
    const costPerWeekInput = document.getElementById('cost_per_week');
    const costPerMonthInput = document.getElementById('cost_per_month');
    
    costPerDayInput.addEventListener('input', calculatePrices);
    costPerWeekInput.addEventListener('input', calculatePrices);
    costPerMonthInput.addEventListener('input', calculatePrices);
    
    // Location autocomplete and map integration
    initializeLocationAutocomplete();
    initializeYandexMap();
    
    // Form validation
    initializeFormValidation();
    
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
    
    // Form progress tracking
    document.addEventListener('input', function() {
        updateFormProgress();
    });
});

// Initialize form
function initializeForm() {
    // Load draft if exists
    loadDraft();
    
    // Initialize progress indicator
    updateFormProgress();
    
    // Show auto-save indicator
    showAutoSaveIndicator('ready', 'Форма готова к заполнению');
}

// Handle form submission
async function handleFormSubmit(event) {
    event.preventDefault();
    
    if (!validateForm()) {
        showMessage('Пожалуйста, исправьте ошибки в форме', 'error');
        return;
    }
    
    const formData = new FormData(event.target);
    const submitButton = event.target.querySelector('button[type="submit"]');
    const originalText = submitButton.innerHTML;
    
    try {
        // Show loading state
        submitButton.innerHTML = '<div class="loading-spinner"></div> Отправка...';
        submitButton.disabled = true;
        
        // Add photos to form data
        uploadedPhotos.forEach((photo, index) => {
            formData.append(`photo_${index}`, photo.file);
        });
        
        // Send data to server
        const response = await fetch('/advertisement/add', {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error('Ошибка при отправке данных');
        }
        
        const result = await response.json();
        
        // Show success message
        showMessage('Объявление успешно создано!', 'success');
        
        // Redirect to the new advertisement page
        window.location.href = `/profile-saler/items`;
        
    } catch (error) {
        console.error('Error submitting form:', error);
        showMessage('Произошла ошибка при создании объявления. Попробуйте еще раз.', 'error');
    } finally {
        // Restore button state
        submitButton.innerHTML = originalText;
        submitButton.disabled = false;
    }
}

// Handle photo upload
function handlePhotoUpload(event) {
    const files = Array.from(event.target.files);
    processPhotoFiles(files);
    // Clear input to allow selecting the same file again
    event.target.value = '';
}

// Handle drag and drop
function handleDragOver(event) {
    event.preventDefault();
    event.stopPropagation();
    const photoUploadArea = document.getElementById('photoUploadArea');
    photoUploadArea.classList.add('dragover');
}

function handleDragLeave(event) {
    event.preventDefault();
    event.stopPropagation();
    const photoUploadArea = document.getElementById('photoUploadArea');
    photoUploadArea.classList.remove('dragover');
}

function handleDrop(event) {
    event.preventDefault();
    event.stopPropagation();
    const photoUploadArea = document.getElementById('photoUploadArea');
    photoUploadArea.classList.remove('dragover');
    
    const files = Array.from(event.dataTransfer.files);
    const imageFiles = files.filter(file => file.type.startsWith('image/'));
    
    if (imageFiles.length > 0) {
        processPhotoFiles(imageFiles);
    }
}

// Process photo files
function processPhotoFiles(files) {
    const maxFiles = 10;
    const maxSize = 5 * 1024 * 1024; // 5MB
    
    if (uploadedPhotos.length + files.length > maxFiles) {
        showMessage(`Максимальное количество фотографий: ${maxFiles}`, 'error');
        return;
    }
    
    files.forEach(file => {
        if (file.size > maxSize) {
            showMessage(`Файл ${file.name} слишком большой. Максимальный размер: 5MB`, 'error');
            return;
        }
        
        if (!file.type.startsWith('image/')) {
            showMessage(`Файл ${file.name} не является изображением`, 'error');
            return;
        }
        
        const reader = new FileReader();
        reader.onload = function(e) {
            const photoData = {
                file: file,
                preview: e.target.result,
                id: Date.now() + Math.random()
            };
            
            uploadedPhotos.push(photoData);
            updatePhotoPreview();
        };
        reader.readAsDataURL(file);
    });
}

// Update photo preview
function updatePhotoPreview() {
    const photoPreviewGrid = document.getElementById('photoPreviewGrid');
    photoPreviewGrid.innerHTML = '';
    
    uploadedPhotos.forEach((photo, index) => {
        const photoItem = document.createElement('div');
        photoItem.className = 'photo-preview-item';
        photoItem.innerHTML = `
            <img src="${photo.preview}" alt="Фото ${index + 1}">
            <button type="button" class="remove-photo" onclick="removePhoto(${photo.id})">
                <i class="fas fa-times"></i>
            </button>
            <div class="photo-order">${index + 1}</div>
        `;
        photoPreviewGrid.appendChild(photoItem);
    });
    
    // Update upload area text
    const photoUploadArea = document.getElementById('photoUploadArea');
    const uploadPlaceholder = photoUploadArea.querySelector('.upload-placeholder p');
    if (uploadedPhotos.length === 0) {
        uploadPlaceholder.textContent = 'Перетащите фотографии сюда или нажмите для выбора';
    } else {
        uploadPlaceholder.textContent = `Загружено ${uploadedPhotos.length} фото. Добавить еще`;
    }
}

// Remove photo - глобальная функция
function removePhoto(photoId) {
    uploadedPhotos = uploadedPhotos.filter(photo => photo.id !== photoId);
    updatePhotoPreview();
}

// Calculate prices
function calculatePrices() {
    const costPerDayInput = document.getElementById('cost_per_day');
    const costPerWeekInput = document.getElementById('cost_per_week');
    const costPerMonthInput = document.getElementById('cost_per_month');
    
    const costPerDay = parseFloat(costPerDayInput.value) || 0;
    const costPerWeek = parseFloat(costPerWeekInput.value) || 0;
    const costPerMonth = parseFloat(costPerMonthInput.value) || 0;
    
    // Auto-calculate weekly price if not set
    if (costPerDay > 0 && costPerWeek === 0) {
        costPerWeekInput.value = Math.round(costPerDay * 7 * 0.9); // 10% discount
    }
    
    // Auto-calculate monthly price if not set
    if (costPerDay > 0 && costPerMonth === 0) {
        costPerMonthInput.value = Math.round(costPerDay * 30 * 0.8); // 20% discount
    }
}

// Initialize Yandex Map
function initializeYandexMap() {
    if (typeof ymaps !== 'undefined') {
        ymaps.ready(function() {
            const mapElement = document.getElementById('map');
            
            // Remove placeholder
            const placeholder = mapElement.querySelector('.map-placeholder');
            if (placeholder) {
                placeholder.remove();
            }
            
            // Create map
            yandexMap = new ymaps.Map(mapElement, {
                center: [55.7558, 37.6176], // Moscow center
                zoom: 10,
                controls: ['zoomControl', 'fullscreenControl']
            });
            
            // Add click handler
            yandexMap.events.add('click', function(e) {
                const coords = e.get('coords');
                updateCoordinates(coords[0], coords[1]);
                updateMapPlacemark(coords);
            });
            
            // Load saved coordinates if exist
            const savedLat = document.getElementById('geolocation_x').value;
            const savedLng = document.getElementById('geolocation_y').value;
            if (savedLat && savedLng) {
                const coords = [parseFloat(savedLat), parseFloat(savedLng)];
                yandexMap.setCenter(coords);
                updateMapPlacemark(coords);
            }
        });
    } else {
        console.warn('Yandex Maps API not loaded');
    }
}

// Update map placemark
function updateMapPlacemark(coords) {
    if (!yandexMap) return;
    
    // Remove existing placemark
    if (mapPlacemark) {
        yandexMap.geoObjects.remove(mapPlacemark);
    }
    
    // Add new placemark
    mapPlacemark = new ymaps.Placemark(coords, {
        balloonContent: 'Выбранное местоположение'
    }, {
        preset: 'islands#redDotIcon'
    });
    
    yandexMap.geoObjects.add(mapPlacemark);
}

// Update coordinates in form
function updateCoordinates(lat, lng) {
    document.getElementById('geolocation_x').value = lat.toFixed(6);
    document.getElementById('geolocation_y').value = lng.toFixed(6);
}

// Initialize location autocomplete with geocoding
function initializeLocationAutocomplete() {
    const locationInput = document.getElementById('location');
    
    locationInput.addEventListener('input', debounce(async (event) => {
        const query = event.target.value;
        if (query.length > 3 && typeof ymaps !== 'undefined') {
            try {
                const geocoder = ymaps.geocode(query, {
                    results: 5
                });
                
                geocoder.then(function(res) {
                    // Here you could show suggestions
                    // For now, we'll just log them
                    console.log('Geocoding results:', res.geoObjects.get(0));
                });
            } catch (error) {
                console.error('Geocoding error:', error);
            }
        }
    }, 500));
    
    // Add button to get current location
    const locationGroup = locationInput.closest('.form-group');
    const locationButton = document.createElement('button');
    locationButton.type = 'button';
    locationButton.className = 'btn btn-outline';
    locationButton.style.marginTop = '0.5rem';
    locationButton.innerHTML = '<i class="fas fa-location-arrow"></i> Определить местоположение';
    locationButton.onclick = getCurrentLocation;
    locationGroup.appendChild(locationButton);
}

// Get current location
function getCurrentLocation() {
    if (navigator.geolocation) {
        navigator.geolocation.getCurrentPosition(
            function(position) {
                const lat = position.coords.latitude;
                const lng = position.coords.longitude;
                
                updateCoordinates(lat, lng);
                
                if (yandexMap) {
                    const coords = [lat, lng];
                    yandexMap.setCenter(coords);
                    updateMapPlacemark(coords);
                }
                
                // Reverse geocoding to get address
                if (typeof ymaps !== 'undefined') {
                    ymaps.geocode([lat, lng]).then(function(res) {
                        const address = res.geoObjects.get(0).getAddressLine();
                        document.getElementById('location').value = address;
                    });
                }
                
                showMessage('Местоположение определено', 'success');
            },
            function(error) {
                showMessage('Не удалось определить местоположение', 'error');
                console.error('Geolocation error:', error);
            }
        );
    } else {
        showMessage('Геолокация не поддерживается в вашем браузере', 'error');
    }
}

// Auto-save functionality
function autoSave() {
    const formData = new FormData(document.getElementById('advertisementForm'));
    const data = Object.fromEntries(formData.entries());
    
    // Add photos info
    data.photos = uploadedPhotos.map(photo => ({
        name: photo.file.name,
        size: photo.file.size
    }));
    
    localStorage.setItem('advertisement_draft', JSON.stringify(data));
    showAutoSaveIndicator('saved', 'Черновик сохранен');
}

// Load draft
function loadDraft() {
    const draft = localStorage.getItem('advertisement_draft');
    if (draft) {
        try {
            const data = JSON.parse(draft);
            
            // Fill form fields
            Object.keys(data).forEach(key => {
                const field = document.querySelector(`[name="${key}"]`);
                if (field && key !== 'photos') {
                    if (field.type === 'checkbox') {
                        field.checked = data[key] === 'true';
                    } else {
                        field.value = data[key];
                    }
                }
            });
            
            showMessage('Черновик загружен', 'info');
        } catch (error) {
            console.error('Error loading draft:', error);
        }
    }
}

// Save draft manually - глобальная функция
function saveDraft() {
    autoSave();
    showMessage('Черновик сохранен', 'success');
}

// Clear draft
function clearDraft() {
    localStorage.removeItem('advertisement_draft');
    showMessage('Черновик удален', 'info');
}

// Form validation
function validateForm() {
    let isValid = true;
    const requiredFields = ['title', 'category', 'description', 'cost_per_day', 'rental_rules', 'location'];
    
    // Clear previous errors
    document.querySelectorAll('.form-group.error').forEach(group => {
        group.classList.remove('error');
        const errorMessage = group.querySelector('.error-message');
        if (errorMessage) errorMessage.remove();
    });
    
    // Validate required fields
    requiredFields.forEach(fieldName => {
        const field = document.querySelector(`[name="${fieldName}"]`);
        if (!field.value.trim()) {
            markFieldAsError(field, 'Это поле обязательно для заполнения');
            isValid = false;
        }
    });
    
    // Validate photos
    if (uploadedPhotos.length === 0) {
        showMessage('Добавьте хотя бы одну фотографию', 'error');
        isValid = false;
    }
    
    // Validate prices
    const costPerDay = parseFloat(document.getElementById('cost_per_day').value);
    if (costPerDay <= 0) {
        markFieldAsError(document.getElementById('cost_per_day'), 'Цена должна быть больше 0');
        isValid = false;
    }
    
    return isValid;
}

// Mark field as error
function markFieldAsError(field, message) {
    const formGroup = field.closest('.form-group');
    formGroup.classList.add('error');
    
    const errorMessage = document.createElement('div');
    errorMessage.className = 'error-message';
    errorMessage.textContent = message;
    formGroup.appendChild(errorMessage);
}

// Initialize form validation
function initializeFormValidation() {
    const fields = document.querySelectorAll('input, select, textarea');
    
    fields.forEach(field => {
        field.addEventListener('blur', () => {
            validateField(field);
        });
    });
}

// Validate single field
function validateField(field) {
    const formGroup = field.closest('.form-group');
    formGroup.classList.remove('error');
    
    const errorMessage = formGroup.querySelector('.error-message');
    if (errorMessage) errorMessage.remove();
    
    if (field.hasAttribute('required') && !field.value.trim()) {
        markFieldAsError(field, 'Это поле обязательно для заполнения');
        return false;
    }
    
    // Specific validations
    if (field.type === 'email' && field.value) {
        const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
        if (!emailRegex.test(field.value)) {
            markFieldAsError(field, 'Введите корректный email');
            return false;
        }
    }
    
    if (field.type === 'number' && field.value) {
        const value = parseFloat(field.value);
        if (field.min && value < parseFloat(field.min)) {
            markFieldAsError(field, `Значение должно быть не менее ${field.min}`);
            return false;
        }
    }
    
    return true;
}

// Update form progress
function updateFormProgress() {
    const requiredFields = ['title', 'category', 'description', 'cost_per_day', 'rental_rules', 'location'];
    const filledFields = requiredFields.filter(fieldName => {
        const field = document.querySelector(`[name="${fieldName}"]`);
        return field && field.value.trim();
    });
    
    const progress = (filledFields.length / requiredFields.length) * 100;
    
    // Create progress bar if it doesn't exist
    let progressBar = document.querySelector('.form-progress');
    if (!progressBar) {
        progressBar = document.createElement('div');
        progressBar.className = 'form-progress';
        progressBar.innerHTML = `
            <div class="progress-bar">
                <div class="progress-fill" style="width: ${progress}%"></div>
            </div>
            <div class="progress-text">Заполнено ${Math.round(progress)}%</div>
        `;
        
        const form = document.getElementById('advertisementForm');
        form.insertBefore(progressBar, form.firstChild);
    } else {
        progressBar.querySelector('.progress-fill').style.width = `${progress}%`;
        progressBar.querySelector('.progress-text').textContent = `Заполнено ${Math.round(progress)}%`;
    }
}

// Show auto-save indicator
function showAutoSaveIndicator(status, message) {
    let indicator = document.querySelector('.auto-save-indicator');
    
    if (!indicator) {
        indicator = document.createElement('div');
        indicator.className = 'auto-save-indicator';
        document.body.appendChild(indicator);
    }
    
    indicator.className = `auto-save-indicator show ${status}`;
    indicator.innerHTML = `
        <i class="fas fa-${status === 'saving' ? 'spinner fa-spin' : status === 'saved' ? 'check' : 'info'}"></i>
        <span>${message}</span>
    `;
    
    if (status === 'saved') {
        setTimeout(() => {
            indicator.classList.remove('show');
        }, 3000);
    }
}

// Show message
function showMessage(message, type = 'info') {
    const messageElement = document.createElement('div');
    messageElement.className = `message ${type}`;
    messageElement.innerHTML = `
        <i class="fas fa-${type === 'success' ? 'check-circle' : type === 'error' ? 'exclamation-circle' : 'info-circle'}"></i>
        <span>${message}</span>
    `;
    
    const form = document.getElementById('advertisementForm');
    form.insertBefore(messageElement, form.firstChild);
    
    setTimeout(() => {
        messageElement.remove();
    }, 5000);
}

// Utility functions
function debounce(func, wait) {
    let timeout;
    return function executedFunction(...args) {
        const later = () => {
            clearTimeout(timeout);
            func(...args);
        };
        clearTimeout(timeout);
        timeout = setTimeout(later, wait);
    };
} 