const key = "access_token";

class SessionService {
  set(token: string) {
    localStorage.setItem(key, token);
  }

  get() {
    return localStorage.getItem(key);
  }

  clear() {
    localStorage.removeItem(key);
  }
}

export const sessionService = new SessionService();
