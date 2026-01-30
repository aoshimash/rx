/**
 * Authentication Token Store
 *
 * Manages Bearer token persistence in localStorage for API authentication.
 * Token is requested lazily when API operations are attempted.
 */

const TOKEN_KEY = 'optel_token';

export const authStore = {
  /**
   * Get current token from localStorage
   */
  getToken(): string | null {
    if (typeof window === 'undefined') return null;
    return localStorage.getItem(TOKEN_KEY);
  },

  /**
   * Save token to localStorage
   */
  setToken(token: string): void {
    if (typeof window === 'undefined') return;
    localStorage.setItem(TOKEN_KEY, token);
  },

  /**
   * Remove token from localStorage
   */
  clearToken(): void {
    if (typeof window === 'undefined') return;
    localStorage.removeItem(TOKEN_KEY);
  },

  /**
   * Check if token exists
   */
  hasToken(): boolean {
    return this.getToken() !== null;
  },
};
