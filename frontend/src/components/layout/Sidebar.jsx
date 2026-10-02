import {
  Files,
  FolderOpen,
  Share2,
  Settings,
  House,
  LogOut,
  X,
} from 'lucide-react';

const items = [
  ['Overview', 'overview', House],
  ['My files', 'files', FolderOpen],
  ['Shared links', 'shares', Share2],
  ['Settings', 'settings', Settings],
];

export default function Sidebar({
  active,
  onSelect,
  user,
  onLogout,
  mobileOpen,
  onClose,
}) {
  return (
    <>
      <div
        className={`sidebar-backdrop ${mobileOpen ? 'visible' : ''}`}
        onClick={onClose}
      />
      <aside className={`sidebar ${mobileOpen ? 'sidebar-open' : ''}`}>
        <div className="sidebar-brand">
          <span className="brand-symbol">
            <i />
            <b />
            <em />
          </span>
          <b>
            REPONO<span>.</span>
          </b>
          <button className="sidebar-close" onClick={onClose}>
            <X />
          </button>
        </div>

        <div className="workspace-label">WORKSPACE</div>

        <nav>
          {items.map(([label, key, Icon]) => (
            <button
              key={key}
              className={active === key ? 'side-link active' : 'side-link'}
              onClick={() => {
                onSelect(key);
                onClose();
              }}
            >
              <Icon size={19} />
              {label}
              {key === 'shares' && <span className="side-count">↗</span>}
            </button>
          ))}
        </nav>

        <div className="sidebar-bottom">
          <div className="storage-mini">
            <div className="storage-title">
              <span>Storage</span>
              <b>Cloud</b>
            </div>
            <div className="storage-track">
              <i />
            </div>
            <small>Repono S3 Storage</small>
          </div>

          <div className="user-chip">
            <div className="avatar">
              {user?.name?.[0]?.toUpperCase() || 'R'}
            </div>
            <div className="user-label">
              <b>{user?.name || 'Repono user'}</b>
              <small>{user?.email}</small>
            </div>
            <button
              className="logout-icon"
              title="Log out"
              onClick={onLogout}
            >
              <LogOut size={17} />
            </button>
          </div>
        </div>
      </aside>
    </>
  );
}
