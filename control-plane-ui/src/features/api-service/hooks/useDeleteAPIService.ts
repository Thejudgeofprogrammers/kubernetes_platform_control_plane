import { useMutation, useQueryClient } from "@tanstack/react-query";
import { deleteAPIService } from "../../../shared/api/apiServiceApi";

export const useDeleteAPIService = () => {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: deleteAPIService,
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["api-services"] });
    },
  });
};
