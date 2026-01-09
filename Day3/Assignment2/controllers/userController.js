import { UserError } from "../errorhandling.js";
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
        throw new UserError("User not found", 404)
    }

    return user;
}

export function createUser(name, email, role) {
    if (!name) throw new UserError('Name is required', 400);
    if (!email) throw new UserError("Email is required", 400);

    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(email)) {
        throw new UserError('Invalid email format', 400);
    }

    const users = getUsers();

    if (users.find(u => u.email === email)) {
        throw new UserError('Email already exists', 400);
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
    if (!newName) throw new UserError('Name is required', 400);

    const users = getUsers();
    const user = users.find(u => u.id === id);

    if (!user) {
        throw new UserError('User not found', 404);
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
        throw new UserError('User not found', 400);
    }

    const updatedUsers = users.filter(u => u.id !== id);
    saveUsers(updatedUsers);
}