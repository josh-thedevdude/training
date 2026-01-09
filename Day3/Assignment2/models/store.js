let usersStore = [];

export function getUsers() {
    return usersStore;
}

export function saveUsers(users) {
    usersStore = users;
}