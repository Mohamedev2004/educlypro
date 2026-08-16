import { axiosInstance } from "../axios"
import type {
  AuthData,
  AuthEnvelope,
  LoginPayload,
  UpdateProfilePayload,
  UpdatePasswordPayload,
  ForgotPasswordPayload,
  ResetPasswordPayload,
  User,
} from "../types/auth.types"

export const AuthService = {
  async login(payload: LoginPayload): Promise<User> {
    const response = await axiosInstance.post<AuthEnvelope<AuthData>>(
      "/auth/login",
      payload
    )
    return response.data.data.user
  },

  async me(): Promise<User> {
    const response = await axiosInstance.get<AuthEnvelope<AuthData>>("/auth/me")
    return response.data.data.user
  },

  async refresh(): Promise<void> {
    await axiosInstance.post("/auth/refresh")
  },

  async logout(): Promise<void> {
    await axiosInstance.post("/auth/logout")
  },

  async updateProfile(payload: UpdateProfilePayload): Promise<User> {
    const response = await axiosInstance.put<AuthEnvelope<User>>(
      "/auth/profile",
      payload
    )
    return response.data.data
  },

  async updatePassword(payload: UpdatePasswordPayload): Promise<void> {
    await axiosInstance.put("/auth/password", payload)
  },

  async forgotPassword(payload: ForgotPasswordPayload): Promise<void> {
    await axiosInstance.post("/auth/forgot-password", payload)
  },

  async resetPassword(payload: ResetPasswordPayload): Promise<void> {
    await axiosInstance.post("/auth/reset-password", payload)
  },
}