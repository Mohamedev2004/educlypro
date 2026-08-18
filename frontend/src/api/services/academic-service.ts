import { axiosInstance } from "../axios"
import type {
  AcademicTree,
  AcademicTreeEnvelope,
  AddGradePayload,
  AddMajorPayload,
  AddSubjectPayload,
  Grade,
  GradeEnvelope,
  Major,
  MajorEnvelope,
  Subject,
  SubjectEnvelope,
} from "../types/academic.types"

export const AcademicService = {
  async tree(): Promise<AcademicTree> {
    const response =
      await axiosInstance.get<AcademicTreeEnvelope>("/academic/grades")
    return response.data.data
  },

  async addGrade(payload: AddGradePayload): Promise<Grade> {
    const response = await axiosInstance.post<GradeEnvelope>(
      "/academic/grades",
      payload
    )
    return response.data.data
  },

  async removeGrade(gradeId: number): Promise<void> {
    await axiosInstance.delete(`/academic/grades/${gradeId}`)
  },

  async addMajor(gradeId: number, payload: AddMajorPayload): Promise<Major> {
    const response = await axiosInstance.post<MajorEnvelope>(
      `/academic/grades/${gradeId}/majors`,
      payload
    )
    return response.data.data
  },

  async removeMajor(majorId: number): Promise<void> {
    await axiosInstance.delete(`/academic/majors/${majorId}`)
  },

  async addSubject(
    majorId: number,
    payload: AddSubjectPayload
  ): Promise<Subject> {
    const response = await axiosInstance.post<SubjectEnvelope>(
      `/academic/majors/${majorId}/subjects`,
      payload
    )
    return response.data.data
  },

  async removeSubject(subjectId: number): Promise<void> {
    await axiosInstance.delete(`/academic/subjects/${subjectId}`)
  },
}
