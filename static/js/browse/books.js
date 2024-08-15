async function DisplayBooks() {
    res = await fetch('http://localhost:2337/api/v1/public/books', {
        method: 'GET',
    });
    data = await res.json();

    body = document.body

    // Create a container for the grid
    const gridContainer = document.createElement('div');
    gridContainer.setAttribute('class', 'grid grid-cols-5 gap-4 p-4'); // Adjust grid-cols to your desired number of columns

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

        div.addEventListener('click', () => {
            window.location.href = `/book/${book.id}`;
        });

        gridContainer.appendChild(div);
    }

// Append the grid container to the body or a specific element
body.appendChild(gridContainer);

}

DisplayBooks();