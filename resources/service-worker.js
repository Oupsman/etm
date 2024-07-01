function urlBase64ToUint8Array(base64String) {
    const padding = '='.repeat((4 - (base64String.length % 4)) % 4)
    const base64 = (base64String + padding).replace(/\-/g, '+').replace(/_/g, '/')
    const rawData = atob(base64)
    const outputArray = new Uint8Array(rawData.length)
    for (let i = 0; i < rawData.length; ++i) {
        outputArray[i] = rawData.charCodeAt(i)
    }
    return outputArray
}

self.addEventListener('activate', async () => {
    // This will be called only once when the service worker is activated.
    try {
        let vapidPubKey = await fetch('/api/v1/getvapidkey')
        vapidPubKey = await vapidPubKey.json()
        console.log('Vapid Key:', vapidPubKey)
        const applicationServerKey = urlBase64ToUint8Array(
            vapidPubKey.public_key
        )
        console.log('Application Server Key:', applicationServerKey, vapidPubKey)
        const options = { applicationServerKey: applicationServerKey, userVisibleOnly: true }
        console.log('Options:', options)
        const subscription = await self.registration.pushManager.subscribe(options)
        console.log("Subscription:", subscription)
        console.log(JSON.stringify(subscription))
    } catch (err) {
        console.log('Error', err)
    }
})

self.addEventListener('push', event => {
    console.log('[Service Worker] Push Received.');
    console.log(`[Service Worker] Push had this data: "${event.data.text()}"`);

    var myNotif = event.data.text();
    console.log('myNotif:', myNotif);
    const title = "SwitchDB Push Notification";
    const options = {
        body: myNotif,
    };
    const promiseChain = self.registration.showNotification(title, options);

    event.waitUntil(promiseChain);
});