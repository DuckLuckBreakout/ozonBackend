const CACHE_NAME = 'OZON-CACHE-v1';


const resources = [];


self.addEventListener('install', (event) => {
    event.waitUntil(
        caches.open(CACHE_NAME)
            .then((cache) => cache.addAll(resources))
            .catch((err) => {
                console.error('smth went wrong with caches.open: ', err);
            }),
    );
});

self.addEventListener('activate', (event) => {
    event.waitUntil(self.clients.claim());
});

self.addEventListener('push', (e) => {
    const data = e.data.json();
    self.registration.showNotification('Уведомление от DuckLuck Market!', {
        body: `Ваш заказ № ${data.order_number.number} приобрёл новый статус: ${data.status}`,
        icon: '/svg/header/logo.svg',
    });
});

self.addEventListener('fetch', (event) => {
    event.respondWith(
        (async() => {
            if (navigator.onLine) {
                const response = await fetch(event.request);
                if (response && response.ok && event.type && event.request.method === 'GET') {
                    const cache = await caches.open(CACHE_NAME);
                    await cache.put(event.request, response.clone());
                }

                return response;
            }
            const cache = await caches.open(CACHE_NAME);
            const response = await cache.match(event.request);
            if (!response || event.request.method !== 'GET') {
                return new Response(null, {
                    headers: {
                        'Content-Type': 'text/html; charset=utf-8',
                    },
                    status: 410,
                });
            } else {
                return response;
            }
        })(),
    );
});
