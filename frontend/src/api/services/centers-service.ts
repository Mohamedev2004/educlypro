import { axiosInstance } from "../axios"
import type { Staff } from "../types/staff.types"
import type { SubCentersListData } from "../types/subcenters.types"
import type {
  AddOwnerPayload,
  AddStaffEnvelope,
  AddStaffPayload,
  Center,
  CenterDetail,
  CenterDetailEnvelope,
  CenterEnvelope,
  CenterOwner,
  CenterOwnerEnvelope,
  CenterSubCentersEnvelope,
  CentersListData,
  CentersListEnvelope,
  CentersListParams,
  CreateCenterPayload,
} from "../types/centers.types"

export const CentersService = {
  async list(params: CentersListParams): Promise<CentersListData> {
    const response = await axiosInstance.get<CentersListEnvelope>("/centers", {
      params: {
        page: params.page,
        per_page: params.perPage,
        q: params.q,
        sort: params.sort,
        direction: params.direction,
      },
    })
    return response.data.data
  },

  async create(payload: CreateCenterPayload): Promise<Center> {
    const response = await axiosInstance.post<CenterEnvelope>(
      "/centers",
      payload
    )
    return response.data.data
  },

  async getBySlug(slug: string): Promise<CenterDetail> {
    const response = await axiosInstance.get<CenterDetailEnvelope>(
      `/centers/${encodeURIComponent(slug)}`
    )
    return response.data.data
  },

  async addOwner(slug: string, payload: AddOwnerPayload): Promise<CenterOwner> {
    const response = await axiosInstance.post<CenterOwnerEnvelope>(
      `/centers/${encodeURIComponent(slug)}/owner`,
      payload
    )
    return response.data.data
  },

  async addStaff(slug: string, payload: AddStaffPayload): Promise<Staff> {
    const response = await axiosInstance.post<AddStaffEnvelope>(
      `/centers/${encodeURIComponent(slug)}/staff`,
      payload
    )
    return response.data.data
  },

  async listSubCenters(slug: string): Promise<SubCentersListData> {
    const response = await axiosInstance.get<CenterSubCentersEnvelope>(
      `/centers/${encodeURIComponent(slug)}/subcenters`
    )
    return response.data.data
  },
}
