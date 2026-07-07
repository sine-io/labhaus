'use client';

import { useState } from 'react';
import { KeyRound, Loader2, LogOut } from 'lucide-react';
import { clearBearerToken, getStoredBearerToken, storeBearerToken } from '../lib/auth-token.mjs';

type AuthMode = 'login' | 'register';

export default function AuthPage() {
  const [mode, setMode] = useState<AuthMode>('login');
  const [email, setEmail] = useState('');
  const [name, setName] = useState('');
  const [password, setPassword] = useState('');
  const [token, setToken] = useState(() =>
    typeof window === 'undefined' ? '' : getStoredBearerToken()
  );
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleSubmit = async () => {
    if (!email.trim() || !password.trim() || (mode === 'register' && !name.trim())) {
      setError('请填写必要信息');
      return;
    }

    setLoading(true);
    setError('');

    try {
      if (mode === 'register') {
        const registerResponse = await fetch('/api/users/register', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            email: email.trim(),
            password,
            name: name.trim(),
          }),
        });

        if (!registerResponse.ok) {
          const data = await registerResponse.json();
          throw new Error(data.error || '注册失败');
        }
      }

      const loginResponse = await fetch('/api/users/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          email: email.trim(),
          password,
        }),
      });

      const data = await loginResponse.json();
      if (!loginResponse.ok) {
        throw new Error(data.error || '登录失败');
      }

      storeBearerToken(data.token);
      setToken(getStoredBearerToken());
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '认证失败');
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    clearBearerToken();
    setToken('');
  };

  return (
    <div className="mx-auto w-full max-w-3xl px-4 py-12 sm:px-6 lg:px-8">
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg glass-card">
            <KeyRound className="h-5 w-5 text-accent-cyan" />
          </div>
          <h1 className="text-3xl font-black">认证</h1>
        </div>
        <p className="text-text-secondary">
          登录或注册后，浏览器会保存 Bearer Token，批量生图和样式推荐会自动带上 Authorization。
        </p>
      </div>

      <div className="glass-card rounded-2xl p-6">
        <div className="flex gap-2 mb-6">
          <button
            className={`glass-button rounded-lg px-4 py-2 ${mode === 'login' ? 'border-border-glow' : ''}`}
            onClick={() => setMode('login')}
            disabled={loading}
          >
            登录
          </button>
          <button
            className={`glass-button rounded-lg px-4 py-2 ${mode === 'register' ? 'border-border-glow' : ''}`}
            onClick={() => setMode('register')}
            disabled={loading}
          >
            注册
          </button>
        </div>

        <div className="space-y-4">
          {mode === 'register' && (
            <div>
              <label className="block text-sm font-bold mb-2">名称</label>
              <input
                value={name}
                onChange={(event) => setName(event.target.value)}
                className="glass-input w-full rounded-xl px-4 py-3"
                disabled={loading}
              />
            </div>
          )}

          <div>
            <label className="block text-sm font-bold mb-2">邮箱</label>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              className="glass-input w-full rounded-xl px-4 py-3"
              disabled={loading}
            />
          </div>

          <div>
            <label className="block text-sm font-bold mb-2">密码</label>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="glass-input w-full rounded-xl px-4 py-3"
              disabled={loading}
            />
          </div>

          {error && <p className="text-sm font-bold text-red-300">{error}</p>}

          <button
            onClick={handleSubmit}
            disabled={loading}
            className="glass-button flex w-full items-center justify-center gap-2 rounded-xl px-6 py-3"
          >
            {loading && <Loader2 className="h-4 w-4 animate-spin" />}
            {mode === 'login' ? '登录并保存 Token' : '注册、登录并保存 Token'}
          </button>
        </div>
      </div>

      <div className="glass-card mt-6 rounded-2xl p-6">
        <div className="flex items-center justify-between gap-4">
          <div>
            <p className="text-sm font-bold text-text-secondary">当前 Token 状态</p>
            <p className="mt-1 font-mono text-sm text-text-muted">
              {token ? `${token.slice(0, 28)}...` : '未保存'}
            </p>
          </div>
          <button
            onClick={handleLogout}
            className="glass-button flex items-center gap-2 rounded-xl px-4 py-2"
            disabled={!token}
          >
            <LogOut className="h-4 w-4" />
            清除
          </button>
        </div>
      </div>
    </div>
  );
}
