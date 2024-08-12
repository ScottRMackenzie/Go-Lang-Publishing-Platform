async function GetUsers() {
    const response = await fetch('http://localhost/client-api/users', {
        method: 'GET',
    });
    const data = await response.json();
    console.log(data);
}

GetUsers();

// function getCookieValue(name) {
//     const value = `; ${document.cookie}`;
//     const parts = value.split(`; ${name}=`);
//     if (parts.length === 2) return parts.pop().split(';').shift();
// }