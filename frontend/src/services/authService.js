
import { apiRequest, setToken, getToken } from "./apiClient.js";

const USER_KEY = "repono_user_v1";

export function parseJwt(token) {
  try {
    if (!token || typeof token !== "string") return null;
    const base64Url = token.split(".")[1];
    if (!base64Url) return null;
    const base64 = base64Url.replace(/-/g, "+").replace(/_/g, "/");
    const jsonPayload = decodeURIComponent(
      atob(base64)
        .split("")
        .map((c) => "%" + ("00" + c.charCodeAt(0).toString(16)).slice(-2))
        .join("")
    );
    return JSON.parse(jsonPayload);
  } catch {
    return null;
  }
}

function normalizeUser(user, fallback = {}, token = null) {
  const tokenToUse = token || getToken();
  const jwtClaims = tokenToUse ? parseJwt(tokenToUse) : null;

  const id =
    user?.id ??
    user?.ID ??
    user?.Id ??
    user?.user_id ??
    jwtClaims?.user_id ??
    fallback?.id ??
    null;

  const username =
    user?.username ??
    user?.Username ??
    fallback?.username ??
    user?.name ??
    fallback?.name ??
    "";

  const name = user?.name ?? fallback?.name ?? username;
  const email = user?.email ?? user?.Email ?? fallback?.email ?? "";

  return {
    id,
    username,
    name,
    email,
  };
}

function saveSession(user, token, fallback = {}) {
  if (!token) {
    throw new Error("Login response did not contain a token.");
  }

  setToken(token);
  const normalizedUser = normalizeUser(user, fallback, token);
  localStorage.setItem(USER_KEY, JSON.stringify(normalizedUser));

  return normalizedUser;
}

export const authService = {
  async signup({ username, name, email, password }) {
    const cleanUsername = (username ?? name ?? "").trim();
    const cleanEmail = email.trim();

    await apiRequest("/users/register", {
      method: "POST",
      body: JSON.stringify({
        username: cleanUsername,
        email: cleanEmail,
        password,
      }),
    });

    // Registration returns UserResponse, not a JWT.
    // Log in after successful registration to obtain the token.
    return this.login({
      email: cleanEmail,
      password,
      username: cleanUsername,
    });
  },

  async login({ email, password, username, name }) {
    const response = await apiRequest("/users/login", {
      method: "POST",
      body: JSON.stringify({
        email: email.trim(),
        password,
      }),
    });

    const payload =
      typeof response === "string" ? JSON.parse(response) : response;

    const token = payload?.token ?? payload?.Token ?? payload?.data?.token;
    const user = payload?.user ?? payload?.User ?? payload?.data?.user;

    if (!token) {
      console.error("Unexpected login response:", payload);
      throw new Error("Login response did not contain a token.");
    }

    return saveSession(user, token, {
      email: email.trim(),
      username: username || name,
    });
  },

  logout() {
    setToken(null);
    localStorage.removeItem(USER_KEY);
  },

  getCurrent() {
    const token = getToken();
    if (!token) {
      return null;
    }

    try {
      const raw = localStorage.getItem(USER_KEY);
      const user = raw ? JSON.parse(raw) : {};
      return normalizeUser(user, {}, token);
    } catch {
      return null;
    }
  },
};