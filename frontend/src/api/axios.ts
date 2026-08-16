import axios, { AxiosError } from "axios"
import { session } from "@/utils/session-utils"

export const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL ?? "http://localhost:8080/api/v1",
  withCredentials: true,
  headers: {
    "Content-Type": "application/json",
  },
})

let refreshPromise: Promise<void> | null = null

function shouldSkipRefresh(url?: string) {
  return Boolean(
    url?.includes("/auth/login") ||
    // url?.includes("/auth/register") ||
    url?.includes("/auth/logout") ||
    url?.includes("/auth/refresh") ||
    url?.includes("/auth/forgot-password") ||
    url?.includes("/auth/reset-password")
  )
}

// ─── Request Interceptor ─────────────────────────────────────────────────────

axiosInstance.interceptors.request.use(
  (config) => {
    config.headers["X-Request-ID"] = crypto.randomUUID()
    return config
  },
  (error) => Promise.reject(error)
)

// ─── Response Interceptor ────────────────────────────────────────────────────

axiosInstance.interceptors.response.use(
  (response) => response,
  async (
    error: AxiosError<{
      message?: string
      code?: string
      errors?: Record<string, string>
    }>
  ) => {
    const originalRequest = error.config as typeof error.config & {
      _retry?: boolean
    }

    const is401 = error.response?.status === 401

    if (
      is401 &&
      originalRequest &&
      !shouldSkipRefresh(originalRequest.url) &&
      !originalRequest._retry
    ) {
      originalRequest._retry = true

      refreshPromise ??= axiosInstance
        .post("/auth/refresh")
        .then(() => undefined)
        .finally(() => {
          refreshPromise = null
        })

      try {
        await refreshPromise;
        // Use the main instance to retry the request
        return axiosInstance(originalRequest); 
      } catch (refreshError) {
        // Only logout if the refresh CALL itself fails (meaning the 7 days are up)
        session.clear();
        window.dispatchEvent(new Event("auth:logout"));
        return Promise.reject(refreshError);
      }
    }

    const message =
      error.response?.data?.message ||
      error.message ||
      "An unexpected error occurred"

    const apiError = {
      message,
      status: error.response?.status,
      code: error.response?.data?.code,
      errors: error.response?.data?.errors,
    }
    return Promise.reject(apiError)
  }
)
