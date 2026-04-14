import { useMutation, useQueryClient } from "@tanstack/react-query";
import { createClient, type CreateClientPayload } from "../../../shared/api/clientApi";

export const useCreateClient = () => {
  const qc = useQueryClient();

  return useMutation({
    mutationFn: (data: CreateClientPayload) => createClient(data),
    onSuccess: () => {
      qc.invalidateQueries({ queryKey: ["clients"] });
    },
  });
};
