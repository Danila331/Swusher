// Initialize map
function initMap() {
    // Create map instance
    const map = new ymaps.Map('map', {
        center: [55.76, 37.64], // Moscow coordinates
        zoom: 10,
        controls: ['zoomControl', 'fullscreenControl']
    });

    // Disable scroll zoom
    map.behaviors.disable('scrollZoom');

    // Remove unnecessary controls
    map.controls.remove('geolocationControl');
    map.controls.remove('searchControl');
    map.controls.remove('trafficControl');
    map.controls.remove('typeSelector');
    map.controls.remove('rulerControl');

    // Set custom style
    map.panes.get('ground').getElement().style.filter = 'grayscale(100%)';

    // Example markers
    const markers = [
        {
            coordinates: [55.76, 37.64],
            title: 'Canon EOS R5',
            price: '1500 ₽/день'
        },
        {
            coordinates: [55.75, 37.65],
            title: 'Велосипед Trek',
            price: '500 ₽/день'
        },
        {
            coordinates: [55.77, 37.63],
            title: 'Дрель Bosch',
            price: '300 ₽/день'
        }
    ];

    // Create markers
    const markerCollection = new ymaps.GeoObjectCollection();

    markers.forEach(marker => {
        const placemark = new ymaps.Placemark(marker.coordinates, {
            balloonContent: `
                <div class="map-balloon">
                    <h3>${marker.title}</h3>
                    <p>${marker.price}</p>
                </div>
            `
        }, {
            preset: 'islands#blueDotIcon',
            balloonOffset: [0, -20]
        });

        markerCollection.add(placemark);
    });

    // Add markers to map
    map.geoObjects.add(markerCollection);

    // Fit map to markers
    map.setBounds(markerCollection.getBounds(), {
        checkZoomRange: true,
        zoomMargin: 50
    });

    // Add click handler
    map.geoObjects.events.add('click', function(e) {
        const target = e.get('target');
        target.balloon.open();
    });

    // Add hover effects
    map.geoObjects.events.add('mouseenter', function(e) {
        const target = e.get('target');
        target.options.set('iconImageSize', [50, 50]);
        target.options.set('iconImageOffset', [-25, -50]);
    });

    map.geoObjects.events.add('mouseleave', function(e) {
        const target = e.get('target');
        target.options.set('iconImageSize', [40, 40]);
        target.options.set('iconImageOffset', [-20, -40]);
    });
}

// Initialize map when API is ready
ymaps.ready(initMap); 