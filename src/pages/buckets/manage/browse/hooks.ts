import api from "@/lib/api";
import {
  useMutation,
  useInfiniteQuery,
  UseMutationOptions,
} from "@tanstack/react-query";
import {
  GetObjectsResult,
  PutObjectPayload,
  UseBrowserObjectOptions,
} from "./types";

export const useBrowseObjects = (
  bucket: string,
  options?: UseBrowserObjectOptions
) => {
  return useInfiniteQuery({
    queryKey: ["browse", bucket, options?.prefix, options?.limit],
    initialPageParam: options?.next,
    queryFn: ({ pageParam }) =>
      api.get<GetObjectsResult>(`/browse/${bucket}`, {
        params: { ...options, next: pageParam },
      }),
    getNextPageParam: (lastPage) => lastPage.nextToken || undefined,
  });
};

export const usePutObject = (
  bucket: string,
  options?: UseMutationOptions<any, Error, PutObjectPayload>
) => {
  return useMutation({
    mutationFn: async (body) => {
      const formData = new FormData();
      if (body.file) {
        formData.append("file", body.file);
      }

      return api.put(`/browse/${bucket}/${body.key}`, { body: formData });
    },
    ...options,
  });
};

export const useDeleteObject = (
  bucket: string,
  options?: UseMutationOptions<any, Error, { key: string; recursive?: boolean }>
) => {
  return useMutation({
    mutationFn: (data) =>
      api.delete(`/browse/${bucket}/${data.key}`, {
        params: { recursive: data.recursive },
      }),
    ...options,
  });
};
