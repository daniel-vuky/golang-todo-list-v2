// Selectors
const usernameInput = document.getElementById("username")
const toDoInput = document.querySelector('.todo-input');
const toDoBtn = document.querySelector('.todo-btn');
const toDoList = document.querySelector('.todo-list');
const standardTheme = document.querySelector('.standard-theme');
const lightTheme = document.querySelector('.light-theme');
const darkerTheme = document.querySelector('.darker-theme');

const itemUrl = `${window.location.origin}/items`;
const logOutUrl = `${window.location.origin}/logout`;
const STATUS_PROCESSING = 1;
const STATUS_COMPLETED = 2;


// Event Listeners

usernameInput.addEventListener('mouseover', usernameMouseOver);
usernameInput.addEventListener('mouseout', usernameMouseOut);
toDoBtn.addEventListener('click', addToDo);
toDoList.addEventListener('click', todoAction);
document.addEventListener("DOMContentLoaded", getTodos);
standardTheme.addEventListener('click', () => changeTheme('standard'));
lightTheme.addEventListener('click', () => changeTheme('light'));
darkerTheme.addEventListener('click', () => changeTheme('darker'));

// Check if one theme has been set previously and apply it (or std theme if not found):
let savedTheme = localStorage.getItem('savedTheme');
savedTheme === null ?
    changeTheme('standard')
    : changeTheme(localStorage.getItem('savedTheme'));

// Functions;
function usernameMouseOver(event) {
    usernameInput.innerText = "Sign Out?";
    usernameInput.style.textDecoration = "underline";
    usernameInput.style.cursor = "pointer";
    usernameInput.addEventListener('click', function () {
        window.location.href = logOutUrl;
    });
}

function usernameMouseOut(event) {
    const item = event.target;
    usernameInput.innerText = item.getAttribute("data-username");
    usernameInput.style.textDecoration = "none";
    usernameInput.style.cursor = "none";
}

function todoAction(event){
    const item = event.target;

    // delete
    if(item.getAttribute("data-action") === 'delete')
    {
        removeItem(item, removeItemElement);
    }

    // check
    if(item.getAttribute("data-action") === 'checked')
    {
        completeItem(item, completeItemElement)
    }
}

// Saving to local storage:
function saveItem(todo, callback){
    fetch(itemUrl, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            title: todo,
            status: 1
        }) // Convert data to JSON string
    })
        .then(response => response.json())
        .then(data => {
            callback()
        })
        .catch(error => {
            console.log(error);
        });
}

function getTodos() {
    toDoList.innerHTML = "";
    fetch(itemUrl)
        .then(response => {
            // Check if the response status is OK (status code 200)
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }

            // Parse the JSON response
            return response.json();
        })
        .then(data => {
            // Handle the data from the response
            for (var key in data) {
                addToDoElement(data[key])
            }
        })
        .catch(error => {
            // Handle errors
            console.error("Error:", error);
        });
}

function addToDo(event) {
    // Prevents form from submitting / Prevents form from reloading;
    event.preventDefault();
    if (toDoInput.value === '') {
        alert("You must write something!");
        return
    }

    // Adding to local storage;
    saveItem(toDoInput.value, getTodos);
    // Clear the input;
    toDoInput.value = '';
}

function addToDoElement(item) {
    // toDo DIV;
    const toDoDiv = document.createElement("div");
    toDoDiv.classList.add('todo', `${savedTheme}-todo`);
    if (item.status === STATUS_COMPLETED) {
        toDoDiv.classList.add('completed')
    }

    // Create LI
    const newToDo = document.createElement('li');

    // newToDo.innerText = "hey";
    newToDo.innerText = item.title;
    newToDo.classList.add('todo-item');
    newToDo.setAttribute("id", `todo-item-${item.item_id}`);
    toDoDiv.appendChild(newToDo);

    // check btn;
    const checked = document.createElement('button');
    checked.innerHTML = '<i class="fas fa-check"></i>';
    checked.classList.add('check-btn', `${savedTheme}-button`);
    checked.setAttribute("data-item-id", item.item_id);
    checked.setAttribute("data-action", "checked");
    toDoDiv.appendChild(checked);
    // delete btn;
    const deleted = document.createElement('button');
    deleted.innerHTML = '<i class="fas fa-trash"></i>';
    deleted.classList.add('delete-btn', `${savedTheme}-button`);
    deleted.setAttribute("data-item-id", item.item_id);
    deleted.setAttribute("data-action", "delete");
    toDoDiv.appendChild(deleted);

    // Append to list;
    toDoList.appendChild(toDoDiv);
}

function completeItem(itemElement, callback){
    let itemId = itemElement.getAttribute("data-item-id"),
        currentTitle = document.getElementById(`todo-item-${itemId}`),
        completed = itemElement.parentElement.classList.contains("completed")
        itemUrlEncoded = itemUrl + "/" + itemId;
    fetch(itemUrlEncoded, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            title: currentTitle.innerText,
            status: !!completed ? 1 : 2
        }) // Convert data to JSON string
    })
        .then(response => response.json())
        .then(data => {
            callback(itemElement)
            getTodos()
        })
        .catch(error => {
            console.log(error);
        });
}

function completeItemElement(itemElement) {
    itemElement.parentElement.classList.toggle("completed");
}

function removeItem(itemElement, callback){
    let itemId = itemElement.getAttribute("data-item-id"),
        itemUrlEncoded = itemUrl + "/" + itemId;
    fetch(itemUrlEncoded, {
        method: 'DELETE'
    })
        .then(response => response.json())
        .then(data => {
            callback(itemElement)
        })
        .catch(error => {
            console.log(error);
        });
}

function removeItemElement(itemElement) {
    // item.parentElement.remove();
    // animation
    itemElement.parentElement.classList.add("fall");

    itemElement.parentElement.addEventListener('transitionend', function(){
        itemElement.parentElement.remove();
    });
}

// Change theme function:
function changeTheme(color) {
    localStorage.setItem('savedTheme', color);
    savedTheme = localStorage.getItem('savedTheme');

    document.body.className = color;
    // Change blinking cursor for darker theme:
    color === 'darker' ?
        document.getElementById('title').classList.add('darker-title')
        : document.getElementById('title').classList.remove('darker-title');

    document.querySelector('input').className = `${color}-input`;
    // Change todo color without changing their status (completed or not):
    document.querySelectorAll('.todo').forEach(todo => {
        Array.from(todo.classList).some(item => item === 'completed') ?
            todo.className = `todo ${color}-todo completed`
            : todo.className = `todo ${color}-todo`;
    });
    // Change buttons color according to their type (todo, check or delete):
    document.querySelectorAll('button').forEach(button => {
        Array.from(button.classList).some(item => {
            if (item === 'check-btn') {
                button.className = `check-btn ${color}-button`;
            } else if (item === 'delete-btn') {
                button.className = `delete-btn ${color}-button`;
            } else if (item === 'todo-btn') {
                button.className = `todo-btn ${color}-button`;
            }
        });
    });
}