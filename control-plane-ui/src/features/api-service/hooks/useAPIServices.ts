import { useQuery } from "@tanstack/react-query";
import { getAPIServices } from "../../../shared/api/apiServiceApi";
import type { APIService } from "../../../shared/types/apiService";

export const useAPIServices = () => {
  return useQuery<APIService[]>({
    queryKey: ["api-services"],
    queryFn: getAPIServices,
  });
};
