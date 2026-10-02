
import { useEffect, useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import {
  ArrowLeft,
  ArrowRight,
  Eye,
  EyeOff,
  ShieldCheck,
} from 'lucide-react';

import Brand from '../components/common/Brand.jsx';
import { useAuth } from '../hooks/useAuth.js';
import { validateAuth } from '../utils/validators.js';

export default function AuthPage({ mode }) {
  const signup = mode === 'signup';
  const navigate = useNavigate();
  const { login, loading } = useAuth();

  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [show, setShow] = useState(false);
  const [error, setError] = useState('');

  // Clear form fields when switching between signup and login.
  useEffect(() => {
    setName('');
    setEmail('');
    setPassword('');
    setError('');
    setShow(false);
  }, [mode]);

  async function submit(e) {
    e.preventDefault();

    const err = validateAuth({ name, email, password }, mode);

    if (err) {
      setError(err);
      return;
    }

    setError('');

    try {
      await login(
        {
          username: name.trim(),
          email: email.trim(),
          password,
        },
        mode
      );

      setPassword('');
      navigate('/app');
    } catch (err) {
      if (
        err.message?.toLowerCase().includes('email already exist') ||
        err.message?.toLowerCase().includes('email already exists')
      ) {
        setError('This email is already registered. Please log in.');
      } else if (
        err.message?.toLowerCase().includes('invalid credentials') ||
        err.message?.toLowerCase().includes('invalid password')
      ) {
        setError('Incorrect email or password. Please try again.');
      } else {
        setError(err.message || 'Something went wrong. Please try again.');
      }
    }
  }

  return (
    <main className="auth-page">
      <Link to="/" className="auth-back">
        <ArrowLeft size={17} />
        Back to home
      </Link>

      <section className="auth-aside">
        <div className="auth-geometry">
          <i />
          <b />
          <em />
        </div>

        <p className="eyebrow">YOUR FILES. YOUR WORLD.</p>

        <h1>
          {signup ? (
            <>
              START YOUR
              <br />
              <em>ADVENTURE.</em>
            </>
          ) : (
            <>
              WELCOME
              <br />
              <em>BACK, HERO.</em>
            </>
          )}
        </h1>

        <p>
          Your personal space for storing, sharing, and keeping your work in
          motion.
        </p>

        <div className="auth-note">
          <ShieldCheck />
          Private by default. Yours to manage.
        </div>
      </section>

      <section className="auth-panel">
        <div className="auth-form-wrap">
          <Brand />

          <p className="eyebrow">
            {signup ? 'CREATE YOUR ACCOUNT' : 'GOOD TO SEE YOU AGAIN'}
          </p>

          <h2>
            {signup ? 'Make room for more.' : 'Pick up where you left off.'}
          </h2>

          <p className="auth-sub">
            {signup
              ? 'Set up your Repono workspace in a few steps.'
              : 'Sign in to open your personal file space.'}
          </p>

          <form onSubmit={submit}>
            {signup && (
              <label className="form-field">
                Full name
                <input
                  value={name}
                  onChange={(e) => setName(e.target.value)}
                  placeholder="Jatin Nayak"
                  autoComplete="name"
                  required
                />
              </label>
            )}

            <label className="form-field">
              Email address
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="you@example.com"
                autoComplete="email"
                required
              />
            </label>

            <label className="form-field">
              Password
              <div className="password-field">
                <input
                  type={show ? 'text' : 'password'}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  placeholder="At least 8 characters"
                  autoComplete={signup ? 'new-password' : 'current-password'}
                  required
                />

                <button
                  type="button"
                  onClick={() => setShow(!show)}
                  aria-label={show ? 'Hide password' : 'Show password'}
                >
                  {show ? <EyeOff size={17} /> : <Eye size={17} />}
                </button>
              </div>
            </label>

            {error && <p className="form-error">{error}</p>}

            <button
              className="btn btn-red full auth-submit"
              type="submit"
              disabled={loading}
            >
              {loading
                ? 'Please wait...'
                : signup
                  ? 'Create account'
                  : 'Log in'}
              <ArrowRight size={17} />
            </button>
          </form>

          <p className="auth-switch">
            {signup ? 'Already have an account?' : 'New to Repono?'}{' '}
            <Link to={signup ? '/login' : '/signup'}>
              {signup ? 'Log in' : 'Create an account'}
            </Link>
          </p>

          <p className="demo-warning">
            Connected to Repono Authentication. Protected with secure password hashing and JWT sessions.
          </p>
        </div>
      </section>
    </main>
  );
}