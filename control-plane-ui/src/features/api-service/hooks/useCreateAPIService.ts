import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createAPIService } from "../../../shared/api/apiServiceApi";

export const useCreateAPIService = () => {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: createAPIService,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["api-services"] });
    },
  });
};
