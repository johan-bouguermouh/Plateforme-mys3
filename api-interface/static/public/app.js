// const $ = (selector) => document.querySelector(selector);
// const container = $('#users');
// const API_ENDPOINT = '/api/v1/users';

// const listUsers = async () => {
// 	const response = await fetch(API_ENDPOINT);
// 	const data = await response.json();
// 	const users = data.users.reverse();

// 	for (let index = 0; index < users.length; index++) {
// 		const child = document.createElement('li');
// 		child.className = 'list-group-item';
// 		child.innerText = users[index].name;

// 		container.appendChild(child);
// 	}
// };

// $('#add_user').addEventListener('click', async (e) => {
// 	e.preventDefault();
// 	const user = $('#user').value;

// 	if (!user) return;

// 	const form = new FormData();
// 	form.append('user', user);

// 	const response = await fetch(API_ENDPOINT, {
// 		method: 'POST',
// 		body: form,
// 	});

// 	const data = await response.json();

// 	const child = document.createElement('li');
// 	child.className = 'list-group-item';
// 	child.innerText = data.user.name;

// 	container.insertBefore(child, container.firstChild);

// 	$('#user').value = '';
// });

// document.addEventListener('DOMContentLoaded', listUsers);
// JavaScript for handling user registration and login

// document.getElementById('register-form')?.addEventListener('submit', async (e) => {
// 	e.preventDefault();
// 	const username = document.getElementById('reg-username').value;
// 	const password = document.getElementById('reg-password').value;
  
// 	const response = await fetch('http://localhost:3000/register', {
// 	  method: 'POST',
// 	  headers: {
// 		'Content-Type': 'application/json'
// 	  },
// 	  body: JSON.stringify({ name: username, password: password })
// 	});
  
// 	const result = await response.text();
// 	document.getElementById('register-response').textContent = result;
//   });
  
//  // JavaScript for handling user registration and login

// document.getElementById('register-form')?.addEventListener('submit', async (e) => {
// 	e.preventDefault();
// 	const username = document.getElementById('reg-username').value;
// 	const password = document.getElementById('reg-password').value;
  
// 	const response = await fetch('http://localhost:3000/register', {
// 	  method: 'POST',
// 	  headers: {
// 		'Content-Type': 'application/json'
// 	  },
// 	  body: JSON.stringify({ name: username, password: password })
// 	});
  
// 	const result = await response.text();
// 	document.getElementById('register-response').textContent = result;
//   });
  
//   document.getElementById('login-form')?.addEventListener('submit', async (e) => {
// 	e.preventDefault();
// 	const username = document.getElementById('username').value;
// 	const password = document.getElementById('password').value;
  
// 	const response = await fetch('http://localhost:3000/login', {
// 	  method: 'POST',
// 	  headers: {
// 		'Content-Type': 'application/json'
// 	  },
// 	  body: JSON.stringify({ name: username, password: password })
// 	});
  
// 	const result = await response.json();
// 	if (response.ok) {
// 	  document.getElementById('login-response').textContent = `Token: ${result.token}`;
	
// 	  localStorage.setItem('authToken', result.token);
// 	} else {
// 	  document.getElementById('login-response').textContent = `Error: ${result.message}`;
// 	}
//   });
  
document.getElementById('register-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('reg-username').value;
    const password = document.getElementById('reg-password').value;

    const response = await fetch('http://localhost:3000/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: username, password: password })
    });

    const result = await response.text();
    document.getElementById('register-response').textContent = result;
});

document.getElementById('login-form')?.addEventListener('submit', async (e) => {
    e.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    const response = await fetch('http://localhost:3000/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: username, password: password })
    });

    const result = await response.json();
    if (response.ok) {
        document.getElementById('login-response').textContent = `Token: ${result.token}`;
        localStorage.setItem('authToken', result.token); // Stocke le token
    } else {
        document.getElementById('login-response').textContent = `Error: ${result.message}`;
    }
});

async function accessProtectedRoute() {
    const token = localStorage.getItem('authToken'); // Récupère le token
    if (!token) {
        console.log('No token found');
        return;
    }

    const response = await fetch('http://localhost:3000/bucket/myNewBucket/files/', {
        method: 'GET',
        headers: {
            'Authorization': `Bearer ${token}` // Inclut le token dans l'en-tête
        }
    });

    if (response.ok) {
        const data = await response.json();
        console.log(data);
    } else {
        console.log('Error:', response.statusText);
    }
}


accessProtectedRoute();
