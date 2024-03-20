import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    scenarios: {
        user: {
            executor: 'per-vu-iterations',
            exec: 'createUser',
            vus: 10,
            iterations: 100,
        },
        admin: {
            executor: 'per-vu-iterations',
            exec: 'createAdmin',
            vus: 10,
            iterations: 100,
        },
        authorize: {
            executor: 'per-vu-iterations',
            exec: 'authorizeOnly',
            vus: 10,
            iterations: 100,
        },
        post: {
            executor: 'per-vu-iterations',
            exec: 'post',
            vus: 10,
            iterations: 100,
        },
    },
};

// Create a random string of given length
function randomString(length, charset = '') {
    if (!charset) charset = 'abcdefghijklmnopqrstuvwxyz';
    let res = '';
    while (length--) res += charset[(Math.random() * charset.length) | 0];
    return res;
}

const PASSWORD = 'secret';
const BASE_URL = 'http://localhost:8080/api/v1';

export function createUser() {
    const res = http.post(`${BASE_URL}/signup`, {
        email: `${randomString(10)}@example.com`,
        password: PASSWORD,
        username: randomString(10),
    })

    check(res, {'created user': (r) => r.status === 201});
};


export function createAdmin() {
    const adminLogin = http.post(`${BASE_URL}/login`, {
        email: 'admin@outlook.com',
        password: 'secret',
    })

    check(adminLogin, { 'login admin': (r) => r.status === 302});

    const authToken = adminLogin.json('token');
    const formData = {
        email: `${randomString(10)}@example.com`,
        password: PASSWORD,
        username: randomString(10),
    };
    const headers = {
        'Authorization': 'Bearer ' + authToken,
    };

    const create = http.post(`${BASE_URL}/signup-admin`, formData, { headers: headers });

    check(create, {'create admin': (r) => r.status === 201});
};

export function authorizeOnly() {
    const adminLogin = http.post(`${BASE_URL}/login`, {
        email: 'admin@outlook.com',
        password: 'secret',
    })

    check(adminLogin, { 'login admin': (r) => r.status === 302});

    const authToken = adminLogin.json('token');
    const headers = {
        'Authorization': 'Bearer ' + authToken,
    };

    const user = 'user';
    const normalUsers = http.get(`${BASE_URL}/users/${user}`, { headers: headers });
    check(normalUsers, {'normal users': (r) => r.status === 200});

    const admin = 'admin';
    const adminUsers = http.get(`${BASE_URL}/users/${admin}`, { headers: headers });
    check(adminUsers, {'admin users': (r) => r.status === 200});

    const usernameContain = 'a';
    const usersByUsername = http.get(`${BASE_URL}/users?username=${usernameContain}`, { headers: headers });
    check(usersByUsername, {'users by username': (r) => r.status === 200});
};

// Sample arrays of words for generating random titles
const adjectives = ['Amazing', 'Awesome', 'Fantastic', 'Incredible', 'Superb', 'Wonderful'];
const nouns = ['Adventure', 'Journey', 'Experience', 'Discovery', 'Quest', 'Mission'];

// Function to generate a random title
function generateRandomTitle() {
    const adjective = adjectives[Math.floor(Math.random() * adjectives.length)];
    const noun = nouns[Math.floor(Math.random() * nouns.length)];
    return `${adjective} ${noun}`;
}

const contentPhrases = [
    'Lorem ipsum dolor sit amet, consectetur adipiscing elit.',
    'Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.',
    'Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.',
    'Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.',
    'Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.',
    'Integer posuere erat a ante venenatis dapibus posuere velit aliquet.',
    'Maecenas sed diam eget risus varius blandit sit amet non magna.',
    'Nullam quis risus eget urna mollis ornare vel eu leo.',
];

// Function to generate random content
function generateRandomContent() {
    const numberOfSentences = Math.floor(Math.random() * 5) + 1; // Generate between 1 and 5 sentences
    let content = '';
    for (let i = 0; i < numberOfSentences; i++) {
        const randomIndex = Math.floor(Math.random() * contentPhrases.length);
        content += contentPhrases[randomIndex] + ' ';
    }
    return content.trim(); // Trim extra whitespace
}

// Example usage
function generateRandomTitleAndContent() {
    const title = generateRandomTitle();
    const content = generateRandomContent();
    return { title, content };
}

function getFirstWord(title) {
    // Split the title string into an array of words
    const words = title.split(' ');

    // Return the first word
    return words[0];
}

export function post() {
    const userID = 1;
    const { title, content } = generateRandomTitleAndContent();
    const formData = {
        title: title,
        content: content,
        id: userID,
    };
    const firstWord = getFirstWord(title);

    const createPost = http.post(`${BASE_URL}/posts`, formData);
    check(createPost, {'create post': (r) => r.status === 201});

    const getPostByUserID = http.get(`${BASE_URL}/posts/${userID}`, formData);
    check(getPostByUserID, {'get post by user ID': (r) => r.status === 302});

    const getPostByKeyword = http.get(`${BASE_URL}/posts?keyword=${firstWord}`, formData);
    check(getPostByKeyword, {'get post keyword': (r) => r.status === 302});
};
