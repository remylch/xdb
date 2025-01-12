import { useQuery } from "react-query";
import { NodeInfoResponse } from "../type/graph.ts";
import { fetchNodeInfos } from "../adapter/node.adapter.ts";

export function useNode() {
  const { data, error, isFetching, isLoading } = useQuery<NodeInfoResponse>(
    ["/node"],
    {
      queryFn: fetchNodeInfos,
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
