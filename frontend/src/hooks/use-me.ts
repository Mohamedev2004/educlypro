import type { User } from "@/api/types/auth.types"
import { AuthService } from "@/api/services/auth-service"
import { session } from "@/utils/session-utils"
import { useQuery } from "@tanstack/react-query"

export function useMe() {
  return useQuery<User | null>({
    queryKey: ["me"],
    queryFn: async () => {
      // Remove the try/catch. Let the error bubble to the interceptor.
      return await AuthService.me()
    },
    enabled: session.exists(),
    staleTime: 1000 * 60 * 5,
    retry: false, // Keep this false so it doesn't loop infinitely
  })
}
