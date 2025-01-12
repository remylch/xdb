import { useQuery } from "react-query";
import { GraphInfoResponse } from "../type/graph.ts";
import { fetchGraphInfos } from "../adapter/graph.adapter.ts";

export function useGraph() {
  const { data, error, isFetching, isLoading } = useQuery<GraphInfoResponse>(
    ["/graph"],
    {
      queryFn: fetchGraphInfos,
      staleTime: 1000 * 60 * 60 * 60,
    }
  );

  return {
    data,
    error,
    isFetching,
    isLoading,
  };
}
