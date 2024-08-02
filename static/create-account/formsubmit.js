// async function submitForm() {
//     console.log('submitting form');
//     let formData = new FormData(document.getElementById('create-account-form'));

//     const res = await fetch ('/api/users/create', {
//         method: 'POST',
//         body: formData
//     })

//     if (res.ok) {
//         window.location.href = '/success';
//     }
//     else {
//         alert('Failed to create account');
//     }
// }

// const send = document.querySelector("#send");
// send.addEventListener("click", submitForm);


const form = document.getElementById("create-account-form");

async function sendData() {
    const formData = new FormData(form);
    const responseData = Object.fromEntries(formData.entries());
    console.log(responseData);

    const response = await fetch("http://127.0.0.1:2337/api/users/create", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify(responseData),
    });

    if (response.ok) {
        window.location.href = "/success";
    } else {
        alert("Failed to create account");
    }
}

// Take over form submission
form.addEventListener("submit", (event) => {
    event.preventDefault();
    sendData();
});
