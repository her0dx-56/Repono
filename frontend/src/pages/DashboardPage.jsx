import { useMemo, useState } from "react";
import {
  CloudUpload,
  Files,
  HardDrive,
  Share2,
  ArrowDownToLine,
  Plus,
  CheckCircle2,
  Settings,
  ShieldCheck,
  LogOut,
  UserRound,
  RefreshCw,
} from "lucide-react";
import {
  Route,
  Routes,
  useNavigate,
  useLocation,
  Navigate,
} from "react-router-dom";

import { useAuth } from "../hooks/useAuth.js";
import { useFiles } from "../hooks/useFiles.js";
import { useToast } from "../hooks/useToast.js";

import AppLayout from "../components/layout/AppLayout.jsx";
import UploadZone from "../components/files/UploadZone.jsx";
import FileTable from "../components/files/FileTable.jsx";
import FileCard from "../components/files/FileCard.jsx";
import ShareDialog from "../components/sharing/ShareDialog.jsx";
import SharedLinksTable from "../components/sharing/SharedLinksTable.jsx";
import UploadAnimation from "../components/mascot/UploadAnimation.jsx";
import MascotReaction from "../components/mascot/MascotReaction.jsx";
import Modal from "../components/common/Modal.jsx";
import ConfirmDialog from "../components/common/ConfirmDialog.jsx";
import Toast from "../components/common/Toast.jsx";

import { fileService } from "../services/fileService.js";
import {
  formatFileSize,
  formatDate,
} from "../utils/formatFileSize.js";
import { validateFile } from "../utils/validators.js";

const tabInfo = {
  overview: [
    "Overview",
    "A quick look at your digital workspace.",
  ],
  files: [
    "My files",
    "Everything you’ve stored, all in one place.",
  ],
  shares: [
    "Shared links",
    "Manage the files you’ve passed along.",
  ],
  settings: [
    "Settings",
    "Your profile and workspace preferences.",
  ],
};

export default function DashboardPage() {
  const { user, logout } = useAuth();

  const {
    files,
    shares,
    loading,
    upload,
    remove,
    createShare,
    addShare,
    revoke,
  } = useFiles();

  const { toast, notify, clear } = useToast();

  const navigate = useNavigate();
  const loc = useLocation();

  const [active, setActive] = useState(
    loc.pathname.split("/")[2] || "overview"
  );
  const [search, setSearch] = useState("");
  const [modal, setModal] = useState("");
  const [selected, setSelected] = useState(null);

  const [uploadStatus, setUploadStatus] = useState("idle");
  const [uploadName, setUploadName] = useState("");
  const [pendingFile, setPendingFile] = useState(null);
  const [motion, setMotion] = useState(true);

  const current = tabInfo[active] ? active : "overview";

  const filtered = useMemo(
    () =>
      files.filter((f) =>
        f.name.toLowerCase().includes(search.toLowerCase())
      ),
    [files, search]
  );

  function selectTab(key) {
    setActive(key);
    setSearch("");

    navigate(key === "overview" ? "/app" : `/app/${key}`);
  }

  function openUpload() {
    setUploadStatus("idle");
    setUploadName("");
    setPendingFile(null);
    setModal("upload");
  }

  async function handleUploads(items) {
    const file = items[0];

    if (items.length > 1) {
      notify(
        "Please upload one file at a time; other files were not added."
      );
    }

    const err = validateFile(file);

    if (err) {
      notify(err);
      return;
    }

    setPendingFile(file);
    setUploadName(file.name);
    setUploadStatus("uploading");
    setModal("upload");

    try {
      await upload(file);

      setUploadStatus("success");
      notify("File delivered to your workspace");
    } catch (e) {
      setUploadStatus("error");
      notify(e.message || "Upload failed");
    }
  }

  async function download(file) {
    try {
      notify("Download started");
      await fileService.download(file);
    } catch (e) {
      notify(e.message || "Could not download file");
    }
  }

  async function deleteFile() {
    if (!selected) return;

    try {
      await remove(selected.id);
      notify("File removed from workspace");
    } catch (e) {
      notify(e.message || "Failed to remove file");
    } finally {
      setModal("");
      setSelected(null);
    }
  }

  // Create public or private shares.
  async function makeShare(file, access, recipient) {
    try {
      if (access === "public") {
        const result = await fileService.createPublicShare(
          file.id,
          null,
          file
        );

        const token = result?.token ?? result?.Token;
        if (!token) {
          throw new Error(
            "The server did not return a public share token."
          );
        }

        if (typeof addShare === "function") {
          addShare(result);
        }

        notify("Public share created successfully");

        return {
          ...result,
          access: "public",
          url: fileService.getPublicShareUrl(token),
        };
      }

      const email =
        typeof recipient === "string"
          ? recipient.trim()
          : recipient?.email?.trim();

      if (!email) {
        throw new Error("A recipient email is required.");
      }

      const result = await createShare(
        file,
        access,
        email
      );

      notify("Private share created successfully");

      return {
        ...result,
        access: "private",
        recipient: email,
      };
    } catch (error) {
      notify(error.message || "Could not create share");
      throw error;
    }
  }

  function share(file) {
    setSelected(file);
    setModal("share");
  }

  function deletePrompt(file) {
    setSelected(file);
    setModal("delete");
  }

  function details(file) {
    setSelected(file);
    setModal("details");
  }

  const title = tabInfo[current][0];
  const subtitle = tabInfo[current][1];

  const recent = [...files]
    .sort(
      (a, b) =>
        new Date(b.modified || b.createdAt || 0) -
        new Date(a.modified || a.createdAt || 0)
    )
    .slice(0, 4);

  return (
    <>
      <AppLayout
        active={current}
        onSelect={selectTab}
        user={user}
        onLogout={() => {
          logout();
          navigate("/");
        }}
        title={title}
        subtitle={subtitle}
        onUpload={openUpload}
        search={search}
        setSearch={setSearch}
      >
        {/* OVERVIEW */}
        {current === "overview" && (
          <>
            <section className="welcome-banner">
              <div className="welcome-copy">
                <span className="eyebrow">
                  YOUR SPACE IS READY
                </span>

                <h2>
                  Hey,{" "}
                  {user?.name?.split(" ")[0] || "there"}!
                  <br />
                  <em>
                    What are we sending today?
                  </em>
                </h2>

                <p>
                  Keep your work moving. Your files are
                  just a few clicks away.
                </p>

                <button
                  className="btn btn-dark"
                  onClick={openUpload}
                >
                  <Plus size={17} />
                  Upload a file
                </button>
              </div>

              <div className="welcome-shapes">
                <span />
                <i />
                <b />
              </div>

              {motion && (
                <MascotReaction action="ready" />
              )}
            </section>

            <div className="stats-grid">
              <article>
                <span className="stat-icon red">
                  <Files />
                </span>
                <small>Total files</small>
                <strong>{files.length}</strong>
                <em>In your workspace</em>
              </article>

              <article>
                <span className="stat-icon blue">
                  <Share2 />
                </span>
                <small>Active shares</small>
                <strong>
                  {shares.filter((s) => s.active).length}
                </strong>
                <em>Links and invites</em>
              </article>

              <article>
                <span className="stat-icon yellow">
                  <HardDrive />
                </span>
                <small>Storage used</small>
                <strong>
                  {formatFileSize(
                    files.reduce(
                      (n, f) => n + f.size,
                      0
                    )
                  )}
                </strong>
                <em>Total file volume</em>
              </article>
            </div>

            <section className="content-card">
              <div className="section-bar">
                <div>
                  <h2>All files</h2>
                  <p>Your files and files shared with you</p>
                </div>

                <button
                  className="text-button"
                  onClick={() => selectTab("files")}
                >
                  Manage files
                  <ArrowDownToLine size={15} />
                </button>
              </div>

              {loading ? (
                <p className="loading-line">
                  Loading files...
                </p>
              ) : (
                <FileTable
                  files={filtered}
                  onDownload={download}
                  onShare={share}
                  onDelete={deletePrompt}
                  onPreview={details}
                />
              )}
            </section>

            <section className="upload-callout">
              <UploadZone onFiles={handleUploads} />
            </section>
          </>
        )}

        {/* FILES */}
        {current === "files" && (
          <>
            <div className="content-card">
              <div className="section-bar">
                <div>
                  <h2>
                    All files{" "}
                    <span className="count-badge">
                      {filtered.length}
                    </span>
                  </h2>
                  <p>Manage your stored items</p>
                </div>

                <button
                  className="btn btn-red"
                  onClick={openUpload}
                >
                  <Plus size={17} />
                  Upload
                </button>
              </div>

              {loading ? (
                <p className="loading-line">
                  Loading files...
                </p>
              ) : (
                <FileTable
                  files={filtered}
                  onDownload={download}
                  onShare={share}
                  onDelete={deletePrompt}
                  onPreview={details}
                />
              )}
            </div>

            <UploadZone onFiles={handleUploads} />
          </>
        )}

        {/* SHARES */}
        {current === "shares" && (
          <section className="content-card">
            <div className="section-bar">
              <div>
                <h2>Active shares</h2>
                <p>
                  Revoke access when a link or invitation
                  is no longer needed.
                </p>
              </div>

              <span className="count-badge">
                {shares.filter((s) => s.active).length}{" "}
                active
              </span>
            </div>

            <SharedLinksTable
              shares={shares}
              files={files}
              onRevoke={async (id) => {
                await revoke(id);
                notify("Share revoked successfully");
              }}
              notify={notify}
            />
          </section>
        )}

        {/* SETTINGS */}
        {current === "settings" && (
          <div className="settings-grid">
            <section className="content-card settings-card">
              <div className="settings-heading">
                <span className="stat-icon red">
                  <UserRound />
                </span>

                <div>
                  <h2>Profile</h2>
                  <p>Your account details</p>
                </div>
              </div>

              <div className="profile-summary">
                <div className="profile-avatar">
                  {user?.name?.[0]?.toUpperCase() || "R"}
                </div>

                <div>
                  <b>{user?.name}</b>
                  <small>{user?.email}</small>
                  <span>Repono member</span>
                </div>
              </div>

              <div className="setting-line">
                <span>Display name</span>
                <b>{user?.name}</b>
              </div>

              <div className="setting-line">
                <span>Email address</span>
                <b>{user?.email}</b>
              </div>

              <button
                className="btn btn-outline"
                onClick={() =>
                  notify(
                    "Profile editing will be added with backend integration"
                  )
                }
              >
                Edit profile
              </button>
            </section>

            <section className="content-card settings-card">
              <div className="settings-heading">
                <span className="stat-icon blue">
                  <ShieldCheck />
                </span>

                <div>
                  <h2>Security & access</h2>
                  <p>Workspace behavior</p>
                </div>
              </div>

              <div className="setting-line">
                <span>Authentication</span>
                <b>JWT Session</b>
              </div>

              <div className="setting-line">
                <span>File storage</span>
                <b>Repono API · Connected</b>
              </div>

              <div className="setting-line">
                <span>Sharing</span>
                <b>Public & Private Active</b>
              </div>

              <div className="setting-line">
                <span>Dashboard mascot</span>
                <label className="toggle">
                  <input
                    type="checkbox"
                    checked={motion}
                    onChange={(e) =>
                      setMotion(e.target.checked)
                    }
                  />
                  <i />
                </label>
              </div>

              <div className="settings-tip">
                Your workspace is connected to the Repono backend with secure
                JWT authentication and cloud file storage.
              </div>

              <button
                className="btn btn-outline"
                onClick={() => {
                  logout();
                  navigate("/");
                }}
              >
                <LogOut size={16} />
                Log out
              </button>
            </section>
          </div>
        )}

        {/* FOOTER */}
        <footer className="dash-footer">
          <span>© 2026 REPONO</span>
          <span>YOUR FILES. YOUR WORLD.</span>
          <span className="demo-label">
            <i />
            CONNECTED
          </span>
        </footer>
      </AppLayout>

      {/* UPLOAD MODAL */}
      {modal === "upload" && (
        <Modal
          title="Send a file"
          onClose={() => {
            setModal("");
            if (uploadStatus !== "uploading") {
              setUploadStatus("idle");
            }
          }}
        >
          <UploadAnimation
            status={uploadStatus}
            filename={uploadName}
            onRetry={() =>
              pendingFile && handleUploads([pendingFile])
            }
          />

          {(uploadStatus === "idle" || uploadStatus === "error") && (
            <UploadZone onFiles={handleUploads} />
          )}

          <p className="form-hint">
            Files are uploaded directly to your Repono cloud storage.
          </p>

          {uploadStatus !== "uploading" && (
            <button
              className="btn btn-outline full"
              onClick={() => {
                setModal("");
                setUploadStatus("idle");
              }}
            >
              Close
            </button>
          )}
        </Modal>
      )}

      {/* SHARE MODAL */}
      {modal === "share" && selected && (
        <ShareDialog
          file={selected}
          onClose={() => {
            setModal("");
            setSelected(null);
          }}
          onCreate={makeShare}
        />
      )}

      {/* DELETE MODAL */}
      {modal === "delete" && selected && (
        <ConfirmDialog
          title="Delete this file?"
          onCancel={() => {
            setModal("");
            setSelected(null);
          }}
          onConfirm={deleteFile}
          confirmText="Delete file"
        >
          <p>
            <b>{selected.name}</b> will be deleted from your workspace.
          </p>
          <small>
            This action will permanently delete the file from storage.
          </small>
        </ConfirmDialog>
      )}

      {/* FILE DETAILS MODAL */}
      {modal === "details" && selected && (
        <Modal
          title="File details"
          onClose={() => {
            setModal("");
            setSelected(null);
          }}
        >
          <div className="details-file">
            <h3>{selected.name}</h3>
            <p>{formatFileSize(selected.size)}</p>
          </div>

          <div className="detail-lines">
            <span>
              Type <b>{selected.type || "Unknown"}</b>
            </span>

            <span>
              Modified{" "}
              <b>{formatDate(selected.modified)}</b>
            </span>

            <span>
              Owner{" "}
              <b>
                {selected.owner ||
                  (selected.isSharedWithMe ? "Shared with you" : "You")}
              </b>
            </span>

            <span>
              Access{" "}
              <b>
                {selected.shared ? "Shared" : "Private"}
              </b>
            </span>
          </div>

          <div className="dialog-actions">
            {!selected.isSharedWithMe && (
              <button
                className="btn btn-outline"
                onClick={() => share(selected)}
              >
                <Share2 size={16} />
                Share
              </button>
            )}

            <button
              className="btn btn-red"
              onClick={() => download(selected)}
            >
              Download
            </button>
          </div>
        </Modal>
      )}

      <Toast message={toast} onClose={clear} />
    </>
  );
}