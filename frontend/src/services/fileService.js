
import { apiRequest, API_BASE_URL, getToken } from "./apiClient.js";
import { parseJwt } from "./authService.js";

// Clean up legacy global share keys from older builds to prevent cross-account pollution
try {
  localStorage.removeItem("repono_user_shares_v1");
  localStorage.removeItem("repono_user_shares_guest_v1");
} catch {}

function getCurrentUserId() {
  try {
    const token = getToken();
    if (token) {
      const claims = parseJwt(token);
      if (claims?.user_id) return claims.user_id;
    }
    const raw = localStorage.getItem("repono_user_v1");
    if (raw) {
      const user = JSON.parse(raw);
      if (user?.id) return user.id;
      if (user?.email) return user.email;
    }
    return null;
  } catch {
    return null;
  }
}

function getUserSharesStorageKey() {
  const uid = getCurrentUserId();
  return uid ? `repono_user_shares_${uid}_v1` : null;
}

function getUserShares() {
  const key = getUserSharesStorageKey();
  if (!key) return [];
  try {
    const raw = localStorage.getItem(key);
    return raw ? JSON.parse(raw) : [];
  } catch {
    return [];
  }
}

function saveUserShare(share) {
  const key = getUserSharesStorageKey();
  if (!key) return;
  try {
    const existing = getUserShares().filter((s) => s.id !== share.id);
    localStorage.setItem(key, JSON.stringify([share, ...existing]));
  } catch (e) {
    console.warn("Failed to persist user share:", e);
  }
}

function markShareInactive(id) {
  const key = getUserSharesStorageKey();
  if (!key) return;
  try {
    const updated = getUserShares().map((s) =>
      s.id === id || s.token === id || s.shareId === id ? { ...s, active: false } : s
    );
    localStorage.setItem(key, JSON.stringify(updated));
  } catch (e) {
    console.warn("Failed to update user shares:", e);
  }
}

function markShareInactiveByFileId(fileId, access = null) {
  const key = getUserSharesStorageKey();
  if (!key) return;
  try {
    const updated = getUserShares().map((s) => {
      if (s.fileId === fileId && (!access || s.access === access)) {
        return { ...s, active: false };
      }
      return s;
    });
    localStorage.setItem(key, JSON.stringify(updated));
  } catch (e) {
    console.warn("Failed to update user shares:", e);
  }
}

// Normalize Go JSON responses for the existing React frontend.
function normalizeFile(file) {
  if (!file) return null;
  const name = file.name ?? "Untitled";
  const type = file.type ?? file.content_type ?? file.contentType ?? "";
  const createdAt = file.created_at ?? file.createdAt ?? file.modified ?? new Date().toISOString();

  return {
    ...file,
    id: file.id,
    name,
    size: Number(file.size) || 0,
    type,
    contentType: type,
    modified: file.modified ?? createdAt,
    createdAt,
    owner: file.owner ?? "You",
    shared: Boolean(file.shared),
  };
}

function normalizeShare(share, extra = {}) {
  if (!share) return null;
  const id = share.share_id ?? share.id ?? share.ID ?? share.token ?? share.Token ?? extra.id;
  const fileId = share.file_id ?? share.fileId ?? share.FileID ?? extra.fileId;
  const fileName = share.name ?? share.fileName ?? extra.fileName ?? "Shared File";
  const type = share.type ?? share.content_type ?? share.contentType ?? extra.type ?? "";
  const token = share.token ?? share.Token ?? extra.token ?? "";
  const access = share.access ?? (token ? "public" : "private");
  const recipient = share.recipient ?? share.shared_with_user_email ?? extra.recipient ?? "";
  const created =
    share.shared_at ??
    share.sharedAt ??
    share.created_at ??
    share.createdAt ??
    share.created ??
    share.CreatedAt ??
    new Date().toISOString();
  const isRevoked = Boolean(share.revoked_at || share.RevokedAt);

  return {
    ...share,
    id: id || token || fileId,
    shareId: share.share_id ?? id,
    fileId,
    fileName,
    name: fileName,
    size: Number(share.size ?? extra.size) || 0,
    type,
    contentType: type,
    created,
    createdAt: created,
    sharedAt: created,
    access,
    recipient,
    token,
    url: share.url || (token ? fileService.getPublicShareUrl(token) : ""),
    active: share.active !== false && !isRevoked,
  };
}

export const fileService = {
  // Get files owned by the authenticated user plus files shared with the user.
  async listFiles() {
    const [ownedResult, sharedResult] = await Promise.allSettled([
      apiRequest("/files"),
      apiRequest("/files/shared"),
    ]);

    const ownedFiles =
      ownedResult.status === "fulfilled" && Array.isArray(ownedResult.value)
        ? ownedResult.value.map(normalizeFile).filter(Boolean)
        : [];

    const sharedFiles =
      sharedResult.status === "fulfilled" && Array.isArray(sharedResult.value)
        ? sharedResult.value
            .map((s) => {
              const fileId = s.file_id ?? s.fileID ?? s.id;
              const name = s.name ?? s.fileName ?? "Shared File";
              const type = s.content_type ?? s.contentType ?? "";
              const createdAt =
                s.created_at ?? s.createdAt ?? s.shared_at ?? new Date().toISOString();
              const modified = s.shared_at ?? s.sharedAt ?? createdAt;
              return {
                id: fileId,
                name,
                size: Number(s.size) || 0,
                type,
                contentType: type,
                modified,
                createdAt,
                owner: "Shared with you",
                shared: true,
                isSharedWithMe: true,
                shareId: s.share_id ?? s.id,
                permission: s.permission ?? "read",
              };
            })
            .filter(Boolean)
        : [];

    const ownedIds = new Set(ownedFiles.map((f) => f.id));
    const uniqueShared = sharedFiles.filter((f) => !ownedIds.has(f.id));

    return [...ownedFiles, ...uniqueShared];
  },

  // Get all active shares created by the authenticated user.
  async listShares() {
    const localShares = getUserShares()
      .map((s) => normalizeShare(s))
      .filter((s) => s && s.active !== false);

    return localShares;
  },

  // Upload a file using multipart/form-data.
  async upload(file) {
    const formData = new FormData();
    formData.append("file", file);

    const response = await apiRequest("/files", {
      method: "POST",
      body: formData,
    });

    return normalizeFile(response);
  },

  // Download a file owned by or shared with the authenticated user.
  async download(fileOrId) {
    const file =
      typeof fileOrId === "object" && fileOrId !== null
        ? fileOrId
        : { id: fileOrId };

    const blob = await this.downloadFile(file.id);
    const { downloadBlob } = await import("../utils/fileHelpers.js");
    downloadBlob(blob, file.name || blob.name || "download");
    return true;
  },

  async downloadFile(id) {
    const token = getToken();

    const response = await fetch(
      `${API_BASE_URL}/files/${encodeURIComponent(id)}`,
      {
        method: "GET",
        headers: token ? { Authorization: `Bearer ${token}` } : {},
      }
    );

    if (!response.ok) {
      const message = await response.text();
      throw new Error(message ? message.trim() : "Failed to download file");
    }

    let filename = "";
    const disposition = response.headers.get("content-disposition");
    if (disposition) {
      const match = disposition.match(/filename="?([^";]+)"?/i);
      if (match && match[1]) {
        filename = match[1];
      }
    }

    const blob = await response.blob();
    if (filename) {
      blob.name = filename;
    }
    return blob;
  },

  // Delete a file owned by the authenticated user.
  async deleteFile(id) {
    await apiRequest(`/files/${encodeURIComponent(id)}`, {
      method: "DELETE",
    });

    markShareInactiveByFileId(id);
  },

  // Share a file with another registered user by email.
  async createShare(file, access, recipient) {
    const email =
      typeof recipient === "string"
        ? recipient.trim()
        : recipient?.email?.trim();

    if (!file?.id) {
      throw new Error("A valid file is required.");
    }

    if (access !== "private") {
      throw new Error("Use createPublicShare for public sharing.");
    }

    if (!email) {
      throw new Error("A recipient email is required.");
    }

    const response = await apiRequest(
      `/files/${encodeURIComponent(file.id)}/shares`,
      {
        method: "POST",
        body: JSON.stringify({
          shared_with_user_email: email,
        }),
      }
    );

    const normalized = normalizeShare(response, {
      fileId: file.id,
      fileName: file.name,
      size: file.size,
      type: file.type,
      access: "private",
      recipient: email,
    });

    saveUserShare(normalized);
    return normalized;
  },

  // Revoke user-to-user share or delegate to public revocation.
  async revokeShare(id) {
    const localShares = getUserShares();
    const target = localShares.find(
      (s) => s.id === id || s.token === id || s.shareId === id
    );

    if (target?.access === "public" && target.fileId) {
      return this.revokePublicShare(target.fileId);
    }

    markShareInactive(id);
  },

  // Create a public link for a file.
  async createPublicShare(fileId, expiresAt = null, fileDetails = null) {
    const response = await apiRequest(
      `/files/${encodeURIComponent(fileId)}/public-share`,
      {
        method: "POST",
        body: JSON.stringify({
          expires_at: expiresAt,
        }),
      }
    );

    const token = response?.token ?? response?.Token;
    const normalized = normalizeShare(response, {
      fileId,
      fileName: fileDetails?.name,
      size: fileDetails?.size,
      type: fileDetails?.type,
      token,
      access: "public",
    });

    saveUserShare(normalized);
    return normalized;
  },

  // Download a publicly shared file using its token.
  async downloadPublicShare(token) {
    const response = await fetch(
      `${API_BASE_URL}/share/${encodeURIComponent(token)}`
    );

    if (!response.ok) {
      const message = await response.text();
      throw new Error(message ? message.trim() : "Unable to download shared file");
    }

    let filename = "";
    const disposition = response.headers.get("content-disposition");
    if (disposition) {
      const match = disposition.match(/filename="?([^";]+)"?/i);
      if (match && match[1]) {
        filename = match[1];
      }
    }

    const blob = await response.blob();
    if (filename) {
      blob.name = filename;
    }
    return blob;
  },

  // Revoke a file's public share link.
  async revokePublicShare(fileId) {
    await apiRequest(
      `/files/${encodeURIComponent(fileId)}/public-share`,
      {
        method: "DELETE",
      }
    );

    markShareInactiveByFileId(fileId, "public");
  },

  // Fetch metadata of a publicly shared file using HEAD request
  async getPublicShareInfo(token) {
    const local = getUserShares().find((s) => s.token === token || s.id === token);

    try {
      const controller = new AbortController();
      const timeoutId = setTimeout(() => controller.abort(), 6000);

      const response = await fetch(
        `${API_BASE_URL}/share/${encodeURIComponent(token)}`,
        {
          method: "HEAD",
          signal: controller.signal,  
        }
      );
      clearTimeout(timeoutId);

      if (!response.ok) {
        if (response.status === 404) {
          throw new Error("This share link does not exist or has expired.");
        }
        throw new Error(`Unable to load share details (${response.status})`);
      }

      let name = local?.fileName || local?.name || "Shared File";
      const disposition = response.headers.get("content-disposition");
      if (disposition) {
        const utfMatch = disposition.match(/filename\*=UTF-8''([^;]+)/i);
        const standardMatch = disposition.match(/filename="?([^";]+)"?/i);
        if (utfMatch && utfMatch[1]) {
          try {
            name = decodeURIComponent(utfMatch[1]);
          } catch {
            name = utfMatch[1];
          }
        } else if (standardMatch && standardMatch[1]) {
          name = standardMatch[1];
        }
      }

      const sizeHeader = response.headers.get("content-length");
      const size = sizeHeader ? parseInt(sizeHeader, 10) : (local?.size || 0);
      const type = response.headers.get("content-type") || local?.type || "";

      return {
        name,
        size,
        type,
        token,
      };
    } catch (err) {
      if (local && err.name === "AbortError") {
        return {
          name: local.fileName || local.name || "Shared File",
          size: local.size || 0,
          type: local.type || "",
          token,
        };
      }
      throw err;
    }
  },

  // Construct the frontend public share page URL.
  getPublicShareUrl(token) {
    return `${window.location.origin}/share/${encodeURIComponent(token)}`;
  },

  // Direct backend download URL for public share (triggers instant browser Save As dialogue)
  getPublicDownloadUrl(token) {
    return `${API_BASE_URL}/share/${encodeURIComponent(token)}?download=true`;
  },
};