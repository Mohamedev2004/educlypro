import { axiosInstance } from "../axios"
import type {
  CreateTeacherPayload,
  Teacher,
  TeacherEnvelope,
  TeachersListData,
  TeachersListEnvelope,
  TeachersListParams,
  UpdateTeacherClassesPayload,
  UpdateTeacherPayload,
  UpdateTeacherSubjectsPayload,
} from "../types/teachers.types"

export const TeachersService = {
  async list(params: TeachersListParams): Promise<TeachersListData> {
    const response = await axiosInstance.get<TeachersListEnvelope>(
      "/teachers",
      {
        params: {
          page: params.page,
          per_page: params.perPage,
          q: params.q,
          sort: params.sort,
          direction: params.direction,
        },
      }
    )
    return response.data.data
  },

  async getBySlug(slug: string): Promise<Teacher> {
    const response = await axiosInstance.get<TeacherEnvelope>(
      `/teachers/${slug}`
    )
    return response.data.data
  },

  async create(payload: CreateTeacherPayload): Promise<Teacher> {
    const response = await axiosInstance.post<TeacherEnvelope>(
      "/teachers",
      payload
    )
    return response.data.data
  },

  async update(id: number, payload: UpdateTeacherPayload): Promise<Teacher> {
    const response = await axiosInstance.put<TeacherEnvelope>(
      `/teachers/${id}`,
      payload
    )
    return response.data.data
  },

  async updateSubjects(
    id: number,
    payload: UpdateTeacherSubjectsPayload
  ): Promise<Teacher> {
    const response = await axiosInstance.put<TeacherEnvelope>(
      `/teachers/${id}/subjects`,
      payload
    )
    return response.data.data
  },

  async updateClasses(
    id: number,
    payload: UpdateTeacherClassesPayload
  ): Promise<Teacher> {
    const response = await axiosInstance.put<TeacherEnvelope>(
      `/teachers/${id}/classes`,
      payload
    )
    return response.data.data
  },

  async delete(id: number): Promise<void> {
    await axiosInstance.delete(`/teachers/${id}`)
  },
}
