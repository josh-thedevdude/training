import express from 'express';
import userRouter from './routes/v1/userRouter.js';

const app = express();
const PORT = process.env.PORT || 3000;

app.use(express.json());

export let users = [];

app.get('/health', (req, res) => {
    res.status(200).json({ status: 'ok' });
});

app.use("/api/v1", userRouter)

app.use((req, res) => {
    res.status(404).json({ error: 'Route not found' });
});

app.use((err, req, res, next) => {
    console.error(err.stack)
    res.status(err.statusCode).json({ error: err.message });
});

app.listen(PORT, () => {
    console.log(`Server is running on http://localhost:${PORT}`);
});