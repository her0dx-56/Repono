import { createContext, useEffect, useState, useCallback } from 'react';
import { fileService } from '../services/fileService.js';
import { useAuth } from '../hooks/useAuth.js';
import { getToken } from '../services/apiClient.js';

export const FileContext = createContext(null);

export function FileProvider({ children }) {
  const { user } = useAuth();
  const [files, setFiles] = useState([]);
  const [shares, setShares] = useState([]);
  const [loading, setLoading] = useState(false);

  const refresh = useCallback(async () => {
    if (!getToken()) {
      setFiles([]);
      setShares([]);
      setLoading(false);
      return;
    }

    setLoading(true);
    try {
      const [f, s] = await Promise.all([
        fileService.listFiles(),
        fileService.listShares(),
      ]);

      const activeFileShareIds = new Set(
        s.filter((item) => item.active).map((item) => item.fileId)
      );

      const mappedFiles = f.map((file) => ({
        ...file,
        shared: Boolean(file.shared || activeFileShareIds.has(file.id)),
      }));

      setFiles(mappedFiles);
      setShares(s);
    } catch (err) {
      console.error("Failed to load files and shares:", err);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    if (user) {
      refresh();
    } else {
      setFiles([]);
      setShares([]);
      setLoading(false);
    }
  }, [user?.id, user?.email, refresh]);

  async function upload(file) {
    const item = await fileService.upload(file);
    setFiles((old) => [item, ...old.filter((f) => f.id !== item.id)]);
    return item;
  }

  async function remove(id) {
    await fileService.deleteFile(id);
    setFiles((old) => old.filter((f) => f.id !== id));
    setShares((old) =>
      old.map((s) => (s.fileId === id ? { ...s, active: false } : s))
    );
  }

  async function createShare(file, access, recipient) {
    const s = await fileService.createShare(file, access, recipient);
    setShares((old) => [s, ...old.filter((x) => x.id !== s.id)]);
    setFiles((old) =>
      old.map((f) => (f.id === file.id ? { ...f, shared: true } : f))
    );
    return s;
  }

  function addShare(share) {
    if (!share) return;
    setShares((old) => [share, ...old.filter((s) => s.id !== share.id)]);
    if (share.fileId) {
      setFiles((old) =>
        old.map((f) => (f.id === share.fileId ? { ...f, shared: true } : f))
      );
    }
  }

  async function revoke(id) {
    await fileService.revokeShare(id);

    setShares((old) =>
      old.map((s) =>
        s.id === id || s.token === id || s.shareId === id
          ? { ...s, active: false }
          : s
      )
    );

    const target = shares.find(
      (s) => s.id === id || s.token === id || s.shareId === id
    );

    if (target?.fileId) {
      setFiles((old) =>
        old.map((f) => {
          if (f.id !== target.fileId) return f;
          const stillShared = shares.some(
            (s) =>
              s.fileId === f.id &&
              s.active &&
              s.id !== id &&
              s.token !== id &&
              s.shareId !== id
          );
          return { ...f, shared: stillShared };
        })
      );
    }
  }

  return (
    <FileContext.Provider
      value={{
        files,
        shares,
        loading,
        refresh,
        upload,
        remove,
        createShare,
        addShare,
        revoke,
      }}
    >
      {children}
    </FileContext.Provider>
  );
}
