import { Copy, Link2, LockKeyhole, Trash2 } from 'lucide-react';
import EmptyState from '../common/EmptyState.jsx';

export default function SharedLinksTable({ shares, files, onRevoke, notify }) {
  const active = shares.filter((s) => s.active);

  if (!active.length) {
    return (
      <EmptyState
        title="No active shares"
        body="Create a share from any file to see it here."
      />
    );
  }

  return (
    <div className="share-list">
      {active.map((s) => {
        const f = files.find((x) => x.id === s.fileId);
        const shareKey = s.token || s.shareId || s.id;
        const url = s.url || `${window.location.origin}/share/${shareKey}`;
        const dateStr = s.created && !isNaN(new Date(s.created).getTime())
          ? new Date(s.created).toLocaleDateString()
          : 'Recently';

        return (
          <article className="share-row" key={s.id || shareKey}>
            <span className="share-type">
              {s.access === 'public' ? <Link2 /> : <LockKeyhole />}
            </span>

            <div className="share-info">
              <b>{f?.name || s.fileName || s.name || 'Shared file'}</b>
              <small>
                {s.access === 'public'
                  ? 'Anyone with link'
                  : `Private · ${s.recipient || 'Recipient'}`}{' '}
                · {dateStr}
              </small>
            </div>

            <span className="share-label">
              {s.access === 'public' ? 'PUBLIC' : 'PRIVATE'}
            </span>

            {s.access === 'public' && (
              <button
                className="icon-button"
                title="Copy link"
                onClick={() => {
                  navigator.clipboard?.writeText(url);
                  notify('Share link copied');
                }}
              >
                <Copy size={17} />
              </button>
            )}

            <button
              className="icon-button danger-hover"
              title="Revoke share"
              onClick={() => onRevoke(s.id || shareKey)}
            >
              <Trash2 size={17} />
            </button>
          </article>
        );
      })}
    </div>
  );
}
