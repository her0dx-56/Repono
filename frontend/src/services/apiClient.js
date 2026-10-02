
export const API_BASE_URL =
  import.meta.env.VITE_API_BASE_URL !== undefined &&
  import.meta.env.VITE_API_BASE_URL !== ""
    ? import.meta.env.VITE_API_BASE_URL
    : "";

const TOKEN_KEY = "repono_access_token_v1";

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token) {
  if (token) {
    localStorage.setItem(TOKEN_KEY, token);
  } else {
    localStorage.removeItem(TOKEN_KEY);
  }
}

export async function apiRequest(path, options = {}) {
  const token = getToken();
  const headers = new Headers(options.headers || {});

  if (token) {
    headers.set("Authorization", `Bearer ${token}`);
  }

  const isFormData = options.body instanceof FormData;

  if (options.body != null && !isFormData && !headers.has("Content-Type")) {
    headers.set("Content-Type", "application/json");
  }

  const response = await fetch(`${API_BASE_URL}${path}`, {
    ...options,
    headers,
  });

  if (response.status === 204) return null;

  const rawText = await response.text();
  let body = null;
  try {
    body = rawText ? JSON.parse(rawText) : null;
  } catch {
    body = rawText ? rawText.trim() : null;
  }

  if (!response.ok) {
    let message;
    if (typeof body === "object" && body !== null) {
      message = body.error || body.message || `Request failed (${response.status})`;
    } else if (typeof body === "string" && body.trim()) {
      message = body.trim();
    } else {
      message = `Request failed (${response.status})`;
    }

    if (response.status === 401) {
      setToken(null);
      localStorage.removeItem("repono_user_v1");
    }

    throw new Error(message);
  }

  return body;
}