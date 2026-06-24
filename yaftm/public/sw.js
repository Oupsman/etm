/* global clients */
self.addEventListener('push', event => {
  let data = { title: 'ETM', body: 'You have a notification' }
  if (event.data) {
    try { data = event.data.json() } catch { data.body = event.data.text() }
  }
  event.waitUntil(
    self.registration.showNotification(data.title, {
      body: data.body,
      icon: '/favicon.ico',
    })
  )
})

self.addEventListener('notificationclick', event => {
  event.notification.close()
  event.waitUntil(clients.openWindow('/'))
})
