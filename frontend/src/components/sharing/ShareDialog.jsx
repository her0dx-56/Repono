
import { useState } from 'react';
import {
  Link2,
  LockKeyhole,
  Mail,
  Copy,
  Check,
} from 'lucide-react';

import Modal from '../common/Modal.jsx';

export default function ShareDialog({ file, onClose, onCreate }) {
  const [access, setAccess] = useState('public');
  const [recipient, setRecipient] = useState('');
  const [created, setCreated] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);
  const [copied, setCopied] = useState(false);

  async function submit(e) {
    e.preventDefault();
    setError('');

    const email = recipient.trim();

    if (
      access === 'private' &&
      !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)
    ) {
      setError('Enter a valid recipient email.');
      return;
    }

    try {
      setLoading(true);

      const result = await onCreate(
        file,
        access,
        access === 'private' ? email : ''
      );

      setCreated({
        ...result,
        access,
        recipient: email,
      });
    } catch (err) {
      setError(err.message || 'Could not create share.');
    } finally {
      setLoading(false);
    }
  }

  const url = created?.url || '';

  async function copyLink() {
    if (!url) {
      setError('Public share link is unavailable.');
      return;
    }

    try {
      await navigator.clipboard.writeText(url);
      setCopied(true);
    } catch {
      setError('Could not copy the link. Please copy it manually.');
    }
  }

  return (
    <Modal
      title={created ? 'Your share is ready' : 'Share a file'}
      onClose={onClose}
    >
      <div className="share-file">
        <span className="file-type pdf">
          <Link2 />
        </span>

        <div>
          <b>{file.name}</b>
          <small>Private by default</small>
        </div>
      </div>

      {created ? (
        <div className="share-success">
          <div className="success-mark">
            <Check />
          </div>

          <h3>Ready to pass it on!</h3>

          <p>
            {created.access === 'public'
              ? 'Anyone with this link can access the shared file.'
              : `Private share created for ${created.recipient}.`}
          </p>

          {created.access === 'public' && (
            <div className="copy-field">
              <input
                readOnly
                value={url}
                aria-label="Public share link"
                onFocus={(e) => e.target.select()}
              />

              <button
                type="button"
                className="btn btn-blue"
                onClick={copyLink}
              >
                {copied ? <Check size={16} /> : <Copy size={16} />}
                {copied ? 'Copied' : 'Copy'}
              </button>
            </div>
          )}

          {error && <p className="form-error">{error}</p>}

          <button
            type="button"
            className="btn btn-dark full"
            onClick={onClose}
          >
            Done
          </button>
        </div>
      ) : (
        <form className="share-form" onSubmit={submit}>
          <label
            className={`choice ${access === 'public' ? 'selected' : ''}`}
          >
            <input
              type="radio"
              name="access"
              value="public"
              checked={access === 'public'}
              onChange={() => {
                setAccess('public');
                setError('');
              }}
            />

            <span className="choice-icon">
              <Link2 />
            </span>

            <span>
              <b>Public link</b>
              <small>
                Anyone with the link can access this file.
              </small>
            </span>
          </label>

          <label
            className={`choice ${access === 'private' ? 'selected' : ''}`}
          >
            <input
              type="radio"
              name="access"
              value="private"
              checked={access === 'private'}
              onChange={() => {
                setAccess('private');
                setError('');
              }}
            />

            <span className="choice-icon">
              <LockKeyhole />
            </span>

            <span>
              <b>Private share</b>
              <small>
                Share with a specific person by email.
              </small>
            </span>
          </label>

          {access === 'private' && (
            <label className="form-field">
              Recipient email

              <div className="recipient-input">
                <Mail size={17} />
                <input
                  type="email"
                  value={recipient}
                  onChange={(e) => setRecipient(e.target.value)}
                  placeholder="name@example.com"
                  autoComplete="email"
                  required
                />
              </div>
            </label>
          )}

          {error && <p className="form-error">{error}</p>}

          <p className="form-hint">
            {access === 'private'
              ? 'The recipient must have a registered Repono account.'
              : 'Anyone with the generated link can access this file.'}
          </p>

          <button
            className="btn btn-red full"
            type="submit"
            disabled={loading}
          >
            {loading ? 'Creating...' : 'Create share'}
            {!loading && <Link2 size={17} />}
          </button>
        </form>
      )}
    </Modal>
  );
}