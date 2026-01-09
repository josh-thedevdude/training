import { getUsers, saveUsers } from "../models/store.js";
import { v4 as uuidv4 } from "uuid"

export function getAllUsers(role) {
    const users = getUsers();

    if (role) {
        const filteredUsers = users.filter(user => user.role === role);
        return filteredUsers;
    }

    return users;
}

export function getUserById(id) {
    const users = getUsers();
    const user = users.find(u => u.id === id);

    if (!user) {
        throw new Error({ error: "User not found", status: 404 })
    }

    return user;
}

export function createUser(name, email, role) {
    if (!name) throw new Error({ error: 'Name is required', status: 400 });
    if (!email) throw new Error({ error: "Email is required", status: 400 });

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        throw new Error({ error: 'Invalid email format', status: 400 });
    }

    const users = getUsers();

    if (users.find(u => u.email === email)) {
        throw new Error({ error: 'Email already exists', status: 400 });
    }

    const newUser = {
        id: uuidv4(),
        name,
        email,
        role: role || 'user'
    };

    users.push(newUser);
    saveUsers(users);

    return newUser;
}

export function updateUserName(id, newName) {
    if (!newName) throw new Error({ error: 'Name is required', status: 400 });

    const users = getUsers();
    const user = users.find(u => u.id === id);

    if (!user) {
        throw new Error({ error: 'User not found', status: 404 });
    }

    user.name = newName;
    saveUsers(users);

    return user;
}

export function deleteUser(id) {
    const users = getUsers();
    console.log(users)
    const userExists = users.find(u => u.id === id);
    console.log(userExists)

    if (!userExists) {
        throw new Error({ error: 'User not found', status: 400 });
    }

    const updatedUsers = users.filter(u => u.id !== id);
    saveUsers(updatedUsers);
}