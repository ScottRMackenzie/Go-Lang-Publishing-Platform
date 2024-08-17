let start = 0;
const count = 10;
let isLoading = false;

async function DisplayBooks() {
    if (isLoading) return;
    isLoading = true;

    const url = BASE_API_URL + '/api/v1/public/books/sorted'; 
    const params = {
        start: start,
        count: count,
        sort: "created_at",
        is_accending: false
    };
    const res = await fetch(url, {
        method: 'POST',
        body: JSON.stringify(params),
        headers: {
            'Content-Type': 'application/json'
        }
    });
    const data = await res.json();

    const body = document.body;

    if (data.length === 0 || data.length < count) {
        if (document.getElementById('load-more-button-container')) {
            body.removeChild(document.getElementById('load-more-button-container'));
        }
        return;
    }

    // Create a container for the grid if it doesn't exist
    let gridContainer = document.getElementById('grid-container');
    if (!gridContainer) {
        gridContainer = document.createElement('div');
        gridContainer.setAttribute('id', 'grid-container');
        gridContainer.setAttribute('class', 'grid grid-cols-5 gap-4 p-4'); // Adjust grid-cols to your desired number of columns
        body.appendChild(gridContainer);
    }

    for (let i = 0; i < data.length; i++) {
        const book = data[i];
        const div = document.createElement('div');

        div.setAttribute('class', 'bg-gray-200 p-4 rounded-md flex flex-col items-center cursor-pointer hover:bg-grey-400'); // Center elements
        div.innerHTML = `
            <img src="${book.cover_img_url}" alt="${book.title}" class="max-w-16 h-auto rounded-md object-cover object-center mb-2">
            <h1 class="text-xl font-bold text-center">${book.title}</h1>
            <h2 class="text-lg font-medium text-center">${book.author}</h2>
            <p class="text-gray-600 text-center">${book.genre}</p>
        `;

        // Add click event listener to redirect to /book/{id}
        div.addEventListener('click', () => {
            window.location.href = `/book/${book.id}`;
        });

        gridContainer.appendChild(div);
    }

    // Update the start for the next batch of books
    start += count;
    isLoading = false;

    if (!document.getElementById('load-more-button-container')) {
        const loadMoreButtonContainer = document.createElement('div');
        loadMoreButtonContainer.setAttribute('id', 'load-more-button-container');
        loadMoreButtonContainer.setAttribute('class', 'flex justify-center mt-4');

        const loadMoreButton = document.createElement('button');
        loadMoreButton.setAttribute('id', 'load-more-button');
        loadMoreButton.setAttribute('class', 'bg-blue-500 text-white p-2 rounded-md cursor-pointer hover:bg-blue-600 mb-4');
        loadMoreButton.innerText = 'Load More';
        loadMoreButton.addEventListener('click', DisplayBooks);

        loadMoreButtonContainer.appendChild(loadMoreButton);
        body.appendChild(loadMoreButtonContainer);
    }
}

// Initial load
DisplayBooks();

// async function GetUsers() {
//     res = await fetch('http://api.tb-books.local:2337/api/users', {
//         method: 'GET',
//         credentials: "include",
//         // credentials: 'include': This will include the cookies in the request if the api server is on a different port and/or domain
//         // credentials: 'same-origin': This will include the cookies in the request is the api server is on the same port and domain
//     });
//     data = await res.json();

//     console.log(data);
// }

// GetUsers();

