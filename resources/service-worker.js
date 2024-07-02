let token = ""

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

self.addEventListener('message', function (e) {
    token = e.data
})

self.addEventListener('activate', async () => {
    // This will be called only once when the service worker is activated.
    try {
        // Wait until token is populated
        while (token === "")  {
            await new Promise(r => setTimeout(r, 2000));
        }
        let vapidPubKey = await fetch('/api/v1/getvapidkey', {
            method: "GET",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
        })
        vapidPubKey = await vapidPubKey.json()
        const applicationServerKey = urlBase64ToUint8Array(
            vapidPubKey.public_key
        )
        const options = { applicationServerKey: applicationServerKey, userVisibleOnly: true }
        const subscription = await self.registration.pushManager.subscribe(options)

        const data = {
            subscription: subscription
        }
        console.log("Data to send: ", JSON.stringify(data))
        await fetch("/api/v1/user/updatesubscription", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json',
                'Authorization': 'Bearer ' + token
            },
            body: JSON.stringify({
                subscription: JSON.stringify(subscription)
            })
        }).then( )
        console.log(JSON.stringify(subscription))
    } catch (err) {
        console.log('Error', JSON.stringify(err))
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
