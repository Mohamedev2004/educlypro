import { axiosInstance } from "../axios"
import type {
  CreateStaffPayload,
  Staff,
  StaffEnvelope,
  StaffListData,
  StaffListEnvelope,
  StaffListParams,
  UpdateStaffPayload,
} from "../types/staff.types"

export const StaffService = {
  async list(params: StaffListParams): Promise<StaffListData> {
    const response = await axiosInstance.get<StaffListEnvelope>("/staff", {
      params: {
        page: params.page,
        per_page: params.perPage,
        role: params.role,
        q: params.q,
        sort: params.sort,
        direction: params.direction,
      },
    })
    return response.data.data
  },

  async create(payload: CreateStaffPayload): Promise<Staff> {
    const response = await axiosInstance.post<StaffEnvelope>("/staff", payload)
    return response.data.data
  },

  async update(id: number, payload: UpdateStaffPayload): Promise<Staff> {
    const response = await axiosInstance.put<StaffEnvelope>(`/staff/${id}`, payload)
    return response.data.data
  },

  async delete(id: number): Promise<void> {
    await axiosInstance.delete(`/staff/${id}`)
  },
}
