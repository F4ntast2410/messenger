import { useState, type SubmitEvent } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { register } from '../api/auth';
import { ApiError } from '../api/client';

export function RegisterPage() {
    const [name, setName] = useState('');
    const [username, setUsername] = useState('')
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState<string | null>(null);
    const [isSubmitting, setIsSubmitting] = useState(false);

    const navigate = useNavigate();

    async function handleSubmit(e: SubmitEvent<HTMLFormElement>) {
        e.preventDefault();
        setError(null);
        setIsSubmitting(true);

        try {
        await register(email, username, name, password);
        navigate('/login');
        } catch (err) {
        if (err instanceof ApiError && err.status === 400) {
            setError('Пользователь с такой почтой или username уже существует');
        } else {
            console.error(err)
            setError('Что-то пошло не так, попробуй ещё раз');
        }
        } finally {
        setIsSubmitting(false);
        }
    }

    return (
        <div>
        <h1>Регистрация</h1>
        <form onSubmit={handleSubmit}>
            <div>
            <label htmlFor="name">Имя</label>
            <input id="name" type="text" value={name}
                    onChange={(e) => setName(e.target.value)} required />
            </div>
            <div>
            <label htmlFor="username">Username</label>
            <input id="username" type="text" value={username}
                    onChange={(e) => setUsername(e.target.value)} required />
            </div>
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
            {isSubmitting ? 'Регистрация...' : 'Зарегистрироваться'}
            </button>
        </form>

        <p>Есть аккаунт? <Link to="/login">Войти</Link></p>
        </div>
    );
}