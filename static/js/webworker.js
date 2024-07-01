let token = ""

// Check if api endpoint respond with the error code 200
async function checkApi() {

    // Check every 2 seconds until the token is not empty
    while (token === "") {
        await new Promise(r => setTimeout(r, 2000));
    }

    await new Promise(r => setTimeout(r, 5000));

    let response = await fetch('/api/v1/', {
        method: 'GET',
        headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + token,
        }

    });

    if (response.status !== 200) {
        console.error("API is not responding")
        postMessage("login")
    }

    checkApi();
}

onmessage = function (e) {
    token = e.data
}
/*

*/
checkApi();