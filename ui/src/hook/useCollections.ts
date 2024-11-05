import { useQuery } from "react-query";
import { CollectionsResponse } from "../type/collection";
import { fetchCollections } from "../adapter/collections.adapter";

export const useCollections = () => {
  const { data, error, isFetching, isLoading } = useQuery<CollectionsResponse>(
    ["/collections"],
    {
      queryFn: fetchCollections,
      staleTime: 1000 * 60 * 60 * 60,
    }
  );

  return {
        location: data?.location,
        collections: data?.collections ?? [],
        isLoading: isFetching || isLoading,
        error,
  };
};
