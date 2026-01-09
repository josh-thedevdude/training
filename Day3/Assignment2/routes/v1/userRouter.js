import { Router } from "express";
import { createUser, deleteUser, getAllUsers, getUserById, updateUserName } from "../../controllers/userController.js";

const userRouter = Router()


userRouter.get('/users', (req, res) => {
    const { role } = req.query
    const users = getAllUsers(role)
    res.status(200).json(users);
});

userRouter.get('/users/:id', (req, res) => {
    const { id } = req.params
    const user = getUserById(id);
    res.status(200).json(user);
});

userRouter.post('/users', (req, res) => {
    const { name, email, role } = req.body
    const newUser = createUser(name, email, role);
    res.status(201).json(newUser);
});

userRouter.patch("/users/:id", (req, res) => {
    const { id } = req.params;
    const { name } = req.body;
    const updatedUser = updateUserName(id, name);
    res.status(200).json(updatedUser);
})


userRouter.delete("/users/:id", (req, res) => {
    const { id } = req.params;
    deleteUser(id);
    res.status(204).json({});
})

export default userRouter;