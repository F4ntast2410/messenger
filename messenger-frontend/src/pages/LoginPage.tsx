import { useState, type SubmitEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { login } from '../api/auth';
import { ApiError } from '../api/client';
import { useAuth } from '../context/AuthContext';

export function LoginPage() {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const navigate = useNavigate();
    const { refetchUser } = useAuth();

    async function handleSubmit(e: SubmitEvent<HTMLFormElement>) {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        try {
        await login(email, password);
        await refetchUser();
        navigate('/chats');
        } catch (err) {
        if (err instanceof ApiError && err.status === 400) {
            setError('Неверная почта или пароль');
        } else {
            setError('Что-то пошло не так, попробуй ещё раз');
        }
        } finally {
        setIsSubmitting(false);
        }
    }

    return (
        <div>
        <h1>Вход</h1>
        <form onSubmit={handleSubmit}>
            <div>
            <label htmlFor="email">Email</label>
            <input id="email" type="email" value={email}
                    onChange={(e) => setEmail(e.target.value)} required />
            </div>
            <div>
            <label htmlFor="password">Пароль</label>
            <input id="password" type="password" value={password}
                    onChange={(e) => setPassword(e.target.value)} required />
            </div>

            {error && <p style={{ color: 'red' }}>{error}</p>}

            <button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Входим...' : 'Войти'}
            </button>
        </form>

        <p>Нет аккаунта? <Link to="/register">Зарегистрироваться</Link></p>
        </div>
    );
}