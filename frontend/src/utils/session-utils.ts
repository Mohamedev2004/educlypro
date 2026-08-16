/**
 * Session related utility functions.
 * 
 * Responsibility: Manage session cookies.
 * Layer: Utils
 */

export const session = {
  /**
   * Checks if the session cookie exists.
   */
  exists: () =>
    typeof document !== "undefined" &&
    document.cookie.includes("session_exists=true"),

  /**
   * Clears the session cookie.
   */
  clear: () => {
    document.cookie =
      "session_exists=; path=/; expires=Thu, 01 Jan 1970 00:00:01 GMT;"
  },
}
