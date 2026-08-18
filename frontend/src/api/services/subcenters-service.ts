import { axiosInstance } from "../axios"
import type {
  CreateSubCenterPayload,
  SubCenter,
  SubCenterEnvelope,
  SubCentersListData,
  SubCentersListEnvelope,
  UpdateSubCenterPayload,
} from "../types/subcenters.types"

export const SubCentersService = {
  async list(): Promise<SubCentersListData> {
    const response =
      await axiosInstance.get<SubCentersListEnvelope>("/subcenters")
    return response.data.data
  },

  async create(payload: CreateSubCenterPayload): Promise<SubCenter> {
    const response = await axiosInstance.post<SubCenterEnvelope>(
      "/subcenters",
      payload
    )
    return response.data.data
  },

  async update(
    id: number,
    payload: UpdateSubCenterPayload
  ): Promise<SubCenter> {
    const response = await axiosInstance.put<SubCenterEnvelope>(
      `/subcenters/${id}`,
      payload
    )
    return response.data.data
  },

  async delete(id: number): Promise<void> {
    await axiosInstance.delete(`/subcenters/${id}`)
  },
}
