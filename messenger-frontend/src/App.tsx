import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { ProtectedRoute } from './components/ProtectedRoute';
import { LoginPage } from './pages/LoginPage';
import { RegisterPage } from './pages/RegisterPage';
import { ChatsPage } from './pages/ChatsPage';
import { ChatPage } from './pages/ChatPage';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/register" element={<RegisterPage />} />

        <Route element={<ProtectedRoute />}>
          <Route path="/chats" element={<ChatsPage />}>
            <Route path=":chatId" element={<ChatPage />} />
          </Route>
        </Route>

        <Route path="*" element={<Navigate to="/chats" replace />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;